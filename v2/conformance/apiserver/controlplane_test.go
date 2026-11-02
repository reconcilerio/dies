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
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
	"k8s.io/client-go/rest"
)

// The control plane is etcd and kube-apiserver, each in a container, so the
// suite needs Docker and nothing else: no binaries on the machine, and the
// same versions everywhere.
//
// The API server is the one for the Kubernetes version that the dies are
// generated from, the version of k8s.io/api the module requires, unless
// KUBE_APISERVER_IMAGE names another. ETCD_IMAGE overrides etcd's.

const defaultEtcdImage = "registry.k8s.io/etcd:3.6.16-0"

// errNoDocker is the cause when there's no Docker to run the control plane
// in.
var errNoDocker = errors.New("docker is not available")

type controlPlane struct {
	config     *rest.Config
	containers []testcontainers.Container
	network    *testcontainers.DockerNetwork
}

func startControlPlane(ctx context.Context) (cp *controlPlane, err error) {
	if err := dockerHealthy(ctx); err != nil {
		return nil, fmt.Errorf("%w: %v", errNoDocker, err)
	}
	cp = &controlPlane{}
	defer func() {
		if err != nil {
			cp.stop()
		}
	}()

	cp.network, err = network.New(ctx)
	if err != nil {
		return nil, err
	}

	etcd, err := testcontainers.Run(ctx, imageFromEnv("ETCD_IMAGE", defaultEtcdImage),
		network.WithNetwork([]string{"etcd"}, cp.network),
		testcontainers.WithEntrypoint("etcd"),
		testcontainers.WithCmd(
			"--name=etcd",
			"--data-dir=/var/lib/etcd",
			"--listen-client-urls=http://0.0.0.0:2379",
			"--advertise-client-urls=http://etcd:2379",
			// the data is thrown away with the container
			"--unsafe-no-fsync",
		),
		testcontainers.WithTmpfs(map[string]string{"/var/lib/etcd": "rw"}),
		testcontainers.WithWaitStrategy(wait.ForLog("ready to serve client requests").WithStartupTimeout(time.Minute)),
	)
	if etcd != nil {
		cp.containers = append(cp.containers, etcd)
	}
	if err != nil {
		return nil, withLogs(ctx, etcd, fmt.Errorf("starting etcd: %w", err))
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	saKey := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	token := randomToken()
	tokens := fmt.Sprintf("%s,conformance,conformance,\"system:masters\"\n", token)

	image, err := apiServerImage()
	if err != nil {
		return nil, err
	}
	apiServer, err := testcontainers.Run(ctx, image,
		network.WithNetwork([]string{"apiserver"}, cp.network),
		testcontainers.WithEntrypoint("kube-apiserver"),
		testcontainers.WithCmd(
			"--etcd-servers=http://etcd:2379",
			"--bind-address=0.0.0.0",
			"--secure-port=6443",
			// the serving certificate is self-signed, and not verified
			"--cert-dir=/tmp/certs",
			"--token-auth-file=/etc/kubernetes/tokens.csv",
			"--authorization-mode=RBAC",
			"--service-account-issuer=https://kubernetes.default.svc",
			"--service-account-key-file=/etc/kubernetes/sa.key",
			"--service-account-signing-key-file=/etc/kubernetes/sa.key",
			"--service-cluster-ip-range=10.0.0.0/24",
			// there's no controller manager to create service accounts
			"--disable-admission-plugins=ServiceAccount",
			// every API version is served, alpha and beta included, so every
			// die has a schema to conform to
			"--runtime-config=api/all=true",
		),
		testcontainers.WithFiles(
			testcontainers.ContainerFile{Reader: bytes.NewReader(saKey), ContainerFilePath: "/etc/kubernetes/sa.key", FileMode: 0o600},
			testcontainers.ContainerFile{Reader: strings.NewReader(tokens), ContainerFilePath: "/etc/kubernetes/tokens.csv", FileMode: 0o600},
		),
		testcontainers.WithExposedPorts("6443/tcp"),
		testcontainers.WithWaitStrategy(wait.ForHTTP("/readyz").WithPort("6443/tcp").
			WithTLS(true, &tls.Config{InsecureSkipVerify: true}).
			WithHeaders(map[string]string{"Authorization": "Bearer " + token}).
			WithStartupTimeout(2*time.Minute)),
	)
	if apiServer != nil {
		cp.containers = append(cp.containers, apiServer)
	}
	if err != nil {
		return nil, withLogs(ctx, apiServer, fmt.Errorf("starting kube-apiserver %s: %w", image, err))
	}

	endpoint, err := apiServer.PortEndpoint(ctx, "6443/tcp", "https")
	if err != nil {
		return nil, err
	}
	cp.config = &rest.Config{
		Host:            endpoint,
		BearerToken:     token,
		TLSClientConfig: rest.TLSClientConfig{Insecure: true},
		QPS:             100,
		Burst:           200,
	}
	return cp, nil
}

func (cp *controlPlane) stop() {
	for i := len(cp.containers) - 1; i >= 0; i-- {
		_ = testcontainers.TerminateContainer(cp.containers[i])
	}
	if cp.network != nil {
		_ = cp.network.Remove(context.Background())
	}
}

func dockerHealthy(ctx context.Context) (err error) {
	// the provider panics, rather than failing, on some broken setups
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	provider, err := testcontainers.ProviderDocker.GetProvider()
	if err != nil {
		return err
	}
	defer provider.Close()
	return provider.Health(ctx)
}

// apiServerImage is the kube-apiserver for the version of k8s.io/api that
// the dies are built with: v0.37.0 is Kubernetes v1.37.0.
func apiServerImage() (string, error) {
	if image := os.Getenv("KUBE_APISERVER_IMAGE"); image != "" {
		return image, nil
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "", errors.New("no build info to find the version of k8s.io/api in")
	}
	for _, dep := range info.Deps {
		if dep.Path != "k8s.io/api" {
			continue
		}
		minor, ok := strings.CutPrefix(dep.Version, "v0.")
		if !ok {
			return "", fmt.Errorf("k8s.io/api %s is not a Kubernetes release; set KUBE_APISERVER_IMAGE", dep.Version)
		}
		return "registry.k8s.io/kube-apiserver:v1." + minor, nil
	}
	return "", errors.New("k8s.io/api is not a dependency")
}

func imageFromEnv(name, image string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return image
}

func randomToken() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// withLogs adds a container's logs to the error it failed to start with.
func withLogs(ctx context.Context, c testcontainers.Container, err error) error {
	if c == nil {
		return err
	}
	r, lerr := c.Logs(ctx)
	if lerr != nil {
		return err
	}
	defer r.Close()
	logs, _ := io.ReadAll(r)
	const tail = 4096
	if len(logs) > tail {
		logs = logs[len(logs)-tail:]
	}
	return fmt.Errorf("%w\n%s", err, logs)
}
