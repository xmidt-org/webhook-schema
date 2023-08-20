// SPDX-FileCopyrightText: 2022 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package webhook

import (
	"fmt"
	"time"
)

var (
	ErrInvalidInput = fmt.Errorf("invalid input")
)

// DeliveryConfig is a Webhook substructure with data related to event delivery.
type DeliveryConfig struct {
	// URL is the HTTP URL to deliver messages to.
	ReceiverURL string `json:"url"`

	// ContentType is content type value to set WRP messages to (unless already specified in the WRP).
	ContentType string `json:"content_type"`

	// Secret is the string value for the SHA1 HMAC.
	// (Optional, set to "" to disable behavior).
	Secret string `json:"secret,omitempty"`

	// AlternativeURLs is a list of explicit URLs that should be round robin through on failure cases to the main URL.
	AlternativeURLs []string `json:"alt_urls,omitempty"`
}

// MetadataMatcherConfig is Webhook substructure with config to match event metadata.
type MetadataMatcherConfig struct {
	// DeviceID is the list of regular expressions to match device id type against.
	DeviceID []string `json:"device_id"`
}

// Registration is a special struct for unmarshaling a webhook as part of
// a webhook registration request.  The only difference between this struct and
// the Webhook struct is the Duration field.
type Registration struct {
	// Address is the subscription request origin HTTP Address.
	Address string `json:"registered_from_address"`

	// Config contains data to inform how events are delivered.
	Config DeliveryConfig `json:"config"`

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

	// now is a function that returns the current time.  It is used for testing.
	nowFunc func() time.Time `json:"-"`
}

type Option interface {
	fmt.Stringer
	Validate(*Registration) error
}

// Validate is a method on Registration that validates the registration
// against a list of options.
func (r *Registration) Validate(opts ...Option) error {
	for _, opt := range opts {
		if opt != nil {
			if err := opt.Validate(r); err != nil {
				return err
			}
		}
	}
	return nil
}
