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

// Package schemes registers the types of an API's kinds, their immutable
// dies and their mutable dies, each in a scheme of its own: a
// [runtime.Scheme] maps a kind to one Go type, so the three can't share
// one. A [Composite] holds the three, and answers for any of them: the kinds
// of an object of any category, the scheme it's registered in, to pass to
// clients and tools that take one, and an object of each category for a
// kind.
//
// A die package's SchemeBuilder is a [Builder], whose AddToScheme registers
// its mutable dies, as a client decodes into them. AddToSchemeAs registers a
// category of the package's, and AddToComposite each:
//
//	cs := schemes.NewComposite()
//	if err := dieappsv1.SchemeBuilder.AddToComposite(cs); err != nil { ... }
//	d, err := cs.New(dieappsv1.SchemeGroupVersion.WithKind("Deployment"), schemes.DieMuts)
package schemes

import (
	"errors"
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Category is a category of the Go types of a kind.
type Category int

const (
	// Types are the Go types the dies are of, such as appsv1.Deployment.
	// A type that isn't a runtime.Object, as a cloned type isn't, isn't
	// registered.
	Types Category = iota
	// Dies are the immutable dies, such as *DeploymentDie. A list die,
	// which is one type, is registered as the list's die of both
	// categories of dies.
	Dies
	// DieMuts are the mutable dies, such as *DeploymentDieMut, which a
	// client decodes into.
	DieMuts
)

// categories are every category, in the order a [Composite] looks for the
// scheme of an object in.
var categories = [...]Category{Types, DieMuts, Dies}

func (c Category) String() string {
	switch c {
	case Types:
		return "Types"
	case Dies:
		return "Dies"
	case DieMuts:
		return "DieMuts"
	}
	return fmt.Sprintf("Category(%d)", int(c))
}

// valid reports whether c is a category.
func (c Category) valid() bool { return c >= Types && c <= DieMuts }

// ErrCategory is the error for a category that isn't one.
var ErrCategory = errors.New("schemes: not a category")

// CheckCategory returns an error, wrapping [ErrCategory], unless c is a
// category. A registration of a category the package doesn't know returns
// it.
func CheckCategory(c Category) error {
	if !c.valid() {
		return fmt.Errorf("%w: %d", ErrCategory, int(c))
	}
	return nil
}

// AddFunc registers the types of a category with a scheme.
type AddFunc func(s *runtime.Scheme, c Category) error

// Builder registers the types of each category, as a [runtime.SchemeBuilder]
// registers one: a die package's SchemeBuilder is one, and a group's or an
// API's is made of its packages', by their Add methods:
//
//	var SchemeBuilder = schemes.NewBuilder(appsv1.SchemeBuilder.Add, appsv1beta1.SchemeBuilder.Add)
type Builder []AddFunc

// NewBuilder returns a builder of the functions.
func NewBuilder(fns ...AddFunc) Builder { return Builder(fns) }

// Register adds functions to the builder.
func (b *Builder) Register(fns ...AddFunc) { *b = append(*b, fns...) }

// Add registers the types of category c with a scheme. It is an
// [AddFunc], so that a builder is made of others.
func (b Builder) Add(s *runtime.Scheme, c Category) error {
	if err := CheckCategory(c); err != nil {
		return err
	}
	for _, fn := range b {
		if err := fn(s, c); err != nil {
			return err
		}
	}
	return nil
}

// AddToScheme registers the mutable dies with a scheme, as clients decode
// into them: it is AddToSchemeAs(DieMuts), as a die package's AddToScheme
// has always been.
func (b Builder) AddToScheme(s *runtime.Scheme) error { return b.Add(s, DieMuts) }

// AddToSchemeAs returns a function that registers the types of category c
// with a scheme, of the form of a [runtime.SchemeBuilder]'s AddToScheme, for
// clients and tools that take one.
func (b Builder) AddToSchemeAs(c Category) func(*runtime.Scheme) error {
	return func(s *runtime.Scheme) error { return b.Add(s, c) }
}

// AddToComposite registers the types of every category with a composite
// scheme, each with its scheme.
func (b Builder) AddToComposite(cs *Composite) error {
	for _, c := range categories {
		if err := b.Add(cs.Scheme(c), c); err != nil {
			return err
		}
	}
	return nil
}

// Composite is a scheme of each category, for the types of the same kinds:
// the Go types, the immutable dies and the mutable dies. Each is a
// [runtime.Scheme], to pass to what takes one. It is a [runtime.ObjectTyper]
// of objects of any category.
type Composite struct {
	schemes [3]*runtime.Scheme
}

var _ runtime.ObjectTyper = (*Composite)(nil)

// NewComposite returns a composite scheme of empty schemes.
func NewComposite() *Composite {
	cs := &Composite{}
	for _, c := range categories {
		cs.schemes[c] = runtime.NewScheme()
	}
	return cs
}

// Scheme returns the scheme of category c. It panics if c isn't one.
func (cs *Composite) Scheme(c Category) *runtime.Scheme {
	if err := CheckCategory(c); err != nil {
		panic(err)
	}
	return cs.schemes[c]
}

// SchemeFor returns the scheme an object's type is registered in, and its
// category, to pass to a client or tool that takes a scheme, for it to work
// with objects of that type. A type registered in more than one, such as a
// list die, which is the list's die of both categories of dies, or
// apimachinery's meta types, which each scheme has, is found in the first
// that has it of Types, DieMuts and Dies. An unstructured object's is the
// scheme of Types, as a scheme takes an unstructured object's kind from the
// object.
func (cs *Composite) SchemeFor(obj runtime.Object) (*runtime.Scheme, Category, error) {
	for _, c := range categories {
		if _, _, err := cs.schemes[c].ObjectKinds(obj); err == nil {
			return cs.schemes[c], c, nil
		}
	}
	return nil, 0, runtime.NewNotRegisteredErrForType("schemes.Composite", reflect.TypeOf(obj))
}

// ObjectKinds returns the kinds of an object of any category, as the scheme
// it's registered in returns them.
func (cs *Composite) ObjectKinds(obj runtime.Object) ([]schema.GroupVersionKind, bool, error) {
	s, _, err := cs.SchemeFor(obj)
	if err != nil {
		return nil, false, err
	}
	return s.ObjectKinds(obj)
}

// Recognizes reports whether a kind is registered in any category.
func (cs *Composite) Recognizes(gvk schema.GroupVersionKind) bool {
	for _, c := range categories {
		if cs.schemes[c].Recognizes(gvk) {
			return true
		}
	}
	return false
}

// New returns a new object of a kind's type of category c: its Go type, its
// immutable die or its mutable die.
func (cs *Composite) New(gvk schema.GroupVersionKind, c Category) (runtime.Object, error) {
	if err := CheckCategory(c); err != nil {
		return nil, err
	}
	return cs.schemes[c].New(gvk)
}

// TypeOf returns a kind's type of category c, or false if it has none.
func (cs *Composite) TypeOf(gvk schema.GroupVersionKind, c Category) (reflect.Type, bool) {
	if !c.valid() {
		return nil, false
	}
	t, ok := cs.schemes[c].AllKnownTypes()[gvk]
	return t, ok
}
