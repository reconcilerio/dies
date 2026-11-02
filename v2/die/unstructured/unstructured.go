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

// Package unstructured has dies for objects with no Go type, held as
// apimachinery's unstructured objects hold them, and nodes to traverse them.
package unstructured

import (
	"reflect"
	"slices"

	"reconciler.io/dies/v2/die/base"
	"reconciler.io/dies/v2/die/internal/core"
	"reconciler.io/dies/v2/die/internal/dies"
	"reconciler.io/dies/v2/die/lists"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// Object is a die for an object with no Go type, held as
// unstructured.Unstructured holds one. It behaves as any die: immutable or
// mutable, with intent recorded and patches from it. It is an object, a
// runtime.Object and a metav1.Object, but has no field accessors: its values
// are read and written through Nodes, from Node.
type Object struct {
	base.DieObjectBase[map[string]any, *Object]
}

func (*Object) DieFrom(b base.DieBase[map[string]any, *Object]) *Object {
	return &Object{base.AsObject(b)}
}

// Blank is an empty, immutable unstructured die.
var Blank = base.Blank[map[string]any, *Object]()

// BlankFor returns an immutable, empty unstructured die of a kind: a blank,
// with only its type metadata, and no mold. Molds makes blanks in a mold.
func BlankFor(apiVersion, kind string) *Object {
	return Blank.DieWith(func(d *Object) {
		d.Node().Field("apiVersion").Set(apiVersion)
		d.Node().Field("kind").Set(kind)
	}).DieCommit()
}

// DieFeedUnstructured feeds an unstructured object. Unlike other object dies,
// which decode it, an unstructured die copies its content, so its numbers
// stay int64 where they are.
func (d *Object) DieFeedUnstructured(u *unstructured.Unstructured) *Object {
	return d.DieFeed(u.DeepCopy().Object)
}

// DieReleaseUnstructured releases the object as an unstructured object.
func (d *Object) DieReleaseUnstructured() *unstructured.Unstructured {
	return &unstructured.Unstructured{Object: d.DieRelease()}
}

// An unstructured die is a runtime.Unstructured, so clients and libraries that
// handle unstructured objects, such as controller-runtime's client, take it
// directly.
var _ runtime.Unstructured = (*Object)(nil)

// NewEmptyInstance returns a mutable, empty die of the same kind.
func (d *Object) NewEmptyInstance() runtime.Unstructured {
	empty := &Object{}
	empty.GetObjectKind().SetGroupVersionKind(d.GetObjectKind().GroupVersionKind())
	// of the same kind, with the same schema, if it has one
	h, _ := dies.MustStateOf(empty)
	src, _ := dies.MustStateOf(d)
	h.Schema = src.Schema
	return empty
}

// UnstructuredContent returns a copy of the object's content. Unlike an
// unstructured.Unstructured's, changing it does not change the die; write
// through the die instead.
func (d *Object) UnstructuredContent() map[string]any { return d.DieRelease() }

// SetUnstructuredContent replaces the base of a mutable die and keeps its
// intent, as decoding into it does. On a die within another, such as an item
// of a list, it sets the content as intent.
func (d *Object) SetUnstructuredContent(content map[string]any) {
	h, at := dies.MustStateOf(d)
	core.Check(setUnstructuredContent(h, at, content))
}

func setUnstructuredContent(h *core.Holder, at base.Path, content map[string]any) error {
	if !h.Mutable {
		return core.Errorf("die: SetUnstructuredContent on an immutable die, which would lose the change; use a mutable die, such as within DieWith")
	}
	content = runtime.DeepCopyJSON(content)
	if len(at) > 0 || h.View != nil {
		_, err := h.Write(at, nil, core.SetOp(core.OpSet, reflect.ValueOf(content)), false)
		return err
	}
	h.Base = reflect.ValueOf(content)
	return nil
}

// IsList reports whether the content is a list, with items, as an
// unstructured.Unstructured does.
func (d *Object) IsList() bool {
	_, ok := d.Node().Field("items").Get().([]any)
	return ok
}

// ToList converts the die to a list, when its content is a list. The list is
// a copy, whose items are mutable dies.
func (d *Object) ToList() (*ObjectList, bool) {
	if !d.IsList() {
		return nil, false
	}
	l := &ObjectList{}
	l.SetUnstructuredContent(d.DieRelease())
	return l, true
}

// EachListItem calls fn with each item of a list, as an unstructured die
// within the list: like an unstructured.Unstructured's items, reads and
// writes go to the list, so a change to an item changes the list. An item is
// located by its index, and a change is recorded as intent on the whole list.
// The items of an immutable list are read-only. It fails if the content is
// not a list, or an item is not an object.
func (d *Object) EachListItem(fn func(runtime.Object) error) error {
	h, at := dies.MustStateOf(d)
	return eachListItem(h, at, false, fn)
}

// EachListItemWithAlloc calls fn with each item of a list, as a mutable
// copy: changing it does not change the list.
func (d *Object) EachListItemWithAlloc(fn func(runtime.Object) error) error {
	h, at := dies.MustStateOf(d)
	return eachListItem(h, at, true, fn)
}

func eachListItem(h *core.Holder, at base.Path, alloc bool, fn func(runtime.Object) error) error {
	v, err := h.ValueAt(core.ChildPath(at, core.FieldStep("items")))
	if err != nil {
		return err
	}
	items, ok := v.Interface().([]any)
	if !ok {
		return core.Errorf("die: the content of an unstructured die is not a list")
	}
	for i, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			return core.Errorf("die: an item of an unstructured die is not an object")
		}
		var obj *Object
		if alloc {
			var err error
			if obj, err = Blank.DieFeedE(m); err != nil {
				return err
			}
			obj = obj.DieImmutable(false)
		} else {
			obj = listItem(h, at, i)
		}
		if err := fn(obj); err != nil {
			return err
		}
	}
	return nil
}

