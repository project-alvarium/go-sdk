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

const FailureResult = annotator.FailureKind

// Failure defines the structure that encapsulates this assessor's result.
type Failure struct {
	Result       string `json:"result"`
	ErrorMessage string `json:"errorMessage"`
}

// NewFailure is a factory function that returns an initialized Failure.
func NewFailure(errorMessage string) *Failure {
	return &Failure{
		Result:       FailureResult,
		ErrorMessage: errorMessage,
	}
}

// Kind returns the type of concrete implementation.
func (f *Failure) Kind() string {
	return Kind
}
