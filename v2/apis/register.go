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

package apis

import (
	"k8s.io/apimachinery/pkg/runtime"

	"reconciler.io/dies/v2/apis/admission"
	"reconciler.io/dies/v2/apis/admissionregistration"
	"reconciler.io/dies/v2/apis/apidiscovery"
	"reconciler.io/dies/v2/apis/apiextensions"
	"reconciler.io/dies/v2/apis/apiregistration"
	"reconciler.io/dies/v2/apis/apiserverinternal"
	"reconciler.io/dies/v2/apis/apps"
	"reconciler.io/dies/v2/apis/authentication"
	"reconciler.io/dies/v2/apis/authorization"
	"reconciler.io/dies/v2/apis/autoscaling"
	"reconciler.io/dies/v2/apis/batch"
	"reconciler.io/dies/v2/apis/certificates"
	"reconciler.io/dies/v2/apis/coordination"
	"reconciler.io/dies/v2/apis/core"
	"reconciler.io/dies/v2/apis/discovery"
	"reconciler.io/dies/v2/apis/events"
	"reconciler.io/dies/v2/apis/extensions"
	"reconciler.io/dies/v2/apis/flowcontrol"
	"reconciler.io/dies/v2/apis/imagepolicy"
	"reconciler.io/dies/v2/apis/lifecycle"
	"reconciler.io/dies/v2/apis/meta"
	"reconciler.io/dies/v2/apis/networking"
	"reconciler.io/dies/v2/apis/node"
	"reconciler.io/dies/v2/apis/policy"
	"reconciler.io/dies/v2/apis/rbac"
	"reconciler.io/dies/v2/apis/resource"
	"reconciler.io/dies/v2/apis/scheduling"
	"reconciler.io/dies/v2/apis/storage"
	"reconciler.io/dies/v2/apis/storagemigration"
)

// SchemeBuilder registers the dies of every API group.
var SchemeBuilder = runtime.NewSchemeBuilder(
	admission.AddToScheme,
	admissionregistration.AddToScheme,
	apidiscovery.AddToScheme,
	apiextensions.AddToScheme,
	apiregistration.AddToScheme,
	apiserverinternal.AddToScheme,
	apps.AddToScheme,
	authentication.AddToScheme,
	authorization.AddToScheme,
	autoscaling.AddToScheme,
	batch.AddToScheme,
	certificates.AddToScheme,
	coordination.AddToScheme,
	core.AddToScheme,
	discovery.AddToScheme,
	events.AddToScheme,
	extensions.AddToScheme,
	flowcontrol.AddToScheme,
	imagepolicy.AddToScheme,
	lifecycle.AddToScheme,
	meta.AddToScheme,
	networking.AddToScheme,
	node.AddToScheme,
	policy.AddToScheme,
	rbac.AddToScheme,
	resource.AddToScheme,
	scheduling.AddToScheme,
	storage.AddToScheme,
	storagemigration.AddToScheme,
)

// AddToScheme registers the dies of every API group with a scheme.
var AddToScheme = SchemeBuilder.AddToScheme
