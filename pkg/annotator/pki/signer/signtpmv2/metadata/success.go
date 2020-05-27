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

import "github.com/project-alvarium/go-sdk/pkg/annotator"

const SuccessResult = annotator.SuccessKind

// CapabilityProperties defines type to contain the TPM capability properties.
type CapabilityProperties map[string]string

// Success is the metadata specific to this signer implementation that results from an annotator event.
type Success struct {
	Result               string               `json:"result"`
	ReducerHash          string               `json:"reducerHash"`
	CapabilityProperties CapabilityProperties `json:"capabilityProperties"`
}

// NewSuccess is a factory function that returns an initialized Success.
func NewSuccess(reducerHash string, capabilityProperties CapabilityProperties) *Success {
	return &Success{
		Result:               SuccessResult,
		ReducerHash:          reducerHash,
		CapabilityProperties: capabilityProperties,
	}
}

// Kind returns the type of concrete implementation.
func (*Success) Kind() string {
	return Kind
}
