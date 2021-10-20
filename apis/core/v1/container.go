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
type _ = corev1.Container

type container interface {
	AddEnv(env ...corev1.EnvVar) ContainerDie
	AddEnvFrom(env ...corev1.EnvFromSource) ContainerDie
	AddVolumeMounts(volumeMounts ...corev1.VolumeMount) ContainerDie
}

func (d *containerDie) AddEnv(env ...corev1.EnvVar) ContainerDie {
	return d.DieStamp(func(r *corev1.Container) {
		for _, e := range env {
			found := false
			for i := range r.Env {
				if e.Name == r.Env[i].Name {
					found = true
					r.Env[i] = e
				}
			}
			if !found {
				r.Env = append(r.Env, e)
			}
		}
	})
}

func (d *containerDie) AddEnvFrom(env ...corev1.EnvFromSource) ContainerDie {
	return d.DieStamp(func(r *corev1.Container) {
		for _, e := range env {
			found := false
			for i := range r.EnvFrom {
				if e.Prefix == r.EnvFrom[i].Prefix {
					found = true
					r.EnvFrom[i] = e
				}
			}
			if !found {
				r.EnvFrom = append(r.EnvFrom, e)
			}
		}
	})
}

func (d *containerDie) AddVolumeMounts(volumeMounts ...corev1.VolumeMount) ContainerDie {
	return d.DieStamp(func(r *corev1.Container) {
		for _, m := range volumeMounts {
			found := false
			for i := range r.VolumeMounts {
				if m.Name == r.VolumeMounts[i].Name {
					found = true
					r.VolumeMounts[i] = m
				}
			}
			if !found {
				r.VolumeMounts = append(r.VolumeMounts, m)
			}
		}
	})
}
