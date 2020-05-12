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

package signpkcs1v15

import (
	"crypto"
	"testing"

	testInternal "github.com/project-alvarium/go-sdk/internal/pkg/test"
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/pki/verifier/verifypkcs1v15"
	pkcsSignerMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signpkcs1v15/metadata"
	"github.com/project-alvarium/go-sdk/pkg/hashprovider"
	"github.com/project-alvarium/go-sdk/pkg/hashprovider/passthrough"
	"github.com/project-alvarium/go-sdk/pkg/hashprovider/sha256"
	identityProvider "github.com/project-alvarium/go-sdk/pkg/identityprovider/hash"
	"github.com/project-alvarium/go-sdk/pkg/test"

	"github.com/stretchr/testify/assert"
)

// newSUT returns a new system under test.
func newSUT(hash crypto.Hash, privateKey, publicKey []byte, hashProvider hashprovider.Contract) *signer {
	return New(hash, privateKey, publicKey, hashProvider)
}

// TestSigner_New tests signpkcs1v15.New.
func TestSigner_New(t *testing.T) {
	type testCase struct {
		name       string
		privateKey []byte
		publicKey  []byte
		hash       crypto.Hash
		expected   *signer
	}

	cases := []testCase{
		func() testCase {
			return testCase{
				name:       "invalid (nil private key)",
				privateKey: []byte(nil),
				publicKey:  testInternal.ValidPublicKey,
				hash:       crypto.SHA256,
				expected:   nil,
			}
		}(),
		func() testCase {
			return testCase{
				name:       "invalid (invalid private key)",
				privateKey: testInternal.InvalidPrivateKey,
				publicKey:  testInternal.ValidPublicKey,
				hash:       crypto.SHA256,
				expected:   nil,
			}
		}(),
	}

	for i := range cases {
		t.Run(
			cases[i].name,
			func(t *testing.T) {
				result := newSUT(cases[i].hash, cases[i].privateKey, cases[i].publicKey, sha256.New())

				assert.Equal(t, cases[i].expected, result)
			},
		)
	}
}

// TestSigner_Sign tests signpkcs1v15.Sign.
func TestSigner_Sign(t *testing.T) {
	type testCase struct {
		name string
		test func(t *testing.T)
	}

	cases := []testCase{
		{
			name: "valid (signatures not nil)",
			test: func(t *testing.T) {
				hashProvider := sha256.New()
				data := test.FactoryRandomByteSlice()
				sut := newSUT(crypto.SHA256, testInternal.ValidPrivateKey, testInternal.ValidPublicKey, hashProvider)

				identitySignature, dataSignature := sut.Sign(
					identityProvider.New(hashProvider).Derive(data).Binary(),
					data,
				)

				assert.NotNil(t, testInternal.Encode(identitySignature))
				assert.NotNil(t, testInternal.Encode(dataSignature))
			},
		},
		{
			name: "valid (different data = different signatures)",
			test: func(t *testing.T) {
				hashProvider := sha256.New()
				idProvider := identityProvider.New(hashProvider)
				data1 := test.FactoryRandomByteSlice()
				data2 := test.FactoryRandomByteSlice()
				sut := newSUT(crypto.SHA256, testInternal.ValidPrivateKey, testInternal.ValidPublicKey, hashProvider)

				identitySignature1, dataSignature1 := sut.Sign(idProvider.Derive(data1).Binary(), data1)
				identitySignature2, dataSignature2 := sut.Sign(idProvider.Derive(data2).Binary(), data2)

				assert.NotEqual(t, identitySignature1, identitySignature2)
				assert.NotEqual(t, dataSignature1, dataSignature2)
			},
		},
		{
			name: "valid (signatures verified)",
			test: func(t *testing.T) {
				publicKey := testInternal.ValidPublicKey
				h := crypto.SHA256
				hashProvider := sha256.New()
				data := test.FactoryRandomByteSlice()
				id := identityProvider.New(hashProvider).Derive(data).Binary()
				sut := newSUT(h, testInternal.ValidPrivateKey, publicKey, hashProvider)

				identitySignature, dataSignature := sut.Sign(id, data)

				v := verifypkcs1v15.New(h, hashProvider)
				assert.True(t, v.VerifyIdentity(id, identitySignature, publicKey))
				assert.True(t, v.VerifyData(data, dataSignature, publicKey))
			},
		},
	}

	for i := range cases {
		t.Run(cases[i].name, cases[i].test)
	}
}

// TestSigner_SetUp tests signpkcs1v15.SetUp.
func TestSigner_SetUp(t *testing.T) {
	sut := newSUT(crypto.SHA256, testInternal.ValidPrivateKey, testInternal.ValidPublicKey, sha256.New())

	// for coverage; no assertion
	sut.SetUp()
}

// TestSigner_TearDown tests signpkcs1v15.TearDown.
func TestSigner_TearDown(t *testing.T) {
	sut := newSUT(crypto.SHA256, testInternal.ValidPrivateKey, testInternal.ValidPublicKey, sha256.New())

	// for coverage; no assertion
	sut.TearDown()
}

// TestSigner_PublicKey tests signpkcs1v15.PublicKey.
func TestSigner_PublicKey(t *testing.T) {
	publicKey := test.FactoryRandomByteSlice()
	sut := newSUT(crypto.SHA256, testInternal.ValidPrivateKey, publicKey, sha256.New())

	assert.Equal(t, publicKey, sut.PublicKey())
}

// TestSigner_Metadata tests signpkcs1v15.Metadata.
func TestSigner_Metadata(t *testing.T) {
	type testCase struct {
		name string
		test func(t *testing.T)
	}

	cases := []testCase{
		{
			name: "Success",
			test: func(t *testing.T) {
				signerHash := crypto.SHA256
				reducerHash := sha256.New()
				sut := newSUT(signerHash, testInternal.ValidPrivateKey, testInternal.ValidPublicKey, reducerHash)

				result := sut.Metadata()

				assert.Equal(t, pkcsSignerMetadata.NewSuccess(signerHash, reducerHash.Kind()), result)
			},
		},
		{
			name: "Failure",
			test: func(t *testing.T) {
				sut := newSUT(
					crypto.SHA256,
					testInternal.ValidPrivateKey,
					testInternal.ValidPublicKey,
					passthrough.New(test.FactoryRandomString()),
				)
				identitySignature, dataSignature := sut.Sign(
					test.FactoryRandomFixedLengthByteSlice(1024, test.AlphanumericCharset),
					test.FactoryRandomFixedLengthByteSlice(1024, test.AlphanumericCharset),
				)
				assert.Nil(t, identitySignature, dataSignature)

				result := sut.Metadata()

				assert.Equal(
					t,
					testInternal.Marshal(t, pkcsSignerMetadata.NewFailure("crypto/rsa: input must be hashed message")),
					testInternal.Marshal(t, result),
				)
			},
		},
	}

	for i := range cases {
		t.Run(cases[i].name, cases[i].test)
	}
}
