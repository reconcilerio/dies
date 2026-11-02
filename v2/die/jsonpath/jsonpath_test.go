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

package jsonpath_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	upstreamcorev1 "k8s.io/api/core/v1"
	upstreamrbacv1 "k8s.io/api/rbac/v1"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	diecorev1 "reconciler.io/dies/v2/apis/core/v1"
	dierbacv1 "reconciler.io/dies/v2/apis/rbac/v1"
	"reconciler.io/dies/v2/die"
	"reconciler.io/dies/v2/die/jsonpath"
	"reconciler.io/dies/v2/die/unstructured"
)

// comparisons are on numbers
func TestFilterComparisons(t *testing.T) {
	d := dieappsv1.DeploymentBlank.DieWith(func(d *dieappsv1.DeploymentDieMut) {
		ports := d.Spec().Template().Spec().Containers().At("app").Ports()
		ports.At(diecorev1.ContainerPortKey{ContainerPort: 80}).Name().Set("http")
		ports.At(diecorev1.ContainerPortKey{ContainerPort: 8443}).Name().Set("https")
	})
	for expr, want := range map[string][]string{
		"$..ports[?(@.containerPort > 1024)].name": {".spec.template.spec.containers[name=app].ports[containerPort=8443,protocol=TCP].name"},
		"$..ports[?(@.containerPort <= 80)].name":  {".spec.template.spec.containers[name=app].ports[containerPort=80,protocol=TCP].name"},
		"$..ports[?(@.containerPort >= 80)].name":  {".spec.template.spec.containers[name=app].ports[containerPort=80,protocol=TCP].name", ".spec.template.spec.containers[name=app].ports[containerPort=8443,protocol=TCP].name"},
		"$..ports[?(@.name < 'z')].name":           {},
	} {
		if got := paths(t, d, expr); !reflect.DeepEqual(got, want) && !(len(got) == 0 && len(want) == 0) {
			t.Errorf("%s\n got %v\nwant %v", expr, got, want)
		}
	}
}

func web() *dieappsv1.DeploymentDie {
	return dieappsv1.DeploymentBlank.DieWith(func(d *dieappsv1.DeploymentDieMut) {
		d.Metadata().Name().Set("web")
		d.Metadata().Labels().Set("app", "web")
		spec := d.Spec().Template().Spec()
		spec.Containers().At("app").Image().Set("app:1")
		spec.Containers().At("proxy").Image().Set("proxy:1")
		spec.Tolerations().Set([]upstreamcorev1.Toleration{{Key: "a"}, {Key: "b"}})
	}).DieCommit()
}

func paths(t *testing.T, d die.Die, expr string) []string {
	t.Helper()
	ps, err := jsonpath.FindE(d, expr)
	if err != nil {
		t.Fatalf("%s: %v", expr, err)
	}
	out := make([]string, len(ps))
	for i, p := range ps {
		out[i] = p.String()
	}
	return out
}

func assertJSON(t *testing.T, got []byte, err error, want string) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
	var g, w any
	_ = json.Unmarshal(got, &g)
	_ = json.Unmarshal([]byte(want), &w)
	if !reflect.DeepEqual(g, w) {
		t.Errorf("got %s\nwant %s", got, want)
	}
}

