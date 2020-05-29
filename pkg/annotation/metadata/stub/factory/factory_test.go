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
	metadataStub "github.com/project-alvarium/go-sdk/pkg/annotation/metadata/stub"
	"github.com/project-alvarium/go-sdk/pkg/test"

	"github.com/stretchr/testify/assert"
)

// newSUT returns a new system under test.
func newSUT(metadataStub *metadataStub.Instance) *instance {
	return New(metadataStub)
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
				sut := newSUT(metadataStub.NewNullObject())

				result := sut.Create(test.FactoryRandomString(), test.FactoryRandomByteSlice())

				assert.Nil(t, result)
			},
		},
		{
			name: "Known name",
			test: func(t *testing.T) {
				kind := test.FactoryRandomString()
				value := test.FactoryRandomByteSlice()
				stub := metadataStub.New(kind, value)
				marshalledStub := testInternal.Marshal(t, stub)
				sut := newSUT(stub)

				result := sut.Create(kind, json.RawMessage(marshalledStub))

				assert.NotNil(t, result)
				assert.Equal(t, marshalledStub, testInternal.Marshal(t, result))
			},
		},
	}

	for i := range cases {
		t.Run(cases[i].name, cases[i].test)
	}
}
