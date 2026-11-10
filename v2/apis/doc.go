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

//go:generate go tool diegen -header ../../hack/boilerplate.go.txt ./...

// Package apis registers the dies of every API group. Registered dies are
// the Go types of their kinds, in place of the upstream types, so a scheme
// should not also have the upstream types of the same group versions.
package apis
