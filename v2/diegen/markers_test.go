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
