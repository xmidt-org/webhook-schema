// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package errortrace_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	wrpeventstream "github.com/xmidt-org/webhook-schema"
	errortrace "github.com/xmidt-org/webhook-schema/errortrace"
	"github.com/xmidt-org/webhook-schema/registry"
	"github.com/xmidt-org/webhook-schema/stream"
	"go.uber.org/multierr"
)

var (
	errRequest      = errors.New("request `/api/examples` failed")
	errRequestTrace = errortrace.New(errRequest,
		errortrace.Level(errortrace.ErrorLevel),
		errortrace.Tag("request"))
	// nolint:unused
	errProcessData      = errors.New("data processing failure")
	errProcessDataTrace = errortrace.New(errRequest,
		errortrace.Level(errortrace.ErrorLevel),
		errortrace.Tag("processData"))
	// nolint:unused
	errTransformerA      = errors.New("data transformer A failed")
	errTransformerATrace = errortrace.New(errRequest,
		errortrace.Level(errortrace.ErrorLevel),
		errortrace.Tag("transformerA"))
	// nolint:unused
	errTransformerB      = errors.New("data transformer B failed")
	errTransformerBTrace = errortrace.New(errRequest,
		errortrace.Level(errortrace.InfoLevel),
		errortrace.Tag("transformerB"))
	// nolint:unused
	errTransformerC      = errors.New("data transformer C failed")
	errTransformerCTrace = errortrace.New(errRequest,
		errortrace.Level(errortrace.WarningLevel),
		errortrace.Tag("transformerC"))
)

func encode(any) (errs error) {
	return errors.New("response encoding failed")
}

func request(data any) (errs error) {
	errs = multierr.Append(errs, extractMetadata(data))
	errs = multierr.Append(errs, processData(data))
	if errs != nil {
		return errRequestTrace.AppendDetail(errs)
	}

	return nil
}

func processData(data any) (errs error) {
	transformers := []func(any) error{transformerA, transformerB, transformerC}
	for _, t := range transformers {
		errs = multierr.Append(errs, t(data))
	}

	if errs != nil {
		return errProcessDataTrace.AppendDetail(errs)
	}

	return nil
}

func transformerA(any) error {
	errs := []error{
		errors.New("error 1"),
		errors.New("error 2"),
		errors.New("error ..."),
		errors.New("error 1020"),
	}

	return errTransformerATrace.AppendDetail(errs...)
}
func transformerB(any) error {
	err := errors.New("some randome error")

	return errTransformerBTrace.AppendDetail(err)
}

func transformerC(any) error {
	err := errors.New("some randome error")

	return errTransformerCTrace.AppendDetail(err)
}

func extractMetadata(any) error {
	return errors.New("failed to extract metadata from request")
}

func ExampleTraceError() {
	var errs error
	errs = multierr.Append(errs, request(struct{}{}))
	errs = multierr.Append(errs, encode(struct{}{}))
	trace := errortrace.TraceError(errs)

	fmt.Printf("Error Trace:\n%s", trace)

	// Output: Error Trace:
	// [START OF ERROR TRACE]
	// |-err0->[err="request `/api/examples` failed", level=error, tag="request"]
	// |   |-err0.0->[err="failed to extract metadata from request"]
	// |   |-err0.1->[err="request `/api/examples` failed", level=error, tag="processData"]
	// |   |   |-err0.1.0->[err="request `/api/examples` failed", level=error, tag="transformerA"]
	// |   |   |   |-err0.1.0.0->[err="error 1"]
	// |   |   |   |-err0.1.0.1->[err="error 2"]
	// |   |   |   |-err0.1.0.2->[err="error ..."]
	// |   |   |   |-err0.1.0.3->[err="error 1020"]
	// |   |   |-err0.1.1->[err="request `/api/examples` failed", level=info, tag="transformerB"]
	// |   |   |   |-err0.1.1.0->[err="some randome error"]
	// |   |   |-err0.1.2->[err="request `/api/examples` failed", level=warning, tag="transformerC"]
	// |   |   |   |-err0.1.2.0->[err="some randome error"]
	// |-err.1->[err="response encoding failed"]
	// |
	// |
	// [END OF ERROR TRACE]
}

