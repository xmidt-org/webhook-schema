// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package v2schema

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmidt-org/urlegit"
	"github.com/xmidt-org/webhook-schema/stream/sink/webhook"
)

type optionTest struct {
	description string
	schema      Schema
	val         Option
	expectedErr error
}

func TestAlwaysOption(t *testing.T) {
	tests := []optionTest{
		{
			description: "success",
			val:         AlwaysValid(),
			schema:      Schema{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			assert := assert.New(t)
			err := tc.val.Apply(&tc.schema)
			if tc.expectedErr == nil {
				assert.NoError(err)
			} else {
				assert.ErrorIs(err, tc.expectedErr)
			}
		})
	}
}

func TestReceiverURLValidator(t *testing.T) {
	checker, err := urlegit.New(urlegit.OnlyAllowSchemes("https"))
	require.NoError(t, err)
	require.NotNil(t, checker)
	tests := []optionTest{
		{
			description: "success, with checker - Schema",
			val:         ReceiverURLValidator(checker),
			schema:      Schema{Webhooks: []webhook.V2Schema{{ReceiverURLs: []string{"https://example.com", "https://example.com"}}}},
		}, {
			description: "failure, with checker - Schema",
			val:         ReceiverURLValidator(checker),
			schema:      Schema{Webhooks: []webhook.V2Schema{{ReceiverURLs: []string{"https://example.com", "http://example.com"}}}},
			expectedErr: ErrInvalidReceiverURL,
		},
	}
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			assert := assert.New(t)
			err := tc.val.Apply(&tc.schema)
			if tc.expectedErr == nil {
				assert.NoError(err)
			} else {
				assert.ErrorIs(err, tc.expectedErr)
			}
		})
	}
}

func TestEventRegexValidator(t *testing.T) {
	tests := []optionTest{
		{
			description: "the regex compiles - Schema",
			val:         EventRegexValidator(),
			schema:      Schema{Matcher: []FieldRegex{{Field: "canonical_name", Regex: "webpa"}}},
		},
		{
			description: "multiple matchers - Schema",
			val:         EventRegexValidator(),
			schema:      Schema{Matcher: []FieldRegex{{Field: "canonical_name", Regex: "webpa"}, {Field: "address", Regex: "www.example.com"}}},
		},
		{
			description: "failure - Schema",
			val:         EventRegexValidator(),
			schema:      Schema{Matcher: []FieldRegex{{Regex: "("}}},
			expectedErr: ErrEventRegexCompilerFailure,
		},
	}
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			assert := assert.New(t)
			err := tc.val.Apply(&tc.schema)
			if tc.expectedErr == nil {
				assert.NoError(err)
			} else {
				assert.ErrorIs(err, tc.expectedErr)
			}
		})
	}
}

func TestDurationValidator(t *testing.T) {
	now := func() time.Time {
		return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	tests := []optionTest{
		{
			description: "failure, exipred - Schema",
			schema:      Schema{Expires: now()},
			val:         ExpiresValidator(),
			expectedErr: ErrExpired,
		},
	}
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			assert := assert.New(t)
			err := tc.val.Apply(&tc.schema)
			if tc.expectedErr == nil {
				assert.NoError(err)
			} else {
				assert.ErrorIs(err, tc.expectedErr)
			}
		})
	}
}

func TestFailureURLValidator(t *testing.T) {
	checker, err := urlegit.New(urlegit.OnlyAllowSchemes("https"))
	require.NoError(t, err)
	require.NotNil(t, checker)
	tests := []optionTest{
		{
			description: "success, with checker - Schema",
			val:         FailureURLValidator(checker),
			schema:      Schema{FailureURL: "https://example.com"},
		}, {
			description: "failure, with checker - Schema",
			val:         FailureURLValidator(checker),
			schema:      Schema{FailureURL: "http://example.com"},
			expectedErr: ErrInvalidFailureURL,
		},
	}
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			assert := assert.New(t)
			err := tc.val.Apply(&tc.schema)
			if tc.expectedErr == nil {
				assert.NoError(err)
			} else {
				assert.ErrorIs(err, tc.expectedErr)
			}
		})
	}
}
