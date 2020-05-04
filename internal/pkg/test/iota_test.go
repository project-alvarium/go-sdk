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

	"github.com/iotaledger/iota.go/guards/validators"
	"github.com/iotaledger/iota.go/trinary"
	"github.com/stretchr/testify/assert"
)

// TestFactoryRandomSeedString tests FactoryRandomSeedString
func TestFactoryRandomSeedString(t *testing.T) {
	type testCase struct {
		name string
		test func(t *testing.T)
	}

	cases := []testCase{
		{
			name: "returns a valid seed string",
			test: func(t *testing.T) {
				result := FactoryRandomSeedString()
				seedTrytes, _ := trinary.NewTrytes(result)
				err := validators.Validate(validators.ValidateSeed(seedTrytes))

				assert.Nil(t, err)
			},
		},
		{
			name: "returns valid length",
			test: func(t *testing.T) {
				result := FactoryRandomSeedString()
				seedTrytes, _ := trinary.NewTrytes(result)

				assert.Equal(t, seedSize, len(seedTrytes))
			},
		},
		{
			name: "returns valid charset",
			test: func(t *testing.T) {
				result := FactoryRandomSeedString()

				for i := range result {
					assert.True(t, strings.Contains(trytesCharset, string(result[i])))
				}
			},
		},
	}

	for i := range cases {
		t.Run(cases[i].name, cases[i].test)
	}
}

// TestFactoryRandomFixedLengthTrytesString tests FactoryRandomFixedLengthTrytesString.
func TestFactoryRandomFixedLengthTrytesString(t *testing.T) {
	type testCase struct {
		name string
		test func(t *testing.T)
	}

	cases := []testCase{
		{
			name: "returns a valid trytes string",
			test: func(t *testing.T) {
				length := rand.Intn(1024)
				result := FactoryRandomFixedLengthTrytesString(length)
				err := trinary.ValidTrytes(result)

				assert.Nil(t, err)
			},
		},
		{
			name: "returns fixed length",
			test: func(t *testing.T) {
				length := rand.Intn(1024)
				result := FactoryRandomFixedLengthTrytesString(length)

				assert.Equal(t, length, len(result))
			},
		},
		{
			name: "returns valid charset",
			test: func(t *testing.T) {
				length := rand.Intn(1024)
				result := FactoryRandomFixedLengthTrytesString(length)

				for i := range result {
					assert.True(t, strings.Contains(trytesCharset, string(result[i])))
				}
			},
		},
	}

	for i := range cases {
		t.Run(cases[i].name, cases[i].test)
	}
}

// TestFactoryRandomAddressTrytesString tests FactoryRandomAddressTrytesString.
func TestFactoryRandomAddressTrytesString(t *testing.T) {
	type testCase struct {
		name string
		test func(t *testing.T)
	}

	cases := []testCase{
		{
			name: "returns fixed size",
			test: func(t *testing.T) {
				result := FactoryRandomAddressTrytesString()

				assert.Equal(t, addressSize, len(result))
			},
		},
		{
			name: "returns valid charset",
			test: func(t *testing.T) {
				result := FactoryRandomAddressTrytesString()

				for i := range result {
					assert.True(t, strings.Contains(trytesCharset, string(result[i])))
				}
			},
		},
		{
			name: "returns varying content",
			test: func(t *testing.T) {
				result1 := FactoryRandomAddressTrytesString()
				result2 := FactoryRandomAddressTrytesString()

				assert.NotEqual(t, result1, result2)
			},
		},
	}

	for i := range cases {
		t.Run(cases[i].name, cases[i].test)
	}
}

// TestFactoryRandomFixedSizeBundle tests FactoryRandomFixedSizeBundle.
func TestFactoryRandomFixedSizeBundle(t *testing.T) {
	type testCase struct {
		name string
		test func(t *testing.T)
	}

	cases := []testCase{
		{
			name: "fixed bundle size",
			test: func(t *testing.T) {
				size := rand.Intn(1024)
				result := FactoryRandomFixedSizeBundle(size)

				assert.Equal(t, size, len(result))
			},
		},
		{
			name: "returns valid charset",
			test: func(t *testing.T) {
				size := rand.Intn(1024)
				result := FactoryRandomFixedSizeBundle(size)

				for i := range result {
					fields := []string{result[i].Address, result[i].Hash, result[i].Tag}
					for f := range fields {
						for c := range fields[f] {
							assert.True(t, strings.Contains(trytesCharset, string(fields[f][c])))
						}
					}
				}

			},
		},
		{
			name: "random transaction field values",
			test: func(t *testing.T) {
				size := rand.Intn(1024)
				result := FactoryRandomFixedSizeBundle(size)

				for i := range result {
					assert.NotEqual(t, result[i].Address, result[i].Hash)
					assert.NotEqual(t, result[i].Address, result[i].Tag)
					assert.NotEqual(t, result[i].Hash, result[i].Tag)
				}
			},
		},
		{
			name: "transaction fields varying length",
			test: func(t *testing.T) {
				size := rand.Intn(1024)
				result1 := FactoryRandomFixedSizeBundle(size)
				result2 := FactoryRandomFixedSizeBundle(size)

				for i := range result1 {
					assert.NotEqual(t, len(result1[i].Address), len(result2[i].Address))
					assert.NotEqual(t, len(result1[i].Hash), len(result2[i].Hash))
					assert.NotEqual(t, len(result1[i].Tag), len(result2[i].Tag))
				}

			},
		},
	}

	for i := range cases {
		t.Run(cases[i].name, cases[i].test)
	}
}
