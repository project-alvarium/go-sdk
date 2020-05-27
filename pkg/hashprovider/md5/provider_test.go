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

package md5

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
			expected: []byte{0xac, 0xbd, 0x18, 0xdb, 0x4c, 0xc2, 0xf8, 0x5c, 0xed, 0xef, 0x65, 0x4f, 0xcc, 0xc4, 0xa4, 0xd8},
		},
		{
			name:     "text variation 2",
			data:     []byte("bar"),
			expected: []byte{0x37, 0xb5, 0x1d, 0x19, 0x4a, 0x75, 0x13, 0xe4, 0x5b, 0x56, 0xf6, 0x52, 0x4f, 0x2d, 0x51, 0xf2},
		},
		{
			name:     "text variation 3",
			data:     []byte("baz"),
			expected: []byte{0x73, 0xfe, 0xff, 0xa4, 0xb7, 0xf6, 0xbb, 0x68, 0xe4, 0x4c, 0xf9, 0x84, 0xc8, 0x5f, 0x6e, 0x88},
		},
		{
			name:     "byte sequence",
			data:     []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
			expected: []byte{0x7f, 0x63, 0xcb, 0x6d, 0x6, 0x79, 0x72, 0xc3, 0xf3, 0x4f, 0x9, 0x4b, 0xb7, 0xe7, 0x76, 0xa8},
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

// TestProvider_Kind tests provider.Kind.
func TestProvider_Kind(t *testing.T) {
	sut := newSUT()
	assert.Equal(t, Kind, sut.Kind())
}
