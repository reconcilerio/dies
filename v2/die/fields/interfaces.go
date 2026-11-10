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

package fields

import "reconciler.io/dies/v2/die/base"

// DieTypedOf carries an immutable typed die's types for generic code, as
// [base.DieOf] does a die's, with its type metadata's accessors.
type DieTypedOf[T any, D base.DieOf[T, D, M], M base.DieMutOf[T, D, M]] interface {
	base.DieOf[T, D, M]
	base.DieTyped
	APIVersion() Value[string]
	Kind() Value[string]
}

// DieTypedMutOf carries a mutable typed die's types for generic code, as
// [base.DieMutOf] does a die's, with its type metadata's accessors.
type DieTypedMutOf[T any, D base.DieOf[T, D, M], M base.DieMutOf[T, D, M]] interface {
	base.DieMutOf[T, D, M]
	base.DieTyped
	APIVersion() ValueMut[string]
	Kind() ValueMut[string]
}

// DieObjectOf carries an immutable object die's types for generic code, as
// [DieTypedOf] does a typed die's, as a [base.DieObject]. It doesn't name
// the die of the object's metadata, which is generated for a version of
// apimachinery's meta/v1 and so isn't the die API's to choose: a package of
// dies for meta/v1, such as reconciler.io/dies/v2/apis/meta/v1, has an
// interface that adds its Metadata(). An unstructured die is a
// [base.DieObject], but not a DieObjectOf: its fields are reached through its
// nodes, not accessors.
type DieObjectOf[T any, D base.DieOf[T, D, M], M base.DieMutOf[T, D, M]] interface {
	DieTypedOf[T, D, M]
	base.DieObject
}

// DieObjectMutOf carries a mutable object die's types for generic code, as
// [DieObjectOf] does an immutable one's.
type DieObjectMutOf[T any, D base.DieOf[T, D, M], M base.DieMutOf[T, D, M]] interface {
	DieTypedMutOf[T, D, M]
	base.DieObject
}
