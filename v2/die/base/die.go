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

// Package base has the bases dies embed, which give every die its methods,
// and blanks.
//
// A die is an immutable base value plus a persistent overlay of intent. Writes
// record intent (set, clear, empty) by wire path and copy only the overlay
// nodes along that path; the base is never copied. Reads apply the overlay
// along the path being read. Patches are encoded from the overlay, so they
// contain exactly what was stamped.
//
// A child die returned by an accessor is attached to its root: it addresses
// into the root by path, reads never create anything, and the first write
// creates whatever the path needs.
package base

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	utiljson "k8s.io/apimachinery/pkg/util/json"
	"sigs.k8s.io/yaml"

	"reconciler.io/dies/v2/die/internal/core"
)

// DieOf carries a die's types for generic code: T is the die's value type,
// and D its own pointer type, so its methods return the concrete die. It has
// every method of a die. Generic code constrains a die by it,
// D DieOf[T, D]; code that takes a die of any type takes a Die.
type DieOf[T any, D DieOf[T, D]] interface {
	Die
	LocatedAs[T]
	DieEmbedded() *DieBase[T, D]
	// DieFrom wraps state in a new die. It is called on a nil receiver.
	DieFrom(DieBase[T, D]) D

	DieClear() D
	DieClearE() (D, error)
	DieRelease() T
	DieReleaseE() (T, error)
	DieReleasePtr() *T
	DieReleasePtrE() (*T, error)
	DieFeed(r T) D
	DieFeedE(r T) (D, error)
	DieFeedPtr(r *T) D
	DieFeedPtrE(r *T) (D, error)
	DieFeedJSON(data []byte) D
	DieFeedJSONE(data []byte) (D, error)
	DieFeedYAML(data []byte) D
	DieFeedYAMLE(data []byte) (D, error)
	DieFeedRawExtension(raw runtime.RawExtension) D
	DieFeedRawExtensionE(raw runtime.RawExtension) (D, error)
	DieFeedApplyPatch(patch []byte) D
	DieFeedApplyPatchE(patch []byte) (D, error)
	DieFeedMergePatch(patch []byte) D
	DieFeedMergePatchE(patch []byte) (D, error)
	DieStamp(fn func(r *T)) D
	DieStampE(fn func(r *T)) (D, error)
	DieWith(fns ...func(d D)) D
	DieWithE(fns ...func(d D) error) (D, error)
	DieImmutable(immutable bool) D
	DieImmutableE(immutable bool) (D, error)
	DeepCopy() D
	DeepCopyE() (D, error)
	DieCommit() D
	DieCommitE() (D, error)
	DieOnly(at ...Located) D
	DieOnlyE(at ...Located) (D, error)
	DieExcept(at ...Located) D
	DieExceptE(at ...Located) (D, error)
	DieSubresource(name string) D
	DieSubresourceE(name string) (D, error)
	DieMain() D
	DieMainE() (D, error)
	Equal(o D) bool
	EqualE(o D) (bool, error)
}

// Die is a die of any type: the methods every die has that don't name its
// types, for functions that take any die. Handles, collections and nodes
// have some of them, but aren't dies. DieOf has the rest, for generic code.
type Die interface {
	// MarshalJSON encodes the die's value as its presence says.
	json.Marshaler

	DiePath() Path
	DieRootType() reflect.Type
	DieRoot() Die
	DieMutable() bool
	DieIntent() string
	DieFields() []FieldInfo
	DieHas() bool
	DieHasE() (bool, error)
	DieReleaseJSON() []byte
	DieReleaseJSONE() ([]byte, error)
	DieReleaseYAML() []byte
	DieReleaseYAMLE() ([]byte, error)
	DieReleaseRawExtension() runtime.RawExtension
	DieReleaseRawExtensionE() (runtime.RawExtension, error)
	DieReleaseApplyPatch() []byte
	DieReleaseApplyPatchE() ([]byte, error)
	DieReleaseMergePatch() []byte
	DieReleaseMergePatchE() ([]byte, error)
}

// DieBase holds the state of a die and provides the methods shared by every die.
// Dies embed it by value.
//
// The zero value is a mutable, empty die, so a die allocated by a decoder can
// be unmarshaled into.
type DieBase[T any, D DieOf[T, D]] struct {
	h  *core.Holder
	at Path
}

