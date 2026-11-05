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

// Package validation makes apimachinery's field errors at locations within
// dies, for the DieValidate methods that [reconciler.io/dies/v2/die.ValidateCreate]
// calls. A location is a die or a handle, so an error names the field it is
// about, and its path from the root is the location's:
//
//	if d.Replicas().Get() < 0 {
//		errs = append(errs, validation.Invalid(d.Replicas(), "must be non-negative"))
//	}
//
// Errors that take a value read it at the location. For the errors of
// apimachinery's that have no helper here, [Path] converts a location to a
// field path.
package validation

import (
	"reflect"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"reconciler.io/dies/v2/die/base"
	"reconciler.io/dies/v2/die/internal/core"
)

type (
	ErrorList = field.ErrorList
)

// Path returns the field path of a location: a die, a handle or a [base.Path].
// A list-map element is its index in the list, as the API server reports it,
// when the location is a die or handle whose list has the element; otherwise
// it is its keys, such as containers[name=app]. A map entry is its key.
func Path(l base.Located) *field.Path {
	return core.FieldPath(holderOf(l), l.DiePath())
}

// holderOf returns the state a location is within, or nil for a bare path.
func holderOf(l base.Located) *core.Holder {
	if h, _, ok := core.StateOf(l); ok {
		return h
	}
	if r, ok := l.(interface{ DieRoot() base.Die }); ok {
		if root := r.DieRoot(); root != nil {
			if h, _, ok := core.StateOf(root); ok {
				return h
			}
		}
	}
	return nil
}

// value returns the value at a location, or nil for a bare path.
func value(l base.Located) any {
	h := holderOf(l)
	if h == nil {
		return nil
	}
	v, err := h.ValueAt(l.DiePath())
	if err != nil || !v.IsValid() {
		return nil
	}
	return v.Interface()
}

// length returns the length of the list or map at a location.
func length(l base.Located) int {
	v := reflect.ValueOf(value(l))
	switch v.Kind() {
	case reflect.Slice, reflect.Map, reflect.String:
		return v.Len()
	}
	return 0
}

// Required is the error of a required value that is missing.
func Required(l base.Located, detail string) *field.Error {
	return field.Required(Path(l), detail)
}

// Forbidden is the error of a value that may not be set.
func Forbidden(l base.Located, detail string) *field.Error {
	return field.Forbidden(Path(l), detail)
}

// Invalid is the error of a value, the location's, that is malformed or out
// of range.
func Invalid(l base.Located, detail string) *field.Error {
	return field.Invalid(Path(l), value(l), detail)
}

// NotSupported is the error of a value, the location's, that isn't one of
// the supported values.
func NotSupported[T ~string](l base.Located, supported []T) *field.Error {
	return field.NotSupported(Path(l), value(l), supported)
}

// Duplicate is the error of a value, the location's, that must be unique.
func Duplicate(l base.Located) *field.Error {
	return field.Duplicate(Path(l), value(l))
}

// NotFound is the error of a value, the location's, that refers to
// something that doesn't exist.
func NotFound(l base.Located) *field.Error {
	return field.NotFound(Path(l), value(l))
}

// TooLong is the error of a string or list, the location's, longer than max.
func TooLong(l base.Located, max int) *field.Error {
	return field.TooLong(Path(l), value(l), max)
}

// TooMany is the error of a list or map, the location's, with more than max
// items.
func TooMany(l base.Located, max int) *field.Error {
	return field.TooMany(Path(l), length(l), max)
}

// TooFew is the error of a list or map, the location's, with fewer than min
// items.
func TooFew(l base.Located, min int) *field.Error {
	return field.TooFew(Path(l), length(l), min)
}

// InternalError is the error of a failure to validate the location, rather
// than of its value.
func InternalError(l base.Located, err error) *field.Error {
	return field.InternalError(Path(l), err)
}
