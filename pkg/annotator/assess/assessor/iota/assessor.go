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
	"encoding/json"
	"github.com/iotaledger/iota.go/address"
	"github.com/iotaledger/iota.go/api"
	"github.com/iotaledger/iota.go/transaction"
	"github.com/iotaledger/iota.go/trinary"
	"github.com/project-alvarium/go-sdk/pkg/annotation"
	"github.com/project-alvarium/go-sdk/pkg/annotation/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/iota/client"
	iotaAssessorMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/iota/metadata"
	publisherMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/publish/metadata"
	iota2 "github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/iota"
	iotaAnnotatorMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/iota/metadata"
	ipfsAnnotatorMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/ipfs/metadata"
)

const name = iota2.Name

// assessor is a receiver that encapsulates required dependencies.
type assessor struct {
	client client.Contract
	ipfsMetadata *ipfsAnnotatorMetadata.Instance
}

// New is a factory function that returns an initialized assessor.
func New(client client.Contract) *assessor {
	return &assessor{
		client: client,
	}
}

// SetUp is called once when the assessor is instantiated.
func (*assessor) SetUp() {}

// TearDown is called once when assessor is terminated.
func (*assessor) TearDown() {}

func (*assessor) invalidTransactionAssessment(annotation *annotation.Instance) *iotaAssessorMetadata.Instance {
	return iotaAssessorMetadata.New(false, []string{annotation.Unique})
}

// Assess accepts data and returns associated assessments.
func (a *assessor) Assess(annotations []*annotation.Instance) metadata.Contract {
	uniques := make([]string, 0)
	for i := range annotations {
		if annotations[i].MetadataKind != publisherMetadata.Kind {
			continue
		}

		switch m := annotations[i].Metadata.(*publisherMetadata.Success).PublisherMetadata.(type) {
		case *ipfsAnnotatorMetadata.Instance:
			a.ipfsMetadata = m
		case *iotaAnnotatorMetadata.Instance:
			checksum, err := address.Checksum(m.Address)
			if err != nil {
				return a.invalidTransactionAssessment(annotations[i])
			}
			transactions, err := a.client.FindTransactionObjects(api.FindTransactionsQuery{
				Addresses: trinary.Hashes{m.Address + checksum},
			})
			if err != nil {
				return a.invalidTransactionAssessment(annotations[i])
			}

			txJson, err := transaction.ExtractJSON(transactions)
			if err != nil || len(transactions) != 1 {
				return a.invalidTransactionAssessment(annotations[i])
			}

			txIpfsMetadata := make(map[string]json.RawMessage)
			_ = json.Unmarshal([]byte(txJson), &txIpfsMetadata)
			var txIpfsPublisherMetadata *ipfsAnnotatorMetadata.Instance
			for k, v := range txIpfsMetadata {
				if k == publisherMetadata.PublisherMetadata{
					_ = json.Unmarshal(v, &txIpfsPublisherMetadata)
					continue
				}
			}

			if txIpfsPublisherMetadata.CID != a.ipfsMetadata.CID {
				return a.invalidTransactionAssessment(annotations[i])
			}

		}

		uniques = append(uniques, annotations[i].Unique)
	}

	return iotaAssessorMetadata.New(true, uniques)
}

// Kind returns an implementation mnemonic
func (*assessor) Kind() string {
	return name
}
