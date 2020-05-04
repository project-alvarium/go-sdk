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

// publish implements an annotator that captures an publisher's assessment.
package publish

import (
	"fmt"

	"github.com/project-alvarium/go-sdk/pkg/annotation"
	"github.com/project-alvarium/go-sdk/pkg/annotation/store"
	"github.com/project-alvarium/go-sdk/pkg/annotation/uniqueprovider"
	"github.com/project-alvarium/go-sdk/pkg/annotator/filter"
	"github.com/project-alvarium/go-sdk/pkg/annotator/provenance"
	publishMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/publish/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/published"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher"
	"github.com/project-alvarium/go-sdk/pkg/identityprovider"
	"github.com/project-alvarium/go-sdk/pkg/status"
)

// annotator is a receiver that encapsulates required dependencies.
type annotator struct {
	provenance       provenance.Contract
	uniqueProvider   uniqueprovider.Contract
	identityProvider identityprovider.Contract
	store            store.Contract
	publisher        publisher.Contract
	filter           filter.Contract
}

// New is a factory function that returns an initialized annotator.
func New(
	provenance provenance.Contract,
	uniqueProvider uniqueprovider.Contract,
	identityProvider identityprovider.Contract,
	store store.Contract,
	publisher publisher.Contract,
	filter filter.Contract) *annotator {

	return &annotator{
		provenance:       provenance,
		uniqueProvider:   uniqueProvider,
		identityProvider: identityProvider,
		store:            store,
		publisher:        publisher,
		filter:           filter,
	}
}

// SetUp is called once when the signer is instantiated.
func (a *annotator) SetUp() {
	a.publisher.SetUp()
}

// TearDown is called once when signer is terminated.
func (a *annotator) TearDown() {
	a.publisher.TearDown()
}

// failureFindByIdentity returns annotations for failure case; separated to facilitate unit testing.
func (*annotator) failureFindByIdentity(result status.Value) *publishMetadata.PublishedFailure {
	return publishMetadata.NewFailure(fmt.Sprintf("FindByIdentity returned %d", result))
}

// publish delegates to publisher's publish method, stores publish result as annotation, and returns status.
func (a *annotator) publish(data []byte) *status.Contract {
	var publishResult published.Contract

	id := a.identityProvider.Derive(data)
	annotations, result := a.store.FindByIdentity(id)
	switch result {
	case status.Success:
		publishResult = a.publisher.Publish(a.filter.Do(annotations))
	default:
		publishResult = a.failureFindByIdentity(result)
	}

	m := annotation.New(
		a.uniqueProvider.Get(),
		id,
		nil,
		publishMetadata.New(a.provenance, a.publisher.Kind(), publishResult),
	)
	result = a.store.Append(id, m)
	if result == status.NotFound {
		result = a.store.Create(id, m)
	}
	return status.New(a.provenance, result)
}

// Create evaluates newly-created data.
func (a *annotator) Create(data []byte) *status.Contract {
	return a.publish(data)
}

// Mutate evaluates mutated data.
func (a *annotator) Mutate(_, newData []byte) *status.Contract {
	return a.publish(newData)
}
