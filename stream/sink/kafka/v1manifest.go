// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package kafka

import (
	"iter"

	"github.com/xmidt-org/webhook-schema/stream/sink"
)

func IterateV1Manifest(e []V1Schema) iter.Seq2[int, sink.Manifest] {
	return func(yield func(int, sink.Manifest) bool) {
		for i := 0; i <= len(e)-1; i++ {
			if !yield(i, v1Manifest{Sink: e[i]}) {
				return
			}
		}
	}
}

type v1Manifest struct {
	Sink V1Schema
}

func (v2m v1Manifest) GetName() string {
	return "kafkaV1"
}

func (v1m v1Manifest) GetSink() any {
	return v1m.Sink
}

func (v1m v1Manifest) GetAcceptType() string {
	return v1m.Sink.Accept
}

func (v1m v1Manifest) GetBootstrapServers() []string {
	return v1m.Sink.BootstrapServers
}

func (v1m v1Manifest) GetKafkaProducer() any {
	return v1m.Sink.KafkaProducer
}
