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
	pkcsSignerMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signpkcs1v15/metadata"
)

// instance is a receiver that encapsulates required dependencies.
type instance struct{}

// New is a factory function that returns an initialized instance.
func New() *instance {
	return &instance{}
}

// Create returns a contract implementation based on the provided metadata.
func (i *instance) Create(kind string, data json.RawMessage) metadata.Contract {
	if kind != pkcsSignerMetadata.Kind {
		return nil
	}

	type instance struct {
		Result string `json:"result"`
	}

	var value instance
	if err := json.Unmarshal(data, &value); err != nil {
		return nil
	}

	switch value.Result {
	case pkcsSignerMetadata.FailureResult:
		var concrete pkcsSignerMetadata.Failure
		if err := json.Unmarshal(data, &concrete); err == nil {
			return &concrete
		}
	case pkcsSignerMetadata.SuccessResult:
		var concrete pkcsSignerMetadata.Success
		if err := json.Unmarshal(data, &concrete); err == nil {
			return &concrete
		}
	}
	return nil
}
