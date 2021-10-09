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

package v1_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	diebatchv1 "github.com/scothis/dies/apis/batch/v1"
	diemetav1 "github.com/scothis/dies/apis/meta/v1"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

func TestCronJob(t *testing.T) {
	tests := []struct {
		name     string
		die      *diebatchv1.CronJobDie
		expected batchv1.CronJob
	}{
		{
			name:     "empty",
			die:      diebatchv1.CronJobBlank,
			expected: batchv1.CronJob{},
		},
		{
			name: "object metadata",
			die: diebatchv1.CronJobBlank.
				MetadataDie(func(d *diemetav1.ObjectMetaDie) {
					d.Namespace("my-namespace")
					d.Name("my-name")
				}),
			expected: batchv1.CronJob{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "my-namespace",
					Name:      "my-name",
				},
			},
		},
		{
			name: "spec job template",
			die: diebatchv1.CronJobBlank.
				SpecDie(func(d *diebatchv1.CronJobSpecDie) {
					d.JobTemplateDie(func(d *diebatchv1.JobDie) {
						d.MetadataDie(func(d *diemetav1.ObjectMetaDie) {
							d.AddLabel("key", "value")
						})
						d.SpecDie(func(d *diebatchv1.JobSpecDie) {
							d.Parallelism(pointer.Int32(1))
						})
					})
				}),
			expected: batchv1.CronJob{
				Spec: batchv1.CronJobSpec{
					JobTemplate: batchv1.JobTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"key": "value",
							},
						},
						Spec: batchv1.JobSpec{
							Parallelism: pointer.Int32(1),
						},
					},
				},
			},
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			actual := c.die.DieRelease()
			if diff := cmp.Diff(c.expected, actual); diff != "" {
				t.Errorf("(-expected, +actual): %s", diff)
			}
		})
	}
}
