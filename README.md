# Dies - immutable, fluent, builders for Kubernetes resources <!-- omit in toc -->

![CI](https://github.com/scothis/dies/workflows/CI/badge.svg?branch=main)
[![codecov](https://codecov.io/gh/scothis/dies/branch/main/graph/badge.svg?token=vyXLcPBdV3)](https://codecov.io/gh/scothis/dies)
[![Go Reference](https://pkg.go.dev/badge/github.com/scothis/dies.svg)](https://pkg.go.dev/github.com/scothis/dies)
[![Go Report Card](https://goreportcard.com/badge/github.com/scothis/dies)](https://goreportcard.com/report/github.com/scothis/dies)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](CODE_OF_CONDUCT.md)

- [Using dies](#using-dies)
	- [Common methods](#common-methods)
- [Creating dies](#creating-dies)
	- [diegen](#diegen)
	- [die markers](#die-markers)
		- [+die](#die)
		- [+die:field](#diefield)
		- [+die:scheme](#diescheme)

---

This project contains dies for many of the most common built-in Kubernetes types, and tools to create dies for custom resources.

## Using dies

Dies start with a blank object that is stamped to add state. All dies are immutable, each stamping returns a new instance containing the mutated state while the original instance is not modified.

```go
import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	dieappsv1 "github.com/scothis/dies/apis/apps/v1"
	diecorev1 "github.com/scothis/dies/apis/core/v1"
	diemetav1 "github.com/scothis/dies/apis/meta/v1"
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
var <T>Blank = (&<T>Die{}).DieFeed(<T>{})

// DieStamp returns a new die with the resource passed to the callback
// function. The resource is mutable.
func (d *<T>Die) DieStamp(fn func(r *<T>)) *<T>Die

// DieFeed returns a new die with the provided resource.
func (d *<T>Die) DieFeed(r <T>) *D<T>Die

// DieRelease returns the resource managed by the die.
func (d *<T>Die) DieRelease() <T>

// DieImmutable returns a new die for the current die's state that is either
// mutable (`false`) or immutable (`true`). 
func (d *<T>Die) DieImmutable(immutable bool) *<T>Die

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a
// mutable die.
func (d *<T>Die) DeepCopy() *<T>Die
```

Dies marked as implementing `metav1.Object` and `runtime.Object`  generate
additional methods.

```go
// DeepCopyObject returns a deep copy of the resource.
func (d *<T>Die) DeepCopyObject() runtime.Object

// GetObjectKind returns the resources's ObjectKind.
func (d *<T>Die) GetObjectKind() schema.ObjectKind

// MarshalJSON returns the die's resource as JSON.
func (d *<T>Die) MarshalJSON() ([]byte, error)

// UnmarshalJSON sets the die's resource from JSON.
func (d *<T>Die) UnmarshalJSON(b []byte) error

// MetadataDie mutates the resource's ObjectMeta field with a mutable die.
func (d *<T>Die) MetadataDie(fn func(d *metav1.ObjectMetaDie)) *<T>Die

// SpecDie mutates the resource's spec field with a mutable die. This method is
// only created if `<T>SpecDie` is defined.
func (d *<T>Die) SpecDie(fn func(d *<T>SpecDie)) *<T>Die

// StatusDie mutates the resource's status field with a mutable die. This
// method is only created if `<T>StatusDie` is defined.
func (d *<T>Die) StatusDie(fn func(d *<T>StatusDie)) *<T>Die
```

## Creating dies

Dies are primarily generated for types from [die markers](#die-markers) using
[`diegen`](#diegen).

Additional methods can be added to a die to enrich its behavior.

Example dispatching to a nested die to managed the template
`corev1.PodTemplateSpec` for `DeploymentSpec`:

```go
func (d *DeploymentSpecDie) TemplateDie(fn func(d *diecorev1.PodTemplateSpecDie)) *DeploymentSpecDie {
	return d.DieStamp(func(r *appsv1.DeploymentSpec) {
		d := diecorev1.PodTemplateSpecBlank.DieImmutable(false).DieFeed(r.Template)
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
go install github.com/scothis/dies/cmd/diegen
```

Create or update the generated dies:

```sh
diegen die:headerFile="hack/boilerplate.go.txt" paths="./..."
```

All generated content is created within `zz_generated.die.go` in each package where die markers are found.

### die markers

#### +die

```go
// +die:target=k8s.io/api/apps/v1.Deployment,object=true
// +die:target=k8s.io/api/apps/v1.DeploymentSpec
// +die:target=k8s.io/api/apps/v1.DeploymentStatus
```

Properties:
- **target** `string`: type the die manages
- **object** `bool` (optional): indicates the target type implements `metav1.Object` and `runtime.Object`

#### +die:field

```go
// +die:field:receiver=DeploymentSpecDie,name=Template,type=k8s.io/api/core/v1.PodTemplateSpec
```

Properties:
- **receiver** `string`: the die type managing the field
- **name** `string`: the field name on the managed resource
- **type** `string`: the field type on the managed resource


#### +die:scheme

```go
// +die:scheme:group=apps,version=v1
```

Properties:
- **group** `string`: api group for dies in this package
- **version** `string`: version for dies in this package
