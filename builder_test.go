// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package wrpeventstream

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/xmidt-org/webhook-schema/errortrace"
	"github.com/xmidt-org/webhook-schema/registry"
	"github.com/xmidt-org/webhook-schema/stream"
	"go.uber.org/multierr"
)

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

func Example() {
	var builder Builder
	if err := json.Unmarshal(jsonBuilder, &builder); err != nil {
		panic(err)
	}

	jsonStream := []byte(`{
		"metadata": {
			"version": "v0.1.0-alpha",
			"total_demos": 11
		},
		"flag": true,
		"someImplementationSpecificName": {
			"contact_info": {
				"name": "",
				"phone": "",
				"email": ""
			},
			"canonical_name": "com.examples.webhook-tests",
			"registered_from_address": "",
			"webhooks": [
				{
					"accept": "",
					"accept_encoding": "",
					"secret_hash": "",
					"payload_only": false,
					"receiver_urls": [
						"https://example.com",
						"https://example.com"
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
					"field": "canonical_name",
					"regex": "webpa"
				}
			],
			"expires": "0001-01-01T00:00:00Z"
			}
	}`)
	stream, err := New(Build(builder), ClientStream(&StreamDemo{}), UnmarshalJSON(jsonStream))
	if err != nil {
		panic(err)
	}

	jsonStream, _ = json.MarshalIndent(stream, "", "   ")
	fmt.Printf("Valid StreamDemo stream example\nstream:\n%s\n\n", jsonStream)

	// Convert a stream into a record to store it in a wrpeventstream registry (argus)

	record, err := stream.ToRecord()
	if err != nil {
		panic(err)
	}
	jsonRecord, _ := json.MarshalIndent(record, "", "   ")
	fmt.Printf("Convert StreamDemo stream into a record for storage in a wrpeventstream registry (argus):\n%s", jsonRecord)

	// Output: Valid StreamDemo stream example
	// stream:
	// {
	//    "someImplementationSpecificName": {
	//       "contact_info": {
	//          "name": "",
	//          "phone": "",
	//          "email": ""
	//       },
	//       "canonical_name": "com.examples.webhook-tests",
	//       "registered_from_address": "",
	//       "webhooks": [
	//          {
	//             "accept": "",
	//             "accept_encoding": "",
	//             "secret_hash": "",
	//             "payload_only": false,
	//             "receiver_urls": [
	//                "https://example.com",
	//                "https://example.com"
	//             ],
	//             "dns_srv_record": {
	//                "fqdns": null,
	//                "load_balancing_scheme": ""
	//             },
	//             "retry_hint": {
	//                "retry_each_url": 0,
	//                "max_retry": 0
	//             }
	//          }
	//       ],
	//       "hash": {
	//          "field": "",
	//          "regex": ""
	//       },
	//       "batch_hints": {
	//          "max_linger_duration": 0,
	//          "max_messages": 0
	//       },
	//       "failure_url": "",
	//       "matcher": [
	//          {
	//             "field": "canonical_name",
	//             "regex": "webpa"
	//          }
	//       ],
	//       "expires": "0001-01-01T00:00:00Z"
	//    },
	//    "metadata": {
	//       "version": "v0.1.0-alpha",
	//       "total_demos": 11
	//    },
	//    "flag": true
	// }
	//
	// Convert StreamDemo stream into a record for storage in a wrpeventstream registry (argus):
	// {
	//    "id": "3ec92689678b86b573efade5728343bc750adb34d81e42d173eceab78c7dbadb",
	//    "data": {
	//       "flag": true,
	//       "metadata": {
	//          "total_demos": 11,
	//          "version": "v0.1.0-alpha"
	//       },
	//       "someImplementationSpecificName": {
	//          "batch_hints": {
	//             "max_linger_duration": 0,
	//             "max_messages": 0
	//          },
	//          "canonical_name": "com.examples.webhook-tests",
	//          "contact_info": {
	//             "email": "",
	//             "name": "",
	//             "phone": ""
	//          },
	//          "expires": "0001-01-01T00:00:00Z",
	//          "failure_url": "",
	//          "hash": {
	//             "field": "",
	//             "regex": ""
	//          },
	//          "matcher": [
	//             {
	//                "field": "canonical_name",
	//                "regex": "webpa"
	//             }
	//          ],
	//          "registered_from_address": "",
	//          "webhooks": [
	//             {
	//                "accept": "",
	//                "accept_encoding": "",
	//                "dns_srv_record": {
	//                   "fqdns": null,
	//                   "load_balancing_scheme": ""
	//                },
	//                "payload_only": false,
	//                "receiver_urls": [
	//                   "https://example.com",
	//                   "https://example.com"
	//                ],
	//                "retry_hint": {
	//                   "max_retry": 0,
	//                   "retry_each_url": 0
	//                },
	//                "secret_hash": ""
	//             }
	//          ]
	//       }
	//    },
	//    "ttl": 0
	// }
}