// listItem is the item of a list at an index, as a die within the list.
func listItem(h *core.Holder, at base.Path, i int) *Object {
	return dies.WrapDie[map[string]any, *Object](h, core.ChildPath(core.ChildPath(at, core.FieldStep("items")), core.IndexStep(i)))
}

// Node returns the node for the whole of this die's object, from which its
// fields are reached: d.Node().Field("spec"). The die is the object, with its
// lifecycle, metadata and patches; nodes read and write its values.
func (d *Object) Node() *Node {
	h, at := dies.MustStateOf(d)
	return &Node{h: h, at: at}
}

// ObjectList is a list of unstructured dies, as
// unstructured.UnstructuredList is a list of unstructured objects. Like every
// list die, it is a plain list: see lists.List. It is also a runtime.Unstructured,
// so clients list into it, when it is registered as a kind's list type for
// those that need it. Fields of the list other than its type and list
// metadata and items are not kept.
type ObjectList struct {
	lists.List[*Object]
}

var _ runtime.Unstructured = (*ObjectList)(nil)
var _ metav1.ListInterface = (*ObjectList)(nil)

// NewList returns an unstructured list of a kind, with items.
func NewList(apiVersion, kind string, items ...*Object) *ObjectList {
	return &ObjectList{lists.NewList(apiVersion, kind, items...)}
}

// ListFrom returns a list of an unstructured list's content. Its
// items are mutable, as decoded items are.
func ListFrom(u *unstructured.UnstructuredList) *ObjectList { return core.Must(ListFromE(u)) }

// ListFromE is ListFrom, returning an error rather than panicking.
func ListFromE(u *unstructured.UnstructuredList) (*ObjectList, error) {
	l := &ObjectList{}
	if err := l.setContent(u.UnstructuredContent()); err != nil {
		return nil, err
	}
	return l, nil
}

// UnstructuredList returns the list as an unstructured list.
func (l *ObjectList) UnstructuredList() *unstructured.UnstructuredList {
	u := &unstructured.UnstructuredList{}
	u.SetUnstructuredContent(l.UnstructuredContent())
	return u
}

// DeepCopy returns a copy of the list, with its items copied.
func (l *ObjectList) DeepCopy() *ObjectList {
	return &ObjectList{l.DieCopy()}
}

func (l *ObjectList) DeepCopyObject() runtime.Object { return l.DeepCopy() }

