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

package unstructured_test

import (
	"encoding/json"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/openapi"

	"reconciler.io/dies/v2/die"
	"reconciler.io/dies/v2/die/base"
	"reconciler.io/dies/v2/die/dietest"
	dieunstructured "reconciler.io/dies/v2/die/unstructured"
)

// thingSchema is a group version's OpenAPI v3 document, as a server publishes
// it, with a kind whose fields merge in each of the ways a server merges.
const thingSchema = `{"components": {"schemas": {
	"dev.acme.v1.Thing": {
		"type": "object",
		"x-kubernetes-group-version-kind": [{"group": "acme.dev", "version": "v1", "kind": "Thing"}],
		"properties": {
			"apiVersion": {"type": "string"},
			"kind": {"type": "string"},
			"metadata": {"allOf": [{"$ref": "#/components/schemas/io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta"}]},
			"spec": {
				"type": "object",
				"properties": {
					"ports": {
						"type": "array",
						"x-kubernetes-list-type": "map",
						"x-kubernetes-list-map-keys": ["name"],
						"items": {"type": "object", "properties": {"name": {"type": "string"}, "port": {"type": "integer"}}}
					},
					"args": {"type": "array", "items": {"type": "string"}},
					"ref": {"allOf": [{"$ref": "#/components/schemas/dev.acme.v1.Ref"}]},
					"granularRef": {"allOf": [{"$ref": "#/components/schemas/dev.acme.v1.Ref"}], "x-kubernetes-map-type": "granular"},
					"selector": {"type": "object", "additionalProperties": {"type": "string"}},
					"minReady": {"type": "integer", "nullable": true},
					"count": {"type": "integer"},
					"template": {"type": "object", "x-kubernetes-preserve-unknown-fields": true}
				}
			}
		}
	},
	"dev.acme.v1.Ref": {
		"type": "object",
		"x-kubernetes-map-type": "atomic",
		"properties": {"name": {"type": "string"}, "key": {"type": "string"}}
	},
	"io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta": {
		"type": "object",
		"properties": {
			"name": {"type": "string"},
			"namespace": {"type": "string"},
			"labels": {"type": "object", "additionalProperties": {"type": "string"}}
		}
	}
}}}`

type fakeOpenAPI struct {
	docs  map[string]string
	calls int
}

func (f *fakeOpenAPI) Paths() (map[string]openapi.GroupVersion, error) {
	f.calls++
	paths := map[string]openapi.GroupVersion{}
	for path, doc := range f.docs {
		paths[path] = fakeGroupVersion{path: path, doc: doc}
	}
	return paths, nil
}

type fakeGroupVersion struct{ path, doc string }

func (g fakeGroupVersion) Schema(string) ([]byte, error) { return []byte(g.doc), nil }
func (g fakeGroupVersion) ServerRelativeURL() string     { return "/openapi/v3/" + g.path }

var thingGVK = schema.GroupVersionKind{Group: "acme.dev", Version: "v1", Kind: "Thing"}

func thingMolds() *dieunstructured.Molds {
	return dieunstructured.NewMolds(&fakeOpenAPI{docs: map[string]string{"apis/acme.dev/v1": thingSchema}})
}

// a thing as the server has it
func serverThing(t *testing.T, d *dieunstructured.Object) *dieunstructured.Object {
	t.Helper()
	m := d.DieMutable()
	dietest.NoError(t, die.FeedJSON(m, []byte(`{
		"apiVersion": "acme.dev/v1", "kind": "Thing",
		"metadata": {"name": "web", "namespace": "default", "labels": {"team": "a"}},
		"spec": {
			"ports": [{"name": "http", "port": 80}],
			"args": ["a", "x"],
			"ref": {"name": "settings", "key": "color"},
			"granularRef": {"name": "settings", "key": "color"},
			"selector": {"app": "web"},
			"minReady": 5,
			"template": {"spec": {"replicas": 1}}
		}
	}`)))
	return m.DieImmutable()
}

