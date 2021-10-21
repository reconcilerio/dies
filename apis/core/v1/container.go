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
	"k8s.io/apimachinery/pkg/api/resource"
)

// +die
type _ = corev1.Container

type container interface {
	PortsDie(ports ...ContainerPortDie) ContainerDie
	EnvFromDie(prefix string, fn func(d EnvFromSourceDie)) ContainerDie
	EnvDie(name string, fn func(d EnvVarDie)) ContainerDie
	ResourcesDie(fn func(d ResourceRequirementsDie)) ContainerDie
	VolumeMountDie(name string, fn func(d VolumeMountDie)) ContainerDie
	LivenessProbeDie(fn func(d ProbeDie)) ContainerDie
	ReadinessProbeDie(fn func(d ProbeDie)) ContainerDie
	StartupProbeDie(fn func(d ProbeDie)) ContainerDie
	LifecycleDie(fn func(d LifecycleDie)) ContainerDie
	SecurityContextDie(fn func(d SecurityContextDie)) ContainerDie
}

func (d *containerDie) PortsDie(ports ...ContainerPortDie) ContainerDie {
	return d.DieStamp(func(r *corev1.Container) {
		r.Ports = make([]corev1.ContainerPort, len(ports))
		for i := range ports {
			r.Ports[i] = ports[i].DieRelease()
		}
	})
}

func (d *containerDie) EnvFromDie(prefix string, fn func(d EnvFromSourceDie)) ContainerDie {
	return d.DieStamp(func(r *corev1.Container) {
		for i := range r.EnvFrom {
			if prefix == r.EnvFrom[i].Prefix {
				d := EnvFromSourceBlank.DieImmutable(false).DieFeed(r.EnvFrom[i])
				fn(d)
				r.EnvFrom[i] = d.DieRelease()
				return
			}
		}

		d := EnvFromSourceBlank.DieImmutable(false).DieFeed(corev1.EnvFromSource{Prefix: prefix})
		fn(d)
		r.EnvFrom = append(r.EnvFrom, d.DieRelease())
	})
}

func (d *containerDie) EnvDie(name string, fn func(d EnvVarDie)) ContainerDie {
	return d.DieStamp(func(r *corev1.Container) {
		for i := range r.Env {
			if name == r.Env[i].Name {
				d := EnvVarBlank.DieImmutable(false).DieFeed(r.Env[i])
				fn(d)
				r.Env[i] = d.DieRelease()
				return
			}
		}

		d := EnvVarBlank.DieImmutable(false).DieFeed(corev1.EnvVar{Name: name})
		fn(d)
		r.Env = append(r.Env, d.DieRelease())
	})
}

func (d *containerDie) ResourcesDie(fn func(d ResourceRequirementsDie)) ContainerDie {
	return d.DieStamp(func(r *corev1.Container) {
		d := ResourceRequirementsBlank.DieImmutable(false).DieFeed(r.Resources)
		fn(d)
		r.Resources = d.DieRelease()
	})
}

func (d *containerDie) VolumeMountDie(name string, fn func(d VolumeMountDie)) ContainerDie {
	return d.DieStamp(func(r *corev1.Container) {
		for i := range r.VolumeMounts {
			if name == r.VolumeMounts[i].Name {
				d := VolumeMountBlank.DieImmutable(false).DieFeed(r.VolumeMounts[i])
				fn(d)
				r.VolumeMounts[i] = d.DieRelease()
				return
			}
		}

		d := VolumeMountBlank.DieImmutable(false).DieFeed(corev1.VolumeMount{Name: name})
		fn(d)
		r.VolumeMounts = append(r.VolumeMounts, d.DieRelease())
	})
}

func (d *containerDie) VolumeDeviceDie(name string, fn func(d VolumeDeviceDie)) ContainerDie {
	return d.DieStamp(func(r *corev1.Container) {
		for i := range r.VolumeDevices {
			if name == r.VolumeDevices[i].Name {
				d := VolumeDeviceBlank.DieImmutable(false).DieFeed(r.VolumeDevices[i])
				fn(d)
				r.VolumeDevices[i] = d.DieRelease()
				return
			}
		}

		d := VolumeDeviceBlank.DieImmutable(false).DieFeed(corev1.VolumeDevice{Name: name})
		fn(d)
		r.VolumeDevices = append(r.VolumeDevices, d.DieRelease())
	})
}

