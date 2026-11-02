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

	"k8s.io/apimachinery/pkg/runtime/schema"
)

type ObjectKind struct {
	H  *Holder
	At Path
}

// GroupVersionKind implements schema.ObjectKind, which can't return an
// error, so it panics with one, as a read through a jig can fail.
func (k *ObjectKind) GroupVersionKind() schema.GroupVersionKind {
	apiVersion, err := GetAt[string](k.H, k.At, Path{FieldStep("apiVersion")})
	if err != nil {
		panic(err)
	}
	kind, err := GetAt[string](k.H, k.At, Path{FieldStep("kind")})
	if err != nil {
		panic(err)
	}
	return schema.FromAPIVersionAndKind(apiVersion, kind)
}

// SetGroupVersionKind sets the type metadata in the base, not as intent:
// libraries set it, and strip it, as they encode and decode, such as a
// client decoding into a typed object, and it is not something to patch. It
// implements schema.ObjectKind, which can't return an error, so it panics
// with one, on an immutable die.
func (k *ObjectKind) SetGroupVersionKind(gvk schema.GroupVersionKind) {
	if err := k.setGroupVersionKind(gvk); err != nil {
		panic(err)
	}
}

func (k *ObjectKind) setGroupVersionKind(gvk schema.GroupVersionKind) error {
	h := k.H
	if !h.Mutable {
		return Errorf("die: SetGroupVersionKind on an immutable die, which would lose the change; use a mutable die, such as within DieWith")
	}
	apiVersion, kind := gvk.ToAPIVersionAndKind()
	var paths []Path
	for field, v := range map[string]string{"apiVersion": apiVersion, "kind": kind} {
		// stripped, it is absent, as it would be had it not been on the wire
		op := SetOp(OpSet, reflect.ValueOf(v))
		if v == "" {
			op = SetOp(OpClear, reflect.Value{})
		}
		if _, err := h.Write(k.At, Path{FieldStep(field)}, op, true); err != nil {
			return err
		}
		paths = append(paths, k.At.join(Path{FieldStep(field)}))
	}
	fold := Only(h.Ov, paths)
	base, err := (&Holder{Base: h.Base, Ov: fold}).Materialize()
	if err != nil {
		return err
	}
	h.Base, h.Wire, h.Ov = base, CommitWire(h.Wire, fold), Except(h.Ov, paths)
	return nil
}

// SetAt sets a value in place, which requires a mutable die.
func SetAt[V any](h *Holder, at Path, method string, rel Path, v V) error {
	if !h.Mutable {
		return Errorf("die: %s on an immutable die, which would lose the change; use a mutable die, such as within DieWith", method)
	}
	_, err := h.Write(at, rel, SetOp(OpSet, DeepCopy(reflect.ValueOf(&v).Elem())), true)
	return err
}

// GetAt reads a value as V.
func GetAt[V any](h *Holder, at Path, rel Path) (V, error) {
	var zero V
	v, err := h.ValueAt(at.join(rel))
	if err != nil {
		return zero, err
	}
	if t := reflect.TypeFor[V](); v.Type() != t {
		// an unstructured value, read as V
		if v.Kind() == reflect.Interface && v.IsNil() {
			return zero, nil
		}
		uv, err := fromUnstructured(v.Interface(), t)
		if err != nil {
			return zero, err
		}
		return uv.Interface().(V), nil
	}
	return v.Interface().(V), nil
}
