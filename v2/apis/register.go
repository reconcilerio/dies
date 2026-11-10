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
	"reconciler.io/dies/v2/die/schemes"
)

// SchemeBuilder registers the dies of every API group.
var SchemeBuilder = schemes.NewBuilder(
	admission.SchemeBuilder.Add,
	admissionregistration.SchemeBuilder.Add,
	apidiscovery.SchemeBuilder.Add,
	apiextensions.SchemeBuilder.Add,
	apiregistration.SchemeBuilder.Add,
	apiserverinternal.SchemeBuilder.Add,
	apps.SchemeBuilder.Add,
	authentication.SchemeBuilder.Add,
	authorization.SchemeBuilder.Add,
	autoscaling.SchemeBuilder.Add,
	batch.SchemeBuilder.Add,
	certificates.SchemeBuilder.Add,
	coordination.SchemeBuilder.Add,
	core.SchemeBuilder.Add,
	discovery.SchemeBuilder.Add,
	events.SchemeBuilder.Add,
	extensions.SchemeBuilder.Add,
	flowcontrol.SchemeBuilder.Add,
	imagepolicy.SchemeBuilder.Add,
	lifecycle.SchemeBuilder.Add,
	meta.SchemeBuilder.Add,
	networking.SchemeBuilder.Add,
	node.SchemeBuilder.Add,
	policy.SchemeBuilder.Add,
	rbac.SchemeBuilder.Add,
	resource.SchemeBuilder.Add,
	scheduling.SchemeBuilder.Add,
	storage.SchemeBuilder.Add,
	storagemigration.SchemeBuilder.Add,
)

// AddToScheme registers the dies of every API group with a scheme.
var AddToScheme = SchemeBuilder.AddToScheme
