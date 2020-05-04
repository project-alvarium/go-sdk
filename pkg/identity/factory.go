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

package identity

import identityHash "github.com/project-alvarium/go-sdk/pkg/identity/hash"

// Factory returns a contract implementation based on the provided metadata.
func Factory(kind string, implementation interface{}) Contract {
	switch kind {
	case identityHash.Name:
		if m, ok := implementation.(*identityHash.Identity); ok {
			return m
		}
	}
	return nil
}
