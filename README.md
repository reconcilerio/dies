# Dies - immutable, fluent builders for Kubernetes resources <!-- omit in toc -->

![CI](https://github.com/reconcilerio/dies/workflows/CI/badge.svg?branch=main)
[![codecov](https://codecov.io/gh/reconcilerio/dies/branch/main/graph/badge.svg?token=vyXLcPBdV3)](https://codecov.io/gh/reconcilerio/dies)
[![Go Reference](https://pkg.go.dev/badge/reconciler.io/dies/v2.svg)](https://pkg.go.dev/reconciler.io/dies/v2)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](CODE_OF_CONDUCT.md)

Dies build and change Kubernetes resources fluently, in tests and in controllers. This project has dies for every built-in Kubernetes API, and a generator, `diegen`, that creates dies for custom resources.

A die records what you change, not just the resulting value. It knows whether you left a field alone, set it to its zero value, or cleared it. So it can send exactly that as a server-side apply or a JSON merge patch, without sending stale values back to the server or clearing what other writers own.

In a workshop, a die is the tool in a press that gives material its shape. A blank is fed in, stamped, and released as a formed part. The names here follow that:

- **Blank:** every die starts from one, such as `DeploymentBlank`.
- **Feed:** `DieFeed` and the `die.Feed…` functions give a die its material, such as an object, JSON, YAML, or what the server has.
- **Stamp:** `DieStamp`, the field accessors and `die.StampValue` shape a die, and it records each change as intent.
- **Release:** `DieRelease` and the `die.Release…` functions give you the part. `die.ReleaseApplyPatch` and `die.ReleaseMergePatch` release only what was stamped, as a patch.
- **Jig:** a jig holds dies of different shapes the same way, so the same tools work on each.

v2 is a rewrite, and its API is not compatible with v0. The v0 API is documented at [pkg.go.dev/reconciler.io/dies](https://pkg.go.dev/reconciler.io/dies). We skipped v1 so that both versions can be used in a project during the transition.


- [Getting started](#getting-started)
- [Using dies](#using-dies)
	- [Fields](#fields)
	- [Immutable and mutable dies](#immutable-and-mutable-dies)
	- [Feeding and releasing](#feeding-and-releasing)
	- [Errors](#errors)
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
- [Conversion](#conversion)
- [Observers](#observers)
- [JSONPath](#jsonpath)
- [Validation and defaulting](#validation-and-defaulting)
	- [Checks from markers](#checks-from-markers)
	- [Updates](#updates)
	- [Field errors](#field-errors)
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

deployment, err := dieappsv1.DeploymentBlank.DieStampE(func(d *dieappsv1.DeploymentDieMut) error {
	app := d.Spec().Template().Spec().Containers().At("app")
	return errors.Join(
		d.Metadata().Namespace().Set("default"),
		d.Metadata().Name().Set("web"),
		d.Spec().Replicas().Set(2),
		app.Image().Set("registry.example/app:1"),
		app.Env().At("MODE").Value().Set("production"),
		app.Ports().At(diecorev1.ContainerPortKey{ContainerPort: 8080}).Name().Set("http"),
	)
})
// ...

obj := deployment.DieRelease() // a *appsv1.Deployment
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

An immutable die's handles read: `die.Value`, `die.Optional`, `die.ListMap` and `die.Map`. A mutable die's handles, `die.ValueMut`, `die.OptionalMut`, `die.ListMapMut` and `die.MapMut`, also write (see [Immutable and mutable dies](#immutable-and-mutable-dies)). Writes return only an error, which you check (see [Errors](#errors)).

Every leaf handle has `Get` and `Has`; a mutable die's also have `Set` and `Clear`, and lists and maps also have `Empty`:

- `Has` reports whether the field is present, as the server sent it or as you set it. A zero value can be present.
- `Set` records the value, even when it's the zero value.
- `Clear` records that the field is absent.
- `Empty` records a list or map that is present and empty.

A list map's elements are found by their keys. `At` returns the element's die, and its first write creates the element with its keys set. List maps also have `Has`, `Keys`, `Len` and `All`, and a mutable die's have `Remove`. A list map with more than one key, such as container ports keyed by port and protocol, takes a generated key type, such as `ContainerPortKey`. A key with a default, such as a port's protocol, which defaults to `TCP`, matches whether it's set or not.

Lists that aren't list maps, such as a container's `args`, are values, which are read and written whole.

### Immutable and mutable dies

Every die type comes in two: an immutable die, such as `DeploymentDie`, and a mutable die, such as `DeploymentDieMut`. An immutable die only reads. It has no methods that write, and nothing writes to it, so a die you hold never changes underneath you. A mutable die also writes, in place.

Blanks are immutable. `DieStamp` passes a mutable copy of a die to each function, and returns an immutable die of the result, leaving the original unchanged. That lets you capture a die partway through and derive variations from it, which is what makes dies useful for test fixtures. In a test, `dietesting.NoError` checks the writes:

```go
base := dieappsv1.DeploymentBlank.DieStamp(func(d *dieappsv1.DeploymentDieMut) {
	dietesting.NoError(t,
		d.Metadata().Name().Set("web"),
		d.Spec().Template().Spec().Containers().At("app").Image().Set("app:1"),
	)
})

scaled := base.DieStamp(func(d *dieappsv1.DeploymentDieMut) {
	dietesting.NoError(t, d.Spec().Replicas().Set(3))
})
// base is unchanged
```

The children of an immutable die are immutable, and those of a mutable die are mutable. On a mutable die, `DieStamp` passes the die itself, which groups writes to part of it, as the condition in [Server-side apply](#server-side-apply) does. `DieStamp` on a die within an immutable one makes the copy a new root, holding the die's current value.

- `DieMutable()` returns a mutable copy of an immutable die, and `DieImmutable()` an immutable copy of a mutable one. Only an immutable die has `DieMutable`, and only a mutable die `DieImmutable`. Both always copy: an immutable die never sees writes to a mutable die it was copied from, or to one copied from it.
- `die.ImmutableDie` and `die.MutableDie` are an immutable and a mutable die of any type, for functions that take only one kind, such as a jig's `Hold` and `HoldMut`. Passing the other kind doesn't compile.
- `DeepCopy` copies a die, as the same kind of die.
- `die.StampValue(d, func(r *T))` passes a copy of a mutable die's Go value to a function, for code written for the Go type. The die works out what the function changed, and records it as intent.

Both kinds of object die implement `runtime.Object` and `metav1.Object`. Those interfaces' setters can't return an error, so on an immutable die they panic with `die.ErrImmutable`. A scheme registers the mutable dies, which decode, and the items of a list die are mutable dies.

### Feeding and releasing

Dies convert to and from the forms Kubernetes code uses:

Feeding is a write, so the feeds are methods of mutable dies, and return only an error. Every die releases. Feed within `DieStampE` for an immutable die of the value:

```go
d, err := dieappsv1.DeploymentBlank.DieStampE(func(d *dieappsv1.DeploymentDieMut) error {
	return d.DieFeed(deployment)
})
```

A die's own methods feed and release its Go value, by pointer: `DieFeed(*T)`, where nil feeds the zero value, and `DieRelease() *T`, a copy. The other forms are functions in the `die` package, which take the die. They're functions rather than methods so that their code is compiled once, rather than for each type of die:

| Feed                                                                | Release                                                                        |
| ------------------------------------------------------------------- | ------------------------------------------------------------------------------ |
| `d.DieFeed(*T)`                                                     | `d.DieRelease() *T`                                                            |
| `die.FeedJSON(d, []byte)`                                           | `die.ReleaseJSON(d) ([]byte, error)`                                           |
| `die.FeedYAML(d, []byte)`                                           | `die.ReleaseYAML(d) ([]byte, error)`                                           |
| `die.FeedRawExtension(d, runtime.RawExtension)`                     | `die.ReleaseRawExtension(d) (runtime.RawExtension, error)`                     |
| `die.FeedUnstructured(d, *unstructured.Unstructured)`, objects only | `die.ReleaseUnstructured(d) (*unstructured.Unstructured, error)`, objects only |
| `die.FeedApplyPatch(d, []byte)`, root dies only                     | `die.ReleaseApplyPatch(d) ([]byte, error)`, root dies only                     |
| `die.FeedMergePatch(d, []byte)`, root dies only                     | `die.ReleaseMergePatch(d) ([]byte, error)`, root dies only                     |

Feeds take a mutable die and return only an error, as writes do. Releases take any die, and return an error for a value that doesn't encode. Each has a `Must` form that panics instead, such as `die.MustReleaseJSON(d)`. An unstructured object is an object, so only object dies feed and release one; a raw extension carries any other die's value.

When a die decodes JSON or YAML, it records which fields were present and keeps any fields its Go type doesn't have, such as fields added in a newer version of the API. Encoding it again gives back what it read: an update from a die doesn't drop fields it doesn't know about.

Dies are `json.Marshaler`s. Mutable dies are also `json.Unmarshaler`s; an immutable die's `UnmarshalJSON` fails with `die.ErrImmutable`. Dies for objects also implement `runtime.Object` and `metav1.Object`, so Kubernetes libraries and clients take them directly.

For test fixtures, `die.FeedYAMLFile(DeploymentBlank, "testdata/deployment.yaml")` returns an immutable die of a file, and `die.MustFeedYAMLFile` panics rather than returning an error, for fixtures in package variables.

### Errors

A die's reads can't fail. Content is checked where it enters a die: a feed, decode or write that would give a list map two elements with the same key is refused, as is a jig of a source it can't hold. So `Get`, `Has`, `DieRelease` and the other reads always have an answer, and return no error.

Writes return only an error: `Set`, `Clear`, `Empty` and `Remove`, `DieFeed` and the `die.Feed…` functions, `die.StampValue` and `DieClear`. A write fails for content the die can't hold: JSON or YAML that doesn't decode as the die's type, a duplicate list-map key, a value an unstructured die can't convert, or a change a jig's transform can't represent in its source. Check the error, as linters such as errcheck require. `DieStampE` takes functions that return errors, such as a client operation or a block of writes, and returns the first one, along with the die as the functions left it. `errors.Join` checks a block of writes at once:

```go
updated, err := current.DieStampE(func(d *dieappsv1.DeploymentDieMut) error {
	return errors.Join(
		d.Spec().Replicas().Set(3),
		d.Spec().Template().Spec().Containers().At("app").Image().Set(image),
	)
})
```

In a test, `dietesting.NoError(t, errs...)` fails the test with the first error, so fixtures check their writes inline: `dietesting.NoError(t, d.Spec().Replicas().Set(3))`.

Functions that can fail for other reasons return an error too: the release functions, such as `die.ReleaseJSON`, for a value JSON can't encode, and functions that take an expression or content from elsewhere, such as `jsonpath.Find`, `dieunstructured.GetAs`, molds' `Blank` and `Cast`, and a jig's `Clamp`, `Hold`, `HoldMut` and `HoldUnstructured`. Each has a `Must` form that panics with the error instead, for content you know, as Go's `regexp.MustCompile` does: `die.MustReleaseJSON`, `jsonpath.MustFind`, `jig.MustHold`. A die's method is `DieMust…`, such as `DieMustReleaseApplyConfiguration`. Defining a jig, usually in a package variable, uses `MustClamp`, as a mistake in a definition is a mistake in the code. Defining a [conversion](#conversion) has only forms that panic, for the same reason.

Mistakes in the code panic, as Go's own APIs do. They don't depend on the data, so a test catches them the first time it runs:
- **A root's operation on a die within another, or on a jig,** such as releasing or feeding a patch.
- **A scope's location within a die of another type.**
- **A jig location the jig doesn't bind for its source, written,** or an aggregate list written whole, or an element added where it could go to more than one list. `jigs.IsBound` says whether a jig binds a location; a location it doesn't bind reads as absent.
- **An unstructured write that steps into a value of another kind,** such as a field of a string, or past the end of a list, as reflect's methods panic. A read of such a place reads as absent, and `Type`, `Has` and `Len` say what a node holds, so careful code never panics. A key that isn't unique in an unstructured list, which has no schema to refuse it with when fed, panics where `Item` uses it; `Find` with a filter finds every element with a key.
- **Interface implementations,** such as `metav1.Object`'s setters, can't return an error. They panic on an immutable die. Write through a mutable die's own accessors instead.
- **A broken die definition,** whether generated or hand-written, such as an unknown die tag, or a field table that doesn't match its type. Generated dies' tests catch them.

A jig checks its source once, when it's made. A change made to the source directly, not through the jig, that leaves the jig unable to hold it, such as the same container name in two lists an aggregate binds, makes the jig's reads panic: take care when changing the same content through two dies.

Use `errors.As` with `*die.Error` to tell a die's error from another, and `errors.Is` with `die.ErrImmutable` for a write to an immutable die, such as through `metav1.Object`'s setters. A cause, such as a JSON syntax error, is wrapped.

## Patch with intent

A die is a value plus a record of what you changed: its intent. Patches are encoded from the intent, never from the whole Go value, so:

- a field you didn't touch is never sent, even when its Go type would encode it as `null`
- a zero value you set, such as `replicas: 0` or `paused: false`, is sent
- a field you cleared is sent as removed

Reads see the value with the intent applied. Feeding a root die replaces its value and discards its intent. `DieIntent()` describes the intent, one change per line, which helps when debugging.

Patches can be fed too. `die.FeedApplyPatch` and `die.FeedMergePatch` take a patch as the die's intent, over the value it already has, so releasing the same kind of patch sends it back. Use them to inspect a patch from elsewhere through a die's accessors, or to replay one against another value. A merge patch fed to the value it was made against reproduces the change, lists included, which a merge patch sends whole. Fields the die's type doesn't have aren't fed. An unstructured die follows its mold where it has one; without one, its lists are set whole.

### Server-side apply

`die.ReleaseApplyPatch(d)` encodes the intent as a server-side apply request, identified by the object's `apiVersion`, `kind`, `name` and `namespace`. The manager applying it owns exactly the fields it sent:

- `Clear` on a field sends `null` where the server accepts it: the manager owns the field, and it's absent. A server that converts objects through Go types accepts `null` for any field; that covers built-in types, and aggregated APIs such as sample-apiserver's. A custom resource's validation rejects `null` unless the field is `nullable`, so for a custom resource, `Clear` sends `null` only for fields marked `+nullable`, and otherwise leaves the field out, which removes it only if this manager owned it. Typed dies treat types from `k8s.io` as served from Go types. Unstructured dies with the server's schema tell from the schema itself.
- `Empty` sends `[]` or `{}`: the manager owns an empty list or map.
- Removing a list-map element or a map entry leaves it out of the apply. An apply can't remove what another manager owns. The server removes what this manager owned and no longer applies. Sending `null` for a map entry would set it to `""` instead.
- A struct the server replaces whole, such as a `LabelSelector` or an `ObjectReference` (`+structType=atomic`), keeps its child die, but is sent whole.

```go
status, err := dieappsv1.DeploymentBlank.DieStampE(func(d *dieappsv1.DeploymentDieMut) error {
	available := d.Status().Conditions().At(dieappsv1.DeploymentAvailable)
	return errors.Join(
		d.Metadata().Namespace().Set(namespace),
		d.Metadata().Name().Set(name),
		d.Status().ObservedGeneration().Set(generation),
		available.Status().Set(diecorev1.ConditionTrue),
		available.Reason().Set("Reconciled"),
	)
})
// ...

body, err := die.ReleaseApplyPatch(status, die.Subresource("status"))
// ...
// the object as the server has it after the apply, such as its new resource
// version: tared first, so the server's response isn't overlaid with the
// apply's intent
applied, err := status.DieTare().DieStampE(func(d *dieappsv1.DeploymentDieMut) error {
	return c.Status().Patch(ctx, d, client.RawPatch(types.ApplyPatchType, body), client.FieldOwner("my-controller"))
})
```

A client decodes the server's response into the object you pass it, so that object has to be mutable, and you need a reference to it to see the result. `DieStampE` runs a client operation on a mutable copy of an immutable die, and returns the result as a new immutable die, with the operation's error. Passing an immutable die to a client directly fails: it isn't the type the scheme registers, and it can't be decoded into.

Decoding replaces the die's value but keeps its intent, applied on top. To see exactly what the server returned, pass a die without intent: one read from the server, or one tared with `DieTare`. This matters after an apply, because the server doesn't always do what the apply asked. A field the apply left out, for example, stays as another manager set it.

An apply is all of a manager's intent. If a manager applies from server state with one change, it gives up every other field it owned, and the server removes them. Build applies from a blank, as a reconciler does above, or make the change as a different manager.

### Merge patches

`die.ReleaseMergePatch(d)` encodes the intent as a JSON merge patch (RFC 7386) against the value the die was fed. Only what you changed is sent, so concurrent changes to other fields survive. A merge patch can only replace a list whole, so a change to a list sends the whole list as the die has it.

```go
// the current state, read from the server
current, err := dieappsv1.DeploymentBlank.DieStampE(func(d *dieappsv1.DeploymentDieMut) error {
	return c.Get(ctx, key, d)
})
// ...

// the update, a change to the current state
update, err := current.DieStampE(func(d *dieappsv1.DeploymentDieMut) error {
	return errors.Join(
		d.Spec().Replicas().Set(0),
		d.Metadata().Labels().Remove("canary"),
	)
})
// ...

// the patch, from the update's intent against the current state
body, err := die.ReleaseMergePatch(update)
// ...

// the patched state, as the server has it after the patch
patched, err := current.DieStampE(func(d *dieappsv1.DeploymentDieMut) error {
	return c.Patch(ctx, d, client.RawPatch(types.MergePatchType, body))
})
```

The patch is sent with the current die, which has no intent, so the patched die is exactly what the server returned. Each step is a die of its own, and none of them changes another.

### Choosing what to send

A patch is released from all of a die's intent, or from the part that scopes keep. A scope is part of releasing a patch, not of the die, so the die keeps all of its intent:

- `die.Subresource("status")` keeps only the intent for a subresource.
- `die.MainResource()` drops the intent for every subresource.
- `die.Only(locations...)` and `die.Except(locations...)` keep or drop the intent at locations, such as `die.Only(d.Spec().Replicas(), d.Metadata().Labels())`.

```go
status := die.Subresource("status")
body, err := die.ReleaseApplyPatch(d, status)
ac, err := die.ReleaseApplyConfiguration(d, status)
```

`die.ReleaseApplyPatch`, `die.ReleaseMergePatch` and `die.ReleaseApplyConfiguration` take scopes, and a scope can be kept in a variable for each of them. Scopes passed together each narrow what the one before kept. Intent a scope leaves out is folded into the value the patch is made against, identity included, so only the kept intent is sent.

Two methods of the die help too:

- `DieTare()` folds the intent into the value, so the current value is what later writes are measured from, as taring a scale makes what's on it the new zero. Tare after a successful write, or before changes that should be intent of their own.
- `DieFields()` lists the fields a die manages. Each entry is a location too, for a scope.

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

d, err := dieappsv1.DeploymentBlank.DieStampE(func(d *dieappsv1.DeploymentDieMut) error {
	return c.Get(ctx, client.ObjectKey{Namespace: "default", Name: "web"}, d)
})
// ...
replicas := d.Spec().Replicas().Get()
```

A client strips the type metadata of a typed object as it decodes it. Dies don't treat that as intent, and an apply patch fills in the object's `apiVersion` and `kind` from its die type.

### Apply configurations

Dies for types that client-go has apply configurations for have methods to convert between the two:

- `DieReleaseApplyConfiguration()` returns the intent as client-go's apply configuration, for clients and tools that take one, and `DieMustReleaseApplyConfiguration()` panics rather than returning an error.
- `DieFeedApplyConfiguration(ac)` feeds one in as intent, such as what client-go's `Extract` functions find a manager owns.
- `DieStampApplyConfiguration(fn)` changes the intent through one.

Each die package aliases its apply configuration types and `Extract` functions, too. For any root die, with or without an apply configuration type of its own, `die.ReleaseApplyConfiguration(d, scopes...)` adapts its intent for controller-runtime's `client.Apply`, and `Result()` has what the server returned.

As with other client operations, the client decodes the server's response into the apply configuration you pass, so keep a reference to it:

```go
ac, err := die.ReleaseApplyConfiguration(deployment)
// ...
err = c.Apply(ctx, ac, client.FieldOwner("my-controller"))
// ...
// the object as the server has it after the apply
actual, err := dieappsv1.DeploymentBlank.DieStampE(func(d *dieappsv1.DeploymentDieMut) error {
	return die.FeedJSON(d, ac.Result())
})
```

### Lists

A list die, such as `dieappsv1.DeploymentListDie`, is a plain Kubernetes list whose items are mutable dies, as the scheme registers them. Clients list into it, and `meta.ExtractList` and `meta.SetList` work with it.

```go
list := diecorev1.NewConfigMapList()
err := c.List(ctx, list, client.InNamespace("default"))
for _, cm := range list.Items {
	_ = cm.Data().At("color")
}
```

## Unstructured dies

`reconciler.io/dies/v2/die/unstructured` has dies for objects with no Go type, such as custom resources you don't have types for. An unstructured die is a `runtime.Unstructured`, so controller-runtime's client takes it without any registration. It's immutable, `dieunstructured.Object`, or mutable, `dieunstructured.ObjectMut`, records intent and encodes patches like any other die. You traverse it through nodes: an `Object`'s `Node` reads, and an `ObjectMut`'s `NodeMut` also writes.

```go
d, err := dieunstructured.BlankFor("example.dev/v1", "Widget").DieStampE(func(d *dieunstructured.ObjectMut) error {
	d.SetNamespace("default")
	d.SetName("web")
	ports := d.Node().Field("spec").Field("ports")
	return ports.Item("name", "http").Field("port").Set(int64(80))
})
```

- `Field` steps into an object field, `Item` into a list element by its keys, and `Index` into a list element by position.
- `Get`, `Has`, `Keys`, `Len` and `Type` read a node, and `Set`, `Clear`, `EmptyObject` and `EmptyList` write a `NodeMut`.
- `dieunstructured.JigAs[*diecorev1.ContainerDie](node)` gives an immutable typed die for part of an immutable unstructured object, and `JigAsMut[*diecorev1.ContainerDieMut](node)` a mutable one for part of a mutable object, whose writes go to that object. `GetAs` gives a typed copy.
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

blank, err := molds.Blank(schema.GroupVersionKind{Group: "example.dev", Version: "v1", Kind: "Widget"})
// ...
widget, err := blank.DieStampE(func(d *dieunstructured.ObjectMut) error {
	d.SetNamespace("default")
	d.SetName("web")
	return d.Node().Field("spec").Field("match").Field("note").Clear()
})
```

Dies make no requests of their own: `Molds` reads the schemas through the OpenAPI client you give it, such as a discovery client's `OpenAPIV3()`, or client-go's cached client in front of it. It reads each group version's schema once, when it first makes a blank of one of its kinds, and keeps it until the server changes it. `molds.Cast(d)` returns an immutable unstructured die in its kind's mold, and `molds.CastMut(d)` gives a mutable one its mold in place, such as an item of a list a client listed into. `dieunstructured.BlankFor(apiVersion, kind)` makes a blank without a mold.

Caveats:

- **A die's mold is the schema the server served when its blank was made.** Make a blank again to pick up a changed CRD.
- **A CRD's schema is published shortly after the CRD is established,** not with it, so making a blank right after installing a CRD can fail until the schema appears.
- **The die still guesses where the schema doesn't say,** such as below `x-kubernetes-preserve-unknown-fields`.
- **Only apply patches use the mold.** Merge patches don't need it: RFC 7386 merges objects and maps alike, and replaces lists.
- **Defaults aren't used,** so a list-map key left to its default, such as a port's protocol, doesn't match the element with the default set.
- **Subresources aren't known,** so `die.Subresource` and `die.MainResource` have no effect.
- **A client zeroes the object it decodes into, mold included.** `DieStamp` and `DieStampE` keep the mold. A die passed to a client directly loses it, as do the items of a list a client lists into; `For` gives them their mold back.

## Jigs

A jig holds material in place so that a tool can work on it. A jig here holds dies of different types as one type, such as anything with a pod template, so the same code works on all of them. Each binding pairs a location of the jig's type with a location of a source type. Both are navigated from blanks, so a mismatched binding doesn't compile:

```go
var PodSpecableJig = jigs.NewJig(PodSpecableBlank).
	MustClamp(dieappsv1.DeploymentBlank,
		jigs.Bind(PodSpecableBlank.Metadata(), dieappsv1.DeploymentBlank.Metadata()),
		jigs.Bind(PodSpecableBlank.Template(), dieappsv1.DeploymentBlank.Spec().Template()),
	).
	MustClampUnknownKinds(
		jigs.BindPath(PodSpecableBlank.Metadata(), ".metadata"),
		jigs.BindPath(PodSpecableBlank.Template(), ".spec.template"),
	)

ps, ok, err := PodSpecableJig.HoldMut(workload)
// ...
if ok {
	if err := ps.Template().Spec().Containers().At("app").Image().Set("app:2"); err != nil {
		return err
	}
}
```

A jig holds nothing of its own. Its writes go to the source, and the source's patches carry them. `Hold` takes a `die.ImmutableDie` and returns an immutable die of the jig's type, which reads; `HoldMut` takes a `die.MutableDie` and returns a mutable one, which also writes. A read-only jig of a mutable die would see its writes, so `Hold` doesn't take one. Either takes unstructured dies too, and `HoldUnstructured` takes an `*unstructured.Unstructured`, which the jig changes in place. `jigs.BindVia` converts between shapes that differ. See [`v2/examples/jigs/v1`](v2/examples/jigs/v1) for a complete example.

## Conversion

A custom resource with more than one version needs a conversion webhook, which the API server calls to convert objects between versions. With dies, a conversion between two types is a jig: it binds the locations of one version to those of the other. `reconciler.io/dies/v2/die/conversion` builds the conversion from those bindings.

One version is the hub, usually the storage version, and each other version is a spoke. A spoke is converted to and from the hub, so a conversion between two spokes goes through the hub. You define a converter from the hub's blank, with a spoke for each other version:

```go
var Conversion = conversion.New(sprocketv2.SprocketBlank).
	AddSpoke(sprocketv1.SprocketBlank,
		jigs.Bind(sprocketv1.SprocketBlank.Spec().Size(), sprocketv2.SprocketBlank.Spec().Teeth()),
		jigs.BindVia(sprocketv1.SprocketBlank.Spec().Image(), sprocketv2.SprocketBlank.Spec().Image(), ImageReference),
		jigs.Bind(sprocketv1.SprocketBlank.Status().ObservedGeneration(), sprocketv2.SprocketBlank.Status().ObservedGeneration()),
	)
```

One set of bindings converts both ways. `jigs.Bind` binds a value that's the same in both versions, even if it's renamed. `jigs.BindVia` binds a value whose shape changes, through a transform that must be lossless, as a jig's are. The conversion copies a binding only where the source has a value, so a field left unset stays unset, and a field set to its zero value stays set. The metadata is bound as it is, because the API server lets a conversion change only labels and annotations.

A conversion needs every version of the resource, so define it in a central location, like the API group's package. A version's package then needn't load the other versions, and different kinds in the group can have different hub versions without an import cycle.

The bindings needn't cover every field, as a field of one version may have no place in the other. What a conversion drops, it stashes in an annotation of the converted object, named for the version it came from: `conversion.dies.reconciler.io/v1` holds what a conversion from v1 dropped. Converting back restores the stash and removes the annotation, so a round trip changes nothing. A restore fills only fields the conversion left unset, so a stale stash can't overwrite a value changed in the other version. `WithStashAnnotation` sets another prefix for the annotations. The prefix becomes part of the stored data: a stash under an old prefix isn't restored after you change it.

A stash keeps whole values. If a list's elements lose a field in the other version, the whole list is stashed. The stash isn't restored over the list the conversion back writes, so that field is lost. See [`v2/examples/sprocket`](v2/examples/sprocket) for a complete example, with a renamed field, a field whose shape changes, and a field in each version that the other doesn't have.

Each version's group, version and kind are those its die package records from the type's `+die:object` marker, so only object dies can be versions. A definition that can't convert panics: a spoke of another kind, a second conversion of the same version, or a binding to another type. Each of these is a mistake in the code. `AddSpoke` and `WithStashAnnotation` return a new converter, leaving the one they're called on unchanged. `Convert` returns an error for problems at run time: an object of another kind, a version the converter doesn't know, or a value a transform can't represent.

`reconciler.io/dies/v2/die/conversion/webhook` serves a converter as the webhook:

```go
mux.Handle("/convert/sprockets", webhook.Handler(sprocket.Conversion))
```

Register it in the CRD with the `Webhook` conversion strategy, `conversionReviewVersions: [v1]`, and the CA that signs the webhook's serving certificate.

## Observers

An observer is a function a die calls after each change to it. `DieObserve` observes the die itself. The function gets the die, mutable, as the change left it, so it can write in response, and a `die.Change` holding immutable dies of it before and after. They make it explicit that a change is only read. They read as dies do, so they tell an unset field from one set to zero, and each is within its root as it was then.

```go
err := deployment.Spec().DieObserve(func(spec *dieappsv1.DeploymentSpecDieMut, c die.Change[*dieappsv1.DeploymentSpecDie]) error {
	replicas := c.After.Replicas()
	return spec.Paused().Set(replicas.Has() && replicas.Get() == 0)
})
```

`die.ObserveAt` observes a location within a mutable die, which can be a field or a collection as well as a die. Its second argument returns that location from the die, and the observer is called with a mutable handle to it. The handle's `DieRoot` is the changed die, so the observer can write anywhere in it. The change holds the location's immutable form before and after, as its `DieReadOnly` returns it: an immutable die, or a handle such as a `die.Value` or `die.ListMap`, whose `Has` says whether it's present. The callback's types follow from the location's, so an editor can complete them. A function that takes the handle and the change can observe dies of any type:

```go
err := die.ObserveAt(deployment, func(m *dieappsv1.DeploymentDieMut) die.ValueMut[int32] {
	return m.Spec().MinReadySeconds()
}, func(v die.ValueMut[int32], c die.Change[die.Value[int32]]) error {
	if !c.After.Has() {
		return v.Set(10) // cleared, not set to zero
	}
	return nil
})
```

A mutable list map's `DieObserveEach` observes each of its elements. The observer is called once for each element a change added, removed or changed, with the element's die and its own change, so elements come matched by key, and elements that didn't change aren't called. Its change holds immutable dies of the element: for an added element, `Before.DieHas()` is false, and for a removed one, `After.DieHas()` is false and the die it's called with reads as absent. Register it on a mutable die, or for an immutable die, within a stamp of its root.

The condition sets in [`v2/examples/conditions`](v2/examples/conditions) keep a happy condition, such as Ready, the aggregate of the conditions it depends on, observing the whole list. `conditions.LastTransitionTime` observes each condition, and updates its `lastTransitionTime` when it changes:

```go
ready := conditions.NewConditionSet("Ready", "Configured", "Deployed")
widgetConditions := func(m *widgetv1.WidgetDieMut) *conditions.Conditions { return m.Status().Conditions() }

blank, err := widgetv1.WidgetBlank.DieStampE(func(m *widgetv1.WidgetDieMut) error {
	return errors.Join(
		die.ObserveAt(m, widgetConditions, ready.Observe),
		m.Status().Conditions().DieObserveEach(conditions.LastTransitionTime(metav1.Now)),
	)
})

widget := blank.DieStamp(func(d *widgetv1.WidgetDieMut) {
	// Ready becomes Unknown, waiting for Deployed
	err := die.StampValue(d.Status().Conditions().At("Configured"), func(c *metav1.Condition) {
		c.Status, c.Reason = metav1.ConditionTrue, "Applied"
	})
	// ...
})
```

- **A change is compared, not just written.** After each write, each observer whose location the write may have reached reads it again. The observer is called only if the value differs from what it last saw. The change's `Before` is what it last saw, so writes it didn't see individually arrive together.
- **Observers cascade.** An observer's writes notify the die's other observers in turn. With condition sets where A and B make C, and C and D make E, marking B updates C and then E. Observers run in the order they're registered, and the result shouldn't depend on that order.
- **Observers settle.** An observer isn't called again for its own writes while it runs. Once it returns, it's called again if its location changed since it was called, until the location stops changing. An observer that changes its location every time, or observers that undo each other, make the write fail.
- **Derived dies keep observers.** A die's observers are copied to its stamps, its mutable and immutable copies, deep copies and tares. Observing an immutable die returns a new die, and its observers run within the dies derived from it, as in `DieStamp` above. A copy of a die within another keeps the observers registered at or below it.
- **A part of an immutable die is observed through its root.** `DieObserve` on an immutable die within another, or on an immutable jig, panics: it would observe a detached copy. Observe the part within a stamp of its root: `root.DieStamp(func(m M) { m.Status().DieObserve(fn) })`. `die.ObserveAt` and `DieObserveEach` take only mutable dies, so they're registered on an immutable die the same way. A mutable die within another shares its root's state, so it can be observed directly.
- **A jig's observers see writes through the jig,** not writes made to its source directly.

Observers are notified as each write is made, so they see a change made in several writes part by part. Writing a condition's status and then its reason runs them twice. `die.StampValue` writes several fields as one change. `DieStampAtomic` holds notifications until each of its functions returns, so each function's writes are one change. Within a function, what the observers would write isn't there yet.

`DieStampAtomic` is used for rollback. Nothing is undone implicitly: a write whose observer fails is kept, along with what its observers wrote, and returns the observer's error. Within `DieStampAtomic`, a function that returns `nil` commits its writes. A function that returns an error, or whose observers fail, rolls back the whole stamp: every function's writes, what the observers wrote and saw, and any observer registered within it. It returns the error. What an observer did outside the die isn't undone. A stamp nested in another's function is held until that function returns, and if it fails, it rolls back only its own changes. A jig can't be rolled back, because its writes go to its source: stamp the source instead.

```go
updated, err := current.DieStampAtomic(func(d *widgetv1.WidgetDieMut) error {
	c := d.Status().Conditions().At("Deployed")
	return errors.Join(c.Status().Set(metav1.ConditionFalse), c.Reason().Set("ImagePullBackOff"))
})
```

An atomic stamp covers every write to the die while it runs, through any reference to it, not only the die it passes to its functions. A mutable die isn't safe for concurrent use, with or without observers, so write it from one goroutine at a time.

## JSONPath

`reconciler.io/dies/v2/die/jsonpath` finds locations within a die with kubectl's JSONPath expressions. A match is a live location within the die, not a copy, so writes through a match are recorded as intent. `$` is the root die, and `.` is the die you query, so a function written for a `PodSpecDie` can also reach the object the pod spec is part of:

```go
images, err := jsonpath.FindValues[string](podSpec, ".containers[*].image")
namespace, err := jsonpath.FindValues[string](podSpec, "$.metadata.namespace")
containers, err := jsonpath.FindDies[*diecorev1.ContainerDie](podSpec, ".containers[?(@.name=='app')]")
```

`FindValues` returns handles that read, and `FindDies` immutable dies of an immutable die. For a mutable die, `FindValuesMut` returns handles that also write, and `FindDiesMut` mutable dies, such as `*diecorev1.ContainerDieMut`.

A list-map element is found by its keys, however the expression found it, so a write through it merges by key. A path with no wildcards or filters resolves even when its value isn't set, so a write through it creates it. A field the type doesn't have is an error, not a silent miss.

## Validation and defaulting

A die validates or defaults only its own fields. Implement `DieValidate` on the immutable dies, or `DieDefault` on the mutable dies, of the types that need logic of their own, and `die.ValidateCreate` or `die.Default` walks every die within the one you pass it, calling each die that has one. The dies in between need nothing, so a nested type's rules apply wherever it's used, without plumbing them through each type above it:

```go
// in the die package, beside the generated dies
func (d *WidgetPortDie) DieValidate(ctx context.Context, op die.Operation) (validation.ErrorList, error) {
	if d.Port().Get() == 22 {
		return validation.ErrorList{validation.Invalid(d.Port(), "is reserved for ssh")}, nil
	}
	return nil, nil
}

func (d *WidgetPortDieMut) DieDefault(ctx context.Context) error {
	if d.Port().Get() == 0 {
		return d.Port().Set(80)
	}
	return nil
}
```

```go
errs, err := die.ValidateCreate(ctx, widget) // a field.ErrorList
// spec.ports[1].port: Invalid value: 22: is reserved for ssh

err := die.Default(ctx, mutableWidget)
```

The walk reaches each struct field held by value, or by pointer when it's present, each element of a list map, and each item of an atomic list whose type has a die. Errors are located from the root, as the die knows where it is: a list-map element is at its index, as the API server reports it. A die within another reports the paths its root has, so validating a pod spec within a deployment reports `spec.template.spec.…`.

`die.ValidateCreate`, `die.ValidateUpdate` and `die.ValidateDelete` each pass a `die.Operation` saying what's being validated, as an admission request does. Validation returns apimachinery's `field.ErrorList`, which a webhook or the API server reports as it is.

`DieValidate` returns the errors of the value, and an error of its own when it can't validate the value at all, such as when a lookup it needs fails. That error stops validation: the `die.Validate…` functions return it, with the location of the die that returned it, along with the value's errors found until then.

Defaulting is top down: each die is defaulted before the dies within it are walked, so what a die's defaults create is defaulted too. It stops at the first error, with the location of the die that returned it.

### Checks from markers

The markers diegen reads are checked without any code: those of Kubernetes' declarative validation, of kubebuilder, and of the OpenAPI schema. The checks are apimachinery's declarative validation, so their errors are the API server's:

- **Required fields:** `+required` and `+kubebuilder:validation:Required` need the field present; `+k8s:required` also needs it not empty, unless it's held by pointer. A list-map element's key without a default can't be empty.
- **Values:** `+k8s:format` (`k8s-short-name`, `k8s-long-name`, `k8s-label-key` and the rest), `+k8s:minimum` and `+k8s:maximum`, `+k8s:maxLength`, `+k8s:maxBytes`, `+k8s:maxItems`, and `+k8s:eachVal` for a list's items or a map's values; kubebuilder's `Minimum`, `Maximum`, `ExclusiveMinimum`, `ExclusiveMaximum`, `MultipleOf`, `MinLength`, `MaxLength`, `Pattern`, `Enum`, `Format` (OpenAPI's formats, such as `date-time`), `MinItems`, `MaxItems`, `UniqueItems`, `MinProperties`, `MaxProperties` and `items:…`.
- **Enums:** a field whose type is marked `+enum` or `+k8s:enum` may only be one of its type's constants, such as a container port's protocol.
- **Sets:** a `+listType=set` or `+k8s:unique=set` list may not have an item twice.
- **Unions:** exactly one member of a `+oneOf` or `+k8s:unionMember` union must be set, and at most one of a `+zeroOrOneOf` union, or a `+die:union` unless it's `oneOf=true`.

A marker on a type holds for every field of that type, or of a list of it.

`die.Default` sets each absent field that has a `+default` or `+kubebuilder:default` to it, as the API server does for a custom resource, before calling the die's `DieDefault`.

A package of custom resources has kubebuilder's convention: a field without `omitempty` or `omitzero` is required, unless it's `+optional`, or the package is `+kubebuilder:validation:Optional`. diegen recognizes such a package by the markers controller-gen generates CRDs from, such as `+kubebuilder:object:root` and `+kubebuilder:subresource`. Markers that only validate don't count, so the few in `k8s.io/api` and apimachinery's `meta/v1`, for the custom resources that use their types, don't make those packages custom resources.

Not checked: CEL rules (`+kubebuilder:validation:XValidation`), feature-gated markers (`+k8s:alpha`, `+k8s:ifEnabled`), and unions with a discriminator (`+union` and `+unionDiscriminator`).

### Updates

`die.ValidateUpdate(ctx, d, old)` validates a die that replaces an old one of the same type, such as in an admission webhook of the update, or before sending one. Each die it walks to is paired with the old die at the same location: a field by its name, a list-map element by its key, and an atomic list's item by its index. `die.OldOf(d, op)` returns it as the die's own type:

```go
func (d *WidgetPortDie) DieValidate(ctx context.Context, op die.Operation) (validation.ErrorList, error) {
	if old, ok := die.OldOf(d, op); ok && old.Port().Get() != d.Port().Get() {
		return validation.ErrorList{validation.Forbidden(d.Port(), "a named port's number may not change")}, nil
	}
	return nil, nil
}
```

The markers' checks ratchet, as the API server's do: a field unchanged from the old value isn't checked again, so values that were valid when they were written don't block an unrelated change. The markers of an update are checked too: `+k8s:immutable`, and `+k8s:update`'s `NoSet`, `NoUnset`, `NoModify`, `NoAddItem` and `NoRemoveItem`.

`die.ValidateDelete` calls only the dies' `DieValidate`, with `op.Type` set to `die.Delete`: a value being deleted has no schema to meet.

### Field errors

`reconciler.io/dies/v2/die/validation` makes apimachinery's field errors at a location, a die or a handle, rather than at a path you spell out: `Required`, `Forbidden`, `Invalid`, `NotSupported`, `Duplicate`, `NotFound`, `TooLong`, `TooMany`, `TooFew` and `InternalError`. Those that report a value read it at the location, so `validation.Invalid(d.Port(), "…")` reports the port's number.

For apimachinery's other validation helpers, `validation.Path(location)` converts a location to a `*field.Path`. A list-map element is its index in the list; one the list doesn't have, or one of a bare `die.Path`, is its keys, such as `containers[name=app]`.

## Die methods

Every die has these methods, but those that write, which only mutable dies have. `die.Die` is any die, immutable or mutable, whatever its type, for functions that take one, and `die.ImmutableDie` and `die.MutableDie` one of either kind; `die.DieTyped` is any die with type metadata, which is also a `runtime.Object`; `die.DieObject` is any object die, which is also a `metav1.Object`; and `die.DieList` is any list die. Generic code constrains them by `die.DieOf[T, D, M]` and `die.DieMutOf[T, D, M]`, `die.DieTypedOf[T, D, M]` and `die.DieTypedMutOf[T, D, M]` (with the `APIVersion()` and `Kind()` accessors), `die.DieObjectOf[T, D, M]` and `die.DieObjectMutOf[T, D, M]` (with `Metadata()`), and `die.DieListOf[D]` (with `Metadata()`, the list's own `*metav1.ListMeta`, and `DieItems()`), which carry their types: `T` the value's, `D` the immutable die's and `M` the mutable die's. An unstructured die is a `die.DieObject`, but has no typed accessors. Fields' accessors use the field names, so die methods start with `Die`, apart from `DeepCopy`, `Equal` and the methods of the Kubernetes interfaces that object dies implement. Writes return only an error, and reads can't fail; see [Errors](#errors).

| Method                                                          | Does                                                                                                                                 |
| --------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| `DieStamp(fns...)`                                              | Passes a mutable die to each function: a copy of an immutable die, returning an immutable die of the result, or a mutable die itself |
| `DieStampE(fn)`                                                 | Passes a mutable die to a function that can fail, such as a client operation, and returns the result with its error                  |
| `DieMutable()`, immutable only / `DieImmutable()`, mutable only | Return a mutable copy of an immutable die, or an immutable copy of a mutable one                                                     |
| `DieFeed(*T)`, mutable only / `DieRelease()`                    | Convert to and from Go values; the `die.Feed…` and `die.Release…` functions convert to and from JSON, YAML, patches and the rest     |
| `DieHas()` / `DieClear()`, mutable only                         | Report and change whether the die's value is present                                                                                 |
| `DieTare()`                                                     | Fold the intent into the value                                                                                                       |
| `DieIntent()`                                                   | Describe the intent                                                                                                                  |
| `DieFields()`                                                   | List the fields the die manages                                                                                                      |
| `DiePath()` / `DieRootType()`                                   | The die's location within its root, and the type of the root's value                                                                 |
| `DieRoot()`                                                     | The root die the die is within, as a `die.Die`, sharing its state, such as `d.DieRoot().(metav1.Object)` for an object's identity    |
| `DeepCopy()` / `Equal(d)`                                       | Copy a die, or compare current values                                                                                                |

On an immutable die, the methods that return a die, such as `DieTare`, return a new one. On a mutable die, they change it in place.

Object dies implement `runtime.Object` and `metav1.Object`, and have `APIVersion()`, `Kind()` and `Metadata()` accessors.

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

A die package can be one that the common deepcopy generators, controller-gen's `object` and Kubernetes' deepcopy-gen, generate deepcopies for, such as one of your own types with `+kubebuilder:object:generate=true` or `+k8s:deepcopy-gen=package`. They generate deepcopies of the package's types, but not of the dies, which copy themselves. `zz_generated.dies.go` has the build constraint `//go:build !ignore_autogenerated`, as these generators' own files do, so they don't load it, and each type diegen declares there is marked `+kubebuilder:object:generate=false` and `+k8s:deepcopy-gen=false`, for a generator that does. Types in `zz_generated.dies.types.go`, such as cloned types, are generated for as any others are.

### Dies for types from another package

A die package for another package's types needs only a `doc.go`. The `+die:types` marker says where the types come from:

```go
// +die:types:package=github.com/example/widgets/api/v1
package v1
```

- **Aliased** (`mode=alias`, the default): the die package aliases the upstream package's types and constants, so one import gives you both.
- **Cloned** (`mode=clone`): the die package copies the types, with their docs, markers and constants, so it doesn't depend on the upstream package or its dependencies. The package is resolved at a version given with it, `package=…@v1.2.3`, or in the module a `from=` path names.

Each type that needs a die is marked as the upstream package marks it. The markers that change a die's API or its patches, such as `+listType`, `+listMapKey`, `+mapType`, `+structType`, `+default`, `+nullable` and kubebuilder's `+kubebuilder:subresource:status`, are read from the upstream source, in their plain, `k8s:` and `kubebuilder:` forms, as are the markers that validate and default it (see [Checks from markers](#checks-from-markers)). What upstream doesn't say is added with metadata declared on the blank identifier, which adds markers to a type without declaring anything:

```go
// +die:field:name=selector,structType=atomic
type _ = v1.WidgetSpec
```

If the upstream package provides apply configurations, such as from applyconfiguration-gen, `+die:applyconfigurations:package=…` aliases them and adds the apply configuration methods to its dies. For the built-in APIs, client-go's apply configurations are used.

### Dies for your own types

Mark your own types directly. A type's markers can be in its doc comment, or, to keep them out of its godoc, in a comment of their own just above it, separated from the doc comment by a blank line:

```go
// +die:object:apiVersion=example.dev/v1,kind=Widget
// +kubebuilder:subresource:status

// Widget runs a pod template with named ports.
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

Only the comment just above the doc comment is read, as controller-gen reads it, so diegen stops with an error if a `+die` marker is in any other comment above a type.

An object without a `kind`, `+die:object` alone or with only an `apiVersion`, such as a duck type, has an object die, but isn't registered in a scheme: it's a view of other kinds, not a resource of its own.

A field whose type is from another package uses that type's die: from the built-in dies, from a die package an import marker names, or from a type marked `+die` in the field type's own package. A field whose type has no die is a value.

### Markers

| Marker                                              | On      | Does                                                                                                                                                                                                                                                                                         |
| --------------------------------------------------- | ------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `+die`                                              | type    | Generates a die for the type                                                                                                                                                                                                                                                                 |
| `+die:object:apiVersion=…,kind=…`                   | type    | Generates a die for a top-level object, with its blank's type metadata, scheme registration and a list die; without a `kind`, it isn't registered                                                                                                                                            |
| `+die:types:package=…[,mode=alias\|clone][,from=…]` | package | Brings in another package's types, aliased or cloned                                                                                                                                                                                                                                         |
| `+die:applyconfigurations:package=…`                | package | Aliases a package's apply configurations, for the dies' apply configuration methods                                                                                                                                                                                                          |
| `+die:import:package=…,as=…`                        | package | Says where the dies for another package's types are                                                                                                                                                                                                                                          |
| `+die:field:name=…,…`                               | type    | Sets a field's `listType`, `listMapKey`, `mapType`, `structType`, `default` or `nullable`, renames its accessor with `accessor=`, or sets how it's validated: `required=true` or `nonzero`, or a check such as `minimum=0`, `format=k8s-short-name`, `immutable=true` or `each.maxLength=63` |
| `+die:union:name=…,members={…}[,oneOf=true]`        | type    | Declares a union, whose members clear the others when set; at most one may be set, or with `oneOf=true` exactly one                                                                                                                                                                          |
| `+die:subresource:name=…`                           | type    | Declares a subresource, for `die.Subresource` and `die.MainResource`                                                                                                                                                                                                                         |
| `+die:listMap:keys=…`                               | type    | Names the list maps of the type, keyed by those fields, for a type listed by another package, as `metav1.Condition` is: `ConditionListMap` and `ConditionListMapMut`                                                                                                                         |
| `+die:shared`                                       | method  | On a hand-written method of an immutable die, generates the same method on its mutable die, reading through `DieReadOnly`                                                                                                                                                                    |

A die package names the list maps of the element types its list-map fields list, such as `ContainerListMap` and `ContainerListMapMut`, which its accessors return. An accessor of another package's elements, keyed by one field, returns that package's, such as `*diemetav1.ConditionListMapMut`. An element type listed by more than one set of keys isn't named.

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
- a conversion webhook built from dies converts objects the server stores in one version and serves in another, through create, update and server-side apply. The test process serves the webhook, and the API server's container reaches it through a port testcontainers forwards from the host, over TLS with a CA made for the run

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
