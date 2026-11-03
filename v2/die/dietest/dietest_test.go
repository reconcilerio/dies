/*
Copyright 2026 the original author or authors.

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

package dietest_test

import (
	"strings"
	"testing"

	dieappsv1 "reconciler.io/dies/v2/apis/apps/v1"
	"reconciler.io/dies/v2/die"
	"reconciler.io/dies/v2/die/dietest"
)

var web = dietest.FeedYAMLFile(dieappsv1.DeploymentBlank, "testdata/web.yaml")

func TestFeedYAMLFile(t *testing.T) {
	if web.GetName() != "web" || web.Spec().Template().Spec().Containers().At("app").Image().Get() != "app:1" || !web.Spec().Paused().Has() {
		t.Errorf("fed %s", die.ReleaseJSON(web))
	}
	if !strings.Contains(string(die.ReleaseYAML(web)), "futureField: from a newer API") {
		t.Errorf("released %s", die.ReleaseYAML(web))
	}
	defer func() {
		if r := recover(); r == nil || !strings.Contains(r.(error).Error(), "no such file") {
			t.Errorf("recovered %v", r)
		}
	}()
	dietest.FeedYAMLFile(dieappsv1.DeploymentBlank, "testdata/missing.yaml")
}
