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

package die_test

import (
	"context"
	"testing"

	"k8s.io/apimachinery/pkg/util/validation/field"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	diecorev1 "reconciler.io/dies/v2/apis/core/v1"
	diemetav1 "reconciler.io/dies/v2/apis/meta/v1"
	"reconciler.io/dies/v2/die"
	dietesting "reconciler.io/dies/v2/die/testing"
	"reconciler.io/dies/v2/die/unstructured"
	"reconciler.io/dies/v2/die/validation"
)

func deployment(t *testing.T, yaml string) *dieappsv1.DeploymentDie {
	t.Helper()
	m := dieappsv1.DeploymentBlank.DieMutable()
	if err := die.FeedYAML(m, []byte(yaml)); err != nil {
		t.Fatal(err)
	}
	return m.DieImmutable()
}

func errorStrings(errs field.ErrorList) []string {
	var out []string
	for _, err := range errs {
		out = append(out, err.Error())
	}
	return out
}

func assertErrors(t *testing.T, errs field.ErrorList, want ...string) {
	t.Helper()
	got := errorStrings(errs)
	if len(got) != len(want) {
		t.Fatalf("got %q, want %q", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("error %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestValidateSchema(t *testing.T) {
	d := deployment(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web
spec:
  template:
    spec:
      containers:
      - name: app
        image: app:1
        ports:
        - name: http
      - name: ""
        image: sidecar:1
      volumes:
      - name: config
        configMap:
          name: config
        secret:
          secretName: config
`)
	assertErrors(t, validated(t)(die.ValidateCreate(context.Background(), d)),
		// +required
		"spec.selector: Required value",
		// a union, in field order, with declarative validation's error
		"spec.template.spec.volumes[0]: Invalid value: \"{secret, configMap}\": must specify at most one of: `hostPath`, `emptyDir`, `gcePersistentDisk`, `awsElasticBlockStore`, `gitRepo`, `secret`, `nfs`, `iscsi`, `glusterfs`, `persistentVolumeClaim`, `rbd`, `flexVolume`, `cinder`, `cephfs`, `flocker`, `downwardAPI`, `fc`, `azureFile`, `configMap`, `vsphereVolume`, `quobyte`, `azureDisk`, `photonPersistentDisk`, `projected`, `portworxVolume`, `scaleIO`, `storageos`, `csi`, `ephemeral`, `image`",
		// a list map's key, not defaulted; protocol is
		"spec.template.spec.containers[0].ports[0].containerPort: Required value",
		"spec.template.spec.containers[1].name: Required value",
	)
}

func TestValidateWithin(t *testing.T) {
	d := deployment(t, `
spec:
  selector: {}
  template:
    spec:
      containers:
      - image: app:1
`)
	// a die within another reports the paths from its root, and validates
	// only what it holds
	assertErrors(t, validated(t)(die.ValidateCreate(context.Background(), d.Spec().Template().Spec())),
		"spec.template.spec.containers[0].name: Required value",
	)
	if errs := validated(t)(die.ValidateCreate(context.Background(), d.Metadata())); len(errs) > 0 {
		t.Errorf("metadata: %v", errs)
	}
}

func TestValidationPath(t *testing.T) {
	d := deployment(t, `
spec:
  template:
    metadata:
      labels:
        app: web
    spec:
      containers:
      - name: init
      - name: app
`)
	containers := d.Spec().Template().Spec().Containers()
	tests := []struct {
		l    die.Located
		want string
	}{
		{d.Spec().Replicas(), "spec.replicas"},
		{containers.At("app"), "spec.template.spec.containers[1]"},
		{containers.At("app").Image(), "spec.template.spec.containers[1].image"},
		// an element the list doesn't have is its keys
		{containers.At("missing").Image(), "spec.template.spec.containers[name=missing].image"},
		// as is one of a bare path, which has no list to look in
		{containers.At("app").Image().DiePath(), "spec.template.spec.containers[name=app].image"},
		{d.Spec().Template().Metadata().Labels(), "spec.template.metadata.labels"},
	}
	for _, tt := range tests {
		if got := validation.Path(tt.l).String(); got != tt.want {
			t.Errorf("%s: got %s, want %s", tt.l.DiePath(), got, tt.want)
		}
	}
}

func TestValidationErrors(t *testing.T) {
	d := deployment(t, `
spec:
  replicas: -1
  template:
    spec:
      restartPolicy: Sometimes
      containers:
      - name: app
      - name: sidecar
`)
	spec := d.Spec().Template().Spec()
	assertErrors(t, field.ErrorList{
		validation.Invalid(d.Spec().Replicas(), "must be non-negative"),
		validation.NotSupported(spec.RestartPolicy(), []string{"Always", "OnFailure", "Never"}),
		validation.TooMany(spec.Containers(), 1),
		validation.Required(spec.Containers().At("sidecar").Image(), "an image is required"),
	},
		"spec.replicas: Invalid value: -1: must be non-negative",
		`spec.template.spec.restartPolicy: Unsupported value: "Sometimes": supported values: "Always", "OnFailure", "Never"`,
		"spec.template.spec.containers: Too many: 2: must have at most 1 item",
		"spec.template.spec.containers[1].image: Required value: an image is required",
	)
}

func TestValidateUnstructured(t *testing.T) {
	// an unstructured die has no schema of its own to check, and no dies
	// within it to walk
	ctx := context.Background()
	if errs := validated(t)(die.ValidateCreate(ctx, unstructured.BlankFor("example.dev/v1", "Widget"))); len(errs) > 0 {
		t.Errorf("got %v", errs)
	}
	if err := die.Default(ctx, unstructured.BlankFor("example.dev/v1", "Widget").DieMutable()); err != nil {
		t.Errorf("got %v", err)
	}
}

func TestValidateMarkers(t *testing.T) {
	d := deployment(t, `
metadata:
  finalizers: [a, b, a]
spec:
  selector: {}
  template:
    spec:
      restartPolicy: Sometimes
      containers:
      - name: app
        ports:
        - containerPort: 80
          protocol: HTTP
        restartPolicyRules:
        - action: Restart
`)
	assertErrors(t, validated(t)(die.ValidateCreate(context.Background(), d)),
		// a set: +listType=set
		`metadata.finalizers[2]: Duplicate value: "a"`,
		// +enum: the constants of the field's type
		`spec.template.spec.restartPolicy: Unsupported value: "Sometimes": supported values: "Always", "Never", "OnFailure"`,
		`spec.template.spec.containers[0].ports[0].protocol: Unsupported value: "HTTP": supported values: "SCTP", "TCP", "UDP"`,
		// +oneOf, of an item of an atomic list: exactly one member
		"spec.template.spec.containers[0].restartPolicyRules[0]: Invalid value: \"\": must specify one of: `exitCodes`",
	)

	// +enum of a list's items
	m := dieappsv1.StatefulSetBlank.DieMutable()
	dietesting.NoError(t, die.FeedYAML(m, []byte(`
spec:
  selector: {}
  template: {}
  volumeClaimTemplates:
  - spec:
      accessModes: [ReadWriteOnce, ReadWriteSometimes]
`)))
	errs := validated(t)(die.ValidateCreate(context.Background(), m.DieImmutable()))
	assertErrors(t, errs,
		`spec.volumeClaimTemplates[0].spec.accessModes[1]: Unsupported value: "ReadWriteSometimes": supported values: "ReadOnlyMany", "ReadWriteMany", "ReadWriteOnce", "ReadWriteOncePod"`,
	)
}

func TestValidateKubebuilderMarkers(t *testing.T) {
	// metav1.Condition is validated by kubebuilder's markers, for custom
	// resources that use it
	c := diemetav1.ConditionBlank.DieStamp(func(d *diemetav1.ConditionDieMut) {
		dietesting.NoError(t, d.Type().Set("Ready!"))
		dietesting.NoError(t, d.Status().Set("Maybe"))
		dietesting.NoError(t, d.ObservedGeneration().Set(-1))
	})
	assertErrors(t, validated(t)(die.ValidateCreate(context.Background(), c)),
		`type: Invalid value: "Ready!": must match the regular expression "^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])$"`,
		`status: Unsupported value: "Maybe": supported values: "False", "True", "Unknown"`,
		"observedGeneration: Invalid value: -1: must be greater than or equal to 0",
		"reason: Too short: must be at least 1 character",
	)
}

func TestDefaultMarkers(t *testing.T) {
	m := deployment(t, `
spec:
  template:
    spec:
      containers:
      - name: app
        ports:
        - containerPort: 80
      volumes:
      - name: disk
        azureDisk:
          diskName: disk
          diskURI: uri
`).DieMutable()
	dietesting.NoError(t, die.Default(context.Background(), m))
	spec := m.Spec().Template().Spec()
	// +default, set where the field is absent: a list map's defaulted key is
	// the same element, with its key set
	if got := spec.Containers().At("app").Ports().At(diecorev1.ContainerPortKey{ContainerPort: 80}).Protocol().Get(); got != "TCP" {
		t.Errorf("protocol %q", got)
	}
	disk := spec.Volumes().At("disk").AzureDisk()
	if disk.FSType().Get() != "ext4" || !disk.ReadOnly().Has() || disk.ReadOnly().Get() {
		t.Errorf("azure disk %s", die.MustReleaseJSON(disk))
	}
}

func TestValidateUpdateRatchets(t *testing.T) {
	ctx := context.Background()
	old := deployment(t, `
spec:
  selector: {}
  template:
    spec:
      restartPolicy: Sometimes
      containers:
      - name: init
      - name: app
        ports:
        - containerPort: 80
          protocol: HTTP
      volumes:
      - name: config
        configMap:
          name: config
        secret:
          secretName: config
`)
	if errs := validated(t)(die.ValidateCreate(ctx, old)); len(errs) != 3 {
		t.Errorf("create: %v", errs)
	}
	// values invalid when they were written stay valid while unchanged, a
	// union's members included, though their elements move
	unchanged := old.DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietesting.NoError(t, d.Spec().Replicas().Set(3))
		dietesting.NoError(t, d.Spec().Template().Spec().Containers().Remove("init"))
	})
	if errs := validated(t)(die.ValidateUpdate(ctx, unchanged, old)); len(errs) > 0 {
		t.Errorf("unchanged: %v", errs)
	}
	changed := old.DieStamp(func(d *dieappsv1.DeploymentDieMut) {
		dietesting.NoError(t, d.Spec().Template().Spec().RestartPolicy().Set("Seldom"))
	})
	assertErrors(t, validated(t)(die.ValidateUpdate(ctx, changed, old)),
		`spec.template.spec.restartPolicy: Unsupported value: "Seldom": supported values: "Always", "Never", "OnFailure"`,
	)
}

// validated fails the test if validation fails, and returns its errors.
func validated(t *testing.T) func(field.ErrorList, error) field.ErrorList {
	t.Helper()
	return func(errs field.ErrorList, err error) field.ErrorList {
		t.Helper()
		if err != nil {
			t.Fatalf("validating: %v", err)
		}
		return errs
	}
}
