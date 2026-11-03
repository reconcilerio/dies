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
	"go/token"
	"go/types"
	"path/filepath"
	"reflect"
	"slices"
	"strings"

	"golang.org/x/tools/go/packages"
)

const diesFile = "zz_generated.dies.go"

// reserved method names that a field accessor cannot use.
var reserved = []string{"DeepCopy", "Equal", "MarshalJSON", "UnmarshalJSON", "GetObjectKind", "DeepCopyObject"}

// pkgGen generates the dies of one package.
type pkgGen struct {
	*pkgMeta
	g     *generator
	f     *file
	keys  map[string][]wfield
	specs map[*types.Named]*typeSpec
	// dies is the other die packages the generated dies refer to
	dies map[string]bool
	// checks are the interfaces beyond DieOf and DieMutOf that typed and
	// object dies are checked to satisfy, by the generated test
	checks []dieCheck
	err    error
}

// dieCheck is a die's pair of interfaces to satisfy.
type dieCheck struct {
	name, check, mutCheck string
}

func (g *generator) diesFor(p *packages.Package) error {
	pg := &pkgGen{pkgMeta: g.meta[p.PkgPath], g: g, f: newFile(p.Name), keys: map[string][]wfield{}, specs: map[*types.Named]*typeSpec{}, dies: map[string]bool{}}
	out := filepath.Join(filepath.Dir(p.GoFiles[0]), diesFile)
	if len(pg.targets) == 0 {
		if err := removeFile(filepath.Dir(out), diesFile); err != nil {
			return err
		}
		return removeFile(filepath.Dir(out), acFile)
	}
	if err := pg.emit(); err != nil {
		return err
	}
	if pg.err != nil {
		return pg.err
	}
	if err := g.checkDiePackages(p, pg.dies); err != nil {
		return err
	}
	if err := pg.f.write(filepath.Dir(out), diesFile); err != nil {
		return err
	}
	if err := pg.emitTest(filepath.Dir(out)); err != nil {
		return err
	}
	return pg.emitApplyConfigurations(filepath.Dir(out))
}

// checkDiePackages fails when generated dies would refer to a die package
// that does not exist, such as the default for a Kubernetes group that has
// no dies yet, rather than leaving it to fail to compile.
func (g *generator) checkDiePackages(p *packages.Package, paths map[string]bool) error {
	var missing []string
	for path := range paths {
		if !g.dieExists(p, path) {
			missing = append(missing, path)
		}
	}
	if len(missing) == 0 {
		return nil
	}
	slices.Sort(missing)
	return fmt.Errorf("%s refers to dies in packages that don't exist: %v; create them (a die package with +die:types), or map their upstream packages elsewhere with +die:import", p.PkgPath, missing)
}

func (g *generator) dieExists(from *packages.Package, path string) bool {
	if _, ok := g.exists[path]; !ok {
		pkgs, err := packages.Load(&packages.Config{Mode: packages.NeedName | packages.NeedFiles, Dir: filepath.Dir(from.GoFiles[0])}, path)
		g.exists[path] = err == nil && len(pkgs) == 1 && len(pkgs[0].Errors) == 0 && len(pkgs[0].GoFiles) > 0
	}
	return g.exists[path]
}

// emitTest writes a test that each die still matches its type. It catches
// an upstream dependency that moved forward without the dies being
// regenerated, which would otherwise compile and silently miss new fields.
//
// The test is in the package's external test package, so that go test
// compiles the package once, rather than again with the test in it.
func (pg *pkgGen) emitTest(dir string) error {
	f := newFile(pg.p.Name + "_test")
	f.use("reconciler.io/dies/v2/die")
	f.use("testing")
	self := f.use(pg.p.PkgPath)
	f.p("func TestDiesMatchTypes(t *testing.T) {")
	f.p("\tfor _, check := range []func() error{")
	for _, tg := range pg.targets {
		if _, ok := pg.listItem(tg); ok {
			f.p("\t\tdie.CheckList[%s.%s],", self, tg.name())
			continue
		}
		f.p("\t\tdie.CheckManaged[%s.%s],", self, tg.name())
	}
	f.p("\t} {")
	f.p("\t\tif err := check(); err != nil {")
	f.p("\t\t\tt.Error(err)")
	f.p("\t\t}")
	f.p("\t}")
	f.p("}")
	var lists []string
	for _, tg := range pg.targets {
		if items, ok := pg.listItemType(tg); ok {
			path, name, _ := pg.dieOf(items)
			lists = append(lists, fmt.Sprintf("\t_ die.DieListOf[*%s.%sDieMut] = (*%s.%sDie)(nil)", f.use(path), name, self, tg.name()))
		}
	}
	if len(pg.checks) > 0 || len(lists) > 0 {
		f.p("")
		f.p("// typed, object and list dies have the methods of their interfaces")
		f.p("var (")
		for _, c := range pg.checks {
			params := fmt.Sprintf("%s.%s, *%s.%sDie, *%s.%sDieMut", self, c.name, self, c.name, self, c.name)
			f.p("\t_ %s[%s] = (*%s.%sDie)(nil)", c.check, params, self, c.name)
			f.p("\t_ %s[%s] = (*%s.%sDieMut)(nil)", c.mutCheck, params, self, c.name)
		}
		for _, l := range lists {
			f.p("%s", l)
		}
		f.p(")")
	}
	if up, ok := pg.p.Imports[pg.upstream]; ok && pg.hasObjects() && up.Types.Scope().Lookup("AddToScheme") != nil {
		f.p("")
		f.p("func TestDiesRegisterUpstreamKinds(t *testing.T) {")
		f.p("	if err := die.CheckRegistered(%s.AddToScheme, %s.AddToScheme, %s.SchemeGroupVersion); err != nil {", self, f.use(pg.upstream), self)
		f.p("		t.Error(err)")
		f.p("	}")
		f.p("}")
	}
	return f.write(dir, "zz_generated.dies_test.go")
}

