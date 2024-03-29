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

package annotator

import "github.com/project-alvarium/go-sdk/pkg/status"

const (
	SuccessKind = "success"
	FailureKind = "failure"
)

// Contract defines the annotator abstraction.
type Contract interface {
	// SetUp is called once when the annotator is instantiated.
	SetUp()

	// TearDown is called once when annotator is terminated.
	TearDown()

	// Create evaluates newly-created data.
	Create(data []byte) *status.Contract

	// Mutate evaluates mutated data.
	Mutate(oldData, newData []byte) *status.Contract
}
