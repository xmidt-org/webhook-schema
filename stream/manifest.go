// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package stream

import (
	"encoding/json"
	"iter"
	"time"

	"github.com/xmidt-org/webhook-schema/stream/sink"
)

// Manifest defines the functionality required to access information to standup streams for event processing.
// Manifest is stream version and sink agnostic until `Manifest.GetSinks()` or `Manifest.GetStream()` are called.
type Manifest interface {
	json.Unmarshaler
	json.Marshaler

	GetID() string
	GetFailureUrl() string
	GetExpired() time.Time
	GetTTLSeconds(func() time.Time) int64
	GetBatchMaxMesasges() int
	GetBatchLinger() time.Duration
	ObfuscateSecrets()
	Validate() error
	SetDefaults() error
	GetStream() any
	GetSinks() (iter.Seq2[int, sink.Manifest], error)
	GetName() string
}
