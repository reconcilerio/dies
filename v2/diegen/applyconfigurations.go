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
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"slices"
	"strings"

	"golang.org/x/tools/go/packages"
)

// Apply configurations. A die package says where the apply configurations
// of its types come from, for types whose publishers provide them, such as
// with applyconfiguration-gen:
//
//	// +die:applyconfigurations:package=example.dev/client/applyconfiguration/widgets/v1
//
// Its apply configuration types are aliased (mode=alias, the default), as are
// its Extract functions; its constructors aren't, as their names are those of
// the types. With no marker, the apply configurations Kubernetes publishes for
// the die package's types are aliased, and other die packages have none:
// their dies are apply configurations enough.
//
// The aliases are declared in the die package, and a die whose type has one
// has DieReleaseApplyConfiguration, DieFeedApplyConfiguration and
// DieStampApplyConfiguration.

const acFile = "zz_generated.dies.applyconfigurations.go"

// acSet is a die package's apply configurations.
type acSet struct {
	mode string // "" or "alias"
	// pkg is the package aliased
	pkg      string
	names    map[string]bool
	extracts []string
	err      error
}

// acsOf returns the apply configurations of a die package, whose types are
// in typesPath.
func (g *generator) acsOf(diePkg, typesPath string) *acSet {
	if set, ok := g.acs[diePkg]; ok {
		return set
	}
	set := &acSet{names: map[string]bool{}}
	g.acs[diePkg] = set
	meta := g.meta[diePkg]
	var marker dieMarker
	var marked bool
	if meta != nil {
		marker, marked = findDieMarker(packageDoc(meta.p), "die:applyconfigurations")
	}
	switch {
	case marked && marker.arg("package") == "":
		set.err = fmt.Errorf("%s: +die:applyconfigurations needs the package of the apply configurations to alias", diePkg)
		return set
	case marked:
		switch mode := marker.arg("mode"); mode {
		case "", "alias":
		default:
			set.err = fmt.Errorf("%s: +die:applyconfigurations: unknown mode %q, want alias", diePkg, mode)
			return set
		}
		set.pkg = marker.arg("package")
	default:
		if meta != nil && meta.upstream != "" {
			typesPath = meta.upstream
		}
		pkg, ok := defaultApplyConfigurations(typesPath)
		if !ok {
			return set
		}
		set.pkg = pkg
	}
	dir := "."
	if meta != nil {
		dir = filepath.Dir(meta.p.GoFiles[0])
	}
	set.names, set.extracts = scanApplyConfigurations(dir, set.pkg)
	switch {
	case len(set.names) > 0:
		set.mode = "alias"
	case marked:
		set.err = fmt.Errorf("%s: +die:applyconfigurations: %s has no apply configurations", diePkg, set.pkg)
	}
	return set
}

// defaultApplyConfigurations is where the apply configurations of a package's
// types are, for the packages that publish them.
func defaultApplyConfigurations(path string) (string, bool) {
	if rest, ok := strings.CutPrefix(path, "k8s.io/api/"); ok {
		return "k8s.io/client-go/applyconfigurations/" + rest, true
	}
	switch path {
	case "k8s.io/apimachinery/pkg/apis/meta/v1":
		return "k8s.io/client-go/applyconfigurations/meta/v1", true
	}
	for _, group := range []struct{ types, acs string }{
		{"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/", "k8s.io/apiextensions-apiserver/pkg/client/applyconfiguration/apiextensions/"},
		{"k8s.io/kube-aggregator/pkg/apis/apiregistration/", "k8s.io/kube-aggregator/pkg/client/applyconfiguration/apiregistration/"},
	} {
		if version, ok := strings.CutPrefix(path, group.types); ok {
			return group.acs + version, true
		}
	}
	return "", false
}

// scanApplyConfigurations reads the apply configuration types and Extract
// functions of a package from its source, as resolved from dir.
func scanApplyConfigurations(dir, acPkg string) (map[string]bool, []string) {
	names := map[string]bool{}
	var extracts []string
	pkgs, err := packages.Load(&packages.Config{Mode: packages.NeedName | packages.NeedFiles, Dir: dir}, acPkg)
	if err != nil || len(pkgs) != 1 || len(pkgs[0].Errors) > 0 {
		return names, nil
	}
	fset := token.NewFileSet()
	for _, file := range pkgs[0].GoFiles {
		f, err := parser.ParseFile(fset, file, nil, parser.SkipObjectResolution)
		if err != nil {
			continue
		}
		for _, d := range f.Decls {
			switch d := d.(type) {
			case *ast.GenDecl:
				if d.Tok != token.TYPE {
					continue
				}
				for _, s := range d.Specs {
					if n := s.(*ast.TypeSpec).Name.Name; strings.HasSuffix(n, "ApplyConfiguration") && ast.IsExported(n) {
						names[n] = true
					}
				}
			case *ast.FuncDecl:
				if d.Recv == nil && strings.HasPrefix(d.Name.Name, "Extract") && d.Type.TypeParams == nil {
					extracts = append(extracts, d.Name.Name)
				}
			}
		}
	}
	slices.Sort(extracts)
	return names, extracts
}

