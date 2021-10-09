/*
Copyright 2021 the original author or authors.

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

package v1

import (
	corev1 "k8s.io/api/core/v1"
)

// +die:target=k8s.io/api/core/v1.Secret,object=true
// +die:field:receiver=SecretDie,name=Data,type=map[string][]byte

func (d *SecretDie) AddData(key string, value []byte) *SecretDie {
	return d.DieStamp(func(r *corev1.Secret) {
		if r.Data == nil {
			r.Data = map[string][]byte{}
		}
		r.Data[key] = value
	})
}
