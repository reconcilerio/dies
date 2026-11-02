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

package fields

import (
	"sync"
)

// Table is a die's field table, built by Build on first use. Generated dies
// declare each as a package variable whose Build captures nothing, which the
// linker lays out statically, so a die package does no work at init.
type Table[F any] struct {
	Build func() F
	once  sync.Once
	f     F
}

// Get returns the table, building it the first time.
func (t *Table[F]) Get() *F {
	t.once.Do(t.build)
	return &t.f
}

func (t *Table[F]) build() { t.f = t.Build() }