func Example_v1Stream() {
	var builder Builder
	if err := json.Unmarshal(jsonBuilder, &builder); err != nil {
		panic(err)
	}

	jsonStream := []byte(`{
		"someImplementationSpecificName": {
			"registered_from_address": "",
			"config": {
				"url": "https://example.com",
				"content_type": ""
			},
			"failure_url": "",
			"events": [
				"event.*"
			],
			"matcher": {
				"device_id": null
			},
			"duration": "0s",
			"until": "0001-01-01T00:00:00Z"
			}
   	}`)
	stream, err := New(Build(builder), ClientStream(&StreamDemo{}), UnmarshalJSON(jsonStream))
	if err != nil {
		panic(err)
	}

	jsonStream, _ = json.MarshalIndent(stream, "", "   ")
	fmt.Printf("Valid v1 stream example\nstream:\n%s\n\n", jsonStream)

	// Convert a stream into a record to store it in a wrpeventstream registry (argus)

	record, err := stream.ToRecord()
	if err != nil {
		panic(err)
	}
	jsonRecord, _ := json.MarshalIndent(record, "", "   ")
	fmt.Printf("Convert StreamDemo stream into a record for storage in a wrpeventstream registry (argus):\n%s", jsonRecord)

	// Output: Valid v1 stream example
	// stream:
	// {
	//    "someImplementationSpecificName": {
	//       "registered_from_address": "",
	//       "config": {
	//          "url": "https://example.com",
	//          "content_type": ""
	//       },
	//       "failure_url": "",
	//       "events": [
	//          "event.*"
	//       ],
	//       "matcher": {
	//          "device_id": null
	//       },
	//       "duration": "0s",
	//       "until": "0001-01-01T00:00:00Z"
	//    },
	//    "metadata": {
	//       "version": "",
	//       "total_demos": 0
	//    },
	//    "flag": false
	// }
	//
	// Convert StreamDemo stream into a record for storage in a wrpeventstream registry (argus):
	// {
	//    "id": "100680ad546ce6a577f42f52df33b4cfdca756859e664b8d7de329b150d09ce9",
	//    "data": {
	//       "flag": false,
	//       "metadata": {
	//          "total_demos": 0,
	//          "version": ""
	//       },
	//       "someImplementationSpecificName": {
	//          "config": {
	//             "content_type": "",
	//             "url": "https://example.com"
	//          },
	//          "duration": "0s",
	//          "events": [
	//             "event.*"
	//          ],
	//          "failure_url": "",
	//          "matcher": {
	//             "device_id": null
	//          },
	//          "registered_from_address": "",
	//          "until": "0001-01-01T00:00:00Z"
	//       }
	//    },
	//    "ttl": 0
	// }
}

