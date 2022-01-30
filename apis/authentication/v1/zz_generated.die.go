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
	authenticationv1 "k8s.io/api/authentication/v1"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
)

var TokenReviewBlank = (&TokenReviewDie{}).DieFeed(authenticationv1.TokenReview{})

type TokenReviewDie struct {
	metav1.FrozenObjectMeta
	mutable bool
	r       authenticationv1.TokenReview
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *TokenReviewDie) DieImmutable(immutable bool) *TokenReviewDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *TokenReviewDie) DieFeed(r authenticationv1.TokenReview) *TokenReviewDie {
	if d.mutable {
		d.FrozenObjectMeta = metav1.FreezeObjectMeta(r.ObjectMeta)
		d.r = r
		return d
	}
	return &TokenReviewDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *TokenReviewDie) DieFeedPtr(r *authenticationv1.TokenReview) *TokenReviewDie {
	if r == nil {
		r = &authenticationv1.TokenReview{}
	}
	return d.DieFeed(*r)
}

// DieRelease returns the resource managed by the die.
func (d *TokenReviewDie) DieRelease() authenticationv1.TokenReview {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *TokenReviewDie) DieReleasePtr() *authenticationv1.TokenReview {
	r := d.DieRelease()
	return &r
}

// DieReleaseUnstructured returns the resource managed by the die as an unstructured object.
func (d *TokenReviewDie) DieReleaseUnstructured() runtime.Unstructured {
	r := d.DieReleasePtr()
	u, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(r)
	return &unstructured.Unstructured{
		Object: u,
	}
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *TokenReviewDie) DieStamp(fn func(r *authenticationv1.TokenReview)) *TokenReviewDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *TokenReviewDie) DeepCopy() *TokenReviewDie {
	r := *d.r.DeepCopy()
	return &TokenReviewDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

var _ runtime.Object = (*TokenReviewDie)(nil)

func (d *TokenReviewDie) DeepCopyObject() runtime.Object {
	return d.r.DeepCopy()
}

func (d *TokenReviewDie) GetObjectKind() schema.ObjectKind {
	r := d.DieRelease()
	return r.GetObjectKind()
}

func (d *TokenReviewDie) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.r)
}

func (d *TokenReviewDie) UnmarshalJSON(b []byte) error {
	if d == TokenReviewBlank {
		return fmtx.Errorf("cannot unmarshal into the blank die, create a copy first")
	}
	if !d.mutable {
		return fmtx.Errorf("cannot unmarshal into immutable dies, create a mutable version first")
	}
	r := &authenticationv1.TokenReview{}
	err := json.Unmarshal(b, r)
	*d = *d.DieFeed(*r)
	return err
}

// MetadataDie stamps the resource's ObjectMeta field with a mutable die.
func (d *TokenReviewDie) MetadataDie(fn func(d *metav1.ObjectMetaDie)) *TokenReviewDie {
	return d.DieStamp(func(r *authenticationv1.TokenReview) {
		d := metav1.ObjectMetaBlank.DieImmutable(false).DieFeed(r.ObjectMeta)
		fn(d)
		r.ObjectMeta = d.DieRelease()
	})
}

// Spec holds information about the request being evaluated
func (d *TokenReviewDie) Spec(v authenticationv1.TokenReviewSpec) *TokenReviewDie {
	return d.DieStamp(func(r *authenticationv1.TokenReview) {
		r.Spec = v
	})
}

// Status is filled in by the server and indicates whether the request can be authenticated.
func (d *TokenReviewDie) Status(v authenticationv1.TokenReviewStatus) *TokenReviewDie {
	return d.DieStamp(func(r *authenticationv1.TokenReview) {
		r.Status = v
	})
}

var TokenRequestSpecBlank = (&TokenRequestSpecDie{}).DieFeed(authenticationv1.TokenRequestSpec{})

type TokenRequestSpecDie struct {
	mutable bool
	r       authenticationv1.TokenRequestSpec
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *TokenRequestSpecDie) DieImmutable(immutable bool) *TokenRequestSpecDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *TokenRequestSpecDie) DieFeed(r authenticationv1.TokenRequestSpec) *TokenRequestSpecDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &TokenRequestSpecDie{
		mutable: d.mutable,
		r:       r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *TokenRequestSpecDie) DieFeedPtr(r *authenticationv1.TokenRequestSpec) *TokenRequestSpecDie {
	if r == nil {
		r = &authenticationv1.TokenRequestSpec{}
	}
	return d.DieFeed(*r)
}

// DieRelease returns the resource managed by the die.
func (d *TokenRequestSpecDie) DieRelease() authenticationv1.TokenRequestSpec {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *TokenRequestSpecDie) DieReleasePtr() *authenticationv1.TokenRequestSpec {
	r := d.DieRelease()
	return &r
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *TokenRequestSpecDie) DieStamp(fn func(r *authenticationv1.TokenRequestSpec)) *TokenRequestSpecDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *TokenRequestSpecDie) DeepCopy() *TokenRequestSpecDie {
	r := *d.r.DeepCopy()
	return &TokenRequestSpecDie{
		mutable: d.mutable,
		r:       r,
	}
}

// Audiences are the intendend audiences of the token. A recipient of a token must identitfy themself with an identifier in the list of audiences of the token, and otherwise should reject the token. A token issued for multiple audiences may be used to authenticate against any of the audiences listed but implies a high degree of trust between the target audiences.
func (d *TokenRequestSpecDie) Audiences(v ...string) *TokenRequestSpecDie {
	return d.DieStamp(func(r *authenticationv1.TokenRequestSpec) {
		r.Audiences = v
	})
}

