// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package v1schema

import (
	"time"

	"github.com/xmidt-org/webhook-schema/stream/sink/webhook"
)

// Schema is the schema for only v1 stream registration and used to unmarshal registration requests.
// Deprecated: This package should only be used for backwards compatibility
// reasons. Use v2 instead.
type Schema struct {
	// Address is the subscription request origin HTTP Address.
	Address string `json:"registered_from_address"`

	// Webhook contains data to inform how events are delivered.
	// Note, `json:"config"` is used for backwards compatibility reasons.
	// nolint:staticcheck
	Webhook webhook.V1Schema `json:"config"`

	// FailureURL is the URL used to notify subscribers when they've been cut off due to event overflow.
	// Optional, set to "" to disable notifications.
	FailureURL string `json:"failure_url"`

	// Events is the list of regular expressions to match an event type against.
	Events []string `json:"events"`

	// Matcher type contains values to match against the metadata.
	Matcher MetadataMatcherConfig `json:"matcher,omitempty"`

	// Duration describes how long the subscription lasts once added.
	Duration CustomDuration `json:"duration"`

	// Until describes the time this subscription expires.
	Until time.Time `json:"until"`
}

// MetadataMatcherConfig is Webhook substructure with config to match event metadata.
// Deprecated: This package should only be used for backwards compatibility
// reasons. Use v2 instead.
type MetadataMatcherConfig struct {
	// DeviceID is the list of regular expressions to match device id type against.
	DeviceID []string `json:"device_id"`
}
