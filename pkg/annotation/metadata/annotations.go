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
	"github.com/project-alvarium/go-sdk/internal/pkg/datetime"
	"github.com/project-alvarium/go-sdk/pkg/identity"
)

// Annotations is the standard metadata returned by all Annotate() methods.
type Annotations struct {
	Unique               string      `json:"unique"`
	Created              string      `json:"created"`
	CurrentIdentityKind  string      `json:"identityCurrentType"`
	CurrentIdentity      interface{} `json:"identityCurrent"`
	PreviousIdentityKind string      `json:"identityPreviousType"`
	PreviousIdentity     interface{} `json:"identityPrevious"`
	MetadataKind         string      `json:"metadataType"`
	Metadata             interface{} `json:"metadata"`
}

// New is a factory function that returns an initialized Annotations.
func New(
	unique string,
	currentIdentity identity.Contract,
	previousIdentity identity.Contract,
	metadataKind string,
	metadata interface{}) *Annotations {

	return &Annotations{
		Unique:               unique,
		Created:              datetime.Created(),
		CurrentIdentityKind:  currentIdentity.Kind(),
		CurrentIdentity:      currentIdentity,
		PreviousIdentityKind: currentIdentity.Kind(),
		PreviousIdentity:     previousIdentity,
		MetadataKind:         metadataKind,
		Metadata:             metadata,
	}
}
