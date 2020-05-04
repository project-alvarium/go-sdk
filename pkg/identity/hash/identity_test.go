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

package hash

import (
	"testing"

	"github.com/project-alvarium/go-sdk/pkg/test"

	"github.com/stretchr/testify/assert"
)

// newSUT returns a new system under test.
func newSUT(hash []byte) *identity {
	return New(hash)
}

// TestIdentity_Binary tests identity.Binary.
func TestIdentity_Binary(t *testing.T) {
	h := test.FactoryRandomByteSlice()
	sut := newSUT(h)

	result := sut.Binary()

	assert.Equal(t, h, result)
}

// TestIdentity_Printable tests identity.Printable.
func TestIdentity_Printable(t *testing.T) {
	h := test.FactoryRandomByteSlice()
	sut := newSUT(h)

	result := sut.Printable()

	assert.Equal(t, toPrintable(h), result)
}

// TestIdentity_Kind tests identity.Kind.
func TestIdentity_Kind(t *testing.T) {
	sut := newSUT(test.FactoryRandomByteSlice())

	result := sut.Kind()

	assert.Equal(t, name, result)
}