// With the server's schema, an unstructured die's apply follows the schema as
// a typed die's follows its type.
func TestMoldsApplyPatches(t *testing.T) {
	blank, err := thingMolds().BlankE(thingGVK)
	if err != nil {
		t.Fatal(err)
	}
	d := serverThing(t, blank).DieStamp(func(d *dieunstructured.ObjectMut) {
		spec := d.Node().Field("spec")
		// an object's nullable field: null; a custom resource's validation
		// rejects null for another, which is left out
		dietest.NoError(t, spec.Field("minReady").Clear())
		dietest.NoError(t, spec.Field("count").Clear())
		// a map entry: left out
		dietest.NoError(t, d.Node().Field("metadata").Field("labels").Field("team").Clear())
		dietest.NoError(t, spec.Field("selector").Field("app").Clear())
		// an atomic list: whole
		dietest.NoError(t, spec.Field("args").Index(1).Set("b"))
		// an atomic struct: whole, a granular one: only what changed
		dietest.NoError(t, spec.Field("ref").Field("key").Set("colour"))
		dietest.NoError(t, spec.Field("granularRef").Field("key").Set("colour"))
		// a list map: by key
		dietest.NoError(t, spec.Field("ports").Item("name", "metrics").Field("port").Set(int64(9090)))
		// below preserve-unknown-fields, as without a schema: a map
		dietest.NoError(t, spec.Field("template").Field("spec").Field("replicas").Clear())
	})
	got, err := die.ReleaseApplyPatchE(d)
	assertJSON(t, got, err, `{
		"apiVersion": "acme.dev/v1", "kind": "Thing",
		"metadata": {"name": "web", "namespace": "default"},
		"spec": {
			"minReady": null,
			"args": ["a", "b"],
			"ref": {"name": "settings", "key": "colour"},
			"granularRef": {"key": "colour"},
			"ports": [{"name": "metrics", "port": 9090}]
		}
	}`)

	// without a schema, the same intent is guessed at: the field is a map
	// entry, left out, and so is every other removal; the atomic list is sent
	// as the element written, and the atomic struct as the field written
	plain := serverThing(t, dieunstructured.Blank).DieStamp(func(d *dieunstructured.ObjectMut) {
		spec := d.Node().Field("spec")
		dietest.NoError(t, spec.Field("minReady").Clear())
		dietest.NoError(t, spec.Field("ref").Field("key").Set("colour"))
	})
	got, err = die.ReleaseApplyPatchE(plain)
	assertJSON(t, got, err, `{
		"apiVersion": "acme.dev/v1", "kind": "Thing",
		"metadata": {"name": "web", "namespace": "default"},
		"spec": {"ref": {"key": "colour"}}
	}`)
}

// The schema is part of the die: feeding, copying and the client-shaped
// operations keep it.
func TestMoldsAreKept(t *testing.T) {
	molds := thingMolds()
	blank, err := molds.BlankE(thingGVK)
	if err != nil {
		t.Fatal(err)
	}
	clearsAsNull := func(t *testing.T, d *dieunstructured.Object) {
		t.Helper()
		got, err := die.ReleaseApplyPatchE(d.DieStamp(func(d *dieunstructured.ObjectMut) {
			dietest.NoError(t, d.Node().Field("spec").Field("minReady").Clear())
		}))
		if err != nil || !strings.Contains(string(got), `"minReady":null`) {
			t.Errorf("apply %s, %v; want the schema's null", got, err)
		}
	}
	clearsAsNull(t, serverThing(t, blank))
	clearsAsNull(t, serverThing(t, blank).DieMutable().DieImmutable())
	clearsAsNull(t, serverThing(t, blank).DeepCopy())

	// a client zeroes the object it decodes into, holder and all
	read, err := blank.DieStampE(func(d *dieunstructured.ObjectMut) error {
		*d = dieunstructured.ObjectMut{}
		return json.Unmarshal([]byte(`{"apiVersion": "acme.dev/v1", "kind": "Thing", "spec": {"minReady": 5}}`), d)
	})
	if err != nil {
		t.Fatal(err)
	}
	clearsAsNull(t, read)

	// For gives a die without one its kind's schema, such as a listed item
	listed := serverThing(t, dieunstructured.Blank)
	withSchema, err := molds.ForE(listed)
	if err != nil {
		t.Fatal(err)
	}
	clearsAsNull(t, withSchema)
}

// A server that converts a kind through Go types, as for built-in kinds and
// aggregated APIs, takes null as absent for any field, so a die of one sends
// null for any field cleared, as a typed die of one does. A CRD's validation
// rejects it, unless the field is nullable. Which a kind is shows in its
// definition's name: a Go package's, or the CRD's group, version and kind.
func TestMoldsGoTypedKinds(t *testing.T) {
	for _, c := range []struct {
		definition string
		null       bool
	}{
		{"io.k8s.api.acme.v1.Thing", true},
		{"io.example.thing-apiserver.pkg.apis.acme.v1.Thing", true},
		{"dev.acme.v1.Thing", false},
	} {
		t.Run(c.definition, func(t *testing.T) {
			doc := strings.ReplaceAll(thingSchema, "dev.acme.v1.Thing", c.definition)
			molds := dieunstructured.NewMolds(&fakeOpenAPI{docs: map[string]string{"apis/acme.dev/v1": doc}})
			blank, err := molds.BlankE(thingGVK)
			if err != nil {
				t.Fatal(err)
			}
			got, err := die.ReleaseApplyPatchE(serverThing(t, blank).DieStamp(func(d *dieunstructured.ObjectMut) {
				dietest.NoError(t, d.Node().Field("spec").Field("count").Clear())
			}))
			if err != nil || strings.Contains(string(got), `"count":null`) != c.null {
				t.Errorf("apply %s, %v; want null %v", got, err, c.null)
			}
		})
	}
}

