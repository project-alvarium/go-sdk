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
	"testing"

	testInternal "github.com/project-alvarium/go-sdk/internal/pkg/test"
	"github.com/project-alvarium/go-sdk/pkg/annotation"
	"github.com/project-alvarium/go-sdk/pkg/annotation/metadata"
	metadataStub "github.com/project-alvarium/go-sdk/pkg/annotation/metadata/stub"
	publisherMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/iota/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/iota/sdk"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/iota/sdk/iota"
	clientStub "github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/iota/sdk/iota/client/stub"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/iota/sdk/stub"
	"github.com/project-alvarium/go-sdk/pkg/hashprovider/sha256"
	identityProvider "github.com/project-alvarium/go-sdk/pkg/identityprovider/hash"
	"github.com/project-alvarium/go-sdk/pkg/test"

	"github.com/stretchr/testify/assert"
)

// newSUT returns a system under test.
func newSUT(sdk sdk.Contract) *publisher {
	return newWithIOTA(
		testInternal.FactoryRandomSeedString(),
		test.FactoryRandomUint64(),
		test.FactoryRandomUint64(),
		sdk,
	)
}

// TestNew tests publisher.New.
func TestNew(t *testing.T) {
	sut := New(
		testInternal.FactoryRandomSeedString(),
		test.FactoryRandomUint64(),
		test.FactoryRandomUint64(),
		clientStub.New(nil, nil, nil, nil),
	)

	assert.NotNil(t, sut)
	assert.IsType(t, &publisher{}, sut)
}

// TestPublisher_SetUp tests publisher.SetUp.
func TestPublisher_SetUp(t *testing.T) {
	sut := New(
		testInternal.FactoryRandomSeedString(),
		test.FactoryRandomUint64(),
		test.FactoryRandomUint64(),
		clientStub.New(nil, nil, nil, nil),
	)

	sut.SetUp()

	// nothing to assert; for coverage only
	sut.TearDown()
}

// TestPublisher_TearDown tests publisher.TearDown.
func TestPublisher_TearDown(t *testing.T) {
	sut := New(
		testInternal.FactoryRandomSeedString(),
		test.FactoryRandomUint64(),
		test.FactoryRandomUint64(),
		clientStub.New(nil, nil, nil, nil),
	)
	sut.SetUp()

	sut.TearDown()

	// nothing to assert; for coverage only
}

// TestPublisher_Publish tests publisher.Publish.
func TestPublisher_Publish(t *testing.T) {
	type testCase struct {
		name           string
		annotations    []*annotation.Instance
		sdk            sdk.Contract
		expectedResult func(sut *publisher) metadata.Contract
	}

	cases := []testCase{
		func() testCase {
			return testCase{
				name:        "no annotations",
				annotations: []*annotation.Instance{},
				sdk:         iota.New(Name, clientStub.New(nil, nil, nil, nil)),
				expectedResult: func(sut *publisher) metadata.Contract {
					return sut.failureNoAnnotations()
				},
			}
		}(),
		func() testCase {
			resultTx := testInternal.FactoryRandomFixedSizeBundle(1)[0]
			result := publisherMetadata.New(Name, resultTx.Address, resultTx.Hash, resultTx.Tag)
			return testCase{
				name: "single annotation",
				annotations: []*annotation.Instance{
					annotation.New(
						test.FactoryRandomString(),
						identityProvider.New(sha256.New()).Derive(test.FactoryRandomByteSlice()),
						nil,
						metadataStub.New(test.FactoryRandomString(), test.FactoryRandomString()),
					),
				},
				sdk: stub.New(result),
				expectedResult: func(sut *publisher) metadata.Contract {
					return result
				},
			}
		}(),
		func() testCase {
			resultTx := testInternal.FactoryRandomFixedSizeBundle(1)[0]
			result := publisherMetadata.New(Name, resultTx.Address, resultTx.Hash, resultTx.Tag)
			return testCase{
				name: "single annotation as structure",
				annotations: []*annotation.Instance{
					annotation.New(
						test.FactoryRandomString(),
						identityProvider.New(sha256.New()).Derive(test.FactoryRandomByteSlice()),
						nil,
						metadataStub.New(
							test.FactoryRandomString(),
							struct {
								Name  string
								Value int
							}{
								Name:  test.FactoryRandomString(),
								Value: test.FactoryRandomInt(),
							},
						),
					),
				},
				sdk: stub.New(result),
				expectedResult: func(sut *publisher) metadata.Contract {
					return result
				},
			}
		}(),
		func() testCase {
			resultTx := testInternal.FactoryRandomFixedSizeBundle(1)[0]
			result := publisherMetadata.New(Name, resultTx.Address, resultTx.Hash, resultTx.Tag)
			idProvider := identityProvider.New(sha256.New())
			data := test.FactoryRandomByteSlice()
			id1 := idProvider.Derive(test.FactoryRandomByteSlice())
			id2 := idProvider.Derive(data)
			kind := test.FactoryRandomString()
			return testCase{
				name: "two annotations",
				annotations: []*annotation.Instance{
					annotation.New(
						test.FactoryRandomString(),
						id2,
						id1,
						metadataStub.New(kind, test.FactoryRandomString()),
					),
					annotation.New(
						test.FactoryRandomString(),
						id1,
						nil,
						metadataStub.New(kind, test.FactoryRandomString()),
					),
				},
				sdk: stub.New(result),
				expectedResult: func(sut *publisher) metadata.Contract {
					return result
				},
			}
		}(),
	}

	for i := range cases {
		t.Run(cases[i].name, func(t *testing.T) {
			sut := newSUT(cases[i].sdk)

			result := sut.Publish(cases[i].annotations)

			assert.Equal(t, cases[i].expectedResult(sut), result)
		})
	}

}

// TestPublisher_Kind tests publisher.Kind.
func TestPublisher_Kind(t *testing.T) {
	sut := newSUT(stub.New(nil))

	result := sut.Kind()

	assert.Equal(t, Name, result)
}

// TestKind tests Kind.
func TestKind(t *testing.T) {
	assert.Equal(t, Name, Kind())
}
