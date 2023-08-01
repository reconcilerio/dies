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
	unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	jsonpath "k8s.io/client-go/util/jsonpath"
	"k8s.io/kube-aggregator/pkg/apis/apiregistration"
	osx "os"
	reflectx "reflect"
	yaml "sigs.k8s.io/yaml"
)

var APIServiceBlank = (&APIServiceDie{}).DieFeed(apiregistration.APIService{})

type APIServiceDie struct {
	metav1.FrozenObjectMeta
	mutable bool
	r       apiregistration.APIService
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *APIServiceDie) DieImmutable(immutable bool) *APIServiceDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *APIServiceDie) DieFeed(r apiregistration.APIService) *APIServiceDie {
	if d.mutable {
		d.FrozenObjectMeta = metav1.FreezeObjectMeta(r.ObjectMeta)
		d.r = r
		return d
	}
	return &APIServiceDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *APIServiceDie) DieFeedPtr(r *apiregistration.APIService) *APIServiceDie {
	if r == nil {
		r = &apiregistration.APIService{}
	}
	return d.DieFeed(*r)
}

// DieFeedJSON returns a new die with the provided JSON. Panics on error.
func (d *APIServiceDie) DieFeedJSON(j []byte) *APIServiceDie {
	r := apiregistration.APIService{}
	if err := json.Unmarshal(j, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAML returns a new die with the provided YAML. Panics on error.
func (d *APIServiceDie) DieFeedYAML(y []byte) *APIServiceDie {
	r := apiregistration.APIService{}
	if err := yaml.Unmarshal(y, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAMLFile returns a new die loading YAML from a file path. Panics on error.
func (d *APIServiceDie) DieFeedYAMLFile(name string) *APIServiceDie {
	y, err := osx.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return d.DieFeedYAML(y)
}

// DieFeedRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *APIServiceDie) DieFeedRawExtension(raw runtime.RawExtension) *APIServiceDie {
	j, err := json.Marshal(raw)
	if err != nil {
		panic(err)
	}
	return d.DieFeedJSON(j)
}

// DieRelease returns the resource managed by the die.
func (d *APIServiceDie) DieRelease() apiregistration.APIService {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *APIServiceDie) DieReleasePtr() *apiregistration.APIService {
	r := d.DieRelease()
	return &r
}

// DieReleaseUnstructured returns the resource managed by the die as an unstructured object. Panics on error.
func (d *APIServiceDie) DieReleaseUnstructured() *unstructured.Unstructured {
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
func (d *APIServiceDie) DieReleaseJSON() []byte {
	r := d.DieReleasePtr()
	j, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return j
}

// DieReleaseYAML returns the resource managed by the die as YAML. Panics on error.
func (d *APIServiceDie) DieReleaseYAML() []byte {
	r := d.DieReleasePtr()
	y, err := yaml.Marshal(r)
	if err != nil {
		panic(err)
	}
	return y
}

// DieReleaseRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *APIServiceDie) DieReleaseRawExtension() runtime.RawExtension {
	j := d.DieReleaseJSON()
	raw := runtime.RawExtension{}
	if err := json.Unmarshal(j, &raw); err != nil {
		panic(err)
	}
	return raw
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *APIServiceDie) DieStamp(fn func(r *apiregistration.APIService)) *APIServiceDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// Experimental: DieStampAt uses a JSON path (http://goessner.net/articles/JsonPath/) expression to stamp portions of the resource. The callback is invoked with each JSON path match. Panics if the callback function does not accept a single argument of the same type or a pointer to that type as found on the resource at the target location.
//
// Future iterations will improve type coercion from the resource to the callback argument.
func (d *APIServiceDie) DieStampAt(jp string, fn interface{}) *APIServiceDie {
	return d.DieStamp(func(r *apiregistration.APIService) {
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
func (d *APIServiceDie) DieWith(fn func(d *APIServiceDie)) *APIServiceDie {
	nd := APIServiceBlank.DieFeed(d.DieRelease()).DieImmutable(false)
	fn(nd)
	return d.DieFeed(nd.DieRelease())
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *APIServiceDie) DeepCopy() *APIServiceDie {
	r := *d.r.DeepCopy()
	return &APIServiceDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

var _ runtime.Object = (*APIServiceDie)(nil)

func (d *APIServiceDie) DeepCopyObject() runtime.Object {
	return d.r.DeepCopy()
}

func (d *APIServiceDie) GetObjectKind() schema.ObjectKind {
	r := d.DieRelease()
	return r.GetObjectKind()
}

func (d *APIServiceDie) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.r)
}

func (d *APIServiceDie) UnmarshalJSON(b []byte) error {
	if d == APIServiceBlank {
		return fmtx.Errorf("cannot unmarshal into the blank die, create a copy first")
	}
	if !d.mutable {
		return fmtx.Errorf("cannot unmarshal into immutable dies, create a mutable version first")
	}
	r := &apiregistration.APIService{}
	err := json.Unmarshal(b, r)
	*d = *d.DieFeed(*r)
	return err
}

// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
func (d *APIServiceDie) APIVersion(v string) *APIServiceDie {
	return d.DieStamp(func(r *apiregistration.APIService) {
		r.APIVersion = v
	})
}

// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
func (d *APIServiceDie) Kind(v string) *APIServiceDie {
	return d.DieStamp(func(r *apiregistration.APIService) {
		r.Kind = v
	})
}

// MetadataDie stamps the resource's ObjectMeta field with a mutable die.
func (d *APIServiceDie) MetadataDie(fn func(d *metav1.ObjectMetaDie)) *APIServiceDie {
	return d.DieStamp(func(r *apiregistration.APIService) {
		d := metav1.ObjectMetaBlank.DieImmutable(false).DieFeed(r.ObjectMeta)
		fn(d)
		r.ObjectMeta = d.DieRelease()
	})
}

// SpecDie stamps the resource's spec field with a mutable die.
func (d *APIServiceDie) SpecDie(fn func(d *APIServiceSpecDie)) *APIServiceDie {
	return d.DieStamp(func(r *apiregistration.APIService) {
		d := APIServiceSpecBlank.DieImmutable(false).DieFeed(r.Spec)
		fn(d)
		r.Spec = d.DieRelease()
	})
}

// StatusDie stamps the resource's status field with a mutable die.
func (d *APIServiceDie) StatusDie(fn func(d *APIServiceStatusDie)) *APIServiceDie {
	return d.DieStamp(func(r *apiregistration.APIService) {
		d := APIServiceStatusBlank.DieImmutable(false).DieFeed(r.Status)
		fn(d)
		r.Status = d.DieRelease()
	})
}

// Spec contains information for locating and communicating with a server
func (d *APIServiceDie) Spec(v apiregistration.APIServiceSpec) *APIServiceDie {
	return d.DieStamp(func(r *apiregistration.APIService) {
		r.Spec = v
	})
}

// Status contains derived information about an API server
func (d *APIServiceDie) Status(v apiregistration.APIServiceStatus) *APIServiceDie {
	return d.DieStamp(func(r *apiregistration.APIService) {
		r.Status = v
	})
}

var APIServiceSpecBlank = (&APIServiceSpecDie{}).DieFeed(apiregistration.APIServiceSpec{})

type APIServiceSpecDie struct {
	mutable bool
	r       apiregistration.APIServiceSpec
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *APIServiceSpecDie) DieImmutable(immutable bool) *APIServiceSpecDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *APIServiceSpecDie) DieFeed(r apiregistration.APIServiceSpec) *APIServiceSpecDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &APIServiceSpecDie{
		mutable: d.mutable,
		r:       r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *APIServiceSpecDie) DieFeedPtr(r *apiregistration.APIServiceSpec) *APIServiceSpecDie {
	if r == nil {
		r = &apiregistration.APIServiceSpec{}
	}
	return d.DieFeed(*r)
}

// DieFeedJSON returns a new die with the provided JSON. Panics on error.
func (d *APIServiceSpecDie) DieFeedJSON(j []byte) *APIServiceSpecDie {
	r := apiregistration.APIServiceSpec{}
	if err := json.Unmarshal(j, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAML returns a new die with the provided YAML. Panics on error.
func (d *APIServiceSpecDie) DieFeedYAML(y []byte) *APIServiceSpecDie {
	r := apiregistration.APIServiceSpec{}
	if err := yaml.Unmarshal(y, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAMLFile returns a new die loading YAML from a file path. Panics on error.
func (d *APIServiceSpecDie) DieFeedYAMLFile(name string) *APIServiceSpecDie {
	y, err := osx.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return d.DieFeedYAML(y)
}

// DieFeedRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *APIServiceSpecDie) DieFeedRawExtension(raw runtime.RawExtension) *APIServiceSpecDie {
	j, err := json.Marshal(raw)
	if err != nil {
		panic(err)
	}
	return d.DieFeedJSON(j)
}

// DieRelease returns the resource managed by the die.
func (d *APIServiceSpecDie) DieRelease() apiregistration.APIServiceSpec {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *APIServiceSpecDie) DieReleasePtr() *apiregistration.APIServiceSpec {
	r := d.DieRelease()
	return &r
}

// DieReleaseJSON returns the resource managed by the die as JSON. Panics on error.
func (d *APIServiceSpecDie) DieReleaseJSON() []byte {
	r := d.DieReleasePtr()
	j, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return j
}

// DieReleaseYAML returns the resource managed by the die as YAML. Panics on error.
func (d *APIServiceSpecDie) DieReleaseYAML() []byte {
	r := d.DieReleasePtr()
	y, err := yaml.Marshal(r)
	if err != nil {
		panic(err)
	}
	return y
}

// DieReleaseRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *APIServiceSpecDie) DieReleaseRawExtension() runtime.RawExtension {
	j := d.DieReleaseJSON()
	raw := runtime.RawExtension{}
	if err := json.Unmarshal(j, &raw); err != nil {
		panic(err)
	}
	return raw
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *APIServiceSpecDie) DieStamp(fn func(r *apiregistration.APIServiceSpec)) *APIServiceSpecDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// Experimental: DieStampAt uses a JSON path (http://goessner.net/articles/JsonPath/) expression to stamp portions of the resource. The callback is invoked with each JSON path match. Panics if the callback function does not accept a single argument of the same type or a pointer to that type as found on the resource at the target location.
//
// Future iterations will improve type coercion from the resource to the callback argument.
func (d *APIServiceSpecDie) DieStampAt(jp string, fn interface{}) *APIServiceSpecDie {
	return d.DieStamp(func(r *apiregistration.APIServiceSpec) {
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
func (d *APIServiceSpecDie) DieWith(fn func(d *APIServiceSpecDie)) *APIServiceSpecDie {
	nd := APIServiceSpecBlank.DieFeed(d.DieRelease()).DieImmutable(false)
	fn(nd)
	return d.DieFeed(nd.DieRelease())
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *APIServiceSpecDie) DeepCopy() *APIServiceSpecDie {
	r := *d.r.DeepCopy()
	return &APIServiceSpecDie{
		mutable: d.mutable,
		r:       r,
	}
}

// Service is a reference to the service for this API server.  It must communicate on port 443. If the Service is nil, that means the handling for the API groupversion is handled locally on this server. The call will simply delegate to the normal handler chain to be fulfilled.
func (d *APIServiceSpecDie) Service(v *apiregistration.ServiceReference) *APIServiceSpecDie {
	return d.DieStamp(func(r *apiregistration.APIServiceSpec) {
		r.Service = v
	})
}

// Group is the API group name this server hosts
func (d *APIServiceSpecDie) Group(v string) *APIServiceSpecDie {
	return d.DieStamp(func(r *apiregistration.APIServiceSpec) {
		r.Group = v
	})
}

// Version is the API version this server hosts.  For example, "v1"
func (d *APIServiceSpecDie) Version(v string) *APIServiceSpecDie {
	return d.DieStamp(func(r *apiregistration.APIServiceSpec) {
		r.Version = v
	})
}

// InsecureSkipTLSVerify disables TLS certificate verification when communicating with this server. This is strongly discouraged.  You should use the CABundle instead.
func (d *APIServiceSpecDie) InsecureSkipTLSVerify(v bool) *APIServiceSpecDie {
	return d.DieStamp(func(r *apiregistration.APIServiceSpec) {
		r.InsecureSkipTLSVerify = v
	})
}

// CABundle is a PEM encoded CA bundle which will be used to validate an API server's serving certificate. If unspecified, system trust roots on the apiserver are used.
func (d *APIServiceSpecDie) CABundle(v []byte) *APIServiceSpecDie {
	return d.DieStamp(func(r *apiregistration.APIServiceSpec) {
		r.CABundle = v
	})
}

// GroupPriorityMinimum is the priority this group should have at least. Higher priority means that the group is preferred by clients over lower priority ones. Note that other versions of this group might specify even higher GroupPriorityMininum values such that the whole group gets a higher priority. The primary sort is based on GroupPriorityMinimum, ordered highest number to lowest (20 before 10). The secondary sort is based on the alphabetical comparison of the name of the object.  (v1.bar before v1.foo) We'd recommend something like: *.k8s.io (except extensions) at 18000 and PaaSes (OpenShift, Deis) are recommended to be in the 2000s
func (d *APIServiceSpecDie) GroupPriorityMinimum(v int32) *APIServiceSpecDie {
	return d.DieStamp(func(r *apiregistration.APIServiceSpec) {
		r.GroupPriorityMinimum = v
	})
}

// VersionPriority controls the ordering of this API version inside of its group.  Must be greater than zero. The primary sort is based on VersionPriority, ordered highest to lowest (20 before 10). Since it's inside of a group, the number can be small, probably in the 10s. In case of equal version priorities, the version string will be used to compute the order inside a group. If the version string is "kube-like", it will sort above non "kube-like" version strings, which are ordered lexicographically. "Kube-like" versions start with a "v", then are followed by a number (the major version), then optionally the string "alpha" or "beta" and another number (the minor version). These are sorted first by GA > beta > alpha (where GA is a version with no suffix such as beta or alpha), and then by comparing major version, then minor version. An example sorted list of versions: v10, v2, v1, v11beta2, v10beta3, v3beta1, v12alpha1, v11alpha2, foo1, foo10.
func (d *APIServiceSpecDie) VersionPriority(v int32) *APIServiceSpecDie {
	return d.DieStamp(func(r *apiregistration.APIServiceSpec) {
		r.VersionPriority = v
	})
}

var ServiceReferenceBlank = (&ServiceReferenceDie{}).DieFeed(apiregistration.ServiceReference{})

type ServiceReferenceDie struct {
	mutable bool
	r       apiregistration.ServiceReference
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *ServiceReferenceDie) DieImmutable(immutable bool) *ServiceReferenceDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *ServiceReferenceDie) DieFeed(r apiregistration.ServiceReference) *ServiceReferenceDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &ServiceReferenceDie{
		mutable: d.mutable,
		r:       r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *ServiceReferenceDie) DieFeedPtr(r *apiregistration.ServiceReference) *ServiceReferenceDie {
	if r == nil {
		r = &apiregistration.ServiceReference{}
	}
	return d.DieFeed(*r)
}

// DieFeedJSON returns a new die with the provided JSON. Panics on error.
func (d *ServiceReferenceDie) DieFeedJSON(j []byte) *ServiceReferenceDie {
	r := apiregistration.ServiceReference{}
	if err := json.Unmarshal(j, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAML returns a new die with the provided YAML. Panics on error.
func (d *ServiceReferenceDie) DieFeedYAML(y []byte) *ServiceReferenceDie {
	r := apiregistration.ServiceReference{}
	if err := yaml.Unmarshal(y, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAMLFile returns a new die loading YAML from a file path. Panics on error.
func (d *ServiceReferenceDie) DieFeedYAMLFile(name string) *ServiceReferenceDie {
	y, err := osx.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return d.DieFeedYAML(y)
}

// DieFeedRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *ServiceReferenceDie) DieFeedRawExtension(raw runtime.RawExtension) *ServiceReferenceDie {
	j, err := json.Marshal(raw)
	if err != nil {
		panic(err)
	}
	return d.DieFeedJSON(j)
}

// DieRelease returns the resource managed by the die.
func (d *ServiceReferenceDie) DieRelease() apiregistration.ServiceReference {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *ServiceReferenceDie) DieReleasePtr() *apiregistration.ServiceReference {
	r := d.DieRelease()
	return &r
}

// DieReleaseJSON returns the resource managed by the die as JSON. Panics on error.
func (d *ServiceReferenceDie) DieReleaseJSON() []byte {
	r := d.DieReleasePtr()
	j, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return j
}

// DieReleaseYAML returns the resource managed by the die as YAML. Panics on error.
func (d *ServiceReferenceDie) DieReleaseYAML() []byte {
	r := d.DieReleasePtr()
	y, err := yaml.Marshal(r)
	if err != nil {
		panic(err)
	}
	return y
}

// DieReleaseRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *ServiceReferenceDie) DieReleaseRawExtension() runtime.RawExtension {
	j := d.DieReleaseJSON()
	raw := runtime.RawExtension{}
	if err := json.Unmarshal(j, &raw); err != nil {
		panic(err)
	}
	return raw
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *ServiceReferenceDie) DieStamp(fn func(r *apiregistration.ServiceReference)) *ServiceReferenceDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// Experimental: DieStampAt uses a JSON path (http://goessner.net/articles/JsonPath/) expression to stamp portions of the resource. The callback is invoked with each JSON path match. Panics if the callback function does not accept a single argument of the same type or a pointer to that type as found on the resource at the target location.
//
// Future iterations will improve type coercion from the resource to the callback argument.
func (d *ServiceReferenceDie) DieStampAt(jp string, fn interface{}) *ServiceReferenceDie {
	return d.DieStamp(func(r *apiregistration.ServiceReference) {
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
func (d *ServiceReferenceDie) DieWith(fn func(d *ServiceReferenceDie)) *ServiceReferenceDie {
	nd := ServiceReferenceBlank.DieFeed(d.DieRelease()).DieImmutable(false)
	fn(nd)
	return d.DieFeed(nd.DieRelease())
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *ServiceReferenceDie) DeepCopy() *ServiceReferenceDie {
	r := *d.r.DeepCopy()
	return &ServiceReferenceDie{
		mutable: d.mutable,
		r:       r,
	}
}

// Namespace is the namespace of the service
func (d *ServiceReferenceDie) Namespace(v string) *ServiceReferenceDie {
	return d.DieStamp(func(r *apiregistration.ServiceReference) {
		r.Namespace = v
	})
}

// Name is the name of the service
func (d *ServiceReferenceDie) Name(v string) *ServiceReferenceDie {
	return d.DieStamp(func(r *apiregistration.ServiceReference) {
		r.Name = v
	})
}

// If specified, the port on the service that hosting the service. Default to 443 for backward compatibility. `port` should be a valid port number (1-65535, inclusive).
func (d *ServiceReferenceDie) Port(v int32) *ServiceReferenceDie {
	return d.DieStamp(func(r *apiregistration.ServiceReference) {
		r.Port = v
	})
}

var APIServiceStatusBlank = (&APIServiceStatusDie{}).DieFeed(apiregistration.APIServiceStatus{})

type APIServiceStatusDie struct {
	mutable bool
	r       apiregistration.APIServiceStatus
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *APIServiceStatusDie) DieImmutable(immutable bool) *APIServiceStatusDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *APIServiceStatusDie) DieFeed(r apiregistration.APIServiceStatus) *APIServiceStatusDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &APIServiceStatusDie{
		mutable: d.mutable,
		r:       r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *APIServiceStatusDie) DieFeedPtr(r *apiregistration.APIServiceStatus) *APIServiceStatusDie {
	if r == nil {
		r = &apiregistration.APIServiceStatus{}
	}
	return d.DieFeed(*r)
}

// DieFeedJSON returns a new die with the provided JSON. Panics on error.
func (d *APIServiceStatusDie) DieFeedJSON(j []byte) *APIServiceStatusDie {
	r := apiregistration.APIServiceStatus{}
	if err := json.Unmarshal(j, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAML returns a new die with the provided YAML. Panics on error.
func (d *APIServiceStatusDie) DieFeedYAML(y []byte) *APIServiceStatusDie {
	r := apiregistration.APIServiceStatus{}
	if err := yaml.Unmarshal(y, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAMLFile returns a new die loading YAML from a file path. Panics on error.
func (d *APIServiceStatusDie) DieFeedYAMLFile(name string) *APIServiceStatusDie {
	y, err := osx.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return d.DieFeedYAML(y)
}

// DieFeedRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *APIServiceStatusDie) DieFeedRawExtension(raw runtime.RawExtension) *APIServiceStatusDie {
	j, err := json.Marshal(raw)
	if err != nil {
		panic(err)
	}
	return d.DieFeedJSON(j)
}

// DieRelease returns the resource managed by the die.
func (d *APIServiceStatusDie) DieRelease() apiregistration.APIServiceStatus {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *APIServiceStatusDie) DieReleasePtr() *apiregistration.APIServiceStatus {
	r := d.DieRelease()
	return &r
}

// DieReleaseJSON returns the resource managed by the die as JSON. Panics on error.
func (d *APIServiceStatusDie) DieReleaseJSON() []byte {
	r := d.DieReleasePtr()
	j, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return j
}

// DieReleaseYAML returns the resource managed by the die as YAML. Panics on error.
func (d *APIServiceStatusDie) DieReleaseYAML() []byte {
	r := d.DieReleasePtr()
	y, err := yaml.Marshal(r)
	if err != nil {
		panic(err)
	}
	return y
}

// DieReleaseRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *APIServiceStatusDie) DieReleaseRawExtension() runtime.RawExtension {
	j := d.DieReleaseJSON()
	raw := runtime.RawExtension{}
	if err := json.Unmarshal(j, &raw); err != nil {
		panic(err)
	}
	return raw
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *APIServiceStatusDie) DieStamp(fn func(r *apiregistration.APIServiceStatus)) *APIServiceStatusDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// Experimental: DieStampAt uses a JSON path (http://goessner.net/articles/JsonPath/) expression to stamp portions of the resource. The callback is invoked with each JSON path match. Panics if the callback function does not accept a single argument of the same type or a pointer to that type as found on the resource at the target location.
//
// Future iterations will improve type coercion from the resource to the callback argument.
func (d *APIServiceStatusDie) DieStampAt(jp string, fn interface{}) *APIServiceStatusDie {
	return d.DieStamp(func(r *apiregistration.APIServiceStatus) {
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
func (d *APIServiceStatusDie) DieWith(fn func(d *APIServiceStatusDie)) *APIServiceStatusDie {
	nd := APIServiceStatusBlank.DieFeed(d.DieRelease()).DieImmutable(false)
	fn(nd)
	return d.DieFeed(nd.DieRelease())
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *APIServiceStatusDie) DeepCopy() *APIServiceStatusDie {
	r := *d.r.DeepCopy()
	return &APIServiceStatusDie{
		mutable: d.mutable,
		r:       r,
	}
}

// Current service state of apiService.
func (d *APIServiceStatusDie) Conditions(v ...apiregistration.APIServiceCondition) *APIServiceStatusDie {
	return d.DieStamp(func(r *apiregistration.APIServiceStatus) {
		r.Conditions = v
	})
}
