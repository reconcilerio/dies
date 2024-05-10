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
	rbacv1 "k8s.io/api/rbac/v1"
)

// +die:object=true,apiVersion=rbac.authorization.k8s.io/v1,kind=Role
type _ = rbacv1.Role

func (d *RoleDie) RulesDie(rules ...*PolicyRuleDie) *RoleDie {
	return d.DieStamp(func(r *rbacv1.Role) {
		r.Rules = make([]rbacv1.PolicyRule, len(rules))
		for i := range rules {
			r.Rules[i] = rules[i].DieRelease()
		}
	})
}

func (d *RoleDie) AddRuleDie(rule *PolicyRuleDie) *RoleDie {
	return d.DieStamp(func(r *rbacv1.Role) {
		r.Rules = append(r.Rules, rule.DieRelease())
	})
}

// +die
type _ = rbacv1.PolicyRule

func (d *PolicyRuleDie) AddVerbs(verbs ...string) *PolicyRuleDie {
	return d.DieStamp(func(r *rbacv1.PolicyRule) {
		r.Verbs = append(r.Verbs, verbs...)
	})
}

func (d *PolicyRuleDie) AddAPIGroups(apiGroups ...string) *PolicyRuleDie {
	return d.DieStamp(func(r *rbacv1.PolicyRule) {
		r.APIGroups = append(r.APIGroups, apiGroups...)
	})
}

func (d *PolicyRuleDie) AddAResources(resources ...string) *PolicyRuleDie {
	return d.DieStamp(func(r *rbacv1.PolicyRule) {
		r.Resources = append(r.Resources, resources...)
	})
}

func (d *PolicyRuleDie) AddResourceNames(resourceNames ...string) *PolicyRuleDie {
	return d.DieStamp(func(r *rbacv1.PolicyRule) {
		r.ResourceNames = append(r.ResourceNames, resourceNames...)
	})
}

func (d *PolicyRuleDie) AddNonResourceURLs(nonResourceURLs ...string) *PolicyRuleDie {
	return d.DieStamp(func(r *rbacv1.PolicyRule) {
		r.NonResourceURLs = append(r.NonResourceURLs, nonResourceURLs...)
	})
}
