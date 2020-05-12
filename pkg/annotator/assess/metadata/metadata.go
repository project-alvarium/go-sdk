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

package metadata

import (
	"encoding/json"
	"errors"

	"github.com/project-alvarium/go-sdk/pkg/annotation/metadata"
	metadataFactory "github.com/project-alvarium/go-sdk/pkg/annotation/metadata/factory"
	"github.com/project-alvarium/go-sdk/pkg/annotator/provenance"
)

const Kind = "assessment"

// Instance is the annotator-specific metadata.
type Instance struct {
	Provenance       provenance.Contract `json:"provenance"`
	AssessorKind     string              `json:"assessorType"`
	AssessorMetadata metadata.Contract   `json:"assessorMetadata"`

	assessorFactories []metadataFactory.Contract
}

// New is a factory function that returns an initialized Instance.
func New(provenance provenance.Contract, assessorMetadata metadata.Contract) *Instance {
	return &Instance{
		Provenance:       provenance,
		AssessorKind:     assessorMetadata.Kind(),
		AssessorMetadata: assessorMetadata,
	}
}

// Kind returns the type of concrete implementation.
func (*Instance) Kind() string {
	return Kind
}

// SetAssessorFactories provides for method injection of required factory to unmarshal metadata JSON.
func (i *Instance) SetAssessorFactories(assessorFactories []metadataFactory.Contract) {
	i.assessorFactories = assessorFactories
}

// UnmarshalJSON converts JSON into appropriate contract implementations.
func (i *Instance) UnmarshalJSON(data []byte) error {
	if i.assessorFactories == nil {
		return errors.New("uninitialized assessor factories")
	}

	type instance struct {
		Provenance       provenance.Contract `json:"provenance"`
		AssessorKind     string              `json:"assessorType"`
		AssessorMetadata json.RawMessage     `json:"assessorMetadata"`
	}

	var value instance
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	i.Provenance = value.Provenance
	i.AssessorKind = value.AssessorKind

	for f := range i.assessorFactories {
		if result := i.assessorFactories[f].Create(value.AssessorKind, value.AssessorMetadata); result != nil {
			i.AssessorMetadata = result
			break
		}
	}

	return nil
}
