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

package store

import (
	"github.com/project-alvarium/go-sdk/pkg/identity"
	"github.com/project-alvarium/go-sdk/pkg/status"
)

// Contract defines the store abstraction.
type Contract interface {
	// FindByIdentity returns model and status corresponding to identity.
	FindByIdentity(id identity.Contract) ([]interface{}, status.Value)

	// Create stores model corresponding to a new identity and returns status.
	Create(id identity.Contract, model interface{}) status.Value

	// Append stores model corresponding to identity and returns status.
	Append(id identity.Contract, model interface{}) status.Value
}
