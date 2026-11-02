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
	base.DieBase[Deployment, *DeploymentDie, *DeploymentDieMut]
}

type DeploymentDieMut struct {
	base.DieMutBase[Deployment, *DeploymentDie, *DeploymentDieMut]
}

func (*DeploymentDie) DieFrom(b base.DieBase[Deployment, *DeploymentDie, *DeploymentDieMut]) *DeploymentDie {
	return &DeploymentDie{b}
}

func (*DeploymentDieMut) DieFrom(b base.DieMutBase[Deployment, *DeploymentDie, *DeploymentDieMut]) *DeploymentDieMut {
	return &DeploymentDieMut{b}
}

var DeploymentBlank = base.Blank[Deployment, *DeploymentDie, *DeploymentDieMut]().DieMutable().
	APIVersion().Set("apps/v1").
	Kind().Set("Deployment").
	DieCommit().
	DieImmutable()

var deploymentFields = struct {
	APIVersion fields.Field[Deployment, string]
	Kind       fields.Field[Deployment, string]
	Metadata   fields.Struct[Deployment, metav1.ObjectMeta, *metav1.ObjectMetaDie, *metav1.ObjectMetaDieMut]
	Spec       fields.Struct[Deployment, DeploymentSpec, *DeploymentSpecDie, *DeploymentSpecDieMut]
	Status     fields.Struct[Deployment, DeploymentStatus, *DeploymentStatusDie, *DeploymentStatusDieMut]
}{
	APIVersion: fields.NewField[Deployment, string]("apiVersion"),
	Kind:       fields.NewField[Deployment, string]("kind"),
	Metadata:   fields.NewStruct[Deployment, metav1.ObjectMeta, *metav1.ObjectMetaDie, *metav1.ObjectMetaDieMut]("metadata"),
	Spec:       fields.NewStruct[Deployment, DeploymentSpec, *DeploymentSpecDie, *DeploymentSpecDieMut]("spec"),
	Status:     fields.NewStruct[Deployment, DeploymentStatus, *DeploymentStatusDie, *DeploymentStatusDieMut]("status"),
}

func (d *DeploymentDie) APIVersion() die.Value[string] {
	return fields.ValueOf(d, deploymentFields.APIVersion)
}

func (d *DeploymentDieMut) APIVersion() die.ValueMut[string, *DeploymentDieMut] {
	return fields.ValueMutOf(d, deploymentFields.APIVersion)
}

func (d *DeploymentDie) Kind() die.Value[string] {
	return fields.ValueOf(d, deploymentFields.Kind)
}

func (d *DeploymentDieMut) Kind() die.ValueMut[string, *DeploymentDieMut] {
	return fields.ValueMutOf(d, deploymentFields.Kind)
}

func (d *DeploymentDie) Metadata() *metav1.ObjectMetaDie {
	return fields.Child(d, deploymentFields.Metadata)
}

func (d *DeploymentDieMut) Metadata() *metav1.ObjectMetaDieMut {
	return fields.ChildMut(d, deploymentFields.Metadata)
}

func (d *DeploymentDie) Spec() *DeploymentSpecDie { return fields.Child(d, deploymentFields.Spec) }

func (d *DeploymentDieMut) Spec() *DeploymentSpecDieMut {
	return fields.ChildMut(d, deploymentFields.Spec)
}

func (d *DeploymentDie) Status() *DeploymentStatusDie {
	return fields.Child(d, deploymentFields.Status)
}

func (d *DeploymentDieMut) Status() *DeploymentStatusDieMut {
	return fields.ChildMut(d, deploymentFields.Status)
}

// DeploymentSpec

type DeploymentSpecDie struct {
	base.DieBase[DeploymentSpec, *DeploymentSpecDie, *DeploymentSpecDieMut]
}

type DeploymentSpecDieMut struct {
	base.DieMutBase[DeploymentSpec, *DeploymentSpecDie, *DeploymentSpecDieMut]
}

func (*DeploymentSpecDie) DieFrom(b base.DieBase[DeploymentSpec, *DeploymentSpecDie, *DeploymentSpecDieMut]) *DeploymentSpecDie {
	return &DeploymentSpecDie{b}
}

func (*DeploymentSpecDieMut) DieFrom(b base.DieMutBase[DeploymentSpec, *DeploymentSpecDie, *DeploymentSpecDieMut]) *DeploymentSpecDieMut {
	return &DeploymentSpecDieMut{b}
}

