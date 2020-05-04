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

package annotation

import (
	"github.com/project-alvarium/go-sdk/internal/pkg/datetime"
	"github.com/project-alvarium/go-sdk/pkg/annotation/metadata"
	"github.com/project-alvarium/go-sdk/pkg/identity"
)

// Instance is the standard metadata returned by all Annotate() methods.
type Instance struct {
	Unique               string            `json:"unique"`
	Created              string            `json:"created"`
	CurrentIdentityKind  string            `json:"identityCurrentType"`
	CurrentIdentity      identity.Contract `json:"identityCurrent"`
	PreviousIdentityKind string            `json:"identityPreviousType"`
	PreviousIdentity     identity.Contract `json:"identityPrevious"`
	MetadataKind         string            `json:"metadataType"`
	Metadata             metadata.Contract `json:"metadata"`
}

// New is a factory function that returns an initialized Instance.
func New(
	unique string,
	currentIdentity identity.Contract,
	previousIdentity identity.Contract,
	metadata metadata.Contract) *Instance {

	currentIdentityKind := currentIdentity.Kind()
	previousIdentityKind := currentIdentityKind
	if previousIdentity != nil {
		previousIdentityKind = previousIdentity.Kind()
	}

	return &Instance{
		Unique:               unique,
		Created:              datetime.Created(),
		CurrentIdentityKind:  currentIdentityKind,
		CurrentIdentity:      currentIdentity,
		PreviousIdentityKind: previousIdentityKind,
		PreviousIdentity:     previousIdentity,
		MetadataKind:         metadata.Kind(),
		Metadata:             metadata,
	}
}