func wrap[T any, D DieOf[T, D]](h *core.Holder, at Path) D {
	if len(at) == 0 {
		h.SetRoot(func(h *core.Holder) any { return wrap[T, D](h, nil) })
	}
	var d D
	return d.DieFrom(DieBase[T, D]{h: h, at: at})
}

// BlankOf returns an immutable die of a value, such as an object with only
// its type metadata, for a blank.
func BlankOf[T any, D DieOf[T, D]](v T) D {
	return wrap[T, D](&core.Holder{Base: reflect.ValueOf(v)}, nil)
}

// Blank returns an immutable, empty die.
func Blank[T any, D DieOf[T, D]]() D {
	return wrap[T, D](&core.Holder{Base: reflect.Zero(reflect.TypeFor[T]())}, nil)
}

func (b *DieBase[T, D]) DieEmbedded() *DieBase[T, D] { return b }

// DieValue marks a die as a location of a T, for LocatedAs. It does nothing.
func (*DieBase[T, D]) DieValue(*T) {}

// dieState is how the packages of handles, jigs and unstructured dies reach
// a die's state, through the functions installed below.
type dieState interface {
	hold() *core.Holder
	dieAt() Path
	setDieState(h *core.Holder, at Path)
	dieValueType() reflect.Type
}

func (b *DieBase[T, D]) dieAt() Path { return b.at }

func (b *DieBase[T, D]) setDieState(h *core.Holder, at Path) { b.h, b.at = h, at }

func (b *DieBase[T, D]) dieValueType() reflect.Type { return reflect.TypeFor[T]() }

func init() {
	core.StateOf = func(d any) (*core.Holder, Path, bool) {
		s, ok := d.(dieState)
		if !ok {
			return nil, nil, false
		}
		return s.hold(), s.dieAt(), true
	}
	core.SetState = func(d any, h *core.Holder, at Path) {
		if len(at) == 0 {
			// the die of the root, of whatever type it is
			t := reflect.TypeOf(d).Elem()
			h.SetRoot(func(h *core.Holder) any {
				root := reflect.New(t).Interface()
				root.(dieState).setDieState(h, nil)
				return root
			})
		}
		d.(dieState).setDieState(h, at)
	}
	core.ValueTypeOf = func(d any) reflect.Type { return d.(dieState).dieValueType() }
}

func (b *DieBase[T, D]) hold() *core.Holder {
	if b.h == nil {
		var zero T
		b.h = &core.Holder{Base: reflect.ValueOf(zero), Mutable: true}
	}
	return b.h
}

func (b *DieBase[T, D]) self(h *core.Holder) D {
	if h != b.h {
		// a new root
		return wrap[T, D](h, nil)
	}
	return wrap[T, D](h, b.at)
}

func (b *DieBase[T, D]) root(method string) (*core.Holder, error) {
	if len(b.at) > 0 {
		return nil, core.Errorf("die: %s is only supported on a root die, not at %s", method, b.at)
	}
	if b.hold().View != nil {
		return nil, core.Errorf("die: %s is not supported on a jig; use it on the jig's source", method)
	}
	return b.hold(), nil
}

// DiePath returns the location of this die within its root.
func (b *DieBase[T, D]) DiePath() Path { return b.at }

// DieRoot returns the root die this die is within, sharing its state: writes
// through it are writes to this die's tree. A root returns itself. Its type
// is the root's die type, such as *DeploymentDie for a die within a
// Deployment. Assert it as a type, or as metav1.Object for an object's
// identity whatever its type: d.DieRoot().(metav1.Object).GetNamespace().
func (b *DieBase[T, D]) DieRoot() Die {
	if len(b.at) == 0 {
		return any(b.self(b.hold())).(Die)
	}
	return RootOf(b.hold())
}

// RootOf returns the root die over a holder, or nil if none was made.
func RootOf(h *core.Holder) Die {
	root, _ := h.Root().(Die)
	return root
}

// DieRootType returns the type of the root value this die is within.
func (b *DieBase[T, D]) DieRootType() reflect.Type { return b.hold().Base.Type() }

// DieMutable reports whether writes to the die change it in place.
func (b *DieBase[T, D]) DieMutable() bool { return b.hold().Mutable }

