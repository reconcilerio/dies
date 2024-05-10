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
type _ = apiextensionsv1.CustomResourceDefinitionSpec

func (d *CustomResourceDefinitionSpecDie) NamesDie(fn func(d *CustomResourceDefinitionNamesDie)) *CustomResourceDefinitionSpecDie {
	return d.DieStamp(func(r *apiextensionsv1.CustomResourceDefinitionSpec) {
		d := CustomResourceDefinitionNamesBlank.DieImmutable(false).DieFeed(r.Names)
		fn(d)
		r.Names = d.DieRelease()
	})
}

func (d *CustomResourceDefinitionSpecDie) VersionDie(name string, fn func(d *CustomResourceDefinitionVersionDie)) *CustomResourceDefinitionSpecDie {
	return d.DieStamp(func(r *apiextensionsv1.CustomResourceDefinitionSpec) {
		for i := range r.Versions {
			if name == r.Versions[i].Name {
				d := CustomResourceDefinitionVersionBlank.DieImmutable(false).DieFeed(r.Versions[i])
				fn(d)
				r.Versions[i] = d.DieRelease()
				return
			}
		}

		d := CustomResourceDefinitionVersionBlank.DieImmutable(false).DieFeed(apiextensionsv1.CustomResourceDefinitionVersion{Name: name})
		fn(d)
		r.Versions = append(r.Versions, d.DieRelease())
	})
}

func (d *CustomResourceDefinitionSpecDie) ConversionDie(fn func(d *CustomResourceConversionDie)) *CustomResourceDefinitionSpecDie {
	return d.DieStamp(func(r *apiextensionsv1.CustomResourceDefinitionSpec) {
		d := CustomResourceConversionBlank.DieImmutable(false).DieFeedPtr(r.Conversion)
		fn(d)
		r.Conversion = d.DieReleasePtr()
	})
}

// +die
type _ apiextensionsv1.CustomResourceDefinitionVersion

func (d *CustomResourceDefinitionVersionDie) SchemaDie(fn func(d *CustomResourceValidationDie)) *CustomResourceDefinitionVersionDie {
	return d.DieStamp(func(r *apiextensionsv1.CustomResourceDefinitionVersion) {
		d := CustomResourceValidationBlank.DieImmutable(false).DieFeedPtr(r.Schema)
		fn(d)
		r.Schema = d.DieReleasePtr()
	})
}

func (d *CustomResourceDefinitionVersionDie) SubresourcesDie(fn func(d *CustomResourceSubresourcesDie)) *CustomResourceDefinitionVersionDie {
	return d.DieStamp(func(r *apiextensionsv1.CustomResourceDefinitionVersion) {
		d := CustomResourceSubresourcesBlank.DieImmutable(false).DieFeedPtr(r.Subresources)
		fn(d)
		r.Subresources = d.DieReleasePtr()
	})
}

func (d *CustomResourceDefinitionVersionDie) AdditionalPrinterColumnDie(name string, fn func(d *CustomResourceColumnDefinitionDie)) *CustomResourceDefinitionVersionDie {
	return d.DieStamp(func(r *apiextensionsv1.CustomResourceDefinitionVersion) {
		for i := range r.AdditionalPrinterColumns {
			if name == r.AdditionalPrinterColumns[i].Name {
				d := CustomResourceColumnDefinitionBlank.DieImmutable(false).DieFeed(r.AdditionalPrinterColumns[i])
				fn(d)
				r.AdditionalPrinterColumns[i] = d.DieRelease()
				return
			}
		}

		d := CustomResourceColumnDefinitionBlank.DieImmutable(false).DieFeed(apiextensionsv1.CustomResourceColumnDefinition{Name: name})
		fn(d)
		r.AdditionalPrinterColumns = append(r.AdditionalPrinterColumns, d.DieRelease())
	})
}

func (d *CustomResourceDefinitionVersionDie) SelectableFieldsDie(fields ...*SelectableFieldDie) *CustomResourceDefinitionVersionDie {
	return d.DieStamp(func(r *apiextensionsv1.CustomResourceDefinitionVersion) {
		r.SelectableFields = make([]apiextensionsv1.SelectableField, len(fields))
		for i := range fields {
			r.SelectableFields[i] = fields[i].DieRelease()
		}
	})
}

// +die
type _ apiextensionsv1.CustomResourceValidation

// +die
type _ apiextensionsv1.CustomResourceSubresources

func (d *CustomResourceSubresourcesDie) ScaleDie(fn func(d *CustomResourceSubresourceScaleDie)) *CustomResourceSubresourcesDie {
	return d.DieStamp(func(r *apiextensionsv1.CustomResourceSubresources) {
		d := CustomResourceSubresourceScaleBlank.DieImmutable(false).DieFeedPtr(r.Scale)
		fn(d)
		r.Scale = d.DieReleasePtr()
	})
}

// +die
type _ apiextensionsv1.CustomResourceSubresourceScale

// +die
type _ apiextensionsv1.CustomResourceColumnDefinition

// +die
type _ apiextensionsv1.CustomResourceConversion

func (d *CustomResourceConversionDie) WebhookDie(fn func(d *WebhookConversionDie)) *CustomResourceConversionDie {
	return d.DieStamp(func(r *apiextensionsv1.CustomResourceConversion) {
		d := WebhookConversionBlank.DieImmutable(false).DieFeedPtr(r.Webhook)
		fn(d)
		r.Webhook = d.DieReleasePtr()
	})
}

// +die
type _ apiextensionsv1.WebhookConversion

func (d *WebhookConversionDie) ClientConfigDie(fn func(d *WebhookClientConfigDie)) *WebhookConversionDie {
	return d.DieStamp(func(r *apiextensionsv1.WebhookConversion) {
		d := WebhookClientConfigBlank.DieImmutable(false).DieFeedPtr(r.ClientConfig)
		fn(d)
		r.ClientConfig = d.DieReleasePtr()
	})
}

// +die
type _ apiextensionsv1.WebhookClientConfig

func (d *WebhookClientConfigDie) ServiceDie(fn func(d *ServiceReferenceDie)) *WebhookClientConfigDie {
	return d.DieStamp(func(r *apiextensionsv1.WebhookClientConfig) {
		d := ServiceReferenceBlank.DieImmutable(false).DieFeedPtr(r.Service)
		fn(d)
		r.Service = d.DieReleasePtr()
	})
}

// +die
type _ apiextensionsv1.ServiceReference

// +die
type _ apiextensionsv1.SelectableField

// +die
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

func (d *CustomResourceDefinitionStatusDie) AcceptedNamesDie(fn func(d *CustomResourceDefinitionNamesDie)) *CustomResourceDefinitionStatusDie {
	return d.DieStamp(func(r *apiextensionsv1.CustomResourceDefinitionStatus) {
		d := CustomResourceDefinitionNamesBlank.DieImmutable(false).DieFeed(r.AcceptedNames)
		fn(d)
		r.AcceptedNames = d.DieRelease()
	})
}

// +die
type _ = apiextensionsv1.CustomResourceDefinitionNames