// NewEmptyInstance returns an empty list of the same kind.
func (l *ObjectList) NewEmptyInstance() runtime.Unstructured {
	return &ObjectList{lists.List[*Object]{TypeMeta: l.TypeMeta}}
}

// UnstructuredContent returns the list's content, items included.
func (l *ObjectList) UnstructuredContent() map[string]any {
	content, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&struct {
		metav1.TypeMeta `json:",inline"`
		metav1.ListMeta `json:"metadata,omitempty"`
	}{l.TypeMeta, l.ListMeta})
	if err != nil {
		panic(core.Wrap(err))
	}
	items := make([]any, len(l.Items))
	for i, item := range l.Items {
		items[i] = item.UnstructuredContent()
	}
	content["items"] = items
	return content
}

// SetUnstructuredContent replaces the list with content, as decoding into it
// does. Its items are mutable, as decoded items are.
func (l *ObjectList) SetUnstructuredContent(content map[string]any) {
	core.Check(l.setContent(content))
}

func (l *ObjectList) setContent(content map[string]any) error {
	var meta struct {
		metav1.TypeMeta `json:",inline"`
		metav1.ListMeta `json:"metadata,omitempty"`
	}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(content, &meta); err != nil {
		return core.Errorf("die: the content of an unstructured list: %w", err)
	}
	l.TypeMeta, l.ListMeta = meta.TypeMeta, meta.ListMeta
	raw, _ := content["items"].([]any)
	l.Items = make([]*Object, 0, len(raw))
	for _, item := range raw {
		if m, ok := item.(map[string]any); ok {
			obj, err := Blank.DieFeedE(runtime.DeepCopyJSON(m))
			if err != nil {
				return err
			}
			l.Items = append(l.Items, obj.DieImmutable(false))
		}
	}
	return nil
}

// IsList is true: a list die is a list.
func (l *ObjectList) IsList() bool { return true }

// EachListItem calls fn with each item.
func (l *ObjectList) EachListItem(fn func(runtime.Object) error) error {
	for _, item := range l.Items {
		if err := fn(item); err != nil {
			return err
		}
	}
	return nil
}

// EachListItemWithAlloc calls fn with a copy of each item.
func (l *ObjectList) EachListItemWithAlloc(fn func(runtime.Object) error) error {
	for _, item := range l.Items {
		if err := fn(item.DeepCopy()); err != nil {
			return err
		}
	}
	return nil
}

// Node is a location within an unstructured die: an object, a
// list, or a value. Like a child die, it is attached to its root, reads
// never create anything, and the first write creates what its path needs.
// Writes through a node of an immutable die fail; write within DieWith.
type Node struct {
	h  *core.Holder
	at base.Path
}

// Field returns the node for a field of this object.
func (n *Node) Field(name string) *Node {
	return &Node{h: n.h, at: core.ChildPath(n.at, core.FieldStep(name))}
}

// Item returns the node for the element of this list with key fields, given
// as name and value pairs: Item("name", "app"). A list without keys of its
// own is addressed by the key fields its elements are known by. A key that
// is not unique is ambiguous.
func (n *Node) Item(keysAndValues ...any) *Node { return core.Must(n.ItemE(keysAndValues...)) }

// ItemE is Item, returning an error rather than panicking.
func (n *Node) ItemE(keysAndValues ...any) (*Node, error) {
	if len(keysAndValues) == 0 || len(keysAndValues)%2 != 0 {
		return nil, core.Errorf("die: Item takes key names and values, in pairs")
	}
	var keys []string
	var vals []any
	for i := 0; i < len(keysAndValues); i += 2 {
		k, ok := keysAndValues[i].(string)
		if !ok {
			return nil, core.Errorf("die: Item key %v is not a string", keysAndValues[i])
		}
		v, err := core.ToUnstructured(reflect.ValueOf(keysAndValues[i+1]))
		if err != nil {
			return nil, err
		}
		keys = append(keys, k)
		vals = append(vals, v)
	}
	return &Node{h: n.h, at: core.ChildPath(n.at, core.UnstructuredKeyStep(keys, vals))}, nil
}

// Index returns the node for the element of this list at an index. Unlike
// Item, it is for lists whose elements have no keys: the element is found by
// its position in the list as it is, and a write through it sets the whole
// list. Clearing the element removes it from the list. Writing an element
// past the end of the list is an error.
func (n *Node) Index(i int) *Node { return core.Must(n.IndexE(i)) }

