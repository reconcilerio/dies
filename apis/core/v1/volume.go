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
type _ = corev1.Volume

type volumeDieExtension interface {
	HostPathDie(fn func(d HostPathVolumeSourceDie)) VolumeDie
	EmptyDirDie(fn func(d EmptyDirVolumeSourceDie)) VolumeDie
	GCEPersistentDiskDie(fn func(d GCEPersistentDiskVolumeSourceDie)) VolumeDie
	AWSElasticBlockStoreDie(fn func(d AWSElasticBlockStoreVolumeSourceDie)) VolumeDie
	GitRepoDie(fn func(d GitRepoVolumeSourceDie)) VolumeDie
	SecretDie(fn func(d SecretVolumeSourceDie)) VolumeDie
	NFSDie(fn func(d NFSVolumeSourceDie)) VolumeDie
	ISCSIDie(fn func(d ISCSIVolumeSourceDie)) VolumeDie
	GlusterfsDie(fn func(d GlusterfsVolumeSourceDie)) VolumeDie
	PersistentVolumeClaimDie(fn func(d PersistentVolumeClaimVolumeSourceDie)) VolumeDie
	RBDDie(fn func(d RBDVolumeSourceDie)) VolumeDie
	FlexVolumeDie(fn func(d FlexVolumeSourceDie)) VolumeDie
	CinderDie(fn func(d CinderVolumeSourceDie)) VolumeDie
	CephFSDie(fn func(d CephFSVolumeSourceDie)) VolumeDie
	FlockerDie(fn func(d FlockerVolumeSourceDie)) VolumeDie
	DownwardAPIDie(fn func(d DownwardAPIVolumeSourceDie)) VolumeDie
	FCDie(fn func(d FCVolumeSourceDie)) VolumeDie
	AzureFileDie(fn func(d AzureFileVolumeSourceDie)) VolumeDie
	ConfigMapDie(fn func(d ConfigMapVolumeSourceDie)) VolumeDie
	VsphereVolumeDie(fn func(d VsphereVirtualDiskVolumeSourceDie)) VolumeDie
	QuobyteDie(fn func(d QuobyteVolumeSourceDie)) VolumeDie
	AzureDiskDie(fn func(d AzureDiskVolumeSourceDie)) VolumeDie
	PhotonPersistentDiskDie(fn func(d PhotonPersistentDiskVolumeSourceDie)) VolumeDie
	ProjectedDie(fn func(d ProjectedVolumeSourceDie)) VolumeDie
	PortworxVolumeDie(fn func(d PortworxVolumeSourceDie)) VolumeDie
	ScaleIODie(fn func(d ScaleIOVolumeSourceDie)) VolumeDie
	StorageOSDie(fn func(d StorageOSVolumeSourceDie)) VolumeDie
	CSIDie(fn func(d CSIVolumeSourceDie)) VolumeDie
	EphemeralDie(fn func(d EphemeralVolumeSourceDie)) VolumeDie
}

func (d *volumeDie) HostPathDie(fn func(d HostPathVolumeSourceDie)) VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := HostPathVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.HostPath)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			HostPath: d.DieReleasePtr(),
		}
	})
}

func (d *volumeDie) EmptyDirDie(fn func(d EmptyDirVolumeSourceDie)) VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := EmptyDirVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.EmptyDir)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			EmptyDir: d.DieReleasePtr(),
		}
	})
}

func (d *volumeDie) GCEPersistentDiskDie(fn func(d GCEPersistentDiskVolumeSourceDie)) VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := GCEPersistentDiskVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.GCEPersistentDisk)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			GCEPersistentDisk: d.DieReleasePtr(),
		}
	})
}

func (d *volumeDie) AWSElasticBlockStoreDie(fn func(d AWSElasticBlockStoreVolumeSourceDie)) VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := AWSElasticBlockStoreVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.AWSElasticBlockStore)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			AWSElasticBlockStore: d.DieReleasePtr(),
		}
	})
}

