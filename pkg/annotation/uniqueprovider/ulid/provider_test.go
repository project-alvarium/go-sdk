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

package ulid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// newSUT returns a new system under test.
func newSUT() *provider {
	return New()
}

// TestProvider_Get tests provider.Get.
func TestProvider_Get(t *testing.T) {
	cases := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "not nil",
			test: func(t *testing.T) {
				sut := newSUT()

				assert.NotNil(t, sut.Get())
			},
		},
		{
			name: "not empty string",
			test: func(t *testing.T) {
				sut := newSUT()

				assert.NotEqual(t, "", sut.Get())
			},
		},
		{
			name: "not same string",
			test: func(t *testing.T) {
				sut := newSUT()

				assert.NotEqual(t, sut.Get(), sut.Get())
			},
		},
	}

	for i := range cases {
		t.Run(cases[i].name, cases[i].test)
	}
}
