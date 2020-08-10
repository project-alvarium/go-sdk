![Project Alvarium](README.assets/ProjectAlvarium.png)

# SDK (Golang)

This repository contains a software development kit (SDK) to track data provenance across mutations, derive related annotations, and provide a generic architecture to assess those annotations to reach one or more conclusions.



## Table of Contents

- [Installation](#installation)
- [Getting Started](#getting-started)
  - [Ubiquitous Language](#ubiquitous-language)
  - [Data Provenance](#data-provenance)
  - [Data Identity](#data-identity)
  - [Data Sovereignty](#data-sovereignty)
  - [Annotators and Annotations](#annotators-and-annotations)
- [Basic SDK Usage](#basic-sdk-usage)
- [Example Code](#example-code)
- [Project Structure](#project-structure)
- [Opportunities for Improvement](#opportunities-for-improvement)



## Installation

This repository contains a module written in Golang.  

It assumes a minimum Go version of 1.13.  

To add it to your project: `go get github.com/project-alvarium/go-sdk`



## Getting Started

### Ubiquitous Language

Defines the terms used to communicate SDK concepts (which are echoed in implementation).

- **Annotation**. Information (metadata) derived from (and related to) specific data. 
- **Annotator**. Evaluates data to create  annotations.
- **Assessment**.  Conclusion(s) derived from specific data and related annotations.
- **Assessor**.  Assesses annotations and data to create an assessment.
- **Data**.  Plural; two or more data points.
- **Data Lineage**.  Description of the lifecycle of data to the current time.  May Include details about the data's origin, who/what accessed or mutated the data, and where the data has transited.  This is a key attribute of data provenance (that is, why-provenance).
- **Data Point**.  Singular; a piece of information.
- **Data Provenance**.  A historical record of the data and its origins.
- **Identity**.  A specific data's unique identifier.  
- **Identity Provider**.  Derives identity from data.
- **Publisher**.  Publishes annotations to external systems or persistence.



### Data Provenance

The SDK tracks the provenance of any type of data (for example, JSON, XML, or binary content) as long as the data can be represented as a Golang `[]byte` and be reduced to an identity.  



### Data Identity

Identity can be derived algorithmically (for example, the result of a hash function applied to the data) or from introspection (for example, use of a known unique identifier within the data).  

The SDK defines an [identity abstraction](pkg/identity/contract.go) to encapsulate data identity.

The SDK defines an [identity provider abstraction](pkg/identityprovider/contract.go) to encapsulate the creation of identity from data.



### Data Sovereignty

The SDK tracks annotations associated with data via the data's identity.  

As implemented, the SDK does not modify or maintain a copy of the data it tracks.  There is nothing to prevent a custom annotator from keeping a copy of the data it tracks.  The annotators currently provided with the SDK do not do so.



### Annotators and Annotations

The SDK revolves around the concept of annotators and annotations.  Pass data to the SDK and it delegates to the list of configured annotators to create zero or more annotations.  

Related abstractions and SDK implementations are detailed in the [Annotator documentation](pkg/annotator/README.md).



## Basic SDK Usage

The SDK provides a minimal API -- New(), Create(), Mutate(), and Close().



### New()

```go
func New(annotators []annotator.Contract) *instance
```

Used to instantiate a new SDK instance with the specified list of annotators.

Takes a list of annotators and returns an SDK instance.



### Create()

```go
func (sdk *instance) Create(data []byte) []*status.Contract
```

Used to register creation of new data with the SDK.  Passes data through the SDK instance's list of annotators.

SDK instance method.  Takes the data to annotate and returns a status.  

Returns nil (and does not annotate) if `Close()` was previously called for the instance.



### Mutate()

```go
func (sdk *instance) Mutate(oldData, newData []byte) []*status.Contract
```

Used to register mutation of existing data with the SDK.  Passes data through the SDK instance's list of annotators.

SDK instance method.  Takes the data to annotate (and its prior version) and returns a status.  

Returns nil (and does not annotate) if `Close()` was previously called for the instance.



### Close()

```go
func (sdk *instance) Close()
```

SDK instance method.  Call to ensure proper shutdown of the SDK.



## Example Code

Simple example that uses the SDK's PKI annotator, PKI verifier, and example publisher.

```go
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
```

Executing this code results in JSON output resembling the following:

```json
[
  {
    "unique": "01E7JQF86JMQJHBF4QX1EFD6Y3",
    "created": "2020-05-05T15:32:34.6429938Z",
    "identityCurrentType": "hash",
    "identityCurrent": {
      "hash": "Vc1WRGhjPV5SoXMDfRBazQMrAcvUVzAJi0+9x+hnjbE="
    },
    "identityPreviousType": "hash",
    "identityPrevious": null,
    "metadataType": "pki",
    "metadata": {
      "provenance": {
        "Node": "origin"
      },
      "identitySignature": "tFohwFE4I1AXi4JLkrir5GkJzts+JJzicHgpL71AMuDvEWiG6edmsk2VE46lX9mnBfV1ja56+kKhHFISQFESgyvB9H/dwoUngMxrGKKwII6SoHqXpKjGrZ6qfV8R+sjgogDQ0IKcNbB6ouGaK1+0pMAsht1vorlYGWVC9qnWpzIEHAHjyC0T6VUzoiY4TcCS2SIiIdTGFdkgxd+gGyFAmJ2pd4kSUJrnRZbWr22tsVNMkgDC/Yqca4OLYDQv0ulNHBZK10M3QP9pK7yIutFyHjo0YTzoiNZGjDxKmurwaO10+lduo1qIiEqG2U4pTdUfL4d7V0h81OZqYus/azJpPQ==",
      "dataSignature": "30oRl1lXT2nASM9X9LYkKqQL0FWpjN8dm95H9czDya4OIsru7wwVb3PFiOgTBtAMT+V//XrXBSkr1Hhqes1K41Q03Qv3AJGYUBeImyOV60dGj0+DxKZEbNhvH1ElyzVhhDWNz2994WPbpQ5XbEzK76XGlCCBR/AIiRMKPjb1zr8DNdwlUyLz7yI5Pqt8eZxpBN4TbI/y1eCwRgRbMxIgr1bht0B38rw5mdhfy0I854HztbDmvxiQqahL4oooGWDHFrh7w0rtWHGBf3ZrtQPR16fxqu0QvRQ0N/tS7q+mqjGtkG3AlWzBJD32M11amGj9JXb/JP3vriY4t/rh1HnJ/Q==",
      "publicKey": "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUErQ2wrYVByWFlybUJMUGo3TWZsVgpTbkMwQ25iKzA4QzgrVXlReC9sbFdiZHpic09wcFVpeDJndk1lVUllcWJyRGdndTFhYzhvZ3dqUnVSZGlBTDBzCnlnQ3V0Q01ENGpCLzJkMHgyRkVndktiUmR1ZkQydGJPUTJQbmlxb0pyOUliWUhNSllXQUtES1dyQ1hjVitnNUoKay9oS1NDdG1Ecmd2R04zZ21iYXJndFdtRkdMYklHY0o5SHFmQWVMYzR0dUNaOTFpL0kxY3ZLVWFxZmRBSmtKZQpFYnJIdzBQdUdEeXdwcWZQbERDY2xKWCtjT09QVzRqQ0dXZDREYjQydmxTakNCZnE0WjBNQlgrK3U0UDh1eXg2ClNDMkxybEVkUk80M3lCL3J5QlIxN0lpb0NTeHVraGF2ZlpIYllQRXdyTFRZYXJ4OEZZNjJGVWdjZmZ6Y3NtaHgKWndJREFRQUIKLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tCg==",
      "signerType": "x509.pkcsv1",
      "signerMetadata": {
        "result": "success",
        "signerHash": "sha256",
        "reducerHash": "sha256"
      }
    }
  },
  {
    "unique": "01E7JQF86JMQJHBF4QX3A9X93K",
    "created": "2020-05-05T15:32:34.6429938Z",
    "identityCurrentType": "hash",
    "identityCurrent": {
      "hash": "Vc1WRGhjPV5SoXMDfRBazQMrAcvUVzAJi0+9x+hnjbE="
    },
    "identityPreviousType": "hash",
    "identityPrevious": null,
    "metadataType": "assessment",
    "metadata": {
      "provenance": {
        "Node": "evaluation"
      },
      "assessorType": "x509.pkcsv1",
      "assessorMetadata": {
        "result": "success",
        "validSignature": true,
        "unique": [
          "01E7JQF86JMQJHBF4QX1EFD6Y3"
        ]
      }
    }
  }
]
```

Full executable source is at [`cmd/examples/readme/main.go`](cmd/examples/readme/main.go).



## Project Structure

```
README.md                                This file
README.assets/                           Images and assets included in README.md

LICENSE                                  Project's license

cmd/                            
    examples/
        multistage/
            main.go                      Sample SDK usage (multiple stages)
        readme/
            main.go                      Sample SDK usage (as presented above)

internal/
    pkg/                        
        datetime/                        Date and time stamp implementation            
        test/                            Test-related implementation
            metadata/                    Annotation-specific assertions

pkg/
    annotation/                          Annotations
        metadata/                        Common annotation definitions
        store/                           Annotation store implementation
            contract.go                  Annotation store abstraction
            memory/                      In-process in-memory store implementation
        uniqueprovider/                  Unique provider
            contract.go                  Unique provider abstraction
            ulid/                        ULID-based implementation

    annotator/                           Annotators
        README.md                        Annotator documentation
        README.assets/                   Images and assets included in README.md
        contract.go                      Annotator abstraction

    hashprovider/                        Hash Provider (reduce data to unique hash)
        contract.go                      Hash provider abstraction
        md5/                             MD5-based implementation
        passthrough/                     passthrough implementation
        sha256/                          SHA256-based implementation

    identity/                            Identity
        contract.go                      Identity abstraction
        hash/                            Hash-based identity implementation

    identityprovider/                    Identity provider
        contract.go                      Identity provider abstraction
        hash/                            Hash-based identity provider implementation

    sdk/                                 Public SDK API
        close.go                         SDK Close() implementation
        create.go                        SDK Create() implementation
        mutate.go                        SDK Mutate() implementation
        sdk.go                           SDK factory function implementation

    status/
        contract.go                      Return value abstraction
```



## Opportunities for Improvement

The current SDK implementation is incomplete and limited.  It was intended as a simple reference implementation.  Opportunities for improvement exist -- in no particular order:

1. It is not optimized or performant.  
2. While part of its guiding vision, the current implementation does not implement trust scoring.  A scoring implementation would leverage the annotations created by the SDK. 
3. It does not include a non-transient annotation persistence implementation.  Only an example in-memory store service is provided.
4. While the annotation store contract implies immutability, there are no restrictions on implementation to enforce it. 
5. It does not version, sign, or encrypt individual annotations.  Annotations created by the current implementation are not signed or secured against tampering.
6. It does not conform to existing annotation standards.  SDK annotations use bespoke JSON schema.
7. It is currently limited to storing, retrieving, and processing only the annotations it originates. This precludes accessing and evaluating metadata originated and stored outside of Alvarium -- particularly limiting when Alvarium is not the primary annotation mechanism (as is expected in most use-cases).
8. It does not currently delineate -- must less track and attest -- authorship.  Authorship is the identity of the entity that recorded the annotation.  Authorship could be captured using the current SDK implementation as ad-hoc context.  This is less desirable/flexible than adding explicit authorship as an Alvarium annotation property.
9. It does not currently delineate -- much less track and attest -- data ownership.  Ownership is the identity of the entity currently responsible for the data and which acts as the source of truth for that data.  Ownership could be captured using the current SDK implementation as ad-hoc context.  This is less desirable/flexible than adding explicit ownership as an Alvarium annotation property.
10. It does not evaluate conformance of data to a specification.  It could be made to do so via a new purpose-specific annotator.
11. It does not evaluate data consistency to specified dynamic tolerances.  This would require definition of a data set -- a concept not currently recognized by the Alvarium architecture -- to bound evaluation across a subset of data.
12. It does not evaluate data consistency to specified static tolerances. It could be made to do so through a new purpose-specific annotator.
13. While Alvarium annotations are treated as immutable historical events and assessments can be successfully repeated (given an identical configuration of the SDK and no new related annotations have been added), Alvarium does not currently implement arbitrary assessment for a given point in time.
14. It implements trust factor assessment as a function via purpose-built annotators. The function implementing assessment is not portable.
15. It uses RSA PKCS#1 v1.5 for signatures and verification.  Version 1 has know flaws.  The SDK should be updated to use version 2 (i.e.  OAEP and PSS).
16. It does not include annotators that retrieve published annotations (from IPFS or IOTA Tangle).
17. The IOTA Tangle publisher implementation uses a single transaction for storage which limits maximum annotation storage size.  
18. The IOTA Tangle publisher implementation uses a unique address for each transaction.