func marked(doc []string) bool {
	for _, m := range dieMarkers(doc) {
		if m.name == "die" || m.name == "die:object" {
			return true
		}
	}
	return false
}

func (pg *pkgGen) mapping(upstream string) (string, bool) {
	if out, ok := pg.mappings[upstream]; ok {
		return out, true
	}
	return defaultMapping(upstream)
}

// dieOf locates the die for a type: declared here, in a mapped die package,
// or with a type marked +die in its own package.
func (pg *pkgGen) dieOf(t types.Type) (pkgPath, name string, ok bool) {
	n, _, ok := structOf(t)
	if !ok {
		return "", "", false
	}
	if tg := pg.byNamed[n]; tg != nil {
		return pg.p.PkgPath, tg.name(), true
	}
	path := n.Obj().Pkg().Path()
	if out, ok := pg.mapping(path); ok && out != pg.p.PkgPath {
		return out, n.Obj().Name(), true
	}
	if marked(pg.g.typeDocs[n.Obj()]) {
		return path, n.Obj().Name(), true
	}
	return "", "", false
}

func (pg *pkgGen) hasDie(t types.Type) bool {
	_, _, ok := pg.dieOf(t)
	return ok
}

// typeExpr prints a type as seen from the die package. Types of mapped
// packages are referenced through their die package's aliases.
func (pg *pkgGen) typeExpr(t types.Type) string {
	switch t := t.(type) {
	case *types.Alias:
		return pg.typeExpr(types.Unalias(t))
	case *types.Basic:
		return t.Name()
	case *types.Pointer:
		return "*" + pg.typeExpr(t.Elem())
	case *types.Slice:
		return "[]" + pg.typeExpr(t.Elem())
	case *types.Map:
		return "map[" + pg.typeExpr(t.Key()) + "]" + pg.typeExpr(t.Elem())
	case *types.Named:
		if tg := pg.byNamed[t]; tg != nil {
			return tg.name()
		}
		obj := t.Obj()
		if obj.Pkg() == nil {
			return obj.Name()
		}
		path := obj.Pkg().Path()
		if out, ok := pg.mapping(path); ok {
			path = out
		}
		if path == pg.p.PkgPath {
			return obj.Name()
		}
		return pg.f.use(path) + "." + obj.Name()
	case *types.Interface:
		if t.Empty() {
			return "any"
		}
	}
	panic(fmt.Sprintf("diegen: unsupported type %s", t))
}

// dieMutExpr is the mutable die of a type.
func (pg *pkgGen) dieMutExpr(t types.Type) string { return pg.dieExpr(t) + "Mut" }

func (pg *pkgGen) dieExpr(t types.Type) string {
	path, name, ok := pg.dieOf(t)
	if !ok {
		panic(fmt.Sprintf("diegen: %s has no die", t))
	}
	if path == pg.p.PkgPath {
		return "*" + name + "Die"
	}
	pg.dies[path] = true
	return "*" + pg.f.use(path) + "." + name + "Die"
}

// specOf returns the structural markers of a type: those of a target, or
// the normalized upstream markers of any other type.
func (pg *pkgGen) specOf(n *types.Named, st *types.Struct) *typeSpec {
	if tg := pg.byNamed[n]; tg != nil {
		return tg.spec
	}
	if s, ok := pg.specs[n]; ok {
		return s
	}
	// a type with a die in another package has its spec there, including
	// that package's own metadata; the runtime schema comes from the same
	// place, through that package's annotations
	// a die package that doesn't exist falls back to upstream's markers here,
	// and is reported by checkDiePackages
	var s *typeSpec
	if path, _, ok := pg.dieOf(n); ok && path != pg.p.PkgPath && pg.g.dieExists(pg.p, path) {
		m, err := pg.g.metaFor(path)
		if err != nil {
			pg.err = errors.Join(pg.err, fmt.Errorf("reading the die metadata for %s: %w", n.Obj().Name(), err))
		} else {
			s = m.byName[qualified(n)]
		}
	}
	if s == nil {
		s = pg.g.upstreamSpec(n, st)
	}
	pg.specs[n] = s
	return s
}

