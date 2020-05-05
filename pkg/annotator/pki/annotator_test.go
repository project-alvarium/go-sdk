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

package pki

import (
	"crypto"
	"testing"

	testInternal "github.com/project-alvarium/go-sdk/internal/pkg/test"
	testMetadata "github.com/project-alvarium/go-sdk/internal/pkg/test/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotation"
	"github.com/project-alvarium/go-sdk/pkg/annotation/store"
	"github.com/project-alvarium/go-sdk/pkg/annotation/store/memory"
	"github.com/project-alvarium/go-sdk/pkg/annotation/uniqueprovider/ulid"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/fail"
	failMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/fail/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signpkcs1v15"
	"github.com/project-alvarium/go-sdk/pkg/annotator/provenance"
	"github.com/project-alvarium/go-sdk/pkg/hashprovider/sha256"
	"github.com/project-alvarium/go-sdk/pkg/identityprovider"
	identityProvider "github.com/project-alvarium/go-sdk/pkg/identityprovider/hash"
	"github.com/project-alvarium/go-sdk/pkg/status"
	"github.com/project-alvarium/go-sdk/pkg/test"

	"github.com/stretchr/testify/assert"
)

// newSUT returns a new system under test.
func newSUT(
	provenance provenance.Contract,
	identityProvider identityprovider.Contract,
	store store.Contract,
	signer signer.Contract) *annotator {

	return New(provenance, ulid.New(), identityProvider, store, signer)
}

// TestAnnotator_SetUp tests annotator.SetUp.
func TestAnnotator_SetUp(t *testing.T) {
	hashProvider := sha256.New()
	s := fail.New()
	sut := newSUT(test.FactoryRandomString(), identityProvider.New(hashProvider), memory.New(), s)

	sut.SetUp()

	assert.True(t, s.SetUpCalled)
}

// TestAnnotator_TearDown tests annotator.TearDown.
func TestAnnotator_TearDown(t *testing.T) {
	hashProvider := sha256.New()
	s := fail.New()
	sut := newSUT(test.FactoryRandomString(), identityProvider.New(hashProvider), memory.New(), s)

	sut.TearDown()

	assert.True(t, s.TearDownCalled)
}

