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

	"github.com/project-alvarium/go-sdk/pkg/annotation"
	"github.com/project-alvarium/go-sdk/pkg/annotation/metadata"
	metadataFactory "github.com/project-alvarium/go-sdk/pkg/annotation/metadata/factory"
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/iota/client"
	iotaAssessorMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/iota/metadata"
	publishMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/publish/metadata"
	iotaPublisherMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/iota/metadata"

	"github.com/iotaledger/iota.go/address"
	"github.com/iotaledger/iota.go/api"
	"github.com/iotaledger/iota.go/trinary"
)

// assessor is a receiver that encapsulates required dependencies.
type assessor struct {
	client  client.Contract
	factory metadataFactory.Contract
}

// New is a factory function that returns an initialized assessor.
func New(client client.Contract, factory metadataFactory.Contract) *assessor {
	return &assessor{
		client:  client,
		factory: factory,
	}
}

// SetUp is called once when the assessor is instantiated.
func (*assessor) SetUp() {}

// TearDown is called once when assessor is terminated.
func (*assessor) TearDown() {}

// failureAnnotationMatch returns an annotation with a failure case; separated to facilitate unit testing.
func (*assessor) failureAnnotationMatch() *iotaAssessorMetadata.Failure {
	return iotaAssessorMetadata.NewFailure("failed to find an IOTA assessor annotation")
}

// failureTransactionResultError returns an annotation with a failure case; separated to facilitate unit testing.
func (*assessor) failureTransactionResultError() *iotaAssessorMetadata.Failure {
	return iotaAssessorMetadata.NewFailure("expected 1 IOTA annotation result")
}

// failureAddressChecksumError returns an annotation with a failure case; separated to facilitate unit testing.
func (*assessor) failureAddressChecksumError(errorMessage string) *iotaAssessorMetadata.Failure {
	return iotaAssessorMetadata.NewFailure(fmt.Sprintf("address.Checksum() returned \"%s\"", errorMessage))
}

// failureFindTransactionError returns an annotation with a failure case; separated to facilitate unit testing.
func (*assessor) failureFindTransactionError(errorMessage string) *iotaAssessorMetadata.Failure {
	return iotaAssessorMetadata.NewFailure(fmt.Sprintf("client.FindTransactionObjects returned \"%s\"", errorMessage))
}

// Assess accepts data and returns associated assessments.
func (a *assessor) Assess(annotations []*annotation.Instance) metadata.Contract {
	uniques := make([]string, 0)
	for i := range annotations {
		if annotations[i].MetadataKind != publishMetadata.Kind {
			continue
		}

		if m, ok := annotations[i].Metadata.(*publishMetadata.Instance); ok {
			if pm, pOk := m.PublisherMetadata.(*iotaPublisherMetadata.Success); pOk {
				addrChecksum, err := address.Checksum(pm.Address)
				if err != nil {
					return a.failureAddressChecksumError(err.Error())
				}

				txs, err := a.client.FindTransactionObjects(
					api.FindTransactionsQuery{
						Addresses: []trinary.Hash{
							pm.Address + addrChecksum,
						},
					},
				)
				if err != nil {
					return a.failureFindTransactionError(err.Error())
				}

				if len(txs) != 1 {
					return a.failureTransactionResultError()
				}

				uniques = append(uniques, annotations[i].Unique)
			}

			if _, pOk := m.PublisherMetadata.(*iotaPublisherMetadata.Failure); pOk {
				return iotaAssessorMetadata.NewSuccess(false, []string{annotations[i].Unique})
			}
		}
	}

	if len(annotations) == 0 || len(uniques) == 0 {
		return a.failureAnnotationMatch()
	}

	return iotaAssessorMetadata.NewSuccess(true, uniques)
}

// Failure creates a publisher-specific failure annotation.
func (*assessor) Failure(errorMessage string) metadata.Contract {
	return iotaAssessorMetadata.NewFailure(errorMessage)
}

// Kind returns an implementation mnemonic
func (*assessor) Kind() string {
	return iotaPublisherMetadata.Kind
}
