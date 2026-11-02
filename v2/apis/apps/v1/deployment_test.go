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
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	upstreamappsv1 "k8s.io/api/apps/v1"
	upstreamcorev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"

	appsv1 "reconciler.io/dies/v2/apis/apps/v1"
	corev1 "reconciler.io/dies/v2/apis/core/v1"
	"reconciler.io/dies/v2/die"
	"reconciler.io/dies/v2/die/base"
	"reconciler.io/dies/v2/die/codegen"
)

// upstreamDeployment exercises much of the shape of a deployment, built with
// the upstream types.
func upstreamDeployment() *upstreamappsv1.Deployment {
	return &upstreamappsv1.Deployment{
		TypeMeta: metav1.TypeMeta{APIVersion: "apps/v1", Kind: "Deployment"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-app",
			Namespace: "default",
			Labels:    map[string]string{"app": "my-app"},
			OwnerReferences: []metav1.OwnerReference{
				{APIVersion: "v1", Kind: "ConfigMap", Name: "owner", UID: "1234", Controller: ptr.To(true)},
			},
			Finalizers: []string{"example.dev/cleanup"},
		},
		Spec: upstreamappsv1.DeploymentSpec{
			Replicas: ptr.To[int32](3),
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": "my-app"}},
			Strategy: upstreamappsv1.DeploymentStrategy{
				Type: upstreamappsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &upstreamappsv1.RollingUpdateDeployment{
					MaxSurge: ptr.To(intstr.FromString("25%")),
				},
			},
			Template: upstreamcorev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "my-app"}},
				Spec: upstreamcorev1.PodSpec{
					Containers: []upstreamcorev1.Container{{
						Name:  "app",
						Image: "busybox:1",
						// protocol is left to its default
						Ports: []upstreamcorev1.ContainerPort{{Name: "http", ContainerPort: 8080}},
						Env: []upstreamcorev1.EnvVar{
							{Name: "A", Value: "1"},
							{Name: "PASSWORD", ValueFrom: &upstreamcorev1.EnvVarSource{
								SecretKeyRef: &upstreamcorev1.SecretKeySelector{
									LocalObjectReference: upstreamcorev1.LocalObjectReference{Name: "creds"},
									Key:                  "password",
								},
							}},
						},
						Resources: upstreamcorev1.ResourceRequirements{
							Limits: upstreamcorev1.ResourceList{upstreamcorev1.ResourceCPU: resource.MustParse("500m")},
						},
						VolumeMounts: []upstreamcorev1.VolumeMount{{Name: "creds", MountPath: "/creds"}},
						ReadinessProbe: &upstreamcorev1.Probe{
							ProbeHandler: upstreamcorev1.ProbeHandler{
								HTTPGet: &upstreamcorev1.HTTPGetAction{Path: "/ready", Port: intstr.FromInt32(8080)},
							},
						},
					}},
					Volumes: []upstreamcorev1.Volume{{
						Name:         "creds",
						VolumeSource: upstreamcorev1.VolumeSource{Secret: &upstreamcorev1.SecretVolumeSource{SecretName: "creds"}},
					}},
				},
			},
		},
		Status: upstreamappsv1.DeploymentStatus{
			ObservedGeneration: 2,
			Conditions: []upstreamappsv1.DeploymentCondition{{
				Type:           upstreamappsv1.DeploymentAvailable,
				Status:         upstreamcorev1.ConditionTrue,
				LastUpdateTime: metav1.Unix(1700000000, 0),
			}},
		},
	}
}

func fedDie(t *testing.T) *appsv1.DeploymentDie {
	t.Helper()
	b, err := json.Marshal(upstreamDeployment())
	if err != nil {
		t.Fatal(err)
	}
	// the zero value die is mutable, as a decoder would allocate it
	d := &appsv1.DeploymentDieMut{}
	if err := json.Unmarshal(b, d); err != nil {
		t.Fatal(err)
	}
	return d.DieImmutable()
}

