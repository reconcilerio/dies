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

package v1alpha1

import (
	admissionregistrationv1alpha1 "k8s.io/api/admissionregistration/v1alpha1"
)

// +die:object=true,apiVersion=admissionregistration.k8s.io/v1alpha1,kind=MutatingAdmissionPolicy
type _ = admissionregistrationv1alpha1.MutatingAdmissionPolicy

// +die
// +die:field:name=ParamKind,die=ParamKindDie,pointer=true
// +die:field:name=MatchConstraints,die=MatchResourcesDie,pointer=true
// +die:field:name=Variables,die=VariableDie,listType=atomic
// +die:field:name=Mutations,die=MutationDie,listType=atomic
// +die:field:name=MatchConditions,die=MatchConditionDie,listType=map
type _ = admissionregistrationv1alpha1.MutatingAdmissionPolicySpec

// +die
type _ = admissionregistrationv1alpha1.ParamKind

// +die
// +die:field:name=NamespaceSelector,die=LabelSelectorDie,package=_/meta/v1,pointer=true
// +die:field:name=ObjectSelector,die=LabelSelectorDie,package=_/meta/v1,pointer=true
// +die:field:name=ResourceRules,die=NamedRuleWithOperationsDie,listType=atomic
// +die:field:name=ExcludeResourceRules,die=NamedRuleWithOperationsDie,listType=atomic
type _ = admissionregistrationv1alpha1.MatchResources

// +die
type _ = admissionregistrationv1alpha1.MatchCondition

// Name is an identifier for this match condition, used for strategic merging of MatchConditions,
// as well as providing an identifier for logging purposes. A good name should be descriptive of
// the associated expression.
// Name must be a qualified name consisting of alphanumeric characters, '-', '_' or '.', and
// must start and end with an alphanumeric character (e.g. 'MyName',  or 'my.name',  or
// '123-abc', regex used for validation is '([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]') with an
// optional DNS subdomain prefix and '/' (e.g. 'example.com/MyName')
//
// Required.
func (d *MatchConditionDie) Name(v string) *MatchConditionDie {
	return d.DieStamp(func(r *admissionregistrationv1alpha1.MatchCondition) {
		r.Name = v
	})
}

// Expression represents the expression which will be evaluated by CEL. Must evaluate to bool.
// CEL expressions have access to the contents of the AdmissionRequest and Authorizer, organized into CEL variables:
//
// 'object' - The object from the incoming request. The value is null for DELETE requests.
// 'oldObject' - The existing object. The value is null for CREATE requests.
// 'request' - Attributes of the admission request(/pkg/apis/admission/types.go#AdmissionRequest).
// 'authorizer' - A CEL Authorizer. May be used to perform authorization checks for the principal (user or service account) of the request.
//
//	See https://pkg.go.dev/k8s.io/apiserver/pkg/cel/library#Authz
//
// 'authorizer.requestResource' - A CEL ResourceCheck constructed from the 'authorizer' and configured with the
//
//	request resource.
//
// Documentation on CEL: https://kubernetes.io/docs/reference/using-api/cel/
//
// Required.
func (d *MatchConditionDie) Expression(v string) *MatchConditionDie {
	return d.DieStamp(func(r *admissionregistrationv1alpha1.MatchCondition) {
		r.Expression = v
	})
}

// +die
type _ = admissionregistrationv1alpha1.Variable

// +die
// +die:field:name=ApplyConfiguration,die=ApplyConfigurationDie,pointer=true
// +die:field:name=JSONPatch,die=JSONPatchDie,pointer=true
type _ = admissionregistrationv1alpha1.Mutation

// +die
// +die:field:name=RuleWithOperations,die=RuleWithOperationsDie,package=_/admissionregistration/v1
type _ = admissionregistrationv1alpha1.NamedRuleWithOperations

// +die
type _ = admissionregistrationv1alpha1.ApplyConfiguration

// +die
type _ = admissionregistrationv1alpha1.JSONPatch
