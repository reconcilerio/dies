//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2021-2022 the original author or authors.

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
	unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
)

var ClusterRoleBlank = (&ClusterRoleDie{}).DieFeed(rbacv1.ClusterRole{})

type ClusterRoleDie struct {
	metav1.FrozenObjectMeta
	mutable bool
	r       rbacv1.ClusterRole
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *ClusterRoleDie) DieImmutable(immutable bool) *ClusterRoleDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *ClusterRoleDie) DieFeed(r rbacv1.ClusterRole) *ClusterRoleDie {
	if d.mutable {
		d.FrozenObjectMeta = metav1.FreezeObjectMeta(r.ObjectMeta)
		d.r = r
		return d
	}
	return &ClusterRoleDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *ClusterRoleDie) DieFeedPtr(r *rbacv1.ClusterRole) *ClusterRoleDie {
	if r == nil {
		r = &rbacv1.ClusterRole{}
	}
	return d.DieFeed(*r)
}

// DieRelease returns the resource managed by the die.
func (d *ClusterRoleDie) DieRelease() rbacv1.ClusterRole {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *ClusterRoleDie) DieReleasePtr() *rbacv1.ClusterRole {
	r := d.DieRelease()
	return &r
}

// DieReleaseUnstructured returns the resource managed by the die as an unstructured object.
func (d *ClusterRoleDie) DieReleaseUnstructured() runtime.Unstructured {
	r := d.DieReleasePtr()
	u, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(r)
	return &unstructured.Unstructured{
		Object: u,
	}
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *ClusterRoleDie) DieStamp(fn func(r *rbacv1.ClusterRole)) *ClusterRoleDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *ClusterRoleDie) DeepCopy() *ClusterRoleDie {
	r := *d.r.DeepCopy()
	return &ClusterRoleDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

var _ runtime.Object = (*ClusterRoleDie)(nil)

func (d *ClusterRoleDie) DeepCopyObject() runtime.Object {
	return d.r.DeepCopy()
}

func (d *ClusterRoleDie) GetObjectKind() schema.ObjectKind {
	r := d.DieRelease()
	return r.GetObjectKind()
}

func (d *ClusterRoleDie) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.r)
}

func (d *ClusterRoleDie) UnmarshalJSON(b []byte) error {
	if d == ClusterRoleBlank {
		return fmtx.Errorf("cannot unmarshal into the blank die, create a copy first")
	}
	if !d.mutable {
		return fmtx.Errorf("cannot unmarshal into immutable dies, create a mutable version first")
	}
	r := &rbacv1.ClusterRole{}
	err := json.Unmarshal(b, r)
	*d = *d.DieFeed(*r)
	return err
}

// MetadataDie stamps the resource's ObjectMeta field with a mutable die.
func (d *ClusterRoleDie) MetadataDie(fn func(d *metav1.ObjectMetaDie)) *ClusterRoleDie {
	return d.DieStamp(func(r *rbacv1.ClusterRole) {
		d := metav1.ObjectMetaBlank.DieImmutable(false).DieFeed(r.ObjectMeta)
		fn(d)
		r.ObjectMeta = d.DieRelease()
	})
}

// Rules holds all the PolicyRules for this ClusterRole
func (d *ClusterRoleDie) Rules(v ...rbacv1.PolicyRule) *ClusterRoleDie {
	return d.DieStamp(func(r *rbacv1.ClusterRole) {
		r.Rules = v
	})
}

// AggregationRule is an optional field that describes how to build the Rules for this ClusterRole. If AggregationRule is set, then the Rules are controller managed and direct changes to Rules will be stomped by the controller.
func (d *ClusterRoleDie) AggregationRule(v *rbacv1.AggregationRule) *ClusterRoleDie {
	return d.DieStamp(func(r *rbacv1.ClusterRole) {
		r.AggregationRule = v
	})
}

var AggregationRuleBlank = (&AggregationRuleDie{}).DieFeed(rbacv1.AggregationRule{})

