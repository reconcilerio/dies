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

package corev1

import (
	"reconciler.io/dies/v2/die"
	"reconciler.io/dies/v2/die/base"
	"reconciler.io/dies/v2/die/fields"
	"reconciler.io/dies/v2/spike/metav1"
)

// PodTemplateSpec

type PodTemplateSpecDie struct {
	base.DieBase[PodTemplateSpec, *PodTemplateSpecDie]
}

func (*PodTemplateSpecDie) DieFrom(b base.DieBase[PodTemplateSpec, *PodTemplateSpecDie]) *PodTemplateSpecDie {
	return &PodTemplateSpecDie{b}
}

var PodTemplateSpecBlank = base.Blank[PodTemplateSpec, *PodTemplateSpecDie]()

var podTemplateSpecFields = struct {
	Metadata fields.Struct[PodTemplateSpec, metav1.ObjectMeta, *metav1.ObjectMetaDie]
	Spec     fields.Struct[PodTemplateSpec, PodSpec, *PodSpecDie]
}{
	Metadata: fields.NewStruct[PodTemplateSpec, metav1.ObjectMeta, *metav1.ObjectMetaDie]("metadata"),
	Spec:     fields.NewStruct[PodTemplateSpec, PodSpec, *PodSpecDie]("spec"),
}

func (d *PodTemplateSpecDie) Metadata() *metav1.ObjectMetaDie {
	return fields.Child(d, podTemplateSpecFields.Metadata)
}

func (d *PodTemplateSpecDie) Spec() *PodSpecDie { return fields.Child(d, podTemplateSpecFields.Spec) }

// PodSpec

type PodSpecDie struct {
	base.DieBase[PodSpec, *PodSpecDie]
}

func (*PodSpecDie) DieFrom(b base.DieBase[PodSpec, *PodSpecDie]) *PodSpecDie {
	return &PodSpecDie{b}
}

var PodSpecBlank = base.Blank[PodSpec, *PodSpecDie]()

var podSpecFields = struct {
	Volumes                      fields.ListMapField[PodSpec, Volume, *VolumeDie, string]
	InitContainers               fields.ListMapField[PodSpec, Container, *ContainerDie, string]
	Containers                   fields.ListMapField[PodSpec, Container, *ContainerDie, string]
	NodeSelector                 fields.MapField[PodSpec, string, string]
	ServiceAccountName           fields.Field[PodSpec, string]
	AutomountServiceAccountToken fields.Field[PodSpec, *bool]
}{
	Volumes:                      fields.NewListMap[PodSpec, Volume, *VolumeDie]("volumes", byName, func(v Volume) string { return v.Name }),
	InitContainers:               fields.NewListMap[PodSpec, Container, *ContainerDie]("initContainers", byName, containerName),
	Containers:                   fields.NewListMap[PodSpec, Container, *ContainerDie]("containers", byName, containerName),
	NodeSelector:                 fields.NewMap[PodSpec, string, string]("nodeSelector"),
	ServiceAccountName:           fields.NewField[PodSpec, string]("serviceAccountName"),
	AutomountServiceAccountToken: fields.NewField[PodSpec, *bool]("automountServiceAccountToken"),
}

func byName(name string) []any         { return []any{name} }
func containerName(c Container) string { return c.Name }

func (d *PodSpecDie) Volumes() *die.ListMap[Volume, *VolumeDie, string] {
	return fields.ListMapOf(d, podSpecFields.Volumes)
}

func (d *PodSpecDie) InitContainers() *die.ListMap[Container, *ContainerDie, string] {
	return fields.ListMapOf(d, podSpecFields.InitContainers)
}

func (d *PodSpecDie) Containers() *die.ListMap[Container, *ContainerDie, string] {
	return fields.ListMapOf(d, podSpecFields.Containers)
}

func (d *PodSpecDie) NodeSelector() *die.Map[string, string] {
	return fields.MapOf(d, podSpecFields.NodeSelector)
}

func (d *PodSpecDie) ServiceAccountName() die.Value[string, *PodSpecDie] {
	return fields.ValueOf(d, podSpecFields.ServiceAccountName)
}

func (d *PodSpecDie) AutomountServiceAccountToken() die.Optional[bool, *PodSpecDie] {
	return fields.OptionalOf(d, podSpecFields.AutomountServiceAccountToken)
}

// Container

type ContainerDie struct {
	base.DieBase[Container, *ContainerDie]
}

