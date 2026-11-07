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

// Package webhook serves a [conversion.Converter] as a CRD's conversion
// webhook, which the API server calls with apiextensions.k8s.io/v1
// conversion reviews.
package webhook

import (
	"encoding/json"
	"fmt"
	"net/http"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"reconciler.io/dies/v2/die/conversion"
)

// Review answers a conversion review with a converter: each object is
// converted to the desired group and version, or the response fails with why the
// first that couldn't be wasn't.
func Review(c *conversion.Converter, req *apiextensionsv1.ConversionRequest) *apiextensionsv1.ConversionResponse {
	resp := &apiextensionsv1.ConversionResponse{UID: req.UID}
	to, err := schema.ParseGroupVersion(req.DesiredAPIVersion)
	if err != nil {
		resp.Result = metav1.Status{Status: metav1.StatusFailure, Message: err.Error()}
		return resp
	}
	for i, obj := range req.Objects {
		out, err := c.Convert(obj.Raw, to)
		if err != nil {
			resp.ConvertedObjects = nil
			resp.Result = metav1.Status{Status: metav1.StatusFailure, Message: fmt.Sprintf("object %d: %v", i, err)}
			return resp
		}
		resp.ConvertedObjects = append(resp.ConvertedObjects, runtime.RawExtension{Raw: out})
	}
	resp.Result = metav1.Status{Status: metav1.StatusSuccess}
	return resp
}

// Handler serves a converter as a CRD's conversion webhook.
func Handler(c *conversion.Converter) http.Handler { return handler{c} }

type handler struct{ c *conversion.Converter }

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var review apiextensionsv1.ConversionReview
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if review.Request == nil {
		http.Error(w, "webhook: the review has no request", http.StatusBadRequest)
		return
	}
	out := apiextensionsv1.ConversionReview{
		APIVersion: apiextensionsv1.SchemeGroupVersion.String(),
		Kind:       "ConversionReview",
		Response:   Review(h.c, review.Request),
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}