// Matches are found as accessors find them.
func TestFind(t *testing.T) {
	d := web()
	for expr, want := range map[string][]string{
		// elements of a list map, by their keys, however they're selected
		"$.spec.template.spec.containers[0].image":                    {".spec.template.spec.containers[name=app].image"},
		"$.spec.template.spec.containers[-1].image":                   {".spec.template.spec.containers[name=proxy].image"},
		"$.spec.template.spec.containers[?(@.name=='proxy')].image":   {".spec.template.spec.containers[name=proxy].image"},
		"{.spec.template.spec.containers[*].image}":                   {".spec.template.spec.containers[name=app].image", ".spec.template.spec.containers[name=proxy].image"},
		"$.spec.template.spec.containers[?(@.image != 'app:1')].name": {".spec.template.spec.containers[name=proxy].name"},
		// recursive descent only finds present values
		"$..image": {".spec.template.spec.containers[name=app].image", ".spec.template.spec.containers[name=proxy].image"},
		// an atomic list, by index; a map, by key
		"$.spec.template.spec.tolerations[1].key": {".spec.template.spec.tolerations[1].key"},
		"$.metadata.labels.app":                   {`.metadata.labels["app"]`},
		"$.metadata['labels']['app']":             {`.metadata.labels["app"]`},
		// a definite path resolves whether or not it's set
		"$.spec.minReadySeconds": {".spec.minReadySeconds"},
		// kubectl's grammar: slices, unions, comparisons and existence
		"$.spec.template.spec.containers[0:1].image":                              {".spec.template.spec.containers[name=app].image"},
		"$.spec.template.spec.tolerations[1:].key":                                {".spec.template.spec.tolerations[1].key"},
		"$.spec.template.spec.containers[0,1].name":                               {".spec.template.spec.containers[name=app].name", ".spec.template.spec.containers[name=proxy].name"},
		"$.metadata['name','namespace']":                                          {".metadata.name", ".metadata.namespace"},
		"$.spec.template.spec.containers[?(@.image)].name":                        {".spec.template.spec.containers[name=app].name", ".spec.template.spec.containers[name=proxy].name"},
		"$.spec.template.spec.containers[?(@.ports)].name":                        {},
		"$.spec.template.spec.containers[?(@.terminationGracePeriodSeconds > 1)]": {},
		// a filter that matches nothing, or an index past the end
		"$.spec.template.spec.containers[?(@.name=='sidecar')]": {},
		"$.spec.template.spec.containers[5]":                    {},
	} {
		if got := paths(t, d, expr); !reflect.DeepEqual(got, want) && !(len(got) == 0 && len(want) == 0) {
			t.Errorf("%s\n got %v\nwant %v", expr, got, want)
		}
	}

	for expr, want := range map[string]string{
		"$.spec.replica":             `has no field "replica"`,
		"$.spec[":                    "unterminated array",
		"$.spec.containers[?(name)]": "a filter compares a field",
		"{.metadata}{.spec}":         "want one expression, not 2",
		"name: {.metadata.name}":     "text outside of an expression",
		"{range .items[*]}{end}":     "unsupported",
	} {
		if _, err := jsonpath.FindE(d, expr); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: error %v, want %q", expr, err, want)
		}
	}
}

// Writes through matches are intent, keyed as accessors key them.
func TestWritesThroughMatches(t *testing.T) {
	d := web().DieWith(func(d *dieappsv1.DeploymentDieMut) {
		images, err := jsonpath.ValuesMutE[string](d, "$.spec.template.spec.containers[?(@.name=='app')].image")
		if err != nil || len(images) != 1 {
			t.Fatalf("images %v, %v", images, err)
		}
		images[0].Set("app:2")

		// a definite path finds a value that isn't set, to set it
		ready, err := jsonpath.ValuesMutE[int32](d, "$.spec.minReadySeconds")
		if err != nil {
			t.Fatal(err)
		}
		ready[0].Set(5)

		containers, err := jsonpath.DiesMutE[*diecorev1.ContainerDieMut](d, "$..containers[*]")
		if err != nil || len(containers) != 2 {
			t.Fatalf("containers %v, %v", containers, err)
		}
		containers[1].Args().Set([]string{"--verbose"})
	})
	got, err := d.DieReleaseApplyPatchE()
	assertJSON(t, got, err, `{
		"apiVersion": "apps/v1", "kind": "Deployment", "metadata": {"name": "web"},
		"spec": {"minReadySeconds": 5, "template": {"spec": {"containers": [
			{"name": "app", "image": "app:2"},
			{"name": "proxy", "args": ["--verbose"]}
		]}}}
	}`)

	// matches of another type are refused
	if _, err := jsonpath.DiesE[*diecorev1.ContainerDie](d, "$.spec"); err == nil || !strings.Contains(err.Error(), "not a v1.Container") {
		t.Errorf("error %v", err)
	}
	if _, err := jsonpath.ValuesE[int64](d, "$.spec.minReadySeconds"); err == nil || !strings.Contains(err.Error(), "is a int32, not a int64") {
		t.Errorf("error %v", err)
	}
}