// wfield is a field as it appears on the wire. Fields of embedded structs
// with an empty json name are flattened into their parent, as encoding/json
// does.
type wfield struct {
	v       *types.Var
	wire    string
	jsonTag string
	spec    *fieldSpec
	doc     []string
	// accessor is the die method name.
	accessor string
}

func (pg *pkgGen) flatten(n *types.Named, st *types.Struct) []wfield {
	spec := pg.specOf(n, st)
	var out []wfield
	for i := range st.NumFields() {
		f := st.Field(i)
		jsonTag := reflect.StructTag(st.Tag(i)).Get("json")
		name, opts, _ := strings.Cut(jsonTag, ",")
		if name == "-" && opts == "" {
			continue
		}
		if f.Embedded() && name == "" {
			en, est, ok := structOf(deref(f.Type()))
			if !ok {
				panic(fmt.Sprintf("diegen: cannot flatten embedded %s", f.Type()))
			}
			// the embedded type's annotations shape this type's schema, so
			// its die package must be initialized first
			if path, _, ok := pg.dieOf(en); ok && path != pg.p.PkgPath {
				pg.f.blank(path)
			}
			out = append(out, pg.flatten(en, est)...)
			continue
		}
		if !f.Exported() {
			continue
		}
		if name == "" {
			name = f.Name()
		}
		accessor := f.Name()
		if f.Embedded() {
			// an embedded type with a json name, such as ObjectMeta at
			// "metadata", is named after its wire name
			accessor = strings.ToUpper(name[:1]) + name[1:]
		}
		if fs := spec.fields[name]; fs != nil && fs.accessor != "" {
			accessor = fs.accessor
		}
		out = append(out, wfield{v: f, wire: name, jsonTag: jsonTag, spec: spec.fields[name], doc: docText(pg.g.fieldDocs[f]), accessor: accessor})
	}
	return out
}

// directFields are a struct's own fields, as declared rather than flattened.
func (g *generator) directFields(st *types.Struct) []wfield {
	var out []wfield
	for i := range st.NumFields() {
		fv := st.Field(i)
		if !fv.Exported() && !fv.Embedded() {
			continue
		}
		jsonTag := reflect.StructTag(st.Tag(i)).Get("json")
		wf := wfield{v: fv, wire: fv.Name(), jsonTag: jsonTag}
		if name, _, _ := strings.Cut(jsonTag, ","); name != "" {
			wf.wire = name
		}
		out = append(out, wf)
	}
	return out
}

// fkind is how a die exposes a field.
type fkind int

const (
	fValue fkind = iota
	fOptional
	fChild
	fListMap
	fMap
)

func (pg *pkgGen) kind(f wfield) fkind {
	t := types.Unalias(f.v.Type())
	// a struct with a die is a child, atomic or not: an atomic struct is
	// applied whole, which its die tag tells the runtime
	if p, ok := t.(*types.Pointer); ok {
		if pg.hasDie(p.Elem()) {
			return fChild
		}
		return fOptional
	}
	if pg.hasDie(t) {
		return fChild
	}
	if f.spec.atomic() {
		return fValue
	}
	switch u := t.Underlying().(type) {
	case *types.Slice:
		if f.spec.listMap() && pg.hasDie(u.Elem()) {
			return fListMap
		}
	case *types.Map:
		if b, ok := u.Key().Underlying().(*types.Basic); ok && b.Kind() == types.String && !pg.hasDie(u.Elem()) {
			return fMap
		}
	}
	return fValue
}

// dieTag expresses a field's structural markers as a die tag for the
// runtime.
func dieTag(owner *target, f wfield) string {
	var opts []string
	if fs := f.spec; fs != nil {
		if fs.listMap() {
			opts = append(opts, "listMap="+strings.Join(fs.listMapKey, ","))
		}
		if fs.atomic() {
			opts = append(opts, "atomic")
		}
		if fs.def != "" {
			opts = append(opts, "default="+fs.def)
		}
		if fs.union != "" {
			opts = append(opts, "union="+fs.union)
		}
		if fs.nullable {
			opts = append(opts, "nullable")
		}
	}
	if f.wire == "status" && slices.Contains(owner.spec.subresources, "status") {
		opts = append(opts, "subresource=status")
	}
	return strings.Join(opts, ";")
}

