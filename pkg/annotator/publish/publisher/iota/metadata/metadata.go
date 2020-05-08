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

package metadata

import "github.com/project-alvarium/go-sdk/pkg/annotator"

const kind = "iota"

// Instance defines the structure that encapsulates this publisher's result.
type Instance struct {
	kind    string
	Result  string `json:"result"`
	Address string `json:"address"`
	Hash    string `json:"hash"`
	Tag     string `json:"tag"`
}

// New is a factory function that returns an initialized Instance.
func New(kind string, address, hash, tag string) *Instance {
	// TODO: Figure out what to do with kind here
	return &Instance{
		kind:    kind,
		Result:  annotator.SuccessKind,
		Address: address,
		Hash:    hash,
		Tag:     tag,
	}
}

// Kind returns the type of concrete implementation.
func (*Instance) Kind() string {
	return Kind()
}

// Kind returns the type of concrete implementation.
func Kind() string {
	return kind
}
