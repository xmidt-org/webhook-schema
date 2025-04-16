// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package registry

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/multierr"
)

var (
	ErrMarshal   = fmt.Errorf("`%T` marshal error", Record{})
	ErrUnmarshal = fmt.Errorf("`%T` unmarshal error", Record{})
)

// Registrable defines the functionality required to handle and process json objects containing a stream configuration, among other metadata,
// in and out and within a wrpeventstream registration service, while being stream version and sink agnostic.
type Registrable interface {
	// Validate validates the stream.
	Validate() error
	// GetID returns a unique identifier used for storage and retrieval of stream.
	GetID() string
	// GetTTLSeconds returns the TTL of the stream.
	GetTTLSeconds(func() time.Time) int64
}

func New(stream Registrable, opts ...Option) (Record, error) {
	if err := stream.Validate(); err != nil {
		return Record{}, multierr.Append(ErrMarshal, err)
	}

	TTLSeconds := stream.GetTTLSeconds(time.Now)
	jsonStream, err := json.Marshal(stream)
	if err != nil {
		return Record{}, multierr.Append(ErrMarshal, err)
	}

	var data map[string]any
	if err = json.Unmarshal(jsonStream, &data); err != nil {
		err = multierr.Append(fmt.Errorf("failed to marshal `%T`", data), err)
		return Record{}, multierr.Append(ErrMarshal, err)
	}

	r := Record{
		Data: data,
		ID:   fmt.Sprintf("%x", sha256.Sum256([]byte(stream.GetID()))),
		TTL:  &TTLSeconds,
	}

	defaultsVal := Options{
		IDValidator(),
		DataValidator(),
	}
	opts = append(opts, defaultsVal...)
	if err := Options(opts).Apply(&r); err != nil {
		return Record{}, multierr.Append(ErrMarshal, err)
	}

	return r, nil
}

// Key defines the field mapping to retrieve an item from storage.
type Key struct {
	// Bucket is the name for a collection or partition to which an item belongs.
	Bucket string `json:"bucket"`

	// ID is an item's identifier. Note that different buckets may have
	// items with the same ID.
	ID string `json:"id"`
}

// Record is the registry record used by argus, the wrpeventstream registry component, to store registered wrpeventstream configurations.
type Record struct {
	// ID is the unique ID identifying this item. It is recommended this value is the resulting
	// value of a SHA256 calculation, using the unique attributes of the object being represented
	// (e.g. SHA256(<common_name>)). This will be used by argus to determine uniqueness of objects being stored or updated.
	ID string `json:"id"`

	// Data is the JSON object to be stored. Opaque to argus.
	Data map[string]any `json:"data"`

	// TTL is the time to live in storage, specified in seconds.
	// Optional. When not set, items don't expire.
	TTL *int64 `json:"ttl,omitempty"`
}