func (pg *pkgGen) emit() error {
	f := pg.f
	f.use("reconciler.io/dies/v2/die")

	// markers can't be struct tags on aliased or hand-written types, so they
	// are registered, with the package's field tables: a die package does no
	// work at init, its types are annotated before any of their schemas is
	// built, and its tables are built on first use
	var tables []string
	annotations := pg.annotations()
	defer func() {
		if pg.err == nil {
			pg.emitPackage(tables, annotations)
		}
	}()

	for _, tg := range pg.targets {
		if item, ok := pg.listItem(tg); ok {
			if err := pg.emitList(tg, item); err != nil {
				return err
			}
			continue
		}
		name, dname := tg.name(), tg.name()+"Die"
		mname := dname + "Mut"
		fields := pg.flatten(tg.named, tg.st)
		if err := pg.checkAccessors(tg, fields); err != nil {
			return err
		}

		ob := pg.objectBase(tg, name, dname)
		params := fmt.Sprintf("%s, *%s, *%s", name, dname, mname)
		if doc := docText(tg.doc); len(doc) > 0 {
			f.doc(doc)
			f.p("//")
		}
		f.p("// %s is immutable: it reads.", dname)
		f.p("// [%s], its mutable die, also writes, as DieStamp passes it.", mname)
		f.p("type %s struct {", dname)
		f.p("\tdie.%s[%s]", ob.base, params)
		f.p("}")
		f.p("")
		f.p("// %s is the mutable die of [%s]: see [%s].", mname, name, dname)
		f.p("type %s struct {", mname)
		f.p("\tdie.%s[%s]", ob.mutBase, params)
		f.p("}")
		f.p("")
		f.p("func (*%s) DieFrom(b die.DieBase[%s]) *%s { return &%s{%s} }", dname, params, dname, dname, ob.wrap)
		f.p("")
		f.p("func (*%s) DieFrom(b die.DieMutBase[%s]) *%s { return &%s{%s} }", mname, params, mname, mname, ob.mutWrap)
		f.p("")
		// the blank's type arguments check the die is a DieOf and its
		// mutable die a DieMutOf. Checks of more, which would compile method
		// wrappers for each die here, are made by the generated test.
		if ob.check != "die.DieOf" {
			pg.checks = append(pg.checks, dieCheck{name: name, check: ob.check, mutCheck: ob.mutCheck})
		}
		blank := fmt.Sprintf("die.Blank[%s]()", params)
		if m, ok := findDieMarker(tg.doc, "die:object"); ok {
			if typeMeta(tg.st) {
				// a literal, so that a blank costs no field table at init
				blank = fmt.Sprintf("die.BlankOf[%s](%s{TypeMeta: %s.TypeMeta{APIVersion: %q, Kind: %q}})", params, name, f.use("k8s.io/apimachinery/pkg/apis/meta/v1"), m.arg("apiVersion"), m.arg("kind"))
			} else {
				blank = fmt.Sprintf("die.TypedBlank[%s](%q, %q)", params, m.arg("apiVersion"), m.arg("kind"))
			}
		}
		f.p("var %sBlank = %s", name, blank)
		f.p("")
		if ac, ok := pg.applyConfiguration(tg); ok {
			pg.emitApplyConfiguration(dname, mname, ac)
		}

		// the field table is laid out statically and built on first use
		table := lowerFirst(name) + "Fields"
		tables = append(tables, table)
		f.p("type %sTable struct {", table)
		var inits []string
		for _, wf := range fields {
			decl, init := pg.fieldRef(tg, wf)
			f.p("\t%s %s", wf.accessor, decl)
			inits = append(inits, fmt.Sprintf("\tf.%s = %s", wf.accessor, init))
		}
		f.p("}")
		f.p("")
		f.p("var %s = &die.Table[%sTable]{Build: func() (f %sTable) {", table, table, table)
		if len(annotations) > 0 {
			// so that the package is annotated before any of its code that
			// runs at init uses a table
			f.p("\t_ = annotated")
		}
		for _, l := range inits {
			f.p("%s", l)
		}
		f.p("\treturn")
		f.p("}}")
		f.p("")

		for _, wf := range fields {
			pg.accessor(dname, table, wf)
		}
		for _, wf := range fields {
			pg.accessorMut(mname, table, wf)
		}
	}

	if err := pg.emitRegister(); err != nil {
		return err
	}

	names := slices.Sorted(func(yield func(string) bool) {
		for k := range pg.keys {
			if !yield(k) {
				return
			}
		}
	})
	for _, name := range names {
		// a name the generated dies declared last time is not a collision
		if obj := pg.p.Types.Scope().Lookup(name); obj != nil && filepath.Base(pg.p.Fset.File(obj.Pos()).Name()) != diesFile {
			return fmt.Errorf("%s: generated key %s collides with a declared name", pg.p.PkgPath, name)
		}
		f.p("// %s identifies an element of a list map.", name)
		f.p("type %s struct {", name)
		for _, kf := range pg.keys[name] {
			f.p("\t%s %s", kf.v.Name(), pg.typeExpr(kf.v.Type()))
		}
		f.p("}")
		f.p("")
	}
	return nil
}

func (pg *pkgGen) hasObjects() bool {
	for _, tg := range pg.targets {
		if _, ok := findDieMarker(tg.doc, "die:object"); ok {
			return true
		}
	}
	return false
}