func assertJSON(t *testing.T, got []byte, err error, want string) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var g, w any
	if err := json.Unmarshal(got, &g); err != nil {
		t.Fatalf("invalid json %s: %v", got, err)
	}
	if err := json.Unmarshal([]byte(want), &w); err != nil {
		t.Fatalf("invalid expected json %s: %v", want, err)
	}
	if !reflect.DeepEqual(g, w) {
		t.Errorf("json mismatch\n got: %s\nwant: %s", got, want)
	}
}

// Dies reference the upstream types, so values cross into client-go and other
// libraries without conversion, and the aliases make one import enough.
func TestUpstreamTypesInAndOut(t *testing.T) {
	var in *upstreamappsv1.Deployment = upstreamDeployment()
	d := appsv1.DeploymentBlank.DieWith(func(d *appsv1.DeploymentDieMut) { d.DieFeedPtr(in) })

	var out upstreamappsv1.Deployment = d.DieRelease()
	if !reflect.DeepEqual(*in, out) {
		t.Errorf("released value differs from fed value")
	}

	// the same type and constants, through the die package alone
	var _ appsv1.Deployment = out
	var _ upstreamcorev1.ResourceName = corev1.ResourceCPU
	// upstream methods come with the type
	if got := d.Spec().Template().Spec().Containers().At("app").Resources().DieRelease().Limits; got.Cpu().String() != "500m" {
		t.Errorf("cpu %s", got.Cpu())
	}
}

// A die's JSON decodes as the upstream object it was fed. It omits what the
// upstream encoding adds for absent values, such as a zero metav1.Time's null.
func TestRoundTripsUpstreamJSON(t *testing.T) {
	got, err := json.Marshal(fedDie(t))
	if err != nil {
		t.Fatal(err)
	}
	var decoded upstreamappsv1.Deployment
	if err := json.Unmarshal(got, &decoded); err != nil {
		t.Fatal(err)
	}
	if want := upstreamDeployment(); !reflect.DeepEqual(&decoded, want) {
		t.Errorf("decoded %+v, want %+v", decoded, want)
	}
	if want, _ := json.Marshal(upstreamDeployment()); !strings.Contains(string(want), `"lastTransitionTime":null`) || strings.Contains(string(got), "null") {
		t.Errorf("got %s", got)
	}
}

func TestReadsThroughGeneratedAccessors(t *testing.T) {
	d := fedDie(t)
	c := d.Spec().Template().Spec().Containers().At("app")
	if got := c.Resources().Limits().At(corev1.ResourceCPU); got.String() != "500m" {
		t.Errorf("cpu limit %s", got.String())
	}
	if got := c.Env().At("PASSWORD").ValueFrom().SecretKeyRef().Name().Get(); got != "creds" {
		t.Errorf("secret name %q, through the inlined LocalObjectReference", got)
	}
	if got := c.ReadinessProbe().HTTPGet().Path().Get(); got != "/ready" {
		t.Errorf("probe path %q, through the inlined ProbeHandler", got)
	}
	if got := d.Spec().Template().Spec().Volumes().At("creds").Secret().SecretName().Get(); got != "creds" {
		t.Errorf("secret volume %q, through the inlined VolumeSource", got)
	}
	if got := d.Spec().Strategy().RollingUpdate().MaxSurge().Get(); got.String() != "25%" {
		t.Errorf("max surge %s", got.String())
	}
	// keyed by types.UID, from markers on the reused apimachinery type
	if !d.Metadata().OwnerReferences().Has("1234") {
		t.Errorf("owner reference not found by uid")
	}
	// the legacy condition type is kept, with its extra timestamp
	if got := d.Status().Conditions().At(appsv1.DeploymentAvailable).LastUpdateTime().Get(); got.Unix() != 1700000000 {
		t.Errorf("last update time %v", got)
	}
}

