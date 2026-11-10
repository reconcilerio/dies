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
	"context"
	"fmt"
	"reflect"
	"slices"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"reconciler.io/dies/v2/die/internal/core"
)

// Validator is implemented by an [ImmutableDie] that validates its value,
// such as with a hand-written method on a generated die:
//
//	func (d *WidgetPortDie) DieValidate(ctx context.Context, op die.Operation) (field.ErrorList, error) {
//		if old, ok := die.OldOf(d, op); ok && d.Port().Get() < old.Port().Get() {
//			return field.ErrorList{validation.Invalid(d.Port(), "may not decrease")}, nil
//		}
//		return nil, nil
//	}
//
// [ValidateCreate], [ValidateUpdate] and [ValidateDelete] call it for each die
// they walk to, so a die validates only its own fields: those that are dies
// validate themselves, wherever they are, and the dies between need no
// method. The die knows its location from the root, so the helpers of
// [reconciler.io/dies/v2/die/validation] make errors at a field from its
// handle. The operation says what is validated, and, for an update, holds
// the die's old value at its location, which [OldOf] returns.
//
// The errors it returns are the value's. It returns an error when it can't
// validate the value, such as when a lookup it needs fails, which stops
// validation.
type Validator interface {
	ImmutableDie
	DieValidate(ctx context.Context, op Operation) (field.ErrorList, error)
}

// Defaulter is implemented by a [MutableDie] that defaults its value, as a
// [Validator] validates it. [Default] calls it for each die it walks to.
type Defaulter interface {
	MutableDie
	DieDefault(ctx context.Context) error
}

// OperationType is what a validation is of.
type OperationType string

const (
	// Create is the validation of a value to be created.
	Create OperationType = "Create"
	// Update is the validation of a value to replace an old one: an update,
	// whether by an admission webhook of the update, or before sending it.
	Update OperationType = "Update"
	// Delete is the validation of a value to be deleted.
	Delete OperationType = "Delete"
)

// Operation is what a [Validator] is called for: the type of operation, and
// for an update, the die's old value.
type Operation struct {
	Type OperationType
	// old is the old die at the die's location, for an update where the old
	// value has it
	old ImmutableDie
}

// Old returns the old die at the validated die's location, for an update
// where the old value has it, such as an element that was in the old list
// map with the same key. Use [OldOf] for it as the die's type.
func (op Operation) Old() (ImmutableDie, bool) { return op.old, op.old != nil }

// OldOf returns the old value of a die being validated, as the die's type:
// see [Operation.Old].
func OldOf[D ImmutableDie](d D, op Operation) (D, bool) {
	old, ok := op.old.(D)
	return old, ok
}

// ValidateCreate validates a die to be created, and every die within it:
// each struct field held by value, or by pointer when present, each element
// of a list map, and each item of an atomic list. For each die, from the die
// down, it returns the errors of its schema, then its own if it is a
// [Validator]. Errors are located from the root, so validating a die within
// another reports the paths the root has.
//
// The schema's errors are those of its markers, which diegen reads, as the
// API server reports them:
//
//   - a required field that is absent: one marked +required or
//     +kubebuilder:validation:Required, or, in a package of custom resources
//     by kubebuilder's convention, one without omitempty or omitzero that
//     isn't +optional
//   - a field marked +k8s:required that is absent, or empty unless held by
//     pointer, as is a list-map element's key without a default
//   - a union, of +oneOf, +k8s:unionMember or a +die:union with oneOf=true,
//     without exactly one member set, or of +zeroOrOneOf or a +die:union
//     with more than one
//   - a present value that fails its checks, of declarative validation's
//     markers, such as +k8s:format=k8s-short-name, +k8s:minimum=0 or
//     +k8s:maxItems=8, of kubebuilder's, such as
//     +kubebuilder:validation:Pattern, or of its type's, such as +enum, whose
//     values are its type's constants, or a +listType=set with an item twice
//
// CEL rules, such as +kubebuilder:validation:XValidation, are not checked.
// A list-map element is at its index, as the API server reports it.
//
// It returns an error if a [Validator] does, which stops validation, with
// the location of the die that returned it, and the errors of the value
// found until then.
func ValidateCreate(ctx context.Context, d ImmutableDie) (field.ErrorList, error) {
	return validateAs(ctx, d, nil, Create)
}

// ValidateUpdate validates a die to replace an old one, of the same type, as
// [ValidateCreate] validates a die to be created, such as in an admission
// webhook of the update, or before sending it. Each die walked to is paired
// with the old die at its location, if the old value has it: a field by its
// name, a list-map element by its key and an atomic list's item by its
// index. A [Validator] gets it from the [Operation].
//
// As the API server does, the schema's checks ratchet: a field that is
// unchanged from the old value isn't checked, so that values valid when they
// were written stay valid. The markers of an update are checked too:
// +k8s:immutable, and +k8s:update's NoSet, NoUnset, NoModify, NoAddItem and
// NoRemoveItem.
func ValidateUpdate[D ImmutableDie](ctx context.Context, d, old D) (field.ErrorList, error) {
	return validateAs(ctx, d, old, Update)
}

