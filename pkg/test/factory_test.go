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
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFactoryRandomFixedLengthString tests FactoryRandomFixedLengthString
func TestFactoryRandomFixedLengthString(t *testing.T) {
	type testCase struct {
		name string
		test func(t *testing.T)
	}

	cases := []testCase{
		{
			name: "returns string of expected length",
			test: func(t *testing.T) {
				for i := 0; i < 10; i++ {
					length := rand.Intn(1024)

					result := FactoryRandomFixedLengthString(length, "*")

					assert.Len(t, result, length)
				}
			},
		},
		{
			name: "returns string of expected content",
			test: func(t *testing.T) {
				charset := "*"
				length := rand.Intn(1024)

				result := FactoryRandomFixedLengthString(length, charset)

				assert.Equal(t, length, strings.Count(result, charset))
			},
		},
	}

	for i := range cases {
		t.Run(cases[i].name, cases[i].test)
	}
}

// TestFactoryRandomString tests FactoryRandomString
func TestFactoryRandomString(t *testing.T) {
	type testCase struct {
		name string
		test func(t *testing.T)
	}

	cases := []testCase{
		{
			name: "returns string of varying length",
			test: func(t *testing.T) {
				result1 := FactoryRandomString()
				result2 := FactoryRandomString()

				assert.NotEqual(t, len(result1), len(result2))
			},
		},
		{
			name: "returns string of varying content",
			test: func(t *testing.T) {
				result1 := FactoryRandomString()
				result2 := FactoryRandomString()

				assert.NotEqual(t, result1, result2)
			},
		},
		{
			name: "returns string limited to expected charset",
			test: func(t *testing.T) {
				result := FactoryRandomString()

				for i := range result {
					assert.True(t, strings.Contains(AlphanumericCharset, string(result[i])))
				}
			},
		},
	}

	for i := range cases {
		t.Run(cases[i].name, cases[i].test)
	}
}

// TestFactoryRandomFixedLengthAlphanumericString tests FactoryRandomFixedLengthAlphanumericString
func TestFactoryRandomFixedLengthAlphanumericString(t *testing.T) {
	type testCase struct {
		name string
		test func(t *testing.T)
	}

	cases := []testCase{
		{
			name: "returns a string of fixed length",
			test: func(t *testing.T) {
				length := rand.Intn(1024)
				result := FactoryRandomFixedLengthAlphanumericString(length)

				assert.Equal(t, length, len(result))
			},
		},
		{
			name: "returns a string of varying content",
			test: func(t *testing.T) {
				length := rand.Intn(1024)
				result1 := FactoryRandomFixedLengthAlphanumericString(length)
				result2 := FactoryRandomFixedLengthAlphanumericString(length)

				assert.NotEqual(t, result1, result2)
			},
		},
		{
			name: "returns a string limited to the expected charset",
			test: func(t *testing.T) {
				length := rand.Intn(1024)
				result := FactoryRandomFixedLengthAlphanumericString(length)

				for i := range result {
					assert.True(t, strings.Contains(AlphanumericCharset, string(result[i])))
				}
			},
		},
	}

	for i := range cases {
		t.Run(cases[i].name, cases[i].test)
	}
}

// TestFactoryRandomFixedLengthByteSlice tests FactoryRandomFixedLengthByteSlice
func TestFactoryRandomFixedLengthByteSlice(t *testing.T) {
	type testCase struct {
		name string
		test func(t *testing.T)
	}

	cases := []testCase{
		{
			name: "returns []byte of fixed length",
			test: func(t *testing.T) {
				length := rand.Intn(1024)
				result := FactoryRandomFixedLengthByteSlice(length, AlphanumericCharset)

				assert.Equal(t, length, len(result))
			},
		},
		{
			name: "returns []bytes of varying content",
			test: func(t *testing.T) {
				length := rand.Intn(1024)
				result1 := FactoryRandomFixedLengthByteSlice(length, AlphanumericCharset)
				result2 := FactoryRandomFixedLengthByteSlice(length, AlphanumericCharset)

				assert.NotEqual(t, result1, result2)
			},
		},
		{
			name: "returns []bytes of limited to expected charset",
			test: func(t *testing.T) {
				length := rand.Intn(1024)
				result := FactoryRandomFixedLengthByteSlice(length, AlphanumericCharset)

				for i := range result {
					assert.True(t, strings.Contains(AlphanumericCharset, string(result[i])))
				}
			},
		},
	}

	for i := range cases {
		t.Run(cases[i].name, cases[i].test)
	}
}

