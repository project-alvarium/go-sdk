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

package hash

import (
	"github.com/project-alvarium/go-sdk/pkg/hashprovider"
	"github.com/project-alvarium/go-sdk/pkg/identity"
	identityHash "github.com/project-alvarium/go-sdk/pkg/identity/hash"
)

// provider is a receiver that encapsulates required dependencies.
type provider struct {
	hashProvider hashprovider.Contract
}

// New is a factory function that returns an initialized provider.
func New(hashProvider hashprovider.Contract) *provider {
	return &provider{
		hashProvider: hashProvider,
	}
}

// Derive converts data to an identity value.
func (p *provider) Derive(data []byte) identity.Contract {
	return identityHash.New(p.hashProvider.Derive(data))
}
