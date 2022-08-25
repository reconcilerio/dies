//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2021-2022 the original author or authors.

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

// Code generated by diegen. DO NOT EDIT.

package v1beta1

import (
	"dies.dev/apis/meta/v1"
	json "encoding/json"
	fmtx "fmt"
	storagev1beta1 "k8s.io/api/storage/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
)

var CSIStorageCapacityBlank = (&CSIStorageCapacityDie{}).DieFeed(storagev1beta1.CSIStorageCapacity{})

type CSIStorageCapacityDie struct {
	v1.FrozenObjectMeta
	mutable bool
	r       storagev1beta1.CSIStorageCapacity
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *CSIStorageCapacityDie) DieImmutable(immutable bool) *CSIStorageCapacityDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *CSIStorageCapacityDie) DieFeed(r storagev1beta1.CSIStorageCapacity) *CSIStorageCapacityDie {
	if d.mutable {
		d.FrozenObjectMeta = v1.FreezeObjectMeta(r.ObjectMeta)
		d.r = r
		return d
	}
	return &CSIStorageCapacityDie{
		FrozenObjectMeta: v1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *CSIStorageCapacityDie) DieFeedPtr(r *storagev1beta1.CSIStorageCapacity) *CSIStorageCapacityDie {
	if r == nil {
		r = &storagev1beta1.CSIStorageCapacity{}
	}
	return d.DieFeed(*r)
}

// DieFeedRawExtension returns the resource managed by the die as an raw extension.
func (d *CSIStorageCapacityDie) DieFeedRawExtension(raw runtime.RawExtension) *CSIStorageCapacityDie {
	b, _ := json.Marshal(raw)
	r := storagev1beta1.CSIStorageCapacity{}
	_ = json.Unmarshal(b, &r)
	return d.DieFeed(r)
}

// DieRelease returns the resource managed by the die.
func (d *CSIStorageCapacityDie) DieRelease() storagev1beta1.CSIStorageCapacity {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *CSIStorageCapacityDie) DieReleasePtr() *storagev1beta1.CSIStorageCapacity {
	r := d.DieRelease()
	return &r
}

// DieReleaseUnstructured returns the resource managed by the die as an unstructured object.
func (d *CSIStorageCapacityDie) DieReleaseUnstructured() *unstructured.Unstructured {
	r := d.DieReleasePtr()
	u, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(r)
	return &unstructured.Unstructured{
		Object: u,
	}
}

// DieReleaseRawExtension returns the resource managed by the die as an raw extension.
func (d *CSIStorageCapacityDie) DieReleaseRawExtension() runtime.RawExtension {
	r := d.DieReleasePtr()
	b, _ := json.Marshal(r)
	raw := runtime.RawExtension{}
	_ = json.Unmarshal(b, &raw)
	return raw
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *CSIStorageCapacityDie) DieStamp(fn func(r *storagev1beta1.CSIStorageCapacity)) *CSIStorageCapacityDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *CSIStorageCapacityDie) DeepCopy() *CSIStorageCapacityDie {
	r := *d.r.DeepCopy()
	return &CSIStorageCapacityDie{
		FrozenObjectMeta: v1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

var _ runtime.Object = (*CSIStorageCapacityDie)(nil)

func (d *CSIStorageCapacityDie) DeepCopyObject() runtime.Object {
	return d.r.DeepCopy()
}

func (d *CSIStorageCapacityDie) GetObjectKind() schema.ObjectKind {
	r := d.DieRelease()
	return r.GetObjectKind()
}

func (d *CSIStorageCapacityDie) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.r)
}

func (d *CSIStorageCapacityDie) UnmarshalJSON(b []byte) error {
	if d == CSIStorageCapacityBlank {
		return fmtx.Errorf("cannot unmarshal into the blank die, create a copy first")
	}
	if !d.mutable {
		return fmtx.Errorf("cannot unmarshal into immutable dies, create a mutable version first")
	}
	r := &storagev1beta1.CSIStorageCapacity{}
	err := json.Unmarshal(b, r)
	*d = *d.DieFeed(*r)
	return err
}

// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
func (d *CSIStorageCapacityDie) APIVersion(v string) *CSIStorageCapacityDie {
	return d.DieStamp(func(r *storagev1beta1.CSIStorageCapacity) {
		r.APIVersion = v
	})
}

// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
func (d *CSIStorageCapacityDie) Kind(v string) *CSIStorageCapacityDie {
	return d.DieStamp(func(r *storagev1beta1.CSIStorageCapacity) {
		r.Kind = v
	})
}

// MetadataDie stamps the resource's ObjectMeta field with a mutable die.
func (d *CSIStorageCapacityDie) MetadataDie(fn func(d *v1.ObjectMetaDie)) *CSIStorageCapacityDie {
	return d.DieStamp(func(r *storagev1beta1.CSIStorageCapacity) {
		d := v1.ObjectMetaBlank.DieImmutable(false).DieFeed(r.ObjectMeta)
		fn(d)
		r.ObjectMeta = d.DieRelease()
	})
}

// NodeTopology defines which nodes have access to the storage for which capacity was reported. If not set, the storage is not accessible from any node in the cluster. If empty, the storage is accessible from all nodes. This field is immutable.
func (d *CSIStorageCapacityDie) NodeTopology(v *metav1.LabelSelector) *CSIStorageCapacityDie {
	return d.DieStamp(func(r *storagev1beta1.CSIStorageCapacity) {
		r.NodeTopology = v
	})
}

// The name of the StorageClass that the reported capacity applies to. It must meet the same requirements as the name of a StorageClass object (non-empty, DNS subdomain). If that object no longer exists, the CSIStorageCapacity object is obsolete and should be removed by its creator. This field is immutable.
func (d *CSIStorageCapacityDie) StorageClassName(v string) *CSIStorageCapacityDie {
	return d.DieStamp(func(r *storagev1beta1.CSIStorageCapacity) {
		r.StorageClassName = v
	})
}

// Capacity is the value reported by the CSI driver in its GetCapacityResponse for a GetCapacityRequest with topology and parameters that match the previous fields.
//
// The semantic is currently (CSI spec 1.2) defined as: The available capacity, in bytes, of the storage that can be used to provision volumes. If not set, that information is currently unavailable.
func (d *CSIStorageCapacityDie) Capacity(v *resource.Quantity) *CSIStorageCapacityDie {
	return d.DieStamp(func(r *storagev1beta1.CSIStorageCapacity) {
		r.Capacity = v
	})
}

// MaximumVolumeSize is the value reported by the CSI driver in its GetCapacityResponse for a GetCapacityRequest with topology and parameters that match the previous fields.
//
// This is defined since CSI spec 1.4.0 as the largest size that may be used in a CreateVolumeRequest.capacity_range.required_bytes field to create a volume with the same parameters as those in GetCapacityRequest. The corresponding value in the Kubernetes API is ResourceRequirements.Requests in a volume claim.
func (d *CSIStorageCapacityDie) MaximumVolumeSize(v *resource.Quantity) *CSIStorageCapacityDie {
	return d.DieStamp(func(r *storagev1beta1.CSIStorageCapacity) {
		r.MaximumVolumeSize = v
	})
}
