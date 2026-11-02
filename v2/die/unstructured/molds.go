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

package unstructured

import (
	"encoding/json"
	"reflect"
	"slices"
	"strings"
	"sync"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/openapi"

	"reconciler.io/dies/v2/die/base"
	"reconciler.io/dies/v2/die/internal/core"
	"reconciler.io/dies/v2/die/internal/dies"
)

// Molds makes unstructured blanks in the shape of the kinds an API server
// serves. A mold is a kind's OpenAPI v3 schema, as the server publishes it:
// as a mold gives material its shape, it gives an unstructured die the
// structure the server merges by, so that its apply patches merge as the
// server merges them, as a typed die's do. Atomic objects, maps and lists are
// sent whole, and Clear on an object's field sends null where the server
// accepts it, which is for any field of a kind served from Go types, as
// built-in kinds and aggregated APIs are, and for a nullable field of a custom
// resource; elsewhere, it's left out.
//
// An unstructured die without a mold can only guess: every object is a map,
// so a field it clears is left out of an apply, and a list addressed by keys
// is a list map, which, if the server's list is atomic, replaces the list
// with what the die sent.
//
// Molds reads the schemas through an OpenAPI client, such as a discovery
// client's OpenAPIV3(), or client-go's cached client in front of it; dies
// make no requests of their own. Each group version's schema is read once,
// when a blank of one of its kinds is first made, and kept, unless the server
// changes it.
//
// Caveats:
//
//   - A die's mold is the schema the server served when its blank was made.
//     Make a blank again to pick up a changed CRD.
//   - A CRD's schema is published shortly after the CRD is established, not
//     with it, so making a blank right after installing a CRD can fail.
//   - Below x-kubernetes-preserve-unknown-fields, and wherever else the schema
//     doesn't say, the die guesses as it does without a mold.
//   - Only apply patches use the mold. Merge patches don't need it: RFC 7386
//     merges objects and maps alike, and replaces lists.
//   - Defaults aren't used: a list-map key left to its default, such as a
//     port's protocol, doesn't match the element with the default set.
//   - Subresources aren't known: DieSubresource and DieMain have no effect.
//   - A client zeroes the object it decodes into, mold included. DieWith and
//     DieWithE keep it; a die passed to a client directly loses it, as do the
//     items of a list a client lists into. For gives them their mold.
type Molds struct {
	client openapi.Client

	mu sync.Mutex
	// paths are the group versions the server has schemas for
	paths map[string]openapi.GroupVersion
	// docs are the parsed schemas of a group version, by its URL, which
	// changes when the schema does
	docs map[string]map[schema.GroupVersionKind]*core.Schema
}

// NewMolds returns the molds of the schemas an OpenAPI client reads.
func NewMolds(client openapi.Client) *Molds {
	return &Molds{client: client, docs: map[string]map[schema.GroupVersionKind]*core.Schema{}}
}

// Blank returns an immutable, empty die of a kind, in its mold. It panics if
// the server has no schema for the kind.
func (m *Molds) Blank(gvk schema.GroupVersionKind) *Object { return core.Must(m.BlankE(gvk)) }

// BlankE is [Molds.Blank], returning an error rather than panicking.
func (m *Molds) BlankE(gvk schema.GroupVersionKind) (*Object, error) {
	sc, err := m.schemaOf(gvk)
	if err != nil {
		return nil, err
	}
	apiVersion, kind := gvk.ToAPIVersionAndKind()
	d := &Object{}
	core.SetState(d, &core.Holder{Base: reflect.ValueOf(map[string]any{"apiVersion": apiVersion, "kind": kind}), Schema: sc}, nil)
	return d, nil
}

// For returns an immutable unstructured die in the mold of its kind,
// keeping its value and intent. It panics if the server has no schema for
// the die's kind.
func (m *Molds) For(d *Object) *Object { return core.Must(m.ForE(d)) }

// ForE is [Molds.For], returning an error rather than panicking.
func (m *Molds) ForE(d *Object) (*Object, error) {
	h, sc, err := m.moldOf(d)
	if err != nil {
		return nil, err
	}
	out := &Object{}
	core.SetState(out, &core.Holder{Base: h.Base, Ov: h.Ov, Wire: h.Wire, Schema: sc}, nil)
	return out, nil
}

