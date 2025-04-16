// SPDX-FileCopyrightText: 2022 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package kafka

// V1Schema is a substructure with data related to event delivery.
type V1Schema struct {
	// Accept is the encoding type of outgoing events. The following encoding types are supported, otherwise
	// a 406 response code is returned: application/octet-stream, application/json, application/jsonl, application/msgpack.
	// Note: An `Accept` of application/octet-stream or application/json will result in a single response for batch sizes of 0 or 1
	// and batch sizes greater than 1 will result in a multipart response. An `Accept` of application/jsonl or application/msgpack
	// will always result in a single response with a list of batched events for any batch size.
	Accept string `json:"accept"`

	// BootstrapServers is a list of kafka broker addresses.
	BootstrapServers []string `json:"bootstrap_servers"`

	// TODO: figure out which kafka configuration substructures we want to expose to users (to be set by users)
	// going to be based on https://pkg.go.dev/github.com/IBM/sarama#Config
	// this substructures also includes auth related secrets, noted `MaxOpenRequests` will be excluded since it's already exposed
	KafkaProducer V1KafkaProducer `json:"kafka_producer"`
}

type V1KafkaProducer struct {
	// Secret is the string value.
	// (Optional, set to "" to disable behavior).
	Secret string `json:"secret,omitempty"`
}
