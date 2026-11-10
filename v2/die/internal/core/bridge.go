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

import "reflect"

// The state of a die is unexported, so that its methods are only the die's
// own. The die package installs these at init, for the packages of handles,
// views and unstructured dies to reach that state.
var (
	// StateOf returns a die's holder and its location within it, or false
	// for a value that isn't a die.
	StateOf func(d any) (*Holder, Path, bool)
	// SetState sets the state of a die that has none yet.
	SetState func(d any, h *Holder, at Path)
	// ValueTypeOf returns the type of a die's value.
	ValueTypeOf func(d any) reflect.Type
	// IsMutable reports whether a die is a mutable die.
	IsMutable func(d any) bool
	// HandleStateOf returns a handle's holder and its location within it,
	// or false for a value that isn't a handle. The package of handles
	// installs it.
	HandleStateOf func(l any) (*Holder, Path, bool)
)