// emitRegister emits the package's scheme registration: each object and
// list die as the Go type of its kind, as an upstream API package registers
// its types. A package with no object dies has a registration that does
// nothing, so group packages can register every version alike.
func (pg *pkgGen) emitRegister() error {
	f := pg.f
	type known struct{ kind, die string }
	var apiVersion string
	var kinds []known
	for _, tg := range pg.targets {
		m, ok := findDieMarker(tg.doc, "die:object")
		if !ok {
			continue
		}
		av := m.arg("apiVersion")
		if apiVersion != "" && av != apiVersion {
			return fmt.Errorf("%s: object dies of more than one group version, %s and %s", pg.p.PkgPath, apiVersion, av)
		}
		apiVersion = av
		// an object's registered type is its mutable die, which decodes; a
		// list die is one type
		d := tg.name() + "DieMut"
		if _, ok := pg.listItem(tg); ok {
			d = tg.name() + "Die"
		}
		kinds = append(kinds, known{m.arg("kind"), d})
	}
	for _, name := range []string{"SchemeGroupVersion", "SchemeBuilder", "AddToScheme", "addKnownTypes"} {
		if obj := pg.p.Types.Scope().Lookup(name); obj != nil && filepath.Base(pg.p.Fset.File(obj.Pos()).Name()) != diesFile {
			return fmt.Errorf("%s: generated scheme registration %s collides with a declared name", pg.p.PkgPath, name)
		}
	}
	runtime := f.use("k8s.io/apimachinery/pkg/runtime")
	if len(kinds) == 0 {
		f.p("// SchemeBuilder registers this package's object and list dies, of which")
		f.p("// it has none.")
		f.p("var SchemeBuilder = %s.NewSchemeBuilder()", runtime)
		f.p("")
		f.p("// AddToScheme registers this package's dies with a scheme: here, nothing.")
		f.p("var AddToScheme = SchemeBuilder.AddToScheme")
		f.p("")
		return nil
	}
	group, version, ok := strings.Cut(apiVersion, "/")
	if !ok {
		group, version = "", apiVersion
	}
	schema := f.use("k8s.io/apimachinery/pkg/runtime/schema")
	metav1 := f.use("k8s.io/apimachinery/pkg/apis/meta/v1")
	f.p("// SchemeGroupVersion is the group version of this package's object dies.")
	f.p("var SchemeGroupVersion = %s.GroupVersion{Group: %q, Version: %q}", schema, group, version)
	f.p("")
	f.p("// SchemeBuilder registers this package's mutable object dies, and list")
	f.p("// dies, as the Go types of their kinds, in place of the upstream types.")
	f.p("var SchemeBuilder = %s.NewSchemeBuilder(addKnownTypes)", runtime)
	f.p("")
	f.p("// AddToScheme registers this package's dies with a scheme.")
	f.p("var AddToScheme = SchemeBuilder.AddToScheme")
	f.p("")
	f.p("func addKnownTypes(scheme *%s.Scheme) error {", runtime)
	for _, k := range kinds {
		f.p("	scheme.AddKnownTypeWithName(SchemeGroupVersion.WithKind(%q), &%s{})", k.kind, k.die)
	}
	f.p("	%s.AddToGroupVersion(scheme, SchemeGroupVersion)", metav1)
	f.p("	return nil")
	f.p("}")
	f.p("")
	return nil
}

// listItem returns the die of a list's items, for an object with list
// metadata and items that have dies.
func (pg *pkgGen) listItem(tg *target) (string, bool) {
	items, ok := pg.listItemType(tg)
	if !ok {
		return "", false
	}
	return pg.dieMutExpr(items), true
}

// listItemType returns the type of a list's items, for an object with list
// metadata and items that have dies.
func (pg *pkgGen) listItemType(tg *target) (types.Type, bool) {
	if _, ok := findDieMarker(tg.doc, "die:object"); !ok {
		return nil, false
	}
	var listMeta bool
	var items types.Type
	for _, wf := range pg.g.directFields(tg.st) {
		switch wf.wire {
		case "metadata":
			n, ok := types.Unalias(wf.v.Type()).(*types.Named)
			listMeta = ok && n.Obj().Pkg() != nil && n.Obj().Pkg().Path() == "k8s.io/apimachinery/pkg/apis/meta/v1" && n.Obj().Name() == "ListMeta"
		case "items":
			if sl, ok := types.Unalias(wf.v.Type()).(*types.Slice); ok {
				items = sl.Elem()
			}
		}
	}
	if !listMeta || items == nil || !pg.hasDie(items) {
		return nil, false
	}
	return items, true
}

