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

package v1

import (
	"strings"

	diecorev1 "reconciler.io/dies/v2/apis/core/v1"
	"reconciler.io/dies/v2/die/jigs"
)

// ImageReferenceTransform presents an image string in parts. It is lossless:
// a string that doesn't parse into parts that format back to it exactly, such
// as "repo:" with an empty tag, is presented whole as the repository.
var ImageReferenceTransform = jigs.Transform[string, ImageReference]{
	To: func(image string) ImageReference {
		ref := parseImage(image)
		if ref.String() != image {
			return ImageReference{Repository: image}
		}
		return ref
	},
	From: func(ref ImageReference, _ string) string {
		// the parts are the whole image, so nothing else needs keeping
		return ref.String()
	},
}

func parseImage(image string) ImageReference {
	var ref ImageReference
	ref.Repository, ref.Digest, _ = strings.Cut(image, "@")
	// a tag follows the last colon, unless that colon is a registry's port
	if i := strings.LastIndex(ref.Repository, ":"); i > strings.LastIndex(ref.Repository, "/") {
		ref.Repository, ref.Tag = ref.Repository[:i], ref.Repository[i+1:]
	}
	return ref
}

func (r ImageReference) String() string {
	s := r.Repository
	if r.Tag != "" {
		s += ":" + r.Tag
	}
	if r.Digest != "" {
		s += "@" + r.Digest
	}
	return s
}

// ContainerImageJig presents a container with its image in parts.
var ContainerImageJig = jigs.NewJig[ContainerImage, *ContainerImageDie]().
	For(diecorev1.ContainerBlank,
		jigs.Bind(ContainerImageBlank.Name(), diecorev1.ContainerBlank.Name()),
		jigs.BindVia(ContainerImageBlank.Image(), diecorev1.ContainerBlank.Image(), ImageReferenceTransform),
	)
