package iota

import (
	"testing"

	testInternal "github.com/project-alvarium/go-sdk/internal/pkg/test"
	envelope "github.com/project-alvarium/go-sdk/pkg/annotation/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/published"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/iota/metadata"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/iota/sdk"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/iota/sdk/iota"
	clientStub "github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/iota/sdk/iota/client/stub"
	"github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/iota/sdk/stub"
	"github.com/project-alvarium/go-sdk/pkg/hashprovider/sha256"
	identityProvider "github.com/project-alvarium/go-sdk/pkg/identityprovider/hash"
	"github.com/project-alvarium/go-sdk/pkg/test"

	"github.com/stretchr/testify/assert"
)

// newSUT returns a system under test.
func newSUT(sdk sdk.Contract) *publisher {
	return newWithIOTA(
		testInternal.FactoryRandomSeedString(),
		test.FactoryRandomUint64(),
		test.FactoryRandomUint64(),
		sdk,
	)
}

// TestNew tests publisher.New.
func TestNew(t *testing.T) {
	sut := New(
		testInternal.FactoryRandomSeedString(),
		test.FactoryRandomUint64(),
		test.FactoryRandomUint64(),
		clientStub.New(nil, nil, nil, nil),
	)

	assert.NotNil(t, sut)
	assert.IsType(t, &publisher{}, sut)
}

// TestPublisher_SetUp tests publisher.SetUp.
func TestPublisher_SetUp(t *testing.T) {
	sut := New(
		testInternal.FactoryRandomSeedString(),
		test.FactoryRandomUint64(),
		test.FactoryRandomUint64(),
		clientStub.New(nil, nil, nil, nil),
	)

	sut.SetUp()

	// nothing to assert; for coverage only
	sut.TearDown()
}

// TestPublisher_TearDown tests publisher.TearDown.
func TestPublisher_TearDown(t *testing.T) {
	sut := New(
		testInternal.FactoryRandomSeedString(),
		test.FactoryRandomUint64(),
		test.FactoryRandomUint64(),
		clientStub.New(nil, nil, nil, nil),
	)
	sut.SetUp()

	sut.TearDown()

	// nothing to assert; for coverage only
}

// TestPublisher_Publish tests publisher.Publish.
func TestPublisher_Publish(t *testing.T) {
	type testCase struct {
		name           string
		annotations    []*envelope.Annotations
		sdk            sdk.Contract
		expectedResult func(sut *publisher) published.Contract
	}

	cases := []testCase{
		func() testCase {
			return testCase{
				name:        "no annotations",
				annotations: []*envelope.Annotations{},
				sdk:         iota.New(clientStub.New(nil, nil, nil, nil)),
				expectedResult: func(sut *publisher) published.Contract {
					return sut.failureNoAnnotations()
				},
			}
		}(),
		func() testCase {
			resultTx := testInternal.FactoryRandomFixedSizeBundle(1)[0]
			result := metadata.New(resultTx.Address, resultTx.Hash, resultTx.Tag)
			return testCase{
				name: "single annotation",
				annotations: []*envelope.Annotations{
					envelope.New(
						test.FactoryRandomString(),
						identityProvider.New(sha256.New()).Derive(test.FactoryRandomByteSlice()),
						nil,
						test.FactoryRandomString(),
						test.FactoryRandomString(),
					),
				},
				sdk: stub.New(result),
				expectedResult: func(sut *publisher) published.Contract {
					return result
				},
			}
		}(),
		func() testCase {
			resultTx := testInternal.FactoryRandomFixedSizeBundle(1)[0]
			result := metadata.New(resultTx.Address, resultTx.Hash, resultTx.Tag)
			return testCase{
				name: "single annotation as structure",
				annotations: []*envelope.Annotations{
					envelope.New(
						test.FactoryRandomString(),
						identityProvider.New(sha256.New()).Derive(test.FactoryRandomByteSlice()),
						nil,
						test.FactoryRandomString(),
						struct {
							Name  string
							Value int
						}{
							Name:  test.FactoryRandomString(),
							Value: test.FactoryRandomInt(),
						},
					),
				},
				sdk: stub.New(result),
				expectedResult: func(sut *publisher) published.Contract {
					return result
				},
			}
		}(),
		func() testCase {
			resultTx := testInternal.FactoryRandomFixedSizeBundle(1)[0]
			result := metadata.New(resultTx.Address, resultTx.Hash, resultTx.Tag)
			idProvider := identityProvider.New(sha256.New())
			data := test.FactoryRandomByteSlice()
			id1 := idProvider.Derive(test.FactoryRandomByteSlice())
			id2 := idProvider.Derive(data)
			kind := test.FactoryRandomString()
			return testCase{
				name: "two annotations",
				annotations: []*envelope.Annotations{
					envelope.New(test.FactoryRandomString(), id2, id1, kind, test.FactoryRandomString()),
					envelope.New(test.FactoryRandomString(), id1, nil, kind, test.FactoryRandomString()),
				},
				sdk: stub.New(result),
				expectedResult: func(sut *publisher) published.Contract {
					return result
				},
			}
		}(),
	}

	for i := range cases {
		t.Run(cases[i].name, func(t *testing.T) {
			sut := newSUT(cases[i].sdk)

			result := sut.Publish(cases[i].annotations)

			assert.Equal(t, cases[i].expectedResult(sut), result)
		})
	}

}

// TestPublisher_Kind tests publisher.Kind.
func TestPublisher_Kind(t *testing.T) {
	sut := newSUT(stub.New(nil))

	result := sut.Kind()

	assert.Equal(t, name, result)
}

// TestKind tests Kind.
func TestKind(t *testing.T) {
	assert.Equal(t, name, Kind())
}