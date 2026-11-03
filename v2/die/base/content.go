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
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"

	"reconciler.io/dies/v2/die/internal/core"
	"reconciler.io/dies/v2/die/internal/refs"
)

// Content is what a die holds, as plain values: its current value and its
// intent. Two dies are [DieBase.Equal] when their contents are deeply equal,
// so a diff of their contents, such as dietest.Transform makes with go-cmp,
// is empty exactly when they are equal.
type Content struct {
	// Value is the current value as its JSON decodes: objects as
	// map[string]any, lists as []any, and numbers as int64 where they are
	// integers, float64 where that keeps them exactly, and json.Number
	// otherwise. Fields the die's type doesn't have are included.
	Value any
	// Intent is the recorded intent, one operation per line, at paths
	// within the die, as [DieBase.DieIntent] describes it: for a jig, its
	// sources' intent through its bindings, so that jigs backed by different
	// types compare by what they bind.
	Intent []string
}

// ExtractContentOf returns the content of a die, or nil for a nil die.
func ExtractContentOf(d Die) (*Content, error) {
	if d == nil {
		return nil, nil
	}
	if v := reflect.ValueOf(d); v.Kind() == reflect.Pointer && v.IsNil() {
		return nil, nil
	}
	r, ok := refs.Of(d)
	if !ok {
		return nil, core.Errorf("die: %T is not a die", d)
	}
	return content(r)
}

// MustExtractContentOf is [ExtractContentOf], panicking rather than returning
// an error.
func MustExtractContentOf(d Die) *Content { return core.Must(ExtractContentOf(d)) }

func content(r refs.Ref) (*Content, error) {
	data, err := r.MarshalJSON()
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return nil, core.Errorf("die: decoding the value of %s: %w", r.T, err)
	}
	var intent []string
	describeIntent(core.NodeAt(core.IntentOf(r.H), r.At), nil, &intent)
	return &Content{Value: plainNumbers(v), Intent: intent}, nil
}

// plainNumbers replaces the json.Numbers in a decoded value with int64 or
// float64 where those hold them exactly, so that they print as numbers.
func plainNumbers(v any) any {
	switch v := v.(type) {
	case map[string]any:
		for k, e := range v {
			v[k] = plainNumbers(e)
		}
	case []any:
		for i, e := range v {
			v[i] = plainNumbers(e)
		}
	case json.Number:
		if i, err := v.Int64(); err == nil {
			return i
		}
		if f, err := v.Float64(); err == nil && strconv.FormatFloat(f, 'g', -1, 64) == v.String() {
			return f
		}
	}
	return v
}

func equal(a refs.Ref, b Die) (bool, error) {
	mine, err := content(a)
	if err != nil {
		return false, err
	}
	theirs, err := ExtractContentOf(b)
	if err != nil {
		return false, err
	}
	return reflect.DeepEqual(mine, theirs), nil
}

// ReleaseOf returns the current value of a die as its type has it, as
// [DieBase.DieRelease] does, but boxed as a value; or nil for a nil die.
// Unlike [Content], it holds neither intent nor fields the type doesn't have.
func ReleaseOf(d Die) any {
	if d == nil {
		return nil
	}
	if v := reflect.ValueOf(d); v.Kind() == reflect.Pointer && v.IsNil() {
		return nil
	}
	r, ok := refs.Of(d)
	if !ok {
		panic(core.Errorf("die: %T is not a die", d))
	}
	return core.Must(r.Value()).Interface()
}
