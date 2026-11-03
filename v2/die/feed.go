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

package die

import (
	"os"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"reconciler.io/dies/v2/die/base"
	"reconciler.io/dies/v2/die/internal/core"
	"reconciler.io/dies/v2/die/internal/refs"
)

// Feeding and releasing convert dies to and from the forms Kubernetes code
// uses. They are functions of a die, rather than its methods, so that they
// are compiled once rather than for each type of die. A die's own methods
// feed and release its Go value: DieFeed and DieRelease.
//
// Feeding is a write: a feed takes a mutable die, and returns only an error.
// Releasing is a read of any die, and fails only if the die can't be encoded.

func refOf(d Die) (refs.Ref, error) {
	r, ok := refs.Of(d)
	if !ok {
		return refs.Ref{}, core.Errorf("die: %T is not a die", d)
	}
	return r, nil
}

// FeedJSON feeds a die a value decoded from JSON. On a root die, as decoding
// does, the values the JSON has are present, zero or not, and fields the
// die's type doesn't have are kept, to be encoded again. Otherwise it feeds
// the value as DieFeed does: a root die discards its intent and a die within
// another sets the value as intent. It fails if the JSON doesn't decode as
// the die's type.
func FeedJSON(d MutableDie, data []byte) error {
	r, err := refOf(d)
	if err != nil {
		return err
	}
	return r.FeedJSON(data)
}

// FeedYAML feeds a die a value decoded from YAML, as [FeedJSON] does.
func FeedYAML(d MutableDie, data []byte) error {
	r, err := refOf(d)
	if err != nil {
		return err
	}
	return r.FeedYAML(data)
}

// FeedYAMLFile returns an immutable die of the YAML, or JSON, in a file, fed
// to a copy of d as [FeedYAML] feeds it. It returns an error if the file
// can't be read or doesn't decode. [MustFeedYAMLFile] panics instead, so that
// fixtures can be package variables:
//
//	var web = die.MustFeedYAMLFile(DeploymentBlank, "testdata/web.yaml")
func FeedYAMLFile[T any, D base.DieOf[T, D, M], M base.DieMutOf[T, D, M]](d D, name string) (D, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		var zero D
		return zero, core.Errorf("die: %w", err)
	}
	return d.DieStampE(func(m M) error { return FeedYAML(m, data) })
}

// MustFeedYAMLFile is [FeedYAMLFile], panicking rather than returning an
// error.
func MustFeedYAMLFile[T any, D base.DieOf[T, D, M], M base.DieMutOf[T, D, M]](d D, name string) D {
	return core.Must(FeedYAMLFile(d, name))
}

// FeedRawExtension feeds a die the value of a raw extension, its object or
// its raw JSON, as [FeedJSON] does.
func FeedRawExtension(d MutableDie, raw runtime.RawExtension) error {
	r, err := refOf(d)
	if err != nil {
		return err
	}
	return r.FeedRawExtension(raw)
}

// FeedUnstructured feeds a mutable object die an unstructured object, as
// decoding it would, as [FeedJSON] does. An unstructured die copies its
// content instead, so its numbers stay int64 where they are. Feed a die that
// isn't an object a raw extension instead, with [FeedRawExtension].
func FeedUnstructured(d interface {
	MutableDie
	DieObject
}, u *unstructured.Unstructured) error {
	r, err := refOf(d)
	if err != nil {
		return err
	}
	return r.FeedUnstructured(u)
}

// FeedApplyPatch feeds a root die a server-side apply request as intent,
// over the die's value, which is kept: each field it has is set, null clears
// a field, [] and {} are an empty list or map, list-map elements are found by
// their keys, and an atomic struct is set whole. The intent replaces the
// die's, so [ReleaseApplyPatch] sends the request back. Fields the die's type
// doesn't have aren't fed. An unstructured die follows its mold, if it has
// one; without, its lists are set whole. It fails if the request isn't a
// JSON object.
func FeedApplyPatch(d MutableDie, patch []byte) error {
	r, err := refOf(d)
	if err != nil {
		return err
	}
	return r.FeedPatch("FeedApplyPatch", core.ApplyPatch, patch)
}

