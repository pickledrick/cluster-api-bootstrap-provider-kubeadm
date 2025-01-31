/*
Copyright 2019 The Kubernetes Authors.

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

package cloudinit

import (
	"bytes"
	"testing"

	infrav1 "sigs.k8s.io/cluster-api-bootstrap-provider-kubeadm/api/v1alpha2"
	"sigs.k8s.io/cluster-api-bootstrap-provider-kubeadm/certs"
)

func TestNewInitControlPlaneAdditionalFileEncodings(t *testing.T) {
	cpinput := &ControlPlaneInput{
		BaseUserData: BaseUserData{
			Header:              "test",
			PreKubeadmCommands:  nil,
			PostKubeadmCommands: nil,
			AdditionalFiles: []infrav1.File{
				{
					Path:     "/tmp/my-path",
					Encoding: infrav1.Base64,
					Content:  "aGk=",
				},
				{
					Path:    "/tmp/my-other-path",
					Content: "hi",
				},
			},
			WriteFiles: nil,
			Users:      nil,
			NTP:        nil,
		},
		Certificates: certs.Certificates{
			ClusterCA: &certs.KeyPair{
				Cert: []byte("ca cert"),
				Key:  []byte("ca key"),
			},
			EtcdCA: &certs.KeyPair{
				Cert: []byte("etcd ca cert"),
				Key:  []byte("etcd ca key"),
			},
			FrontProxyCA: &certs.KeyPair{
				Cert: []byte("front proxy ca cert"),
				Key:  []byte("front proxy ca key"),
			},
			ServiceAccount: &certs.KeyPair{
				Cert: []byte("service account ca cert"),
				Key:  []byte("service account ca key"),
			},
		},
		ClusterConfiguration: "my-cluster-config",
		InitConfiguration:    "my-init-config",
	}

	out, err := NewInitControlPlane(cpinput)
	if err != nil {
		t.Fatal(err)
	}
	expectedFiles := []string{
		`-   path: /tmp/my-path
    encoding: "base64"
    content: |
      aGk=`,
		`-   path: /tmp/my-other-path
    content: |
      hi`,
	}
	for _, f := range expectedFiles {
		if !bytes.Contains(out, []byte(f)) {
			t.Errorf("%s\ndid not contain\n%s", out, f)
		}
	}
}
