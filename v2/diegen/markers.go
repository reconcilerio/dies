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

package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"reflect"
	"slices"
	"strings"

	"golang.org/x/tools/go/packages"
)

type generator struct {
	typeDocs  map[*types.TypeName][]string
	fieldDocs map[*types.Var][]string
	indexed   map[*packages.Package]bool
	// meta is the metadata of every die package read, by import path
	meta map[string]*pkgMeta
	// ignored counts type errors outside die metadata
	ignored int
	// exists records whether die packages referred to exist
	exists map[string]bool
	// acs has the apply configurations of each die package
	acs map[string]*acSet
}

func newGenerator() *generator {
	return &generator{
		typeDocs:  map[*types.TypeName][]string{},
		fieldDocs: map[*types.Var][]string{},
		indexed:   map[*packages.Package]bool{},
		meta:      map[string]*pkgMeta{},
		exists:    map[string]bool{},
		acs:       map[string]*acSet{},
	}
}

// index records the doc comments, and so the markers, of every type and
// struct field in a package and its dependencies.
func (g *generator) index(p *packages.Package) {
	if g.indexed[p] {
		return
	}
	g.indexed[p] = true
	for _, imp := range p.Imports {
		g.index(imp)
	}
	if p.TypesInfo == nil {
		return
	}
	for _, f := range p.Syntax {
		for _, decl := range f.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}
			for _, spec := range gd.Specs {
				ts := spec.(*ast.TypeSpec)
				obj, ok := p.TypesInfo.Defs[ts.Name].(*types.TypeName)
				if !ok {
					continue
				}
				doc := ts.Doc
				if doc == nil && len(gd.Specs) == 1 {
					doc = gd.Doc
				}
				g.typeDocs[obj] = commentLines(doc)
				sty, ok := ts.Type.(*ast.StructType)
				if !ok || ts.Assign.IsValid() {
					continue
				}
				st := obj.Type().Underlying().(*types.Struct)
				i := 0
				for _, f := range sty.Fields.List {
					for range max(len(f.Names), 1) {
						g.fieldDocs[st.Field(i)] = commentLines(f.Doc)
						i++
					}
				}
			}
		}
	}
}

func commentLines(cg *ast.CommentGroup) []string {
	if cg == nil {
		return nil
	}
	var lines []string
	for _, c := range cg.List {
		l := strings.TrimPrefix(c.Text, "//")
		lines = append(lines, strings.TrimPrefix(l, " "))
	}
	return lines
}

// markers returns the upstream +markers in doc lines. The plain and k8s:
// forms of a marker are equivalent, and upstream often states both. Markers
// gated on a stability level, such as +k8s:alpha(since: "1.37")=..., are
// validation only and are ignored.
func markers(lines []string) map[string][]string {
	m := map[string][]string{}
	for _, l := range lines {
		if !strings.HasPrefix(l, "+") {
			continue
		}
		k, v, _ := strings.Cut(l[1:], "=")
		k = strings.TrimPrefix(k, "k8s:")
		if strings.Contains(k, "(") || slices.Contains(m[k], v) {
			continue
		}
		m[k] = append(m[k], v)
	}
	return m
}

// dieMarker is a +die marker: +die, or +die:<name>:<key>=<value>,<key>=<value>.
// A value of {a,b} is a list.
type dieMarker struct {
	name string
	args map[string][]string
}

// arg returns the first value of an argument.
func (m dieMarker) arg(key string) string {
	if v := m.args[key]; len(v) > 0 {
		return v[0]
	}
	return ""
}

func dieMarkers(lines []string) []dieMarker {
	var out []dieMarker
	for _, l := range lines {
		if l != "+die" && !strings.HasPrefix(l, "+die:") {
			continue
		}
		rest := l[1:]
		m := dieMarker{name: rest, args: map[string][]string{}}
		if i := strings.Index(rest, "="); i >= 0 {
			head := rest[:i]
			j := strings.LastIndex(head, ":")
			m.name = head[:j]
			for _, kv := range splitTop(rest[j+1:]) {
				k, v, _ := strings.Cut(kv, "=")
				// a default is JSON, whose objects are not lists
				if k != "default" && strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}") {
					m.args[k] = append(m.args[k], splitTop(v[1:len(v)-1])...)
				} else {
					m.args[k] = append(m.args[k], v)
				}
			}
		}
		out = append(out, m)
	}
	return out
}

// splitTop splits on commas that are not within braces, brackets or quotes,
// so values such as default="a,b" and {x,y} stay whole.
func splitTop(s string) []string {
	var out []string
	depth, quoted, start := 0, false, 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case quoted && c == '\\':
			i++
		case c == '"':
			quoted = !quoted
		case quoted:
		case c == '{' || c == '[':
			depth++
		case c == '}' || c == ']':
			depth--
		case c == ',' && depth == 0:
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		out = append(out, s[start:])
	}
	return out
}

