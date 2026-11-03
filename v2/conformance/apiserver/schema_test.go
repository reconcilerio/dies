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

package apiserver_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"sort"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"

	"reconciler.io/dies/v2/apis"
	"reconciler.io/dies/v2/die/base"
)

// The server's OpenAPI v3 schema says how it merges each field in a
// server-side apply: lists as atomic, sets or maps by keys; structs and maps
// as atomic or granular. A die encodes its patches from its own markers, so
// where they disagree, a patch from a die doesn't do what it says. This
// compares every die, from every object kind the server serves, with the
// schema.

type schemaProps struct {
	Ref                  string                  `json:"$ref"`
	AllOf                []schemaProps           `json:"allOf"`
	Type                 string                  `json:"type"`
	Properties           map[string]*schemaProps `json:"properties"`
	Items                *schemaProps            `json:"items"`
	AdditionalProperties *schemaProps            `json:"additionalProperties"`
	ListType             string                  `json:"x-kubernetes-list-type"`
	ListMapKeys          []string                `json:"x-kubernetes-list-map-keys"`
	MapType              string                  `json:"x-kubernetes-map-type"`
}

// ref is the definition a property refers to, if any.
func (p *schemaProps) ref() string {
	if p.Ref != "" {
		return p.Ref
	}
	for _, a := range p.AllOf {
		if a.Ref != "" {
			return a.Ref
		}
	}
	return ""
}

// schemas returns every definition the server publishes, by name.
func schemas(t *testing.T) map[string]*schemaProps {
	t.Helper()
	dc, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	paths, err := dc.OpenAPIV3().Paths()
	if err != nil {
		t.Fatal(err)
	}
	defs := map[string]*schemaProps{}
	for _, gv := range paths {
		raw, err := gv.Schema("application/json")
		if err != nil {
			t.Fatal(err)
		}
		var doc struct {
			Components struct {
				Schemas map[string]*schemaProps `json:"schemas"`
			} `json:"components"`
		}
		if err := json.Unmarshal(raw, &doc); err != nil {
			t.Fatal(err)
		}
		for name, s := range doc.Components.Schemas {
			defs[name] = s
		}
	}
	return defs
}

// definitionName is the name the server gives a Go type's schema:
// k8s.io/api/apps/v1.Deployment is io.k8s.api.apps.v1.Deployment.
func definitionName(t reflect.Type) string {
	host, rest, _ := strings.Cut(t.PkgPath(), "/")
	labels := strings.Split(host, ".")
	slices.Reverse(labels)
	return strings.Join(labels, ".") + "." + strings.ReplaceAll(rest, "/", ".") + "." + t.Name()
}

func TestDiesConformToServerSchema(t *testing.T) {
	defs := schemas(t)
	scheme := runtime.NewScheme()
	if err := apis.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}

	c := &schemaCheck{defs: defs, seen: map[reflect.Type]bool{}}
	objects := 0
	for _, gvk := range sortedKinds(scheme) {
		dieType := scheme.AllKnownTypes()[gvk]
		release, ok := reflect.PointerTo(dieType).MethodByName("DieRelease")
		if !ok || strings.HasSuffix(gvk.Kind, "List") {
			continue
		}
		// DieRelease returns a pointer to a copy of the value
		goType := release.Type.Out(0).Elem()
		def := defs[definitionName(goType)]
		if def == nil {
			// not served, such as AdmissionReview, which is only a request
			// body
			continue
		}
		objects++
		c.check(goType, def)
	}
	if objects == 0 {
		t.Fatal("no object kinds were compared")
	}
	t.Logf("compared %d object kinds and %d types with the server's schema", objects, len(c.seen))
	for _, m := range c.mismatches {
		t.Error(m)
	}
	if len(c.sets) > 0 {
		t.Logf("%d fields are sets, which dies send whole, as atomic lists:\n  %s", len(c.sets), strings.Join(c.sets, "\n  "))
	}
}

type schemaCheck struct {
	defs       map[string]*schemaProps
	seen       map[reflect.Type]bool
	mismatches []string
	sets       []string
}

// check compares a struct's fields, reported as fields of its type, such as
// meta.v1.ObjectMeta.finalizers, wherever the type is used.
func (c *schemaCheck) check(t reflect.Type, def *schemaProps) {
	if c.seen[t] {
		return
	}
	c.seen[t] = true
	name := strings.TrimPrefix(strings.TrimPrefix(definitionName(t), "io.k8s.api."), "io.k8s.")
	for _, f := range base.ManagedFieldsOf(t) {
		at := name + "." + f.Wire
		prop := def.Properties[f.Wire]
		if prop == nil {
			c.mismatches = append(c.mismatches, fmt.Sprintf("%s: not in the server's schema", at))
			continue
		}
		field, ok := wireField(t, f.Wire)
		if !ok {
			c.mismatches = append(c.mismatches, fmt.Sprintf("%s: not a field of %s", at, t))
			continue
		}
		c.checkField(f, field.Type, prop, at)
	}
}

