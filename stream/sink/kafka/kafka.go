// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package kafka

import "github.com/xmidt-org/webhook-schema/stream/sink"

type Manifest interface {
	sink.Manifest
	GetBootstrapServers() []string
	// TODO we need to define KafkaProducer
	GetKafkaProducer() any
}
