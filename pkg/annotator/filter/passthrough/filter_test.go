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

package passthrough

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
func newSUT() *filter {
	return New()
}

// TestFilter_Do tests passthrough.Do.
func TestFilter_Do(t *testing.T) {
	for i := 0; i < 10; i++ {
		annotations := []*annotation.Instance{
			annotation.New(
				ulid.New().Get(),
				identityProvider.New(sha256.New()).Derive(test.FactoryRandomByteSlice()),
				identityProvider.New(sha256.New()).Derive(test.FactoryRandomByteSlice()),
				test.FactoryRandomString(),
				nil,
			),
		}
		sut := newSUT()

		result := sut.Do(annotations)

		assert.Equal(t, annotations, result)
	}
}