func (d *containerDie) LivenessProbeDie(fn func(d ProbeDie)) ContainerDie {
	return d.DieStamp(func(r *corev1.Container) {
		d := ProbeBlank.DieImmutable(false).DieFeedPtr(r.LivenessProbe)
		fn(d)
		r.LivenessProbe = d.DieReleasePtr()
	})
}

func (d *containerDie) ReadinessProbeDie(fn func(d ProbeDie)) ContainerDie {
	return d.DieStamp(func(r *corev1.Container) {
		d := ProbeBlank.DieImmutable(false).DieFeedPtr(r.ReadinessProbe)
		fn(d)
		r.ReadinessProbe = d.DieReleasePtr()
	})
}

func (d *containerDie) StartupProbeDie(fn func(d ProbeDie)) ContainerDie {
	return d.DieStamp(func(r *corev1.Container) {
		d := ProbeBlank.DieImmutable(false).DieFeedPtr(r.StartupProbe)
		fn(d)
		r.StartupProbe = d.DieReleasePtr()
	})
}

func (d *containerDie) LifecycleDie(fn func(d LifecycleDie)) ContainerDie {
	return d.DieStamp(func(r *corev1.Container) {
		d := LifecycleBlank.DieImmutable(false).DieFeedPtr(r.Lifecycle)
		fn(d)
		r.Lifecycle = d.DieReleasePtr()
	})
}

func (d *containerDie) SecurityContextDie(fn func(d SecurityContextDie)) ContainerDie {
	return d.DieStamp(func(r *corev1.Container) {
		d := SecurityContextBlank.DieImmutable(false).DieFeedPtr(r.SecurityContext)
		fn(d)
		r.SecurityContext = d.DieReleasePtr()
	})
}

// +die
type _ = corev1.ContainerPort

// +die
type _ = corev1.EnvFromSource

type envFromSource interface {
	ConfigMapRefDie(fn func(d ConfigMapEnvSourceDie)) EnvFromSourceDie
	SecretRefDie(fn func(d SecretEnvSourceDie)) EnvFromSourceDie
}

func (d *envFromSourceDie) ConfigMapRefDie(fn func(d ConfigMapEnvSourceDie)) EnvFromSourceDie {
	return d.DieStamp(func(r *corev1.EnvFromSource) {
		d := ConfigMapEnvSourceBlank.DieImmutable(false).DieFeedPtr(r.ConfigMapRef)
		fn(d)
		r.ConfigMapRef = d.DieReleasePtr()
	})
}

func (d *envFromSourceDie) SecretRefDie(fn func(d SecretEnvSourceDie)) EnvFromSourceDie {
	return d.DieStamp(func(r *corev1.EnvFromSource) {
		d := SecretEnvSourceBlank.DieImmutable(false).DieFeedPtr(r.SecretRef)
		fn(d)
		r.SecretRef = d.DieReleasePtr()
	})
}

// +die
type _ = corev1.ConfigMapEnvSource

type configMapEnvSource interface {
	Name(v string) ConfigMapEnvSourceDie
}

func (d *configMapEnvSourceDie) Name(v string) ConfigMapEnvSourceDie {
	return d.DieStamp(func(r *corev1.ConfigMapEnvSource) {
		r.Name = v
	})
}

// +die
type _ = corev1.SecretEnvSource

type secretEnvSource interface {
	Name(v string) SecretEnvSourceDie
}

func (d *secretEnvSourceDie) Name(v string) SecretEnvSourceDie {
	return d.DieStamp(func(r *corev1.SecretEnvSource) {
		r.Name = v
	})
}

// +die
type _ = corev1.EnvVar

type envVar interface {
	ValueFromDie(fn func(d EnvVarSourceDie)) EnvVarDie
}

func (d *envVarDie) ValueFromDie(fn func(d EnvVarSourceDie)) EnvVarDie {
	return d.DieStamp(func(r *corev1.EnvVar) {
		d := EnvVarSourceBlank.DieImmutable(false).DieFeedPtr(r.ValueFrom)
		fn(d)
		r.ValueFrom = d.DieReleasePtr()
	})
}

// +die
type _ = corev1.EnvVarSource

type envVarSource interface {
	FieldRefDie(fn func(ObjectFieldSelectorDie)) EnvVarSourceDie
	ResourceFieldRefDie(fn func(ResourceFieldSelectorDie)) EnvVarSourceDie
	ConfigMapKeyRefDie(fn func(ConfigMapKeySelectorDie)) EnvVarSourceDie
	SecretKeyRefDie(fn func(SecretKeySelectorDie)) EnvVarSourceDie
}

