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

	"github.com/project-alvarium/go-sdk/pkg/annotation/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/pki/verifier"
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/pki/verifier/verifypkcs1v15"
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/pki/verifier/verifytpmv2"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/reducer"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signpkcs1v15/hash"
	pkcsSignerMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signpkcs1v15/metadata"
	tpmSignerMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signtpmv2/metadata"
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
	instanceName := pkcsSignerMetadata.Kind + signerHash + reducerHash
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
	instanceName := tpmSignerMetadata.Kind + reducerHash
	if _, ok := f.instances[instanceName]; !ok {
		reducerHash := reducer.To(reducerHash)
		if reducerHash == nil {
			return nil
		}
		f.instances[instanceName] = verifytpmv2.New(reducerHash)
	}
	return f.instances[instanceName]
}

// Create returns a contract implementation based on the provided metadata.
func (f *Factory) Create(m metadata.Contract) verifier.Contract {
	f.m.Lock()
	defer f.m.Unlock()

	switch m.Kind() {
	case pkcsSignerMetadata.Kind:
		if m, ok := m.(*pkcsSignerMetadata.Success); ok {
			return f.SignPKCS1v15Verifier(m.SignerHash, m.ReducerHash)
		}
	case tpmSignerMetadata.Kind:
		if m, ok := m.(*tpmSignerMetadata.Success); ok {
			return f.SignTPMv2Verifier(m.ReducerHash)
		}
	}
	return nil
}
