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

// assess implements an annotator that captures an assessor's assessment.
package assess

import (
	"fmt"

	envelope "github.com/project-alvarium/go-sdk/pkg/annotation/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotation/store"
	"github.com/project-alvarium/go-sdk/pkg/annotation/uniqueprovider"
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessment"
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor"
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotator/filter"
	"github.com/project-alvarium/go-sdk/pkg/annotator/provenance"
	"github.com/project-alvarium/go-sdk/pkg/identityprovider"
	"github.com/project-alvarium/go-sdk/pkg/status"
)

// annotator is a receiver that encapsulates required dependencies.
type annotator struct {
	provenance       provenance.Contract
	uniqueProvider   uniqueprovider.Contract
	identityProvider identityprovider.Contract
	store            store.Contract
	assessor         assessor.Contract
	filter           filter.Contract
}

// New is a factory function that returns an initialized annotator.
func New(
	provenance provenance.Contract,
	uniqueProvider uniqueprovider.Contract,
	identityProvider identityprovider.Contract,
	store store.Contract,
	assessor assessor.Contract,
	filter filter.Contract) *annotator {

	return &annotator{
		provenance:       provenance,
		uniqueProvider:   uniqueProvider,
		identityProvider: identityProvider,
		store:            store,
		assessor:         assessor,
		filter:           filter,
	}
}

// SetUp is called once when the signer is instantiated.
func (a *annotator) SetUp() {
	a.assessor.SetUp()
}

// TearDown is called once when signer is terminated.
func (a *annotator) TearDown() {
	a.assessor.TearDown()
}

// failureFindByIdentity returns annotations for failure case; separated to facilitate unit testing.
func (*annotator) failureFindByIdentity(result status.Value) *metadata.AssessFailure {
	return metadata.NewFailure(fmt.Sprintf("FindByIdentity returned %d", result))
}

// assess delegates to assessor's assess method, stores resulting assessment as annotation, and returns status.
func (a *annotator) assess(newData []byte) *status.Contract {
	var assessResult assessment.Contract

	id := a.identityProvider.Derive(newData)
	annotations, result := a.store.FindByIdentity(id)
	switch result {
	case status.Success:
		assessResult = a.assessor.Assess(a.filter.Do(annotations))
	default:
		assessResult = a.failureFindByIdentity(result)
	}

	m := envelope.New(
		a.uniqueProvider.Get(),
		id,
		nil,
		metadata.Kind,
		metadata.New(a.provenance, a.assessor.Kind(), assessResult),
	)
	result = a.store.Append(id, m)
	if result == status.NotFound {
		result = a.store.Create(id, m)
	}
	return status.New(a.provenance, result)
}

// Create evaluates newly-created data.
func (a *annotator) Create(data []byte) *status.Contract {
	return a.assess(data)
}

// Mutate evaluates mutated data.
func (a *annotator) Mutate(_, newData []byte) *status.Contract {
	return a.assess(newData)
}
