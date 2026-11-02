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

package core

import (
	"encoding/json"
	"reflect"
	"strings"
	"sync"
)

// The schema is derived by reflection from json tags plus a `die` struct tag
// that stands in for +listType/+listMapKey/+union markers. A generator would
// emit these descriptors directly.

// kind classifies how a field participates in intent tracking.
type kind uint8

const (
	// kindLeaf values are replaced as a whole: scalars, atomic lists and
	// types with custom JSON encoding.
	kindLeaf kind = iota
	// kindStruct values are granular structs, held by value or by pointer.
	KindStruct
	// kindListMap values are lists of structs keyed by one or more fields.
	KindListMap
	// kindMap values are granular maps with string keys.
	KindMap
)

type FieldDesc struct {
	Wire        string
	index       []int
	Typ         reflect.Type
	Kind        kind
	Keys        []string
	Union       string
	Subresource string
	// Atomic is set for a struct that the server replaces whole, such as an
	// ObjectReference: it has a die, but is applied whole.
	Atomic bool
	// Nullable is set for a field whose schema accepts null, which a custom
	// resource's validation otherwise rejects in an apply.
	Nullable bool
	// omit is set for an omitempty or omitzero field
	omit bool
	// def is the +default value, used to normalize list-map keys.
	def reflect.Value
}

type structDesc struct {
	typ    reflect.Type
	Fields []*FieldDesc
	byWire map[string]*FieldDesc
}

func (sd *structDesc) Field(wire string) *FieldDesc {
	fd, ok := sd.byWire[wire]
	if !ok {
		panic(Errorf("die: %s has no field %q", sd.typ, wire))
	}
	return fd
}

var (
	descs       sync.Map
	annotations sync.Map
)

// annotateType supplies die tags for a type.
func AnnotateType(t reflect.Type, tags map[string]string) {
	if _, ok := descs.Load(t); ok {
		panic(Errorf("die: %s annotated after its schema was built", t))
	}
	annotations.Store(t, tags)
}

// describe returns the wire schema of a struct type. Descriptors are built
// lazily, so recursive types are fine.
func Describe(t reflect.Type) *structDesc {
	t = Deref(t)
	if sd, ok := descs.Load(t); ok {
		return sd.(*structDesc)
	}
	if t.Kind() != reflect.Struct {
		panic(Errorf("die: %s is not a struct", t))
	}
	sd := &structDesc{typ: t, byWire: map[string]*FieldDesc{}}
	sd.collect(t, nil)
	actual, _ := descs.LoadOrStore(t, sd)
	return actual.(*structDesc)
}

// collect walks struct fields, flattening embedded structs whose json name is
// empty, as encoding/json does.
func (sd *structDesc) collect(t reflect.Type, prefix []int) {
	annotate(t)
	var annotated map[string]string
	if a, ok := annotations.Load(t); ok {
		annotated = a.(map[string]string)
	}
	for i := range t.NumField() {
		sf := t.Field(i)
		name, opts, _ := strings.Cut(sf.Tag.Get("json"), ",")
		if name == "-" && opts == "" {
			continue
		}
		index := append(append([]int{}, prefix...), i)
		if sf.Anonymous && name == "" {
			if sf.Type.Kind() != reflect.Struct {
				panic(Errorf("die: embedded %s in %s is not supported", sf.Type, sd.typ))
			}
			sd.collect(sf.Type, index)
			continue
		}
		if !sf.IsExported() {
			continue
		}
		if name == "" {
			name = sf.Name
		}
		fd := &FieldDesc{Wire: name, index: index, Typ: sf.Type}
		for o := range strings.SplitSeq(opts, ",") {
			fd.omit = fd.omit || o == "omitempty" || o == "omitzero"
		}
		tag := sf.Tag.Get("die")
		if a, ok := annotated[sf.Name]; ok {
			tag = a
		}
		atomic := fd.parseTag(tag)
		fd.Kind = classify(fd, atomic)
		if _, dup := sd.byWire[name]; dup {
			panic(Errorf("die: duplicate wire field %q in %s", name, sd.typ))
		}
		sd.Fields = append(sd.Fields, fd)
		sd.byWire[name] = fd
	}
}

func (fd *FieldDesc) parseTag(tag string) (atomic bool) {
	for opt := range strings.SplitSeq(tag, ";") {
		k, v, _ := strings.Cut(strings.TrimSpace(opt), "=")
		switch k {
		case "":
		case "atomic":
			atomic = true
		case "nullable":
			fd.Nullable = true
		case "listMap":
			fd.Keys = strings.Split(v, ",")
		case "union":
			fd.Union = v
		case "subresource":
			fd.Subresource = v
		case "default":
			def := reflect.New(fd.Typ)
			if err := json.Unmarshal([]byte(v), def.Interface()); err != nil {
				panic(Errorf("die: invalid default for %q: %w", fd.Wire, err))
			}
			fd.def = def.Elem()
		default:
			panic(Errorf("die: unknown die tag option %q", k))
		}
	}
	return atomic
}

var marshalerType = reflect.TypeFor[json.Marshaler]()

func customJSON(t reflect.Type) bool {
	return t.Implements(marshalerType) || reflect.PointerTo(t).Implements(marshalerType)
}

func classify(fd *FieldDesc, atomic bool) kind {
	t := fd.Typ
	if customJSON(t) {
		return kindLeaf
	}
	if atomic {
		if s := Deref(t); s.Kind() == reflect.Struct && !customJSON(s) {
			fd.Atomic = true
			return KindStruct
		}
		return kindLeaf
	}
	switch t.Kind() {
	case reflect.Struct:
		return KindStruct
	case reflect.Pointer:
		if t.Elem().Kind() == reflect.Struct && !customJSON(t.Elem()) {
			return KindStruct
		}
	case reflect.Slice:
		if len(fd.Keys) > 0 {
			return KindListMap
		}
	case reflect.Map:
		if t.Key().Kind() == reflect.String {
			return KindMap
		}
	}
	return kindLeaf
}

// keyValue returns the value of a list-map key field, defaulted when unset.
// An element without a defaulted key and one with the default value set
// explicitly are the same element.
func (fd *FieldDesc) keyValue(v reflect.Value) reflect.Value {
	if fd.def.IsValid() && v.IsZero() {
		return fd.def
	}
	return v
}

func Deref(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}

// typeAt returns the static Go type at a path from the root type.
// Within an unstructured value, every type is any.
func TypeAt(t reflect.Type, p Path) reflect.Type {
	for _, s := range p {
		t = Deref(t)
		switch {
		case t.Kind() == reflect.Interface:
			return t
		case s.kind == stepField && t.Kind() == reflect.Struct:
			t = Describe(t).Field(s.name).Typ
		default:
			t = t.Elem()
		}
	}
	return t
}

var anyType = reflect.TypeFor[any]()