type AggregationRuleDie struct {
	mutable bool
	r       rbacv1.AggregationRule
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *AggregationRuleDie) DieImmutable(immutable bool) *AggregationRuleDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *AggregationRuleDie) DieFeed(r rbacv1.AggregationRule) *AggregationRuleDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &AggregationRuleDie{
		mutable: d.mutable,
		r:       r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *AggregationRuleDie) DieFeedPtr(r *rbacv1.AggregationRule) *AggregationRuleDie {
	if r == nil {
		r = &rbacv1.AggregationRule{}
	}
	return d.DieFeed(*r)
}

// DieRelease returns the resource managed by the die.
func (d *AggregationRuleDie) DieRelease() rbacv1.AggregationRule {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *AggregationRuleDie) DieReleasePtr() *rbacv1.AggregationRule {
	r := d.DieRelease()
	return &r
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *AggregationRuleDie) DieStamp(fn func(r *rbacv1.AggregationRule)) *AggregationRuleDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *AggregationRuleDie) DeepCopy() *AggregationRuleDie {
	r := *d.r.DeepCopy()
	return &AggregationRuleDie{
		mutable: d.mutable,
		r:       r,
	}
}

// ClusterRoleSelectors holds a list of selectors which will be used to find ClusterRoles and create the rules. If any of the selectors match, then the ClusterRole's permissions will be added
func (d *AggregationRuleDie) ClusterRoleSelectors(v ...apismetav1.LabelSelector) *AggregationRuleDie {
	return d.DieStamp(func(r *rbacv1.AggregationRule) {
		r.ClusterRoleSelectors = v
	})
}

var ClusterRoleBindingBlank = (&ClusterRoleBindingDie{}).DieFeed(rbacv1.ClusterRoleBinding{})

type ClusterRoleBindingDie struct {
	metav1.FrozenObjectMeta
	mutable bool
	r       rbacv1.ClusterRoleBinding
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *ClusterRoleBindingDie) DieImmutable(immutable bool) *ClusterRoleBindingDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *ClusterRoleBindingDie) DieFeed(r rbacv1.ClusterRoleBinding) *ClusterRoleBindingDie {
	if d.mutable {
		d.FrozenObjectMeta = metav1.FreezeObjectMeta(r.ObjectMeta)
		d.r = r
		return d
	}
	return &ClusterRoleBindingDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *ClusterRoleBindingDie) DieFeedPtr(r *rbacv1.ClusterRoleBinding) *ClusterRoleBindingDie {
	if r == nil {
		r = &rbacv1.ClusterRoleBinding{}
	}
	return d.DieFeed(*r)
}

// DieRelease returns the resource managed by the die.
func (d *ClusterRoleBindingDie) DieRelease() rbacv1.ClusterRoleBinding {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *ClusterRoleBindingDie) DieReleasePtr() *rbacv1.ClusterRoleBinding {
	r := d.DieRelease()
	return &r
}

// DieReleaseUnstructured returns the resource managed by the die as an unstructured object.
func (d *ClusterRoleBindingDie) DieReleaseUnstructured() runtime.Unstructured {
	r := d.DieReleasePtr()
	u, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(r)
	return &unstructured.Unstructured{
		Object: u,
	}
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *ClusterRoleBindingDie) DieStamp(fn func(r *rbacv1.ClusterRoleBinding)) *ClusterRoleBindingDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *ClusterRoleBindingDie) DeepCopy() *ClusterRoleBindingDie {
	r := *d.r.DeepCopy()
	return &ClusterRoleBindingDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

var _ runtime.Object = (*ClusterRoleBindingDie)(nil)

func (d *ClusterRoleBindingDie) DeepCopyObject() runtime.Object {
	return d.r.DeepCopy()
}

func (d *ClusterRoleBindingDie) GetObjectKind() schema.ObjectKind {
	r := d.DieRelease()
	return r.GetObjectKind()
}

func (d *ClusterRoleBindingDie) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.r)
}

