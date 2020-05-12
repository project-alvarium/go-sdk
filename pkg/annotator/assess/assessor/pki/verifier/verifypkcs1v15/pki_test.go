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

package verifypkcs1v15

import (
	"crypto"
	"testing"

	testInternal "github.com/project-alvarium/go-sdk/internal/pkg/test"
	pkcsSigner "github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signpkcs1v15"
	"github.com/project-alvarium/go-sdk/pkg/hashprovider"
	"github.com/project-alvarium/go-sdk/pkg/hashprovider/sha256"
	identityProvider "github.com/project-alvarium/go-sdk/pkg/identityprovider/hash"
	"github.com/project-alvarium/go-sdk/pkg/test"

	"github.com/stretchr/testify/assert"
)

// newSUT returns a new system under test.
func newSUT(signerHash crypto.Hash, reducerHash hashprovider.Contract) *verifier {
	return New(signerHash, reducerHash)
}

// TestVerifier_VerifyIdentity tests verifypkcs1v15.VerifyIdentity
func TestVerifier_VerifyIdentity(t *testing.T) {
	type testCase struct {
		name        string
		signerHash  crypto.Hash
		reducerHash hashprovider.Contract
		identity    []byte
		signature   []byte
		publicKey   []byte
		expected    bool
	}

	cases := []testCase{
		func() testCase {
			signerHash := crypto.SHA256
			hashProvider := sha256.New()
			publicKey := testInternal.ValidPublicKey
			data := test.FactoryRandomByteSlice()
			id := identityProvider.New(hashProvider).Derive(data)
			s := pkcsSigner.New(signerHash, testInternal.ValidPrivateKey, publicKey, hashProvider)
			signature, _ := s.Sign(id.Binary(), data)
			return testCase{
				name:        "valid",
				signerHash:  signerHash,
				reducerHash: hashProvider,
				identity:    id.Binary(),
				signature:   signature,
				publicKey:   publicKey,
				expected:    true,
			}
		}(),
		func() testCase {
			return testCase{
				name:        "invalid (nil)",
				signerHash:  crypto.SHA256,
				reducerHash: sha256.New(),
				identity:    nil,
				signature:   nil,
				publicKey:   testInternal.ValidPublicKey,
				expected:    false,
			}
		}(),
		func() testCase {
			return testCase{
				name:        "invalid (not hex)",
				signerHash:  crypto.SHA256,
				reducerHash: sha256.New(),
				identity:    test.FactoryRandomByteSlice(),
				signature:   nil,
				publicKey:   testInternal.ValidPublicKey,
				expected:    false,
			}
		}(),
		func() testCase {
			return testCase{
				name:        "invalid (too long)",
				signerHash:  crypto.SHA256,
				reducerHash: sha256.New(),
				identity:    test.FactoryRandomFixedLengthByteSlice(1024, test.HexCharset),
				signature:   nil,
				publicKey:   testInternal.ValidPublicKey,
				expected:    false,
			}
		}(),
		func() testCase {
			hashProvider := sha256.New()
			data := test.FactoryRandomByteSlice()
			id := identityProvider.New(hashProvider).Derive(data)
			return testCase{
				name:        "invalid (nil signature key)",
				signerHash:  crypto.SHA256,
				reducerHash: hashProvider,
				identity:    id.Binary(),
				signature:   nil,
				publicKey:   testInternal.ValidPublicKey,
				expected:    false,
			}
		}(),
		func() testCase {
			signerHash := crypto.SHA256
			hashProvider := sha256.New()
			data := test.FactoryRandomByteSlice()
			id := identityProvider.New(hashProvider).Derive(data)
			s := pkcsSigner.New(signerHash, testInternal.ValidPrivateKey, testInternal.ValidPublicKey, hashProvider)
			signature, _ := s.Sign(id.Binary(), data)
			return testCase{
				name:        "invalid (nil public key)",
				signerHash:  signerHash,
				reducerHash: hashProvider,
				identity:    id.Binary(),
				signature:   signature,
				publicKey:   nil,
				expected:    false,
			}
		}(),
		func() testCase {
			hashProvider := sha256.New()
			data := test.FactoryRandomByteSlice()
			id := identityProvider.New(hashProvider).Derive(data)
			return testCase{
				name:        "invalid (invalid public key)",
				signerHash:  crypto.SHA256,
				reducerHash: hashProvider,
				identity:    id.Binary(),
				signature:   nil,
				publicKey:   test.FactoryRandomByteSlice(),
				expected:    false,
			}
		}(),
		func() testCase {
			hashProvider := sha256.New()
			data := test.FactoryRandomByteSlice()
			id := identityProvider.New(hashProvider).Derive(data)
			return testCase{
				name:        "invalid (invalid public key)",
				signerHash:  crypto.SHA256,
				reducerHash: hashProvider,
				identity:    id.Binary(),
				signature:   nil,
				publicKey:   testInternal.ValidPrivateKey,
				expected:    false,
			}
		}(),
		func() testCase {
			signerHash := crypto.SHA256
			hashProvider := sha256.New()
			data := test.FactoryRandomByteSlice()
			id := identityProvider.New(hashProvider).Derive(data)
			s := pkcsSigner.New(signerHash, testInternal.ValidPrivateKey, testInternal.ValidPublicKey, hashProvider)
			signature, _ := s.Sign(id.Binary(), data)
			return testCase{
				name:        "invalid (invalid public key)",
				signerHash:  signerHash,
				reducerHash: hashProvider,
				identity:    id.Binary(),
				signature:   signature,
				publicKey:   testInternal.InvalidPublicKey,
				expected:    false,
			}
		}(),
	}

	for i := range cases {
		t.Run(
			cases[i].name,
			func(t *testing.T) {
				sut := newSUT(cases[i].signerHash, cases[i].reducerHash)

				result := sut.VerifyIdentity(cases[i].identity, cases[i].signature, cases[i].publicKey)

				assert.Equal(t, cases[i].expected, result)
			},
		)
	}
}

