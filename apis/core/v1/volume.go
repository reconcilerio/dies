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
type Volume = corev1.Volume

func (d *VolumeDie) HostPathDie(fn func(d *HostPathVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := HostPathVolumeSourceBlank.DieImmutable(false)
		if v := r.HostPath; v != nil {
			d.DieFeed(*v)
		}
		fn(d)
		v := d.DieRelease()
		r.VolumeSource = corev1.VolumeSource{
			HostPath: &v,
		}
	})
}

func (d *VolumeDie) EmptyDirDie(fn func(d *EmptyDirVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := EmptyDirVolumeSourceBlank.DieImmutable(false)
		if v := r.EmptyDir; v != nil {
			d.DieFeed(*v)
		}
		fn(d)
		v := d.DieRelease()
		r.VolumeSource = corev1.VolumeSource{
			EmptyDir: &v,
		}
	})
}

func (d *VolumeDie) GCEPersistentDiskDie(fn func(d *GCEPersistentDiskVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := GCEPersistentDiskVolumeSourceBlank.DieImmutable(false)
		if v := r.GCEPersistentDisk; v != nil {
			d.DieFeed(*v)
		}
		fn(d)
		v := d.DieRelease()
		r.VolumeSource = corev1.VolumeSource{
			GCEPersistentDisk: &v,
		}
	})
}

func (d *VolumeDie) AWSElasticBlockStoreDie(fn func(d *AWSElasticBlockStoreVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := AWSElasticBlockStoreVolumeSourceBlank.DieImmutable(false)
		if v := r.AWSElasticBlockStore; v != nil {
			d.DieFeed(*v)
		}
		fn(d)
		v := d.DieRelease()
		r.VolumeSource = corev1.VolumeSource{
			AWSElasticBlockStore: &v,
		}
	})
}

func (d *VolumeDie) GitRepoDie(fn func(d *GitRepoVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := GitRepoVolumeSourceBlank.DieImmutable(false)
		if v := r.GitRepo; v != nil {
			d.DieFeed(*v)
		}
		fn(d)
		v := d.DieRelease()
		r.VolumeSource = corev1.VolumeSource{
			GitRepo: &v,
		}
	})
}

func (d *VolumeDie) SecretDie(fn func(d *SecretVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := SecretVolumeSourceBlank.DieImmutable(false)
		if v := r.Secret; v != nil {
			d.DieFeed(*v)
		}
		fn(d)
		v := d.DieRelease()
		r.VolumeSource = corev1.VolumeSource{
			Secret: &v,
		}
	})
}

func (d *VolumeDie) NFSDie(fn func(d *NFSVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := NFSVolumeSourceBlank.DieImmutable(false)
		if v := r.NFS; v != nil {
			d.DieFeed(*v)
		}
		fn(d)
		v := d.DieRelease()
		r.VolumeSource = corev1.VolumeSource{
			NFS: &v,
		}
	})
}

func (d *VolumeDie) ISCSIDie(fn func(d *ISCSIVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := ISCSIVolumeSourceBlank.DieImmutable(false)
		if v := r.ISCSI; v != nil {
			d.DieFeed(*v)
		}
		fn(d)
		v := d.DieRelease()
		r.VolumeSource = corev1.VolumeSource{
			ISCSI: &v,
		}
	})
}

func (d *VolumeDie) GlusterfsDie(fn func(d *GlusterfsVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := GlusterfsVolumeSourceBlank.DieImmutable(false)
		if v := r.Glusterfs; v != nil {
			d.DieFeed(*v)
		}
		fn(d)
		v := d.DieRelease()
		r.VolumeSource = corev1.VolumeSource{
			Glusterfs: &v,
		}
	})
}

func (d *VolumeDie) PersistentVolumeClaimDie(fn func(d *PersistentVolumeClaimVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := PersistentVolumeClaimVolumeSourceBlank.DieImmutable(false)
		if v := r.PersistentVolumeClaim; v != nil {
			d.DieFeed(*v)
		}
		fn(d)
		v := d.DieRelease()
		r.VolumeSource = corev1.VolumeSource{
			PersistentVolumeClaim: &v,
		}
	})
}

func (d *VolumeDie) RBDDie(fn func(d *RBDVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := RBDVolumeSourceBlank.DieImmutable(false)
		if v := r.RBD; v != nil {
			d.DieFeed(*v)
		}
		fn(d)
		v := d.DieRelease()
		r.VolumeSource = corev1.VolumeSource{
			RBD: &v,
		}
	})
}

