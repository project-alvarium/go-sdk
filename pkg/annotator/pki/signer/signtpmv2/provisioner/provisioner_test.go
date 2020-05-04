/*******************************************************************************
 * Copyright 2020 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package provisioner

import (
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/project-alvarium/go-sdk/internal/pkg/test"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signtpmv2/factory"

	"github.com/google/go-tpm/tpm2"
	"github.com/stretchr/testify/assert"
)

// TestMarshalPublicKey tests MarshalPublicKey.
func TestMarshalPublicKey(t *testing.T) {
	type testCase struct {
		name string
		test func(t *testing.T)
	}
	cases := []testCase{
		{
			name: "Valid key",
			test: func(t *testing.T) {
				block, _ := pem.Decode(test.ValidPublicKey)
				if block == nil {
					assert.FailNow(t, "pem.Decode failed")
					return
				}
				rawKey, err := x509.ParsePKIXPublicKey(block.Bytes)
				if err != nil {
					assert.FailNow(t, "x509.ParsePKIXPublicKey failed")
					return
				}

				result := MarshalPublicKey(rawKey)

				assert.NotNil(t, result)
			},
		},
		{
			name: "Invalid Key",
			test: func(t *testing.T) {
				result := MarshalPublicKey("Not A Key")

				assert.Nil(t, result)
			},
		},
	}

	for i := range cases {
		t.Run(cases[i].name, cases[i].test)
	}
}

// TestGenerateNewKeyPair tests GenerateNewKeyPair.
func TestGenerateNewKeyPair(t *testing.T) {
	type testCase struct {
		name string
		test func(t *testing.T)
	}
	cases := []testCase{
		{
			name: "Valid rwc",
			test: func(t *testing.T) {
				rwc, err := factory.TPM(Path)
				if err != nil {
					assert.FailNow(t, "factory.TPM failed")
					return
				}

				handle, publicKey, err := GenerateNewKeyPair(rwc)

				assert.NotEqual(t, InvalidHandle, handle)
				assert.NotNil(t, publicKey)
				assert.Nil(t, err)

				_ = tpm2.FlushContext(rwc, handle)
				_ = rwc.Close()
			},
		},
		{
			name: "Invalid rwc",
			test: func(t *testing.T) {
				handle, publicKey, err := GenerateNewKeyPair(nil)

				assert.Equal(t, InvalidHandle, handle)
				assert.Nil(t, publicKey)
				assert.NotNil(t, err)
			},
		},
	}

	for i := range cases {
		t.Run(cases[i].name, cases[i].test)
	}
}
