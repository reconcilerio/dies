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
	"errors"
	"fmt"
	"go/ast"
	"go/types"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

// aliases is the first phase: for each die package with a +die:types marker,
// alias the upstream package's types and constants, marking the types that
// need dies.
func (g *generator) aliases(patterns []string) error {
	// only the die packages' docs are needed; their generated files may be
	// stale or missing, so they are not type checked
	progress.logf("finding die packages in %s", strings.Join(patterns, " "))
	pkgs, err := packages.Load(&packages.Config{Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax}, patterns...)
	if err != nil {
		return err
	}
	type aliasing struct {
		pkg      *packages.Package
		upstream string
	}
	var todo []aliasing
	var upstreams []string
	var clones []clone
	for _, p := range pkgs {
		m, mode, ok, err := typesMarker(packageDoc(p))
		if err != nil {
			return fmt.Errorf("%s: %w", p.PkgPath, err)
		}
		switch {
		case !ok:
		case mode == "alias":
			todo = append(todo, aliasing{pkg: p, upstream: m.arg("package")})
			upstreams = append(upstreams, m.arg("package"))
		case mode == "clone":
			upstream, version, _ := strings.Cut(m.arg("package"), "@")
			from := m.arg("from")
			if version != "" && from != "" {
				return fmt.Errorf("%s: a clone takes a version with its package or a from module, not both", p.PkgPath)
			}
			if from == "" {
				from = "."
			}
			clones = append(clones, clone{pkg: p, upstream: upstream, version: version, from: filepath.Join(filepath.Dir(p.GoFiles[0]), from)})
		}
	}
	if len(clones) > 0 {
		ups, err := g.loadClones(clones)
		if err != nil {
			return err
		}
		for _, c := range clones {
			if err := g.emitClone(c.pkg, ups[c.upstream]); err != nil {
				return err
			}
		}
	}
	if len(todo) == 0 {
		return nil
	}
	// type checked from source, as with the dies, rather than from export data
	progress.logf("loading %d %s to alias", len(upstreams), plural(len(upstreams), "package", "packages"))
	ups, err := packages.Load(&packages.Config{Mode: packages.NeedName | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedImports | packages.NeedDeps}, upstreams...)
	if err != nil {
		return err
	}
	if packages.PrintErrors(ups) > 0 {
		return errors.New("failed to load upstream packages")
	}
	byPath := map[string]*packages.Package{}
	for _, u := range ups {
		g.index(u)
		byPath[u.PkgPath] = u
	}
	for _, a := range todo {
		if err := g.emitAliases(a.pkg, byPath[a.upstream]); err != nil {
			return err
		}
	}
	return nil
}

// declaredByHand returns the names a die package declares outside its
// generated files. Phase 1 does not alias them, so that a hand-written alias
// can override the generated one, with markers of its own.
func declaredByHand(p *packages.Package) map[string]bool {
	names := map[string]bool{}
	for _, f := range p.Syntax {
		switch filepath.Base(p.Fset.File(f.Pos()).Name()) {
		case typesFile, diesFile, acFile:
			continue
		}
		for _, decl := range f.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, spec := range gd.Specs {
				switch spec := spec.(type) {
				case *ast.TypeSpec:
					names[spec.Name.Name] = true
				case *ast.ValueSpec:
					for _, n := range spec.Names {
						names[n.Name] = true
					}
				}
			}
		}
	}
	return names
}

