/*
Copyright 2021 the original author or authors.

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

package v1

import (
	corev1 "k8s.io/api/core/v1"
)

// +die
// +die:field:name=Resources,die=ResourceRequirementsDie
// +die:field:name=LivenessProbe,die=ProbeDie,pointer=true
// +die:field:name=ReadinessProbe,die=ProbeDie,pointer=true
// +die:field:name=StartupProbe,die=ProbeDie,pointer=true
// +die:field:name=Lifecycle,die=LifecycleDie,pointer=true
// +die:field:name=SecurityContext,die=SecurityContextDie,pointer=true
// +die:field:name=Ports,die=ContainerPortDie,listType=atomic
// +die:field:name=EnvFrom,die=EnvFromSourceDie,listType=map,listMapKey=Prefix
// +die:field:name=Env,die=EnvVarDie,listType=map
// +die:field:name=ResizePolicy,die=ContainerResizePolicyDie,listType=map,listMapKey=ResourceName,listMapKeyPackage=k8s.io/api/core/v1,listMapKeyType=ResourceName
// +die:field:name=VolumeMounts,die=VolumeMountDie,listType=map
// +die:field:name=VolumeDevices,die=VolumeDeviceDie,listType=map
type _ = corev1.Container

// +die
type _ = corev1.ContainerPort

// +die
// +die:field:name=ConfigMapRef,die=ConfigMapEnvSourceDie,pointer=true
// +die:field:name=SecretRef,die=SecretEnvSourceDie,pointer=true
type _ = corev1.EnvFromSource

// +die
type _ = corev1.ConfigMapEnvSource

func (d *ConfigMapEnvSourceDie) Name(v string) *ConfigMapEnvSourceDie {
	return d.DieStamp(func(r *corev1.ConfigMapEnvSource) {
		r.Name = v
	})
}

// +die
type _ = corev1.SecretEnvSource

func (d *SecretEnvSourceDie) Name(v string) *SecretEnvSourceDie {
	return d.DieStamp(func(r *corev1.SecretEnvSource) {
		r.Name = v
	})
}

// +die
// +die:field:name=ValueFrom,die=EnvVarSourceDie,pointer=true
type _ = corev1.EnvVar

// +die
// +die:field:name=FieldRef,die=ObjectFieldSelectorDie,pointer=true
// +die:field:name=ResourceFieldRef,die=ResourceFieldSelectorDie,pointer=true
// +die:field:name=ConfigMapKeyRef,die=ConfigMapKeySelectorDie,pointer=true
// +die:field:name=SecretKeyRef,die=SecretKeySelectorDie,pointer=true
type _ = corev1.EnvVarSource

// +die
type _ = corev1.ObjectFieldSelector

// +die
type _ = corev1.ResourceFieldSelector

// +die
type _ = corev1.ConfigMapKeySelector

func (d *ConfigMapKeySelectorDie) Name(v string) *ConfigMapKeySelectorDie {
	return d.DieStamp(func(r *corev1.ConfigMapKeySelector) {
		r.Name = v
	})
}

// +die
type _ = corev1.SecretKeySelector

func (d *SecretKeySelectorDie) Name(v string) *SecretKeySelectorDie {
	return d.DieStamp(func(r *corev1.SecretKeySelector) {
		r.Name = v
	})
}

// +die
// +die:field:name=Claims,die=ResourceClaimDie,listType=atomic
// +die:field:name=Claims,die=ResourceClaimDie,listType=map
type _ = corev1.ResourceRequirements

// +die
type _ = corev1.ResourceClaim

// +die
type _ = corev1.ContainerResizePolicy

// +die
type _ = corev1.VolumeMount

// +die
type _ = corev1.VolumeDevice

// +die
// +die:field:name=ProbeHandler,die=ProbeHandlerDie
type _ = corev1.Probe

func (d *ProbeDie) ExecDie(fn func(d *ExecActionDie)) *ProbeDie {
	return d.DieStamp(func(r *corev1.Probe) {
		d := ExecActionBlank.DieImmutable(false).DieFeedPtr(r.Exec)
		fn(d)
		r.ProbeHandler = corev1.ProbeHandler{
			Exec: d.DieReleasePtr(),
		}
	})
}

func (d *ProbeDie) HTTPGetDie(fn func(d *HTTPGetActionDie)) *ProbeDie {
	return d.DieStamp(func(r *corev1.Probe) {
		d := HTTPGetActionBlank.DieImmutable(false).DieFeedPtr(r.HTTPGet)
		fn(d)
		r.ProbeHandler = corev1.ProbeHandler{
			HTTPGet: d.DieReleasePtr(),
		}
	})
}

func (d *ProbeDie) TCPSocketDie(fn func(d *TCPSocketActionDie)) *ProbeDie {
	return d.DieStamp(func(r *corev1.Probe) {
		d := TCPSocketActionBlank.DieImmutable(false).DieFeedPtr(r.TCPSocket)
		fn(d)
		r.ProbeHandler = corev1.ProbeHandler{
			TCPSocket: d.DieReleasePtr(),
		}
	})
}

// +die
// +die:field:name=PostStart,die=LifecycleHandlerDie,pointer=true
// +die:field:name=PreStop,die=LifecycleHandlerDie,pointer=true
type _ = corev1.Lifecycle

// +die
// +die:field:name=Exec,die=ExecActionDie,pointer=true
// +die:field:name=HTTPGet,die=HTTPGetActionDie,pointer=true
// +die:field:name=TCPSocket,die=TCPSocketActionDie,pointer=true
// +die:field:name=Sleep,die=SleepActionDie,pointer=true
type _ = corev1.LifecycleHandler

// +die
// +die:field:name=Exec,die=ExecActionDie,pointer=true
// +die:field:name=HTTPGet,die=HTTPGetActionDie,pointer=true
// +die:field:name=TCPSocket,die=TCPSocketActionDie,pointer=true
// +die:field:name=GRPC,die=GRPCActionDie,pointer=true
type _ = corev1.ProbeHandler

// +die
type _ = corev1.ExecAction

// +die
// +die:field:name=HTTPHeaders,die=HTTPHeaderDie,listType=atomic
type _ = corev1.HTTPGetAction

// +die
type _ = corev1.HTTPHeader

// +die
type _ = corev1.TCPSocketAction

// +die
type _ = corev1.GRPCAction

// +die
type _ = corev1.SleepAction

// +die
// +die:field:name=Capabilities,die=CapabilitiesDie,pointer=true
// +die:field:name=SELinuxOptions,die=SELinuxOptionsDie,pointer=true
// +die:field:name=WindowsOptions,die=WindowsSecurityContextOptionsDie,pointer=true
// +die:field:name=SeccompProfile,die=SeccompProfileDie,pointer=true
// +die:field:name=AppArmorProfile,die=AppArmorProfileDie,pointer=true
type _ = corev1.SecurityContext

// +die
type _ = corev1.Capabilities

// +die
type _ = corev1.SELinuxOptions

// +die
type _ = corev1.WindowsSecurityContextOptions

// +die
type _ = corev1.SeccompProfile

// +die
type _ = corev1.AppArmorProfile

// +die
// +die:field:name=State,die=ContainerStateDie
// +die:field:name=LastTerminationState,die=ContainerStateDie
// +die:field:name=Resources,die=ResourceRequirementsDie,pointer=true
// +die:field:name=VolumeMounts,die=VolumeMountStatusDie,listType=map
// +die:field:name=User,die=ContainerUserDie,pointer=true
// +die:field:name=AllocatedResourcesStatus,die=ResourceStatusDie,listType=map,listMapKeyPackage=k8s.io/api/core/v1,listMapKeyType=ResourceName
type _ = corev1.ContainerStatus

// +die
// +die:field:name=Linux,die=LinuxContainerUserDie,pointer=true
type _ = corev1.ContainerUser

// +die
type _ = corev1.LinuxContainerUser

// +die
// +die:field:name=Resources,die=ResourceHealthDie,listType=map,listMapKeyPackage=k8s.io/api/core/v1,listMapKeyType=ResourceID,listMapKey=ResourceID
type _ = corev1.ResourceStatus

// +die
type _ = corev1.ResourceHealth

// +die
// +die:field:name=Waiting,die=ContainerStateWaitingDie,pointer=true
// +die:field:name=Running,die=ContainerStateRunningDie,pointer=true
// +die:field:name=Terminated,die=ContainerStateTerminatedDie,pointer=true
type _ = corev1.ContainerState

// +die
type _ = corev1.ContainerStateWaiting

// +die
type _ = corev1.ContainerStateRunning

// +die
type _ = corev1.ContainerStateTerminated

// +die
type _ = corev1.VolumeMountStatus
