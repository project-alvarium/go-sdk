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
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/metadata"
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

// TestInstance_Mutate tests instance.Mutate.
func TestInstance_Mutate(t *testing.T) {
	type testCase struct {
		name           string
		provenance     provenance.Contract
		annotator      annotator.Contract
		oldData        []byte
		newData        []byte
		preCondition   func(t *testing.T, sut *instance)
		postCondition  func(t *testing.T, sut *instance)
		expectedStatus status.Value
	}

	cases := []testCase{
		func() testCase {
			prov := test.FactoryRandomString()
			data := test.FactoryRandomByteSlice()
			s := status.Success
			return testCase{
				name:         "Nil after close (stub)",
				provenance:   prov,
				annotator:    stub.NewWithResult(status.New(prov, s)),
				oldData:      data,
				newData:      data,
				preCondition: func(t *testing.T, sut *instance) {},
				postCondition: func(t *testing.T, sut *instance) {
					sut.Close()
					assert.Nil(t, sut.Mutate(data, test.FactoryRandomByteSlice()))
				},
				expectedStatus: s,
			}
		}(),
		func() testCase {
			prov := test.FactoryRandomString()
			persistence := store.New(memory.New())
			h := sha256.New()
			idProvider := identityProvider.New(h)
			data := test.FactoryRandomByteSlice()
			id := idProvider.Derive(data)
			publicKey := testInternal.ValidPublicKey
			s := signpkcs1v15.New(crypto.SHA256, testInternal.ValidPrivateKey, publicKey, h)
			idSignature, dataSignature := s.Sign(id.Binary(), data)
			return testCase{
				name:       "Mutate Once Same",
				provenance: prov,
				annotator:  pki.New(prov, ulid.New(), idProvider, persistence, s),
				oldData:    data,
				newData:    data,
				preCondition: func(t *testing.T, sut *instance) {
					assert.Equal(
						t,
						testInternal.Marshal(t, []*status.Contract{status.New(prov, status.Success)}),
						testInternal.Marshal(t, sut.Create(data)),
					)
				},
				postCondition: func(t *testing.T, sut *instance) {
					testMetadata.Assert(
						t,
						[]*envelope.Annotations{
							envelope.New(
								test.FactoryRandomString(),
								id,
								nil,
								metadata.Kind,
								metadata.New(prov, idSignature, dataSignature, publicKey, s.Kind(), s.Metadata()),
							),
							envelope.New(
								test.FactoryRandomString(),
								id,
								id,
								metadata.Kind,
								metadata.New(prov, idSignature, dataSignature, publicKey, s.Kind(), s.Metadata()),
							),
						},
						idProvider.Derive(data),
						persistence,
					)
				},
				expectedStatus: status.Success,
			}
		}(),
		func() testCase {
			prov := test.FactoryRandomString()
			persistence := store.New(memory.New())
			h := sha256.New()
			idProvider := identityProvider.New(h)
			data1 := test.FactoryRandomByteSlice()
			id1 := idProvider.Derive(data1)
			data2 := test.FactoryRandomByteSlice()
			id2 := idProvider.Derive(data2)
			publicKey := testInternal.ValidPublicKey
			s := signpkcs1v15.New(crypto.SHA256, testInternal.ValidPrivateKey, publicKey, h)
			idSignature1, dataSignature1 := s.Sign(id1.Binary(), data1)
			idSignature2, dataSignature2 := s.Sign(id2.Binary(), data2)
			return testCase{
				name:       "Mutate Once Different",
				provenance: prov,
				annotator:  pki.New(prov, ulid.New(), idProvider, persistence, s),
				oldData:    data1,
				newData:    data2,
				preCondition: func(t *testing.T, sut *instance) {
					assert.Equal(
						t,
						testInternal.Marshal(t, []*status.Contract{status.New(prov, status.Success)}),
						testInternal.Marshal(t, sut.Create(data1)),
					)
				},
				postCondition: func(t *testing.T, sut *instance) {
					testMetadata.Assert(
						t,
						[]*envelope.Annotations{
							envelope.New(
								test.FactoryRandomString(),
								id2,
								id1,
								metadata.Kind,
								metadata.New(prov, idSignature2, dataSignature2, publicKey, s.Kind(), s.Metadata()),
							),
							envelope.New(
								test.FactoryRandomString(),
								id1,
								nil,
								metadata.Kind,
								metadata.New(prov, idSignature1, dataSignature1, publicKey, s.Kind(), s.Metadata()),
							),
						},
						idProvider.Derive(data2),
						persistence,
					)

					testMetadata.Assert(
						t,
						[]*envelope.Annotations{
							envelope.New(
								test.FactoryRandomString(),
								id1,
								nil,
								metadata.Kind,
								metadata.New(prov, idSignature1, dataSignature1, publicKey, s.Kind(), s.Metadata()),
							),
						},
						idProvider.Derive(data1),
						persistence,
					)
				},
				expectedStatus: status.Success,
			}
		}(),
		func() testCase {
			prov := test.FactoryRandomString()
			persistence := store.New(memory.New())
			h := sha256.New()
			idProvider := identityProvider.New(h)
			data1 := test.FactoryRandomByteSlice()
			id1 := idProvider.Derive(data1)
			data2 := test.FactoryRandomByteSlice()
			id2 := idProvider.Derive(data2)
			data3 := test.FactoryRandomByteSlice()
			id3 := idProvider.Derive(data3)
			publicKey := testInternal.ValidPublicKey
			s := signpkcs1v15.New(crypto.SHA256, testInternal.ValidPrivateKey, publicKey, h)
			idSignature1, dataSignature1 := s.Sign(id1.Binary(), data1)
			idSignature2, dataSignature2 := s.Sign(id2.Binary(), data2)
			idSignature3, dataSignature3 := s.Sign(id3.Binary(), data3)
			return testCase{
				name:       "Mutate Twice",
				provenance: prov,
				annotator:  pki.New(prov, ulid.New(), idProvider, persistence, s),
				oldData:    data1,
				newData:    data2,
				preCondition: func(t *testing.T, sut *instance) {
					assert.Equal(
						t,
						testInternal.Marshal(t, []*status.Contract{status.New(prov, status.Success)}),
						testInternal.Marshal(t, sut.Create(data1)),
					)
				},
				postCondition: func(t *testing.T, sut *instance) {
					assert.Equal(
						t,
						testInternal.Marshal(t, []*status.Contract{status.New(prov, status.Success)}),
						testInternal.Marshal(t, sut.Mutate(data2, data3)),
					)

					testMetadata.Assert(
						t,
						[]*envelope.Annotations{
							envelope.New(
								test.FactoryRandomString(),
								id3,
								id2,
								metadata.Kind,
								metadata.New(prov, idSignature3, dataSignature3, publicKey, s.Kind(), s.Metadata()),
							),
							envelope.New(
								test.FactoryRandomString(),
								id2,
								id1,
								metadata.Kind,
								metadata.New(prov, idSignature2, dataSignature2, publicKey, s.Kind(), s.Metadata()),
							),
							envelope.New(
								test.FactoryRandomString(),
								id1,
								nil,
								metadata.Kind,
								metadata.New(prov, idSignature1, dataSignature1, publicKey, s.Kind(), s.Metadata()),
							),
						},
						idProvider.Derive(data3),
						persistence,
					)

					testMetadata.Assert(
						t,
						[]*envelope.Annotations{
							envelope.New(
								test.FactoryRandomString(),
								id2,
								id1,
								metadata.Kind,
								metadata.New(prov, idSignature2, dataSignature2, publicKey, s.Kind(), s.Metadata()),
							),
							envelope.New(
								test.FactoryRandomString(),
								id1,
								nil,
								metadata.Kind,
								metadata.New(prov, idSignature1, dataSignature1, publicKey, s.Kind(), s.Metadata()),
							),
						},
						idProvider.Derive(data2),
						persistence,
					)

					testMetadata.Assert(
						t,
						[]*envelope.Annotations{
							envelope.New(
								test.FactoryRandomString(),
								id1,
								nil,
								metadata.Kind,
								metadata.New(prov, idSignature1, dataSignature1, publicKey, s.Kind(), s.Metadata()),
							),
						},
						idProvider.Derive(data1),
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
				cases[i].preCondition(t, sut)

				result := sut.Mutate(cases[i].oldData, cases[i].newData)

				assert.Equal(
					t,
					testInternal.Marshal(
						t,
						[]*status.Contract{status.New(cases[i].provenance, cases[i].expectedStatus)},
					),
					testInternal.Marshal(t, result),
				)
				cases[i].postCondition(t, sut)
				sut.Close()
			},
		)
	}
}
