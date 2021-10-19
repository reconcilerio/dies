//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by diegen. DO NOT EDIT.

package v1

import (
	metav1 "dies.dev/apis/meta/v1"
	json "encoding/json"
	fmtx "fmt"
	rbacv1 "k8s.io/api/rbac/v1"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
)

type ClusterRoleDie interface {
	// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
	DieStamp(fn func(r *rbacv1.ClusterRole)) ClusterRoleDie
	// DieFeed returns a new die with the provided resource.
	DieFeed(r rbacv1.ClusterRole) ClusterRoleDie
	// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
	DieFeedPtr(r *rbacv1.ClusterRole) ClusterRoleDie
	// DieRelease returns the resource managed by the die.
	DieRelease() rbacv1.ClusterRole
	// DieReleasePtr returns a pointer to the resource managed by the die.
	DieReleasePtr() *rbacv1.ClusterRole
	// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
	DieImmutable(immutable bool) ClusterRoleDie
	// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
	DeepCopy() ClusterRoleDie

	clusterRole
	// MetadataDie stamps the resource's ObjectMeta field with a mutable die.
	MetadataDie(fn func(d metav1.ObjectMetaDie)) ClusterRoleDie
	// Rules holds all the PolicyRules for this ClusterRole
	Rules(Rules ...rbacv1.PolicyRule) ClusterRoleDie
	// AggregationRule is an optional field that describes how to build the Rules for this ClusterRole. If AggregationRule is set, then the Rules are controller managed and direct changes to Rules will be stomped by the controller.
	AggregationRule(AggregationRule *rbacv1.AggregationRule) ClusterRoleDie

	runtime.Object
	apismetav1.Object
	apismetav1.ObjectMetaAccessor
}

var _ ClusterRoleDie = (*clusterRoleDie)(nil)
var ClusterRoleBlank = (&clusterRoleDie{}).DieFeed(rbacv1.ClusterRole{})

type clusterRoleDie struct {
	metav1.FrozenObjectMeta
	mutable bool
	r       rbacv1.ClusterRole
}

func (d *clusterRoleDie) DieImmutable(immutable bool) ClusterRoleDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy().(*clusterRoleDie)
	d.mutable = !immutable
	return d
}

func (d *clusterRoleDie) DieFeed(r rbacv1.ClusterRole) ClusterRoleDie {
	if d.mutable {
		d.FrozenObjectMeta = metav1.FreezeObjectMeta(r.ObjectMeta)
		d.r = r
		return d
	}
	return &clusterRoleDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

func (d *clusterRoleDie) DieFeedPtr(r *rbacv1.ClusterRole) ClusterRoleDie {
	if r == nil {
		r = &rbacv1.ClusterRole{}
	}
	return d.DieFeed(*r)
}

func (d *clusterRoleDie) DieRelease() rbacv1.ClusterRole {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

func (d *clusterRoleDie) DieReleasePtr() *rbacv1.ClusterRole {
	r := d.DieRelease()
	return &r
}

func (d *clusterRoleDie) DieStamp(fn func(r *rbacv1.ClusterRole)) ClusterRoleDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

func (d *clusterRoleDie) DeepCopy() ClusterRoleDie {
	r := *d.r.DeepCopy()
	return &clusterRoleDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

func (d *clusterRoleDie) DeepCopyObject() runtime.Object {
	return d.r.DeepCopy()
}

func (d *clusterRoleDie) GetObjectKind() schema.ObjectKind {
	r := d.DieRelease()
	return r.GetObjectKind()
}

func (d *clusterRoleDie) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.r)
}

func (d *clusterRoleDie) UnmarshalJSON(b []byte) error {
	if d == ClusterRoleBlank {
		return fmtx.Errorf("cannot unmarshal into the root object, create a copy first")
	}
	r := &rbacv1.ClusterRole{}
	err := json.Unmarshal(b, r)
	*d = *d.DieFeed(*r).(*clusterRoleDie)
	return err
}

func (d *clusterRoleDie) MetadataDie(fn func(d metav1.ObjectMetaDie)) ClusterRoleDie {
	return d.DieStamp(func(r *rbacv1.ClusterRole) {
		d := metav1.ObjectMetaBlank.DieImmutable(false).DieFeed(r.ObjectMeta)
		fn(d)
		r.ObjectMeta = d.DieRelease()
	})
}

func (d *clusterRoleDie) Rules(v ...rbacv1.PolicyRule) ClusterRoleDie {
	return d.DieStamp(func(r *rbacv1.ClusterRole) {
		r.Rules = v
	})
}

func (d *clusterRoleDie) AggregationRule(v *rbacv1.AggregationRule) ClusterRoleDie {
	return d.DieStamp(func(r *rbacv1.ClusterRole) {
		r.AggregationRule = v
	})
}

type AggregationRuleDie interface {
	// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
	DieStamp(fn func(r *rbacv1.AggregationRule)) AggregationRuleDie
	// DieFeed returns a new die with the provided resource.
	DieFeed(r rbacv1.AggregationRule) AggregationRuleDie
	// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
	DieFeedPtr(r *rbacv1.AggregationRule) AggregationRuleDie
	// DieRelease returns the resource managed by the die.
	DieRelease() rbacv1.AggregationRule
	// DieReleasePtr returns a pointer to the resource managed by the die.
	DieReleasePtr() *rbacv1.AggregationRule
	// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
	DieImmutable(immutable bool) AggregationRuleDie
	// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
	DeepCopy() AggregationRuleDie

	aggregationRule
	// ClusterRoleSelectors holds a list of selectors which will be used to find ClusterRoles and create the rules. If any of the selectors match, then the ClusterRole's permissions will be added
	ClusterRoleSelectors(ClusterRoleSelectors ...apismetav1.LabelSelector) AggregationRuleDie
}

var _ AggregationRuleDie = (*aggregationRuleDie)(nil)
var AggregationRuleBlank = (&aggregationRuleDie{}).DieFeed(rbacv1.AggregationRule{})

type aggregationRuleDie struct {
	mutable bool
	r       rbacv1.AggregationRule
}

func (d *aggregationRuleDie) DieImmutable(immutable bool) AggregationRuleDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy().(*aggregationRuleDie)
	d.mutable = !immutable
	return d
}

