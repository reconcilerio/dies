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
	"bytes"
	"encoding/json"
	"go/constant"
	"go/types"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/tools/go/packages"
)

// The markers that validate and default a field, rather than shape its die,
// become die tag options that die.ValidateCreate, die.ValidateUpdate and
// die.Default read: Kubernetes'
// declarative validation markers, such as +k8s:minimum=0, and kubebuilder's,
// such as +kubebuilder:validation:Minimum=0, alike. They are read from the
// field, and from its type, or its item type for a list or map, as a marker
// on a type holds for every field of it.

// checkKeys are the die tag options of value checks, in the order they are
// written.
var checkKeys = []string{
	"enum", "format", "pattern",
	"minimum", "exclusiveMinimum", "maximum", "exclusiveMaximum", "multipleOf",
	"minLength", "maxLength", "maxBytes",
	"minItems", "maxItems", "unique",
	"minProperties", "maxProperties",
	"immutable", "update",
}

// checkKey reports whether a die tag option is a value check, of the value
// or, prefixed each., of its items.
func checkKey(k string) bool {
	return slices.Contains(checkKeys, strings.TrimPrefix(k, "each."))
}

// setCheck returns the checks with the option for key replaced, or added in
// order, without changing the checks it's given, which a spec copy shares.
func setCheck(checks []string, key, opt string) []string {
	out := slices.DeleteFunc(slices.Clone(checks), func(o string) bool {
		k, _, _ := strings.Cut(o, "=")
		return k == key
	})
	if opt != "" {
		out = append(out, opt)
	}
	rank := func(o string) int {
		k, _, _ := strings.Cut(o, "=")
		r := slices.Index(checkKeys, strings.TrimPrefix(k, "each."))
		if strings.HasPrefix(k, "each.") {
			r += len(checkKeys)
		}
		return r
	}
	slices.SortStableFunc(out, func(a, b string) int { return rank(a) - rank(b) })
	return out
}

// kubebuilder's markers of a check, and the die tag option of each
var kubebuilderChecks = map[string]string{
	"Enum": "enum", "Format": "format", "Pattern": "pattern",
	"Minimum": "minimum", "ExclusiveMinimum": "exclusiveMinimum", "Maximum": "maximum", "ExclusiveMaximum": "exclusiveMaximum", "MultipleOf": "multipleOf",
	"MinLength": "minLength", "MaxLength": "maxLength",
	"MinItems": "minItems", "MaxItems": "maxItems", "UniqueItems": "unique",
	"MinProperties": "minProperties", "MaxProperties": "maxProperties",
}

// declarative validation's markers of a check, with the k8s: prefix that
// markers removes
var declarativeChecks = []string{"format", "minimum", "maximum", "minLength", "maxLength", "maxBytes", "minItems", "maxItems", "minProperties", "maxProperties"}

// valueChecks returns the die tag options of a field's value checks: those
// of its type's markers, then its own over them.
func (g *generator) valueChecks(m map[string][]string, t types.Type) []string {
	var checks []string
	n := namedOf(t)
	if n != nil {
		checks = g.addChecks(checks, markers(g.typeDocs[n.Obj()]), n, "")
	}
	item := itemOf(t)
	if in := namedOf(item); in != nil {
		checks = g.addChecks(checks, markers(g.typeDocs[in.Obj()]), in, "each.")
	}
	checks = g.addChecks(checks, m, n, "")
	for _, v := range m["eachVal"] {
		checks = g.addChecks(checks, markers([]string{v}), namedOf(item), "each.")
	}
	return checks
}

