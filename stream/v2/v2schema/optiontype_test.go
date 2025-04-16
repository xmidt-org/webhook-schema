// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package v2schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeUnmarshalling(t *testing.T) {
	tests := []struct {
		description string
		config      []byte
		invalid     bool
	}{
		{
			description: "UnknownType valid",
			config:      []byte("unknown"),
		},
		{
			description: "AlwaysValidType valid",
			config:      []byte("always_valid"),
		},
		{
			description: "NotEmptyValidatorType valid",
			config:      []byte("not_empty"),
		},
		{
			description: "OnlyWebhooksValidatorType valid",
			config:      []byte("only_webhooks"),
		},
		{
			description: "EventRegexType valid",
			config:      []byte("event_regex"),
		},
		{
			description: "ExpiresType valid",
			config:      []byte("expires"),
		},
		{
			description: "ReceiverURLType valid",
			config:      []byte("receiver_url"),
		},
		{
			description: "FailureURLType valid",
			config:      []byte("failure_url"),
		},
		{
			description: "Nonexistent type invalid",
			config:      []byte("FOOBAR"),
			invalid:     true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			assert := assert.New(t)
			var l OptionType

			err := l.UnmarshalText(tc.config)
			assert.NotEmpty(l.getKeys())
			if tc.invalid {
				assert.Error(err)
			} else {
				assert.NoError(err)
				assert.Equal(string(tc.config), l.String())
			}
		})
	}
}

func TestTypeState(t *testing.T) {
	tests := []struct {
		description string
		val         OptionType
		expectedVal string
		invalid     bool
		empty       bool
	}{
		{
			description: "UnknownType valid",
			val:         UnknownType,
			expectedVal: "unknown",
			empty:       true,
			invalid:     true,
		},
		{
			description: "AlwaysValidType valid",
			val:         AlwaysValidType,
			expectedVal: "always_valid",
		},
		{
			description: "NotEmptyValidatorType valid",
			val:         NotEmptyValidatorType,
			expectedVal: "not_empty",
		},
		{
			description: "OnlyWebhooksValidatorType valid",
			val:         OnlyWebhooksValidatorType,
			expectedVal: "only_webhooks",
		},
		{
			description: "EventRegexType valid",
			val:         EventRegexValidatorType,
			expectedVal: "event_regex",
		},
		{
			description: "ExpiresType valid",
			val:         ExpiresValidatorType,
			expectedVal: "expires",
		},
		{
			description: "ReceiverURLType valid",
			val:         ReceiverURLValidatorType,
			expectedVal: "receiver_url",
		},
		{
			description: "FailureURLType valid",
			val:         FailureURLValidatorType,
			expectedVal: "failure_url",
		},
		{
			description: "lastLevel valid",
			val:         lastType,
			expectedVal: "unknown",
			invalid:     true,
		},
		{
			description: "Nonexistent positive Level invalid",
			val:         lastType + 1,
			expectedVal: "unknown",
			invalid:     true,
		},
		{
			description: "Nonexistent negative Level invalid",
			val:         UnknownType - 1,
			expectedVal: "unknown",
			invalid:     true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			assert := assert.New(t)
			assert.Equal(tc.expectedVal, tc.val.String())
			assert.Equal(!tc.invalid, tc.val.IsValid())
			assert.Equal(tc.empty, tc.val.IsEmpty())
		})
	}
}
