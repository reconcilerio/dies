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
	"go/token"
	"go/types"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

// aliases is the first phase: for each die package with a +die:types marker,
// alias the upstream package's types and constants, and re-export its
// functions and variables, marking the types that need dies.
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
			if fd, ok := decl.(*ast.FuncDecl); ok && fd.Recv == nil {
				names[fd.Name.Name] = true
			}
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

	docs := valueDocs(up)
	f.p("// The types, constants, functions and variables of [%s], so that", up.PkgPath)
	f.p("// importing this package provides both the types and their dies. Each type is")
	f.p("// followed by its values. A function is called through one of the same")
	f.p("// signature, and a variable is copied as it is when this package is")
	f.p("// initialized. The scheme registration isn't re-exported: this package's")
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
		case *types.Var:
			if !obj.Exported() || schemeRegistration[name] {
				continue
			}
			f.p("// %s is [%s.%s].", name, pkg, name)
			f.value("var %s = %s.%s", []any{name, pkg, name}, deprecation(docs[obj]))
			f.p("")
		case *types.Func:
			if !obj.Exported() || schemeRegistration[name] {
				continue
			}
			if f.forward(obj, pkg, deprecation(docs[obj])) {
				f.p("")
			}
		}
	}
	return f.write(filepath.Dir(p.GoFiles[0]), typesFile)
}

// schemeRegistration is the names of an API package's scheme registration,
// which a die package declares for its dies, in place of the upstream
// package's.
var schemeRegistration = map[string]bool{"SchemeGroupVersion": true, "SchemeBuilder": true, "AddToScheme": true}

// valueDocs returns the doc lines of a package's functions and variables.
func valueDocs(p *packages.Package) map[types.Object][]string {
	docs := map[types.Object][]string{}
	for _, f := range p.Syntax {
		for _, decl := range f.Decls {
			switch decl := decl.(type) {
			case *ast.FuncDecl:
				if decl.Recv == nil {
					docs[p.TypesInfo.Defs[decl.Name]] = commentLines(decl.Doc)
				}
			case *ast.GenDecl:
				if decl.Tok != token.VAR {
					continue
				}
				for _, spec := range decl.Specs {
					vs := spec.(*ast.ValueSpec)
					doc := vs.Doc
					if doc == nil && len(decl.Specs) == 1 {
						doc = decl.Doc
					}
					for _, n := range vs.Names {
						docs[p.TypesInfo.Defs[n]] = commentLines(doc)
					}
				}
			}
		}
	}
	return docs
}

// forward emits a function that calls an upstream function of the same
// signature, as Go has no alias of a function, and reports whether it could:
// one whose signature names what can't be named here, such as an unexported
// type, isn't re-exported.
func (f *file) forward(fn *types.Func, pkg string, deprecated []string) bool {
	sig := fn.Signature()
	if !nameable(sig, map[types.Type]bool{}) {
		return false
	}
	qualify := func(p *types.Package) string { return f.use(p.Path()) }
	typeString := func(t types.Type) string { return types.TypeString(t, qualify) }
	var tparams, targs []string
	for i := range sig.TypeParams().Len() {
		tp := sig.TypeParams().At(i)
		tparams = append(tparams, tp.Obj().Name()+" "+typeString(tp.Constraint()))
		targs = append(targs, tp.Obj().Name())
	}
	var params, results []string
	for i := range sig.Params().Len() {
		t := sig.Params().At(i).Type()
		if sig.Variadic() && i == sig.Params().Len()-1 {
			params = append(params, "..."+typeString(t.(*types.Slice).Elem()))
		} else {
			params = append(params, typeString(t))
		}
	}
	for i := range sig.Results().Len() {
		results = append(results, typeString(sig.Results().At(i).Type()))
	}
	// parameters keep their upstream names, unless one would hide a package
	// or type parameter the function names
	taken := map[string]bool{}
	for a := range f.aliases {
		taken[a] = true
	}
	for _, n := range targs {
		taken[n] = true
	}
	var args []string
	for i := range sig.Params().Len() {
		n := sig.Params().At(i).Name()
		if n == "" || n == "_" || taken[n] {
			n = fmt.Sprintf("p%d", i)
		}
		for taken[n] {
			n += "_"
		}
		taken[n] = true
		params[i] = n + " " + params[i]
		args = append(args, n)
	}
	if sig.Variadic() {
		args[len(args)-1] += "..."
	}
	decl := fn.Name()
	call := pkg + "." + fn.Name()
	if len(tparams) > 0 {
		decl += "[" + strings.Join(tparams, ", ") + "]"
		call += "[" + strings.Join(targs, ", ") + "]"
	}
	decl += "(" + strings.Join(params, ", ") + ")"
	switch len(results) {
	case 0:
	case 1:
		decl += " " + results[0]
	default:
		decl += " (" + strings.Join(results, ", ") + ")"
	}
	call += "(" + strings.Join(args, ", ") + ")"
	if len(results) > 0 {
		call = "return " + call
	}
	f.p("// %s calls [%s.%s].", fn.Name(), pkg, fn.Name())
	if deprecated != nil {
		f.p("//")
		for _, l := range deprecated {
			f.p("// %s", l)
		}
	}
	f.p("func %s {", decl)
	if deprecated != nil {
		f.p("	//lint:ignore SA1019 a function is called through, deprecated or not")
		f.p("	%s //nolint:staticcheck // a function is called through, deprecated or not", call)
	} else {
		f.p("	%s", call)
	}
	f.p("}")
	return true
}

