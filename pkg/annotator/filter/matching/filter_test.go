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

package matching

import (
	"testing"

	"github.com/project-alvarium/go-sdk/pkg/annotation"
	"github.com/project-alvarium/go-sdk/pkg/annotation/uniqueprovider/ulid"
	"github.com/project-alvarium/go-sdk/pkg/hashprovider/sha256"
	identityProvider "github.com/project-alvarium/go-sdk/pkg/identityprovider/hash"
	"github.com/project-alvarium/go-sdk/pkg/test"

	"github.com/stretchr/testify/assert"
)

// newSUT returns a new system under test.
func newSUT(compare Compare) *filter {
	return New(compare)
}

// TestFilter_Do tests filter.Do.
func TestFilter_Do(t *testing.T) {
	type testCase struct {
		name           string
		compare        Compare
		annotations    []*annotation.Instance
		expectedResult []*annotation.Instance
	}

	cases := []testCase{
		func() testCase {
			return testCase{
				name: "no annotations",
				compare: func(annotation *annotation.Instance) bool {
					return true
				},
				annotations:    []*annotation.Instance{},
				expectedResult: []*annotation.Instance{},
			}
		}(),
		func() testCase {
			idProvider := identityProvider.New(sha256.New())
			id1 := idProvider.Derive(test.FactoryRandomByteSlice())
			id2 := idProvider.Derive(test.FactoryRandomByteSlice())
			metadataKind := test.FactoryRandomString()
			unique := ulid.New().Get()
			match := annotation.New(
				unique,
				id1,
				nil,
				metadataKind,
				test.FactoryRandomString(),
			)

			return testCase{
				name: "single match",
				compare: func(annotation *annotation.Instance) bool {
					return annotation.MetadataKind == metadataKind
				},
				annotations: []*annotation.Instance{
					annotation.New(
						unique,
						id2,
						id1,
						test.FactoryRandomString(),
						test.FactoryRandomString(),
					),
					match,
				},
				expectedResult: []*annotation.Instance{
					match,
				},
			}
		}(),
		func() testCase {
			idProvider := identityProvider.New(sha256.New())
			id1 := idProvider.Derive(test.FactoryRandomByteSlice())
			id2 := idProvider.Derive(test.FactoryRandomByteSlice())
			id3 := idProvider.Derive(test.FactoryRandomByteSlice())
			metadataKind := test.FactoryRandomString()
			unique := ulid.New().Get()
			match1 := annotation.New(
				unique,
				id1,
				nil,
				metadataKind,
				test.FactoryRandomString(),
			)
			match2 := annotation.New(
				unique,
				id2,
				id1,
				metadataKind,
				test.FactoryRandomString(),
			)

			return testCase{
				name: "multiple matches",
				compare: func(annotation *annotation.Instance) bool {
					return annotation.MetadataKind == metadataKind
				},
				annotations: []*annotation.Instance{
					match2,
					annotation.New(
						unique,
						id3,
						id2,
						test.FactoryRandomString(),
						test.FactoryRandomString(),
					),
					match1,
				},
				expectedResult: []*annotation.Instance{
					match2,
					match1,
				},
			}
		}(),
	}

	for i := range cases {
		t.Run(
			cases[i].name,
			func(t *testing.T) {
				sut := newSUT(cases[i].compare)

				result := sut.Do(cases[i].annotations)

				assert.Equal(t, cases[i].expectedResult, result)
			})
	}
}