func (d *VolumeDie) FlexVolumeDie(fn func(d *FlexVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := FlexVolumeSourceBlank.DieImmutable(false)
		if v := r.FlexVolume; v != nil {
			d.DieFeed(*v)
		}
		fn(d)
		v := d.DieRelease()
		r.VolumeSource = corev1.VolumeSource{
			FlexVolume: &v,
		}
	})
}

func (d *VolumeDie) CinderDie(fn func(d *CinderVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := CinderVolumeSourceBlank.DieImmutable(false)
		if v := r.Cinder; v != nil {
			d.DieFeed(*v)
		}
		fn(d)
		v := d.DieRelease()
		r.VolumeSource = corev1.VolumeSource{
			Cinder: &v,
		}
	})
}

func (d *VolumeDie) CephFSDie(fn func(d *CephFSVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := CephFSVolumeSourceBlank.DieImmutable(false)
		if v := r.CephFS; v != nil {
			d.DieFeed(*v)
		}
		fn(d)
		v := d.DieRelease()
		r.VolumeSource = corev1.VolumeSource{
			CephFS: &v,
		}
	})
}

func (d *VolumeDie) FlockerDie(fn func(d *FlockerVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := FlockerVolumeSourceBlank.DieImmutable(false)
		if v := r.Flocker; v != nil {
			d.DieFeed(*v)
		}
		fn(d)
		v := d.DieRelease()
		r.VolumeSource = corev1.VolumeSource{
			Flocker: &v,
		}
	})
}

func (d *VolumeDie) DownwardAPIDie(fn func(d *DownwardAPIVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := DownwardAPIVolumeSourceBlank.DieImmutable(false)
		if v := r.DownwardAPI; v != nil {
			d.DieFeed(*v)
		}
		fn(d)
		v := d.DieRelease()
		r.VolumeSource = corev1.VolumeSource{
			DownwardAPI: &v,
		}
	})
}

func (d *VolumeDie) FCDie(fn func(d *FCVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := FCVolumeSourceBlank.DieImmutable(false)
		if v := r.FC; v != nil {
			d.DieFeed(*v)
		}
		fn(d)
		v := d.DieRelease()
		r.VolumeSource = corev1.VolumeSource{
			FC: &v,
		}
	})
}

func (d *VolumeDie) AzureFileDie(fn func(d *AzureFileVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := AzureFileVolumeSourceBlank.DieImmutable(false)
		if v := r.AzureFile; v != nil {
			d.DieFeed(*v)
		}
		fn(d)
		v := d.DieRelease()
		r.VolumeSource = corev1.VolumeSource{
			AzureFile: &v,
		}
	})
}

func (d *VolumeDie) ConfigMapDie(fn func(d *ConfigMapVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := ConfigMapVolumeSourceBlank.DieImmutable(false)
		if v := r.ConfigMap; v != nil {
			d.DieFeed(*v)
		}
		fn(d)
		v := d.DieRelease()
		r.VolumeSource = corev1.VolumeSource{
			ConfigMap: &v,
		}
	})
}

func (d *VolumeDie) VsphereVolumeDie(fn func(d *VsphereVirtualDiskVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := VsphereVirtualDiskVolumeSourceBlank.DieImmutable(false)
		if v := r.VsphereVolume; v != nil {
			d.DieFeed(*v)
		}
		fn(d)
		v := d.DieRelease()
		r.VolumeSource = corev1.VolumeSource{
			VsphereVolume: &v,
		}
	})
}

func (d *VolumeDie) QuobyteDie(fn func(d *QuobyteVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := QuobyteVolumeSourceBlank.DieImmutable(false)
		if v := r.Quobyte; v != nil {
			d.DieFeed(*v)
		}
		fn(d)
		v := d.DieRelease()
		r.VolumeSource = corev1.VolumeSource{
			Quobyte: &v,
		}
	})
}

func (d *VolumeDie) AzureDiskDie(fn func(d *AzureDiskVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := AzureDiskVolumeSourceBlank.DieImmutable(false)
		if v := r.AzureDisk; v != nil {
			d.DieFeed(*v)
		}
		fn(d)
		v := d.DieRelease()
		r.VolumeSource = corev1.VolumeSource{
			AzureDisk: &v,
		}
	})
}

func (d *VolumeDie) PhotonPersistentDiskDie(fn func(d *PhotonPersistentDiskVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := PhotonPersistentDiskVolumeSourceBlank.DieImmutable(false)
		if v := r.PhotonPersistentDisk; v != nil {
			d.DieFeed(*v)
		}
		fn(d)
		v := d.DieRelease()
		r.VolumeSource = corev1.VolumeSource{
			PhotonPersistentDisk: &v,
		}
	})
}

