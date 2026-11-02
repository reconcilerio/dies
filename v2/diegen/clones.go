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
	"go/constant"
	"go/types"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/tools/go/packages"
)

// typesFile has a die package's types: aliases of the upstream package's, or
// clones of them.
const typesFile = "zz_generated.dies.types.go"

// Cloning copies an upstream package's types into a die package, so that
// nothing depends on the upstream package and its transitive dependencies.
// Types, their docs and markers, and constants are cloned; methods,
// functions and variables are not.
//
// The upstream package is resolved at a version given with the package, in a
// temporary module:
//
//	// +die:types:package=carvel.dev/vendir/pkg/vendir/versions/v1alpha1@v0.46.2,mode=clone
//
// or from the module in a from directory, relative to the die package, when
// its versions are managed as a module of their own:
//
//	// +die:types:package=carvel.dev/kapp-controller/pkg/apis/packaging/v1alpha1,mode=clone,from=../../sources
//
// Either way the upstream package is only needed when generating.
//
// A cloned type that refers to another third-party package needs that
// package cloned too, and imported with +die:import. Types from Kubernetes
// and the standard library are referenced where they are.

// clone describes a die package that clones an upstream package.
type clone struct {
	pkg      *packages.Package
	upstream string
	version  string
	from     string
}

// loadClones loads the upstream packages to clone, from the module each is
// resolved in.
func (g *generator) loadClones(clones []clone) (map[string]*packages.Package, error) {
	byFrom := map[string][]string{}
	// versionedDir is the temporary module the versioned clones resolve in
	var versionedDir string
	var versioned []clone
	for _, c := range clones {
		if c.version != "" {
			versioned = append(versioned, c)
			continue
		}
		byFrom[c.from] = append(byFrom[c.from], c.upstream)
	}
	if len(versioned) > 0 {
		dir, err := versionedModule(versioned)
		if dir != "" {
			defer os.RemoveAll(dir)
		}
		if err != nil {
			return nil, err
		}
		versionedDir = dir
		for _, c := range versioned {
			byFrom[dir] = append(byFrom[dir], c.upstream)
		}
	}
	ups := map[string]*packages.Package{}
	for from, paths := range byFrom {
		where := "from the module in " + progress.rel(from)
		if from == versionedDir {
			where = "at the versions asked for"
		}
		progress.logf("loading %d %s to clone, %s", len(paths), plural(len(paths), "package", "packages"), where)
		pkgs, err := packages.Load(&packages.Config{
			Mode: packages.NeedName | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedImports | packages.NeedDeps,
			Dir:  from,
			// the from module is a module of its own, outside any workspace
			Env: append(os.Environ(), "GOWORK=off"),
		}, paths...)
		if err != nil {
			return nil, err
		}
		if packages.PrintErrors(pkgs) > 0 {
			return nil, fmt.Errorf("failed to load packages to clone from %s", from)
		}
		for _, p := range pkgs {
			g.index(p)
			ups[p.PkgPath] = p
		}
	}
	return ups, nil
}

// versionedModule makes a temporary module that requires each clone's
// package at its version. One module for all of them keeps the clones
// consistent with each other: a clone's view of another cloned package is
// the version that is cloned. Kubernetes modules are kept at the versions the
// die packages use, as clones refer to Kubernetes types directly.
func versionedModule(clones []clone) (string, error) {
	dir, err := os.MkdirTemp("", "diegen-clones-")
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module diegen.local/clones\n\ngo 1.27\n"), 0o644); err != nil {
		return dir, err
	}
	progress.logf("resolving %d versioned %s to clone", len(clones), plural(len(clones), "package", "packages"))
	gets := []string{"get"}
	for _, c := range clones {
		gets = append(gets, c.upstream+"@"+c.version)
	}
	for _, m := range []string{"k8s.io/api", "k8s.io/apimachinery"} {
		if v := moduleVersion(filepath.Dir(clones[0].pkg.GoFiles[0]), m); v != "" {
			gets = append(gets, m+"@"+v)
		}
	}
	if out, err := goCmd(dir, true, gets...); err != nil {
		return dir, fmt.Errorf("resolving packages to clone: %w\n%s", err, out)
	}
	// another clone's requirements may have selected a newer version than
	// the one asked for, which would silently clone something else
	for _, c := range clones {
		out, err := goCmd(dir, true, "list", "-f", "{{.Module.Version}}", c.upstream)
		if err != nil {
			return dir, fmt.Errorf("resolving %s: %w\n%s", c.upstream, err, out)
		}
		if got := strings.TrimSpace(out); got != c.version {
			return dir, fmt.Errorf("%s@%s was asked for, but %s is selected by the requirements of the other clones; ask for %s, or clone from a module with from=", c.upstream, c.version, got, got)
		}
	}
	return dir, nil
}

