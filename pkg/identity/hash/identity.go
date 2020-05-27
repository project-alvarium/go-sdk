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
	"encoding/json"
	"strings"
)

const Kind = "hash"

// Identity is a receiver that encapsulates required dependencies.
type Identity struct {
	Hash []byte `json:"hash"`
}

// New is a factory function that returns an initialized Identity.
func New(hash []byte) *Identity {
	return &Identity{
		Hash: hash,
	}
}

// Binary returns a unique key based on identity used within the SDK.
func (i *Identity) Binary() []byte {
	return i.Hash
}

// toPrintable converts []byte to printable string.
func toPrintable(b []byte) string {
	value, _ := json.Marshal(b)
	return strings.TrimPrefix(strings.TrimSuffix(string(value), "\""), "\"")
}

// Printable returns a unique key based on identity used within the SDK.
func (i *Identity) Printable() string {
	return toPrintable(i.Hash)
}

// Kind returns the type of concrete implementation.
func (*Identity) Kind() string {
	return Kind
}