// ForMut gives a mutable unstructured die the mold of its kind, in place,
// such as an item of a list a client listed into, keeping its value and
// intent. It panics if the server has no schema for the die's kind.
func (m *Molds) ForMut(d *ObjectMut) *ObjectMut { return core.Must(m.ForMutE(d)) }

// ForMutE is [Molds.ForMut], returning an error rather than panicking.
func (m *Molds) ForMutE(d *ObjectMut) (*ObjectMut, error) {
	h, sc, err := m.moldOf(d)
	if err != nil {
		return nil, err
	}
	h.Schema = sc
	return d, nil
}

func (m *Molds) moldOf(d base.DieObject) (*core.Holder, *core.Schema, error) {
	h, at := dies.MustStateOf(d)
	if len(at) > 0 || h.View != nil {
		return nil, nil, core.Errorf("die: For is only supported on a root die, not at %s", at)
	}
	sc, err := m.schemaOf(d.GetObjectKind().GroupVersionKind())
	if err != nil {
		return nil, nil, err
	}
	return h, sc, nil
}

func (m *Molds) schemaOf(gvk schema.GroupVersionKind) (*core.Schema, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	path := "apis/" + gvk.Group + "/" + gvk.Version
	if gvk.Group == "" {
		path = "api/" + gvk.Version
	}
	gv, ok := m.paths[path]
	if !ok {
		// the server may have started serving it since, as for a new CRD
		paths, err := m.client.Paths()
		if err != nil {
			return nil, core.Errorf("die: reading the server's OpenAPI paths: %w", err)
		}
		m.paths = paths
		if gv, ok = paths[path]; !ok {
			return nil, core.Errorf("die: the server has no OpenAPI schema for %s; a CRD's schema is published shortly after the CRD is established", gvk.GroupVersion())
		}
	}
	key := gv.ServerRelativeURL()
	kinds, ok := m.docs[key]
	if !ok {
		raw, err := gv.Schema("application/json")
		if err != nil {
			return nil, core.Errorf("die: reading the OpenAPI schema for %s: %w", gvk.GroupVersion(), err)
		}
		if kinds, err = parseSchemas(raw); err != nil {
			return nil, core.Errorf("die: the OpenAPI schema for %s: %w", gvk.GroupVersion(), err)
		}
		m.docs[key] = kinds
	}
	sc, ok := kinds[gvk]
	if !ok {
		return nil, core.Errorf("die: the OpenAPI schema for %s has no kind %s", gvk.GroupVersion(), gvk.Kind)
	}
	return sc, nil
}

// openAPISchema is the part of an OpenAPI v3 schema that says how a
// server-side apply merges a value.
type openAPISchema struct {
	Ref                  string                    `json:"$ref"`
	AllOf                []*openAPISchema          `json:"allOf"`
	Type                 string                    `json:"type"`
	Properties           map[string]*openAPISchema `json:"properties"`
	Items                *openAPISchema            `json:"items"`
	AdditionalProperties json.RawMessage           `json:"additionalProperties"`
	ListType             string                    `json:"x-kubernetes-list-type"`
	ListMapKeys          []string                  `json:"x-kubernetes-list-map-keys"`
	MapType              string                    `json:"x-kubernetes-map-type"`
	IntOrString          bool                      `json:"x-kubernetes-int-or-string"`
	Nullable             bool                      `json:"nullable"`
	GroupVersionKinds    []schema.GroupVersionKind `json:"x-kubernetes-group-version-kind"`
}

func (p *openAPISchema) ref() string {
	ref := p.Ref
	for _, a := range p.AllOf {
		if a.Ref != "" {
			ref = a.Ref
		}
	}
	return strings.TrimPrefix(ref, "#/components/schemas/")
}

