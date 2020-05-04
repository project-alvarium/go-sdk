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

package verifier

import (
	"crypto"
	"testing"

	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/reducer"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signpkcs1v15"
	pkcsHash "github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signpkcs1v15/hash"
	pkcsMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signpkcs1v15/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signtpmv2"
	tpmMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signtpmv2/metadata"
	"github.com/project-alvarium/go-sdk/pkg/hashprovider/sha256"

	"github.com/stretchr/testify/assert"
)

// newSUT returns a new system under test.
func newSUT() *Factory {
	return New()
}

// TestFactory_Factory tests verifier.Factory.
func TestFactory_Factory(t *testing.T) {
	type testCase struct {
		name string
		test func(t *testing.T)
	}

	cases := []testCase{
		{
			name: "invalid (signer)",
			test: func(t *testing.T) {
				sut := newSUT()

				result := sut.Factory("invalid", nil)

				assert.Nil(t, result)
			},
		},
		{
			name: "invalid (pkcs1, signerHash)",
			test: func(t *testing.T) {
				sut := newSUT()

				result := sut.Factory(signpkcs1v15.Name, pkcsMetadata.New(0, sha256.New().Name()))

				assert.Nil(t, result)
			},
		},
		{
			name: "invalid (pkcs1, reducerHash)",
			test: func(t *testing.T) {
				sut := newSUT()

				result := sut.Factory(signpkcs1v15.Name, pkcsMetadata.New(crypto.SHA256, "unknown"))

				assert.Nil(t, result)
			},
		},
		{
			name: "valid (pkcs1)",
			test: func(t *testing.T) {
				for _, signerHash := range pkcsHash.SupportedSigner() {
					for _, reducerHash := range reducer.Supported() {
						sut := newSUT()

						result := sut.Factory(signpkcs1v15.Name, pkcsMetadata.New(signerHash, reducerHash.Name()))

						assert.NotNil(t, result)
						assert.Equal(
							t,
							result,
							sut.Factory(signpkcs1v15.Name, pkcsMetadata.New(signerHash, reducerHash.Name())),
						)
					}
				}
			},
		},
		{
			name: "invalid (tpm, reducerHash)",
			test: func(t *testing.T) {
				sut := newSUT()

				result := sut.Factory(signtpmv2.Name, tpmMetadata.New("unknown", nil))

				assert.Nil(t, result)
			},
		},
		{
			name: "valid (tpm)",
			test: func(t *testing.T) {
				for _, reducerHash := range reducer.Supported() {
					sut := newSUT()

					result := sut.Factory(signtpmv2.Name, tpmMetadata.New(reducerHash.Name(), nil))

					assert.NotNil(t, result)
					assert.Equal(t, result, sut.Factory(signtpmv2.Name, tpmMetadata.New(reducerHash.Name(), nil)))
				}
			},
		},
	}

	for i := range cases {
		t.Run(cases[i].name, cases[i].test)
	}
}
