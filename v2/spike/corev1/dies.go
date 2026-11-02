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
	base.DieBase[PodTemplateSpec, *PodTemplateSpecDie, *PodTemplateSpecDieMut]
}

type PodTemplateSpecDieMut struct {
	base.DieMutBase[PodTemplateSpec, *PodTemplateSpecDie, *PodTemplateSpecDieMut]
}

func (*PodTemplateSpecDie) DieFrom(b base.DieBase[PodTemplateSpec, *PodTemplateSpecDie, *PodTemplateSpecDieMut]) *PodTemplateSpecDie {
	return &PodTemplateSpecDie{b}
}

func (*PodTemplateSpecDieMut) DieFrom(b base.DieMutBase[PodTemplateSpec, *PodTemplateSpecDie, *PodTemplateSpecDieMut]) *PodTemplateSpecDieMut {
	return &PodTemplateSpecDieMut{b}
}

var PodTemplateSpecBlank = base.Blank[PodTemplateSpec, *PodTemplateSpecDie, *PodTemplateSpecDieMut]()

var podTemplateSpecFields = struct {
	Metadata fields.Struct[PodTemplateSpec, metav1.ObjectMeta, *metav1.ObjectMetaDie, *metav1.ObjectMetaDieMut]
	Spec     fields.Struct[PodTemplateSpec, PodSpec, *PodSpecDie, *PodSpecDieMut]
}{
	Metadata: fields.NewStruct[PodTemplateSpec, metav1.ObjectMeta, *metav1.ObjectMetaDie, *metav1.ObjectMetaDieMut]("metadata"),
	Spec:     fields.NewStruct[PodTemplateSpec, PodSpec, *PodSpecDie, *PodSpecDieMut]("spec"),
}

func (d *PodTemplateSpecDie) Metadata() *metav1.ObjectMetaDie {
	return fields.Child(d, podTemplateSpecFields.Metadata)
}

func (d *PodTemplateSpecDieMut) Metadata() *metav1.ObjectMetaDieMut {
	return fields.ChildMut(d, podTemplateSpecFields.Metadata)
}

func (d *PodTemplateSpecDie) Spec() *PodSpecDie { return fields.Child(d, podTemplateSpecFields.Spec) }

func (d *PodTemplateSpecDieMut) Spec() *PodSpecDieMut {
	return fields.ChildMut(d, podTemplateSpecFields.Spec)
}

// PodSpec

type PodSpecDie struct {
	base.DieBase[PodSpec, *PodSpecDie, *PodSpecDieMut]
}

type PodSpecDieMut struct {
	base.DieMutBase[PodSpec, *PodSpecDie, *PodSpecDieMut]
}

func (*PodSpecDie) DieFrom(b base.DieBase[PodSpec, *PodSpecDie, *PodSpecDieMut]) *PodSpecDie {
	return &PodSpecDie{b}
}

func (*PodSpecDieMut) DieFrom(b base.DieMutBase[PodSpec, *PodSpecDie, *PodSpecDieMut]) *PodSpecDieMut {
	return &PodSpecDieMut{b}
}

var PodSpecBlank = base.Blank[PodSpec, *PodSpecDie, *PodSpecDieMut]()

var podSpecFields = struct {
	Volumes                      fields.ListMapField[PodSpec, Volume, *VolumeDie, *VolumeDieMut, string]
	InitContainers               fields.ListMapField[PodSpec, Container, *ContainerDie, *ContainerDieMut, string]
	Containers                   fields.ListMapField[PodSpec, Container, *ContainerDie, *ContainerDieMut, string]
	NodeSelector                 fields.MapField[PodSpec, string, string]
	ServiceAccountName           fields.Field[PodSpec, string]
	AutomountServiceAccountToken fields.Field[PodSpec, *bool]
}{
	Volumes:                      fields.NewListMap[PodSpec, Volume, *VolumeDie, *VolumeDieMut]("volumes", byName, func(v Volume) string { return v.Name }),
	InitContainers:               fields.NewListMap[PodSpec, Container, *ContainerDie, *ContainerDieMut]("initContainers", byName, containerName),
	Containers:                   fields.NewListMap[PodSpec, Container, *ContainerDie, *ContainerDieMut]("containers", byName, containerName),
	NodeSelector:                 fields.NewMap[PodSpec, string, string]("nodeSelector"),
	ServiceAccountName:           fields.NewField[PodSpec, string]("serviceAccountName"),
	AutomountServiceAccountToken: fields.NewField[PodSpec, *bool]("automountServiceAccountToken"),
}