// Unstructured dies find nodes.
func TestUnstructuredFind(t *testing.T) {
	d := unstructured.BlankFor("acme.dev/v1", "Workload").DieWith(func(d *unstructured.ObjectMut) {
		d.Node().Field("spec").Field("containers").Set([]any{
			map[string]any{"name": "app", "image": "app:1"},
			map[string]any{"name": "proxy", "image": "proxy:1"},
		})
	}).DieCommit()

	find := func(d *unstructured.Object, expr string) []*unstructured.Node {
		t.Helper()
		nodes, err := d.Node().FindE(expr)
		if err != nil {
			t.Fatalf("%s: %v", expr, err)
		}
		return nodes
	}
	findMut := func(d *unstructured.ObjectMut, expr string) []*unstructured.NodeMut {
		t.Helper()
		nodes, err := d.Node().FindE(expr)
		if err != nil {
			t.Fatalf("%s: %v", expr, err)
		}
		return nodes
	}
	// a filter matching one element by one field finds it by that field, as
	// Item does; otherwise, by index
	if got := find(d, "$.spec.containers[?(@.name=='app')].image"); len(got) != 1 || got[0].DiePath().String() != ".spec.containers[name=app].image" || got[0].Get() != "app:1" {
		t.Errorf("app image %v", got)
	}
	if got := find(d, "$.spec.containers[*].image"); len(got) != 2 || got[1].DiePath().String() != ".spec.containers[1].image" {
		t.Errorf("images %v", got)
	}
	if got := find(d, "$..image"); len(got) != 2 {
		t.Errorf("images %v", got)
	}

	d = d.DieWith(func(d *unstructured.ObjectMut) {
		findMut(d, "$.spec.containers[?(@.name=='app')].image")[0].Set("app:2")
		// a definite path finds a node with no value, to set it
		findMut(d, "$.spec.replicas")[0].Set(3)
	})
	got, err := d.DieReleaseMergePatchE()
	assertJSON(t, got, err, `{"spec": {"replicas": 3, "containers": [{"name": "app", "image": "app:2"}, {"name": "proxy", "image": "proxy:1"}]}}`)

	// the items of an unstructured list, as dies
	list := unstructured.BlankFor("acme.dev/v1", "WorkloadList").DieWith(func(l *unstructured.ObjectMut) {
		l.Node().Field("items").Set([]any{map[string]any{"metadata": map[string]any{"name": "a"}}})
	})
	items, err := jsonpath.DiesE[*unstructured.Object](list, "$.items[*]")
	if err != nil || len(items) != 1 || items[0].GetName() != "a" {
		t.Errorf("items %v, %v", items, err)
	}
}

// . is the die queried and $ its root: the paths found are within the root
// either way, as every path is.
func TestFindFromAChildDie(t *testing.T) {
	spec := web().Spec().Template().Spec()
	for expr, want := range map[string][]string{
		".containers[0].image":   {".spec.template.spec.containers[name=app].image"},
		"{.containers[0].image}": {".spec.template.spec.containers[name=app].image"},
		"..image":                {".spec.template.spec.containers[name=app].image", ".spec.template.spec.containers[name=proxy].image"},
		"$.metadata.name":        {".metadata.name"},
		"$..image":               {".spec.template.spec.containers[name=app].image", ".spec.template.spec.containers[name=proxy].image"},
	} {
		if got := paths(t, spec, expr); !reflect.DeepEqual(got, want) {
			t.Errorf("%s\n got %v\nwant %v", expr, got, want)
		}
	}
	if _, err := jsonpath.FindE(spec, ".metadata"); err == nil || !strings.Contains(err.Error(), `has no field "metadata"`) {
		t.Errorf("error %v", err)
	}
	if _, err := jsonpath.FindE(spec, "$.containers"); err == nil || !strings.Contains(err.Error(), `has no field "containers"`) {
		t.Errorf("error %v", err)
	}
}

