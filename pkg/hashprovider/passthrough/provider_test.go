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

	"github.com/project-alvarium/go-sdk/pkg/test"

	"github.com/stretchr/testify/assert"
)

// newSUT returns a new system under test.
func newSUT() *provider {
	return New()
}

// TestProvider_Derive tests provider.Derive.
func TestProvider_Derive(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Run(
			"variation "+string(i),
			func(t *testing.T) {
				data := test.FactoryRandomByteSlice()
				sut := newSUT()

				result := sut.Derive(data)

				assert.Equal(t, data, result)
			},
		)
	}
}

// TestProvider_Name tests provider.Name.
func TestProvider_Name(t *testing.T) {
	sut := newSUT()
	assert.Equal(t, name, sut.Name())
}

// TestName tests Name.
func TestName(t *testing.T) {
	assert.Equal(t, name, Name())
}
