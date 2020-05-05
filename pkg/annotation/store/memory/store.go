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

package memory

import (
	"bytes"
	"sync"

	"github.com/project-alvarium/go-sdk/pkg/annotation"
	"github.com/project-alvarium/go-sdk/pkg/identity"
	"github.com/project-alvarium/go-sdk/pkg/status"
)

// data defines the map used to provide generic storage.
type data map[string][]*annotation.Instance

// instance is a receiver that encapsulates required dependencies.
type instance struct {
	m    sync.Mutex
	data data
}

// New is a factory function that returns instance.
func New() *instance {
	return &instance{
		m:    sync.Mutex{},
		data: make(data),
	}
}

// recursiveFind recursively traverses a chain of custody and returns its annotations.
func (i *instance) recursiveFind(id identity.Contract, annotations *[]*annotation.Instance) status.Value {
	m, exists := i.data[id.Printable()]
	if !exists {
		return status.NotFound
	}

	var p identity.Contract = nil
	for i := range m {
		if p == nil && m[i].PreviousIdentity != nil && !bytes.Equal(m[i].PreviousIdentity.Binary(), id.Binary()) {
			p = m[i].PreviousIdentity
		}
		*annotations = append(*annotations, m[i])
	}
	if p != nil {
		i.recursiveFind(p, annotations)
	}

	return status.Success
}

// FindByIdentity returns annotations and status corresponding to identity.
func (i *instance) FindByIdentity(id identity.Contract) ([]*annotation.Instance, status.Value) {
	i.m.Lock()
	defer i.m.Unlock()

	annotations := make([]*annotation.Instance, 0)
	result := i.recursiveFind(id, &annotations)
	return annotations, result
}

// Create stores annotations corresponding to a new identity and returns status.
func (i *instance) Create(id identity.Contract, m *annotation.Instance) status.Value {
	i.m.Lock()
	defer i.m.Unlock()

	idAsString := id.Printable()
	_, exists := i.data[idAsString]
	if exists {
		return status.Exists
	}
	i.data[idAsString] = []*annotation.Instance{m}
	return status.Success
}

// Append stores annotations corresponding to identity and returns status.
func (i *instance) Append(id identity.Contract, m *annotation.Instance) status.Value {
	i.m.Lock()
	defer i.m.Unlock()

	idAsString := id.Printable()
	if _, exists := i.data[idAsString]; !exists {
		return status.NotFound
	}
	i.data[idAsString] = append(i.data[idAsString], m)
	return status.Success
}