// fieldSpec is what a die needs to know about a field beyond its Go type.
type fieldSpec struct {
	listType   string
	listMapKey []string
	mapType    string
	structType string
	def        string
	union      string
	// nullable is set for a field whose schema accepts null, which a custom
	// resource's validation otherwise rejects in an apply
	nullable bool
	// accessor renames the field's die method, such as to avoid a collision
	accessor string
}

func (f *fieldSpec) atomic() bool {
	return f != nil && (f.listType == "atomic" || f.listType == "set" || f.mapType == "atomic" || f.structType == "atomic")
}

func (f *fieldSpec) listMap() bool { return f != nil && f.listType == "map" }

// typeSpec is the structural semantics of a type: the markers that change a
// die's API or its patches. Validation markers are not part of it.
type typeSpec struct {
	fields       map[string]*fieldSpec
	subresources []string
}

func newTypeSpec() *typeSpec { return &typeSpec{fields: map[string]*fieldSpec{}} }

func (s *typeSpec) field(wire string) *fieldSpec {
	if s.fields[wire] == nil {
		s.fields[wire] = &fieldSpec{}
	}
	return s.fields[wire]
}

func (s *typeSpec) subresource(name string) {
	if !slices.Contains(s.subresources, name) {
		s.subresources = append(s.subresources, name)
		slices.Sort(s.subresources)
	}
}

// upstreamSpec normalizes the upstream markers of a type: the plain, k8s:
// and kubebuilder: forms.
func (g *generator) upstreamSpec(n *types.Named, st *types.Struct) *typeSpec {
	s := newTypeSpec()
	for _, wf := range g.directFields(st) {
		m := markers(g.fieldDocs[wf.v])
		fs := fieldSpec{listType: first(m["listType"]), listMapKey: m["listMapKey"], mapType: first(m["mapType"]), structType: first(m["structType"])}
		// a marker on the field's type, such as +structType=atomic on
		// ObjectReference, holds for every field of that type, unless the
		// field says otherwise
		if tm := g.typeMarkers(wf.v.Type()); tm != nil {
			if fs.listType == "" {
				fs.listType = first(tm["listType"])
			}
			if fs.listType == "map" && len(fs.listMapKey) == 0 {
				fs.listMapKey = tm["listMapKey"]
			}
			if fs.mapType == "" {
				fs.mapType = first(tm["mapType"])
			}
			if fs.structType == "" {
				fs.structType = first(tm["structType"])
			}
		}
		if d := first(m["default"]); json.Valid([]byte(d)) {
			fs.def = d
		}
		if _, ok := m["nullable"]; ok {
			fs.nullable = true
		}
		if !reflect.ValueOf(fs).IsZero() {
			*s.field(wf.wire) = fs
		}
	}
	tm := markers(g.typeDocs[n.Obj()])
	for _, v := range tm["supportsSubresource"] {
		s.subresource(strings.Trim(v, `"/`))
	}
	if _, ok := tm["kubebuilder:subresource:status"]; ok {
		s.subresource("status")
	}
	return s
}

// typeMarkers returns the markers of a field's named type, through a
// pointer, if it has a declaration that was indexed.
func (g *generator) typeMarkers(t types.Type) map[string][]string {
	t = types.Unalias(t)
	if p, ok := t.(*types.Pointer); ok {
		t = types.Unalias(p.Elem())
	}
	n, ok := t.(*types.Named)
	if !ok {
		return nil
	}
	doc, ok := g.typeDocs[n.Obj()]
	if !ok {
		return nil
	}
	return markers(doc)
}

// dieSpec reads the canonical +die markers of a type, over a base spec.
func dieSpec(base *typeSpec, doc []string) *typeSpec {
	s := newTypeSpec()
	if base != nil {
		for k, v := range base.fields {
			c := *v
			s.fields[k] = &c
		}
		s.subresources = slices.Clone(base.subresources)
	}
	for _, m := range dieMarkers(doc) {
		switch m.name {
		case "die:field":
			fs := s.field(m.arg("name"))
			for k, v := range m.args {
				switch k {
				case "listType":
					fs.listType = v[0]
				case "listMapKey":
					fs.listMapKey = v
				case "mapType":
					fs.mapType = v[0]
				case "structType":
					fs.structType = v[0]
				case "default":
					fs.def = v[0]
				case "accessor":
					fs.accessor = v[0]
				case "nullable":
					fs.nullable = v[0] != "false"
				}
			}
		case "die:union":
			for _, member := range m.args["members"] {
				s.field(member).union = m.arg("name")
			}
		case "die:subresource":
			s.subresource(m.arg("name"))
		}
	}
	return s
}

