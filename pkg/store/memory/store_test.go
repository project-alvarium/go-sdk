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

package memory

import (
	"testing"

	testInternal "github.com/project-alvarium/go-sdk/internal/pkg/test"
	"github.com/project-alvarium/go-sdk/pkg/identity"
	identityHash "github.com/project-alvarium/go-sdk/pkg/identity/hash"
	"github.com/project-alvarium/go-sdk/pkg/status"
	"github.com/project-alvarium/go-sdk/pkg/test"

	"github.com/stretchr/testify/assert"
)

// newSUT returns a new system under test.
func newSUT() *Persistence {
	return New()
}

// TestStore_FindByIdentity tests store.FindByIdentity.
func TestStore_FindByIdentity(t *testing.T) {
	type testCase struct {
		name           string
		identity       identity.Contract
		preCondition   func(t *testing.T, sut *Persistence)
		expectedModels []interface {
		}
		expectedStatus status.Value
	}

	cases := []testCase{
		func() testCase {
			return testCase{
				name:           "does not exist",
				identity:       identityHash.New(test.FactoryRandomByteSlice()),
				preCondition:   func(_ *testing.T, _ *Persistence) {},
				expectedModels: ([]interface{})(nil),
				expectedStatus: status.NotFound,
			}
		}(),
		func() testCase {
			id := identityHash.New(test.FactoryRandomByteSlice())
			model := interface{}(test.FactoryRandomString())
			return testCase{
				name:     "exists",
				identity: id,
				preCondition: func(t *testing.T, sut *Persistence) {
					assert.Equal(t, status.Success, sut.Create(id, model))
				},
				expectedModels: []interface{}{model},
				expectedStatus: status.Success,
			}
		}(),
	}

	for i := range cases {
		t.Run(
			cases[i].name,
			func(t *testing.T) {
				sut := newSUT()
				cases[i].preCondition(t, sut)

				models, result := sut.FindByIdentity(cases[i].identity)

				assert.Equal(t, testInternal.Marshal(t, cases[i].expectedModels), testInternal.Marshal(t, models))
				assert.Equal(t, cases[i].expectedStatus, result)
			},
		)
	}
}

// TestStore_Create tests store.Create.
func TestStore_Create(t *testing.T) {
	type testCase struct {
		name     string
		identity identity.Contract
		model    interface {
		}
		postCondition  func(t *testing.T, sut *Persistence)
		expectedModels []interface {
		}
	}

	cases := []testCase{
		func() testCase {
			model := test.FactoryRandomString()
			return testCase{
				name:           "create once",
				identity:       identityHash.New(test.FactoryRandomByteSlice()),
				model:          model,
				postCondition:  func(_ *testing.T, _ *Persistence) {},
				expectedModels: []interface{}{model},
			}
		}(),
		func() testCase {
			id := identityHash.New(test.FactoryRandomByteSlice())
			model1 := test.FactoryRandomString()
			model2 := test.FactoryRandomString()
			return testCase{
				name:     "create twice",
				identity: id,
				model:    model1,
				postCondition: func(t *testing.T, sut *Persistence) {
					assert.Equal(t, status.Exists, sut.Create(id, model2))
				},
				expectedModels: []interface{}{model1},
			}
		}(),
	}

	for i := range cases {
		t.Run(
			cases[i].name,
			func(t *testing.T) {
				sut := newSUT()

				result := sut.Create(cases[i].identity, cases[i].model)

				cases[i].postCondition(t, sut)
				savedModel, result := sut.FindByIdentity(cases[i].identity)
				assert.Equal(t, status.Success, result)
				assert.Equal(t, testInternal.Marshal(t, cases[i].expectedModels), testInternal.Marshal(t, savedModel))
			},
		)
	}
}

// TestStore_Append tests store.Append.
func TestStore_Append(t *testing.T) {
	type testCase struct {
		name           string
		identity       identity.Contract
		model          interface{}
		preCondition   func(t *testing.T, sut *Persistence)
		postCondition  func(t *testing.T, sut *Persistence)
		expectedStatus status.Value
		expectedModels []interface{}
	}

	cases := []testCase{
		func() testCase {
			id := identityHash.New(test.FactoryRandomByteSlice())
			return testCase{
				name:           "append to non-existent",
				identity:       id,
				model:          test.FactoryRandomString(),
				preCondition:   func(t *testing.T, sut *Persistence) {},
				postCondition:  func(_ *testing.T, _ *Persistence) {},
				expectedStatus: status.NotFound,
				expectedModels: []interface{}(nil),
			}
		}(),
		func() testCase {
			id := identityHash.New(test.FactoryRandomByteSlice())
			model1 := test.FactoryRandomString()
			model2 := test.FactoryRandomString()
			return testCase{
				name:     "append once",
				identity: id,
				model:    model2,
				preCondition: func(t *testing.T, sut *Persistence) {
					assert.Equal(t, status.Success, sut.Create(id, model1))
				},
				postCondition:  func(_ *testing.T, _ *Persistence) {},
				expectedStatus: status.Success,
				expectedModels: []interface{}{model1, model2},
			}
		}(),
		func() testCase {
			id := identityHash.New(test.FactoryRandomByteSlice())
			model1 := test.FactoryRandomString()
			model2 := test.FactoryRandomString()
			model3 := test.FactoryRandomString()
			return testCase{
				name:     "append twice",
				identity: id,
				model:    model2,
				preCondition: func(t *testing.T, sut *Persistence) {
					assert.Equal(t, status.Success, sut.Create(id, model1))
				},
				postCondition: func(t *testing.T, sut *Persistence) {
					assert.Equal(t, status.Success, sut.Append(id, model3))
				},
				expectedStatus: status.Success,
				expectedModels: []interface{}{model1, model2, model3},
			}
		}(),
		func() testCase {
			id := identityHash.New(test.FactoryRandomByteSlice())
			model1 := test.FactoryRandomString()
			model2 := test.FactoryRandomString()
			model3 := test.FactoryRandomString()
			model4 := test.FactoryRandomString()
			return testCase{
				name:     "append thrice",
				identity: id,
				model:    model2,
				preCondition: func(t *testing.T, sut *Persistence) {
					assert.Equal(t, status.Success, sut.Create(id, model1))
				},
				postCondition: func(t *testing.T, sut *Persistence) {
					assert.Equal(t, status.Success, sut.Append(id, model3))
					assert.Equal(t, status.Success, sut.Append(id, model4))
				},
				expectedStatus: status.Success,
				expectedModels: []interface{}{model1, model2, model3, model4},
			}
		}(),
	}

	for i := range cases {
		t.Run(
			cases[i].name,
			func(t *testing.T) {
				sut := newSUT()
				cases[i].preCondition(t, sut)

				result := sut.Append(cases[i].identity, cases[i].model)

				cases[i].postCondition(t, sut)
				savedModel, result := sut.FindByIdentity(cases[i].identity)
				assert.Equal(t, testInternal.Marshal(t, cases[i].expectedStatus), testInternal.Marshal(t, result))
				assert.Equal(t, testInternal.Marshal(t, cases[i].expectedModels), testInternal.Marshal(t, savedModel))
			},
		)
	}
}
