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

package harness

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	dieapiextensionsv1 "reconciler.io/dies/v2/apis/apiextensions/v1"
)

// CRDs are the CRDs the suites install, built with dies.
func CRDs() []*apiextensionsv1.CustomResourceDefinition {
	return []*apiextensionsv1.CustomResourceDefinition{WidgetCRD()}
}

// WidgetCRD is the CRD for examples/widget/v1, with the list maps its markers
// declare, so the server merges them as the dies do.
func WidgetCRD() *apiextensionsv1.CustomResourceDefinition {
	schema := dieapiextensionsv1.JSONSchemaPropsBlank.DieWith(func(d *dieapiextensionsv1.JSONSchemaPropsDieMut) {
		d.DieFeedYAML([]byte(`
type: object
properties:
  spec:
    type: object
    properties:
      ports:
        type: array
        x-kubernetes-list-type: map
        x-kubernetes-list-map-keys: [name]
        items:
          type: object
          required: [name]
          properties:
            name: {type: string}
            port: {type: integer, format: int32}
      template:
        type: object
        x-kubernetes-preserve-unknown-fields: true
      gadget:
        type: object
        x-kubernetes-preserve-unknown-fields: true
      match:
        type: object
        properties:
          equal: {type: boolean}
          note: {type: string, nullable: true}
  status:
    type: object
    properties:
      observedGeneration: {type: integer, format: int64}
      conditions:
        type: array
        x-kubernetes-list-type: map
        x-kubernetes-list-map-keys: [type]
        items:
          type: object
          required: [type, status]
          x-kubernetes-preserve-unknown-fields: true
          properties:
            type: {type: string}
            status: {type: string}
`))
	}).DieRelease()
	version := dieapiextensionsv1.CustomResourceDefinitionVersionBlank.DieWith(func(d *dieapiextensionsv1.CustomResourceDefinitionVersionDieMut) {
		d.Name().Set("v1")
		d.Served().Set(true)
		d.Storage().Set(true)
		d.Schema().OpenAPIV3Schema().DieFeed(schema)
		d.Subresources().Status().DieFeed(apiextensionsv1.CustomResourceSubresourceStatus{})
	}).DieRelease()
	return dieapiextensionsv1.CustomResourceDefinitionBlank.DieWith(func(d *dieapiextensionsv1.CustomResourceDefinitionDieMut) {
		d.Metadata().Name().Set("widgets.example.dev")
		d.Spec().Group().Set("example.dev")
		d.Spec().Scope().Set(dieapiextensionsv1.NamespaceScoped)
		d.Spec().Names().DieWith(func(n *dieapiextensionsv1.CustomResourceDefinitionNamesDieMut) {
			n.Kind().Set("Widget")
			n.ListKind().Set("WidgetList")
			n.Plural().Set("widgets")
			n.Singular().Set("widget")
		})
		d.Spec().Versions().Set([]apiextensionsv1.CustomResourceDefinitionVersion{version})
	}).DieReleasePtr()
}
