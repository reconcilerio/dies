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

// Package lists has list dies: plain lists of dies, as Kubernetes lists are
// lists of objects.
package lists

import (
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"reconciler.io/dies/v2/die/base"
)

// List is the base of a list die. A list is read, not patched, so unlike
// other dies a list die is a plain list type, as Kubernetes lists are: type
// and list metadata, and items, which are dies. Libraries that fill and read
// lists by reflecting on their Items, such as apimachinery's meta.SetList and
// meta.ExtractList, which controller-runtime's cache and fake client use,
// take it as they take any list.
//
// List dies embed it, and declare DeepCopyObject. List itself implements
// [metav1.ListInterface], by its list metadata.
type List[D ListItem[D]] struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []D `json:"items"`
}

// ListItem is a die that can be an item of a list: an object die that copies
// as itself.
type ListItem[D any] interface {
	base.Die
	runtime.Object
	DeepCopy() D
}

// DieList is a list die of any type, as the Kubernetes list interfaces see a
// list: meta.ExtractList and meta.EachListItem read its items. [DieListOf] has
// them typed, for generic code.
type DieList interface {
	runtime.Object
	metav1.ListInterface
	metav1.ListMetaAccessor
}

// DieListOf is a list die of D items, for generic code.
type DieListOf[D ListItem[D]] interface {
	DieList
	Metadata() *metav1.ListMeta
	DieItems() []D
	DieCopy() List[D]
}

// Metadata returns the list's metadata, its own: a list isn't a die, so its
// metadata is changed in place, as its items are.
func (l *List[D]) Metadata() *metav1.ListMeta { return &l.ListMeta }

// DieItems returns the list's items, the dies themselves, not copies.
func (l *List[D]) DieItems() []D { return l.Items }

// NewList returns a list of a kind, with items.
func NewList[D ListItem[D]](apiVersion, kind string, items ...D) List[D] {
	return List[D]{TypeMeta: metav1.TypeMeta{APIVersion: apiVersion, Kind: kind}, Items: items}
}

// DieCopy returns a copy of the list, with its items copied.
func (l *List[D]) DieCopy() List[D] {
	out := List[D]{TypeMeta: l.TypeMeta}
	l.ListMeta.DeepCopyInto(&out.ListMeta)
	if l.Items != nil {
		out.Items = make([]D, len(l.Items))
		for i, item := range l.Items {
			if !reflect.ValueOf(item).IsNil() {
				out.Items[i] = item.DeepCopy()
			}
		}
	}
	return out
}
