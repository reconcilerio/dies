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
                    d.AddEnv("MY_VAR", "my-value")
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
                    d.AddEnv("MY_VAR", "some-other-value")
                })
            })
        })
    }).DieRelease())
```

While the outer die is immutable, returning a new die for each call that can be chained together. The dies passed to callbacks are mutable.

Additional methods will be added to dies over time to make common operations easier and safer.

### Common methods

```go
// for managed type `<T>`

// <T>Blank is an empty die that mutations can be stamped from. All die blanks
// are immutable.
var <T>Blank <T>Die

type <T>Die interface {
    // DieStamp returns a new die with the resource passed to the callback
    // function. The resource is mutable.
    func DieStamp(fn func(r *<T>)) <T>Die

    // DieFeed returns a new die with the provided resource.
    func DieFeed(r <T>) <T>Die

    // DieFeedPtr returns a new die with the provided resource pointer. If the
    // resource is nil, the empty value is used instead.
    func DieFeedPtr(r *<T>) <T>Die

    // DieRelease returns the resource managed by the die.
    func DieRelease() <T>

    // DieReleasePtr returns a pointer to the resource managed by the die.
    func DieReleasePtr() *<T>

    // DieImmutable returns a new die for the current die's state that is
    // either mutable (`false`) or immutable (`true`). 
    func DieImmutable(immutable bool) <T>Die

    // DeepCopy returns a new die with equivalent state. Useful for
    // snapshotting a mutable die.
    func DeepCopy() <T>Die
}
```

For each exported field `<F>` on `<T>`, a method is registered to set that field.

```go
type <T>Die interface {
    // continued

    <F>(<T>) <T>Die
}
```

Dies marked as implementing `metav1.Object` and `runtime.Object`  generate
additional methods.

```go
type <T>Die interface {
    // continued

    runtime.Object
    metav1.Object
    metav1.ObjectMetaAccessor

    // MetadataDie stamps the resource's ObjectMeta field with a mutable die.
    func MetadataDie(fn func(d diemetav1.ObjectMetaDie)) <T>Die

    // SpecDie stamps the resource's spec field with a mutable die. This method
    // is only created if `<T>SpecDie` is defined.
    func SpecDie(fn func(d <T>SpecDie)) <T>Die

    // StatusDie stamps the resource's status field with a mutable die. This
    // method is only created if `<T>StatusDie` is defined.
    func StatusDie(fn func(d <T>StatusDie)) <T>Die
}
```

## Creating dies

Dies are primarily generated for types from [die markers](#die-markers) using
[`diegen`](#diegen).

Additional methods can be added to a die to enrich its behavior. Each additional method needs to be defined in a non-exported interface matching the target type.

Example dispatching to a nested die to managed the template
`corev1.PodTemplateSpec` for `DeploymentSpec`:

```go
type deploymentSpec interface {
    TemplateDie(fn func(d diecorev1.PodTemplateSpecDie)) DeploymentSpecDie
    // add other methods to contributed to DeploymentSpecDie
}

func (d *deploymentSpecDie) TemplateDie(fn func(d diecorev1.PodTemplateSpecDie)) DeploymentSpecDie {
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
type deploymentStatus interface {
    ConditionsDie(conditions ...diemetav1.ConditionDie) DeploymentStatusDie
    // add other methods to contributed to DeploymentStatusDie
}

func (d *deploymentStatusDie) ConditionsDie(conditions ...diemetav1.ConditionDie) DeploymentStatusDie {
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

All generated content is created within `zz_generated.die.go` in each package where die markers are found.

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

Properties:
- **object** `bool` (optional): indicates the target type implements `metav1.Object` and `runtime.Object`
- **ignores** `[]string` (optional): set of fields to ignore on the type