var deploymentSpecFields = struct {
	Replicas fields.Field[DeploymentSpec, *int32]
	Paused   fields.Field[DeploymentSpec, bool]
	Template fields.Struct[DeploymentSpec, corev1.PodTemplateSpec, *corev1.PodTemplateSpecDie, *corev1.PodTemplateSpecDieMut]
}{
	Replicas: fields.NewField[DeploymentSpec, *int32]("replicas"),
	Paused:   fields.NewField[DeploymentSpec, bool]("paused"),
	Template: fields.NewStruct[DeploymentSpec, corev1.PodTemplateSpec, *corev1.PodTemplateSpecDie, *corev1.PodTemplateSpecDieMut]("template"),
}

func (d *DeploymentSpecDie) Replicas() die.Optional[int32] {
	return fields.OptionalOf(d, deploymentSpecFields.Replicas)
}

func (d *DeploymentSpecDieMut) Replicas() die.OptionalMut[int32, *DeploymentSpecDieMut] {
	return fields.OptionalMutOf(d, deploymentSpecFields.Replicas)
}

func (d *DeploymentSpecDie) Paused() die.Value[bool] {
	return fields.ValueOf(d, deploymentSpecFields.Paused)
}

func (d *DeploymentSpecDieMut) Paused() die.ValueMut[bool, *DeploymentSpecDieMut] {
	return fields.ValueMutOf(d, deploymentSpecFields.Paused)
}

func (d *DeploymentSpecDie) Template() *corev1.PodTemplateSpecDie {
	return fields.Child(d, deploymentSpecFields.Template)
}

func (d *DeploymentSpecDieMut) Template() *corev1.PodTemplateSpecDieMut {
	return fields.ChildMut(d, deploymentSpecFields.Template)
}

// DeploymentStatus

type DeploymentStatusDie struct {
	base.DieBase[DeploymentStatus, *DeploymentStatusDie, *DeploymentStatusDieMut]
}

type DeploymentStatusDieMut struct {
	base.DieMutBase[DeploymentStatus, *DeploymentStatusDie, *DeploymentStatusDieMut]
}

func (*DeploymentStatusDie) DieFrom(b base.DieBase[DeploymentStatus, *DeploymentStatusDie, *DeploymentStatusDieMut]) *DeploymentStatusDie {
	return &DeploymentStatusDie{b}
}

func (*DeploymentStatusDieMut) DieFrom(b base.DieMutBase[DeploymentStatus, *DeploymentStatusDie, *DeploymentStatusDieMut]) *DeploymentStatusDieMut {
	return &DeploymentStatusDieMut{b}
}

var deploymentStatusFields = struct {
	ObservedGeneration fields.Field[DeploymentStatus, int64]
	Replicas           fields.Field[DeploymentStatus, int32]
	Conditions         fields.ListMapField[DeploymentStatus, metav1.Condition, *metav1.ConditionDie, *metav1.ConditionDieMut, string]
}{
	ObservedGeneration: fields.NewField[DeploymentStatus, int64]("observedGeneration"),
	Replicas:           fields.NewField[DeploymentStatus, int32]("replicas"),
	Conditions: fields.NewListMap[DeploymentStatus, metav1.Condition, *metav1.ConditionDie, *metav1.ConditionDieMut]("conditions",
		func(t string) []any { return []any{t} },
		func(c metav1.Condition) string { return c.Type }),
}

func (d *DeploymentStatusDie) ObservedGeneration() die.Value[int64] {
	return fields.ValueOf(d, deploymentStatusFields.ObservedGeneration)
}

func (d *DeploymentStatusDieMut) ObservedGeneration() die.ValueMut[int64, *DeploymentStatusDieMut] {
	return fields.ValueMutOf(d, deploymentStatusFields.ObservedGeneration)
}

func (d *DeploymentStatusDie) Replicas() die.Value[int32] {
	return fields.ValueOf(d, deploymentStatusFields.Replicas)
}

func (d *DeploymentStatusDieMut) Replicas() die.ValueMut[int32, *DeploymentStatusDieMut] {
	return fields.ValueMutOf(d, deploymentStatusFields.Replicas)
}

func (d *DeploymentStatusDie) Conditions() *die.ListMap[metav1.Condition, *metav1.ConditionDie, string] {
	return fields.ListMapOf(d, deploymentStatusFields.Conditions)
}
func (d *DeploymentStatusDieMut) Conditions() *die.ListMapMut[metav1.Condition, *metav1.ConditionDieMut, string] {
	return fields.ListMapMutOf(d, deploymentStatusFields.Conditions)
}
