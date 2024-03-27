//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2021-2022 the original author or authors.

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

// Code generated by diegen. DO NOT EDIT.

package v1

import (
	json "encoding/json"
	fmtx "fmt"
	certificatesv1 "k8s.io/api/certificates/v1"
	unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	jsonpath "k8s.io/client-go/util/jsonpath"
	osx "os"
	metav1 "reconciler.io/dies/apis/meta/v1"
	reflectx "reflect"
	yaml "sigs.k8s.io/yaml"
)

var CertificateSigningRequestBlank = (&CertificateSigningRequestDie{}).DieFeed(certificatesv1.CertificateSigningRequest{})

type CertificateSigningRequestDie struct {
	metav1.FrozenObjectMeta
	mutable bool
	r       certificatesv1.CertificateSigningRequest
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *CertificateSigningRequestDie) DieImmutable(immutable bool) *CertificateSigningRequestDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *CertificateSigningRequestDie) DieFeed(r certificatesv1.CertificateSigningRequest) *CertificateSigningRequestDie {
	if d.mutable {
		d.FrozenObjectMeta = metav1.FreezeObjectMeta(r.ObjectMeta)
		d.r = r
		return d
	}
	return &CertificateSigningRequestDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *CertificateSigningRequestDie) DieFeedPtr(r *certificatesv1.CertificateSigningRequest) *CertificateSigningRequestDie {
	if r == nil {
		r = &certificatesv1.CertificateSigningRequest{}
	}
	return d.DieFeed(*r)
}

