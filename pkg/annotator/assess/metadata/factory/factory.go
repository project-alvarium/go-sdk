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
	pkiAssessorFactory "github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/pki/metadata/factory"
	assessMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/assess/metadata"
)

// instance is a receiver that encapsulates required dependencies.
type instance struct {
	assessorFactories []factory.Contract
}

// New is a factory function that returns an initialized instance.
func New(assessorFactories []factory.Contract) *instance {
	return &instance{
		assessorFactories: assessorFactories,
	}
}

// New is a factory function that returns an initialized instance.
func NewDefault() *instance {
	return New(
		[]factory.Contract{
			pkiAssessorFactory.New(),
		},
	)
}

// Create returns a contract implementation based on the provided metadata.
func (i *instance) Create(kind string, data json.RawMessage) metadata.Contract {
	switch kind {
	case assessMetadata.Kind:
		if string(data) == "null" {
			return nil
		}

		var concrete assessMetadata.Instance
		concrete.SetAssessorFactories(i.assessorFactories)
		if err := json.Unmarshal(data, &concrete); err == nil {
			return &concrete
		}
	}
	return nil
}
