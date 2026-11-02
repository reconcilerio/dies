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

// Package objects has the interface of object dies for generic code,
// DieObjectOf, apart from the base of dies, as it names the die of an
// object's metadata, from apis/meta/v1.
package objects

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	diemetav1 "reconciler.io/dies/v2/apis/meta/v1"
	"reconciler.io/dies/v2/die/base"
	"reconciler.io/dies/v2/die/fields"
)

// DieObjectOf carries an object die's types for generic code, as
// fields.DieTypedOf does a typed die's, with its metadata's accessor.
// Generated object dies check they satisfy it; die.DieObjectOf is the same
// interface. An unstructured
// die is a base.DieObject, but not a DieObjectOf: its fields are reached
// through its nodes, not accessors.
type DieObjectOf[T any, D base.DieOf[T, D]] interface {
	fields.DieTypedOf[T, D]
	base.DieObject
	Metadata() *diemetav1.ObjectMetaDie
	DieFeedUnstructured(u *unstructured.Unstructured) D
	DieFeedUnstructuredE(u *unstructured.Unstructured) (D, error)
}
