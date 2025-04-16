// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package v2schema

import (
	"time"

	"github.com/xmidt-org/webhook-schema/stream/sink/kafka"
	"github.com/xmidt-org/webhook-schema/stream/sink/webhook"
)

// Schema is the schema for only v2 stream registration and used to unmarshal registration requests.
type Schema struct {
	// ContactInfo contains contact information used to reach the owner of the registration.
	// (Optional).
	ContactInfo ContactInfo `json:"contact_info"`

	// CanonicalName is the canonical name of the registration request.
	// Reusing a CanonicalName will override the configurations set in that previous
	// registration request with the same CanonicalName.
	CanonicalName string `json:"canonical_name"`

	// Address is the subscription request origin HTTP Address.
	Address string `json:"registered_from_address"`

	// Webhooks contains data to inform how events are delivered to multiple urls.
	Webhooks []webhook.V2Schema `json:"webhooks,omitempty"`

	// Kafkas contains data to inform how events are delivered to multiple kafkas.
	Kafkas []kafka.V1Schema `json:"kafkas,omitempty"`

	// Hash is a substructure for configuration related to distributing events among sinks.
	// Note. Any failures due to a bad regex feild or regex expression will result in a silent failure.
	Hash FieldRegex `json:"hash"`

	// BatchHint is the substructure for configuration related to event batching.
	// (Optional, if omited then batches of singal events will be sent)
	// Default value will disable batch. All zeros will also disable batch.
	BatchHint BatchHint `json:"batch_hints"`

	// FailureURL is the URL used to notify subscribers when they've been cut off due to event overflow.
	// Optional, set to "" to disable notifications.
	FailureURL string `json:"failure_url"`

	// Matcher is the list of regular expressions to match incoming events against to.
	// Note. Any failures due to a bad regex field or regex expression will result in a silent failure.
	Matcher []FieldRegex `json:"matcher,omitempty"`

	// Expires describes the time this subscription expires.
	// TODO: list of supported formats
	Expires time.Time `json:"expires"`
}

type ContactInfo struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

type BatchHint struct {
	// MaxLingerDuration is the maximum delay for batching if MaxMesasges has not been reached.
	// Default value will set no maximum value.
	MaxLingerDuration time.Duration `json:"max_linger_duration"`
	// MaxMesasges is the maximum number of events that will be sent in a single batch.
	// Default value will set no maximum value.
	MaxMesasges int `json:"max_messages"`
}

// FieldRegex is a substructure with data related to regular expressions.
type FieldRegex struct {
	// Field is the wrp field to be used for regex.
	// All wrp field can be used, refer to the schema for examples.
	Field string `json:"field"`

	// FieldRegex is the regular expression to match `Field` against to.
	Regex string `json:"regex"`
}
