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

package codegen

import (
	"encoding/json"

	utiljson "k8s.io/apimachinery/pkg/util/json"

	"reconciler.io/dies/v2/die/base"
	"reconciler.io/dies/v2/die/internal/core"
	"reconciler.io/dies/v2/die/internal/dies"
)

// Dies of types with apply configurations, such as client-go's, have
// DieReleaseApplyConfiguration, DieFeedApplyConfiguration and
// DieStampApplyConfiguration, which call these.

// ReleaseApplyConfiguration returns a die's intent as an apply configuration
// of type AC: on a root, what an apply sends, the object's identity
// included; on a die within another, its intent.
func ReleaseApplyConfiguration[AC any](d base.Die) *AC {
	return core.Must(ReleaseApplyConfigurationE[AC](d))
}

// ReleaseApplyConfigurationE is ReleaseApplyConfiguration, returning an error
// rather than panicking.
func ReleaseApplyConfigurationE[AC any](d base.Die) (*AC, error) {
	h, at := dies.MustStateOf(d)
	data, err := core.IntentJSON(h, at)
	if err != nil {
		return nil, core.Errorf("die: encoding the intent of %T: %w", d, err)
	}
	ac := new(AC)
	if err := json.Unmarshal(data, ac); err != nil {
		return nil, core.Errorf("die: decoding the intent of %T as a %T: %w", d, ac, err)
	}
	return ac, nil
}

// FeedApplyConfiguration replaces a die's intent with an apply
// configuration's: each value it sets becomes intent, over the die's value.
// A nil apply configuration clears the intent. As for any write, a root that
// is immutable returns a new die, and a die within an immutable one fails.
func FeedApplyConfiguration[D base.Die](d D, ac any) D {
	return core.Must(FeedApplyConfigurationE(d, ac))
}

// FeedApplyConfigurationE is FeedApplyConfiguration, returning an error
// rather than panicking.
func FeedApplyConfigurationE[D base.Die](d D, ac any) (D, error) {
	var zero D
	h, at := dies.MustStateOf(d)
	data, err := json.Marshal(ac)
	if err != nil {
		return zero, core.Errorf("die: encoding a %T: %w", ac, err)
	}
	var raw any
	// numbers decode as int64 where they can, as unstructured values hold them
	if err := utiljson.Unmarshal(data, &raw); err != nil {
		return zero, core.Errorf("die: decoding a %T: %w", ac, err)
	}
	nh, err := core.FeedIntent(h, at, raw)
	if err != nil {
		return zero, err
	}
	if nh != h {
		return dies.Wrap[D](nh, nil), nil
	}
	return dies.Wrap[D](h, at), nil
}

// StampApplyConfiguration passes a die's intent, as an apply configuration,
// to fn, and feeds back what fn makes of it.
func StampApplyConfiguration[AC any, D base.Die](d D, fn func(ac *AC)) D {
	return core.Must(StampApplyConfigurationE(d, fn))
}

// StampApplyConfigurationE is StampApplyConfiguration, returning an error
// rather than panicking.
func StampApplyConfigurationE[AC any, D base.Die](d D, fn func(ac *AC)) (D, error) {
	ac, err := ReleaseApplyConfigurationE[AC](d)
	if err != nil {
		var zero D
		return zero, err
	}
	fn(ac)
	return FeedApplyConfigurationE(d, ac)
}