func (d *volumeDie) GitRepoDie(fn func(d GitRepoVolumeSourceDie)) VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := GitRepoVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.GitRepo)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			GitRepo: d.DieReleasePtr(),
		}
	})
}

func (d *volumeDie) SecretDie(fn func(d SecretVolumeSourceDie)) VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := SecretVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.Secret)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			Secret: d.DieReleasePtr(),
		}
	})
}

func (d *volumeDie) NFSDie(fn func(d NFSVolumeSourceDie)) VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := NFSVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.NFS)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			NFS: d.DieReleasePtr(),
		}
	})
}

func (d *volumeDie) ISCSIDie(fn func(d ISCSIVolumeSourceDie)) VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := ISCSIVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.ISCSI)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			ISCSI: d.DieReleasePtr(),
		}
	})
}

func (d *volumeDie) GlusterfsDie(fn func(d GlusterfsVolumeSourceDie)) VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := GlusterfsVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.Glusterfs)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			Glusterfs: d.DieReleasePtr(),
		}
	})
}

func (d *volumeDie) PersistentVolumeClaimDie(fn func(d PersistentVolumeClaimVolumeSourceDie)) VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := PersistentVolumeClaimVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.PersistentVolumeClaim)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			PersistentVolumeClaim: d.DieReleasePtr(),
		}
	})
}

func (d *volumeDie) RBDDie(fn func(d RBDVolumeSourceDie)) VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := RBDVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.RBD)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			RBD: d.DieReleasePtr(),
		}
	})
}

func (d *volumeDie) FlexVolumeDie(fn func(d FlexVolumeSourceDie)) VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := FlexVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.FlexVolume)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			FlexVolume: d.DieReleasePtr(),
		}
	})
}

func (d *volumeDie) CinderDie(fn func(d CinderVolumeSourceDie)) VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := CinderVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.Cinder)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			Cinder: d.DieReleasePtr(),
		}
	})
}

func (d *volumeDie) CephFSDie(fn func(d CephFSVolumeSourceDie)) VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := CephFSVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.CephFS)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			CephFS: d.DieReleasePtr(),
		}
	})
}

func (d *volumeDie) FlockerDie(fn func(d FlockerVolumeSourceDie)) VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := FlockerVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.Flocker)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			Flocker: d.DieReleasePtr(),
		}
	})
}

func (d *volumeDie) DownwardAPIDie(fn func(d DownwardAPIVolumeSourceDie)) VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := DownwardAPIVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.DownwardAPI)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			DownwardAPI: d.DieReleasePtr(),
		}
	})
}

func (d *volumeDie) FCDie(fn func(d FCVolumeSourceDie)) VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := FCVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.FC)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			FC: d.DieReleasePtr(),
		}
	})
}

func (d *volumeDie) AzureFileDie(fn func(d AzureFileVolumeSourceDie)) VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := AzureFileVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.AzureFile)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			AzureFile: d.DieReleasePtr(),
		}
	})
}

func (d *volumeDie) ConfigMapDie(fn func(d ConfigMapVolumeSourceDie)) VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := ConfigMapVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.ConfigMap)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			ConfigMap: d.DieReleasePtr(),
		}
	})
}

func (d *volumeDie) VsphereVolumeDie(fn func(d VsphereVirtualDiskVolumeSourceDie)) VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := VsphereVirtualDiskVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.VsphereVolume)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			VsphereVolume: d.DieReleasePtr(),
		}
	})
}

func (d *volumeDie) QuobyteDie(fn func(d QuobyteVolumeSourceDie)) VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := QuobyteVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.Quobyte)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			Quobyte: d.DieReleasePtr(),
		}
	})
}

func (d *volumeDie) AzureDiskDie(fn func(d AzureDiskVolumeSourceDie)) VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := AzureDiskVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.AzureDisk)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			AzureDisk: d.DieReleasePtr(),
		}
	})
}

