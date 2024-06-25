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
	nodev1 "k8s.io/api/node/v1"
	diecorev1 "reconciler.io/dies/apis/core/v1"
)

// +die:object=true,apiVersion=node.k8s.io/v1,kind=RuntimeClass
type _ = nodev1.RuntimeClass

// +die
type _ = nodev1.Overhead

// +die
// +die:field:name=Tolerations,package=_/core/v1,die=TolerationDie,listMapKey=Key
type _ = nodev1.Scheduling

// deprecated?: use TolerationDie
func (d *SchedulingDie) TolerationsDie(key string, fn func(d *diecorev1.TolerationDie)) *SchedulingDie {
	return d.TolerationDie(key, fn)
}