// ExpirationSeconds is the requested duration of validity of the request. The token issuer may return a token with a different validity duration so a client needs to check the 'expiration' field in a response.
func (d *TokenRequestSpecDie) ExpirationSeconds(v *int64) *TokenRequestSpecDie {
	return d.DieStamp(func(r *authenticationv1.TokenRequestSpec) {
		r.ExpirationSeconds = v
	})
}

// BoundObjectRef is a reference to an object that the token will be bound to. The token will only be valid for as long as the bound object exists. NOTE: The API server's TokenReview endpoint will validate the BoundObjectRef, but other audiences may not. Keep ExpirationSeconds small if you want prompt revocation.
func (d *TokenRequestSpecDie) BoundObjectRef(v *authenticationv1.BoundObjectReference) *TokenRequestSpecDie {
	return d.DieStamp(func(r *authenticationv1.TokenRequestSpec) {
		r.BoundObjectRef = v
	})
}

var BoundObjectReferenceBlank = (&BoundObjectReferenceDie{}).DieFeed(authenticationv1.BoundObjectReference{})

type BoundObjectReferenceDie struct {
	mutable bool
	r       authenticationv1.BoundObjectReference
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *BoundObjectReferenceDie) DieImmutable(immutable bool) *BoundObjectReferenceDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *BoundObjectReferenceDie) DieFeed(r authenticationv1.BoundObjectReference) *BoundObjectReferenceDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &BoundObjectReferenceDie{
		mutable: d.mutable,
		r:       r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *BoundObjectReferenceDie) DieFeedPtr(r *authenticationv1.BoundObjectReference) *BoundObjectReferenceDie {
	if r == nil {
		r = &authenticationv1.BoundObjectReference{}
	}
	return d.DieFeed(*r)
}

// DieRelease returns the resource managed by the die.
func (d *BoundObjectReferenceDie) DieRelease() authenticationv1.BoundObjectReference {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *BoundObjectReferenceDie) DieReleasePtr() *authenticationv1.BoundObjectReference {
	r := d.DieRelease()
	return &r
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *BoundObjectReferenceDie) DieStamp(fn func(r *authenticationv1.BoundObjectReference)) *BoundObjectReferenceDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *BoundObjectReferenceDie) DeepCopy() *BoundObjectReferenceDie {
	r := *d.r.DeepCopy()
	return &BoundObjectReferenceDie{
		mutable: d.mutable,
		r:       r,
	}
}

// Kind of the referent. Valid kinds are 'Pod' and 'Secret'.
func (d *BoundObjectReferenceDie) Kind(v string) *BoundObjectReferenceDie {
	return d.DieStamp(func(r *authenticationv1.BoundObjectReference) {
		r.Kind = v
	})
}

// API version of the referent.
func (d *BoundObjectReferenceDie) APIVersion(v string) *BoundObjectReferenceDie {
	return d.DieStamp(func(r *authenticationv1.BoundObjectReference) {
		r.APIVersion = v
	})
}

// Name of the referent.
func (d *BoundObjectReferenceDie) Name(v string) *BoundObjectReferenceDie {
	return d.DieStamp(func(r *authenticationv1.BoundObjectReference) {
		r.Name = v
	})
}

// UID of the referent.
func (d *BoundObjectReferenceDie) UID(v types.UID) *BoundObjectReferenceDie {
	return d.DieStamp(func(r *authenticationv1.BoundObjectReference) {
		r.UID = v
	})
}

var TokenRequestStatusBlank = (&TokenRequestStatusDie{}).DieFeed(authenticationv1.TokenRequestStatus{})

type TokenRequestStatusDie struct {
	mutable bool
	r       authenticationv1.TokenRequestStatus
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *TokenRequestStatusDie) DieImmutable(immutable bool) *TokenRequestStatusDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *TokenRequestStatusDie) DieFeed(r authenticationv1.TokenRequestStatus) *TokenRequestStatusDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &TokenRequestStatusDie{
		mutable: d.mutable,
		r:       r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *TokenRequestStatusDie) DieFeedPtr(r *authenticationv1.TokenRequestStatus) *TokenRequestStatusDie {
	if r == nil {
		r = &authenticationv1.TokenRequestStatus{}
	}
	return d.DieFeed(*r)
}

// DieRelease returns the resource managed by the die.
func (d *TokenRequestStatusDie) DieRelease() authenticationv1.TokenRequestStatus {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *TokenRequestStatusDie) DieReleasePtr() *authenticationv1.TokenRequestStatus {
	r := d.DieRelease()
	return &r
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *TokenRequestStatusDie) DieStamp(fn func(r *authenticationv1.TokenRequestStatus)) *TokenRequestStatusDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *TokenRequestStatusDie) DeepCopy() *TokenRequestStatusDie {
	r := *d.r.DeepCopy()
	return &TokenRequestStatusDie{
		mutable: d.mutable,
		r:       r,
	}
}

// Token is the opaque bearer token.
func (d *TokenRequestStatusDie) Token(v string) *TokenRequestStatusDie {
	return d.DieStamp(func(r *authenticationv1.TokenRequestStatus) {
		r.Token = v
	})
}

// ExpirationTimestamp is the time of expiration of the returned token.
func (d *TokenRequestStatusDie) ExpirationTimestamp(v apismetav1.Time) *TokenRequestStatusDie {
	return d.DieStamp(func(r *authenticationv1.TokenRequestStatus) {
		r.ExpirationTimestamp = v
	})
}
