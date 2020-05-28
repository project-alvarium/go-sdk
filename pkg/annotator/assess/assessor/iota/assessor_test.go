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
	annotationMetadata "github.com/project-alvarium/go-sdk/pkg/annotation/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotation/metadata/factory"
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/iota/client"
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/iota/client/stub"
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/iota/metadata"
	iotaAssessorMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/iota/metadata"
	publishMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/publish/metadata"
	publishMetadataFactory "github.com/project-alvarium/go-sdk/pkg/annotator/publish/metadata/factory"
	iotaPublisherMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/iota/metadata"
	iotaMetadataFactory "github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/iota/metadata/factory"
	ipfsPublisherMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/ipfs/metadata"
	"github.com/project-alvarium/go-sdk/pkg/hashprovider/sha256"
	identityProvider "github.com/project-alvarium/go-sdk/pkg/identityprovider/hash"
	"github.com/project-alvarium/go-sdk/pkg/test"

	"github.com/iotaledger/iota.go/transaction"
	"github.com/stretchr/testify/assert"
)

// newSUT returns a new system under test.
func newSUT(client client.Contract) *assessor {
	return New(client, publishMetadataFactory.New([]factory.Contract{iotaMetadataFactory.New()}))
}

// newDefaultSUT returns an initialized system under test.
func newDefaultSUT() *assessor {
	return newSUT(stub.New([]transaction.Transactions{{}}, nil))
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
		expectedResult func(sut *assessor) annotationMetadata.Contract
	}

	cases := []testCase{
		func() testCase {
			return testCase{
				name:        "empty annotations slice",
				annotations: []*annotation.Instance{},
				client:      stub.New([]transaction.Transactions{{}}, nil),
				expectedResult: func(sut *assessor) annotationMetadata.Contract {
					return sut.failureAnnotationMatch()
				},
			}
		}(),
		func() testCase {
			return testCase{
				name: "non-IOTA publisher annotation",
				annotations: []*annotation.Instance{
					annotation.New(
						test.FactoryRandomString(),
						identityProvider.New(sha256.New()).Derive(test.FactoryRandomByteSlice()),
						nil,
						publishMetadata.New(
							test.FactoryRandomString(),
							ipfsPublisherMetadata.NewSuccess(test.FactoryRandomString()),
						),
					),
				},
				client: stub.New([]transaction.Transactions{{}}, nil),
				expectedResult: func(sut *assessor) annotationMetadata.Contract {
					return sut.failureAnnotationMatch()
				},
			}
		}(),
		func() testCase {
			unique := test.FactoryRandomString()
			publisherMetadata := publishMetadata.New(
				test.FactoryRandomString(),
				iotaPublisherMetadata.NewFailure(test.FactoryRandomString()),
			)
			return testCase{
				name: "IOTA publisher failure annotation",
				annotations: []*annotation.Instance{
					annotation.New(
						unique,
						identityProvider.New(sha256.New()).Derive(test.FactoryRandomByteSlice()),
						nil,
						publisherMetadata,
					),
				},
				client: stub.New(
					[]transaction.Transactions{
						{
							internalTest.FactoryAnnotationTransaction(t, unique, publisherMetadata),
						},
					},
					nil,
				),
				expectedResult: func(sut *assessor) annotationMetadata.Contract {
					return iotaAssessorMetadata.NewSuccess(false, []string{unique})
				},
			}
		}(),
		func() testCase {
			unique := test.FactoryRandomString()
			return testCase{
				name: "IOTA publisher annotation not found in IOTA Tangle",
				annotations: []*annotation.Instance{
					annotation.New(
						unique,
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
				client: stub.New([]transaction.Transactions{{}}, nil),
				expectedResult: func(sut *assessor) annotationMetadata.Contract {
					return sut.failureTransactionQueryError()
				},
			}
		}(),
		func() testCase {
			unique := test.FactoryRandomString()
			publisherMetadata := iotaPublisherMetadata.NewSuccess(
				internalTest.FactoryRandomAddressTrytesString(),
				test.FactoryRandomString(),
				test.FactoryRandomString(),
			)
			return testCase{
				name: "multiple IOTA Tangle transactions found for a single IOTA publisher annotation",
				annotations: []*annotation.Instance{
					annotation.New(
						unique,
						identityProvider.New(sha256.New()).Derive(test.FactoryRandomByteSlice()),
						nil,
						publishMetadata.New(
							test.FactoryRandomString(),
							publisherMetadata,
						),
					),
				},
				client: stub.New(
					[]transaction.Transactions{
						{
							internalTest.FactoryAnnotationTransaction(t, unique, publisherMetadata),
							internalTest.FactoryAnnotationTransaction(t, unique, publisherMetadata),
						},
					},
					nil,
				),
				expectedResult: func(sut *assessor) annotationMetadata.Contract {
					return sut.failureTransactionQueryError()
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
				client: stub.New([]transaction.Transactions{{}}, nil),
				expectedResult: func(sut *assessor) annotationMetadata.Contract {
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
				client: stub.New([]transaction.Transactions{{}}, err),
				expectedResult: func(sut *assessor) annotationMetadata.Contract {
					return sut.failureFindTransactionError(err.Error())
				},
			}
		}(),
		func() testCase {
			unique := test.FactoryRandomString()
			publisherMetadata := iotaPublisherMetadata.NewSuccess(
				internalTest.FactoryRandomAddressTrytesString(),
				test.FactoryRandomString(),
				test.FactoryRandomString(),
			)
			return testCase{
				name: "valid assessment of a single IOTA publisher annotation",
				annotations: []*annotation.Instance{
					annotation.New(
						unique,
						identityProvider.New(sha256.New()).Derive(test.FactoryRandomByteSlice()),
						nil,
						publishMetadata.New(
							test.FactoryRandomString(),
							publisherMetadata,
						),
					),
				},
				client: stub.New(
					[]transaction.Transactions{
						{
							internalTest.FactoryAnnotationTransaction(t, unique, publisherMetadata),
						},
					},
					nil,
				),
				expectedResult: func(sut *assessor) annotationMetadata.Contract {
					return iotaAssessorMetadata.NewSuccess(true, []string{unique})
				},
			}
		}(),
		func() testCase {
			unique1 := test.FactoryRandomString()
			publisherMetadata1 := iotaPublisherMetadata.NewSuccess(
				internalTest.FactoryRandomAddressTrytesString(),
				test.FactoryRandomString(),
				test.FactoryRandomString(),
			)
			unique2 := test.FactoryRandomString()
			publisherMetadata2 := iotaPublisherMetadata.NewSuccess(
				internalTest.FactoryRandomAddressTrytesString(),
				test.FactoryRandomString(),
				test.FactoryRandomString(),
			)
			return testCase{
				name: "valid assessment of multiple IOTA publisher annotation",
				annotations: []*annotation.Instance{
					annotation.New(
						unique1,
						identityProvider.New(sha256.New()).Derive(test.FactoryRandomByteSlice()),
						nil,
						publishMetadata.New(
							test.FactoryRandomString(),
							publisherMetadata1,
						),
					),
					annotation.New(
						unique2,
						identityProvider.New(sha256.New()).Derive(test.FactoryRandomByteSlice()),
						nil,
						publishMetadata.New(
							test.FactoryRandomString(),
							publisherMetadata2,
						),
					),
				},
				client: stub.New(
					[]transaction.Transactions{
						{
							internalTest.FactoryAnnotationTransaction(t, unique1, publisherMetadata1),
						},
						{
							internalTest.FactoryAnnotationTransaction(t, unique2, publisherMetadata2),
						},
					},
					nil,
				),
				expectedResult: func(sut *assessor) annotationMetadata.Contract {
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

// TestAssessor_Failure tests assessor.Failure
func TestAssessor_Failure(t *testing.T) {
	sut := newDefaultSUT()

	err := test.FactoryRandomString()
	result := sut.Failure(err)

	assert.Equal(t, metadata.NewFailure(err), result)
}

func TestAssessor_Kind(t *testing.T) {
	sut := newDefaultSUT()

	result := sut.Kind()

	assert.Equal(t, metadata.Kind, result)
}
