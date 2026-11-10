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

// Schema is the structure of an unstructured value, as a server's OpenAPI
// schema describes it: how a server-side apply merges each part of it. A typed
// die has its Go type's descriptions instead. A nil Schema is unknown, such as
// below x-kubernetes-preserve-unknown-fields, and the value is encoded as an
// unstructured die with no schema encodes it.
type Schema struct {
	Kind SchemaKind
	// Fields are an object's fields.
	Fields map[string]*Schema
	// Elem is a map's values, or a list map's elements.
	Elem *Schema
	// Keys are a list map's key fields.
	Keys []string
	// Nullable is set for a value whose schema accepts null.
	Nullable bool
	// GoTypes is set on the schema of a kind its server converts through Go
	// types, where null is absent, as a built-in kind's or an aggregated
	// API's, so it accepts null for any field; a CRD's validation rejects null
	// for a field that isn't nullable.
	GoTypes bool
}

type SchemaKind int

const (
	// SchemaScalar is a string, number, boolean or int-or-string.
	SchemaScalar SchemaKind = iota + 1
	// SchemaObject is an object whose fields are merged one by one.
	SchemaObject
	// SchemaMap is a map whose entries are merged one by one.
	SchemaMap
	// SchemaListMap is a list whose elements are merged by their keys.
	SchemaListMap
	// SchemaAtomic is a value replaced whole: an atomic object, map or list,
	// or a set.
	SchemaAtomic
)

func (s *Schema) field(name string) *Schema {
	if s == nil || s.Kind != SchemaObject {
		return nil
	}
	return s.Fields[name]
}

func (s *Schema) elem() *Schema {
	if s == nil {
		return nil
	}
	return s.Elem
}

func (s *Schema) atomic() bool { return s != nil && s.Kind == SchemaAtomic }

// At returns the schema of the value at a path below s, if it's known.
func (s *Schema) At(p Path) *Schema {
	for _, st := range p {
		if s == nil {
			return nil
		}
		switch st.kind {
		case stepField, stepMapKey:
			switch s.Kind {
			case SchemaObject:
				s = s.Fields[st.name]
			case SchemaMap:
				s = s.Elem
			default:
				return nil
			}
		default:
			s = s.elem()
		}
	}
	return s
}
