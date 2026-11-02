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

// Package v1 has dies for knative's duck/v1, whose
// types are cloned so that nothing depends on knative.
//
// +die:types:package=knative.dev/pkg/apis/duck/v1,mode=clone,from=../../sources
// +die:import:package=knative.dev/pkg/apis,as=reconciler.io/dies/v2/examples/knative/apis
// +die:import:package=knative.dev/pkg/tracker,as=reconciler.io/dies/v2/examples/knative/tracker
package v1