func (d *ClusterRoleBindingDie) UnmarshalJSON(b []byte) error {
	if d == ClusterRoleBindingBlank {
		return fmtx.Errorf("cannot unmarshal into the blank die, create a copy first")
	}
	if !d.mutable {
		return fmtx.Errorf("cannot unmarshal into immutable dies, create a mutable version first")
	}
	r := &rbacv1.ClusterRoleBinding{}
	err := json.Unmarshal(b, r)
	*d = *d.DieFeed(*r)
	return err
}

// MetadataDie stamps the resource's ObjectMeta field with a mutable die.
func (d *ClusterRoleBindingDie) MetadataDie(fn func(d *metav1.ObjectMetaDie)) *ClusterRoleBindingDie {
	return d.DieStamp(func(r *rbacv1.ClusterRoleBinding) {
		d := metav1.ObjectMetaBlank.DieImmutable(false).DieFeed(r.ObjectMeta)
		fn(d)
		r.ObjectMeta = d.DieRelease()
	})
}

// Subjects holds references to the objects the role applies to.
func (d *ClusterRoleBindingDie) Subjects(v ...rbacv1.Subject) *ClusterRoleBindingDie {
	return d.DieStamp(func(r *rbacv1.ClusterRoleBinding) {
		r.Subjects = v
	})
}

// RoleRef can only reference a ClusterRole in the global namespace. If the RoleRef cannot be resolved, the Authorizer must return an error.
func (d *ClusterRoleBindingDie) RoleRef(v rbacv1.RoleRef) *ClusterRoleBindingDie {
	return d.DieStamp(func(r *rbacv1.ClusterRoleBinding) {
		r.RoleRef = v
	})
}

var RoleBlank = (&RoleDie{}).DieFeed(rbacv1.Role{})

type RoleDie struct {
	metav1.FrozenObjectMeta
	mutable bool
	r       rbacv1.Role
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *RoleDie) DieImmutable(immutable bool) *RoleDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *RoleDie) DieFeed(r rbacv1.Role) *RoleDie {
	if d.mutable {
		d.FrozenObjectMeta = metav1.FreezeObjectMeta(r.ObjectMeta)
		d.r = r
		return d
	}
	return &RoleDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *RoleDie) DieFeedPtr(r *rbacv1.Role) *RoleDie {
	if r == nil {
		r = &rbacv1.Role{}
	}
	return d.DieFeed(*r)
}

// DieRelease returns the resource managed by the die.
func (d *RoleDie) DieRelease() rbacv1.Role {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *RoleDie) DieReleasePtr() *rbacv1.Role {
	r := d.DieRelease()
	return &r
}

// DieReleaseUnstructured returns the resource managed by the die as an unstructured object.
func (d *RoleDie) DieReleaseUnstructured() runtime.Unstructured {
	r := d.DieReleasePtr()
	u, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(r)
	return &unstructured.Unstructured{
		Object: u,
	}
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *RoleDie) DieStamp(fn func(r *rbacv1.Role)) *RoleDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *RoleDie) DeepCopy() *RoleDie {
	r := *d.r.DeepCopy()
	return &RoleDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

var _ runtime.Object = (*RoleDie)(nil)

func (d *RoleDie) DeepCopyObject() runtime.Object {
	return d.r.DeepCopy()
}

func (d *RoleDie) GetObjectKind() schema.ObjectKind {
	r := d.DieRelease()
	return r.GetObjectKind()
}

func (d *RoleDie) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.r)
}

func (d *RoleDie) UnmarshalJSON(b []byte) error {
	if d == RoleBlank {
		return fmtx.Errorf("cannot unmarshal into the blank die, create a copy first")
	}
	if !d.mutable {
		return fmtx.Errorf("cannot unmarshal into immutable dies, create a mutable version first")
	}
	r := &rbacv1.Role{}
	err := json.Unmarshal(b, r)
	*d = *d.DieFeed(*r)
	return err
}