func (d *envVarSourceDie) FieldRefDie(fn func(ObjectFieldSelectorDie)) EnvVarSourceDie {
	return d.DieStamp(func(r *corev1.EnvVarSource) {
		d := ObjectFieldSelectorBlank.DieImmutable(false).DieFeedPtr(r.FieldRef)
		fn(d)
		r.FieldRef = d.DieReleasePtr()
	})
}

func (d *envVarSourceDie) ResourceFieldRefDie(fn func(ResourceFieldSelectorDie)) EnvVarSourceDie {
	return d.DieStamp(func(r *corev1.EnvVarSource) {
		d := ResourceFieldSelectorBlank.DieImmutable(false).DieFeedPtr(r.ResourceFieldRef)
		fn(d)
		r.ResourceFieldRef = d.DieReleasePtr()
	})
}

func (d *envVarSourceDie) ConfigMapKeyRefDie(fn func(ConfigMapKeySelectorDie)) EnvVarSourceDie {
	return d.DieStamp(func(r *corev1.EnvVarSource) {
		d := ConfigMapKeySelectorBlank.DieImmutable(false).DieFeedPtr(r.ConfigMapKeyRef)
		fn(d)
		r.ConfigMapKeyRef = d.DieReleasePtr()
	})
}

func (d *envVarSourceDie) SecretKeyRefDie(fn func(SecretKeySelectorDie)) EnvVarSourceDie {
	return d.DieStamp(func(r *corev1.EnvVarSource) {
		d := SecretKeySelectorBlank.DieImmutable(false).DieFeedPtr(r.SecretKeyRef)
		fn(d)
		r.SecretKeyRef = d.DieReleasePtr()
	})
}

// +die
type _ = corev1.ObjectFieldSelector

// +die
type _ = corev1.ResourceFieldSelector

// +die
type _ = corev1.ConfigMapKeySelector

type configMapKeySelector interface {
	Name(v string) ConfigMapKeySelectorDie
}

func (d *configMapKeySelectorDie) Name(v string) ConfigMapKeySelectorDie {
	return d.DieStamp(func(r *corev1.ConfigMapKeySelector) {
		r.Name = v
	})
}

// +die
type _ = corev1.SecretKeySelector

type secretKeySelector interface {
	Name(v string) SecretKeySelectorDie
}

func (d *secretKeySelectorDie) Name(v string) SecretKeySelectorDie {
	return d.DieStamp(func(r *corev1.SecretKeySelector) {
		r.Name = v
	})
}

// +die
type _ = corev1.ResourceRequirements

type resourceRequirements interface {
	AddLimit(name corev1.ResourceName, quantity resource.Quantity) ResourceRequirementsDie
	AddLimitString(name corev1.ResourceName, quantity string) ResourceRequirementsDie
	AddRequest(name corev1.ResourceName, quantity resource.Quantity) ResourceRequirementsDie
	AddRequestString(name corev1.ResourceName, quantity string) ResourceRequirementsDie
}

func (d *resourceRequirementsDie) AddLimit(name corev1.ResourceName, quantity resource.Quantity) ResourceRequirementsDie {
	return d.DieStamp(func(r *corev1.ResourceRequirements) {
		r.Limits[name] = quantity
	})
}

func (d *resourceRequirementsDie) AddLimitString(name corev1.ResourceName, quantity string) ResourceRequirementsDie {
	return d.AddLimit(name, resource.MustParse(quantity))
}

func (d *resourceRequirementsDie) AddRequest(name corev1.ResourceName, quantity resource.Quantity) ResourceRequirementsDie {
	return d.DieStamp(func(r *corev1.ResourceRequirements) {
		r.Requests[name] = quantity
	})
}

func (d *resourceRequirementsDie) AddRequestString(name corev1.ResourceName, quantity string) ResourceRequirementsDie {
	return d.AddRequest(name, resource.MustParse(quantity))
}

// +die
type _ = corev1.VolumeMount

// +die
type _ = corev1.VolumeDevice

// +die
type _ = corev1.Probe

type probe interface {
	HandlerDie(fn func(d HandlerDie)) ProbeDie
	ExecDie(fn func(d ExecActionDie)) ProbeDie
	HTTPGetDie(fn func(d HTTPGetActionDie)) ProbeDie
	TCPSocketDie(fn func(d TCPSocketActionDie)) ProbeDie
}

func (d *probeDie) HandlerDie(fn func(d HandlerDie)) ProbeDie {
	return d.DieStamp(func(r *corev1.Probe) {
		d := HandlerBlank.DieImmutable(false).DieFeed(r.Handler)
		fn(d)
		r.Handler = d.DieRelease()
	})
}

