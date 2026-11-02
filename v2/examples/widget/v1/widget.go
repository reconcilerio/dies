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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Ready marks the widget's Ready condition true. It is a hand-written method
// on a generated die, which doesn't compile until the die is generated; that
// must not keep the die from being generated.
func (d *WidgetDie) Ready(reason string) *WidgetDie {
	return d.DieWith(func(d *WidgetDie) {
		d.Status().Conditions().At("Ready").Status().Set(metav1.ConditionTrue).Reason().Set(reason)
	})
}