func (*ContainerDie) DieFrom(b base.DieBase[Container, *ContainerDie]) *ContainerDie {
	return &ContainerDie{b}
}

var ContainerBlank = base.Blank[Container, *ContainerDie]()

// PortKey identifies a container port.
type PortKey struct {
	ContainerPort int32
	Protocol      string
}

var containerFields = struct {
	Name         fields.Field[Container, string]
	Image        fields.Field[Container, string]
	Args         fields.Field[Container, []string]
	Ports        fields.ListMapField[Container, ContainerPort, *ContainerPortDie, PortKey]
	Env          fields.ListMapField[Container, EnvVar, *EnvVarDie, string]
	VolumeMounts fields.ListMapField[Container, VolumeMount, *VolumeMountDie, string]
}{
	Name:  fields.NewField[Container, string]("name"),
	Image: fields.NewField[Container, string]("image"),
	Args:  fields.NewField[Container, []string]("args"),
	Ports: fields.NewListMap[Container, ContainerPort, *ContainerPortDie]("ports",
		func(k PortKey) []any { return []any{k.ContainerPort, k.Protocol} },
		func(p ContainerPort) PortKey { return PortKey{p.ContainerPort, p.Protocol} }),
	Env: fields.NewListMap[Container, EnvVar, *EnvVarDie]("env", byName,
		func(e EnvVar) string { return e.Name }),
	VolumeMounts: fields.NewListMap[Container, VolumeMount, *VolumeMountDie]("volumeMounts",
		func(path string) []any { return []any{path} },
		func(m VolumeMount) string { return m.MountPath }),
}

func (d *ContainerDie) Name() die.Value[string, *ContainerDie] {
	return fields.ValueOf(d, containerFields.Name)
}

func (d *ContainerDie) Image() die.Value[string, *ContainerDie] {
	return fields.ValueOf(d, containerFields.Image)
}

func (d *ContainerDie) Args() die.Value[[]string, *ContainerDie] {
	return fields.ValueOf(d, containerFields.Args)
}

func (d *ContainerDie) Ports() *die.ListMap[ContainerPort, *ContainerPortDie, PortKey] {
	return fields.ListMapOf(d, containerFields.Ports)
}

func (d *ContainerDie) Env() *die.ListMap[EnvVar, *EnvVarDie, string] {
	return fields.ListMapOf(d, containerFields.Env)
}

func (d *ContainerDie) VolumeMounts() *die.ListMap[VolumeMount, *VolumeMountDie, string] {
	return fields.ListMapOf(d, containerFields.VolumeMounts)
}

// ContainerPort

type ContainerPortDie struct {
	base.DieBase[ContainerPort, *ContainerPortDie]
}

func (*ContainerPortDie) DieFrom(b base.DieBase[ContainerPort, *ContainerPortDie]) *ContainerPortDie {
	return &ContainerPortDie{b}
}

var containerPortFields = struct {
	Name          fields.Field[ContainerPort, string]
	ContainerPort fields.Field[ContainerPort, int32]
	Protocol      fields.Field[ContainerPort, string]
}{
	Name:          fields.NewField[ContainerPort, string]("name"),
	ContainerPort: fields.NewField[ContainerPort, int32]("containerPort"),
	Protocol:      fields.NewField[ContainerPort, string]("protocol"),
}

func (d *ContainerPortDie) Name() die.Value[string, *ContainerPortDie] {
	return fields.ValueOf(d, containerPortFields.Name)
}

func (d *ContainerPortDie) ContainerPort() die.Value[int32, *ContainerPortDie] {
	return fields.ValueOf(d, containerPortFields.ContainerPort)
}

func (d *ContainerPortDie) Protocol() die.Value[string, *ContainerPortDie] {
	return fields.ValueOf(d, containerPortFields.Protocol)
}

// EnvVar

type EnvVarDie struct {
	base.DieBase[EnvVar, *EnvVarDie]
}

func (*EnvVarDie) DieFrom(b base.DieBase[EnvVar, *EnvVarDie]) *EnvVarDie { return &EnvVarDie{b} }

var envVarFields = struct {
	Name  fields.Field[EnvVar, string]
	Value fields.Field[EnvVar, string]
}{
	Name:  fields.NewField[EnvVar, string]("name"),
	Value: fields.NewField[EnvVar, string]("value"),
}

func (d *EnvVarDie) Name() die.Value[string, *EnvVarDie] {
	return fields.ValueOf(d, envVarFields.Name)
}

