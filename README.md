# Dies - immutable, fluent, builders for Kubernetes resources <!-- omit in toc -->

![CI](https://github.com/reconcilerio/dies/workflows/CI/badge.svg?branch=main)
[![codecov](https://codecov.io/gh/reconcilerio/dies/branch/main/graph/badge.svg?token=vyXLcPBdV3)](https://codecov.io/gh/reconcilerio/dies)
[![Go Reference](https://pkg.go.dev/badge/reconciler.io/dies.svg)](https://pkg.go.dev/reconciler.io/dies)
[![Go Report Card](https://goreportcard.com/badge/reconciler.io/dies)](https://goreportcard.com/report/reconciler.io/dies)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](CODE_OF_CONDUCT.md)

- [Using dies](#using-dies)
  - [Common methods](#common-methods)
  - [Helper methods](#helper-methods)
    - [k8s.io/apimachinery/pkg/util/intstr#IntOrString](#k8sioapimachinerypkgutilintstrintorstring)
    - [k8s.io/apimachinery/pkg/api/resource#Quantity](#k8sioapimachinerypkgapiresourcequantity)
    - [k8s.io/api/core/v1#ResourceList](#k8sioapicorev1resourcelist)
- [Creating dies](#creating-dies)
  - [diegen](#diegen)
  - [die markers](#die-markers)
    - [+die](#die)
    - [+die:field](#diefield)
- [Community](#community)
  - [Code of Conduct](#code-of-conduct)
  - [Communication](#communication)
  - [Contributing](#contributing)
- [Acknowledgements](#acknowledgements)
- [License](#license)

---

This project contains dies for many of the most common built-in Kubernetes types, and tools to create dies for custom resources.

## Using dies

Dies start with a blank object that is stamped to add state. All dies are immutable, each stamping returns a new instance containing the mutated state while the original instance is not modified.

```go
import (
    appsv1 "k8s.io/api/apps/v1"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    dieappsv1 "reconciler.io/dies/apis/apps/v1"
    diecorev1 "reconciler.io/dies/apis/core/v1"
    diemetav1 "reconciler.io/dies/apis/meta/v1"
)
```

```go
die := dieappsv1.DeploymentBlank.
    MetadataDie(func(d *diemetav1.ObjectMetaDie) {
        d.Name("my-name")
    }).
    SpecDie(func(d *dieappsv1.DeploymentSpecDie) {
        d.TemplateDie(func(d *diecorev1.PodTemplateSpecDie) {
            d.SpecDie(func(d *diecorev1.PodSpecDie) {
                d.ContainerDie("app", func(d *diecorev1.ContainerDie) {
                    d.Image("registry.example/image:latest")
                    d.EnvDie("MY_VAR", func(d *diecorev1.EnvVarDie) {
                        d.Value("my-value")
                    })
                })
            })
        })
    })
deployment := die.DieRelease()
```

Is equivalent to:

```go
deployment := &appsv1.Deployment{
    ObjectMeta: metav1.ObjectMeta{
        Name: "my-name",
    },
    Spec: appsv1.DeploymentSpec{
        Template: corev1.PodTemplateSpec{
            Spec: corev1.PodSpec{
                Containers: []corev1.Container{
                    {
                        Name:  "app",
                        Image: "registry.example/image:latest",
                        Env:   []corev1.EnvVar{
                            {
                                Name: "MY_VAR",
                                Value: "my-value",
                            },
                        },
                    },
                },
            },
        },
    },
}
```

Unlike when defining an instance of a struct, the die can be captured in an intermediate state with mutations applied incrementally. This behavior is particularly useful for tests that use a common object base with individual mutations applied.

```go
altDeployment := die.
    SpecDie(func(d *dieappsv1.DeploymentSpecDie) {
        d.TemplateDie(func(d *diecorev1.PodTemplateSpecDie) {
            d.SpecDie(func(d *diecorev1.PodSpecDie) {
                d.ContainerDie("app", func(d *diecorev1.ContainerDie) {
                    d.EnvDie("MY_VAR", func(d *diecorev1.EnvVarDie) {
                        d.Value("some-other-value")
                    })
                })
            })
        })
    }).DieRelease())
```

While the outer die is immutable, returning a new die for each call that can be chained together. The dies passed to callbacks are mutable.

Additional methods will be added to dies over time to make common operations easier and safer.

### Common methods

```go
// for managed type `MyResource`

// MyResourceBlank is an empty die that mutations can be stamped from. All die blanks
// are immutable.
var MyResourceBlank *MyResourceDie

type MyResourceDie interface {
    // DieStamp returns a new die with the resource passed to the callback
    // function. The resource is mutable.
    DieStamp(fn func(r *MyResource)) *MyResourceDie

    // Experimental: DieStampAt uses a JSON path (http://goessner.net/articles/JsonPath/)
    // expression to stamp portions of the resource. The callback is invoked with each
    // JSON path match. Panics if the callback function does not accept a single argument
    // of the same type or a pointer to that type as found on the resource at the target
    // location.
    //
    // Future iterations will improve type coercion from the resource to the callback
    // argument.
    DieStampAt(jp string, fn interface{}) *MyResourceDie

    // DieFeed returns a new die with the provided resource.
    DieFeed(r MyResource) *MyResourceDie

    // DieFeedPtr returns a new die with the provided resource pointer. If the
    // resource is nil, the empty value is used instead.
    DieFeedPtr(r *MyResource) *MyResourceDie

    // DieFeedJSON returns a new die with the provided JSON. Panics on error.
    DieFeedJSON(j []byte) *MyResourceDie

    // DieFeedYAML returns a new die with the provided YAML. Panics on error.
    DieFeedYAML(y []byte) *MyResourceDie

    // DieFeedYAMLFile returns a new die loading YAML from a file path. Panics on error.
    DieFeedYAMLFile(name string) *MyResourceDie

    // DieFeedRawExtension returns a new die with the provided raw extension. Panics on error.
    DieFeedRawExtension(raw runtime.RawExtension) *MyResourceDie

    // DieRelease returns the resource managed by the die.
    DieRelease() MyResource

    // DieReleasePtr returns a pointer to the resource managed by the die.
    DieReleasePtr() *MyResource

    // DieReleaseJSON returns the resource managed by the die as JSON. Panics on error.
    DieReleaseJSON() []byte

    // DieReleaseYAML returns the resource managed by the die as YAML. Panics on error.
    DieReleaseYAML() []byte

    // DieReleaseRawExtension returns the resource managed by the die as an
    // raw extension. Panics on error.
    DieReleaseRawExtension() runtime.RawExtension

    // DieImmutable returns a new die for the current die's state that is
    // either mutable (`false`) or immutable (`true`). 
    DieImmutable(immutable bool) *MyResourceDie

    // DieWith returns a new die after passing the current die to the callback
    // function. The passed die is mutable.
    DieWith(fn ...func(d *MyResourceDie)) *MyResourceDie

    // DieSeal returns a new die for the current die's state that is sealed for comparison in future diff and patch operations.
    DieSeal() *MyResourceDie

    // DieSealFeed returns a new die for the current die's state that uses a specific resource for comparison in future diff and patch operations.
    DieSealFeed(r MyResource) *MyResourceDie

    // DieSealFeedPtr returns a new die for the current die's state that uses a specific resource for comparison in future diff and patch operations.
    DieSealFeedPtr(r *MyResource) *MyResourceDie

    // DieSealRelease returns the sealed resource managed by the die.
    DieSealRelease() MyResource

    // DieSealReleasePtr returns the sealed resource pointer managed by the die.
    DieSealReleasePtr() *MyResource

    // DieDiff uses cmp.Diff to compare the current value of the die with the sealed value.
    DieDiff(opts ...cmp.Option) string

    // DiePatch generates a patch between the current value of the die and the sealed value.
    DiePatch(patchType types.PatchType) ([]byte, error)

    // DeepCopy returns a new die with equivalent state. Useful for
    // snapshotting a mutable die.
    DeepCopy() *MyResourceDie
}
```

For each exported field `MyField` on `MyResource`, a method is registered to set that field.

```go
type MyResourceDie interface {
    // continued

    MyField(*MyResource) *MyResourceDie
}
```

Dies marked as implementing `metav1.Object` and `runtime.Object`  generate
additional methods.

```go
import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    runtime "k8s.io/apimachinery/pkg/runtime"
)

type MyResourceDie interface {
    // continued

    runtime.Object
    metav1.Object
    metav1.ObjectMetaAccessor

    // DieReleaseUnstructured returns the resource managed by the die as an
    // unstructured object.
    DieReleaseUnstructured() *unstructured.Unstructured

    // DieDefaultTypeMetadata sets the APIVersion and Kind
    DieDefaultTypeMetadata() *MyResourceDie

    // APIVersion defines the versioned schema of this representation of an
    // object. Servers should convert recognized schemas to the latest internal
    // value, and may reject unrecognized values. More info:
    // https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
    APIVersion(v string) *MyResourceDie

    // Kind is a string value representing the REST resource this object
    // represents. Servers may infer this from the endpoint the client submits
    // requests to. Cannot be updated. In CamelCase. More info:
    // https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
    Kind(v string) *MyResourceDie

    // TypeMetadata standard object's type metadata.
    TypeMetadata(v metav1.TypeMeta) *MyResourceDie

    // TypeMetaDie stamps the resource's TypeMeta field with a mutable die.
    TypeMetaDie(fn func(d *diemetav1.TypeMetaDie)) *MyResourceDie

    // Metadata standard object's metadata.
    Metadata(v metav1.ObjectMeta) *MyResourceDie

    // MetadataDie stamps the resource's ObjectMeta field with a mutable die.
    MetadataDie(fn func(d *diemetav1.ObjectMetaDie)) *MyResourceDie

    // SpecDie stamps the resource's spec field with a mutable die. This method
    // is only created if `MyResourceSpecDie` is defined.
    SpecDie(fn func(d *MyResourceSpecDie)) *MyResourceDie

    // StatusDie stamps the resource's status field with a mutable die. This
    // method is only created if `MyResourceStatusDie` is defined.
    StatusDie(fn func(d *MyResourceStatusDie)) *MyResourceDie
}
```

### Helper methods

Each exported field on the resource backing the die has a setter method defined. Some types offer additional helper methods to make common conversions easier.

####  k8s.io/apimachinery/pkg/util/intstr#IntOrString

```go
type MyResourceDie interface {
    // continued

    Port(port intstr.IntOrString) *MyResourceDie

    // PortInt sets Port with the int value.
    PortInt(port int) *MyResourceDie

    // PortString sets Port with the string value.
    PortString(port string) *MyResourceDie
}
```

#### k8s.io/apimachinery/pkg/api/resource#Quantity

```go
type MyResourceDie interface {
    // continued

    Memory(memory resource.Quantity) *MyResourceDie

    // MemoryString sets Memory by parsing the string as a Quantity. Panics if
    // the string is not parsable.
    MemoryString(memory string) *MyResourceDie
}
```

#### k8s.io/api/core/v1#ResourceList

```go
type MyResourceDie interface {
    // continued

    Limits(limits corev1.ResourceList) *MyResourceDie

    // AddLimit sets a single quantity on the Limits resource list.
    AddLimit(name string, limit resource.Quantity) *MyResourceDie

    // AddLimitString parses the quantity setting a single value on the Limits
    // resource list. Panics if the string is not parsable.
    AddLimitString(name string, limit string) *MyResourceDie
}

```

## Creating dies

Dies are primarily generated for types from [die markers](#die-markers) using
[`diegen`](#diegen).

Additional methods can be added to a die to enrich its behavior. Each additional method is added to the die struct.

Example dispatching to a nested die to managed the template
`corev1.PodTemplateSpec` for `DeploymentSpec`:

```go
func (d *DeploymentSpecDie) TemplateDie(fn func(d *diecorev1.PodTemplateSpecDie)) *DeploymentSpecDie {
    return d.DieStamp(func(r *appsv1.DeploymentSpec) {
        d := diecorev1.PodTemplateSpecBlank.
            DieImmutable(false).
            DieFeed(r.Template)
        fn(d)
        r.Template = d.DieRelease()
    })
}
```

Example adapting `metav1.Condition` dies to `appsv1.DeploymentCondition`:

```go
func (d *DeploymentStatusDie) ConditionsDie(conditions ...*diemetav1.ConditionDie) *DeploymentStatusDie {
    return d.DieStamp(func(r *appsv1.DeploymentStatus) {
        r.Conditions = make([]appsv1.DeploymentCondition, len(conditions))
        for i := range conditions {
            c := conditions[i].DieRelease()
            // coerce metav1.Condition to appsv1.DeploymentCondition
            r.Conditions[i] = appsv1.DeploymentCondition{
                Type:               appsv1.DeploymentConditionType(c.Type),
                Status:             corev1.ConditionStatus(c.Status),
                Reason:             c.Reason,
                Message:            c.Message,
                LastTransitionTime: c.LastTransitionTime,
            }
        }
    })
}
```



### diegen

Install diegen:

```sg
go install reconciler.io/dies/diegen
```

Create or update the generated dies:

```sh
diegen die:headerFile="hack/boilerplate.go.txt" paths="./..."
```

All generated content is created within `zz_generated.die.go` and `zz_generated.die_test.go` in each package where die markers are found.

### die markers

#### +die

```go
import (
    appsv1 "k8s.io/api/apps/v1"
)

// +die:object=true,apiVersion=apps/v1,kind=Deployment
type _ = appsv1.Deployment

// +die
type _ = appsv1.DeploymentSpec

// +die
type _ = appsv1.DeploymentStatus
```

For packages you control, dies can be created in the same package as the resource they model by adding the markers to existing types.

```go
// +die:object=true,apiVersion=example.reconciler.io/v1,kind=MyResource
type MyResource struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Spec MyResourceSpec `json:"spec"`
    // +optional
    Status MyResourceStatus `json:"status"`
}

// +die
type MyResourceSpec struct {
    // add fields
}

// +die
type MyResourceStatus struct {
    // add fields
}
```

Properties:
- **object** `bool` (optional): indicates the target type implements `metav1.Object` and `runtime.Object`
- **apiVersion** `string` (optional): defaults the blanks die's APIVersion (only for objects)
- **kind** `string` (optional): defaults the blank die's Kind (only for objects)
- **ignore** `[]string` (optional): set of fields to ignore on the type

#### +die:field

```go
// +die
// +die:field:name=Selector,package=reconciler.io/dies/apis/meta/v1,die=LabelSelectorDie,pointer=true
// +die:field:name=Template,package=reconciler.io/dies/apis/core/v1,die=PodTemplateSpecDie
type MyResourceSpec struct {
    Selector *metav1.LabelSelector  `json:"selector`"
    Template corev1.PodTemplateSpec `json:"template"`
}
```

Results in the standard die method and two additional mutation methods for the defined fields:

```go
// SelectorDie mutates Selector as a die
//
// Label selector for pods. It must match the pod template's labels.
func (d *MyResourceSpecDie) SelectorDie(fn func(d *metav1.LabelSelectorDie)) *MyResourceSpecDie {
	return d.DieStamp(func(r *appsv1.DeploymentSpec) {
		d := metav1.LabelSelectorBlank.DieImmutable(false).DieFeedPtr(r.Selector)
		fn(d)
		r.Selector = d.DieReleasePtr()
	})
}

// TemplateDie mutates Template as a die
//
// Template describes the pods that will be created.
func (d *MyResourceSpecDie) TemplateDie(fn func(d *corev1.PodTemplateSpecDie)) *MyResourceSpecDie {
	return d.DieStamp(func(r *appsv1.DeploymentSpec) {
		d := corev1.PodTemplateSpecBlank.DieImmutable(false).DieFeed(r.Template)
		fn(d)
		r.Template = d.DieRelease()
	})
}
```

Mutation methods for fields backed by slices require additional metadata. Depending on the content of the list it the whole listing can be replaced, or a single item in the list can be mutated.

Atomic lists cannot be mutated incrementally, the whole content is replaced. This pattern is common where there is no unique field inside the list to find a specific entry for.

```go
// +die:field:name=Ports,die=ContainerPortDie,listType=atomic
```

```go
// PortsDie replaces Ports by collecting the released value from each die passed.
func (d *ContainerDie) PortsDie(ports ...*ContainerPortDie) *ContainerDie {
	...snip...
}
```
[Full source](https://pkg.go.dev/reconciler.io/dies/apis/core/v1#ContainerDie.PortsDie)

The env field is a map list type with the default key. The named env var in the list will be loaded into the EnvVarDie and passed to the callback with the mutated value updated in the slice.

Because a single item is mutated the method name uses a singular form. If the conversion from plural to singular produces a non-ideal name, the default method name can be overridden with the **method** parameter.

```go
// +die:field:name=Env,die=EnvSourceDie,listType=map
```

```go
// EnvDie mutates a single item in Env matched by the nested field Name, appending a new item if no match is found.
func (d *ContainerDie) EnvDie(name string, fn func(d *EnvVarDie)) *ContainerDie {
    ...snip...
}
```
[Full source](https://pkg.go.dev/reconciler.io/dies/apis/core/v1#ContainerDie.EnvDie)

The env from field is a map list type with a custom key. The default map key is `Name`.

```go
// +die:field:name=EnvFrom,die=EnvFromSourceDie,listType=map,listMapKey=Prefix
```

```go
// EnvFromDie mutates a single item in EnvFrom matched by the nested field Prefix, appending a new item if no match is found.
func (d *ContainerDie) EnvFromDie(prefix string, fn func(d *EnvFromSourceDie)) *ContainerDie {
    ...snip...
}
```
[Full source](https://pkg.go.dev/reconciler.io/dies/apis/core/v1#ContainerDie.EnvFromDie)


The resize policy field is also a map list type, but the map key field is defined with a custom type.

```go
// +die:field:name=ResizePolicy,die=ContainerResizePolicyDie,listType=map,listMapKey=ResourceName,listMapKeyPackage=k8s.io/api/core/v1,listMapKeyType=ResourceName
```

```go
// ResizePolicyDie mutates a single item in ResizePolicy matched by the nested field ResourceName, appending a new item if no match is found.
func (d *ContainerDie) ResizePolicyDie(resourceName corev1.ResourceName, fn func(d *ContainerResizePolicyDie)) *ContainerDie {
    ...snip...
}
```
[Full source](https://pkg.go.dev/reconciler.io/dies/apis/core/v1#ContainerDie.ResizePolicyDie)


Properties:
- **name** `string`: name of the field on the target resource to mutate
- **method** `string` (optional): name of the generated mutation method on the host die (defaults to the name with "Die" appended, for list type map the name is made singular)
- **package** `string` (optional): go package containing the target die and blank. `_/` is an alias to the package prefix for the built in dies, for example `_/core/v1` expands to `reconciler.io/dies/apis/core/v1`.
- **die** `string`: go type for the die backing the field's struct
- **blank** `string` (optional): go type for a blank of the die's type (defaults to the die replacing the "Die" suffix with "Blank")
- **pointer** `bool` (optional): the field is defined as a pointer reference, uses DieFeedPtr and DieReleasePtr
- **listType** `string` (optional):
  - ``: field is not a list (default)
  - `map`: a single value in the list is mutated, or a new item appended
  - `atomic`: the list is replaced
- **listMapKey** `string` (optional): defaults to `Name` for map list types (only for list type map)
- **listMapKeyPackage** `string` (optional): the go package containing the type representing the key field on the target struct, defaults to the current package (only for list type map)
- **listMapKeyType** `string` (optional): the go type representing the key field on the target struct, defaults to `string` (only for list type map)

## Community

### Code of Conduct

The reconciler.io projects follow the [Contributor Covenant Code of Conduct](./CODE_OF_CONDUCT.md). In short, be kind and treat others with respect.

### Communication

General discussion and questions about the project can occur either on the Kubernetes Slack [#reconcilerio](https://kubernetes.slack.com/archives/C07J5G9NDHR) channel, or in the project's [GitHub discussions](https://github.com/orgs/reconcilerio/discussions). Use the channel you find most comfortable.

### Contributing

The reconciler.io dies project team welcomes contributions from the community. A contributor license agreement (CLA) is not required. You own full rights to your contribution and agree to license the work to the community under the Apache License v2.0, via a [Developer Certificate of Origin (DCO)](https://developercertificate.org). For more detailed information, refer to [CONTRIBUTING.md](CONTRIBUTING.md).

## Acknowledgements

`dies` was conceived by [Scott Andrews](https://github.com/scothis) drawing inspiration from previous work in [`reconciler.io/runtime`](https://github.com/reconcilerio/runtime/).

## License

Apache License v2.0: see [LICENSE](./LICENSE) for details.
