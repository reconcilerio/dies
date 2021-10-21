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
	diecorev1 "dies.dev/apis/core/v1"
	diemetav1 "dies.dev/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// +die:object=true
type _ = appsv1.DaemonSet

// +die
type _ = appsv1.DaemonSetSpec

type daemonSetSpec interface {
	SelectorDie(fn func(d diemetav1.LabelSelectorDie)) DaemonSetSpecDie
	TemplateDie(fn func(d diecorev1.PodTemplateSpecDie)) DaemonSetSpecDie
	UpdateStrategyDie(fn func(d DaemonSetUpdateStrategyDie)) DaemonSetSpecDie
}

func (d *daemonSetSpecDie) SelectorDie(fn func(d diemetav1.LabelSelectorDie)) DaemonSetSpecDie {
	return d.DieStamp(func(r *appsv1.DaemonSetSpec) {
		d := diemetav1.LabelSelectorBlank.DieImmutable(false).DieFeedPtr(r.Selector)
		fn(d)
		r.Selector = d.DieReleasePtr()
	})
}

func (d *daemonSetSpecDie) TemplateDie(fn func(d diecorev1.PodTemplateSpecDie)) DaemonSetSpecDie {
	return d.DieStamp(func(r *appsv1.DaemonSetSpec) {
		d := diecorev1.PodTemplateSpecBlank.DieImmutable(false).DieFeed(r.Template)
		fn(d)
		r.Template = d.DieRelease()
	})
}

func (d *daemonSetSpecDie) UpdateStrategyDie(fn func(d DaemonSetUpdateStrategyDie)) DaemonSetSpecDie {
	return d.DieStamp(func(r *appsv1.DaemonSetSpec) {
		d := DaemonSetUpdateStrategyBlank.DieImmutable(false).DieFeed(r.UpdateStrategy)
		fn(d)
		r.UpdateStrategy = d.DieRelease()
	})
}

// +die
type _ = appsv1.DaemonSetUpdateStrategy

type daemonSetUpdateStrategy interface {
	OnDelete() DaemonSetUpdateStrategyDie
	RollingUpdateDie(fn func(d RollingUpdateDaemonSetDie)) DaemonSetUpdateStrategyDie
}

func (d *daemonSetUpdateStrategyDie) OnDelete() DaemonSetUpdateStrategyDie {
	return d.DieStamp(func(r *appsv1.DaemonSetUpdateStrategy) {
		r.Type = appsv1.OnDeleteDaemonSetStrategyType
		r.RollingUpdate = nil
	})
}

func (d *daemonSetUpdateStrategyDie) RollingUpdateDie(fn func(d RollingUpdateDaemonSetDie)) DaemonSetUpdateStrategyDie {
	return d.DieStamp(func(r *appsv1.DaemonSetUpdateStrategy) {
		r.Type = appsv1.RollingUpdateDaemonSetStrategyType
		d := RollingUpdateDaemonSetBlank.DieImmutable(false).DieFeedPtr(r.RollingUpdate)
		fn(d)
		r.RollingUpdate = d.DieReleasePtr()
	})
}

// +die
type _ = appsv1.RollingUpdateDaemonSet

// +die
type _ = appsv1.DaemonSetStatus

type daemonSetStatus interface {
	ConditionsDie(conditions ...diemetav1.ConditionDie) DaemonSetStatusDie
}

func (d *daemonSetStatusDie) ConditionsDie(conditions ...diemetav1.ConditionDie) DaemonSetStatusDie {
	return d.DieStamp(func(r *appsv1.DaemonSetStatus) {
		r.Conditions = make([]appsv1.DaemonSetCondition, len(conditions))
		for i := range conditions {
			c := conditions[i].DieRelease()
			r.Conditions[i] = appsv1.DaemonSetCondition{
				Type:               appsv1.DaemonSetConditionType(c.Type),
				Status:             corev1.ConditionStatus(c.Status),
				Reason:             c.Reason,
				Message:            c.Message,
				LastTransitionTime: c.LastTransitionTime,
			}
		}
	})
}
