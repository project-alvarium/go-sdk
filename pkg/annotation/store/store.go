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

package store

import (
	"bytes"

	envelope "github.com/project-alvarium/go-sdk/pkg/annotation/metadata"
	"github.com/project-alvarium/go-sdk/pkg/identity"
	"github.com/project-alvarium/go-sdk/pkg/status"
	"github.com/project-alvarium/go-sdk/pkg/store"
)

// Persistence is a receiver that encapsulates required dependencies.
type Persistence struct {
	persistence store.Contract
}

// New is a factory function that returns persistence.
func New(persistence store.Contract) *Persistence {
	return &Persistence{
		persistence: persistence,
	}
}

// previous returns a non-nil non-id-matching identity if one exists.
func (*Persistence) previous(m envelope.Annotations, p, id identity.Contract) identity.Contract {
	idContract := identity.Factory(m.PreviousIdentityKind, m.PreviousIdentity)
	if p == nil && m.PreviousIdentity != nil && !bytes.Equal(idContract.Binary(), id.Binary()) {
		p = idContract
	}
	return p
}

// recursiveGet recursively traverses a chain of custody and returns its annotations.
func (p *Persistence) recursiveGet(id identity.Contract, annotations *[]*envelope.Annotations) status.Value {
	m, result := p.persistence.FindByIdentity(id)
	if result != status.Success {
		return result
	}

	var previousIdentity identity.Contract = nil
	for i := range m {
		if a, ok := m[i].(*envelope.Annotations); ok {
			previousIdentity = p.previous(*a, previousIdentity, id)
			*annotations = append(*annotations, a)
		}
	}
	if previousIdentity != nil {
		return p.recursiveGet(previousIdentity, annotations)
	}
	return status.Success
}

// FindByIdentity returns annotations and status corresponding to identity.
func (p *Persistence) FindByIdentity(id identity.Contract) ([]*envelope.Annotations, status.Value) {
	annotations := make([]*envelope.Annotations, 0)
	result := p.recursiveGet(id, &annotations)
	return annotations, result
}

// Create stores annotations corresponding to a new identity and returns status.
func (p *Persistence) Create(id identity.Contract, m *envelope.Annotations) status.Value {
	return p.persistence.Create(id, m)
}

// Append stores annotations corresponding to identity and returns status.
func (p *Persistence) Append(id identity.Contract, m *envelope.Annotations) status.Value {
	return p.persistence.Append(id, m)
}
