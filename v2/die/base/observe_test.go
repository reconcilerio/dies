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

package base_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	diecorev1 "reconciler.io/dies/v2/apis/core/v1"
	"reconciler.io/dies/v2/die"
	"reconciler.io/dies/v2/die/fields"
	dietesting "reconciler.io/dies/v2/die/testing"
)

func minReady(m *dieappsv1.DeploymentDieMut) fields.ValueMut[int32] {
	return m.Spec().MinReadySeconds()
}

func labels(m *dieappsv1.DeploymentDieMut) *fields.MapMut[string, string] {
	return m.Metadata().Labels()
}

func annotations(m *dieappsv1.DeploymentDieMut) *fields.MapMut[string, string] {
	return m.Metadata().Annotations()
}

// show formats a change of an int32 as {before after had has}.
func show(c die.Change[fields.Value[int32]]) string {
	return fmt.Sprintf("{%d %d %t %t}", c.Before.Get(), c.After.Get(), c.Before.Has(), c.After.Has())
}

// root returns the mutable deployment a handle is within.
func root(l die.Located) *dieappsv1.DeploymentDieMut {
	return l.(interface{ DieRoot() die.Die }).DieRoot().(*dieappsv1.DeploymentDieMut)
}

func TestObserverIsCalledWithTheChange(t *testing.T) {
	d := dieappsv1.DeploymentBlank.DieMutable()
	var changes []string
	dietesting.NoError(t, die.ObserveAt(d, minReady, func(r fields.ValueMut[int32], c die.Change[fields.Value[int32]]) error {
		changes = append(changes, show(c))
		// the handle is of the die as the mutation left it
		if r.Get() != c.After.Get() {
			t.Errorf("handle reads %d, the change is to %d", r.Get(), c.After.Get())
		}
		return nil
	}))

	dietesting.NoError(t, d.Spec().MinReadySeconds().Set(1))
	dietesting.NoError(t, d.Spec().MinReadySeconds().Set(1)) // no change
	dietesting.NoError(t, d.Spec().Paused().Set(true))       // not the target
	dietesting.NoError(t, d.Spec().MinReadySeconds().Set(3))
	// a clear is told from a zero by presence
	dietesting.NoError(t, d.Spec().MinReadySeconds().Clear())
	if got := fmt.Sprint(changes); got != "[{0 1 false true} {1 3 true true} {3 0 true false}]" {
		t.Errorf("changes %s", got)
	}
}

// An observer writes anywhere in the die, and what it writes notifies the
// die's other observers in turn.
func TestObserversCascade(t *testing.T) {
	d := dieappsv1.DeploymentBlank.DieMutable()
	// minReadySeconds -> a label -> an annotation
	dietesting.NoError(t, die.ObserveAt(d, minReady, func(r fields.ValueMut[int32], c die.Change[fields.Value[int32]]) error {
		return root(r).Metadata().Labels().Set("size", fmt.Sprint(c.After.Get()))
	}))
	dietesting.NoError(t, die.ObserveAt(d, labels, func(l *fields.MapMut[string, string], c die.Change[*fields.Map[string, string]]) error {
		return root(l).Metadata().Annotations().Set("seen-size", c.After.At("size"))
	}))
	dietesting.NoError(t, d.Spec().MinReadySeconds().Set(2))
	if got := d.Metadata().Annotations().At("seen-size"); got != "2" {
		t.Errorf("annotation %q", got)
	}
}