func (d *aggregationRuleDie) DieFeed(r rbacv1.AggregationRule) AggregationRuleDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &aggregationRuleDie{
		mutable: d.mutable,
		r:       r,
	}
}

func (d *aggregationRuleDie) DieFeedPtr(r *rbacv1.AggregationRule) AggregationRuleDie {
	if r == nil {
		r = &rbacv1.AggregationRule{}
	}
	return d.DieFeed(*r)
}

func (d *aggregationRuleDie) DieRelease() rbacv1.AggregationRule {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

func (d *aggregationRuleDie) DieReleasePtr() *rbacv1.AggregationRule {
	r := d.DieRelease()
	return &r
}

func (d *aggregationRuleDie) DieStamp(fn func(r *rbacv1.AggregationRule)) AggregationRuleDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

func (d *aggregationRuleDie) DeepCopy() AggregationRuleDie {
	r := *d.r.DeepCopy()
	return &aggregationRuleDie{
		mutable: d.mutable,
		r:       r,
	}
}

func (d *aggregationRuleDie) ClusterRoleSelectors(v ...apismetav1.LabelSelector) AggregationRuleDie {
	return d.DieStamp(func(r *rbacv1.AggregationRule) {
		r.ClusterRoleSelectors = v
	})
}

type ClusterRoleBindingDie interface {
	// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
	DieStamp(fn func(r *rbacv1.ClusterRoleBinding)) ClusterRoleBindingDie
	// DieFeed returns a new die with the provided resource.
	DieFeed(r rbacv1.ClusterRoleBinding) ClusterRoleBindingDie
	// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
	DieFeedPtr(r *rbacv1.ClusterRoleBinding) ClusterRoleBindingDie
	// DieRelease returns the resource managed by the die.
	DieRelease() rbacv1.ClusterRoleBinding
	// DieReleasePtr returns a pointer to the resource managed by the die.
	DieReleasePtr() *rbacv1.ClusterRoleBinding
	// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
	DieImmutable(immutable bool) ClusterRoleBindingDie
	// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
	DeepCopy() ClusterRoleBindingDie

	clusterRoleBinding
	// MetadataDie stamps the resource's ObjectMeta field with a mutable die.
	MetadataDie(fn func(d metav1.ObjectMetaDie)) ClusterRoleBindingDie
	// Subjects holds references to the objects the role applies to.
	Subjects(Subjects ...rbacv1.Subject) ClusterRoleBindingDie
	// RoleRef can only reference a ClusterRole in the global namespace. If the RoleRef cannot be resolved, the Authorizer must return an error.
	RoleRef(RoleRef rbacv1.RoleRef) ClusterRoleBindingDie

	runtime.Object
	apismetav1.Object
	apismetav1.ObjectMetaAccessor
}

var _ ClusterRoleBindingDie = (*clusterRoleBindingDie)(nil)
var ClusterRoleBindingBlank = (&clusterRoleBindingDie{}).DieFeed(rbacv1.ClusterRoleBinding{})

type clusterRoleBindingDie struct {
	metav1.FrozenObjectMeta
	mutable bool
	r       rbacv1.ClusterRoleBinding
}

func (d *clusterRoleBindingDie) DieImmutable(immutable bool) ClusterRoleBindingDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy().(*clusterRoleBindingDie)
	d.mutable = !immutable
	return d
}