func (d *volumeDie) PhotonPersistentDiskDie(fn func(d PhotonPersistentDiskVolumeSourceDie)) VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := PhotonPersistentDiskVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.PhotonPersistentDisk)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			PhotonPersistentDisk: d.DieReleasePtr(),
		}
	})
}

func (d *volumeDie) ProjectedDie(fn func(d ProjectedVolumeSourceDie)) VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := ProjectedVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.Projected)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			Projected: d.DieReleasePtr(),
		}
	})
}

func (d *volumeDie) PortworxVolumeDie(fn func(d PortworxVolumeSourceDie)) VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := PortworxVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.PortworxVolume)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			PortworxVolume: d.DieReleasePtr(),
		}
	})
}

func (d *volumeDie) ScaleIODie(fn func(d ScaleIOVolumeSourceDie)) VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := ScaleIOVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.ScaleIO)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			ScaleIO: d.DieReleasePtr(),
		}
	})
}

func (d *volumeDie) StorageOSDie(fn func(d StorageOSVolumeSourceDie)) VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := StorageOSVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.StorageOS)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			StorageOS: d.DieReleasePtr(),
		}
	})
}

func (d *volumeDie) CSIDie(fn func(d CSIVolumeSourceDie)) VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := CSIVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.CSI)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			CSI: d.DieReleasePtr(),
		}
	})
}

func (d *volumeDie) EphemeralDie(fn func(d EphemeralVolumeSourceDie)) VolumeDie {
	return d.DieStamp(func(r *corev1.Volume) {
		d := EphemeralVolumeSourceBlank.DieImmutable(false).DieFeedPtr(r.Ephemeral)
		fn(d)
		r.VolumeSource = corev1.VolumeSource{
			Ephemeral: d.DieReleasePtr(),
		}
	})
}

// +die
type _ = corev1.HostPathVolumeSource

// +die
type _ = corev1.EmptyDirVolumeSource

// +die
type _ = corev1.GCEPersistentDiskVolumeSource

// +die
type _ = corev1.AWSElasticBlockStoreVolumeSource

// +die
type _ = corev1.GitRepoVolumeSource

// +die
type _ = corev1.SecretVolumeSource

type secretVolumeSourceDieExtension interface {
	ItemDie(key string, fn func(d KeyToPathDie)) SecretVolumeSourceDie
}

func (d *secretVolumeSourceDie) ItemDie(key string, fn func(d KeyToPathDie)) SecretVolumeSourceDie {
	return d.DieStamp(func(r *corev1.SecretVolumeSource) {
		for i := range r.Items {
			if key == r.Items[i].Key {
				d := KeyToPathBlank.DieImmutable(false).DieFeed(r.Items[i])
				fn(d)
				r.Items[i] = d.DieRelease()
				return
			}
		}

		d := KeyToPathBlank.DieImmutable(false).DieFeed(corev1.KeyToPath{Key: key})
		fn(d)
		r.Items = append(r.Items, d.DieRelease())
	})
}

// +die
type _ = corev1.NFSVolumeSource

// +die
type _ = corev1.ISCSIVolumeSource

type iscsiVolumeSourceDieExtension interface {
	SecretRefDie(fn func(d LocalObjectReferenceDie)) ISCSIVolumeSourceDie
}

func (d *iscsiVolumeSourceDie) SecretRefDie(fn func(d LocalObjectReferenceDie)) ISCSIVolumeSourceDie {
	return d.DieStamp(func(r *corev1.ISCSIVolumeSource) {
		d := LocalObjectReferenceBlank.DieImmutable(false).DieFeedPtr(r.SecretRef)
		fn(d)
		r.SecretRef = d.DieReleasePtr()
	})
}

// +die
type _ = corev1.GlusterfsVolumeSource

