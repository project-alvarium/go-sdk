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

package sdk

import (
	"testing"

	metadataStub "github.com/project-alvarium/go-sdk/pkg/annotation/metadata/stub"
	"github.com/project-alvarium/go-sdk/pkg/annotation/store"
	"github.com/project-alvarium/go-sdk/pkg/annotation/uniqueprovider/ulid"
	"github.com/project-alvarium/go-sdk/pkg/annotator"
	"github.com/project-alvarium/go-sdk/pkg/annotator/filter/passthrough"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/fail"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish"
	publishStub "github.com/project-alvarium/go-sdk/pkg/annotator/publish/stub"
	annotatorStub "github.com/project-alvarium/go-sdk/pkg/annotator/stub"
	"github.com/project-alvarium/go-sdk/pkg/hashprovider/sha256"
	identityProvider "github.com/project-alvarium/go-sdk/pkg/identityprovider/hash"
	"github.com/project-alvarium/go-sdk/pkg/store/memory"
	"github.com/project-alvarium/go-sdk/pkg/test"

	"github.com/stretchr/testify/assert"
)

// newSUT returns a new system under test.
func newSUT(annotators []annotator.Contract) *instance {
	return New(annotators)
}

// TestInstance_SetUpCalled tests New.
func TestInstance_SetUpCalled(t *testing.T) {
	type testCase struct {
		name string
		test func(t *testing.T)
	}

	cases := []testCase{
		{
			name: "SetUp called (stub annotator)",
			test: func(t *testing.T) {
				a := annotatorStub.New()
				sut := newSUT([]annotator.Contract{a})

				assert.Equal(t, 1, a.SetUpCalled)
				sut.Close()
			},
		},
		{
			name: "SetUp called (annotator - pki signer)",
			test: func(t *testing.T) {
				s := fail.New()
				a := pki.New(
					test.FactoryRandomString(),
					ulid.New(),
					identityProvider.New(sha256.New()),
					store.New(memory.New()),
					s,
				)
				sut := newSUT([]annotator.Contract{a})

				assert.True(t, s.SetUpCalled)
				sut.Close()
			},
		},
		{
			name: "SetUp called (annotator - stub publisher)",
			test: func(t *testing.T) {
				s := store.New(memory.New())
				p := publishStub.New(test.FactoryRandomString(), metadataStub.NewNullObject())
				a := publish.New(
					test.FactoryRandomString(),
					ulid.New(),
					identityProvider.New(sha256.New()),
					s,
					p,
					passthrough.New(),
				)
				sut := newSUT([]annotator.Contract{a})

				assert.True(t, p.SetUpCalled)
				sut.Close()
			},
		},
	}

	for i := range cases {
		t.Run(cases[i].name, cases[i].test)
	}
}
