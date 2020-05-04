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

package verifypkcs1v15

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer"
	"github.com/project-alvarium/go-sdk/pkg/hashprovider"
)

// verifier is a receiver that encapsulates required dependencies.
type verifier struct {
	signerHash  crypto.Hash
	reducerHash hashprovider.Contract
}

// New is a factory function that returns verifier.
func New(signerHash crypto.Hash, reducerHash hashprovider.Contract) *verifier {
	return &verifier{
		signerHash:  signerHash,
		reducerHash: reducerHash,
	}
}

func (v *verifier) verify(data, signature, publicKey []byte) bool {
	var rsaPublicKey *rsa.PublicKey
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return false
	}
	if block.Type != signer.PublicKeyType {
		return false
	}

	rawKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false
	}
	switch t := rawKey.(type) {
	case *rsa.PublicKey:
		rsaPublicKey = t
	default:
		return false
	}

	return rsa.VerifyPKCS1v15(rsaPublicKey, v.signerHash, data[:], signature) == nil
}

// VerifyIdentity returns whether the given identity can be verified by the given signature.
func (v *verifier) VerifyIdentity(identity, signature, publicKey []byte) bool {
	return v.verify(v.reducerHash.Derive(identity), signature, publicKey)
}

// VerifyData returns whether the given data can be verified by the given signature.
func (v *verifier) VerifyData(data, signature, publicKey []byte) bool {
	return v.verify(v.reducerHash.Derive(data), signature, publicKey)
}