// emitList emits a list die: a plain list of item dies, as a Kubernetes list
// is, which libraries that reflect on a list's Items fill and read.
func (pg *pkgGen) emitList(tg *target, item string) error {
	f := pg.f
	name, dname := tg.name(), tg.name()+"Die"
	ctor := "New" + name
	if obj := pg.p.Types.Scope().Lookup(ctor); obj != nil && filepath.Base(pg.p.Fset.File(obj.Pos()).Name()) != diesFile {
		return fmt.Errorf("%s: generated constructor %s collides with a declared name", pg.p.PkgPath, ctor)
	}
	m, _ := findDieMarker(tg.doc, "die:object")
	runtime := f.use("k8s.io/apimachinery/pkg/runtime")
	f.doc(docText(tg.doc))
	f.p("type %s struct {", dname)
	f.p("	die.List[%s]", item)
	f.p("}")
	f.p("")
	f.p("// %s returns a list with items.", ctor)
	f.p("func %s(items ...%s) *%s {", ctor, item, dname)
	f.p("	return &%s{die.NewList(%q, %q, items...)}", dname, m.arg("apiVersion"), m.arg("kind"))
	f.p("}")
	f.p("")
	f.p("// DeepCopy returns a copy of the list, with its items copied.")
	f.p("func (l *%s) DeepCopy() *%s { return &%s{l.DieCopy()} }", dname, dname, dname)
	f.p("")
	f.p("func (l *%s) DeepCopyObject() %s.Object { return l.DeepCopy() }", dname, runtime)
	f.p("")
	return nil
}

// typeMeta reports whether a struct embeds metav1.TypeMeta.
func typeMeta(st *types.Struct) bool {
	for i := range st.NumFields() {
		v := st.Field(i)
		n, ok := types.Unalias(v.Type()).(*types.Named)
		if v.Embedded() && ok && n.Obj().Pkg() != nil && n.Obj().Pkg().Path() == "k8s.io/apimachinery/pkg/apis/meta/v1" && n.Obj().Name() == "TypeMeta" {
			return true
		}
	}
	return false
}

// dieBase is what a die and its mutable die embed, how DieFrom wraps their
// state in it, and the interfaces of their kind they're checked against.
type dieBase struct {
	base, mutBase   string
	wrap, mutWrap   string
	check, mutCheck string
}

// objectBase is the base a die embeds, and the interface of its kind it's
// checked against: an object die implements runtime.Object, metav1.Object
// by its metadata, and DieObjectOf; a typed die runtime.Object and
// DieTypedOf; and any other die DieOf. The mutable die has the mutable
// counterparts.
func (pg *pkgGen) objectBase(tg *target, name, dname string) dieBase {
	if _, ok := findDieMarker(tg.doc, "die:object"); !ok {
		return dieBase{"DieBase", "DieMutBase", "b", "b", "die.DieOf", "die.DieMutOf"}
	}
	for _, wf := range pg.g.directFields(tg.st) {
		n, ok := types.Unalias(wf.v.Type()).(*types.Named)
		if wf.wire != "metadata" || !ok || n.Obj().Pkg() == nil || n.Obj().Pkg().Path() != "k8s.io/apimachinery/pkg/apis/meta/v1" {
			continue
		}
		if n.Obj().Name() == "ObjectMeta" {
			// DieObjectOf names the dies of meta/v1, which can't import it
			if pg.p.PkgPath == "reconciler.io/dies/v2/apis/meta/v1" {
				return dieBase{"DieObjectBase", "DieObjectMutBase", "die.AsObject(b)", "die.AsObjectMut(b)", "die.DieTypedOf", "die.DieTypedMutOf"}
			}
			return dieBase{"DieObjectBase", "DieObjectMutBase", "die.AsObject(b)", "die.AsObjectMut(b)", "die.DieObjectOf", "die.DieObjectMutOf"}
		}
	}
	return dieBase{"DieTypedBase", "DieTypedMutBase", "die.AsTyped(b)", "die.AsTypedMut(b)", "die.DieTypedOf", "die.DieTypedMutOf"}
}

// emitPackage registers the package with the runtime: the packages declaring
// its types, their annotations, and how to build its field tables. The
// annotations and tables are registered apart, as each table refers to the
// annotations, to be initialized after them.
func (pg *pkgGen) emitPackage(tables, annotations []string) {
	f := pg.f
	var paths []string
	for _, tg := range pg.targets {
		if p := tg.named.Obj().Pkg().Path(); !slices.Contains(paths, p) {
			paths = append(paths, p)
		}
	}
	slices.Sort(paths)
	types := func() {
		f.p("\tTypes: []string{")
		for _, p := range paths {
			f.p("\t\t%q,", p)
		}
		f.p("\t},")
	}
	if len(annotations) > 0 {
		f.p("var annotated = die.Register(die.Package{")
		types()
		f.p("\tAnnotate: func() {")
		for _, l := range annotations {
			f.p("%s", l)
		}
		f.p("\t},")
		f.p("})")
		f.p("")
	}
	if len(tables) > 0 {
		f.p("var _ = die.Register(die.Package{")
		types()
		f.p("\tTables: func() {")
		for _, t := range tables {
			f.p("\t\t%s.Get()", t)
		}
		f.p("\t},")
		f.p("})")
		f.p("")
	}
}

