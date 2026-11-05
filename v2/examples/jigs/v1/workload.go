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

package v1

import (
	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	diecorev1 "reconciler.io/dies/v2/apis/core/v1"
	"reconciler.io/dies/v2/die/jigs"
)

// WorkloadJig presents workloads as their containers and volumes. A
// Deployment's containers and init containers are one list: changing a
// container goes back to the list it is in, but adding one is ambiguous.
var WorkloadJig = jigs.NewJig(WorkloadBlank).
	MustClamp(dieappsv1.DeploymentBlank,
		jigs.BindAll(WorkloadBlank.Containers(),
			jigs.Each(dieappsv1.DeploymentBlank.Spec().Template().Spec().Containers()),
			jigs.Each(dieappsv1.DeploymentBlank.Spec().Template().Spec().InitContainers()),
		),
		jigs.Bind(WorkloadBlank.Volumes(), dieappsv1.DeploymentBlank.Spec().Template().Spec().Volumes()),
	).
	// any unstructured workload, by the paths ServiceBinding uses by default
	MustClampUnknownKinds(
		jigs.BindAll(WorkloadBlank.Containers(),
			jigs.EachAt[diecorev1.Container](".spec.template.spec.containers"),
			jigs.EachAt[diecorev1.Container](".spec.template.spec.initContainers"),
		),
		jigs.BindPath(WorkloadBlank.Volumes(), ".spec.template.spec.volumes"),
	).
	// a pipeline's steps are its only list of containers, so a new
	// container is added to them
	// a batch's containers are an atomic list: its elements are found by
	// index, and a change is written back as the whole list
	MustClamp(BatchBlank,
		jigs.BindAll(WorkloadBlank.Containers(),
			jigs.Each(BatchBlank.Containers()),
			jigs.One(BatchBlank.Init()),
		),
	).
	MustClamp(PipelineBlank,
		jigs.BindAll(WorkloadBlank.Containers(),
			jigs.Each(PipelineBlank.Steps()),
			jigs.One(PipelineBlank.Sidecar()),
			jigs.OneVia(PipelineBlank.Builder(), BuilderContainerJig),
		),
	)

// BuilderContainerJig presents a [Builder] as a Container.
var BuilderContainerJig = jigs.NewJig(diecorev1.ContainerBlank).
	MustClamp(BuilderBlank,
		jigs.Bind(diecorev1.ContainerBlank.Name(), BuilderBlank.Name()),
		jigs.Bind(diecorev1.ContainerBlank.Image(), BuilderBlank.Image()),
		jigs.Bind(diecorev1.ContainerBlank.Env(), BuilderBlank.Env()),
	)
