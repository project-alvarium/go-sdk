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

	testInternal "github.com/project-alvarium/go-sdk/internal/pkg/test"
	"github.com/project-alvarium/go-sdk/pkg/identity/hash"
	"github.com/project-alvarium/go-sdk/pkg/test"

	"github.com/stretchr/testify/assert"
)

// newSUT returns a new system under test.
func newSUT() *instance {
	return New()
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
				sut := newSUT()

				result := sut.Create(test.FactoryRandomString(), test.FactoryRandomByteSlice())

				assert.Nil(t, result)
			},
		},
		{
			name: "Valid (hash)",
			test: func(t *testing.T) {
				sut := newSUT()
				value := test.FactoryRandomByteSlice()

				result := sut.Create(hash.Kind, []byte(testInternal.Marshal(t, hash.New(value))))

				assert.NotNil(t, result)
				assert.IsType(t, &hash.Identity{}, result)
				assert.Equal(t, value, result.Binary())
			},
		},
	}

	for i := range cases {
		t.Run(cases[i].name, cases[i].test)
	}
}