// annotations are the die tags of the package's types, for die.Annotate,
// and the identity of its objects, for die.Object.
func (pg *pkgGen) annotations() []string {
	var lines []string
	for _, tg := range pg.targets {
		var ann []string
		for _, wf := range pg.g.directFields(tg.st) {
			wf.spec = tg.spec.fields[wf.wire]
			if d := dieTag(tg, wf); d != "" {
				ann = append(ann, fmt.Sprintf("\t\t%q: %q,", wf.v.Name(), d))
			}
		}
		if m, ok := findDieMarker(tg.doc, "die:object"); ok {
			lines = append(lines, fmt.Sprintf("\tdie.Object[%s](%q, %q)", tg.name(), m.arg("apiVersion"), m.arg("kind")))
		}
		if len(ann) == 0 {
			continue
		}
		lines = append(lines, fmt.Sprintf("\tdie.Annotate[%s](map[string]string{", tg.name()))
		lines = append(lines, ann...)
		lines = append(lines, "\t})")
	}
	return lines
}

// objectMethods are the methods of metav1.Object and metav1.ListInterface,
// which an object or list die's accessors cannot use.
var objectMethods = []string{
	"GetNamespace", "SetNamespace", "GetName", "SetName", "GetGenerateName", "SetGenerateName",
	"GetUID", "SetUID", "GetResourceVersion", "SetResourceVersion", "GetGeneration", "SetGeneration",
	"GetSelfLink", "SetSelfLink", "GetCreationTimestamp", "SetCreationTimestamp",
	"GetDeletionTimestamp", "SetDeletionTimestamp", "GetDeletionGracePeriodSeconds", "SetDeletionGracePeriodSeconds",
	"GetLabels", "SetLabels", "GetAnnotations", "SetAnnotations", "GetFinalizers", "SetFinalizers",
	"GetOwnerReferences", "SetOwnerReferences", "GetManagedFields", "SetManagedFields",
	"GetContinue", "SetContinue", "GetRemainingItemCount", "SetRemainingItemCount",
}

// checkAccessors rejects accessor names that collide, saying how to rename
// them.
func (pg *pkgGen) checkAccessors(tg *target, fields []wfield) error {
	base := pg.objectBase(tg, tg.name(), tg.name()+"Die").base
	rename := func(wf wfield, why string) error {
		how := "on the type"
		if tg.obj.IsAlias() {
			how = fmt.Sprintf("with metadata, type _ = %s", pg.typeExpr(tg.named))
		}
		return fmt.Errorf("%s.%s: the accessor for field %s %s; rename it %s:\n\t// +die:field:name=%s,accessor=<name>", pg.p.PkgPath, tg.name(), wf.wire, why, how, wf.wire)
	}
	seen := map[string]bool{}
	for _, wf := range fields {
		a := wf.accessor
		switch {
		case !token.IsIdentifier(a) || !token.IsExported(a):
			return rename(wf, fmt.Sprintf("is %q, which is not an exported Go name", a))
		case strings.HasPrefix(a, "Die") || slices.Contains(reserved, a):
			return rename(wf, fmt.Sprintf("is %s, which collides with a die method", a))
		case base != "DieBase" && slices.Contains(objectMethods, a):
			return rename(wf, fmt.Sprintf("is %s, which collides with a metav1 object method", a))
		case seen[a]:
			return rename(wf, fmt.Sprintf("is %s, which another field's accessor already is", a))
		}
		seen[a] = true
	}
	return nil
}

// fieldRef returns the declaration and initializer of a field table entry.
func (pg *pkgGen) fieldRef(tg *target, wf wfield) (string, string) {
	t := tg.name()
	ft := types.Unalias(wf.v.Type())
	switch pg.kind(wf) {
	case fChild:
		c, cd, cdm := pg.typeExpr(deref(ft)), pg.dieExpr(deref(ft)), pg.dieMutExpr(deref(ft))
		return fmt.Sprintf("die.Struct[%s, %s, %s, %s]", t, c, cd, cdm),
			fmt.Sprintf("die.NewStruct[%s, %s, %s, %s](%q)", t, c, cd, cdm, wf.wire)
	case fListMap:
		elem := ft.Underlying().(*types.Slice).Elem()
		c, cd, cdm := pg.typeExpr(elem), pg.dieExpr(elem), pg.dieMutExpr(elem)
		k, keyFn, keyOf := pg.listMapKey(elem, wf.spec.listMapKey)
		return fmt.Sprintf("die.ListMapField[%s, %s, %s, %s, %s]", t, c, cd, cdm, k),
			fmt.Sprintf("die.NewListMap[%s, %s, %s, %s](%q, %s, %s)", t, c, cd, cdm, wf.wire, keyFn, keyOf)
	case fMap:
		m := ft.Underlying().(*types.Map)
		k, v := pg.typeExpr(m.Key()), pg.typeExpr(m.Elem())
		return fmt.Sprintf("die.MapField[%s, %s, %s]", t, k, v),
			fmt.Sprintf("die.NewMap[%s, %s, %s](%q)", t, k, v, wf.wire)
	}
	v := pg.typeExpr(ft)
	return fmt.Sprintf("die.Field[%s, %s]", t, v),
		fmt.Sprintf("die.NewField[%s, %s](%q)", t, v, wf.wire)
}

