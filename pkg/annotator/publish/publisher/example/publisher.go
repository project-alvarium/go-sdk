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

package example

import (
	"encoding/json"
	"io"

	"github.com/project-alvarium/go-sdk/pkg/annotation"
	"github.com/project-alvarium/go-sdk/pkg/annotation/metadata"
	examplePublisherMetadata "github.com/project-alvarium/go-sdk/pkg/annotator/publish/publisher/example/metadata"
)

// publisher is a receiver that encapsulates required dependencies.
type publisher struct {
	writer io.Writer
}

// New is a factory function that returns an initialized publisher.
func New(writer io.Writer) *publisher {
	return &publisher{
		writer: writer,
	}
}

// SetUp is called once when the publisher is instantiated.
func (p *publisher) SetUp() {}

// TearDown is called once when publisher is terminated.
func (p *publisher) TearDown() {}

// Format converts struct into []byte; separated into own function for use by unit tests.
func (p *publisher) Format(s interface{}) []byte {
	b, _ := json.MarshalIndent(s, "", "  ")
	return b
}

// Publish retrieves and "publishes" annotations.
func (p *publisher) Publish(annotations []*annotation.Instance) metadata.Contract {
	if _, err := p.writer.Write(p.Format(annotations)); err != nil {
		return examplePublisherMetadata.NewFailure(err.Error())
	}
	return examplePublisherMetadata.NewSuccess()
}

// Failure creates a publisher-specific failure annotation.
func (p *publisher) Failure(errorMessage string) metadata.Contract {
	return examplePublisherMetadata.NewFailure(errorMessage)
}

// Kind returns an implementation mnemonic.
func (*publisher) Kind() string {
	return examplePublisherMetadata.Kind
}
