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
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
)

// +die:object=true
type _ = admissionregistrationv1.ValidatingWebhookConfiguration

func (d *ValidatingWebhookConfigurationDie) WebhookDie(name string, fn func(d *ValidatingWebhookDie)) *ValidatingWebhookConfigurationDie {
	return d.DieStamp(func(r *admissionregistrationv1.ValidatingWebhookConfiguration) {
		for i := range r.Webhooks {
			if name == r.Webhooks[i].Name {
				d := ValidatingWebhookBlank.DieImmutable(false).DieFeed(r.Webhooks[i])
				fn(d)
				r.Webhooks[i] = d.DieRelease()
				return
			}
		}

		d := ValidatingWebhookBlank.DieImmutable(false).DieFeed(admissionregistrationv1.ValidatingWebhook{Name: name})
		fn(d)
		r.Webhooks = append(r.Webhooks, d.DieRelease())
	})
}

// +die
type _ = admissionregistrationv1.ValidatingWebhook

func (d *ValidatingWebhookDie) ClientConfigDie(fn func(d *WebhookClientConfigDie)) *ValidatingWebhookDie {
	return d.DieStamp(func(r *admissionregistrationv1.ValidatingWebhook) {
		d := WebhookClientConfigBlank.DieImmutable(false).DieFeed(r.ClientConfig)
		fn(d)
		r.ClientConfig = d.DieRelease()
	})
}

func (d *ValidatingWebhookDie) RulesDie(rules ...*RuleWithOperationsDie) *ValidatingWebhookDie {
	return d.DieStamp(func(r *admissionregistrationv1.ValidatingWebhook) {
		r.Rules = make([]admissionregistrationv1.RuleWithOperations, len(rules))
		for i := range rules {
			r.Rules[i] = rules[i].DieRelease()
		}
	})
}

func (d *ValidatingWebhookDie) NamespaceSelectorDie(fn func(d *diemetav1.LabelSelectorDie)) *ValidatingWebhookDie {
	return d.DieStamp(func(r *admissionregistrationv1.ValidatingWebhook) {
		d := diemetav1.LabelSelectorBlank.DieImmutable(false).DieFeedPtr(r.NamespaceSelector)
		fn(d)
		r.NamespaceSelector = d.DieReleasePtr()
	})
}

func (d *ValidatingWebhookDie) ObjectSelectorDie(fn func(d *diemetav1.LabelSelectorDie)) *ValidatingWebhookDie {
	return d.DieStamp(func(r *admissionregistrationv1.ValidatingWebhook) {
		d := diemetav1.LabelSelectorBlank.DieImmutable(false).DieFeedPtr(r.ObjectSelector)
		fn(d)
		r.ObjectSelector = d.DieReleasePtr()
	})
}

func (d *ValidatingWebhookDie) MatchConditionDie(name string, fn func(d *MatchConditionDie)) *ValidatingWebhookDie {
	return d.DieStamp(func(r *admissionregistrationv1.ValidatingWebhook) {
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
