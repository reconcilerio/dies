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

package json

import (
	"encoding/json"

	"gomodules.xyz/jsonpatch/v3"
	"k8s.io/apimachinery/pkg/types"
	"reconciler.io/dies/patch"
)

func init() {
	patch.Register(types.JSONPatchType, JSONPatchFunc)
}

func JSONPatchFunc(original, modified interface{}) ([]byte, error) {
	originalJSON, err := json.Marshal(original)
	if err != nil {
		return nil, err
	}

	modifiedJSON, err := json.Marshal(modified)
	if err != nil {
		return nil, err
	}

	ops, err := jsonpatch.CreatePatch(originalJSON, modifiedJSON)
	if err != nil {
		return nil, err
	}
	return json.Marshal(ops)
}