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

package appsv1

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"k8s.io/utils/ptr"

	"reconciler.io/dies/v2/die"
	"reconciler.io/dies/v2/die/dietest"
	"reconciler.io/dies/v2/spike/corev1"
	"reconciler.io/dies/v2/spike/metav1"
)

// serverDeployment is the state of a deployment as read from the API server.
func serverDeployment() Deployment {
	return Deployment{
		TypeMeta: metav1.TypeMeta{APIVersion: "apps/v1", Kind: "Deployment"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-app",
			Namespace: "default",
			Labels:    map[string]string{"team": "a"},
		},
		Spec: DeploymentSpec{
			Replicas: ptr.To[int32](3),
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{Name: "app", Image: "busybox:1", Env: []corev1.EnvVar{{Name: "A", Value: "1"}}},
						{Name: "sidecar", Image: "proxy"},
					},
					Volumes: []corev1.Volume{
						{Name: "creds", VolumeSource: corev1.VolumeSource{Secret: &corev1.SecretVolumeSource{SecretName: "s"}}},
					},
				},
			},
		},
	}
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

func assertIntent(t *testing.T, got string, want ...string) {
	t.Helper()
	if w := strings.Join(want, "\n"); got != w {
		t.Errorf("intent mismatch\n got:\n%s\nwant:\n%s", got, w)
	}
}

func TestGetChainCreatesOnWrite(t *testing.T) {
	d := DeploymentBlank.DieStamp(func(d *DeploymentDieMut) {
		c := d.Spec().Template().Spec().Containers().At("app")
		// reads have no side effects, and a missing element reads with its key
		if c.Name().Get() != "app" || c.Image().Get() != "" {
			t.Errorf("unexpected read: name=%q image=%q", c.Name().Get(), c.Image().Get())
		}
		if d.DieIntent() != "" {
			t.Errorf("reads recorded intent: %s", d.DieIntent())
		}
		dietest.NoError(t, c.Image().Set("busybox"))
	})

	want := []corev1.Container{{Name: "app", Image: "busybox"}}
	if got := d.DieRelease().Spec.Template.Spec.Containers; !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
	assertIntent(t, d.DieIntent(), `set .spec.template.spec.containers[name=app].image = "busybox"`)
	if got := *d.Spec().Template().Spec().DieRelease(); !reflect.DeepEqual(got.Containers, want) {
		t.Errorf("child release got %+v", got)
	}
}

func TestImmutability(t *testing.T) {
	a := DeploymentBlank.DieStamp(func(d *DeploymentDieMut) { dietest.NoError(t, d.Metadata().Name().Set("a")) })
	b := a.DieStamp(func(d *DeploymentDieMut) { dietest.NoError(t, d.Metadata().Name().Set("b")) })
	if DeploymentBlank.Metadata().Name().Get() != "" || a.Metadata().Name().Get() != "a" || b.Metadata().Name().Get() != "b" {
		t.Errorf("names: blank=%q a=%q b=%q", DeploymentBlank.Metadata().Name().Get(), a.Metadata().Name().Get(), b.Metadata().Name().Get())
	}

	// an immutable die has no writes: its mutable copy writes, without
	// reaching it
	c := corev1.ContainerBlank.DieMutable()
	dietest.NoError(t, c.Name().Set("app"))
	dietest.NoError(t, c.Image().Set("busybox"))
	if corev1.ContainerBlank.Name().Get() != "" || c.Name().Get() != "app" || c.Image().Get() != "busybox" {
		t.Errorf("unexpected container %+v", *c.DieRelease())
	}
	// nor does an immutable copy of it see its later writes
	frozen := c.DieImmutable()
	dietest.NoError(t, c.Image().Set("nginx"))
	if frozen.Image().Get() != "busybox" {
		t.Errorf("an immutable copy saw a later write: %q", frozen.Image().Get())
	}

	// the zero value of a mutable die is a mutable blank, ready for a decoder
	var z DeploymentDieMut
	dietest.NoError(t, z.Metadata().Name().Set("z"))
	if z.Metadata().Name().Get() != "z" {
		t.Errorf("zero value die is not mutable")
	}
}

