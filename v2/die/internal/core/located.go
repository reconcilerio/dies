/*
Copyright 2026 the original author or authors.

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

package core

import (
	"reflect"
)

// Located is anything that identifies a location within a root die: a Path,
// a child die or a handle.
//
// Navigating a blank builds a typed location without a die at hand:
// DeploymentBlank.Spec().Template() is the *PodTemplateSpecDie at
// .spec.template. Reads have no side effects, so this is free.
type Located interface {
	DiePath() Path
}

// rooted locations know the type of their root, so a location taken from a
// different type is rejected rather than silently matching nothing. Only dies
// and handles implement it; a bare Path is unchecked.
type Rooted interface {
	DieRootType() reflect.Type
}
