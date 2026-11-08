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
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	diecorev1 "reconciler.io/dies/v2/apis/core/v1"
	diemetav1 "reconciler.io/dies/v2/apis/meta/v1"
	dierbacv1 "reconciler.io/dies/v2/apis/rbac/v1"
	"reconciler.io/dies/v2/die"
	dietesting "reconciler.io/dies/v2/die/testing"
)

func sameJSON(t *testing.T, got []byte, want string) {
	t.Helper()
	var g, w any
	if err := json.Unmarshal(got, &g); err != nil {
		t.Fatalf("%s: %v", got, err)
	}
	if err := json.Unmarshal([]byte(want), &w); err != nil {
		t.Fatalf("%s: %v", want, err)
	}
	if !reflect.DeepEqual(g, w) {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

// current is a deployment read from the API server, with labels, replicas
// and a container.
func current(t *testing.T) *dieappsv1.DeploymentDieMut {
	d := dieappsv1.DeploymentBlank.DieMutable()
	dietesting.NoError(t, d.DieFeed(&appsv1.Deployment{}))
	dietesting.NoError(t, json.Unmarshal([]byte(`{
		"metadata": {"name": "web", "labels": {"app": "web", "stale": "yes"}},
		"spec": {"replicas": 3, "template": {"spec": {"containers": [{"name": "app", "image": "app:1", "ports": [{"containerPort": 80}]}]}}}
	}`), d))
	return d
}

// A map is replaced whole: keys the source doesn't have are removed, and
// the patch says so.
func TestMapReplace(t *testing.T) {
	c := current(t)
	desired := dieappsv1.DeploymentBlank.DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietesting.NoError(t, d.Metadata().Labels().Set("app", "web"))
		dietesting.NoError(t, d.Metadata().Labels().Set("tier", "frontend"))
	})
	dietesting.NoError(t, c.Metadata().Labels().Replace(desired.Metadata().Labels()))
	if got, want := c.Metadata().Labels().DieReadOnly().Keys(), []string{"app", "tier"}; !reflect.DeepEqual(got, want) {
		t.Errorf("keys %v, want %v", got, want)
	}
	patch, err := die.ReleaseMergePatch(c)
	dietesting.NoError(t, err)
	sameJSON(t, patch, `{"metadata": {"labels": {"stale": null, "tier": "frontend"}}}`)

	// from an absent map, the map is removed
	dietesting.NoError(t, c.Metadata().Labels().Replace(dieappsv1.DeploymentBlank.Metadata().Labels()))
	if c.Metadata().Labels().Has("app") || c.DieImmutable().Metadata().Labels().Len() != 0 {
		t.Errorf("labels %v, want none", c.Metadata().Labels().DieReadOnly().Keys())
	}
	patch, err = die.ReleaseMergePatch(c)
	dietesting.NoError(t, err)
	sameJSON(t, patch, `{"metadata": {"labels": null}}`)
}

// A child die is replaced whole from another root's, as the issue's
// aggregation rule is, or cleared when the source has none.
func TestDieReplace(t *testing.T) {
	selector := diemetav1.LabelSelector{MatchLabels: map[string]string{"aggregate": "true"}}
	desired := dierbacv1.ClusterRoleBlank.DieStamp(func(d *dierbacv1.ClusterRoleDieMut) {
		dietesting.NoError(t, d.AggregationRule().ClusterRoleSelectors().Set([]diemetav1.LabelSelector{selector}))
	})
	c := dierbacv1.ClusterRoleBlank.DieMutable()
	dietesting.NoError(t, c.AggregationRule().DieReplace(desired.AggregationRule()))
	if got := c.DieRelease().AggregationRule; !reflect.DeepEqual(got, &rbacv1.AggregationRule{ClusterRoleSelectors: []diemetav1.LabelSelector{selector}}) {
		t.Errorf("aggregation rule %v", got)
	}
	// the source is copied, not shared
	dietesting.NoError(t, desired.DieMutable().AggregationRule().ClusterRoleSelectors().Clear())
	if c.AggregationRule().ClusterRoleSelectors().Get() == nil {
		t.Error("the replaced die shares its source's value")
	}

	dietesting.NoError(t, c.AggregationRule().DieReplace(dierbacv1.ClusterRoleBlank.AggregationRule()))
	if got := c.DieRelease().AggregationRule; got != nil {
		t.Errorf("aggregation rule %v, want none", got)
	}
}

