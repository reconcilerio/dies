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

package apiserver_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"time"

	"github.com/testcontainers/testcontainers-go"

	"reconciler.io/dies/v2/die/conversion/webhook"
	"reconciler.io/dies/v2/examples/sprocket"
)

// The webhooks the suite registers are served by the test process, on the
// host. The API server, in a container, reaches them at testcontainers' name
// for the host, through a port it forwards, over TLS with a certificate for
// that name, signed by a CA made for the run, which the registrations trust.

type webhooks struct {
	port     int
	caBundle []byte
	server   *http.Server
}

// url is the URL of a webhook's path, as the API server reaches it.
func (w *webhooks) url(path string) string {
	return fmt.Sprintf("https://%s:%d%s", testcontainers.HostInternal, w.port, path)
}

func (w *webhooks) stop() { _ = w.server.Close() }

func startWebhooks() (*webhooks, error) {
	caCert, caKey, caPEM, err := newCA()
	if err != nil {
		return nil, err
	}
	cert, err := newServingCert(caCert, caKey, testcontainers.HostInternal)
	if err != nil {
		return nil, err
	}
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	mux.Handle("/convert/sprockets", webhook.Handler(sprocket.Conversion))
	w := &webhooks{
		port:     l.Addr().(*net.TCPAddr).Port,
		caBundle: caPEM,
		server: &http.Server{
			Handler:           mux,
			TLSConfig:         &tls.Config{Certificates: []tls.Certificate{cert}, MinVersion: tls.VersionTLS12},
			ReadHeaderTimeout: 10 * time.Second,
		},
	}
	go func() { _ = w.server.ServeTLS(l, "", "") }()
	return w, nil
}

func newCA() (*x509.Certificate, *ecdsa.PrivateKey, []byte, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, nil, err
	}
	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "dies conformance webhooks"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		return nil, nil, nil, err
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return nil, nil, nil, err
	}
	return cert, key, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), nil
}

func newServingCert(ca *x509.Certificate, caKey *ecdsa.PrivateKey, host string) (tls.Certificate, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, err
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject:      pkix.Name{CommonName: host},
		DNSNames:     []string{host},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, ca, &key.PublicKey, caKey)
	if err != nil {
		return tls.Certificate{}, err
	}
	return tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key}, nil
}
