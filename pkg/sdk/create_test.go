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

package sdk

import (
	"crypto"
	"testing"

	testInternal "github.com/project-alvarium/go-sdk/internal/pkg/test"
	testMetadata "github.com/project-alvarium/go-sdk/internal/pkg/test/metadata"
	envelope "github.com/project-alvarium/go-sdk/pkg/annotation/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotation/store"
	"github.com/project-alvarium/go-sdk/pkg/annotation/uniqueprovider/ulid"
	"github.com/project-alvarium/go-sdk/pkg/annotator"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki"
	pkiMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/pki/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signpkcs1v15"
	"github.com/project-alvarium/go-sdk/pkg/annotator/provenance"
	"github.com/project-alvarium/go-sdk/pkg/annotator/stub"
	"github.com/project-alvarium/go-sdk/pkg/hashprovider/sha256"
	identityProvider "github.com/project-alvarium/go-sdk/pkg/identityprovider/hash"
	"github.com/project-alvarium/go-sdk/pkg/status"
	"github.com/project-alvarium/go-sdk/pkg/store/memory"
	"github.com/project-alvarium/go-sdk/pkg/test"

	"github.com/stretchr/testify/assert"
)

// TestInstance_Create tests instance.Create.
func TestInstance_Create(t *testing.T) {
	type testCase struct {
		name           string
		provenance     provenance.Contract
		data           []byte
		annotator      annotator.Contract
		postCondition  func(t *testing.T, sut *instance)
		expectedStatus status.Value
	}

	cases := []testCase{
		func() testCase {
			prov := test.FactoryRandomString()
			data := test.FactoryRandomByteSlice()
			s := status.Success
			return testCase{
				name:       "Nil after close (stub)",
				provenance: prov,
				data:       data,
				annotator:  stub.NewWithResult(status.New(prov, s)),
				postCondition: func(t *testing.T, sut *instance) {
					sut.Close()
					assert.Nil(t, sut.Create(test.FactoryRandomByteSlice()))
				},
				expectedStatus: s,
			}
		}(),
		func() testCase {
			prov := test.FactoryRandomString()
			h := sha256.New()
			idProvider := identityProvider.New(h)
			data := test.FactoryRandomByteSlice()
			id := idProvider.Derive(data)
			publicKey := testInternal.ValidPublicKey
			persistence := store.New(memory.New())
			s := signpkcs1v15.New(crypto.SHA256, testInternal.ValidPrivateKey, publicKey, h)
			a := pki.New(prov, ulid.New(), idProvider, persistence, s)
			idSignature, dataSignature := s.Sign(id.Binary(), data)
			return testCase{
				name:       "Success (One)",
				provenance: prov,
				data:       data,
				annotator:  a,
				postCondition: func(t *testing.T, sut *instance) {
					testMetadata.Assert(
						t,
						[]*envelope.Annotations{
							envelope.New(
								test.FactoryRandomString(),
								id,
								nil,
								pkiMetadata.Kind,
								pkiMetadata.New(prov, idSignature, dataSignature, publicKey, s.Kind(), s.Metadata()),
							),
						},
						id,
						persistence,
					)
				},
				expectedStatus: status.Success,
			}
		}(),
		func() testCase {
			prov := test.FactoryRandomString()
			h := sha256.New()
			idProvider := identityProvider.New(h)
			data := test.FactoryRandomByteSlice()
			id := idProvider.Derive(data)
			publicKey := testInternal.ValidPublicKey
			persistence := store.New(memory.New())
			s := signpkcs1v15.New(crypto.SHA256, testInternal.ValidPrivateKey, publicKey, h)
			a := pki.New(prov, ulid.New(), idProvider, persistence, s)
			idSignature, dataSignature := s.Sign(id.Binary(), data)
			return testCase{
				name:       "Success (Two)",
				provenance: prov,
				data:       data,
				annotator:  a,
				postCondition: func(t *testing.T, sut *instance) {
					assert.Equal(
						t,
						testInternal.Marshal(t, []*status.Contract{status.New(prov, status.Exists)}),
						testInternal.Marshal(t, sut.Create(data)),
					)

					testMetadata.Assert(
						t,
						[]*envelope.Annotations{
							envelope.New(
								test.FactoryRandomString(),
								id,
								nil,
								pkiMetadata.Kind,
								pkiMetadata.New(prov, idSignature, dataSignature, publicKey, s.Kind(), s.Metadata()),
							),
						},
						id,
						persistence,
					)
				},
				expectedStatus: status.Success,
			}
		}(),
	}

	for i := range cases {
		t.Run(
			cases[i].name,
			func(t *testing.T) {
				sut := newSUT([]annotator.Contract{cases[i].annotator})

				result := sut.Create(cases[i].data)

				assert.Equal(
					t,
					testInternal.Marshal(t, []*status.Contract{status.New(cases[i].provenance, cases[i].expectedStatus)}),
					testInternal.Marshal(t, result),
				)
				cases[i].postCondition(t, sut)
				sut.Close()
			},
		)
	}
}
