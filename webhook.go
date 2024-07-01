// SPDX-FileCopyrightText: 2022 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package webhook

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/xmidt-org/urlegit"
)

var (
	ErrInvalidInput = fmt.Errorf("invalid input")
)

type Register interface {
	GetId() string
	GetUntil() time.Time
}

// Deprecated: This substructure should only be used for backwards compatibility
// matching. Use Webhook instead.
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

// Deprecated: This structure should only be used for backwards compatibility
// matching. Use RegistrationV2 instead.
// RegistrationV1 is a special struct for unmarshaling a webhook as part of a webhook registration request.
type RegistrationV1 struct {
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
}

type RetryHint struct {
	//RetryEachUrl is the amount of times a URL should be retried given a failed response until the next URL in the request is tried.
	//Default value will be set to none
	RetryEachUrl int `json:"retry_each_url"`

	//MaxRetry is the total amount times a request will be retried.
	MaxRetry int `json:"max_retry"`
}

// Webhook is a substructure with data related to event delivery.
type Webhook struct {
	// Accept is the encoding type of outgoing events. The following encoding types are supported, otherwise
	// a 406 response code is returned: application/octet-stream, application/json, application/jsonl, application/msgpack.
	// Note: An `Accept` of application/octet-stream or application/json will result in a single response for batch sizes of 0 or 1
	// and batch sizes greater than 1 will result in a multipart response. An `Accept` of application/jsonl or application/msgpack
	// will always result in a single response with a list of batched events for any batch size.
	Accept string `json:"accept"`

	// AcceptEncoding is the content type of outgoing events. The following content types are supported, otherwise
	// a 406 response code is returned: gzip.
	AcceptEncoding string `json:"accept_encoding"`

	// Secret is the string value.
	// (Optional, set to "" to disable behavior).
	Secret string `json:"secret,omitempty"`

	// SecretHash is the hash algorithm to be used. Only sha256 HMAC and sha512 HMAC are supported.
	// (Optional).
	// The Default value is the largest sha HMAC supported, sha512 HMAC.
	SecretHash string `json:"secret_hash"`

	// If true, response will use the device content-type and wrp payload as its body
	// Otherwise, response will Accecpt as the content-type and wrp message as its body
	// Default: False (the entire wrp message is sent)
	PayloadOnly bool `json:"payload_only"`

	// ReceiverUrls is the list of receiver urls that will be used where as if the first url fails,
	// then the second url would be used and so on.
	// Note: either `ReceiverURLs` or `DNSSrvRecord` must be used but not both.
	ReceiverURLs []string `json:"receiver_urls"`

	// DNSSrvRecord is the substructure for configuration related to load balancing.
	// Note: either `ReceiverURLs` or `DNSSrvRecord` must be used but not both.
	DNSSrvRecord struct {
		// FQDNs is a list of FQDNs pointing to dns srv records
		FQDNs []string `json:"fqdns"`

		// LoadBalancingScheme is the scheme to use for load balancing. Either the
		// srv record attribute `weight` or `priortiy` can be used.
		LoadBalancingScheme string `json:"load_balancing_scheme"`
	} `json:"dns_srv_record"`

	//RetryHint is the substructure for configuration related to retrying requests.
	// (Optional, if omited then retries will be based on default values defined by server)
	RetryHint RetryHint `json:"retry_hint"`
}

// Kafka is a substructure with data related to event delivery.
type Kafka struct {
	// Accept is content type value to set WRP messages to (unless already specified in the WRP).
	Accept string `json:"accept"`

	// BootstrapServers is a list of kafka broker addresses.
	BootstrapServers []string `json:"bootstrap_servers"`

	// TODO: figure out which kafka configuration substructures we want to expose to users (to be set by users)
	// going to be based on https://pkg.go.dev/github.com/IBM/sarama#Config
	// this substructures also includes auth related secrets, noted `MaxOpenRequests` will be excluded since it's already exposed
	KafkaProducer struct{} `json:"kafka_producer"`

	//RetryHint is the substructure for configuration related to retrying requests.
	// (Optional, if omited then retries will be based on default values defined by server)
	RetryHint RetryHint `json:"retry_hint"`
}

// FieldRegex is a substructure with data related to regular expressions.
type FieldRegex struct {
	// Field is the wrp field to be used for regex.
	// All wrp field can be used, refer to the schema for examples.
	Field string `json:"field"`

	// FieldRegex is the regular expression to match `Field` against to.
	Regex string `json:"regex"`
}

type BatchHint struct {
	// MaxLingerDuration is the maximum delay for batching if MaxMesasges has not been reached.
	// Default value will set no maximum value.
	MaxLingerDuration time.Duration `json:"max_linger_duration"`
	// MaxMesasges is the maximum number of events that will be sent in a single batch.
	// Default value will set no maximum value.
	MaxMesasges int `json:"max_messages"`
}

type ContactInfo struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

