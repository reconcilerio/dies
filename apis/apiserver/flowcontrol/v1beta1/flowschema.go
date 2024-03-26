/*
Copyright 2022 the original author or authors.

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

package v1beta1

import (
	flowcontrolv1beta1 "k8s.io/api/flowcontrol/v1beta1"
	diemetav1 "reconciler.io/dies/apis/meta/v1"
)

// +die:object=true
type _ = flowcontrolv1beta1.FlowSchema

// +die
type _ = flowcontrolv1beta1.FlowSchemaSpec

func (d *FlowSchemaSpecDie) PriorityLevelConfigurationDie(fn func(d *PriorityLevelConfigurationReferenceDie)) *FlowSchemaSpecDie {
	return d.DieStamp(func(r *flowcontrolv1beta1.FlowSchemaSpec) {
		d := PriorityLevelConfigurationReferenceBlank.DieImmutable(false).DieFeed(r.PriorityLevelConfiguration)
		fn(d)
		r.PriorityLevelConfiguration = d.DieRelease()
	})
}

func (d *FlowSchemaSpecDie) DistinguisherMethodDie(fn func(d *FlowDistinguisherMethodDie)) *FlowSchemaSpecDie {
	return d.DieStamp(func(r *flowcontrolv1beta1.FlowSchemaSpec) {
		d := FlowDistinguisherMethodBlank.DieImmutable(false).DieFeedPtr(r.DistinguisherMethod)
		fn(d)
		r.DistinguisherMethod = d.DieReleasePtr()
	})
}

func (d *FlowSchemaSpecDie) RulesDie(rules ...*PolicyRulesWithSubjectsDie) *FlowSchemaSpecDie {
	return d.DieStamp(func(r *flowcontrolv1beta1.FlowSchemaSpec) {
		r.Rules = make([]flowcontrolv1beta1.PolicyRulesWithSubjects, len(rules))
		for i := range rules {
			r.Rules[i] = rules[i].DieRelease()
		}
	})
}

// +die
type _ = flowcontrolv1beta1.FlowSchemaStatus

func (d *FlowSchemaStatusDie) ConditionsDie(conditions ...*diemetav1.ConditionDie) *FlowSchemaStatusDie {
	return d.DieStamp(func(r *flowcontrolv1beta1.FlowSchemaStatus) {
		r.Conditions = make([]flowcontrolv1beta1.FlowSchemaCondition, len(conditions))
		for i := range conditions {
			c := conditions[i].DieRelease()
			r.Conditions[i] = flowcontrolv1beta1.FlowSchemaCondition{
				Type:               flowcontrolv1beta1.FlowSchemaConditionType(c.Type),
				Status:             flowcontrolv1beta1.ConditionStatus(c.Status),
				Reason:             c.Reason,
				Message:            c.Message,
				LastTransitionTime: c.LastTransitionTime,
			}
		}
	})
}

// +die
type _ = flowcontrolv1beta1.PriorityLevelConfigurationReference

// +die
type _ = flowcontrolv1beta1.FlowDistinguisherMethod

// +die
type _ = flowcontrolv1beta1.PolicyRulesWithSubjects

func (d *PolicyRulesWithSubjectsDie) SubjectsDie(subjects ...*SubjectDie) *PolicyRulesWithSubjectsDie {
	return d.DieStamp(func(r *flowcontrolv1beta1.PolicyRulesWithSubjects) {
		r.Subjects = make([]flowcontrolv1beta1.Subject, len(subjects))
		for i := range subjects {
			r.Subjects[i] = subjects[i].DieRelease()
		}
	})
}

func (d *PolicyRulesWithSubjectsDie) ResourceRulesDie(rules ...*ResourcePolicyRuleDie) *PolicyRulesWithSubjectsDie {
	return d.DieStamp(func(r *flowcontrolv1beta1.PolicyRulesWithSubjects) {
		r.ResourceRules = make([]flowcontrolv1beta1.ResourcePolicyRule, len(rules))
		for i := range rules {
			r.ResourceRules[i] = rules[i].DieRelease()
		}
	})
}

func (d *PolicyRulesWithSubjectsDie) NonResourcePolicyRuleDie(rules ...*NonResourcePolicyRuleDie) *PolicyRulesWithSubjectsDie {
	return d.DieStamp(func(r *flowcontrolv1beta1.PolicyRulesWithSubjects) {
		r.NonResourceRules = make([]flowcontrolv1beta1.NonResourcePolicyRule, len(rules))
		for i := range rules {
			r.NonResourceRules[i] = rules[i].DieRelease()
		}
	})
}

// +die
type _ = flowcontrolv1beta1.Subject

func (d *SubjectDie) UserDie(fn func(d *UserSubjectDie)) *SubjectDie {
	return d.DieStamp(func(r *flowcontrolv1beta1.Subject) {
		d := UserSubjectBlank.DieImmutable(false).DieFeedPtr(r.User)
		fn(d)
		r.User = d.DieReleasePtr()
		r.Kind = flowcontrolv1beta1.SubjectKindUser
	})
}

func (d *SubjectDie) GroupDie(fn func(d *GroupSubjectDie)) *SubjectDie {
	return d.DieStamp(func(r *flowcontrolv1beta1.Subject) {
		d := GroupSubjectBlank.DieImmutable(false).DieFeedPtr(r.Group)
		fn(d)
		r.Group = d.DieReleasePtr()
		r.Kind = flowcontrolv1beta1.SubjectKindGroup
	})
}

func (d *SubjectDie) ServiceAccountDie(fn func(d *ServiceAccountSubjectDie)) *SubjectDie {
	return d.DieStamp(func(r *flowcontrolv1beta1.Subject) {
		d := ServiceAccountSubjectBlank.DieImmutable(false).DieFeedPtr(r.ServiceAccount)
		fn(d)
		r.ServiceAccount = d.DieReleasePtr()
		r.Kind = flowcontrolv1beta1.SubjectKindServiceAccount
	})
}

// +die
type _ = flowcontrolv1beta1.UserSubject

// +die
type _ = flowcontrolv1beta1.GroupSubject

// +die
type _ = flowcontrolv1beta1.ServiceAccountSubject

// +die
type _ = flowcontrolv1beta1.ResourcePolicyRule

// +die
type _ = flowcontrolv1beta1.NonResourcePolicyRule