func TestApplyFromServerState(t *testing.T) {
	d := fedDie(t).DieWith(func(d *appsv1.DeploymentDieMut) {
		c := d.Spec().Template().Spec().Containers().At("app")
		c.Image().Set("busybox:2")
		c.Env().At("B").Value().Set("2")
		// the existing port has no protocol; its default identifies it
		c.Ports().At(corev1.ContainerPortKey{ContainerPort: 8080}).Name().Set("web")
		c.Resources().Limits().Set(corev1.ResourceMemory, resource.MustParse("64Mi"))
		d.Metadata().Labels().Set("tier", "web")
	})

	got, err := d.DieMain().DieReleaseApplyPatchE()
	assertJSON(t, got, err, `{
		"apiVersion": "apps/v1",
		"kind": "Deployment",
		"metadata": {"name": "my-app", "namespace": "default", "labels": {"tier": "web"}},
		"spec": {"template": {"spec": {"containers": [{
			"name": "app",
			"image": "busybox:2",
			"env": [{"name": "B", "value": "2"}],
			"ports": [{"containerPort": 8080, "protocol": "TCP", "name": "web"}],
			"resources": {"limits": {"memory": "64Mi"}}
		}]}}}
	}`)
	if ports := d.DieRelease().Spec.Template.Spec.Containers[0].Ports; len(ports) != 1 {
		t.Errorf("expected the defaulted port to be updated in place, got %+v", ports)
	}
}

func TestStatusSubresourceFromMarkers(t *testing.T) {
	d := fedDie(t).DieWith(func(d *appsv1.DeploymentDieMut) {
		d.Spec().Replicas().Set(5)
		d.Status().ObservedGeneration().Set(3)
	})
	got, err := d.DieSubresource("status").DieReleaseApplyPatchE()
	assertJSON(t, got, err, `{
		"apiVersion": "apps/v1",
		"kind": "Deployment",
		"metadata": {"name": "my-app", "namespace": "default"},
		"status": {"observedGeneration": 3}
	}`)
}

func TestDieEnumeratesItsFields(t *testing.T) {
	d := fedDie(t)

	fields := map[string]base.FieldInfo{}
	var wires []string
	for _, f := range d.Spec().DieFields() {
		fields[f.Wire] = f
		wires = append(wires, f.Wire)
	}
	if want := codegen.WireFields[appsv1.DeploymentSpec](); !reflect.DeepEqual(wires, want) {
		t.Errorf("fields %v, want %v", wires, want)
	}
	if f := fields["template"]; f.Kind != base.FieldStruct || f.DiePath().String() != ".spec.template" {
		t.Errorf("template %+v at %s", f, f.DiePath())
	}
	if f := fields["replicas"]; f.Kind != base.FieldValue {
		t.Errorf("replicas %+v", f)
	}

	c := d.Spec().Template().Spec().Containers().At("app")
	for _, f := range c.DieFields() {
		if f.Wire == "ports" && (f.Kind != base.FieldListMap || !reflect.DeepEqual(f.Keys, []string{"containerPort", "protocol"})) {
			t.Errorf("ports %+v", f)
		}
	}
	for _, f := range base.ManagedFields[corev1.Volume]() {
		if f.Wire == "secret" && f.Union != "source" {
			t.Errorf("secret %+v", f)
		}
	}

	// fields are locations, so they can select intent
	d = d.DieWith(func(d *appsv1.DeploymentDieMut) {
		d.Spec().Replicas().Set(5)
		d.Spec().Paused().Set(true)
		d.Spec().Template().Spec().Containers().At("app").Image().Set("busybox:2")
	})
	var values []die.Located
	for _, f := range d.Spec().DieFields() {
		if f.Kind == base.FieldValue {
			values = append(values, f)
		}
	}
	got, err := d.DieOnly(values...).DieReleaseApplyPatchE()
	assertJSON(t, got, err, `{
		"apiVersion": "apps/v1",
		"kind": "Deployment",
		"metadata": {"name": "my-app", "namespace": "default"},
		"spec": {"replicas": 5, "paused": true}
	}`)
}

