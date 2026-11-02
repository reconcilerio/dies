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

// Object is an immutable die for an object with no Go type, held as
// unstructured.Unstructured holds one, and ObjectMut is its mutable die.
// They behave as any die, with intent recorded and patches from it. They are
// objects, runtime.Objects and metav1.Objects, but have no field accessors:
// their values are read, and an ObjectMut's written, through nodes, from
// Node.
type Object struct {
	base.DieObjectBase[map[string]any, *Object, *ObjectMut]
}

func (*Object) DieFrom(b base.DieBase[map[string]any, *Object, *ObjectMut]) *Object {
	return &Object{base.AsObject(b)}
}

// ObjectMut is the mutable die of an object with no Go type: see Object.
type ObjectMut struct {
	base.DieObjectMutBase[map[string]any, *Object, *ObjectMut]
}

func (*ObjectMut) DieFrom(b base.DieMutBase[map[string]any, *Object, *ObjectMut]) *ObjectMut {
	return &ObjectMut{base.AsObjectMut(b)}
}

var (
	_ base.DieOf[map[string]any, *Object, *ObjectMut]    = (*Object)(nil)
	_ base.DieMutOf[map[string]any, *Object, *ObjectMut] = (*ObjectMut)(nil)
	_ base.DieObject                                     = (*Object)(nil)
	_ base.DieObject                                     = (*ObjectMut)(nil)
)

// Blank is an empty, immutable unstructured die.
var Blank = base.Blank[map[string]any, *Object, *ObjectMut]()

// BlankFor returns an immutable, empty unstructured die of a kind: a blank,
// with only its type metadata, and no mold. Molds makes blanks in a mold.
func BlankFor(apiVersion, kind string) *Object {
	return Blank.DieWith(func(d *ObjectMut) {
		d.Node().Field("apiVersion").Set(apiVersion)
		d.Node().Field("kind").Set(kind)
	}).DieCommit()
}

// DieFeedUnstructured feeds an unstructured object. Unlike other object dies,
// which decode it, an unstructured die copies its content, so its numbers
// stay int64 where they are.
func (d *ObjectMut) DieFeedUnstructured(u *unstructured.Unstructured) *ObjectMut {
	return d.DieFeed(u.DeepCopy().Object)
}

// DieFeedUnstructuredE is DieFeedUnstructured, returning an error rather
// than panicking.
func (d *ObjectMut) DieFeedUnstructuredE(u *unstructured.Unstructured) (*ObjectMut, error) {
	return d.DieFeedE(u.DeepCopy().Object)
}

// DieReleaseUnstructured releases the object as an unstructured object.
func (d *Object) DieReleaseUnstructured() *unstructured.Unstructured {
	return &unstructured.Unstructured{Object: d.DieRelease()}
}

// DieReleaseUnstructured releases the object as an unstructured object.
func (d *ObjectMut) DieReleaseUnstructured() *unstructured.Unstructured {
	return &unstructured.Unstructured{Object: d.DieRelease()}
}

// An unstructured die is a runtime.Unstructured, so clients and libraries that
// handle unstructured objects, such as controller-runtime's client, take it
// directly. An immutable die's SetUnstructuredContent panics.
var (
	_ runtime.Unstructured = (*Object)(nil)
	_ runtime.Unstructured = (*ObjectMut)(nil)
)

// NewEmptyInstance returns a mutable, empty die of the same kind.
func (d *Object) NewEmptyInstance() runtime.Unstructured { return newEmptyInstance(d) }

// NewEmptyInstance returns a mutable, empty die of the same kind.
func (d *ObjectMut) NewEmptyInstance() runtime.Unstructured { return newEmptyInstance(d) }

func newEmptyInstance(d runtime.Object) *ObjectMut {
	empty := &ObjectMut{}
	empty.GetObjectKind().SetGroupVersionKind(d.GetObjectKind().GroupVersionKind())
	// of the same kind, with the same schema, if it has one
	h, _ := dies.MustStateOf(empty)
	src, _ := dies.MustStateOf(d)
	h.Schema = src.Schema
	return empty
}

// UnstructuredContent returns a copy of the object's content. Unlike an
// unstructured.Unstructured's, changing it does not change the die.
func (d *Object) UnstructuredContent() map[string]any { return d.DieRelease() }

