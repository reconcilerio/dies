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
	"slices"

	"reconciler.io/dies/v2/die/internal/core"
)

// The fields a die manages are recorded as its field table is built, so a
// die can enumerate the fields it was generated for. That can differ from
// the fields its type has now, if the type moved on without the dies being
// regenerated.

// FieldKind is how a die exposes a field.
type FieldKind string

const (
	// FieldValue is a leaf, read and written whole: Value or Optional.
	FieldValue FieldKind = "value"
	// FieldStruct is a struct with a die of its own.
	FieldStruct FieldKind = "struct"
	// FieldListMap is a list of elements identified by keys.
	FieldListMap FieldKind = "listMap"
	// FieldMap is a map with string keys.
	FieldMap FieldKind = "map"
)

// FieldInfo describes a field a die manages. It is also a location within
// the die, so it can be passed to DieOnly and DieExcept.
type FieldInfo struct {
	// Wire is the field's name on the wire.
	Wire string
	Kind FieldKind
	// Keys are the key fields of a list map.
	Keys []string
	// Union is the union the field is a member of, if any.
	Union string
	// Atomic is set for a struct the server replaces whole: it has a die,
	// but is applied whole.
	Atomic bool

	path Path
	root reflect.Type
}

func (f FieldInfo) DiePath() Path { return f.path }

func (f FieldInfo) DieRootType() reflect.Type { return f.root }

func fieldInfos(t reflect.Type, root reflect.Type, at Path) []FieldInfo {
	core.Built(t)
	core.ManagedMu.Lock()
	fds := slices.Clone(core.Managed[t])
	core.ManagedMu.Unlock()
	infos := make([]FieldInfo, len(fds))
	for i, fd := range fds {
		kind := FieldValue
		switch fd.Kind {
		case core.KindStruct:
			kind = FieldStruct
		case core.KindListMap:
			kind = FieldListMap
		case core.KindMap:
			kind = FieldMap
		}
		infos[i] = FieldInfo{Wire: fd.Wire, Kind: kind, Keys: fd.Keys, Union: fd.Union, Atomic: fd.Atomic, path: core.ChildPath(at, core.FieldStep(fd.Wire)), root: root}
	}
	return infos
}

// ManagedFields returns the fields the die for T manages, in order.
func ManagedFields[T any]() []FieldInfo {
	t := reflect.TypeFor[T]()
	return fieldInfos(t, t, nil)
}

// ManagedFieldsOf returns the fields the die for a type manages, in order,
// for tools that walk types reflectively, such as to compare dies with a
// server's schema. A type with no die manages none.
func ManagedFieldsOf(t reflect.Type) []FieldInfo {
	return fieldInfos(t, t, nil)
}

// DieFields returns the fields this die manages, in order, as locations
// within it.
func (b *DieBase[T, D]) DieFields() []FieldInfo {
	return fieldInfos(reflect.TypeFor[T](), b.DieRootType(), b.at)
}
