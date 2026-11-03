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
	"errors"
	"fmt"
)

// Error is a die's error. Writes return one, for content a die can't hold, as
// do functions that encode, or take an expression or content from elsewhere;
// their Must forms panic with it. A mistake in the code, such as a write the
// types don't allow, panics with one.
type Error struct {
	err error
}

func (e *Error) Error() string { return e.err.Error() }

func (e *Error) Unwrap() error { return e.err }

// Errorf returns an [Error], formatted as [fmt.Errorf] formats, so %w wraps.
func Errorf(format string, args ...any) *Error {
	return &Error{err: fmt.Errorf(format, args...)}
}

// Wrap returns err as an [Error], unless it is one.
func Wrap(err error) *Error {
	if e, ok := errors.AsType[*Error](err); ok {
		return e
	}
	return &Error{err: err}
}

// Must returns v, or panics with err: the Must form of a function that
// returns an error, or a read of a die, which can't fail but for a mistake in
// the code.
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// Check panics with err, if it's set: the Must form of a function that
// returns only an error.
func Check(err error) {
	if err != nil {
		panic(err)
	}
}
