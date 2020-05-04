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

package assessor

import (
	"github.com/project-alvarium/go-sdk/pkg/annotation"
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessment"
)

// Contract defines the assessor abstraction.
type Contract interface {
	// SetUp is called once when the assessor is instantiated.
	SetUp()

	// TearDown is called once when assessor is terminated.
	TearDown()

	// Assess accepts data and returns associated assessments.
	Assess(annotations []*annotation.Instance) assessment.Contract

	// Kind returns an implementation mnemonic.
	Kind() string
}
