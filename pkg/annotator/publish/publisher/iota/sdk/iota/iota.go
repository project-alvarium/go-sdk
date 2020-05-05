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

package iota

import (
	"fmt"

	"github.com/project-alvarium/go-sdk/pkg/annotation/metadata"
	publishMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/publish/metadata"
	publisherMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/iota/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/iota/sdk/iota/client"

	iotaAPI "github.com/iotaledger/iota.go/api"
	"github.com/iotaledger/iota.go/bundle"
	"github.com/iotaledger/iota.go/converter"
	"github.com/iotaledger/iota.go/trinary"
)

// instance is a receiver that encapsulates required dependencies.
type instance struct {
	kind   string
	client client.Contract
}

// New is a factory function that returns an initialized instance.
func New(kind string, client client.Contract) *instance {
	return &instance{
		kind:   kind,
		client: client,
	}
}

// createTransfer is a factory function that returns an initialized instance of a Transfer.
func (*instance) createTransfer(address trinary.Hash, message trinary.Trytes) bundle.Transfer {
	return bundle.Transfer{
		Address: address,
		Message: message,
		Value:   0,
	}
}

// composeAPIError returns annotations for failure case; separated to facilitate unit testing.
func (i *instance) composeAPIError(message string) *publishMetadata.Failure {
	return publishMetadata.NewFailure(i.kind, fmt.Sprintf("iotaAPI.ComposeApi() returned \"%s\"", message))
}

// convertToTrytesError returns annotations for failure case; separated to facilitate unit testing.
func (i *instance) convertToTrytesError(message string) *publishMetadata.Failure {
	return publishMetadata.NewFailure(i.kind, fmt.Sprintf("converter.ASCIIToTrytes() returned \"%s\"", message))
}

// getNewAddressError returns annotations for failure case; separated to facilitate unit testing.
func (i *instance) getNewAddressError(message string) *publishMetadata.Failure {
	return publishMetadata.NewFailure(i.kind, fmt.Sprintf("iotaAPI.GetNewAddress() returned \"%s\"", message))
}

// sendTransferError returns annotations for failure case; separated to facilitate unit testing.
func (i *instance) sendTransferError(message string) *publishMetadata.Failure {
	return publishMetadata.NewFailure(i.kind, fmt.Sprintf("iotaAPI.sendTransfer() returned \"%s\"", message))
}

// invalidResultSetError returns annotations for failure case; separated to facilitate unit testing.
func (i *instance) invalidResultSetError() *publishMetadata.Failure {
	return i.sendTransferError("Expected result Bundle to contain size of 1")
}

// Send is called to send annotations to an IOTA Tangle.
func (i *instance) Send(seed string, depth uint64, mwm uint64, annotations []byte) metadata.Contract {
	addr, err := i.client.GetNewAddress(seed, iotaAPI.GetNewAddressOptions{})
	if err != nil {
		return i.getNewAddressError(err.Error())
	}

	messageTrytes, err := converter.ASCIIToTrytes(string(annotations))
	if err != nil {
		return i.convertToTrytesError(err.Error())
	}

	res, err := i.client.SendTransfer(
		seed,
		depth,
		mwm,
		[]bundle.Transfer{i.createTransfer(addr[0], messageTrytes)},
		&iotaAPI.SendTransfersOptions{},
	)
	if err != nil {
		return i.sendTransferError(err.Error())
	}
	if len(res) != 1 {
		return i.invalidResultSetError()
	}
	return publisherMetadata.New(i.kind, res[0].Address, res[0].Hash, res[0].Tag)
}
