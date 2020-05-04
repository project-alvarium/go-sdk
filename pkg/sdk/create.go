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

import "github.com/project-alvarium/go-sdk/pkg/status"

// Create calls the Create method on each registered annotator and returns a set of status results.
func (sdk *instance) Create(data []byte) []*status.Contract {
	if sdk.closed {
		return nil
	}

	result := make([]*status.Contract, 0)
	for i := range sdk.annotators {
		result = append(result, sdk.annotators[i].Create(data))
	}
	return result
}
