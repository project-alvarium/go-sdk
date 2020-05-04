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

package stub

import (
	"github.com/iotaledger/iota.go/api"
	"github.com/iotaledger/iota.go/bundle"
	"github.com/iotaledger/iota.go/trinary"
)

// instance is a receiver that encapsulates required dependencies.
type Instance struct {
	addressResultValue trinary.Hashes
	addressResultError error
	sendResultValue    bundle.Bundle
	sendResultError    error
}

// New is a factory function that returns an initialized instance.
func New(
	addressResultValue trinary.Hashes,
	addressResultError error,
	sendResultValue bundle.Bundle,
	sendResultError error) *Instance {

	return &Instance{
		addressResultValue: addressResultValue,
		addressResultError: addressResultError,
		sendResultValue:    sendResultValue,
		sendResultError:    sendResultError,
	}
}

// GetNewAddress is called to retrieve a new available address from the IOTA Tangle.
func (i *Instance) GetNewAddress(_ trinary.Trytes, _ api.GetNewAddressOptions) (trinary.Hashes, error) {
	return i.addressResultValue, i.addressResultError
}

// SendTransfer is called to send transactions to an IOTA Tangle.
func (i *Instance) SendTransfer(
	_ trinary.Trytes,
	_ uint64,
	_ uint64,
	_ bundle.Transfers,
	_ *api.SendTransfersOptions) (bundle.Bundle, error) {

	return i.sendResultValue, i.sendResultError
}