// MetadataDie stamps the resource's ObjectMeta field with a mutable die.
func (d *RoleDie) MetadataDie(fn func(d *metav1.ObjectMetaDie)) *RoleDie {
	return d.DieStamp(func(r *rbacv1.Role) {
		d := metav1.ObjectMetaBlank.DieImmutable(false).DieFeed(r.ObjectMeta)
		fn(d)
		r.ObjectMeta = d.DieRelease()
	})
}

// Rules holds all the PolicyRules for this Role
func (d *RoleDie) Rules(v ...rbacv1.PolicyRule) *RoleDie {
	return d.DieStamp(func(r *rbacv1.Role) {
		r.Rules = v
	})
}

var PolicyRuleBlank = (&PolicyRuleDie{}).DieFeed(rbacv1.PolicyRule{})

type PolicyRuleDie struct {
	mutable bool
	r       rbacv1.PolicyRule
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *PolicyRuleDie) DieImmutable(immutable bool) *PolicyRuleDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *PolicyRuleDie) DieFeed(r rbacv1.PolicyRule) *PolicyRuleDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &PolicyRuleDie{
		mutable: d.mutable,
		r:       r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *PolicyRuleDie) DieFeedPtr(r *rbacv1.PolicyRule) *PolicyRuleDie {
	if r == nil {
		r = &rbacv1.PolicyRule{}
	}
	return d.DieFeed(*r)
}

// DieRelease returns the resource managed by the die.
func (d *PolicyRuleDie) DieRelease() rbacv1.PolicyRule {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *PolicyRuleDie) DieReleasePtr() *rbacv1.PolicyRule {
	r := d.DieRelease()
	return &r
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *PolicyRuleDie) DieStamp(fn func(r *rbacv1.PolicyRule)) *PolicyRuleDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *PolicyRuleDie) DeepCopy() *PolicyRuleDie {
	r := *d.r.DeepCopy()
	return &PolicyRuleDie{
		mutable: d.mutable,
		r:       r,
	}
}

// Verbs is a list of Verbs that apply to ALL the ResourceKinds contained in this rule. '*' represents all verbs.
func (d *PolicyRuleDie) Verbs(v ...string) *PolicyRuleDie {
	return d.DieStamp(func(r *rbacv1.PolicyRule) {
		r.Verbs = v
	})
}

// APIGroups is the name of the APIGroup that contains the resources.  If multiple API groups are specified, any action requested against one of the enumerated resources in any API group will be allowed.
func (d *PolicyRuleDie) APIGroups(v ...string) *PolicyRuleDie {
	return d.DieStamp(func(r *rbacv1.PolicyRule) {
		r.APIGroups = v
	})
}

// Resources is a list of resources this rule applies to. '*' represents all resources.
func (d *PolicyRuleDie) Resources(v ...string) *PolicyRuleDie {
	return d.DieStamp(func(r *rbacv1.PolicyRule) {
		r.Resources = v
	})
}

// ResourceNames is an optional white list of names that the rule applies to.  An empty set means that everything is allowed.
func (d *PolicyRuleDie) ResourceNames(v ...string) *PolicyRuleDie {
	return d.DieStamp(func(r *rbacv1.PolicyRule) {
		r.ResourceNames = v
	})
}

// NonResourceURLs is a set of partial urls that a user should have access to.  *s are allowed, but only as the full, final step in the path Since non-resource URLs are not namespaced, this field is only applicable for ClusterRoles referenced from a ClusterRoleBinding. Rules can either apply to API resources (such as "pods" or "secrets") or non-resource URL paths (such as "/api"),  but not both.
func (d *PolicyRuleDie) NonResourceURLs(v ...string) *PolicyRuleDie {
	return d.DieStamp(func(r *rbacv1.PolicyRule) {
		r.NonResourceURLs = v
	})
}