// An observer that writes its own target is called again until it settles;
// one that never settles fails the write that called it.
func TestObserversSettle(t *testing.T) {
	d := dieappsv1.DeploymentBlank.DieMutable()
	calls := 0
	// minReadySeconds are rounded up to an even number
	dietesting.NoError(t, die.ObserveAt(d, minReady, func(r fields.ValueMut[int32], c die.Change[fields.Value[int32]]) error {
		calls++
		if c.After.Get()%2 == 1 {
			return r.Set(c.After.Get() + 1)
		}
		return nil
	}))
	dietesting.NoError(t, d.Spec().MinReadySeconds().Set(3))
	if got := d.Spec().MinReadySeconds().Get(); got != 4 || calls != 2 {
		t.Errorf("minReadySeconds %d after %d calls, want 4 after 2", got, calls)
	}

	err := die.ObserveAt(d, labels, func(l *fields.MapMut[string, string], c die.Change[*fields.Map[string, string]]) error {
		return l.Set("n", c.After.At("n")+"+")
	})
	dietesting.NoError(t, err)
	if err := d.Metadata().Labels().Set("n", "0"); err == nil || !strings.Contains(err.Error(), "doesn't settle") {
		t.Errorf("error %v, want that the observer doesn't settle", err)
	}

}

// Two observers that undo each other never settle.
func TestObserversThatCycleFail(t *testing.T) {
	d := dieappsv1.DeploymentBlank.DieMutable()
	dietesting.NoError(t, die.ObserveAt(d, labels, func(l *fields.MapMut[string, string], c die.Change[*fields.Map[string, string]]) error {
		return root(l).Metadata().Annotations().Set("a", c.After.At("l")+"a")
	}))
	dietesting.NoError(t, die.ObserveAt(d, annotations, func(a *fields.MapMut[string, string], c die.Change[*fields.Map[string, string]]) error {
		return root(a).Metadata().Labels().Set("l", c.After.At("a")+"l")
	}))
	if err := d.Metadata().Labels().Set("l", ""); err == nil || !strings.Contains(err.Error(), "settle") {
		t.Errorf("error %v, want that observers don't settle", err)
	}
}

// An immutable die's observers are called within the dies derived from it.
func TestObserversOfAnImmutableDie(t *testing.T) {
	calls := 0
	blank := dieappsv1.DeploymentBlank
	observed := blank.DieStamp(func(m *dieappsv1.DeploymentDieMut) {
		dietesting.NoError(t, die.ObserveAt(m, minReady, func(r fields.ValueMut[int32], c die.Change[fields.Value[int32]]) error {
			calls++
			return root(r).Metadata().Labels().Set("size", fmt.Sprint(c.After.Get()))
		}))
	})

	// the die it was derived from isn't observed
	blank.DieStamp(func(d *dieappsv1.DeploymentDieMut) { dietesting.NoError(t, d.Spec().MinReadySeconds().Set(1)) })
	if calls != 0 {
		t.Errorf("the blank is observed")
	}

	stamped := observed.DieStamp(func(d *dieappsv1.DeploymentDieMut) { dietesting.NoError(t, d.Spec().MinReadySeconds().Set(2)) })
	if got := stamped.Metadata().Labels().At("size"); got != "2" {
		t.Errorf("size %q", got)
	}
	// and the dies derived from those
	for name, m := range map[string]*dieappsv1.DeploymentDieMut{
		"stamp": stamped.DieMutable(),
		"tare":  stamped.DieTare().DieMutable(),
		"copy":  stamped.DeepCopy().DieMutable(),
	} {
		dietesting.NoError(t, m.Spec().MinReadySeconds().Set(5))
		if got := m.Metadata().Labels().At("size"); got != "5" {
			t.Errorf("%s: size %q", name, got)
		}
	}
	// a stamped die's observer last saw what its stamp left
	if got := stamped.Metadata().Labels().At("size"); got != "2" {
		t.Errorf("a derived die's write reached its source: size %q", got)
	}
}

