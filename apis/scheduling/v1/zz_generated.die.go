//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2021-2022 the original author or authors.

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

// Code generated by diegen. DO NOT EDIT.

package v1

import (
	metav1 "dies.dev/apis/meta/v1"
	json "encoding/json"
	fmtx "fmt"
	corev1 "k8s.io/api/core/v1"
	schedulingv1 "k8s.io/api/scheduling/v1"
	unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	jsonpath "k8s.io/client-go/util/jsonpath"
	osx "os"
	reflectx "reflect"
	yaml "sigs.k8s.io/yaml"
)

var PriorityClassBlank = (&PriorityClassDie{}).DieFeed(schedulingv1.PriorityClass{})

type PriorityClassDie struct {
	metav1.FrozenObjectMeta
	mutable bool
	r       schedulingv1.PriorityClass
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *PriorityClassDie) DieImmutable(immutable bool) *PriorityClassDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *PriorityClassDie) DieFeed(r schedulingv1.PriorityClass) *PriorityClassDie {
	if d.mutable {
		d.FrozenObjectMeta = metav1.FreezeObjectMeta(r.ObjectMeta)
		d.r = r
		return d
	}
	return &PriorityClassDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *PriorityClassDie) DieFeedPtr(r *schedulingv1.PriorityClass) *PriorityClassDie {
	if r == nil {
		r = &schedulingv1.PriorityClass{}
	}
	return d.DieFeed(*r)
}