var RoleBindingBlank = (&RoleBindingDie{}).DieFeed(rbacv1.RoleBinding{})

type RoleBindingDie struct {
	metav1.FrozenObjectMeta
	mutable bool
	r       rbacv1.RoleBinding
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *RoleBindingDie) DieImmutable(immutable bool) *RoleBindingDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *RoleBindingDie) DieFeed(r rbacv1.RoleBinding) *RoleBindingDie {
	if d.mutable {
		d.FrozenObjectMeta = metav1.FreezeObjectMeta(r.ObjectMeta)
		d.r = r
		return d
	}
	return &RoleBindingDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *RoleBindingDie) DieFeedPtr(r *rbacv1.RoleBinding) *RoleBindingDie {
	if r == nil {
		r = &rbacv1.RoleBinding{}
	}
	return d.DieFeed(*r)
}

// DieRelease returns the resource managed by the die.
func (d *RoleBindingDie) DieRelease() rbacv1.RoleBinding {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *RoleBindingDie) DieReleasePtr() *rbacv1.RoleBinding {
	r := d.DieRelease()
	return &r
}

// DieReleaseUnstructured returns the resource managed by the die as an unstructured object.
func (d *RoleBindingDie) DieReleaseUnstructured() runtime.Unstructured {
	r := d.DieReleasePtr()
	u, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(r)
	return &unstructured.Unstructured{
		Object: u,
	}
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *RoleBindingDie) DieStamp(fn func(r *rbacv1.RoleBinding)) *RoleBindingDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *RoleBindingDie) DeepCopy() *RoleBindingDie {
	r := *d.r.DeepCopy()
	return &RoleBindingDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

var _ runtime.Object = (*RoleBindingDie)(nil)

func (d *RoleBindingDie) DeepCopyObject() runtime.Object {
	return d.r.DeepCopy()
}

func (d *RoleBindingDie) GetObjectKind() schema.ObjectKind {
	r := d.DieRelease()
	return r.GetObjectKind()
}

func (d *RoleBindingDie) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.r)
}

func (d *RoleBindingDie) UnmarshalJSON(b []byte) error {
	if d == RoleBindingBlank {
		return fmtx.Errorf("cannot unmarshal into the blank die, create a copy first")
	}
	if !d.mutable {
		return fmtx.Errorf("cannot unmarshal into immutable dies, create a mutable version first")
	}
	r := &rbacv1.RoleBinding{}
	err := json.Unmarshal(b, r)
	*d = *d.DieFeed(*r)
	return err
}

// MetadataDie stamps the resource's ObjectMeta field with a mutable die.
func (d *RoleBindingDie) MetadataDie(fn func(d *metav1.ObjectMetaDie)) *RoleBindingDie {
	return d.DieStamp(func(r *rbacv1.RoleBinding) {
		d := metav1.ObjectMetaBlank.DieImmutable(false).DieFeed(r.ObjectMeta)
		fn(d)
		r.ObjectMeta = d.DieRelease()
	})
}

// Subjects holds references to the objects the role applies to.
func (d *RoleBindingDie) Subjects(v ...rbacv1.Subject) *RoleBindingDie {
	return d.DieStamp(func(r *rbacv1.RoleBinding) {
		r.Subjects = v
	})
}

// RoleRef can reference a Role in the current namespace or a ClusterRole in the global namespace. If the RoleRef cannot be resolved, the Authorizer must return an error.
func (d *RoleBindingDie) RoleRef(v rbacv1.RoleRef) *RoleBindingDie {
	return d.DieStamp(func(r *rbacv1.RoleBinding) {
		r.RoleRef = v
	})
}

var SubjectBlank = (&SubjectDie{}).DieFeed(rbacv1.Subject{})

