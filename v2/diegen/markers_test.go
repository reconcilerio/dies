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
	"go/parser"
	"go/token"
	"reflect"
	"testing"

	"golang.org/x/tools/go/packages"
)

func TestDieMarkers(t *testing.T) {
	tests := []struct {
		line string
		want dieMarker
	}{
		{"+die", dieMarker{name: "die", args: map[string][]string{}}},
		{"+die:object:apiVersion=apps/v1,kind=Deployment", dieMarker{name: "die:object", args: map[string][]string{
			"apiVersion": {"apps/v1"}, "kind": {"Deployment"},
		}}},
		{"+die:field:name=ports,listType=map,listMapKey={containerPort,protocol}", dieMarker{name: "die:field", args: map[string][]string{
			"name": {"ports"}, "listType": {"map"}, "listMapKey": {"containerPort", "protocol"},
		}}},
		// quoted values and JSON arrays keep their commas
		{`+die:field:name=protocol,default="a,b"`, dieMarker{name: "die:field", args: map[string][]string{
			"name": {"protocol"}, "default": {`"a,b"`},
		}}},
		{`+die:field:name=ports,default=[1,2]`, dieMarker{name: "die:field", args: map[string][]string{
			"name": {"ports"}, "default": {"[1,2]"},
		}}},
		// a default is JSON, so its braces are an object, not a list
		{`+die:field:name=disruptionMode,default={"single": {}}`, dieMarker{name: "die:field", args: map[string][]string{
			"name": {"disruptionMode"}, "default": {`{"single": {}}`},
		}}},
		{"+die:import:package=k8s.io/apimachinery/pkg/apis/meta/v1,as=reconciler.io/dies/v2/apis/meta/v1", dieMarker{name: "die:import", args: map[string][]string{
			"package": {"k8s.io/apimachinery/pkg/apis/meta/v1"}, "as": {"reconciler.io/dies/v2/apis/meta/v1"},
		}}},
	}
	for _, tt := range tests {
		got := dieMarkers([]string{tt.line})
		if len(got) != 1 || !reflect.DeepEqual(got[0], tt.want) {
			t.Errorf("%s\n got %+v\nwant %+v", tt.line, got, tt.want)
		}
	}
	if got := dieMarkers([]string{"+dies", "+listType=map", "not a marker"}); len(got) != 0 {
		t.Errorf("unexpected markers %+v", got)
	}
}

func TestSpecRoundTrip(t *testing.T) {
	s := newTypeSpec()
	*s.field("ports") = fieldSpec{listType: "map", listMapKey: []string{"containerPort", "protocol"}}
	*s.field("protocol") = fieldSpec{def: `"TCP"`, accessor: "Proto"}
	s.field("secret").union = "source"
	s.field("emptyDir").union = "source"
	s.subresource("status")

	lines := s.render([]string{"ports", "protocol", "secret", "emptyDir"})
	if got := dieSpec(nil, lines); !reflect.DeepEqual(got, s) {
		t.Errorf("round trip through\n%v\n got %+v\nwant %+v", lines, got, s)
	}
}

