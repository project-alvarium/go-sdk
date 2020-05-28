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
	"flag"
	"fmt"
	"io"
	"os"

	testInternal "github.com/project-alvarium/go-sdk/internal/pkg/test"
	"github.com/project-alvarium/go-sdk/pkg/annotation"
	metadataFactory "github.com/project-alvarium/go-sdk/pkg/annotation/metadata/factory"
	"github.com/project-alvarium/go-sdk/pkg/annotation/store/memory"
	"github.com/project-alvarium/go-sdk/pkg/annotation/uniqueprovider/ulid"
	"github.com/project-alvarium/go-sdk/pkg/annotator"
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess"
	iotaAssessor "github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/iota"
	pkiAssessor "github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/pki"
	"github.com/project-alvarium/go-sdk/pkg/annotator/assess/assessor/pki/factory/verifier"
	filterFactory "github.com/project-alvarium/go-sdk/pkg/annotator/filter/matching"
	"github.com/project-alvarium/go-sdk/pkg/annotator/filter/passthrough"
	pkiAnnotator "github.com/project-alvarium/go-sdk/pkg/annotator/pki"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signpkcs1v15"
	tpmSigner "github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signtpmv2"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signtpmv2/factory"
	"github.com/project-alvarium/go-sdk/pkg/annotator/pki/signer/signtpmv2/provisioner"
	"github.com/project-alvarium/go-sdk/pkg/annotator/provenance"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish"
	publishMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/publish/metadata"
	publisherMetadataFactory "github.com/project-alvarium/go-sdk/pkg/annotator/publish/metadata/factory"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/example"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/example/writer/testwriter"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/iota"
	iotaPublisherMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/iota/metadata"
	iotaMetadataFactory "github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/iota/metadata/factory"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/ipfs"
	ipfsPublisherMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/ipfs/metadata"
	"github.com/project-alvarium/go-sdk/pkg/hashprovider/sha256"
	identityProvider "github.com/project-alvarium/go-sdk/pkg/identityprovider/hash"
	"github.com/project-alvarium/go-sdk/pkg/sdk"
	"github.com/project-alvarium/go-sdk/pkg/test"

	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpmutil"
	"github.com/iotaledger/iota.go/api"
)

// exampleData defines the structure of the example data.
type exampleData struct {
	Name  string
	Value int
}

// tpmSetUp creates a temporary key and certificate to use for the example.
func tpmSetUp(path string) (io.ReadWriteCloser, tpmutil.Handle, string, []byte, func()) {
	rwc, err := factory.TPM(path)
	if err != nil {
		fmt.Println("Unable to factory TPM instance")
		os.Exit(1)
	}

	handle, publicKey, err := provisioner.GenerateNewKeyPair(rwc)
	if err != nil {
		fmt.Println("Unable to generate new key pair")
		os.Exit(1)
	}

	return rwc,
		handle,
		path,
		provisioner.MarshalPublicKey(publicKey),
		func() {
			provisioner.Flush(rwc, handle)
		}
}

// newExampleData is a factory function that returns an initialized exampleData.
func newExampleData() *exampleData {
	return &exampleData{
		Name:  test.FactoryRandomString(),
		Value: test.FactoryRandomInt(),
	}
}

// newProvenance is a factory function that returns a provenance.Contract.
func newProvenance(node string) provenance.Contract {
	return &struct {
		Node string `json:"node"`
	}{
		Node: node,
	}
}

// newClient is a factory function that returns an api.API reference.
func newClient(url string) *api.API {
	client, err := api.ComposeAPI(api.HTTPClientSettings{URI: url})
	if err != nil {
		fmt.Println("Unable to factory IOTA API instance")
		os.Exit(1)
	}
	return client
}