func (d *clusterRoleBindingDie) DieFeed(r rbacv1.ClusterRoleBinding) ClusterRoleBindingDie {
	if d.mutable {
		d.FrozenObjectMeta = metav1.FreezeObjectMeta(r.ObjectMeta)
		d.r = r
		return d
	}
	return &clusterRoleBindingDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

func (d *clusterRoleBindingDie) DieFeedPtr(r *rbacv1.ClusterRoleBinding) ClusterRoleBindingDie {
	if r == nil {
		r = &rbacv1.ClusterRoleBinding{}
	}
	return d.DieFeed(*r)
}

func (d *clusterRoleBindingDie) DieRelease() rbacv1.ClusterRoleBinding {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

func (d *clusterRoleBindingDie) DieReleasePtr() *rbacv1.ClusterRoleBinding {
	r := d.DieRelease()
	return &r
}

func (d *clusterRoleBindingDie) DieStamp(fn func(r *rbacv1.ClusterRoleBinding)) ClusterRoleBindingDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

func (d *clusterRoleBindingDie) DeepCopy() ClusterRoleBindingDie {
	r := *d.r.DeepCopy()
	return &clusterRoleBindingDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

func (d *clusterRoleBindingDie) DeepCopyObject() runtime.Object {
	return d.r.DeepCopy()
}

func (d *clusterRoleBindingDie) GetObjectKind() schema.ObjectKind {
	r := d.DieRelease()
	return r.GetObjectKind()
}

func (d *clusterRoleBindingDie) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.r)
}

func (d *clusterRoleBindingDie) UnmarshalJSON(b []byte) error {
	if d == ClusterRoleBindingBlank {
		return fmtx.Errorf("cannot unmarshal into the root object, create a copy first")
	}
	r := &rbacv1.ClusterRoleBinding{}
	err := json.Unmarshal(b, r)
	*d = *d.DieFeed(*r).(*clusterRoleBindingDie)
	return err
}

func (d *clusterRoleBindingDie) MetadataDie(fn func(d metav1.ObjectMetaDie)) ClusterRoleBindingDie {
	return d.DieStamp(func(r *rbacv1.ClusterRoleBinding) {
		d := metav1.ObjectMetaBlank.DieImmutable(false).DieFeed(r.ObjectMeta)
		fn(d)
		r.ObjectMeta = d.DieRelease()
	})
}

func (d *clusterRoleBindingDie) Subjects(v ...rbacv1.Subject) ClusterRoleBindingDie {
	return d.DieStamp(func(r *rbacv1.ClusterRoleBinding) {
		r.Subjects = v
	})
}

func (d *clusterRoleBindingDie) RoleRef(v rbacv1.RoleRef) ClusterRoleBindingDie {
	return d.DieStamp(func(r *rbacv1.ClusterRoleBinding) {
		r.RoleRef = v
	})
}

type RoleDie interface {
	// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
	DieStamp(fn func(r *rbacv1.Role)) RoleDie
	// DieFeed returns a new die with the provided resource.
	DieFeed(r rbacv1.Role) RoleDie
	// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
	DieFeedPtr(r *rbacv1.Role) RoleDie
	// DieRelease returns the resource managed by the die.
	DieRelease() rbacv1.Role
	// DieReleasePtr returns a pointer to the resource managed by the die.
	DieReleasePtr() *rbacv1.Role
	// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
	DieImmutable(immutable bool) RoleDie
	// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
	DeepCopy() RoleDie

	role
	// MetadataDie stamps the resource's ObjectMeta field with a mutable die.
	MetadataDie(fn func(d metav1.ObjectMetaDie)) RoleDie
	// Rules holds all the PolicyRules for this Role
	Rules(Rules ...rbacv1.PolicyRule) RoleDie

	runtime.Object
	apismetav1.Object
	apismetav1.ObjectMetaAccessor
}

var _ RoleDie = (*roleDie)(nil)
var RoleBlank = (&roleDie{}).DieFeed(rbacv1.Role{})

type roleDie struct {
	metav1.FrozenObjectMeta
	mutable bool
	r       rbacv1.Role
}

func (d *roleDie) DieImmutable(immutable bool) RoleDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy().(*roleDie)
	d.mutable = !immutable
	return d
}

