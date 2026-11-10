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

package apis

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VolatileTime is declared by hand instead of cloned: upstream wraps a
// [metav1.Time] and encodes exactly as it does.
type VolatileTime = metav1.Time