func byName(name string) []any         { return []any{name} }
func containerName(c Container) string { return c.Name }

func (d *PodSpecDie) Volumes() *die.ListMap[Volume, *VolumeDie, string] {
	return fields.ListMapOf(d, podSpecFields.Volumes)
}

func (d *PodSpecDieMut) Volumes() *die.ListMapMut[Volume, *VolumeDieMut, string] {
	return fields.ListMapMutOf(d, podSpecFields.Volumes)
}

func (d *PodSpecDie) InitContainers() *die.ListMap[Container, *ContainerDie, string] {
	return fields.ListMapOf(d, podSpecFields.InitContainers)
}

func (d *PodSpecDieMut) InitContainers() *die.ListMapMut[Container, *ContainerDieMut, string] {
	return fields.ListMapMutOf(d, podSpecFields.InitContainers)
}

func (d *PodSpecDie) Containers() *die.ListMap[Container, *ContainerDie, string] {
	return fields.ListMapOf(d, podSpecFields.Containers)
}

func (d *PodSpecDieMut) Containers() *die.ListMapMut[Container, *ContainerDieMut, string] {
	return fields.ListMapMutOf(d, podSpecFields.Containers)
}

func (d *PodSpecDie) NodeSelector() *die.Map[string, string] {
	return fields.MapOf(d, podSpecFields.NodeSelector)
}

func (d *PodSpecDieMut) NodeSelector() *die.MapMut[string, string] {
	return fields.MapMutOf(d, podSpecFields.NodeSelector)
}

func (d *PodSpecDie) ServiceAccountName() die.Value[string] {
	return fields.ValueOf(d, podSpecFields.ServiceAccountName)
}

func (d *PodSpecDieMut) ServiceAccountName() die.ValueMut[string, *PodSpecDieMut] {
	return fields.ValueMutOf(d, podSpecFields.ServiceAccountName)
}

func (d *PodSpecDie) AutomountServiceAccountToken() die.Optional[bool] {
	return fields.OptionalOf(d, podSpecFields.AutomountServiceAccountToken)
}

func (d *PodSpecDieMut) AutomountServiceAccountToken() die.OptionalMut[bool, *PodSpecDieMut] {
	return fields.OptionalMutOf(d, podSpecFields.AutomountServiceAccountToken)
}

// Container

type ContainerDie struct {
	base.DieBase[Container, *ContainerDie, *ContainerDieMut]
}

type ContainerDieMut struct {
	base.DieMutBase[Container, *ContainerDie, *ContainerDieMut]
}

func (*ContainerDie) DieFrom(b base.DieBase[Container, *ContainerDie, *ContainerDieMut]) *ContainerDie {
	return &ContainerDie{b}
}

func (*ContainerDieMut) DieFrom(b base.DieMutBase[Container, *ContainerDie, *ContainerDieMut]) *ContainerDieMut {
	return &ContainerDieMut{b}
}

var ContainerBlank = base.Blank[Container, *ContainerDie, *ContainerDieMut]()

// PortKey identifies a container port.
type PortKey struct {
	ContainerPort int32
	Protocol      string
}

var containerFields = struct {
	Name         fields.Field[Container, string]
	Image        fields.Field[Container, string]
	Args         fields.Field[Container, []string]
	Ports        fields.ListMapField[Container, ContainerPort, *ContainerPortDie, *ContainerPortDieMut, PortKey]
	Env          fields.ListMapField[Container, EnvVar, *EnvVarDie, *EnvVarDieMut, string]
	VolumeMounts fields.ListMapField[Container, VolumeMount, *VolumeMountDie, *VolumeMountDieMut, string]
}{
	Name:  fields.NewField[Container, string]("name"),
	Image: fields.NewField[Container, string]("image"),
	Args:  fields.NewField[Container, []string]("args"),
	Ports: fields.NewListMap[Container, ContainerPort, *ContainerPortDie, *ContainerPortDieMut]("ports",
		func(k PortKey) []any { return []any{k.ContainerPort, k.Protocol} },
		func(p ContainerPort) PortKey { return PortKey{p.ContainerPort, p.Protocol} }),
	Env: fields.NewListMap[Container, EnvVar, *EnvVarDie, *EnvVarDieMut]("env", byName,
		func(e EnvVar) string { return e.Name }),
	VolumeMounts: fields.NewListMap[Container, VolumeMount, *VolumeMountDie, *VolumeMountDieMut]("volumeMounts",
		func(path string) []any { return []any{path} },
		func(m VolumeMount) string { return m.MountPath }),
}

