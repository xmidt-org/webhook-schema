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

	// Secret is the string value.
	// (Optional, set to "" to disable behavior).
	Secret string `json:"secret,omitempty"`

	// AlternativeURLs is a list of explicit URLs that should be round robin through on failure cases to the main URL.
	AlternativeURLs []string `json:"alt_urls,omitempty"`

	// ID is the configured webhook's name used to map hashed events to.
	// Refer to the Hash substructure configuration for more details.
	ID string `json:"id"`

	// SecretHash is the hash algorithm to be used. Only sha256 HMAC and sha512 HMAC are supported.
	// (Optional).
	// Deprecated: The Default value is the sha1 HMAC for backwards compatibility.
	SecretHash string `json:"secret_hash"`

	// Batch is the substructure for configuration related to event batching.
	// (Optional, if omited then batches of singal events will be sent)
	// Default value will disable batch. All zeros will also disable batch.
	Batch struct {
		// MaxLingerDuration is the maximum delay for batching if MaxMesasges has not been reached.
		// Default value will set no maximum value.
		MaxLingerDuration time.Duration `json:"max_linger_duration"`
		// MaxMesasges is the maximum number of events that will be sent in a single batch.
		// Default value will set no maximum value.
		MaxMesasges int `json:"max_messages"`
	} `json:"batch"`

	// DNSSrvRecord is the substructure for configuration related to load balancing.
	DNSSrvRecord struct {
		// FQDNs is a list of FQDNs pointing to dns srv records
		FQDNs []string `json:"fqdns"`
		// LoadBalancingScheme is the scheme to use for load balancing. Either the
		// srv record attribute `weight` or `priortiy` can be used.
		LoadBalancingScheme string `json:"load_balancing_scheme"`
	} `json:"dns_srv_record"`

}

// KafkaConfig is a Kafka substructure with data related to event delivery.
type KafkaConfig struct {
	// ID is the configured kafka's name used to map hashed events to.
	// Refer to the Hash substructure configuration for more details.
	ID string `json:"id"`

	// Accept is content type value to set WRP messages to (unless already specified in the WRP).
	Accept string `json:"accept"`

	// BootstrapServers is a list of kafka broker addresses.
	BootstrapServers []string `json:"bootstrap_servers"`

	// TODO: figure out which kafka configuration substructures we want to expose to users (to be set by users)
	// going to be based on https://pkg.go.dev/github.com/IBM/sarama#Config
	// this substructures also includes auth related secrets, noted `MaxOpenRequests` will be excluded since it's already exposed
	KafkaProducer struct{} `json:"kafka_producer"`

}

// MetadataMatcherConfig is Webhook substructure with config to match event metadata.
type MetadataMatcherConfig struct {
	// DeviceID is the list of regular expressions to match device id type against.
	DeviceID []string `json:"device_id"`

	// Account is the list of regular expressions to match account type against.
	Account []string `json:"metadata:/account"`

	// Model is the list of regular expressions to match model type against.
	Model []string `json:"metadata:/hw-model"`

	// FirmwareName is the list of regular expressions to match firmware type against.
	FirmwareName []string `json:"metadata:/fw-name"`
}

// Registration is a special struct for unmarshaling a webhook as part of
// a webhook registration request.  The only difference between this struct and
// the Webhook struct is the Duration field.
type Registration struct {
	// CanonicalName is the canonical name of the registration request.
	// Reusing a CanonicalName will override the configurations set in that previous
	// registration request with the same CanonicalName.
	CanonicalName string `json:"canonical_name"`

	// Address is the subscription request origin HTTP Address.
	Address string `json:"registered_from_address"`

	// Deprecated: This field should only be used for backwards compatibility
	// matching. Use ConfigWebhooks instead.
	// Config contains data to inform how events are delivered to single url.
	Config WebhookConfig `json:"config"`

	// Webhooks contains data to inform how events are delivered to multiple urls.
	Webhooks []WebhookConfig `json:"webhooks"`

	// Kafkas contains data to inform how events are delivered to multiple kafkas.
	Kafkas []KafkaConfig `json:"kafkas"`

	// Hash is a substructure for configuration related to distributing events among sinks (kafka and webhooks)
	Hash struct {
		// Field is the wrp field to be used for hashing.
		// Either "device_id" or "account" can be used
		Field string `json:"field"`

		// FieldRegex is the regular expression to match `Field` type against.
		FieldRegex string `json:"field_regex"`

		// IDs is the list of configured webhooks' and kafkas' names that hashed events to be sent to.
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