// UnstructuredContent returns a copy of the object's content. Unlike an
// unstructured.Unstructured's, changing it does not change the die; write
// through the die instead.
func (d *ObjectMut) UnstructuredContent() map[string]any { return d.DieRelease() }

// SetUnstructuredContent panics with ErrImmutable: the die is immutable.
func (d *Object) SetUnstructuredContent(map[string]any) {
	panic(core.Errorf("die: SetUnstructuredContent: %w", core.ErrImmutable))
}

// SetUnstructuredContent replaces the base of the die and keeps its intent,
// as decoding into it does. On a die within another, such as an item of a
// list, it sets the content as intent.
func (d *ObjectMut) SetUnstructuredContent(content map[string]any) {
	h, at := dies.MustStateOf(d)
	core.Check(setUnstructuredContent(h, at, content))
}

func setUnstructuredContent(h *core.Holder, at base.Path, content map[string]any) error {
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
func (d *Object) IsList() bool { return isList(d) }

// IsList reports whether the content is a list, with items, as an
// unstructured.Unstructured does.
func (d *ObjectMut) IsList() bool { return isList(d) }

func isList(d Location) bool {
	h, at := d.unstructuredAt()
	items := (&node{h: h, at: at}).field("items")
	_, ok := items.Get().([]any)
	return ok
}

// ToList converts the die to a list, when its content is a list. The list is
// a copy, whose items are mutable dies.
func (d *Object) ToList() (*ObjectList, bool) { return toList(d) }

// ToList converts the die to a list, when its content is a list. The list is
// a copy, whose items are mutable dies.
func (d *ObjectMut) ToList() (*ObjectList, bool) { return toList(d) }

func toList(d interface {
	Location
	DieRelease() map[string]any
}) (*ObjectList, bool) {
	if !isList(d) {
		return nil, false
	}
	l := &ObjectList{}
	l.SetUnstructuredContent(d.DieRelease())
	return l, true
}

// EachListItem calls fn with each item of a list, as an immutable
// unstructured die within the list. It fails if the content is not a list,
// or an item is not an object.
func (d *Object) EachListItem(fn func(runtime.Object) error) error {
	h, at := dies.MustStateOf(d)
	return eachListItem(h, at, false, fn)
}

// EachListItem calls fn with each item of a list, as a mutable unstructured
// die within the list: like an unstructured.Unstructured's items, reads and
// writes go to the list, so a change to an item changes the list. An item is
// located by its index, and a change is recorded as intent on the whole list.
// It fails if the content is not a list, or an item is not an object.
func (d *ObjectMut) EachListItem(fn func(runtime.Object) error) error {
	h, at := dies.MustStateOf(d)
	return eachListItem(h, at, false, fn)
}

// EachListItemWithAlloc calls fn with each item of a list, as a mutable
// copy: changing it does not change the list.
func (d *Object) EachListItemWithAlloc(fn func(runtime.Object) error) error {
	h, at := dies.MustStateOf(d)
	return eachListItem(h, at, true, fn)
}

// EachListItemWithAlloc calls fn with each item of a list, as a mutable
// copy: changing it does not change the list.
func (d *ObjectMut) EachListItemWithAlloc(fn func(runtime.Object) error) error {
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
		var obj runtime.Object
		switch {
		case alloc:
			var err error
			if obj, err = Blank.DieMutable().DieFeedE(m); err != nil {
				return err
			}
		case h.Mutable:
			obj = dies.WrapDieMut[map[string]any, *Object](h, itemAt(at, i))
		default:
			obj = dies.WrapDie[map[string]any, *Object, *ObjectMut](h, itemAt(at, i))
		}
		if err := fn(obj); err != nil {
			return err
		}
	}
	return nil
}

// itemAt is the location of the item of a list at an index.
func itemAt(at base.Path, i int) base.Path {
	return core.ChildPath(core.ChildPath(at, core.FieldStep("items")), core.IndexStep(i))
}

// Node returns the node for the whole of this die's object, from which its
// fields are read: d.Node().Field("spec"). The die is the object, with its
// lifecycle, metadata and patches; nodes read its values.
func (d *Object) Node() *Node {
	h, at := dies.MustStateOf(d)
	return &Node{node{h: h, at: at}}
}

