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

package die

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"slices"
	"strings"

	"k8s.io/apimachinery/pkg/api/operation"
	"k8s.io/apimachinery/pkg/api/validate"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/kube-openapi/pkg/validation/strfmt"

	"reconciler.io/dies/v2/die/internal/core"
)

// The checks of values from their validation markers are made with
// apimachinery's declarative validation, as the API server makes them, so
// their errors are the server's. Validation is of a value as it is, as on a
// create: checks of an update, such as +k8s:immutable, need the old value.

var create = operation.Operation{Type: operation.Create}

// checkValue returns the errors of a present value's checks.
func checkValue(ctx context.Context, fp *field.Path, v reflect.Value, c *core.Checks) field.ErrorList {
	for v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	var errs field.ErrorList
	switch v.Kind() {
	case reflect.String:
		errs = checkString(ctx, fp, v.String(), c)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		errs = checkInt(ctx, fp, v.Int(), c)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v.Uint() <= math.MaxInt64 {
			errs = checkInt(ctx, fp, int64(v.Uint()), c)
		}
	case reflect.Float32, reflect.Float64:
		errs = checkNumber(fp, v.Float(), c)
	case reflect.Slice:
		n := v.Len()
		items := make([]struct{}, n)
		if c.MinItems != nil {
			errs = append(errs, validate.MinItems(ctx, create, fp, items, nil, *c.MinItems)...)
		}
		if c.MaxItems != nil {
			errs = append(errs, validate.MaxItems(ctx, create, fp, items, nil, *c.MaxItems)...)
		}
		if c.Unique {
			for j := range n {
				for i := range j {
					if reflect.DeepEqual(v.Index(i).Interface(), v.Index(j).Interface()) {
						errs = append(errs, field.Duplicate(fp.Index(j), v.Index(j).Interface()))
						break
					}
				}
			}
		}
		if c.Each != nil {
			for i := range n {
				errs = append(errs, checkValue(ctx, fp.Index(i), v.Index(i), c.Each)...)
			}
		}
	case reflect.Map:
		props := make(map[int]struct{}, v.Len())
		for i := range v.Len() {
			props[i] = struct{}{}
		}
		if c.MinProperties != nil {
			errs = append(errs, validate.MinProperties(ctx, create, fp, props, nil, *c.MinProperties)...)
		}
		if c.MaxProperties != nil {
			errs = append(errs, validate.MaxProperties(ctx, create, fp, props, nil, *c.MaxProperties)...)
		}
		if c.Each != nil && v.Type().Key().Kind() == reflect.String {
			keys := v.MapKeys()
			slices.SortFunc(keys, func(a, b reflect.Value) int {
				if a.String() < b.String() {
					return -1
				}
				return 1
			})
			for _, k := range keys {
				errs = append(errs, checkValue(ctx, fp.Key(k.String()), v.MapIndex(k), c.Each)...)
			}
		}
	}
	return errs
}

// checkString returns the errors of a string's checks. An empty string has
// no enum, format or pattern to check: a required field that must not be
// empty says so.
func checkString(ctx context.Context, fp *field.Path, s string, c *core.Checks) field.ErrorList {
	var errs field.ErrorList
	if c.MinLength != nil {
		errs = append(errs, validate.MinLength(ctx, create, fp, &s, nil, *c.MinLength)...)
	}
	if c.MaxLength != nil {
		errs = append(errs, validate.MaxLength(ctx, create, fp, &s, nil, *c.MaxLength)...)
	}
	if c.MaxBytes != nil {
		errs = append(errs, validate.MaxBytes(ctx, create, fp, &s, nil, *c.MaxBytes)...)
	}
	if s == "" {
		return errs
	}
	if len(c.Enum) > 0 {
		errs = append(errs, validate.Enum(ctx, create, fp, &s, nil, sets.New(c.Enum...), nil)...)
	}
	if c.Format != "" {
		errs = append(errs, checkFormat(ctx, fp, s, c.Format)...)
	}
	if c.Pattern != nil && !c.Pattern.MatchString(s) {
		errs = append(errs, field.Invalid(fp, s, fmt.Sprintf("must match the regular expression %q", c.Pattern)).WithOrigin("pattern"))
	}
	return errs
}

