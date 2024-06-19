/*
Copyright 2024 the original author or authors.

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

package patch

import (
	"fmt"

	"k8s.io/apimachinery/pkg/types"
)

type PatchFunc = func(original, modified interface{}) ([]byte, error)

type patchFuncRegistry map[types.PatchType]PatchFunc

func (r patchFuncRegistry) Register(patchType types.PatchType, patchFunc PatchFunc) {
	r[patchType] = patchFunc
}

func (r patchFuncRegistry) Lookup(patchType types.PatchType) (PatchFunc, error) {
	if patchFunc := r[patchType]; patchFunc != nil {
		return patchFunc, nil
	}

	// provide hints for how to register known patch types
	var importPath string
	switch patchType {
	case types.JSONPatchType:
		importPath = "reconciler.io/dies/patch/json"
	case types.MergePatchType:
		importPath = "reconciler.io/dies/patch/merge"
	case types.StrategicMergePatchType:
		importPath = "reconciler.io/dies/patch/strategicmerge"
	case types.ApplyPatchType:
		importPath = "reconciler.io/dies/patch/apply"
	}

	if importPath != "" {
		return nil, fmt.Errorf("unknown patch type %q, load support with `import _ %q`", patchType, importPath)
	}
	return nil, fmt.Errorf("unknown patch type %q", patchType)
}

var registry = patchFuncRegistry{}

func Register(patchType types.PatchType, patchFunc PatchFunc) {
	registry.Register(patchType, patchFunc)
}

// Create a patch that will convert the original value into the modified value based on the
// specified patch type. Both values must marshal to JSON and should be of the same type.
func Create(original, modified interface{}, patchType types.PatchType) ([]byte, error) {
	patchFunc, err := registry.Lookup(patchType)
	if err != nil {
		return nil, err
	}

	return patchFunc(original, modified)
}
