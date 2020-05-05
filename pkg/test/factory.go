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
	"time"
)

const (
	HexCharset          = "abcdef0123456789"
	AlphanumericCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	maxInt              = 1024
)

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// FactoryRandomFixedLengthString returns a string of fixed length with a random value.
func FactoryRandomFixedLengthString(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// FactoryRandomString returns a string with a random value.
func FactoryRandomString() string {
	return FactoryRandomFixedLengthString(FactoryRandomInt(), AlphanumericCharset)
}

// FactoryRandomFixedLengthAlphanumericString returns a string of fixed length with an random alphanumeric value.
func FactoryRandomFixedLengthAlphanumericString(length int) string {
	return FactoryRandomFixedLengthString(length, AlphanumericCharset)
}

// FactoryRandomFixedLengthByteSlice returns a []byte with a random value.
func FactoryRandomFixedLengthByteSlice(length int, charset string) []byte {
	return []byte(FactoryRandomFixedLengthString(length, charset))
}

// FactoryRandomByteSlice returns a []byte with a random value.
func FactoryRandomByteSlice() []byte {
	return []byte(FactoryRandomString())
}

// FactoryRandomFixedLengthAlphanumericByteSlice returns a []byte of fixed length and random alphanumeric value.
func FactoryRandomFixedLengthAlphanumericByteSlice(length int) []byte {
	return FactoryRandomFixedLengthByteSlice(length, AlphanumericCharset)
}

// FactoryRandomInt returns an int with a random value.
func FactoryRandomInt() int {
	return rand.Intn(maxInt)
}

// FactoryRandomUint64 returns a uint64 with a random value.
func FactoryRandomUint64() uint64 {
	return rand.Uint64()
}
