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
	"encoding/json"
	"errors"

	"github.com/project-alvarium/go-sdk/internal/pkg/datetime"
	"github.com/project-alvarium/go-sdk/pkg/annotation/metadata"
	metadataFactory "github.com/project-alvarium/go-sdk/pkg/annotation/metadata/factory"
	"github.com/project-alvarium/go-sdk/pkg/identity"
	identityFactory "github.com/project-alvarium/go-sdk/pkg/identity/factory"
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

	identityFactory identityFactory.Contract
	metadataFactory metadataFactory.Contract
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

// SetIdentityFactory provides for method injection of required factory to unmarshal identity JSON.
func (i *Instance) SetIdentityFactory(identityFactory identityFactory.Contract) {
	i.identityFactory = identityFactory
}

// SetMetadataFactory provides for method injection of required factory to unmarshal metadata JSON.
func (i *Instance) SetMetadataFactory(metadataFactory metadataFactory.Contract) {
	i.metadataFactory = metadataFactory
}

// UnmarshalJSON converts JSON into appropriate contract implementations.
func (i *Instance) UnmarshalJSON(data []byte) error {
	if i.identityFactory == nil {
		return errors.New("identityFactory not set")
	}
	if i.metadataFactory == nil {
		return errors.New("metadataFactory not set")
	}

	type instance struct {
		Unique               string          `json:"unique"`
		Created              string          `json:"created"`
		CurrentIdentityKind  string          `json:"identityCurrentType"`
		CurrentIdentity      json.RawMessage `json:"identityCurrent"`
		PreviousIdentityKind string          `json:"identityPreviousType"`
		PreviousIdentity     json.RawMessage `json:"identityPrevious"`
		MetadataKind         string          `json:"metadataType"`
		Metadata             json.RawMessage `json:"metadata"`
	}

	var value instance
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	i.Unique = value.Unique
	i.Created = value.Created
	i.CurrentIdentityKind = value.CurrentIdentityKind
	i.CurrentIdentity = i.identityFactory.Create(value.CurrentIdentityKind, value.CurrentIdentity)
	i.PreviousIdentityKind = value.PreviousIdentityKind
	i.PreviousIdentity = i.identityFactory.Create(value.PreviousIdentityKind, value.PreviousIdentity)
	i.MetadataKind = value.MetadataKind
	i.Metadata = i.metadataFactory.Create(value.MetadataKind, value.Metadata)

	return nil
}