func TestApplyPatchFromBlank(t *testing.T) {
	d := DeploymentBlank.DieStamp(func(d *DeploymentDieMut) {
		dietest.NoError(t, d.Metadata().Name().Set("my-app"))
		dietest.NoError(t, d.Metadata().Namespace().Set("default"))
		// explicit zero values, including an omitempty bool
		dietest.NoError(t, d.Spec().Replicas().Set(0))
		dietest.NoError(t, d.Spec().Paused().Set(false))
		ps := d.Spec().Template().Spec()
		dietest.NoError(t, ps.NodeSelector().Clear())
		dietest.NoError(t, ps.InitContainers().Empty())
		dietest.NoError(t, ps.Containers().At("app").Image().Set("busybox"))
	})

	// containers is not omitempty, yet untouched fields never appear as null
	got, err := die.ReleaseApplyPatchE(d)
	assertJSON(t, got, err, `{
		"apiVersion": "apps/v1",
		"kind": "Deployment",
		"metadata": {"name": "my-app", "namespace": "default"},
		"spec": {
			"replicas": 0,
			"paused": false,
			"template": {"spec": {
				"nodeSelector": null,
				"initContainers": [],
				"containers": [{"name": "app", "image": "busybox"}]
			}}
		}
	}`)
}

func TestPatchesFromServerState(t *testing.T) {
	d := DeploymentBlank.DieStamp(func(d *DeploymentDieMut) { dietest.NoError(t, d.DieFeed(new(serverDeployment()))) }).DieStamp(func(d *DeploymentDieMut) {
		dietest.NoError(t, d.Spec().Template().Spec().Containers().At("app").Env().At("FOO").Value().Set("bar"))
		dietest.NoError(t, d.Metadata().Labels().Remove("team"))
	})

	// apply only states what was stamped
	got, err := die.ReleaseApplyPatchE(d)
	// an apply can't remove a label: null would set it to "", so it's left
	// out, and the server removes it if this manager owned it
	assertJSON(t, got, err, `{
		"apiVersion": "apps/v1",
		"kind": "Deployment",
		"metadata": {"name": "my-app", "namespace": "default"},
		"spec": {"template": {"spec": {"containers": [
			{"name": "app", "env": [{"name": "FOO", "value": "bar"}]}
		]}}}
	}`)

	// merge patches can only replace lists, so the whole list is sent
	got, err = die.ReleaseMergePatchE(d)
	assertJSON(t, got, err, `{
		"metadata": {"labels": {"team": null}},
		"spec": {"template": {"spec": {"containers": [
			{"name": "app", "image": "busybox:1", "env": [{"name": "A", "value": "1"}, {"name": "FOO", "value": "bar"}]},
			{"name": "sidecar", "image": "proxy"}
		]}}}
	}`)
}

func TestUnionOverInlinedStruct(t *testing.T) {
	d := DeploymentBlank.DieStamp(func(d *DeploymentDieMut) { dietest.NoError(t, d.DieFeed(new(serverDeployment()))) }).DieStamp(func(d *DeploymentDieMut) {
		v := d.Spec().Template().Spec().Volumes().At("creds")
		if !v.Secret().DieHas() {
			t.Errorf("expected secret volume")
		}
		dietest.NoError(t, v.EmptyDir().Medium().Set("Memory"))
		if v.Secret().DieHas() || !v.EmptyDir().DieHas() {
			t.Errorf("setting emptyDir did not clear secret")
		}
	})

	// inlined fields are addressed by their wire path, without .volumeSource
	assertIntent(t, d.DieIntent(),
		`clear .spec.template.spec.volumes[name=creds].secret`,
		`set .spec.template.spec.volumes[name=creds].emptyDir.medium = "Memory"`,
	)
	want := []corev1.Volume{{Name: "creds", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{Medium: "Memory"}}}}
	if got := d.DieRelease().Spec.Template.Spec.Volumes; !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
	got, err := die.ReleaseApplyPatchE(d)
	assertJSON(t, got, err, `{
		"apiVersion": "apps/v1",
		"kind": "Deployment",
		"metadata": {"name": "my-app", "namespace": "default"},
		"spec": {"template": {"spec": {"volumes": [
			{"name": "creds", "secret": null, "emptyDir": {"medium": "Memory"}}
		]}}}
	}`)
}

func TestMultiKeyListMap(t *testing.T) {
	c := corev1.ContainerBlank.DieStamp(func(c *corev1.ContainerDieMut) {
		dietest.NoError(t, c.Ports().At(corev1.PortKey{ContainerPort: 80, Protocol: "TCP"}).Name().Set("http"))
	})
	want := []corev1.ContainerPort{{Name: "http", ContainerPort: 80, Protocol: "TCP"}}
	if got := c.DieRelease().Ports; !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
	assertIntent(t, c.DieIntent(), `set .ports[containerPort=80,protocol=TCP].name = "http"`)
	if keys := c.Ports().Keys(); !reflect.DeepEqual(keys, []corev1.PortKey{{ContainerPort: 80, Protocol: "TCP"}}) {
		t.Errorf("keys %+v", keys)
	}
}