// value emits a declaration of a value, below the doc lines already emitted
// for it, deprecated as what it is of is, and exempt from the linters that
// would report its use, as an alias is.
func (f *file) value(format string, args []any, deprecated []string) {
	if deprecated == nil {
		f.p(format, args...)
		return
	}
	f.p("//")
	for _, l := range deprecated {
		f.p("// %s", l)
	}
	f.p("//")
	f.p("//lint:ignore SA1019 a variable is of its upstream variable, deprecated or not")
	f.p(format+" //nolint:staticcheck // a variable is of its upstream variable, deprecated or not", args...)
}

// nameable reports whether a type can be written in another package: every
// type it names is exported, and not from an internal package.
func nameable(t types.Type, seen map[types.Type]bool) bool {
	if seen[t] {
		return true
	}
	seen[t] = true
	named := func(obj *types.TypeName) bool {
		if obj.Pkg() == nil {
			return true
		}
		p := obj.Pkg().Path()
		internal := p == "internal" || strings.HasPrefix(p, "internal/") || strings.HasSuffix(p, "/internal") || strings.Contains(p, "/internal/")
		return obj.Exported() && !internal
	}
	tuple := func(tu *types.Tuple) bool {
		for i := range tu.Len() {
			if !nameable(tu.At(i).Type(), seen) {
				return false
			}
		}
		return true
	}
	switch t := t.(type) {
	case *types.Basic, *types.TypeParam:
		return true
	case *types.Alias:
		if !named(t.Obj()) {
			return false
		}
		for i := range t.TypeArgs().Len() {
			if !nameable(t.TypeArgs().At(i), seen) {
				return false
			}
		}
		return true
	case *types.Named:
		if !named(t.Obj()) {
			return false
		}
		for i := range t.TypeArgs().Len() {
			if !nameable(t.TypeArgs().At(i), seen) {
				return false
			}
		}
		return true
	case *types.Pointer:
		return nameable(t.Elem(), seen)
	case *types.Slice:
		return nameable(t.Elem(), seen)
	case *types.Array:
		return nameable(t.Elem(), seen)
	case *types.Chan:
		return nameable(t.Elem(), seen)
	case *types.Map:
		return nameable(t.Key(), seen) && nameable(t.Elem(), seen)
	case *types.Signature:
		for i := range t.TypeParams().Len() {
			if !nameable(t.TypeParams().At(i).Constraint(), seen) {
				return false
			}
		}
		return tuple(t.Params()) && tuple(t.Results())
	case *types.Struct:
		for i := range t.NumFields() {
			if !t.Field(i).Exported() || !nameable(t.Field(i).Type(), seen) {
				return false
			}
		}
		return true
	case *types.Interface:
		for i := range t.NumExplicitMethods() {
			if !t.ExplicitMethod(i).Exported() || !nameable(t.ExplicitMethod(i).Type(), seen) {
				return false
			}
		}
		for i := range t.NumEmbeddeds() {
			if !nameable(t.EmbeddedType(i), seen) {
				return false
			}
		}
		return true
	case *types.Union:
		for i := range t.Len() {
			if !nameable(t.Term(i).Type(), seen) {
				return false
			}
		}
		return true
	}
	return false
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
