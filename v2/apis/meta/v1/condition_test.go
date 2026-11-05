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

package v1_test

import (
	"testing"

	dieadmissionregistrationv1 "reconciler.io/dies/v2/apis/admissionregistration/v1"
	diemetav1 "reconciler.io/dies/v2/apis/meta/v1"
	dietesting "reconciler.io/dies/v2/die/testing"
)

// conditionHelpers are the helpers of a condition die, immutable or mutable.
type conditionHelpers interface {
	IsTrue() bool
	IsFalse() bool
	IsUnknown() bool
	IsStale() bool
}

// both returns a condition's immutable die and its mutable die.
func both(d *diemetav1.ConditionDie) map[string]conditionHelpers {
	return map[string]conditionHelpers{"immutable": d, "mutable": d.DieMutable()}
}

func TestConditionStatus(t *testing.T) {
	tests := []struct {
		name                   string
		status                 diemetav1.ConditionStatus
		set                    bool
		isTrue, isFalse, isUnk bool
	}{
		{name: "unset", isUnk: true},
		{name: "True", status: diemetav1.ConditionTrue, set: true, isTrue: true},
		{name: "False", status: diemetav1.ConditionFalse, set: true, isFalse: true},
		{name: "Unknown", status: diemetav1.ConditionUnknown, set: true, isUnk: true},
		{name: "empty", status: "", set: true, isUnk: true},
		{name: "another value", status: "Maybe", set: true, isUnk: true},
	}
	for _, tt := range tests {
		c := diemetav1.ConditionBlank.DieStamp(func(d *diemetav1.ConditionDieMut) {
			dietesting.NoError(t, d.Type().Set("Ready"))
			if tt.set {
				dietesting.NoError(t, d.Status().Set(tt.status))
			}
		})
		for kind, d := range both(c) {
			t.Run(tt.name+", "+kind, func(t *testing.T) {
				if got := d.IsTrue(); got != tt.isTrue {
					t.Errorf("IsTrue %t, want %t", got, tt.isTrue)
				}
				if got := d.IsFalse(); got != tt.isFalse {
					t.Errorf("IsFalse %t, want %t", got, tt.isFalse)
				}
				if got := d.IsUnknown(); got != tt.isUnk {
					t.Errorf("IsUnknown %t, want %t", got, tt.isUnk)
				}
			})
		}
	}
}

func TestConditionIsStale(t *testing.T) {
	tests := []struct {
		name     string
		observed int64
		set      bool
		stale    bool
	}{
		{name: "no observed generation"},
		// the wire can't tell an observed generation of 0 from none
		{name: "observed zero", observed: 0, set: true},
		{name: "an older generation", observed: 1, set: true, stale: true},
		{name: "the object's generation", observed: 2, set: true},
		{name: "a newer generation", observed: 3, set: true, stale: true},
	}
	for _, tt := range tests {
		policy := dieadmissionregistrationv1.ValidatingAdmissionPolicyBlank.DieStamp(func(d *dieadmissionregistrationv1.ValidatingAdmissionPolicyDieMut) {
			dietesting.NoError(t, d.Metadata().Generation().Set(2))
			c := d.Status().Conditions().At("Ready")
			dietesting.NoError(t, c.Status().Set(diemetav1.ConditionTrue))
			if tt.set {
				dietesting.NoError(t, c.ObservedGeneration().Set(tt.observed))
			}
		})
		dies := map[string]conditionHelpers{
			"immutable": policy.Status().Conditions().At("Ready"),
			"mutable":   policy.DieMutable().Status().Conditions().At("Ready"),
		}
		for kind, d := range dies {
			t.Run(tt.name+", "+kind, func(t *testing.T) {
				if got := d.IsStale(); got != tt.stale {
					t.Errorf("IsStale %t, want %t", got, tt.stale)
				}
			})
		}
	}
}

// An object the API server hasn't seen has no generation, so nothing its
// conditions observe is stale.
func TestConditionOfAnUnseenObjectIsFresh(t *testing.T) {
	policy := dieadmissionregistrationv1.ValidatingAdmissionPolicyBlank.DieStamp(func(d *dieadmissionregistrationv1.ValidatingAdmissionPolicyDieMut) {
		dietesting.NoError(t, d.Status().Conditions().At("Ready").ObservedGeneration().Set(3))
	})
	dies := map[string]conditionHelpers{
		"immutable": policy.Status().Conditions().At("Ready"),
		"mutable":   policy.DieMutable().Status().Conditions().At("Ready"),
	}
	for kind, d := range dies {
		if d.IsStale() {
			t.Errorf("%s: stale, want fresh", kind)
		}
	}
}

// A condition not within an object has no generation to compare with.
func TestConditionNotOfAnObjectIsFresh(t *testing.T) {
	c := diemetav1.ConditionBlank.DieStamp(func(d *diemetav1.ConditionDieMut) {
		dietesting.NoError(t, d.ObservedGeneration().Set(3))
	})
	for kind, d := range both(c) {
		if d.IsStale() {
			t.Errorf("%s: stale, want fresh", kind)
		}
	}
}

// The helpers of a mutable die read its current state, including writes made
// after it was obtained, and, within an object, the object's generation.
func TestConditionHelpersOfAMutableDieReadItsWrites(t *testing.T) {
	policy := dieadmissionregistrationv1.ValidatingAdmissionPolicyBlank.DieMutable()
	dietesting.NoError(t, policy.Metadata().Generation().Set(1))
	c := policy.Status().Conditions().At("Ready")

	dietesting.NoError(t, c.Status().Set(diemetav1.ConditionFalse))
	dietesting.NoError(t, c.ObservedGeneration().Set(1))
	if !c.IsFalse() || c.IsTrue() || c.IsUnknown() || c.IsStale() {
		t.Errorf("condition %+v", *c.DieRelease())
	}

	dietesting.NoError(t, c.Status().Set(diemetav1.ConditionTrue))
	if !c.IsTrue() {
		t.Error("not True after it was set True")
	}
	dietesting.NoError(t, policy.Metadata().Generation().Set(2))
	if !c.IsStale() {
		t.Error("fresh after the object's generation moved on")
	}
}
