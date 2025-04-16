// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package webhook

import (
	"iter"

	"github.com/xmidt-org/webhook-schema/stream/sink"
)

func IterateV2Manifest(e []V2Schema) iter.Seq2[int, sink.Manifest] {
	return func(yield func(int, sink.Manifest) bool) {
		for i := 0; i <= len(e)-1; i++ {
			if !yield(i, v2Manifest{V2Schema: e[i]}) {
				return
			}
		}
	}
}

type v2Manifest struct {
	V2Schema
}

func (v2m v2Manifest) GetName() string {
	return "webhookV2"
}

func (v2m v2Manifest) GetAcceptType() string {
	return v2m.Accept
}

func (v2m v2Manifest) GetSink() any {
	return v2m.V2Schema
}

func (v2m v2Manifest) GetReceiverURLs() []string {
	return v2m.ReceiverURLs
}

func (v2m v2Manifest) GetSecret() string {
	return v2m.Secret
}
