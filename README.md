# Dies - immutable, fluent, builders for Kubernetes resources <!-- omit in toc -->

![CI](https://github.com/scothis/dies/workflows/CI/badge.svg?branch=main)
[![codecov](https://codecov.io/gh/scothis/dies/branch/main/graph/badge.svg?token=vyXLcPBdV3)](https://codecov.io/gh/scothis/dies)
[![Go Reference](https://pkg.go.dev/badge/dies.dev.svg)](https://pkg.go.dev/dies.dev)
[![Go Report Card](https://goreportcard.com/badge/dies.dev)](https://goreportcard.com/report/dies.dev)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](CODE_OF_CONDUCT.md)

- [Using dies](#using-dies)
  - [Common methods](#common-methods)
- [Creating dies](#creating-dies)
  - [diegen](#diegen)
  - [die markers](#die-markers)
    - [+die](#die)

---

This project contains dies for many of the most common built-in Kubernetes types, and tools to create dies for custom resources.

## Using dies

Dies start with a blank object that is stamped to add state. All dies are immutable, each stamping returns a new instance containing the mutated state while the original instance is not modified.

```go
import (
    dieappsv1 "dies.dev/apis/apps/v1"
    diecorev1 "dies.dev/apis/core/v1"
    diemetav1 "dies.dev/apis/meta/v1"
    appsv1 "k8s.io/api/apps/v1"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
go install dies.dev/diegen
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

// +die:object=true
type _ = appsv1.Deployment

// +die
type _ = appsv1.DeploymentSpec

// +die
type _ = appsv1.DeploymentStatus
```

For packages you control, dies can be created in the same package as the resource they model by adding the markers to existing types.

```go
// +die:object=true
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
- **ignores** `[]string` (optional): set of fields to ignore on the type