// DieFeedJSON returns a new die with the provided JSON. Panics on error.
func (d *CertificateSigningRequestDie) DieFeedJSON(j []byte) *CertificateSigningRequestDie {
	r := certificatesv1.CertificateSigningRequest{}
	if err := json.Unmarshal(j, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAML returns a new die with the provided YAML. Panics on error.
func (d *CertificateSigningRequestDie) DieFeedYAML(y []byte) *CertificateSigningRequestDie {
	r := certificatesv1.CertificateSigningRequest{}
	if err := yaml.Unmarshal(y, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAMLFile returns a new die loading YAML from a file path. Panics on error.
func (d *CertificateSigningRequestDie) DieFeedYAMLFile(name string) *CertificateSigningRequestDie {
	y, err := osx.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return d.DieFeedYAML(y)
}

// DieFeedRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *CertificateSigningRequestDie) DieFeedRawExtension(raw runtime.RawExtension) *CertificateSigningRequestDie {
	j, err := json.Marshal(raw)
	if err != nil {
		panic(err)
	}
	return d.DieFeedJSON(j)
}

// DieRelease returns the resource managed by the die.
func (d *CertificateSigningRequestDie) DieRelease() certificatesv1.CertificateSigningRequest {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *CertificateSigningRequestDie) DieReleasePtr() *certificatesv1.CertificateSigningRequest {
	r := d.DieRelease()
	return &r
}

// DieReleaseUnstructured returns the resource managed by the die as an unstructured object. Panics on error.
func (d *CertificateSigningRequestDie) DieReleaseUnstructured() *unstructured.Unstructured {
	r := d.DieReleasePtr()
	u, err := runtime.DefaultUnstructuredConverter.ToUnstructured(r)
	if err != nil {
		panic(err)
	}
	return &unstructured.Unstructured{
		Object: u,
	}
}

// DieReleaseJSON returns the resource managed by the die as JSON. Panics on error.
func (d *CertificateSigningRequestDie) DieReleaseJSON() []byte {
	r := d.DieReleasePtr()
	j, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return j
}

// DieReleaseYAML returns the resource managed by the die as YAML. Panics on error.
func (d *CertificateSigningRequestDie) DieReleaseYAML() []byte {
	r := d.DieReleasePtr()
	y, err := yaml.Marshal(r)
	if err != nil {
		panic(err)
	}
	return y
}

// DieReleaseRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *CertificateSigningRequestDie) DieReleaseRawExtension() runtime.RawExtension {
	j := d.DieReleaseJSON()
	raw := runtime.RawExtension{}
	if err := json.Unmarshal(j, &raw); err != nil {
		panic(err)
	}
	return raw
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *CertificateSigningRequestDie) DieStamp(fn func(r *certificatesv1.CertificateSigningRequest)) *CertificateSigningRequestDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// Experimental: DieStampAt uses a JSON path (http://goessner.net/articles/JsonPath/) expression to stamp portions of the resource. The callback is invoked with each JSON path match. Panics if the callback function does not accept a single argument of the same type or a pointer to that type as found on the resource at the target location.
//
// Future iterations will improve type coercion from the resource to the callback argument.
func (d *CertificateSigningRequestDie) DieStampAt(jp string, fn interface{}) *CertificateSigningRequestDie {
	return d.DieStamp(func(r *certificatesv1.CertificateSigningRequest) {
		if ni := reflectx.ValueOf(fn).Type().NumIn(); ni != 1 {
			panic(fmtx.Errorf("callback function must have 1 input parameters, found %d", ni))
		}
		if no := reflectx.ValueOf(fn).Type().NumOut(); no != 0 {
			panic(fmtx.Errorf("callback function must have 0 output parameters, found %d", no))
		}

		cp := jsonpath.New("")
		if err := cp.Parse(fmtx.Sprintf("{%s}", jp)); err != nil {
			panic(err)
		}
		cr, err := cp.FindResults(r)
		if err != nil {
			// errors are expected if a path is not found
			return
		}
		for _, cv := range cr[0] {
			arg0t := reflectx.ValueOf(fn).Type().In(0)

			var args []reflectx.Value
			if cv.Type().AssignableTo(arg0t) {
				args = []reflectx.Value{cv}
			} else if cv.CanAddr() && cv.Addr().Type().AssignableTo(arg0t) {
				args = []reflectx.Value{cv.Addr()}
			} else {
				panic(fmtx.Errorf("callback function must accept value of type %q, found type %q", cv.Type(), arg0t))
			}

			reflectx.ValueOf(fn).Call(args)
		}
	})
}

// DieWith returns a new die after passing the current die to the callback function. The passed die is mutable.
func (d *CertificateSigningRequestDie) DieWith(fns ...func(d *CertificateSigningRequestDie)) *CertificateSigningRequestDie {
	nd := CertificateSigningRequestBlank.DieFeed(d.DieRelease()).DieImmutable(false)
	for _, fn := range fns {
		if fn != nil {
			fn(nd)
		}
	}
	return d.DieFeed(nd.DieRelease())
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *CertificateSigningRequestDie) DeepCopy() *CertificateSigningRequestDie {
	r := *d.r.DeepCopy()
	return &CertificateSigningRequestDie{
		FrozenObjectMeta: metav1.FreezeObjectMeta(r.ObjectMeta),
		mutable:          d.mutable,
		r:                r,
	}
}

var _ runtime.Object = (*CertificateSigningRequestDie)(nil)

func (d *CertificateSigningRequestDie) DeepCopyObject() runtime.Object {
	return d.r.DeepCopy()
}

func (d *CertificateSigningRequestDie) GetObjectKind() schema.ObjectKind {
	r := d.DieRelease()
	return r.GetObjectKind()
}

func (d *CertificateSigningRequestDie) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.r)
}

func (d *CertificateSigningRequestDie) UnmarshalJSON(b []byte) error {
	if d == CertificateSigningRequestBlank {
		return fmtx.Errorf("cannot unmarshal into the blank die, create a copy first")
	}
	if !d.mutable {
		return fmtx.Errorf("cannot unmarshal into immutable dies, create a mutable version first")
	}
	r := &certificatesv1.CertificateSigningRequest{}
	err := json.Unmarshal(b, r)
	*d = *d.DieFeed(*r)
	return err
}

// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
func (d *CertificateSigningRequestDie) APIVersion(v string) *CertificateSigningRequestDie {
	return d.DieStamp(func(r *certificatesv1.CertificateSigningRequest) {
		r.APIVersion = v
	})
}

// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
func (d *CertificateSigningRequestDie) Kind(v string) *CertificateSigningRequestDie {
	return d.DieStamp(func(r *certificatesv1.CertificateSigningRequest) {
		r.Kind = v
	})
}

// MetadataDie stamps the resource's ObjectMeta field with a mutable die.
func (d *CertificateSigningRequestDie) MetadataDie(fn func(d *metav1.ObjectMetaDie)) *CertificateSigningRequestDie {
	return d.DieStamp(func(r *certificatesv1.CertificateSigningRequest) {
		d := metav1.ObjectMetaBlank.DieImmutable(false).DieFeed(r.ObjectMeta)
		fn(d)
		r.ObjectMeta = d.DieRelease()
	})
}

// SpecDie stamps the resource's spec field with a mutable die.
func (d *CertificateSigningRequestDie) SpecDie(fn func(d *CertificateSigningRequestSpecDie)) *CertificateSigningRequestDie {
	return d.DieStamp(func(r *certificatesv1.CertificateSigningRequest) {
		d := CertificateSigningRequestSpecBlank.DieImmutable(false).DieFeed(r.Spec)
		fn(d)
		r.Spec = d.DieRelease()
	})
}

// StatusDie stamps the resource's status field with a mutable die.
func (d *CertificateSigningRequestDie) StatusDie(fn func(d *CertificateSigningRequestStatusDie)) *CertificateSigningRequestDie {
	return d.DieStamp(func(r *certificatesv1.CertificateSigningRequest) {
		d := CertificateSigningRequestStatusBlank.DieImmutable(false).DieFeed(r.Status)
		fn(d)
		r.Status = d.DieRelease()
	})
}

// spec contains the certificate request, and is immutable after creation.
//
// Only the request, signerName, expirationSeconds, and usages fields can be set on creation.
//
// Other fields are derived by Kubernetes and cannot be modified by users.
func (d *CertificateSigningRequestDie) Spec(v certificatesv1.CertificateSigningRequestSpec) *CertificateSigningRequestDie {
	return d.DieStamp(func(r *certificatesv1.CertificateSigningRequest) {
		r.Spec = v
	})
}

// status contains information about whether the request is approved or denied,
//
// and the certificate issued by the signer, or the failure condition indicating signer failure.
func (d *CertificateSigningRequestDie) Status(v certificatesv1.CertificateSigningRequestStatus) *CertificateSigningRequestDie {
	return d.DieStamp(func(r *certificatesv1.CertificateSigningRequest) {
		r.Status = v
	})
}

var CertificateSigningRequestSpecBlank = (&CertificateSigningRequestSpecDie{}).DieFeed(certificatesv1.CertificateSigningRequestSpec{})

type CertificateSigningRequestSpecDie struct {
	mutable bool
	r       certificatesv1.CertificateSigningRequestSpec
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *CertificateSigningRequestSpecDie) DieImmutable(immutable bool) *CertificateSigningRequestSpecDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *CertificateSigningRequestSpecDie) DieFeed(r certificatesv1.CertificateSigningRequestSpec) *CertificateSigningRequestSpecDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &CertificateSigningRequestSpecDie{
		mutable: d.mutable,
		r:       r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *CertificateSigningRequestSpecDie) DieFeedPtr(r *certificatesv1.CertificateSigningRequestSpec) *CertificateSigningRequestSpecDie {
	if r == nil {
		r = &certificatesv1.CertificateSigningRequestSpec{}
	}
	return d.DieFeed(*r)
}