// Object dies are runtime.Objects and metav1.Objects, and list dies are
// metav1.ListInterfaces, so Kubernetes libraries can use them directly.
func TestObjectInterfaces(t *testing.T) {
	var obj runtime.Object = appsv1.DeploymentBlank
	if gvk := obj.GetObjectKind().GroupVersionKind(); gvk != (schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}) {
		t.Errorf("gvk %v", gvk)
	}

	d := fedDie(t).DieMutable()
	accessor, err := meta.Accessor(d)
	if err != nil {
		t.Fatal(err)
	}
	if accessor.GetName() != "my-app" || accessor.GetOwnerReferences()[0].UID != "1234" {
		t.Errorf("name %q, owners %v", accessor.GetName(), accessor.GetOwnerReferences())
	}
	// a setter changes the die in place, and is recorded as intent
	accessor.SetLabels(map[string]string{"tier": "web"})
	accessor.SetResourceVersion("42")
	if d.Metadata().Labels().At("tier") != "web" || d.GetResourceVersion() != "42" {
		t.Errorf("labels %v, resource version %q", d.GetLabels(), d.GetResourceVersion())
	}
	got, err := d.DieReleaseMergePatchE()
	assertJSON(t, got, err, `{"metadata": {"labels": {"app": null, "tier": "web"}, "resourceVersion": "42"}}`)

	// a copy is a die, independent of the original
	c := d.DeepCopyObject().(*appsv1.DeploymentDieMut)
	c.SetName("copy")
	if d.GetName() != "my-app" || c.GetName() != "copy" {
		t.Errorf("names %q and %q", d.GetName(), c.GetName())
	}

	// setters have nowhere to put a change to an immutable die
	assertPanics(t, func() { appsv1.DeploymentBlank.SetName("x") })
	assertPanics(t, func() { appsv1.DeploymentBlank.GetObjectKind().SetGroupVersionKind(schema.GroupVersionKind{}) })

	list := appsv1.NewDeploymentList(fedDie(t).DieMutable())
	la, err := meta.ListAccessor(list)
	if err != nil {
		t.Fatal(err)
	}
	la.SetContinue("next")
	la.SetRemainingItemCount(ptr.To[int64](3))
	if list.GetContinue() != "next" || *list.GetRemainingItemCount() != 3 {
		t.Errorf("continue %q, remaining %v", list.GetContinue(), list.GetRemainingItemCount())
	}
	if gvk := list.GetObjectKind().GroupVersionKind(); gvk.Kind != "DeploymentList" {
		t.Errorf("list gvk %v", gvk)
	}

	// apimachinery reads and fills a list by reflecting on its items, as for
	// any list
	objs, err := meta.ExtractList(list)
	if err != nil || len(objs) != 1 || objs[0] != runtime.Object(list.Items[0]) {
		t.Fatalf("extracted %v, %v", objs, err)
	}
	other := appsv1.DeploymentBlank.DieWith(func(d *appsv1.DeploymentDieMut) { d.Metadata().Name().Set("other") }).DieMutable()
	if err := meta.SetList(list, append(objs, other)); err != nil {
		t.Fatal(err)
	}
	if len(list.Items) != 2 || list.Items[1].GetName() != "other" {
		t.Errorf("items %v", list.Items)
	}

	// copies are independent, and lists round trip as JSON, items as dies
	copied := list.DeepCopy()
	copied.Items[0].SetName("copy")
	if list.Items[0].GetName() != "my-app" {
		t.Errorf("name %q", list.Items[0].GetName())
	}
	body, err := json.Marshal(list)
	if err != nil {
		t.Fatal(err)
	}
	var decoded appsv1.DeploymentListDie
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Kind != "DeploymentList" || decoded.GetContinue() != "next" || len(decoded.Items) != 2 || decoded.Items[1].GetName() != "other" {
		t.Errorf("decoded %s", body)
	}
}

func assertPanics(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Errorf("expected panic")
		}
	}()
	fn()
}