// A leaf, an optional leaf and a list map are replaced from dies of either
// mutability, and from elsewhere in the same root.
func TestHandleReplace(t *testing.T) {
	c := current(t)
	desired := dieappsv1.DeploymentBlank.DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietesting.NoError(t, d.Spec().Replicas().Set(5))
		d.Spec().Template().Spec().Containers().At("app").DieStamp(func(c *diecorev1.ContainerDieMut) {
			dietesting.NoError(t, c.Image().Set("app:2"))
			dietesting.NoError(t, c.Ports().At(diecorev1.ContainerPortKey{ContainerPort: 8080}).Name().Set("http"))
		})
	}).DieMutable()

	dietesting.NoError(t, c.Spec().Replicas().Replace(desired.Spec().Replicas()))
	if got := c.Spec().Replicas().Get(); got != 5 {
		t.Errorf("replicas %d, want 5", got)
	}
	dietesting.NoError(t, c.Spec().Replicas().Replace(dieappsv1.DeploymentBlank.Spec().Replicas()))
	if c.Spec().Replicas().Has() {
		t.Errorf("replicas %d, want none", c.Spec().Replicas().Get())
	}

	app := c.Spec().Template().Spec().Containers().At("app")
	desiredApp := desired.Spec().Template().Spec().Containers().At("app")
	dietesting.NoError(t, app.Image().Replace(desiredApp.Image()))
	dietesting.NoError(t, app.Ports().Replace(desiredApp.Ports()))
	if got := app.Image().Get(); got != "app:2" {
		t.Errorf("image %q, want app:2", got)
	}
	if got := app.Ports().DieReadOnly().Keys(); !reflect.DeepEqual(got, []diecorev1.ContainerPortKey{{ContainerPort: 8080, Protocol: diecorev1.ProtocolTCP}}) {
		t.Errorf("ports %v, want only 8080", got)
	}

	// within the root: compatible paths of the same type
	dietesting.NoError(t, c.Metadata().Name().Replace(app.Image()))
	if got := c.Metadata().Name().Get(); got != "app:2" {
		t.Errorf("name %q, want app:2", got)
	}
}

// A root die is replaced as intent, as a child is: its base stays what was
// read, and the intent is the diff from it to the source, field by field,
// as writes of each change would record it.
func TestRootDieReplaceRecordsTheDiff(t *testing.T) {
	c := current(t)
	desired := dieappsv1.DeploymentBlank.DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietesting.NoError(t, d.Metadata().Name().Set("web"))
		dietesting.NoError(t, d.Metadata().Labels().Set("app", "web"))
		dietesting.NoError(t, d.Metadata().Labels().Set("tier", "frontend"))
		dietesting.NoError(t, d.Spec().Replicas().Set(3))
	})
	dietesting.NoError(t, c.DieReplace(desired))
	// what is the same, the name, the app label and the replicas, isn't
	// written; the containers desired doesn't have are cleared
	want := strings.Join([]string{
		`set .kind = "Deployment"`,
		`set .apiVersion = "apps/v1"`,
		`set .metadata.labels["tier"] = "frontend"`,
		`clear .metadata.labels["stale"]`,
		`clear .spec.template.spec.containers`,
	}, "\n")
	if got := c.DieIntent(); got != want {
		t.Errorf("intent:\n%s\nwant:\n%s", got, want)
	}
	patch, err := die.ReleaseMergePatch(c)
	dietesting.NoError(t, err)
	sameJSON(t, patch, `{"apiVersion": "apps/v1", "kind": "Deployment", "metadata": {"labels": {"stale": null, "tier": "frontend"}}, "spec": {"template": {"spec": {"containers": null}}}}`)
	if got := c.DieRelease(); !reflect.DeepEqual(got, desired.DieRelease()) {
		t.Errorf("got %v, want %v", got, desired.DieRelease())
	}
}