// +die
type _ = corev1.PersistentVolumeClaimVolumeSource

// +die
type _ = corev1.RBDVolumeSource

type rbdVolumeSourceDieExtension interface {
	SecretRefDie(fn func(d LocalObjectReferenceDie)) RBDVolumeSourceDie
}

func (d *rbdVolumeSourceDie) SecretRefDie(fn func(d LocalObjectReferenceDie)) RBDVolumeSourceDie {
	return d.DieStamp(func(r *corev1.RBDVolumeSource) {
		d := LocalObjectReferenceBlank.DieImmutable(false).DieFeedPtr(r.SecretRef)
		fn(d)
		r.SecretRef = d.DieReleasePtr()
	})
}

// +die
type _ = corev1.FlexVolumeSource

type flexVolumeSourceDieExtension interface {
	SecretRefDie(fn func(d LocalObjectReferenceDie)) FlexVolumeSourceDie
}

func (d *flexVolumeSourceDie) SecretRefDie(fn func(d LocalObjectReferenceDie)) FlexVolumeSourceDie {
	return d.DieStamp(func(r *corev1.FlexVolumeSource) {
		d := LocalObjectReferenceBlank.DieImmutable(false).DieFeedPtr(r.SecretRef)
		fn(d)
		r.SecretRef = d.DieReleasePtr()
	})
}

// +die
type _ = corev1.CinderVolumeSource

type cinderVolumeSourceDieExtension interface {
	SecretRefDie(fn func(d LocalObjectReferenceDie)) CinderVolumeSourceDie
}

func (d *cinderVolumeSourceDie) SecretRefDie(fn func(d LocalObjectReferenceDie)) CinderVolumeSourceDie {
	return d.DieStamp(func(r *corev1.CinderVolumeSource) {
		d := LocalObjectReferenceBlank.DieImmutable(false).DieFeedPtr(r.SecretRef)
		fn(d)
		r.SecretRef = d.DieReleasePtr()
	})
}

// +die
type _ = corev1.CephFSVolumeSource

type cephFSVolumeSourceDieExtension interface {
	SecretRefDie(fn func(d LocalObjectReferenceDie)) CephFSVolumeSourceDie
}

func (d *cephFSVolumeSourceDie) SecretRefDie(fn func(d LocalObjectReferenceDie)) CephFSVolumeSourceDie {
	return d.DieStamp(func(r *corev1.CephFSVolumeSource) {
		d := LocalObjectReferenceBlank.DieImmutable(false).DieFeedPtr(r.SecretRef)
		fn(d)
		r.SecretRef = d.DieReleasePtr()
	})
}

// +die
type _ = corev1.FlockerVolumeSource

// +die
type _ = corev1.DownwardAPIVolumeSource

type downwardAPIVolumeSourceDieExtension interface {
	ItemDie(path string, fn func(d DownwardAPIVolumeFileDie)) DownwardAPIVolumeSourceDie
}

func (d *downwardAPIVolumeSourceDie) ItemDie(path string, fn func(d DownwardAPIVolumeFileDie)) DownwardAPIVolumeSourceDie {
	return d.DieStamp(func(r *corev1.DownwardAPIVolumeSource) {
		for i := range r.Items {
			if path == r.Items[i].Path {
				d := DownwardAPIVolumeFileBlank.DieImmutable(false).DieFeed(r.Items[i])
				fn(d)
				r.Items[i] = d.DieRelease()
				return
			}
		}

		d := DownwardAPIVolumeFileBlank.DieImmutable(false).DieFeed(corev1.DownwardAPIVolumeFile{Path: path})
		fn(d)
		r.Items = append(r.Items, d.DieRelease())
	})
}

// +die
type _ = corev1.DownwardAPIVolumeFile

type downwardAPIVolumeFileDieExtension interface {
	FieldRefDie(fn func(d ObjectFieldSelectorDie)) DownwardAPIVolumeFileDie
	ResourceFieldRefDie(fn func(d ResourceFieldSelectorDie)) DownwardAPIVolumeFileDie
}

