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

package apis_test

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"reconciler.io/dies/v2/apis"
	dieapiextensionsv1 "reconciler.io/dies/v2/apis/apiextensions/v1"
	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	diecorev1 "reconciler.io/dies/v2/apis/core/v1"
)

func TestAddToSchemeRegistersEveryGroup(t *testing.T) {
	scheme := runtime.NewScheme()
	if err := apis.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	for gvk, want := range map[schema.GroupVersionKind]any{
		{Group: "apps", Version: "v1", Kind: "Deployment"}:                               &dieappsv1.DeploymentDie{},
		{Group: "apps", Version: "v1", Kind: "DeploymentList"}:                           &dieappsv1.DeploymentListDie{},
		{Version: "v1", Kind: "Pod"}:                                                     &diecorev1.PodDie{},
		{Group: "apiextensions.k8s.io", Version: "v1", Kind: "CustomResourceDefinition"}: &dieapiextensionsv1.CustomResourceDefinitionDie{},
	} {
		obj, err := scheme.New(gvk)
		if err != nil {
			t.Errorf("%s: %v", gvk, err)
			continue
		}
		if reflect.TypeOf(obj) != reflect.TypeOf(want) {
			t.Errorf("%s is %T, want %T", gvk, obj, want)
		}
	}
	// every registered type is a die, or one of apimachinery's meta types
	for gvk, typ := range scheme.AllKnownTypes() {
		if pkg := typ.PkgPath(); !strings.HasPrefix(pkg, "reconciler.io/dies/v2/") && pkg != "k8s.io/apimachinery/pkg/apis/meta/v1" {
			t.Errorf("%s is registered as %s", gvk, typ)
		}
	}
}

// The group packages and this one are written by hand, so this checks they
// register every version package and every group.
func TestRegistrationIsComplete(t *testing.T) {
	groups, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	var want []string
	for _, g := range groups {
		if !g.IsDir() {
			continue
		}
		want = append(want, "reconciler.io/dies/v2/apis/"+g.Name())
		versions, err := filepath.Glob(filepath.Join(g.Name(), "*", "zz_generated.dies.go"))
		if err != nil {
			t.Fatal(err)
		}
		var wantVersions []string
		for _, v := range versions {
			wantVersions = append(wantVersions, "reconciler.io/dies/v2/apis/"+filepath.ToSlash(filepath.Dir(v)))
		}
		assertImports(t, filepath.Join(g.Name(), "register.go"), wantVersions)
	}
	assertImports(t, "register.go", want)
}

func assertImports(t *testing.T, file string, want []string) {
	t.Helper()
	f, err := parser.ParseFile(token.NewFileSet(), file, nil, parser.ImportsOnly)
	if err != nil {
		t.Error(err)
		return
	}
	var have []string
	for _, imp := range f.Imports {
		p, _ := strconv.Unquote(imp.Path.Value)
		if strings.HasPrefix(p, "reconciler.io/dies/v2/apis/") {
			have = append(have, p)
		}
	}
	slices.Sort(have)
	slices.Sort(want)
	if !slices.Equal(have, want) {
		t.Errorf("%s registers %v, want %v", file, short(have), short(want))
	}
}

func short(paths []string) []string {
	out := make([]string, len(paths))
	for i, p := range paths {
		out[i] = strings.TrimPrefix(p, "reconciler.io/dies/v2/apis/")
	}
	return out
}