func (d *roleDie) DieFeed(r rbacv1.Role) RoleDie {
	if d.mutable {
		d.FrozenObjectMeta = metav1.FreezeObjectMeta(r.ObjectMeta)
		d.r = r
		return d
	}
	return &roleDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

func (d *roleDie) DieFeedPtr(r *rbacv1.Role) RoleDie {
	if r == nil {
		r = &rbacv1.Role{}
	}
	return d.DieFeed(*r)
}

func (d *roleDie) DieRelease() rbacv1.Role {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

func (d *roleDie) DieReleasePtr() *rbacv1.Role {
	r := d.DieRelease()
	return &r
}

func (d *roleDie) DieStamp(fn func(r *rbacv1.Role)) RoleDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

func (d *roleDie) DeepCopy() RoleDie {
	r := *d.r.DeepCopy()
	return &roleDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

func (d *roleDie) DeepCopyObject() runtime.Object {
	return d.r.DeepCopy()
}

func (d *roleDie) GetObjectKind() schema.ObjectKind {
	r := d.DieRelease()
	return r.GetObjectKind()
}

func (d *roleDie) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.r)
}

func (d *roleDie) UnmarshalJSON(b []byte) error {
	if d == RoleBlank {
		return fmtx.Errorf("cannot unmarshal into the root object, create a copy first")
	}
	r := &rbacv1.Role{}
	err := json.Unmarshal(b, r)
	*d = *d.DieFeed(*r).(*roleDie)
	return err
}

func (d *roleDie) MetadataDie(fn func(d metav1.ObjectMetaDie)) RoleDie {
	return d.DieStamp(func(r *rbacv1.Role) {
		d := metav1.ObjectMetaBlank.DieImmutable(false).DieFeed(r.ObjectMeta)
		fn(d)
		r.ObjectMeta = d.DieRelease()
	})
}

func (d *roleDie) Rules(v ...rbacv1.PolicyRule) RoleDie {
	return d.DieStamp(func(r *rbacv1.Role) {
		r.Rules = v
	})
}

type PolicyRuleDie interface {
	// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
	DieStamp(fn func(r *rbacv1.PolicyRule)) PolicyRuleDie
	// DieFeed returns a new die with the provided resource.
	DieFeed(r rbacv1.PolicyRule) PolicyRuleDie
	// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
	DieFeedPtr(r *rbacv1.PolicyRule) PolicyRuleDie
	// DieRelease returns the resource managed by the die.
	DieRelease() rbacv1.PolicyRule
	// DieReleasePtr returns a pointer to the resource managed by the die.
	DieReleasePtr() *rbacv1.PolicyRule
	// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
	DieImmutable(immutable bool) PolicyRuleDie
	// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
	DeepCopy() PolicyRuleDie

	policyRule
	// Verbs is a list of Verbs that apply to ALL the ResourceKinds and AttributeRestrictions contained in this rule. '*' represents all verbs.
	Verbs(Verbs ...string) PolicyRuleDie
	// APIGroups is the name of the APIGroup that contains the resources.  If multiple API groups are specified, any action requested against one of the enumerated resources in any API group will be allowed.
	APIGroups(APIGroups ...string) PolicyRuleDie
	// Resources is a list of resources this rule applies to. '*' represents all resources.
	Resources(Resources ...string) PolicyRuleDie
	// ResourceNames is an optional white list of names that the rule applies to.  An empty set means that everything is allowed.
	ResourceNames(ResourceNames ...string) PolicyRuleDie
	// NonResourceURLs is a set of partial urls that a user should have access to.  *s are allowed, but only as the full, final step in the path Since non-resource URLs are not namespaced, this field is only applicable for ClusterRoles referenced from a ClusterRoleBinding. Rules can either apply to API resources (such as "pods" or "secrets") or non-resource URL paths (such as "/api"),  but not both.
	NonResourceURLs(NonResourceURLs ...string) PolicyRuleDie
}

var _ PolicyRuleDie = (*policyRuleDie)(nil)
var PolicyRuleBlank = (&policyRuleDie{}).DieFeed(rbacv1.PolicyRule{})

type policyRuleDie struct {
	mutable bool
	r       rbacv1.PolicyRule
}

func (d *policyRuleDie) DieImmutable(immutable bool) PolicyRuleDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy().(*policyRuleDie)
	d.mutable = !immutable
	return d
}

