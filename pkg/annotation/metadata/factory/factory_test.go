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
	"testing"

	factoryStub "github.com/project-alvarium/go-sdk/pkg/annotation/metadata/factory/stub"
	metadataStub "github.com/project-alvarium/go-sdk/pkg/annotation/metadata/stub"
	"github.com/project-alvarium/go-sdk/pkg/test"

	"github.com/stretchr/testify/assert"
)

// newSUT returns a new system under test.
func newSUT(annotatorFactories []Contract) *instance {
	return New(annotatorFactories)
}

// TestInstance_Create tests instance.Create.
func TestInstance_Create(t *testing.T) {
	type testCase struct {
		name string
		test func(t *testing.T)
	}

	cases := []testCase{
		{
			name: "Unknown name",
			test: func(t *testing.T) {
				sut := newSUT([]Contract{})

				result := sut.Create(test.FactoryRandomString(), test.FactoryRandomByteSlice())

				assert.Nil(t, result)
			},
		},
		{
			name: "Known name",
			test: func(t *testing.T) {
				kind := test.FactoryRandomString()
				expected := metadataStub.New(test.FactoryRandomString(), test.FactoryRandomByteSlice())
				sut := newSUT([]Contract{factoryStub.New(kind, expected)})

				result := sut.Create(kind, test.FactoryRandomByteSlice())

				assert.NotNil(t, result)
				assert.Equal(t, expected, result)
			},
		},
	}

	for i := range cases {
		t.Run(cases[i].name, cases[i].test)
	}
}
