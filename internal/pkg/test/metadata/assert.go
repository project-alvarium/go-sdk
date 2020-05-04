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

package metadata

import (
	"testing"
	"time"

	"github.com/project-alvarium/go-sdk/internal/pkg/datetime"
	"github.com/project-alvarium/go-sdk/pkg/annotation"
	"github.com/project-alvarium/go-sdk/pkg/annotation/store"
	"github.com/project-alvarium/go-sdk/pkg/identity"
	"github.com/project-alvarium/go-sdk/pkg/status"

	"github.com/stretchr/testify/assert"
)

// Assert compares metadata actual against expected and allows for a +/- one second variance in the Created value.
func Assert(t *testing.T, expected []*annotation.Instance, id identity.Contract, store store.Contract) {
	actual, result := store.FindByIdentity(id)
	if !assert.Equal(t, status.Success, result) {
		assert.FailNow(t, "FindByIdentity faulted")
	}

	if !assert.Equal(t, len(expected), len(actual)) {
		assert.FailNow(t, "expected != actual")
	}

	for i := range expected {
		assert.Equal(t, expected[i].CurrentIdentity, actual[i].CurrentIdentity)
		assert.Equal(t, expected[i].PreviousIdentity, actual[i].PreviousIdentity)

		expectedCreated := datetime.TimeFromCreated(expected[i].Created)
		actualCreated := datetime.TimeFromCreated(actual[i].Created)
		if expectedCreated != nil && actualCreated != nil {
			diff := expectedCreated.Sub(*actualCreated).Nanoseconds()
			if diff < 0 {
				diff = -diff
			}
			assert.True(t, diff <= time.Second.Nanoseconds())
		} else {
			assert.FailNow(t, "unable to parse created.")
		}

		assert.Equal(t, expected[i].MetadataKind, actual[i].MetadataKind)
		assert.Equal(t, expected[i].Metadata, actual[i].Metadata)
	}
}
