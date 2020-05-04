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

package stub

// Instance is a receiver that encapsulates required dependencies.
type Instance struct {
	CapturedURL         string
	CapturedAnnotations []byte
	resultValue         string
	resultError         error
}

// New is a factory function that returns an initialized Instance.
func New(resultValue string, resultError error) *Instance {
	return &Instance{
		resultValue: resultValue,
		resultError: resultError,
	}
}

// Add is called to add annotations to the IPFS instance that resides at url.
func (i *Instance) Add(url string, annotations []byte) (string, error) {
	i.CapturedURL = url
	i.CapturedAnnotations = annotations
	return i.resultValue, i.resultError
}