func (d *VolumeDie) ProjectedDie(fn func(d *ProjectedVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := ProjectedVolumeSourceBlank.DieImmutable(false)
		if v := r.Projected; v != nil {
			d.DieFeed(*v)
		}
		fn(d)
		v := d.DieRelease()
		r.VolumeSource = corev1.VolumeSource{
			Projected: &v,
		}
	})
}

func (d *VolumeDie) PortworxVolumeDie(fn func(d *PortworxVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := PortworxVolumeSourceBlank.DieImmutable(false)
		if v := r.PortworxVolume; v != nil {
			d.DieFeed(*v)
		}
		fn(d)
		v := d.DieRelease()
		r.VolumeSource = corev1.VolumeSource{
			PortworxVolume: &v,
		}
	})
}

func (d *VolumeDie) ScaleIODie(fn func(d *ScaleIOVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := ScaleIOVolumeSourceBlank.DieImmutable(false)
		if v := r.ScaleIO; v != nil {
			d.DieFeed(*v)
		}
		fn(d)
		v := d.DieRelease()
		r.VolumeSource = corev1.VolumeSource{
			ScaleIO: &v,
		}
	})
}

func (d *VolumeDie) StorageOSDie(fn func(d *StorageOSVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := StorageOSVolumeSourceBlank.DieImmutable(false)
		if v := r.StorageOS; v != nil {
			d.DieFeed(*v)
		}
		fn(d)
		v := d.DieRelease()
		r.VolumeSource = corev1.VolumeSource{
			StorageOS: &v,
		}
	})
}

func (d *VolumeDie) CSIDie(fn func(d *CSIVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := CSIVolumeSourceBlank.DieImmutable(false)
		if v := r.CSI; v != nil {
			d.DieFeed(*v)
		}
		fn(d)
		v := d.DieRelease()
		r.VolumeSource = corev1.VolumeSource{
			CSI: &v,
		}
	})
}

func (d *VolumeDie) EphemeralDie(fn func(d *EphemeralVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := EphemeralVolumeSourceBlank.DieImmutable(false)
		if v := r.Ephemeral; v != nil {
			d.DieFeed(*v)
		}
		fn(d)
		v := d.DieRelease()
		r.VolumeSource = corev1.VolumeSource{
			Ephemeral: &v,
		}
	})
}

// +die
type HostPathVolumeSource = corev1.HostPathVolumeSource

// +die
type EmptyDirVolumeSource = corev1.EmptyDirVolumeSource

// +die
type GCEPersistentDiskVolumeSource = corev1.GCEPersistentDiskVolumeSource

// +die
type AWSElasticBlockStoreVolumeSource = corev1.AWSElasticBlockStoreVolumeSource

// +die
type GitRepoVolumeSource = corev1.GitRepoVolumeSource

// +die
type SecretVolumeSource = corev1.SecretVolumeSource

// +die
type NFSVolumeSource = corev1.NFSVolumeSource

// +die
type ISCSIVolumeSource = corev1.ISCSIVolumeSource

// +die
type GlusterfsVolumeSource = corev1.GlusterfsVolumeSource

// +die
type PersistentVolumeClaimVolumeSource = corev1.PersistentVolumeClaimVolumeSource

// +die
type RBDVolumeSource = corev1.RBDVolumeSource

// +die
type FlexVolumeSource = corev1.FlexVolumeSource

// +die
type CinderVolumeSource = corev1.CinderVolumeSource

// +die
type CephFSVolumeSource = corev1.CephFSVolumeSource

// +die
type FlockerVolumeSource = corev1.FlockerVolumeSource

// +die
type DownwardAPIVolumeSource = corev1.DownwardAPIVolumeSource

// +die
type FCVolumeSource = corev1.FCVolumeSource

// +die
type AzureFileVolumeSource = corev1.AzureFileVolumeSource

// +die
type ConfigMapVolumeSource = corev1.ConfigMapVolumeSource

// +die
type VsphereVirtualDiskVolumeSource = corev1.VsphereVirtualDiskVolumeSource

// +die
type QuobyteVolumeSource = corev1.QuobyteVolumeSource

// +die
type AzureDiskVolumeSource = corev1.AzureDiskVolumeSource

// +die
type PhotonPersistentDiskVolumeSource = corev1.PhotonPersistentDiskVolumeSource

// +die
type ProjectedVolumeSource = corev1.ProjectedVolumeSource

// +die
type PortworxVolumeSource = corev1.PortworxVolumeSource

// +die
type ScaleIOVolumeSource = corev1.ScaleIOVolumeSource

// +die
type StorageOSVolumeSource = corev1.StorageOSVolumeSource

// +die
type CSIVolumeSource = corev1.CSIVolumeSource

// +die
type EphemeralVolumeSource = corev1.EphemeralVolumeSource