// A copy of a die within another keeps the observers registered at or below
// it, and a mutable die's feed notifies them.
func TestObserversOfACopyOfAChild(t *testing.T) {
	d := dieappsv1.DeploymentBlank.DieMutable()
	paused := 0
	dietesting.NoError(t, die.ObserveAt(d.Spec(), func(m *dieappsv1.DeploymentSpecDieMut) fields.ValueMut[bool] { return m.Paused() },
		func(p fields.ValueMut[bool], c die.Change[fields.Value[bool]]) error {
			paused++
			return nil
		}))
	dietesting.NoError(t, die.ObserveAt(d, minReady, func(r fields.ValueMut[int32], c die.Change[fields.Value[int32]]) error {
		t.Error("the deployment's observer is called within a copy of its spec")
		return nil
	}))

	spec := d.Spec().DieImmutable().DieMutable()
	dietesting.NoError(t, spec.Paused().Set(true))
	if paused != 1 {
		t.Errorf("called %d times within the copy, want 1", paused)
	}
	dietesting.NoError(t, spec.MinReadySeconds().Set(1))

	dietesting.NoError(t, spec.DieFeed(&appsv1.DeploymentSpec{Paused: false}))
	if paused != 2 {
		t.Errorf("called %d times after a feed, want 2", paused)
	}
}

// changesOf observes a deployment's minReadySeconds, recording each change,
// and labels it with them, as an observer that writes would.
func changesOf(t *testing.T, d *dieappsv1.DeploymentDieMut, fail int32) *[]string {
	t.Helper()
	var changes []string
	dietesting.NoError(t, die.ObserveAt(d, minReady, func(r fields.ValueMut[int32], c die.Change[fields.Value[int32]]) error {
		if c.After.Get() == fail {
			return fmt.Errorf("%d is refused", fail)
		}
		changes = append(changes, show(c))
		return root(r).Metadata().Labels().Set("min-ready", fmt.Sprint(c.After.Get()))
	}))
	return &changes
}

// An atomic stamp notifies after each fn, of what it wrote.
func TestDieStampAtomicNotifiesAfterEachFunc(t *testing.T) {
	d := dieappsv1.DeploymentBlank.DieMutable()
	changes := changesOf(t, d, -1)
	set := func(vs ...int32) func(m *dieappsv1.DeploymentDieMut) error {
		return func(m *dieappsv1.DeploymentDieMut) error {
			for _, v := range vs {
				if err := m.Spec().MinReadySeconds().Set(v); err != nil {
					return err
				}
			}
			// observers haven't run yet for this fn's writes
			if got := m.Metadata().Labels().At("min-ready"); got == fmt.Sprint(vs[len(vs)-1]) {
				t.Errorf("observed within the fn: %s", got)
			}
			return nil
		}
	}
	_, err := d.DieStampAtomic(set(1, 3), set(4), set(5))
	dietesting.NoError(t, err)
	if got := fmt.Sprint(*changes); got != "[{0 3 false true} {3 4 true true} {4 5 true true}]" || d.Metadata().Labels().At("min-ready") != "5" {
		t.Errorf("changes %s, label %q", got, d.Metadata().Labels().At("min-ready"))
	}
}

// What observers write as an atomic stamp notifies them notifies as it is
// written.
func TestDieStampAtomicObserversWriteImmediately(t *testing.T) {
	d := dieappsv1.DeploymentBlank.DieMutable()
	dietesting.NoError(t, die.ObserveAt(d, minReady, func(r fields.ValueMut[int32], c die.Change[fields.Value[int32]]) error {
		m := root(r)
		if err := m.Metadata().Labels().Set("min-ready", fmt.Sprint(c.After.Get())); err != nil {
			return err
		}
		// the label's observer has run
		if got := m.Metadata().Annotations().At("label"); got != fmt.Sprint(c.After.Get()) {
			t.Errorf("annotation %q, want the label's observer to have written %d", got, c.After.Get())
		}
		return nil
	}))
	dietesting.NoError(t, die.ObserveAt(d, labels, func(l *fields.MapMut[string, string], c die.Change[*fields.Map[string, string]]) error {
		return root(l).Metadata().Annotations().Set("label", c.After.At("min-ready"))
	}))
	_, err := d.DieStampAtomic(func(m *dieappsv1.DeploymentDieMut) error { return m.Spec().MinReadySeconds().Set(7) })
	dietesting.NoError(t, err)
}

