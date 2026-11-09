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

// Package sprocket registers the dies of every version of the example.dev
// sprockets, and converts between them: the conversion is here, rather than
// in a version's package, as it needs every version, and a version's package
// needn't load the others.
package sprocket

import (
	"reconciler.io/dies/v2/die/schemes"
	sprocketv1 "reconciler.io/dies/v2/examples/sprocket/v1"
	sprocketv2 "reconciler.io/dies/v2/examples/sprocket/v2"
)

// SchemeBuilder registers the dies of every version, in order of version
// priority.
var SchemeBuilder = schemes.NewBuilder(
	sprocketv2.SchemeBuilder.Add,
	sprocketv1.SchemeBuilder.Add,
)

// AddToScheme registers the dies of every version with a scheme.
var AddToScheme = SchemeBuilder.AddToScheme
