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

package v1beta1

import (
	diecorev1 "dies.dev/apis/core/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
)

// +die:object=true
type _ = policyv1beta1.PodSecurityPolicy

// +die
type _ = policyv1beta1.PodSecurityPolicySpec

func (d *PodSecurityPolicySpecDie) SELinuxDie(fn func(d *SELinuxStrategyOptionsDie)) *PodSecurityPolicySpecDie {
	return d.DieStamp(func(r *policyv1beta1.PodSecurityPolicySpec) {
		d := SELinuxStrategyOptionsBlank.DieImmutable(false).DieFeed(r.SELinux)
		fn(d)
		r.SELinux = d.DieRelease()
	})
}

func (d *PodSecurityPolicySpecDie) RunAsUserDie(fn func(d *RunAsUserStrategyOptionsDie)) *PodSecurityPolicySpecDie {
	return d.DieStamp(func(r *policyv1beta1.PodSecurityPolicySpec) {
		d := RunAsUserStrategyOptionsBlank.DieImmutable(false).DieFeed(r.RunAsUser)
		fn(d)
		r.RunAsUser = d.DieRelease()
	})
}

func (d *PodSecurityPolicySpecDie) RunAsGroupDie(fn func(d *RunAsGroupStrategyOptionsDie)) *PodSecurityPolicySpecDie {
	return d.DieStamp(func(r *policyv1beta1.PodSecurityPolicySpec) {
		d := RunAsGroupStrategyOptionsBlank.DieImmutable(false).DieFeedPtr(r.RunAsGroup)
		fn(d)
		r.RunAsGroup = d.DieReleasePtr()
	})
}

func (d *PodSecurityPolicySpecDie) SupplementalGroupsDie(fn func(d *SupplementalGroupsStrategyOptionsDie)) *PodSecurityPolicySpecDie {
	return d.DieStamp(func(r *policyv1beta1.PodSecurityPolicySpec) {
		d := SupplementalGroupsStrategyOptionsBlank.DieImmutable(false).DieFeed(r.SupplementalGroups)
		fn(d)
		r.SupplementalGroups = d.DieRelease()
	})
}

func (d *PodSecurityPolicySpecDie) FSGroupDie(fn func(d *FSGroupStrategyOptionsDie)) *PodSecurityPolicySpecDie {
	return d.DieStamp(func(r *policyv1beta1.PodSecurityPolicySpec) {
		d := FSGroupStrategyOptionsBlank.DieImmutable(false).DieFeed(r.FSGroup)
		fn(d)
		r.FSGroup = d.DieRelease()
	})
}

func (d *PodSecurityPolicySpecDie) AllowedHostPathsDie(paths ...*AllowedHostPathDie) *PodSecurityPolicySpecDie {
	return d.DieStamp(func(r *policyv1beta1.PodSecurityPolicySpec) {
		r.AllowedHostPaths = make([]policyv1beta1.AllowedHostPath, len(paths))
		for i := range paths {
			r.AllowedHostPaths[i] = paths[i].DieRelease()
		}
	})
}

func (d *PodSecurityPolicySpecDie) AllowedFlexVolumesDie(volumes ...*AllowedFlexVolumeDie) *PodSecurityPolicySpecDie {
	return d.DieStamp(func(r *policyv1beta1.PodSecurityPolicySpec) {
		r.AllowedFlexVolumes = make([]policyv1beta1.AllowedFlexVolume, len(volumes))
		for i := range volumes {
			r.AllowedFlexVolumes[i] = volumes[i].DieRelease()
		}
	})
}

func (d *PodSecurityPolicySpecDie) AllowedCSIDriversDie(drivers ...*AllowedCSIDriverDie) *PodSecurityPolicySpecDie {
	return d.DieStamp(func(r *policyv1beta1.PodSecurityPolicySpec) {
		r.AllowedCSIDrivers = make([]policyv1beta1.AllowedCSIDriver, len(drivers))
		for i := range drivers {
			r.AllowedCSIDrivers[i] = drivers[i].DieRelease()
		}
	})
}

func (d *PodSecurityPolicySpecDie) RuntimeClassDie(fn func(d *RuntimeClassStrategyOptionsDie)) *PodSecurityPolicySpecDie {
	return d.DieStamp(func(r *policyv1beta1.PodSecurityPolicySpec) {
		d := RuntimeClassStrategyOptionsBlank.DieImmutable(false).DieFeedPtr(r.RuntimeClass)
		fn(d)
		r.RuntimeClass = d.DieReleasePtr()
	})
}

// +die
type _ = policyv1beta1.HostPortRange

// +die
type _ = policyv1beta1.SELinuxStrategyOptions

func (d *SELinuxStrategyOptionsDie) RuntimeClassDie(fn func(d *diecorev1.SELinuxOptionsDie)) *SELinuxStrategyOptionsDie {
	return d.DieStamp(func(r *policyv1beta1.SELinuxStrategyOptions) {
		d := diecorev1.SELinuxOptionsBlank.DieImmutable(false).DieFeedPtr(r.SELinuxOptions)
		fn(d)
		r.SELinuxOptions = d.DieReleasePtr()
	})
}

// +die
type _ = policyv1beta1.RunAsUserStrategyOptions

func (d *RunAsUserStrategyOptionsDie) RangesDie(ranges ...*IDRangeDie) *RunAsUserStrategyOptionsDie {
	return d.DieStamp(func(r *policyv1beta1.RunAsUserStrategyOptions) {
		r.Ranges = make([]policyv1beta1.IDRange, len(ranges))
		for i := range ranges {
			r.Ranges[i] = ranges[i].DieRelease()
		}
	})
}

// +die
type _ = policyv1beta1.RunAsGroupStrategyOptions

func (d *RunAsGroupStrategyOptionsDie) RangesDie(ranges ...*IDRangeDie) *RunAsGroupStrategyOptionsDie {
	return d.DieStamp(func(r *policyv1beta1.RunAsGroupStrategyOptions) {
		r.Ranges = make([]policyv1beta1.IDRange, len(ranges))
		for i := range ranges {
			r.Ranges[i] = ranges[i].DieRelease()
		}
	})
}

// +die
type _ = policyv1beta1.SupplementalGroupsStrategyOptions

func (d *SupplementalGroupsStrategyOptionsDie) RangesDie(ranges ...*IDRangeDie) *SupplementalGroupsStrategyOptionsDie {
	return d.DieStamp(func(r *policyv1beta1.SupplementalGroupsStrategyOptions) {
		r.Ranges = make([]policyv1beta1.IDRange, len(ranges))
		for i := range ranges {
			r.Ranges[i] = ranges[i].DieRelease()
		}
	})
}

// +die
type _ = policyv1beta1.FSGroupStrategyOptions

func (d *FSGroupStrategyOptionsDie) RangesDie(ranges ...*IDRangeDie) *FSGroupStrategyOptionsDie {
	return d.DieStamp(func(r *policyv1beta1.FSGroupStrategyOptions) {
		r.Ranges = make([]policyv1beta1.IDRange, len(ranges))
		for i := range ranges {
			r.Ranges[i] = ranges[i].DieRelease()
		}
	})
}

// +die
type _ = policyv1beta1.AllowedHostPath

// +die
type _ = policyv1beta1.AllowedFlexVolume

// +die
type _ = policyv1beta1.AllowedCSIDriver

// +die
type _ = policyv1beta1.RuntimeClassStrategyOptions

// +die
type _ = policyv1beta1.IDRange
