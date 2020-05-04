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

import (
	"github.com/project-alvarium/go-sdk/pkg/annotation"
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessment"
)

// annotator is a receiver that encapsulates required dependencies.
type annotator struct {
	kind           string
	assessment     assessment.Contract
	SetUpCalled    bool
	TearDownCalled bool
}

// New is a factory function that returns an initialized annotator.
func New(kind string, assessment assessment.Contract) *annotator {
	return &annotator{
		kind:       kind,
		assessment: assessment,
	}
}

// SetUp is called once when the annotator is instantiated.
func (a *annotator) SetUp() {
	a.SetUpCalled = true
}

// TearDown is called once when annotator is terminated.
func (a *annotator) TearDown() {
	a.TearDownCalled = true
}

// Assess accepts data and returns associated assessments.
func (a *annotator) Assess(_ []*annotation.Instance) assessment.Contract {
	return a.assessment
}

// Kind returns an implementation mnemonic.
func (a *annotator) Kind() string {
	return a.kind
}