func (d *EnvVarDie) Value() die.Value[string, *EnvVarDie] {
	return fields.ValueOf(d, envVarFields.Value)
}

// VolumeMount

type VolumeMountDie struct {
	base.DieBase[VolumeMount, *VolumeMountDie]
}

func (*VolumeMountDie) DieFrom(b base.DieBase[VolumeMount, *VolumeMountDie]) *VolumeMountDie {
	return &VolumeMountDie{b}
}

var volumeMountFields = struct {
	Name      fields.Field[VolumeMount, string]
	ReadOnly  fields.Field[VolumeMount, bool]
	MountPath fields.Field[VolumeMount, string]
}{
	Name:      fields.NewField[VolumeMount, string]("name"),
	ReadOnly:  fields.NewField[VolumeMount, bool]("readOnly"),
	MountPath: fields.NewField[VolumeMount, string]("mountPath"),
}

func (d *VolumeMountDie) Name() die.Value[string, *VolumeMountDie] {
	return fields.ValueOf(d, volumeMountFields.Name)
}

func (d *VolumeMountDie) ReadOnly() die.Value[bool, *VolumeMountDie] {
	return fields.ValueOf(d, volumeMountFields.ReadOnly)
}

func (d *VolumeMountDie) MountPath() die.Value[string, *VolumeMountDie] {
	return fields.ValueOf(d, volumeMountFields.MountPath)
}

// Volume: the members of the inlined VolumeSource are accessed directly.

type VolumeDie struct {
	base.DieBase[Volume, *VolumeDie]
}

func (*VolumeDie) DieFrom(b base.DieBase[Volume, *VolumeDie]) *VolumeDie { return &VolumeDie{b} }

var volumeFields = struct {
	Name     fields.Field[Volume, string]
	Secret   fields.Struct[Volume, SecretVolumeSource, *SecretVolumeSourceDie]
	EmptyDir fields.Struct[Volume, EmptyDirVolumeSource, *EmptyDirVolumeSourceDie]
}{
	Name:     fields.NewField[Volume, string]("name"),
	Secret:   fields.NewStruct[Volume, SecretVolumeSource, *SecretVolumeSourceDie]("secret"),
	EmptyDir: fields.NewStruct[Volume, EmptyDirVolumeSource, *EmptyDirVolumeSourceDie]("emptyDir"),
}

func (d *VolumeDie) Name() die.Value[string, *VolumeDie] { return fields.ValueOf(d, volumeFields.Name) }

func (d *VolumeDie) Secret() *SecretVolumeSourceDie { return fields.Child(d, volumeFields.Secret) }

func (d *VolumeDie) EmptyDir() *EmptyDirVolumeSourceDie {
	return fields.Child(d, volumeFields.EmptyDir)
}

type SecretVolumeSourceDie struct {
	base.DieBase[SecretVolumeSource, *SecretVolumeSourceDie]
}

func (*SecretVolumeSourceDie) DieFrom(b base.DieBase[SecretVolumeSource, *SecretVolumeSourceDie]) *SecretVolumeSourceDie {
	return &SecretVolumeSourceDie{b}
}

var secretVolumeSourceFields = struct {
	SecretName fields.Field[SecretVolumeSource, string]
}{
	SecretName: fields.NewField[SecretVolumeSource, string]("secretName"),
}

func (d *SecretVolumeSourceDie) SecretName() die.Value[string, *SecretVolumeSourceDie] {
	return fields.ValueOf(d, secretVolumeSourceFields.SecretName)
}

type EmptyDirVolumeSourceDie struct {
	base.DieBase[EmptyDirVolumeSource, *EmptyDirVolumeSourceDie]
}

func (*EmptyDirVolumeSourceDie) DieFrom(b base.DieBase[EmptyDirVolumeSource, *EmptyDirVolumeSourceDie]) *EmptyDirVolumeSourceDie {
	return &EmptyDirVolumeSourceDie{b}
}

var emptyDirVolumeSourceFields = struct {
	Medium fields.Field[EmptyDirVolumeSource, string]
}{
	Medium: fields.NewField[EmptyDirVolumeSource, string]("medium"),
}

func (d *EmptyDirVolumeSourceDie) Medium() die.Value[string, *EmptyDirVolumeSourceDie] {
	return fields.ValueOf(d, emptyDirVolumeSourceFields.Medium)
}
