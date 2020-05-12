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

package factory

import (
	"encoding/json"

	"github.com/project-alvarium/go-sdk/pkg/annotation/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotation/metadata/factory"
	publishMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/publish/metadata"
	exampleMetadataFactory "github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/example/metadata/factory"
	iotaMetadataFactory "github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/iota/metadata/factory"
	ipfsMetadataFactory "github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/ipfs/metadata/factory"
)

// instance is a receiver that encapsulates required dependencies.
type instance struct {
	publisherFactories []factory.Contract
}

// New is a factory function that returns an initialized instance.
func New(publisherFactories []factory.Contract) *instance {
	return &instance{
		publisherFactories: publisherFactories,
	}
}

// New is a factory function that returns an initialized instance.
func NewDefault() *instance {
	return New(
		[]factory.Contract{
			exampleMetadataFactory.New(),
			iotaMetadataFactory.New(),
			ipfsMetadataFactory.New(),
		},
	)
}

// Create returns a contract implementation based on the provided metadata.
func (i *instance) Create(kind string, data json.RawMessage) metadata.Contract {
	switch kind {
	case publishMetadata.Kind:
		var concrete publishMetadata.Instance
		concrete.SetPublisherFactories(i.publisherFactories)
		if err := json.Unmarshal(data, &concrete); err == nil {
			return &concrete
		}
	}
	return nil
}