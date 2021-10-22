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
type _ = admissionregistrationv1.MutatingWebhookConfiguration

type mutatingWebhookConfigurationDieExtension interface {
	WebhookDie(name string, fn func(d MutatingWebhookDie)) MutatingWebhookConfigurationDie
}

func (d *mutatingWebhookConfigurationDie) WebhookDie(name string, fn func(d MutatingWebhookDie)) MutatingWebhookConfigurationDie {
	return d.DieStamp(func(r *admissionregistrationv1.MutatingWebhookConfiguration) {
		for i := range r.Webhooks {
			if name == r.Webhooks[i].Name {
				d := MutatingWebhookBlank.DieImmutable(false).DieFeed(r.Webhooks[i])
				fn(d)
				r.Webhooks[i] = d.DieRelease()
				return
			}
		}

		d := MutatingWebhookBlank.DieImmutable(false).DieFeed(admissionregistrationv1.MutatingWebhook{Name: name})
		fn(d)
		r.Webhooks = append(r.Webhooks, d.DieRelease())
	})
}

// +die
type _ = admissionregistrationv1.MutatingWebhook

type mutatingWebhookDieExtension interface {
	ClientConfigDie(fn func(d WebhookClientConfigDie)) MutatingWebhookDie
	RulesDie(rules ...RuleWithOperationsDie) MutatingWebhookDie
	NamespaceSelectorDie(fn func(d diemetav1.LabelSelectorDie)) MutatingWebhookDie
	ObjectSelectorDie(fn func(d diemetav1.LabelSelectorDie)) MutatingWebhookDie
}

func (d *mutatingWebhookDie) ClientConfigDie(fn func(d WebhookClientConfigDie)) MutatingWebhookDie {
	return d.DieStamp(func(r *admissionregistrationv1.MutatingWebhook) {
		d := WebhookClientConfigBlank.DieImmutable(false).DieFeed(r.ClientConfig)
		fn(d)
		r.ClientConfig = d.DieRelease()
	})
}

func (d *mutatingWebhookDie) RulesDie(rules ...RuleWithOperationsDie) MutatingWebhookDie {
	return d.DieStamp(func(r *admissionregistrationv1.MutatingWebhook) {
		r.Rules = make([]admissionregistrationv1.RuleWithOperations, len(rules))
		for i := range rules {
			r.Rules[i] = rules[i].DieRelease()
		}
	})
}

func (d *mutatingWebhookDie) NamespaceSelectorDie(fn func(d diemetav1.LabelSelectorDie)) MutatingWebhookDie {
	return d.DieStamp(func(r *admissionregistrationv1.MutatingWebhook) {
		d := diemetav1.LabelSelectorBlank.DieImmutable(false).DieFeedPtr(r.NamespaceSelector)
		fn(d)
		r.NamespaceSelector = d.DieReleasePtr()
	})
}

func (d *mutatingWebhookDie) ObjectSelectorDie(fn func(d diemetav1.LabelSelectorDie)) MutatingWebhookDie {
	return d.DieStamp(func(r *admissionregistrationv1.MutatingWebhook) {
		d := diemetav1.LabelSelectorBlank.DieImmutable(false).DieFeedPtr(r.ObjectSelector)
		fn(d)
		r.ObjectSelector = d.DieReleasePtr()
	})
}
