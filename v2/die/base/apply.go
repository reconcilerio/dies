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
	"encoding/json"

	"k8s.io/apimachinery/pkg/runtime"

	"reconciler.io/dies/v2/die/internal/core"
)

// ApplyConfiguration is a die's intent as an apply configuration, for clients
// that take one, such as controller-runtime's client.Apply: what
// DieReleaseApplyPatch sends, the object's identity included. Any root die has one,
// with or without an apply configuration type of its own.
//
// It is the intent when it was made, not after. A client decodes what the
// server returns into it; Result has that.
type ApplyConfiguration struct {
	patch    []byte
	identity struct {
		APIVersion *string `json:"apiVersion"`
		Kind       *string `json:"kind"`
		Metadata   struct {
			Name      *string `json:"name"`
			Namespace *string `json:"namespace"`
		} `json:"metadata"`
	}
	result []byte
}

var _ runtime.ApplyConfiguration = (*ApplyConfiguration)(nil)

// ApplyConfigurationOf returns a root die's intent as an apply configuration.
// It panics on a die within another, as DieReleaseApplyPatch does.
func ApplyConfigurationOf(d Die) *ApplyConfiguration {
	return core.Must(ApplyConfigurationOfE(d))
}

// ApplyConfigurationOfE is ApplyConfigurationOf, returning an error rather
// than panicking.
func ApplyConfigurationOfE(d Die) (*ApplyConfiguration, error) {
	h, at, ok := core.StateOf(d)
	if !ok {
		return nil, core.Errorf("die: %T is not a die", d)
	}
	if len(at) > 0 {
		return nil, core.Errorf("die: ApplyConfigurationOf is only supported on a root die, not at %s", at)
	}
	patch, err := h.ApplyPatch()
	if err != nil {
		return nil, core.Errorf("die: encoding the intent of %T: %w", d, err)
	}
	ac := &ApplyConfiguration{patch: patch}
	if err := json.Unmarshal(patch, &ac.identity); err != nil {
		return nil, core.Errorf("die: decoding the identity of %T: %w", d, err)
	}
	return ac, nil
}

func (*ApplyConfiguration) IsApplyConfiguration() {}

func (ac *ApplyConfiguration) GetAPIVersion() *string { return ac.identity.APIVersion }
func (ac *ApplyConfiguration) GetKind() *string       { return ac.identity.Kind }
func (ac *ApplyConfiguration) GetName() *string       { return ac.identity.Metadata.Name }
func (ac *ApplyConfiguration) GetNamespace() *string  { return ac.identity.Metadata.Namespace }

// MarshalJSON encodes the intent, as an apply request's body.
func (ac *ApplyConfiguration) MarshalJSON() ([]byte, error) { return ac.patch, nil }

// UnmarshalJSON keeps what a client decodes into the apply configuration,
// the object as the server returned it, for Result.
func (ac *ApplyConfiguration) UnmarshalJSON(data []byte) error {
	ac.result = append([]byte(nil), data...)
	return nil
}

// Result is the object as the server returned it from the apply, as JSON, or
// nil before it is applied: to feed a die, d.DieFeedJSON(ac.Result()).
func (ac *ApplyConfiguration) Result() []byte { return ac.result }
