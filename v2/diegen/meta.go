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
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"golang.org/x/tools/go/packages"
)

// dies is the second phase: generate a die for every type marked +die.
func (g *generator) dies(patterns []string) error {
	progress.logf("type checking die packages")
	pkgs, err := g.load(patterns)
	if err != nil {
		return err
	}
	for _, p := range pkgs {
		if g.meta[p.PkgPath] == nil {
			continue
		}
		if n := len(g.meta[p.PkgPath].targets); n > 0 {
			progress.generated(p.PkgPath, n)
		}
		if err := g.diesFor(p); err != nil {
			return err
		}
	}
	if g.ignored > 0 {
		fmt.Fprintf(os.Stderr, "diegen: ignored %d type errors outside die metadata, such as hand-written code that refers to dies not generated yet\n", g.ignored)
	}
	return nil
}

// load type checks packages for phase 2 and reads the metadata of every die
// package among them and their dependencies.
//
// Generated dies may be stale or not generated yet, and hand-written code may
// refer to dies that don't exist, so type errors in die packages are
// tolerated, unless they are in the die metadata: a +die declaration or a
// type _ metadata declaration. Parse errors in a die package, and any error
// in another package, mean the metadata cannot be read, and fail.
func (g *generator) load(patterns []string) ([]*packages.Package, error) {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedImports | packages.NeedDeps,
	}, patterns...)
	if err != nil {
		return nil, err
	}
	var errs []error
	packages.Visit(pkgs, nil, func(p *packages.Package) {
		g.index(p)
	})
	packages.Visit(pkgs, nil, func(p *packages.Package) {
		if !isDiePackage(p) {
			for _, e := range p.Errors {
				// a missing package imported by generated dies is reported
				// against that package
				if inGeneratedDies(e) {
					g.ignored++
					continue
				}
				errs = append(errs, e)
			}
			return
		}
		if _, ok := g.meta[p.PkgPath]; ok {
			return
		}
		// type errors in the metadata are the cause of any problem collecting
		// it, so they are reported first
		m, err := g.collect(p)
		errs = append(errs, g.checkErrors(m)...)
		if err != nil {
			errs = append(errs, err)
			return
		}
		g.meta[p.PkgPath] = m
	})
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return pkgs, nil
}

// metaFor returns the metadata of a die package, loading it if it was not
// loaded with the packages being generated.
func (g *generator) metaFor(path string) (*pkgMeta, error) {
	if m, ok := g.meta[path]; ok {
		return m, nil
	}
	if _, err := g.load([]string{path}); err != nil {
		return nil, err
	}
	m, ok := g.meta[path]
	if !ok {
		return nil, fmt.Errorf("%s is not a die package", path)
	}
	return m, nil
}

// isDiePackage reports whether a package has die metadata: +die markers in
// the comments above its package clause or on its types.
func isDiePackage(p *packages.Package) bool {
	if len(dieMarkers(packageDoc(p))) > 0 {
		return true
	}
	found := false
	eachTypeSpec(p, func(_ *ast.GenDecl, _ *ast.TypeSpec, doc []string, unread []*ast.CommentGroup) {
		found = found || len(dieMarkers(doc)) > 0 || len(unreadDieMarkers(unread)) > 0
	})
	return found
}

// unreadDieMarkers returns the comments with +die markers among those that
// are not read as a type's markers.
func unreadDieMarkers(unread []*ast.CommentGroup) []*ast.CommentGroup {
	var out []*ast.CommentGroup
	for _, cg := range unread {
		if len(dieMarkers(commentLines(cg))) > 0 {
			out = append(out, cg)
		}
	}
	return out
}

func eachTypeSpec(p *packages.Package, fn func(gd *ast.GenDecl, ts *ast.TypeSpec, doc []string, unread []*ast.CommentGroup)) {
	for _, file := range p.Syntax {
		for _, d := range file.Decls {
			gd, ok := d.(*ast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}
			for _, spec := range gd.Specs {
				ts := spec.(*ast.TypeSpec)
				doc, unread := typeDoc(file, gd, ts)
				fn(gd, ts, doc, unread)
			}
		}
	}
}

// pkgMeta is the metadata of a die package.
type pkgMeta struct {
	p        *packages.Package
	mappings map[string]string
	// upstream is the package a die package aliases, if any
	upstream string
	targets  []*target
	byNamed  map[*types.Named]*target
	// byName has the spec of each target, by the qualified name of the type
	// it is for, so other loads can find it
	byName map[string]*typeSpec
	// spans are the positions of the metadata declarations
	spans [][2]token.Pos
}

// target is a type marked +die: an alias of an upstream type, or a type
// declared in the die package.
type target struct {
	obj   *types.TypeName
	named *types.Named
	st    *types.Struct
	// doc is the declaration's doc, followed by the aliased type's
	doc []string
	// spec is the type's upstream markers, normalized, with its +die
	// markers and any extra metadata applied over them
	spec *typeSpec
}

func (t *target) name() string { return t.obj.Name() }

func qualified(n *types.Named) string { return n.Obj().Pkg().Path() + "." + n.Obj().Name() }

