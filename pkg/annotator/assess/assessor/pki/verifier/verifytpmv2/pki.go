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

package verifytpmv2

import (
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/pki/verifier"
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/pki/verifier/verifypkcs1v15"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signtpmv2/provisioner"
	"github.com/project-alvarium/go-sdk/pkg/hashprovider"
)

// New is a factory function that returns verifier.
func New(reducerHash hashprovider.Contract) verifier.Contract {
	return verifypkcs1v15.New(provisioner.CryptoHash, reducerHash)
}