var (
	jsonBuilder = []byte(`{
		"v1": {
			"disable": false,
			"options": [
				{
					"type": "receiver_url",
					"level": "warning",
					"url_checker": {
						"domain": {
							"allow_special_use_domains": true
						},
						"url": {
							"schemes": ["https"]
						}
					}
				},
				{
					"type": "event_regex",
					"level": "error"
				},
				{
					"type": "until",
					"level": "error",
					"disable": true
				}
			]
		},
		"v2": {
			"disable": false,
			"options": [
				{
					"type": "receiver_url",
					"level": "warning",
					"url_checker": {
						"domain": {
							"allow_special_use_domains": true
						},
						"url": {
							"schemes": ["https"]
						}
					}
				},
				{
					"type": "event_regex",
					"level": "error"
				},
				{
					"type": "expires",
					"level": "error",
					"disable": true
				}
			]
		}
	}`)
)

type RandomStream struct {
	StreamDemo
}

type StreamDemo struct {
	Stream stream.Manifest `json:"someImplementationSpecificName"`

	OtherMetadata struct {
		DemoV1        string `json:"version"`
		NumberOfDemos int    `json:"total_demos"`
	} `json:"metadata"`

	SomeFlag bool `json:"flag"`
}

func (s *StreamDemo) ToRecord() (registry.Record, error) {
	return registry.New(s)
}

func (s *StreamDemo) SetStream(stream stream.Manifest) {
	s.Stream = stream
}

func (s *StreamDemo) GetStream() stream.Manifest {
	return s.Stream
}

func (s *StreamDemo) GetID() string {
	return s.Stream.GetID()
}

func (s *StreamDemo) GetTTLSeconds(now func() time.Time) int64 {
	return s.Stream.GetTTLSeconds(now)
}

func (s *StreamDemo) ObfuscateSecrets() {
	s.Stream.ObfuscateSecrets()
}

func (s *StreamDemo) Validate() error {
	return s.Stream.Validate()
}

