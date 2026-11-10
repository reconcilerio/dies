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

// Package v1 is a second example die package, whose types the widget
// example embeds. Its metadata must be read from here, not from markers
// alone, for the widget's dies to agree with the runtime.
package v1

// GadgetSpec's parts are a list map by their upstream-style markers, which
// the +die marker overrides to atomic.
//
// +die
// +die:field:name=parts,listType=atomic
type GadgetSpec struct {
	// +listType=map
	// +listMapKey=name
	Parts []GadgetPart `json:"parts,omitempty"`
}

// +die
type GadgetPart struct {
	Name string `json:"name"`
}
