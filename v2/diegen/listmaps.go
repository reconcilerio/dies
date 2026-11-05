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
	"go/types"
	"path/filepath"
	"slices"
	"strings"

	"golang.org/x/tools/go/packages"
)

// A die package names the list maps of its element types, so that code
// that takes one, such as an observer of a resource's conditions, names it
// by its element:
//
//	type ConditionListMap = die.ListMap[Condition, *ConditionDie, string]
//	type ConditionListMapMut = die.ListMapMut[Condition, *ConditionDie, *ConditionDieMut, string]
//
// The names depend on the keys a list is keyed by, which are the list
// field's, not the element type's. They are given for an element type that
// a list-map field of the package lists, and for one marked with the keys it
// is listed by elsewhere, as a metav1.Condition is listed by its type:
//
//	// +die:listMap:keys=type
//	type _ = Condition
//
// An element type listed by more than one set of keys has no names, as
// which keys a name means would be ambiguous. A list-map field of an element
// type of another package uses that package's names when the list is keyed
// by one field; a list keyed by several has a key struct of its own
// package's, so it doesn't.

// listMapAlias is the names of the list maps of an element type.
type listMapAlias struct {
	elem types.Type
	keys []string
	// ambiguous is set when the element type is listed by other keys too
	ambiguous bool
}

// listMapAliases returns the list maps to name, by element die name.
func (pg *pkgGen) listMapAliases() map[string]*listMapAlias {
	out := map[string]*listMapAlias{}
	add := func(elem types.Type, keys []string) {
		path, name, ok := pg.dieOf(elem)
		if !ok || path != pg.p.PkgPath {
			return
		}
		if a, ok := out[name]; ok {
			if !slices.Equal(a.keys, keys) {
				a.ambiguous = true
			}
			return
		}
		out[name] = &listMapAlias{elem: elem, keys: keys}
	}
	for _, tg := range pg.targets {
		if m, ok := findDieMarker(tg.markers, "die:listMap"); ok {
			add(tg.named, m.args["keys"])
		}
		if _, ok := pg.listItem(tg); ok {
			continue
		}
		for _, wf := range pg.flatten(tg.named, tg.st) {
			if pg.kind(wf) == fListMap {
				add(types.Unalias(wf.v.Type()).Underlying().(*types.Slice).Elem(), wf.spec.listMapKey)
			}
		}
	}
	return out
}

// listMapAlias returns the name of the list map of an element type, keyed
// by keys, if the element's die package names it, qualified if that is
// another package; its mutable list map's is the name with Mut.
func (pg *pkgGen) listMapAlias(elem types.Type, keys []string) (string, bool) {
	path, name, ok := pg.dieOf(elem)
	if !ok {
		return "", false
	}
	aliases, p := pg.aliases, pg.p
	if path != pg.p.PkgPath {
		// a list keyed by several fields has a key struct of its package's
		// own, not the element's package's
		if len(keys) != 1 {
			return "", false
		}
		m := pg.g.listMapsOf(path)
		if m == nil {
			return "", false
		}
		aliases, p = pg.g.listMaps[path], m.p
	}
	a, ok := aliases[name]
	if !ok || a.ambiguous || !slices.Equal(a.keys, keys) || declared(p, name+"ListMap") || declared(p, name+"ListMapMut") {
		return "", false
	}
	if path != pg.p.PkgPath {
		return pg.f.use(path) + "." + name + "ListMap", true
	}
	return name + "ListMap", true
}

// listMapsOf works out the list maps another die package names, returning
// its metadata, or nil if it isn't a die package.
func (g *generator) listMapsOf(path string) *pkgMeta {
	m, err := g.metaFor(path)
	if err != nil {
		return nil
	}
	if _, ok := g.listMaps[path]; !ok {
		other := &pkgGen{pkgMeta: m, g: g, f: newFile(m.p.Name), keys: map[string][]wfield{}, specs: map[*types.Named]*typeSpec{}, dies: map[string]bool{}}
		g.listMaps[path] = other.listMapAliases()
	}
	return m
}

// declared reports whether a package declares a name outside its generated
// dies.
func declared(p *packages.Package, name string) bool {
	obj := p.Types.Scope().Lookup(name)
	return obj != nil && filepath.Base(p.Fset.File(obj.Pos()).Name()) != diesFile
}

// emitListMapAliases writes the names of the package's list maps.
func (pg *pkgGen) emitListMapAliases() {
	names := make([]string, 0, len(pg.aliases))
	for name := range pg.aliases {
		names = append(names, name)
	}
	slices.Sort(names)
	f := pg.f
	for _, name := range names {
		a := pg.aliases[name]
		alias, ok := pg.listMapAlias(a.elem, a.keys)
		if !ok {
			continue
		}
		k, _, _ := pg.listMapKey(a.elem, a.keys)
		elem, keys := pg.typeExpr(a.elem), strings.Join(a.keys, ", ")
		f.p("// %s is the list map of an immutable die's %s elements, keyed by %s.", alias, elem, keys)
		f.p("type %s = die.ListMap[%s, %s, %s]", alias, elem, pg.dieExpr(a.elem), k)
		f.p("")
		f.p("// %sMut is the list map of a mutable die's %s elements, keyed by %s.", alias, elem, keys)
		f.p("type %sMut = die.ListMapMut[%s, %s, %s, %s]", alias, elem, pg.dieExpr(a.elem), pg.dieMutExpr(a.elem), k)
		f.p("")
	}
}
