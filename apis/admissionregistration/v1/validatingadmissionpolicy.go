/*
Copyright 2024 the original author or authors.

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

// +die:object=true,apiVersion=admissionregistration.k8s.io/v1,kind=ValidatingAdmissionPolicy
type _ = admissionregistrationv1.ValidatingAdmissionPolicy

// +die
// +die:field:name=ParamKind,die=ParamKindDie,pointer=true
// +die:field:name=MatchConstraints,die=MatchResourcesDie,pointer=true
// +die:field:name=Validations,die=ValidationDie,listType=atomic
// +die:field:name=AuditAnnotations,die=AuditAnnotationDie,listType=atomic
// +die:field:name=MatchConditions,die=MatchConditionDie,listType=map
// +die:field:name=Variables,die=VariableDie,listType=map
type _ = admissionregistrationv1.ValidatingAdmissionPolicySpec

// deprecated: use VariableDie instead
func (d *ValidatingAdmissionPolicySpecDie) VariablesDie(name string, fn func(d *VariableDie)) *ValidatingAdmissionPolicySpecDie {
	return d.VariableDie(name, fn)
}

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
type _ = admissionregistrationv1.Validation

// +die
type _ = admissionregistrationv1.AuditAnnotation

// +die
type _ = admissionregistrationv1.Variable

// +die
// +die:field:name=TypeChecking,die=TypeCheckingDie,pointer=true
// +die:field:name=Conditions,package=_/meta/v1,die=ConditionDie,listType=atomic
type _ = admissionregistrationv1.ValidatingAdmissionPolicyStatus

// +die
// +die:field:name=ExpressionWarnings,die=ExpressionWarningDie,listType=atomic
type _ = admissionregistrationv1.TypeChecking

// +die
type _ = admissionregistrationv1.ExpressionWarning
