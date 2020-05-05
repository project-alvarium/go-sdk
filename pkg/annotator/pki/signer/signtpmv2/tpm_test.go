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

package signtpmv2

import (
	"io"
	"testing"

	testInternal "github.com/project-alvarium/go-sdk/internal/pkg/test"
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/pki/verifier/verifypkcs1v15"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signtpmv2/factory"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signtpmv2/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signtpmv2/provisioner"
	"github.com/project-alvarium/go-sdk/pkg/hashprovider"
	"github.com/project-alvarium/go-sdk/pkg/hashprovider/sha256"
	identityProvider "github.com/project-alvarium/go-sdk/pkg/identityprovider/hash"
	"github.com/project-alvarium/go-sdk/pkg/test"

	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpmutil"
	"github.com/stretchr/testify/assert"
)

// isTPMAvailable checks if the test environment has a tpm device. If found, a new primary key is created in the device
// and its key handle, public key, readwriteCloser interface to the tmp device and a found flag are returned.
// If not found, empty values of handle, public key and readwriteCloser interface along with false flag are returned.
func isTPMAvailable(t *testing.T, path string) (string, tpmutil.Handle, []byte, io.ReadWriteCloser, bool) {
	rwc, err := factory.TPM(path)
	if err != nil {
		return path, provisioner.InvalidHandle, nil, nil, false
	}

	handle, publicKey, err := provisioner.GenerateNewKeyPair(rwc)
	if err != nil {
		assert.FailNow(t, "generateNewKeyPair failed", err)
	}

	return path, handle, provisioner.MarshalPublicKey(publicKey), rwc, true
}

// newSUT returns a new system under test with tpm readWriteCloser interface initialized.
func newSUT(
	hashProvider hashprovider.Contract,
	publicKey []byte,
	handle tpmutil.Handle,
	path string,
	requestedCapabilityProperties RequestedCapabilityProperties,
	rwc io.ReadWriteCloser) *signer {

	return NewWithRWC(hashProvider, publicKey, handle, path, requestedCapabilityProperties, rwc)
}

// newSUTWithoutRWC return a new system under test with tpm readWriteCloser interface uninitialized.
func newSUTWithoutRWC(
	hashProvider hashprovider.Contract,
	publicKey []byte,
	handle tpmutil.Handle,
	path string,
	requestedCapabilityProperties RequestedCapabilityProperties) *signer {

	return New(hashProvider, publicKey, handle, path, requestedCapabilityProperties)
}

// TestSigner_SetUp tests signtpmv2.SetUp.
func TestSigner_SetUp(t *testing.T) {
	if path, handle, publicKey, rwc, ok := isTPMAvailable(t, provisioner.Path); ok {
		defer provisioner.FlushAndClose(rwc, handle)

		sut := newSUTWithoutRWC(sha256.New(), publicKey, handle, path, RequestedCapabilityProperties{})
		defer sut.TearDown()

		// coverage only; no assertions
		sut.SetUp()
	}
}

// TestSigner_TearDown tests signtpmv2.TearDown.
func TestSigner_TearDown(t *testing.T) {
	if path, handle, publicKey, rwc, ok := isTPMAvailable(t, provisioner.Path); ok {
		defer provisioner.FlushAndClose(rwc, handle)
		sut := newSUTWithoutRWC(sha256.New(), publicKey, handle, path, RequestedCapabilityProperties{})
		sut.SetUp()

		// for coverage; no assertion
		sut.TearDown()
	}
}

// TestSigner_PublicKey tests signtpmv2.PublicKey.
func TestSigner_PublicKey(t *testing.T) {
	if path, handle, publicKey, rwc, ok := isTPMAvailable(t, provisioner.Path); ok {
		defer provisioner.FlushAndClose(rwc, handle)
		sut := newSUT(sha256.New(), publicKey, handle, path, RequestedCapabilityProperties{}, rwc)

		assert.Equal(t, publicKey, sut.PublicKey())
	}
}