// A write whose observers fail returns their error, and is kept, with what
// they wrote; within an atomic stamp, both roll back.
func TestObserversThatFail(t *testing.T) {
	d := dieappsv1.DeploymentBlank.DieMutable()
	changesOf(t, d, -1)
	dietesting.NoError(t, die.ObserveAt(d, labels, func(l *fields.MapMut[string, string], c die.Change[*fields.Map[string, string]]) error {
		if c.After.At("min-ready") == "13" {
			return fmt.Errorf("13 is refused")
		}
		return nil
	}))
	dietesting.NoError(t, d.Spec().MinReadySeconds().Set(2))
	before, intent := d.DieImmutable(), d.DieIntent()

	_, err := d.DieStampAtomic(func(m *dieappsv1.DeploymentDieMut) error { return m.Spec().MinReadySeconds().Set(13) })
	if err == nil || !strings.Contains(err.Error(), "13 is refused") {
		t.Errorf("error %v", err)
	}
	if !d.DieImmutable().Equal(before) || d.DieIntent() != intent {
		t.Errorf("the atomic stamp isn't rolled back:\n%s", d.DieIntent())
	}

	if err := d.Spec().MinReadySeconds().Set(13); err == nil || !strings.Contains(err.Error(), "13 is refused") {
		t.Errorf("error %v", err)
	}
	if d.Spec().MinReadySeconds().Get() != 13 || d.Metadata().Labels().At("min-ready") != "13" {
		t.Errorf("the write isn't kept: minReadySeconds %d, label %q", d.Spec().MinReadySeconds().Get(), d.Metadata().Labels().At("min-ready"))
	}
}

// An atomic stamp that fails, or whose observers do, leaves the die as it was.
func TestDieStampAtomicRollsBack(t *testing.T) {
	d := dieappsv1.DeploymentBlank.DieMutable()
	changes := changesOf(t, d, 13)
	dietesting.NoError(t, d.Spec().MinReadySeconds().Set(2))
	before, intent := d.DieImmutable(), d.DieIntent()

	for name, fn := range map[string]func(m *dieappsv1.DeploymentDieMut) error{
		"a fn fails": func(m *dieappsv1.DeploymentDieMut) error {
			if err := m.Spec().MinReadySeconds().Set(5); err != nil {
				return err
			}
			return fmt.Errorf("give up")
		},
		"an observer fails": func(m *dieappsv1.DeploymentDieMut) error {
			if err := m.Spec().Paused().Set(true); err != nil {
				return err
			}
			return m.Spec().MinReadySeconds().Set(13)
		},
		"an observer registered within fails": func(m *dieappsv1.DeploymentDieMut) error {
			if err := die.ObserveAt(m, labels, func(*fields.MapMut[string, string], die.Change[*fields.Map[string, string]]) error {
				return fmt.Errorf("no labels")
			}); err != nil {
				return err
			}
			return m.Spec().MinReadySeconds().Set(4)
		},
	} {
		if _, err := d.DieStampAtomic(fn); err == nil {
			t.Errorf("%s: no error", name)
		}
		if !d.DieImmutable().Equal(before) || d.DieIntent() != intent {
			t.Errorf("%s: not rolled back:\n%s", name, d.DieIntent())
		}
	}

	// what the observers last saw is rolled back too, and the one
	// registered within a stamp that failed is gone; what an observer did
	// outside the die, as recording its changes here, isn't undone
	*changes = nil
	dietesting.NoError(t, d.Spec().MinReadySeconds().Set(6))
	if got := fmt.Sprint(*changes); got != "[{2 6 true true}]" {
		t.Errorf("changes %s", got)
	}
}

