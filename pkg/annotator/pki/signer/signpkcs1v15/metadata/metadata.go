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
	"crypto"

	"github.com/project-alvarium/go-sdk/pkg/annotator"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signpkcs1v15/hash"
)

// Success is the metadata specific to this signer implementation that results from an annotator event.
type Success struct {
	Kind        string `json:"type"`
	SignerHash  string `json:"signerHash"`
	ReducerHash string `json:"reducerHash"`
}

// New is a factory function that returns an initialized Success.
func New(signerHash crypto.Hash, reducerHash string) *Success {
	return &Success{
		Kind:        annotator.SuccessKind,
		SignerHash:  hash.FromSigner(signerHash),
		ReducerHash: reducerHash,
	}
}

// Failure defines the structure that encapsulates this signer's result.
type Failure struct {
	Kind         string `json:"type"`
	ErrorMessage string `json:"errorMessage"`
}

// NewFailure is a factory function that returns an initialized Failure.
func NewFailure(errorMessage string) *Failure {
	return &Failure{
		Kind:         annotator.FailureKind,
		ErrorMessage: errorMessage,
	}
}