func TestMoldsErrors(t *testing.T) {
	api := &fakeOpenAPI{docs: map[string]string{"apis/acme.dev/v1": thingSchema}}
	molds := dieunstructured.NewMolds(api)
	if _, err := molds.BlankE(schema.GroupVersionKind{Group: "acme.dev", Version: "v1", Kind: "Other"}); err == nil || !strings.Contains(err.Error(), "no kind Other") {
		t.Errorf("err %v", err)
	}
	// a group version the server didn't serve when the paths were read is
	// looked for again, as for a CRD installed since
	if _, err := molds.BlankE(schema.GroupVersionKind{Group: "new.dev", Version: "v1", Kind: "Thing"}); err == nil || !strings.Contains(err.Error(), "no OpenAPI schema for new.dev/v1") {
		t.Errorf("err %v", err)
	}
	api.docs["apis/new.dev/v1"] = strings.ReplaceAll(thingSchema, `"group": "acme.dev"`, `"group": "new.dev"`)
	if _, err := molds.BlankE(schema.GroupVersionKind{Group: "new.dev", Version: "v1", Kind: "Thing"}); err != nil {
		t.Errorf("err %v", err)
	}
	// once at first, then once for each miss; a hit reads nothing
	if _, err := molds.BlankE(thingGVK); err != nil {
		t.Errorf("err %v", err)
	}
	if api.calls != 3 {
		t.Errorf("paths read %d times, want 3", api.calls)
	}
}

// Patches fed to an unstructured die follow its mold, and release again as
// the same patch; without a mold, lists are set whole.
func TestMoldsFeedPatches(t *testing.T) {
	blank, err := thingMolds().BlankE(thingGVK)
	if err != nil {
		t.Fatal(err)
	}
	apply := []byte(`{
		"apiVersion": "acme.dev/v1", "kind": "Thing",
		"metadata": {"name": "web", "namespace": "default"},
		"spec": {
			"minReady": null,
			"args": ["a", "b"],
			"ref": {"name": "settings", "key": "colour"},
			"ports": [{"name": "metrics", "port": 9090}]
		}
	}`)
	fed := serverThing(t, blank).DieMutable()
	dietest.NoError(t, die.FeedApplyPatch(fed, apply))
	got, err := die.ReleaseApplyPatchE(fed)
	assertJSON(t, got, err, string(apply))
	// the list map's other elements are kept, as the element is found by key
	ports := fed.Node().Field("spec").Field("ports")
	if ports.Len() != 2 || ports.Item("name", "http").Field("port").Get() != int64(80) {
		t.Errorf("ports %v", ports.Get())
	}

	merge := []byte(`{"metadata": {"labels": {"team": null}}, "spec": {"minReady": 7, "ports": [{"name": "http", "port": 8080}]}}`)
	fed = serverThing(t, blank).DieMutable()
	dietest.NoError(t, die.FeedMergePatch(fed, merge))
	got, err = die.ReleaseMergePatchE(fed)
	assertJSON(t, got, err, string(merge))
	if fed.Node().Field("metadata").Field("labels").Field("team").Has() || fed.Node().Field("spec").Field("ports").Len() != 1 {
		t.Errorf("fed %v", fed.Node().Get())
	}

	// without a mold, the list map's element can't be found by key, so the
	// list is set whole
	plain := serverThing(t, dieunstructured.Blank).DieMutable()
	dietest.NoError(t, die.FeedApplyPatch(plain, apply))
	if ports := plain.Node().Field("spec").Field("ports"); ports.Len() != 1 {
		t.Errorf("ports %v, want the list set whole", ports.Get())
	}
}

// The unstructured API's fluent reads panic with a die's error, and their E
// variants return it, as writes do.
func TestUnstructuredEVariants(t *testing.T) {
	d := serverThing(t, dieunstructured.Blank)
	if _, err := d.Node().Field("spec").Field("args").IndexE(-1); err == nil {
		t.Error("IndexE(-1) returned no error")
	}
	if _, err := d.Node().ItemE("name"); err == nil {
		t.Error("ItemE with an odd number of arguments returned no error")
	}
	if _, err := d.Node().FindE("{.spec[}"); err == nil {
		t.Error("FindE of a bad expression returned no error")
	}
	if _, err := thingMolds().BlankE(schema.GroupVersionKind{Group: "acme.dev", Version: "v1", Kind: "Other"}); err == nil {
		t.Error("BlankE of an unknown kind returned no error")
	}
	defer func() {
		if _, ok := recover().(*base.Error); !ok {
			t.Error("Blank of an unknown kind didn't panic with a die's error")
		}
	}()
	thingMolds().Blank(schema.GroupVersionKind{Group: "acme.dev", Version: "v1", Kind: "Other"})
}
