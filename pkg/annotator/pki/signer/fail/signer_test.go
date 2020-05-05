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

package fail

import (
	"testing"

	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/fail/metadata"
	"github.com/project-alvarium/go-sdk/pkg/test"

	"github.com/stretchr/testify/assert"
)

// TestSigner_Sign tests fail.Sign.
func TestSigner_Sign(t *testing.T) {
	sut := New()

	identitySignature, dataSignature := sut.Sign(test.FactoryRandomByteSlice(), test.FactoryRandomByteSlice())

	assert.Nil(t, identitySignature)
	assert.Nil(t, dataSignature)
}

// TestSigner_SetUp tests fail.SetUp.
func TestSigner_SetUp(t *testing.T) {
	sut := New()

	// for coverage; no assertion
	sut.SetUp()
}

// TestSigner_TearDown tests fail.TearDown.
func TestSigner_TearDown(t *testing.T) {
	sut := New()

	// for coverage; no assertion
	sut.TearDown()
}

// TestSigner_PublicKey tests fail.PublicKey.
func TestSigner_PublicKey(t *testing.T) {
	sut := New()

	assert.Nil(t, sut.PublicKey())
}

// TestSigner_Metadata tests fail.Metadata.
func TestSigner_Metadata(t *testing.T) {
	sut := New()

	assert.Equal(t, metadata.New(), sut.Metadata())
}