// ValidateDelete validates a die to be deleted. Only each die that is a
// [Validator] validates it: a value being deleted has no schema to meet.
func ValidateDelete(ctx context.Context, d ImmutableDie) (field.ErrorList, error) {
	return validateAs(ctx, d, nil, Delete)
}

func validateAs(ctx context.Context, d, old ImmutableDie, op OperationType) (field.ErrorList, error) {
	h, at, t := stateOf(d)
	n := node{h: h, at: at, t: t, fp: core.FieldPath(h, at), d: d}
	if old != nil {
		oh, oat, ot := stateOf(old)
		if ot != t {
			panic(core.Errorf("die: validating an update of %s from a %s", t, ot))
		}
		n.old = &node{h: oh, at: oat, t: ot, d: old}
	}
	var errs field.ErrorList
	err := walk(n, false, func(n node) error {
		vop := Operation{Type: op}
		if n.old != nil {
			vop.old = n.old.d.(ImmutableDie)
		}
		if op != Delete {
			errs = append(errs, schemaErrors(ctx, n, op)...)
		}
		if v, ok := n.d.(Validator); ok {
			verrs, err := v.DieValidate(ctx, vop)
			errs = append(errs, verrs...)
			if err != nil {
				return fmt.Errorf("validating %s: %w", n.at, err)
			}
		}
		return nil
	})
	return errs, err
}

// Default defaults a die and every die within it, as [ValidateCreate] walks
// them: from the die down, it sets each absent field with a +default or
// +kubebuilder:default to it, as the API server does for a custom resource,
// then calls the die if it is a [Defaulter], so what either creates is then
// walked too. A die's defaults are its own: those of a field that is a die
// are its die's, wherever it is. It stops at the first error, with the
// location of the die that returned it.
//
// A die that changes the keys of the list-map element it is should not have
// children to default: they are walked at the key the element had.
func Default(ctx context.Context, m MutableDie) error {
	h, at, t := stateOf(m)
	return walk(node{h: h, at: at, t: t, fp: core.FieldPath(h, at), d: m}, true, func(n node) error {
		if err := applyDefaults(n); err != nil {
			return fmt.Errorf("defaulting %s: %w", n.at, err)
		}
		if d, ok := n.d.(Defaulter); ok {
			if err := d.DieDefault(ctx); err != nil {
				return fmt.Errorf("defaulting %s: %w", n.at, err)
			}
		}
		return nil
	})
}

func stateOf(d Die) (*core.Holder, Path, reflect.Type) {
	h, at, ok := core.StateOf(d)
	if !ok {
		panic(core.Errorf("die: %T is not a die", d))
	}
	return h, at, core.Deref(core.ValueTypeOf(d))
}

// node is a die that is walked to, with what its schema checks need, and
// the old die at its location, for an update where the old value has it.
type node struct {
	h  *core.Holder
	at Path
	t  reflect.Type
	fp *field.Path
	// keys are the key fields of a list-map element
	keys []string
	d    any
	old  *node
}

// child returns the node of a die within this one, at a step from it, paired
// with the old die at the same step, if the old value has it.
func (n node) child(fd *core.FieldDesc, steps []core.Step, t reflect.Type, fp *field.Path, keys []string, mutable bool, inOld func(old node) bool) node {
	p := core.JoinPath(n.at, steps)
	c := node{h: n.h, at: p, t: t, fp: fp, keys: keys, d: fd.Dies(n.h, p, mutable)}
	if n.old != nil && inOld(*n.old) {
		op := core.JoinPath(n.old.at, steps)
		c.old = &node{h: n.old.h, at: op, t: t, keys: keys, d: fd.Dies(n.old.h, op, false)}
	}
	return c
}