// applyConfiguration returns the apply configuration type for a target's
// type, if the die package has one.
func (pg *pkgGen) applyConfiguration(tg *target) (string, bool) {
	if _, ok := pg.listItem(tg); ok {
		return "", false
	}
	set := pg.g.acsOf(pg.p.PkgPath, pg.upstream)
	name := tg.name() + "ApplyConfiguration"
	if !set.names[name] {
		return "", false
	}
	return "*" + name, true
}

// emitApplyConfigurations writes the die package's apply configurations,
// aliases of a package's.
func (pg *pkgGen) emitApplyConfigurations(dir string) error {
	set := pg.g.acsOf(pg.p.PkgPath, pg.upstream)
	if set.err != nil {
		return set.err
	}
	f := newFile(pg.p.Name)
	switch set.mode {
	case "":
		return removeFile(dir, acFile)
	case "alias":
		acPkg := f.use(set.pkg)
		f.p("// The apply configuration types of [%s], and its Extract functions, so", set.pkg)
		f.p("// that importing this package provides them with the dies. Its")
		f.p("// constructors are not aliased: their names are those of the types.")
		f.p("")
		for _, name := range slices.Sorted(func(yield func(string) bool) {
			for n := range set.names {
				if !yield(n) {
					return
				}
			}
		}) {
			if err := pg.checkACName(name); err != nil {
				return err
			}
			f.p("type %s = %s.%s", name, acPkg, name)
		}
		if len(set.extracts) > 0 {
			f.p("")
			f.p("var (")
			for _, name := range set.extracts {
				if err := pg.checkACName(name); err != nil {
					return err
				}
				f.p("\t%s = %s.%s", name, acPkg, name)
			}
			f.p(")")
		}
	}
	return f.write(dir, acFile)
}

func (pg *pkgGen) checkACName(name string) error {
	if obj := pg.p.Types.Scope().Lookup(name); obj != nil && filepath.Base(pg.p.Fset.File(obj.Pos()).Name()) != acFile {
		return fmt.Errorf("%s: generated apply configuration %s collides with a declared name", pg.p.PkgPath, name)
	}
	return nil
}

// emitApplyConfiguration emits the methods that convert a die's intent to
// and from its type's apply configuration.
func (pg *pkgGen) emitApplyConfiguration(dname, mname, ac string) {
	f := pg.f
	for _, n := range []string{dname, mname} {
		f.p("// DieReleaseApplyConfiguration returns the die's intent as an apply")
		f.p("// configuration: on a root, what an apply sends, the object's identity")
		f.p("// included; on a die within another, its intent.")
		f.p("func (d *%s) DieReleaseApplyConfiguration() %s {", n, ac)
		f.p("	return die.ReleaseApplyConfiguration[%s](d)", ac[1:])
		f.p("}")
		f.p("")
		f.p("// DieReleaseApplyConfigurationE is [%s.DieReleaseApplyConfiguration],", n)
		f.p("// returning an error rather than panicking.")
		f.p("func (d *%s) DieReleaseApplyConfigurationE() (%s, error) {", n, ac)
		f.p("	return die.ReleaseApplyConfigurationE[%s](d)", ac[1:])
		f.p("}")
		f.p("")
	}
	f.p("// DieFeedApplyConfiguration replaces the die's intent with an apply")
	f.p("// configuration's, such as the fields a manager owns that Extract returns:")
	f.p("// each value it sets becomes intent, over the die's value.")
	f.p("func (d *%s) DieFeedApplyConfiguration(ac %s) *%s {", mname, ac, mname)
	f.p("	return die.FeedApplyConfiguration(d, ac)")
	f.p("}")
	f.p("")
	f.p("// DieStampApplyConfiguration passes the die's intent, as an apply")
	f.p("// configuration, to fn, and feeds back what fn makes of it.")
	f.p("func (d *%s) DieStampApplyConfiguration(fn func(ac %s)) *%s {", mname, ac, mname)
	f.p("	return die.StampApplyConfiguration(d, fn)")
	f.p("}")
	f.p("")
	f.p("// DieFeedApplyConfigurationE is [%s.DieFeedApplyConfiguration],", mname)
	f.p("// returning an error rather than panicking.")
	f.p("func (d *%s) DieFeedApplyConfigurationE(ac %s) (*%s, error) {", mname, ac, mname)
	f.p("	return die.FeedApplyConfigurationE(d, ac)")
	f.p("}")
	f.p("")
	f.p("// DieStampApplyConfigurationE is [%s.DieStampApplyConfiguration],", mname)
	f.p("// returning an error rather than panicking.")
	f.p("func (d *%s) DieStampApplyConfigurationE(fn func(ac %s)) (*%s, error) {", mname, ac, mname)
	f.p("	return die.StampApplyConfigurationE(d, fn)")
	f.p("}")
	f.p("")
}