// RegistrationV2 is a special struct for unmarshaling sink information as part of a sink registration request.
type RegistrationV2 struct {
	// ContactInfo contains contact information used to reach the owner of the registration.
	// (Optional).
	ContactInfo ContactInfo `json:"contact_info,omitempty"`

	// CanonicalName is the canonical name of the registration request.
	// Reusing a CanonicalName will override the configurations set in that previous
	// registration request with the same CanonicalName.
	CanonicalName string `json:"canonical_name"`

	// Address is the subscription request origin HTTP Address.
	Address string `json:"registered_from_address"`

	// Webhooks contains data to inform how events are delivered to multiple urls.
	Webhooks []Webhook `json:"webhooks"`

	// Kafkas contains data to inform how events are delivered to multiple kafkas.
	Kafkas []Kafka `json:"kafkas"`

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
	// Note. Any failures due to a bad regex feild or regex expression will result in a silent failure.
	Matcher []FieldRegex `json:"matcher,omitempty"`

	// Expires describes the time this subscription expires.
	// TODO: list of supported formats
	Expires time.Time `json:"expires"`
}

type Option interface {
	fmt.Stringer
	Validate(Validator) error
}

// Validate is a method on Registration that validates the registration
// against a list of options.
func Validate(v Validator, opts []Option) error {
	var errs error
	for _, opt := range opts {
		if opt != nil {
			if err := opt.Validate(v); err != nil {
				errs = errors.Join(errs, err)
			}
		}
	}
	return errs
}

func (v1 *RegistrationV1) ValidateOneEvent() error {
	if len(v1.Events) == 0 {
		return fmt.Errorf("%w: cannot have zero events", ErrInvalidInput)
	}
	return nil
}

func (v1 *RegistrationV1) ValidateEventRegex() error {
	var errs error
	for _, e := range v1.Events {
		_, err := regexp.Compile(e)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("%w: unable to compile matching", ErrInvalidInput))
		}
	}
	return errs
}

func (v1 *RegistrationV1) ValidateDeviceId() error {
	var errs error
	for _, e := range v1.Matcher.DeviceID {
		_, err := regexp.Compile(e)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("%w: unable to compile matching", ErrInvalidInput))
		}
	}
	return errs
}

func (v1 *RegistrationV1) ValidateDuration(ttl time.Duration) error {
	var errs error
	if ttl <= 0 {
		ttl = time.Duration(0)
	}

	if ttl != 0 && ttl < time.Duration(v1.Duration) {
		errs = errors.Join(errs, fmt.Errorf("%w: the registration is for too long", ErrInvalidInput))
	}

	if v1.Until.IsZero() && v1.Duration == 0 {
		errs = errors.Join(errs, fmt.Errorf("%w: either Duration or Until must be set", ErrInvalidInput))
	}

	if !v1.Until.IsZero() && v1.Duration != 0 {
		errs = errors.Join(errs, fmt.Errorf("%w: only one of Duration or Until may be set", ErrInvalidInput))
	}

	if !v1.Until.IsZero() {
		nowFunc := time.Now
		// if v1.nowFunc != nil {
		// 	nowFunc = v1.nowFunc
		// }

		now := nowFunc()
		if ttl != 0 && v1.Until.After(now.Add(ttl)) {
			errs = errors.Join(errs, fmt.Errorf("%w: the registration is for too long", ErrInvalidInput))
		}

		if v1.Until.Before(now) {
			errs = errors.Join(errs, fmt.Errorf("%w: the registration has already expired", ErrInvalidInput))
		}
	}

	return errs
}

func (v1 *RegistrationV1) ValidateFailureURL(c *urlegit.Checker) error {
	if v1.FailureURL != "" {
		if err := c.Text(v1.FailureURL); err != nil {
			return fmt.Errorf("%w: failure url is invalid", ErrInvalidInput)
		}
	}
	return nil
}

func (v1 *RegistrationV1) ValidateReceiverURL(c *urlegit.Checker) error {
	if v1.Config.ReceiverURL != "" {
		if err := c.Text(v1.Config.ReceiverURL); err != nil {
			return fmt.Errorf("%w: failure url is invalid", ErrInvalidInput)
		}
	}
	return nil
}

func (v1 *RegistrationV1) ValidateAltURL(c *urlegit.Checker) error {
	var errs error
	for _, url := range v1.Config.AlternativeURLs {
		if err := c.Text(url); err != nil {
			errs = errors.Join(errs, fmt.Errorf("%w: failure url is invalid", ErrInvalidInput))
		}
	}
	return errs
}

func (v1 *RegistrationV1) ValidateNoUntil() error {
	if !v1.Until.IsZero() {
		return fmt.Errorf("%w: Until is not allowed", ErrInvalidInput)
	}
	return nil
}

func (v1 *RegistrationV1) ValidateUntil(jitter time.Duration, maxTTL time.Duration, now func() time.Time) error {
	if now == nil {
		now = time.Now
	}
	if maxTTL < 0 {
		return ErrInvalidInput
	} else if jitter < 0 {
		return ErrInvalidInput
	}

	if v1.Until.IsZero() {
		return nil
	}
	limit := (now().Add(maxTTL)).Add(jitter)
	proposed := (v1.Until)
	if proposed.After(limit) {
		return fmt.Errorf("%w: %v after %v",
			ErrInvalidInput, proposed.String(), limit.String())
	}
	return nil

}