func ExampleTraceError_wrpeventstream() {
	var builder wrpeventstream.Builder
	if err := json.Unmarshal(jsonBuilder, &builder); err != nil {
		panic(err)
	}

	jsonV1Stream := []byte(`{
		"someImplementationSpecificName": {
			"registered_from_address": "",
			"config": {
				"url": "http://example.com",
				"content_type": ""
			},
			"failure_url": "",
			"events": [
				"("
			],
			"matcher": {
				"device_id": null
			},
			"duration": "0s",
			"until": "0001-01-01T00:00:00Z"
			}
   	}`)

	jsonV2Stream := []byte(`{
		"someImplementationSpecificName": {
			"contact_info": {
				"name": "",
				"phone": "",
				"email": ""
			},
			"canonical_name": "",
			"registered_from_address": "",
			"webhooks": [
				{
					"accept": "",
					"accept_encoding": "",
					"secret_hash": "",
					"payload_only": false,
					"receiver_urls": [
						"https://example.com",
						"http://example.com"
					],
					"dns_srv_record": {
						"fqdns": null,
						"load_balancing_scheme": ""
					},
					"retry_hint": {
						"retry_each_url": 0,
						"max_retry": 0
					}
				}
			],
			"kafkas": null,
			"hash": {
				"field": "",
				"regex": ""
			},
			"batch_hints": {
				"max_linger_duration": 0,
				"max_messages": 0
			},
			"failure_url": "",
			"matcher": [
				{
					"field": "",
					"regex": "("
				}
			],
			"expires": "0001-01-01T00:00:00Z"
			}
   	}`)

	_, err1 := wrpeventstream.New(wrpeventstream.Build(builder), wrpeventstream.ClientStream(&StreamDemo{}), wrpeventstream.UnmarshalJSON(jsonV1Stream))
	_, err2 := wrpeventstream.New(wrpeventstream.Build(builder), wrpeventstream.ClientStream(&RandomStream{}), wrpeventstream.UnmarshalJSON(jsonV2Stream))

	trace := errortrace.TraceError(errors.Join(err1, err2))

	fmt.Printf("Error Trace:\n%s", trace)

	// Output: Error Trace:
	// [START OF ERROR TRACE]
	// |-err0->[err="failed to build stream", level=error, tag="wrpeventstream.New"]
	// |   |-err0.0->[err="failed to build stream", level=error, tag="wrpeventstream.Builder.Build"]
	// |   |   |-err0.0.0->[err="all supported stream version unmarshals failed", level=error, tag="wrpeventstream.Builder.unmarshaller"]
	// |   |   |   |-err0.0.0.0->[err="client stream config `*errortrace_test.StreamDemo` is not a match to `v1schema.Schema`: unmarshalling failure", level=error, tag="*v1.manifest.UnmarshalJSON"]
	// |   |   |   |   |-err0.0.0.0.0->[err="failed to json unmarshal a wrpeventstream configuration", level=error, tag="UnmarshalJSON"]
	// |   |   |   |   |   |-err0.0.0.0.0.0->[err="`v1schema.Schema` unmarshal error", level=error, tag="v1.manifest.UnmarshalJSON"]
	// |   |   |   |   |   |   |-err0.0.0.0.0.0.0->[err="`v1schema.Schema` failed validation", level=error, tag="v1.manifest.Validate"]
	// |   |   |   |   |   |   |   |-err0.0.0.0.0.0.0.0->[err="option failure; invaild `*v1schema.Schema` stream, receiver_url validator error", level=warning, tag="receiver_url"]
	// |   |   |   |   |   |   |   |   |-err0.0.0.0.0.0.0.0.0->[err="invalid ReceiverURL: `http://example.com`: scheme not allowed"]
	// |   |   |   |   |   |   |   |-err0.0.0.0.0.0.0.1->[err="option failure; invaild `*v1schema.Schema` stream, event_regex validator error", level=error, tag="event_regex"]
	// |   |   |   |   |   |   |   |   |-err0.0.0.0.0.0.0.1.0->[err="failed to compile Event regexp: `(`: error parsing regexp: missing closing ): `(`"]
	// |   |   |   |   |   |   |   |-err0.0.0.0.0.0.0.2->[err="option failure; invaild `*v1schema.Schema` stream, receiver_url validator error", level=warning, tag="receiver_url"]
	// |   |   |   |   |   |   |   |   |-err0.0.0.0.0.0.0.2.0->[err="invalid ReceiverURL: `http://example.com`: scheme not allowed"]
	// |   |   |   |   |   |   |   |-err0.0.0.0.0.0.0.3->[err="option failure; invaild `*v1schema.Schema` stream, event_regex validator error", level=error, tag="event_regex"]
	// |   |   |   |   |   |   |   |   |-err0.0.0.0.0.0.0.3.0->[err="failed to compile Event regexp: `(`: error parsing regexp: missing closing ): `(`"]
	// |   |   |   |-err0.0.0.1->[err="client stream config `*errortrace_test.StreamDemo` is not a match to `v2schema.Schema`: unmarshalling failure", level=error, tag="*v2.manifest.UnmarshalJSON"]
	// |   |   |   |   |-err0.0.0.1.0->[err="failed to json unmarshal a wrpeventstream configuration", level=error, tag="UnmarshalJSON"]
	// |   |   |   |   |   |-err0.0.0.1.0.0->[err="`v2schema.Schema` unmarshal error", level=error, tag="v2.manifest.UnmarshalJSON"]
	// |   |   |   |   |   |   |-err0.0.0.1.0.0.0->[err="json: cannot unmarshal object into Go struct field Schema.matcher of type []v2schema.FieldRegex"]
	// |-err1->[err="failed to build stream", level=error, tag="wrpeventstream.New"]
	// |   |-err1.0->[err="failed to build stream", level=error, tag="wrpeventstream.Builder.Build"]
	// |   |   |-err1.0.0->[err="all supported stream version unmarshals failed", level=error, tag="wrpeventstream.Builder.unmarshaller"]
	// |   |   |   |-err1.0.0.0->[err="client stream config `*errortrace_test.RandomStream` is not a match to `v1schema.Schema`: unmarshalling failure", level=error, tag="*v1.manifest.UnmarshalJSON"]
	// |   |   |   |   |-err1.0.0.0.0->[err="failed to json unmarshal a wrpeventstream configuration", level=error, tag="UnmarshalJSON"]
	// |   |   |   |   |   |-err1.0.0.0.0.0->[err="`v1schema.Schema` unmarshal error", level=error, tag="v1.manifest.UnmarshalJSON"]
	// |   |   |   |   |   |   |-err1.0.0.0.0.0.0->[err="json: cannot unmarshal array into Go struct field Schema.matcher of type v1schema.MetadataMatcherConfig"]
	// |   |   |   |-err1.0.0.1->[err="client stream config `*errortrace_test.RandomStream` is not a match to `v2schema.Schema`: unmarshalling failure", level=error, tag="*v2.manifest.UnmarshalJSON"]
	// |   |   |   |   |-err1.0.0.1.0->[err="failed to json unmarshal a wrpeventstream configuration", level=error, tag="UnmarshalJSON"]
	// |   |   |   |   |   |-err1.0.0.1.0.0->[err="`v2schema.Schema` unmarshal error", level=error, tag="v2.manifest.UnmarshalJSON"]
	// |   |   |   |   |   |   |-err1.0.0.1.0.0.0->[err="`v2schema.Schema` failed validation", level=error, tag="v2.manifest.Validate"]
	// |   |   |   |   |   |   |   |-err1.0.0.1.0.0.0.0->[err="option failure; invaild `*v2schema.Schema` stream, not_empty validator error", level=error, tag="not_empty"]
	// |   |   |   |   |   |   |   |   |-err1.0.0.1.0.0.0.0.0->[err="empty schema error"]
	// |   |   |   |   |   |   |   |-err1.0.0.1.0.0.0.1->[err="option failure; invaild `*v2schema.Schema` stream, receiver_url validator error", level=warning, tag="receiver_url"]
	// |   |   |   |   |   |   |   |   |-err1.0.0.1.0.0.0.1.0->[err="invalid ReceiverURL: webhook 0 `http://example.com`:scheme not allowed"]
	// |   |   |   |   |   |   |   |-err1.0.0.1.0.0.0.2->[err="option failure; invaild `*v2schema.Schema` stream, event_regex validator error", level=error, tag="event_regex"]
	// |   |   |   |   |   |   |   |   |-err1.0.0.1.0.0.0.2.0->[err="failed to compile Matcher regexp: `(`: error parsing regexp: missing closing ): `(`"]
	// |   |   |   |   |   |   |   |-err1.0.0.1.0.0.0.3->[err="option failure; invaild `*v2schema.Schema` stream, receiver_url validator error", level=warning, tag="receiver_url"]
	// |   |   |   |   |   |   |   |   |-err1.0.0.1.0.0.0.3.0->[err="invalid ReceiverURL: webhook 0 `http://example.com`:scheme not allowed"]
	// |   |   |   |   |   |   |   |-err1.0.0.1.0.0.0.4->[err="option failure; invaild `*v2schema.Schema` stream, event_regex validator error", level=error, tag="event_regex"]
	// |   |   |   |   |   |   |   |   |-err1.0.0.1.0.0.0.4.0->[err="failed to compile Matcher regexp: `(`: error parsing regexp: missing closing ): `(`"]
	// |
	// |
	// [END OF ERROR TRACE]
}
