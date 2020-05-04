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

package ipfs

import (
	"errors"
	"testing"

	envelope "github.com/project-alvarium/go-sdk/pkg/annotation/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/published"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/ipfs/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/ipfs/sdk"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/ipfs/sdk/stub"
	"github.com/project-alvarium/go-sdk/pkg/hashprovider/sha256"
	identityProvider "github.com/project-alvarium/go-sdk/pkg/identityprovider/hash"
	"github.com/project-alvarium/go-sdk/pkg/test"
	"github.com/stretchr/testify/assert"
)

// newSUT returns a new system under test.
func newSUT(url string, sdk sdk.Contract) *publisher {
	return newWithIPFS(url, sdk)
}

// TestNew tests publisher.New.
func TestNew(t *testing.T) {
	sut := New(test.FactoryRandomString())

	assert.NotNil(t, sut)
	assert.IsType(t, &publisher{}, sut)
}

// TestPublisher_SetUp tests publisher.SetUp.
func TestPublisher_SetUp(t *testing.T) {
	sut := New(test.FactoryRandomString())

	sut.SetUp()

	// nothing to assert; for coverage only
	sut.TearDown()
}

// TestPublisher_TearDown tests publisher.TearDown.
func TestPublisher_TearDown(t *testing.T) {
	sut := New(test.FactoryRandomString())
	sut.SetUp()

	sut.TearDown()

	// nothing to assert; for coverage only
}

// TestPublisher_Publish tests publisher.Publish.
func TestPublisher_Publish(t *testing.T) {
	type testCase struct {
		name           string
		sdk            sdk.Contract
		annotations    []*envelope.Annotations
		expectedResult func(sut *publisher) published.Contract
	}

	cases := []testCase{
		func() testCase {
			message := test.FactoryRandomString()
			return testCase{
				name: "add failure",
				sdk:  stub.New(test.FactoryRandomString(), errors.New(message)),
				annotations: []*envelope.Annotations{
					envelope.New(
						test.FactoryRandomString(),
						identityProvider.New(sha256.New()).Derive(test.FactoryRandomByteSlice()),
						nil,
						test.FactoryRandomString(),
						test.FactoryRandomString(),
					),
				},
				expectedResult: func(sut *publisher) published.Contract {
					return sut.failureAdd(message)
				},
			}
		}(),
		func() testCase {
			cid := test.FactoryRandomString()
			return testCase{
				name: "single annotation as string",
				sdk:  stub.New(cid, nil),
				annotations: []*envelope.Annotations{
					envelope.New(
						test.FactoryRandomString(),
						identityProvider.New(sha256.New()).Derive(test.FactoryRandomByteSlice()),
						nil,
						test.FactoryRandomString(),
						test.FactoryRandomString(),
					),
				},
				expectedResult: func(sut *publisher) published.Contract {
					return metadata.New(cid)
				},
			}
		}(),
		func() testCase {
			cid := test.FactoryRandomString()
			return testCase{
				name: "single annotation as structure",
				sdk:  stub.New(cid, nil),
				annotations: []*envelope.Annotations{
					envelope.New(
						test.FactoryRandomString(),
						identityProvider.New(sha256.New()).Derive(test.FactoryRandomByteSlice()),
						nil,
						test.FactoryRandomString(),
						struct {
							Name  string
							Value int
						}{
							Name:  test.FactoryRandomString(),
							Value: test.FactoryRandomInt(),
						},
					),
				},
				expectedResult: func(sut *publisher) published.Contract {
					return metadata.New(cid)
				},
			}
		}(),
		func() testCase {
			cid := test.FactoryRandomString()
			idProvider := identityProvider.New(sha256.New())
			data := test.FactoryRandomByteSlice()
			id1 := idProvider.Derive(test.FactoryRandomByteSlice())
			id2 := idProvider.Derive(data)
			kind := test.FactoryRandomString()
			return testCase{
				name: "two annotations",
				sdk:  stub.New(cid, nil),
				annotations: []*envelope.Annotations{
					envelope.New(test.FactoryRandomString(), id2, id1, kind, test.FactoryRandomString()),
					envelope.New(test.FactoryRandomString(), id1, nil, kind, test.FactoryRandomString()),
				},
				expectedResult: func(sut *publisher) published.Contract {
					return metadata.New(cid)
				},
			}
		}(),
	}

	for i := range cases {
		t.Run(
			cases[i].name,
			func(t *testing.T) {
				url := test.FactoryRandomString()
				sut := newSUT(url, cases[i].sdk)

				result := sut.Publish(cases[i].annotations)

				assert.Equal(t, cases[i].expectedResult(sut), result)
			},
		)
	}
}

// TestPublisher_Kind tests publisher.Kind.
func TestPublisher_Kind(t *testing.T) {
	sut := newSUT(test.FactoryRandomString(), stub.New(test.FactoryRandomString(), nil))

	result := sut.Kind()

	assert.Equal(t, name, result)
}

// TestKind tests Kind.
func TestKind(t *testing.T) {
	assert.Equal(t, name, Kind())
}
