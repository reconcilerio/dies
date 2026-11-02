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
// Every die type comes in two: an immutable die, such as DeploymentDie, which
// only reads, and a mutable die, such as DeploymentDieMut, which also
// writes, in place. Blanks are immutable; DieWith passes a mutable copy to
// functions and returns an immutable die of the result. DieMutable and
// DieImmutable convert, always by copying, so an immutable die never sees a
// write.
//
// A child die returned by an accessor is attached to its root: it addresses
// into the root by path, reads never create anything, and the first write
// creates whatever the path needs. The children of an immutable die are
// immutable, and those of a mutable die mutable.
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

// dieOf is what immutable and mutable dies of a T share: reading.
type dieOf[T any] interface {
	Die
	LocatedAs[T]

	DieRelease() T
	DieReleaseE() (T, error)
	DieReleasePtr() *T
	DieReleasePtrE() (*T, error)
}

// DieOf carries an immutable die's types for generic code: T is the die's
// value type, D its own pointer type, so its methods return the concrete
// die, and M the pointer type of its mutable counterpart. It has every method
// of an immutable die. Generic code constrains a die by it,
// D DieOf[T, D, M]; code that takes a die of any type takes a Die.
type DieOf[T any, D DieOf[T, D, M], M DieMutOf[T, D, M]] interface {
	dieOf[T]
	// DieFrom wraps state in a new die. It is called on a nil receiver.
	DieFrom(DieBase[T, D, M]) D

	DieWith(fns ...func(m M)) D
	DieWithE(fns ...func(m M) error) (D, error)
	DieMutable() M
	DieMutableE() (M, error)
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

// DieMutOf carries a mutable die's types for generic code, as DieOf does an
// immutable die's. It has every method of a mutable die.
type DieMutOf[T any, D DieOf[T, D, M], M DieMutOf[T, D, M]] interface {
	dieOf[T]
	json.Unmarshaler
	// DieFrom wraps state in a new die. It is called on a nil receiver.
	DieFrom(DieMutBase[T, D, M]) M

	DieClear() M
	DieClearE() (M, error)
	DieFeed(r T) M
	DieFeedE(r T) (M, error)
	DieFeedPtr(r *T) M
	DieFeedPtrE(r *T) (M, error)
	DieFeedJSON(data []byte) M
	DieFeedJSONE(data []byte) (M, error)
	DieFeedYAML(data []byte) M
	DieFeedYAMLE(data []byte) (M, error)
	DieFeedRawExtension(raw runtime.RawExtension) M
	DieFeedRawExtensionE(raw runtime.RawExtension) (M, error)
	DieFeedApplyPatch(patch []byte) M
	DieFeedApplyPatchE(patch []byte) (M, error)
	DieFeedMergePatch(patch []byte) M
	DieFeedMergePatchE(patch []byte) (M, error)
	DieStamp(fn func(r *T)) M
	DieStampE(fn func(r *T)) (M, error)
	DieWith(fns ...func(m M)) M
	DieWithE(fns ...func(m M) error) (M, error)
	DieImmutable() D
	DieImmutableE() (D, error)
	DeepCopy() M
	DeepCopyE() (M, error)
	DieCommit() M
	DieCommitE() (M, error)
	DieOnly(at ...Located) M
	DieOnlyE(at ...Located) (M, error)
	DieExcept(at ...Located) M
	DieExceptE(at ...Located) (M, error)
	DieSubresource(name string) M
	DieSubresourceE(name string) (M, error)
	DieMain() M
	DieMainE() (M, error)
	Equal(o M) bool
	EqualE(o M) (bool, error)
}

// Die is a die of any type, immutable or mutable: the methods every die has
// that don't name its types, for functions that take any die. Handles,
// collections and nodes have some of them, but aren't dies. DieOf and
// DieMutOf have the rest, for generic code.
type Die interface {
	// MarshalJSON encodes the die's value as its presence says.
	json.Marshaler

	DiePath() Path
	DieRootType() reflect.Type
	DieRoot() Die
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

// ImmutableDie is an immutable die of any type, for functions that take
// only immutable dies, such as jigs.Jig.Of. Only the dies of DieBase, and the
// typed and object bases that embed it, are ImmutableDies.
type ImmutableDie interface {
	Die
	dieImmutableMarker()
}

// MutableDie is a mutable die of any type, for functions that take only
// mutable dies, such as jigs.Jig.OfMut. Only the dies of DieMutBase, and the
// typed and object bases that embed it, are MutableDies.
type MutableDie interface {
	Die
	dieMutableMarker()
}

// IsMutable reports whether a die is a mutable die.
func IsMutable(d any) bool {
	_, ok := d.(MutableDie)
	return ok
}

// state is the state of a die, immutable or mutable, with the methods both
// have, which read. DieBase and DieMutBase embed it, rather than one the
// other, so that each has only its own methods: DieMutable only on an
// immutable die, and DieImmutable only on a mutable one.
type state[T any] struct {
	h  *core.Holder
	at Path
}

// DieBase holds the state of an immutable die and provides its methods,
// which read. Immutable dies embed it by value.
type DieBase[T any, D DieOf[T, D, M], M DieMutOf[T, D, M]] struct {
	state[T]
}

// DieMutBase holds the state of a mutable die and provides its methods. Its
// writes change the die in place, and return it, so they chain. Mutable dies
// embed it by value.
//
// The zero value is a mutable, empty die, so a die allocated by a decoder can
// be unmarshaled into.
type DieMutBase[T any, D DieOf[T, D, M], M DieMutOf[T, D, M]] struct {
	state[T]
}

func wrap[T any, D DieOf[T, D, M], M DieMutOf[T, D, M]](h *core.Holder, at Path) D {
	if len(at) == 0 {
		h.SetRoot(func(h *core.Holder) any { return wrap[T, D, M](h, nil) })
	}
	var d D
	return d.DieFrom(DieBase[T, D, M]{state[T]{h: h, at: at}})
}

func wrapMut[T any, D DieOf[T, D, M], M DieMutOf[T, D, M]](h *core.Holder, at Path) M {
	if len(at) == 0 {
		h.SetRoot(func(h *core.Holder) any { return wrapMut[T, D, M](h, nil) })
	}
	var m M
	return m.DieFrom(DieMutBase[T, D, M]{state[T]{h: h, at: at}})
}

// BlankOf returns an immutable die of a value, such as an object with only
// its type metadata, for a blank.
func BlankOf[T any, D DieOf[T, D, M], M DieMutOf[T, D, M]](v T) D {
	return wrap[T, D, M](&core.Holder{Base: reflect.ValueOf(v)}, nil)
}

// Blank returns an immutable, empty die.
func Blank[T any, D DieOf[T, D, M], M DieMutOf[T, D, M]]() D {
	return wrap[T, D, M](&core.Holder{Base: reflect.Zero(reflect.TypeFor[T]())}, nil)
}

// DieValue marks a die as a location of a T, for LocatedAs. It does nothing.
func (*state[T]) DieValue(*T) {}

// dieState is how the packages of handles, jigs and unstructured dies reach
// a die's state, through the functions installed below.
type dieState interface {
	hold() *core.Holder
	dieAt() Path
	setDieState(h *core.Holder, at Path)
	dieValueType() reflect.Type
}

func (b *state[T]) dieAt() Path { return b.at }

func (b *state[T]) setDieState(h *core.Holder, at Path) { b.h, b.at = h, at }

func (b *state[T]) dieValueType() reflect.Type { return reflect.TypeFor[T]() }

func (*DieBase[T, D, M]) dieImmutableMarker() {}

func (*DieMutBase[T, D, M]) dieMutableMarker() {}

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
	core.IsMutable = IsMutable
}

// hold returns the die's holder. A zero die's is made on first use, mutable,
// for a decoder; an immutable die never writes to it.
func (b *state[T]) hold() *core.Holder {
	if b.h == nil {
		var zero T
		b.h = &core.Holder{Base: reflect.ValueOf(zero), Mutable: true}
	}
	return b.h
}

func (b *state[T]) root(method string) (*core.Holder, error) {
	if len(b.at) > 0 {
		return nil, core.Errorf("die: %s is only supported on a root die, not at %s", method, b.at)
	}
	if b.hold().View != nil {
		return nil, core.Errorf("die: %s is not supported on a jig; use it on the jig's source", method)
	}
	return b.hold(), nil
}

// copyOf returns a new root holder with a copy of the die: a root's value
// and intent, or the current value of a die within another or presented by a
// jig, without intent, so the copy is unlinked from its source.
func (b *state[T]) copyOf(mutable bool) (*core.Holder, error) {
	h := b.hold()
	if len(b.at) > 0 || h.View != nil {
		r, err := b.DieReleaseE()
		if err != nil {
			return nil, err
		}
		return &core.Holder{Base: reflect.ValueOf(r), Mutable: mutable, Schema: h.Schema.At(b.at)}, nil
	}
	return &core.Holder{Base: h.Base, Ov: h.Ov, Wire: h.Wire, Mutable: mutable, Schema: h.Schema}, nil
}

// DiePath returns the location of this die within its root.
func (b *state[T]) DiePath() Path { return b.at }

// DieRoot returns the root die this die is within, sharing its state. A root
// returns itself. Its type is the root's die type, such as *DeploymentDie
// for a die within an immutable Deployment, or *DeploymentDieMut within a
// mutable one, whose writes are writes to this die's tree. Assert it as a
// type, or as metav1.Object for an object's identity whatever its type:
// d.DieRoot().(metav1.Object).GetNamespace().
func (b *DieBase[T, D, M]) DieRoot() Die {
	if len(b.at) == 0 {
		return wrap[T, D, M](b.hold(), nil)
	}
	return RootOf(b.hold())
}

// DieRoot returns the root die this die is within: see DieBase.DieRoot.
func (b *DieMutBase[T, D, M]) DieRoot() Die {
	if len(b.at) == 0 {
		return wrapMut[T, D, M](b.hold(), nil)
	}
	return RootOf(b.hold())
}

// RootOf returns the root die over a holder, or nil if none was made.
func RootOf(h *core.Holder) Die {
	root, _ := h.Root().(Die)
	return root
}

// DieRootType returns the type of the root value this die is within.
func (b *state[T]) DieRootType() reflect.Type { return b.hold().Base.Type() }

// DieHas reports whether the value is present. A root is always present.
func (b *state[T]) DieHas() bool { return core.Must(b.DieHasE()) }

// DieHasE is DieHas, returning an error rather than panicking.
func (b *state[T]) DieHasE() (bool, error) {
	if len(b.at) == 0 {
		return true, nil
	}
	return b.hold().Has(b.at)
}

// DieRelease returns a copy of the current value.
func (b *state[T]) DieRelease() T { return core.Must(b.DieReleaseE()) }

// DieReleaseE is DieRelease, returning an error rather than panicking.
func (b *state[T]) DieReleaseE() (T, error) {
	v, err := b.hold().ValueAt(b.at)
	if err != nil {
		var zero T
		return zero, err
	}
	return as[T](v), nil
}

// DieReleasePtr returns a pointer to a copy of the current value.
func (b *state[T]) DieReleasePtr() *T { return core.Must(b.DieReleasePtrE()) }

// DieReleasePtrE is DieReleasePtr, returning an error rather than panicking.
func (b *state[T]) DieReleasePtrE() (*T, error) {
	r, err := b.DieReleaseE()
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// DieReleaseJSON returns the current value as JSON, as MarshalJSON encodes
// it. It panics if the value doesn't encode.
func (b *state[T]) DieReleaseJSON() []byte { return core.Must(b.DieReleaseJSONE()) }

// DieReleaseJSONE is DieReleaseJSON, returning an error rather than
// panicking.
func (b *state[T]) DieReleaseJSONE() ([]byte, error) {
	data, err := b.MarshalJSON()
	if err != nil {
		return nil, core.Errorf("die: encoding %s: %w", reflect.TypeFor[T](), err)
	}
	return data, nil
}

// DieReleaseYAML returns the current value as YAML, encoded as
// DieReleaseJSON encodes it.
func (b *state[T]) DieReleaseYAML() []byte { return core.Must(b.DieReleaseYAMLE()) }

// DieReleaseYAMLE is DieReleaseYAML, returning an error rather than
// panicking.
func (b *state[T]) DieReleaseYAMLE() ([]byte, error) {
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

// DieReleaseRawExtension returns the current value as a raw extension holding
// its JSON, encoded as DieReleaseJSON encodes it.
func (b *state[T]) DieReleaseRawExtension() runtime.RawExtension {
	return core.Must(b.DieReleaseRawExtensionE())
}

// DieReleaseRawExtensionE is DieReleaseRawExtension, returning an error
// rather than panicking.
func (b *state[T]) DieReleaseRawExtensionE() (runtime.RawExtension, error) {
	data, err := b.DieReleaseJSONE()
	return runtime.RawExtension{Raw: data}, err
}

// DieWith passes a mutable copy of the die to each fn, and returns an
// immutable die of the result. The copy of a root has its intent; the copy
// of a die within another, or presented by a jig, is a new root holding its
// current value, as DeepCopy makes it.
func (b *DieBase[T, D, M]) DieWith(fns ...func(m M)) D {
	return core.Must(b.DieWithE(withE(fns)...))
}

// DieWithE is DieWith for functions that can fail, such as a client reading
// into a die: it passes a mutable copy to each fn, as DieWith does, until one
// returns an error, and returns that error with an immutable die of the
// copy, as the functions left it.
//
//	current, err := DeploymentBlank.DieWithE(func(d *DeploymentDieMut) error {
//		return c.Get(ctx, key, d)
//	})
func (b *DieBase[T, D, M]) DieWithE(fns ...func(m M) error) (D, error) {
	h, err := b.copyOf(true)
	if err != nil {
		var zero D
		return zero, err
	}
	m := wrapMut[T, D, M](h, nil)
	err = runWith(m, fns)
	// the result is what m holds now, which may not be the holder it was
	// given: a client zeroes the object it decodes into, holder and all
	r := held(m, h)
	return wrap[T, D, M](&core.Holder{Base: r.Base, Ov: r.Ov, Wire: r.Wire, Schema: r.Schema}, nil), err
}

func withE[M any](fns []func(m M)) []func(M) error {
	efns := make([]func(M) error, len(fns))
	for i, fn := range fns {
		efns[i] = func(m M) error { fn(m); return nil }
	}
	return efns
}

func runWith[M any](m M, fns []func(M) error) error {
	for _, fn := range fns {
		if err := fn(m); err != nil {
			return err
		}
	}
	return nil
}

// held returns the holder of a die passed to functions, restoring the schema
// of the holder it was given if its state was replaced, as a client zeroes
// the object it decodes into.
func held(d any, was *core.Holder) *core.Holder {
	m, _, _ := core.StateOf(d)
	if m != was && m.Schema == nil {
		m.Schema = was.Schema
	}
	return m
}

// DieMutable returns a mutable copy of the die: of a root, its value and
// intent; of a die within another, or presented by a jig, a new root holding
// its current value, as DeepCopy makes it. Writes to the copy never reach
// the die.
func (b *DieBase[T, D, M]) DieMutable() M { return core.Must(b.DieMutableE()) }

// DieMutableE is DieMutable, returning an error rather than panicking.
func (b *DieBase[T, D, M]) DieMutableE() (M, error) {
	h, err := b.copyOf(true)
	if err != nil {
		var zero M
		return zero, err
	}
	return wrapMut[T, D, M](h, nil), nil
}

// DeepCopy returns an independent copy. The base and overlay are shared, as
// neither is ever mutated, so this is cheap. A child die, or a die presented
// by another, such as a jig or a list item, is copied as a new root holding
// its current value, without intent: the copy is unlinked from its source.
func (b *DieBase[T, D, M]) DeepCopy() D { return core.Must(b.DeepCopyE()) }

// DeepCopyE is DeepCopy, returning an error rather than panicking.
func (b *DieBase[T, D, M]) DeepCopyE() (D, error) {
	h, err := b.copyOf(false)
	if err != nil {
		var zero D
		return zero, err
	}
	return wrap[T, D, M](h, nil), nil
}

// DieCommit returns the die with its intent folded into its base, as after a
// successful write.
func (b *DieBase[T, D, M]) DieCommit() D { return core.Must(b.DieCommitE()) }

// DieCommitE is DieCommit, returning an error rather than panicking.
func (b *DieBase[T, D, M]) DieCommitE() (D, error) {
	h, err := b.commit()
	if err != nil {
		var zero D
		return zero, err
	}
	return wrap[T, D, M](h, nil), nil
}

// commit returns a new holder of the die's value with its intent folded in.
func (b *state[T]) commit() (*core.Holder, error) {
	h, err := b.root("DieCommit")
	if err != nil {
		return nil, err
	}
	v, err := h.Materialize()
	if err != nil {
		return nil, err
	}
	return &core.Holder{Base: v, Wire: core.CommitWire(h.Wire, h.Ov), Mutable: h.Mutable, Schema: h.Schema}, nil
}

// DieOnly returns the die keeping only the intent at and below the
// locations, such as d.DieOnly(d.Status()). The rest of the intent is folded
// into the base, so the value is unchanged (identity included) but patches
// only carry the kept intent.
func (b *DieBase[T, D, M]) DieOnly(at ...Located) D { return core.Must(b.DieOnlyE(at...)) }

// DieOnlyE is DieOnly, returning an error rather than panicking.
func (b *DieBase[T, D, M]) DieOnlyE(at ...Located) (D, error) {
	h, err := b.keep("DieOnly", at, true)
	return wrapped(h, err, wrap[T, D, M])
}

// DieExcept returns the die with the intent at and below the locations
// folded into its base.
func (b *DieBase[T, D, M]) DieExcept(at ...Located) D { return core.Must(b.DieExceptE(at...)) }

// DieExceptE is DieExcept, returning an error rather than panicking.
func (b *DieBase[T, D, M]) DieExceptE(at ...Located) (D, error) {
	h, err := b.keep("DieExcept", at, false)
	return wrapped(h, err, wrap[T, D, M])
}

// DieSubresource returns the die keeping only the intent for a subresource,
// such as "status".
func (b *DieBase[T, D, M]) DieSubresource(name string) D { return core.Must(b.DieSubresourceE(name)) }

// DieSubresourceE is DieSubresource, returning an error rather than
// panicking.
func (b *DieBase[T, D, M]) DieSubresourceE(name string) (D, error) {
	return b.DieOnlyE(subresourcePaths[T](name)...)
}

// DieMain returns the die without the intent for any subresource, leaving
// what is written through the main resource.
func (b *DieBase[T, D, M]) DieMain() D { return core.Must(b.DieMainE()) }

// DieMainE is DieMain, returning an error rather than panicking.
func (b *DieBase[T, D, M]) DieMainE() (D, error) {
	return b.DieExceptE(mainPaths[T]()...)
}

func wrapped[X any](h *core.Holder, err error, wrap func(*core.Holder, Path) X) (X, error) {
	if err != nil {
		var zero X
		return zero, err
	}
	return wrap(h, nil), nil
}

// keep returns a new holder keeping the intent at the locations, or all but
// it, with the rest folded into the base.
func (b *state[T]) keep(method string, at []Located, only bool) (*core.Holder, error) {
	h, err := b.root(method)
	if err != nil {
		return nil, err
	}
	paths, err := pathsOf(h.Base.Type(), at)
	if err != nil {
		return nil, err
	}
	keep, fold := core.Only(h.Ov, paths), core.Except(h.Ov, paths)
	if !only {
		keep, fold = fold, keep
	}
	base, err := (&core.Holder{Base: h.Base, Ov: fold}).Materialize()
	if err != nil {
		return nil, err
	}
	return &core.Holder{Base: base, Ov: keep, Wire: core.CommitWire(h.Wire, fold), Mutable: h.Mutable, Schema: h.Schema}, nil
}

func subresourcePaths[T any](name string) []Located {
	return subresourcesWhere(reflect.TypeFor[T](), func(s string) bool { return s == name })
}

func mainPaths[T any]() []Located {
	return subresourcesWhere(reflect.TypeFor[T](), func(s string) bool { return s != "" })
}

func subresourcesWhere(t reflect.Type, match func(string) bool) []Located {
	var paths []Located
	for _, fd := range core.Describe(t).Fields {
		if match(fd.Subresource) {
			paths = append(paths, Path{core.FieldStep(fd.Wire)})
		}
	}
	return paths
}

// DieIntent describes the recorded intent, one operation per line.
func (b *state[T]) DieIntent() string {
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
func (b *state[T]) DieReleaseApplyPatch() []byte {
	return core.Must(b.DieReleaseApplyPatchE())
}

// DieReleaseApplyPatchE is DieReleaseApplyPatch, returning an error rather
// than panicking.
func (b *state[T]) DieReleaseApplyPatchE() ([]byte, error) {
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

// DieReleaseMergePatch encodes the intent as a JSON merge patch against the
// base.
func (b *state[T]) DieReleaseMergePatch() []byte {
	return core.Must(b.DieReleaseMergePatchE())
}

// DieReleaseMergePatchE is DieReleaseMergePatch, returning an error rather
// than panicking.
func (b *state[T]) DieReleaseMergePatchE() ([]byte, error) {
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

// Equal compares current values. Intent is not compared.
func (b *DieBase[T, D, M]) Equal(o D) bool { return core.Must(b.EqualE(o)) }

// EqualE is Equal, returning an error rather than panicking.
func (b *DieBase[T, D, M]) EqualE(o D) (bool, error) { return equal[T](b, o) }

func equal[T any](a, b interface{ DieReleaseE() (T, error) }) (bool, error) {
	mine, err := a.DieReleaseE()
	if err != nil {
		return false, err
	}
	theirs, err := b.DieReleaseE()
	if err != nil {
		return false, err
	}
	return reflect.DeepEqual(mine, theirs), nil
}

// MarshalJSON encodes the current value as its presence says: values that
// were decoded as they were on the wire, fields its type doesn't have
// included, and written values as Go has them, nil and omitempty zero values
// omitted.
func (b *state[T]) MarshalJSON() ([]byte, error) {
	v, err := b.hold().EncodeAt(b.at)
	if err != nil {
		return nil, err
	}
	return json.Marshal(v)
}

// UnmarshalJSON fails with ErrImmutable: an immutable die can't be decoded
// into. Decode into its mutable die, such as within DieWithE.
func (b *DieBase[T, D, M]) UnmarshalJSON([]byte) error {
	return core.Errorf("die: decoding into an immutable %s die: %w", reflect.TypeFor[T](), core.ErrImmutable)
}

// The methods of mutable dies: those that write, which change the die in
// place, and those of DieBase that return a die, returning a mutable one.

// self returns the mutable die at this die's location.
func (b *DieMutBase[T, D, M]) self() M { return wrapMut[T, D, M](b.hold(), b.at) }

// write returns this die after a write, or the write's error.
func (b *DieMutBase[T, D, M]) write(_ *core.Holder, err error) (M, error) {
	if err != nil {
		var zero M
		return zero, err
	}
	return b.self(), nil
}

// DieClear records that the value is absent: null in patches, or a removed
// element. On a root it resets to the zero value.
func (b *DieMutBase[T, D, M]) DieClear() M { return core.Must(b.DieClearE()) }

// DieClearE is DieClear, returning an error rather than panicking.
func (b *DieMutBase[T, D, M]) DieClearE() (M, error) {
	if len(b.at) == 0 {
		var zero T
		return b.DieFeedE(zero)
	}
	return b.write(b.hold().Write(b.at, nil, core.SetOp(core.OpClear, reflect.Value{}), false))
}

// DieFeed replaces the value. On a root die the value becomes the new base
// and the intent is discarded. On a child die the value is set as intent.
func (b *DieMutBase[T, D, M]) DieFeed(r T) M { return core.Must(b.DieFeedE(r)) }

// DieFeedE is DieFeed, returning an error rather than panicking.
func (b *DieMutBase[T, D, M]) DieFeedE(r T) (M, error) {
	h := b.hold()
	v := core.DeepCopy(reflect.ValueOf(&r).Elem())
	if len(b.at) == 0 && h.View == nil {
		h.Base, h.Ov, h.Wire = v, nil, nil
		return b.self(), nil
	}
	return b.write(h.Write(b.at, nil, core.SetOp(core.OpSet, fit(core.TypeAt(h.Base.Type(), b.at), v)), true))
}

// DieFeedPtr is DieFeed for a pointer; nil feeds the zero value.
func (b *DieMutBase[T, D, M]) DieFeedPtr(r *T) M { return core.Must(b.DieFeedPtrE(r)) }

// DieFeedPtrE is DieFeedPtr, returning an error rather than panicking.
func (b *DieMutBase[T, D, M]) DieFeedPtrE(r *T) (M, error) {
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
func (b *DieMutBase[T, D, M]) DieFeedJSON(data []byte) M { return core.Must(b.DieFeedJSONE(data)) }

// DieFeedJSONE is DieFeedJSON, returning an error rather than panicking.
func (b *DieMutBase[T, D, M]) DieFeedJSONE(data []byte) (M, error) {
	var zero M
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
	h.Base, h.Ov, h.Wire = reflect.ValueOf(r), nil, w
	return b.self(), nil
}

// DieFeedYAML feeds a value decoded from YAML, as DieFeedJSON does. It panics
// if the YAML doesn't decode as a T.
func (b *DieMutBase[T, D, M]) DieFeedYAML(data []byte) M { return core.Must(b.DieFeedYAMLE(data)) }

// DieFeedYAMLE is DieFeedYAML, returning an error rather than panicking.
func (b *DieMutBase[T, D, M]) DieFeedYAMLE(data []byte) (M, error) {
	j, err := yaml.YAMLToJSON(data)
	if err != nil {
		var zero M
		return zero, core.Errorf("die: feeding %s from YAML: %w", reflect.TypeFor[T](), err)
	}
	return b.DieFeedJSONE(j)
}

// DieFeedRawExtension feeds the value of a raw extension, its object or its
// raw JSON, as DieFeedJSON does.
func (b *DieMutBase[T, D, M]) DieFeedRawExtension(raw runtime.RawExtension) M {
	return core.Must(b.DieFeedRawExtensionE(raw))
}

// DieFeedRawExtensionE is DieFeedRawExtension, returning an error rather than
// panicking.
func (b *DieMutBase[T, D, M]) DieFeedRawExtensionE(raw runtime.RawExtension) (M, error) {
	data, err := raw.MarshalJSON()
	if err != nil {
		var zero M
		return zero, core.Errorf("die: feeding %s from a raw extension: %w", reflect.TypeFor[T](), err)
	}
	return b.DieFeedJSONE(data)
}

// DieStamp passes a copy of the value to fn and records what fn changed as
// intent.
func (b *DieMutBase[T, D, M]) DieStamp(fn func(r *T)) M { return core.Must(b.DieStampE(fn)) }

// DieStampE is DieStamp, returning an error rather than panicking.
func (b *DieMutBase[T, D, M]) DieStampE(fn func(r *T)) (M, error) {
	var zero M
	h := b.hold()
	before, err := b.DieReleaseE()
	if err != nil {
		return zero, err
	}
	after := as[T](core.DeepCopy(reflect.ValueOf(&before).Elem()))
	fn(&after)
	if core.FirstIndex(b.at) >= 0 {
		// within a list element found by its index, which a write sets whole
		if reflect.DeepEqual(before, after) {
			return b.self(), nil
		}
		return b.write(h.Write(b.at, nil, core.SetOp(core.OpSet, reflect.ValueOf(after)), false))
	}
	ov, err := core.DiffInto(h.Ov, b.at, reflect.TypeFor[T](), reflect.ValueOf(before), reflect.ValueOf(after))
	if err != nil {
		return zero, err
	}
	return b.write(h.Replace(b.at, ov))
}

// DieWith passes the die itself to each fn, so their writes change it, and
// returns it.
func (b *DieMutBase[T, D, M]) DieWith(fns ...func(m M)) M {
	return core.Must(b.DieWithE(withE(fns)...))
}

// DieWithE is DieWith for functions that can fail: it passes the die to each
// fn until one returns an error, and returns that error with the die, as the
// functions left it.
func (b *DieMutBase[T, D, M]) DieWithE(fns ...func(m M) error) (M, error) {
	h := b.hold()
	m := b.self()
	err := runWith(m, fns)
	// a client zeroes the object it decodes into, holder and all: keep what
	// it decoded in this die's holder, which its children share
	if r := held(m, h); r != h {
		h.Base, h.Ov, h.Wire = r.Base, r.Ov, r.Wire
		core.SetState(m, h, b.at)
	}
	return m, err
}

// DieImmutable returns an immutable copy of the die: of a root, its value
// and intent; of a die within another, or presented by a jig, a new root
// holding its current value, as DeepCopy makes it. Later writes to the die
// never reach the copy.
func (b *DieMutBase[T, D, M]) DieImmutable() D { return core.Must(b.DieImmutableE()) }

// DieImmutableE is DieImmutable, returning an error rather than panicking.
func (b *DieMutBase[T, D, M]) DieImmutableE() (D, error) {
	h, err := b.copyOf(false)
	if err != nil {
		var zero D
		return zero, err
	}
	return wrap[T, D, M](h, nil), nil
}

// DeepCopy returns an independent, mutable copy: see DieBase.DeepCopy.
func (b *DieMutBase[T, D, M]) DeepCopy() M { return core.Must(b.DeepCopyE()) }

// DeepCopyE is DeepCopy, returning an error rather than panicking.
func (b *DieMutBase[T, D, M]) DeepCopyE() (M, error) {
	h, err := b.copyOf(true)
	if err != nil {
		var zero M
		return zero, err
	}
	return wrapMut[T, D, M](h, nil), nil
}

// DieCommit folds the intent into the base, in place, as after a successful
// write.
func (b *DieMutBase[T, D, M]) DieCommit() M { return core.Must(b.DieCommitE()) }

// DieCommitE is DieCommit, returning an error rather than panicking.
func (b *DieMutBase[T, D, M]) DieCommitE() (M, error) {
	return b.replace(b.commit())
}

// replace replaces the die's state with a holder's, in place.
func (b *DieMutBase[T, D, M]) replace(h *core.Holder, err error) (M, error) {
	if err != nil {
		var zero M
		return zero, err
	}
	m := b.hold()
	m.Base, m.Ov, m.Wire = h.Base, h.Ov, h.Wire
	return b.self(), nil
}

// DieOnly keeps only the intent at and below the locations, in place: see
// DieBase.DieOnly.
func (b *DieMutBase[T, D, M]) DieOnly(at ...Located) M { return core.Must(b.DieOnlyE(at...)) }

// DieOnlyE is DieOnly, returning an error rather than panicking.
func (b *DieMutBase[T, D, M]) DieOnlyE(at ...Located) (M, error) {
	return b.replace(b.keep("DieOnly", at, true))
}

// DieExcept folds the intent at and below the locations into the base, in
// place.
func (b *DieMutBase[T, D, M]) DieExcept(at ...Located) M { return core.Must(b.DieExceptE(at...)) }

// DieExceptE is DieExcept, returning an error rather than panicking.
func (b *DieMutBase[T, D, M]) DieExceptE(at ...Located) (M, error) {
	return b.replace(b.keep("DieExcept", at, false))
}

// DieSubresource keeps only the intent for a subresource, such as "status",
// in place.
func (b *DieMutBase[T, D, M]) DieSubresource(name string) M {
	return core.Must(b.DieSubresourceE(name))
}

// DieSubresourceE is DieSubresource, returning an error rather than
// panicking.
func (b *DieMutBase[T, D, M]) DieSubresourceE(name string) (M, error) {
	return b.DieOnlyE(subresourcePaths[T](name)...)
}

// DieMain drops the intent for every subresource, in place, leaving what is
// written through the main resource.
func (b *DieMutBase[T, D, M]) DieMain() M { return core.Must(b.DieMainE()) }

// DieMainE is DieMain, returning an error rather than panicking.
func (b *DieMutBase[T, D, M]) DieMainE() (M, error) {
	return b.DieExceptE(mainPaths[T]()...)
}

// DieFeedApplyPatch feeds a server-side apply request as intent, over the
// die's value, which is kept: each field it has is set, null clears a field,
// [] and {} are an empty list or map, list-map elements are found by their
// keys, and an atomic struct is set whole. The intent replaces the die's, so
// DieReleaseApplyPatch sends the request back. Fields the die's type doesn't
// have aren't fed. An unstructured die follows its mold, if it has one;
// without, its lists are set whole. It panics if the request isn't a JSON
// object.
func (b *DieMutBase[T, D, M]) DieFeedApplyPatch(patch []byte) M {
	return core.Must(b.DieFeedApplyPatchE(patch))
}

// DieFeedApplyPatchE is DieFeedApplyPatch, returning an error rather than
// panicking.
func (b *DieMutBase[T, D, M]) DieFeedApplyPatchE(patch []byte) (M, error) {
	return b.feedPatch("DieFeedApplyPatch", core.ApplyPatch, patch)
}

// DieFeedMergePatch feeds a JSON merge patch, RFC 7386, as intent over the
// die's value, which is kept: null clears a field, objects and maps are
// merged field by field, and anything else, lists included, is set whole. The
// intent replaces the die's, so DieReleaseMergePatch sends the patch back,
// but for what the patch sets that the value already has. Fields the die's
// type doesn't have aren't fed. It panics if the patch isn't a JSON object.
func (b *DieMutBase[T, D, M]) DieFeedMergePatch(patch []byte) M {
	return core.Must(b.DieFeedMergePatchE(patch))
}

// DieFeedMergePatchE is DieFeedMergePatch, returning an error rather than
// panicking.
func (b *DieMutBase[T, D, M]) DieFeedMergePatchE(patch []byte) (M, error) {
	return b.feedPatch("DieFeedMergePatch", core.MergePatch, patch)
}

func (b *DieMutBase[T, D, M]) feedPatch(method string, kind core.PatchKind, patch []byte) (M, error) {
	var zero M
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
func (b *DieMutBase[T, D, M]) Equal(o M) bool { return core.Must(b.EqualE(o)) }

// EqualE is Equal, returning an error rather than panicking.
func (b *DieMutBase[T, D, M]) EqualE(o M) (bool, error) { return equal[T](b, o) }

// UnmarshalJSON replaces the base and keeps the intent, so the intent is
// reapplied on top of the decoded value.
func (b *DieMutBase[T, D, M]) UnmarshalJSON(data []byte) error {
	h := b.hold()
	if len(b.at) > 0 || h.View != nil {
		return core.Errorf("die: can only unmarshal into a root die")
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
// with ErrImmutable, ErrImmutableChild or ErrReadOnlyJig for those.
//
// Methods that implement Kubernetes' and Go's interfaces, such as
// metav1.Object's setters, can't return an error: they panic with one, such
// as on an immutable die.
type Error = core.Error

var (
	// ErrImmutable is the error of a write to an immutable die through an
	// interface it implements, such as metav1.Object's setters: write
	// within DieWith, or to a DieMutable copy.
	ErrImmutable error = core.ErrImmutable
	// ErrImmutableChild is the error of a write through a child of an
	// immutable die, which would have nowhere to go: write within DieWith.
	ErrImmutableChild error = core.ErrImmutableChild
	// ErrReadOnlyJig is the error of a write through a jig of an immutable
	// die.
	ErrReadOnlyJig error = core.ErrReadOnlyJig
)