// formats are Kubernetes' formats, as +k8s:format names them.
var formats = map[string]func(context.Context, operation.Operation, *field.Path, *string, *string) field.ErrorList{
	"k8s-short-name": validate.ShortName[string],
	"k8s-long-name":  validate.LongName[string],
	// deprecated for new fields, but upstream's fields still have it
	"k8s-long-name-caseless":            validate.LongNameCaseless[string], //nolint:staticcheck // upstream fields are still this format
	"k8s-label-key":                     validate.LabelKey[string],
	"k8s-prefixed-label-key":            validate.PrefixedLabelKey[string],
	"k8s-label-value":                   validate.LabelValue[string],
	"k8s-path-segment-name":             validate.PathSegmentName[string],
	"k8s-uuid":                          validate.UUID[string],
	"k8s-resource-pool-name":            validate.ResourcePoolName[string],
	"k8s-extended-resource-name":        validate.ExtendedResourceName[string],
	"k8s-resource-fully-qualified-name": validate.ResourceFullyQualifiedName[string],
}

// checkFormat returns the errors of a string's format: one of Kubernetes',
// or of OpenAPI's, as a custom resource's schema checks it. A format of
// neither isn't checked.
func checkFormat(ctx context.Context, fp *field.Path, s, format string) field.ErrorList {
	if fn, ok := formats[format]; ok {
		return fn(ctx, create, fp, &s, nil)
	}
	if strfmt.Default.ContainsName(format) && !strfmt.Default.Validates(format, s) {
		return field.ErrorList{field.Invalid(fp, s, "must be of format "+format).WithOrigin("format=" + format)}
	}
	return nil
}

// checkInt returns the errors of an integer's checks, with declarative
// validation's for an inclusive integer bound.
func checkInt(ctx context.Context, fp *field.Path, n int64, c *core.Checks) field.ErrorList {
	var errs field.ErrorList
	rest := *c
	if b := c.Minimum; b != nil && !c.ExclusiveMinimum && *b == math.Trunc(*b) {
		errs = append(errs, validate.Minimum(ctx, create, fp, &n, nil, int64(*b))...)
		rest.Minimum = nil
	}
	if b := c.Maximum; b != nil && !c.ExclusiveMaximum && *b == math.Trunc(*b) {
		errs = append(errs, validate.Maximum(ctx, create, fp, &n, nil, int64(*b))...)
		rest.Maximum = nil
	}
	return append(errs, checkNumber(fp, float64(n), &rest)...)
}

// checkNumber returns the errors of a number's checks, as a custom
// resource's schema makes them.
func checkNumber(fp *field.Path, f float64, c *core.Checks) field.ErrorList {
	var errs field.ErrorList
	if b := c.Minimum; b != nil {
		if c.ExclusiveMinimum && f <= *b {
			errs = append(errs, field.Invalid(fp, f, fmt.Sprintf("must be greater than %v", *b)).WithOrigin("minimum"))
		} else if f < *b {
			errs = append(errs, field.Invalid(fp, f, fmt.Sprintf("must be greater than or equal to %v", *b)).WithOrigin("minimum"))
		}
	}
	if b := c.Maximum; b != nil {
		if c.ExclusiveMaximum && f >= *b {
			errs = append(errs, field.Invalid(fp, f, fmt.Sprintf("must be less than %v", *b)).WithOrigin("maximum"))
		} else if f > *b {
			errs = append(errs, field.Invalid(fp, f, fmt.Sprintf("must be less than or equal to %v", *b)).WithOrigin("maximum"))
		}
	}
	if m := c.MultipleOf; m != nil && *m != 0 && math.Mod(f, *m) != 0 {
		errs = append(errs, field.Invalid(fp, f, fmt.Sprintf("must be a multiple of %v", *m)).WithOrigin("multipleOf"))
	}
	return errs
}

// checkUnion returns the errors of a union: exactly one member set, or at
// most one. For an update, a union whose members are set as they were isn't
// checked, as declarative validation ratchets it.
func checkUnion(ctx context.Context, fp *field.Path, members []string, oneOf bool, isSet func(string) bool, old *node) field.ErrorList {
	// the extractors' value says whether to read the old die, or this one
	var um []validate.UnionMember
	var fns []validate.ExtractorFn[bool, bool]
	for _, m := range members {
		um = append(um, validate.NewUnionMember(m))
		fns = append(fns, func(inOld bool) bool {
			if inOld {
				return old.has(m)
			}
			return isSet(m)
		})
	}
	u := validate.NewUnionMembership(um...)
	op := create
	if old != nil {
		op = operation.Operation{Type: operation.Update}
	}
	if oneOf {
		return validate.Union(ctx, op, fp, false, old != nil, u, fns...)
	}
	return validate.ZeroOrOneOfUnion(ctx, op, fp, false, old != nil, u, fns...)
}

