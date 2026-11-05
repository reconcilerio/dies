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

package core

import (
	"reflect"
	"slices"
	"sync"
)

var (
	ManagedMu sync.Mutex
	Managed   = map[reflect.Type][]*FieldDesc{}
)

type registered struct {
	annotateFn, tablesFn func()
	annotate, tables     sync.Once
}

var (
	packagesMu sync.Mutex
	packages   = map[string][]*registered{}
)

// registerPackage registers what to do when the types of the packages at
// paths are first used.
func RegisterPackage(paths []string, annotate, tables func()) {
	r := &registered{annotateFn: annotate, tablesFn: tables}
	packagesMu.Lock()
	defer packagesMu.Unlock()
	for _, path := range paths {
		packages[path] = append(packages[path], r)
	}
}

func registeredFor(t reflect.Type) []*registered {
	packagesMu.Lock()
	defer packagesMu.Unlock()
	return packages[t.PkgPath()]
}

// annotate runs the annotations of the die packages for t's package.
func annotate(t reflect.Type) {
	for _, r := range registeredFor(t) {
		if r.annotateFn != nil {
			r.annotate.Do(r.annotateFn)
		}
	}
}

// built builds the field tables of the die packages for t's package.
func Built(t reflect.Type) {
	for _, r := range registeredFor(t) {
		if r.tablesFn != nil {
			r.tables.Do(r.tablesFn)
		}
	}
}

// identity is an object type's apiVersion and kind.
type identity struct{ apiVersion, kind string }

var identities sync.Map

// SetIdentity records an object type's apiVersion and kind.
func SetIdentity(t reflect.Type, apiVersion, kind string) {
	identities.Store(t, identity{apiVersion: apiVersion, kind: kind})
}

// identityOf returns an object type's apiVersion and kind, if its die package
// recorded them.
func identityOf(t reflect.Type) (identity, bool) {
	annotate(t)
	id, ok := identities.Load(t)
	if !ok {
		return identity{}, false
	}
	return id.(identity), true
}

func Manage(t reflect.Type, fd *FieldDesc) {
	ManagedMu.Lock()
	defer ManagedMu.Unlock()
	if !slices.Contains(Managed[t], fd) {
		Managed[t] = append(Managed[t], fd)
	}
}

// ManagedFields returns the fields the die for t manages, its field tables
// built.
func ManagedFields(t reflect.Type) []*FieldDesc {
	Built(t)
	ManagedMu.Lock()
	defer ManagedMu.Unlock()
	return slices.Clone(Managed[t])
}
