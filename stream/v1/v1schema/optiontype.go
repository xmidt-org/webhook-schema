// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package v1schema

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type OptionType int

// reasons. Use v2 instead.
const (
	UnknownType OptionType = iota
	AddressDefaultType
	MatcherType
	SetSchemaType
	AlwaysValidType
	NotEmptyValidatorType
	AtleastOneEventValidatorType
	EventRegexValidatorType
	DeviceIDRegexValidatorType
	DurationValidatorType
	CheckUntilValidatorType
	ReceiverURLValidatorType
	FailureURLValidatorType
	AlternativeURLValidatorType
	UntilValidatorType
	lastType
)

var ErrOptionTypeInvalid = errors.New("schema option type is invalid")

var (
	OptionTypeUnmarshal = map[string]OptionType{
		"unknown":           UnknownType,
		"address_default":   AddressDefaultType,
		"matcher_default":   MatcherType,
		"set_schema":        SetSchemaType,
		"always_valid":      AlwaysValidType,
		"not_empty":         NotEmptyValidatorType,
		"atleast_one_event": AtleastOneEventValidatorType,
		"event_regex":       EventRegexValidatorType,
		"device_id_regex":   DeviceIDRegexValidatorType,
		"duration":          DurationValidatorType,
		"check_until":       CheckUntilValidatorType,
		"receiver_url":      ReceiverURLValidatorType,
		"failure_url":       FailureURLValidatorType,
		"alt_url":           AlternativeURLValidatorType,
		"until":             UntilValidatorType,
	}
	OptionTypeMarshal = map[OptionType]string{
		UnknownType:                  "unknown",
		AddressDefaultType:           "address_default",
		MatcherType:                  "matcher_default",
		SetSchemaType:                "set_schema",
		AlwaysValidType:              "always_valid",
		NotEmptyValidatorType:        "not_empty",
		AtleastOneEventValidatorType: "atleast_one_event",
		EventRegexValidatorType:      "event_regex",
		DeviceIDRegexValidatorType:   "device_id_regex",
		DurationValidatorType:        "duration",
		CheckUntilValidatorType:      "check_until",
		ReceiverURLValidatorType:     "receiver_url",
		FailureURLValidatorType:      "failure_url",
		AlternativeURLValidatorType:  "alt_url",
		UntilValidatorType:           "until",
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
