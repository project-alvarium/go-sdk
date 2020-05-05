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

package pki

import (
	"github.com/project-alvarium/go-sdk/pkg/annotation"
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessment"
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/pki/factory"
	pkiAssessorMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/pki/metadata"
	pkiAnnotatorMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/pki/metadata"
)

const name = "verifier"

// assessor is a receiver that encapsulates required dependencies.
type assessor struct {
	factory factory.Contract
}

// New is a factory function that returns an initialized assessor.
func New(factory factory.Contract) *assessor {
	return &assessor{
		factory: factory,
	}
}

// SetUp is called once when the assessor is instantiated.
func (*assessor) SetUp() {}

// TearDown is called once when assessor is terminated.
func (*assessor) TearDown() {}

// Assess accepts data and returns associated assessments.
func (a *assessor) Assess(annotations []*annotation.Instance) assessment.Contract {
	uniques := make([]string, 0)
	for i := range annotations {
		if annotations[i].MetadataKind != pkiAnnotatorMetadata.Kind() {
			continue
		}

		m := annotations[i].Metadata.(*pkiAnnotatorMetadata.Instance)
		v := a.factory.Create(m.SignerMetadata)
		if v == nil || v.VerifyIdentity(annotations[i].CurrentIdentity.Binary(), m.IdentitySignature, m.PublicKey) == false {
			return pkiAssessorMetadata.New(false, []string{annotations[i].Unique})
		}
		uniques = append(uniques, annotations[i].Unique)
	}
	return pkiAssessorMetadata.New(true, uniques)
}

// Kind returns an implementation mnemonic.
func (*assessor) Kind() string {
	return name
}