// checkUpdate returns the errors of an update of a field that changed:
// +k8s:immutable, and +k8s:update's constraints, with declarative
// validation's errors.
func checkUpdate(fd *core.FieldDesc, fp *field.Path, present bool, value reflect.Value, oldPresent bool, oldValue reflect.Value) field.ErrorList {
	c := fd.Checks
	var errs field.ErrorList
	if c.Immutable {
		errs = append(errs, field.Invalid(fp, nil, "field is immutable").WithOrigin("immutable"))
	}
	set, wasSet := present && !empty(value), oldPresent && !empty(oldValue)
	for _, u := range c.Update {
		switch u {
		case "NoSet":
			if !wasSet && set {
				errs = append(errs, field.Invalid(fp, nil, "field cannot be set once created").WithOrigin("update"))
			}
		case "NoUnset":
			if wasSet && !set {
				errs = append(errs, field.Invalid(fp, nil, "field cannot be cleared once set").WithOrigin("update"))
			}
		case "NoModify":
			if wasSet && set {
				errs = append(errs, field.Invalid(fp, nil, "field cannot be modified once set").WithOrigin("update"))
			}
		case "NoAddItem", "NoRemoveItem":
			errs = append(errs, checkItems(fd, fp, u, value, oldValue)...)
		}
	}
	return errs
}

// checkItems returns the errors of items added to, or removed from, a list
// or map: a list map's elements are matched by key, a map's entries by key,
// and another list's items by value.
func checkItems(fd *core.FieldDesc, fp *field.Path, constraint string, value, oldValue reflect.Value) field.ErrorList {
	for _, v := range []*reflect.Value{&value, &oldValue} {
		for v.Kind() == reflect.Pointer && !v.IsNil() {
			*v = v.Elem()
		}
	}
	var errs field.ErrorList
	switch value.Kind() {
	case reflect.Map:
		in := func(m reflect.Value, k reflect.Value) bool { return m.IsValid() && m.MapIndex(k).IsValid() }
		if constraint == "NoAddItem" {
			for _, k := range sortedKeys(value) {
				if !in(oldValue, k) {
					errs = append(errs, field.Forbidden(fp.Key(k.String()), "item may not be added").WithOrigin("update"))
				}
			}
		} else if oldValue.Kind() == reflect.Map {
			for _, k := range sortedKeys(oldValue) {
				if !in(value, k) {
					errs = append(errs, field.Forbidden(fp.Key(k.String()), "item may not be removed").WithOrigin("update"))
				}
			}
		}
	case reflect.Slice:
		find := func(list reflect.Value, item reflect.Value) bool {
			if !list.IsValid() || list.Kind() != reflect.Slice {
				return false
			}
			if len(fd.Keys) > 0 {
				i, _ := core.FindElem(list, core.KeyStepOf(list.Type().Elem(), fd.Keys, item))
				return i >= 0
			}
			for i := range list.Len() {
				if reflect.DeepEqual(list.Index(i).Interface(), item.Interface()) {
					return true
				}
			}
			return false
		}
		if constraint == "NoAddItem" {
			for i := range value.Len() {
				if !find(oldValue, value.Index(i)) {
					errs = append(errs, field.Forbidden(fp.Index(i), "item may not be added").WithOrigin("update"))
				}
			}
		} else if oldValue.Kind() == reflect.Slice {
			for i := range oldValue.Len() {
				if !find(value, oldValue.Index(i)) {
					errs = append(errs, field.Forbidden(fp, "item may not be removed").WithOrigin("update"))
					break
				}
			}
		}
	}
	return errs
}

func sortedKeys(m reflect.Value) []reflect.Value {
	keys := m.MapKeys()
	slices.SortFunc(keys, func(a, b reflect.Value) int { return strings.Compare(a.String(), b.String()) })
	return keys
}
