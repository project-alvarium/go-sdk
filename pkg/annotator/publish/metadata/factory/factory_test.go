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

package factory

import (
	"encoding/json"
	"testing"

	testInternal "github.com/project-alvarium/go-sdk/internal/pkg/test"
	"github.com/project-alvarium/go-sdk/pkg/annotation/metadata/factory"
	publishMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/publish/metadata"
	examplePublisherMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/example/metadata"
	iotaPublisherMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/iota/metadata"
	ipfsPublisherMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/ipfs/metadata"
	"github.com/project-alvarium/go-sdk/pkg/test"

	"github.com/stretchr/testify/assert"
)

// newSUT returns a new system under test.
func newSUT(publisherFactories []factory.Contract) *instance {
	return New(publisherFactories)
}

// newDefaultSUT returns a new system under test.
func newDefaultSUT() *instance {
	return NewDefault()
}

// TestInstance_Create tests instance.Create.
func TestInstance_Create(t *testing.T) {
	type testCase struct {
		name string
		test func(t *testing.T)
	}

	cases := []testCase{
		{
			name: "Valid name, empty publisher factory slice",
			test: func(t *testing.T) {
				sut := newSUT([]factory.Contract{})

				result := sut.Create(ipfsPublisherMetadata.Kind, test.FactoryRandomByteSlice())

				assert.Nil(t, result)
			},
		},
		{
			name: "Unknown name",
			test: func(t *testing.T) {
				sut := newDefaultSUT()

				result := sut.Create(test.FactoryRandomString(), test.FactoryRandomByteSlice())

				assert.Nil(t, result)
			},
		},
		{
			name: "Valid (example publisher failure)",
			test: func(t *testing.T) {
				sut := newDefaultSUT()
				value := publishMetadata.New(
					test.FactoryRandomByteSlice(),
					examplePublisherMetadata.NewFailure(test.FactoryRandomString()),
				)

				result := sut.Create(publishMetadata.Kind, json.RawMessage(testInternal.Marshal(t, value)))

				assert.NotNil(t, result)
				assert.IsType(t, &publishMetadata.Instance{}, result)
				assert.Equal(t, testInternal.Marshal(t, value), testInternal.Marshal(t, result))
			},
		},
		{
			name: "Valid (example publisher success)",
			test: func(t *testing.T) {
				sut := newDefaultSUT()
				value := publishMetadata.New(
					test.FactoryRandomByteSlice(),
					examplePublisherMetadata.NewSuccess(),
				)

				result := sut.Create(publishMetadata.Kind, json.RawMessage(testInternal.Marshal(t, value)))

				assert.NotNil(t, result)
				assert.IsType(t, &publishMetadata.Instance{}, result)
				assert.Equal(t, testInternal.Marshal(t, value), testInternal.Marshal(t, result))
			},
		},
		{
			name: "Valid (iota publisher failure)",
			test: func(t *testing.T) {
				sut := newDefaultSUT()
				value := publishMetadata.New(
					test.FactoryRandomByteSlice(),
					iotaPublisherMetadata.NewFailure(test.FactoryRandomString()),
				)

				result := sut.Create(publishMetadata.Kind, json.RawMessage(testInternal.Marshal(t, value)))

				assert.NotNil(t, result)
				assert.IsType(t, &publishMetadata.Instance{}, result)
				assert.Equal(t, testInternal.Marshal(t, value), testInternal.Marshal(t, result))
			},
		},
		{
			name: "Valid (iota publisher success)",
			test: func(t *testing.T) {
				sut := newDefaultSUT()
				value := publishMetadata.New(
					test.FactoryRandomByteSlice(),
					iotaPublisherMetadata.NewSuccess(
						test.FactoryRandomString(),
						test.FactoryRandomString(),
						test.FactoryRandomString(),
					),
				)

				result := sut.Create(publishMetadata.Kind, json.RawMessage(testInternal.Marshal(t, value)))

				assert.NotNil(t, result)
				assert.IsType(t, &publishMetadata.Instance{}, result)
				assert.Equal(t, testInternal.Marshal(t, value), testInternal.Marshal(t, result))
			},
		},
		{
			name: "Valid (ipfs publisher failure)",
			test: func(t *testing.T) {
				sut := newDefaultSUT()
				value := publishMetadata.New(
					test.FactoryRandomByteSlice(),
					ipfsPublisherMetadata.NewFailure(test.FactoryRandomString()),
				)

				result := sut.Create(publishMetadata.Kind, json.RawMessage(testInternal.Marshal(t, value)))

				assert.NotNil(t, result)
				assert.IsType(t, &publishMetadata.Instance{}, result)
				assert.Equal(t, testInternal.Marshal(t, value), testInternal.Marshal(t, result))
			},
		},
		{
			name: "Valid (ipfs publisher success)",
			test: func(t *testing.T) {
				sut := newDefaultSUT()
				value := publishMetadata.New(
					test.FactoryRandomByteSlice(),
					ipfsPublisherMetadata.NewSuccess(test.FactoryRandomString()),
				)

				result := sut.Create(publishMetadata.Kind, json.RawMessage(testInternal.Marshal(t, value)))

				assert.NotNil(t, result)
				assert.IsType(t, &publishMetadata.Instance{}, result)
				assert.Equal(t, testInternal.Marshal(t, value), testInternal.Marshal(t, result))
			},
		},
	}

	for i := range cases {
		t.Run(cases[i].name, cases[i].test)
	}
}
