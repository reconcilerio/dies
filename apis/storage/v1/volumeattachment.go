/*
Copyright 2022 the original author or authors.

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
	storagev1 "k8s.io/api/storage/v1"
	diecorev1 "reconciler.io/dies/apis/core/v1"
)

// +die:object=true,apiVersion=storage.k8s.io/v1,kind=VolumeAttachment
type _ = storagev1.VolumeAttachment

// +die
type _ = storagev1.VolumeAttachmentSpec

func (d *VolumeAttachmentSpecDie) SourceDie(fn func(d *VolumeAttachmentSourceDie)) *VolumeAttachmentSpecDie {
	return d.DieStamp(func(r *storagev1.VolumeAttachmentSpec) {
		d := VolumeAttachmentSourceBlank.DieImmutable(false).DieFeed(r.Source)
		fn(d)
		r.Source = d.DieRelease()
	})
}

// +die
type _ = storagev1.VolumeAttachmentSource

func (d *VolumeAttachmentSourceDie) SourceDie(fn func(d *diecorev1.PersistentVolumeSpecDie)) *VolumeAttachmentSourceDie {
	return d.DieStamp(func(r *storagev1.VolumeAttachmentSource) {
		d := diecorev1.PersistentVolumeSpecBlank.DieImmutable(false).DieFeedPtr(r.InlineVolumeSpec)
		fn(d)
		r.InlineVolumeSpec = d.DieReleasePtr()
	})
}

// +die
type _ = storagev1.VolumeAttachmentStatus

func (d *VolumeAttachmentStatusDie) AttachErrorDie(fn func(d *VolumeErrorDie)) *VolumeAttachmentStatusDie {
	return d.DieStamp(func(r *storagev1.VolumeAttachmentStatus) {
		d := VolumeErrorBlank.DieImmutable(false).DieFeedPtr(r.AttachError)
		fn(d)
		r.AttachError = d.DieReleasePtr()
	})
}

func (d *VolumeAttachmentStatusDie) DetachErrorDie(fn func(d *VolumeErrorDie)) *VolumeAttachmentStatusDie {
	return d.DieStamp(func(r *storagev1.VolumeAttachmentStatus) {
		d := VolumeErrorBlank.DieImmutable(false).DieFeedPtr(r.DetachError)
		fn(d)
		r.DetachError = d.DieReleasePtr()
	})
}

// +die
type _ = storagev1.VolumeError
