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
	"strings"

	"k8s.io/apimachinery/pkg/util/validation/field"
)

// FieldPath converts a path to apimachinery's field path, as validation
// reports locations. A list-map element is its index in the list, as the
// API server reports it, when h is given and has the element; otherwise it is
// its keys, such as [name=app]. A map entry is its key.
func FieldPath(h *Holder, p Path) *field.Path {
	var fp *field.Path
	for i, s := range p {
		switch s.kind {
		case stepField:
			if fp == nil {
				fp = field.NewPath(s.name)
			} else {
				fp = fp.Child(s.name)
			}
		case stepIndex:
			fp = fp.Index(s.index)
		case stepMapKey:
			fp = fp.Key(s.name)
		case stepKey:
			fp = elemPath(fp, h, p[:i], s)
		}
	}
	return fp
}

func elemPath(fp *field.Path, h *Holder, list Path, s Step) *field.Path {
	if h != nil {
		if l, err := h.ValueAt(list); err == nil {
			if i, err := FindElem(l, s); err == nil && i >= 0 {
				return fp.Index(i)
			}
		}
	}
	return fp.Key(strings.TrimSuffix(strings.TrimPrefix(s.String(), "["), "]"))
}