// TestVerifier_VerifyData tests verifypkcs1v15.VerifyData
func TestVerifier_VerifyData(t *testing.T) {
	type testCase struct {
		name        string
		signerHash  crypto.Hash
		reducerHash hashprovider.Contract
		data        []byte
		signature   []byte
		publicKey   []byte
		expected    bool
	}

	cases := []testCase{
		func() testCase {
			signerHash := crypto.SHA256
			hashProvider := sha256.New()
			publicKey := testInternal.ValidPublicKey
			data := test.FactoryRandomByteSlice()
			s := pkcsSigner.New(signerHash, testInternal.ValidPrivateKey, publicKey, hashProvider)
			_, signature := s.Sign(hashProvider.Derive(data), data)
			return testCase{
				name:        "valid",
				signerHash:  signerHash,
				reducerHash: hashProvider,
				data:        data,
				signature:   signature,
				publicKey:   publicKey,
				expected:    true,
			}
		}(),
		func() testCase {
			return testCase{
				name:        "invalid (nil)",
				signerHash:  crypto.SHA256,
				reducerHash: sha256.New(),
				data:        nil,
				signature:   nil,
				publicKey:   testInternal.ValidPublicKey,
				expected:    false,
			}
		}(),
		func() testCase {
			return testCase{
				name:        "invalid (not hex)",
				signerHash:  crypto.SHA256,
				reducerHash: sha256.New(),
				data:        test.FactoryRandomByteSlice(),
				signature:   nil,
				publicKey:   testInternal.ValidPublicKey,
				expected:    false,
			}
		}(),
		func() testCase {
			return testCase{
				name:        "invalid (too long)",
				signerHash:  crypto.SHA256,
				reducerHash: sha256.New(),
				data:        []byte(test.FactoryRandomFixedLengthString(1024, test.HexCharset)),
				signature:   nil,
				publicKey:   testInternal.ValidPublicKey,
				expected:    false,
			}
		}(),
		func() testCase {
			return testCase{
				name:        "invalid (nil signature key)",
				signerHash:  crypto.SHA256,
				reducerHash: sha256.New(),
				data:        test.FactoryRandomByteSlice(),
				signature:   nil,
				publicKey:   testInternal.ValidPublicKey,
				expected:    false,
			}
		}(),
		func() testCase {
			signerHash := crypto.SHA256
			hashProvider := sha256.New()
			data := test.FactoryRandomByteSlice()
			s := pkcsSigner.New(signerHash, testInternal.ValidPrivateKey, testInternal.ValidPublicKey, hashProvider)
			_, signature := s.Sign(hashProvider.Derive(data), data)
			return testCase{
				name:        "invalid (nil public key)",
				signerHash:  signerHash,
				reducerHash: hashProvider,
				data:        test.FactoryRandomByteSlice(),
				signature:   signature,
				publicKey:   nil,
				expected:    false,
			}
		}(),
		func() testCase {
			return testCase{
				name:        "invalid (invalid public key)",
				signerHash:  crypto.SHA256,
				reducerHash: sha256.New(),
				data:        test.FactoryRandomByteSlice(),
				signature:   nil,
				publicKey:   test.FactoryRandomByteSlice(),
				expected:    false,
			}
		}(),
		func() testCase {
			return testCase{
				name:        "invalid (invalid public key)",
				signerHash:  crypto.SHA256,
				reducerHash: sha256.New(),
				data:        test.FactoryRandomByteSlice(),
				signature:   nil,
				publicKey:   testInternal.ValidPrivateKey,
				expected:    false,
			}
		}(),
		func() testCase {
			signerHash := crypto.SHA256
			hashProvider := sha256.New()
			data := test.FactoryRandomByteSlice()
			s := pkcsSigner.New(signerHash, testInternal.ValidPrivateKey, testInternal.ValidPublicKey, hashProvider)
			signature, _ := s.Sign(hashProvider.Derive(data), data)
			return testCase{
				name:        "invalid (invalid public key)",
				signerHash:  signerHash,
				reducerHash: hashProvider,
				data:        data,
				signature:   signature,
				publicKey:   testInternal.InvalidPublicKey,
				expected:    false,
			}
		}(),
	}

	for i := range cases {
		t.Run(
			cases[i].name,
			func(t *testing.T) {
				sut := newSUT(cases[i].signerHash, cases[i].reducerHash)

				result := sut.VerifyData(cases[i].data, cases[i].signature, cases[i].publicKey)

				assert.Equal(t, cases[i].expected, result)
			},
		)
	}
}