func TestPackageDocMarkers(t *testing.T) {
	tests := []struct {
		name, src string
		want      []string
	}{
		{"in the doc", "// Package v1 has dies.\n//\n// +die:types:package=k8s.io/api/core/v1\npackage v1\n", []string{"+die:types:package=k8s.io/api/core/v1"}},
		{"apart from the doc", "/* license */\n\n// +die:types:package=k8s.io/api/core/v1\n\n// Package v1 has dies.\npackage v1\n", []string{"+die:types:package=k8s.io/api/core/v1"}},
		{"without markers", "// Package v1 has dies.\npackage v1\n", nil},
		{"after the package clause", "package v1\n\n// +die:types:package=k8s.io/api/core/v1\n", nil},
	}
	for _, tt := range tests {
		f, err := parser.ParseFile(token.NewFileSet(), "doc.go", tt.src, parser.ParseComments)
		if err != nil {
			t.Fatal(err)
		}
		var got []string
		for _, m := range dieMarkers(packageDoc(&packages.Package{Syntax: []*ast.File{f}})) {
			got = append(got, "+"+m.name+":package="+m.arg("package"))
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%s: got %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestTypeDocMarkers(t *testing.T) {
	tests := []struct {
		name, src string
		want      []string
		unread    int
	}{
		{"in the doc", "package v1\n\n// Duck quacks.\n//\n// +die\ntype Duck struct{}\n", []string{"Duck quacks.", "", "+die"}, 0},
		{"apart from the doc", "package v1\n\n// +kubebuilder:object:root=true\n// +die\n\n// Duck quacks.\ntype Duck struct{}\n", []string{"+kubebuilder:object:root=true", "+die", "Duck quacks."}, 0},
		{"apart from no doc", "package v1\n\n// +die\n\ntype Duck struct{}\n", []string{"+die"}, 0},
		{"apart, in parentheses", "package v1\n\ntype (\n\tGoose struct{}\n\n\t// +die\n\n\t// Duck quacks.\n\tDuck struct{}\n)\n", []string{"+die", "Duck quacks."}, 0},
		{"not only markers", "package v1\n\n// Ducks.\n// +die\n\n// Duck quacks.\ntype Duck struct{}\n", []string{"Duck quacks."}, 1},
		{"in two comments", "package v1\n\n// +die:object:apiVersion=example.dev/v1,kind=Duck\n\n// +kubebuilder:subresource:status\n\n// Duck quacks.\ntype Duck struct{}\n", []string{"+kubebuilder:subresource:status", "Duck quacks."}, 1},
		{"other markers in two comments", "package v1\n\n// +genduck\n\n// +die\n\n// Duck quacks.\ntype Duck struct{}\n", []string{"+die", "Duck quacks."}, 0},
		{"of the previous declaration", "package v1\n\nfunc f() {\n\t// +die\n}\n\ntype Duck struct{}\n", nil, 0},
		{"of the previous spec", "package v1\n\ntype (\n\t// +die\n\tGoose struct{}\n\n\tDuck struct{}\n)\n", nil, 0},
	}
	for _, tt := range tests {
		f, err := parser.ParseFile(token.NewFileSet(), "types.go", tt.src, parser.ParseComments)
		if err != nil {
			t.Fatal(err)
		}
		var got []string
		var unread []*ast.CommentGroup
		for _, d := range f.Decls {
			if gd, ok := d.(*ast.GenDecl); ok && gd.Tok == token.TYPE {
				for _, spec := range gd.Specs {
					if ts := spec.(*ast.TypeSpec); ts.Name.Name == "Duck" {
						got, unread = typeDoc(f, gd, ts)
					}
				}
			}
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%s: got %q, want %q", tt.name, got, tt.want)
		}
		if n := len(unreadDieMarkers(unread)); n != tt.unread {
			t.Errorf("%s: got %d unread comments with +die markers, want %d", tt.name, n, tt.unread)
		}
	}
}

func TestValueChecks(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		want  []string
	}{
		{"declarative", []string{"+k8s:format=k8s-short-name", "+k8s:maxLength=63", "+k8s:maximum=1000 # Highest"}, []string{"format=k8s-short-name", "maximum=1000", "maxLength=63"}},
		{"kubebuilder", []string{"+kubebuilder:validation:Minimum=0.5", "+kubebuilder:validation:ExclusiveMinimum=true", "+kubebuilder:validation:UniqueItems=true"}, []string{"minimum=0.5", "exclusiveMinimum", "unique"}},
		{"kubebuilder enum and pattern", []string{`+kubebuilder:validation:Enum=a;"b c"`, "+kubebuilder:validation:Pattern=`^a;b$`"}, []string{`enum=["a","b c"]`, `pattern="^a;b$"`}},
		{"items", []string{"+k8s:eachVal=+k8s:format=k8s-label-key", "+kubebuilder:validation:items:MaxLength=63", "+k8s:maxItems=4"}, []string{"maxItems=4", "each.format=k8s-label-key", "each.maxLength=63"}},
		{"sets", []string{"+k8s:unique=set"}, []string{"unique"}},
		{"not a number", []string{"+k8s:minimum=ref(Min)"}, nil},
	}
	g := newGenerator()
	for _, tt := range tests {
		got := g.valueChecks(markers(tt.lines), nil)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%s: got %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestKubebuilderPackage(t *testing.T) {
	tests := []struct {
		name, src string
		want      *kubebuilderPackage
	}{
		{"custom resources", "package v1\n\n// +kubebuilder:object:root=true\ntype Duck struct{}\n", &kubebuilderPackage{required: true}},
		{"optional by default", "// +kubebuilder:validation:Optional\n\n// Package v1 has ducks.\npackage v1\n\n// +kubebuilder:subresource:status\ntype Duck struct{}\n", &kubebuilderPackage{required: false}},
		// markers that only validate, as metav1.Condition's do
		{"validation only", "package v1\n\ntype Condition struct {\n\t// +kubebuilder:validation:Required\n\tType string `json:\"type\"`\n}\n", nil},
	}
	for _, tt := range tests {
		f, err := parser.ParseFile(token.NewFileSet(), "types.go", tt.src, parser.ParseComments)
		if err != nil {
			t.Fatal(err)
		}
		g := newGenerator()
		g.indexKubebuilder(&packages.Package{PkgPath: "example.dev/v1", Syntax: []*ast.File{f}})
		if got := g.kubebuilder["example.dev/v1"]; !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%s: got %+v, want %+v", tt.name, got, tt.want)
		}
	}
}
