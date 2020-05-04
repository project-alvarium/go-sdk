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

package client

import (
	"github.com/iotaledger/iota.go/api"
	"github.com/iotaledger/iota.go/bundle"
	"github.com/iotaledger/iota.go/trinary"
)

// Contract defines the contract used to encapsulate the IOTA Client; separated to facilitate unit testing.
type Contract interface {
	// GetNewAddress is called to retrieve a new available address from the IOTA Tangle.
	GetNewAddress(seed trinary.Trytes, options api.GetNewAddressOptions) (trinary.Hashes, error)

	// SendTransfer is called to send transactions to an IOTA Tangle.
	SendTransfer(
		seed trinary.Trytes,
		depth uint64,
		mwm uint64,
		transfers bundle.Transfers,
		options *api.SendTransfersOptions) (bundle.Bundle, error)
}
