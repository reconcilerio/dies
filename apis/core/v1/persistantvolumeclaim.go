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
	diemetav1 "reconciler.io/dies/apis/meta/v1"
)

// +die:object=true,apiVersion=v1,kind=PersistentVolumeClaim
type _ = corev1.PersistentVolumeClaim

// +die
// +die:field:name=Selector,package=_/meta/v1,die=LabelSelectorDie,pointer=true
// +die:field:name=Resources,die=VolumeResourceRequirementsDie
// +die:field:name=DataSource,die=TypedLocalObjectReferenceDie,pointer=true
// +die:field:name=DataSourceRef,die=TypedObjectReferenceDie,pointer=true
type _ = corev1.PersistentVolumeClaimSpec

// +die
type _ = corev1.VolumeResourceRequirements

// +die:ignore={AllocatedResourceStatuses}
// +die:field:name=ModifyVolumeStatus,die=ModifyVolumeStatusDie,pointer=true
type _ = corev1.PersistentVolumeClaimStatus

func (d *PersistentVolumeClaimStatusDie) ConditionsDie(conditions ...*diemetav1.ConditionDie) *PersistentVolumeClaimStatusDie {
	return d.DieStamp(func(r *corev1.PersistentVolumeClaimStatus) {
		r.Conditions = make([]corev1.PersistentVolumeClaimCondition, len(conditions))
		for i := range conditions {
			c := conditions[i].DieRelease()
			r.Conditions[i] = corev1.PersistentVolumeClaimCondition{
				Type:               corev1.PersistentVolumeClaimConditionType(c.Type),
				Status:             corev1.ConditionStatus(c.Status),
				Reason:             c.Reason,
				Message:            c.Message,
				LastTransitionTime: c.LastTransitionTime,
			}
		}
	})
}

// allocatedResourceStatuses stores status of resource being resized for the given PVC.
// Key names follow standard Kubernetes label syntax. Valid values are either:
//   - Un-prefixed keys:
//   - storage - the capacity of the volume.
//   - Custom resources must use implementation-defined prefixed names such as "example.com/my-custom-resource"
//
// Apart from above values - keys that are unprefixed or have kubernetes.io prefix are considered
// reserved and hence may not be used.
//
// ClaimResourceStatus can be in any of following states:
//   - ControllerResizeInProgress:
//     State set when resize controller starts resizing the volume in control-plane.
//   - ControllerResizeFailed:
//     State set when resize has failed in resize controller with a terminal error.
//   - NodeResizePending:
//     State set when resize controller has finished resizing the volume but further resizing of
//     volume is needed on the node.
//   - NodeResizeInProgress:
//     State set when kubelet starts resizing the volume.
//   - NodeResizeFailed:
//     State set when resizing has failed in kubelet with a terminal error. Transient errors don't set
//     NodeResizeFailed.
//
// For example: if expanding a PVC for more capacity - this field can be one of the following states:
//   - pvc.status.allocatedResourceStatus['storage'] = "ControllerResizeInProgress"
//   - pvc.status.allocatedResourceStatus['storage'] = "ControllerResizeFailed"
//   - pvc.status.allocatedResourceStatus['storage'] = "NodeResizePending"
//   - pvc.status.allocatedResourceStatus['storage'] = "NodeResizeInProgress"
//   - pvc.status.allocatedResourceStatus['storage'] = "NodeResizeFailed"
//
// When this field is not set, it means that no resize operation is in progress for the given PVC.
//
// A controller that receives PVC update with previously unknown resourceName or ClaimResourceStatus
// should ignore the update for the purpose it was designed. For example - a controller that
// only is responsible for resizing capacity of the volume, should ignore PVC updates that change other valid
// resources associated with PVC.
//
// This is an alpha field and requires enabling RecoverVolumeExpansionFailure feature.
func (d *PersistentVolumeClaimStatusDie) AllocatedResourceStatuses(v map[corev1.ResourceName]corev1.ClaimResourceStatus) *PersistentVolumeClaimStatusDie {
	return d.DieStamp(func(r *corev1.PersistentVolumeClaimStatus) {
		r.AllocatedResourceStatuses = v
	})
}

func (d *PersistentVolumeClaimStatusDie) AddAllocatedResourceStatus(name corev1.ResourceName, status corev1.ClaimResourceStatus) *PersistentVolumeClaimStatusDie {
	return d.DieStamp(func(r *corev1.PersistentVolumeClaimStatus) {
		if r.AllocatedResourceStatuses == nil {
			r.AllocatedResourceStatuses = map[corev1.ResourceName]corev1.ClaimResourceStatus{}
		}
		r.AllocatedResourceStatuses[name] = status
	})
}

// +die
type _ corev1.ModifyVolumeStatus

// +die
// +die:field:name=ObjectMeta,package=_/meta/v1,die=ObjectMetaDie
// +die:field:name=Spec,die=PersistentVolumeClaimSpecDie
type _ corev1.PersistentVolumeClaimTemplate
