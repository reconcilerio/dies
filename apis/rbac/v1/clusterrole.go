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
	diemetav1 "dies.dev/apis/meta/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +die:object=true
type ClusterRole = rbacv1.ClusterRole

func (d *ClusterRoleDie) RulesDie(rules ...*PolicyRuleDie) *ClusterRoleDie {
	return d.DieStamp(func(r *rbacv1.ClusterRole) {
		r.Rules = make([]rbacv1.PolicyRule, len(rules))
		for i := range rules {
			r.Rules[i] = rules[i].DieRelease()
		}
	})
}

func (d *ClusterRoleDie) AddRuleDie(rule *PolicyRuleDie) *ClusterRoleDie {
	return d.DieStamp(func(r *rbacv1.ClusterRole) {
		r.Rules = append(r.Rules, rule.DieRelease())
	})
}

func (d *ClusterRoleDie) AggregationRuleDie(fn func(d *AggregationRuleDie)) *ClusterRoleDie {
	return d.DieStamp(func(r *rbacv1.ClusterRole) {
		d := AggregationRuleBlank.DieImmutable(false).DieFeedPtr(r.AggregationRule)
		fn(d)
		r.AggregationRule = d.DieReleasePtr()
	})
}

// +die
type AggregationRule = rbacv1.AggregationRule

func (d *AggregationRuleDie) ClusterRoleSelectorsDie(selectors ...*diemetav1.LabelSelectorDie) *AggregationRuleDie {
	return d.DieStamp(func(r *rbacv1.AggregationRule) {
		r.ClusterRoleSelectors = make([]metav1.LabelSelector, len(selectors))
		for i := range selectors {
			r.ClusterRoleSelectors[i] = selectors[i].DieRelease()
		}
	})
}
