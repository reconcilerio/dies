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
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
)

// +die
// +die:field:name=Service,die=ServiceReferenceDie,pointer=true
type _ = admissionregistrationv1.WebhookClientConfig

// +die
type _ admissionregistrationv1.ServiceReference

// +die
type _ admissionregistrationv1.RuleWithOperations

func (d *RuleWithOperationsDie) APIGroups(v ...string) *RuleWithOperationsDie {
	return d.DieStamp(func(r *admissionregistrationv1.RuleWithOperations) {
		r.APIGroups = v
	})
}

func (d *RuleWithOperationsDie) APIVersions(v ...string) *RuleWithOperationsDie {
	return d.DieStamp(func(r *admissionregistrationv1.RuleWithOperations) {
		r.APIVersions = v
	})
}

func (d *RuleWithOperationsDie) Resources(v ...string) *RuleWithOperationsDie {
	return d.DieStamp(func(r *admissionregistrationv1.RuleWithOperations) {
		r.Resources = v
	})
}

func (d *RuleWithOperationsDie) Scope(v *admissionregistrationv1.ScopeType) *RuleWithOperationsDie {
	return d.DieStamp(func(r *admissionregistrationv1.RuleWithOperations) {
		r.Scope = v
	})
}

// +die
type _ admissionregistrationv1.Rule

// +die
type _ = admissionregistrationv1.MatchCondition

// +die
type _ = admissionregistrationv1.ParamKind

// +die
// +die:field:name=NamespaceSelector,package=_/meta/v1,die=LabelSelectorDie,pointer=true
// +die:field:name=ObjectSelector,package=_/meta/v1,die=LabelSelectorDie,pointer=true
// +die:field:name=ResourceRules,die=NamedRuleWithOperationsDie,listType=atomic
// +die:field:name=ExcludeResourceRules,die=NamedRuleWithOperationsDie,listType=atomic
type _ = admissionregistrationv1.MatchResources

// +die
type _ = admissionregistrationv1.NamedRuleWithOperations

func (d *NamedRuleWithOperationsDie) Operations(v ...admissionregistrationv1.OperationType) *NamedRuleWithOperationsDie {
	return d.DieStamp(func(r *admissionregistrationv1.NamedRuleWithOperations) {
		r.Operations = v
	})
}

func (d *NamedRuleWithOperationsDie) Rule(v admissionregistrationv1.Rule) *NamedRuleWithOperationsDie {
	return d.DieStamp(func(r *admissionregistrationv1.NamedRuleWithOperations) {
		r.Rule = v
	})
}

func (d *NamedRuleWithOperationsDie) APIGroups(v ...string) *NamedRuleWithOperationsDie {
	return d.DieStamp(func(r *admissionregistrationv1.NamedRuleWithOperations) {
		r.APIGroups = v
	})
}

func (d *NamedRuleWithOperationsDie) APIVersions(v ...string) *NamedRuleWithOperationsDie {
	return d.DieStamp(func(r *admissionregistrationv1.NamedRuleWithOperations) {
		r.APIVersions = v
	})
}

func (d *NamedRuleWithOperationsDie) Resources(v ...string) *NamedRuleWithOperationsDie {
	return d.DieStamp(func(r *admissionregistrationv1.NamedRuleWithOperations) {
		r.Resources = v
	})
}

func (d *NamedRuleWithOperationsDie) Scope(v *admissionregistrationv1.ScopeType) *NamedRuleWithOperationsDie {
	return d.DieStamp(func(r *admissionregistrationv1.NamedRuleWithOperations) {
		r.Scope = v
	})
}

// +die
type _ = admissionregistrationv1.Variable

// +die
// +die:field:name=Selector,package=_/meta/v1,die=LabelSelectorDie,pointer=true
type _ = admissionregistrationv1.ParamRef