func Example_v2Stream() {
	var builder Builder
	if err := json.Unmarshal(jsonBuilder, &builder); err != nil {
		panic(err)
	}

	jsonStream := []byte(`{
		"metadata": {
			"version": "v0.1.0-alpha",
			"total_demos": 11
		},
		"flag": true,
		"someImplementationSpecificName": {
			"contact_info": {
				"name": "",
				"phone": "",
				"email": ""
			},
			"canonical_name": "com.examples.webhook-tests",
			"registered_from_address": "",
			"webhooks": [
				{
					"accept": "",
					"accept_encoding": "",
					"secret_hash": "",
					"payload_only": false,
					"receiver_urls": [
						"https://example.com",
						"https://example.com"
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
					"field": "canonical_name",
					"regex": "webpa"
				}
			],
			"expires": "0001-01-01T00:00:00Z"
			}
	}`)
	stream, err := New(Build(builder), ClientStream(&StreamDemo{}), UnmarshalJSON(jsonStream))
	if err != nil {
		return
	}

	jsonStream, _ = json.MarshalIndent(stream, "", "   ")
	fmt.Printf("Valid v2 stream example\nstream:\n%s\n\n", jsonStream)

	record, err := stream.ToRecord()
	if err != nil {
		panic(err)
	}
	jsonRecord, _ := json.MarshalIndent(record, "", "   ")
	fmt.Printf("Convert StreamDemo stream into a record for storage in a wrpeventstream registry (argus):\n%s", jsonRecord)

	// Output: Valid v2 stream example
	// stream:
	// {
	//    "someImplementationSpecificName": {
	//       "contact_info": {
	//          "name": "",
	//          "phone": "",
	//          "email": ""
	//       },
	//       "canonical_name": "com.examples.webhook-tests",
	//       "registered_from_address": "",
	//       "webhooks": [
	//          {
	//             "accept": "",
	//             "accept_encoding": "",
	//             "secret_hash": "",
	//             "payload_only": false,
	//             "receiver_urls": [
	//                "https://example.com",
	//                "https://example.com"
	//             ],
	//             "dns_srv_record": {
	//                "fqdns": null,
	//                "load_balancing_scheme": ""
	//             },
	//             "retry_hint": {
	//                "retry_each_url": 0,
	//                "max_retry": 0
	//             }
	//          }
	//       ],
	//       "hash": {
	//          "field": "",
	//          "regex": ""
	//       },
	//       "batch_hints": {
	//          "max_linger_duration": 0,
	//          "max_messages": 0
	//       },
	//       "failure_url": "",
	//       "matcher": [
	//          {
	//             "field": "canonical_name",
	//             "regex": "webpa"
	//          }
	//       ],
	//       "expires": "0001-01-01T00:00:00Z"
	//    },
	//    "metadata": {
	//       "version": "v0.1.0-alpha",
	//       "total_demos": 11
	//    },
	//    "flag": true
	// }
	//
	// Convert StreamDemo stream into a record for storage in a wrpeventstream registry (argus):
	// {
	//    "id": "3ec92689678b86b573efade5728343bc750adb34d81e42d173eceab78c7dbadb",
	//    "data": {
	//       "flag": true,
	//       "metadata": {
	//          "total_demos": 11,
	//          "version": "v0.1.0-alpha"
	//       },
	//       "someImplementationSpecificName": {
	//          "batch_hints": {
	//             "max_linger_duration": 0,
	//             "max_messages": 0
	//          },
	//          "canonical_name": "com.examples.webhook-tests",
	//          "contact_info": {
	//             "email": "",
	//             "name": "",
	//             "phone": ""
	//          },
	//          "expires": "0001-01-01T00:00:00Z",
	//          "failure_url": "",
	//          "hash": {
	//             "field": "",
	//             "regex": ""
	//          },
	//          "matcher": [
	//             {
	//                "field": "canonical_name",
	//                "regex": "webpa"
	//             }
	//          ],
	//          "registered_from_address": "",
	//          "webhooks": [
	//             {
	//                "accept": "",
	//                "accept_encoding": "",
	//                "dns_srv_record": {
	//                   "fqdns": null,
	//                   "load_balancing_scheme": ""
	//                },
	//                "payload_only": false,
	//                "receiver_urls": [
	//                   "https://example.com",
	//                   "https://example.com"
	//                ],
	//                "retry_hint": {
	//                   "max_retry": 0,
	//                   "retry_each_url": 0
	//                },
	//                "secret_hash": ""
	//             }
	//          ]
	//       }
	//    },
	//    "ttl": 0
	// }
}

