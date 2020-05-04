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
	envelope "github.com/project-alvarium/go-sdk/pkg/annotation/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/published"
)

// annotator is a receiver that encapsulates required dependencies.
type publisher struct {
	kind           string
	published      published.Contract
	SetUpCalled    bool
	TearDownCalled bool
}

// New is a factory function that returns an initialized annotator.
func New(kind string, published published.Contract) *publisher {
	return &publisher{
		kind:      kind,
		published: published,
	}
}

// SetUp is called once when the annotator is instantiated.
func (p *publisher) SetUp() {
	p.SetUpCalled = true
}

// TearDown is called once when annotator is terminated.
func (p *publisher) TearDown() {
	p.TearDownCalled = true
}

// Assess accepts data and returns associated assessments.
func (p *publisher) Publish([]*envelope.Annotations) published.Contract {
	return p.published
}

// Kind returns an implementation mnemonic.
func (p *publisher) Kind() string {
	return p.kind
}
