/*
Copyright 2026 the original author or authors.

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

package base

import (
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utiljson "k8s.io/apimachinery/pkg/util/json"
	"reconciler.io/dies/v2/die/internal/core"
)

// Dies for top-level objects embed one of these instead of DieBase, so that
// they can be used where Kubernetes libraries expect objects:
//
//   - DieTypedBase implements runtime.Object
//   - DieObjectBase also implements metav1.Object and metav1.ObjectMetaAccessor,
//     through the object's metadata
//
// Lists are not dies of their own, but lists of dies: see List.
//
// Their setters return nothing, so they change a die in place: they require
// a mutable die, and panic on an immutable one rather than lose the change.

// DieTyped is a die of any type with type metadata: a Die that is a
// runtime.Object, as a typed die's base makes it. die.DieTypedOf has the
// rest, for generic code.
type DieTyped interface {
	Die
	runtime.Object
}

// DieObject is an object die of any type: a DieTyped that is a
// metav1.Object, as an object die's base makes it. die.DieObjectOf has the
// rest, for generic code.
type DieObject interface {
	DieTyped
	metav1.Object
	metav1.ObjectMetaAccessor

	DieReleaseUnstructured() *unstructured.Unstructured
	DieReleaseUnstructuredE() (*unstructured.Unstructured, error)
}

// DieTypedBase is the base of a die for a type with type metadata.
type DieTypedBase[T any, D DieOf[T, D]] struct {
	DieBase[T, D]
}

var _ runtime.Object = (*DieTypedBase[struct{}, *typedCheck])(nil)

// AsTyped wraps the state of a die in a DieTypedBase.
func AsTyped[T any, D DieOf[T, D]](b DieBase[T, D]) DieTypedBase[T, D] {
	return DieTypedBase[T, D]{DieBase: b}
}

// GetObjectKind returns the die's apiVersion and kind.
func (b *DieTypedBase[T, D]) GetObjectKind() schema.ObjectKind {
	return &core.ObjectKind{H: b.hold(), At: b.at}
}

// DeepCopyObject returns a copy of the die.
func (b *DieTypedBase[T, D]) DeepCopyObject() runtime.Object {
	return any(b.DeepCopy()).(runtime.Object)
}

func metadata(field string) Path { return Path{core.FieldStep("metadata"), core.FieldStep(field)} }

// DieObjectBase is the base of a die for an object with object metadata.
type DieObjectBase[T any, D DieOf[T, D]] struct {
	DieTypedBase[T, D]
}

var _ runtime.Object = (*DieObjectBase[struct{}, *objectCheck])(nil)
var _ metav1.Object = (*DieObjectBase[struct{}, *objectCheck])(nil)
var _ metav1.ObjectMetaAccessor = (*DieObjectBase[struct{}, *objectCheck])(nil)

// AsObject wraps the state of a die in a DieObjectBase.
func AsObject[T any, D DieOf[T, D]](b DieBase[T, D]) DieObjectBase[T, D] {
	return DieObjectBase[T, D]{DieTypedBase: AsTyped(b)}
}

// DieFeedUnstructured feeds an unstructured object, as decoding it would:
// the values it has are present, zero or not, and fields T doesn't have are
// kept, to be encoded again. As with DieFeed, a root die discards its intent
// and a die within another sets the value as intent. It panics if the object
// doesn't decode as a T.
func (b *DieObjectBase[T, D]) DieFeedUnstructured(u *unstructured.Unstructured) D {
	return core.Must(b.DieFeedUnstructuredE(u))
}

// DieFeedUnstructuredE is DieFeedUnstructured, returning an error rather
// than panicking.
func (b *DieObjectBase[T, D]) DieFeedUnstructuredE(u *unstructured.Unstructured) (D, error) {
	data, err := u.MarshalJSON()
	if err != nil {
		var zero D
		return zero, core.Errorf("die: encoding an unstructured object: %w", err)
	}
	return b.DieFeedJSONE(data)
}

// DieReleaseUnstructured releases the object as an unstructured object,
// encoded as MarshalJSON encodes it: fields T doesn't have that were fed
// included, and absent values omitted.
func (b *DieObjectBase[T, D]) DieReleaseUnstructured() *unstructured.Unstructured {
	return core.Must(b.DieReleaseUnstructuredE())
}

// DieReleaseUnstructuredE is DieReleaseUnstructured, returning an error
// rather than panicking.
func (b *DieObjectBase[T, D]) DieReleaseUnstructuredE() (*unstructured.Unstructured, error) {
	data, err := b.DieReleaseJSONE()
	if err != nil {
		return nil, err
	}
	u := &unstructured.Unstructured{Object: map[string]any{}}
	// numbers decode as int64 where they can, as unstructured objects hold them
	if err := utiljson.Unmarshal(data, &u.Object); err != nil {
		return nil, core.Errorf("die: decoding %s as an unstructured object: %w", reflect.TypeFor[T](), err)
	}
	return u, nil
}

func (b *DieObjectBase[T, D]) GetObjectMeta() metav1.Object {
	return b
}

func (b *DieObjectBase[T, D]) GetNamespace() string {
	return core.Must(core.GetAt[string](b.hold(), b.at, metadata("namespace")))
}
func (b *DieObjectBase[T, D]) SetNamespace(namespace string) {
	core.Check(core.SetAt(b.hold(), b.at, "SetNamespace", metadata("namespace"), namespace))
}

func (b *DieObjectBase[T, D]) GetName() string {
	return core.Must(core.GetAt[string](b.hold(), b.at, metadata("name")))
}
func (b *DieObjectBase[T, D]) SetName(name string) {
	core.Check(core.SetAt(b.hold(), b.at, "SetName", metadata("name"), name))
}

func (b *DieObjectBase[T, D]) GetGenerateName() string {
	return core.Must(core.GetAt[string](b.hold(), b.at, metadata("generateName")))
}
func (b *DieObjectBase[T, D]) SetGenerateName(name string) {
	core.Check(core.SetAt(b.hold(), b.at, "SetGenerateName", metadata("generateName"), name))
}

func (b *DieObjectBase[T, D]) GetUID() types.UID {
	return core.Must(core.GetAt[types.UID](b.hold(), b.at, metadata("uid")))
}
func (b *DieObjectBase[T, D]) SetUID(uid types.UID) {
	core.Check(core.SetAt(b.hold(), b.at, "SetUID", metadata("uid"), uid))
}

func (b *DieObjectBase[T, D]) GetResourceVersion() string {
	return core.Must(core.GetAt[string](b.hold(), b.at, metadata("resourceVersion")))
}
func (b *DieObjectBase[T, D]) SetResourceVersion(version string) {
	core.Check(core.SetAt(b.hold(), b.at, "SetResourceVersion", metadata("resourceVersion"), version))
}

func (b *DieObjectBase[T, D]) GetGeneration() int64 {
	return core.Must(core.GetAt[int64](b.hold(), b.at, metadata("generation")))
}
func (b *DieObjectBase[T, D]) SetGeneration(generation int64) {
	core.Check(core.SetAt(b.hold(), b.at, "SetGeneration", metadata("generation"), generation))
}

func (b *DieObjectBase[T, D]) GetSelfLink() string {
	return core.Must(core.GetAt[string](b.hold(), b.at, metadata("selfLink")))
}
func (b *DieObjectBase[T, D]) SetSelfLink(selfLink string) {
	core.Check(core.SetAt(b.hold(), b.at, "SetSelfLink", metadata("selfLink"), selfLink))
}

func (b *DieObjectBase[T, D]) GetCreationTimestamp() metav1.Time {
	return core.Must(core.GetAt[metav1.Time](b.hold(), b.at, metadata("creationTimestamp")))
}
func (b *DieObjectBase[T, D]) SetCreationTimestamp(timestamp metav1.Time) {
	core.Check(core.SetAt(b.hold(), b.at, "SetCreationTimestamp", metadata("creationTimestamp"), timestamp))
}

func (b *DieObjectBase[T, D]) GetDeletionTimestamp() *metav1.Time {
	return core.Must(core.GetAt[*metav1.Time](b.hold(), b.at, metadata("deletionTimestamp")))
}
func (b *DieObjectBase[T, D]) SetDeletionTimestamp(timestamp *metav1.Time) {
	core.Check(core.SetAt(b.hold(), b.at, "SetDeletionTimestamp", metadata("deletionTimestamp"), timestamp))
}

func (b *DieObjectBase[T, D]) GetDeletionGracePeriodSeconds() *int64 {
	return core.Must(core.GetAt[*int64](b.hold(), b.at, metadata("deletionGracePeriodSeconds")))
}
func (b *DieObjectBase[T, D]) SetDeletionGracePeriodSeconds(seconds *int64) {
	core.Check(core.SetAt(b.hold(), b.at, "SetDeletionGracePeriodSeconds", metadata("deletionGracePeriodSeconds"), seconds))
}

func (b *DieObjectBase[T, D]) GetLabels() map[string]string {
	return core.Must(core.GetAt[map[string]string](b.hold(), b.at, metadata("labels")))
}
func (b *DieObjectBase[T, D]) SetLabels(labels map[string]string) {
	core.Check(core.SetAt(b.hold(), b.at, "SetLabels", metadata("labels"), labels))
}

func (b *DieObjectBase[T, D]) GetAnnotations() map[string]string {
	return core.Must(core.GetAt[map[string]string](b.hold(), b.at, metadata("annotations")))
}
func (b *DieObjectBase[T, D]) SetAnnotations(annotations map[string]string) {
	core.Check(core.SetAt(b.hold(), b.at, "SetAnnotations", metadata("annotations"), annotations))
}

func (b *DieObjectBase[T, D]) GetFinalizers() []string {
	return core.Must(core.GetAt[[]string](b.hold(), b.at, metadata("finalizers")))
}
func (b *DieObjectBase[T, D]) SetFinalizers(finalizers []string) {
	core.Check(core.SetAt(b.hold(), b.at, "SetFinalizers", metadata("finalizers"), finalizers))
}

func (b *DieObjectBase[T, D]) GetOwnerReferences() []metav1.OwnerReference {
	return core.Must(core.GetAt[[]metav1.OwnerReference](b.hold(), b.at, metadata("ownerReferences")))
}
func (b *DieObjectBase[T, D]) SetOwnerReferences(references []metav1.OwnerReference) {
	core.Check(core.SetAt(b.hold(), b.at, "SetOwnerReferences", metadata("ownerReferences"), references))
}

func (b *DieObjectBase[T, D]) GetManagedFields() []metav1.ManagedFieldsEntry {
	return core.Must(core.GetAt[[]metav1.ManagedFieldsEntry](b.hold(), b.at, metadata("managedFields")))
}
func (b *DieObjectBase[T, D]) SetManagedFields(managedFields []metav1.ManagedFieldsEntry) {
	core.Check(core.SetAt(b.hold(), b.at, "SetManagedFields", metadata("managedFields"), managedFields))
}

// bases implement their interfaces. Generated dies take the same shape.
type (
	typedCheck struct {
		DieTypedBase[struct{}, *typedCheck]
	}
	objectCheck struct {
		DieObjectBase[struct{}, *objectCheck]
	}
)

func (*typedCheck) DieFrom(b DieBase[struct{}, *typedCheck]) *typedCheck {
	return &typedCheck{AsTyped(b)}
}

func (*objectCheck) DieFrom(b DieBase[struct{}, *objectCheck]) *objectCheck {
	return &objectCheck{AsObject(b)}
}