func (d *downwardAPIVolumeFileDie) FieldRefDie(fn func(d ObjectFieldSelectorDie)) DownwardAPIVolumeFileDie {
	return d.DieStamp(func(r *corev1.DownwardAPIVolumeFile) {
		d := ObjectFieldSelectorBlank.DieImmutable(false).DieFeedPtr(r.FieldRef)
		fn(d)
		r.FieldRef = d.DieReleasePtr()
	})
}

func (d *downwardAPIVolumeFileDie) ResourceFieldRefDie(fn func(d ResourceFieldSelectorDie)) DownwardAPIVolumeFileDie {
	return d.DieStamp(func(r *corev1.DownwardAPIVolumeFile) {
		d := ResourceFieldSelectorBlank.DieImmutable(false).DieFeedPtr(r.ResourceFieldRef)
		fn(d)
		r.ResourceFieldRef = d.DieReleasePtr()
	})
}

// +die
type _ = corev1.FCVolumeSource

// +die
type _ = corev1.AzureFileVolumeSource

// +die
type _ = corev1.ConfigMapVolumeSource

type configMapVolumeSourceDieExtension interface {
	Name(v string) ConfigMapVolumeSourceDie
	ItemDie(key string, fn func(d KeyToPathDie)) ConfigMapVolumeSourceDie
}

func (d *configMapVolumeSourceDie) Name(v string) ConfigMapVolumeSourceDie {
	return d.DieStamp(func(r *corev1.ConfigMapVolumeSource) {
		r.Name = v
	})
}

func (d *configMapVolumeSourceDie) ItemDie(key string, fn func(d KeyToPathDie)) ConfigMapVolumeSourceDie {
	return d.DieStamp(func(r *corev1.ConfigMapVolumeSource) {
		for i := range r.Items {
			if key == r.Items[i].Key {
				d := KeyToPathBlank.DieImmutable(false).DieFeed(r.Items[i])
				fn(d)
				r.Items[i] = d.DieRelease()
				return
			}
		}

		d := KeyToPathBlank.DieImmutable(false).DieFeed(corev1.KeyToPath{Key: key})
		fn(d)
		r.Items = append(r.Items, d.DieRelease())
	})
}

// +die
type _ = corev1.VsphereVirtualDiskVolumeSource

// +die
type _ = corev1.QuobyteVolumeSource

// +die
type _ = corev1.AzureDiskVolumeSource

// +die
type _ = corev1.PhotonPersistentDiskVolumeSource

// +die
type _ = corev1.ProjectedVolumeSource

type projectedVolumeSourceDieExtension interface {
	SourcesDie(sources ...VolumeProjectionDie) ProjectedVolumeSourceDie
}

func (d *projectedVolumeSourceDie) SourcesDie(sources ...VolumeProjectionDie) ProjectedVolumeSourceDie {
	return d.DieStamp(func(r *corev1.ProjectedVolumeSource) {
		r.Sources = make([]corev1.VolumeProjection, len(sources))
		for i := range sources {
			r.Sources[i] = sources[i].DieRelease()
		}
	})
}

// +die
type _ = corev1.VolumeProjection

type volumeProjectionDieExtension interface {
	SecretDie(fn func(d SecretProjectionDie)) VolumeProjectionDie
	DownwardAPIDie(fn func(d DownwardAPIProjectionDie)) VolumeProjectionDie
	ConfigMapDie(fn func(d ConfigMapProjectionDie)) VolumeProjectionDie
	ServiceAccountTokenDie(fn func(d ServiceAccountTokenProjectionDie)) VolumeProjectionDie
}

func (d *volumeProjectionDie) SecretDie(fn func(d SecretProjectionDie)) VolumeProjectionDie {
	return d.DieStamp(func(r *corev1.VolumeProjection) {
		d := SecretProjectionBlank.DieImmutable(false).DieFeedPtr(r.Secret)
		fn(d)
		r.Secret = d.DieReleasePtr()
	})
}

