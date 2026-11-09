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

package schemes_test

import (
	"errors"
	"reflect"
	"slices"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"reconciler.io/dies/v2/apis"
	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	"reconciler.io/dies/v2/die/schemes"
)

var deployment = schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}

func composite(t *testing.T) *schemes.Composite {
	t.Helper()
	cs := schemes.NewComposite()
	if err := apis.SchemeBuilder.AddToComposite(cs); err != nil {
		t.Fatal(err)
	}
	return cs
}

// A kind has a type of each category, each in its own scheme.
func TestCompositeNewOfEachCategory(t *testing.T) {
	cs := composite(t)
	for c, want := range map[schemes.Category]runtime.Object{
		schemes.Types:   &appsv1.Deployment{},
		schemes.Dies:    &dieappsv1.DeploymentDie{},
		schemes.DieMuts: &dieappsv1.DeploymentDieMut{},
	} {
		got, err := cs.New(deployment, c)
		if err != nil {
			t.Fatalf("%s: %v", c, err)
		}
		if reflect.TypeOf(got) != reflect.TypeOf(want) {
			t.Errorf("%s: new %T, want %T", c, got, want)
		}
		if typ, ok := cs.TypeOf(deployment, c); !ok || typ != reflect.TypeOf(want).Elem() {
			t.Errorf("%s: type %v, want %v", c, typ, reflect.TypeOf(want).Elem())
		}
	}
	// a list die is the list's die of both categories of dies
	list := deployment.GroupVersion().WithKind("DeploymentList")
	for c, want := range map[schemes.Category]reflect.Type{
		schemes.Types:   reflect.TypeFor[appsv1.DeploymentList](),
		schemes.Dies:    reflect.TypeFor[dieappsv1.DeploymentListDie](),
		schemes.DieMuts: reflect.TypeFor[dieappsv1.DeploymentListDie](),
	} {
		if typ, ok := cs.TypeOf(list, c); !ok || typ != want {
			t.Errorf("%s: list type %v, want %v", c, typ, want)
		}
	}
	if _, err := cs.New(deployment, schemes.Category(7)); !errors.Is(err, schemes.ErrCategory) {
		t.Errorf("error %v, want ErrCategory", err)
	}
	if _, ok := cs.TypeOf(deployment, schemes.Category(-1)); ok {
		t.Error("a type of a category that isn't one")
	}
}

// An object of any category has its kind, and the scheme it is registered
// in, which a client of that type takes.
func TestCompositeSchemeFor(t *testing.T) {
	cs := composite(t)
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(deployment)
	for _, tc := range []struct {
		name string
		obj  runtime.Object
		c    schemes.Category
		kind schema.GroupVersionKind
	}{
		{"a type", &appsv1.Deployment{}, schemes.Types, deployment},
		{"an immutable die", dieappsv1.DeploymentBlank, schemes.Dies, deployment},
		{"a mutable die", dieappsv1.DeploymentBlank.DieMutable(), schemes.DieMuts, deployment},
		{"a list die", dieappsv1.NewDeploymentList(), schemes.DieMuts, deployment.GroupVersion().WithKind("DeploymentList")},
		{"a meta type", &metav1.Status{}, schemes.Types, schema.GroupVersionKind{Version: "v1", Kind: "Status"}},
		{"an unstructured object", u, schemes.Types, deployment},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s, c, err := cs.SchemeFor(tc.obj)
			if err != nil {
				t.Fatal(err)
			}
			if c != tc.c || s != cs.Scheme(tc.c) {
				t.Errorf("category %s, want %s", c, tc.c)
			}
			kinds, _, err := cs.ObjectKinds(tc.obj)
			if err != nil {
				t.Fatal(err)
			}
			if !slices.Contains(kinds, tc.kind) {
				t.Errorf("kinds %v, want %v", kinds, tc.kind)
			}
			// the scheme alone knows the object, as a client would ask it
			if _, _, err := s.ObjectKinds(tc.obj); err != nil {
				t.Errorf("the scheme for the object doesn't know it: %v", err)
			}
		})
	}

	if _, _, err := cs.SchemeFor(&dieappsv1.DeploymentList{}); err != nil {
		t.Errorf("an upstream list: %v", err)
	}
	if _, _, err := cs.ObjectKinds(&unregistered{}); err == nil || !runtime.IsNotRegisteredError(err) {
		t.Errorf("error %v, want not registered", err)
	}
	if !cs.Recognizes(deployment) || cs.Recognizes(schema.GroupVersionKind{Group: "example.dev", Version: "v1", Kind: "Widget"}) {
		t.Error("Recognizes is wrong")
	}
}

// unregistered is an object no scheme has.
type unregistered struct{ metav1.TypeMeta }

func (u *unregistered) DeepCopyObject() runtime.Object { c := *u; return &c }

// A builder registers one category with a scheme, as a runtime.Scheme's
// AddToScheme does; AddToScheme itself registers the mutable dies, as
// before.
func TestBuilderAddsACategory(t *testing.T) {
	for c, want := range map[schemes.Category]reflect.Type{
		schemes.Types:   reflect.TypeFor[appsv1.Deployment](),
		schemes.Dies:    reflect.TypeFor[dieappsv1.DeploymentDie](),
		schemes.DieMuts: reflect.TypeFor[dieappsv1.DeploymentDieMut](),
	} {
		s := runtime.NewScheme()
		if err := apis.SchemeBuilder.AddToSchemeAs(c)(s); err != nil {
			t.Fatal(err)
		}
		if got := s.AllKnownTypes()[deployment]; got != want {
			t.Errorf("%s: %v, want %v", c, got, want)
		}
	}
	s := runtime.NewScheme()
	if err := apis.AddToScheme(s); err != nil {
		t.Fatal(err)
	}
	if got := s.AllKnownTypes()[deployment]; got != reflect.TypeFor[dieappsv1.DeploymentDieMut]() {
		t.Errorf("AddToScheme registers %v, want the mutable die", got)
	}

	// a builder of builders, as a group's is of its versions'
	var b schemes.Builder
	b.Register(dieappsv1.SchemeBuilder.Add)
	cs := schemes.NewComposite()
	if err := b.AddToComposite(cs); err != nil {
		t.Fatal(err)
	}
	if _, err := cs.New(deployment, schemes.Dies); err != nil {
		t.Error(err)
	}
	if err := b.Add(runtime.NewScheme(), schemes.Category(3)); !errors.Is(err, schemes.ErrCategory) {
		t.Errorf("error %v, want ErrCategory", err)
	}
}
