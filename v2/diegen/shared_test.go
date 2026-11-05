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
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"strings"
	"testing"

	"golang.org/x/tools/go/packages"
)

// sharedGen returns the generator of a package of one die, Duck's, with the
// methods of src written by hand.
func sharedGen(t *testing.T, src string) *pkgGen {
	t.Helper()
	fset := token.NewFileSet()
	var files []*ast.File
	for name, s := range map[string]string{
		diesFile:  "package v1\n\ntype Duck struct{}\ntype DuckDie struct{}\ntype DuckDieMut struct{}\n",
		"duck.go": "package v1\n\nimport \"time\"\n\nvar _ time.Duration\n\n" + src,
	} {
		f, err := parser.ParseFile(fset, name, s, parser.ParseComments)
		if err != nil {
			t.Fatal(err)
		}
		files = append(files, f)
	}
	info := &types.Info{Defs: map[*ast.Ident]types.Object{}, Uses: map[*ast.Ident]types.Object{}, Types: map[ast.Expr]types.TypeAndValue{}}
	pkg, err := (&types.Config{Importer: importer.Default()}).Check("example.dev/duck/v1", fset, files, info)
	if err != nil {
		t.Fatal(err)
	}
	p := &packages.Package{PkgPath: pkg.Path(), Fset: fset, Syntax: files, Types: pkg, TypesInfo: info}
	duck := pkg.Scope().Lookup("Duck").(*types.TypeName)
	return &pkgGen{pkgMeta: &pkgMeta{p: p, targets: []*target{{obj: duck}}}, f: newFile("v1")}
}

func TestSharedMethods(t *testing.T) {
	pg := sharedGen(t, `
// Quacks reports whether the duck quacks.
//
// +die:shared
func (d *DuckDie) Quacks() bool { return true }

// Waddle isn't shared.
func (d *DuckDie) Waddle() {}

// +die:shared
func (w *DuckDie) Swim(d time.Duration, _ int, strokes ...string) (time.Duration, error) { return 0, nil }

// +die:shared
func (DuckDie) Rest() {}
`)
	shared, err := pg.sharedMethods()
	if err != nil {
		t.Fatal(err)
	}
	for _, fn := range shared["DuckDie"] {
		pg.emitShared("DuckDie", "DuckDieMut", fn)
	}
	want := `// Quacks is [DuckDie.Quacks], shared by +die:shared.
func (d *DuckDieMut) Quacks() bool {
	return d.DieReadOnly().Quacks()
}

// Swim is [DuckDie.Swim], shared by +die:shared.
func (d_ *DuckDieMut) Swim(d time.Duration, p1 int, strokes ...string) (time.Duration, error) {
	return d_.DieReadOnly().Swim(d, p1, strokes...)
}

// Rest is [DuckDie.Rest], shared by +die:shared.
func (d *DuckDieMut) Rest() {
	d.DieReadOnly().Rest()
}

`
	if got := pg.f.buf.String(); got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
	if pg.f.imports["time"] != "time" {
		t.Errorf("imports %v, want time", pg.f.imports)
	}
}

func TestSharedMethodsMisplaced(t *testing.T) {
	tests := []struct{ name, src, want string }{
		{"on the mutable die", "// +die:shared\nfunc (d *DuckDieMut) Quacks() bool { return true }\n", "only a method of an immutable die, such as DuckDie, is shared"},
		{"on another type", "type Pond struct{}\n\n// +die:shared\nfunc (p *Pond) Deep() bool { return true }\n", "+die:shared on Pond.Deep: only a method of an immutable die"},
		{"declared by hand too", "// +die:shared\nfunc (d *DuckDie) Quacks() bool { return true }\n\nfunc (d *DuckDieMut) Quacks() bool { return true }\n", "DuckDieMut.Quacks is declared by hand"},
		{"generic", "// +die:shared\nfunc (d *DuckDie) Call[T any](t T) T { return t }\n", "a generic method isn't shared yet"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := sharedGen(t, tt.src).sharedMethods()
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("error %v, want it to contain %q", err, tt.want)
			}
		})
	}
}