// addChecks adds the checks of markers, for a value of type n if it's
// named, with their options prefixed.
func (g *generator) addChecks(checks []string, m map[string][]string, n *types.Named, prefix string) []string {
	set := func(key, v string) {
		if v == "" {
			checks = setCheck(checks, prefix+key, "")
			return
		}
		checks = setCheck(checks, prefix+key, prefix+key+"="+v)
	}
	flag := func(key, v string) {
		if v == "" || v == "true" {
			checks = setCheck(checks, prefix+key, prefix+key)
		} else {
			checks = setCheck(checks, prefix+key, "")
		}
	}
	for _, k := range declarativeChecks {
		if v, ok := m[k]; ok {
			set(k, checkValue(k, first(v)))
		}
	}
	if _, ok := m["enum"]; ok && n != nil {
		if vs := enumValues(n); len(vs) > 0 {
			set("enum", jsonOf(vs))
		}
	}
	if v, ok := m["unique"]; ok && first(v) == "set" {
		flag("unique", "")
	}
	// the markers of an update
	if _, ok := m["immutable"]; ok {
		flag("immutable", "")
	}
	if v, ok := m["update"]; ok {
		var known []string
		for _, u := range v {
			if slices.Contains([]string{"NoSet", "NoUnset", "NoModify", "NoAddItem", "NoRemoveItem"}, u) {
				known = append(known, u)
			}
		}
		set("update", strings.Join(known, ","))
	}
	for mk, v := range m {
		name, ok := strings.CutPrefix(mk, "kubebuilder:validation:")
		if !ok {
			continue
		}
		if item, ok := strings.CutPrefix(name, "items:"); ok {
			if _, ok := kubebuilderChecks[item]; ok {
				checks = g.addChecks(checks, map[string][]string{"kubebuilder:validation:" + item: v}, nil, prefix+"each.")
			}
			continue
		}
		key, ok := kubebuilderChecks[name]
		if !ok {
			continue
		}
		switch key {
		case "exclusiveMinimum", "exclusiveMaximum", "unique":
			flag(key, first(v))
		case "enum":
			var vs []string
			for e := range strings.SplitSeq(first(v), ";") {
				vs = append(vs, unquote(e))
			}
			set(key, jsonOf(vs))
		case "pattern":
			set(key, jsonOf(unquote(first(v))))
		default:
			set(key, checkValue(key, unquote(first(v))))
		}
	}
	return checks
}

// checkValue returns a check's value as its die tag option has it: a number
// for a bound, or the value, or "" if it isn't valid.
func checkValue(key, v string) string {
	if key == "format" {
		return v
	}
	// a bound may be followed by a comment: +k8s:maximum=1000 # Highest
	v, _, _ = strings.Cut(v, "#")
	v = strings.TrimSpace(v)
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return ""
	}
	switch key {
	case "minimum", "maximum", "multipleOf":
		return strconv.FormatFloat(f, 'g', -1, 64)
	}
	if f != float64(int(f)) || f < 0 {
		return ""
	}
	return strconv.Itoa(int(f))
}

// enumValues returns the values of an enum type: the string constants of
// its type in its package, as declarative validation's +k8s:enum has them.
func enumValues(n *types.Named) []string {
	scope := n.Obj().Pkg().Scope()
	var vs []string
	for _, name := range scope.Names() {
		c, ok := scope.Lookup(name).(*types.Const)
		if ok && types.Identical(c.Type(), n) && c.Val().Kind() == constant.String {
			if v := constant.StringVal(c.Val()); !slices.Contains(vs, v) {
				vs = append(vs, v)
			}
		}
	}
	slices.Sort(vs)
	return vs
}

// defaultOf returns a field's default as JSON: +default's, which is JSON,
// or kubebuilder's, which may be a string without quotes.
func defaultOf(m map[string][]string, t types.Type) string {
	if d := first(m["default"]); json.Valid([]byte(d)) {
		return d
	}
	for _, k := range []string{"kubebuilder:default", "kubebuilder:default:"} {
		v, ok := m[k]
		if !ok {
			continue
		}
		d := first(v)
		if json.Valid([]byte(d)) {
			return d
		}
		if b, ok := types.Unalias(deref(t)).Underlying().(*types.Basic); ok && b.Info()&types.IsString != 0 {
			return jsonOf(unquote(d))
		}
	}
	return ""
}