func Example_v1StreamFailure() {
	var builder Builder
	if err := json.Unmarshal(jsonBuilder, &builder); err != nil {
		panic(err)
	}

	jsonStream := []byte(`{
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

	stream, err := New(Build(builder), ClientStream(&StreamDemo{}), UnmarshalJSON(jsonStream))
	warnings, errs := getErrors(err)
	jsonStream, _ = json.MarshalIndent(stream, "", "   ")
	fmt.Printf("invalid v1 stream\nwarnings:\n%+v\n\nerrors:\n%+v\nstream:\n%s\n", warnings, errs, jsonStream)

	// Output:invalid v1 stream
	// warnings:
	// the following errors occurred:
	//  -  [err="option failure; invaild `*v1schema.Schema` stream, receiver_url validator error", level=warning, tag="receiver_url"]
	//  -  [err="option failure; invaild `*v1schema.Schema` stream, receiver_url validator error", level=warning, tag="receiver_url"]
	//
	// errors:
	// the following errors occurred:
	//  -  [err="failed to build stream", level=error, tag="wrpeventstream.New"]
	//  -  [err="failed to build stream", level=error, tag="wrpeventstream.Builder.Build"]
	//  -  [err="all supported stream version unmarshals failed", level=error, tag="wrpeventstream.Builder.unmarshaller"]
	//  -  [err="client stream config `*wrpeventstream.StreamDemo` is not a match to `v1schema.Schema`: unmarshalling failure", level=error, tag="*v1.manifest.UnmarshalJSON"]
	//  -  [err="failed to json unmarshal a wrpeventstream configuration", level=error, tag="UnmarshalJSON"]
	//  -  [err="`v1schema.Schema` unmarshal error", level=error, tag="v1.manifest.UnmarshalJSON"]
	//  -  [err="`v1schema.Schema` failed validation", level=error, tag="v1.manifest.Validate"]
	//  -  [err="option failure; invaild `*v1schema.Schema` stream, event_regex validator error", level=error, tag="event_regex"]
	//  -  [err="option failure; invaild `*v1schema.Schema` stream, event_regex validator error", level=error, tag="event_regex"]
	//  -  [err="client stream config `*wrpeventstream.StreamDemo` is not a match to `v2schema.Schema`: unmarshalling failure", level=error, tag="*v2.manifest.UnmarshalJSON"]
	//  -  [err="failed to json unmarshal a wrpeventstream configuration", level=error, tag="UnmarshalJSON"]
	//  -  [err="`v2schema.Schema` unmarshal error", level=error, tag="v2.manifest.UnmarshalJSON"]
	// stream:
	// null
}

func Example_v2StreamFailure() {
	var builder Builder
	if err := json.Unmarshal(jsonBuilder, &builder); err != nil {
		panic(err)
	}

	jsonStream := []byte(`{
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

	stream, err := New(Build(builder), ClientStream(&StreamDemo{}), UnmarshalJSON(jsonStream))
	warnings, errs := getErrors(err)
	jsonStream, _ = json.MarshalIndent(stream, "", "   ")
	fmt.Printf("invalid v2 stream\nwarnings:\n%+v\n\nerrors:\n%+v\nstream:\n%s\n", warnings, errs, jsonStream)

	// Output:invalid v2 stream
	// warnings:
	// the following errors occurred:
	//  -  [err="option failure; invaild `*v2schema.Schema` stream, receiver_url validator error", level=warning, tag="receiver_url"]
	//  -  [err="option failure; invaild `*v2schema.Schema` stream, receiver_url validator error", level=warning, tag="receiver_url"]
	//
	// errors:
	// the following errors occurred:
	//  -  [err="failed to build stream", level=error, tag="wrpeventstream.New"]
	//  -  [err="failed to build stream", level=error, tag="wrpeventstream.Builder.Build"]
	//  -  [err="all supported stream version unmarshals failed", level=error, tag="wrpeventstream.Builder.unmarshaller"]
	//  -  [err="client stream config `*wrpeventstream.StreamDemo` is not a match to `v1schema.Schema`: unmarshalling failure", level=error, tag="*v1.manifest.UnmarshalJSON"]
	//  -  [err="failed to json unmarshal a wrpeventstream configuration", level=error, tag="UnmarshalJSON"]
	//  -  [err="`v1schema.Schema` unmarshal error", level=error, tag="v1.manifest.UnmarshalJSON"]
	//  -  [err="client stream config `*wrpeventstream.StreamDemo` is not a match to `v2schema.Schema`: unmarshalling failure", level=error, tag="*v2.manifest.UnmarshalJSON"]
	//  -  [err="failed to json unmarshal a wrpeventstream configuration", level=error, tag="UnmarshalJSON"]
	//  -  [err="`v2schema.Schema` unmarshal error", level=error, tag="v2.manifest.UnmarshalJSON"]
	//  -  [err="`v2schema.Schema` failed validation", level=error, tag="v2.manifest.Validate"]
	//  -  [err="option failure; invaild `*v2schema.Schema` stream, not_empty validator error", level=error, tag="not_empty"]
	//  -  [err="option failure; invaild `*v2schema.Schema` stream, event_regex validator error", level=error, tag="event_regex"]
	//  -  [err="option failure; invaild `*v2schema.Schema` stream, event_regex validator error", level=error, tag="event_regex"]
	// stream:
	// null
}

func getErrors(actualErr error) (warnings, errs error) {
	errss := multierr.Errors(actualErr)
	for _, err := range errss {
		oerr, ok := err.(errortrace.Trace)
		if !ok {
			continue
		}
		switch oerr.Level() {
		case errortrace.WarningLevel:
			warnings = multierr.Append(warnings, oerr)
		case errortrace.ErrorLevel:
			errs = multierr.Append(errs, oerr)
		}
	}

	return warnings, errs
}
