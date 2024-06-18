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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	diemetav1 "reconciler.io/dies/apis/meta/v1"
)

// +die:object=true,apiVersion=admissionregistration.k8s.io/v1,kind=ValidatingAdmissionPolicy
type _ = admissionregistrationv1.ValidatingAdmissionPolicy

// +die
type _ = admissionregistrationv1.ValidatingAdmissionPolicySpec

func (d *ValidatingAdmissionPolicySpecDie) ParamKindDie(fn func(d *ParamKindDie)) *ValidatingAdmissionPolicySpecDie {
	return d.DieStamp(func(r *admissionregistrationv1.ValidatingAdmissionPolicySpec) {
		d := ParamKindBlank.DieImmutable(false).DieFeedPtr(r.ParamKind)
		fn(d)
		r.ParamKind = d.DieReleasePtr()
	})
}

func (d *ValidatingAdmissionPolicySpecDie) MatchConstraintsDie(fn func(d *MatchResourcesDie)) *ValidatingAdmissionPolicySpecDie {
	return d.DieStamp(func(r *admissionregistrationv1.ValidatingAdmissionPolicySpec) {
		d := MatchResourcesBlank.DieImmutable(false).DieFeedPtr(r.MatchConstraints)
		fn(d)
		r.MatchConstraints = d.DieReleasePtr()
	})
}

func (d *ValidatingAdmissionPolicySpecDie) ValidationsDie(validations ...*ValidationDie) *ValidatingAdmissionPolicySpecDie {
	return d.DieStamp(func(r *admissionregistrationv1.ValidatingAdmissionPolicySpec) {
		r.Validations = make([]admissionregistrationv1.Validation, len(validations))
		for i := range validations {
			r.Validations[i] = validations[i].DieRelease()
		}
	})
}

func (d *ValidatingAdmissionPolicySpecDie) AuditAnnotationsDie(annotations ...*AuditAnnotationDie) *ValidatingAdmissionPolicySpecDie {
	return d.DieStamp(func(r *admissionregistrationv1.ValidatingAdmissionPolicySpec) {
		r.AuditAnnotations = make([]admissionregistrationv1.AuditAnnotation, len(annotations))
		for i := range annotations {
			r.AuditAnnotations[i] = annotations[i].DieRelease()
		}
	})
}

func (d *ValidatingAdmissionPolicySpecDie) MatchConditionDie(name string, fn func(d *MatchConditionDie)) *ValidatingAdmissionPolicySpecDie {
	return d.DieStamp(func(r *admissionregistrationv1.ValidatingAdmissionPolicySpec) {
		for i := range r.MatchConditions {
			if name == r.MatchConditions[i].Name {
				d := MatchConditionBlank.DieImmutable(false).DieFeed(r.MatchConditions[i])
				fn(d)
				r.MatchConditions[i] = d.DieRelease()
				return
			}
		}

		d := MatchConditionBlank.DieImmutable(false).DieFeed(admissionregistrationv1.MatchCondition{Name: name})
		fn(d)
		r.MatchConditions = append(r.MatchConditions, d.DieRelease())
	})
}

func (d *ValidatingAdmissionPolicySpecDie) VariablesDie(name string, fn func(d *VariableDie)) *ValidatingAdmissionPolicySpecDie {
	return d.DieStamp(func(r *admissionregistrationv1.ValidatingAdmissionPolicySpec) {
		for i := range r.Variables {
			if name == r.Variables[i].Name {
				d := VariableBlank.DieImmutable(false).DieFeed(r.Variables[i])
				fn(d)
				r.Variables[i] = d.DieRelease()
				return
			}
		}

		d := VariableBlank.DieImmutable(false).DieFeed(admissionregistrationv1.Variable{Name: name})
		fn(d)
		r.Variables = append(r.Variables, d.DieRelease())
	})
}

// +die
type _ = admissionregistrationv1.ParamKind

// +die
type _ = admissionregistrationv1.MatchResources

func (d *MatchResourcesDie) NamespaceSelectorDie(fn func(d *diemetav1.LabelSelectorDie)) *MatchResourcesDie {
	return d.DieStamp(func(r *admissionregistrationv1.MatchResources) {
		d := diemetav1.LabelSelectorBlank.DieImmutable(false).DieFeedPtr(r.NamespaceSelector)
		fn(d)
		r.NamespaceSelector = d.DieReleasePtr()
	})
}

func (d *MatchResourcesDie) ObjectSelectorDie(fn func(d *diemetav1.LabelSelectorDie)) *MatchResourcesDie {
	return d.DieStamp(func(r *admissionregistrationv1.MatchResources) {
		d := diemetav1.LabelSelectorBlank.DieImmutable(false).DieFeedPtr(r.ObjectSelector)
		fn(d)
		r.ObjectSelector = d.DieReleasePtr()
	})
}

func (d *MatchResourcesDie) ResourceRulesDie(rules ...*NamedRuleWithOperationsDie) *MatchResourcesDie {
	return d.DieStamp(func(r *admissionregistrationv1.MatchResources) {
		r.ResourceRules = make([]admissionregistrationv1.NamedRuleWithOperations, len(rules))
		for i := range rules {
			r.ResourceRules[i] = rules[i].DieRelease()
		}
	})
}

func (d *MatchResourcesDie) ExcludeResourceRulesDie(rules ...*NamedRuleWithOperationsDie) *MatchResourcesDie {
	return d.DieStamp(func(r *admissionregistrationv1.MatchResources) {
		r.ExcludeResourceRules = make([]admissionregistrationv1.NamedRuleWithOperations, len(rules))
		for i := range rules {
			r.ExcludeResourceRules[i] = rules[i].DieRelease()
		}
	})
}

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
type _ = admissionregistrationv1.ValidatingAdmissionPolicyStatus

func (d *ValidatingAdmissionPolicyStatusDie) TypeCheckingDie(fn func(d *TypeCheckingDie)) *ValidatingAdmissionPolicyStatusDie {
	return d.DieStamp(func(r *admissionregistrationv1.ValidatingAdmissionPolicyStatus) {
		d := TypeCheckingBlank.DieImmutable(false).DieFeedPtr(r.TypeChecking)
		fn(d)
		r.TypeChecking = d.DieReleasePtr()
	})
}

func (d *ValidatingAdmissionPolicyStatusDie) ConditionsDie(conditions ...*diemetav1.ConditionDie) *ValidatingAdmissionPolicyStatusDie {
	return d.DieStamp(func(r *admissionregistrationv1.ValidatingAdmissionPolicyStatus) {
		r.Conditions = make([]metav1.Condition, len(conditions))
		for i := range conditions {
			r.Conditions[i] = conditions[i].DieRelease()
		}
	})
}

// +die
type _ = admissionregistrationv1.TypeChecking

func (d *TypeCheckingDie) ExpressionWarningsDie(warnings ...*ExpressionWarningDie) *TypeCheckingDie {
	return d.DieStamp(func(r *admissionregistrationv1.TypeChecking) {
		r.ExpressionWarnings = make([]admissionregistrationv1.ExpressionWarning, len(warnings))
		for i := range warnings {
			r.ExpressionWarnings[i] = warnings[i].DieRelease()
		}
	})
}

// +die
type _ = admissionregistrationv1.ExpressionWarning
