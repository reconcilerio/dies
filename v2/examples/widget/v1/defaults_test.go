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

package v1

import (
	"context"
	"errors"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"reconciler.io/dies/v2/die"
)

// DieDefault fails for a note of "fail", for tests of defaulting's errors.
func (d *WidgetMatchDieMut) DieDefault(context.Context) error {
	if d.Note().Get() == "fail" {
		return errors.New("a note of fail")
	}
	return nil
}

// DieValidate fails for a note of "unknown", for tests of validation's
// errors.
func (d *WidgetMatchDie) DieValidate(context.Context, die.Operation) (field.ErrorList, error) {
	if d.Note().Get() == "unknown" {
		return nil, errors.New("a note of unknown can't be validated")
	}
	return nil, nil
}
