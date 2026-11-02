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

package sources_test

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	upstreamkappctrlv1alpha1 "carvel.dev/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	upstreampackagingv1alpha1 "carvel.dev/kapp-controller/pkg/apis/packaging/v1alpha1"
	vendirv1alpha1 "carvel.dev/vendir/pkg/vendir/versions/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kappctrlv1alpha1 "reconciler.io/dies/v2/examples/kapp/kappctrl/v1alpha1"
	packagingv1alpha1 "reconciler.io/dies/v2/examples/kapp/packaging/v1alpha1"
)

// This module is the one place that depends on both the upstream types and
// their clones, so it checks that the clones keep the upstream wire shape.

func TestClonesMatchUpstreamWireShape(t *testing.T) {
	for _, pair := range [][2]reflect.Type{
		{reflect.TypeFor[packagingv1alpha1.PackageInstall](), reflect.TypeFor[upstreampackagingv1alpha1.PackageInstall]()},
		{reflect.TypeFor[packagingv1alpha1.PackageRepository](), reflect.TypeFor[upstreampackagingv1alpha1.PackageRepository]()},
		{reflect.TypeFor[kappctrlv1alpha1.App](), reflect.TypeFor[upstreamkappctrlv1alpha1.App]()},
	} {
		compareShape(t, pair[1].Name(), pair[0], pair[1], map[[2]reflect.Type]bool{})
	}
}

func compareShape(t *testing.T, path string, clone, upstream reflect.Type, seen map[[2]reflect.Type]bool) {
	t.Helper()
	if clone == upstream || seen[[2]reflect.Type{clone, upstream}] {
		return
	}
	seen[[2]reflect.Type{clone, upstream}] = true
	if clone.Kind() != upstream.Kind() {
		t.Errorf("%s: kind %s, upstream %s", path, clone.Kind(), upstream.Kind())
		return
	}
	switch clone.Kind() {
	case reflect.Pointer, reflect.Slice:
		compareShape(t, path+"[]", clone.Elem(), upstream.Elem(), seen)
	case reflect.Map:
		compareShape(t, path+"{}", clone.Elem(), upstream.Elem(), seen)
	case reflect.Struct:
		c, u := wireFields(clone), wireFields(upstream)
		for name, cf := range c {
			uf, ok := u[name]
			if !ok {
				t.Errorf("%s.%s: not upstream", path, name)
				continue
			}
			if cf.Tag.Get("json") != uf.Tag.Get("json") {
				t.Errorf("%s.%s: json tag %q, upstream %q", path, name, cf.Tag.Get("json"), uf.Tag.Get("json"))
			}
			compareShape(t, path+"."+name, cf.Type, uf.Type, seen)
		}
		for name := range u {
			if _, ok := c[name]; !ok {
				t.Errorf("%s.%s: missing from the clone", path, name)
			}
		}
	}
}

func wireFields(t reflect.Type) map[string]reflect.StructField {
	out := map[string]reflect.StructField{}
	for i := range t.NumField() {
		f := t.Field(i)
		name, _, _ := strings.Cut(f.Tag.Get("json"), ",")
		if f.Anonymous && name == "" {
			for k, v := range wireFields(f.Type) {
				out[k] = v
			}
			continue
		}
		if f.IsExported() && name != "-" {
			out[name] = f
		}
	}
	return out
}

func TestUpstreamRoundTripsThroughClonedDies(t *testing.T) {
	upstream := upstreampackagingv1alpha1.PackageInstall{
		TypeMeta:   metav1.TypeMeta{APIVersion: "packaging.carvel.dev/v1alpha1", Kind: "PackageInstall"},
		ObjectMeta: metav1.ObjectMeta{Name: "cert-manager", Namespace: "packages"},
		Spec: upstreampackagingv1alpha1.PackageInstallSpec{
			ServiceAccountName: "installer",
			PackageRef: &upstreampackagingv1alpha1.PackageRef{
				RefName:          "cert-manager.community.tanzu.vmware.com",
				VersionSelection: &vendirv1alpha1.VersionSelectionSemver{Constraints: ">=1.0.0"},
			},
			Values: []upstreampackagingv1alpha1.PackageInstallValues{{
				SecretRef: &upstreampackagingv1alpha1.PackageInstallValuesSecretRef{Name: "values"},
			}},
		},
	}
	want, _ := json.Marshal(upstream)

	d := &packagingv1alpha1.PackageInstallDie{}
	if err := json.Unmarshal(want, d); err != nil {
		t.Fatal(err)
	}
	if got := d.Spec().PackageRef().VersionSelection().Constraints().Get(); got != ">=1.0.0" {
		t.Errorf("constraints %q, read through the cloned vendir die", got)
	}
	// the die's JSON decodes as the upstream object; it omits the nulls the
	// upstream encoding adds for absent values, such as a nil conditions list
	got, _ := json.Marshal(d)
	var back upstreampackagingv1alpha1.PackageInstall
	if err := json.Unmarshal(got, &back); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(back, upstream) {
		t.Errorf("round trip\n got %s\nwant %s", got, want)
	}
}
