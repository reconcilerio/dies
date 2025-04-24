/*
Copyright 2025 the original author or authors.

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
	flowcontrolv1 "k8s.io/api/flowcontrol/v1"
	diemetav1 "reconciler.io/dies/apis/meta/v1"
)

// +die:object=true,apiVersion=flowcontrol.apiserver.k8s.io/v1,kind=FlowSchema
type _ = flowcontrolv1.FlowSchema

// +die
// +die:field:name=PriorityLevelConfiguration,die=PriorityLevelConfigurationReferenceDie
// +die:field:name=DistinguisherMethod,die=FlowDistinguisherMethodDie,pointer=true
// +die:field:name=Rules,die=PolicyRulesWithSubjectsDie,listType=atomic
type _ = flowcontrolv1.FlowSchemaSpec

// +die
type _ = flowcontrolv1.FlowSchemaStatus

func (d *FlowSchemaStatusDie) ConditionsDie(conditions ...*diemetav1.ConditionDie) *FlowSchemaStatusDie {
	return d.DieStamp(func(r *flowcontrolv1.FlowSchemaStatus) {
		r.Conditions = make([]flowcontrolv1.FlowSchemaCondition, len(conditions))
		for i := range conditions {
			c := conditions[i].DieRelease()
			r.Conditions[i] = flowcontrolv1.FlowSchemaCondition{
				Type:               flowcontrolv1.FlowSchemaConditionType(c.Type),
				Status:             flowcontrolv1.ConditionStatus(c.Status),
				Reason:             c.Reason,
				Message:            c.Message,
				LastTransitionTime: c.LastTransitionTime,
			}
		}
	})
}

// +die
type _ = flowcontrolv1.PriorityLevelConfigurationReference

// +die
type _ = flowcontrolv1.FlowDistinguisherMethod

// +die
// +die:field:name=Subjects,die=SubjectDie,listType=atomic
// +die:field:name=ResourceRules,die=ResourcePolicyRuleDie,listType=atomic
// +die:field:name=NonResourceRules,die=NonResourcePolicyRuleDie,listType=atomic
type _ = flowcontrolv1.PolicyRulesWithSubjects

// deprecated: use NonResourceRulesDie
func (d *PolicyRulesWithSubjectsDie) NonResourcePolicyRuleDie(rules ...*NonResourcePolicyRuleDie) *PolicyRulesWithSubjectsDie {
	return d.NonResourceRulesDie(rules...)
}

// +die
type _ = flowcontrolv1.Subject

func (d *SubjectDie) UserDie(fn func(d *UserSubjectDie)) *SubjectDie {
	return d.DieStamp(func(r *flowcontrolv1.Subject) {
		d := UserSubjectBlank.DieImmutable(false).DieFeedPtr(r.User)
		fn(d)
		r.User = d.DieReleasePtr()
		r.Kind = flowcontrolv1.SubjectKindUser
	})
}

func (d *SubjectDie) GroupDie(fn func(d *GroupSubjectDie)) *SubjectDie {
	return d.DieStamp(func(r *flowcontrolv1.Subject) {
		d := GroupSubjectBlank.DieImmutable(false).DieFeedPtr(r.Group)
		fn(d)
		r.Group = d.DieReleasePtr()
		r.Kind = flowcontrolv1.SubjectKindGroup
	})
}

func (d *SubjectDie) ServiceAccountDie(fn func(d *ServiceAccountSubjectDie)) *SubjectDie {
	return d.DieStamp(func(r *flowcontrolv1.Subject) {
		d := ServiceAccountSubjectBlank.DieImmutable(false).DieFeedPtr(r.ServiceAccount)
		fn(d)
		r.ServiceAccount = d.DieReleasePtr()
		r.Kind = flowcontrolv1.SubjectKindServiceAccount
	})
}

// +die
type _ = flowcontrolv1.UserSubject

// +die
type _ = flowcontrolv1.GroupSubject

// +die
type _ = flowcontrolv1.ServiceAccountSubject

// +die
type _ = flowcontrolv1.ResourcePolicyRule

// +die
type _ = flowcontrolv1.NonResourcePolicyRule