// parseSchemas reads the schema of each kind in a group version's OpenAPI v3
// document.
func parseSchemas(raw []byte) (map[schema.GroupVersionKind]*core.Schema, error) {
	var doc struct {
		Components struct {
			Schemas map[string]*openAPISchema `json:"schemas"`
		} `json:"components"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, err
	}
	c := &converter{defs: doc.Components.Schemas, done: map[string]*core.Schema{}}
	kinds := map[schema.GroupVersionKind]*core.Schema{}
	for name, def := range c.defs {
		for _, gvk := range def.GroupVersionKinds {
			sc := c.def(name, "")
			if sc == nil {
				continue
			}
			// a copy, as the definition may be another kind's too
			kind := *sc
			kind.GoTypes = !crdDefinition(name, gvk)
			kinds[gvk] = &kind
		}
	}
	return kinds, nil
}

// crdDefinition reports whether a kind's definition is a CRD's, whose
// validation rejects null for a field that isn't nullable, rather than a type
// a server converts through Go types, where null is absent: a built-in type,
// or an aggregated API server's. A CRD's definition is named for its group,
// version and kind, as dev.example.v1.Widget is for example.dev/v1 Widget,
// even for a *.k8s.io group, as Gateway API's are; a Go type's for its
// package, as io.k8s.api.apps.v1.Deployment and
// io.k8s.sample-apiserver.pkg.apis.wardle.v1alpha1.Flunder are.
func crdDefinition(name string, gvk schema.GroupVersionKind) bool {
	labels := strings.Split(gvk.Group, ".")
	slices.Reverse(labels)
	return name == strings.Join(append(labels, gvk.Version, gvk.Kind), ".")
}

type converter struct {
	defs map[string]*openAPISchema
	// done are the definitions converted, by name and map type, which also
	// ends a recursive type's recursion
	done map[string]*core.Schema
}

// def converts a definition, as a property of its type with a map type of
// its own, if any, sees it.
func (c *converter) def(name, mapType string) *core.Schema {
	key := name + "/" + mapType
	if sc, ok := c.done[key]; ok {
		return sc
	}
	p := c.defs[name]
	if p == nil {
		return nil
	}
	// stored before it's filled in, for the types that refer to themselves
	sc := &core.Schema{}
	c.done[key] = sc
	if got := c.convert(p, mapType); got != nil {
		*sc = *got
	} else {
		c.done[key] = nil
		return nil
	}
	return sc
}

// convert converts a schema; mapType, if set, is a property's own map type,
// such as granular on a field of an atomic type.
func (c *converter) convert(p *openAPISchema, mapType string) *core.Schema {
	sc := c.convertType(p, mapType)
	if sc != nil && p.Nullable && !sc.Nullable {
		n := *sc
		n.Nullable = true
		sc = &n
	}
	return sc
}

func (c *converter) convertType(p *openAPISchema, mapType string) *core.Schema {
	if p == nil {
		return nil
	}
	if mapType == "" {
		mapType = p.MapType
	}
	if ref := p.ref(); ref != "" {
		return c.def(ref, mapType)
	}
	if p.IntOrString {
		return &core.Schema{Kind: core.SchemaScalar}
	}
	switch p.Type {
	case "array":
		if p.ListType == "map" {
			return &core.Schema{Kind: core.SchemaListMap, Keys: p.ListMapKeys, Elem: c.convert(p.Items, "")}
		}
		// atomic, the default, or a set, which the die sends whole
		return &core.Schema{Kind: core.SchemaAtomic}
	case "string", "integer", "number", "boolean":
		return &core.Schema{Kind: core.SchemaScalar}
	}
	if len(p.Properties) > 0 {
		if mapType == "atomic" {
			return &core.Schema{Kind: core.SchemaAtomic}
		}
		sc := &core.Schema{Kind: core.SchemaObject, Fields: map[string]*core.Schema{}}
		for name, prop := range p.Properties {
			sc.Fields[name] = c.convert(prop, "")
		}
		return sc
	}
	var values openAPISchema
	if len(p.AdditionalProperties) > 0 && json.Unmarshal(p.AdditionalProperties, &values) == nil {
		if mapType == "atomic" {
			return &core.Schema{Kind: core.SchemaAtomic}
		}
		return &core.Schema{Kind: core.SchemaMap, Elem: c.convert(&values, "")}
	}
	// an object of any shape, such as below x-kubernetes-preserve-unknown-fields
	return nil
}
