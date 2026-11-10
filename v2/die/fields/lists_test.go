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

package fields_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"

	dierbacv1 "reconciler.io/dies/v2/apis/rbac/v1"
	"reconciler.io/dies/v2/die"
	dietesting "reconciler.io/dies/v2/die/testing"
	"reconciler.io/dies/v2/die/validation"
)

func rule(t *testing.T, verbs ...string) *dierbacv1.PolicyRuleDie {
	t.Helper()
	return dierbacv1.PolicyRuleBlank.DieStamp(func(d *dierbacv1.PolicyRuleDieMut) {
		dietesting.NoError(t, errors.Join(
			d.APIGroups().SetValue("apps"),
			d.Resources().SetValue("deployments"),
			d.Verbs().SetValue(verbs...),
		))
	})
}

// A list is set from dies and handles of either mutability, mixed, or from
// values, and appended to alike.
func TestListSetAndAppend(t *testing.T) {
	role := dierbacv1.ClusterRoleBlank.DieMutable()
	other := dierbacv1.ClusterRoleBlank.DieStamp(func(d *dierbacv1.ClusterRoleDieMut) {
		dietesting.NoError(t, d.Rules().SetValue(rbacv1.PolicyRule{Verbs: []string{"watch"}}))
	})
	dietesting.NoError(t, role.Rules().Set(rule(t, "get"), rule(t, "list").DieMutable(), other.Rules().At(0)))
	if got := verbsOf(role.Rules().Get()); !reflect.DeepEqual(got, []string{"get", "list", "watch"}) {
		t.Errorf("verbs %v", got)
	}
	dietesting.NoError(t, role.Rules().Append(rule(t, "create")))
	dietesting.NoError(t, role.Rules().AppendValue(rbacv1.PolicyRule{Verbs: []string{"update"}}, rbacv1.PolicyRule{Verbs: []string{"patch"}}))
	if got := verbsOf(role.Rules().Get()); !reflect.DeepEqual(got, []string{"get", "list", "watch", "create", "update", "patch"}) {
		t.Errorf("verbs %v", got)
	}
	dietesting.NoError(t, role.Rules().SetValue(rbacv1.PolicyRule{Verbs: []string{"delete"}}))
	if got := verbsOf(role.Rules().Get()); !reflect.DeepEqual(got, []string{"delete"}) {
		t.Errorf("verbs %v", got)
	}

	// appending to an absent list creates it; setting no items empties it
	fresh := dierbacv1.ClusterRoleBlank.DieMutable()
	dietesting.NoError(t, fresh.Rules().AppendValue(rbacv1.PolicyRule{Verbs: []string{"get"}}))
	if fresh.Rules().Len() != 1 {
		t.Errorf("len %d, want 1", fresh.Rules().Len())
	}
	dietesting.NoError(t, fresh.Rules().Set())
	if !fresh.Rules().Has() || fresh.Rules().Len() != 0 {
		t.Errorf("has %v, len %d; want present and empty", fresh.Rules().Has(), fresh.Rules().Len())
	}
	patch, err := die.ReleaseMergePatch(fresh)
	dietesting.NoError(t, err)
	sameJSON(t, patch, `{"rules": []}`)

	// a scalar list, by its values or other lists' elements
	r := dierbacv1.PolicyRuleBlank.DieMutable()
	dietesting.NoError(t, r.Verbs().SetValue("get", "list"))
	dietesting.NoError(t, r.Resources().Set(r.Verbs().At(1), r.Verbs().At(0)))
	if got := r.Resources().Get(); !reflect.DeepEqual(got, []string{"list", "get"}) {
		t.Errorf("resources %v", got)
	}
}

func verbsOf(rules []rbacv1.PolicyRule) []string {
	var verbs []string
	for _, r := range rules {
		verbs = append(verbs, r.Verbs...)
	}
	return verbs
}

