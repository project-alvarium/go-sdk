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

package signpkcs1v15

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/project-alvarium/go-sdk/pkg/annotation/metadata"
	pkiSigner "github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer"
	pkcsSignerMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signpkcs1v15/metadata"
	"github.com/project-alvarium/go-sdk/pkg/hashprovider"
)

// signer is a receiver that encapsulates required dependencies.
type signer struct {
	hash         crypto.Hash
	privateKey   *rsa.PrivateKey
	publicKey    []byte
	hashProvider hashprovider.Contract
	signerError  error
}

// New is a factory function that returns signer.
func New(hash crypto.Hash, privateKey, publicKey []byte, hashProvider hashprovider.Contract) *signer {
	block, _ := pem.Decode(privateKey)
	if block == nil || block.Type != pkiSigner.RSAPrivateKeyType {
		return nil
	}

	rsaPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil
	}

	return &signer{
		hash:         hash,
		privateKey:   rsaPrivateKey,
		publicKey:    publicKey,
		hashProvider: hashProvider,
	}
}

// SetUp is called once when the signer is instantiated.
func (*signer) SetUp() {}

// TearDown is called once when signer is terminated.
func (*signer) TearDown() {}

// PublicKey returns the associated public key.
func (s *signer) PublicKey() []byte {
	return s.publicKey
}

// sign implements the common signature implementation.
func (s *signer) sign(hash []byte) (signature []byte) {
	signature, s.signerError = rsa.SignPKCS1v15(rand.Reader, s.privateKey, s.hash, hash[:])
	return signature
}

// Sign returns a signature for the given identity and data.
func (s *signer) Sign(identity, data []byte) (identitySignature, dataSignature []byte) {
	return s.sign(s.hashProvider.Derive(identity)), s.sign(s.hashProvider.Derive(data))
}

// Metadata returns implementation-specific metadata.
func (s *signer) Metadata() metadata.Contract {
	if s.signerError != nil {
		return pkcsSignerMetadata.NewFailure(s.signerError.Error())
	}
	return pkcsSignerMetadata.NewSuccess(s.hash, s.hashProvider.Kind())
}