func (d *policyRuleDie) DieFeed(r rbacv1.PolicyRule) PolicyRuleDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &policyRuleDie{
		mutable: d.mutable,
		r:       r,
	}
}

func (d *policyRuleDie) DieFeedPtr(r *rbacv1.PolicyRule) PolicyRuleDie {
	if r == nil {
		r = &rbacv1.PolicyRule{}
	}
	return d.DieFeed(*r)
}

func (d *policyRuleDie) DieRelease() rbacv1.PolicyRule {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

func (d *policyRuleDie) DieReleasePtr() *rbacv1.PolicyRule {
	r := d.DieRelease()
	return &r
}

func (d *policyRuleDie) DieStamp(fn func(r *rbacv1.PolicyRule)) PolicyRuleDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

func (d *policyRuleDie) DeepCopy() PolicyRuleDie {
	r := *d.r.DeepCopy()
	return &policyRuleDie{
		mutable: d.mutable,
		r:       r,
	}
}

func (d *policyRuleDie) Verbs(v ...string) PolicyRuleDie {
	return d.DieStamp(func(r *rbacv1.PolicyRule) {
		r.Verbs = v
	})
}

func (d *policyRuleDie) APIGroups(v ...string) PolicyRuleDie {
	return d.DieStamp(func(r *rbacv1.PolicyRule) {
		r.APIGroups = v
	})
}

func (d *policyRuleDie) Resources(v ...string) PolicyRuleDie {
	return d.DieStamp(func(r *rbacv1.PolicyRule) {
		r.Resources = v
	})
}

func (d *policyRuleDie) ResourceNames(v ...string) PolicyRuleDie {
	return d.DieStamp(func(r *rbacv1.PolicyRule) {
		r.ResourceNames = v
	})
}

func (d *policyRuleDie) NonResourceURLs(v ...string) PolicyRuleDie {
	return d.DieStamp(func(r *rbacv1.PolicyRule) {
		r.NonResourceURLs = v
	})
}

type RoleBindingDie interface {
	// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
	DieStamp(fn func(r *rbacv1.RoleBinding)) RoleBindingDie
	// DieFeed returns a new die with the provided resource.
	DieFeed(r rbacv1.RoleBinding) RoleBindingDie
	// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
	DieFeedPtr(r *rbacv1.RoleBinding) RoleBindingDie
	// DieRelease returns the resource managed by the die.
	DieRelease() rbacv1.RoleBinding
	// DieReleasePtr returns a pointer to the resource managed by the die.
	DieReleasePtr() *rbacv1.RoleBinding
	// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
	DieImmutable(immutable bool) RoleBindingDie
	// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
	DeepCopy() RoleBindingDie

	roleBinding
	// MetadataDie stamps the resource's ObjectMeta field with a mutable die.
	MetadataDie(fn func(d metav1.ObjectMetaDie)) RoleBindingDie
	// Subjects holds references to the objects the role applies to.
	Subjects(Subjects ...rbacv1.Subject) RoleBindingDie
	// RoleRef can reference a Role in the current namespace or a ClusterRole in the global namespace. If the RoleRef cannot be resolved, the Authorizer must return an error.
	RoleRef(RoleRef rbacv1.RoleRef) RoleBindingDie

	runtime.Object
	apismetav1.Object
	apismetav1.ObjectMetaAccessor
}

var _ RoleBindingDie = (*roleBindingDie)(nil)
var RoleBindingBlank = (&roleBindingDie{}).DieFeed(rbacv1.RoleBinding{})

type roleBindingDie struct {
	metav1.FrozenObjectMeta
	mutable bool
	r       rbacv1.RoleBinding
}

func (d *roleBindingDie) DieImmutable(immutable bool) RoleBindingDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy().(*roleBindingDie)
	d.mutable = !immutable
	return d
}