// IndexE is Index, returning an error rather than panicking.
func (n *Node) IndexE(i int) (*Node, error) {
	if i < 0 {
		return nil, core.Errorf("die: index %d is negative", i)
	}
	return &Node{h: n.h, at: core.ChildPath(n.at, core.IndexStep(i))}, nil
}

// Find returns the nodes that a JSONPath expression matches, as the
// jsonpath package resolves them: $ is the die's root, and an expression that
// doesn't start with it starts at this node; an element of a list is found by
// the field a filter matched it by when it matched one element by one field,
// otherwise by its index, and a definite path finds a node whether or not it
// has a value, so that setting it creates it.
func (n *Node) Find(expr string) []*Node { return core.Must(n.FindE(expr)) }

// FindE is Find, returning an error rather than panicking.
func (n *Node) FindE(expr string) ([]*Node, error) {
	paths, err := core.ResolveJSONPath(n.h, n.at, expr)
	if err != nil {
		return nil, err
	}
	nodes := make([]*Node, len(paths))
	for i, p := range paths {
		nodes[i] = &Node{h: n.h, at: p}
	}
	return nodes, nil
}

func (n *Node) DiePath() base.Path { return n.at }

// DieRoot returns the root die this is within: see base.DieBase.DieRoot.
func (n *Node) DieRoot() base.Die { return base.RootOf(n.h) }

// DieRootType returns the type of the root value this is within.
func (n *Node) DieRootType() reflect.Type { return n.h.Base.Type() }

// Get returns a copy of the value, in unstructured form, or nil.
func (n *Node) Get() any { return core.Must(n.GetE()) }

// GetE is Get, returning an error rather than panicking.
func (n *Node) GetE() (any, error) {
	v, err := n.h.ValueAt(n.at)
	if err != nil || !v.IsValid() || (v.Kind() == reflect.Interface && v.IsNil()) {
		return nil, err
	}
	return v.Interface(), nil
}

// Location is a location in an unstructured die: the die itself, or a
// node within it.
type Location interface {
	unstructuredAt() (*core.Holder, base.Path)
}

func (d *Object) unstructuredAt() (*core.Holder, base.Path) { return dies.MustStateOf(d) }

func (n *Node) unstructuredAt() (*core.Holder, base.Path) { return n.h, n.at }

// GetAs reads the value of an unstructured die or node as a T. A die is read
// as a detached copy; to keep it linked, use ViewAs.
func GetAs[T any](u Location) T { return core.Must(GetAsE[T](u)) }

// GetAsE is GetAs, returning an error rather than panicking.
func GetAsE[T any](u Location) (T, error) {
	h, at := u.unstructuredAt()
	return core.GetAt[T](h, at, nil)
}

// JigAs presents an unstructured die or node as a typed die, D, such as a
// *DeploymentDie or a *ContainerDie. Like any jig, it holds no state of its
// own: reads come from the unstructured die, and writes go to it as intent,
// changing only the fields written, so what D's type does not describe, such
// as fields newer than it, is kept. A jig of an immutable die is read-only.
func JigAs[D interface{ DieMutable() bool }](u Location) D {
	h, at := u.unstructuredAt()
	d := reflect.New(reflect.TypeFor[D]().Elem()).Interface().(D)
	t := core.ValueTypeOf(d)
	src := core.DieSource{H: h}
	core.SetState(d, &core.Holder{Base: reflect.Zero(t), Mutable: src.Mutable(), View: core.WholeView(src, at, t)}, nil)
	return d
}

// Has reports whether the value is present.
func (n *Node) Has() bool { return core.Must(n.HasE()) }

// HasE is Has, returning an error rather than panicking.
func (n *Node) HasE() (bool, error) { return n.h.Has(n.at) }

// Set records the value as intent. Typed values, such as structs, are
// converted to unstructured form.
func (n *Node) Set(v any) *Node { return core.Must(n.SetE(v)) }

