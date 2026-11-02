# Dies - immutable, fluent builders for Kubernetes resources <!-- omit in toc -->

![CI](https://github.com/reconcilerio/dies/workflows/CI/badge.svg?branch=main)
[![codecov](https://codecov.io/gh/reconcilerio/dies/branch/main/graph/badge.svg?token=vyXLcPBdV3)](https://codecov.io/gh/reconcilerio/dies)
[![Go Reference](https://pkg.go.dev/badge/reconciler.io/dies/v2.svg)](https://pkg.go.dev/reconciler.io/dies/v2)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](CODE_OF_CONDUCT.md)

Dies build and change Kubernetes resources fluently, in tests and in controllers. This project has dies for every built-in Kubernetes API, and a generator, `diegen`, that creates dies for custom resources.

A die records what you change, not just the resulting value. It knows whether you left a field alone, set it to its zero value, or cleared it. So it can send exactly that as a server-side apply or a JSON merge patch, without sending stale values back to the server or clearing what other writers own.

In a workshop, a die is the tool in a press that gives material its shape. A blank is fed in, stamped, and released as a formed part. The names here follow that:

- **Blank:** every die starts from one, such as `DeploymentBlank`.
- **Feed:** `DieFeed…` gives a die its material, such as an object, JSON, YAML, or what the server has.
- **Stamp:** `DieWith`, `DieStamp` and the field accessors shape a die, and it records each change as intent.
- **Release:** `DieRelease…` gives you the part. `DieReleaseApplyPatch` and `DieReleaseMergePatch` release only what was stamped, as a patch.
- **Jig:** a jig holds dies of different shapes the same way, so the same tools work on each.

v2 is a rewrite, and its API is not compatible with v0. The v0 API is documented at [pkg.go.dev/reconciler.io/dies](https://pkg.go.dev/reconciler.io/dies). We skipped v1 so that both versions can be used in a project during the transition.


- [Getting started](#getting-started)
- [Using dies](#using-dies)
	- [Fields](#fields)
	- [Immutable and mutable dies](#immutable-and-mutable-dies)
	- [Feeding and releasing](#feeding-and-releasing)
	- [Errors: fluent or E](#errors-fluent-or-e)
- [Patch with intent](#patch-with-intent)
	- [Server-side apply](#server-side-apply)
	- [Merge patches](#merge-patches)
	- [Choosing what to send](#choosing-what-to-send)
- [Dies with Kubernetes clients](#dies-with-kubernetes-clients)
	- [Apply configurations](#apply-configurations)
	- [Lists](#lists)
- [Unstructured dies](#unstructured-dies)
	- [Molds: unstructured dies with the server's schema](#molds-unstructured-dies-with-the-servers-schema)
- [Jigs](#jigs)
- [JSONPath](#jsonpath)
- [Die methods](#die-methods)
- [Creating dies](#creating-dies)
	- [diegen](#diegen)
	- [Dies for types from another package](#dies-for-types-from-another-package)
	- [Dies for your own types](#dies-for-your-own-types)
	- [Markers](#markers)
	- [Generated tests](#generated-tests)
- [Conformance](#conformance)
- [Community](#community)
	- [Code of Conduct](#code-of-conduct)
	- [Communication](#communication)
	- [Contributing](#contributing)
- [Acknowledgements](#acknowledgements)
- [License](#license)


## Getting started

```sh
go get reconciler.io/dies/v2
```

Dies for the built-in APIs are in `reconciler.io/dies/v2/apis/<group>/<version>`, following the layout of `k8s.io/api`. Each die package aliases the upstream package's types and constants, so one import gives you both the types and their dies.

```go
import (
	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	diecorev1 "reconciler.io/dies/v2/apis/core/v1"
)

deployment := dieappsv1.DeploymentBlank.DieWith(func(d *dieappsv1.DeploymentDie) {
	d.Metadata().Namespace().Set("default")
	d.Metadata().Name().Set("web")
	d.Spec().Replicas().Set(2)
	d.Spec().Template().Spec().Containers().At("app").DieWith(func(c *diecorev1.ContainerDie) {
		c.Image().Set("registry.example/app:1")
		c.Env().At("MODE").Value().Set("production")
		c.Ports().At(diecorev1.ContainerPortKey{ContainerPort: 8080}).Name().Set("http")
	})
})

obj := deployment.DieRelease() // an appsv1.Deployment
```

## Using dies

Every type with a die has a blank, `<Type>Blank`, to start from. Blanks for objects already have their `apiVersion` and `kind`.

### Fields

Each field has one accessor method, named after the field. What it returns depends on the field's type:

| Field                                  | Accessor returns                                   | For example                                 |
| -------------------------------------- | -------------------------------------------------- | ------------------------------------------- |
| A struct with a die                    | The child die                                      | `d.Spec().Template()`                       |
| A list map, such as containers by name | A `die.ListMap`, whose elements are dies           | `podSpec.Containers().At("app")`            |
| A map with string keys                 | A `die.Map`                                        | `d.Metadata().Labels().Set("tier", "web")`  |
| Anything else, held by value           | A `die.Value`                                      | `c.Image().Get()`, `c.Image().Set("app:2")` |
| Anything else, held by pointer         | A `die.Optional`, which reads and writes the value | `d.Spec().Replicas().Set(2)`, `.Ptr()`      |

Every leaf handle has `Get`, `Has`, `Set` and `Clear`, and lists and maps also have `Empty`:

- `Has` reports whether the field is present, as the server sent it or as you set it. A zero value can be present.
- `Set` records the value, even when it's the zero value.
- `Clear` records that the field is absent.
- `Empty` records a list or map that is present and empty.

A list map's elements are found by their keys. `At` returns the element's die, and its first write creates the element with its keys set. List maps also have `Has`, `Remove`, `Keys`, `Len` and `All`. A list map with more than one key, such as container ports keyed by port and protocol, takes a generated key type, such as `ContainerPortKey`. A key with a default, such as a port's protocol, which defaults to `TCP`, matches whether it's set or not.

Lists that aren't list maps, such as a container's `args`, are values, which are read and written whole.

### Immutable and mutable dies

Blanks, and the dies they return, are immutable: a change returns a new die and leaves the original unchanged. That lets you capture a die partway through and derive variations from it, which is what makes dies useful for test fixtures:

```go
base := dieappsv1.DeploymentBlank.DieWith(func(d *dieappsv1.DeploymentDie) {
	d.Metadata().Name().Set("web")
	d.Spec().Template().Spec().Containers().At("app").Image().Set("app:1")
})

scaled := base.DieWith(func(d *dieappsv1.DeploymentDie) {
	d.Spec().Replicas().Set(3)
})
// base is unchanged
```

The die passed to a `DieWith` callback is mutable, so you can change it and the fields within it in place. Writing to a field within an immutable die panics, because the change would have nowhere to go. A field directly on an immutable root can be written, and returns a new root, such as `diecorev1.ConfigMapBlank.Immutable().Set(true)`.

- `DieImmutable(false)` returns a mutable copy of a die, and `DieImmutable(true)` an immutable one.
- `DeepCopy` snapshots a die.
- `DieStamp(func(r *T))` changes the underlying Go value directly. The die works out what changed, and records it as intent.

### Feeding and releasing

Dies convert to and from the forms Kubernetes code uses:

| Feed                                                            | Release                                                             |
| --------------------------------------------------------------- | ------------------------------------------------------------------- |
| `DieFeed(T)`, `DieFeedPtr(*T)`                                  | `DieRelease() T`, `DieReleasePtr() *T`                              |
| `DieFeedJSON([]byte)`                                           | `DieReleaseJSON() []byte`                                           |
| `DieFeedYAML([]byte)`                                           | `DieReleaseYAML() []byte`                                           |
| `DieFeedRawExtension(runtime.RawExtension)`                     | `DieReleaseRawExtension() runtime.RawExtension`                     |
| `DieFeedUnstructured(*unstructured.Unstructured)`, objects only | `DieReleaseUnstructured() *unstructured.Unstructured`, objects only |
| `DieFeedApplyPatch([]byte)`, root dies only                     | `DieReleaseApplyPatch() []byte`, root dies only                     |
| `DieFeedMergePatch([]byte)`, root dies only                     | `DieReleaseMergePatch() []byte`, root dies only                     |

When a die decodes JSON or YAML, it records which fields were present and keeps any fields its Go type doesn't have, such as fields added in a newer version of the API. Encoding it again gives back what it read: an update from a die doesn't drop fields it doesn't know about.

Dies are `json.Marshaler`s and `json.Unmarshaler`s. Dies for objects also implement `runtime.Object` and `metav1.Object`, so Kubernetes libraries and clients take them directly.

For tests, `dietest.FeedYAMLFile(d, "testdata/deployment.yaml")` feeds a die from a file.

### Errors: fluent or E

Every method that can fail comes in two forms:
- **E variant,** such as `SetE`, which returns the error. Use it in a controller, for content you don't control.
- **Fluent,** such as `Set`, which calls the E variant and panics with its error. Use it when you know the content you're working on, as in tests and fixtures, where reading well matters most.

```go
// E variant: returns the error
if _, err := d.Spec().Replicas().SetE(3); err != nil {
	return err
}

// fluent: panics if the write can't be made
d.Spec().Replicas().Set(3)
```

Nothing recovers a panic: the E variant does the work and returns its error, and the fluent form only adds the panic. `DieWithE` takes functions that return errors, such as a client operation or a block of E calls, and returns the first one, along with the die as the functions left it:

```go
updated, err := current.DieWithE(func(d *dieappsv1.DeploymentDie) error {
	if _, err := d.Spec().Replicas().SetE(3); err != nil {
		return err
	}
	_, err := d.Spec().Template().Spec().Containers().At("app").Image().SetE(image)
	return err
})
```

Use `errors.As` with `*die.Error` to tell a die's error from another, and `errors.Is` with `die.ErrImmutableChild` or `die.ErrReadOnlyJig` for those two. A cause, such as a JSON syntax error, is wrapped.

Two kinds of failure have no E variant:
- **Interface implementations,** such as `metav1.Object`'s setters, can't return an error. They panic on an immutable die. Write through the die's own accessors, which have E variants, instead.
- **A broken die definition panics,** whether generated or hand-written. Examples are an unknown die tag, or a field table that doesn't match its type. These are bugs in the die code, not in the data it's given, and generated dies' tests catch them.

## Patch with intent

A die is a value plus a record of what you changed: its intent. Patches are encoded from the intent, never from the whole Go value, so:

- a field you didn't touch is never sent, even when its Go type would encode it as `null`
- a zero value you set, such as `replicas: 0` or `paused: false`, is sent
- a field you cleared is sent as removed

Reads see the value with the intent applied. Feeding a root die replaces its value and discards its intent. `DieIntent()` describes the intent, one change per line, which helps when debugging.

Patches can be fed too. `DieFeedApplyPatch` and `DieFeedMergePatch` take a patch as the die's intent, over the value it already has, so releasing the same kind of patch sends it back. Use them to inspect a patch from elsewhere through a die's accessors, or to replay one against another value. A merge patch fed to the value it was made against reproduces the change, lists included, which a merge patch sends whole. Fields the die's type doesn't have aren't fed. An unstructured die follows its mold where it has one; without one, its lists are set whole.

### Server-side apply

`DieReleaseApplyPatch()` encodes the intent as a server-side apply request, identified by the object's `apiVersion`, `kind`, `name` and `namespace`. The manager applying it owns exactly the fields it sent:

- `Clear` on a field sends `null` where the server accepts it: the manager owns the field, and it's absent. A server that converts objects through Go types accepts `null` for any field; that covers built-in types, and aggregated APIs such as sample-apiserver's. A custom resource's validation rejects `null` unless the field is `nullable`, so for a custom resource, `Clear` sends `null` only for fields marked `+nullable`, and otherwise leaves the field out, which removes it only if this manager owned it. Typed dies treat types from `k8s.io` as served from Go types. Unstructured dies with the server's schema tell from the schema itself.
- `Empty` sends `[]` or `{}`: the manager owns an empty list or map.
- Removing a list-map element or a map entry leaves it out of the apply. An apply can't remove what another manager owns. The server removes what this manager owned and no longer applies. Sending `null` for a map entry would set it to `""` instead.
- A struct the server replaces whole, such as a `LabelSelector` or an `ObjectReference` (`+structType=atomic`), keeps its child die, but is sent whole.

```go
status := dieappsv1.DeploymentBlank.DieWith(func(d *dieappsv1.DeploymentDie) {
	d.Metadata().Namespace().Set(namespace)
	d.Metadata().Name().Set(name)
	d.Status().ObservedGeneration().Set(generation)
	d.Status().Conditions().At(dieappsv1.DeploymentAvailable).DieWith(func(c *dieappsv1.DeploymentConditionDie) {
		c.Status().Set(diecorev1.ConditionTrue)
		c.Reason().Set("Reconciled")
	})
}).DieSubresource("status")

body, err := status.DieReleaseApplyPatchE()
// ...
// the object as the server has it after the apply, such as its new resource
// version: committed first, so the server's response isn't overlaid with the
// apply's intent
applied, err := status.DieCommit().DieWithE(func(d *dieappsv1.DeploymentDie) error {
	return c.Status().Patch(ctx, d, client.RawPatch(types.ApplyPatchType, body), client.FieldOwner("my-controller"))
})
```

A client decodes the server's response into the object you pass it, so that object has to be mutable, and you need a reference to it to see the result. `DieWithE` runs a client operation on a mutable copy of an immutable die, and returns the result as a new immutable die, with the operation's error.

Decoding replaces the die's value but keeps its intent, applied on top. To see exactly what the server returned, pass a die without intent: one read from the server, or one committed with `DieCommit`. This matters after an apply, because the server doesn't always do what the apply asked. A field the apply left out, for example, stays as another manager set it.

An apply is all of a manager's intent. If a manager applies from server state with one change, it gives up every other field it owned, and the server removes them. Build applies from a blank, as a reconciler does above, or make the change as a different manager.

### Merge patches

`DieReleaseMergePatch()` encodes the intent as a JSON merge patch (RFC 7386) against the value the die was fed. Only what you changed is sent, so concurrent changes to other fields survive. A merge patch can only replace a list whole, so a change to a list sends the whole list as the die has it.

```go
// the current state, read from the server
current, err := dieappsv1.DeploymentBlank.DieWithE(func(d *dieappsv1.DeploymentDie) error {
	return c.Get(ctx, key, d)
})
// ...

// the update, a change to the current state
update := current.DieWith(func(d *dieappsv1.DeploymentDie) {
	d.Spec().Replicas().Set(0)
	d.Metadata().Labels().Remove("canary")
})

// the patch, from the update's intent against the current state
body, err := update.DieReleaseMergePatchE()
// ...

// the patched state, as the server has it after the patch
patched, err := current.DieWithE(func(d *dieappsv1.DeploymentDie) error {
	return c.Patch(ctx, d, client.RawPatch(types.MergePatchType, body))
})
```

The patch is sent with the current die, which has no intent, so the patched die is exactly what the server returned. Each step is a die of its own, and none of them changes another.

### Choosing what to send

These methods change which part of the intent a die sends:

- `DieSubresource("status")` keeps only the intent for a subresource.
- `DieMain()` drops the intent for every subresource.
- `DieOnly(locations...)` and `DieExcept(locations...)` keep or drop the intent at locations, such as `d.DieOnly(d.Spec().Replicas(), d.Metadata().Labels())`.
- `DieCommit()` folds the intent into the value, as after a successful write.
- `DieFields()` lists the fields a die manages. Each entry is a location too.

Intent that isn't kept is folded into the value, so the die's value is unchanged, identity included, and only the kept intent is sent.

## Dies with Kubernetes clients

Each die package registers its object dies, and the list dies for them, in a scheme. Use `AddToScheme` from one version package, such as `dieappsv1.AddToScheme`, from one group, such as `reconciler.io/dies/v2/apis/apps`, or from every group, with `reconciler.io/dies/v2/apis`. With the dies in a client's scheme, you read and write objects as dies:

```go
scheme := runtime.NewScheme()
_ = apis.AddToScheme(scheme)

config := rest.CopyConfig(restConfig)
// controller-runtime sends the types client-go knows as protobuf, unless told
// otherwise; dies are JSON
config.ContentType = runtime.ContentTypeJSON
c, err := client.New(config, client.Options{Scheme: scheme})

d, err := dieappsv1.DeploymentBlank.DieWithE(func(d *dieappsv1.DeploymentDie) error {
	return c.Get(ctx, client.ObjectKey{Namespace: "default", Name: "web"}, d)
})
// ...
replicas := d.Spec().Replicas().Get()
```

A client strips the type metadata of a typed object as it decodes it. Dies don't treat that as intent, and an apply patch fills in the object's `apiVersion` and `kind` from its die type.

### Apply configurations

Dies for types that client-go has apply configurations for have methods to convert between the two:

- `DieReleaseApplyConfiguration()` returns the intent as client-go's apply configuration, for clients and tools that take one.
- `DieFeedApplyConfiguration(ac)` feeds one in as intent, such as what client-go's `Extract` functions find a manager owns.
- `DieStampApplyConfiguration(fn)` changes the intent through one.

Each die package aliases its apply configuration types and `Extract` functions, too. For any root die, with or without an apply configuration type of its own, `base.ApplyConfigurationOf(d)` adapts its intent for controller-runtime's `client.Apply`, and `Result()` has what the server returned.

As with other client operations, the client decodes the server's response into the apply configuration you pass, so keep a reference to it:

```go
ac := base.ApplyConfigurationOf(deployment)
err := c.Apply(ctx, ac, client.FieldOwner("my-controller"))
// ...
// the object as the server has it after the apply
actual := dieappsv1.DeploymentBlank.DieFeedJSON(ac.Result())
```

### Lists

A list die, such as `dieappsv1.DeploymentListDie`, is a plain Kubernetes list whose items are dies. Clients list into it, and `meta.ExtractList` and `meta.SetList` work with it.

```go
list := diecorev1.NewConfigMapList()
err := c.List(ctx, list, client.InNamespace("default"))
for _, cm := range list.Items {
	_ = cm.Data().At("color")
}
```

## Unstructured dies

`reconciler.io/dies/v2/die/unstructured` has dies for objects with no Go type, such as custom resources you don't have types for. An unstructured die is a `runtime.Unstructured`, so controller-runtime's client takes it without any registration. It's immutable or mutable, records intent and encodes patches like any other die. You traverse it through `Node`s:

```go
d := dieunstructured.BlankFor("example.dev/v1", "Widget").DieWith(func(d *dieunstructured.Object) {
	d.SetNamespace("default")
	d.SetName("web")
	ports := d.Node().Field("spec").Field("ports")
	ports.Item("name", "http").Field("port").Set(int64(80))
})
```

- `Field` steps into an object field, `Item` into a list element by its keys, and `Index` into a list element by position.
- `Get`, `Has`, `Set`, `Clear`, `EmptyObject`, `EmptyList`, `Keys`, `Len` and `Type` read and write a node.
- `dieunstructured.JigAs[*diecorev1.ContainerDie](node)` gives a typed die for part of an unstructured object, whose writes go to that object, and `GetAs` gives a typed copy.
- `dieunstructured.NewList` and `ListFrom` make unstructured list dies.

### Molds: unstructured dies with the server's schema

Without a mold, an unstructured die has to guess how the server merges each part of the object in an apply:
- Every object is treated as a map, so a field it clears is left out.
- A list addressed by keys is treated as a list map, so if the server's list is atomic, the apply replaces the whole list with what the die sent.
- It can't recognize atomic structs, so it sends only the fields that changed, and the server replaces the struct with them.

A mold gives material its shape. Here, a mold is a kind's OpenAPI v3 schema as the server publishes it, and it gives an unstructured blank the structure the server merges by. `dieunstructured.Molds` makes blanks in their kind's mold, so their applies merge as a typed die's do:

```go
dc, err := discovery.NewDiscoveryClientForConfig(restConfig)
// ...
molds := dieunstructured.NewMolds(dc.OpenAPIV3())

blank, err := molds.BlankE(schema.GroupVersionKind{Group: "example.dev", Version: "v1", Kind: "Widget"})
// ...
widget := blank.DieWith(func(d *dieunstructured.Object) {
	d.SetNamespace("default")
	d.SetName("web")
	d.Node().Field("spec").Field("match").Field("note").Clear()
})
```

Dies make no requests of their own: `Molds` reads the schemas through the OpenAPI client you give it, such as a discovery client's `OpenAPIV3()`, or client-go's cached client in front of it. It reads each group version's schema once, when it first makes a blank of one of its kinds, and keeps it until the server changes it. `molds.For(d)` gives an unstructured die you already have, such as an item of a list a client listed into, its kind's mold. `dieunstructured.BlankFor(apiVersion, kind)` makes a blank without a mold.

Caveats:

- **A die's mold is the schema the server served when its blank was made.** Make a blank again to pick up a changed CRD.
- **A CRD's schema is published shortly after the CRD is established,** not with it, so making a blank right after installing a CRD can fail until the schema appears.
- **The die still guesses where the schema doesn't say,** such as below `x-kubernetes-preserve-unknown-fields`.
- **Only apply patches use the mold.** Merge patches don't need it: RFC 7386 merges objects and maps alike, and replaces lists.
- **Defaults aren't used,** so a list-map key left to its default, such as a port's protocol, doesn't match the element with the default set.
- **Subresources aren't known,** so `DieSubresource` and `DieMain` have no effect.
- **A client zeroes the object it decodes into, mold included.** `DieWith` and `DieWithE` keep the mold. A die passed to a client directly loses it, as do the items of a list a client lists into; `For` gives them their mold back.

## Jigs

A jig holds material in place so that a tool can work on it. A jig here holds dies of different types as one type, such as anything with a pod template, so the same code works on all of them. Each binding pairs a location of the jig's type with a location of a source type. Both are navigated from blanks, so a mismatched binding doesn't compile:

```go
var PodSpecableJig = jigs.NewJig[PodSpecable, *PodSpecableDie]().
	For(dieappsv1.DeploymentBlank,
		jigs.Bind(PodSpecableBlank.Metadata(), dieappsv1.DeploymentBlank.Metadata()),
		jigs.Bind(PodSpecableBlank.Template(), dieappsv1.DeploymentBlank.Spec().Template()),
	).
	ForUnknownKinds(
		jigs.BindPath(PodSpecableBlank.Metadata(), ".metadata"),
		jigs.BindPath(PodSpecableBlank.Template(), ".spec.template"),
	)

if ps, ok := PodSpecableJig.Of(workload); ok {
	ps.Template().Spec().Containers().At("app").Image().Set("app:2")
}
```

A jig holds nothing of its own. Its writes go to the source, and the source's patches carry them. `Of` takes any die, unstructured dies included, and `OfUnstructured` takes an `*unstructured.Unstructured`, which the jig changes in place. `jigs.BindVia` converts between shapes that differ. See [`v2/examples/jigs/v1`](v2/examples/jigs/v1) for a complete example.

## JSONPath

`reconciler.io/dies/v2/die/jsonpath` finds locations within a die with kubectl's JSONPath expressions. A match is a live location within the die, not a copy, so writes through a match are recorded as intent. `$` is the root die, and `.` is the die you query, so a function written for a `PodSpecDie` can also reach the object the pod spec is part of:

```go
images, err := jsonpath.ValuesE[string](podSpec, ".containers[*].image")
namespace, err := jsonpath.ValuesE[string](podSpec, "$.metadata.namespace")
containers, err := jsonpath.DiesE[*diecorev1.ContainerDie](podSpec, ".containers[?(@.name=='app')]")
```

A list-map element is found by its keys, however the expression found it, so a write through it merges by key. A path with no wildcards or filters resolves even when its value isn't set, so a write through it creates it. A field the type doesn't have is an error, not a silent miss.

## Die methods

Every die has these methods. `die.Die` is any die, whatever its type, for functions that take one; `die.DieTyped` is any die with type metadata, which is also a `runtime.Object`; `die.DieObject` is any object die, which is also a `metav1.Object`; and `die.DieList` is any list die. Generic code constrains them by `base.DieOf[T, D]`, `die.DieTypedOf[T, D]` (with the `APIVersion()` and `Kind()` accessors), `die.DieObjectOf[T, D]` (with `Metadata()`) and `die.DieListOf[D]` (with `Metadata()`, the list's own `*metav1.ListMeta`, and `DieItems()`), which carry their types. An unstructured die is a `die.DieObject`, but has no typed accessors. Fields' accessors use the field names, so die methods start with `Die`, apart from `DeepCopy`, `Equal` and the methods of the Kubernetes interfaces that object dies implement. Each method that can fail also has an E variant that returns its error; see [Errors: fluent or E](#errors-fluent-or-e).

| Method                                | Does                                                                                                                |
| ------------------------------------- | ------------------------------------------------------------------------------------------------------------------- |
| `DieWith(fns...)`                     | Passes a mutable die to each function, and returns the result                                                       |
| `DieWithE(fn)`                        | Passes a mutable die to a function that can fail, such as a client operation, and returns the result with its error |
| `DieStamp(fn)`                        | Passes the underlying value to a function, and records what changed as intent                                       |
| `DieImmutable(bool)`                  | Returns an immutable or mutable copy                                                                                |
| `DieFeed…` / `DieRelease…`            | Convert to and from Go values, JSON, YAML, raw extensions and unstructured objects                                  |
| `DieHas()` / `DieClear()`             | Report and change whether the die's value is present                                                                |
| `DieReleaseApplyPatch()` / `DieReleaseMergePatch()` | Encode the intent as a server-side apply or a JSON merge patch, which `DieFeedApplyPatch` and `DieFeedMergePatch` feed back as intent |
| `DieSubresource(name)` / `DieMain()`  | Keep the intent for a subresource, or drop it                                                                       |
| `DieOnly(at...)` / `DieExcept(at...)` | Keep or drop the intent at locations                                                                                |
| `DieCommit()`                         | Fold the intent into the value                                                                                      |
| `DieIntent()`                         | Describe the intent                                                                                                 |
| `DieFields()`                         | List the fields the die manages                                                                                     |
| `DiePath()` / `DieRootType()`             | The die's location within its root, and the type of the root's value                                                |
| `DieRoot()`                               | The root die the die is within, as a `die.Die`, sharing its state, such as `d.DieRoot().(metav1.Object)` for an object's identity |
| `DieMutable()`                        | Report whether the die is mutable                                                                                   |
| `DeepCopy()` / `Equal(d)`             | Copy a die, or compare current values                                                                               |

Object dies also have `DieFeedUnstructured` and `DieReleaseUnstructured`, implement `runtime.Object` and `metav1.Object`, and have `APIVersion()`, `Kind()` and `Metadata()` accessors.

## Creating dies

Dies are generated by `diegen`, from Go types and the markers on them.

### diegen

Add diegen as a tool of your module, so it's always the same version as the dies it generates code for:

```sh
go get -tool reconciler.io/dies/v2/diegen
```

Run it with `go generate`, from a directive in a package above your die packages:

```go
//go:generate go tool diegen -header ../hack/boilerplate.go.txt ./...
```

diegen generates dies for the packages matching its arguments. It prints its progress, and only rewrites files that change. `-header` names a file of Go comments, such as a license, to put at the top of each generated file, above the `// Code generated by diegen. DO NOT EDIT.` line.

In each die package, it writes:

- `zz_generated.dies.go`: the dies, their blanks and accessors, scheme registration and list dies
- `zz_generated.dies.types.go`: the package's types, aliased or cloned from another package (if any)
- `zz_generated.dies.applyconfigurations.go`: the package's apply configuration types and `Extract` functions, aliased (if any)
- `zz_generated.dies_test.go`: [tests](#generated-tests) of the dies against their types

### Dies for types from another package

A die package for another package's types needs only a `doc.go`. The `+die:types` marker says where the types come from:

```go
// +die:types:package=github.com/example/widgets/api/v1
package v1
```

- **Aliased** (`mode=alias`, the default): the die package aliases the upstream package's types and constants, so one import gives you both.
- **Cloned** (`mode=clone`): the die package copies the types, with their docs, markers and constants, so it doesn't depend on the upstream package or its dependencies. The package is resolved at a version given with it, `package=…@v1.2.3`, or in the module a `from=` path names.

Each type that needs a die is marked as the upstream package marks it. The markers that change a die's API or its patches, such as `+listType`, `+listMapKey`, `+mapType`, `+structType`, `+default`, `+nullable` and kubebuilder's `+kubebuilder:subresource:status`, are read from the upstream source, in their plain, `k8s:` and `kubebuilder:` forms. What upstream doesn't say is added with metadata declared on the blank identifier, which adds markers to a type without declaring anything:

```go
// +die:field:name=selector,structType=atomic
type _ = v1.WidgetSpec
```

If the upstream package provides apply configurations, such as from applyconfiguration-gen, `+die:applyconfigurations:package=…` aliases them and adds the apply configuration methods to its dies. For the built-in APIs, client-go's apply configurations are used.

### Dies for your own types

Mark your own types directly:

```go
// Widget runs a pod template with named ports.
//
// +die:object:apiVersion=example.dev/v1,kind=Widget
// +kubebuilder:subresource:status
type Widget struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WidgetSpec   `json:"spec,omitempty"`
	Status WidgetStatus `json:"status,omitempty"`
}

// +die
type WidgetSpec struct {
	// +listType=map
	// +listMapKey=name
	Ports    []WidgetPort           `json:"ports,omitempty"`
	Template corev1.PodTemplateSpec `json:"template"`
}
```

A field whose type is from another package uses that type's die: from the built-in dies, from a die package an import marker names, or from a type marked `+die` in the field type's own package. A field whose type has no die is a value.

### Markers

| Marker                                              | On      | Does                                                                                                                                |
| --------------------------------------------------- | ------- | ----------------------------------------------------------------------------------------------------------------------------------- |
| `+die`                                              | type    | Generates a die for the type                                                                                                        |
| `+die:object:apiVersion=…,kind=…`                   | type    | Generates a die for a top-level object, with its blank's type metadata, scheme registration and a list die                          |
| `+die:types:package=…[,mode=alias\|clone][,from=…]` | package | Brings in another package's types, aliased or cloned                                                                                |
| `+die:applyconfigurations:package=…`                | package | Aliases a package's apply configurations, for the dies' apply configuration methods                                                 |
| `+die:import:package=…,as=…`                        | package | Says where the dies for another package's types are                                                                                 |
| `+die:field:name=…,…`                               | type    | Sets a field's `listType`, `listMapKey`, `mapType`, `structType`, `default` or `nullable`, or renames its accessor with `accessor=` |
| `+die:union:name=…,members={…}`                     | type    | Declares a union, whose members clear the others when set                                                                           |
| `+die:subresource:name=…`                           | type    | Declares a subresource, for `DieSubresource` and `DieMain`                                                                          |

Fields in `+die:field` and `+die:union` are named by their wire names. A type's `+structType=atomic` applies to every field of that type, unless the field says `+structType=granular`.

### Generated tests

The generated `zz_generated.dies_test.go` checks that each die still matches its type. A dependency update that adds fields still compiles, but fails these tests, naming the type and its new fields. The tests also check that each version package registers the same kinds as the upstream package.

## Conformance

[`v2/conformance`](v2/conformance) tests dies against real Kubernetes API servers, in two suites.

**[`apiserver`](v2/conformance/apiserver)** checks the dies themselves. It runs etcd and kube-apiserver in Docker, for the Kubernetes version the dies are generated from, with every API version enabled and no controllers, and connects envtest to them. It checks that:

- every die agrees with the server's OpenAPI schema on how each field merges
- applies own exactly what they set, and merge list maps by key
- merge patches carry only intent
- objects round-trip through dies exactly as the server has them
- unstructured dies in the server's molds apply as typed dies do

It needs only Docker, and is skipped without it, unless `CI` is set:

```sh
cd v2/conformance && go test ./apiserver/...
```

**[`cluster`](v2/conformance/cluster)** checks dies with a real cluster, of any Kubernetes version or distribution. It doesn't care which APIs the cluster serves. It checks that server-side apply and merge patches from dies work as they should, and that unstructured dies find the schemas of the kinds the cluster serves: built-in kinds, custom resources, and an API served through an APIService. If dies behave strangely on your distribution, this is the suite to run against it.

It changes the cluster, so it only runs against a cluster you name. It installs a CRD, which it removes afterwards, and creates a labeled namespace for each test, which it deletes afterwards. That takes a cluster admin's permissions:

```sh
cd v2/conformance && DIES_CLUSTER_CONTEXT=my-cluster go test -count=1 ./cluster/...
```

The tests of an aggregated API need Kubernetes' sample-apiserver, which serves `wardle.example.com` behind an APIService, and are skipped without it. `./hack/conformance-cluster.sh` creates a kind cluster with sample-apiserver installed, runs the suite and deletes the cluster. It keeps the cluster's kubeconfig apart from yours, and `--keep` keeps the cluster.

## Community

### Code of Conduct

The reconciler.io projects follow the [Contributor Covenant Code of Conduct](./CODE_OF_CONDUCT.md). In short, be kind and treat others with respect.

### Communication

General discussion and questions about the project can occur either on the Kubernetes Slack [#reconcilerio](https://kubernetes.slack.com/archives/C07J5G9NDHR) channel, or in the project's [GitHub discussions](https://github.com/orgs/reconcilerio/discussions). Use the channel you find most comfortable.

### Contributing

The reconciler.io dies project team welcomes contributions from the community. A contributor license agreement (CLA) is not required. You own full rights to your contribution and agree to license the work to the community under the Apache License v2.0, via a [Developer Certificate of Origin (DCO)](https://developercertificate.org). For more detailed information, refer to [CONTRIBUTING.md](CONTRIBUTING.md).

To regenerate the dies in this repository, tidy its modules and run its tests, run `./hack/update-codegen.sh`.

## Acknowledgements

`dies` was conceived by [Scott Andrews](https://github.com/scothis) drawing inspiration from previous work in [`reconciler.io/runtime`](https://github.com/reconcilerio/runtime/).

## License

Apache License v2.0: see [LICENSE](./LICENSE) for details.
