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

import "github.com/project-alvarium/go-sdk/pkg/annotation/metadata"

// instance is a receiver that encapsulates required dependencies.
type instance struct {
	result metadata.Contract
}

// New is a factory function which returns an initialized instance.
func New(result metadata.Contract) *instance {
	return &instance{
		result: result,
	}
}

// Send is called to send annotations to an IOTA Tangle.
func (i *instance) Send(_ string, _ uint64, _ uint64, _ []byte) metadata.Contract {
	return i.result
}