// FeedMergePatch feeds a root die a JSON merge patch, RFC 7386, as intent
// over the die's value, which is kept: null clears a field, objects and maps
// are merged field by field, and anything else, lists included, is set whole.
// The intent replaces the die's, so [ReleaseMergePatch] sends the patch back,
// but for what the patch sets that the value already has. Fields the die's
// type doesn't have aren't fed. It fails if the patch isn't a JSON object.
func FeedMergePatch(d MutableDie, patch []byte) error {
	r, err := refOf(d)
	if err != nil {
		return err
	}
	return r.FeedPatch("FeedMergePatch", core.MergePatch, patch)
}

// ReleaseJSON returns a die's current value as JSON, as its MarshalJSON
// encodes it. It returns an error if the value doesn't encode.
func ReleaseJSON(d Die) ([]byte, error) {
	r, err := refOf(d)
	if err != nil {
		return nil, err
	}
	return r.ReleaseJSON()
}

// MustReleaseJSON is [ReleaseJSON], panicking rather than returning an error.
func MustReleaseJSON(d Die) []byte { return core.Must(ReleaseJSON(d)) }

// ReleaseYAML returns a die's current value as YAML, encoded as
// [ReleaseJSON] encodes it.
func ReleaseYAML(d Die) ([]byte, error) {
	r, err := refOf(d)
	if err != nil {
		return nil, err
	}
	return r.ReleaseYAML()
}

// MustReleaseYAML is [ReleaseYAML], panicking rather than returning an error.
func MustReleaseYAML(d Die) []byte { return core.Must(ReleaseYAML(d)) }

// ReleaseRawExtension returns a die's current value as a raw extension holding
// its JSON, encoded as [ReleaseJSON] encodes it.
func ReleaseRawExtension(d Die) (runtime.RawExtension, error) {
	r, err := refOf(d)
	if err != nil {
		return runtime.RawExtension{}, err
	}
	return r.ReleaseRawExtension()
}

// MustReleaseRawExtension is [ReleaseRawExtension], panicking rather than returning an
// error.
func MustReleaseRawExtension(d Die) runtime.RawExtension { return core.Must(ReleaseRawExtension(d)) }

// ReleaseUnstructured returns an object die's current value as an unstructured
// object, encoded as [ReleaseJSON] encodes it: fields the die's type
// doesn't have that were fed included, and absent values omitted. An
// unstructured die's is a copy of its content. Release a die that isn't an
// object as a raw extension instead, with [ReleaseRawExtension].
func ReleaseUnstructured(d DieObject) (*unstructured.Unstructured, error) {
	r, err := refOf(d)
	if err != nil {
		return nil, err
	}
	return r.ReleaseUnstructured()
}

// MustReleaseUnstructured is [ReleaseUnstructured], panicking rather than returning an
// error.
func MustReleaseUnstructured(d DieObject) *unstructured.Unstructured {
	return core.Must(ReleaseUnstructured(d))
}

// ReleaseApplyPatch encodes a root die's intent as a server-side apply
// request, or the part of it scopes keep, such as [Subresource]("status").
func ReleaseApplyPatch(d Die, scopes ...Scope) ([]byte, error) {
	r, err := scoped(d, "ReleaseApplyPatch", scopes)
	if err != nil {
		return nil, err
	}
	return r.ReleaseApplyPatch()
}

// MustReleaseApplyPatch is [ReleaseApplyPatch], panicking rather than returning an
// error.
func MustReleaseApplyPatch(d Die, scopes ...Scope) []byte {
	return core.Must(ReleaseApplyPatch(d, scopes...))
}

// ReleaseMergePatch encodes a root die's intent as a JSON merge patch against
// its base, or the part of it scopes keep, such as [Only](d.Spec()).
func ReleaseMergePatch(d Die, scopes ...Scope) ([]byte, error) {
	r, err := scoped(d, "ReleaseMergePatch", scopes)
	if err != nil {
		return nil, err
	}
	return r.ReleaseMergePatch()
}

// MustReleaseMergePatch is [ReleaseMergePatch], panicking rather than returning an
// error.
func MustReleaseMergePatch(d Die, scopes ...Scope) []byte {
	return core.Must(ReleaseMergePatch(d, scopes...))
}