// main is the example entry point.
func main() {
	var iotaURL string
	var iotaDepth uint64
	var iotaMWM uint64
	var ipfsURL string
	flag.StringVar(&iotaURL, "iotaURL", "http://localhost:14265", "IOTA Tangle URL (http://localhost:14265)")
	flag.Uint64Var(&iotaDepth, "iotaDepth", 3, "IOTA Tangle Depth (3)")
	flag.Uint64Var(&iotaMWM, "iotaMWM", 9, "IOTA Tangle Minimum Weight Magnitude (9)")
	flag.StringVar(&ipfsURL, "ipfsURL", "localhost:5001", "IPFS URL (localhost:5001)")
	flag.Parse()

	hashProvider := sha256.New()
	uniqueProvider := ulid.New()
	idProvider := identityProvider.New(hashProvider)
	persistence := memory.New()
	passthroughFilter := passthrough.New()

	// create new TPM keys
	rwc, tpmHandle, tpmPath, publicKey, cleanUp := tpmSetUp(provisioner.Path)

	// create SDK instance for annotation and assessment.
	p := newProvenance("origin")
	sdkInstance := sdk.New(
		[]annotator.Contract{
			pkiAnnotator.New(
				p,
				uniqueProvider,
				idProvider,
				persistence,
				tpmSigner.NewWithRWC(
					hashProvider,
					publicKey,
					tpmHandle,
					tpmPath,
					tpmSigner.RequestedCapabilityProperties{
						"Version":      tpm2.FamilyIndicator,
						"Manufacturer": tpm2.Manufacturer,
					},
					rwc),
			),
			assess.New(
				p,
				uniqueProvider,
				idProvider,
				persistence,
				pkiAssessor.New(verifier.New()),
				passthroughFilter,
			),
		},
	)

	// register data creation.
	data := newExampleData()
	dataAsBytes, _ := json.Marshal(data)
	_ = sdkInstance.Create(dataAsBytes)

	// cleanup TPM resources and close SDK instance.
	cleanUp()
	sdkInstance.Close()

	// create SDK instance for annotation and assessment.
	p = newProvenance("transit-1")
	sdkInstance = sdk.New(
		[]annotator.Contract{
			pkiAnnotator.New(
				p,
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
			assess.New(
				p,
				uniqueProvider,
				idProvider,
				persistence,
				pkiAssessor.New(verifier.New()),
				passthroughFilter,
			),
		},
	)

	// modify data; register data mutation.
	data.Value += 1
	newDataAsBytes, _ := json.Marshal(data)
	_ = sdkInstance.Mutate(dataAsBytes, newDataAsBytes)

	// close SDK instance.
	sdkInstance.Close()

	// create SDK instance for annotation and assessment.
	p = newProvenance("transit-2")
	sdkInstance = sdk.New(
		[]annotator.Contract{
			pkiAnnotator.New(
				p,
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
			assess.New(
				p,
				uniqueProvider,
				idProvider,
				persistence,
				pkiAssessor.New(verifier.New()),
				passthroughFilter,
			),
		},
	)

	// even though no mutation occurred, register data mutation to capture transit event.
	_ = sdkInstance.Mutate(newDataAsBytes, newDataAsBytes)

	// close SDK instance.
	sdkInstance.Close()

	// create SDK instance for publishing.
	p = newProvenance("publisher")
	w := testwriter.New()
	sdkInstance = sdk.New(
		[]annotator.Contract{
			publish.New(p, uniqueProvider, idProvider, persistence, ipfs.New(ipfsURL), passthroughFilter),
			publish.New(
				p,
				uniqueProvider,
				idProvider,
				persistence,
				iota.New(testInternal.FactoryRandomSeedString(), iotaDepth, iotaMWM, newClient(iotaURL)),
				filterFactory.New(
					func(annotation *annotation.Instance) bool {
						t, ok := annotation.Metadata.(*publishMetadata.Instance)
						return ok && t.PublisherKind == ipfsPublisherMetadata.Kind
					},
				),
			),
			assess.New(
				p,
				uniqueProvider,
				idProvider,
				persistence,
				iotaAssessor.New(
					newClient(iotaURL),
					publisherMetadataFactory.New([]metadataFactory.Contract{iotaMetadataFactory.New()}),
				),
				filterFactory.New(
					func(annotation *annotation.Instance) bool {
						t, ok := annotation.Metadata.(*publishMetadata.Instance)
						return ok && t.PublisherKind == iotaPublisherMetadata.Kind
					},
				),
			),
			publish.New(p, uniqueProvider, idProvider, persistence, example.New(w), passthroughFilter),
		},
	)

	// publish result.
	_ = sdkInstance.Create(newDataAsBytes)

	// display it.
	fmt.Printf("%s\n", w.Get())

	// close SDK instance
	sdkInstance.Close()
}
