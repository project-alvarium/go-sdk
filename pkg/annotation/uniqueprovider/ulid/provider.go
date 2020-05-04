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

package ulid

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

// provider is a receiver that encapsulates required dependencies.
type provider struct {
	e *ulid.MonotonicEntropy
}

// New is a factory function that returns an initialized provider.
func New() *provider {
	return &provider{
		e: ulid.Monotonic(rand.New(rand.NewSource(time.Unix(1000000, 0).UnixNano())), 0),
	}
}

// Get returns a globally unique identifier.
func (p *provider) Get() string {
	return ulid.MustNew(ulid.Now(), p.e).String()
}
