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
	envelope "github.com/project-alvarium/go-sdk/pkg/annotation/metadata"
	"github.com/project-alvarium/go-sdk/pkg/identity"
	"github.com/project-alvarium/go-sdk/pkg/status"
)

// Contract defines the pki/store abstraction.
type Contract interface {
	// FindByIdentity returns annotations and status corresponding to identity.
	FindByIdentity(id identity.Contract) ([]*envelope.Annotations, status.Value)

	// Create stores annotations corresponding to a new identity and returns status.
	Create(id identity.Contract, m *envelope.Annotations) status.Value

	// Append stores annotations corresponding to identity and returns status.
	Append(id identity.Contract, m *envelope.Annotations) status.Value
}
