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

package reducer

import (
	"github.com/project-alvarium/go-sdk/pkg/hashprovider"
	"github.com/project-alvarium/go-sdk/pkg/hashprovider/md5"
	"github.com/project-alvarium/go-sdk/pkg/hashprovider/sha256"
)

// To converts a string representation to a hashprovider.Contract (or nil if unable to do so).
func To(name string) (hash hashprovider.Contract) {
	switch name {
	case md5.Name():
		hash = md5.New()
	case sha256.Name():
		hash = sha256.New()
	default:
		hash = nil
	}
	return
}

// Supported returns a slice of supported hashprovider.Contract.
func Supported() []hashprovider.Contract {
	return []hashprovider.Contract{
		To(md5.Name()),
		To(sha256.Name()),
	}
}
