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

package v1

import (
	authorizationv1 "k8s.io/api/authorization/v1"
)

// +die:object=true,apiVersion=authorization.k8s.io/v1,kind=SubjectAccessReview
type _ = authorizationv1.SubjectAccessReview

// TODO(scothis) fix import for maps with struct values, ignore the 'Extra' field until then

// +die:ignore={Extra}
// +die:field:name=ResourceAttributes,die=ResourceAttributesDie,pointer=true
// +die:field:name=NonResourceAttributes,die=NonResourceAttributesDie,pointer=true
// +die:field:name=AuthorizationOptions,die=AuthorizationOptionsDie,pointer=true
type _ = authorizationv1.SubjectAccessReviewSpec

// Extra corresponds to the user.Info.GetExtra() method from the authenticator.  Since that is input to the authorizer it needs a reflection here.
func (d *SubjectAccessReviewSpecDie) Extra(v map[string]authorizationv1.ExtraValue) *SubjectAccessReviewSpecDie {
	return d.DieStamp(func(r *authorizationv1.SubjectAccessReviewSpec) {
		r.Extra = v
	})
}

func (d *SubjectAccessReviewSpecDie) AddExtra(key string, value authorizationv1.ExtraValue) *SubjectAccessReviewSpecDie {
	return d.DieStamp(func(r *authorizationv1.SubjectAccessReviewSpec) {
		if r.Extra == nil {
			r.Extra = map[string]authorizationv1.ExtraValue{}
		}
		r.Extra[key] = value
	})
}

// +die
// +die:field:name=FieldSelector,die=FieldSelectorAttributesDie,pointer=true
// +die:field:name=LabelSelector,die=LabelSelectorAttributesDie,pointer=true
type _ = authorizationv1.ResourceAttributes

// +die
type _ = authorizationv1.NonResourceAttributes

// +die
type _ = authorizationv1.AuthorizationOptions

// +die
// +die:field:name=ConditionalDecision,die=ConditionsAwareDecisionDie,pointer=true
type _ = authorizationv1.SubjectAccessReviewStatus

// +die
// +die:field:name=Deny,die=UnconditionalDecisionDie,pointer=true
// +die:field:name=NoOpinion,die=UnconditionalDecisionDie,pointer=true
// +die:field:name=Allow,die=UnconditionalDecisionDie,pointer=true
// +die:field:name=ConditionsMap,die=ConditionsMapDie,pointer=true
// +die:field:name=Union,die=NamedConditionsAwareDecisionDie,listType=map,listMapKey=AuthorizerName
type _ = authorizationv1.ConditionsAwareDecision

// +die
type _ = authorizationv1.UnconditionalDecision

// +die
// +die:field:name=DenyConditions,die=ConditionDie,listType=map,listMapKey=ID
// +die:field:name=NoOpinionConditions,die=ConditionDie,listType=map,listMapKey=ID
// +die:field:name=AllowConditions,die=ConditionDie,listType=map,listMapKey=ID
type _ = authorizationv1.ConditionsMap

// +die
type _ = authorizationv1.Condition

// +die
// +die:field:name=Decision,die=ConditionsAwareDecisionDie
type _ = authorizationv1.NamedConditionsAwareDecision

// +die
// +die:field:name=Requirements,package=_/meta/v1,die=FieldSelectorRequirementDie,listType=atomic
type _ = authorizationv1.FieldSelectorAttributes

// +die
// +die:field:name=Requirements,package=_/meta/v1,die=LabelSelectorRequirementDie,listType=atomic
type _ = authorizationv1.LabelSelectorAttributes
