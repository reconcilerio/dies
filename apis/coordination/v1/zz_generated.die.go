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
	coordinationv1 "k8s.io/api/coordination/v1"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	jsonpath "k8s.io/client-go/util/jsonpath"
	osx "os"
	reflectx "reflect"
	yaml "sigs.k8s.io/yaml"
)

var LeaseBlank = (&LeaseDie{}).DieFeed(coordinationv1.Lease{})

type LeaseDie struct {
	metav1.FrozenObjectMeta
	mutable bool
	r       coordinationv1.Lease
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *LeaseDie) DieImmutable(immutable bool) *LeaseDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *LeaseDie) DieFeed(r coordinationv1.Lease) *LeaseDie {
	if d.mutable {
		d.FrozenObjectMeta = metav1.FreezeObjectMeta(r.ObjectMeta)
		d.r = r
		return d
	}
	return &LeaseDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *LeaseDie) DieFeedPtr(r *coordinationv1.Lease) *LeaseDie {
	if r == nil {
		r = &coordinationv1.Lease{}
	}
	return d.DieFeed(*r)
}

// DieFeedJSON returns a new die with the provided JSON. Panics on error.
func (d *LeaseDie) DieFeedJSON(j []byte) *LeaseDie {
	r := coordinationv1.Lease{}
	if err := json.Unmarshal(j, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAML returns a new die with the provided YAML. Panics on error.
func (d *LeaseDie) DieFeedYAML(y []byte) *LeaseDie {
	r := coordinationv1.Lease{}
	if err := yaml.Unmarshal(y, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAMLFile returns a new die loading YAML from a file path. Panics on error.
func (d *LeaseDie) DieFeedYAMLFile(name string) *LeaseDie {
	y, err := osx.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return d.DieFeedYAML(y)
}

// DieFeedRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *LeaseDie) DieFeedRawExtension(raw runtime.RawExtension) *LeaseDie {
	j, err := json.Marshal(raw)
	if err != nil {
		panic(err)
	}
	return d.DieFeedJSON(j)
}

// DieRelease returns the resource managed by the die.
func (d *LeaseDie) DieRelease() coordinationv1.Lease {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *LeaseDie) DieReleasePtr() *coordinationv1.Lease {
	r := d.DieRelease()
	return &r
}

// DieReleaseUnstructured returns the resource managed by the die as an unstructured object. Panics on error.
func (d *LeaseDie) DieReleaseUnstructured() *unstructured.Unstructured {
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
func (d *LeaseDie) DieReleaseJSON() []byte {
	r := d.DieReleasePtr()
	j, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return j
}

// DieReleaseYAML returns the resource managed by the die as YAML. Panics on error.
func (d *LeaseDie) DieReleaseYAML() []byte {
	r := d.DieReleasePtr()
	y, err := yaml.Marshal(r)
	if err != nil {
		panic(err)
	}
	return y
}

// DieReleaseRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *LeaseDie) DieReleaseRawExtension() runtime.RawExtension {
	j := d.DieReleaseJSON()
	raw := runtime.RawExtension{}
	if err := json.Unmarshal(j, &raw); err != nil {
		panic(err)
	}
	return raw
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *LeaseDie) DieStamp(fn func(r *coordinationv1.Lease)) *LeaseDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// Experimental: DieStampAt uses a JSON path (http://goessner.net/articles/JsonPath/) expression to stamp portions of the resource. The callback is invoked with each JSON path match. Panics if the callback function does not accept a single argument of the same type or a pointer to that type as found on the resource at the target location.
//
// Future iterations will improve type coercion from the resource to the callback argument.
func (d *LeaseDie) DieStampAt(jp string, fn interface{}) *LeaseDie {
	return d.DieStamp(func(r *coordinationv1.Lease) {
		if ni := reflectx.ValueOf(fn).Type().NumIn(); ni != 1 {
			panic(fmtx.Errorf("callback function must have 1 input parameters, found %d", ni))
		}
		if no := reflectx.ValueOf(fn).Type().NumOut(); no != 0 {
			panic(fmtx.Errorf("callback function must have 0 output parameters, found %d", no))
		}

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
			arg0t := reflectx.ValueOf(fn).Type().In(0)

			var args []reflectx.Value
			if cv.Type().AssignableTo(arg0t) {
				args = []reflectx.Value{cv}
			} else if cv.CanAddr() && cv.Addr().Type().AssignableTo(arg0t) {
				args = []reflectx.Value{cv.Addr()}
			} else {
				panic(fmtx.Errorf("callback function must accept value of type %q, found type %q", cv.Type(), arg0t))
			}

			reflectx.ValueOf(fn).Call(args)
		}
	})
}

// DieWith returns a new die after passing the current die to the callback function. The passed die is mutable.
func (d *LeaseDie) DieWith(fns ...func(d *LeaseDie)) *LeaseDie {
	nd := LeaseBlank.DieFeed(d.DieRelease()).DieImmutable(false)
	for _, fn := range fns {
		if fn != nil {
			fn(nd)
		}
	}
	return d.DieFeed(nd.DieRelease())
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *LeaseDie) DeepCopy() *LeaseDie {
	r := *d.r.DeepCopy()
	return &LeaseDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

var _ runtime.Object = (*LeaseDie)(nil)

func (d *LeaseDie) DeepCopyObject() runtime.Object {
	return d.r.DeepCopy()
}

func (d *LeaseDie) GetObjectKind() schema.ObjectKind {
	r := d.DieRelease()
	return r.GetObjectKind()
}

func (d *LeaseDie) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.r)
}

func (d *LeaseDie) UnmarshalJSON(b []byte) error {
	if d == LeaseBlank {
		return fmtx.Errorf("cannot unmarshal into the blank die, create a copy first")
	}
	if !d.mutable {
		return fmtx.Errorf("cannot unmarshal into immutable dies, create a mutable version first")
	}
	r := &coordinationv1.Lease{}
	err := json.Unmarshal(b, r)
	*d = *d.DieFeed(*r)
	return err
}

// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
func (d *LeaseDie) APIVersion(v string) *LeaseDie {
	return d.DieStamp(func(r *coordinationv1.Lease) {
		r.APIVersion = v
	})
}

// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
func (d *LeaseDie) Kind(v string) *LeaseDie {
	return d.DieStamp(func(r *coordinationv1.Lease) {
		r.Kind = v
	})
}

// MetadataDie stamps the resource's ObjectMeta field with a mutable die.
func (d *LeaseDie) MetadataDie(fn func(d *metav1.ObjectMetaDie)) *LeaseDie {
	return d.DieStamp(func(r *coordinationv1.Lease) {
		d := metav1.ObjectMetaBlank.DieImmutable(false).DieFeed(r.ObjectMeta)
		fn(d)
		r.ObjectMeta = d.DieRelease()
	})
}

// SpecDie stamps the resource's spec field with a mutable die.
func (d *LeaseDie) SpecDie(fn func(d *LeaseSpecDie)) *LeaseDie {
	return d.DieStamp(func(r *coordinationv1.Lease) {
		d := LeaseSpecBlank.DieImmutable(false).DieFeed(r.Spec)
		fn(d)
		r.Spec = d.DieRelease()
	})
}

// spec contains the specification of the Lease.
//
// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
func (d *LeaseDie) Spec(v coordinationv1.LeaseSpec) *LeaseDie {
	return d.DieStamp(func(r *coordinationv1.Lease) {
		r.Spec = v
	})
}

var LeaseSpecBlank = (&LeaseSpecDie{}).DieFeed(coordinationv1.LeaseSpec{})

type LeaseSpecDie struct {
	mutable bool
	r       coordinationv1.LeaseSpec
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *LeaseSpecDie) DieImmutable(immutable bool) *LeaseSpecDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *LeaseSpecDie) DieFeed(r coordinationv1.LeaseSpec) *LeaseSpecDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &LeaseSpecDie{
		mutable: d.mutable,
		r:       r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *LeaseSpecDie) DieFeedPtr(r *coordinationv1.LeaseSpec) *LeaseSpecDie {
	if r == nil {
		r = &coordinationv1.LeaseSpec{}
	}
	return d.DieFeed(*r)
}

