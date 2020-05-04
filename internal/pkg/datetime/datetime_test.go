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

package datetime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestCreated tests Created.
func TestCreated(t *testing.T) {
	result := Created()

	assert.NotNil(t, TimeFromCreated(result))
}

// TestTimeFromCreated tests TimeFromCreated.
func TestTimeFromCreated(t *testing.T) {
	type testCase struct {
		name string
		test func(t *testing.T)
	}

	cases := []testCase{
		{
			name: "Success",
			test: func(t *testing.T) {
				result := TimeFromCreated(Created())

				assert.NotNil(t, result)
				assert.True(t, time.Now().UTC().Sub(*result).Nanoseconds() <= time.Second.Nanoseconds())
			},
		},
		{
			name: "Failure",
			test: func(t *testing.T) {
				result := TimeFromCreated("badDateString")

				assert.Nil(t, result)
			},
		},
	}

	for i := range cases {
		t.Run(cases[i].name, cases[i].test)
	}
}
