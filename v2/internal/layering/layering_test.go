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

// Package layering checks the layers of the module: the die API and the
// generator depend on no package of dies for Kubernetes' APIs, so that a
// user can generate their own dies for any of those APIs, of any version,
// without the module's dies, compiled against another, coming along.
package layering

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/tools/go/packages"
)

const module = "reconciler.io/dies/v2/"

// The packages, other than their tests, that mustn't depend on the module's
// dies for Kubernetes' APIs, and those dies.
var (
	layers = []string{module + "die/...", module + "diegen/..."}
	apis   = module + "apis/"
)

func TestDieAPIDoesNotDependOnAPIs(t *testing.T) {
	pkgs, err := packages.Load(&packages.Config{Mode: packages.NeedName | packages.NeedFiles | packages.NeedImports | packages.NeedDeps}, layers...)
	if err != nil {
		t.Fatal(err)
	}
	if len(pkgs) == 0 {
		t.Fatalf("no packages match %v", layers)
	}
	seen := map[string]bool{}
	var visit func(p *packages.Package, chain []string)
	visit = func(p *packages.Package, chain []string) {
		chain = append(chain, p.PkgPath)
		if strings.HasPrefix(p.PkgPath, apis) {
			t.Errorf("%s depends on %s: %s", chain[0], p.PkgPath, strings.Join(chain, " -> "))
			return
		}
		if seen[p.PkgPath] {
			return
		}
		seen[p.PkgPath] = true
		if strings.HasPrefix(p.PkgPath, module) {
			track(t, p)
		}
		for _, imp := range p.Imports {
			visit(imp, chain)
		}
	}
	for _, p := range pkgs {
		for _, e := range p.Errors {
			t.Errorf("%s: %v", p.PkgPath, e)
		}
		visit(p, nil)
	}
}

// track reads a package's directory and stats its files, as go test caches a
// test's result by the files it reads: packages.Load reads them in go list,
// which the cache doesn't see, so that a new import would pass from the
// cache.
func track(t *testing.T, p *packages.Package) {
	t.Helper()
	dirs := map[string]bool{}
	for _, f := range append(p.GoFiles, p.IgnoredFiles...) {
		dirs[filepath.Dir(f)] = true
		if _, err := os.Stat(f); err != nil {
			t.Error(err)
		}
	}
	for d := range dirs {
		if _, err := os.ReadDir(d); err != nil {
			t.Error(err)
		}
	}
}