func (d *probeDie) ExecDie(fn func(d ExecActionDie)) ProbeDie {
	return d.DieStamp(func(r *corev1.Probe) {
		d := ExecActionBlank.DieImmutable(false).DieFeedPtr(r.Exec)
		fn(d)
		r.Handler = corev1.Handler{
			Exec: d.DieReleasePtr(),
		}
	})
}

func (d *probeDie) HTTPGetDie(fn func(d HTTPGetActionDie)) ProbeDie {
	return d.DieStamp(func(r *corev1.Probe) {
		d := HTTPGetActionBlank.DieImmutable(false).DieFeedPtr(r.HTTPGet)
		fn(d)
		r.Handler = corev1.Handler{
			HTTPGet: d.DieReleasePtr(),
		}
	})
}

func (d *probeDie) TCPSocketDie(fn func(d TCPSocketActionDie)) ProbeDie {
	return d.DieStamp(func(r *corev1.Probe) {
		d := TCPSocketActionBlank.DieImmutable(false).DieFeedPtr(r.TCPSocket)
		fn(d)
		r.Handler = corev1.Handler{
			TCPSocket: d.DieReleasePtr(),
		}
	})
}

// +die
type _ = corev1.Lifecycle

type lifecycle interface {
	PostStartDie(fn func(d HandlerDie)) LifecycleDie
	PreStopDie(fn func(d HandlerDie)) LifecycleDie
}

func (d *lifecycleDie) PostStartDie(fn func(d HandlerDie)) LifecycleDie {
	return d.DieStamp(func(r *corev1.Lifecycle) {
		d := HandlerBlank.DieImmutable(false).DieFeedPtr(r.PostStart)
		fn(d)
		r.PostStart = d.DieReleasePtr()
	})
}

func (d *lifecycleDie) PreStopDie(fn func(d HandlerDie)) LifecycleDie {
	return d.DieStamp(func(r *corev1.Lifecycle) {
		d := HandlerBlank.DieImmutable(false).DieFeedPtr(r.PreStop)
		fn(d)
		r.PreStop = d.DieReleasePtr()
	})
}

// +die
type _ = corev1.Handler

type handler interface {
	ExecDie(fn func(d ExecActionDie)) HandlerDie
	HTTPGetDie(fn func(d HTTPGetActionDie)) HandlerDie
	TCPSocketDie(fn func(d TCPSocketActionDie)) HandlerDie
}

func (d *handlerDie) ExecDie(fn func(d ExecActionDie)) HandlerDie {
	return d.DieStamp(func(r *corev1.Handler) {
		d := ExecActionBlank.DieImmutable(false).DieFeedPtr(r.Exec)
		fn(d)
		r.Exec = d.DieReleasePtr()
	})
}

func (d *handlerDie) HTTPGetDie(fn func(d HTTPGetActionDie)) HandlerDie {
	return d.DieStamp(func(r *corev1.Handler) {
		d := HTTPGetActionBlank.DieImmutable(false).DieFeedPtr(r.HTTPGet)
		fn(d)
		r.HTTPGet = d.DieReleasePtr()
	})
}

func (d *handlerDie) TCPSocketDie(fn func(d TCPSocketActionDie)) HandlerDie {
	return d.DieStamp(func(r *corev1.Handler) {
		d := TCPSocketActionBlank.DieImmutable(false).DieFeedPtr(r.TCPSocket)
		fn(d)
		r.TCPSocket = d.DieReleasePtr()
	})
}

// +die
type _ = corev1.ExecAction

// +die
type _ = corev1.HTTPGetAction

type hTTPGetAction interface {
	HTTPHeadersDie(headers ...HTTPHeaderDie) HTTPGetActionDie
}

func (d *hTTPGetActionDie) HTTPHeadersDie(headers ...HTTPHeaderDie) HTTPGetActionDie {
	return d.DieStamp(func(r *corev1.HTTPGetAction) {
		r.HTTPHeaders = make([]corev1.HTTPHeader, len(headers))
		for i := range headers {
			r.HTTPHeaders[i] = headers[i].DieRelease()
		}
	})
}

// +die
type _ = corev1.HTTPHeader

// +die
type _ = corev1.TCPSocketAction

// +die
type _ = corev1.SecurityContext

type securityContext interface {
	CapabilitiesDie(fn func(d CapabilitiesDie)) SecurityContextDie
	SELinuxOptionsDie(fn func(d SELinuxOptionsDie)) SecurityContextDie
	WindowsOptionsDie(fn func(d WindowsSecurityContextOptionsDie)) SecurityContextDie
	SeccompProfileDie(fn func(d SeccompProfileDie)) SecurityContextDie
}

