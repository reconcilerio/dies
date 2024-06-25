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
	admissionv1 "k8s.io/api/admission/v1"
)

// +die
// +die:field:name=Kind,package=_/meta/v1,die=GroupVersionKindDie
// +die:field:name=Resource,package=_/meta/v1,die=GroupVersionResourceDie
// +die:field:name=RequestKind,package=_/meta/v1,die=GroupVersionKindDie,pointer=true
// +die:field:name=RequestResource,package=_/meta/v1,die=GroupVersionResourceDie,pointer=true
// +die:field:name=UserInfo,package=_/authentication/v1,die=UserInfoDie
type _ = admissionv1.AdmissionRequest
