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

package assess

import (
	"testing"

	testMetadata "github.com/project-alvarium/go-sdk/internal/pkg/test/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotation"
	metadataStub "github.com/project-alvarium/go-sdk/pkg/annotation/metadata/stub"
	"github.com/project-alvarium/go-sdk/pkg/annotation/store"
	"github.com/project-alvarium/go-sdk/pkg/annotation/store/memory"
	"github.com/project-alvarium/go-sdk/pkg/annotation/uniqueprovider/ulid"
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor"
	assessorStub "github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/stub"
	assessMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/assess/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotator/filter/passthrough"
	"github.com/project-alvarium/go-sdk/pkg/annotator/provenance"
	"github.com/project-alvarium/go-sdk/pkg/hashprovider/sha256"
	"github.com/project-alvarium/go-sdk/pkg/identityprovider"
	identityProvider "github.com/project-alvarium/go-sdk/pkg/identityprovider/hash"
	"github.com/project-alvarium/go-sdk/pkg/status"
	"github.com/project-alvarium/go-sdk/pkg/test"

	"github.com/stretchr/testify/assert"
)

// newSUT returns a new system under test.
func newSUT(
	provenance provenance.Contract,
	identityProvider identityprovider.Contract,
	store store.Contract,
	assessor assessor.Contract) *annotator {

	return New(provenance, ulid.New(), identityProvider, store, assessor, passthrough.New())
}

// TestAnnotator_SetUp tests annotator.SetUp.
func TestAnnotator_SetUp(t *testing.T) {
	a := assessorStub.New(test.FactoryRandomString(), metadataStub.NewNullObject())
	sut := newSUT(test.FactoryRandomString(), identityProvider.New(sha256.New()), memory.New(), a)

	sut.SetUp()

	assert.True(t, a.SetUpCalled)
}

// TestAnnotator_TearDown tests annotator.TearDown.
func TestAnnotator_TearDown(t *testing.T) {
	a := assessorStub.New(test.FactoryRandomString(), metadataStub.NewNullObject())
	sut := newSUT(test.FactoryRandomString(), identityProvider.New(sha256.New()), memory.New(), a)

	sut.TearDown()

	assert.True(t, a.TearDownCalled)
}

// TestAnnotator_Create tests annotator.Create.
func TestAnnotator_Create(t *testing.T) {
	type testCase struct {
		name string
		test func(t *testing.T)
	}

	cases := []testCase{
		{
			name: "no annotations in storage",
			test: func(t *testing.T) {
				prov := test.FactoryRandomString()
				idProvider := identityProvider.New(sha256.New())
				persistence := memory.New()
				data := test.FactoryRandomByteSlice()
				id := idProvider.Derive(data)
				sut := newSUT(
					prov,
					idProvider,
					persistence,
					assessorStub.New(test.FactoryRandomString(), metadataStub.NewNullObject()),
				)

				result := sut.Create(data)

				assert.Equal(t, status.New(prov, status.Success), result)
				testMetadata.Assert(
					t,
					[]*annotation.Instance{
						annotation.New(
							test.FactoryRandomString(),
							id,
							nil,
							assessMetadata.New(prov, sut.assessor.Failure(sut.failureFindByIdentity(status.NotFound))),
						),
					},
					id,
					persistence,
				)
			},
		},
		{
			name: "annotation in storage",
			test: func(t *testing.T) {
				prov := test.FactoryRandomString()
				idProvider := identityProvider.New(sha256.New())
				persistence := memory.New()
				kind := test.FactoryRandomString()
				data := test.FactoryRandomByteSlice()
				id := idProvider.Derive(data)
				m := metadataStub.New(kind, test.FactoryRandomString())
				a := annotation.New(test.FactoryRandomString(), id, nil, metadataStub.New(kind, m))
				assert.Equal(t, status.Success, persistence.Create(id, a))
				sut := newSUT(prov, idProvider, persistence, assessorStub.New(test.FactoryRandomString(), m))

				result := sut.Create(data)

				assert.Equal(t, status.New(prov, status.Success), result)
				testMetadata.Assert(
					t,
					[]*annotation.Instance{
						a,
						annotation.New(test.FactoryRandomString(), id, nil, assessMetadata.New(prov, m)),
					},
					id,
					persistence,
				)
			},
		},
	}

	for i := range cases {
		t.Run(cases[i].name, cases[i].test)
	}
}

// TestAnnotator_Mutate tests annotator.Mutate.
func TestAnnotator_Mutate(t *testing.T) {
	type testCase struct {
		name string
		test func(t *testing.T)
	}

	cases := []testCase{
		{
			name: "no annotations in storage",
			test: func(t *testing.T) {
				prov := test.FactoryRandomString()
				idProvider := identityProvider.New(sha256.New())
				persistence := memory.New()
				data := test.FactoryRandomByteSlice()
				id := idProvider.Derive(data)
				sut := newSUT(
					prov,
					idProvider,
					persistence,
					assessorStub.New(test.FactoryRandomString(), metadataStub.NewNullObject()),
				)

				result := sut.Mutate(data, data)

				assert.Equal(t, status.New(prov, status.Success), result)
				testMetadata.Assert(
					t,
					[]*annotation.Instance{
						annotation.New(
							test.FactoryRandomString(),
							id,
							nil,
							assessMetadata.New(prov, sut.assessor.Failure(sut.failureFindByIdentity(status.NotFound))),
						),
					},
					id,
					persistence,
				)
			},
		},
		{
			name: "annotation in storage",
			test: func(t *testing.T) {
				prov := test.FactoryRandomString()
				idProvider := identityProvider.New(sha256.New())
				persistence := memory.New()
				kind := test.FactoryRandomString()
				data := test.FactoryRandomByteSlice()
				id := idProvider.Derive(data)
				m := metadataStub.New(kind, test.FactoryRandomString())
				a := annotation.New(test.FactoryRandomString(), id, nil, m)
				assert.Equal(t, status.Success, persistence.Create(id, a))
				sut := newSUT(prov, idProvider, persistence, assessorStub.New(kind, m))

				result := sut.Mutate(data, data)

				assert.Equal(t, status.New(prov, status.Success), result)
				testMetadata.Assert(
					t,
					[]*annotation.Instance{
						a,
						annotation.New(test.FactoryRandomString(), id, nil, assessMetadata.New(prov, m)),
					},
					id,
					persistence,
				)
			},
		},
	}

	for i := range cases {
		t.Run(cases[i].name, cases[i].test)
	}
}
