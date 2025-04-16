// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package v1schema

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

func TestAtleastOneEventValidatorOption(t *testing.T) {
	tests := []optionTest{
		{
			description: "there is an event - Schema",
			val:         AtleastOneEventValidator(),
			schema:      Schema{Events: []string{"foo"}},
		}, {
			description: "multiple events - Schema",
			val:         AtleastOneEventValidator(),
			schema:      Schema{Events: []string{"foo", "bar"}},
		}, {
			description: "there are no events - Schema",
			val:         AtleastOneEventValidator(),
			schema:      Schema{},
			expectedErr: ErrEmptyEvents,
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
			schema:      Schema{Events: []string{"event.*"}},
		}, {
			description: "multiple events",
			val:         EventRegexValidator(),
			schema:      Schema{Events: []string{"magic-thing", "event.*"}},
		}, {
			description: "failure - Schema",
			val:         EventRegexValidator(),
			schema:      Schema{Events: []string{"("}},
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

func TestDeviceIDRegexValidator(t *testing.T) {
	tests := []optionTest{
		{
			description: "the regex compiles - v1",
			val:         DeviceIDRegexValidator(),
			schema:      Schema{Matcher: MetadataMatcherConfig{DeviceID: []string{"device.*"}}},
		}, {
			description: "multiple device ids - v1",
			val:         DeviceIDRegexValidator(),
			schema:      Schema{Matcher: MetadataMatcherConfig{DeviceID: []string{"device.*", "magic-thing"}}},
		}, {
			description: "failure - v1",
			val:         DeviceIDRegexValidator(),
			schema:      Schema{Matcher: MetadataMatcherConfig{DeviceID: []string{"("}}},
			expectedErr: ErrDeviceIDRegexCompilerFailure,
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
			description: "success with time schema bounds - Schema",
			val:         DurationValidator(time.Now, 5*time.Minute),
			schema:      Schema{Duration: CustomDuration(4 * time.Minute)},
		}, {
			description: "success with time schema bounds, exactly - Schema",
			val:         DurationValidator(time.Now, 5*time.Minute),
			schema:      Schema{Duration: CustomDuration(5 * time.Minute)},
		}, {
			description: "failure with time out of bounds - Schema",
			val:         DurationValidator(time.Now, 5*time.Minute),
			schema:      Schema{Duration: CustomDuration(6 * time.Minute)},
			expectedErr: ErrInvalidDuration,
		}, {
			description: "success with max ttl ignored - Schema",
			val:         DurationValidator(time.Now, -5*time.Minute),
			schema:      Schema{Duration: CustomDuration(1 * time.Minute)},
		}, {
			description: "success with max ttl ignored, 0 duration - Schema",
			val:         DurationValidator(time.Now, 0),
			schema:      Schema{Duration: CustomDuration(1 * time.Minute)},
		}, {
			description: "success with until schema bounds - Schema",
			val:         DurationValidator(now, 5*time.Minute),
			schema:      Schema{Until: time.Date(2021, 1, 1, 0, 4, 0, 0, time.UTC)},
		}, {
			description: "success with until exactly schema bounds - Schema",
			val:         DurationValidator(now, 5*time.Minute),
			schema:      Schema{Until: time.Date(2021, 1, 1, 0, 5, 0, 0, time.UTC)},
		}, {
			description: "failure with until out of bounds - Schema",
			val:         DurationValidator(now, 5*time.Minute),
			schema:      Schema{Until: time.Date(2021, 1, 1, 0, 6, 0, 0, time.UTC)},
			expectedErr: ErrInvalidDuration,
		}, {
			description: "success with until just needing to be present - Schema",
			val:         DurationValidator(now, 0),
			schema:      Schema{Until: time.Date(2021, 1, 1, 0, 6, 0, 0, time.UTC)},
		}, {
			description: "failure, both expirations set - Schema",
			val:         DurationValidator(time.Now, 5*time.Minute),
			schema:      Schema{Duration: CustomDuration(1 * time.Minute), Until: time.Date(2021, 1, 1, 0, 4, 0, 0, time.UTC)},
			expectedErr: ErrInvalidDuration,
		}, {
			description: "failure, no expiration set - Schema",
			schema:      Schema{},
			val:         DurationValidator(time.Now, 5*time.Minute),
			expectedErr: ErrInvalidDuration,
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

func TestReceiverURLValidator(t *testing.T) {
	checker, err := urlegit.New(urlegit.OnlyAllowSchemes("https"))
	require.NoError(t, err)
	require.NotNil(t, checker)
	// nolint:staticcheck
	tests := []optionTest{
		{
			description: "success, with checker - Schema",
			val:         ReceiverURLValidator(checker),
			schema:      Schema{Webhook: webhook.V1Schema{ReceiverURL: "https://example.com"}},
		}, {
			description: "failure, with checker - Schema",
			val:         ReceiverURLValidator(checker),
			schema:      Schema{Webhook: webhook.V1Schema{ReceiverURL: "http://example.com"}},
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

func TestAlternativeURLValidator(t *testing.T) {
	checker, err := urlegit.New(urlegit.OnlyAllowSchemes("https"))
	require.NoError(t, err)
	require.NotNil(t, checker)
	// nolint:staticcheck
	tests := []optionTest{
		{
			description: "success, with checker",
			val:         AlternativeURLValidator(checker),
			schema:      Schema{Webhook: webhook.V1Schema{AlternativeURLs: []string{"https://example.com"}}},
		}, {
			description: "success, with checker and multiple urls",
			val:         AlternativeURLValidator(checker),
			schema:      Schema{Webhook: webhook.V1Schema{AlternativeURLs: []string{"https://example.com", "https://example.org"}}},
		}, {
			description: "failure, with checker",
			val:         AlternativeURLValidator(checker),
			schema:      Schema{Webhook: webhook.V1Schema{AlternativeURLs: []string{"http://example.com"}}},
			expectedErr: ErrAlternativeURLs,
		}, {
			description: "failure, with checker with multiple urls",
			val:         AlternativeURLValidator(checker),
			schema:      Schema{Webhook: webhook.V1Schema{AlternativeURLs: []string{"https://example.com", "http://example.com"}}},
			expectedErr: ErrAlternativeURLs,
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

func TestUntilValidator(t *testing.T) {
	mockNow := func() time.Time {
		return time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	}
	tests := []optionTest{
		{
			description: "success, until",
			schema: Schema{
				Until: mockNow(),
			},
			val: UntilValidator(time.Duration(1*time.Minute), time.Duration(5*time.Minute), time.Now),
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