func (d *securityContextDie) CapabilitiesDie(fn func(d CapabilitiesDie)) SecurityContextDie {
	return d.DieStamp(func(r *corev1.SecurityContext) {
		d := CapabilitiesBlank.DieImmutable(false).DieFeedPtr(r.Capabilities)
		fn(d)
		r.Capabilities = d.DieReleasePtr()
	})
}

func (d *securityContextDie) SELinuxOptionsDie(fn func(d SELinuxOptionsDie)) SecurityContextDie {
	return d.DieStamp(func(r *corev1.SecurityContext) {
		d := SELinuxOptionsBlank.DieImmutable(false).DieFeedPtr(r.SELinuxOptions)
		fn(d)
		r.SELinuxOptions = d.DieReleasePtr()
	})
}

func (d *securityContextDie) WindowsOptionsDie(fn func(d WindowsSecurityContextOptionsDie)) SecurityContextDie {
	return d.DieStamp(func(r *corev1.SecurityContext) {
		d := WindowsSecurityContextOptionsBlank.DieImmutable(false).DieFeedPtr(r.WindowsOptions)
		fn(d)
		r.WindowsOptions = d.DieReleasePtr()
	})
}

func (d *securityContextDie) SeccompProfileDie(fn func(d SeccompProfileDie)) SecurityContextDie {
	return d.DieStamp(func(r *corev1.SecurityContext) {
		d := SeccompProfileBlank.DieImmutable(false).DieFeedPtr(r.SeccompProfile)
		fn(d)
		r.SeccompProfile = d.DieReleasePtr()
	})
}

// +die
type _ = corev1.Capabilities

// +die
type _ = corev1.SELinuxOptions

// +die
type _ = corev1.WindowsSecurityContextOptions

// +die
type _ = corev1.SeccompProfile

// +die
type _ = corev1.ContainerStatus

type containerStatus interface {
	StateDie(fn func(d ContainerStateDie)) ContainerStatusDie
	LastTerminationStateDie(fn func(d ContainerStateDie)) ContainerStatusDie
}

func (d *containerStatusDie) StateDie(fn func(d ContainerStateDie)) ContainerStatusDie {
	return d.DieStamp(func(r *corev1.ContainerStatus) {
		d := ContainerStateBlank.DieImmutable(false).DieFeed(r.State)
		fn(d)
		r.State = d.DieRelease()
	})
}

func (d *containerStatusDie) LastTerminationStateDie(fn func(d ContainerStateDie)) ContainerStatusDie {
	return d.DieStamp(func(r *corev1.ContainerStatus) {
		d := ContainerStateBlank.DieImmutable(false).DieFeed(r.LastTerminationState)
		fn(d)
		r.LastTerminationState = d.DieRelease()
	})
}

// +die
type _ = corev1.ContainerState

type containerState interface {
	WaitingDie(fn func(d ContainerStateWaitingDie)) ContainerStateDie
	RunningDie(fn func(d ContainerStateRunningDie)) ContainerStateDie
	TerminatedDie(fn func(d ContainerStateTerminatedDie)) ContainerStateDie
}

func (d *containerStateDie) WaitingDie(fn func(d ContainerStateWaitingDie)) ContainerStateDie {
	return d.DieStamp(func(r *corev1.ContainerState) {
		d := ContainerStateWaitingBlank.DieImmutable(false).DieFeedPtr(r.Waiting)
		fn(d)
		r.Waiting = d.DieReleasePtr()
	})
}

func (d *containerStateDie) RunningDie(fn func(d ContainerStateRunningDie)) ContainerStateDie {
	return d.DieStamp(func(r *corev1.ContainerState) {
		d := ContainerStateRunningBlank.DieImmutable(false).DieFeedPtr(r.Running)
		fn(d)
		r.Running = d.DieReleasePtr()
	})
}

func (d *containerStateDie) TerminatedDie(fn func(d ContainerStateTerminatedDie)) ContainerStateDie {
	return d.DieStamp(func(r *corev1.ContainerState) {
		d := ContainerStateTerminatedBlank.DieImmutable(false).DieFeedPtr(r.Terminated)
		fn(d)
		r.Terminated = d.DieReleasePtr()
	})
}

// +die
type _ = corev1.ContainerStateWaiting

// +die
type _ = corev1.ContainerStateRunning

// +die
type _ = corev1.ContainerStateTerminated
