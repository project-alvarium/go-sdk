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

const SuccessResult = annotator.SuccessKind

// Success defines the structure that encapsulates this publisher's result.
type Success struct {
	Result string `json:"result"`
	CID    string `json:"cid"`
}

// NewSuccess is a factory function that returns an initialized Success.
func NewSuccess(cid string) *Success {
	return &Success{
		Result: SuccessResult,
		CID:    cid,
	}
}

// Kind returns the type of concrete implementation.
func (i *Success) Kind() string {
	return Kind
}