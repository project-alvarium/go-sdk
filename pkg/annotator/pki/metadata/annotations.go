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

import "github.com/project-alvarium/go-sdk/pkg/annotator/provenance"

const Kind = "PKI"

// Annotations is the annotator-specific metadata.
type Annotations struct {
	Provenance        provenance.Contract `json:"provenance"`
	IdentitySignature []byte              `json:"identitySignature"`
	DataSignature     []byte              `json:"dataSignature"`
	PublicKey         []byte              `json:"publicKey"`
	SignerKind        string              `json:"signerType"`
	SignerMetadata    interface{}         `json:"signerMetadata"`
}

// New is a factory function that returns an initialized Annotations.
func New(
	provenance provenance.Contract,
	identitySignature []byte,
	dataSignature []byte,
	publicKey []byte,
	signerKind string,
	signerMetadata interface{}) *Annotations {

	return &Annotations{
		Provenance:        provenance,
		IdentitySignature: identitySignature,
		DataSignature:     dataSignature,
		PublicKey:         publicKey,
		SignerKind:        signerKind,
		SignerMetadata:    signerMetadata,
	}
}
