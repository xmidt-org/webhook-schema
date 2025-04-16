// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package v2schema

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type OptionType int

const (
	UnknownType OptionType = iota
	// NotEmptyValidator is a default validator, meaning it's used
	// regardless of the configuration.
	NotEmptyValidatorType
	// Validators are exposed to configuration.
	AlwaysValidType
	OnlyWebhooksValidatorType
	EventRegexValidatorType
	ExpiresValidatorType
	ReceiverURLValidatorType
	FailureURLValidatorType
	// Defaults/setters are not exposed to configuration, but they can be used
	// programmatically.
	AddressDefaultType
	SetSchemaType
	lastType
)

var ErrOptionTypeInvalid = errors.New("schema option type is invalid")

var (
	OptionTypeUnmarshal = map[string]OptionType{
		"unknown":       UnknownType,
		"always_valid":  AlwaysValidType,
		"not_empty":     NotEmptyValidatorType,
		"only_webhooks": OnlyWebhooksValidatorType,
		"event_regex":   EventRegexValidatorType,
		"expires":       ExpiresValidatorType,
		"receiver_url":  ReceiverURLValidatorType,
		"failure_url":   FailureURLValidatorType,
	}
	OptionTypeMarshal = map[OptionType]string{
		UnknownType:               "unknown",
		AlwaysValidType:           "always_valid",
		NotEmptyValidatorType:     "not_empty",
		OnlyWebhooksValidatorType: "only_webhooks",
		EventRegexValidatorType:   "event_regex",
		ExpiresValidatorType:      "expires",
		ReceiverURLValidatorType:  "receiver_url",
		FailureURLValidatorType:   "failure_url",
	}
)

// IsEmpty returns true if the value is UnknownType (the default),
// otherwise false is returned.
func (ot OptionType) IsEmpty() bool {
	return ot == UnknownType
}

// IsValid returns true if the value is valid,
// otherwise false is returned.
func (ot OptionType) IsValid() bool {
	return UnknownType < ot && ot < lastType
}

// String returns a human-readable string representation for an existing OptionType,
// otherwise String returns the `unknown` string value.
func (ot OptionType) String() string {
	if value, ok := OptionTypeMarshal[ot]; ok {
		return value
	}

	return OptionTypeMarshal[UnknownType]
}

// UnmarshalText unmarshals a OptionType's enum value.
func (ot *OptionType) UnmarshalText(b []byte) error {
	s := strings.ToLower(string(b))
	OT, ok := OptionTypeUnmarshal[s]
	if !ok {
		return fmt.Errorf("%w: '%s' does not match any valid options: %s", ErrOptionTypeInvalid,
			s, ot.getKeys())
	}

	*ot = OT
	return nil
}

// getKeys returns the string keys for the OptionType enums.
func (ot OptionType) getKeys() string {
	keys := make([]string, 0, len(OptionTypeUnmarshal))
	for k := range OptionTypeUnmarshal {
		k = "'" + k + "'"
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return strings.Join(keys, ", ")
}