// TestAnnotator_Create tests annotator.Create.
func TestAnnotator_Create(t *testing.T) {
	type testCase struct {
		name             string
		provenance       provenance.Contract
		store            store.Contract
		identityProvider identityprovider.Contract
		signer           signer.Contract
		data             []byte
		postCondition    func(t *testing.T, sut *annotator)
	}

	cases := []testCase{
		func() testCase {
			p := test.FactoryRandomString()
			persistence := memory.New()
			hashProvider := sha256.New()
			idProvider := identityProvider.New(hashProvider)
			s := signpkcs1v15.New(
				crypto.SHA256,
				testInternal.ValidPrivateKey,
				testInternal.ValidPublicKey,
				hashProvider,
			)
			data := test.FactoryRandomByteSlice()
			id := idProvider.Derive(data)
			return testCase{
				name:             "Success (one)",
				provenance:       p,
				store:            persistence,
				identityProvider: idProvider,
				signer:           s,
				data:             data,
				postCondition: func(t *testing.T, sut *annotator) {
					identitySignature, dataSignature := s.Sign(id.Binary(), data)

					testMetadata.Assert(
						t,
						[]*annotation.Instance{
							annotation.New(
								test.FactoryRandomString(),
								id,
								nil,
								metadata.New(
									p,
									identitySignature,
									dataSignature,
									testInternal.ValidPublicKey,
									s.Metadata(),
								),
							),
						},
						idProvider.Derive(data),
						persistence,
					)
				},
			}
		}(),
		func() testCase {
			p := test.FactoryRandomString()
			persistence := memory.New()
			hashProvider := sha256.New()
			idProvider := identityProvider.New(hashProvider)
			s := signpkcs1v15.New(
				crypto.SHA256,
				testInternal.ValidPrivateKey,
				testInternal.ValidPublicKey,
				hashProvider,
			)
			data := test.FactoryRandomByteSlice()
			id := idProvider.Derive(data)
			return testCase{
				name:             "Success (two)",
				provenance:       p,
				store:            persistence,
				identityProvider: idProvider,
				signer:           s,
				data:             data,
				postCondition: func(t *testing.T, sut *annotator) {
					assert.Equal(
						t,
						testInternal.Marshal(t, status.New(p, status.Exists)),
						testInternal.Marshal(t, sut.Create(data)),
					)

					identitySignature, dataSignature := s.Sign(id.Binary(), data)

					testMetadata.Assert(
						t,
						[]*annotation.Instance{
							annotation.New(
								test.FactoryRandomString(),
								id,
								nil,
								metadata.New(
									p,
									identitySignature,
									dataSignature,
									testInternal.ValidPublicKey,
									s.Metadata(),
								),
							),
						},
						idProvider.Derive(data),
						persistence,
					)
				},
			}
		}(),
		func() testCase {
			p := test.FactoryRandomString()
			persistence := memory.New()
			idProvider := identityProvider.New(sha256.New())
			s := fail.New()
			data := test.FactoryRandomByteSlice()
			id := idProvider.Derive(data)
			return testCase{
				name:             "Fail (signer)",
				provenance:       p,
				store:            persistence,
				identityProvider: idProvider,
				signer:           s,
				data:             data,
				postCondition: func(t *testing.T, sut *annotator) {
					assert.Equal(
						t,
						testInternal.Marshal(t, status.New(p, status.Exists)),
						testInternal.Marshal(t, sut.Create(data)),
					)

					identitySignature, dataSignature := s.Sign(id.Binary(), data)

					assert.Nil(t, identitySignature)
					assert.Nil(t, dataSignature)

					testMetadata.Assert(
						t,
						[]*annotation.Instance{
							annotation.New(
								test.FactoryRandomString(),
								id,
								nil,
								metadata.New(p, identitySignature, dataSignature, nil, failMetadata.New()),
							),
						},
						idProvider.Derive(data),
						persistence,
					)
				},
			}
		}(),
	}

	for i := range cases {
		t.Run(
			cases[i].name,
			func(t *testing.T) {
				sut := newSUT(cases[i].provenance, cases[i].identityProvider, cases[i].store, cases[i].signer)

				result := sut.Create(cases[i].data)

				assert.Equal(
					t,
					testInternal.Marshal(t, status.New(cases[i].provenance, status.Success)),
					testInternal.Marshal(t, result),
				)
				cases[i].postCondition(t, sut)
			},
		)
	}
}

