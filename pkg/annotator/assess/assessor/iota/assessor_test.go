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

package iota

import (
	"errors"
	"testing"

	internalTest "github.com/project-alvarium/go-sdk/internal/pkg/test"
	"github.com/project-alvarium/go-sdk/pkg/annotation"
	"github.com/project-alvarium/go-sdk/pkg/annotation/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/iota/client"
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/iota/client/stub"
	iotaAssessorMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/iota/metadata"
	assessMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/assess/metadata"
	publishMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/publish/metadata"
	metadataFactory "github.com/project-alvarium/go-sdk/pkg/annotator/publish/metadata/factory"
	iotaPublisherMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/iota/metadata"
	"github.com/project-alvarium/go-sdk/pkg/hashprovider/sha256"
	identityProvider "github.com/project-alvarium/go-sdk/pkg/identityprovider/hash"

	"github.com/iotaledger/iota.go/transaction"
	"github.com/project-alvarium/go-sdk/pkg/test"
	"github.com/stretchr/testify/assert"
)

// newSUT returns a system under test.
func newSUT(client client.Contract) *assessor {
	return New(client, metadataFactory.NewDefault())
}

// newDefaultSUT returns an initialized system under test.
func newDefaultSUT() *assessor {
	return newSUT(stub.New([]transaction.Transactions{}, nil))
}

// TestAssessor_SetUp tests assessor.SetUp.
func TestAssessor_SetUp(t *testing.T) {
	sut := newDefaultSUT()

	sut.SetUp()
}

// TestAssessor_TearDown tests assessor.TearDown.
func TestAssessor_TearDown(t *testing.T) {
	sut := newDefaultSUT()

	sut.TearDown()

}