func TestDefaultedListMapKey(t *testing.T) {
	// as created by a client, before the server defaults the protocol
	c := corev1.ContainerBlank.DieStamp(func(d *corev1.ContainerDieMut) {
		dietest.NoError(t, d.DieFeed(&corev1.Container{
			Name:  "app",
			Ports: []corev1.ContainerPort{{ContainerPort: 80}},
		}))
	})

	// an unset protocol and the default protocol identify the same element
	for _, k := range []corev1.PortKey{{ContainerPort: 80}, {ContainerPort: 80, Protocol: "TCP"}} {
		if !c.Ports().Has(k) {
			t.Errorf("port %+v not found", k)
		}
	}
	if c.Ports().Has(corev1.PortKey{ContainerPort: 80, Protocol: "UDP"}) {
		t.Errorf("UDP port should not match")
	}
	if keys := c.Ports().Keys(); !reflect.DeepEqual(keys, []corev1.PortKey{{ContainerPort: 80, Protocol: "TCP"}}) {
		t.Errorf("keys %+v", keys)
	}

	c = c.DieStamp(func(c *corev1.ContainerDieMut) {
		// updates the existing element rather than adding a second one
		dietest.NoError(t, c.Ports().At(corev1.PortKey{ContainerPort: 80}).Name().Set("http"))
		// a new element gets the default written explicitly
		dietest.NoError(t, c.Ports().At(corev1.PortKey{ContainerPort: 8080}).Name().Set("alt"))
	})
	want := []corev1.ContainerPort{
		{Name: "http", ContainerPort: 80},
		{Name: "alt", ContainerPort: 8080, Protocol: "TCP"},
	}
	if got := c.DieRelease().Ports; !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
	assertIntent(t, c.DieIntent(),
		`set .ports[containerPort=80,protocol=TCP].name = "http"`,
		`set .ports[containerPort=8080,protocol=TCP].name = "alt"`,
	)

	// keys are always sent in full, so the server never has to default them
	// to identify the element
	got, err := die.ReleaseApplyPatchE(c)
	assertJSON(t, got, err, `{"ports": [
		{"containerPort": 80, "protocol": "TCP", "name": "http"},
		{"containerPort": 8080, "protocol": "TCP", "name": "alt"}
	]}`)

	// diffing a stamp matches elements by their defaulted keys
	s := c.DieTare().DieMutable()
	dietest.NoError(t, die.StampValue(s, func(r *corev1.Container) {
		r.Ports[0].Protocol = "TCP"
		r.Ports[1].Protocol = ""
	}))
	assertIntent(t, s.DieIntent(),
		`set .ports[containerPort=80,protocol=TCP].protocol = "TCP"`,
		`set .ports[containerPort=8080,protocol=TCP].protocol = ""`,
	)
}

func TestStampRecordsIntent(t *testing.T) {
	c := corev1.ContainerBlank.DieStamp(func(d *corev1.ContainerDieMut) {
		dietest.NoError(t, d.DieFeed(&corev1.Container{
			Name:  "app",
			Image: "busybox:1",
			Env:   []corev1.EnvVar{{Name: "A", Value: "1"}, {Name: "B", Value: "2"}},
		}))
	}).DieMutable()
	dietest.NoError(t, die.StampValue(c, func(r *corev1.Container) {
		r.Image = "busybox:2"
		r.Env = r.Env[:1]
	}))
	assertIntent(t, c.DieIntent(),
		`set .image = "busybox:2"`,
		`clear .env[name=B]`,
	)
}

func TestUnmarshalRebasesIntent(t *testing.T) {
	d := DeploymentBlank.DieMutable()
	dietest.NoError(t, d.DieFeed(new(serverDeployment())), d.Spec().Template().Spec().Containers().At("app").Image().Set("busybox:2"))

	// someone else scaled the deployment; decoding the latest state into the
	// die keeps our intent on top of it
	latest := serverDeployment()
	latest.Spec.Replicas = ptr.To[int32](5)
	b, _ := json.Marshal(latest)
	if err := json.Unmarshal(b, d); err != nil {
		t.Fatal(err)
	}
	if got := d.Spec().Replicas().Get(); got != 5 {
		t.Errorf("replicas %d, want 5", got)
	}
	if got := d.Spec().Template().Spec().Containers().At("app").Image().Get(); got != "busybox:2" {
		t.Errorf("image %q, want busybox:2", got)
	}
	assertIntent(t, d.DieIntent(), `set .spec.template.spec.containers[name=app].image = "busybox:2"`)

	// after a successful write the intent is tared
	dietest.NoError(t, d.DieTare())
	assertIntent(t, d.DieIntent())
	if got := d.Spec().Template().Spec().Containers().At("app").Image().Get(); got != "busybox:2" {
		t.Errorf("image %q, want busybox:2", got)
	}

	// an immutable die cannot be decoded into
	if err := json.Unmarshal(b, DeploymentBlank); err == nil {
		t.Errorf("expected error unmarshaling into an immutable die")
	}
}