// DieHas reports whether the value is present. A root is always present.
func (b *DieBase[T, D]) DieHas() bool { return core.Must(b.DieHasE()) }

// DieHasE is DieHas, returning an error rather than panicking.
func (b *DieBase[T, D]) DieHasE() (bool, error) {
	if len(b.at) == 0 {
		return true, nil
	}
	return b.hold().Has(b.at)
}

// DieClear records that the value is absent: null in patches, or a removed
// element. On a root it resets to the zero value.
func (b *DieBase[T, D]) DieClear() D { return core.Must(b.DieClearE()) }

// DieClearE is DieClear, returning an error rather than panicking.
func (b *DieBase[T, D]) DieClearE() (D, error) {
	if len(b.at) == 0 {
		var zero T
		return b.DieFeedE(zero)
	}
	return b.write(b.hold().Write(b.at, nil, core.SetOp(core.OpClear, reflect.Value{}), false))
}

// write returns the die of the holder a write returned, or its error.
func (b *DieBase[T, D]) write(h *core.Holder, err error) (D, error) {
	if err != nil {
		var zero D
		return zero, err
	}
	return b.self(h), nil
}

// DieRelease returns a copy of the current value.
func (b *DieBase[T, D]) DieRelease() T { return core.Must(b.DieReleaseE()) }

// DieReleaseE is DieRelease, returning an error rather than panicking.
func (b *DieBase[T, D]) DieReleaseE() (T, error) {
	v, err := b.hold().ValueAt(b.at)
	if err != nil {
		var zero T
		return zero, err
	}
	return as[T](v), nil
}

// DieReleasePtr returns a pointer to a copy of the current value.
func (b *DieBase[T, D]) DieReleasePtr() *T { return core.Must(b.DieReleasePtrE()) }

