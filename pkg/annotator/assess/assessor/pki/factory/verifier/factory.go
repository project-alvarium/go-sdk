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
	"sync"

	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/pki/verifier"
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/pki/verifier/verifypkcs1v15"
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/pki/verifier/verifytpm2"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/reducer"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signpkcs1v15"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signpkcs1v15/hash"
	pkcs "github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signpkcs1v15/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signtpmv2"
	tpm "github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signtpmv2/metadata"
)

// instances defines the map used to track verifier instances.
type instances map[string]verifier.Contract

// Factory is a receiver that encapsulates required dependencies.
type Factory struct {
	m         sync.Mutex
	instances instances
}

// New is a factory function that returns a factory.
func New() *Factory {
	return &Factory{
		m:         sync.Mutex{},
		instances: make(instances),
	}
}

// SignPKCS1v15Verifier returns a verifier.
func (f *Factory) SignPKCS1v15Verifier(signerHash, reducerHash string) verifier.Contract {
	instanceName := signpkcs1v15.Name + signerHash + reducerHash
	if _, ok := f.instances[instanceName]; !ok {
		signerHash := hash.ToSigner(signerHash)
		if signerHash == 0 {
			return nil
		}
		reducerHash := reducer.To(reducerHash)
		if reducerHash == nil {
			return nil
		}
		f.instances[instanceName] = verifypkcs1v15.New(signerHash, reducerHash)
	}
	return f.instances[instanceName]
}

// SignTPMv2Verifier returns a verifier.
func (f *Factory) SignTPMv2Verifier(reducerHash string) verifier.Contract {
	instanceName := signtpmv2.Name + reducerHash
	if _, ok := f.instances[instanceName]; !ok {
		reducerHash := reducer.To(reducerHash)
		if reducerHash == nil {
			return nil
		}
		f.instances[instanceName] = verifytpm2.New(reducerHash)
	}
	return f.instances[instanceName]
}

// Factory returns a contract implementation based on the provided metadata.
func (f *Factory) Factory(signer string, implementation interface{}) verifier.Contract {
	f.m.Lock()
	defer f.m.Unlock()

	switch signer {
	case signpkcs1v15.Name:
		if m, ok := implementation.(*pkcs.Success); ok {
			return f.SignPKCS1v15Verifier(m.SignerHash, m.ReducerHash)
		}
	case signtpmv2.Name:
		if m, ok := implementation.(*tpm.Success); ok {
			return f.SignTPMv2Verifier(m.ReducerHash)
		}
	}
	return nil
}