// listMapKey returns the key type of a list map with its key functions. A
// single key uses the key field's type; several keys use a key struct,
// generated in this package.
func (pg *pkgGen) listMapKey(elem types.Type, keys []string) (string, string, string) {
	en, est, _ := structOf(elem)
	fields := pg.flatten(en, est)
	var kfs []wfield
	for _, k := range keys {
		i := slices.IndexFunc(fields, func(wf wfield) bool { return wf.wire == k })
		if i < 0 {
			panic(fmt.Sprintf("diegen: %s has no list-map key %q", elem, k))
		}
		kfs = append(kfs, fields[i])
	}
	c := pg.typeExpr(elem)
	if len(kfs) == 1 {
		k := pg.typeExpr(kfs[0].v.Type())
		return k,
			fmt.Sprintf("func(k %s) []any { return []any{k} }", k),
			fmt.Sprintf("func(e %s) %s { return e.%s }", c, k, kfs[0].v.Name())
	}
	_, name, _ := pg.dieOf(elem)
	k := name + "Key"
	pg.keys[k] = kfs
	var vals, inits []string
	for _, kf := range kfs {
		vals = append(vals, "k."+kf.v.Name())
		inits = append(inits, fmt.Sprintf("%s: e.%s", kf.v.Name(), kf.v.Name()))
	}
	return k,
		fmt.Sprintf("func(k %s) []any { return []any{%s} }", k, strings.Join(vals, ", ")),
		fmt.Sprintf("func(e %s) %s { return %s{%s} }", c, k, k, strings.Join(inits, ", "))
}

func (pg *pkgGen) accessor(dname, table string, wf wfield) {
	f := pg.f
	ft := types.Unalias(wf.v.Type())
	f.doc(wf.doc)
	switch pg.kind(wf) {
	case fChild:
		f.p("func (d *%s) %s() %s { return die.Child(d, %s.Get().%s) }", dname, wf.accessor, pg.dieExpr(deref(ft)), table, wf.accessor)
	case fListMap:
		elem := ft.Underlying().(*types.Slice).Elem()
		k, _, _ := pg.listMapKey(elem, wf.spec.listMapKey)
		f.p("func (d *%s) %s() *die.ListMap[%s, %s, %s] { return die.ListMapOf(d, %s.Get().%s) }", dname, wf.accessor, pg.typeExpr(elem), pg.dieExpr(elem), k, table, wf.accessor)
	case fMap:
		m := ft.Underlying().(*types.Map)
		f.p("func (d *%s) %s() *die.Map[%s, %s] { return die.MapOf(d, %s.Get().%s) }", dname, wf.accessor, pg.typeExpr(m.Key()), pg.typeExpr(m.Elem()), table, wf.accessor)
	case fOptional:
		f.p("func (d *%s) %s() die.Optional[%s] { return die.OptionalOf(d, %s.Get().%s) }", dname, wf.accessor, pg.typeExpr(deref(ft)), table, wf.accessor)
	default:
		f.p("func (d *%s) %s() die.Value[%s] { return die.ValueOf(d, %s.Get().%s) }", dname, wf.accessor, pg.typeExpr(ft), table, wf.accessor)
	}
	f.p("")
}

// accessorMut emits a mutable die's accessor for a field, which writes it.
func (pg *pkgGen) accessorMut(mname, table string, wf wfield) {
	f := pg.f
	ft := types.Unalias(wf.v.Type())
	f.doc(wf.doc)
	switch pg.kind(wf) {
	case fChild:
		f.p("func (d *%s) %s() %s { return die.ChildMut(d, %s.Get().%s) }", mname, wf.accessor, pg.dieMutExpr(deref(ft)), table, wf.accessor)
	case fListMap:
		elem := ft.Underlying().(*types.Slice).Elem()
		k, _, _ := pg.listMapKey(elem, wf.spec.listMapKey)
		f.p("func (d *%s) %s() *die.ListMapMut[%s, %s, %s] { return die.ListMapMutOf(d, %s.Get().%s) }", mname, wf.accessor, pg.typeExpr(elem), pg.dieMutExpr(elem), k, table, wf.accessor)
	case fMap:
		m := ft.Underlying().(*types.Map)
		f.p("func (d *%s) %s() *die.MapMut[%s, %s] { return die.MapMutOf(d, %s.Get().%s) }", mname, wf.accessor, pg.typeExpr(m.Key()), pg.typeExpr(m.Elem()), table, wf.accessor)
	case fOptional:
		f.p("func (d *%s) %s() die.OptionalMut[%s] { return die.OptionalMutOf(d, %s.Get().%s) }", mname, wf.accessor, pg.typeExpr(deref(ft)), table, wf.accessor)
	default:
		f.p("func (d *%s) %s() die.ValueMut[%s] { return die.ValueMutOf(d, %s.Get().%s) }", mname, wf.accessor, pg.typeExpr(ft), table, wf.accessor)
	}
	f.p("")
}