type SubjectDie struct {
	mutable bool
	r       rbacv1.Subject
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *SubjectDie) DieImmutable(immutable bool) *SubjectDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *SubjectDie) DieFeed(r rbacv1.Subject) *SubjectDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &SubjectDie{
		mutable: d.mutable,
		r:       r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *SubjectDie) DieFeedPtr(r *rbacv1.Subject) *SubjectDie {
	if r == nil {
		r = &rbacv1.Subject{}
	}
	return d.DieFeed(*r)
}

// DieRelease returns the resource managed by the die.
func (d *SubjectDie) DieRelease() rbacv1.Subject {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *SubjectDie) DieReleasePtr() *rbacv1.Subject {
	r := d.DieRelease()
	return &r
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *SubjectDie) DieStamp(fn func(r *rbacv1.Subject)) *SubjectDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *SubjectDie) DeepCopy() *SubjectDie {
	r := *d.r.DeepCopy()
	return &SubjectDie{
		mutable: d.mutable,
		r:       r,
	}
}

// Kind of object being referenced. Values defined by this API group are "User", "Group", and "ServiceAccount". If the Authorizer does not recognized the kind value, the Authorizer should report an error.
func (d *SubjectDie) Kind(v string) *SubjectDie {
	return d.DieStamp(func(r *rbacv1.Subject) {
		r.Kind = v
	})
}

// APIGroup holds the API group of the referenced subject. Defaults to "" for ServiceAccount subjects. Defaults to "rbac.authorization.k8s.io" for User and Group subjects.
func (d *SubjectDie) APIGroup(v string) *SubjectDie {
	return d.DieStamp(func(r *rbacv1.Subject) {
		r.APIGroup = v
	})
}

// Name of the object being referenced.
func (d *SubjectDie) Name(v string) *SubjectDie {
	return d.DieStamp(func(r *rbacv1.Subject) {
		r.Name = v
	})
}

// Namespace of the referenced object.  If the object kind is non-namespace, such as "User" or "Group", and this value is not empty the Authorizer should report an error.
func (d *SubjectDie) Namespace(v string) *SubjectDie {
	return d.DieStamp(func(r *rbacv1.Subject) {
		r.Namespace = v
	})
}

var RoleRefBlank = (&RoleRefDie{}).DieFeed(rbacv1.RoleRef{})

type RoleRefDie struct {
	mutable bool
	r       rbacv1.RoleRef
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *RoleRefDie) DieImmutable(immutable bool) *RoleRefDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *RoleRefDie) DieFeed(r rbacv1.RoleRef) *RoleRefDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &RoleRefDie{
		mutable: d.mutable,
		r:       r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *RoleRefDie) DieFeedPtr(r *rbacv1.RoleRef) *RoleRefDie {
	if r == nil {
		r = &rbacv1.RoleRef{}
	}
	return d.DieFeed(*r)
}

// DieRelease returns the resource managed by the die.
func (d *RoleRefDie) DieRelease() rbacv1.RoleRef {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *RoleRefDie) DieReleasePtr() *rbacv1.RoleRef {
	r := d.DieRelease()
	return &r
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *RoleRefDie) DieStamp(fn func(r *rbacv1.RoleRef)) *RoleRefDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *RoleRefDie) DeepCopy() *RoleRefDie {
	r := *d.r.DeepCopy()
	return &RoleRefDie{
		mutable: d.mutable,
		r:       r,
	}
}

// APIGroup is the group for the resource being referenced
func (d *RoleRefDie) APIGroup(v string) *RoleRefDie {
	return d.DieStamp(func(r *rbacv1.RoleRef) {
		r.APIGroup = v
	})
}

// Kind is the type of resource being referenced
func (d *RoleRefDie) Kind(v string) *RoleRefDie {
	return d.DieStamp(func(r *rbacv1.RoleRef) {
		r.Kind = v
	})
}

// Name is the name of resource being referenced
func (d *RoleRefDie) Name(v string) *RoleRefDie {
	return d.DieStamp(func(r *rbacv1.RoleRef) {
		r.Name = v
	})
}