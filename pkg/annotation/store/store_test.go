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

package store

import (
	"testing"

	testInternal "github.com/project-alvarium/go-sdk/internal/pkg/test"
	envelope "github.com/project-alvarium/go-sdk/pkg/annotation/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/metadata"
	"github.com/project-alvarium/go-sdk/pkg/identity"
	identityHash "github.com/project-alvarium/go-sdk/pkg/identity/hash"
	"github.com/project-alvarium/go-sdk/pkg/status"
	"github.com/project-alvarium/go-sdk/pkg/store/memory"
	"github.com/project-alvarium/go-sdk/pkg/test"

	"github.com/stretchr/testify/assert"
)

// newSUT returns a new system under test.
func newSUT() *Persistence {
	return New(memory.New())
}

// TestStore_FindByIdentity tests store.FindByIdentity.
func TestStore_FindByIdentity(t *testing.T) {
	type testCase struct {
		name                string
		identity            identity.Contract
		preCondition        func(t *testing.T, sut *Persistence)
		expectedAnnotations []*envelope.Annotations
		expectedStatus      status.Value
	}

	cases := []testCase{
		{
			name:                "does not exist",
			identity:            identityHash.New(test.FactoryRandomByteSlice()),
			preCondition:        func(_ *testing.T, _ *Persistence) {},
			expectedAnnotations: []*envelope.Annotations{},
			expectedStatus:      status.NotFound,
		},
		func() testCase {
			id := identityHash.New(test.FactoryRandomByteSlice())
			m := envelope.New(
				test.FactoryRandomString(),
				id,
				nil,
				metadata.Kind,
				metadata.New(
					test.FactoryRandomString(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomString(),
					test.FactoryEmptyInterface(),
				),
			)
			return testCase{
				name:     "exists",
				identity: id,
				preCondition: func(t *testing.T, sut *Persistence) {
					assert.Equal(t, status.Success, sut.Create(id, m))
				},
				expectedAnnotations: []*envelope.Annotations{m},
				expectedStatus:      status.Success,
			}
		}(),
	}

	for i := range cases {
		t.Run(
			cases[i].name,
			func(t *testing.T) {
				sut := newSUT()
				cases[i].preCondition(t, sut)

				m, result := sut.FindByIdentity(cases[i].identity)

				assert.Equal(t, testInternal.Marshal(t, cases[i].expectedAnnotations), testInternal.Marshal(t, m))
				assert.Equal(t, cases[i].expectedStatus, result)
			},
		)
	}
}

// TestStore_Create tests store.Create.
func TestStore_Create(t *testing.T) {
	type testCase struct {
		name                string
		identity            identity.Contract
		m                   *envelope.Annotations
		postCondition       func(t *testing.T, sut *Persistence)
		expectedAnnotations []*envelope.Annotations
	}

	cases := []testCase{
		func() testCase {
			id := identityHash.New(test.FactoryRandomByteSlice())
			m := envelope.New(
				test.FactoryRandomString(),
				id,
				nil,
				metadata.Kind,
				metadata.New(
					test.FactoryRandomString(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomString(),
					test.FactoryEmptyInterface(),
				),
			)
			return testCase{
				name:                "create once",
				identity:            id,
				m:                   m,
				postCondition:       func(_ *testing.T, _ *Persistence) {},
				expectedAnnotations: []*envelope.Annotations{m},
			}
		}(),
		func() testCase {
			id := identityHash.New(test.FactoryRandomByteSlice())
			m1 := envelope.New(
				test.FactoryRandomString(),
				id,
				nil,
				metadata.Kind,
				metadata.New(
					test.FactoryRandomString(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomString(),
					test.FactoryEmptyInterface(),
				),
			)
			m2 := envelope.New(
				test.FactoryRandomString(),
				id,
				nil,
				metadata.Kind,
				metadata.New(
					test.FactoryRandomString(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomString(),
					test.FactoryEmptyInterface(),
				),
			)
			return testCase{
				name:     "create twice",
				identity: id,
				m:        m1,
				postCondition: func(t *testing.T, sut *Persistence) {
					assert.Equal(t, status.Exists, sut.Create(id, m2))
				},
				expectedAnnotations: []*envelope.Annotations{m1},
			}
		}(),
	}

	for i := range cases {
		t.Run(
			cases[i].name,
			func(t *testing.T) {
				sut := newSUT()

				result := sut.Create(cases[i].identity, cases[i].m)

				cases[i].postCondition(t, sut)
				savedModel, result := sut.FindByIdentity(cases[i].identity)
				assert.Equal(t, status.Success, result)
				assert.Equal(
					t,
					testInternal.Marshal(t, cases[i].expectedAnnotations),
					testInternal.Marshal(t, savedModel),
				)
			},
		)
	}
}

// TestStore_Append tests store.Append.
func TestStore_Append(t *testing.T) {
	type testCase struct {
		name                string
		identity            identity.Contract
		m                   *envelope.Annotations
		preCondition        func(t *testing.T, sut *Persistence)
		postCondition       func(t *testing.T, sut *Persistence)
		expectedAnnotations []*envelope.Annotations
	}

	cases := []testCase{
		func() testCase {
			id := identityHash.New(test.FactoryRandomByteSlice())
			m1 := envelope.New(
				test.FactoryRandomString(),
				id,
				nil,
				metadata.Kind,
				metadata.New(
					test.FactoryRandomString(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomString(),
					test.FactoryEmptyInterface(),
				),
			)
			m2 := envelope.New(
				test.FactoryRandomString(),
				id,
				nil,
				metadata.Kind,
				metadata.New(
					test.FactoryRandomString(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomString(),
					test.FactoryEmptyInterface(),
				),
			)
			return testCase{
				name:     "append once",
				identity: id,
				m:        m2,
				preCondition: func(t *testing.T, sut *Persistence) {
					assert.Equal(t, status.Success, sut.Create(id, m1))
				},
				postCondition:       func(_ *testing.T, _ *Persistence) {},
				expectedAnnotations: []*envelope.Annotations{m1, m2},
			}
		}(),
		func() testCase {
			id := identityHash.New(test.FactoryRandomByteSlice())
			m1 := envelope.New(
				test.FactoryRandomString(),
				id,
				nil,
				metadata.Kind,
				metadata.New(
					test.FactoryRandomString(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomString(),
					test.FactoryEmptyInterface(),
				),
			)
			m2 := envelope.New(
				test.FactoryRandomString(),
				id,
				nil,
				metadata.Kind,
				metadata.New(
					test.FactoryRandomString(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomString(),
					test.FactoryEmptyInterface(),
				),
			)
			m3 := envelope.New(
				test.FactoryRandomString(),
				id,
				nil,
				metadata.Kind,
				metadata.New(
					test.FactoryRandomString(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomString(),
					test.FactoryEmptyInterface(),
				),
			)
			return testCase{
				name:     "append twice",
				identity: id,
				m:        m2,
				preCondition: func(t *testing.T, sut *Persistence) {
					assert.Equal(t, status.Success, sut.Create(id, m1))
				},
				postCondition: func(t *testing.T, sut *Persistence) {
					assert.Equal(t, status.Success, sut.Append(id, m3))
				},
				expectedAnnotations: []*envelope.Annotations{m1, m2, m3},
			}
		}(),
		func() testCase {
			id := identityHash.New(test.FactoryRandomByteSlice())
			m1 := envelope.New(
				test.FactoryRandomString(),
				id,
				nil,
				metadata.Kind,
				metadata.New(
					test.FactoryRandomString(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomString(),
					test.FactoryEmptyInterface(),
				),
			)
			m2 := envelope.New(
				test.FactoryRandomString(),
				id,
				nil,
				metadata.Kind,
				metadata.New(
					test.FactoryRandomString(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomString(),
					test.FactoryEmptyInterface(),
				),
			)
			m3 := envelope.New(
				test.FactoryRandomString(),
				id,
				nil,
				metadata.Kind,
				metadata.New(
					test.FactoryRandomString(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomString(),
					test.FactoryEmptyInterface(),
				),
			)
			m4 := envelope.New(
				test.FactoryRandomString(),
				id,
				nil,
				metadata.Kind,
				metadata.New(
					test.FactoryRandomString(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomByteSlice(),
					test.FactoryRandomString(),
					test.FactoryEmptyInterface(),
				),
			)
			return testCase{
				name:     "append thrice",
				identity: id,
				m:        m2,
				preCondition: func(t *testing.T, sut *Persistence) {
					assert.Equal(t, status.Success, sut.Create(id, m1))
				},
				postCondition: func(t *testing.T, sut *Persistence) {
					assert.Equal(t, status.Success, sut.Append(id, m3))
					assert.Equal(t, status.Success, sut.Append(id, m4))
				},
				expectedAnnotations: []*envelope.Annotations{m1, m2, m3, m4},
			}
		}(),
	}

	for i := range cases {
		t.Run(
			cases[i].name,
			func(t *testing.T) {
				sut := newSUT()
				cases[i].preCondition(t, sut)

				result := sut.Append(cases[i].identity, cases[i].m)

				cases[i].postCondition(t, sut)
				savedMetadata, result := sut.FindByIdentity(cases[i].identity)
				assert.Equal(t, status.Success, result)
				assert.Equal(
					t,
					testInternal.Marshal(t, cases[i].expectedAnnotations),
					testInternal.Marshal(t, savedMetadata),
				)
			},
		)
	}
}
