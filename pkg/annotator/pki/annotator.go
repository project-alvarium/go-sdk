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

// pki implements a PKI-signed annotator.
package pki

import (
	"bytes"

	envelope "github.com/project-alvarium/go-sdk/pkg/annotation/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotation/store"
	"github.com/project-alvarium/go-sdk/pkg/annotation/uniqueprovider"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer"
	"github.com/project-alvarium/go-sdk/pkg/annotator/provenance"
	"github.com/project-alvarium/go-sdk/pkg/identity"
	"github.com/project-alvarium/go-sdk/pkg/identityprovider"
	"github.com/project-alvarium/go-sdk/pkg/status"
)

// annotator is a receiver that encapsulates required dependencies.
type annotator struct {
	provenance       provenance.Contract
	uniqueProvider   uniqueprovider.Contract
	identityProvider identityprovider.Contract
	store            store.Contract
	signer           signer.Contract
}

// New is a factory function that returns an initialized annotator.
func New(
	provenance provenance.Contract,
	uniqueProvider uniqueprovider.Contract,
	identityProvider identityprovider.Contract,
	store store.Contract,
	signer signer.Contract) *annotator {

	return &annotator{
		provenance:       provenance,
		uniqueProvider:   uniqueProvider,
		identityProvider: identityProvider,
		store:            store,
		signer:           signer,
	}
}

// metadata is a private factory function that delegates to metadata.New() and returns Annotate.
func (a *annotator) metadata(
	identity identity.Contract,
	previousIdentity identity.Contract,
	identitySignature []byte,
	dataSignature []byte) *envelope.Annotations {

	return envelope.New(
		a.uniqueProvider.Get(),
		identity,
		previousIdentity,
		metadata.Kind,
		metadata.New(
			a.provenance,
			identitySignature,
			dataSignature,
			a.signer.PublicKey(),
			a.signer.Kind(),
			a.signer.Metadata(),
		),
	)
}

// SetUp is called once when the signer is instantiated.
func (a *annotator) SetUp() {
	a.signer.SetUp()
}

// TearDown is called once when signer is terminated.
func (a *annotator) TearDown() {
	a.signer.TearDown()
}

// sign evaluates data and returns metadata.
func (a *annotator) sign(oldIdentity identity.Contract, data []byte) (identity.Contract, *envelope.Annotations) {
	id := a.identityProvider.Derive(data)
	identitySignature, dataSignature := a.signer.Sign(id.Binary(), data)
	return id, a.metadata(id, oldIdentity, identitySignature, dataSignature)
}

// Create evaluates newly-created data.
func (a *annotator) Create(data []byte) *status.Contract {
	id, m := a.sign(nil, data)
	return status.New(a.provenance, a.store.Create(id, m))
}

// Mutate evaluates mutated data.
func (a *annotator) Mutate(oldData, newData []byte) *status.Contract {
	oldDataIdentity := a.identityProvider.Derive(oldData)
	newDataIdentity, m := a.sign(oldDataIdentity, newData)

	if !bytes.Equal(oldDataIdentity.Binary(), newDataIdentity.Binary()) {
		return status.New(a.provenance, a.store.Create(newDataIdentity, m))
	}
	return status.New(a.provenance, a.store.Append(newDataIdentity, m))
}
