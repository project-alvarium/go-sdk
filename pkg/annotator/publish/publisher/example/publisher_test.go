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

package example

import (
	"io"
	"testing"

	"github.com/project-alvarium/go-sdk/pkg/annotation"
	"github.com/project-alvarium/go-sdk/pkg/annotation/metadata"
	metadataStub "github.com/project-alvarium/go-sdk/pkg/annotation/metadata/stub"
	examplePublisherMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/example/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/example/writer"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/example/writer/failwriter"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/example/writer/testwriter"
	"github.com/project-alvarium/go-sdk/pkg/hashprovider/sha256"
	identityProvider "github.com/project-alvarium/go-sdk/pkg/identityprovider/hash"
	"github.com/project-alvarium/go-sdk/pkg/test"

	"github.com/stretchr/testify/assert"
)

// newSUT returns a new system under test.
func newSUT(writer io.Writer) *publisher {
	return New(writer)
}

// TestNew tests publisher.New.
func TestNew(t *testing.T) {
	sut := New(testwriter.New())

	assert.NotNil(t, sut)
	assert.IsType(t, &publisher{}, sut)
}

// TestPublisher_SetUp tests publisher.SetUp.
func TestPublisher_SetUp(t *testing.T) {
	sut := New(testwriter.New())

	sut.SetUp()

	// nothing to assert; for coverage only
	sut.TearDown()
}

// TestPublisher_TearDown tests publisher.TearDown.
func TestPublisher_TearDown(t *testing.T) {
	sut := New(testwriter.New())
	sut.SetUp()

	sut.TearDown()

	// nothing to assert; for coverage only
}

// TestPublisher_Publish tests publisher.Publish.
func TestPublisher_Publish(t *testing.T) {
	type testCase struct {
		name           string
		writer         writer.Contract
		annotations    []*annotation.Instance
		expectedResult metadata.Contract
		expectedGet    []*annotation.Instance
	}

	cases := []testCase{
		func() testCase {
			idProvider := identityProvider.New(sha256.New())
			data := test.FactoryRandomByteSlice()
			id := idProvider.Derive(data)
			return testCase{
				name:   "writer.Write failure",
				writer: failwriter.New(),
				annotations: []*annotation.Instance{
					annotation.New(
						test.FactoryRandomString(),
						id,
						nil,
						metadataStub.New(test.FactoryRandomString(), test.FactoryRandomString()),
					),
				},
				expectedResult: examplePublisherMetadata.NewFailure(failwriter.WriteErrorMessage),
				expectedGet:    nil,
			}
		}(),
		func() testCase {
			idProvider := identityProvider.New(sha256.New())
			data := test.FactoryRandomByteSlice()
			id := idProvider.Derive(data)
			kind := test.FactoryRandomString()
			annotations := []*annotation.Instance{
				annotation.New(test.FactoryRandomString(), id, nil, metadataStub.New(kind, test.FactoryRandomString())),
			}
			return testCase{
				name:           "single annotation as string",
				writer:         testwriter.New(),
				annotations:    annotations,
				expectedResult: examplePublisherMetadata.NewSuccess(),
				expectedGet:    annotations,
			}
		}(),
		func() testCase {
			idProvider := identityProvider.New(sha256.New())
			data := test.FactoryRandomByteSlice()
			id := idProvider.Derive(data)
			kind := test.FactoryRandomString()
			annotations := []*annotation.Instance{
				annotation.New(
					test.FactoryRandomString(),
					id,
					nil,
					metadataStub.New(
						kind,
						struct {
							Name  string
							Value int
						}{
							Name:  test.FactoryRandomString(),
							Value: test.FactoryRandomInt(),
						},
					),
				),
			}
			return testCase{
				name:           "single annotation as structure",
				writer:         testwriter.New(),
				annotations:    annotations,
				expectedResult: examplePublisherMetadata.NewSuccess(),
				expectedGet:    annotations,
			}
		}(),
		func() testCase {
			idProvider := identityProvider.New(sha256.New())
			data := test.FactoryRandomByteSlice()
			id1 := idProvider.Derive(test.FactoryRandomByteSlice())
			id2 := idProvider.Derive(data)
			kind := test.FactoryRandomString()
			annotations := []*annotation.Instance{
				annotation.New(
					test.FactoryRandomString(),
					id1,
					nil,
					metadataStub.New(kind, test.FactoryRandomString()),
				),
				annotation.New(
					test.FactoryRandomString(),
					id2,
					id1,
					metadataStub.New(kind, test.FactoryRandomString()),
				),
			}
			return testCase{
				name:           "two annotations",
				writer:         testwriter.New(),
				annotations:    annotations,
				expectedResult: examplePublisherMetadata.NewSuccess(),
				expectedGet:    annotations,
			}
		}(),
	}

	for i := range cases {
		t.Run(
			cases[i].name,
			func(t *testing.T) {
				sut := newSUT(cases[i].writer)

				result := sut.Publish(cases[i].annotations)

				assert.Equal(t, cases[i].expectedResult, result)
				if cases[i].expectedGet == nil {
					assert.Nil(t, cases[i].writer.Get())
				} else {
					assert.Equal(t, sut.Format(cases[i].expectedGet), cases[i].writer.Get())
				}
			},
		)
	}
}

// TestPublisher_Kind tests publisher.Kind.
func TestPublisher_Kind(t *testing.T) {
	sut := newSUT(testwriter.New())

	result := sut.Kind()

	assert.Equal(t, examplePublisherMetadata.Kind, result)
}
