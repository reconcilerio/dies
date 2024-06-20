/*
Copyright 2021 the original author or authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/jsonpath"
	"reconciler.io/dies/patch"
	"sigs.k8s.io/yaml"
)

var TypeMetaBlank = (&TypeMetaDie{}).DieFeed(metav1.TypeMeta{})

type TypeMetaDie struct {
	mutable bool
	r       metav1.TypeMeta
	seal    metav1.TypeMeta
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *TypeMetaDie) DieImmutable(immutable bool) *TypeMetaDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *TypeMetaDie) DieFeed(r metav1.TypeMeta) *TypeMetaDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &TypeMetaDie{
		mutable: d.mutable,
		r:       r,
		seal:    d.seal,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *TypeMetaDie) DieFeedPtr(r *metav1.TypeMeta) *TypeMetaDie {
	if r == nil {
		r = &metav1.TypeMeta{}
	}
	return d.DieFeed(*r)
}

// DieFeedJSON returns a new die with the provided JSON. Panics on error.
func (d *TypeMetaDie) DieFeedJSON(j []byte) *TypeMetaDie {
	r := metav1.TypeMeta{}
	if err := json.Unmarshal(j, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAML returns a new die with the provided YAML. Panics on error.
func (d *TypeMetaDie) DieFeedYAML(y []byte) *TypeMetaDie {
	r := metav1.TypeMeta{}
	if err := yaml.Unmarshal(y, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAMLFile returns a new die loading YAML from a file path. Panics on error.
func (d *TypeMetaDie) DieFeedYAMLFile(name string) *TypeMetaDie {
	y, err := os.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return d.DieFeedYAML(y)
}

// DieFeedRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *TypeMetaDie) DieFeedRawExtension(raw runtime.RawExtension) *TypeMetaDie {
	j, err := json.Marshal(raw)
	if err != nil {
		panic(err)
	}
	return d.DieFeedJSON(j)
}

// DieRelease returns the resource managed by the die.
func (d *TypeMetaDie) DieRelease() metav1.TypeMeta {
	if d.mutable {
		return d.r
	}
	return metav1.TypeMeta{
		APIVersion: d.r.APIVersion,
		Kind:       d.r.Kind,
	}
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *TypeMetaDie) DieReleasePtr() *metav1.TypeMeta {
	r := d.DieRelease()
	return &r
}

// DieReleaseJSON returns the resource managed by the die as JSON. Panics on error.
func (d *TypeMetaDie) DieReleaseJSON() []byte {
	r := d.DieReleasePtr()
	j, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return j
}

// DieReleaseYAML returns the resource managed by the die as YAML. Panics on error.
func (d *TypeMetaDie) DieReleaseYAML() []byte {
	r := d.DieReleasePtr()
	y, err := yaml.Marshal(r)
	if err != nil {
		panic(err)
	}
	return y
}

// DieReleaseRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *TypeMetaDie) DieReleaseRawExtension() runtime.RawExtension {
	j := d.DieReleaseJSON()
	raw := runtime.RawExtension{}
	if err := json.Unmarshal(j, &raw); err != nil {
		panic(err)
	}
	return raw
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *TypeMetaDie) DieStamp(fn func(r *metav1.TypeMeta)) *TypeMetaDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// Experimental: DieStampAt uses a JSON path (http://goessner.net/articles/JsonPath/) expression to stamp portions of the resource. The callback is invoked with each JSON path match. Panics if the callback function does not accept a single argument of the same type or a pointer to that type as found on the resource at the target location.
//
// Future iterations will improve type coercion from the resource to the callback argument.
func (d *TypeMetaDie) DieStampAt(jp string, fn interface{}) *TypeMetaDie {
	return d.DieStamp(func(r *metav1.TypeMeta) {
		if ni := reflect.ValueOf(fn).Type().NumIn(); ni != 1 {
			panic(fmt.Errorf("callback function must have 1 input parameters, found %d", ni))
		}
		if no := reflect.ValueOf(fn).Type().NumOut(); no != 0 {
			panic(fmt.Errorf("callback function must have 0 output parameters, found %d", no))
		}

		cp := jsonpath.New("")
		if err := cp.Parse(fmt.Sprintf("{%s}", jp)); err != nil {
			panic(err)
		}
		cr, err := cp.FindResults(r)
		if err != nil {
			// errors are expected if a path is not found
			return
		}
		for _, cv := range cr[0] {
			arg0t := reflect.ValueOf(fn).Type().In(0)

			var args []reflect.Value
			if cv.Type().AssignableTo(arg0t) {
				args = []reflect.Value{cv}
			} else if cv.CanAddr() && cv.Addr().Type().AssignableTo(arg0t) {
				args = []reflect.Value{cv.Addr()}
			} else {
				panic(fmt.Errorf("callback function must accept value of type %q, found type %q", cv.Type(), arg0t))
			}

			reflect.ValueOf(fn).Call(args)
		}
	})
}

// DieWith returns a new die after passing the current die to the callback function. The passed die is mutable.
func (d *TypeMetaDie) DieWith(fns ...func(d *TypeMetaDie)) *TypeMetaDie {
	nd := TypeMetaBlank.DieFeed(d.DieRelease()).DieImmutable(false)
	for _, fn := range fns {
		if fn != nil {
			fn(nd)
		}
	}
	return d.DieFeed(nd.DieRelease())
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *TypeMetaDie) DeepCopy() *TypeMetaDie {
	r := metav1.TypeMeta{
		APIVersion: d.r.APIVersion,
		Kind:       d.r.Kind,
	}
	return &TypeMetaDie{
		mutable: d.mutable,
		r:       r,
		seal:    d.seal,
	}
}

// DieSeal returns a new die for the current die's state that is sealed for comparison in future diff and patch operations.
func (d *TypeMetaDie) DieSeal() *TypeMetaDie {
	return d.DieSealFeed(d.r)
}

// DieSealFeed returns a new die for the current die's state that uses a specific resource for comparison in future diff and patch operations.
func (d *TypeMetaDie) DieSealFeed(r metav1.TypeMeta) *TypeMetaDie {
	if !d.mutable {
		d = d.DeepCopy()
	}
	d.seal = metav1.TypeMeta{
		APIVersion: r.APIVersion,
		Kind:       r.Kind,
	}
	return d
}

// DieSealFeedPtr returns a new die for the current die's state that uses a specific resource pointer for comparison in future diff and patch operations. If the resource is nil, the empty value is used instead.
func (d *TypeMetaDie) DieSealFeedPtr(r *metav1.TypeMeta) *TypeMetaDie {
	if r == nil {
		r = &metav1.TypeMeta{}
	}
	return d.DieSealFeed(*r)
}

// DieSealRelease returns the sealed resource managed by the die.
func (d *TypeMetaDie) DieSealRelease() metav1.TypeMeta {
	return metav1.TypeMeta{
		APIVersion: d.seal.APIVersion,
		Kind:       d.seal.Kind,
	}
}

// DieSealReleasePtr returns the sealed resource pointer managed by the die.
func (d *TypeMetaDie) DieSealReleasePtr() *metav1.TypeMeta {
	r := d.DieSealRelease()
	return &r
}

// DieDiff uses cmp.Diff to compare the current value of the die with the sealed value.
func (d *TypeMetaDie) DieDiff(opts ...cmp.Option) string {
	return cmp.Diff(d.seal, d.r, opts...)
}

// DiePatch generates a patch between the current value of the die and the sealed value.
func (d *TypeMetaDie) DiePatch(patchType types.PatchType) ([]byte, error) {
	return patch.Create(d.seal, d.r, patchType)
}

// Kind is a string value representing the REST resource this object represents.
//
// Servers may infer this from the endpoint the client submits requests to.
//
// Cannot be updated.
//
// In CamelCase.
//
// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
func (d *TypeMetaDie) Kind(v string) *TypeMetaDie {
	return d.DieStamp(func(r *metav1.TypeMeta) {
		r.Kind = v
	})
}

// APIVersion defines the versioned schema of this representation of an object.
//
// # Servers should convert recognized schemas to the latest internal value, and
//
// may reject unrecognized values.
//
// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
func (d *TypeMetaDie) APIVersion(v string) *TypeMetaDie {
	return d.DieStamp(func(r *metav1.TypeMeta) {
		r.APIVersion = v
	})
}

// Group defines the group component of the API version.
func (d *TypeMetaDie) Group(v string) *TypeMetaDie {
	return d.DieStamp(func(r *metav1.TypeMeta) {
		gv, err := schema.ParseGroupVersion(r.APIVersion)
		if err != nil {
			panic(err)
		}
		gv.Group = v
		r.APIVersion = gv.String()
	})
}

// Version defines the version component of the API version.
func (d *TypeMetaDie) Version(v string) *TypeMetaDie {
	return d.DieStamp(func(r *metav1.TypeMeta) {
		gv, err := schema.ParseGroupVersion(r.APIVersion)
		if err != nil {
			panic(err)
		}
		gv.Version = v
		r.APIVersion = gv.String()
	})
}
