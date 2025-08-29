/*
Copyright 2025 the original author or authors.

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
	resourcev1 "k8s.io/api/resource/v1"
)

// +die:object=true,apiVersion=resource.k8s.io/v1,kind=DeviceClass
type _ = resourcev1.DeviceClass

// +die
// +die:field:name=Selectors,die=DeviceSelectorDie,listType=atomic
// +die:field:name=Config,die=DeviceClassConfigurationDie,listType=atomic
type _ = resourcev1.DeviceClassSpec

// +die
// +die:field:name=CEL,die=CELDeviceSelectorDie,pointer=true
type _ = resourcev1.DeviceSelector

// +die
type _ = resourcev1.CELDeviceSelector

// +die
type _ = resourcev1.DeviceClassConfiguration

// OpaqueDie mutates Opaque as a die.
//
// Opaque provides driver-specific configuration parameters.
func (d *DeviceClassConfigurationDie) OpaqueDie(fn func(d *OpaqueDeviceConfigurationDie)) *DeviceClassConfigurationDie {
	return d.DieStamp(func(r *resourcev1.DeviceClassConfiguration) {
		d := OpaqueDeviceConfigurationBlank.DieImmutable(false).DieFeedPtr(r.Opaque)
		fn(d)
		r.Opaque = d.DieReleasePtr()
	})
}

// Opaque provides driver-specific configuration parameters.
func (d *DeviceClassConfigurationDie) Opaque(v *resourcev1.OpaqueDeviceConfiguration) *DeviceClassConfigurationDie {
	return d.DieStamp(func(r *resourcev1.DeviceClassConfiguration) {
		r.Opaque = v
	})
}

// +die
// +die:field:name=Opaque,die=OpaqueDeviceConfigurationDie,pointer=true
type _ = resourcev1.DeviceConfiguration

// +die
type _ = resourcev1.OpaqueDeviceConfiguration
