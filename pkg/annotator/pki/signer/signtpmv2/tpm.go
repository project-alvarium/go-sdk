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

package signtpmv2

import (
	"bytes"
	"io"
	"sync"

	"github.com/project-alvarium/go-sdk/pkg/annotation/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signtpmv2/factory"
	tpmMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signtpmv2/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signtpmv2/provisioner"
	"github.com/project-alvarium/go-sdk/pkg/hashprovider"

	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpmutil"
)

const Name = "tpm2.0"

// RequestedCapabilityProperties identifies the TPM capability properties requested to be included in annotations.
type RequestedCapabilityProperties map[string]tpm2.TPMProp

// signer is a receiver that encapsulates required dependencies.
type signer struct {
	hashProvider                  hashprovider.Contract
	publicKey                     []byte
	rwc                           io.ReadWriteCloser
	handle                        tpmutil.Handle
	scheme                        *tpm2.SigScheme
	path                          string
	RequestedCapabilityProperties RequestedCapabilityProperties
	capabilityProperties          tpmMetadata.CapabilityProperties
	m                             sync.Mutex
	signerError                   error
}

// NewWithRWC return signer with rwc initialized.
func NewWithRWC(
	hashProvider hashprovider.Contract,
	publicKey []byte,
	handle tpmutil.Handle,
	path string,
	requestedCapabilityProperties RequestedCapabilityProperties,
	rwc io.ReadWriteCloser) *signer {

	scheme := &tpm2.SigScheme{
		Alg:  provisioner.Algorithm,
		Hash: provisioner.Hash,
	}
	return &signer{
		hashProvider:                  hashProvider,
		publicKey:                     publicKey,
		rwc:                           rwc,
		handle:                        handle,
		scheme:                        scheme,
		path:                          path,
		RequestedCapabilityProperties: requestedCapabilityProperties,
		m:                             sync.Mutex{},
	}
}

// New is a factory function that returns signer.
func New(
	hashProvider hashprovider.Contract,
	publicKey []byte,
	handle tpmutil.Handle,
	path string,
	requestedCapabilityProperties RequestedCapabilityProperties) *signer {

	return NewWithRWC(hashProvider, publicKey, handle, path, requestedCapabilityProperties, nil)
}

// getCapabilityProperties returns initialized TpmContext.
func (s *signer) getCapabilityProperties() tpmMetadata.CapabilityProperties {
	getProperty := func(requestedProperty tpm2.TPMProp) string {
		var data []interface{}
		var err error

		data, _, err = tpm2.GetCapability(s.rwc, tpm2.CapabilityTPMProperties, 1, uint32(requestedProperty))
		if err != nil {
			return ""
		}

		var b []byte
		if b, err = tpmutil.Pack(data[0].(tpm2.TaggedProperty).Value); err != nil {
			return ""
		}
		return string(bytes.Trim(b, "\x00"))
	}

	contexts := make(tpmMetadata.CapabilityProperties)
	for i := range s.RequestedCapabilityProperties {
		contexts[i] = getProperty(s.RequestedCapabilityProperties[i])
	}
	return contexts
}

// SetUp is called once when the signer is instantiated.
func (s *signer) SetUp() {
	if s.rwc == nil {
		if rwc, err := factory.TPM(s.path); err == nil {
			s.rwc = rwc
		}
	}
	s.capabilityProperties = s.getCapabilityProperties()

}

// TearDown is called once when signer is terminated.
func (s *signer) TearDown() {
	if s.rwc != nil {
		_ = s.rwc.Close()
	}
}

// PublicKey returns the associated public key.o
func (s *signer) PublicKey() []byte {
	return s.publicKey
}

// sign implements the common signature implementation.
func (s *signer) sign(data []byte) []byte {
	s.m.Lock()
	defer s.m.Unlock()

	signature, err := tpm2.Sign(s.rwc, s.handle, "", data, s.scheme)
	if err != nil {
		s.signerError = err
		return nil
	}
	return signature.RSA.Signature
}

// Sign returns a signature for the given identity and data.
func (s *signer) Sign(identity, data []byte) (identitySignature, dataSignature []byte) {
	return s.sign(s.hashProvider.Derive(identity)), s.sign(s.hashProvider.Derive(data))
}

// Metadata returns implementation-specific metadata.
func (s *signer) Metadata() metadata.Contract {
	if s.signerError != nil {
		return tpmMetadata.NewFailure(Name, s.signerError.Error())
	}
	return tpmMetadata.NewSuccess(Name, s.hashProvider.Name(), s.capabilityProperties)
}
