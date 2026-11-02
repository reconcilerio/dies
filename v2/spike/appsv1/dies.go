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

package appsv1

import (
	"reconciler.io/dies/v2/die"
	"reconciler.io/dies/v2/die/base"
	"reconciler.io/dies/v2/die/fields"
	"reconciler.io/dies/v2/spike/corev1"
	"reconciler.io/dies/v2/spike/metav1"
)

// Deployment

type DeploymentDie struct {
	base.DieBase[Deployment, *DeploymentDie]
}

func (*DeploymentDie) DieFrom(b base.DieBase[Deployment, *DeploymentDie]) *DeploymentDie {
	return &DeploymentDie{b}
}

var DeploymentBlank = base.Blank[Deployment, *DeploymentDie]().
	APIVersion().Set("apps/v1").
	Kind().Set("Deployment").
	DieCommit()

var deploymentFields = struct {
	APIVersion fields.Field[Deployment, string]
	Kind       fields.Field[Deployment, string]
	Metadata   fields.Struct[Deployment, metav1.ObjectMeta, *metav1.ObjectMetaDie]
	Spec       fields.Struct[Deployment, DeploymentSpec, *DeploymentSpecDie]
	Status     fields.Struct[Deployment, DeploymentStatus, *DeploymentStatusDie]
}{
	APIVersion: fields.NewField[Deployment, string]("apiVersion"),
	Kind:       fields.NewField[Deployment, string]("kind"),
	Metadata:   fields.NewStruct[Deployment, metav1.ObjectMeta, *metav1.ObjectMetaDie]("metadata"),
	Spec:       fields.NewStruct[Deployment, DeploymentSpec, *DeploymentSpecDie]("spec"),
	Status:     fields.NewStruct[Deployment, DeploymentStatus, *DeploymentStatusDie]("status"),
}

func (d *DeploymentDie) APIVersion() die.Value[string, *DeploymentDie] {
	return fields.ValueOf(d, deploymentFields.APIVersion)
}

func (d *DeploymentDie) Kind() die.Value[string, *DeploymentDie] {
	return fields.ValueOf(d, deploymentFields.Kind)
}

func (d *DeploymentDie) Metadata() *metav1.ObjectMetaDie {
	return fields.Child(d, deploymentFields.Metadata)
}

func (d *DeploymentDie) Spec() *DeploymentSpecDie { return fields.Child(d, deploymentFields.Spec) }

func (d *DeploymentDie) Status() *DeploymentStatusDie {
	return fields.Child(d, deploymentFields.Status)
}

// DeploymentSpec

type DeploymentSpecDie struct {
	base.DieBase[DeploymentSpec, *DeploymentSpecDie]
}

func (*DeploymentSpecDie) DieFrom(b base.DieBase[DeploymentSpec, *DeploymentSpecDie]) *DeploymentSpecDie {
	return &DeploymentSpecDie{b}
}

var deploymentSpecFields = struct {
	Replicas fields.Field[DeploymentSpec, *int32]
	Paused   fields.Field[DeploymentSpec, bool]
	Template fields.Struct[DeploymentSpec, corev1.PodTemplateSpec, *corev1.PodTemplateSpecDie]
}{
	Replicas: fields.NewField[DeploymentSpec, *int32]("replicas"),
	Paused:   fields.NewField[DeploymentSpec, bool]("paused"),
	Template: fields.NewStruct[DeploymentSpec, corev1.PodTemplateSpec, *corev1.PodTemplateSpecDie]("template"),
}

func (d *DeploymentSpecDie) Replicas() die.Optional[int32, *DeploymentSpecDie] {
	return fields.OptionalOf(d, deploymentSpecFields.Replicas)
}

func (d *DeploymentSpecDie) Paused() die.Value[bool, *DeploymentSpecDie] {
	return fields.ValueOf(d, deploymentSpecFields.Paused)
}

func (d *DeploymentSpecDie) Template() *corev1.PodTemplateSpecDie {
	return fields.Child(d, deploymentSpecFields.Template)
}

// DeploymentStatus

type DeploymentStatusDie struct {
	base.DieBase[DeploymentStatus, *DeploymentStatusDie]
}

func (*DeploymentStatusDie) DieFrom(b base.DieBase[DeploymentStatus, *DeploymentStatusDie]) *DeploymentStatusDie {
	return &DeploymentStatusDie{b}
}

var deploymentStatusFields = struct {
	ObservedGeneration fields.Field[DeploymentStatus, int64]
	Replicas           fields.Field[DeploymentStatus, int32]
	Conditions         fields.ListMapField[DeploymentStatus, metav1.Condition, *metav1.ConditionDie, string]
}{
	ObservedGeneration: fields.NewField[DeploymentStatus, int64]("observedGeneration"),
	Replicas:           fields.NewField[DeploymentStatus, int32]("replicas"),
	Conditions: fields.NewListMap[DeploymentStatus, metav1.Condition, *metav1.ConditionDie]("conditions",
		func(t string) []any { return []any{t} },
		func(c metav1.Condition) string { return c.Type }),
}

func (d *DeploymentStatusDie) ObservedGeneration() die.Value[int64, *DeploymentStatusDie] {
	return fields.ValueOf(d, deploymentStatusFields.ObservedGeneration)
}

func (d *DeploymentStatusDie) Replicas() die.Value[int32, *DeploymentStatusDie] {
	return fields.ValueOf(d, deploymentStatusFields.Replicas)
}

func (d *DeploymentStatusDie) Conditions() *die.ListMap[metav1.Condition, *metav1.ConditionDie, string] {
	return fields.ListMapOf(d, deploymentStatusFields.Conditions)
}
