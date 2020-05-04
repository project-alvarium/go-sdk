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

package provisioner

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"io"

	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer"

	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpmutil"
)

const (
	Path          = "/dev/tpm0"
	InvalidHandle = tpmutil.Handle(0x00000000)
	Algorithm     = tpm2.AlgRSASSA
	Hash          = tpm2.AlgSHA256
	CryptoHash    = crypto.SHA256
)

// Flush removes any loaded handles in the tpm. This prevent out-of-memory errors.
func Flush(rwc io.ReadWriteCloser, handle tpmutil.Handle) {
	_ = tpm2.FlushContext(rwc, handle)
}

// FlushAndClose removes any loaded handles in the tpm and closes it. This prevent out-of-memory errors.
func FlushAndClose(rwc io.ReadWriteCloser, handle tpmutil.Handle) {
	Flush(rwc, handle)
	_ = rwc.Close()
}

// MarshalPublicKey converts a crypto.PublicKey value to its pem encoding.
func MarshalPublicKey(key crypto.PublicKey) []byte {
	publicKeyBlob, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return nil
	}
	return pem.EncodeToMemory(&pem.Block{Type: signer.PublicKeyType, Bytes: publicKeyBlob})
}

// GenerateNewKeyPair creates a new primary key inside a tpm and returns the key handle and its public key.
// If GenerateNewKeyPair encounters an error during primary key creation, it returns an invalid tpm handle,
// nil public key and error.
func GenerateNewKeyPair(rwc io.ReadWriteCloser) (tpmutil.Handle, crypto.PublicKey, error) {
	template := tpm2.Public{
		Type:       tpm2.AlgRSA,
		NameAlg:    Hash,
		Attributes: tpm2.FlagSignerDefault & ^tpm2.FlagRestricted,
		RSAParameters: &tpm2.RSAParams{
			Sign: &tpm2.SigScheme{
				Alg:  Algorithm,
				Hash: Hash,
			},
			KeyBits: 2048,
		},
	}
	handle, publicKey, err := tpm2.CreatePrimary(rwc, tpm2.HandleOwner, tpm2.PCRSelection{}, "", "", template)
	if err != nil {
		return InvalidHandle, nil, err
	}
	return handle, publicKey, nil
}
