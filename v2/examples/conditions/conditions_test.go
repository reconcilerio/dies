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

package conditions_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"reconciler.io/dies/v2/die"
	dietesting "reconciler.io/dies/v2/die/testing"
	"reconciler.io/dies/v2/examples/conditions"
	widgetv1 "reconciler.io/dies/v2/examples/widget/v1"
)

// clock ticks a second each time it is read, from 0.
func clock() func() metav1.Time {
	var t int64
	return func() metav1.Time {
		t++
		return metav1.Unix(t, 0)
	}
}

func widgetConditions(m *widgetv1.WidgetDieMut) *conditions.Conditions {
	return m.Status().Conditions()
}

// cascade is a widget blank whose conditions cascade:
//
//	A + B -> C
//	D + E -> F
//	C + F -> G
//
// with the observers registered in order, or in reverse.
func cascade(reverse bool) *widgetv1.WidgetDie {
	sets := func(set conditions.ConditionSet) func(m *widgetv1.WidgetDieMut) error {
		return observes(set)
	}
	observers := []func(m *widgetv1.WidgetDieMut) error{
		sets(conditions.NewConditionSet("C", "A", "B")),
		sets(conditions.NewConditionSet("F", "D", "E")),
		sets(conditions.NewConditionSet("G", "C", "F")),
		eachTransitions(clock()),
	}
	if reverse {
		slices.Reverse(observers)
	}
	return must(widgetv1.WidgetBlank.DieStampE(observers...))
}

// observes observes a widget's conditions with a condition set.
func observes(set conditions.ConditionSet) func(m *widgetv1.WidgetDieMut) error {
	return func(m *widgetv1.WidgetDieMut) error { return die.ObserveAt(m, widgetConditions, set.Observe) }
}

// eachTransitions observes each condition's lastTransitionTime.
func eachTransitions(now func() metav1.Time) func(m *widgetv1.WidgetDieMut) error {
	return func(m *widgetv1.WidgetDieMut) error {
		return m.Status().Conditions().DieObserveEach(conditions.LastTransitionTime(now))
	}
}

func must[D any](d D, err error) D {
	if err != nil {
		panic(err)
	}
	return d
}

// summary is each condition as type=status(reason).
func summary(d *widgetv1.WidgetDie) string {
	var out []string
	for t, c := range d.Status().Conditions().All() {
		out = append(out, fmt.Sprintf("%s=%s(%s)", t, c.Status().Get(), c.Reason().Get()))
	}
	slices.Sort(out)
	return strings.Join(out, " ")
}

// transitioned returns the conditions whose lastTransitionTime differs
// between two dies.
func transitioned(before, after *widgetv1.WidgetDie) string {
	var out []string
	for t, c := range after.Status().Conditions().All() {
		was := before.Status().Conditions().At(t).LastTransitionTime().Get()
		if now := c.LastTransitionTime().Get(); !now.Equal(&was) {
			out = append(out, t)
		}
	}
	slices.Sort(out)
	return strings.Join(out, " ")
}

// mark writes a condition's status and reason in one write.
func mark(t, status, reason string) func(d *widgetv1.WidgetDieMut) {
	return func(d *widgetv1.WidgetDieMut) {
		err := die.StampValue(d.Status().Conditions().At(t), func(c *metav1.Condition) {
			c.Status, c.Reason = metav1.ConditionStatus(status), reason
		})
		if err != nil {
			panic(err)
		}
	}
}

// Marking a condition updates the conditions that aggregate it, and those
// that aggregate them, whatever order the observers are registered in. Only
// the conditions that change transition.
func TestConditionsCascade(t *testing.T) {
	steps := []struct {
		name          string
		mark          func(d *widgetv1.WidgetDieMut)
		want, changed string
	}{
		{"B -> C -> G", mark("B", "False", "Broken"),
			"A=True(Ready) B=False(Broken) C=False(Broken) D=True(Ready) E=True(Ready) F=True(F) G=False(Broken)", "B C G"},
		{"D -> F, and G stays False for C", mark("D", "Unknown", "Waiting"),
			"A=True(Ready) B=False(Broken) C=False(Broken) D=Unknown(Waiting) E=True(Ready) F=Unknown(Waiting) G=False(Broken)", "D F"},
		{"B -> C -> G, Unknown for F", mark("B", "True", "Fixed"),
			"A=True(Ready) B=True(Fixed) C=True(C) D=Unknown(Waiting) E=True(Ready) F=Unknown(Waiting) G=Unknown(Waiting)", "B C G"},
		{"D -> F -> G", mark("D", "True", "Ready"),
			"A=True(Ready) B=True(Fixed) C=True(C) D=True(Ready) E=True(Ready) F=True(F) G=True(G)", "D F G"},
		{"E, again", mark("E", "True", "Ready"),
			"A=True(Ready) B=True(Fixed) C=True(C) D=True(Ready) E=True(Ready) F=True(F) G=True(G)", ""},
	}
	for _, reverse := range []bool{false, true} {
		t.Run(fmt.Sprintf("reverse=%t", reverse), func(t *testing.T) {
			d := cascade(reverse).DieStamp(
				mark("A", "True", "Ready"), mark("B", "True", "Ready"),
				mark("D", "True", "Ready"), mark("E", "True", "Ready"),
			)
			assertSummary(t, d, "A=True(Ready) B=True(Ready) C=True(C) D=True(Ready) E=True(Ready) F=True(F) G=True(G)")
			for _, step := range steps {
				next := d.DieStamp(step.mark)
				if got := summary(next); got != step.want {
					t.Errorf("%s:\n got: %s\nwant: %s", step.name, got, step.want)
				}
				if got := transitioned(d, next); got != step.changed {
					t.Errorf("%s: transitioned %q, want %q", step.name, got, step.changed)
				}
				d = next
			}
		})
	}
}

