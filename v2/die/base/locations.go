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

package base

import (
	"reflect"

	"reconciler.io/dies/v2/die/internal/core"
)

func pathsOf(root reflect.Type, at []Located) ([]Path, error) {
	paths := make([]Path, len(at))
	for i, l := range at {
		if r, ok := l.(core.Rooted); ok && r.DieRootType() != root {
			return nil, core.Errorf("die: location %s is within %s, not %s", l.DiePath(), r.DieRootType(), root)
		}
		paths[i] = l.DiePath()
	}
	return paths, nil
}

// Path is a location in an object's wire schema: a sequence of steps, such
// as .spec.template.spec.containers[name=app].
type Path = core.Path

// Located is anything that identifies a location within a root die: a [Path],
// a die or a handle.
type Located = core.Located

// LocatedAs is a location of a value of type X: a die for X, or a handle to
// a leaf of type X.
type LocatedAs[X any] interface {
	Located
	DieValue(*X)
}
