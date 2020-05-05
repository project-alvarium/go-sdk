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

package main

import (
	"crypto"
	"encoding/json"
	"fmt"

	testInternal "github.com/project-alvarium/go-sdk/internal/pkg/test"
	"github.com/project-alvarium/go-sdk/pkg/annotation/store/memory"
	"github.com/project-alvarium/go-sdk/pkg/annotation/uniqueprovider/ulid"
	"github.com/project-alvarium/go-sdk/pkg/annotator"
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess"
	pkiAssessor "github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/pki"
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/pki/factory/verifier"
	"github.com/project-alvarium/go-sdk/pkg/annotator/filter/passthrough"
	pkiAnnotator "github.com/project-alvarium/go-sdk/pkg/annotator/pki"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signpkcs1v15"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/example"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/example/writer/testwriter"
	"github.com/project-alvarium/go-sdk/pkg/hashprovider/sha256"
	identityProvider "github.com/project-alvarium/go-sdk/pkg/identityprovider/hash"
	"github.com/project-alvarium/go-sdk/pkg/sdk"
	"github.com/project-alvarium/go-sdk/pkg/test"
)

// main is the example entry point.
func main() {
	hashProvider := sha256.New()
	uniqueProvider := ulid.New()
	idProvider := identityProvider.New(hashProvider)
	persistence := memory.New()
	filter := passthrough.New()

	// create SDK instance for annotation.
	sdkInstance := sdk.New(
		[]annotator.Contract{
			pkiAnnotator.New(
				struct{ Node string }{Node: "origin"},
				uniqueProvider,
				idProvider,
				persistence,
				signpkcs1v15.New(
					crypto.SHA256,
					testInternal.ValidPrivateKey,
					testInternal.ValidPublicKey,
					hashProvider,
				),
			),
		},
	)

	// register data creation.
	data, _ := json.Marshal(test.FactoryRandomString())
	_ = sdkInstance.Create(data)

	// close SDK instance.
	sdkInstance.Close()

	// create SDK instance for assessment.
	p := struct{ Node string }{Node: "evaluation"}
	w := testwriter.New()
	sdkInstance = sdk.New(
		[]annotator.Contract{
			assess.New(
				p,
				uniqueProvider,
				idProvider,
				persistence,
				pkiAssessor.New(verifier.New()),
				filter,
			),
			publish.New(p, uniqueProvider, idProvider, persistence, example.New(w), filter),
		},
	)

	// assess and publish result.
	_ = sdkInstance.Create(data)

	// display it.
	fmt.Printf("%s\n", w.Get())

	// close SDK instance
	sdkInstance.Close()
}
