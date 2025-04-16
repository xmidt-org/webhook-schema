// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package webhook

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

func (v1m v1Manifest) GetName() string {
	return "webhookV1"
}

func (v1m v1Manifest) GetSink() any {
	return v1m.Sink
}

// Note, V1Schema.ContentType is known to be poorly named and it's actually
// used as AcceptType.
func (v1m v1Manifest) GetAcceptType() string {
	return v1m.Sink.ContentType
}

func (v1m v1Manifest) GetReceiverURLs() []string {
	return append([]string{v1m.Sink.ReceiverURL}, v1m.Sink.AlternativeURLs...)
}

func (v1m v1Manifest) GetReceiverURL() string {
	return v1m.Sink.ReceiverURL
}

func (v1m v1Manifest) GetAlternativeURLs() []string {
	return v1m.Sink.AlternativeURLs
}

func (v1m v1Manifest) GetSecret() string {
	return v1m.Sink.Secret
}