// required returns how validation requires a field: nonzero for
// +k8s:required, as declarative validation has it, present for +required
// and kubebuilder's, as an OpenAPI schema's required fields are, and, by
// kubebuilder's convention in a package of custom resources, present for a
// field without omitempty or omitzero, unless it's optional.
func (g *generator) required(owner *types.Named, wf wfield, m map[string][]string) string {
	_, opt := m["optional"]
	_, kopt := m["kubebuilder:validation:Optional"]
	if opt || kopt {
		return ""
	}
	if slices.Contains(g.fieldDocs[wf.v], "+k8s:required") {
		return "nonzero"
	}
	_, req := m["required"]
	_, kreq := m["kubebuilder:validation:Required"]
	if req || kreq {
		return "present"
	}
	if kb, ok := g.kubebuilder[owner.Obj().Pkg().Path()]; ok && kb.required {
		name, opts, _ := strings.Cut(wf.jsonTag, ",")
		omit := false
		for o := range strings.SplitSeq(opts, ",") {
			omit = omit || o == "omitempty" || o == "omitzero" || o == "inline"
		}
		if !wf.v.Embedded() && name != "" && name != "-" && !omit {
			return "present"
		}
	}
	return ""
}

// unionOf returns the union a field is a member of, from declarative
// validation's markers, and whether exactly one of its members must be set.
// The unions of +union and +unionDiscriminator, whose discriminator's values
// don't name their members, are not read.
func unionOf(m map[string][]string) (string, bool) {
	if v, ok := m["oneOf"]; ok {
		return first(v), true
	}
	if v, ok := m["zeroOrOneOf"]; ok {
		return first(v), false
	}
	if _, ok := m["unionMember"]; ok {
		return "union", true
	}
	if _, ok := m["zeroOrOneOfMember"]; ok {
		return "union", false
	}
	return "", false
}

// kubebuilderPackage is a package of custom resources whose schemas
// controller-gen generates, which has kubebuilder's conventions.
type kubebuilderPackage struct {
	// required is set unless the package says its fields are optional by
	// default
	required bool
}

// crdMarkers are the markers that make a package one of custom resources:
// those of controller-gen's that generate a CRD. Markers that only
// validate, such as those of metav1.Condition for custom resources that use
// it, don't.
var crdMarkers = []string{
	"+kubebuilder:object:root", "+kubebuilder:resource", "+kubebuilder:subresource",
	"+kubebuilder:printcolumn", "+kubebuilder:storageversion", "+kubebuilder:skipversion",
}

// indexKubebuilder records whether a package is one of custom resources.
func (g *generator) indexKubebuilder(p *packages.Package) {
	kb := &kubebuilderPackage{required: true}
	found := false
	for _, l := range packageDoc(p) {
		switch {
		case strings.HasPrefix(l, "+kubebuilder:validation:Optional"):
			kb.required, found = false, true
		case strings.HasPrefix(l, "+kubebuilder:validation:Required"):
			found = true
		}
	}
	for _, f := range p.Syntax {
		for _, cg := range f.Comments {
			for _, l := range commentLines(cg) {
				found = found || slices.ContainsFunc(crdMarkers, func(m string) bool { return strings.HasPrefix(l, m) })
			}
		}
	}
	if found {
		g.kubebuilder[p.PkgPath] = kb
	}
}

func namedOf(t types.Type) *types.Named {
	if t == nil {
		return nil
	}
	n, _ := types.Unalias(deref(t)).(*types.Named)
	return n
}

// itemOf returns the item type of a list, other than bytes, or the value
// type of a map, or nil.
func itemOf(t types.Type) types.Type {
	if t == nil {
		return nil
	}
	switch u := types.Unalias(deref(t)).Underlying().(type) {
	case *types.Slice:
		if b, ok := u.Elem().Underlying().(*types.Basic); ok && b.Kind() == types.Byte {
			return nil
		}
		return u.Elem()
	case *types.Map:
		return u.Elem()
	}
	return nil
}

func unquote(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '`' && s[len(s)-1] == '`' {
		return s[1 : len(s)-1]
	}
	if u, err := strconv.Unquote(s); err == nil {
		return u
	}
	return s
}

// jsonOf returns a value as JSON, without escaping HTML, so that a pattern
// reads as written.
func jsonOf(v any) string {
	var b bytes.Buffer
	e := json.NewEncoder(&b)
	e.SetEscapeHTML(false)
	_ = e.Encode(v)
	return strings.TrimSuffix(b.String(), "\n")
}
