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
		d := HostPathVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.HostPath)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			HostPath: d.DieReleasePtr(),
		}
	})
}

func (d *VolumeDie) EmptyDirDie(fn func(d *EmptyDirVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := EmptyDirVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.EmptyDir)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			EmptyDir: d.DieReleasePtr(),
		}
	})
}

func (d *VolumeDie) GCEPersistentDiskDie(fn func(d *GCEPersistentDiskVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := GCEPersistentDiskVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.GCEPersistentDisk)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			GCEPersistentDisk: d.DieReleasePtr(),
		}
	})
}

func (d *VolumeDie) AWSElasticBlockStoreDie(fn func(d *AWSElasticBlockStoreVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := AWSElasticBlockStoreVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.AWSElasticBlockStore)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			AWSElasticBlockStore: d.DieReleasePtr(),
		}
	})
}

func (d *VolumeDie) GitRepoDie(fn func(d *GitRepoVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := GitRepoVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.GitRepo)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			GitRepo: d.DieReleasePtr(),
		}
	})
}

func (d *VolumeDie) SecretDie(fn func(d *SecretVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := SecretVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.Secret)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			Secret: d.DieReleasePtr(),
		}
	})
}

func (d *VolumeDie) NFSDie(fn func(d *NFSVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := NFSVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.NFS)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			NFS: d.DieReleasePtr(),
		}
	})
}

func (d *VolumeDie) ISCSIDie(fn func(d *ISCSIVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := ISCSIVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.ISCSI)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			ISCSI: d.DieReleasePtr(),
		}
	})
}

func (d *VolumeDie) GlusterfsDie(fn func(d *GlusterfsVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := GlusterfsVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.Glusterfs)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			Glusterfs: d.DieReleasePtr(),
		}
	})
}

func (d *VolumeDie) PersistentVolumeClaimDie(fn func(d *PersistentVolumeClaimVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := PersistentVolumeClaimVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.PersistentVolumeClaim)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			PersistentVolumeClaim: d.DieReleasePtr(),
		}
	})
}

func (d *VolumeDie) RBDDie(fn func(d *RBDVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := RBDVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.RBD)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			RBD: d.DieReleasePtr(),
		}
	})
}

func (d *VolumeDie) FlexVolumeDie(fn func(d *FlexVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := FlexVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.FlexVolume)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			FlexVolume: d.DieReleasePtr(),
		}
	})
}

func (d *VolumeDie) CinderDie(fn func(d *CinderVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := CinderVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.Cinder)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			Cinder: d.DieReleasePtr(),
		}
	})
}

func (d *VolumeDie) CephFSDie(fn func(d *CephFSVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := CephFSVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.CephFS)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			CephFS: d.DieReleasePtr(),
		}
	})
}

func (d *VolumeDie) FlockerDie(fn func(d *FlockerVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := FlockerVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.Flocker)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			Flocker: d.DieReleasePtr(),
		}
	})
}

func (d *VolumeDie) DownwardAPIDie(fn func(d *DownwardAPIVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := DownwardAPIVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.DownwardAPI)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			DownwardAPI: d.DieReleasePtr(),
		}
	})
}

func (d *VolumeDie) FCDie(fn func(d *FCVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := FCVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.FC)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			FC: d.DieReleasePtr(),
		}
	})
}

func (d *VolumeDie) AzureFileDie(fn func(d *AzureFileVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := AzureFileVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.AzureFile)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			AzureFile: d.DieReleasePtr(),
		}
	})
}

func (d *VolumeDie) ConfigMapDie(fn func(d *ConfigMapVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := ConfigMapVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.ConfigMap)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			ConfigMap: d.DieReleasePtr(),
		}
	})
}

func (d *VolumeDie) VsphereVolumeDie(fn func(d *VsphereVirtualDiskVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := VsphereVirtualDiskVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.VsphereVolume)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			VsphereVolume: d.DieReleasePtr(),
		}
	})
}

func (d *VolumeDie) QuobyteDie(fn func(d *QuobyteVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := QuobyteVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.Quobyte)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			Quobyte: d.DieReleasePtr(),
		}
	})
}

func (d *VolumeDie) AzureDiskDie(fn func(d *AzureDiskVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := AzureDiskVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.AzureDisk)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			AzureDisk: d.DieReleasePtr(),
		}
	})
}

func (d *VolumeDie) PhotonPersistentDiskDie(fn func(d *PhotonPersistentDiskVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := PhotonPersistentDiskVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.PhotonPersistentDisk)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			PhotonPersistentDisk: d.DieReleasePtr(),
		}
	})
}

func (d *VolumeDie) ProjectedDie(fn func(d *ProjectedVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := ProjectedVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.Projected)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			Projected: d.DieReleasePtr(),
		}
	})
}

func (d *VolumeDie) PortworxVolumeDie(fn func(d *PortworxVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := PortworxVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.PortworxVolume)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			PortworxVolume: d.DieReleasePtr(),
		}
	})
}

func (d *VolumeDie) ScaleIODie(fn func(d *ScaleIOVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := ScaleIOVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.ScaleIO)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			ScaleIO: d.DieReleasePtr(),
		}
	})
}

func (d *VolumeDie) StorageOSDie(fn func(d *StorageOSVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := StorageOSVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.StorageOS)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			StorageOS: d.DieReleasePtr(),
		}
	})
}

func (d *VolumeDie) CSIDie(fn func(d *CSIVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := CSIVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.CSI)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			CSI: d.DieReleasePtr(),
		}
	})
}

func (d *VolumeDie) EphemeralDie(fn func(d *EphemeralVolumeSourceDie)) *VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := EphemeralVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.Ephemeral)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			Ephemeral: d.DieReleasePtr(),
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