func (d *volumeProjectionDie) DownwardAPIDie(fn func(d DownwardAPIProjectionDie)) VolumeProjectionDie {
	return d.DieStamp(func(r *corev1.VolumeProjection) {
		d := DownwardAPIProjectionBlank.DieImmutable(false).DieFeedPtr(r.DownwardAPI)
		fn(d)
		r.DownwardAPI = d.DieReleasePtr()
	})
}

func (d *volumeProjectionDie) ConfigMapDie(fn func(d ConfigMapProjectionDie)) VolumeProjectionDie {
	return d.DieStamp(func(r *corev1.VolumeProjection) {
		d := ConfigMapProjectionBlank.DieImmutable(false).DieFeedPtr(r.ConfigMap)
		fn(d)
		r.ConfigMap = d.DieReleasePtr()
	})
}

func (d *volumeProjectionDie) ServiceAccountTokenDie(fn func(d ServiceAccountTokenProjectionDie)) VolumeProjectionDie {
	return d.DieStamp(func(r *corev1.VolumeProjection) {
		d := ServiceAccountTokenProjectionBlank.DieImmutable(false).DieFeedPtr(r.ServiceAccountToken)
		fn(d)
		r.ServiceAccountToken = d.DieReleasePtr()
	})
}

// +die
type _ = corev1.SecretProjection

type secretProjectionDieExtension interface {
	Name(v string) SecretProjectionDie
	ItemDie(key string, fn func(d KeyToPathDie)) SecretProjectionDie
}

func (d *secretProjectionDie) Name(v string) SecretProjectionDie {
	return d.DieStamp(func(r *corev1.SecretProjection) {
		r.Name = v
	})
}

func (d *secretProjectionDie) ItemDie(key string, fn func(d KeyToPathDie)) SecretProjectionDie {
	return d.DieStamp(func(r *corev1.SecretProjection) {
		for i := range r.Items {
			if key == r.Items[i].Key {
				d := KeyToPathBlank.DieImmutable(false).DieFeed(r.Items[i])
				fn(d)
				r.Items[i] = d.DieRelease()
				return
			}
		}

		d := KeyToPathBlank.DieImmutable(false).DieFeed(corev1.KeyToPath{Key: key})
		fn(d)
		r.Items = append(r.Items, d.DieRelease())
	})
}

// +die
type _ = corev1.DownwardAPIProjection

type downwardAPIProjectionDieExtension interface {
	ItemDie(path string, fn func(d DownwardAPIVolumeFileDie)) DownwardAPIProjectionDie
}

func (d *downwardAPIProjectionDie) ItemDie(path string, fn func(d DownwardAPIVolumeFileDie)) DownwardAPIProjectionDie {
	return d.DieStamp(func(r *corev1.DownwardAPIProjection) {
		for i := range r.Items {
			if path == r.Items[i].Path {
				d := DownwardAPIVolumeFileBlank.DieImmutable(false).DieFeed(r.Items[i])
				fn(d)
				r.Items[i] = d.DieRelease()
				return
			}
		}

		d := DownwardAPIVolumeFileBlank.DieImmutable(false).DieFeed(corev1.DownwardAPIVolumeFile{Path: path})
		fn(d)
		r.Items = append(r.Items, d.DieRelease())
	})
}

// +die
type _ = corev1.ConfigMapProjection

type configmapProjectionDieExtension interface {
	Name(v string) ConfigMapProjectionDie
	ItemDie(key string, fn func(d KeyToPathDie)) ConfigMapProjectionDie
}

func (d *configMapProjectionDie) Name(v string) ConfigMapProjectionDie {
	return d.DieStamp(func(r *corev1.ConfigMapProjection) {
		r.Name = v
	})
}