func TestSubresourceSplit(t *testing.T) {
	d := DeploymentBlank.DieStamp(func(d *DeploymentDieMut) {
		dietest.NoError(t, d.Metadata().Name().Set("my-app"))
		dietest.NoError(t, d.Spec().Replicas().Set(2))
		dietest.NoError(t, d.Status().ObservedGeneration().Set(4))
		dietest.NoError(t, d.Status().Conditions().At("Ready").True("Available"))
	})

	got, err := die.ReleaseApplyPatchE(d, die.MainResource())
	assertJSON(t, got, err, `{
		"apiVersion": "apps/v1",
		"kind": "Deployment",
		"metadata": {"name": "my-app"},
		"spec": {"replicas": 2}
	}`)
	got, err = die.ReleaseApplyPatchE(d, die.Subresource("status"))
	assertJSON(t, got, err, `{
		"apiVersion": "apps/v1",
		"kind": "Deployment",
		"metadata": {"name": "my-app"},
		"status": {
			"observedGeneration": 4,
			"conditions": [{"type": "Ready", "status": "True", "reason": "Available"}]
		}
	}`)

	// handles and child dies are locations, so no path needs spelling out
	got, err = die.ReleaseApplyPatchE(d, die.Only(d.Spec().Replicas(), d.Status().Conditions()))
	assertJSON(t, got, err, `{
		"apiVersion": "apps/v1",
		"kind": "Deployment",
		"metadata": {"name": "my-app"},
		"spec": {"replicas": 2},
		"status": {
			"conditions": [{"type": "Ready", "status": "True", "reason": "Available"}]
		}
	}`)
}

func TestPresenceOfOptionalScalars(t *testing.T) {
	d := DeploymentBlank.DieStamp(func(d *DeploymentDieMut) {
		ps := d.Spec().Template().Spec()
		if ps.AutomountServiceAccountToken().Has() {
			t.Errorf("unset field is present")
		}
		dietest.NoError(t, ps.AutomountServiceAccountToken().Set(false))
		if !ps.AutomountServiceAccountToken().Has() || ps.AutomountServiceAccountToken().Get() {
			t.Errorf("explicit false is not present")
		}
	})
	got, err := die.ReleaseApplyPatchE(d)
	assertJSON(t, got, err, `{
		"apiVersion": "apps/v1",
		"kind": "Deployment",
		"spec": {"template": {"spec": {"automountServiceAccountToken": false}}}
	}`)

	d = d.DieStamp(func(d *DeploymentDieMut) {
		dietest.NoError(t, d.Spec().Template().Spec().AutomountServiceAccountToken().Clear())
	})
	if d.Spec().Template().Spec().AutomountServiceAccountToken().Has() {
		t.Errorf("cleared field is present")
	}
	got, err = die.ReleaseApplyPatchE(d)
	assertJSON(t, got, err, `{
		"apiVersion": "apps/v1",
		"kind": "Deployment",
		"spec": {"template": {"spec": {"automountServiceAccountToken": null}}}
	}`)
}

// Equal compares intent as well as values, so that dies that are equal have
// no diff.
func TestEqualComparesIntent(t *testing.T) {
	fed := DeploymentBlank.DieStamp(func(d *DeploymentDieMut) {
		dietest.NoError(t, d.DieFeed(&Deployment{
			TypeMeta:   metav1.TypeMeta{APIVersion: "apps/v1", Kind: "Deployment"},
			ObjectMeta: metav1.ObjectMeta{Name: "x"},
		}))
	})
	stamped := DeploymentBlank.DieStamp(func(d *DeploymentDieMut) { dietest.NoError(t, d.Metadata().Name().Set("x")) })
	if fed.Equal(stamped) {
		t.Errorf("equal values with different intent should not be equal")
	}
	if !reflect.DeepEqual(*fed.DieRelease(), *stamped.DieRelease()) {
		t.Errorf("expected equal values")
	}
	if fed.DieIntent() == stamped.DieIntent() {
		t.Errorf("expected different intent")
	}
}

func TestDeepCopyIsIndependent(t *testing.T) {
	m := DeploymentBlank.DieMutable()
	c := m.DeepCopy()
	dietest.NoError(t, c.Metadata().Name().Set("c"))
	dietest.NoError(t, c.Metadata().Labels().Set("k", "v"))
	if m.Metadata().Name().Get() != "" || m.Metadata().Labels().Len() != 0 {
		t.Errorf("copy shares state with original")
	}
	if c.Metadata().Labels().At("k") != "v" {
		t.Errorf("copy lost its write")
	}
}