// TestSigner_Metadata tests signtpmv2.Metadata.
func TestSigner_Metadata(t *testing.T) {
	assertCapabilityPropertiesHas := func(t *testing.T, result interface{}, name string) {
		actual := result.(*metadata.Success)
		if assert.NotNil(t, actual) {
			if assert.NotNil(t, actual) {
				_, exists := actual.CapabilityProperties[name]
				assert.True(t, exists)
			}
		}
	}

	if path, handle, publicKey, rwc, ok := isTPMAvailable(t, provisioner.Path); ok {
		defer provisioner.FlushAndClose(rwc, handle)

		type testCase struct {
			name string
			test func(t *testing.T)
		}
		cases := []testCase{
			{
				name: "no capability properties",
				test: func(t *testing.T) {
					reducerHash := sha256.New()
					sut := newSUTWithoutRWC(reducerHash, publicKey, handle, path, nil)
					sut.SetUp()

					result := sut.Metadata()

					assert.Equal(
						t,
						testInternal.Marshal(
							t,
							metadata.NewSuccess(Name, reducerHash.Name(), metadata.CapabilityProperties{}),
						),
						testInternal.Marshal(t, result),
					)
				},
			},
			{
				name: "family indicator capability property",
				test: func(t *testing.T) {
					name := "family"
					reducerHash := sha256.New()
					sut := newSUTWithoutRWC(
						reducerHash,
						publicKey,
						handle,
						"",
						RequestedCapabilityProperties{name: tpm2.FamilyIndicator},
					)
					sut.SetUp()
					expectedCapabilityProperties := sut.getCapabilityProperties()

					result := sut.Metadata()

					assert.Equal(
						t,
						testInternal.Marshal(
							t,
							metadata.NewSuccess(Name, reducerHash.Name(), expectedCapabilityProperties),
						),
						testInternal.Marshal(t, result),
					)
					assertCapabilityPropertiesHas(t, result, name)
				},
			},
			{
				name: "manufacturer capability property",
				test: func(t *testing.T) {
					name := "manufacturer"
					reducerHash := sha256.New()
					sut := newSUTWithoutRWC(
						reducerHash,
						publicKey,
						handle,
						"",
						RequestedCapabilityProperties{name: tpm2.Manufacturer},
					)
					sut.SetUp()
					expectedCapabilityProperties := sut.getCapabilityProperties()

					result := sut.Metadata()

					assert.Equal(
						t,
						testInternal.Marshal(
							t,
							metadata.NewSuccess(Name, reducerHash.Name(), expectedCapabilityProperties),
						),
						testInternal.Marshal(t, result),
					)
					assertCapabilityPropertiesHas(t, result, name)
				},
			},
			{
				name: "failure (tpm signing error)",
				test: func(t *testing.T) {
					reducerHash := sha256.New()
					data := test.FactoryRandomByteSlice()
					id := identityProvider.New(reducerHash).Derive(data).Binary()
					sut := newSUT(reducerHash, publicKey, handle, path, RequestedCapabilityProperties{}, nil)
					_, _ = sut.Sign(id, data)

					result := sut.Metadata()

					assert.Equal(
						t,
						testInternal.Marshal(t, metadata.NewFailure(Name, "nil TPM handle")),
						testInternal.Marshal(t, result),
					)
				},
			},
		}

		for i := range cases {
			t.Run(cases[i].name, cases[i].test)
		}
	}
}

// TestSigner_ValidSignature tests signtpmv2.ValidSignature.
func TestSigner_ValidSignature(t *testing.T) {
	if path, handle, publicKey, rwc, ok := isTPMAvailable(t, provisioner.Path); ok {
		defer provisioner.FlushAndClose(rwc, handle)

		type testCase struct {
			name string
			test func(t *testing.T)
		}
		cases := []testCase{
			{
				name: "valid (signature not nil)",
				test: func(t *testing.T) {
					hashProvider := sha256.New()
					data := test.FactoryRandomByteSlice()
					id := identityProvider.New(hashProvider).Derive(data).Binary()
					sut := newSUT(hashProvider, publicKey, handle, path, RequestedCapabilityProperties{}, rwc)

					identitySignature, dataSignature := sut.Sign(id, data)

					assert.NotNil(t, identitySignature)
					assert.NotNil(t, dataSignature)
				},
			},
			{
				name: "valid (valid signature)",
				test: func(t *testing.T) {
					hashProvider := sha256.New()
					data := test.FactoryRandomByteSlice()
					id := identityProvider.New(hashProvider).Derive(data).Binary()
					sut := newSUT(hashProvider, publicKey, handle, path, RequestedCapabilityProperties{}, rwc)

					identitySignature, dataSignature := sut.Sign(id, data)

					v := verifypkcs1v15.New(provisioner.CryptoHash, hashProvider)
					assert.True(t, v.VerifyIdentity(id, identitySignature, publicKey))
					assert.True(t, v.VerifyData(data, dataSignature, publicKey))
				},
			},
			{
				name: "invalid (invalid handle)",
				test: func(t *testing.T) {
					hashProvider := sha256.New()
					data := test.FactoryRandomByteSlice()
					sut := newSUT(
						hashProvider,
						publicKey,
						provisioner.InvalidHandle,
						path,
						RequestedCapabilityProperties{},
						rwc,
					)

					identitySignature, dataSignature := sut.Sign(
						identityProvider.New(hashProvider).Derive(data).Binary(),
						data,
					)

					assert.Nil(t, identitySignature)
					assert.Nil(t, dataSignature)
				},
			},
		}

		for i := range cases {
			t.Run(cases[i].name, cases[i].test)
		}
	}
}
