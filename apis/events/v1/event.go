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
	eventsv1 "k8s.io/api/events/v1"
	diecorev1 "reconciler.io/dies/apis/core/v1"
)

// +die:object=true
type _ = eventsv1.Event

func (d *EventDie) SeriesDie(fn func(d *EventSeriesDie)) *EventDie {
	return d.DieStamp(func(r *eventsv1.Event) {
		d := EventSeriesBlank.DieImmutable(false).DieFeedPtr(r.Series)
		fn(d)
		r.Series = d.DieReleasePtr()
	})
}

func (d *EventDie) RegardingDie(fn func(d *diecorev1.ObjectReferenceDie)) *EventDie {
	return d.DieStamp(func(r *eventsv1.Event) {
		d := diecorev1.ObjectReferenceBlank.DieImmutable(false).DieFeed(r.Regarding)
		fn(d)
		r.Regarding = d.DieRelease()
	})
}

func (d *EventDie) RelatedDie(fn func(d *diecorev1.ObjectReferenceDie)) *EventDie {
	return d.DieStamp(func(r *eventsv1.Event) {
		d := diecorev1.ObjectReferenceBlank.DieImmutable(false).DieFeedPtr(r.Related)
		fn(d)
		r.Related = d.DieReleasePtr()
	})
}

func (d *EventDie) DeprecatedSourceDie(fn func(d *diecorev1.EventSourceDie)) *EventDie {
	return d.DieStamp(func(r *eventsv1.Event) {
		d := diecorev1.EventSourceBlank.DieImmutable(false).DieFeed(r.DeprecatedSource)
		fn(d)
		r.DeprecatedSource = d.DieRelease()
	})
}

// +die
type _ = eventsv1.EventSeries
