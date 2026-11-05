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

	"k8s.io/apimachinery/pkg/util/validation/field"

	metav1 "reconciler.io/dies/v2/apis/meta/v1"
	"reconciler.io/dies/v2/die"
	"reconciler.io/dies/v2/die/validation"
)

// Ready marks the widget's Ready condition true. It is a hand-written method
// on a generated die, which doesn't compile until the die is generated; that
// must not keep the die from being generated. It writes, so it is a method of
// the mutable die, and returns only an error, as writes do.
func (d *WidgetDieMut) Ready(reason string) error {
	ready := d.Status().Conditions().At("Ready")
	if err := ready.Status().Set(metav1.ConditionTrue); err != nil {
		return err
	}
	return ready.Reason().Set(reason)
}

// DieDefault defaults a port's number to 80. A widget's ports are defaulted by
// die.Default wherever they are, so neither WidgetDieMut nor WidgetSpecDieMut
// need a method to reach them.
func (d *WidgetPortDieMut) DieDefault(context.Context) error {
	if d.Port().Get() == 0 {
		return d.Port().Set(80)
	}
	return nil
}

// DieValidate validates a port beyond what its markers do, as die.ValidateCreate
// and die.ValidateUpdate call it for each port, which knows where it is: an error at its Port
// field is at spec.ports[i].port. Its markers' checks, of its name's format
// and its number's range, are made without it. In an update, the operation
// has the port's old value, of the port with the same name.
func (d *WidgetPortDie) DieValidate(_ context.Context, op die.Operation) (field.ErrorList, error) {
	var errs field.ErrorList
	if op.Type != die.Delete && d.Port().Get() == 22 {
		errs = append(errs, validation.Invalid(d.Port(), "is reserved for ssh"))
	}
	if old, ok := die.OldOf(d, op); ok && old.Port().Get() != d.Port().Get() {
		errs = append(errs, validation.Forbidden(d.Port(), "a named port's number may not change"))
	}
	return errs, nil
}
