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
	eventsv1 "k8s.io/api/events/v1"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	jsonpath "k8s.io/client-go/util/jsonpath"
	osx "os"
	reflectx "reflect"
	yaml "sigs.k8s.io/yaml"
)

var EventBlank = (&EventDie{}).DieFeed(eventsv1.Event{})

type EventDie struct {
	metav1.FrozenObjectMeta
	mutable bool
	r       eventsv1.Event
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *EventDie) DieImmutable(immutable bool) *EventDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *EventDie) DieFeed(r eventsv1.Event) *EventDie {
	if d.mutable {
		d.FrozenObjectMeta = metav1.FreezeObjectMeta(r.ObjectMeta)
		d.r = r
		return d
	}
	return &EventDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *EventDie) DieFeedPtr(r *eventsv1.Event) *EventDie {
	if r == nil {
		r = &eventsv1.Event{}
	}
	return d.DieFeed(*r)
}

// DieFeedJSON returns a new die with the provided JSON. Panics on error.
func (d *EventDie) DieFeedJSON(j []byte) *EventDie {
	r := eventsv1.Event{}
	if err := json.Unmarshal(j, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAML returns a new die with the provided YAML. Panics on error.
func (d *EventDie) DieFeedYAML(y []byte) *EventDie {
	r := eventsv1.Event{}
	if err := yaml.Unmarshal(y, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAMLFile returns a new die loading YAML from a file path. Panics on error.
func (d *EventDie) DieFeedYAMLFile(name string) *EventDie {
	y, err := osx.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return d.DieFeedYAML(y)
}

// DieFeedRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *EventDie) DieFeedRawExtension(raw runtime.RawExtension) *EventDie {
	j, err := json.Marshal(raw)
	if err != nil {
		panic(err)
	}
	return d.DieFeedJSON(j)
}

// DieRelease returns the resource managed by the die.
func (d *EventDie) DieRelease() eventsv1.Event {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *EventDie) DieReleasePtr() *eventsv1.Event {
	r := d.DieRelease()
	return &r
}

// DieReleaseUnstructured returns the resource managed by the die as an unstructured object. Panics on error.
func (d *EventDie) DieReleaseUnstructured() *unstructured.Unstructured {
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
func (d *EventDie) DieReleaseJSON() []byte {
	r := d.DieReleasePtr()
	j, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return j
}

// DieReleaseYAML returns the resource managed by the die as YAML. Panics on error.
func (d *EventDie) DieReleaseYAML() []byte {
	r := d.DieReleasePtr()
	y, err := yaml.Marshal(r)
	if err != nil {
		panic(err)
	}
	return y
}

// DieReleaseRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *EventDie) DieReleaseRawExtension() runtime.RawExtension {
	j := d.DieReleaseJSON()
	raw := runtime.RawExtension{}
	if err := json.Unmarshal(j, &raw); err != nil {
		panic(err)
	}
	return raw
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *EventDie) DieStamp(fn func(r *eventsv1.Event)) *EventDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// Experimental: DieStampAt uses a JSON path (http://goessner.net/articles/JsonPath/) expression to stamp portions of the resource. The callback is invoked with each JSON path match. Panics if the callback function does not accept a single argument of the same type or a pointer to that type as found on the resource at the target location.
//
// Future iterations will improve type coercion from the resource to the callback argument.
func (d *EventDie) DieStampAt(jp string, fn interface{}) *EventDie {
	return d.DieStamp(func(r *eventsv1.Event) {
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
func (d *EventDie) DieWith(fns ...func(d *EventDie)) *EventDie {
	nd := EventBlank.DieFeed(d.DieRelease()).DieImmutable(false)
	for _, fn := range fns {
		if fn != nil {
			fn(nd)
		}
	}
	return d.DieFeed(nd.DieRelease())
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *EventDie) DeepCopy() *EventDie {
	r := *d.r.DeepCopy()
	return &EventDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

var _ runtime.Object = (*EventDie)(nil)

func (d *EventDie) DeepCopyObject() runtime.Object {
	return d.r.DeepCopy()
}

func (d *EventDie) GetObjectKind() schema.ObjectKind {
	r := d.DieRelease()
	return r.GetObjectKind()
}

func (d *EventDie) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.r)
}

func (d *EventDie) UnmarshalJSON(b []byte) error {
	if d == EventBlank {
		return fmtx.Errorf("cannot unmarshal into the blank die, create a copy first")
	}
	if !d.mutable {
		return fmtx.Errorf("cannot unmarshal into immutable dies, create a mutable version first")
	}
	r := &eventsv1.Event{}
	err := json.Unmarshal(b, r)
	*d = *d.DieFeed(*r)
	return err
}

// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
func (d *EventDie) APIVersion(v string) *EventDie {
	return d.DieStamp(func(r *eventsv1.Event) {
		r.APIVersion = v
	})
}

// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
func (d *EventDie) Kind(v string) *EventDie {
	return d.DieStamp(func(r *eventsv1.Event) {
		r.Kind = v
	})
}

// MetadataDie stamps the resource's ObjectMeta field with a mutable die.
func (d *EventDie) MetadataDie(fn func(d *metav1.ObjectMetaDie)) *EventDie {
	return d.DieStamp(func(r *eventsv1.Event) {
		d := metav1.ObjectMetaBlank.DieImmutable(false).DieFeed(r.ObjectMeta)
		fn(d)
		r.ObjectMeta = d.DieRelease()
	})
}

// eventTime is the time when this Event was first observed. It is required.
func (d *EventDie) EventTime(v apismetav1.MicroTime) *EventDie {
	return d.DieStamp(func(r *eventsv1.Event) {
		r.EventTime = v
	})
}

// series is data about the Event series this event represents or nil if it's a singleton Event.
func (d *EventDie) Series(v *eventsv1.EventSeries) *EventDie {
	return d.DieStamp(func(r *eventsv1.Event) {
		r.Series = v
	})
}

// reportingController is the name of the controller that emitted this Event, e.g. `kubernetes.io/kubelet`. This field cannot be empty for new Events.
func (d *EventDie) ReportingController(v string) *EventDie {
	return d.DieStamp(func(r *eventsv1.Event) {
		r.ReportingController = v
	})
}

// reportingInstance is the ID of the controller instance, e.g. `kubelet-xyzf`. This field cannot be empty for new Events and it can have at most 128 characters.
func (d *EventDie) ReportingInstance(v string) *EventDie {
	return d.DieStamp(func(r *eventsv1.Event) {
		r.ReportingInstance = v
	})
}

// action is what action was taken/failed regarding to the regarding object. It is machine-readable. This field cannot be empty for new Events and it can have at most 128 characters.
func (d *EventDie) Action(v string) *EventDie {
	return d.DieStamp(func(r *eventsv1.Event) {
		r.Action = v
	})
}

// reason is why the action was taken. It is human-readable. This field cannot be empty for new Events and it can have at most 128 characters.
func (d *EventDie) Reason(v string) *EventDie {
	return d.DieStamp(func(r *eventsv1.Event) {
		r.Reason = v
	})
}

// regarding contains the object this Event is about. In most cases it's an Object reporting controller implements, e.g. ReplicaSetController implements ReplicaSets and this event is emitted because it acts on some changes in a ReplicaSet object.
func (d *EventDie) Regarding(v corev1.ObjectReference) *EventDie {
	return d.DieStamp(func(r *eventsv1.Event) {
		r.Regarding = v
	})
}

// related is the optional secondary object for more complex actions. E.g. when regarding object triggers a creation or deletion of related object.
func (d *EventDie) Related(v *corev1.ObjectReference) *EventDie {
	return d.DieStamp(func(r *eventsv1.Event) {
		r.Related = v
	})
}

// note is a human-readable description of the status of this operation. Maximal length of the note is 1kB, but libraries should be prepared to handle values up to 64kB.
func (d *EventDie) Note(v string) *EventDie {
	return d.DieStamp(func(r *eventsv1.Event) {
		r.Note = v
	})
}

// type is the type of this event (Normal, Warning), new types could be added in the future. It is machine-readable. This field cannot be empty for new Events.
func (d *EventDie) Type(v string) *EventDie {
	return d.DieStamp(func(r *eventsv1.Event) {
		r.Type = v
	})
}

// deprecatedSource is the deprecated field assuring backward compatibility with core.v1 Event type.
func (d *EventDie) DeprecatedSource(v corev1.EventSource) *EventDie {
	return d.DieStamp(func(r *eventsv1.Event) {
		r.DeprecatedSource = v
	})
}

// deprecatedFirstTimestamp is the deprecated field assuring backward compatibility with core.v1 Event type.
func (d *EventDie) DeprecatedFirstTimestamp(v apismetav1.Time) *EventDie {
	return d.DieStamp(func(r *eventsv1.Event) {
		r.DeprecatedFirstTimestamp = v
	})
}

// deprecatedLastTimestamp is the deprecated field assuring backward compatibility with core.v1 Event type.
func (d *EventDie) DeprecatedLastTimestamp(v apismetav1.Time) *EventDie {
	return d.DieStamp(func(r *eventsv1.Event) {
		r.DeprecatedLastTimestamp = v
	})
}

// deprecatedCount is the deprecated field assuring backward compatibility with core.v1 Event type.
func (d *EventDie) DeprecatedCount(v int32) *EventDie {
	return d.DieStamp(func(r *eventsv1.Event) {
		r.DeprecatedCount = v
	})
}

var EventSeriesBlank = (&EventSeriesDie{}).DieFeed(eventsv1.EventSeries{})

type EventSeriesDie struct {
	mutable bool
	r       eventsv1.EventSeries
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *EventSeriesDie) DieImmutable(immutable bool) *EventSeriesDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *EventSeriesDie) DieFeed(r eventsv1.EventSeries) *EventSeriesDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &EventSeriesDie{
		mutable: d.mutable,
		r:       r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *EventSeriesDie) DieFeedPtr(r *eventsv1.EventSeries) *EventSeriesDie {
	if r == nil {
		r = &eventsv1.EventSeries{}
	}
	return d.DieFeed(*r)
}

// DieFeedJSON returns a new die with the provided JSON. Panics on error.
func (d *EventSeriesDie) DieFeedJSON(j []byte) *EventSeriesDie {
	r := eventsv1.EventSeries{}
	if err := json.Unmarshal(j, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAML returns a new die with the provided YAML. Panics on error.
func (d *EventSeriesDie) DieFeedYAML(y []byte) *EventSeriesDie {
	r := eventsv1.EventSeries{}
	if err := yaml.Unmarshal(y, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAMLFile returns a new die loading YAML from a file path. Panics on error.
func (d *EventSeriesDie) DieFeedYAMLFile(name string) *EventSeriesDie {
	y, err := osx.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return d.DieFeedYAML(y)
}

// DieFeedRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *EventSeriesDie) DieFeedRawExtension(raw runtime.RawExtension) *EventSeriesDie {
	j, err := json.Marshal(raw)
	if err != nil {
		panic(err)
	}
	return d.DieFeedJSON(j)
}

// DieRelease returns the resource managed by the die.
func (d *EventSeriesDie) DieRelease() eventsv1.EventSeries {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *EventSeriesDie) DieReleasePtr() *eventsv1.EventSeries {
	r := d.DieRelease()
	return &r
}

// DieReleaseJSON returns the resource managed by the die as JSON. Panics on error.
func (d *EventSeriesDie) DieReleaseJSON() []byte {
	r := d.DieReleasePtr()
	j, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return j
}

// DieReleaseYAML returns the resource managed by the die as YAML. Panics on error.
func (d *EventSeriesDie) DieReleaseYAML() []byte {
	r := d.DieReleasePtr()
	y, err := yaml.Marshal(r)
	if err != nil {
		panic(err)
	}
	return y
}

// DieReleaseRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *EventSeriesDie) DieReleaseRawExtension() runtime.RawExtension {
	j := d.DieReleaseJSON()
	raw := runtime.RawExtension{}
	if err := json.Unmarshal(j, &raw); err != nil {
		panic(err)
	}
	return raw
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *EventSeriesDie) DieStamp(fn func(r *eventsv1.EventSeries)) *EventSeriesDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// Experimental: DieStampAt uses a JSON path (http://goessner.net/articles/JsonPath/) expression to stamp portions of the resource. The callback is invoked with each JSON path match. Panics if the callback function does not accept a single argument of the same type or a pointer to that type as found on the resource at the target location.
//
// Future iterations will improve type coercion from the resource to the callback argument.
func (d *EventSeriesDie) DieStampAt(jp string, fn interface{}) *EventSeriesDie {
	return d.DieStamp(func(r *eventsv1.EventSeries) {
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
func (d *EventSeriesDie) DieWith(fns ...func(d *EventSeriesDie)) *EventSeriesDie {
	nd := EventSeriesBlank.DieFeed(d.DieRelease()).DieImmutable(false)
	for _, fn := range fns {
		if fn != nil {
			fn(nd)
		}
	}
	return d.DieFeed(nd.DieRelease())
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *EventSeriesDie) DeepCopy() *EventSeriesDie {
	r := *d.r.DeepCopy()
	return &EventSeriesDie{
		mutable: d.mutable,
		r:       r,
	}
}

// count is the number of occurrences in this series up to the last heartbeat time.
func (d *EventSeriesDie) Count(v int32) *EventSeriesDie {
	return d.DieStamp(func(r *eventsv1.EventSeries) {
		r.Count = v
	})
}

// lastObservedTime is the time when last Event from the series was seen before last heartbeat.
func (d *EventSeriesDie) LastObservedTime(v apismetav1.MicroTime) *EventSeriesDie {
	return d.DieStamp(func(r *eventsv1.EventSeries) {
		r.LastObservedTime = v
	})
}