// moduleVersion is the version of a module in the build of the package in
// dir, if it is required.
func moduleVersion(dir, module string) string {
	out, err := goCmd(dir, false, "list", "-m", "-f", "{{.Version}}", module)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(out)
}

// goCmd runs the go command in dir. A standalone module is run outside any
// workspace, and may update its requirements.
func goCmd(dir string, standalone bool, args ...string) (string, error) {
	cmd := exec.Command("go", args...)
	cmd.Dir = dir
	if standalone {
		cmd.Env = append(os.Environ(), "GOWORK=off", "GOFLAGS=-mod=mod")
	}
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// codegenMarkers are upstream markers for other generators, which are not
// copied: a clone has no deepcopy methods, clients or protobuf encoding.
var codegenMarkers = []string{"+genclient", "+k8s:deepcopy-gen", "+k8s:openapi-gen", "+k8s:protobuf-gen", "+k8s:prerelease-lifecycle-gen", "+k8s:conversion-gen", "+k8s:defaulter-gen", "+protobuf"}

func cloneDoc(lines []string) []string {
	var out []string
	for _, l := range lines {
		if !slices.ContainsFunc(codegenMarkers, func(m string) bool { return strings.HasPrefix(l, m) }) {
			out = append(out, l)
		}
	}
	return out
}

func (g *generator) emitClone(p, up *packages.Package) error {
	declared := declaredByHand(p)
	imports := map[string]string{}
	for _, m := range dieMarkers(packageDoc(p)) {
		if m.name == "die:import" {
			imports[m.arg("package")] = m.arg("as")
		}
	}
	f := newFile(p.Name)
	var errs []error
	expr := func(t types.Type) string {
		s, err := cloneTypeExpr(f, up.PkgPath, imports, t)
		if err != nil {
			errs = append(errs, err)
		}
		return s
	}
	scope := up.Types.Scope()
	apiVersion := groupVersion(up)
	cloned := typesToClone(up)
	hand := methodsByHand(p)

	f.p("// The types and constants of %s, cloned so that nothing depends on it", up.PkgPath)
	f.p("// or its dependencies. Methods, functions and variables are not cloned.")
	f.p("")
	for _, name := range scope.Names() {
		if declared[name] {
			continue
		}
		switch obj := scope.Lookup(name).(type) {
		case *types.TypeName:
			if !cloned[obj] {
				continue
			}
			if obj.IsAlias() {
				f.doc(cloneDoc(g.typeDocs[obj]))
				f.p("type %s = %s", name, expr(obj.Type()))
				f.p("")
				continue
			}
			n := obj.Type().(*types.Named)
			// a type's JSON methods define its encoding, and are not cloned,
			// so the die package must provide them or declare the type itself
			custom := customJSON(n)
			if custom && !(hand[name]["MarshalJSON"] && hand[name]["UnmarshalJSON"]) {
				errs = append(errs, fmt.Errorf("%s.%s has custom JSON encoding, which is not cloned; in %s, either declare %s by hand as a type with the same encoding (such as a string), or write MarshalJSON and UnmarshalJSON for the cloned type", up.PkgPath, name, p.PkgPath, name))
				continue
			}
			doc := cloneDoc(g.typeDocs[obj])
			st, isStruct := n.Underlying().(*types.Struct)
			if isStruct && !custom && obj.Exported() {
				if len(doc) > 0 {
					doc = append(doc, "")
				}
				if apiVersion != "" && embedsTypeMeta(st) {
					doc = append(doc, fmt.Sprintf("+die:object:apiVersion=%s,kind=%s", apiVersion, name))
				} else {
					doc = append(doc, "+die")
				}
			}
			f.doc(doc)
			if !isStruct {
				f.p("type %s %s", name, expr(n.Underlying()))
			} else {
				f.p("type %s struct {", name)
				for i := range st.NumFields() {
					fv := st.Field(i)
					// unexported fields are not encoded, tagged or not; an
					// embedded unexported struct is, through its exported
					// fields, so it is kept
					if !fv.Exported() && !fv.Embedded() {
						continue
					}
					// nor are fields tagged out of it
					if reflect.StructTag(st.Tag(i)).Get("json") == "-" {
						continue
					}
					tag := ""
					if json, ok := reflect.StructTag(st.Tag(i)).Lookup("json"); ok {
						tag = " `json:" + strconv.Quote(json) + "`"
					}
					f.doc(cloneDoc(g.fieldDocs[fv]))
					if fv.Embedded() {
						f.p("\t%s%s", expr(fv.Type()), tag)
					} else {
						f.p("\t%s %s%s", fv.Name(), expr(fv.Type()), tag)
					}
				}
				f.p("}")
			}
			var values []string
			for _, cn := range scope.Names() {
				if c, ok := scope.Lookup(cn).(*types.Const); ok && c.Exported() && !declared[cn] && types.Identical(c.Type(), n) {
					values = append(values, fmt.Sprintf("\t%s %s = %s", cn, name, constValue(c)))
				}
			}
			if len(values) > 0 {
				f.p("")
				f.p("const (")
				for _, v := range values {
					f.p("%s", v)
				}
				f.p(")")
			}
			f.p("")
		case *types.Const:
			if !obj.Exported() {
				continue
			}
			if n, ok := types.Unalias(obj.Type()).(*types.Named); ok && n.Obj().Pkg() == up.Types {
				// listed with its type
				continue
			}
			if b, ok := obj.Type().(*types.Basic); ok && b.Info()&types.IsUntyped != 0 {
				f.p("const %s = %s", name, constValue(obj))
			} else {
				f.p("const %s %s = %s", name, expr(obj.Type()), constValue(obj))
			}
			f.p("")
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("cloning %s into %s: %w", up.PkgPath, p.PkgPath, errors.Join(errs...))
	}
	return f.write(filepath.Dir(p.GoFiles[0]), typesFile)
}

// typesToClone is the exported types of a package, other than interfaces,
// functions and channels, which describe behavior rather than data, and the
// unexported types they reach through their fields, which are on the wire
// too.
func typesToClone(up *packages.Package) map[*types.TypeName]bool {
	out := map[*types.TypeName]bool{}
	var visit func(t types.Type)
	visit = func(t types.Type) {
		switch t := t.(type) {
		case *types.Alias:
			visit(types.Unalias(t))
		case *types.Pointer:
			visit(t.Elem())
		case *types.Slice:
			visit(t.Elem())
		case *types.Array:
			visit(t.Elem())
		case *types.Map:
			visit(t.Key())
			visit(t.Elem())
		case *types.Named:
			obj := t.Obj()
			if obj.Pkg() != up.Types || out[obj] {
				return
			}
			if behavior(t) {
				return
			}
			out[obj] = true
			if st, ok := t.Underlying().(*types.Struct); ok {
				for i := range st.NumFields() {
					if f := st.Field(i); f.Exported() || f.Embedded() {
						visit(f.Type())
					}
				}
			} else {
				visit(t.Underlying())
			}
		}
	}
	scope := up.Types.Scope()
	for _, name := range scope.Names() {
		obj, ok := scope.Lookup(name).(*types.TypeName)
		if !ok || !obj.Exported() {
			continue
		}
		if obj.IsAlias() {
			if !behavior(types.Unalias(obj.Type())) {
				out[obj] = true
			}
			continue
		}
		visit(obj.Type())
	}
	return out
}

func behavior(t types.Type) bool {
	switch t.Underlying().(type) {
	case *types.Interface, *types.Signature, *types.Chan:
		return true
	}
	return false
}

// methodsByHand returns the methods declared outside generated files of a
// die package, by receiver type name.
func methodsByHand(p *packages.Package) map[string]map[string]bool {
	out := map[string]map[string]bool{}
	for _, f := range p.Syntax {
		switch filepath.Base(p.Fset.File(f.Pos()).Name()) {
		case typesFile, diesFile, acFile:
			continue
		}
		for _, d := range f.Decls {
			fd, ok := d.(*ast.FuncDecl)
			if !ok || fd.Recv == nil || len(fd.Recv.List) == 0 {
				continue
			}
			recv := fd.Recv.List[0].Type
			if star, ok := recv.(*ast.StarExpr); ok {
				recv = star.X
			}
			if id, ok := recv.(*ast.Ident); ok {
				if out[id.Name] == nil {
					out[id.Name] = map[string]bool{}
				}
				out[id.Name][fd.Name.Name] = true
			}
		}
	}
	return out
}

// cloneTypeExpr prints a type for a clone: types of the cloned package are
// local, types of imported clones are in their die package, and Kubernetes
// and standard library types are referenced where they are. Anything else
// would bring back the dependencies cloning removes.
func cloneTypeExpr(f *file, cloned string, imports map[string]string, t types.Type) (string, error) {
	sub := func(t types.Type) (string, error) { return cloneTypeExpr(f, cloned, imports, t) }
	switch t := t.(type) {
	case *types.Alias:
		return sub(types.Unalias(t))
	case *types.Basic:
		return t.Name(), nil
	case *types.Pointer:
		s, err := sub(t.Elem())
		return "*" + s, err
	case *types.Slice:
		s, err := sub(t.Elem())
		return "[]" + s, err
	case *types.Array:
		s, err := sub(t.Elem())
		return fmt.Sprintf("[%d]%s", t.Len(), s), err
	case *types.Map:
		k, err := sub(t.Key())
		if err != nil {
			return "", err
		}
		v, err := sub(t.Elem())
		return "map[" + k + "]" + v, err
	case *types.Named:
		obj := t.Obj()
		if obj.Pkg() == nil {
			return obj.Name(), nil
		}
		path := obj.Pkg().Path()
		switch {
		case path == cloned:
			return obj.Name(), nil
		case imports[path] != "":
			return f.use(imports[path]) + "." + obj.Name(), nil
		case strings.HasPrefix(path, "k8s.io/") || !strings.Contains(strings.Split(path, "/")[0], "."):
			return f.use(path) + "." + obj.Name(), nil
		}
		return "", fmt.Errorf("%s is in %s, which is not cloned; clone it into a die package and add +die:import:package=%s,as=<that die package>", obj.Name(), path, path)
	case *types.Interface:
		if t.Empty() {
			return "any", nil
		}
	}
	return "", fmt.Errorf("unsupported type %s", t)
}

func constValue(c *types.Const) string {
	v := c.Val()
	if v.Kind() == constant.Float {
		f, _ := constant.Float64Val(v)
		return strconv.FormatFloat(f, 'g', -1, 64)
	}
	return v.ExactString()
}

// groupVersion is the apiVersion of a package's objects, from its GroupName
// constant or its SchemeGroupVersion variable.
func groupVersion(p *packages.Package) string {
	scope := p.Types.Scope()
	if c, ok := scope.Lookup("GroupName").(*types.Const); ok {
		group := constant.StringVal(c.Val())
		if group == "" {
			return p.Types.Name()
		}
		return group + "/" + p.Types.Name()
	}
	var group, version string
	for _, f := range p.Syntax {
		ast.Inspect(f, func(n ast.Node) bool {
			vs, ok := n.(*ast.ValueSpec)
			if !ok || len(vs.Names) != 1 || vs.Names[0].Name != "SchemeGroupVersion" || len(vs.Values) != 1 {
				return true
			}
			lit, ok := vs.Values[0].(*ast.CompositeLit)
			if !ok {
				return false
			}
			for _, e := range lit.Elts {
				kv, ok := e.(*ast.KeyValueExpr)
				if !ok {
					continue
				}
				key, _ := kv.Key.(*ast.Ident)
				tv := p.TypesInfo.Types[kv.Value]
				if key == nil || tv.Value == nil || tv.Value.Kind() != constant.String {
					continue
				}
				switch key.Name {
				case "Group":
					group = constant.StringVal(tv.Value)
				case "Version":
					version = constant.StringVal(tv.Value)
				}
			}
			return false
		})
	}
	if version == "" {
		return ""
	}
	if group == "" {
		return version
	}
	return group + "/" + version
}