// DieFeedJSON returns a new die with the provided JSON. Panics on error.
func (d *CertificateSigningRequestSpecDie) DieFeedJSON(j []byte) *CertificateSigningRequestSpecDie {
	r := certificatesv1.CertificateSigningRequestSpec{}
	if err := json.Unmarshal(j, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAML returns a new die with the provided YAML. Panics on error.
func (d *CertificateSigningRequestSpecDie) DieFeedYAML(y []byte) *CertificateSigningRequestSpecDie {
	r := certificatesv1.CertificateSigningRequestSpec{}
	if err := yaml.Unmarshal(y, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAMLFile returns a new die loading YAML from a file path. Panics on error.
func (d *CertificateSigningRequestSpecDie) DieFeedYAMLFile(name string) *CertificateSigningRequestSpecDie {
	y, err := osx.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return d.DieFeedYAML(y)
}

// DieFeedRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *CertificateSigningRequestSpecDie) DieFeedRawExtension(raw runtime.RawExtension) *CertificateSigningRequestSpecDie {
	j, err := json.Marshal(raw)
	if err != nil {
		panic(err)
	}
	return d.DieFeedJSON(j)
}

// DieRelease returns the resource managed by the die.
func (d *CertificateSigningRequestSpecDie) DieRelease() certificatesv1.CertificateSigningRequestSpec {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *CertificateSigningRequestSpecDie) DieReleasePtr() *certificatesv1.CertificateSigningRequestSpec {
	r := d.DieRelease()
	return &r
}

// DieReleaseJSON returns the resource managed by the die as JSON. Panics on error.
func (d *CertificateSigningRequestSpecDie) DieReleaseJSON() []byte {
	r := d.DieReleasePtr()
	j, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return j
}

// DieReleaseYAML returns the resource managed by the die as YAML. Panics on error.
func (d *CertificateSigningRequestSpecDie) DieReleaseYAML() []byte {
	r := d.DieReleasePtr()
	y, err := yaml.Marshal(r)
	if err != nil {
		panic(err)
	}
	return y
}

// DieReleaseRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *CertificateSigningRequestSpecDie) DieReleaseRawExtension() runtime.RawExtension {
	j := d.DieReleaseJSON()
	raw := runtime.RawExtension{}
	if err := json.Unmarshal(j, &raw); err != nil {
		panic(err)
	}
	return raw
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *CertificateSigningRequestSpecDie) DieStamp(fn func(r *certificatesv1.CertificateSigningRequestSpec)) *CertificateSigningRequestSpecDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// Experimental: DieStampAt uses a JSON path (http://goessner.net/articles/JsonPath/) expression to stamp portions of the resource. The callback is invoked with each JSON path match. Panics if the callback function does not accept a single argument of the same type or a pointer to that type as found on the resource at the target location.
//
// Future iterations will improve type coercion from the resource to the callback argument.
func (d *CertificateSigningRequestSpecDie) DieStampAt(jp string, fn interface{}) *CertificateSigningRequestSpecDie {
	return d.DieStamp(func(r *certificatesv1.CertificateSigningRequestSpec) {
		if ni := reflectx.ValueOf(fn).Type().NumIn(); ni != 1 {
			panic(fmtx.Errorf("callback function must have 1 input parameters, found %d", ni))
		}
		if no := reflectx.ValueOf(fn).Type().NumOut(); no != 0 {
			panic(fmtx.Errorf("callback function must have 0 output parameters, found %d", no))
		}

		cp := jsonpath.New("")
		if err := cp.Parse(fmtx.Sprintf("{%s}", jp)); err != nil {
			panic(err)
		}
		cr, err := cp.FindResults(r)
		if err != nil {
			// errors are expected if a path is not found
			return
		}
		for _, cv := range cr[0] {
			arg0t := reflectx.ValueOf(fn).Type().In(0)

			var args []reflectx.Value
			if cv.Type().AssignableTo(arg0t) {
				args = []reflectx.Value{cv}
			} else if cv.CanAddr() && cv.Addr().Type().AssignableTo(arg0t) {
				args = []reflectx.Value{cv.Addr()}
			} else {
				panic(fmtx.Errorf("callback function must accept value of type %q, found type %q", cv.Type(), arg0t))
			}

			reflectx.ValueOf(fn).Call(args)
		}
	})
}

// DieWith returns a new die after passing the current die to the callback function. The passed die is mutable.
func (d *CertificateSigningRequestSpecDie) DieWith(fns ...func(d *CertificateSigningRequestSpecDie)) *CertificateSigningRequestSpecDie {
	nd := CertificateSigningRequestSpecBlank.DieFeed(d.DieRelease()).DieImmutable(false)
	for _, fn := range fns {
		if fn != nil {
			fn(nd)
		}
	}
	return d.DieFeed(nd.DieRelease())
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *CertificateSigningRequestSpecDie) DeepCopy() *CertificateSigningRequestSpecDie {
	r := *d.r.DeepCopy()
	return &CertificateSigningRequestSpecDie{
		mutable: d.mutable,
		r:       r,
	}
}

// request contains an x509 certificate signing request encoded in a "CERTIFICATE REQUEST" PEM block.
//
// When serialized as JSON or YAML, the data is additionally base64-encoded.
func (d *CertificateSigningRequestSpecDie) Request(v []byte) *CertificateSigningRequestSpecDie {
	return d.DieStamp(func(r *certificatesv1.CertificateSigningRequestSpec) {
		r.Request = v
	})
}

// signerName indicates the requested signer, and is a qualified name.
//
// List/watch requests for CertificateSigningRequests can filter on this field using a "spec.signerName=NAME" fieldSelector.
//
// Well-known Kubernetes signers are:
//
// 1. "kubernetes.io/kube-apiserver-client": issues client certificates that can be used to authenticate to kube-apiserver.
//
// Requests for this signer are never auto-approved by kube-controller-manager, can be issued by the "csrsigning" controller in kube-controller-manager.
//
// 2. "kubernetes.io/kube-apiserver-client-kubelet": issues client certificates that kubelets use to authenticate to kube-apiserver.
//
// Requests for this signer can be auto-approved by the "csrapproving" controller in kube-controller-manager, and can be issued by the "csrsigning" controller in kube-controller-manager.
//
// 3. "kubernetes.io/kubelet-serving" issues serving certificates that kubelets use to serve TLS endpoints, which kube-apiserver can connect to securely.
//
// Requests for this signer are never auto-approved by kube-controller-manager, and can be issued by the "csrsigning" controller in kube-controller-manager.
//
// More details are available at https://k8s.io/docs/reference/access-authn-authz/certificate-signing-requests/#kubernetes-signers
//
// Custom signerNames can also be specified. The signer defines:
//
// 1. Trust distribution: how trust (CA bundles) are distributed.
//
// 2. Permitted subjects: and behavior when a disallowed subject is requested.
//
// 3. Required, permitted, or forbidden x509 extensions in the request (including whether subjectAltNames are allowed, which types, restrictions on allowed values) and behavior when a disallowed extension is requested.
//
// 4. Required, permitted, or forbidden key usages / extended key usages.
//
// 5. Expiration/certificate lifetime: whether it is fixed by the signer, configurable by the admin.
//
// 6. Whether or not requests for CA certificates are allowed.
func (d *CertificateSigningRequestSpecDie) SignerName(v string) *CertificateSigningRequestSpecDie {
	return d.DieStamp(func(r *certificatesv1.CertificateSigningRequestSpec) {
		r.SignerName = v
	})
}

// expirationSeconds is the requested duration of validity of the issued
//
// certificate. The certificate signer may issue a certificate with a different
//
// validity duration so a client must check the delta between the notBefore and
//
// and notAfter fields in the issued certificate to determine the actual duration.
//
// The v1.22+ in-tree implementations of the well-known Kubernetes signers will
//
// honor this field as long as the requested duration is not greater than the
//
// maximum duration they will honor per the --cluster-signing-duration CLI
//
// flag to the Kubernetes controller manager.
//
// Certificate signers may not honor this field for various reasons:
//
// 1. Old signer that is unaware of the field (such as the in-tree
//
// implementations prior to v1.22)
//
// 2. Signer whose configured maximum is shorter than the requested duration
//
// 3. Signer whose configured minimum is longer than the requested duration
//
// The minimum valid value for expirationSeconds is 600, i.e. 10 minutes.
func (d *CertificateSigningRequestSpecDie) ExpirationSeconds(v *int32) *CertificateSigningRequestSpecDie {
	return d.DieStamp(func(r *certificatesv1.CertificateSigningRequestSpec) {
		r.ExpirationSeconds = v
	})
}

// usages specifies a set of key usages requested in the issued certificate.
//
// Requests for TLS client certificates typically request: "digital signature", "key encipherment", "client auth".
//
// Requests for TLS serving certificates typically request: "key encipherment", "digital signature", "server auth".
//
// Valid values are:
//
// "signing", "digital signature", "content commitment",
//
// "key encipherment", "key agreement", "data encipherment",
//
// "cert sign", "crl sign", "encipher only", "decipher only", "any",
//
// "server auth", "client auth",
//
// "code signing", "email protection", "s/mime",
//
// "ipsec end system", "ipsec tunnel", "ipsec user",
//
// "timestamping", "ocsp signing", "microsoft sgc", "netscape sgc"
func (d *CertificateSigningRequestSpecDie) Usages(v ...certificatesv1.KeyUsage) *CertificateSigningRequestSpecDie {
	return d.DieStamp(func(r *certificatesv1.CertificateSigningRequestSpec) {
		r.Usages = v
	})
}

// username contains the name of the user that created the CertificateSigningRequest.
//
// Populated by the API server on creation and immutable.
func (d *CertificateSigningRequestSpecDie) Username(v string) *CertificateSigningRequestSpecDie {
	return d.DieStamp(func(r *certificatesv1.CertificateSigningRequestSpec) {
		r.Username = v
	})
}

// uid contains the uid of the user that created the CertificateSigningRequest.
//
// Populated by the API server on creation and immutable.
func (d *CertificateSigningRequestSpecDie) UID(v string) *CertificateSigningRequestSpecDie {
	return d.DieStamp(func(r *certificatesv1.CertificateSigningRequestSpec) {
		r.UID = v
	})
}

// groups contains group membership of the user that created the CertificateSigningRequest.
//
// Populated by the API server on creation and immutable.
func (d *CertificateSigningRequestSpecDie) Groups(v ...string) *CertificateSigningRequestSpecDie {
	return d.DieStamp(func(r *certificatesv1.CertificateSigningRequestSpec) {
		r.Groups = v
	})
}

var CertificateSigningRequestStatusBlank = (&CertificateSigningRequestStatusDie{}).DieFeed(certificatesv1.CertificateSigningRequestStatus{})

type CertificateSigningRequestStatusDie struct {
	mutable bool
	r       certificatesv1.CertificateSigningRequestStatus
}

// DieImmutable returns a new die for the current die's state that is either mutable (`false`) or immutable (`true`).
func (d *CertificateSigningRequestStatusDie) DieImmutable(immutable bool) *CertificateSigningRequestStatusDie {
	if d.mutable == !immutable {
		return d
	}
	d = d.DeepCopy()
	d.mutable = !immutable
	return d
}

// DieFeed returns a new die with the provided resource.
func (d *CertificateSigningRequestStatusDie) DieFeed(r certificatesv1.CertificateSigningRequestStatus) *CertificateSigningRequestStatusDie {
	if d.mutable {
		d.r = r
		return d
	}
	return &CertificateSigningRequestStatusDie{
		mutable: d.mutable,
		r:       r,
	}
}

// DieFeedPtr returns a new die with the provided resource pointer. If the resource is nil, the empty value is used instead.
func (d *CertificateSigningRequestStatusDie) DieFeedPtr(r *certificatesv1.CertificateSigningRequestStatus) *CertificateSigningRequestStatusDie {
	if r == nil {
		r = &certificatesv1.CertificateSigningRequestStatus{}
	}
	return d.DieFeed(*r)
}

// DieFeedJSON returns a new die with the provided JSON. Panics on error.
func (d *CertificateSigningRequestStatusDie) DieFeedJSON(j []byte) *CertificateSigningRequestStatusDie {
	r := certificatesv1.CertificateSigningRequestStatus{}
	if err := json.Unmarshal(j, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAML returns a new die with the provided YAML. Panics on error.
func (d *CertificateSigningRequestStatusDie) DieFeedYAML(y []byte) *CertificateSigningRequestStatusDie {
	r := certificatesv1.CertificateSigningRequestStatus{}
	if err := yaml.Unmarshal(y, &r); err != nil {
		panic(err)
	}
	return d.DieFeed(r)
}

// DieFeedYAMLFile returns a new die loading YAML from a file path. Panics on error.
func (d *CertificateSigningRequestStatusDie) DieFeedYAMLFile(name string) *CertificateSigningRequestStatusDie {
	y, err := osx.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return d.DieFeedYAML(y)
}

// DieFeedRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *CertificateSigningRequestStatusDie) DieFeedRawExtension(raw runtime.RawExtension) *CertificateSigningRequestStatusDie {
	j, err := json.Marshal(raw)
	if err != nil {
		panic(err)
	}
	return d.DieFeedJSON(j)
}

// DieRelease returns the resource managed by the die.
func (d *CertificateSigningRequestStatusDie) DieRelease() certificatesv1.CertificateSigningRequestStatus {
	if d.mutable {
		return d.r
	}
	return *d.r.DeepCopy()
}

// DieReleasePtr returns a pointer to the resource managed by the die.
func (d *CertificateSigningRequestStatusDie) DieReleasePtr() *certificatesv1.CertificateSigningRequestStatus {
	r := d.DieRelease()
	return &r
}

// DieReleaseJSON returns the resource managed by the die as JSON. Panics on error.
func (d *CertificateSigningRequestStatusDie) DieReleaseJSON() []byte {
	r := d.DieReleasePtr()
	j, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return j
}

// DieReleaseYAML returns the resource managed by the die as YAML. Panics on error.
func (d *CertificateSigningRequestStatusDie) DieReleaseYAML() []byte {
	r := d.DieReleasePtr()
	y, err := yaml.Marshal(r)
	if err != nil {
		panic(err)
	}
	return y
}

// DieReleaseRawExtension returns the resource managed by the die as an raw extension. Panics on error.
func (d *CertificateSigningRequestStatusDie) DieReleaseRawExtension() runtime.RawExtension {
	j := d.DieReleaseJSON()
	raw := runtime.RawExtension{}
	if err := json.Unmarshal(j, &raw); err != nil {
		panic(err)
	}
	return raw
}

// DieStamp returns a new die with the resource passed to the callback function. The resource is mutable.
func (d *CertificateSigningRequestStatusDie) DieStamp(fn func(r *certificatesv1.CertificateSigningRequestStatus)) *CertificateSigningRequestStatusDie {
	r := d.DieRelease()
	fn(&r)
	return d.DieFeed(r)
}

// Experimental: DieStampAt uses a JSON path (http://goessner.net/articles/JsonPath/) expression to stamp portions of the resource. The callback is invoked with each JSON path match. Panics if the callback function does not accept a single argument of the same type or a pointer to that type as found on the resource at the target location.
//
// Future iterations will improve type coercion from the resource to the callback argument.
func (d *CertificateSigningRequestStatusDie) DieStampAt(jp string, fn interface{}) *CertificateSigningRequestStatusDie {
	return d.DieStamp(func(r *certificatesv1.CertificateSigningRequestStatus) {
		if ni := reflectx.ValueOf(fn).Type().NumIn(); ni != 1 {
			panic(fmtx.Errorf("callback function must have 1 input parameters, found %d", ni))
		}
		if no := reflectx.ValueOf(fn).Type().NumOut(); no != 0 {
			panic(fmtx.Errorf("callback function must have 0 output parameters, found %d", no))
		}

		cp := jsonpath.New("")
		if err := cp.Parse(fmtx.Sprintf("{%s}", jp)); err != nil {
			panic(err)
		}
		cr, err := cp.FindResults(r)
		if err != nil {
			// errors are expected if a path is not found
			return
		}
		for _, cv := range cr[0] {
			arg0t := reflectx.ValueOf(fn).Type().In(0)

			var args []reflectx.Value
			if cv.Type().AssignableTo(arg0t) {
				args = []reflectx.Value{cv}
			} else if cv.CanAddr() && cv.Addr().Type().AssignableTo(arg0t) {
				args = []reflectx.Value{cv.Addr()}
			} else {
				panic(fmtx.Errorf("callback function must accept value of type %q, found type %q", cv.Type(), arg0t))
			}

			reflectx.ValueOf(fn).Call(args)
		}
	})
}

// DieWith returns a new die after passing the current die to the callback function. The passed die is mutable.
func (d *CertificateSigningRequestStatusDie) DieWith(fns ...func(d *CertificateSigningRequestStatusDie)) *CertificateSigningRequestStatusDie {
	nd := CertificateSigningRequestStatusBlank.DieFeed(d.DieRelease()).DieImmutable(false)
	for _, fn := range fns {
		if fn != nil {
			fn(nd)
		}
	}
	return d.DieFeed(nd.DieRelease())
}

// DeepCopy returns a new die with equivalent state. Useful for snapshotting a mutable die.
func (d *CertificateSigningRequestStatusDie) DeepCopy() *CertificateSigningRequestStatusDie {
	r := *d.r.DeepCopy()
	return &CertificateSigningRequestStatusDie{
		mutable: d.mutable,
		r:       r,
	}
}

// conditions applied to the request. Known conditions are "Approved", "Denied", and "Failed".
func (d *CertificateSigningRequestStatusDie) Conditions(v ...certificatesv1.CertificateSigningRequestCondition) *CertificateSigningRequestStatusDie {
	return d.DieStamp(func(r *certificatesv1.CertificateSigningRequestStatus) {
		r.Conditions = v
	})
}

// certificate is populated with an issued certificate by the signer after an Approved condition is present.
//
// This field is set via the /status subresource. Once populated, this field is immutable.
//
// If the certificate signing request is denied, a condition of type "Denied" is added and this field remains empty.
//
// If the signer cannot issue the certificate, a condition of type "Failed" is added and this field remains empty.
//
// Validation requirements:
//
// 1. certificate must contain one or more PEM blocks.
//
// 2. All PEM blocks must have the "CERTIFICATE" label, contain no headers, and the encoded data
//
// must be a BER-encoded ASN.1 Certificate structure as described in section 4 of RFC5280.
//
// 3. Non-PEM content may appear before or after the "CERTIFICATE" PEM blocks and is unvalidated,
//
// to allow for explanatory text as described in section 5.2 of RFC7468.
//
// # If more than one PEM block is present, and the definition of the requested spec.signerName
//
// does not indicate otherwise, the first block is the issued certificate,
//
// and subsequent blocks should be treated as intermediate certificates and presented in TLS handshakes.
//
// The certificate is encoded in PEM format.
//
// When serialized as JSON or YAML, the data is additionally base64-encoded, so it consists of:
//
// base64(
//
// -----BEGIN CERTIFICATE-----
//
// ...
//
// -----END CERTIFICATE-----
//
// )
func (d *CertificateSigningRequestStatusDie) Certificate(v []byte) *CertificateSigningRequestStatusDie {
	return d.DieStamp(func(r *certificatesv1.CertificateSigningRequestStatus) {
		r.Certificate = v
	})
}