func (d *ContainerDie) Name() die.Value[string] {
	return fields.ValueOf(d, containerFields.Name)
}

func (d *ContainerDieMut) Name() die.ValueMut[string, *ContainerDieMut] {
	return fields.ValueMutOf(d, containerFields.Name)
}

func (d *ContainerDie) Image() die.Value[string] {
	return fields.ValueOf(d, containerFields.Image)
}

func (d *ContainerDieMut) Image() die.ValueMut[string, *ContainerDieMut] {
	return fields.ValueMutOf(d, containerFields.Image)
}

func (d *ContainerDie) Args() die.Value[[]string] {
	return fields.ValueOf(d, containerFields.Args)
}

func (d *ContainerDieMut) Args() die.ValueMut[[]string, *ContainerDieMut] {
	return fields.ValueMutOf(d, containerFields.Args)
}

func (d *ContainerDie) Ports() *die.ListMap[ContainerPort, *ContainerPortDie, PortKey] {
	return fields.ListMapOf(d, containerFields.Ports)
}

func (d *ContainerDieMut) Ports() *die.ListMapMut[ContainerPort, *ContainerPortDieMut, PortKey] {
	return fields.ListMapMutOf(d, containerFields.Ports)
}

func (d *ContainerDie) Env() *die.ListMap[EnvVar, *EnvVarDie, string] {
	return fields.ListMapOf(d, containerFields.Env)
}

func (d *ContainerDieMut) Env() *die.ListMapMut[EnvVar, *EnvVarDieMut, string] {
	return fields.ListMapMutOf(d, containerFields.Env)
}

func (d *ContainerDie) VolumeMounts() *die.ListMap[VolumeMount, *VolumeMountDie, string] {
	return fields.ListMapOf(d, containerFields.VolumeMounts)
}

func (d *ContainerDieMut) VolumeMounts() *die.ListMapMut[VolumeMount, *VolumeMountDieMut, string] {
	return fields.ListMapMutOf(d, containerFields.VolumeMounts)
}

// ContainerPort

type ContainerPortDie struct {
	base.DieBase[ContainerPort, *ContainerPortDie, *ContainerPortDieMut]
}

type ContainerPortDieMut struct {
	base.DieMutBase[ContainerPort, *ContainerPortDie, *ContainerPortDieMut]
}

func (*ContainerPortDie) DieFrom(b base.DieBase[ContainerPort, *ContainerPortDie, *ContainerPortDieMut]) *ContainerPortDie {
	return &ContainerPortDie{b}
}