func (d *roleBindingDie) DieFeed(r rbacv1.RoleBinding) RoleBindingDie {
	if d.mutable {
		d.FrozenObjectMeta = metav1.FreezeObjectMeta(r.ObjectMeta)
		d.r = r
		return d
	}
	return &roleBindingDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

func (d *roleBindingDie) DieFeedPtr(r *rbacv1.RoleBinding) RoleBindingDie {
	if r == nil {
		r = &rbacv1.RoleBinding{}
	}
	return d.DieFeed(*r)
}

func (d *roleBindingDie) DieRelease() rbacv1.RoleBinding {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

func (d *roleBindingDie) DieReleasePtr() *rbacv1.RoleBinding {
	r := d.DieRelease()
	return &r
}

func (d *roleBindingDie) DieStamp(fn func(r *rbacv1.RoleBinding)) RoleBindingDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

func (d *roleBindingDie) DeepCopy() RoleBindingDie {
	r := *d.r.DeepCopy()
	return &roleBindingDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

func (d *roleBindingDie) DeepCopyObject() runtime.Object {
	return d.r.DeepCopy()
}

func (d *roleBindingDie) GetObjectKind() schema.ObjectKind {
	r := d.DieRelease()
	return r.GetObjectKind()
}

func (d *roleBindingDie) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.r)
}

func (d *roleBindingDie) UnmarshalJSON(b []byte) error {
	if d == RoleBindingBlank {
		return fmtx.Errorf("cannot unmarshal into the root object, create a copy first")
	}
	r := &rbacv1.RoleBinding{}
	err := json.Unmarshal(b, r)
	*d = *d.DieFeed(*r).(*roleBindingDie)
	return err
}

func (d *roleBindingDie) MetadataDie(fn func(d metav1.ObjectMetaDie)) RoleBindingDie {
	return d.DieStamp(func(r *rbacv1.RoleBinding) {
		d := metav1.ObjectMetaBlank.DieImmutable(false).DieFeed(r.ObjectMeta)
		fn(d)
		r.ObjectMeta = d.DieRelease()
	})
}

func (d *roleBindingDie) Subjects(v ...rbacv1.Subject) RoleBindingDie {
	return d.DieStamp(func(r *rbacv1.RoleBinding) {
		r.Subjects = v
	})
}

func (d *roleBindingDie) RoleRef(v rbacv1.RoleRef) RoleBindingDie {
	return d.DieStamp(func(r *rbacv1.RoleBinding) {
		r.RoleRef = v
	})
}

type SubjectDie interface {
	// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
	DieStamp(fn func(r *rbacv1.Subject)) SubjectDie
	// DieFeed returns a new die with the provided resource.
	DieFeed(r rbacv1.Subject) SubjectDie
	// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
	DieFeedPtr(r *rbacv1.Subject) SubjectDie
	// DieRelease returns the resource managed by the die.
	DieRelease() rbacv1.Subject
	// DieReleasePtr returns a pointer to the resource managed by the die.
	DieReleasePtr() *rbacv1.Subject
	// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
	DieImmutable(immutable bool) SubjectDie
	// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
	DeepCopy() SubjectDie

	// Kind of object being referenced. Values defined by this API group are "User", "Group", and "ServiceAccount". If the Authorizer does not recognized the kind value, the Authorizer should report an error.
	Kind(Kind string) SubjectDie
	// APIGroup holds the API group of the referenced subject. Defaults to "" for ServiceAccount subjects. Defaults to "rbac.authorization.k8s.io" for User and Group subjects.
	APIGroup(APIGroup string) SubjectDie
	// Name of the object being referenced.
	Name(Name string) SubjectDie
	// Namespace of the referenced object.  If the object kind is non-namespace, such as "User" or "Group", and this value is not empty the Authorizer should report an error.
	Namespace(Namespace string) SubjectDie
}

var _ SubjectDie = (*subjectDie)(nil)
var SubjectBlank = (&subjectDie{}).DieFeed(rbacv1.Subject{})

type subjectDie struct {
	mutable bool
	r       rbacv1.Subject
}

func (d *subjectDie) DieImmutable(immutable bool) SubjectDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy().(*subjectDie)
	d.mutable = !immutable
	return d
}

