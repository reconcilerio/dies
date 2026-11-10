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

package cmp

import (
	"github.com/google/go-cmp/cmp"

	"reconciler.io/dies/v2/die/base"
)

// Transform is a go-cmp option that compares dies by their [base.Content]:
// their current values, as JSON has them, and their intent. Without it,
// go-cmp compares dies by their Equal methods, which agree, but reports a
// difference as the dies' internals. With it, a diff shows the fields and
// the intent that differ, and is empty exactly when the dies are Equal.
//
// It applies to dies wherever they are within the values compared, such as
// the items of a list die or the fields of a test case.
//
// A die's content is not its Go type, so options for Kubernetes types, such
// as cmpopts.IgnoreFields(metav1.ObjectMeta{}, "ResourceVersion") or a
// comparer for [k8s.io/apimachinery/pkg/api/resource.Quantity], don't apply
// to it. Options for the content, such as those of
// [github.com/google/go-cmp/cmp/cmpopts], do:
//
//   - cmpopts.IgnoreFields(base.Content{}, "Intent") compares values only.
//   - cmpopts.IgnoreMapEntries(func(k string, _ any) bool { return k ==
//     "resourceVersion" }) ignores a field by its JSON name, at any depth;
//     [cmp.FilterPath] narrows an option to some paths.
//
// To use options for the Go types instead, use [TransformRelease].
var Transform = cmp.Transformer("die.Content", base.MustExtractContentOf)

// TransformRelease is a go-cmp option that compares dies by their current
// values as their Go types have them, as DieRelease returns them, so that
// options for those types, such as cmpopts.IgnoreFields on a Kubernetes type,
// apply as they do to the types themselves. Intent, and fields the types
// don't have, are not compared, so dies whose diff is empty may not be
// Equal. Use it in place of [Transform], not with it: go-cmp panics when two
// transformers apply to the same value.
var TransformRelease = cmp.Transformer("die.Release", base.ReleaseOf)

// Diff returns a human-readable report of the differences between two
// values, as [cmp.Diff] does with [Transform] and opts: empty when they are
// equal.
func Diff(x, y any, opts ...cmp.Option) string {
	return cmp.Diff(x, y, append([]cmp.Option{Transform}, opts...)...)
}

// DiffRelease is [Diff], comparing dies by their Go values with
// [TransformRelease], so that opts for those types apply.
func DiffRelease(x, y any, opts ...cmp.Option) string {
	return cmp.Diff(x, y, append([]cmp.Option{TransformRelease}, opts...)...)
}

// Equal returns true when [Diff] would be empty for two values, as
// [cmp.Equal] does with [Transform] and opts.
func Equal(x, y any, opts ...cmp.Option) bool {
	return cmp.Equal(x, y, append([]cmp.Option{Transform}, opts...)...)
}

// EqualRelease returns true when [DiffRelease] would be empty for two values,
// as [cmp.Equal] does with [TransformRelease] and opts.
func EqualRelease(x, y any, opts ...cmp.Option) bool {
	return cmp.Equal(x, y, append([]cmp.Option{TransformRelease}, opts...)...)
}
