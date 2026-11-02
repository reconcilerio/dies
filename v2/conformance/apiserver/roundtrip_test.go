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

package apiserver_test

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/intstr"
	utiljson "k8s.io/apimachinery/pkg/util/json"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	diebatchv1 "reconciler.io/dies/v2/apis/batch/v1"
	diecorev1 "reconciler.io/dies/v2/apis/core/v1"
	dierbacv1 "reconciler.io/dies/v2/apis/rbac/v1"
	widgetv1 "reconciler.io/dies/v2/examples/widget/v1"
)

// Objects created from dies, as the server has them after defaulting, read
// into dies and encoded again, are what the server has: nothing lost,
// nothing added.

type roundTrip struct {
	name string
	// object is created, then read into blank
	object, blank client.Object
}

func roundTrips(ns string) []roundTrip {
	return []roundTrip{
		{"Deployment", deployment(ns).DieImmutable(false), dieappsv1.DeploymentBlank.DieImmutable(false)},
		{"Pod", diecorev1.PodBlank.DieWith(func(d *diecorev1.PodDie) {
			d.Metadata().Namespace().Set(ns)
			d.Metadata().Name().Set("web")
			d.Spec().Containers().At("app").DieWith(func(c *diecorev1.ContainerDie) {
				c.Image().Set("app:1")
				c.Ports().At(diecorev1.ContainerPortKey{ContainerPort: 8080}).Name().Set("http")
				c.Env().At("COLOR").ValueFrom().ConfigMapKeyRef().DieWith(func(r *diecorev1.ConfigMapKeySelectorDie) {
					r.Name().Set("settings")
					r.Key().Set("color")
				})
			})
		}).DieImmutable(false), diecorev1.PodBlank.DieImmutable(false)},
		{"Service", diecorev1.ServiceBlank.DieWith(func(d *diecorev1.ServiceDie) {
			d.Metadata().Namespace().Set(ns)
			d.Metadata().Name().Set("web")
			d.Spec().Selector().Set(map[string]string{"app": "web"})
			d.Spec().Ports().At(diecorev1.ServicePortKey{Port: 80}).TargetPort().Set(intstr.FromInt32(8080))
		}).DieImmutable(false), diecorev1.ServiceBlank.DieImmutable(false)},
		{"ConfigMap", diecorev1.ConfigMapBlank.DieWith(func(d *diecorev1.ConfigMapDie) {
			d.Metadata().Namespace().Set(ns)
			d.Metadata().Name().Set("web")
			d.Data().Set("color", "")
			d.BinaryData().Set("blob", []byte{0, 1, 2})
		}).DieImmutable(false), diecorev1.ConfigMapBlank.DieImmutable(false)},
		{"CronJob", diebatchv1.CronJobBlank.DieWith(func(d *diebatchv1.CronJobDie) {
			d.Metadata().Namespace().Set(ns)
			d.Metadata().Name().Set("web")
			d.Spec().Schedule().Set("@hourly")
			d.Spec().JobTemplate().Spec().Template().Spec().DieWith(func(p *diecorev1.PodSpecDie) {
				p.RestartPolicy().Set(diecorev1.RestartPolicyNever)
				p.Containers().At("job").Image().Set("job:1")
			})
		}).DieImmutable(false), diebatchv1.CronJobBlank.DieImmutable(false)},
		{"Role", dierbacv1.RoleBlank.DieWith(func(d *dierbacv1.RoleDie) {
			d.Metadata().Namespace().Set(ns)
			d.Metadata().Name().Set("web")
			d.Rules().Set([]dierbacv1.PolicyRule{{APIGroups: []string{""}, Resources: []string{"configmaps"}, Verbs: []string{"get"}}})
		}).DieImmutable(false), dierbacv1.RoleBlank.DieImmutable(false)},
		{"Widget", widgetv1.WidgetBlank.DieWith(func(d *widgetv1.WidgetDie) {
			d.Metadata().Namespace().Set(ns)
			d.Metadata().Name().Set("web")
			d.Spec().Ports().At("http").Port().Set(80)
			d.Spec().Template().Spec().Containers().At("app").Image().Set("app:1")
			d.Spec().Match().Exact().Set(false)
		}).DieImmutable(false), widgetv1.WidgetBlank.DieImmutable(false)},
	}
}

