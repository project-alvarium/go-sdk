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

package stub

import "github.com/project-alvarium/go-sdk/pkg/test"

// Instance is a receiver that encapsulates required dependencies.
type Instance struct {
	kind  string
	Value interface{} `json:"value"`
}

// New is a factory function that returns an initialized Instance.
func New(kind string, value interface{}) *Instance {
	return &Instance{
		kind:  kind,
		Value: value,
	}
}

// NewNullObject is a factory function that returns an initialized instance.
func NewNullObject() *Instance {
	return New(test.FactoryRandomString(), test.FactoryRandomByteSlice())
}

// Kind returns the type of concrete implementation.
func (i *Instance) Kind() string {
	return i.kind
}
