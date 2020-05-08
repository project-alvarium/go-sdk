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

	"github.com/project-alvarium/go-sdk/pkg/annotation"
	"github.com/project-alvarium/go-sdk/pkg/annotation/metadata"
	publishMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/publish/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/iota/sdk"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/iota/sdk/iota"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/iota/sdk/iota/client"
)

const Name = "iota"

// publisher is a receiver that encapsulates required dependencies.
type publisher struct {
	seed  string
	depth uint64
	mwm   uint64
	sdk   sdk.Contract
}

// newWithIota is a factory function that returns an initialized publisher.
func newWithIOTA(seed string, depth uint64, mwm uint64, sdk sdk.Contract) *publisher {
	return &publisher{
		seed:  seed,
		depth: depth,
		mwm:   mwm,
		sdk:   sdk,
	}
}

// New is a factory function that returns an initialized publisher.
func New(seed string, depth uint64, mwm uint64, client client.Contract) *publisher {
	return newWithIOTA(seed, depth, mwm, iota.New(Name, client))
}

// SetUp is called once when the publisher is instantiated.
func (*publisher) SetUp() {}

// TearDown is called once when publisher is terminated.
func (*publisher) TearDown() {}

// failureNoAnnotations returns annotations for failure case; separated to facilitate unit testing.
func (p *publisher) failureNoAnnotations() *publishMetadata.Failure {
	return publishMetadata.NewFailure(p.Kind(), "IOTA Tangle publisher received 0 annotations")
}

// Publish retrieves and "publishes" annotations.
func (p *publisher) Publish(annotations []*annotation.Instance) metadata.Contract {
	if len(annotations) == 0 {
		return p.failureNoAnnotations()
	}

	marshalledAnnotations, _ := json.Marshal(annotations)
	return p.sdk.Send(p.seed, p.depth, p.mwm, marshalledAnnotations)
}

// Kind returns an implementation mnemonic.
func (*publisher) Kind() string {
	return Kind()
}

// Kind returns an implementation mnemonic.
func Kind() string {
	return Name
}
