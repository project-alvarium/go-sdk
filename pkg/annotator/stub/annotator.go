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

import "github.com/project-alvarium/go-sdk/pkg/status"

// annotator is a receiver that encapsulates required dependencies.
type annotator struct {
	result         *status.Contract
	SetUpCalled    int
	TearDownCalled int
}

// NewWithResult is a factory function that returns an initialized annotator.
func NewWithResult(result *status.Contract) *annotator {
	return &annotator{
		result:         result,
		SetUpCalled:    0,
		TearDownCalled: 0,
	}
}

// New is a factory function that returns an initialized annotator.
func New() *annotator {
	return NewWithResult(nil)
}

// SetUp is called once when the annotator is instantiated.
func (a *annotator) SetUp() {
	a.SetUpCalled++
}

// TearDown is called once when annotator is terminated.
func (a *annotator) TearDown() {
	a.TearDownCalled++
}

// Create evaluates newly-created data.
func (a *annotator) Create(_ []byte) *status.Contract {
	return a.result
}

// Mutate evaluates mutated data.
func (a *annotator) Mutate(_, _ []byte) *status.Contract {
	return a.result
}