// TestAnnotator_Mutate tests annotator.Mutate.
func TestAnnotator_Mutate(t *testing.T) {
	type testCase struct {
		name             string
		provenance       provenance.Contract
		identityProvider identityprovider.Contract
		store            store.Contract
		signer           signer.Contract
		oldData          []byte
		newData          []byte
		preCondition     func(t *testing.T, sut *annotator)
		postCondition    func(t *testing.T, sut *annotator)
	}

	cases := []testCase{
		func() testCase {
			p := test.FactoryRandomString()
			persistence := memory.New()
			hashProvider := sha256.New()
			idProvider := identityProvider.New(hashProvider)
			s := signpkcs1v15.New(
				crypto.SHA256,
				testInternal.ValidPrivateKey,
				testInternal.ValidPublicKey,
				hashProvider,
			)
			data := test.FactoryRandomByteSlice()
			id := idProvider.Derive(data)
			return testCase{
				name:             "Mutate Once Same",
				provenance:       p,
				identityProvider: idProvider,
				store:            persistence,
				signer:           s,
				oldData:          data,
				newData:          data,
				preCondition: func(t *testing.T, sut *annotator) {
					assert.Equal(
						t,
						testInternal.Marshal(t, status.New(p, status.Success)),
						testInternal.Marshal(t, sut.Create(data)),
					)
				},
				postCondition: func(t *testing.T, sut *annotator) {
					identitySignature, dataSignature := s.Sign(id.Binary(), data)

					testMetadata.Assert(
						t,
						[]*annotation.Instance{
							annotation.New(
								test.FactoryRandomString(),
								id,
								nil,
								metadata.New(
									p,
									identitySignature,
									dataSignature,
									testInternal.ValidPublicKey,
									s.Metadata(),
								),
							),
							annotation.New(
								test.FactoryRandomString(),
								id,
								id,
								metadata.New(
									p,
									identitySignature,
									dataSignature,
									testInternal.ValidPublicKey,
									s.Metadata(),
								),
							),
						},
						idProvider.Derive(data),
						persistence,
					)
				},
			}
		}(),
		func() testCase {
			p := test.FactoryRandomString()
			persistence := memory.New()
			hashProvider := sha256.New()
			idProvider := identityProvider.New(hashProvider)
			s := signpkcs1v15.New(
				crypto.SHA256,
				testInternal.ValidPrivateKey,
				testInternal.ValidPublicKey,
				hashProvider,
			)
			data1 := test.FactoryRandomByteSlice()
			id1 := idProvider.Derive(data1)
			data2 := test.FactoryRandomByteSlice()
			id2 := idProvider.Derive(data2)
			return testCase{
				name:             "Mutate Once Different",
				provenance:       p,
				identityProvider: idProvider,
				store:            persistence,
				signer:           s,
				oldData:          data1,
				newData:          data2,
				preCondition: func(t *testing.T, sut *annotator) {
					assert.Equal(
						t,
						testInternal.Marshal(t, status.New(p, status.Success)),
						testInternal.Marshal(t, sut.Create(data1)),
					)
				},
				postCondition: func(t *testing.T, sut *annotator) {
					identitySignature2, dataSignature2 := s.Sign(id2.Binary(), data2)
					identitySignature1, dataSignature1 := s.Sign(id1.Binary(), data1)

					testMetadata.Assert(
						t,
						[]*annotation.Instance{
							annotation.New(
								test.FactoryRandomString(),
								id2,
								id1,
								metadata.New(
									p,
									identitySignature2,
									dataSignature2,
									testInternal.ValidPublicKey,
									s.Metadata(),
								),
							),
							annotation.New(
								test.FactoryRandomString(),
								id1,
								nil,
								metadata.New(
									p,
									identitySignature1,
									dataSignature1,
									testInternal.ValidPublicKey,
									s.Metadata(),
								),
							),
						},
						idProvider.Derive(data2),
						persistence,
					)
				},
			}
		}(),
		func() testCase {
			p := test.FactoryRandomString()
			persistence := memory.New()
			hashProvider := sha256.New()
			idProvider := identityProvider.New(hashProvider)
			s := signpkcs1v15.New(
				crypto.SHA256,
				testInternal.ValidPrivateKey,
				testInternal.ValidPublicKey,
				hashProvider,
			)
			data1 := test.FactoryRandomByteSlice()
			id1 := idProvider.Derive(data1)
			data2 := test.FactoryRandomByteSlice()
			id2 := idProvider.Derive(data2)
			data3 := test.FactoryRandomByteSlice()
			id3 := idProvider.Derive(data3)
			return testCase{
				name:             "Mutate Twice",
				provenance:       p,
				identityProvider: idProvider,
				store:            persistence,
				signer:           s,
				oldData:          data1,
				newData:          data2,
				preCondition: func(t *testing.T, sut *annotator) {
					assert.Equal(
						t,
						testInternal.Marshal(t, status.New(p, status.Success)),
						testInternal.Marshal(t, sut.Create(data1)),
					)
				},
				postCondition: func(t *testing.T, sut *annotator) {
					assert.Equal(
						t,
						testInternal.Marshal(t, status.New(p, status.Success)),
						testInternal.Marshal(t, sut.Mutate(data2, data3)),
					)

					identitySignature3, dataSignature3 := s.Sign(id3.Binary(), data3)
					identitySignature2, dataSignature2 := s.Sign(id2.Binary(), data2)
					identitySignature1, dataSignature1 := s.Sign(id1.Binary(), data1)

					testMetadata.Assert(
						t,
						[]*annotation.Instance{
							annotation.New(
								test.FactoryRandomString(),
								id3,
								id2,
								metadata.New(
									p,
									identitySignature3,
									dataSignature3,
									testInternal.ValidPublicKey,
									s.Metadata(),
								),
							),
							annotation.New(
								test.FactoryRandomString(),
								id2,
								id1,
								metadata.New(
									p,
									identitySignature2,
									dataSignature2,
									testInternal.ValidPublicKey,
									s.Metadata(),
								),
							),
							annotation.New(
								test.FactoryRandomString(),
								id1,
								nil,
								metadata.New(
									p,
									identitySignature1,
									dataSignature1,
									testInternal.ValidPublicKey,
									s.Metadata(),
								),
							),
						},
						idProvider.Derive(data3),
						persistence,
					)
				},
			}
		}(),
		func() testCase {
			p := test.FactoryRandomString()
			persistence := memory.New()
			hashProvider := sha256.New()
			idProvider := identityProvider.New(hashProvider)
			s := fail.New()
			data := test.FactoryRandomByteSlice()
			id := idProvider.Derive(data)
			var expectedAnnotations []*annotation.Instance
			return testCase{
				name:             "Fail (signer)",
				provenance:       p,
				store:            persistence,
				identityProvider: idProvider,
				signer:           s,
				oldData:          data,
				newData:          data,
				preCondition: func(t *testing.T, sut *annotator) {
					// annotator that uses same store as SUT; used for call to Create().
					a := newSUT(
						p,
						idProvider,
						persistence,
						signpkcs1v15.New(
							crypto.SHA256,
							testInternal.ValidPrivateKey,
							testInternal.ValidPublicKey,
							hashProvider,
						),
					)

					assert.Equal(
						t,
						testInternal.Marshal(t, status.New(p, status.Success)),
						testInternal.Marshal(t, a.Create(data)),
					)

					var result status.Value
					expectedAnnotations, result = persistence.FindByIdentity(id)
					assert.Equal(t, status.Success, result)
				},
				postCondition: func(t *testing.T, sut *annotator) {
					assert.Equal(
						t,
						testInternal.Marshal(t, status.New(p, status.Exists)),
						testInternal.Marshal(t, sut.Create(data)),
					)

					identitySignature, dataSignature := s.Sign(id.Binary(), data)

					assert.Nil(t, identitySignature)
					assert.Nil(t, dataSignature)

					expectedAnnotations = append(
						expectedAnnotations,
						annotation.New(
							test.FactoryRandomString(),
							id,
							id,
							metadata.New(p, identitySignature, dataSignature, nil, failMetadata.New()),
						),
					)
					testMetadata.Assert(t, expectedAnnotations, idProvider.Derive(data), persistence)
				},
			}
		}(),
	}

	for i := range cases {
		t.Run(
			cases[i].name,
			func(t *testing.T) {
				sut := newSUT(cases[i].provenance, cases[i].identityProvider, cases[i].store, cases[i].signer)
				cases[i].preCondition(t, sut)

				result := sut.Mutate(cases[i].oldData, cases[i].newData)

				assert.Equal(
					t,
					testInternal.Marshal(t, status.New(cases[i].provenance, status.Success)),
					testInternal.Marshal(t, result),
				)
				cases[i].postCondition(t, sut)
			},
		)
	}
}
