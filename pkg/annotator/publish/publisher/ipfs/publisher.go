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

package ipfs

import (
	"encoding/json"
	"fmt"

	"github.com/project-alvarium/go-sdk/pkg/annotation"
	"github.com/project-alvarium/go-sdk/pkg/annotation/metadata"
	publishMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/publish/metadata"
	publisherMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/ipfs/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/ipfs/sdk"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/ipfs/sdk/ipfs"
)

const name = "ipfs"

// publisher is a receiver that encapsulates required dependencies.
type publisher struct {
	url string
	sdk sdk.Contract
}

// newWithIPFS is a factory function that returns an initialized publisher.
func newWithIPFS(url string, sdk sdk.Contract) *publisher {
	return &publisher{
		url: url,
		sdk: sdk,
	}
}

// New is a factory function that returns an initialized publisher.
func New(url string) *publisher {
	return newWithIPFS(url, ipfs.New())
}

// SetUp is called once when the publisher is instantiated.
func (p *publisher) SetUp() {}

// TearDown is called once when publisher is terminated.
func (p *publisher) TearDown() {}

// failureAdd returns annotations for failure case; separated to facilitate unit testing.
func (p *publisher) failureAdd(message string) *publishMetadata.Failure {
	return publishMetadata.NewFailure(p.Kind(), fmt.Sprintf("Add returned \"%s\"", message))
}

// Publish retrieves and "publishes" annotations.
func (p *publisher) Publish(annotations []*annotation.Instance) metadata.Contract {
	marshaledAnnotations, _ := json.Marshal(annotations)
	cid, err := p.sdk.Add(p.url, marshaledAnnotations)
	if err != nil {
		return p.failureAdd(err.Error())
	}

	return publisherMetadata.New(name, cid)
}

// Kind returns an implementation mnemonic.
func (*publisher) Kind() string {
	return Kind()
}

// Kind returns an implementation mnemonic.
func Kind() string {
	return name
}
