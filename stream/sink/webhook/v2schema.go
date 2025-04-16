// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package webhook

// V2Schema is a substructure with data related to event delivery.
type V2Schema struct {
	// Accept is the encoding type of outgoing events. The following encoding types are supported, otherwise
	// a 406 response code is returned: application/octet-stream, application/json, application/jsonl, application/msgpack.
	/*
		Note:
			application/wrp+json - one json encoded wrp message
			application/wrp+msgpack - one msgpack encoded wrp message
			application/wrp+octet-stream - one message with the wrp payload in the http payload
			application/wrp+jsonl - multiple jsonl encoded wrp messages
			application/wrp+msgpackl - multiple msgpackl encoded wrp messages
	*/
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
	DNSSrvRecord DNSSrvRecordV2 `json:"dns_srv_record"`

	// RetryHint is the substructure for configuration related to retrying requests.
	// (Optional, if omited then retries will be based on default values defined by server)
	RetryHint RetryHintV2 `json:"retry_hint"`
}

// DNSSrvRecordV2 is the substructure for configuration related to load balancing.
type DNSSrvRecordV2 struct {
	// FQDNs is a list of FQDNs pointing to dns srv records
	FQDNs []string `json:"fqdns"`

	// LoadBalancingScheme is the scheme to use for load balancing. Either the
	// srv record attribute `weight` or `priortiy` can be used.
	LoadBalancingScheme string `json:"load_balancing_scheme"`
}

// RetryHintV2 is the substructure for configuration related to retrying requests.
type RetryHintV2 struct {
	//RetryEachUrl is the amount of times a URL should be retried given a failed response until the next URL in the request is tried.
	//Default value will be set to none
	RetryEachUrl int `json:"retry_each_url"`

	//MaxRetry is the total amount times a request will be retried.
	MaxRetry int `json:"max_retry"`
}