// DieFeedJSON returns a new die with the provided JSON. Panics on error.
func (d *LeaseSpecDie) DieFeedJSON(j []byte) *LeaseSpecDie {
	r := coordinationv1.LeaseSpec{}
	if err := json.Unmarshal(j, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAML returns a new die with the provided YAML. Panics on error.
func (d *LeaseSpecDie) DieFeedYAML(y []byte) *LeaseSpecDie {
	r := coordinationv1.LeaseSpec{}
	if err := yaml.Unmarshal(y, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAMLFile returns a new die loading YAML from a file path. Panics on error.
func (d *LeaseSpecDie) DieFeedYAMLFile(name string) *LeaseSpecDie {
	y, err := osx.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return d.DieFeedYAML(y)
}

// DieFeedRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *LeaseSpecDie) DieFeedRawExtension(raw runtime.RawExtension) *LeaseSpecDie {
	j, err := json.Marshal(raw)
	if err != nil {
		panic(err)
	}
	return d.DieFeedJSON(j)
}

// DieRelease returns the resource managed by the die.
func (d *LeaseSpecDie) DieRelease() coordinationv1.LeaseSpec {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *LeaseSpecDie) DieReleasePtr() *coordinationv1.LeaseSpec {
	r := d.DieRelease()
	return &r
}

// DieReleaseJSON returns the resource managed by the die as JSON. Panics on error.
func (d *LeaseSpecDie) DieReleaseJSON() []byte {
	r := d.DieReleasePtr()
	j, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return j
}

// DieReleaseYAML returns the resource managed by the die as YAML. Panics on error.
func (d *LeaseSpecDie) DieReleaseYAML() []byte {
	r := d.DieReleasePtr()
	y, err := yaml.Marshal(r)
	if err != nil {
		panic(err)
	}
	return y
}

// DieReleaseRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *LeaseSpecDie) DieReleaseRawExtension() runtime.RawExtension {
	j := d.DieReleaseJSON()
	raw := runtime.RawExtension{}
	if err := json.Unmarshal(j, &raw); err != nil {
		panic(err)
	}
	return raw
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *LeaseSpecDie) DieStamp(fn func(r *coordinationv1.LeaseSpec)) *LeaseSpecDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// Experimental: DieStampAt uses a JSON path (http://goessner.net/articles/JsonPath/) expression to stamp portions of the resource. The callback is invoked with each JSON path match. Panics if the callback function does not accept a single argument of the same type or a pointer to that type as found on the resource at the target location.
//
// Future iterations will improve type coercion from the resource to the callback argument.
func (d *LeaseSpecDie) DieStampAt(jp string, fn interface{}) *LeaseSpecDie {
	return d.DieStamp(func(r *coordinationv1.LeaseSpec) {
		if ni := reflectx.ValueOf(fn).Type().NumIn(); ni != 1 {
			panic(fmtx.Errorf("callback function must have 1 input parameters, found %d", ni))
		}
		if no := reflectx.ValueOf(fn).Type().NumOut(); no != 0 {
			panic(fmtx.Errorf("callback function must have 0 output parameters, found %d", no))
		}

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
			arg0t := reflectx.ValueOf(fn).Type().In(0)

			var args []reflectx.Value
			if cv.Type().AssignableTo(arg0t) {
				args = []reflectx.Value{cv}
			} else if cv.CanAddr() && cv.Addr().Type().AssignableTo(arg0t) {
				args = []reflectx.Value{cv.Addr()}
			} else {
				panic(fmtx.Errorf("callback function must accept value of type %q, found type %q", cv.Type(), arg0t))
			}

			reflectx.ValueOf(fn).Call(args)
		}
	})
}

// DieWith returns a new die after passing the current die to the callback function. The passed die is mutable.
func (d *LeaseSpecDie) DieWith(fns ...func(d *LeaseSpecDie)) *LeaseSpecDie {
	nd := LeaseSpecBlank.DieFeed(d.DieRelease()).DieImmutable(false)
	for _, fn := range fns {
		if fn != nil {
			fn(nd)
		}
	}
	return d.DieFeed(nd.DieRelease())
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *LeaseSpecDie) DeepCopy() *LeaseSpecDie {
	r := *d.r.DeepCopy()
	return &LeaseSpecDie{
		mutable: d.mutable,
		r:       r,
	}
}

// holderIdentity contains the identity of the holder of a current lease.
func (d *LeaseSpecDie) HolderIdentity(v *string) *LeaseSpecDie {
	return d.DieStamp(func(r *coordinationv1.LeaseSpec) {
		r.HolderIdentity = v
	})
}

// leaseDurationSeconds is a duration that candidates for a lease need
//
// to wait to force acquire it. This is measure against time of last
//
// observed renewTime.
func (d *LeaseSpecDie) LeaseDurationSeconds(v *int32) *LeaseSpecDie {
	return d.DieStamp(func(r *coordinationv1.LeaseSpec) {
		r.LeaseDurationSeconds = v
	})
}

// acquireTime is a time when the current lease was acquired.
func (d *LeaseSpecDie) AcquireTime(v *apismetav1.MicroTime) *LeaseSpecDie {
	return d.DieStamp(func(r *coordinationv1.LeaseSpec) {
		r.AcquireTime = v
	})
}

// renewTime is a time when the current holder of a lease has last
//
// updated the lease.
func (d *LeaseSpecDie) RenewTime(v *apismetav1.MicroTime) *LeaseSpecDie {
	return d.DieStamp(func(r *coordinationv1.LeaseSpec) {
		r.RenewTime = v
	})
}

// leaseTransitions is the number of transitions of a lease between
//
// holders.
func (d *LeaseSpecDie) LeaseTransitions(v *int32) *LeaseSpecDie {
	return d.DieStamp(func(r *coordinationv1.LeaseSpec) {
		r.LeaseTransitions = v
	})
}