func (*ContainerPortDieMut) DieFrom(b base.DieMutBase[ContainerPort, *ContainerPortDie, *ContainerPortDieMut]) *ContainerPortDieMut {
	return &ContainerPortDieMut{b}
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

func (d *ContainerPortDie) Name() die.Value[string] {
	return fields.ValueOf(d, containerPortFields.Name)
}

func (d *ContainerPortDieMut) Name() die.ValueMut[string, *ContainerPortDieMut] {
	return fields.ValueMutOf(d, containerPortFields.Name)
}

func (d *ContainerPortDie) ContainerPort() die.Value[int32] {
	return fields.ValueOf(d, containerPortFields.ContainerPort)
}

func (d *ContainerPortDieMut) ContainerPort() die.ValueMut[int32, *ContainerPortDieMut] {
	return fields.ValueMutOf(d, containerPortFields.ContainerPort)
}

func (d *ContainerPortDie) Protocol() die.Value[string] {
	return fields.ValueOf(d, containerPortFields.Protocol)
}

func (d *ContainerPortDieMut) Protocol() die.ValueMut[string, *ContainerPortDieMut] {
	return fields.ValueMutOf(d, containerPortFields.Protocol)
}

// EnvVar

type EnvVarDie struct {
	base.DieBase[EnvVar, *EnvVarDie, *EnvVarDieMut]
}

type EnvVarDieMut struct {
	base.DieMutBase[EnvVar, *EnvVarDie, *EnvVarDieMut]
}

func (*EnvVarDie) DieFrom(b base.DieBase[EnvVar, *EnvVarDie, *EnvVarDieMut]) *EnvVarDie {
	return &EnvVarDie{b}
}

func (*EnvVarDieMut) DieFrom(b base.DieMutBase[EnvVar, *EnvVarDie, *EnvVarDieMut]) *EnvVarDieMut {
	return &EnvVarDieMut{b}
}

var envVarFields = struct {
	Name  fields.Field[EnvVar, string]
	Value fields.Field[EnvVar, string]
}{
	Name:  fields.NewField[EnvVar, string]("name"),
	Value: fields.NewField[EnvVar, string]("value"),
}

func (d *EnvVarDie) Name() die.Value[string] {
	return fields.ValueOf(d, envVarFields.Name)
}

func (d *EnvVarDieMut) Name() die.ValueMut[string, *EnvVarDieMut] {
	return fields.ValueMutOf(d, envVarFields.Name)
}

func (d *EnvVarDie) Value() die.Value[string] {
	return fields.ValueOf(d, envVarFields.Value)
}

func (d *EnvVarDieMut) Value() die.ValueMut[string, *EnvVarDieMut] {
	return fields.ValueMutOf(d, envVarFields.Value)
}

// VolumeMount

type VolumeMountDie struct {
	base.DieBase[VolumeMount, *VolumeMountDie, *VolumeMountDieMut]
}

type VolumeMountDieMut struct {
	base.DieMutBase[VolumeMount, *VolumeMountDie, *VolumeMountDieMut]
}

func (*VolumeMountDie) DieFrom(b base.DieBase[VolumeMount, *VolumeMountDie, *VolumeMountDieMut]) *VolumeMountDie {
	return &VolumeMountDie{b}
}

func (*VolumeMountDieMut) DieFrom(b base.DieMutBase[VolumeMount, *VolumeMountDie, *VolumeMountDieMut]) *VolumeMountDieMut {
	return &VolumeMountDieMut{b}
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

func (d *VolumeMountDie) Name() die.Value[string] {
	return fields.ValueOf(d, volumeMountFields.Name)
}

func (d *VolumeMountDieMut) Name() die.ValueMut[string, *VolumeMountDieMut] {
	return fields.ValueMutOf(d, volumeMountFields.Name)
}

func (d *VolumeMountDie) ReadOnly() die.Value[bool] {
	return fields.ValueOf(d, volumeMountFields.ReadOnly)
}

func (d *VolumeMountDieMut) ReadOnly() die.ValueMut[bool, *VolumeMountDieMut] {
	return fields.ValueMutOf(d, volumeMountFields.ReadOnly)
}

func (d *VolumeMountDie) MountPath() die.Value[string] {
	return fields.ValueOf(d, volumeMountFields.MountPath)
}

func (d *VolumeMountDieMut) MountPath() die.ValueMut[string, *VolumeMountDieMut] {
	return fields.ValueMutOf(d, volumeMountFields.MountPath)
}

// Volume: the members of the inlined VolumeSource are accessed directly.

type VolumeDie struct {
	base.DieBase[Volume, *VolumeDie, *VolumeDieMut]
}

type VolumeDieMut struct {
	base.DieMutBase[Volume, *VolumeDie, *VolumeDieMut]
}

func (*VolumeDie) DieFrom(b base.DieBase[Volume, *VolumeDie, *VolumeDieMut]) *VolumeDie {
	return &VolumeDie{b}
}

func (*VolumeDieMut) DieFrom(b base.DieMutBase[Volume, *VolumeDie, *VolumeDieMut]) *VolumeDieMut {
	return &VolumeDieMut{b}
}

var volumeFields = struct {
	Name     fields.Field[Volume, string]
	Secret   fields.Struct[Volume, SecretVolumeSource, *SecretVolumeSourceDie, *SecretVolumeSourceDieMut]
	EmptyDir fields.Struct[Volume, EmptyDirVolumeSource, *EmptyDirVolumeSourceDie, *EmptyDirVolumeSourceDieMut]
}{
	Name:     fields.NewField[Volume, string]("name"),
	Secret:   fields.NewStruct[Volume, SecretVolumeSource, *SecretVolumeSourceDie, *SecretVolumeSourceDieMut]("secret"),
	EmptyDir: fields.NewStruct[Volume, EmptyDirVolumeSource, *EmptyDirVolumeSourceDie, *EmptyDirVolumeSourceDieMut]("emptyDir"),
}

func (d *VolumeDie) Name() die.Value[string] { return fields.ValueOf(d, volumeFields.Name) }

func (d *VolumeDieMut) Name() die.ValueMut[string, *VolumeDieMut] {
	return fields.ValueMutOf(d, volumeFields.Name)
}

func (d *VolumeDie) Secret() *SecretVolumeSourceDie { return fields.Child(d, volumeFields.Secret) }

func (d *VolumeDieMut) Secret() *SecretVolumeSourceDieMut {
	return fields.ChildMut(d, volumeFields.Secret)
}

func (d *VolumeDie) EmptyDir() *EmptyDirVolumeSourceDie {
	return fields.Child(d, volumeFields.EmptyDir)
}

func (d *VolumeDieMut) EmptyDir() *EmptyDirVolumeSourceDieMut {
	return fields.ChildMut(d, volumeFields.EmptyDir)
}

type SecretVolumeSourceDie struct {
	base.DieBase[SecretVolumeSource, *SecretVolumeSourceDie, *SecretVolumeSourceDieMut]
}

type SecretVolumeSourceDieMut struct {
	base.DieMutBase[SecretVolumeSource, *SecretVolumeSourceDie, *SecretVolumeSourceDieMut]
}

func (*SecretVolumeSourceDie) DieFrom(b base.DieBase[SecretVolumeSource, *SecretVolumeSourceDie, *SecretVolumeSourceDieMut]) *SecretVolumeSourceDie {
	return &SecretVolumeSourceDie{b}
}

func (*SecretVolumeSourceDieMut) DieFrom(b base.DieMutBase[SecretVolumeSource, *SecretVolumeSourceDie, *SecretVolumeSourceDieMut]) *SecretVolumeSourceDieMut {
	return &SecretVolumeSourceDieMut{b}
}

var secretVolumeSourceFields = struct {
	SecretName fields.Field[SecretVolumeSource, string]
}{
	SecretName: fields.NewField[SecretVolumeSource, string]("secretName"),
}

func (d *SecretVolumeSourceDie) SecretName() die.Value[string] {
	return fields.ValueOf(d, secretVolumeSourceFields.SecretName)
}

func (d *SecretVolumeSourceDieMut) SecretName() die.ValueMut[string, *SecretVolumeSourceDieMut] {
	return fields.ValueMutOf(d, secretVolumeSourceFields.SecretName)
}

type EmptyDirVolumeSourceDie struct {
	base.DieBase[EmptyDirVolumeSource, *EmptyDirVolumeSourceDie, *EmptyDirVolumeSourceDieMut]
}

type EmptyDirVolumeSourceDieMut struct {
	base.DieMutBase[EmptyDirVolumeSource, *EmptyDirVolumeSourceDie, *EmptyDirVolumeSourceDieMut]
}

func (*EmptyDirVolumeSourceDie) DieFrom(b base.DieBase[EmptyDirVolumeSource, *EmptyDirVolumeSourceDie, *EmptyDirVolumeSourceDieMut]) *EmptyDirVolumeSourceDie {
	return &EmptyDirVolumeSourceDie{b}
}

func (*EmptyDirVolumeSourceDieMut) DieFrom(b base.DieMutBase[EmptyDirVolumeSource, *EmptyDirVolumeSourceDie, *EmptyDirVolumeSourceDieMut]) *EmptyDirVolumeSourceDieMut {
	return &EmptyDirVolumeSourceDieMut{b}
}

var emptyDirVolumeSourceFields = struct {
	Medium fields.Field[EmptyDirVolumeSource, string]
}{
	Medium: fields.NewField[EmptyDirVolumeSource, string]("medium"),
}

func (d *EmptyDirVolumeSourceDie) Medium() die.Value[string] {
	return fields.ValueOf(d, emptyDirVolumeSourceFields.Medium)
}
func (d *EmptyDirVolumeSourceDieMut) Medium() die.ValueMut[string, *EmptyDirVolumeSourceDieMut] {
	return fields.ValueMutOf(d, emptyDirVolumeSourceFields.Medium)
}