func (d *subjectDie) DieFeed(r rbacv1.Subject) SubjectDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &subjectDie{
		mutable: d.mutable,
		r:       r,
	}
}

func (d *subjectDie) DieFeedPtr(r *rbacv1.Subject) SubjectDie {
	if r == nil {
		r = &rbacv1.Subject{}
	}
	return d.DieFeed(*r)
}

func (d *subjectDie) DieRelease() rbacv1.Subject {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

func (d *subjectDie) DieReleasePtr() *rbacv1.Subject {
	r := d.DieRelease()
	return &r
}

func (d *subjectDie) DieStamp(fn func(r *rbacv1.Subject)) SubjectDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

func (d *subjectDie) DeepCopy() SubjectDie {
	r := *d.r.DeepCopy()
	return &subjectDie{
		mutable: d.mutable,
		r:       r,
	}
}

func (d *subjectDie) Kind(v string) SubjectDie {
	return d.DieStamp(func(r *rbacv1.Subject) {
		r.Kind = v
	})
}

func (d *subjectDie) APIGroup(v string) SubjectDie {
	return d.DieStamp(func(r *rbacv1.Subject) {
		r.APIGroup = v
	})
}

func (d *subjectDie) Name(v string) SubjectDie {
	return d.DieStamp(func(r *rbacv1.Subject) {
		r.Name = v
	})
}

func (d *subjectDie) Namespace(v string) SubjectDie {
	return d.DieStamp(func(r *rbacv1.Subject) {
		r.Namespace = v
	})
}

type RoleRefDie interface {
	// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
	DieStamp(fn func(r *rbacv1.RoleRef)) RoleRefDie
	// DieFeed returns a new die with the provided resource.
	DieFeed(r rbacv1.RoleRef) RoleRefDie
	// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
	DieFeedPtr(r *rbacv1.RoleRef) RoleRefDie
	// DieRelease returns the resource managed by the die.
	DieRelease() rbacv1.RoleRef
	// DieReleasePtr returns a pointer to the resource managed by the die.
	DieReleasePtr() *rbacv1.RoleRef
	// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
	DieImmutable(immutable bool) RoleRefDie
	// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
	DeepCopy() RoleRefDie

	// APIGroup is the group for the resource being referenced
	APIGroup(APIGroup string) RoleRefDie
	// Kind is the type of resource being referenced
	Kind(Kind string) RoleRefDie
	// Name is the name of resource being referenced
	Name(Name string) RoleRefDie
}

var _ RoleRefDie = (*roleRefDie)(nil)
var RoleRefBlank = (&roleRefDie{}).DieFeed(rbacv1.RoleRef{})

type roleRefDie struct {
	mutable bool
	r       rbacv1.RoleRef
}

func (d *roleRefDie) DieImmutable(immutable bool) RoleRefDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy().(*roleRefDie)
	d.mutable = !immutable
	return d
}

func (d *roleRefDie) DieFeed(r rbacv1.RoleRef) RoleRefDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &roleRefDie{
		mutable: d.mutable,
		r:       r,
	}
}

func (d *roleRefDie) DieFeedPtr(r *rbacv1.RoleRef) RoleRefDie {
	if r == nil {
		r = &rbacv1.RoleRef{}
	}
	return d.DieFeed(*r)
}

func (d *roleRefDie) DieRelease() rbacv1.RoleRef {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

func (d *roleRefDie) DieReleasePtr() *rbacv1.RoleRef {
	r := d.DieRelease()
	return &r
}

func (d *roleRefDie) DieStamp(fn func(r *rbacv1.RoleRef)) RoleRefDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

func (d *roleRefDie) DeepCopy() RoleRefDie {
	r := *d.r.DeepCopy()
	return &roleRefDie{
		mutable: d.mutable,
		r:       r,
	}
}

func (d *roleRefDie) APIGroup(v string) RoleRefDie {
	return d.DieStamp(func(r *rbacv1.RoleRef) {
		r.APIGroup = v
	})
}

func (d *roleRefDie) Kind(v string) RoleRefDie {
	return d.DieStamp(func(r *rbacv1.RoleRef) {
		r.Kind = v
	})
}

func (d *roleRefDie) Name(v string) RoleRefDie {
	return d.DieStamp(func(r *rbacv1.RoleRef) {
		r.Name = v
	})
}
