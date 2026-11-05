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

package core

import (
	"reflect"
	"testing"
)

func TestSplitTag(t *testing.T) {
	got := splitTag(`required;pattern="^a;\"b\"$";enum=["x;y"];maxLength=3`)
	want := []string{"required", `pattern="^a;\"b\"$"`, `enum=["x;y"]`, "maxLength=3"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestParseChecks(t *testing.T) {
	fd := &FieldDesc{Wire: "names", Typ: reflect.TypeFor[[]string]()}
	fd.parseTag(`required=nonzero;oneOf=source;maxItems=2;unique;each.pattern="^[a-z;]+$";each.enum=["a","b"]`)
	if !fd.RequiredNonZero || fd.Required || fd.Union != "source" || !fd.OneOf {
		t.Errorf("got %+v", fd)
	}
	c := fd.Checks
	if c == nil || *c.MaxItems != 2 || !c.Unique || c.Each == nil || c.Each.Pattern.String() != "^[a-z;]+$" || !reflect.DeepEqual(c.Each.Enum, []string{"a", "b"}) {
		t.Errorf("got %+v", c)
	}
	defer func() {
		if r := recover(); r == nil {
			t.Error("an unknown option didn't panic")
		}
	}()
	(&FieldDesc{Wire: "x", Typ: reflect.TypeFor[string]()}).parseTag("maximal=1")
}