func assertSummary(t *testing.T, d *widgetv1.WidgetDie, want string) {
	t.Helper()
	if got := summary(d); got != want {
		t.Errorf("conditions\n got: %s\nwant: %s", got, want)
	}
}

// The happy condition observes the oldest generation its dependents do.
func TestConditionSetObservesTheOldestGeneration(t *testing.T) {
	ready := must(widgetv1.WidgetBlank.DieStampE(observes(conditions.NewConditionSet("Ready", "Configured", "Deployed"))))
	d := ready.DieStamp(func(d *widgetv1.WidgetDieMut) {
		dietesting.NoError(t, d.Metadata().Generation().Set(3))
		for _, c := range []struct {
			t   string
			gen int64
		}{{"Configured", 3}, {"Deployed", 2}} {
			dc := d.Status().Conditions().At(c.t)
			dietesting.NoError(t, dc.Status().Set(metav1.ConditionTrue))
			dietesting.NoError(t, dc.Reason().Set("Done"))
			dietesting.NoError(t, dc.ObservedGeneration().Set(c.gen))
		}
	})
	happy := d.Status().Conditions().At("Ready")
	if !happy.IsTrue() || happy.ObservedGeneration().Get() != 2 || !happy.IsStale() {
		t.Errorf("Ready %+v, want True of generation 2, stale", *happy.DieRelease())
	}
	// a missing dependent makes it Unknown
	d = ready.DieStamp(mark("Configured", "True", "Done"))
	if happy := d.Status().Conditions().At("Ready"); !happy.IsUnknown() || happy.Reason().Get() != "Initializing" {
		t.Errorf("Ready %+v, want Unknown for Deployed", *happy.DieRelease())
	}
}

// The happy condition is the aggregate of its dependents, whatever writes it.
func TestConditionSetKeepsTheHappyCondition(t *testing.T) {
	d := must(widgetv1.WidgetBlank.DieStampE(observes(conditions.NewConditionSet("Ready", "Configured")))).
		DieStamp(mark("Configured", "False", "Invalid"), mark("Ready", "True", "Forced"))
	if happy := d.Status().Conditions().At("Ready"); !happy.IsFalse() || happy.Reason().Get() != "Invalid" {
		t.Errorf("Ready %+v, want False for Configured", *happy.DieRelease())
	}
}

// lastTransitionTime is kept when nothing else changes, and when a change
// sets it.
func TestLastTransitionTime(t *testing.T) {
	ltt := func(d *widgetv1.WidgetDie) int64 {
		return d.Status().Conditions().At("Ready").LastTransitionTime().Get().Unix()
	}
	d := must(widgetv1.WidgetBlank.DieStampE(eachTransitions(clock()))).
		DieStamp(mark("Ready", "True", "Done"))
	if got := ltt(d); got != 1 {
		t.Errorf("added at %d, want 1", got)
	}
	// only the observed generation changes
	d = d.DieStamp(func(d *widgetv1.WidgetDieMut) {
		dietesting.NoError(t, d.Status().Conditions().At("Ready").ObservedGeneration().Set(2))
	})
	if got := ltt(d); got != 1 {
		t.Errorf("transitioned at %d for a new observed generation, want 1", got)
	}
	d = d.DieStamp(mark("Ready", "True", "Redone"))
	if got := ltt(d); got != 2 {
		t.Errorf("transitioned at %d for a new reason, want 2", got)
	}
	// the change sets it
	d = d.DieStamp(func(d *widgetv1.WidgetDieMut) {
		dietesting.NoError(t, die.StampValue(d.Status().Conditions().At("Ready"), func(c *metav1.Condition) {
			c.Status, c.LastTransitionTime = metav1.ConditionFalse, metav1.Unix(100, 0)
		}))
	})
	if got := ltt(d); got != 100 {
		t.Errorf("transitioned at %d, want the 100 the change set", got)
	}
}

// Within an atomic stamp, a condition written field by field is one change:
// it transitions once, and the conditions that aggregate it update once.
func TestConditionsOfAnAtomicStamp(t *testing.T) {
	ready := must(widgetv1.WidgetBlank.DieStampE(observes(conditions.NewConditionSet("Ready", "Configured")), eachTransitions(clock())))
	d, err := ready.DieStampAtomic(func(d *widgetv1.WidgetDieMut) error {
		c := d.Status().Conditions().At("Configured")
		if err := c.Status().Set(metav1.ConditionFalse); err != nil {
			return err
		}
		if err := c.Reason().Set("Invalid"); err != nil {
			return err
		}
		return c.Message().Set("no spec")
	})
	dietesting.NoError(t, err)
	assertSummary(t, d, "Configured=False(Invalid) Ready=False(Invalid)")
	for _, typ := range []string{"Configured", "Ready"} {
		if got := d.Status().Conditions().At(typ).LastTransitionTime().Get().Unix(); got > 2 {
			t.Errorf("%s transitioned at %d: more than once", typ, got)
		}
	}
	if got := d.Status().Conditions().At("Ready").Message().Get(); got != "no spec" {
		t.Errorf("Ready's message %q, want Configured's", got)
	}
}