// SetE is Set, returning an error rather than panicking.
func (n *Node) SetE(v any) (*Node, error) {
	if _, err := n.h.Write(n.at, nil, core.SetOp(core.OpSet, reflect.ValueOf(&v).Elem()), false); err != nil {
		return nil, err
	}
	return n, nil
}

// Clear records that the value is absent: null in patches, or a removed
// list element.
func (n *Node) Clear() *Node { return core.Must(n.ClearE()) }

// ClearE is Clear, returning an error rather than panicking.
func (n *Node) ClearE() (*Node, error) {
	if _, err := n.h.Write(n.at, nil, core.SetOp(core.OpClear, reflect.Value{}), false); err != nil {
		return nil, err
	}
	return n, nil
}

// EmptyObject records that the value is present and an empty object.
func (n *Node) EmptyObject() *Node { return core.Must(n.EmptyObjectE()) }

// EmptyObjectE is EmptyObject, returning an error rather than panicking.
func (n *Node) EmptyObjectE() (*Node, error) {
	if _, err := n.h.Write(n.at, nil, core.SetOp(core.OpEmpty, reflect.ValueOf(map[string]any{})), false); err != nil {
		return nil, err
	}
	return n, nil
}

// EmptyList records that the value is present and an empty list.
func (n *Node) EmptyList() *Node { return core.Must(n.EmptyListE()) }

// EmptyListE is EmptyList, returning an error rather than panicking.
func (n *Node) EmptyListE() (*Node, error) {
	if _, err := n.h.Write(n.at, nil, core.SetOp(core.OpEmpty, reflect.ValueOf([]any{})), false); err != nil {
		return nil, err
	}
	return n, nil
}

// Keys returns the fields of this object, sorted.
func (n *Node) Keys() []string { return core.Must(n.KeysE()) }

// KeysE is Keys, returning an error rather than panicking.
func (n *Node) KeysE() ([]string, error) {
	v, err := n.GetE()
	if err != nil {
		return nil, err
	}
	m, _ := v.(map[string]any)
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys, nil
}

// Len returns the length of this list, or the number of fields of this
// object.
func (n *Node) Len() int { return core.Must(n.LenE()) }

// LenE is Len, returning an error rather than panicking.
func (n *Node) LenE() (int, error) {
	v, err := n.GetE()
	switch v := v.(type) {
	case []any:
		return len(v), err
	case map[string]any:
		return len(v), err
	}
	return 0, err
}

// Type is the JSON type of an unstructured value. Each is held
// as one Go type, which Get returns and GetAs reads without conversion; GetAs
// also reads a value as any other type with the same JSON form, such as an
// object as a struct.
type Type uint8

const (
	// TypeAbsent is a value that is missing or null: nil.
	TypeAbsent Type = iota
	// TypeObject is held as a map[string]any.
	TypeObject
	// TypeList is held as a []any.
	TypeList
	// TypeString is held as a string.
	TypeString
	// TypeBool is held as a bool.
	TypeBool
	// TypeInt is held as an int64.
	TypeInt
	// TypeFloat is held as a float64.
	TypeFloat
)

func (t Type) String() string {
	switch t {
	case TypeObject:
		return "object"
	case TypeList:
		return "list"
	case TypeString:
		return "string"
	case TypeBool:
		return "bool"
	case TypeInt:
		return "int"
	case TypeFloat:
		return "float"
	}
	return "absent"
}

// Type returns the JSON type of this value. Unstructured objects hold
// integers as int64 and other numbers as float64, but a value that was fed in
// another Go form is classified by its kind.
func (n *Node) Type() Type { return core.Must(n.TypeE()) }

// TypeE is Type, returning an error rather than panicking.
func (n *Node) TypeE() (Type, error) {
	v, err := n.GetE()
	if err != nil || v == nil {
		return TypeAbsent, err
	}
	switch reflect.ValueOf(v).Kind() {
	case reflect.Map:
		return TypeObject, nil
	case reflect.Slice, reflect.Array:
		return TypeList, nil
	case reflect.String:
		return TypeString, nil
	case reflect.Bool:
		return TypeBool, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return TypeInt, nil
	case reflect.Float32, reflect.Float64:
		return TypeFloat, nil
	}
	return TypeAbsent, core.Errorf("die: %s holds %T, which is not an unstructured value", n.at, v)
}
