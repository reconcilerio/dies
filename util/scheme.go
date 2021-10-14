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
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func Register(schemeBuilder interface{}, gvk schema.GroupVersionKind, objs ...runtime.Object) {
	// handle runtime.SchemeBuilder
	if sb, ok := schemeBuilder.(runtime.SchemeBuilder); ok {
		sb.Register(func(s *runtime.Scheme) error {
			for _, o := range objs {
				s.AddKnownTypeWithName(gvk, o)
			}
			return nil
		})
		return
	}

	// handle kubebuilder's SchemeBuilder, which has a runtime.SchemeBuilder as a field
	// using reflection to avoid a hard dependency on controller-runtime
	sbv := reflect.ValueOf(schemeBuilder)
	if sbv.Kind() == reflect.Ptr {
		sbf := sbv.Elem().FieldByName("SchemeBuilder")
		if sbf.Type().AssignableTo(reflect.TypeOf(runtime.SchemeBuilder{})) {
			Register(sbf.Interface(), gvk, objs...)
			return
		}
	}

	panic(fmt.Errorf("unknown SchemeBuilder of type %q", sbv.Type()))
}