func (g *generator) emitAliases(p, up *packages.Package) error {
	declared := declaredByHand(p)
	f := newFile(p.Name)
	pkg := f.use(up.PkgPath)
	scope := up.Types.Scope()

	apiVersion := ""
	if !strings.HasPrefix(up.PkgPath, "k8s.io/apimachinery/") {
		// apimachinery's types with TypeMeta, such as Status, are shared by
		// every group, so they are not objects of their own
		apiVersion = groupVersion(up)
	}

	// constants are listed with their type, when it is aliased here
	valuesOf := map[*types.TypeName][]string{}
	var aliasedTypes = map[*types.TypeName]bool{}
	for _, name := range scope.Names() {
		if obj, ok := scope.Lookup(name).(*types.TypeName); ok && obj.Exported() && !declared[name] {
			aliasedTypes[obj] = true
		}
	}
	ownType := func(c *types.Const) *types.TypeName {
		if n, ok := types.Unalias(c.Type()).(*types.Named); ok && aliasedTypes[n.Obj()] {
			return n.Obj()
		}
		return nil
	}
	for _, name := range scope.Names() {
		if c, ok := scope.Lookup(name).(*types.Const); ok && c.Exported() && !declared[name] {
			if t := ownType(c); t != nil {
				valuesOf[t] = append(valuesOf[t], name)
			}
		}
	}

	f.p("// The types and constants of [%s], so that importing this package provides", up.PkgPath)
	f.p("// both the types and their dies. Each type is followed by its values.")
	f.p("// Functions and variables are deliberately not re-exported: this package's")
	f.p("// [AddToScheme] registers its dies, not the upstream types.")
	f.p("")
	for _, name := range scope.Names() {
		if declared[name] {
			continue
		}
		switch obj := scope.Lookup(name).(type) {
		case *types.TypeName:
			if !obj.Exported() {
				continue
			}
			if scope.Lookup(name+"Die") != nil || scope.Lookup(name+"Blank") != nil {
				return fmt.Errorf("%s: generated names for %s collide with upstream names", p.PkgPath, name)
			}
			if n, st, ok := structOf(obj.Type()); ok && n.Obj() == obj {
				if apiVersion != "" && embedsTypeMeta(st) {
					f.p("// +die:object:apiVersion=%s,kind=%s", apiVersion, name)
				} else {
					f.p("// +die")
				}
				// the upstream markers that shape the die, in canonical form
				var order []string
				for _, wf := range g.directFields(st) {
					order = append(order, wf.wire)
				}
				for _, l := range g.upstreamSpec(n, st).render(order) {
					f.p("// %s", l)
				}
			}
			f.alias(name, pkg+"."+name, deprecation(g.typeDocs[obj]))
			if values := valuesOf[obj]; len(values) > 0 {
				f.p("")
				f.p("const (")
				for _, v := range values {
					f.p("\t%s = %s.%s", v, pkg, v)
				}
				f.p(")")
			}
			f.p("")
		case *types.Const:
			if obj.Exported() && ownType(obj) == nil {
				f.p("const %s = %s.%s", name, pkg, name)
				f.p("")
			}
		}
	}
	return f.write(filepath.Dir(p.GoFiles[0]), typesFile)
}

// deprecation returns the paragraph of a doc comment that deprecates what it
// documents, without markers, or nil. As Go's convention has it, and as
// staticcheck reads it, the paragraph starts with "Deprecated: ".
func deprecation(doc []string) []string {
	for i, l := range doc {
		if !strings.HasPrefix(l, "Deprecated: ") || i > 0 && strings.TrimSpace(doc[i-1]) != "" {
			continue
		}
		var lines []string
		for _, l := range doc[i:] {
			if strings.TrimSpace(l) == "" {
				break
			}
			if !strings.HasPrefix(l, "+") {
				lines = append(lines, l)
			}
		}
		return lines
	}
	return nil
}

// alias emits an alias of a type, below the doc lines already emitted for
// it. An alias of a deprecated type is deprecated too, so that code naming
// it is told, and its own use of the type is exempt from the linters that
// would report it: golangci-lint's directive, and staticcheck's.
func (f *file) alias(name, target string, deprecated []string) {
	if deprecated == nil {
		f.p("type %s = %s", name, target)
		return
	}
	f.p("//")
	for _, l := range deprecated {
		f.p("// %s", l)
	}
	f.p("//")
	f.p("//lint:ignore SA1019 an alias is of its type, deprecated or not")
	f.p("type %s = %s //nolint:staticcheck // an alias is of its type, deprecated or not", name, target)
}

func embedsTypeMeta(st *types.Struct) bool {
	for i := range st.NumFields() {
		if f := st.Field(i); f.Embedded() && f.Name() == "TypeMeta" {
			return true
		}
	}
	return false
}
