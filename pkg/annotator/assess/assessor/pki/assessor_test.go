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
	"github.com/project-alvarium/go-sdk/pkg/annotation"
	"github.com/project-alvarium/go-sdk/pkg/annotation/store"
	"github.com/project-alvarium/go-sdk/pkg/annotation/uniqueprovider/ulid"
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/pki/factory"
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/pki/factory/fail"
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/pki/factory/verifier"
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/pki/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki"
	signer "github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signpkcs1v15"
	"github.com/project-alvarium/go-sdk/pkg/hashprovider/sha256"
	identityProvider "github.com/project-alvarium/go-sdk/pkg/identityprovider/hash"
	"github.com/project-alvarium/go-sdk/pkg/status"
	"github.com/project-alvarium/go-sdk/pkg/store/memory"
	"github.com/project-alvarium/go-sdk/pkg/test"

	"github.com/stretchr/testify/assert"
)

// newSUT returns a new system under test.
func newSUT(factory factory.Contract) *assessor {
	return New(factory)
}

// TestAssessor_SetUp tests verifier.SetUp.
func TestAssessor_SetUp(t *testing.T) {
	sut := newSUT(verifier.New())

	// no assertions; called for coverage.
	sut.SetUp()
}

// TestAssessor_TearDown tests verifier.TearDown.
func TestAssessor_TearDown(t *testing.T) {
	sut := newSUT(verifier.New())

	// no assertions; called for coverage.
	sut.TearDown()
}

// TestAssessor_Assess tests verifier.Assess.
func TestAssessor_Assess(t *testing.T) {
	type testCase struct {
		name               string
		h                  crypto.Hash
		verifier           factory.Contract
		preCondition       func(t *testing.T, sut *assessor) []*annotation.Instance
		expectedAssessment func() *metadata.Assessment
	}

	cases := []testCase{
		func() testCase {
			p := test.FactoryRandomString()
			h := crypto.SHA256
			hashProvider := sha256.New()
			idProvider := identityProvider.New(hashProvider)
			s := store.New(memory.New())
			a := pki.New(
				p,
				ulid.New(),
				idProvider,
				s,
				signer.New(h, testInternal.ValidPrivateKey, testInternal.ValidPublicKey, hashProvider),
			)
			data := test.FactoryRandomByteSlice()
			var annotations []*annotation.Instance
			return testCase{
				name:     "Create only",
				h:        h,
				verifier: verifier.New(),
				preCondition: func(t *testing.T, sut *assessor) []*annotation.Instance {
					assert.Equal(
						t,
						testInternal.Marshal(t, status.New(p, status.Success)),
						testInternal.Marshal(t, a.Create(data)),
					)

					var result status.Value
					annotations, result = s.FindByIdentity(idProvider.Derive(data))
					assert.Equal(t, status.Success, result)
					assert.NotNil(t, annotations)
					return annotations
				},
				expectedAssessment: func() *metadata.Assessment {
					var uniques []string
					for i := range annotations {
						uniques = append(uniques, annotations[i].Unique)
					}
					return metadata.New(true, uniques)
				},
			}
		}(),
		func() testCase {
			p := test.FactoryRandomString()
			h := crypto.SHA256
			hashProvider := sha256.New()
			idProvider := identityProvider.New(hashProvider)
			s := store.New(memory.New())
			a := pki.New(
				p,
				ulid.New(),
				idProvider,
				s,
				signer.New(h, testInternal.ValidPrivateKey, testInternal.ValidPublicKey, hashProvider),
			)
			data1 := test.FactoryRandomByteSlice()
			data2 := test.FactoryRandomByteSlice()
			var annotations []*annotation.Instance
			return testCase{
				name:     "Create and Mutate Once",
				h:        h,
				verifier: verifier.New(),
				preCondition: func(t *testing.T, sut *assessor) []*annotation.Instance {
					assert.Equal(
						t,
						testInternal.Marshal(t, status.New(p, status.Success)),
						testInternal.Marshal(t, a.Create(data1)),
					)
					assert.Equal(
						t,
						testInternal.Marshal(t, status.New(p, status.Success)),
						testInternal.Marshal(t, a.Mutate(data1, data2)),
					)

					var result status.Value
					annotations, result = s.FindByIdentity(idProvider.Derive(data2))
					assert.Equal(t, status.Success, result)
					assert.NotNil(t, annotations)
					return annotations
				},
				expectedAssessment: func() *metadata.Assessment {
					var uniques []string
					for i := range annotations {
						uniques = append(uniques, annotations[i].Unique)
					}
					return metadata.New(true, uniques)
				},
			}
		}(),
		func() testCase {
			h := crypto.SHA256
			idProvider := identityProvider.New(sha256.New())
			s := store.New(memory.New())
			data := test.FactoryRandomByteSlice()
			id := idProvider.Derive(data)
			return testCase{
				name:     "Ignore other type",
				h:        h,
				verifier: fail.New(),
				preCondition: func(t *testing.T, sut *assessor) []*annotation.Instance {
					a := annotation.New(ulid.New().Get(), id, nil, "otherType", nil)
					assert.Equal(t, status.Success, s.Create(id, a))

					annotations, result := s.FindByIdentity(id)
					assert.Equal(t, status.Success, result)
					assert.NotNil(t, annotations)
					return annotations
				},
				expectedAssessment: func() *metadata.Assessment {
					return metadata.New(true, []string{})
				},
			}
		}(),
		func() testCase {
			h := crypto.SHA256
			hashProvider := sha256.New()
			idProvider := identityProvider.New(hashProvider)
			s := store.New(memory.New())
			data := test.FactoryRandomByteSlice()
			var annotations []*annotation.Instance
			return testCase{
				name:     "Fail (verifier)",
				h:        h,
				verifier: fail.New(),
				preCondition: func(t *testing.T, sut *assessor) []*annotation.Instance {
					p := test.FactoryRandomString()
					a := pki.New(
						p,
						ulid.New(),
						idProvider,
						s,
						signer.New(h, testInternal.ValidPrivateKey, testInternal.ValidPublicKey, hashProvider),
					)

					assert.Equal(
						t,
						testInternal.Marshal(t, status.New(p, status.Success)),
						testInternal.Marshal(t, a.Create(data)),
					)

					var result status.Value
					annotations, result = s.FindByIdentity(idProvider.Derive(data))
					assert.Equal(t, status.Success, result)
					assert.NotNil(t, annotations)
					return annotations
				},
				expectedAssessment: func() *metadata.Assessment {
					var uniques []string
					for i := range annotations {
						uniques = append(uniques, annotations[i].Unique)
					}
					return metadata.New(false, uniques)
				},
			}
		}(),
	}

	for i := range cases {
		t.Run(
			cases[i].name,
			func(t *testing.T) {
				sut := newSUT(cases[i].verifier)
				annotations := cases[i].preCondition(t, sut)

				s := sut.Assess(annotations)

				assert.Equal(t, testInternal.Marshal(t, cases[i].expectedAssessment()), testInternal.Marshal(t, s))
			},
		)
	}
}

// TestAssessor_Kind tests verifier.Kind.
func TestAssessor_Kind(t *testing.T) {
	sut := newSUT(verifier.New())
	assert.Equal(t, name, sut.Kind())
}
