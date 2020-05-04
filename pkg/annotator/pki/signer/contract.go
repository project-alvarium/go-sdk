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

package signer

const (
	PublicKeyType     = "PUBLIC KEY"
	RSAPrivateKeyType = "RSA PRIVATE KEY"
)

// Contract defines the signer abstraction.
type Contract interface {
	// SetUp is called once when the signer is instantiated.
	SetUp()

	// TearDown is called once when signer is terminated.
	TearDown()

	// PublicKey returns the associated public key.
	PublicKey() []byte

	// Sign returns a signature for the given identity and data.
	Sign(identity, data []byte) (identitySignature, dataSignature []byte)

	// Kind returns an implementation mnemonic; used in assessor when evaluating metadata from multiple implementations.
	Kind() string

	// Metadata returns implementation-specific metadata.
	Metadata() interface{}
}