// walk calls fn for a die, then walks the dies within it, found after fn so
// that it sees what fn wrote.
func walk(n node, mutable bool, fn func(node) error) error {
	if err := fn(n); err != nil {
		return err
	}
	for _, fd := range core.ManagedFields(n.t) {
		if fd.Dies == nil {
			continue
		}
		field := Path{core.FieldStep(fd.Wire)}
		fp := n.fp.Child(fd.Wire)
		switch fd.Kind {
		case core.KindStruct:
			byPointer := fd.Typ.Kind() == reflect.Pointer
			if byPointer && !core.Must(n.h.Has(core.JoinPath(n.at, field))) {
				continue
			}
			c := n.child(fd, field, core.Deref(fd.Typ), fp, nil, mutable, func(old node) bool {
				return !byPointer || core.Must(old.h.Has(core.JoinPath(old.at, field)))
			})
			if err := walk(c, mutable, fn); err != nil {
				return err
			}
		case core.KindLeaf:
			// an atomic list of items with dies, at their indexes
			list := core.Must(n.h.ValueAt(core.JoinPath(n.at, field)))
			for i := range list.Len() {
				c := n.child(fd, Path{field[0], core.IndexStep(i)}, fd.Typ.Elem(), fp.Index(i), nil, mutable, func(old node) bool {
					return i < core.Must(old.h.ValueAt(core.JoinPath(old.at, field))).Len()
				})
				if err := walk(c, mutable, fn); err != nil {
					return err
				}
			}
		case core.KindListMap:
			list := core.Must(n.h.ValueAt(core.JoinPath(n.at, field)))
			elem := fd.Typ.Elem()
			steps := make([]core.Step, list.Len())
			for i := range list.Len() {
				steps[i] = core.KeyStepOf(elem, fd.Keys, list.Index(i))
			}
			for i, s := range steps {
				c := n.child(fd, Path{field[0], s}, elem, fp.Index(i), fd.Keys, mutable, func(old node) bool {
					l := core.Must(old.h.ValueAt(core.JoinPath(old.at, field)))
					return core.Must(core.FindElem(l, s)) >= 0
				})
				if err := walk(c, mutable, fn); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// value returns a die's struct value, or false for a nil pointer or a value
// that isn't a struct.
func (n node) value() (reflect.Value, bool) {
	if n.t.Kind() != reflect.Struct {
		return reflect.Value{}, false
	}
	v := core.Must(n.h.ValueAt(n.at))
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return reflect.Value{}, false
		}
		v = v.Elem()
	}
	return v, true
}

func (n node) has(wire string) bool {
	return core.Must(n.h.Has(core.ChildPath(n.at, core.FieldStep(wire))))
}

// schemaErrors are the errors of a die's schema: see [ValidateCreate] and
// [ValidateUpdate].
func schemaErrors(ctx context.Context, n node, op OperationType) field.ErrorList {
	// only a struct has a schema; an unstructured value has none
	v, ok := n.value()
	if !ok {
		return nil
	}
	var old reflect.Value
	if n.old != nil {
		if old, ok = n.old.value(); !ok {
			n.old = nil
		}
	}
	var errs field.ErrorList
	var unions []string
	members := map[string][]string{}
	oneOf := map[string]bool{}
	for _, fd := range core.Describe(n.t).Fields {
		fp := n.fp.Child(fd.Wire)
		present, value := n.has(fd.Wire), fd.ValueIn(v)
		checked := true
		if n.old != nil {
			oldPresent, oldValue := n.old.has(fd.Wire), fd.ValueIn(old)
			changed := present != oldPresent || !reflect.DeepEqual(value.Interface(), oldValue.Interface())
			if fd.Checks != nil && changed {
				errs = append(errs, checkUpdate(fd, fp, present, value, oldPresent, oldValue)...)
			}
			// ratcheting: a field unchanged from the old value isn't checked
			checked = changed
		}
		if checked {
			// a list-map key without a default identifies its element, so it
			// is required not to be empty
			nonZero := fd.RequiredNonZero || slices.Contains(n.keys, fd.Wire) && !fd.HasDefault()
			switch {
			case (fd.Required || nonZero) && !present, nonZero && empty(value):
				errs = append(errs, field.Required(fp, ""))
			case present && fd.Checks != nil:
				errs = append(errs, checkValue(ctx, fp, value, fd.Checks)...)
			}
		}
		if fd.Union != "" {
			if _, ok := members[fd.Union]; !ok {
				unions = append(unions, fd.Union)
			}
			members[fd.Union] = append(members[fd.Union], fd.Wire)
			oneOf[fd.Union] = oneOf[fd.Union] || fd.OneOf
		}
	}
	for _, u := range unions {
		errs = append(errs, checkUnion(ctx, n.fp, members[u], oneOf[u], n.has, n.old)...)
	}
	return errs
}

// applyDefaults sets each absent field of a die that has a +default to it,
// as the API server does for a custom resource. A default of a union's
// member doesn't clear the others: validation reports both.
func applyDefaults(n node) error {
	if n.t.Kind() != reflect.Struct {
		return nil
	}
	for _, fd := range core.Describe(n.t).Fields {
		if !fd.HasDefault() {
			continue
		}
		rel := Path{core.FieldStep(fd.Wire)}
		if core.Must(n.h.Has(core.JoinPath(n.at, rel))) {
			continue
		}
		if _, err := n.h.Write(n.at, rel, core.SetOp(core.OpSet, fd.Default()), false); err != nil {
			return err
		}
	}
	return nil
}

// empty reports whether a required field's value is missing though present:
// an empty string, list or map, or the zero value of a scalar. A pointer is
// set when it isn't nil, whatever it points to, and a struct when present.
func empty(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Pointer, reflect.Interface:
		return v.IsNil()
	case reflect.String, reflect.Slice, reflect.Map:
		return v.Len() == 0
	case reflect.Struct:
		return false
	}
	return v.IsZero()
}