// An atomic stamp within another that fails rolls back only its own changes.
func TestDieStampAtomicNests(t *testing.T) {
	d := dieappsv1.DeploymentBlank.DieMutable()
	changes := changesOf(t, d, -1)
	_, err := d.DieStampAtomic(func(m *dieappsv1.DeploymentDieMut) error {
		if err := m.Spec().MinReadySeconds().Set(1); err != nil {
			return err
		}
		if _, err := m.DieStampAtomic(func(m *dieappsv1.DeploymentDieMut) error {
			if err := m.Spec().Paused().Set(true); err != nil {
				return err
			}
			return fmt.Errorf("not paused")
		}); err == nil {
			t.Error("the inner stamp didn't fail")
		}
		return m.Spec().MinReadySeconds().Set(2)
	})
	dietesting.NoError(t, err)
	if d.Spec().Paused().Get() || d.Spec().MinReadySeconds().Get() != 2 {
		t.Errorf("paused %t, minReadySeconds %d", d.Spec().Paused().Get(), d.Spec().MinReadySeconds().Get())
	}
	if got := fmt.Sprint(*changes); got != "[{0 2 false true}]" {
		t.Errorf("changes %s", got)
	}
}

// An immutable die's atomic stamp returns the die itself when it fails.
func TestDieStampAtomicOfAnImmutableDie(t *testing.T) {
	calls := 0
	d := dieappsv1.DeploymentBlank.DieStamp(func(m *dieappsv1.DeploymentDieMut) {
		dietesting.NoError(t, die.ObserveAt(m, minReady, func(r fields.ValueMut[int32], c die.Change[fields.Value[int32]]) error {
			calls++
			return nil
		}))
	})
	set := func(v int32, err error) func(m *dieappsv1.DeploymentDieMut) error {
		return func(m *dieappsv1.DeploymentDieMut) error {
			if e := m.Spec().MinReadySeconds().Set(v); e != nil {
				return e
			}
			return err
		}
	}
	got, err := d.DieStampAtomic(set(1, nil), set(2, nil))
	dietesting.NoError(t, err)
	// notified after each fn
	if got.Spec().MinReadySeconds().Get() != 2 || calls != 2 {
		t.Errorf("minReadySeconds %d after %d calls", got.Spec().MinReadySeconds().Get(), calls)
	}
	failed, err := got.DieStampAtomic(set(3, fmt.Errorf("no")))
	if err == nil || failed.Spec().MinReadySeconds().Get() != 2 {
		t.Errorf("error %v, minReadySeconds %d", err, failed.Spec().MinReadySeconds().Get())
	}
}

// A die observes itself: a die within another observes that part of its
// root, and is given itself, as the mutation left it.
func TestDieObserve(t *testing.T) {
	d := dieappsv1.DeploymentBlank.DieMutable()
	var changes []string
	dietesting.NoError(t, d.Spec().DieObserve(func(spec *dieappsv1.DeploymentSpecDieMut, c die.Change[*dieappsv1.DeploymentSpecDie]) error {
		before, after := c.Before.MinReadySeconds().Get(), c.After.MinReadySeconds().Get()
		changes = append(changes, fmt.Sprintf("%d->%d", before, after))
		if spec.MinReadySeconds().Get() != after {
			t.Errorf("the die reads %d, the change is to %d", spec.MinReadySeconds().Get(), after)
		}
		// the change's dies are within their roots as they were
		if name := c.After.DieRoot().(*dieappsv1.DeploymentDie).Metadata().Name().Get(); name != "web" {
			t.Errorf("the change's root is named %q", name)
		}
		return spec.Paused().Set(after > 5)
	}))
	dietesting.NoError(t, d.Metadata().Name().Set("web")) // not the spec
	dietesting.NoError(t, d.Spec().MinReadySeconds().Set(3))
	dietesting.NoError(t, d.Spec().MinReadySeconds().Set(9))
	// its own write to paused is a change too, which it settles
	if got := strings.Join(changes, " "); got != "0->3 3->9 9->9" || !d.Spec().Paused().Get() {
		t.Errorf("changes %s, paused %t", got, d.Spec().Paused().Get())
	}
}