// render writes a spec as canonical +die markers, with fields in declaration
// order.
func (s *typeSpec) render(order []string) []string {
	var lines []string
	unions := map[string][]string{}
	var unionNames []string
	for _, wire := range order {
		fs := s.fields[wire]
		if fs == nil {
			continue
		}
		var args []string
		if fs.listType != "" {
			args = append(args, "listType="+fs.listType)
		}
		if len(fs.listMapKey) == 1 {
			args = append(args, "listMapKey="+fs.listMapKey[0])
		} else if len(fs.listMapKey) > 1 {
			args = append(args, "listMapKey={"+strings.Join(fs.listMapKey, ",")+"}")
		}
		if fs.mapType != "" {
			args = append(args, "mapType="+fs.mapType)
		}
		if fs.structType != "" {
			args = append(args, "structType="+fs.structType)
		}
		if fs.def != "" {
			args = append(args, "default="+fs.def)
		}
		if fs.accessor != "" {
			args = append(args, "accessor="+fs.accessor)
		}
		if fs.nullable {
			args = append(args, "nullable=true")
		}
		if len(args) > 0 {
			lines = append(lines, "+die:field:name="+wire+","+strings.Join(args, ","))
		}
		if fs.union != "" {
			if unions[fs.union] == nil {
				unionNames = append(unionNames, fs.union)
			}
			unions[fs.union] = append(unions[fs.union], wire)
		}
	}
	for _, u := range unionNames {
		lines = append(lines, "+die:union:name="+u+",members={"+strings.Join(unions[u], ",")+"}")
	}
	for _, sub := range s.subresources {
		lines = append(lines, "+die:subresource:name="+sub)
	}
	return lines
}

func first(v []string) string {
	if len(v) > 0 {
		return v[0]
	}
	return ""
}

func findDieMarker(lines []string, name string) (dieMarker, bool) {
	for _, m := range dieMarkers(lines) {
		if m.name == name {
			return m, true
		}
	}
	return dieMarker{}, false
}

// typesMarker returns a die package's +die:types marker, which says where
// the types its dies are for come from and how they're brought in: aliased,
// the default, or cloned.
func typesMarker(lines []string) (m dieMarker, mode string, ok bool, err error) {
	m, ok = findDieMarker(lines, "die:types")
	if !ok {
		return m, "", false, nil
	}
	switch mode = m.arg("mode"); mode {
	case "":
		mode = "alias"
	case "alias", "clone":
	default:
		return m, "", true, fmt.Errorf("+die:types: unknown mode %q, want alias or clone", mode)
	}
	return m, mode, true, nil
}

// docText is the doc without markers or trailing blank lines.
func docText(lines []string) []string {
	var out []string
	for _, l := range lines {
		if !strings.HasPrefix(l, "+") {
			out = append(out, l)
		}
	}
	for len(out) > 0 && strings.TrimSpace(out[len(out)-1]) == "" {
		out = out[:len(out)-1]
	}
	return out
}

// packageDoc returns the comments above the package clause of the first of
// a package's files that has markers there. Markers may be in the package
// doc, or, to keep them out of the godoc, in a comment of their own
// separated from it by a blank line.
func packageDoc(p *packages.Package) []string {
	for _, f := range p.Syntax {
		var lines []string
		for _, cg := range f.Comments {
			if cg.End() >= f.Package {
				break
			}
			lines = append(lines, commentLines(cg)...)
		}
		for _, l := range lines {
			if strings.HasPrefix(l, "+") {
				return lines
			}
		}
	}
	return nil
}

func customJSON(n *types.Named) bool {
	ms := types.NewMethodSet(types.NewPointer(n))
	return ms.Lookup(nil, "MarshalJSON") != nil || ms.Lookup(nil, "UnmarshalJSON") != nil
}

// structOf returns the named struct behind t, if it is one that a die can
// be generated for: a struct without custom JSON encoding.
func structOf(t types.Type) (*types.Named, *types.Struct, bool) {
	n, ok := types.Unalias(t).(*types.Named)
	if !ok {
		return nil, nil, false
	}
	st, ok := n.Underlying().(*types.Struct)
	if !ok || customJSON(n) {
		return nil, nil, false
	}
	return n, st, true
}

func deref(t types.Type) types.Type {
	if p, ok := types.Unalias(t).(*types.Pointer); ok {
		return p.Elem()
	}
	return t
}

func cutPrefix(s, prefix string) (string, bool) { return strings.CutPrefix(s, prefix) }

func lowerFirst(s string) string { return strings.ToLower(s[:1]) + s[1:] }
