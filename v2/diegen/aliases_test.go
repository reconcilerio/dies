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
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/tools/go/packages"
)

// upstreamImporter imports the package checked here by its path, and others
// as std does, so that both are checked against the same packages.
type upstreamImporter struct {
	up  *types.Package
	std types.Importer
}

func (i upstreamImporter) Import(path string) (*types.Package, error) {
	if path == i.up.Path() {
		return i.up, nil
	}
	return i.std.Import(path)
}

// An aliased package re-exports the upstream package's functions, by
// functions of the same signature that call them, and its variables, but not
// its scheme registration, nor what names what can't be named outside it.
func TestAliasesForwardFunctionsAndVariables(t *testing.T) {
	const upstream = `package v1

import (
	"errors"
	"time"
)

type Duck struct{}

type hidden struct{}

// ErrNoDuck is returned when there's no duck.
var ErrNoDuck = errors.New("no duck")

// Deprecated: use ErrNoDuck.
var ErrNoGoose = ErrNoDuck

var SchemeGroupVersion = struct{ Group, Version string }{"ducks.example.dev", "v1"}

func AddToScheme() error { return nil }

var localDucks []Duck

// Hatch hatches a duck at a time.
func Hatch(time time.Time) Duck { return Duck{} }

// Flock gathers ducks.
func Flock(leader *Duck, _ int, followers ...Duck) ([]Duck, error) { return nil, nil }

// Pick picks the first of some values.
func Pick[T any, S ~[]T](s S) (T, bool) { var t T; return t, len(s) > 0 }

// Deprecated: use Hatch.
func Lay() {}

func Hide() hidden { return hidden{} }

func Internal(func(struct{ x int })) {}

func unexported() {}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "duck.go", upstream, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	info := &types.Info{Defs: map[*ast.Ident]types.Object{}, Uses: map[*ast.Ident]types.Object{}, Types: map[ast.Expr]types.TypeAndValue{}}
	std := importer.Default()
	up, err := (&types.Config{Importer: std}).Check("example.dev/duck/v1", fset, []*ast.File{file}, info)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	dp := &packages.Package{Name: "v1", PkgPath: "example.dev/dies/duck/v1", GoFiles: []string{filepath.Join(dir, "doc.go")}}
	upp := &packages.Package{PkgPath: up.Path(), Fset: fset, Syntax: []*ast.File{file}, Types: up, TypesInfo: info}
	if err := newGenerator().emitAliases(dp, upp); err != nil {
		t.Fatal(err)
	}
	src, err := os.ReadFile(filepath.Join(dir, typesFile))
	if err != nil {
		t.Fatal(err)
	}
	out := string(src)
	for _, want := range []string{
		"// ErrNoDuck is [duckv1.ErrNoDuck].\nvar ErrNoDuck = duckv1.ErrNoDuck\n",
		"// ErrNoGoose is [duckv1.ErrNoGoose].\n//\n// Deprecated: use ErrNoDuck.\n//\n//lint:ignore SA1019 a variable is of its upstream variable, deprecated or not\nvar ErrNoGoose = duckv1.ErrNoGoose //nolint:staticcheck",
		"// Hatch calls [duckv1.Hatch].\nfunc Hatch(p0 time.Time) duckv1.Duck {\n\treturn duckv1.Hatch(p0)\n}\n",
		"func Flock(leader *duckv1.Duck, p1 int, followers ...duckv1.Duck) ([]duckv1.Duck, error) {\n\treturn duckv1.Flock(leader, p1, followers...)\n}\n",
		"func Pick[T any, S ~[]T](s S) (T, bool) {\n\treturn duckv1.Pick[T, S](s)\n}\n",
		"// Lay calls [duckv1.Lay].\n//\n// Deprecated: use Hatch.\nfunc Lay() {\n\t//lint:ignore SA1019 a function is called through, deprecated or not\n\tduckv1.Lay() //nolint:staticcheck",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing\n%s\nin\n%s", want, out)
		}
	}
	for _, unwanted := range []string{"SchemeGroupVersion =", "func AddToScheme", "localDucks", "func Hide", "func Internal", "unexported"} {
		if strings.Contains(out, unwanted) {
			t.Errorf("%q re-exported in\n%s", unwanted, out)
		}
	}
	// what it re-exports compiles against the upstream package
	gen, err := parser.ParseFile(fset, typesFile, src, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := (&types.Config{Importer: upstreamImporter{up: up, std: std}}).Check(dp.PkgPath, fset, []*ast.File{gen}, nil); err != nil {
		t.Errorf("%v\n%s", err, out)
	}
}