// An immutable root observes itself within the dies derived from it.
func TestDieObserveOfAnImmutableDie(t *testing.T) {
	observed := dieappsv1.DeploymentBlank.DieObserve(func(m *dieappsv1.DeploymentDieMut, c die.Change[*dieappsv1.DeploymentDie]) error {
		return m.Metadata().Labels().Set("min-ready", fmt.Sprint(c.After.Spec().MinReadySeconds().Get()))
	})
	d := observed.DieStamp(func(m *dieappsv1.DeploymentDieMut) {
		dietesting.NoError(t, m.Spec().MinReadySeconds().Set(4))
	})
	if got := d.Metadata().Labels().At("min-ready"); got != "4" {
		t.Errorf("label %q", got)
	}
}

func panicsWith(t *testing.T, want string, fn func()) {
	t.Helper()
	defer func() {
		t.Helper()
		if r := recover(); r == nil || !strings.Contains(fmt.Sprint(r), want) {
			t.Errorf("panic %v, want it to contain %q", r, want)
		}
	}()
	fn()
}

// An immutable die within another can't be observed, as its copy would be
// detached from its root; it is observed within a stamp of its root.
func TestDieObserveOfAnImmutableChild(t *testing.T) {
	spec := func(s *dieappsv1.DeploymentSpecDieMut, c die.Change[*dieappsv1.DeploymentSpecDie]) error {
		return s.Paused().Set(c.After.MinReadySeconds().Get() > 5)
	}
	panicsWith(t, "DieObserve on an immutable die at .spec would observe a detached copy", func() {
		dieappsv1.DeploymentBlank.Spec().DieObserve(spec)
	})

	observed := dieappsv1.DeploymentBlank.DieStamp(func(m *dieappsv1.DeploymentDieMut) {
		dietesting.NoError(t, m.Spec().DieObserve(spec))
	})
	d := observed.DieStamp(func(m *dieappsv1.DeploymentDieMut) {
		dietesting.NoError(t, m.Spec().MinReadySeconds().Set(9))
	})
	if !d.Spec().Paused().Get() {
		t.Error("the spec, observed within a stamp of its root, isn't observed")
	}
}

// Each element of a list map is observed, matched by key: only the elements
// that change are called, each with its own change.
func TestDieObserveEach(t *testing.T) {
	d := dieappsv1.DeploymentBlank.DieMutable()
	containers := d.Spec().Template().Spec().Containers()
	var calls []string
	dietesting.NoError(t, containers.DieObserveEach(func(c *diecorev1.ContainerDieMut, ch die.Change[*diecorev1.ContainerDie]) error {
		// a removed element's die reads as absent; its change has its key
		had, has := ch.Before.DieHas(), ch.After.DieHas()
		name := ch.After.Name().Get()
		if !has {
			name = ch.Before.Name().Get()
		}
		calls = append(calls, fmt.Sprintf("%s:%t/%t:%s->%s", name, had, has, ch.Before.Image().Get(), ch.After.Image().Get()))
		// an element's observer writes through its die
		if has && ch.After.ImagePullPolicy().Get() == "" {
			return c.ImagePullPolicy().Set(corev1.PullIfNotPresent)
		}
		return nil
	}))
	set := func(name, image string) {
		t.Helper()
		dietesting.NoError(t, containers.At(name).Image().Set(image))
	}

	set("app", "app:1")
	set("sidecar", "proxy:1")
	calls = nil
	set("app", "app:2")
	if got := strings.Join(calls, " "); got != "app:true/true:app:1->app:2" {
		t.Errorf("calls %s, want only app", got)
	}
	if got := containers.At("sidecar").ImagePullPolicy().Get(); got != corev1.PullIfNotPresent {
		t.Errorf("sidecar's pull policy %q", got)
	}

	calls = nil
	dietesting.NoError(t, containers.Remove("sidecar"))
	if got := strings.Join(calls, " "); got != "sidecar:true/false:proxy:1->" {
		t.Errorf("calls %s, want sidecar removed", got)
	}

	// within an atomic stamp, each element is called once for all its fn
	// wrote, and init again for the pull policy its observer set
	calls = nil
	_, err := d.DieStampAtomic(func(m *dieappsv1.DeploymentDieMut) error {
		cs := m.Spec().Template().Spec().Containers()
		return errors.Join(cs.At("app").Image().Set("app:3"), cs.At("app").Image().Set("app:4"), cs.At("init").Image().Set("init:1"))
	})
	dietesting.NoError(t, err)
	if got := strings.Join(calls, " "); got != "app:true/true:app:2->app:4 init:false/true:->init:1 init:true/true:init:1->init:1" {
		t.Errorf("calls %s", got)
	}
}