// TestAssessor_Assess tests assessor.Assess.
func TestAssessor_Assess(t *testing.T) {
	type testCase struct {
		name           string
		annotations    []*annotation.Instance
		client         client.Contract
		expectedResult func(sut *assessor) metadata.Contract
	}

	cases := []testCase{
		func() testCase {
			return testCase{
				name:        "empty annotations slice",
				annotations: []*annotation.Instance{},
				client:      stub.New([]transaction.Transactions{}, nil),
				expectedResult: func(sut *assessor) metadata.Contract {
					return sut.failureAnnotationMatch()
				},
			}
		}(),
		func() testCase {
			return testCase{
				name: "non annotation match (non-empty annotations slice)",
				annotations: []*annotation.Instance{
					annotation.New(
						test.FactoryRandomString(),
						identityProvider.New(sha256.New()).Derive(test.FactoryRandomByteSlice()),
						nil,
						assessMetadata.New(
							test.FactoryRandomString(),
							internalTest.FactoryNonIotaPublisherAnnotation(),
						),
					),
				},
				client: stub.New([]transaction.Transactions{}, nil),
				expectedResult: func(sut *assessor) metadata.Contract {
					return sut.failureAnnotationMatch()
				},
			}
		}(),
		func() testCase {
			return testCase{
				name: "IOTA publisher annotation address checksum error",
				annotations: []*annotation.Instance{
					annotation.New(
						test.FactoryRandomString(),
						identityProvider.New(sha256.New()).Derive(test.FactoryRandomByteSlice()),
						nil,
						publishMetadata.New(
							test.FactoryRandomString(),
							iotaPublisherMetadata.NewSuccess(
								internalTest.FactoryRandomInvalidAddressTrytesString(),
								test.FactoryRandomString(),
								test.FactoryRandomString(),
							),
						),
					),
				},
				client: stub.New([]transaction.Transactions{}, nil),
				expectedResult: func(sut *assessor) metadata.Contract {
					return sut.failureAddressChecksumError("invalid address")
				},
			}
		}(),
		func() testCase {
			err := errors.New(test.FactoryRandomString())
			return testCase{
				name: "find transaction objects error",
				annotations: []*annotation.Instance{
					annotation.New(
						test.FactoryRandomString(),
						identityProvider.New(sha256.New()).Derive(test.FactoryRandomByteSlice()),
						nil,
						publishMetadata.New(
							test.FactoryRandomString(),
							iotaPublisherMetadata.NewSuccess(
								internalTest.FactoryRandomAddressTrytesString(),
								test.FactoryRandomString(),
								test.FactoryRandomString(),
							),
						),
					),
				},
				client: stub.New([]transaction.Transactions{}, err),
				expectedResult: func(sut *assessor) metadata.Contract {
					return sut.failureFindTransactionError(err.Error())
				},
			}
		}(),
		func() testCase {
			unique := test.FactoryRandomString()
			m := iotaPublisherMetadata.NewSuccess(
				internalTest.FactoryRandomAddressTrytesString(),
				test.FactoryRandomString(),
				test.FactoryRandomString(),
			)
			return testCase{
				name: "multiple IOTA publisher annotations found for address",
				annotations: []*annotation.Instance{
					annotation.New(
						unique,
						identityProvider.New(sha256.New()).Derive(test.FactoryRandomByteSlice()),
						nil,
						publishMetadata.New(test.FactoryRandomString(), m),
					),
				},
				client: stub.New(
					[]transaction.Transactions{
						{
							internalTest.FactoryAnnotationTransaction(unique, m),
							internalTest.FactoryAnnotationTransaction(unique, m),
						},
					},
					nil,
				),
				expectedResult: func(sut *assessor) metadata.Contract {
					return sut.failureTransactionResultError()
				},
			}
		}(),
		func() testCase {
			unique := test.FactoryRandomString()
			m := iotaPublisherMetadata.NewFailure(test.FactoryRandomString())
			return testCase{
				name: "IOTA publisher failure annotation",
				annotations: []*annotation.Instance{
					annotation.New(
						unique,
						identityProvider.New(sha256.New()).Derive(test.FactoryRandomByteSlice()),
						nil,
						publishMetadata.New(test.FactoryRandomString(), m),
					),
				},
				client: stub.New(
					[]transaction.Transactions{
						{
							internalTest.FactoryAnnotationTransaction(unique, m),
						},
					},
					nil,
				),
				expectedResult: func(sut *assessor) metadata.Contract {
					return iotaAssessorMetadata.NewSuccess(false, []string{unique})
				},
			}
		}(),
		func() testCase {
			unique := test.FactoryRandomString()
			m := iotaPublisherMetadata.NewSuccess(
				internalTest.FactoryRandomAddressTrytesString(),
				test.FactoryRandomString(),
				test.FactoryRandomString(),
			)
			return testCase{
				name: "IOTA publisher valid single annotation",
				annotations: []*annotation.Instance{
					annotation.New(
						unique,
						identityProvider.New(sha256.New()).Derive(test.FactoryRandomByteSlice()),
						nil,
						publishMetadata.New(test.FactoryRandomString(), m),
					),
				},
				client: stub.New(
					[]transaction.Transactions{
						{
							internalTest.FactoryAnnotationTransaction(unique, m),
						},
					},
					nil,
				),
				expectedResult: func(sut *assessor) metadata.Contract {
					return iotaAssessorMetadata.NewSuccess(true, []string{unique})
				},
			}
		}(),
		func() testCase {
			unique1 := test.FactoryRandomString()
			m1 := iotaPublisherMetadata.NewSuccess(
				internalTest.FactoryRandomAddressTrytesString(),
				test.FactoryRandomString(),
				test.FactoryRandomString(),
			)
			unique2 := test.FactoryRandomString()
			m2 := iotaPublisherMetadata.NewSuccess(
				internalTest.FactoryRandomAddressTrytesString(),
				test.FactoryRandomString(),
				test.FactoryRandomString(),
			)
			id := identityProvider.New(sha256.New()).Derive(test.FactoryRandomByteSlice())
			return testCase{
				name: "IOTA publisher valid multiple annotations",
				annotations: []*annotation.Instance{
					annotation.New(unique1, id, nil, publishMetadata.New(test.FactoryRandomString(), m1)),
					annotation.New(
						unique2,
						identityProvider.New(sha256.New()).Derive(test.FactoryRandomByteSlice()),
						id,
						publishMetadata.New(test.FactoryRandomString(), m2),
					),
				},
				client: stub.New(
					[]transaction.Transactions{
						{
							internalTest.FactoryAnnotationTransaction(unique1, m1),
						},
						{
							internalTest.FactoryAnnotationTransaction(unique2, m2),
						},
					},
					nil,
				),
				expectedResult: func(sut *assessor) metadata.Contract {
					return iotaAssessorMetadata.NewSuccess(true, []string{unique1, unique2})
				},
			}
		}(),
	}

	for i := range cases {
		sut := newSUT(cases[i].client)

		result := sut.Assess(cases[i].annotations)

		assert.Equal(t, cases[i].expectedResult(sut), result)
	}
}

// TestAssessor_Failure tests assessor.Failure.
func TestAssessor_Failure(t *testing.T) {
	sut := newDefaultSUT()

	err := test.FactoryRandomString()
	result := sut.Failure(err)

	assert.Equal(t, iotaAssessorMetadata.NewFailure(err), result)
}

// TestAssessor_Kind tests assessor.Kind.
func TestAssessor_Kind(t *testing.T) {
	sut := newDefaultSUT()

	result := sut.Kind()

	assert.Equal(t, iotaAssessorMetadata.Kind, result)
}
