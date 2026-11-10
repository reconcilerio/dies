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

package sprocket

import (
	"strings"

	"reconciler.io/dies/v2/die/conversion"
	"reconciler.io/dies/v2/die/jigs"
	sprocketv1 "reconciler.io/dies/v2/examples/sprocket/v1"
	sprocketv2 "reconciler.io/dies/v2/examples/sprocket/v2"
)

// ImageReference presents an image as one reference, repository:tag, as v1
// has it. A reference with nothing after its colon, "app:", is a repository
// of that name, so that it round trips.
var ImageReference = jigs.Transform[sprocketv2.SprocketImage, string]{
	To: func(i sprocketv2.SprocketImage) string {
		if i.Tag == "" {
			return i.Repository
		}
		return i.Repository + ":" + i.Tag
	},
	From: func(ref string, _ sprocketv2.SprocketImage) sprocketv2.SprocketImage {
		// a colon before the last slash is a registry's port
		if i := strings.LastIndex(ref, ":"); i > strings.LastIndex(ref, "/") && i < len(ref)-1 {
			return sprocketv2.SprocketImage{Repository: ref[:i], Tag: ref[i+1:]}
		}
		return sprocketv2.SprocketImage{Repository: ref}
	},
}

// Conversion converts a Sprocket between v1 and v2, the hub. v1's lubricant
// has no place in v2, nor v2's paused in v1: each is stashed in the other.
var Conversion = conversion.New(sprocketv2.SprocketBlank).
	AddSpoke(sprocketv1.SprocketBlank,
		jigs.Bind(sprocketv1.SprocketBlank.Spec().Size(), sprocketv2.SprocketBlank.Spec().Teeth()),
		jigs.BindVia(sprocketv1.SprocketBlank.Spec().Image(), sprocketv2.SprocketBlank.Spec().Image(), ImageReference),
		jigs.Bind(sprocketv1.SprocketBlank.Status().ObservedGeneration(), sprocketv2.SprocketBlank.Status().ObservedGeneration()),
	)
