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

package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEncode tests Encode.
func TestEncode(t *testing.T) {
	type testCase struct {
		name     string
		value    []byte
		expected []byte
	}

	cases := []testCase{
		{
			name:     "text variation 1",
			value:    []byte("foo"),
			expected: []byte{0x36, 0x36, 0x36, 0x66, 0x36, 0x66},
		},
		{
			name:     "text variation 2",
			value:    []byte("bar"),
			expected: []byte{0x36, 0x32, 0x36, 0x31, 0x37, 0x32},
		},
		{
			name:     "byte variation 1",
			value:    []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09},
			expected: []byte{0x30, 0x30, 0x30, 0x31, 0x30, 0x32, 0x30, 0x33, 0x30, 0x34, 0x30, 0x35, 0x30, 0x36, 0x30, 0x37, 0x30, 0x38, 0x30, 0x39},
		},
		{
			name:     "byte variation 2",
			value:    []byte{0x09, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 0x00},
			expected: []byte{0x30, 0x39, 0x30, 0x38, 0x30, 0x37, 0x30, 0x36, 0x30, 0x35, 0x30, 0x34, 0x30, 0x33, 0x30, 0x32, 0x30, 0x31, 0x30, 0x30},
		},
	}

	for i := range cases {
		t.Run(
			cases[i].name,
			func(t *testing.T) {
				result := Encode(cases[i].value)

				assert.Equal(t, cases[i].expected, result)
			},
		)
	}
}