// DieFeedJSON returns a new die with the provided JSON. Panics on error.
func (d *PriorityClassDie) DieFeedJSON(j []byte) *PriorityClassDie {
	r := schedulingv1.PriorityClass{}
	if err := json.Unmarshal(j, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAML returns a new die with the provided YAML. Panics on error.
func (d *PriorityClassDie) DieFeedYAML(y []byte) *PriorityClassDie {
	r := schedulingv1.PriorityClass{}
	if err := yaml.Unmarshal(y, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAMLFile returns a new die loading YAML from a file path. Panics on error.
func (d *PriorityClassDie) DieFeedYAMLFile(name string) *PriorityClassDie {
	y, err := osx.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return d.DieFeedYAML(y)
}

// DieFeedRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *PriorityClassDie) DieFeedRawExtension(raw runtime.RawExtension) *PriorityClassDie {
	j, err := json.Marshal(raw)
	if err != nil {
		panic(err)
	}
	return d.DieFeedJSON(j)
}

// DieRelease returns the resource managed by the die.
func (d *PriorityClassDie) DieRelease() schedulingv1.PriorityClass {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *PriorityClassDie) DieReleasePtr() *schedulingv1.PriorityClass {
	r := d.DieRelease()
	return &r
}

// DieReleaseUnstructured returns the resource managed by the die as an unstructured object. Panics on error.
func (d *PriorityClassDie) DieReleaseUnstructured() *unstructured.Unstructured {
	r := d.DieReleasePtr()
	u, err := runtime.DefaultUnstructuredConverter.ToUnstructured(r)
	if err != nil {
		panic(err)
	}
	return &unstructured.Unstructured{
		Object: u,
	}
}

// DieReleaseJSON returns the resource managed by the die as JSON. Panics on error.
func (d *PriorityClassDie) DieReleaseJSON() []byte {
	r := d.DieReleasePtr()
	j, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return j
}

// DieReleaseYAML returns the resource managed by the die as YAML. Panics on error.
func (d *PriorityClassDie) DieReleaseYAML() []byte {
	r := d.DieReleasePtr()
	y, err := yaml.Marshal(r)
	if err != nil {
		panic(err)
	}
	return y
}

// DieReleaseRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *PriorityClassDie) DieReleaseRawExtension() runtime.RawExtension {
	j := d.DieReleaseJSON()
	raw := runtime.RawExtension{}
	if err := json.Unmarshal(j, &raw); err != nil {
		panic(err)
	}
	return raw
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *PriorityClassDie) DieStamp(fn func(r *schedulingv1.PriorityClass)) *PriorityClassDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// Experimental: DieStampAt uses a JSON path (http://goessner.net/articles/JsonPath/) expression to stamp portions of the resource. The callback is invoked with each JSON path match. Panics if the callback function does not accept a single argument of the same type as found on the resource at the target location.
//
// Future iterations will improve type coercion from the resource to the callback argument.
func (d *PriorityClassDie) DieStampAt(jp string, fn interface{}) *PriorityClassDie {
	return d.DieStamp(func(r *schedulingv1.PriorityClass) {
		cp := jsonpath.New("")
		if err := cp.Parse(fmtx.Sprintf("{%s}", jp)); err != nil {
			panic(err)
		}
		cr, err := cp.FindResults(r)
		if err != nil {
			// errors are expected if a path is not found
			return
		}
		for _, cv := range cr[0] {
			args := []reflectx.Value{cv}
			reflectx.ValueOf(fn).Call(args)
		}
	})
}

// DieWith returns a new die after passing the current die to the callback function. The passed die is mutable.
func (d *PriorityClassDie) DieWith(fn func(d *PriorityClassDie)) *PriorityClassDie {
	nd := PriorityClassBlank.DieFeed(d.DieRelease()).DieImmutable(false)
	fn(nd)
	return d.DieFeed(nd.DieRelease())
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *PriorityClassDie) DeepCopy() *PriorityClassDie {
	r := *d.r.DeepCopy()
	return &PriorityClassDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

var _ runtime.Object = (*PriorityClassDie)(nil)

func (d *PriorityClassDie) DeepCopyObject() runtime.Object {
	return d.r.DeepCopy()
}

func (d *PriorityClassDie) GetObjectKind() schema.ObjectKind {
	r := d.DieRelease()
	return r.GetObjectKind()
}

func (d *PriorityClassDie) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.r)
}

func (d *PriorityClassDie) UnmarshalJSON(b []byte) error {
	if d == PriorityClassBlank {
		return fmtx.Errorf("cannot unmarshal into the blank die, create a copy first")
	}
	if !d.mutable {
		return fmtx.Errorf("cannot unmarshal into immutable dies, create a mutable version first")
	}
	r := &schedulingv1.PriorityClass{}
	err := json.Unmarshal(b, r)
	*d = *d.DieFeed(*r)
	return err
}

// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
func (d *PriorityClassDie) APIVersion(v string) *PriorityClassDie {
	return d.DieStamp(func(r *schedulingv1.PriorityClass) {
		r.APIVersion = v
	})
}

// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
func (d *PriorityClassDie) Kind(v string) *PriorityClassDie {
	return d.DieStamp(func(r *schedulingv1.PriorityClass) {
		r.Kind = v
	})
}

// MetadataDie stamps the resource's ObjectMeta field with a mutable die.
func (d *PriorityClassDie) MetadataDie(fn func(d *metav1.ObjectMetaDie)) *PriorityClassDie {
	return d.DieStamp(func(r *schedulingv1.PriorityClass) {
		d := metav1.ObjectMetaBlank.DieImmutable(false).DieFeed(r.ObjectMeta)
		fn(d)
		r.ObjectMeta = d.DieRelease()
	})
}

// value represents the integer value of this priority class. This is the actual priority that pods receive when they have the name of this class in their pod spec.
func (d *PriorityClassDie) Value(v int32) *PriorityClassDie {
	return d.DieStamp(func(r *schedulingv1.PriorityClass) {
		r.Value = v
	})
}

// globalDefault specifies whether this PriorityClass should be considered as the default priority for pods that do not have any priority class. Only one PriorityClass can be marked as `globalDefault`. However, if more than one PriorityClasses exists with their `globalDefault` field set to true, the smallest value of such global default PriorityClasses will be used as the default priority.
func (d *PriorityClassDie) GlobalDefault(v bool) *PriorityClassDie {
	return d.DieStamp(func(r *schedulingv1.PriorityClass) {
		r.GlobalDefault = v
	})
}

// description is an arbitrary string that usually provides guidelines on when this priority class should be used.
func (d *PriorityClassDie) Description(v string) *PriorityClassDie {
	return d.DieStamp(func(r *schedulingv1.PriorityClass) {
		r.Description = v
	})
}

// preemptionPolicy is the Policy for preempting pods with lower priority. One of Never, PreemptLowerPriority. Defaults to PreemptLowerPriority if unset.
func (d *PriorityClassDie) PreemptionPolicy(v *corev1.PreemptionPolicy) *PriorityClassDie {
	return d.DieStamp(func(r *schedulingv1.PriorityClass) {
		r.PreemptionPolicy = v
	})
}
