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

// Dies for top-level objects embed one of these instead of DieBase or
// DieMutBase, so that they can be used where Kubernetes libraries expect
// objects:
//
//   - DieTypedBase and DieTypedMutBase implement runtime.Object
//   - DieObjectBase and DieObjectMutBase also implement metav1.Object and
//     metav1.ObjectMetaAccessor, through the object's metadata
//
// Lists are not dies of their own, but lists of dies: see List.
//
// Their setters return nothing, so they change a die in place: a mutable
// die's do, and an immutable die's panic with ErrImmutable rather than lose
// the change.

// DieTyped is a die of any type with type metadata, immutable or mutable: a
// Die that is a runtime.Object, as a typed die's base makes it.
// die.DieTypedOf and die.DieTypedMutOf have the rest, for generic code.
type DieTyped interface {
	Die
	runtime.Object
}

// DieObject is an object die of any type, immutable or mutable: a DieTyped
// that is a metav1.Object, as an object die's base makes it.
// die.DieObjectOf and die.DieObjectMutOf have the rest, for generic code.
type DieObject interface {
	DieTyped
	metav1.Object
	metav1.ObjectMetaAccessor

	DieReleaseUnstructured() *unstructured.Unstructured
	DieReleaseUnstructuredE() (*unstructured.Unstructured, error)
}

func immutable(method string) error { return core.Errorf("die: %s: %w", method, core.ErrImmutable) }

func metadata(field string) Path { return Path{core.FieldStep("metadata"), core.FieldStep(field)} }

// DieTypedBase is the base of an immutable die for a type with type
// metadata.
type DieTypedBase[T any, D DieOf[T, D, M], M DieMutOf[T, D, M]] struct {
	DieBase[T, D, M]
}

// DieTypedMutBase is the base of a mutable die for a type with type
// metadata.
type DieTypedMutBase[T any, D DieOf[T, D, M], M DieMutOf[T, D, M]] struct {
	DieMutBase[T, D, M]
}

var (
	_ runtime.Object = (*DieTypedBase[struct{}, *typedCheck, *typedMutCheck])(nil)
	_ runtime.Object = (*DieTypedMutBase[struct{}, *typedCheck, *typedMutCheck])(nil)
)

// AsTyped wraps the state of a die in a DieTypedBase.
func AsTyped[T any, D DieOf[T, D, M], M DieMutOf[T, D, M]](b DieBase[T, D, M]) DieTypedBase[T, D, M] {
	return DieTypedBase[T, D, M]{DieBase: b}
}

// AsTypedMut wraps the state of a die in a DieTypedMutBase.
func AsTypedMut[T any, D DieOf[T, D, M], M DieMutOf[T, D, M]](b DieMutBase[T, D, M]) DieTypedMutBase[T, D, M] {
	return DieTypedMutBase[T, D, M]{DieMutBase: b}
}

// GetObjectKind returns the die's apiVersion and kind, which can't be set.
func (b *DieTypedBase[T, D, M]) GetObjectKind() schema.ObjectKind {
	return &core.ObjectKind{H: b.hold(), At: b.at, ReadOnly: true}
}

// DeepCopyObject returns a copy of the die.
func (b *DieTypedBase[T, D, M]) DeepCopyObject() runtime.Object {
	return any(b.DeepCopy()).(runtime.Object)
}

// GetObjectKind returns the die's apiVersion and kind.
func (b *DieTypedMutBase[T, D, M]) GetObjectKind() schema.ObjectKind {
	return &core.ObjectKind{H: b.hold(), At: b.at}
}

// DeepCopyObject returns a copy of the die.
func (b *DieTypedMutBase[T, D, M]) DeepCopyObject() runtime.Object {
	return any(b.DeepCopy()).(runtime.Object)
}

// DieObjectBase is the base of an immutable die for an object with object
// metadata.
type DieObjectBase[T any, D DieOf[T, D, M], M DieMutOf[T, D, M]] struct {
	DieTypedBase[T, D, M]
}

