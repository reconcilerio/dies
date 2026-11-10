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

// Package conditions is an example of observers: conditions that manage
// themselves. A [ConditionSet] keeps a happy condition, such as Ready, the
// aggregate of the conditions it depends on, and [LastTransitionTime] keeps
// each condition's lastTransitionTime, as a die's conditions change.
//
// A condition set observes a die's conditions, registered with
// die.ObserveAt, and [LastTransitionTime] each of them, registered with the
// conditions' DieObserveEach:
//
//	err := die.ObserveAt(widget, func(m *widgetv1.WidgetDieMut) *conditions.Conditions {
//		return m.Status().Conditions()
//	}, ready.Observe)
//	err = widget.Status().Conditions().DieObserveEach(conditions.LastTransitionTime(metav1.Now))
//
// Condition sets cascade: a set whose happy condition is another's
// dependent updates the other in turn.
package conditions

import (
	"fmt"
	"slices"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	diemetav1 "reconciler.io/dies/v2/apis/meta/v1"
	"reconciler.io/dies/v2/die"
)

// Conditions is the handle to a die's conditions that their observers are
// given.
type Conditions = diemetav1.ConditionListMapMut

// ConditionsChange is a change to a die's conditions, as their immutable list
// maps before and after.
type ConditionsChange = die.Change[*diemetav1.ConditionListMap]

// Observer is an observer of a die's conditions, as die.ObserveAt registers
// it.
type Observer = func(conditions *Conditions, c ConditionsChange) error

// ConditionSet aggregates conditions into a happy condition, such as Ready:
// it is True when every dependent is True, False when any is False, and
// Unknown otherwise, as a dependent that is missing is. It observes the
// oldest generation its dependents observe, so it is stale while any of them
// is. Other conditions don't affect it.
type ConditionSet struct {
	happy      string
	dependents []string
}

// NewConditionSet returns the condition set whose happy condition aggregates
// the dependents. It panics if there are none, the happy condition is also a
// dependent, or a dependent is listed twice.
func NewConditionSet(happy string, dependents ...string) ConditionSet {
	if len(dependents) == 0 {
		panic(fmt.Sprintf("condition set: %s has no dependents", happy))
	}
	for i, d := range dependents {
		if d == happy {
			panic(fmt.Sprintf("condition set: %s can't depend on itself", happy))
		}
		if slices.Contains(dependents[:i], d) {
			panic(fmt.Sprintf("condition set: %s is a dependent of %s twice", d, happy))
		}
	}
	return ConditionSet{happy: happy, dependents: slices.Clone(dependents)}
}

// Observe is the condition set as an observer of a die's conditions: when a
// dependent, or the happy condition, changes, it sets the happy condition to
// the dependents' aggregate, in one write. A False or Unknown happy condition
// takes its reason and message from the first dependent, in the set's
// order, that makes it so; a True one has its own type as its reason, and no
// message. It depends only on the dependents as they are, so the order
// observers are registered in, and how a change to the dependents is
// written, don't matter. Its lastTransitionTime is [LastTransitionTime]'s to
// keep.
func (s ConditionSet) Observe(conditions *Conditions, c ConditionsChange) error {
	if !slices.ContainsFunc(append([]string{s.happy}, s.dependents...), func(t string) bool { return changed(c, t) }) {
		return nil
	}
	status, reason, message := metav1.ConditionTrue, s.happy, ""
	var observed int64
	for _, t := range s.dependents {
		if !conditions.Has(t) {
			if status == metav1.ConditionTrue {
				status, reason, message = metav1.ConditionUnknown, "Initializing", fmt.Sprintf("%s is not reported", t)
			}
			continue
		}
		d := conditions.At(t)
		if g := d.ObservedGeneration().Get(); g != 0 && (observed == 0 || g < observed) {
			observed = g
		}
		switch {
		case d.IsFalse() && status != metav1.ConditionFalse:
			status, reason, message = metav1.ConditionFalse, d.Reason().Get(), d.Message().Get()
		case d.IsUnknown() && status == metav1.ConditionTrue:
			status, reason, message = metav1.ConditionUnknown, d.Reason().Get(), d.Message().Get()
		}
	}
	happy := conditions.At(s.happy)
	if conditions.Has(s.happy) && happy.Status().Get() == status && happy.Reason().Get() == reason &&
		happy.Message().Get() == message && happy.ObservedGeneration().Get() == observed {
		return nil
	}
	return die.StampValue(happy, func(c *metav1.Condition) {
		c.Status, c.Reason, c.Message, c.ObservedGeneration = status, reason, message, observed
	})
}

// ConditionObserver is an observer of each of a die's conditions, as
// DieObserveEach registers it.
type ConditionObserver = func(condition *diemetav1.ConditionDieMut, c die.Change[*diemetav1.ConditionDie]) error

// LastTransitionTime returns an observer of each of a die's conditions that
// sets a condition's lastTransitionTime to now when it is added, or its
// status, reason or message changes, unless the change set its
// lastTransitionTime too.
func LastTransitionTime(now func() metav1.Time) ConditionObserver {
	return func(condition *diemetav1.ConditionDieMut, c die.Change[*diemetav1.ConditionDie]) error {
		before, after := c.Before, c.After
		was, is := before.LastTransitionTime().Get(), after.LastTransitionTime().Get()
		switch {
		case !after.DieHas():
			// removed
			return nil
		case before.DieHas() && !was.Equal(&is):
			// the change set it
			return nil
		case before.DieHas() && before.Status().Get() == after.Status().Get() &&
			before.Reason().Get() == after.Reason().Get() && before.Message().Get() == after.Message().Get():
			return nil
		}
		return condition.LastTransitionTime().Set(now())
	}
}

// changed reports whether a change changed a condition, apart from its
// lastTransitionTime, which follows the rest of it.
func changed(c ConditionsChange, t string) bool {
	if c.Before.Has(t) != c.After.Has(t) {
		return true
	}
	before, after := *c.Before.At(t).DieRelease(), *c.After.At(t).DieRelease()
	before.LastTransitionTime, after.LastTransitionTime = metav1.Time{}, metav1.Time{}
	return before != after
}