// validateSubject is written for a subject, wherever it is: a subject's
// namespace defaults to its object's, and is constrained to it unless the
// object is cluster scoped.
func validateSubject(s *dierbacv1.SubjectDieMut) error {
	if s.Kind().Get() != dierbacv1.ServiceAccountKind {
		return nil
	}
	namespaces, err := jsonpath.ValuesE[string](s, "$.metadata.namespace")
	if err != nil {
		return err
	}
	objectNamespace := namespaces[0].Get()
	switch ns := s.Namespace().Get(); {
	case objectNamespace == "":
		// a cluster-scoped object may bind subjects of any namespace
	case ns == "":
		s.Namespace().Set(objectNamespace)
	case ns != objectNamespace:
		return fmt.Errorf("%s: namespace %q must be %q, the namespace of the binding", s.DiePath(), ns, objectNamespace)
	}
	return nil
}

func validateSubjects(d die.MutableDie) error {
	subjects, err := jsonpath.DiesMutE[*dierbacv1.SubjectDieMut](d, "$.subjects[*]")
	if err != nil {
		return err
	}
	var errs []error
	for _, s := range subjects {
		errs = append(errs, validateSubject(s))
	}
	return errors.Join(errs...)
}

func TestValidatingInTheContextOfTheObject(t *testing.T) {
	subjects := []upstreamrbacv1.Subject{
		{Kind: dierbacv1.ServiceAccountKind, Name: "builder"},
		{Kind: dierbacv1.ServiceAccountKind, Name: "deployer", Namespace: "prod"},
	}
	binding := dierbacv1.RoleBindingBlank.DieWith(func(d *dierbacv1.RoleBindingDieMut) {
		d.Metadata().Namespace().Set("prod")
		d.Subjects().Set(subjects)
	}).DieMutable()
	if err := validateSubjects(binding); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	// the default is recorded on the binding
	if got := binding.Subjects().Get()[0].Namespace; got != "prod" {
		t.Errorf("defaulted namespace %q", got)
	}

	binding.Metadata().Namespace().Set("dev")
	if err := validateSubjects(binding); err == nil || !strings.Contains(err.Error(), `.subjects[1]: namespace "prod" must be "dev"`) {
		t.Errorf("error %v", err)
	}

	// a cluster role binding has no namespace, and constrains none
	cluster := dierbacv1.ClusterRoleBindingBlank.DieWith(func(d *dierbacv1.ClusterRoleBindingDieMut) {
		d.Subjects().Set(subjects)
	}).DieMutable()
	if err := validateSubjects(cluster); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got := cluster.Subjects().Get()[0].Namespace; got != "" {
		t.Errorf("namespace %q", got)
	}
}

// A value outside the die queried reads as any does; writing it needs a
// mutable die.
func TestValuesOutsideTheDie(t *testing.T) {
	binding := dierbacv1.RoleBindingBlank.DieWith(func(d *dierbacv1.RoleBindingDieMut) {
		d.Metadata().Namespace().Set("prod")
		d.Subjects().Set([]upstreamrbacv1.Subject{{Kind: dierbacv1.ServiceAccountKind, Name: "builder"}})
	})
	subjects, err := jsonpath.DiesE[*dierbacv1.SubjectDie](binding, "$.subjects[*]")
	if err != nil {
		t.Fatal(err)
	}
	ns, err := jsonpath.ValuesE[string](subjects[0], "$.metadata.namespace")
	if err != nil {
		t.Fatal(err)
	}
	if got := ns[0].Get(); got != "prod" {
		t.Errorf("namespace %q", got)
	}
	// a mutable die writes outside itself, to its root; an immutable die's
	// values only read, as ValuesMut and DiesMut only take mutable dies
	m := binding.DieMutable()
	ms := jsonpath.DiesMut[*dierbacv1.SubjectDieMut](m, "$.subjects[*]")
	jsonpath.ValuesMut[string](ms[0], "$.metadata.namespace")[0].Set("dev")
	if m.GetNamespace() != "dev" || binding.GetNamespace() != "prod" {
		t.Errorf("namespace %q, the immutable die's %q", m.GetNamespace(), binding.GetNamespace())
	}
}
