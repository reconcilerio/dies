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
	diemetav1 "github.com/scothis/dies/apis/meta/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

// +die:target=k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1.CustomResourceDefinition,object=true

// +die:target=k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1.CustomResourceDefinitionSpec
// +die:field:receiver=CustomResourceDefinitionSpecDie,name=Group,type=string
// +die:field:receiver=CustomResourceDefinitionSpecDie,name=Names,type=k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1.CustomResourceDefinitionNames
// +die:field:receiver=CustomResourceDefinitionSpecDie,name=Scope,type=k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1.ResourceScope
// +die:field:receiver=CustomResourceDefinitionSpecDie,name=Versions,type=[]k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1.CustomResourceDefinitionVersion
// +die:field:receiver=CustomResourceDefinitionSpecDie,name=Conversion,type=*k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1.CustomResourceConversion
// +die:field:receiver=CustomResourceDefinitionSpecDie,name=PreserveUnknownFields,type=bool

// +die:target=k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1.CustomResourceDefinitionStatus
// +die:field:receiver=CustomResourceDefinitionStatusDie,name=Conditions,type=[]k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1.CustomResourceDefinitionCondition
// +die:field:receiver=CustomResourceDefinitionStatusDie,name=AcceptedNames,type=k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1.CustomResourceDefinitionNames
// +die:field:receiver=CustomResourceDefinitionStatusDie,name=StoredVersions,type=[]string

func (d *CustomResourceDefinitionStatusDie) ConditionsDie(conditions ...*diemetav1.ConditionDie) *CustomResourceDefinitionStatusDie {
	return d.DieStamp(func(r *apiextensionsv1.CustomResourceDefinitionStatus) {
		r.Conditions = make([]apiextensionsv1.CustomResourceDefinitionCondition, len(conditions))
		for i := range conditions {
			c := conditions[i].DieRelease()
			r.Conditions[i] = apiextensionsv1.CustomResourceDefinitionCondition{
				Type:               apiextensionsv1.CustomResourceDefinitionConditionType(c.Type),
				Status:             apiextensionsv1.ConditionStatus(c.Status),
				Reason:             c.Reason,
				Message:            c.Message,
				LastTransitionTime: c.LastTransitionTime,
			}
		}
	})
}
