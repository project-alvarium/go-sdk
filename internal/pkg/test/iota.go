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

	"github.com/project-alvarium/go-sdk/pkg/test"

	"github.com/iotaledger/iota.go/bundle"
	"github.com/iotaledger/iota.go/transaction"
)

const (
	trytesCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ9"
	seedSize      = 81
	addressSize   = 81
)

// FactoryRandomSeedString returns an IOTA Tangle seed with a random value.
func FactoryRandomSeedString() string {
	return test.FactoryRandomFixedLengthString(seedSize, trytesCharset)
}

// FactoryRandomFixedLengthTrytesString returns a Trytes of a fixed length with a random value.
func FactoryRandomFixedLengthTrytesString(length int) string {
	return test.FactoryRandomFixedLengthString(length, trytesCharset)
}

// FactoryRandomAddressTrytesString returns Trytes for an Address of a fixed length with a random value.
func FactoryRandomAddressTrytesString() string {
	return FactoryRandomFixedLengthTrytesString(addressSize)
}

// FactoryRandomFixedSizeBundle returns a Bundle of a fixed size with random transaction values with a random length.
func FactoryRandomFixedSizeBundle(size int) bundle.Bundle {
	length := rand.Intn(1024)
	txs := make([]transaction.Transaction, size, size)
	for i := 0; i < size; i++ {
		txs[i] = transaction.Transaction{
			Address: FactoryRandomFixedLengthTrytesString(length),
			Hash:    FactoryRandomFixedLengthTrytesString(length),
			Tag:     FactoryRandomFixedLengthTrytesString(length),
		}
	}
	return txs
}