// A child die and a handle are replaced as intent, by what differs: a map
// by key, a list map by element, and nothing where they're the same.
func TestReplaceRecordsTheDiff(t *testing.T) {
	c := current(t)
	desired := dieappsv1.DeploymentBlank.DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietesting.NoError(t, d.Metadata().Labels().Set("app", "web"))
		dietesting.NoError(t, d.Metadata().Labels().Set("tier", "frontend"))
		dietesting.NoError(t, d.Spec().Replicas().Set(3))
		d.Spec().Template().Spec().Containers().At("app").DieStamp(func(c *diecorev1.ContainerDieMut) {
			dietesting.NoError(t, c.Image().Set("app:2"))
			dietesting.NoError(t, c.Ports().At(diecorev1.ContainerPortKey{ContainerPort: 80}).Name().Set("http"))
		})
	})
	dietesting.NoError(t, errors.Join(
		c.Metadata().Labels().Replace(desired.Metadata().Labels()),
		c.Spec().Replicas().Replace(desired.Spec().Replicas()),
		c.Spec().Template().Spec().Containers().Replace(desired.Spec().Template().Spec().Containers()),
	))
	want := strings.Join([]string{
		`set .metadata.labels["tier"] = "frontend"`,
		`clear .metadata.labels["stale"]`,
		`set .spec.template.spec.containers[name=app].image = "app:2"`,
		`set .spec.template.spec.containers[name=app].ports[containerPort=80,protocol=TCP].name = "http"`,
		// the current port has no protocol; desired's key defaults it
		`set .spec.template.spec.containers[name=app].ports[containerPort=80,protocol=TCP].protocol = "TCP"`,
	}, "\n")
	if got := c.DieIntent(); got != want {
		t.Errorf("intent:\n%s\nwant:\n%s", got, want)
	}

	// replacing with the same value writes nothing
	again := current(t)
	dietesting.NoError(t, again.Metadata().DieReplace(current(t).Metadata()))
	if got := again.DieIntent(); got != "" {
		t.Errorf("intent of replacing with the same value:\n%s", got)
	}
}

// What is replaced is a copy: later writes to the source, or to the
// destination, aren't seen by the other, however deep they are.
func TestReplaceIsDecoupledFromItsSource(t *testing.T) {
	src := dierbacv1.ClusterRoleBlank.DieMutable()
	dietesting.NoError(t, errors.Join(
		src.Metadata().Labels().Set("app", "web"),
		src.AggregationRule().ClusterRoleSelectors().Set([]diemetav1.LabelSelector{{MatchLabels: map[string]string{"aggregate": "true"}}}),
	))
	dst := dierbacv1.ClusterRoleBlank.DieMutable()
	dietesting.NoError(t, errors.Join(
		dst.Metadata().Labels().Replace(src.Metadata().Labels()),
		dst.AggregationRule().DieReplace(src.AggregationRule()),
	))
	want := dst.DieRelease()

	// writes to the source, through its dies and in place in what it holds
	dietesting.NoError(t, errors.Join(
		src.Metadata().Labels().Set("app", "changed"),
		src.Metadata().Labels().Set("added", "yes"),
		die.StampValue(src.AggregationRule(), func(r *rbacv1.AggregationRule) {
			r.ClusterRoleSelectors[0].MatchLabels["aggregate"] = "changed"
			r.ClusterRoleSelectors = append(r.ClusterRoleSelectors, metav1.LabelSelector{})
		}),
	))
	if got := dst.DieRelease(); !reflect.DeepEqual(got, want) {
		t.Errorf("a write to the source is seen by the destination: got %v, want %v", got, want)
	}

	// and to the destination
	before := src.DieRelease()
	dietesting.NoError(t, errors.Join(
		dst.Metadata().Labels().Set("app", "dst"),
		die.StampValue(dst.AggregationRule(), func(r *rbacv1.AggregationRule) {
			r.ClusterRoleSelectors[0].MatchLabels["aggregate"] = "dst"
		}),
	))
	if got := src.DieRelease(); !reflect.DeepEqual(got, before) {
		t.Errorf("a write to the destination is seen by the source: got %v, want %v", got, before)
	}

	// an immutable source's value is no more shared than a mutable one's
	frozen := src.DieImmutable()
	dst2 := dierbacv1.ClusterRoleBlank.DieMutable()
	dietesting.NoError(t, dst2.AggregationRule().DieReplace(frozen.AggregationRule()))
	dietesting.NoError(t, die.StampValue(dst2.AggregationRule(), func(r *rbacv1.AggregationRule) {
		r.ClusterRoleSelectors[0].MatchLabels["aggregate"] = "dst2"
	}))
	if got := frozen.AggregationRule().ClusterRoleSelectors().Get()[0].MatchLabels["aggregate"]; got != "changed" {
		t.Errorf("a write to the destination is seen by an immutable source: %q", got)
	}
}