// TestFactoryRandomByteSlice tests FactoryRandomByteSlice
func TestFactoryRandomByteSlice(t *testing.T) {
	type testCase struct {
		name string
		test func(t *testing.T)
	}

	cases := []testCase{
		{
			name: "returns []byte of varying length",
			test: func(t *testing.T) {
				result1 := FactoryRandomByteSlice()
				result2 := FactoryRandomByteSlice()

				assert.NotEqual(t, len(result1), len(result2))
			},
		},
		{
			name: "returns []byte of varying content",
			test: func(t *testing.T) {
				result1 := FactoryRandomByteSlice()
				result2 := FactoryRandomByteSlice()

				assert.NotEqual(t, result1, result2)
			},
		},
		{
			name: "returns []byte limited to expected charset",
			test: func(t *testing.T) {
				result := FactoryRandomByteSlice()

				for i := range result {
					assert.True(t, strings.Contains(AlphanumericCharset, string(result[i])))
				}
			},
		},
	}

	for i := range cases {
		t.Run(cases[i].name, cases[i].test)
	}
}

// TestFactoryRandomFixedLengthAlphanumericByteSlice tests FactoryRandomFixedLengthAlphanumericByteSlice
func TestFactoryRandomFixedLengthAlphanumericByteSlice(t *testing.T) {
	type testCase struct {
		name string
		test func(t *testing.T)
	}

	cases := []testCase{
		{
			name: "returns a byte slice of fixed length",
			test: func(t *testing.T) {
				length := rand.Intn(1024)
				result := FactoryRandomFixedLengthAlphanumericByteSlice(length)

				assert.Equal(t, length, len(result))
			},
		},
		{
			name: "returns a byte[] of varying content",
			test: func(t *testing.T) {
				length := rand.Intn(1024)
				result1 := FactoryRandomFixedLengthAlphanumericByteSlice(length)
				result2 := FactoryRandomFixedLengthAlphanumericByteSlice(length)

				assert.NotEqual(t, result1, result2)
			},
		},
		{
			name: "returns a byte[] limited to the expected charset",
			test: func(t *testing.T) {
				length := rand.Intn(1024)
				result := FactoryRandomFixedLengthAlphanumericByteSlice(length)

				for i := range result {
					assert.True(t, strings.Contains(AlphanumericCharset, string(result[i])))
				}
			},
		},
	}

	for i := range cases {
		t.Run(cases[i].name, cases[i].test)
	}
}

// TestFactoryRandomInt tests FactoryRandomInt
func TestFactoryRandomInt(t *testing.T) {
	type testCase struct {
		name string
		test func(t *testing.T)
	}

	cases := []testCase{
		{
			name: "returns value in range",
			test: func(t *testing.T) {
				for i := 0; i < 1000; i++ {
					result := FactoryRandomInt()

					assert.True(t, result >= 0)
					assert.True(t, result <= maxInt)
				}
			},
		},
	}

	for i := range cases {
		t.Run(cases[i].name, cases[i].test)
	}
}

// TestFactoryRandomUint64 tests FactoryRandomUint64
func TestFactoryRandomUint64(t *testing.T) {
	type testCase struct {
		name string
		test func(t *testing.T)
	}

	cases := []testCase{
		{
			name: "returns value in range",
			test: func(t *testing.T) {
				for i := 0; i < 1000; i++ {
					result := FactoryRandomUint64()

					assert.True(t, result >= 0)
				}
			},
		},
	}

	for i := range cases {
		t.Run(cases[i].name, cases[i].test)
	}
}
