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

package sha256

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// newSUT returns a new system under test.
func newSUT() *provider {
	return New()
}

// TestProvider_Derive tests provider.Derive.
func TestProvider_Derive(t *testing.T) {
	cases := []struct {
		name     string
		data     []byte
		expected []byte
	}{
		{
			name:     "text variation 1",
			data:     []byte("foo"),
			expected: []byte{0x2c, 0x26, 0xb4, 0x6b, 0x68, 0xff, 0xc6, 0x8f, 0xf9, 0x9b, 0x45, 0x3c, 0x1d, 0x30, 0x41, 0x34, 0x13, 0x42, 0x2d, 0x70, 0x64, 0x83, 0xbf, 0xa0, 0xf9, 0x8a, 0x5e, 0x88, 0x62, 0x66, 0xe7, 0xae},
		},
		{
			name:     "text variation 2",
			data:     []byte("bar"),
			expected: []byte{0xfc, 0xde, 0x2b, 0x2e, 0xdb, 0xa5, 0x6b, 0xf4, 0x8, 0x60, 0x1f, 0xb7, 0x21, 0xfe, 0x9b, 0x5c, 0x33, 0x8d, 0x10, 0xee, 0x42, 0x9e, 0xa0, 0x4f, 0xae, 0x55, 0x11, 0xb6, 0x8f, 0xbf, 0x8f, 0xb9},
		},
		{
			name:     "text variation 3",
			data:     []byte("baz"),
			expected: []byte{0xba, 0xa5, 0xa0, 0x96, 0x4d, 0x33, 0x20, 0xfb, 0xc0, 0xc6, 0xa9, 0x22, 0x14, 0x4, 0x53, 0xc8, 0x51, 0x3e, 0xa2, 0x4a, 0xb8, 0xfd, 0x5, 0x77, 0x3, 0x48, 0x4, 0xa9, 0x67, 0x24, 0x80, 0x96},
		},
		{
			name:     "byte sequence",
			data:     []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
			expected: []byte{0x9a, 0x89, 0xc6, 0x8c, 0x4c, 0x5e, 0x28, 0xb8, 0xc4, 0xa5, 0x56, 0x76, 0x73, 0xd4, 0x62, 0xff, 0xf5, 0x15, 0xdb, 0x46, 0x11, 0x6f, 0x99, 0x0, 0x62, 0x4d, 0x9, 0xc4, 0x74, 0xf5, 0x93, 0xfb},
		},
	}

	for i := range cases {
		t.Run(
			cases[i].name,
			func(t *testing.T) {
				sut := newSUT()

				result := sut.Derive(cases[i].data)

				assert.Equal(t, cases[i].expected, result)
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
