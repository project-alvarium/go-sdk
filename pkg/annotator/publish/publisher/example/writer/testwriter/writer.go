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

package testwriter

// instance is a receiver that encapsulates required dependencies.
type instance struct {
	value []byte
}

// New is a factory function that returns an initialized instance.
func New() *instance {
	return &instance{}
}

// Write fulfills the io.Writer contract.
func (i *instance) Write(p []byte) (n int, err error) {
	i.value = p
	return len(p), nil
}

// Get returns the value written via Write().
func (i *instance) Get() []byte {
	return i.value
}
