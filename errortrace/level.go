// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package errortrace

import (
	"fmt"
	"sort"
	"strings"
)

type LevelType int

const (
	UnknownLevel LevelType = iota
	InfoLevel
	WarningLevel
	ErrorLevel
	lastLevel
)

var (
	validatorLevelUnmarshal = map[string]LevelType{
		"unknown": UnknownLevel,
		"info":    InfoLevel,
		"warning": WarningLevel,
		"error":   ErrorLevel,
	}
	validatorLevelMarshal = map[LevelType]string{
		UnknownLevel: "unknown",
		InfoLevel:    "info",
		WarningLevel: "warning",
		ErrorLevel:   "error",
	}
)

// Empty returns true if the value is UnknownLevel (the default),
// otherwise false is returned.
func (ot LevelType) IsEmpty() bool {
	return UnknownLevel == ot
}

func (ot LevelType) IsValid() bool {
	return UnknownLevel < ot && ot < lastLevel
}

// String returns a human-readable string representation for an existing LevelType,
// otherwise String returns the unknown string value.
func (ot LevelType) String() string {
	if value, ok := validatorLevelMarshal[ot]; ok {
		return value
	}

	return validatorLevelMarshal[UnknownLevel]
}

// UnmarshalText unmarshals a LevelType's enum value.
func (ot *LevelType) UnmarshalText(b []byte) error {
	s := strings.ToLower(string(b))
	r, ok := validatorLevelUnmarshal[s]
	if !ok {
		return fmt.Errorf("ValidatorLevel error: '%s' does not match any valid options: %s",
			s, ot.getKeys())
	}

	*ot = r
	return nil
}

// getKeys returns the string keys for the LevelType enums.
func (ot LevelType) getKeys() string {
	keys := make([]string, 0, len(validatorLevelUnmarshal))
	for k := range validatorLevelUnmarshal {
		k = "'" + k + "'"
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return strings.Join(keys, ", ")
}
