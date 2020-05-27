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

const Kind = "pki"

// Instance is the annotator-specific metadata.
type Instance struct {
	Provenance        provenance.Contract `json:"provenance"`
	IdentitySignature []byte              `json:"identitySignature"`
	DataSignature     []byte              `json:"dataSignature"`
	PublicKey         []byte              `json:"publicKey"`
	SignerKind        string              `json:"signerType"`
	SignerMetadata    metadata.Contract   `json:"signerMetadata"`

	signerFactories []metadataFactory.Contract
}

// New is a factory function that returns an initialized Instance.
func New(
	provenance provenance.Contract,
	identitySignature []byte,
	dataSignature []byte,
	publicKey []byte,
	signerMetadata metadata.Contract) *Instance {

	return &Instance{
		Provenance:        provenance,
		IdentitySignature: identitySignature,
		DataSignature:     dataSignature,
		PublicKey:         publicKey,
		SignerKind:        signerMetadata.Kind(),
		SignerMetadata:    signerMetadata,
	}
}

// Kind returns the type of concrete implementation.
func (*Instance) Kind() string {
	return Kind
}

// SetSignerFactories provides for method injection of required factory to unmarshal metadata JSON.
func (i *Instance) SetSignerFactories(signerFactories []metadataFactory.Contract) {
	i.signerFactories = signerFactories
}

// UnmarshalJSON converts JSON into appropriate contract implementations.
func (i *Instance) UnmarshalJSON(data []byte) error {
	if i.signerFactories == nil {
		return errors.New("uninitialized signer factories")
	}

	type instance struct {
		Provenance        provenance.Contract `json:"provenance"`
		IdentitySignature []byte              `json:"identitySignature"`
		DataSignature     []byte              `json:"dataSignature"`
		PublicKey         []byte              `json:"publicKey"`
		SignerKind        string              `json:"signerType"`
		SignerMetadata    json.RawMessage     `json:"signerMetadata"`
	}

	var value instance
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	i.Provenance = value.Provenance
	i.IdentitySignature = value.IdentitySignature
	i.DataSignature = value.DataSignature
	i.PublicKey = value.PublicKey
	i.SignerKind = value.SignerKind

	for f := range i.signerFactories {
		if result := i.signerFactories[f].Create(value.SignerKind, value.SignerMetadata); result != nil {
			i.SignerMetadata = result
			break
		}
	}

	return nil
}
