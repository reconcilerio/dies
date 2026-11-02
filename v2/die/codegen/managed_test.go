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

package codegen_test

import (
	"strings"
	"testing"

	"reconciler.io/dies/v2/die/base"
	"reconciler.io/dies/v2/die/codegen"
	"reconciler.io/dies/v2/die/fields"
)

// widget stands in for a type that gained a field after its die was
// generated: the die's field table only knows about name.
type widget struct {
	Name  string `json:"name"`
	Color string `json:"color,omitempty"`
}

type widgetDie struct {
	base.DieBase[widget, *widgetDie, *widgetDieMut]
}

type widgetDieMut struct {
	base.DieMutBase[widget, *widgetDie, *widgetDieMut]
}

func (*widgetDie) DieFrom(b base.DieBase[widget, *widgetDie, *widgetDieMut]) *widgetDie {
	return &widgetDie{b}
}

func (*widgetDieMut) DieFrom(b base.DieMutBase[widget, *widgetDie, *widgetDieMut]) *widgetDieMut {
	return &widgetDieMut{b}
}

var widgetFields = struct {
	Name fields.Field[widget, string]
}{
	Name: fields.NewField[widget, string]("name"),
}

func TestCheckManaged(t *testing.T) {
	err := codegen.CheckManaged[widget]()
	if err == nil || !strings.Contains(err.Error(), "new fields [color]") {
		t.Errorf("expected the new field to be reported, got %v", err)
	}
	if fields := base.ManagedFields[widget](); len(fields) != 1 || fields[0].Wire != "name" {
		t.Errorf("managed fields %+v", fields)
	}
	_ = widgetFields
}