func TestServerObjectsRoundTripThroughDies(t *testing.T) {
	ns := namespace(t)
	ctx := context.Background()
	for _, rt := range roundTrips(ns) {
		t.Run(rt.name, func(t *testing.T) {
			if err := dies.Create(ctx, rt.object); err != nil {
				t.Fatal(err)
			}
			if err := dies.Get(ctx, client.ObjectKeyFromObject(rt.object), rt.blank); err != nil {
				t.Fatal(err)
			}
			gvk, err := dies.GroupVersionKindFor(rt.blank)
			if err != nil {
				t.Fatal(err)
			}
			server := &unstructured.Unstructured{}
			server.SetGroupVersionKind(gvk)
			if err := upstream.Get(ctx, client.ObjectKeyFromObject(rt.object), server); err != nil {
				t.Fatal(err)
			}

			data, err := json.Marshal(rt.blank)
			if err != nil {
				t.Fatal(err)
			}
			var got map[string]any
			if err := utiljson.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			want := server.Object
			// the client strips the type metadata of a typed object as it
			// decodes it, which the die has no say in
			delete(want, "apiVersion")
			delete(want, "kind")
			delete(got, "apiVersion")
			delete(got, "kind")
			if !reflect.DeepEqual(got, want) {
				g, _ := json.MarshalIndent(got, "", "  ")
				w, _ := json.MarshalIndent(want, "", "  ")
				t.Errorf("the die encodes:\n%s\nthe server has:\n%s", g, w)
			}
		})
	}
}

// Updating from a die sends the object as it was read, with the change: the
// server's defaults and allocations, such as a Service's cluster IP, are
// kept.
func TestUpdateFromADieKeepsWhatTheServerSet(t *testing.T) {
	ns := namespace(t)
	ctx := context.Background()
	svc := diecorev1.ServiceBlank.DieWith(func(d *diecorev1.ServiceDie) {
		d.Metadata().Namespace().Set(ns)
		d.Metadata().Name().Set("web")
		d.Spec().Selector().Set(map[string]string{"app": "web"})
		d.Spec().Ports().At(diecorev1.ServicePortKey{Port: 80}).TargetPort().Set(intstr.FromInt32(8080))
	}).DieImmutable(false)
	if err := dies.Create(ctx, svc); err != nil {
		t.Fatal(err)
	}
	before := svc.Spec().DieRelease()
	if before.ClusterIP == "" {
		t.Fatal("the server allocated no cluster IP")
	}

	read := diecorev1.ServiceBlank.DieImmutable(false)
	if err := dies.Get(ctx, keyOf(ns, "web"), read); err != nil {
		t.Fatal(err)
	}
	read.Metadata().Labels().Set("tier", "web")
	if err := dies.Update(ctx, read); err != nil {
		t.Fatal(err)
	}
	if after := read.Spec().DieRelease(); !reflect.DeepEqual(after, before) {
		t.Errorf("spec changed:\n%+v\nwas:\n%+v", after, before)
	}
	if read.GetLabels()["tier"] != "web" {
		t.Errorf("labels %v", read.GetLabels())
	}
}

// A list die is a list of dies, as clients list into.
func TestListIntoListDies(t *testing.T) {
	ns := namespace(t)
	ctx := context.Background()
	for _, name := range []string{"a", "b"} {
		cm := diecorev1.ConfigMapBlank.DieWith(func(d *diecorev1.ConfigMapDie) {
			d.Metadata().Namespace().Set(ns)
			d.Metadata().Name().Set(name)
			d.Data().Set("name", name)
		}).DieImmutable(false)
		if err := dies.Create(ctx, cm); err != nil {
			t.Fatal(err)
		}
	}
	list := diecorev1.NewConfigMapList()
	if err := dies.List(ctx, list, client.InNamespace(ns)); err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, item := range list.Items {
		// there's no controller manager to publish the namespace's root CA,
		// so these are the only config maps
		names = append(names, item.Data().At("name"))
	}
	if !reflect.DeepEqual(names, []string{"a", "b"}) {
		t.Errorf("items %v", names)
	}
	if list.GetResourceVersion() == "" {
		t.Errorf("the list has no resource version")
	}
}
