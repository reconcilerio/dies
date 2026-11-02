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

// ContentOf returns the content of a die, or nil for a nil die.
func ContentOf(d Die) *Content { return core.Must(ContentOfE(d)) }

// ContentOfE is [ContentOf], returning an error rather than panicking.
func ContentOfE(d Die) (*Content, error) {
	if d == nil {
		return nil, nil
	}
	if v := reflect.ValueOf(d); v.Kind() == reflect.Pointer && v.IsNil() {
		return nil, nil
	}
	c, ok := d.(interface{ content() (*Content, error) })
	if !ok {
		return nil, core.Errorf("die: %T is not a die", d)
	}
	return c.content()
}

func (b *state[T]) content() (*Content, error) {
	data, err := b.MarshalJSON()
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return nil, core.Errorf("die: decoding the value of %s: %w", reflect.TypeFor[T](), err)
	}
	var intent []string
	describeIntent(core.NodeAt(core.IntentOf(b.hold()), b.at), nil, &intent)
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

func equal(a interface{ content() (*Content, error) }, b Die) (bool, error) {
	mine, err := a.content()
	if err != nil {
		return false, err
	}
	theirs, err := ContentOfE(b)
	if err != nil {
		return false, err
	}
	return reflect.DeepEqual(mine, theirs), nil
}

// ReleaseOf returns the current value of a die as its type has it, as
// [DieBase.DieRelease] does, boxed; or nil for a nil die. Unlike [Content],
// it holds neither intent nor fields the type doesn't have.
func ReleaseOf(d Die) any { return core.Must(ReleaseOfE(d)) }

// ReleaseOfE is [ReleaseOf], returning an error rather than panicking.
func ReleaseOfE(d Die) (any, error) {
	if d == nil {
		return nil, nil
	}
	if v := reflect.ValueOf(d); v.Kind() == reflect.Pointer && v.IsNil() {
		return nil, nil
	}
	r, ok := d.(interface{ release() (any, error) })
	if !ok {
		return nil, core.Errorf("die: %T is not a die", d)
	}
	return r.release()
}

func (b *state[T]) release() (any, error) { return b.DieReleaseE() }