// The list of an immutable root is observed within a stamp of the root.
func TestDieObserveEachOfAnImmutableDie(t *testing.T) {
	observed := dieappsv1.DeploymentBlank.DieStamp(func(m *dieappsv1.DeploymentDieMut) {
		dietesting.NoError(t, m.Spec().Template().Spec().Containers().DieObserveEach(func(c *diecorev1.ContainerDieMut, ch die.Change[*diecorev1.ContainerDie]) error {
			if !ch.After.DieHas() {
				return nil
			}
			return c.ImagePullPolicy().Set(corev1.PullAlways)
		}))
	})
	d := observed.DieStamp(func(m *dieappsv1.DeploymentDieMut) {
		dietesting.NoError(t, m.Spec().Template().Spec().Containers().At("app").Image().Set("app:1"))
	})
	if got := d.Spec().Template().Spec().Containers().At("app").ImagePullPolicy().Get(); got != corev1.PullAlways {
		t.Errorf("pull policy %q", got)
	}
}

// A change's dies are frozen: within their roots as they were, and unchanged
// by later writes.
func TestChangesAreOfTheirTime(t *testing.T) {
	d := dieappsv1.DeploymentBlank.DieMutable()
	dietesting.NoError(t, d.Metadata().Name().Set("a"))
	var changes []die.Change[*dieappsv1.DeploymentSpecDie]
	dietesting.NoError(t, d.Spec().DieObserve(func(_ *dieappsv1.DeploymentSpecDieMut, c die.Change[*dieappsv1.DeploymentSpecDie]) error {
		changes = append(changes, c)
		return nil
	}))
	_, err := d.DieStampAtomic(func(m *dieappsv1.DeploymentDieMut) error {
		return errors.Join(m.Metadata().Name().Set("b"), m.Spec().MinReadySeconds().Set(1))
	})
	dietesting.NoError(t, err)
	dietesting.NoError(t, d.Spec().MinReadySeconds().Set(2))

	if len(changes) != 2 {
		t.Fatalf("%d changes", len(changes))
	}
	name := func(spec *dieappsv1.DeploymentSpecDie) string {
		return spec.DieRoot().(*dieappsv1.DeploymentDie).Metadata().Name().Get()
	}
	first := changes[0]
	if name(first.Before) != "a" || name(first.After) != "b" {
		t.Errorf("roots named %q then %q, want a then b", name(first.Before), name(first.After))
	}
	if got := first.After.MinReadySeconds().Get(); got != 1 {
		t.Errorf("the first change's after reads %d, since written to 2", got)
	}
	if before, after := changes[1].Before.MinReadySeconds().Get(), changes[1].After.MinReadySeconds().Get(); before != 1 || after != 2 {
		t.Errorf("second change %d->%d", before, after)
	}
}