func (g *generator) collect(p *packages.Package) (*pkgMeta, error) {
	m := &pkgMeta{p: p, mappings: map[string]string{}, byNamed: map[*types.Named]*target{}, byName: map[string]*typeSpec{}}
	for _, dm := range dieMarkers(packageDoc(p)) {
		switch dm.name {
		case "die:import":
			m.mappings[dm.arg("package")] = dm.arg("as")
		case "die:types":
			if mode := dm.arg("mode"); mode == "" || mode == "alias" {
				m.mappings[dm.arg("package")] = p.PkgPath
				m.upstream, _, _ = strings.Cut(dm.arg("package"), "@")
			}
		}
	}

	// metadata declared with the blank identifier, type _ = corev1.X, adds
	// to a type's markers without declaring it, so the alias itself stays
	// generated
	extra := map[*types.Named][]string{}
	type decl struct {
		obj *types.TypeName
		doc []string
	}
	var decls []decl
	var errs []error
	eachTypeSpec(p, func(gd *ast.GenDecl, ts *ast.TypeSpec, doc []string, unread []*ast.CommentGroup) {
		// a +die marker that would be ignored is a mistake, not a choice
		for _, cg := range unreadDieMarkers(unread) {
			errs = append(errs, fmt.Errorf("%s: +die markers here are not read for %s: a type's markers go in its doc, or in one comment of markers only, separated from the doc by a blank line", p.Fset.Position(cg.Pos()), ts.Name.Name))
		}
		start := ts.Pos()
		if gd.Doc != nil && len(gd.Specs) == 1 {
			start = gd.Pos()
		}
		if ts.Name.Name == "_" {
			if len(dieMarkers(doc)) == 0 {
				return
			}
			m.spans = append(m.spans, [2]token.Pos{start, ts.End()})
			named, _, ok := structOf(p.TypesInfo.TypeOf(ts.Type))
			if !ok || !ts.Assign.IsValid() {
				errs = append(errs, fmt.Errorf("%s: metadata must be declared as type _ = <struct type>", p.Fset.Position(ts.Pos())))
				return
			}
			if marked(doc) {
				errs = append(errs, fmt.Errorf("%s: type _ = %s only adds metadata; mark the alias with +die instead", p.Fset.Position(ts.Pos()), named.Obj().Name()))
				return
			}
			extra[named] = append(extra[named], doc...)
			return
		}
		obj, ok := p.TypesInfo.Defs[ts.Name].(*types.TypeName)
		if ok && marked(g.typeDocs[obj]) {
			m.spans = append(m.spans, [2]token.Pos{start, ts.End()})
			decls = append(decls, decl{obj: obj, doc: g.typeDocs[obj]})
		}
	})
	for _, d := range decls {
		named, st, ok := structOf(d.obj.Type())
		if !ok {
			errs = append(errs, fmt.Errorf("%s is marked +die but is not a struct without custom JSON encoding", d.obj))
			continue
		}
		// a type's markers, then any extra metadata, apply over the upstream
		// markers, so only what upstream lacks needs stating and new
		// upstream fields keep their markers
		spec := dieSpec(g.upstreamSpec(named, st), append(slices.Clone(d.doc), extra[named]...))
		if err := g.checkFields(p, d.obj, st, spec); err != nil {
			errs = append(errs, err)
			continue
		}
		doc := d.doc
		if d.obj.IsAlias() {
			doc = append(slices.Clone(doc), g.typeDocs[named.Obj()]...)
		}
		t := &target{obj: d.obj, named: named, st: st, doc: doc, spec: spec}
		m.targets = append(m.targets, t)
		m.byNamed[named] = t
		m.byName[qualified(named)] = spec
		delete(extra, named)
	}
	for named := range extra {
		errs = append(errs, fmt.Errorf("%s: metadata for %s, which has no die in this package", p.PkgPath, named.Obj().Name()))
	}
	slices.SortFunc(m.targets, func(a, b *target) int { return strings.Compare(a.name(), b.name()) })
	return m, errors.Join(errs...)
}

// checkErrors returns the errors that mean a die package's metadata cannot
// be read, and counts the others.
func (g *generator) checkErrors(m *pkgMeta) []error {
	var errs []error
	for _, e := range m.p.Errors {
		if e.Kind == packages.TypeError {
			continue
		}
		// generated dies are never metadata: an import they can't resolve,
		// say, is fixed by generating them again
		if inGeneratedDies(e) {
			g.ignored++
			continue
		}
		errs = append(errs, e)
	}
	for _, te := range m.p.TypeErrors {
		inMetadata := slices.ContainsFunc(m.spans, func(s [2]token.Pos) bool { return te.Pos >= s[0] && te.Pos < s[1] })
		if inMetadata {
			errs = append(errs, te)
		} else {
			g.ignored++
		}
	}
	return errs
}

func inGeneratedDies(e packages.Error) bool {
	file, _, _ := strings.Cut(e.Pos, ":")
	return filepath.Base(file) == diesFile
}

// checkFields rejects markers that name a field the type does not have,
// such as a union member that upstream removed or renamed.
func (g *generator) checkFields(p *packages.Package, obj *types.TypeName, st *types.Struct, spec *typeSpec) error {
	wires := map[string]bool{}
	for _, wf := range g.directFields(st) {
		wires[wf.wire] = true
	}
	var unknown []string
	for wire, fs := range spec.fields {
		if !wires[wire] {
			unknown = append(unknown, wire)
		}
		if fs.def != "" && !json.Valid([]byte(fs.def)) {
			return fmt.Errorf("%s.%s: the default for %s is not JSON: %s", p.PkgPath, obj.Name(), wire, fs.def)
		}
	}
	if len(unknown) > 0 {
		slices.Sort(unknown)
		return fmt.Errorf("%s.%s: markers name fields the type does not have: %v", p.PkgPath, obj.Name(), unknown)
	}
	return nil
}