// DieObjectMutBase is the base of a mutable die for an object with object
// metadata.
type DieObjectMutBase[T any, D DieOf[T, D, M], M DieMutOf[T, D, M]] struct {
	DieTypedMutBase[T, D, M]
}

var (
	_ DieObject = (*DieObjectBase[struct{}, *objectCheck, *objectMutCheck])(nil)
	_ DieObject = (*DieObjectMutBase[struct{}, *objectCheck, *objectMutCheck])(nil)
)

// AsObject wraps the state of a die in a DieObjectBase.
func AsObject[T any, D DieOf[T, D, M], M DieMutOf[T, D, M]](b DieBase[T, D, M]) DieObjectBase[T, D, M] {
	return DieObjectBase[T, D, M]{DieTypedBase: AsTyped(b)}
}

// AsObjectMut wraps the state of a die in a DieObjectMutBase.
func AsObjectMut[T any, D DieOf[T, D, M], M DieMutOf[T, D, M]](b DieMutBase[T, D, M]) DieObjectMutBase[T, D, M] {
	return DieObjectMutBase[T, D, M]{DieTypedMutBase: AsTypedMut(b)}
}

// DieFeedUnstructured feeds an unstructured object, as decoding it would:
// the values it has are present, zero or not, and fields T doesn't have are
// kept, to be encoded again. As with DieFeed, a root die discards its intent
// and a die within another sets the value as intent. It panics if the object
// doesn't decode as a T.
func (b *DieObjectMutBase[T, D, M]) DieFeedUnstructured(u *unstructured.Unstructured) M {
	return core.Must(b.DieFeedUnstructuredE(u))
}

// DieFeedUnstructuredE is DieFeedUnstructured, returning an error rather
// than panicking.
func (b *DieObjectMutBase[T, D, M]) DieFeedUnstructuredE(u *unstructured.Unstructured) (M, error) {
	data, err := u.MarshalJSON()
	if err != nil {
		var zero M
		return zero, core.Errorf("die: encoding an unstructured object: %w", err)
	}
	return b.DieFeedJSONE(data)
}

// DieReleaseUnstructured releases the object as an unstructured object,
// encoded as MarshalJSON encodes it: fields T doesn't have that were fed
// included, and absent values omitted.
func (b *DieObjectBase[T, D, M]) DieReleaseUnstructured() *unstructured.Unstructured {
	return core.Must(b.DieReleaseUnstructuredE())
}

// DieReleaseUnstructuredE is DieReleaseUnstructured, returning an error
// rather than panicking.
func (b *DieObjectBase[T, D, M]) DieReleaseUnstructuredE() (*unstructured.Unstructured, error) {
	return releaseUnstructured[T](&b.state)
}

// DieReleaseUnstructured releases the object as an unstructured object: see
// DieObjectBase.DieReleaseUnstructured.
func (b *DieObjectMutBase[T, D, M]) DieReleaseUnstructured() *unstructured.Unstructured {
	return core.Must(b.DieReleaseUnstructuredE())
}

// DieReleaseUnstructuredE is DieReleaseUnstructured, returning an error
// rather than panicking.
func (b *DieObjectMutBase[T, D, M]) DieReleaseUnstructuredE() (*unstructured.Unstructured, error) {
	return releaseUnstructured[T](&b.state)
}

