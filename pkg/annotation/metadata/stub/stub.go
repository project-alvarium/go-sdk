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

// instance is a receiver that encapsulates required dependencies.
type instance struct {
	kind  string
	value interface{}
}

// New is a factory function that returns an initialized instance.
func New(kind string, value interface{}) *instance {
	return &instance{
		kind:  kind,
		value: value,
	}
}

// NewNullObject is a factory function that returns an initialized instance.
func NewNullObject() *instance {
	return New(test.FactoryRandomString(), test.FactoryRandomByteSlice())
}

// Kind returns the type of concrete implementation.
func (i *instance) Kind() string {
	return i.kind
}
