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

package metadata

import (
	"crypto"
	"testing"

	"github.com/project-alvarium/go-sdk/pkg/hashprovider/sha256"
	"github.com/project-alvarium/go-sdk/pkg/test"

	"github.com/stretchr/testify/assert"
)

// TestSuccess_Kind tests success.Kind.
func TestSuccess_Kind(t *testing.T) {
	kind := test.FactoryRandomString()
	sut := NewSuccess(kind, crypto.SHA256, sha256.New().Name())

	assert.Equal(t, kind, sut.Kind())
}