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
	"sync"

	"github.com/project-alvarium/go-sdk/pkg/identity"
	"github.com/project-alvarium/go-sdk/pkg/status"
)

// data defines the map used to provide generic storage.
type data map[string][]interface{}

// Persistence is a receiver that encapsulates required dependencies.
type Persistence struct {
	m    sync.Mutex
	data data
}

// New is a factory function that returns persistence.
func New() *Persistence {
	return &Persistence{
		m:    sync.Mutex{},
		data: make(data),
	}
}

// FindByIdentity returns model and status corresponding to identity.
func (p *Persistence) FindByIdentity(id identity.Contract) ([]interface{}, status.Value) {
	p.m.Lock()
	defer p.m.Unlock()
	if models, exists := p.data[id.Printable()]; exists {
		return models, status.Success
	}
	return nil, status.NotFound
}

// Create stores model corresponding to a new identity and returns status.
func (p *Persistence) Create(id identity.Contract, model interface{}) status.Value {
	p.m.Lock()
	defer p.m.Unlock()

	idAsString := id.Printable()
	_, exists := p.data[idAsString]
	if exists {
		return status.Exists
	}
	p.data[idAsString] = []interface{}{model}
	return status.Success
}

// Append stores model corresponding to identity and returns status.
func (p *Persistence) Append(id identity.Contract, model interface{}) status.Value {
	p.m.Lock()
	defer p.m.Unlock()

	idAsString := id.Printable()
	models, exists := p.data[idAsString]
	if !exists {
		return status.NotFound
	}
	p.data[idAsString] = append(models, model)
	return status.Success
}
