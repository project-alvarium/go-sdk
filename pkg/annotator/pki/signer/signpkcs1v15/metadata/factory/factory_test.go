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
	"crypto"
	"encoding/json"
	"testing"

	testInternal "github.com/project-alvarium/go-sdk/internal/pkg/test"
	pkcsSignerMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signpkcs1v15/metadata"
	"github.com/project-alvarium/go-sdk/pkg/hashprovider/sha256"
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
			name: "Valid (signpcks1v15 failure)",
			test: func(t *testing.T) {
				sut := newSUT()
				value := pkcsSignerMetadata.NewFailure(test.FactoryRandomString())

				result := sut.Create(pkcsSignerMetadata.Kind, json.RawMessage(testInternal.Marshal(t, value)))

				assert.NotNil(t, result)
				assert.IsType(t, &pkcsSignerMetadata.Failure{}, result)
				assert.Equal(t, testInternal.Marshal(t, value), testInternal.Marshal(t, result))
			},
		},
		{
			name: "Valid (signpcks1v15 success)",
			test: func(t *testing.T) {
				sut := newSUT()
				value := pkcsSignerMetadata.NewSuccess(crypto.SHA256, sha256.Kind)

				result := sut.Create(pkcsSignerMetadata.Kind, json.RawMessage(testInternal.Marshal(t, value)))

				assert.NotNil(t, result)
				assert.IsType(t, &pkcsSignerMetadata.Success{}, result)
				assert.Equal(t, testInternal.Marshal(t, value), testInternal.Marshal(t, result))
			},
		},
	}

	for i := range cases {
		t.Run(cases[i].name, cases[i].test)
	}
}