func releaseUnstructured[T any](d interface{ DieReleaseJSONE() ([]byte, error) }) (*unstructured.Unstructured, error) {
	data, err := d.DieReleaseJSONE()
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

// The metav1.Object methods of immutable object dies, whose setters panic.

func (b *DieObjectBase[T, D, M]) GetObjectMeta() metav1.Object { return b }

func (b *DieObjectBase[T, D, M]) GetNamespace() string {
	return core.Must(core.GetAt[string](b.hold(), b.at, metadata("namespace")))
}
func (b *DieObjectBase[T, D, M]) SetNamespace(string) { panic(immutable("SetNamespace")) }

func (b *DieObjectBase[T, D, M]) GetName() string {
	return core.Must(core.GetAt[string](b.hold(), b.at, metadata("name")))
}
func (b *DieObjectBase[T, D, M]) SetName(string) { panic(immutable("SetName")) }

func (b *DieObjectBase[T, D, M]) GetGenerateName() string {
	return core.Must(core.GetAt[string](b.hold(), b.at, metadata("generateName")))
}
func (b *DieObjectBase[T, D, M]) SetGenerateName(string) { panic(immutable("SetGenerateName")) }

func (b *DieObjectBase[T, D, M]) GetUID() types.UID {
	return core.Must(core.GetAt[types.UID](b.hold(), b.at, metadata("uid")))
}
func (b *DieObjectBase[T, D, M]) SetUID(types.UID) { panic(immutable("SetUID")) }

func (b *DieObjectBase[T, D, M]) GetResourceVersion() string {
	return core.Must(core.GetAt[string](b.hold(), b.at, metadata("resourceVersion")))
}
func (b *DieObjectBase[T, D, M]) SetResourceVersion(string) { panic(immutable("SetResourceVersion")) }

func (b *DieObjectBase[T, D, M]) GetGeneration() int64 {
	return core.Must(core.GetAt[int64](b.hold(), b.at, metadata("generation")))
}
func (b *DieObjectBase[T, D, M]) SetGeneration(int64) { panic(immutable("SetGeneration")) }

func (b *DieObjectBase[T, D, M]) GetSelfLink() string {
	return core.Must(core.GetAt[string](b.hold(), b.at, metadata("selfLink")))
}
func (b *DieObjectBase[T, D, M]) SetSelfLink(string) { panic(immutable("SetSelfLink")) }

func (b *DieObjectBase[T, D, M]) GetCreationTimestamp() metav1.Time {
	return core.Must(core.GetAt[metav1.Time](b.hold(), b.at, metadata("creationTimestamp")))
}
func (b *DieObjectBase[T, D, M]) SetCreationTimestamp(metav1.Time) {
	panic(immutable("SetCreationTimestamp"))
}

func (b *DieObjectBase[T, D, M]) GetDeletionTimestamp() *metav1.Time {
	return core.Must(core.GetAt[*metav1.Time](b.hold(), b.at, metadata("deletionTimestamp")))
}
func (b *DieObjectBase[T, D, M]) SetDeletionTimestamp(*metav1.Time) {
	panic(immutable("SetDeletionTimestamp"))
}

func (b *DieObjectBase[T, D, M]) GetDeletionGracePeriodSeconds() *int64 {
	return core.Must(core.GetAt[*int64](b.hold(), b.at, metadata("deletionGracePeriodSeconds")))
}
func (b *DieObjectBase[T, D, M]) SetDeletionGracePeriodSeconds(*int64) {
	panic(immutable("SetDeletionGracePeriodSeconds"))
}

func (b *DieObjectBase[T, D, M]) GetLabels() map[string]string {
	return core.Must(core.GetAt[map[string]string](b.hold(), b.at, metadata("labels")))
}
func (b *DieObjectBase[T, D, M]) SetLabels(map[string]string) { panic(immutable("SetLabels")) }

func (b *DieObjectBase[T, D, M]) GetAnnotations() map[string]string {
	return core.Must(core.GetAt[map[string]string](b.hold(), b.at, metadata("annotations")))
}
func (b *DieObjectBase[T, D, M]) SetAnnotations(map[string]string) {
	panic(immutable("SetAnnotations"))
}

func (b *DieObjectBase[T, D, M]) GetFinalizers() []string {
	return core.Must(core.GetAt[[]string](b.hold(), b.at, metadata("finalizers")))
}
func (b *DieObjectBase[T, D, M]) SetFinalizers([]string) { panic(immutable("SetFinalizers")) }

func (b *DieObjectBase[T, D, M]) GetOwnerReferences() []metav1.OwnerReference {
	return core.Must(core.GetAt[[]metav1.OwnerReference](b.hold(), b.at, metadata("ownerReferences")))
}
func (b *DieObjectBase[T, D, M]) SetOwnerReferences([]metav1.OwnerReference) {
	panic(immutable("SetOwnerReferences"))
}

func (b *DieObjectBase[T, D, M]) GetManagedFields() []metav1.ManagedFieldsEntry {
	return core.Must(core.GetAt[[]metav1.ManagedFieldsEntry](b.hold(), b.at, metadata("managedFields")))
}
func (b *DieObjectBase[T, D, M]) SetManagedFields([]metav1.ManagedFieldsEntry) {
	panic(immutable("SetManagedFields"))
}

// The metav1.Object methods of mutable object dies.

func (b *DieObjectMutBase[T, D, M]) GetObjectMeta() metav1.Object { return b }

func (b *DieObjectMutBase[T, D, M]) GetNamespace() string {
	return core.Must(core.GetAt[string](b.hold(), b.at, metadata("namespace")))
}
func (b *DieObjectMutBase[T, D, M]) SetNamespace(namespace string) {
	core.Check(core.SetAt(b.hold(), b.at, "SetNamespace", metadata("namespace"), namespace))
}

func (b *DieObjectMutBase[T, D, M]) GetName() string {
	return core.Must(core.GetAt[string](b.hold(), b.at, metadata("name")))
}
func (b *DieObjectMutBase[T, D, M]) SetName(name string) {
	core.Check(core.SetAt(b.hold(), b.at, "SetName", metadata("name"), name))
}

func (b *DieObjectMutBase[T, D, M]) GetGenerateName() string {
	return core.Must(core.GetAt[string](b.hold(), b.at, metadata("generateName")))
}
func (b *DieObjectMutBase[T, D, M]) SetGenerateName(name string) {
	core.Check(core.SetAt(b.hold(), b.at, "SetGenerateName", metadata("generateName"), name))
}

func (b *DieObjectMutBase[T, D, M]) GetUID() types.UID {
	return core.Must(core.GetAt[types.UID](b.hold(), b.at, metadata("uid")))
}
func (b *DieObjectMutBase[T, D, M]) SetUID(uid types.UID) {
	core.Check(core.SetAt(b.hold(), b.at, "SetUID", metadata("uid"), uid))
}

func (b *DieObjectMutBase[T, D, M]) GetResourceVersion() string {
	return core.Must(core.GetAt[string](b.hold(), b.at, metadata("resourceVersion")))
}
func (b *DieObjectMutBase[T, D, M]) SetResourceVersion(version string) {
	core.Check(core.SetAt(b.hold(), b.at, "SetResourceVersion", metadata("resourceVersion"), version))
}

func (b *DieObjectMutBase[T, D, M]) GetGeneration() int64 {
	return core.Must(core.GetAt[int64](b.hold(), b.at, metadata("generation")))
}
func (b *DieObjectMutBase[T, D, M]) SetGeneration(generation int64) {
	core.Check(core.SetAt(b.hold(), b.at, "SetGeneration", metadata("generation"), generation))
}

func (b *DieObjectMutBase[T, D, M]) GetSelfLink() string {
	return core.Must(core.GetAt[string](b.hold(), b.at, metadata("selfLink")))
}
func (b *DieObjectMutBase[T, D, M]) SetSelfLink(selfLink string) {
	core.Check(core.SetAt(b.hold(), b.at, "SetSelfLink", metadata("selfLink"), selfLink))
}

func (b *DieObjectMutBase[T, D, M]) GetCreationTimestamp() metav1.Time {
	return core.Must(core.GetAt[metav1.Time](b.hold(), b.at, metadata("creationTimestamp")))
}
func (b *DieObjectMutBase[T, D, M]) SetCreationTimestamp(timestamp metav1.Time) {
	core.Check(core.SetAt(b.hold(), b.at, "SetCreationTimestamp", metadata("creationTimestamp"), timestamp))
}

func (b *DieObjectMutBase[T, D, M]) GetDeletionTimestamp() *metav1.Time {
	return core.Must(core.GetAt[*metav1.Time](b.hold(), b.at, metadata("deletionTimestamp")))
}
func (b *DieObjectMutBase[T, D, M]) SetDeletionTimestamp(timestamp *metav1.Time) {
	core.Check(core.SetAt(b.hold(), b.at, "SetDeletionTimestamp", metadata("deletionTimestamp"), timestamp))
}

func (b *DieObjectMutBase[T, D, M]) GetDeletionGracePeriodSeconds() *int64 {
	return core.Must(core.GetAt[*int64](b.hold(), b.at, metadata("deletionGracePeriodSeconds")))
}
func (b *DieObjectMutBase[T, D, M]) SetDeletionGracePeriodSeconds(seconds *int64) {
	core.Check(core.SetAt(b.hold(), b.at, "SetDeletionGracePeriodSeconds", metadata("deletionGracePeriodSeconds"), seconds))
}

func (b *DieObjectMutBase[T, D, M]) GetLabels() map[string]string {
	return core.Must(core.GetAt[map[string]string](b.hold(), b.at, metadata("labels")))
}
func (b *DieObjectMutBase[T, D, M]) SetLabels(labels map[string]string) {
	core.Check(core.SetAt(b.hold(), b.at, "SetLabels", metadata("labels"), labels))
}

func (b *DieObjectMutBase[T, D, M]) GetAnnotations() map[string]string {
	return core.Must(core.GetAt[map[string]string](b.hold(), b.at, metadata("annotations")))
}
func (b *DieObjectMutBase[T, D, M]) SetAnnotations(annotations map[string]string) {
	core.Check(core.SetAt(b.hold(), b.at, "SetAnnotations", metadata("annotations"), annotations))
}

func (b *DieObjectMutBase[T, D, M]) GetFinalizers() []string {
	return core.Must(core.GetAt[[]string](b.hold(), b.at, metadata("finalizers")))
}
func (b *DieObjectMutBase[T, D, M]) SetFinalizers(finalizers []string) {
	core.Check(core.SetAt(b.hold(), b.at, "SetFinalizers", metadata("finalizers"), finalizers))
}

func (b *DieObjectMutBase[T, D, M]) GetOwnerReferences() []metav1.OwnerReference {
	return core.Must(core.GetAt[[]metav1.OwnerReference](b.hold(), b.at, metadata("ownerReferences")))
}
func (b *DieObjectMutBase[T, D, M]) SetOwnerReferences(references []metav1.OwnerReference) {
	core.Check(core.SetAt(b.hold(), b.at, "SetOwnerReferences", metadata("ownerReferences"), references))
}

func (b *DieObjectMutBase[T, D, M]) GetManagedFields() []metav1.ManagedFieldsEntry {
	return core.Must(core.GetAt[[]metav1.ManagedFieldsEntry](b.hold(), b.at, metadata("managedFields")))
}
func (b *DieObjectMutBase[T, D, M]) SetManagedFields(managedFields []metav1.ManagedFieldsEntry) {
	core.Check(core.SetAt(b.hold(), b.at, "SetManagedFields", metadata("managedFields"), managedFields))
}

// bases implement their interfaces. Generated dies take the same shape.
type (
	typedCheck struct {
		DieTypedBase[struct{}, *typedCheck, *typedMutCheck]
	}
	typedMutCheck struct {
		DieTypedMutBase[struct{}, *typedCheck, *typedMutCheck]
	}
	objectCheck struct {
		DieObjectBase[struct{}, *objectCheck, *objectMutCheck]
	}
	objectMutCheck struct {
		DieObjectMutBase[struct{}, *objectCheck, *objectMutCheck]
	}
)

func (*typedCheck) DieFrom(b DieBase[struct{}, *typedCheck, *typedMutCheck]) *typedCheck {
	return &typedCheck{AsTyped(b)}
}

func (*typedMutCheck) DieFrom(b DieMutBase[struct{}, *typedCheck, *typedMutCheck]) *typedMutCheck {
	return &typedMutCheck{AsTypedMut(b)}
}

func (*objectCheck) DieFrom(b DieBase[struct{}, *objectCheck, *objectMutCheck]) *objectCheck {
	return &objectCheck{AsObject(b)}
}

func (*objectMutCheck) DieFrom(b DieMutBase[struct{}, *objectCheck, *objectMutCheck]) *objectMutCheck {
	return &objectMutCheck{AsObjectMut(b)}
}
