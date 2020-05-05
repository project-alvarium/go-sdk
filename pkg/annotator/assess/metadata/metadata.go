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
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessment"
	"github.com/project-alvarium/go-sdk/pkg/annotator/provenance"
)

const kind = "assess"

// Instance is the annotator-specific metadata.
type Instance struct {
	Provenance       provenance.Contract `json:"provenance"`
	AssessorKind     string              `json:"assessorType"`
	AssessorMetadata assessment.Contract `json:"assessorMetadata"`
}

// New is a factory function that returns an initialized Instance.
func New(
	provenance provenance.Contract,
	assessorKind string,
	assessorMetadata assessment.Contract) *Instance {

	return &Instance{
		Provenance:       provenance,
		AssessorKind:     assessorKind,
		AssessorMetadata: assessorMetadata,
	}
}

// Kind returns the type of concrete implementation.
func (*Instance) Kind() string {
	return Kind()
}

// Kind returns the type of concrete implementation.
func Kind() string {
	return kind
}
