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
)

// +die:object=true,apiVersion=storage.k8s.io/v1,kind=VolumeAttachment
type _ = storagev1.VolumeAttachment

// +die
// +die:field:name=Source,die=VolumeAttachmentSourceDie
type _ = storagev1.VolumeAttachmentSpec

// +die
// +die:field:name=InlineVolumeSpec,package=_/core/v1,die=PersistentVolumeSpecDie,pointer=true
type _ = storagev1.VolumeAttachmentSource

// +die
// +die:field:name=AttachError,die=VolumeErrorDie,pointer=true
// +die:field:name=DetachError,die=VolumeErrorDie,pointer=true
type _ = storagev1.VolumeAttachmentStatus

// +die
type _ = storagev1.VolumeError
