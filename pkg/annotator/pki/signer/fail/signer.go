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

package fail

import "github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/fail/metadata"

const name = "fail"

// signer is a receiver that encapsulates required dependencies.
type signer struct {
	SetUpCalled    bool
	TearDownCalled bool
}

// New is a factory function that returns signer.
func New() *signer {
	return &signer{
		SetUpCalled:    false,
		TearDownCalled: false,
	}
}

// SetUp is called once when the signer is instantiated.
func (s *signer) SetUp() {
	s.SetUpCalled = true
}

// TearDown is called once when signer is terminated.
func (s *signer) TearDown() {
	s.TearDownCalled = true
}

// PublicKey returns the associated public key.
func (*signer) PublicKey() []byte {
	return nil
}

// Sign returns a signature for the given data.
func (*signer) Sign(_, _ []byte) (identitySignature, dataSignature []byte) {
	return nil, nil
}

// Kind returns an implementation mnemonic; used in assessor when evaluating metadata from multiple implementations.
func (*signer) Kind() string {
	return Kind()
}

// Kind returns an implementation mnemonic; used in assessor when evaluating metadata from multiple implementations.
func Kind() string {
	return name
}

// Metadata returns implementation-specific metadata.
func (*signer) Metadata() interface{} {
	return metadata.New()
}