// Node returns the node for the whole of this die's object, from which its
// fields are reached: d.Node().Field("spec"). The die is the object, with its
// lifecycle, metadata and patches; nodes read and write its values.
func (d *ObjectMut) Node() *NodeMut {
	h, at := dies.MustStateOf(d)
	return &NodeMut{node{h: h, at: at}}
}

// ObjectList is a list of unstructured dies, as
// unstructured.UnstructuredList is a list of unstructured objects. Like every
// list die, it is a plain list: see lists.List. It is also a runtime.Unstructured,
// so clients list into it, when it is registered as a kind's list type for
// those that need it. Fields of the list other than its type and list
// metadata and items are not kept.
type ObjectList struct {
	lists.List[*ObjectMut]
}

var _ runtime.Unstructured = (*ObjectList)(nil)
var _ metav1.ListInterface = (*ObjectList)(nil)

// NewList returns an unstructured list of a kind, with items.
func NewList(apiVersion, kind string, items ...*ObjectMut) *ObjectList {
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
	return &ObjectList{lists.List[*ObjectMut]{TypeMeta: l.TypeMeta}}
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
	l.Items = make([]*ObjectMut, 0, len(raw))
	for _, item := range raw {
		if m, ok := item.(map[string]any); ok {
			obj, err := Blank.DieMutable().DieFeedE(runtime.DeepCopyJSON(m))
			if err != nil {
				return err
			}
			l.Items = append(l.Items, obj)
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

// node is a location within an unstructured die, with the methods Node and
// NodeMut both have, which read.
type node struct {
	h  *core.Holder
	at base.Path
}

// Node is a location within an immutable unstructured die, which reads it:
// an object, a list, or a value. Like a child die, it is attached to its
// root, and reads never create anything.
type Node struct {
	node
}

// NodeMut is a location within a mutable unstructured die, which also
// writes it. The first write creates what its path needs.
type NodeMut struct {
	node
}

// Field returns the node for a field of this object.
func (n *NodeMut) Field(name string) *NodeMut { return &NodeMut{n.field(name)} }

// Item returns the node for the element of this list with key fields: see
// Node.Item.
func (n *NodeMut) Item(keysAndValues ...any) *NodeMut {
	return core.Must(n.ItemE(keysAndValues...))
}

// ItemE is Item, returning an error rather than panicking.
func (n *NodeMut) ItemE(keysAndValues ...any) (*NodeMut, error) {
	m, err := n.item(keysAndValues)
	if err != nil {
		return nil, err
	}
	return &NodeMut{m}, nil
}

// Index returns the node for the element of this list at an index: see
// Node.Index.
func (n *NodeMut) Index(i int) *NodeMut { return core.Must(n.IndexE(i)) }

// IndexE is Index, returning an error rather than panicking.
func (n *NodeMut) IndexE(i int) (*NodeMut, error) {
	m, err := n.index(i)
	if err != nil {
		return nil, err
	}
	return &NodeMut{m}, nil
}

// Find returns the nodes that a JSONPath expression matches: see Node.Find.
func (n *NodeMut) Find(expr string) []*NodeMut { return core.Must(n.FindE(expr)) }

// FindE is Find, returning an error rather than panicking.
func (n *NodeMut) FindE(expr string) ([]*NodeMut, error) {
	found, err := n.find(expr)
	if err != nil {
		return nil, err
	}
	nodes := make([]*NodeMut, len(found))
	for i, f := range found {
		nodes[i] = &NodeMut{f}
	}
	return nodes, nil
}

// Field returns the node for a field of this object.
func (n *Node) Field(name string) *Node { return &Node{n.field(name)} }

func (n *node) field(name string) node {
	return node{h: n.h, at: core.ChildPath(n.at, core.FieldStep(name))}
}

// Item returns the node for the element of this list with key fields, given
// as name and value pairs: Item("name", "app"). A list without keys of its
// own is addressed by the key fields its elements are known by. A key that
// is not unique is ambiguous.
func (n *Node) Item(keysAndValues ...any) *Node { return core.Must(n.ItemE(keysAndValues...)) }

// ItemE is Item, returning an error rather than panicking.
func (n *Node) ItemE(keysAndValues ...any) (*Node, error) {
	m, err := n.item(keysAndValues)
	if err != nil {
		return nil, err
	}
	return &Node{m}, nil
}

func (n *node) item(keysAndValues []any) (node, error) {
	if len(keysAndValues) == 0 || len(keysAndValues)%2 != 0 {
		return node{}, core.Errorf("die: Item takes key names and values, in pairs")
	}
	var keys []string
	var vals []any
	for i := 0; i < len(keysAndValues); i += 2 {
		k, ok := keysAndValues[i].(string)
		if !ok {
			return node{}, core.Errorf("die: Item key %v is not a string", keysAndValues[i])
		}
		v, err := core.ToUnstructured(reflect.ValueOf(keysAndValues[i+1]))
		if err != nil {
			return node{}, err
		}
		keys = append(keys, k)
		vals = append(vals, v)
	}
	return node{h: n.h, at: core.ChildPath(n.at, core.UnstructuredKeyStep(keys, vals))}, nil
}

// Index returns the node for the element of this list at an index. Unlike
// Item, it is for lists whose elements have no keys: the element is found by
// its position in the list as it is, and a write through it sets the whole
// list. Clearing the element removes it from the list. Writing an element
// past the end of the list is an error.
func (n *Node) Index(i int) *Node { return core.Must(n.IndexE(i)) }

// IndexE is Index, returning an error rather than panicking.
func (n *Node) IndexE(i int) (*Node, error) {
	m, err := n.index(i)
	if err != nil {
		return nil, err
	}
	return &Node{m}, nil
}

func (n *node) index(i int) (node, error) {
	if i < 0 {
		return node{}, core.Errorf("die: index %d is negative", i)
	}
	return node{h: n.h, at: core.ChildPath(n.at, core.IndexStep(i))}, nil
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
	found, err := n.find(expr)
	if err != nil {
		return nil, err
	}
	nodes := make([]*Node, len(found))
	for i, f := range found {
		nodes[i] = &Node{f}
	}
	return nodes, nil
}

func (n *node) find(expr string) ([]node, error) {
	paths, err := core.ResolveJSONPath(n.h, n.at, expr)
	if err != nil {
		return nil, err
	}
	nodes := make([]node, len(paths))
	for i, p := range paths {
		nodes[i] = node{h: n.h, at: p}
	}
	return nodes, nil
}

func (n *node) DiePath() base.Path { return n.at }

// DieRoot returns the root die this is within: see base.DieBase.DieRoot.
func (n *node) DieRoot() base.Die { return base.RootOf(n.h) }

// DieRootType returns the type of the root value this is within.
func (n *node) DieRootType() reflect.Type { return n.h.Base.Type() }

// Get returns a copy of the value, in unstructured form, or nil.
func (n *node) Get() any { return core.Must(n.GetE()) }

// GetE is Get, returning an error rather than panicking.
func (n *node) GetE() (any, error) {
	v, err := n.h.ValueAt(n.at)
	if err != nil || !v.IsValid() || (v.Kind() == reflect.Interface && v.IsNil()) {
		return nil, err
	}
	return v.Interface(), nil
}

// Location is a location in an unstructured die, immutable or mutable: the
// die itself, or a node within it.
type Location interface {
	unstructuredAt() (*core.Holder, base.Path)
}

// ImmutableLocation is a location in an immutable unstructured die: an
// Object, or a Node.
type ImmutableLocation interface {
	Location
	immutableAt()
}

// MutableLocation is a location in a mutable unstructured die: an
// ObjectMut, or a NodeMut.
type MutableLocation interface {
	Location
	mutableAt()
}

func (*Object) immutableAt()  {}
func (*Node) immutableAt()    {}
func (*ObjectMut) mutableAt() {}
func (*NodeMut) mutableAt()   {}

func (d *Object) unstructuredAt() (*core.Holder, base.Path) { return dies.MustStateOf(d) }

func (d *ObjectMut) unstructuredAt() (*core.Holder, base.Path) { return dies.MustStateOf(d) }

func (n *node) unstructuredAt() (*core.Holder, base.Path) { return n.h, n.at }

// GetAs reads the value of an unstructured die or node as a T. A die is read
// as a detached copy; to keep it linked, use JigAs.
func GetAs[T any](u Location) T { return core.Must(GetAsE[T](u)) }

// GetAsE is GetAs, returning an error rather than panicking.
func GetAsE[T any](u Location) (T, error) {
	h, at := u.unstructuredAt()
	return core.GetAt[T](h, at, nil)
}

// JigAs presents an immutable unstructured die or node as an immutable
// typed die, D, such as a *DeploymentDie or a *ContainerDie. Like any jig,
// it holds no state of its own: reads come from the unstructured die, so
// what D's type does not describe, such as fields newer than it, is kept.
// A mutable die or node has JigAsMut.
func JigAs[D base.ImmutableDie](u ImmutableLocation) D { return jigAs[D](u) }

// JigAsMut presents a mutable unstructured die or node as a mutable typed
// die, M, such as a *DeploymentDieMut. Writes go to the unstructured die as
// intent, changing only the fields written.
func JigAsMut[M base.MutableDie](u MutableLocation) M { return jigAs[M](u) }

func jigAs[D base.Die](u Location) D {
	h, at := u.unstructuredAt()
	d := reflect.New(reflect.TypeFor[D]().Elem()).Interface().(D)
	t := core.ValueTypeOf(d)
	src := core.DieSource{H: h}
	core.SetState(d, &core.Holder{Base: reflect.Zero(t), Mutable: src.Mutable(), View: core.WholeView(src, at, t)}, nil)
	return d
}

// Has reports whether the value is present.
func (n *node) Has() bool { return core.Must(n.HasE()) }

// HasE is Has, returning an error rather than panicking.
func (n *node) HasE() (bool, error) { return n.h.Has(n.at) }

// Set records the value as intent. Typed values, such as structs, are
// converted to unstructured form.
func (n *NodeMut) Set(v any) *NodeMut { return core.Must(n.SetE(v)) }

// SetE is Set, returning an error rather than panicking.
func (n *NodeMut) SetE(v any) (*NodeMut, error) {
	if _, err := n.h.Write(n.at, nil, core.SetOp(core.OpSet, reflect.ValueOf(&v).Elem()), false); err != nil {
		return nil, err
	}
	return n, nil
}

// Clear records that the value is absent: null in patches, or a removed
// list element.
func (n *NodeMut) Clear() *NodeMut { return core.Must(n.ClearE()) }

// ClearE is Clear, returning an error rather than panicking.
func (n *NodeMut) ClearE() (*NodeMut, error) {
	if _, err := n.h.Write(n.at, nil, core.SetOp(core.OpClear, reflect.Value{}), false); err != nil {
		return nil, err
	}
	return n, nil
}

// EmptyObject records that the value is present and an empty object.
func (n *NodeMut) EmptyObject() *NodeMut { return core.Must(n.EmptyObjectE()) }

// EmptyObjectE is EmptyObject, returning an error rather than panicking.
func (n *NodeMut) EmptyObjectE() (*NodeMut, error) {
	if _, err := n.h.Write(n.at, nil, core.SetOp(core.OpEmpty, reflect.ValueOf(map[string]any{})), false); err != nil {
		return nil, err
	}
	return n, nil
}

// EmptyList records that the value is present and an empty list.
func (n *NodeMut) EmptyList() *NodeMut { return core.Must(n.EmptyListE()) }

// EmptyListE is EmptyList, returning an error rather than panicking.
func (n *NodeMut) EmptyListE() (*NodeMut, error) {
	if _, err := n.h.Write(n.at, nil, core.SetOp(core.OpEmpty, reflect.ValueOf([]any{})), false); err != nil {
		return nil, err
	}
	return n, nil
}

// Keys returns the fields of this object, sorted.
func (n *node) Keys() []string { return core.Must(n.KeysE()) }

// KeysE is Keys, returning an error rather than panicking.
func (n *node) KeysE() ([]string, error) {
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
func (n *node) Len() int { return core.Must(n.LenE()) }

// LenE is Len, returning an error rather than panicking.
func (n *node) LenE() (int, error) {
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
func (n *node) Type() Type { return core.Must(n.TypeE()) }

// TypeE is Type, returning an error rather than panicking.
func (n *node) TypeE() (Type, error) {
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
