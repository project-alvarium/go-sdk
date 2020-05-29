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

package matching

import "github.com/project-alvarium/go-sdk/pkg/annotation"

// Compare is a function which facilitates filtering determining if an annotation should be included in the result.
type Compare func(annotation *annotation.Instance) bool

// filter is receiver that encapsulates required dependencies.
type filter struct {
	compare Compare
}

// New is a factory function that returns an initialized iota filter instance.
func New(compare Compare) *filter {
	return &filter{
		compare: compare,
	}
}

// Do implements an annotation filter.
func (f *filter) Do(annotations []*annotation.Instance) []*annotation.Instance {
	filteredAnnotations := make([]*annotation.Instance, 0)
	for i := range annotations {
		if f.compare(annotations[i]) {
			filteredAnnotations = append(filteredAnnotations, annotations[i])
		}
	}
	return filteredAnnotations
}
