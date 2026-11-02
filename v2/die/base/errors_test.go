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
	"encoding/json"
	"errors"
	"testing"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	"reconciler.io/dies/v2/die/base"
)

// A fluent method panics with a die's error; its E variant returns it.

func dieError(t *testing.T, fn func()) *base.Error {
	t.Helper()
	var got *base.Error
	func() {
		defer func() {
			r := recover()
			e, ok := r.(*base.Error)
			if !ok {
				t.Fatalf("recovered %#v, want a *base.Error", r)
			}
			got = e
		}()
		fn()
	}()
	return got
}

func TestFluentAndEVariants(t *testing.T) {
	blank := dieappsv1.DeploymentBlank

	// a write to an immutable die through an interface it implements
	panicked := dieError(t, func() { blank.SetName("web") })
	if !errors.Is(panicked, base.ErrImmutable) {
		t.Errorf("panicked with %v", panicked)
	}
	// on a mutable die, the E variant succeeds, as the fluent method does
	m := blank.DieMutable()
	if _, err := m.Spec().Replicas().SetE(2); err != nil || m.Spec().Replicas().Get() != 2 {
		t.Errorf("SetE: %v, replicas %d", err, m.Spec().Replicas().Get())
	}

	// a cause, such as JSON's, is wrapped
	_, err := blank.DieMutable().DieFeedJSONE([]byte(`{"spec": "not an object"}`))
	var typeErr *json.UnmarshalTypeError
	if !errors.As(err, &typeErr) {
		t.Errorf("DieFeedJSONE returned %v, want it to wrap a *json.UnmarshalTypeError", err)
	}
	if _, err := blank.Spec().DieReleaseApplyPatchE(); err == nil {
		t.Errorf("DieReleaseApplyPatchE on a child returned no error")
	}
}

// DieWithE returns the first error of its functions, such as an E method's,
// with the die as they left it.
func TestDieWithEReturnsErrors(t *testing.T) {
	got, err := dieappsv1.DeploymentBlank.DieWithE(
		func(d *dieappsv1.DeploymentDieMut) error {
			_, err := d.Spec().Replicas().SetE(2)
			return err
		},
		func(d *dieappsv1.DeploymentDieMut) error {
			_, err := d.DieFeedJSONE([]byte(`{"spec": "not an object"}`))
			return err
		},
		func(d *dieappsv1.DeploymentDieMut) error { t.Error("ran after an error"); return nil },
	)
	if _, ok := errors.AsType[*base.Error](err); !ok {
		t.Errorf("error %T %v, want a *base.Error", err, err)
	}
	if got.Spec().Replicas().Get() != 2 {
		t.Errorf("replicas %d, want the die as the functions left it", got.Spec().Replicas().Get())
	}

	// on a die within an immutable one, it runs on a copy, a new root
	spec, err := dieappsv1.DeploymentBlank.Spec().DieWithE(func(d *dieappsv1.DeploymentSpecDieMut) error {
		_, err := d.Replicas().SetE(3)
		return err
	})
	if err != nil || spec.Replicas().Get() != 3 || len(spec.DiePath()) != 0 || dieappsv1.DeploymentBlank.Spec().Replicas().Has() {
		t.Errorf("DieWithE on a child of an immutable die: %v, replicas %d at %s", err, spec.Replicas().Get(), spec.DiePath())
	}
}

// A panic that isn't a die's error, such as the caller's own, isn't
// recovered.
func TestEVariantsDontRecoverOtherPanics(t *testing.T) {
	defer func() {
		if r := recover(); r != "the caller's own" {
			t.Errorf("recovered %#v", r)
		}
	}()
	_, _ = dieappsv1.DeploymentBlank.DieMutable().DieStampE(func(*dieappsv1.Deployment) { panic("the caller's own") })
	t.Error("didn't panic")
}