// DieReleasePtrE is DieReleasePtr, returning an error rather than panicking.
func (b *DieBase[T, D]) DieReleasePtrE() (*T, error) {
	r, err := b.DieReleaseE()
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// DieFeed replaces the value. On a root die the value becomes the new base
// and the intent is discarded. On a child die the value is set as intent.
func (b *DieBase[T, D]) DieFeed(r T) D { return core.Must(b.DieFeedE(r)) }

// DieFeedE is DieFeed, returning an error rather than panicking.
func (b *DieBase[T, D]) DieFeedE(r T) (D, error) {
	h := b.hold()
	v := core.DeepCopy(reflect.ValueOf(&r).Elem())
	if len(b.at) == 0 && h.View == nil {
		if h.Mutable {
			h.Base, h.Ov, h.Wire = v, nil, nil
			return b.self(h), nil
		}
		return b.self(&core.Holder{Base: v, Schema: h.Schema}), nil
	}
	return b.write(h.Write(b.at, nil, core.SetOp(core.OpSet, fit(core.TypeAt(h.Base.Type(), b.at), v)), true))
}

// DieFeedPtr is DieFeed for a pointer; nil feeds the zero value.
func (b *DieBase[T, D]) DieFeedPtr(r *T) D { return core.Must(b.DieFeedPtrE(r)) }

// DieFeedPtrE is DieFeedPtr, returning an error rather than panicking.
func (b *DieBase[T, D]) DieFeedPtrE(r *T) (D, error) {
	if r == nil {
		var zero T
		return b.DieFeedE(zero)
	}
	return b.DieFeedE(*r)
}

// DieFeedJSON feeds a value decoded from JSON. On a root die, as decoding
// does, the values the JSON has are present, zero or not, and fields T
// doesn't have are kept, to be encoded again. Otherwise it is DieFeed: a root
// die discards its intent and a die within another sets the value as intent.
// It panics if the JSON doesn't decode as a T.
func (b *DieBase[T, D]) DieFeedJSON(data []byte) D { return core.Must(b.DieFeedJSONE(data)) }

// DieFeedJSONE is DieFeedJSON, returning an error rather than panicking.
func (b *DieBase[T, D]) DieFeedJSONE(data []byte) (D, error) {
	var zero D
	var r T
	// numbers in unstructured values are int64 where they can be, as
	// apimachinery holds them; structs decode as encoding/json decodes them
	if err := utiljson.Unmarshal(data, &r); err != nil {
		return zero, core.Errorf("die: feeding %s from JSON: %w", reflect.TypeFor[T](), err)
	}
	h := b.hold()
	if len(b.at) > 0 || h.View != nil {
		return b.DieFeedE(r)
	}
	w, err := core.DecodeWire(reflect.TypeFor[T](), data)
	if err != nil {
		return zero, core.Errorf("die: feeding %s from JSON: %w", reflect.TypeFor[T](), err)
	}
	if h.Mutable {
		h.Base, h.Ov, h.Wire = reflect.ValueOf(r), nil, w
		return b.self(h), nil
	}
	return b.self(&core.Holder{Base: reflect.ValueOf(r), Wire: w, Schema: h.Schema}), nil
}

// DieReleaseJSON returns the current value as JSON, as MarshalJSON encodes
// it. It panics if the value doesn't encode.
func (b *DieBase[T, D]) DieReleaseJSON() []byte { return core.Must(b.DieReleaseJSONE()) }

// DieReleaseJSONE is DieReleaseJSON, returning an error rather than
// panicking.
func (b *DieBase[T, D]) DieReleaseJSONE() ([]byte, error) {
	data, err := b.MarshalJSON()
	if err != nil {
		return nil, core.Errorf("die: encoding %s: %w", reflect.TypeFor[T](), err)
	}
	return data, nil
}

// DieFeedYAML feeds a value decoded from YAML, as DieFeedJSON does. It panics
// if the YAML doesn't decode as a T.
func (b *DieBase[T, D]) DieFeedYAML(data []byte) D { return core.Must(b.DieFeedYAMLE(data)) }

// DieFeedYAMLE is DieFeedYAML, returning an error rather than panicking.
func (b *DieBase[T, D]) DieFeedYAMLE(data []byte) (D, error) {
	j, err := yaml.YAMLToJSON(data)
	if err != nil {
		var zero D
		return zero, core.Errorf("die: feeding %s from YAML: %w", reflect.TypeFor[T](), err)
	}
	return b.DieFeedJSONE(j)
}

// DieReleaseYAML returns the current value as YAML, encoded as
// DieReleaseJSON encodes it.
func (b *DieBase[T, D]) DieReleaseYAML() []byte { return core.Must(b.DieReleaseYAMLE()) }

// DieReleaseYAMLE is DieReleaseYAML, returning an error rather than
// panicking.
func (b *DieBase[T, D]) DieReleaseYAMLE() ([]byte, error) {
	j, err := b.DieReleaseJSONE()
	if err != nil {
		return nil, err
	}
	y, err := yaml.JSONToYAML(j)
	if err != nil {
		return nil, core.Errorf("die: encoding %s as YAML: %w", reflect.TypeFor[T](), err)
	}
	return y, nil
}

// DieFeedRawExtension feeds the value of a raw extension, its object or its
// raw JSON, as DieFeedJSON does.
func (b *DieBase[T, D]) DieFeedRawExtension(raw runtime.RawExtension) D {
	return core.Must(b.DieFeedRawExtensionE(raw))
}

// DieFeedRawExtensionE is DieFeedRawExtension, returning an error rather than
// panicking.
func (b *DieBase[T, D]) DieFeedRawExtensionE(raw runtime.RawExtension) (D, error) {
	data, err := raw.MarshalJSON()
	if err != nil {
		var zero D
		return zero, core.Errorf("die: feeding %s from a raw extension: %w", reflect.TypeFor[T](), err)
	}
	return b.DieFeedJSONE(data)
}

// DieReleaseRawExtension returns the current value as a raw extension holding
// its JSON, encoded as DieReleaseJSON encodes it.
func (b *DieBase[T, D]) DieReleaseRawExtension() runtime.RawExtension {
	return core.Must(b.DieReleaseRawExtensionE())
}

// DieReleaseRawExtensionE is DieReleaseRawExtension, returning an error
// rather than panicking.
func (b *DieBase[T, D]) DieReleaseRawExtensionE() (runtime.RawExtension, error) {
	data, err := b.DieReleaseJSONE()
	return runtime.RawExtension{Raw: data}, err
}

// DieStamp passes a copy of the value to fn and records what fn changed as
// intent.
func (b *DieBase[T, D]) DieStamp(fn func(r *T)) D { return core.Must(b.DieStampE(fn)) }

// DieStampE is DieStamp, returning an error rather than panicking.
func (b *DieBase[T, D]) DieStampE(fn func(r *T)) (D, error) {
	var zero D
	h := b.hold()
	before, err := b.DieReleaseE()
	if err != nil {
		return zero, err
	}
	after := as[T](core.DeepCopy(reflect.ValueOf(&before).Elem()))
	fn(&after)
	if !h.Mutable && len(b.at) > 0 {
		return zero, core.ErrImmutableChild
	}
	if core.FirstIndex(b.at) >= 0 {
		// within a list element found by its index, which a write sets whole
		if reflect.DeepEqual(before, after) {
			return b.self(h), nil
		}
		return b.write(h.Write(b.at, nil, core.SetOp(core.OpSet, reflect.ValueOf(after)), false))
	}
	ov, err := core.DiffInto(h.Ov, b.at, reflect.TypeFor[T](), reflect.ValueOf(before), reflect.ValueOf(after))
	if err != nil {
		return zero, err
	}
	return b.write(h.Replace(b.at, ov))
}

// DieWith passes a mutable die to each fn. For an immutable die, the
// mutations are made on a copy that is returned as a new immutable die.
func (b *DieBase[T, D]) DieWith(fns ...func(d D)) D {
	efns := make([]func(D) error, len(fns))
	for i, fn := range fns {
		efns[i] = func(d D) error { fn(d); return nil }
	}
	return core.Must(b.DieWithE(efns...))
}

// DieWithE is DieWith for functions that can fail, such as a client reading
// into a die: it passes a mutable die to each fn, as DieWith does, until one
// returns an error, and returns that error with the die, as the functions
// left it. It returns an error, rather than panicking, where DieWith would,
// on a die within an immutable one, or a jig of one. From an immutable die,
// the functions run on a mutable copy, and the result is a new immutable die.
//
//	current, err := DeploymentBlank.DieWithE(func(d *DeploymentDie) error {
//		return c.Get(ctx, key, d)
//	})
func (b *DieBase[T, D]) DieWithE(fns ...func(d D) error) (D, error) {
	h := b.hold()
	run := func(d D) error {
		for _, fn := range fns {
			if err := fn(d); err != nil {
				return err
			}
		}
		return nil
	}
	if h.Mutable {
		d := b.self(h)
		err := run(d)
		keepSchema(d, h)
		return d, err
	}
	var zero D
	if len(b.at) > 0 {
		return zero, core.ErrImmutableChild
	}
	if h.View != nil {
		return zero, core.ErrReadOnlyJig
	}
	d := wrap[T, D](&core.Holder{Base: h.Base, Ov: h.Ov, Wire: h.Wire, Mutable: true, Schema: h.Schema}, nil)
	err := run(d)
	// the result is what d holds now, which may not be the holder it was
	// given: a client zeroes the object it decodes into, holder and all
	keepSchema(d, h)
	m := d.DieEmbedded().hold()
	return wrap[T, D](&core.Holder{Base: m.Base, Ov: m.Ov, Wire: m.Wire, Schema: m.Schema}, nil), err
}

// keepSchema restores the schema of a die whose state was replaced while it
// was passed to a function, as a client zeroes the object it decodes into.
func keepSchema[T any, D DieOf[T, D]](d D, was *core.Holder) {
	if m := d.DieEmbedded().hold(); m != was && m.Schema == nil {
		m.Schema = was.Schema
	}
}

// DieImmutable returns a copy of the die that is immutable or mutable.
func (b *DieBase[T, D]) DieImmutable(immutable bool) D {
	return core.Must(b.DieImmutableE(immutable))
}

// DieImmutableE is DieImmutable, returning an error rather than panicking.
func (b *DieBase[T, D]) DieImmutableE(immutable bool) (D, error) {
	h, err := b.root("DieImmutable")
	if err != nil {
		var zero D
		return zero, err
	}
	return wrap[T, D](&core.Holder{Base: h.Base, Ov: h.Ov, Wire: h.Wire, Mutable: !immutable, Schema: h.Schema}, nil), nil
}

// DeepCopy returns an independent copy. The base and overlay are shared, as
// neither is ever mutated, so this is cheap. A child die, or a die presented
// by another, such as a jig or a list item, is copied as a new root holding
// its current value, without intent: the copy is unlinked, so writes to it no
// longer reach the source.
func (b *DieBase[T, D]) DeepCopy() D { return core.Must(b.DeepCopyE()) }

// DeepCopyE is DeepCopy, returning an error rather than panicking.
func (b *DieBase[T, D]) DeepCopyE() (D, error) {
	h := b.hold()
	if len(b.at) > 0 || h.View != nil {
		r, err := b.DieReleaseE()
		if err != nil {
			var zero D
			return zero, err
		}
		return wrap[T, D](&core.Holder{Base: reflect.ValueOf(r), Mutable: h.Mutable, Schema: h.Schema.At(b.at)}, nil), nil
	}
	return wrap[T, D](&core.Holder{Base: h.Base, Ov: h.Ov, Wire: h.Wire, Mutable: h.Mutable, Schema: h.Schema}, nil), nil
}

// DieCommit folds the intent into the base, as after a successful write.
func (b *DieBase[T, D]) DieCommit() D { return core.Must(b.DieCommitE()) }

// DieCommitE is DieCommit, returning an error rather than panicking.
func (b *DieBase[T, D]) DieCommitE() (D, error) {
	var zero D
	h, err := b.root("DieCommit")
	if err != nil {
		return zero, err
	}
	v, err := h.Materialize()
	if err != nil {
		return zero, err
	}
	w := core.CommitWire(h.Wire, h.Ov)
	if h.Mutable {
		h.Base, h.Ov, h.Wire = v, nil, w
		return b.self(h), nil
	}
	return b.self(&core.Holder{Base: v, Wire: w, Schema: h.Schema}), nil
}

// DieOnly keeps only the intent at and below the locations, such as
// d.DieOnly(d.Status()). The rest of the intent is folded into the base, so
// the value is unchanged (identity included) but patches only carry the kept
// intent.
func (b *DieBase[T, D]) DieOnly(at ...Located) D { return core.Must(b.DieOnlyE(at...)) }

// DieOnlyE is DieOnly, returning an error rather than panicking.
func (b *DieBase[T, D]) DieOnlyE(at ...Located) (D, error) {
	return b.keep("DieOnly", at, true)
}

// DieExcept drops the intent at and below the locations, folding it into the
// base.
func (b *DieBase[T, D]) DieExcept(at ...Located) D { return core.Must(b.DieExceptE(at...)) }

// DieExceptE is DieExcept, returning an error rather than panicking.
func (b *DieBase[T, D]) DieExceptE(at ...Located) (D, error) {
	return b.keep("DieExcept", at, false)
}

// keep keeps the intent at the locations, or all but it, and folds the rest
// into the base.
func (b *DieBase[T, D]) keep(method string, at []Located, only bool) (D, error) {
	var zero D
	h, err := b.root(method)
	if err != nil {
		return zero, err
	}
	paths, err := pathsOf(h.Base.Type(), at)
	if err != nil {
		return zero, err
	}
	keep, fold := core.Only(h.Ov, paths), core.Except(h.Ov, paths)
	if !only {
		keep, fold = fold, keep
	}
	base, err := (&core.Holder{Base: h.Base, Ov: fold}).Materialize()
	if err != nil {
		return zero, err
	}
	return wrap[T, D](&core.Holder{Base: base, Ov: keep, Wire: core.CommitWire(h.Wire, fold), Mutable: h.Mutable, Schema: h.Schema}, nil), nil
}

// DieSubresource keeps only the intent for a subresource, such as "status".
func (b *DieBase[T, D]) DieSubresource(name string) D { return core.Must(b.DieSubresourceE(name)) }

// DieSubresourceE is DieSubresource, returning an error rather than
// panicking.
func (b *DieBase[T, D]) DieSubresourceE(name string) (D, error) {
	return b.DieOnlyE(subresourcePaths(reflect.TypeFor[T](), func(s string) bool { return s == name })...)
}

// DieMain drops the intent for every subresource, leaving what is written
// through the main resource.
func (b *DieBase[T, D]) DieMain() D { return core.Must(b.DieMainE()) }

// DieMainE is DieMain, returning an error rather than panicking.
func (b *DieBase[T, D]) DieMainE() (D, error) {
	return b.DieExceptE(subresourcePaths(reflect.TypeFor[T](), func(s string) bool { return s != "" })...)
}

func subresourcePaths(t reflect.Type, match func(string) bool) []Located {
	var paths []Located
	for _, fd := range core.Describe(t).Fields {
		if match(fd.Subresource) {
			paths = append(paths, Path{core.FieldStep(fd.Wire)})
		}
	}
	return paths
}

// DieIntent describes the recorded intent, one operation per line.
func (b *DieBase[T, D]) DieIntent() string {
	var lines []string
	describeIntent(core.NodeAt(b.hold().Ov, b.at), b.at, &lines)
	return strings.Join(lines, "\n")
}

func describeIntent(n *core.Onode, p Path, lines *[]string) {
	if n == nil {
		return
	}
	switch n.Op {
	case core.OpSet:
		v, _ := json.Marshal(n.Val.Interface())
		*lines = append(*lines, fmt.Sprintf("set %s = %s", p, v))
	case core.OpClear:
		*lines = append(*lines, fmt.Sprintf("clear %s", p))
	case core.OpEmpty:
		*lines = append(*lines, fmt.Sprintf("empty %s", p))
	}
	for _, id := range n.Order {
		describeIntent(n.Kids[id], core.ChildPath(p, n.Steps[id]), lines)
	}
}

// DieReleaseApplyPatch encodes the intent as a server-side apply request.
func (b *DieBase[T, D]) DieReleaseApplyPatch() []byte { return core.Must(b.DieReleaseApplyPatchE()) }

// DieReleaseApplyPatchE is DieReleaseApplyPatch, returning an error rather
// than panicking.
func (b *DieBase[T, D]) DieReleaseApplyPatchE() ([]byte, error) {
	h, err := b.root("DieReleaseApplyPatch")
	if err != nil {
		return nil, err
	}
	patch, err := h.ApplyPatch()
	if err != nil {
		return nil, core.Errorf("die: encoding the intent of %s as an apply: %w", reflect.TypeFor[T](), err)
	}
	return patch, nil
}

// DieReleaseMergePatch encodes the intent as a JSON merge patch against the base.
func (b *DieBase[T, D]) DieReleaseMergePatch() []byte { return core.Must(b.DieReleaseMergePatchE()) }

// DieReleaseMergePatchE is DieReleaseMergePatch, returning an error rather
// than panicking.
func (b *DieBase[T, D]) DieReleaseMergePatchE() ([]byte, error) {
	h, err := b.root("DieReleaseMergePatch")
	if err != nil {
		return nil, err
	}
	patch, err := h.MergePatch()
	if err != nil {
		return nil, core.Errorf("die: encoding the intent of %s as a merge patch: %w", reflect.TypeFor[T](), err)
	}
	return patch, nil
}

// DieFeedApplyPatch feeds a server-side apply request as intent, over the
// die's value, which is kept: each field it has is set, null clears a field,
// [] and {} are an empty list or map, list-map elements are found by their
// keys, and an atomic struct is set whole. The intent replaces the die's, so
// DieReleaseApplyPatch sends the request back. Fields the die's type doesn't
// have aren't fed. An unstructured die follows its mold, if it has one;
// without, its lists are set whole. It panics if the request isn't a JSON
// object.
func (b *DieBase[T, D]) DieFeedApplyPatch(patch []byte) D {
	return core.Must(b.DieFeedApplyPatchE(patch))
}

// DieFeedApplyPatchE is DieFeedApplyPatch, returning an error rather than
// panicking.
func (b *DieBase[T, D]) DieFeedApplyPatchE(patch []byte) (D, error) {
	return b.feedPatch("DieFeedApplyPatch", core.ApplyPatch, patch)
}

// DieFeedMergePatch feeds a JSON merge patch, RFC 7386, as intent over the
// die's value, which is kept: null clears a field, objects and maps are
// merged field by field, and anything else, lists included, is set whole. The
// intent replaces the die's, so DieReleaseMergePatch sends the patch back,
// but for what the patch sets that the value already has. Fields the die's
// type doesn't have aren't fed. It panics if the patch isn't a JSON object.
func (b *DieBase[T, D]) DieFeedMergePatch(patch []byte) D {
	return core.Must(b.DieFeedMergePatchE(patch))
}

// DieFeedMergePatchE is DieFeedMergePatch, returning an error rather than
// panicking.
func (b *DieBase[T, D]) DieFeedMergePatchE(patch []byte) (D, error) {
	return b.feedPatch("DieFeedMergePatch", core.MergePatch, patch)
}

func (b *DieBase[T, D]) feedPatch(method string, kind core.PatchKind, patch []byte) (D, error) {
	var zero D
	h, err := b.root(method)
	if err != nil {
		return zero, err
	}
	var raw map[string]any
	// numbers are int64 where they can be, as unstructured values hold them
	if err := utiljson.Unmarshal(patch, &raw); err != nil {
		return zero, core.Errorf("die: %s: %w", method, err)
	}
	return b.write(core.FeedPatch(h, kind, raw))
}

// Equal compares current values. Intent is not compared.
func (b *DieBase[T, D]) Equal(o D) bool { return core.Must(b.EqualE(o)) }

// EqualE is Equal, returning an error rather than panicking.
func (b *DieBase[T, D]) EqualE(o D) (bool, error) {
	mine, err := b.DieReleaseE()
	if err != nil {
		return false, err
	}
	theirs, err := o.DieReleaseE()
	if err != nil {
		return false, err
	}
	return reflect.DeepEqual(mine, theirs), nil
}

// MarshalJSON encodes the current value as its presence says: values that
// were decoded as they were on the wire, fields its type doesn't have
// included, and written values as Go has them, nil and omitempty zero values
// omitted.
func (b *DieBase[T, D]) MarshalJSON() ([]byte, error) {
	v, err := b.hold().EncodeAt(b.at)
	if err != nil {
		return nil, err
	}
	return json.Marshal(v)
}

// UnmarshalJSON replaces the base and keeps the intent, so the intent is
// reapplied on top of the decoded value.
func (b *DieBase[T, D]) UnmarshalJSON(data []byte) error {
	h := b.hold()
	if !h.Mutable || len(b.at) > 0 || h.View != nil {
		return core.Errorf("die: can only unmarshal into a mutable root die")
	}
	var r T
	// numbers in unstructured values are int64 where they can be, as
	// apimachinery holds them; structs decode as encoding/json decodes them
	if err := utiljson.Unmarshal(data, &r); err != nil {
		return err
	}
	w, err := core.DecodeWire(reflect.TypeFor[T](), data)
	if err != nil {
		return err
	}
	h.Base, h.Wire = reflect.ValueOf(r), w
	return nil
}

func as[T any](v reflect.Value) T {
	if v.Kind() == reflect.Pointer && v.Type().Elem() == reflect.TypeFor[T]() {
		if v.IsNil() {
			var zero T
			return zero
		}
		v = v.Elem()
	}
	return v.Interface().(T)
}

// fit adapts a value to the static type at its location, which may be a
// pointer to it.
func fit(t reflect.Type, v reflect.Value) reflect.Value {
	if t != v.Type() && t.Kind() == reflect.Pointer && t.Elem() == v.Type() {
		p := reflect.New(v.Type())
		p.Elem().Set(v)
		return p
	}
	return v
}

// Error is a die's error. Every method that can fail has an E variant, which
// returns one, and a fluent form, which panics with it: the fluent method for
// content that is known, as in tests and fixtures, and the E variant
// otherwise. Use errors.As to tell a die's error from another, and errors.Is
// with ErrImmutableChild or ErrReadOnlyJig for those two.
//
// Methods that implement Kubernetes' and Go's interfaces, such as
// metav1.Object's setters, can't return an error: they panic with one, such
// as on an immutable die.
type Error = core.Error

var (
	// ErrImmutableChild is the error of a write through a child of an
	// immutable die, which would have nowhere to go: write within DieWith.
	ErrImmutableChild error = core.ErrImmutableChild
	// ErrReadOnlyJig is the error of a write through a jig of an immutable
	// die.
	ErrReadOnlyJig error = core.ErrReadOnlyJig
)
