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

package iota

import (
	"errors"
	"testing"

	testInternal "github.com/project-alvarium/go-sdk/internal/pkg/test"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/published"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/iota/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/iota/sdk/iota/client"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/iota/sdk/iota/client/stub"
	"github.com/project-alvarium/go-sdk/pkg/test"

	"github.com/iotaledger/iota.go/bundle"
	"github.com/iotaledger/iota.go/consts"
	"github.com/iotaledger/iota.go/trinary"
	"github.com/stretchr/testify/assert"
)

// newSUT returns a new system under test.
func newSUT(client client.Contract) *instance {
	return New(client)
}

// TestInstance_Send tests instance.Send.
func TestInstance_Send(t *testing.T) {
	newClient := func(
		addressResultError error,
		sendResultValue bundle.Bundle,
		sendResultError error) *stub.Instance {

		return stub.New(
			[]trinary.Hash{testInternal.FactoryRandomAddressTrytesString()},
			addressResultError,
			sendResultValue,
			sendResultError,
		)
	}

	type testCase struct {
		name           string
		client         client.Contract
		annotations    []byte
		expectedResult func(sut *instance) published.Contract
	}

	cases := []testCase{
		func() testCase {
			res := testInternal.FactoryRandomFixedSizeBundle(1)
			return testCase{
				name:        "send transfer",
				client:      newClient(nil, res, nil),
				annotations: test.FactoryRandomByteSlice(),
				expectedResult: func(_ *instance) published.Contract {
					return metadata.New(res[0].Address, res[0].Hash, res[0].Tag)
				},
			}
		}(),
		func() testCase {
			errMessage := test.FactoryRandomString()
			return testCase{
				name:        "get new address error",
				client:      newClient(errors.New(errMessage), testInternal.FactoryRandomFixedSizeBundle(1), nil),
				annotations: test.FactoryRandomByteSlice(),
				expectedResult: func(sut *instance) published.Contract {
					return sut.getNewAddressError(errMessage)
				},
			}
		}(),
		func() testCase {
			return testCase{
				name:        "convert message to trytes error",
				client:      newClient(nil, testInternal.FactoryRandomFixedSizeBundle(1), nil),
				annotations: []byte("\u2000"),
				expectedResult: func(sut *instance) published.Contract {
					return sut.convertToTrytesError(consts.ErrInvalidASCIIInput.Error())
				},
			}
		}(),
		func() testCase {
			errMessage := test.FactoryRandomString()
			return testCase{
				name:        "send transfer error",
				client:      newClient(nil, testInternal.FactoryRandomFixedSizeBundle(1), errors.New(errMessage)),
				annotations: test.FactoryRandomByteSlice(),
				expectedResult: func(sut *instance) published.Contract {
					return sut.sendTransferError(errMessage)
				},
			}
		}(),
		func() testCase {
			return testCase{
				name:        "send transfer success no results",
				client:      newClient(nil, bundle.Bundle{}, nil),
				annotations: test.FactoryRandomByteSlice(),
				expectedResult: func(sut *instance) published.Contract {
					return sut.invalidResultSetError()
				},
			}
		}(),
		func() testCase {
			return testCase{
				name:        "send transfer success multiple results",
				client:      newClient(nil, testInternal.FactoryRandomFixedSizeBundle(test.FactoryRandomInt()), nil),
				annotations: test.FactoryRandomByteSlice(),
				expectedResult: func(sut *instance) published.Contract {
					return sut.invalidResultSetError()
				},
			}
		}(),
	}

	for i := range cases {
		t.Run(cases[i].name, func(t *testing.T) {
			sut := newSUT(cases[i].client)

			result := sut.Send(
				testInternal.FactoryRandomSeedString(),
				test.FactoryRandomUint64(),
				test.FactoryRandomUint64(),
				cases[i].annotations,
			)

			assert.Equal(t, testInternal.Marshal(t, cases[i].expectedResult(sut)), testInternal.Marshal(t, result))
		})
	}
}