func (c *schemaCheck) checkField(f base.FieldInfo, typ reflect.Type, prop *schemaProps, at string) {
	target := prop
	if ref := prop.ref(); ref != "" {
		target = c.defs[strings.TrimPrefix(ref, "#/components/schemas/")]
		if target == nil {
			c.mismatches = append(c.mismatches, fmt.Sprintf("%s: the server's schema refers to %s, which it doesn't have", at, ref))
			return
		}
	}
	switch {
	case prop.Type == "array":
		listType := prop.ListType
		if listType == "" {
			listType = "atomic"
		}
		switch {
		case listType == "map" && f.Kind != base.FieldListMap:
			c.mismatches = append(c.mismatches, fmt.Sprintf("%s: the server merges by keys %v, the die treats it as %s", at, prop.ListMapKeys, f.Kind))
		case listType != "map" && f.Kind == base.FieldListMap:
			c.mismatches = append(c.mismatches, fmt.Sprintf("%s: the server's list is %s, the die merges by keys %v", at, listType, f.Keys))
		case listType == "map" && !sameKeys(prop.ListMapKeys, f.Keys):
			c.mismatches = append(c.mismatches, fmt.Sprintf("%s: the server merges by keys %v, the die by %v", at, prop.ListMapKeys, f.Keys))
		case listType == "set":
			c.sets = append(c.sets, at)
		}
		if f.Kind == base.FieldListMap && prop.Items != nil {
			c.checkElem(deref(typ).Elem(), prop.Items)
		}
	case target.Properties != nil:
		// a struct
		atomic := mapType(prop, target) == "atomic"
		switch {
		case atomic && f.Kind == base.FieldStruct && !f.Atomic:
			c.mismatches = append(c.mismatches, fmt.Sprintf("%s: the server replaces it whole, the die merges its fields", at))
		case !atomic && f.Atomic:
			c.mismatches = append(c.mismatches, fmt.Sprintf("%s: the server merges its fields, the die applies it whole", at))
		case !atomic && f.Kind != base.FieldStruct:
			c.mismatches = append(c.mismatches, fmt.Sprintf("%s: the server merges its fields, the die treats it as %s", at, f.Kind))
		}
		if f.Kind == base.FieldStruct {
			c.check(deref(typ), target)
		}
	case target.AdditionalProperties != nil:
		// a map
		atomic := mapType(prop, target) == "atomic"
		switch {
		case atomic && f.Kind == base.FieldMap:
			c.mismatches = append(c.mismatches, fmt.Sprintf("%s: the server replaces the map whole, the die merges its entries", at))
		case !atomic && f.Kind != base.FieldMap:
			c.mismatches = append(c.mismatches, fmt.Sprintf("%s: the server merges the map's entries, the die treats it as %s", at, f.Kind))
		}
	}
}

// mapType is how the server merges a struct or map: as the property says,
// such as +structType=granular on a field of an atomic type, or else as its
// type does.
func mapType(prop, target *schemaProps) string {
	if prop.MapType != "" {
		return prop.MapType
	}
	return target.MapType
}

// checkElem checks the element type of a list map.
func (c *schemaCheck) checkElem(typ reflect.Type, items *schemaProps) {
	target := items
	if ref := items.ref(); ref != "" {
		target = c.defs[strings.TrimPrefix(ref, "#/components/schemas/")]
	}
	if target == nil || target.Properties == nil {
		return
	}
	c.check(deref(typ), target)
}

// wireField finds the field of a struct with a wire name, through embedded
// structs whose fields are inlined, as encoding/json does.
func wireField(t reflect.Type, wire string) (reflect.StructField, bool) {
	for i := range t.NumField() {
		sf := t.Field(i)
		name, _, _ := strings.Cut(sf.Tag.Get("json"), ",")
		if sf.Anonymous && name == "" {
			if f, ok := wireField(deref(sf.Type), wire); ok {
				return f, true
			}
			continue
		}
		if name == wire || (name == "" && sf.Name == wire) {
			return sf, true
		}
	}
	return reflect.StructField{}, false
}

func deref(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}

func sameKeys(a, b []string) bool {
	a, b = slices.Clone(a), slices.Clone(b)
	slices.Sort(a)
	slices.Sort(b)
	return slices.Equal(a, b)
}

func sortedKinds(scheme *runtime.Scheme) []schema.GroupVersionKind {
	var out []schema.GroupVersionKind
	for gvk := range scheme.AllKnownTypes() {
		out = append(out, gvk)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].String() < out[j].String() })
	return out
}
