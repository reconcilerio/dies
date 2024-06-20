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
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	diemetav1 "reconciler.io/dies/apis/meta/v1"
)

// +die:object=true,apiVersion=apiextensions.k8s.io/v1,kind=CustomResourceDefinition
type _ = apiextensionsv1.CustomResourceDefinition

// +die
// +die:field:name=Names,die=CustomResourceDefinitionNamesDie
// +die:field:name=Conversion,die=CustomResourceConversionDie,pointer=true
// +die:field:name=Versions,die=CustomResourceDefinitionVersionDie,listType=map
type _ = apiextensionsv1.CustomResourceDefinitionSpec

// +die
// +die:field:name=Schema,die=CustomResourceValidationDie,pointer=true
// +die:field:name=Subresources,die=CustomResourceSubresourcesDie,pointer=true
// +die:field:name=AdditionalPrinterColumns,die=CustomResourceColumnDefinitionDie,listType=map
// +die:field:name=SelectableFields,die=SelectableFieldDie,listType=atomic
type _ apiextensionsv1.CustomResourceDefinitionVersion

// +die
type _ apiextensionsv1.CustomResourceValidation

// +die
// +die:field:name=Scale,die=CustomResourceSubresourceScaleDie,pointer=true
type _ apiextensionsv1.CustomResourceSubresources

// +die
type _ apiextensionsv1.CustomResourceSubresourceScale

// +die
type _ apiextensionsv1.CustomResourceColumnDefinition

// +die
// +die:field:name=Webhook,die=WebhookConversionDie,pointer=true
type _ apiextensionsv1.CustomResourceConversion

// +die
// +die:field:name=ClientConfig,die=WebhookClientConfigDie,pointer=true
type _ apiextensionsv1.WebhookConversion

// +die
// +die:field:name=Service,die=ServiceReferenceDie,pointer=true
type _ apiextensionsv1.WebhookClientConfig

// +die
type _ apiextensionsv1.ServiceReference

// +die
type _ apiextensionsv1.SelectableField

// +die
// +die:field:name=AcceptedNames,die=CustomResourceDefinitionNamesDie
type _ = apiextensionsv1.CustomResourceDefinitionStatus

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

// +die
type _ = apiextensionsv1.CustomResourceDefinitionNames