// An element is reached by its index, and written through its die or
// handle: the list is written whole, as the API server replaces it.
func TestListAt(t *testing.T) {
	role := dierbacv1.ClusterRoleBlank.DieMutable()
	dietesting.NoError(t, role.Rules().Set(rule(t, "get"), rule(t, "list")))
	role.DieTare()

	dietesting.NoError(t, role.Rules().At(1).Verbs().AppendValue("watch"))
	dietesting.NoError(t, role.Rules().At(0).Verbs().At(0).Set("create"))
	if got := verbsOf(role.Rules().Get()); !reflect.DeepEqual(got, []string{"create", "list", "watch"}) {
		t.Errorf("verbs %v", got)
	}
	if got := role.DieReadOnly().Rules().At(1).Verbs().At(1).Get(); got != "watch" {
		t.Errorf("verb %q, want watch", got)
	}
	// the whole list is written
	patch, err := die.ReleaseMergePatch(role)
	dietesting.NoError(t, err)
	sameJSON(t, patch, `{"rules": [
		{"apiGroups": ["apps"], "resources": ["deployments"], "verbs": ["create"]},
		{"apiGroups": ["apps"], "resources": ["deployments"], "verbs": ["list", "watch"]}
	]}`)

	var indexes []int
	for i, r := range role.DieReadOnly().Rules().All() {
		indexes = append(indexes, i)
		if !r.Verbs().Has() {
			t.Errorf("rule %d has no verbs", i)
		}
	}
	if !reflect.DeepEqual(indexes, []int{0, 1}) {
		t.Errorf("indexes %v", indexes)
	}
	// reading past the end is the zero value
	if role.Rules().At(5).Verbs().Len() != 0 {
		t.Error("an element past the end has verbs")
	}
}

// An element's die and handle are at its index, for validation errors.
func TestListElementPath(t *testing.T) {
	role := dierbacv1.ClusterRoleBlank.DieStamp(func(d *dierbacv1.ClusterRoleDieMut) {
		dietesting.NoError(t, d.Rules().Set(rule(t, "get"), rule(t, "frobnicate")))
	})
	err := validation.Invalid(role.Rules().At(1).Verbs().At(0), "is not a verb")
	want := field.Invalid(field.NewPath("rules").Index(1).Child("verbs").Index(0), "frobnicate", "is not a verb")
	if err.Error() != want.Error() {
		t.Errorf("got %v\nwant %v", err, want)
	}
}

// An element is removed by its index; one the list doesn't have is an
// error.
func TestListRemove(t *testing.T) {
	r := dierbacv1.PolicyRuleBlank.DieMutable()
	dietesting.NoError(t, r.Verbs().SetValue("get", "list", "watch"))
	dietesting.NoError(t, r.Verbs().Remove(1))
	if got := r.Verbs().Get(); !reflect.DeepEqual(got, []string{"get", "watch"}) {
		t.Errorf("verbs %v", got)
	}
	if err := r.Verbs().Remove(2); err == nil || !strings.Contains(err.Error(), "the list has 2 elements") {
		t.Errorf("error %v", err)
	}
	dietesting.NoError(t, r.Verbs().Clear())
	if r.Verbs().Has() {
		t.Error("cleared list is present")
	}
}

// What a list is set to is a copy of its items.
func TestListItemsAreCopies(t *testing.T) {
	src := rule(t, "get").DieMutable()
	role := dierbacv1.ClusterRoleBlank.DieMutable()
	dietesting.NoError(t, role.Rules().Set(src))
	dietesting.NoError(t, src.Verbs().At(0).Set("changed"))
	values := []string{"a"}
	dietesting.NoError(t, role.Rules().At(0).Resources().SetValue(values...))
	values[0] = "changed"
	r := role.Rules().At(0)
	if r.Verbs().At(0).Get() != "get" || r.Resources().At(0).Get() != "a" {
		t.Errorf("the list shares its items: %v", role.Rules().Get())
	}
	// the list's own value is a copy too
	got := role.Rules().Get()
	got[0].Verbs[0] = "changed"
	if role.Rules().At(0).Verbs().At(0).Get() != "get" {
		t.Error("Get shares the list")
	}
}