func (d *configMapProjectionDie) ItemDie(key string, fn func(d KeyToPathDie)) ConfigMapProjectionDie {
	return d.DieStamp(func(r *corev1.ConfigMapProjection) {
		for i := range r.Items {
			if key == r.Items[i].Key {
				d := KeyToPathBlank.DieImmutable(false).DieFeed(r.Items[i])
				fn(d)
				r.Items[i] = d.DieRelease()
				return
			}
		}

		d := KeyToPathBlank.DieImmutable(false).DieFeed(corev1.KeyToPath{Key: key})
		fn(d)
		r.Items = append(r.Items, d.DieRelease())
	})
}

// +die
type _ = corev1.ServiceAccountTokenProjection

// +die
type _ = corev1.PortworxVolumeSource

// +die
type _ = corev1.ScaleIOVolumeSource

type scaleIOVolumeSourceDieExtension interface {
	SecretRefDie(fn func(d LocalObjectReferenceDie)) ScaleIOVolumeSourceDie
}

func (d *scaleIOVolumeSourceDie) SecretRefDie(fn func(d LocalObjectReferenceDie)) ScaleIOVolumeSourceDie {
	return d.DieStamp(func(r *corev1.ScaleIOVolumeSource) {
		d := LocalObjectReferenceBlank.DieImmutable(false).DieFeedPtr(r.SecretRef)
		fn(d)
		r.SecretRef = d.DieReleasePtr()
	})
}

// +die
type _ = corev1.StorageOSVolumeSource

type storageOSVolumeSourceDieExtension interface {
	SecretRefDie(fn func(d LocalObjectReferenceDie)) StorageOSVolumeSourceDie
}

func (d *storageOSVolumeSourceDie) SecretRefDie(fn func(d LocalObjectReferenceDie)) StorageOSVolumeSourceDie {
	return d.DieStamp(func(r *corev1.StorageOSVolumeSource) {
		d := LocalObjectReferenceBlank.DieImmutable(false).DieFeedPtr(r.SecretRef)
		fn(d)
		r.SecretRef = d.DieReleasePtr()
	})
}

// +die
type _ = corev1.CSIVolumeSource

type csiVolumeSourceDieExtension interface {
	VolumeAttribute(key, value string) CSIVolumeSourceDie
	NodePublishSecretRefDie(fn func(d LocalObjectReferenceDie)) CSIVolumeSourceDie
}

func (d *csiVolumeSourceDie) VolumeAttribute(key, value string) CSIVolumeSourceDie {
	return d.DieStamp(func(r *corev1.CSIVolumeSource) {
		r.VolumeAttributes[key] = value
	})
}

func (d *csiVolumeSourceDie) NodePublishSecretRefDie(fn func(d LocalObjectReferenceDie)) CSIVolumeSourceDie {
	return d.DieStamp(func(r *corev1.CSIVolumeSource) {
		d := LocalObjectReferenceBlank.DieImmutable(false).DieFeedPtr(r.NodePublishSecretRef)
		fn(d)
		r.NodePublishSecretRef = d.DieReleasePtr()
	})
}

// +die
type _ = corev1.EphemeralVolumeSource

type ephemeralVolumeSourceDieExtension interface {
	VolumeClaimTemplateDie(fn func(d PersistentVolumeClaimTemplateDie)) EphemeralVolumeSourceDie
}

func (d *ephemeralVolumeSourceDie) VolumeClaimTemplateDie(fn func(d PersistentVolumeClaimTemplateDie)) EphemeralVolumeSourceDie {
	return d.DieStamp(func(r *corev1.EphemeralVolumeSource) {
		d := PersistentVolumeClaimTemplateBlank.DieImmutable(false).DieFeedPtr(r.VolumeClaimTemplate)
		fn(d)
		r.VolumeClaimTemplate = d.DieReleasePtr()
	})
}

// +die
type _ = corev1.KeyToPath
