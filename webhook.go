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

// WebhookConfig is a Webhook substructure with data related to event delivery.
type WebhookConfig struct {
	// URL is the HTTP URL to deliver messages to.
	ReceiverURL string `json:"url"`

	// ContentType is content type value to set WRP messages to (unless already specified in the WRP).
	ContentType string `json:"content_type"`

	// Secret is the string value for the SHA1 HMAC.
	// (Optional, set to "" to disable behavior).
	Secret string `json:"secret,omitempty"`

	// AlternativeURLs is a list of explicit URLs that should be round robin through on failure cases to the main URL.
	AlternativeURLs []string `json:"alt_urls,omitempty"`

	// ID is the unique identifier used for hashing and internal identification.
	ID string `json:"id"`

	// SecretHash is the hash algorithm to be used. Only SHA256 HMAC and SHA512 HMAC are supported.
	// (Optional).
	SecretHash string `json:"secret_hash"`

	// Batch is the substructure for configuration related to event batching.
	// (Optional, if omited then batches of singal events will be sent)
	Batch struct {
		// MaxLinger is the maximum delay for batching if MaxMesasges has not been reached.
		MaxLinger float32 `json:"max_linger"`
		// MaxMesasges is the maximum number of events that will be sent in a single batch.
		MaxMesasges int `json:"max_messages"`
		// MaxSize is the maximum batch size in kilobyte that will be sent.
		MaxSize int `json:"max_size"`
	} `json:"batch"`

	// ServiceRecords is the substructure for configuration related to service record
	// load balancing, either using the attribute `weight` or `priortiy` to load balance.
	ServiceRecords struct {
		// FQDNs is a list of FQDNs pointing to service records
		FQDNs []string `json:"fqdns"`
		// LoadBalancingScheme is the scheme to use for load balancing. Either weight or priortiy
		// can be used.
		LoadBalancingScheme string `json:"load_balancing_scheme"`
	} `json:"service_records"`

	// TODO: figure out JWT attribute for QueueDepth.
	// QueueDepth is the maximum number of events that be queued to be sent as a batch of singal or multiple events.
	// QueueDepth will be validated against customer's JWT attribute `PLACEHOLDER`.
	QueueDepth int `json:"queue_depth"`

	// TODO: figure out JWT attribute for MaxOpenRequests.
	// MaxOpenRequests is the maximum number of outstanding batched event requests are allowed before
	// blocked and additional events are queued.
	// MaxOpenRequests will be validated against customer's JWT attribute `PLACEHOLDER`.
	MaxOpenRequests int `json:"max_open_requests"`
}

// KafkaConfig is a Kafka substructure with data related to event delivery.
type KafkaConfig struct {
	// ID is the unique identifier used for hashing and internal identification.
	ID string `json:"id"`

	// Accept is content type value to set WRP messages to (unless already specified in the WRP).
	Accept string `json:"accept"`

	// BootstrapServers is a list of kafka broker addresses.
	BootstrapServers []string `json:"bootstrap_servers"`

	// TODO: figure out which kafka configuration substructures we want to expose to users (to be set by users)
	// going to be based on https://pkg.go.dev/github.com/IBM/sarama#Config
	// this substructures also includes auth related secrets, noted `MaxOpenRequests` will be excluded since it's already exposed
	KafkaProducerConfig struct{}

	// TODO: figure out JWT attribute for QueueDepth.
	// QueueDepth is the maximum number of events that be queued to be sent as a batch of singal or multiple events.
	// QueueDepth will be validated against customer's JWT attribute `PLACEHOLDER`.
	QueueDepth int `json:"queue_depth"`

	// TODO: figure out JWT attribute for MaxOpenRequests.
	// MaxOpenRequests is the maximum number of outstanding batched event requests are allowed before
	// blocked and additional events are queued.
	// MaxOpenRequests will be validated against customer's JWT attribute `PLACEHOLDER`.
	MaxOpenRequests int `json:"max_open_requests"`
}

// MetadataMatcherConfig is Webhook substructure with config to match event metadata.
type MetadataMatcherConfig struct {
	// DeviceID is the list of regular expressions to match device id type against.
	DeviceID []string `json:"device_id"`
	// Accounts is the list of regular expressions to match account type against.
	Accounts []string `json:"metadata:/account"`
}

// Registration is a special struct for unmarshaling a webhook as part of
// a webhook registration request.  The only difference between this struct and
// the Webhook struct is the Duration field.
type Registration struct {
	// Address is the subscription request origin HTTP Address.
	Address string `json:"registered_from_address"`

	// Deprecated: This field should only be used for backwards compatibility
	// matching. Use ConfigWebhooks instead.
	// Config contains data to inform how events are delivered to single url.
	Config WebhookConfig `json:"config"`

	// ConfigWebhooks contains data to inform how events are delivered to multiple urls.
	ConfigWebhooks []WebhookConfig `json:"config_webhooks"`

	// ConfigWebhooks contains data to inform how events are delivered to multiple urls.
	ConfigKafkas []KafkaConfig `json:"config_kafkas"`

	// Hash is a substructure for configuration related to distributing events among sinks (kafka and webhooks)
	Hash struct {
		// Field is the wrp field to be used for hashing.
		// Either "device_id" or "account" can be used
		Field string `json:"field"`

		// FieldRegex is the regular expression to match `Field` type against.
		FieldRegex string `json:"field_regex"`

		// IDs of `WebhookConfig` or `KafkaConfig` configurations to be map to.
		// (Optional, if omited all provided `WebhookConfig` and `KafkaConfig` configurations will be used)
		IDs []string
	}

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
