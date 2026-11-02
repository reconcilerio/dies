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
	"k8s.io/apimachinery/pkg/runtime/schema"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	diebatchv1 "reconciler.io/dies/v2/apis/batch/v1"
	diecorev1 "reconciler.io/dies/v2/apis/core/v1"
	"reconciler.io/dies/v2/die/jigs"
)

// PodSpecableJig presents workloads as [PodSpecable]. Each binding pairs a
// location of the jig with a location of the source, both navigated from
// blanks; binding locations of different types does not compile.
var PodSpecableJig = jigs.NewJig[PodSpecable, *PodSpecableDie]().
	For(dieappsv1.DeploymentBlank,
		jigs.Bind(PodSpecableBlank.Metadata(), dieappsv1.DeploymentBlank.Metadata()),
		jigs.Bind(PodSpecableBlank.Template(), dieappsv1.DeploymentBlank.Spec().Template()),
	).
	For(diebatchv1.CronJobBlank,
		jigs.Bind(PodSpecableBlank.Metadata(), diebatchv1.CronJobBlank.Metadata()),
		jigs.Bind(PodSpecableBlank.Template(), diebatchv1.CronJobBlank.Spec().JobTemplate().Spec().Template()),
	).
	// a pod has no template; the jig's template is assembled from the pod's
	// metadata and spec
	For(diecorev1.PodBlank,
		jigs.Bind(PodSpecableBlank.Metadata(), diecorev1.PodBlank.Metadata()),
		jigs.Bind(PodSpecableBlank.Template().Metadata(), diecorev1.PodBlank.Metadata()),
		jigs.Bind(PodSpecableBlank.Template().Spec(), diecorev1.PodBlank.Spec()),
	).
	// unstructured objects, of kinds that may have no Go types at all: a
	// CronJob's template is deeper, and any other kind is assumed to be
	// PodSpec-able, as ServiceBinding does
	ForKind(schema.GroupVersionKind{Group: "batch", Version: "v1", Kind: "CronJob"},
		jigs.BindPath(PodSpecableBlank.Metadata(), ".metadata"),
		jigs.BindPath(PodSpecableBlank.Template(), ".spec.jobTemplate.spec.template"),
	).
	ForUnknownKinds(
		jigs.BindPath(PodSpecableBlank.Metadata(), ".metadata"),
		jigs.BindPath(PodSpecableBlank.Template(), ".spec.template"),
	)
