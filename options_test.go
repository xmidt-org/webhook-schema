// SPDX-FileCopyrightText: 2022 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package webhook

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmidt-org/urlegit"
)

type optionTest struct {
	description string
	in          any
	opt         Option
	opts        []Option
	str         string
	expectedErr error
}

func TestErrorOption(t *testing.T) {
	run_tests(t, []optionTest{
		{
			description: "success",
			in:          &RegistrationV1{},
			str:         "foo",
		},
		{
			description: "simple error - RegistrationV1",
			opt:         Error(ErrInvalidInput),
			str:         "Error('invalid input')",
			expectedErr: ErrInvalidInput,
			in:          &RegistrationV1{},
		},
		{
			description: "simple error - RegistrationV2",
			opt:         Error(ErrInvalidInput),
			str:         "Error('invalid input')",
			expectedErr: ErrInvalidInput,
			in:          &RegistrationV2{},
		},
		{
			description: "simple nil error",
			opt:         Error(nil),
			str:         "Error(nil)",
		},
	})
}

func TestAtLeastOneEventOption(t *testing.T) {
	run_tests(t, []optionTest{
		{
			description: "there is an event - V1",
			opt:         AtLeastOneEvent(),
			in:          &RegistrationV1{Events: []string{"foo"}},
			str:         "AtLeastOneEvent()",
		}, {
			description: "multiple events - V1",
			opt:         AtLeastOneEvent(),
			in:          &RegistrationV1{Events: []string{"foo", "bar"}},
			str:         "AtLeastOneEvent()",
		}, {
			description: "there are no events - V1",
			opt:         AtLeastOneEvent(),
			in:          &RegistrationV1{},
			expectedErr: ErrInvalidInput,
		},
		{
			description: "invalid type - RegistrationV2",
			opt:         AtLeastOneEvent(),
			in:          &RegistrationV2{},
			expectedErr: ErrInvalidType,
		},
		{
			description: "default case - invalid",
			opt:         AtLeastOneEvent(),
			expectedErr: ErrInvalidType,
		},
	})
}

func TestEventRegexMustCompile(t *testing.T) {
	run_tests(t, []optionTest{
		{
			description: "the regex compiles - V1",
			opt:         EventRegexMustCompile(),
			in:          &RegistrationV1{Events: []string{"event.*"}},
			str:         "EventRegexMustCompile()",
		}, {
			description: "multiple events",
			opt:         EventRegexMustCompile(),
			in:          &RegistrationV1{Events: []string{"magic-thing", "event.*"}},
			str:         "EventRegexMustCompile()",
		}, {
			description: "failure - V1",
			opt:         EventRegexMustCompile(),
			in:          &RegistrationV1{Events: []string{"("}},
			expectedErr: ErrInvalidInput,
		},
		{
			description: "the regex compiles - V2",
			opt:         EventRegexMustCompile(),
			in: &RegistrationV2{Matcher: []FieldRegex{
				{
					Field: "canonical_name",
					Regex: "webpa",
				},
			}},
			str: "EventRegexMustCompile()",
		},
		{
			description: "multiple matchers - V2",
			opt:         EventRegexMustCompile(),
			in: &RegistrationV2{Matcher: []FieldRegex{
				{
					Field: "canonical_name",
					Regex: "webpa",
				},
				{
					Field: "address",
					Regex: "www.example.com",
				},
			}},
			str: "EventRegexMustCompile()",
		},
		{
			description: "failure - V2",
			opt:         EventRegexMustCompile(),
			in: &RegistrationV2{Matcher: []FieldRegex{
				{
					Regex: "(",
				},
			}},
			expectedErr: ErrInvalidInput,
		},
		{
			description: "default case - invalid",
			opt:         EventRegexMustCompile(),
			expectedErr: ErrInvalidType,
		},
	})
}

func TestDeviceIDRegexMustCompile(t *testing.T) {
	run_tests(t, []optionTest{
		{
			description: "the regex compiles",
			opt:         DeviceIDRegexMustCompile(),
			in: &RegistrationV1{
				Matcher: MetadataMatcherConfig{
					DeviceID: []string{"device.*"},
				},
			},
			str: "DeviceIDRegexMustCompile()",
		}, {
			description: "multiple device ids",
			opt:         DeviceIDRegexMustCompile(),
			in: &RegistrationV1{
				Matcher: MetadataMatcherConfig{
					DeviceID: []string{"device.*", "magic-thing"},
				},
			},
			str: "DeviceIDRegexMustCompile()",
		}, {
			description: "failure",
			opt:         DeviceIDRegexMustCompile(),
			in: &RegistrationV1{
				Matcher: MetadataMatcherConfig{
					DeviceID: []string{"("},
				},
			},
			expectedErr: ErrInvalidInput,
		},
		{
			description: "default case - invalid",
			opt:         DeviceIDRegexMustCompile(),
			expectedErr: ErrInvalidType,
		},
	})
}

func TestValidateRegistrationDuration(t *testing.T) {
	now := func() time.Time {
		return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	run_tests(t, []optionTest{
		{
			description: "success with time in bounds - V1",
			opt:         ValidateRegistrationDuration(5 * time.Minute),
			in: &RegistrationV1{
				Duration: CustomDuration(4 * time.Minute),
			},
			str: "ValidateRegistrationDuration(5m0s)",
		}, {
			description: "success with time in bounds, exactly - V1",
			opt:         ValidateRegistrationDuration(5 * time.Minute),
			in: &RegistrationV1{
				Duration: CustomDuration(5 * time.Minute),
			},
		}, {
			description: "failure with time out of bounds - V1",
			opt:         ValidateRegistrationDuration(5 * time.Minute),
			in: &RegistrationV1{
				Duration: CustomDuration(6 * time.Minute),
			},
			expectedErr: ErrInvalidInput,
		}, {
			description: "success with max ttl ignored - V1",
			opt:         ValidateRegistrationDuration(-5 * time.Minute),
			in: &RegistrationV1{
				Duration: CustomDuration(1 * time.Minute),
			},
		}, {
			description: "success with max ttl ignored, 0 duration - V1",
			opt:         ValidateRegistrationDuration(0),
			in: &RegistrationV1{
				Duration: CustomDuration(1 * time.Minute),
			},
		}, {
			description: "success with until in bounds - V1",
			opts: []Option{
				ProvideTimeNowFunc(now),
				ValidateRegistrationDuration(5 * time.Minute),
			},
			in: &RegistrationV1{
				Until: time.Date(2021, 1, 1, 0, 4, 0, 0, time.UTC),
			},
		}, {
			description: "failure due to until being before now - V1",
			opts: []Option{
				ValidateRegistrationDuration(5 * time.Minute),
				ProvideTimeNowFunc(now),
			},
			in: &RegistrationV1{
				Until: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedErr: ErrInvalidInput,
		}, {
			description: "success with until exactly in bounds - V1",
			opts: []Option{
				ProvideTimeNowFunc(now),
				ValidateRegistrationDuration(5 * time.Minute),
			},
			in: &RegistrationV1{
				Until: time.Date(2021, 1, 1, 0, 5, 0, 0, time.UTC),
			},
		}, {
			description: "failure due to the options being out of order - V1",
			opts: []Option{
				ValidateRegistrationDuration(5 * time.Minute),
				ProvideTimeNowFunc(now),
			},
			in: &RegistrationV1{
				Until: time.Date(2021, 1, 1, 0, 5, 0, 0, time.UTC),
			},
			expectedErr: ErrInvalidInput,
		}, {
			description: "failure with until out of bounds - V1",
			opts: []Option{
				ProvideTimeNowFunc(now),
				ValidateRegistrationDuration(5 * time.Minute),
			},
			in: &RegistrationV1{
				Until: time.Date(2021, 1, 1, 0, 6, 0, 0, time.UTC),
			},
			expectedErr: ErrInvalidInput,
		}, {
			description: "success with until just needing to be present - V1",
			opts: []Option{
				ProvideTimeNowFunc(now),
				ValidateRegistrationDuration(0),
			},
			in: &RegistrationV1{
				Until: time.Date(2021, 1, 1, 0, 6, 0, 0, time.UTC),
			},
		}, {
			description: "failure, both expirations set - V1",
			opt:         ValidateRegistrationDuration(5 * time.Minute),
			in: &RegistrationV1{
				Duration: CustomDuration(1 * time.Minute),
				Until:    time.Date(2021, 1, 1, 0, 4, 0, 0, time.UTC),
			},
			expectedErr: ErrInvalidInput,
		}, {
			description: "failure, no expiration set - V1",
			in:          &RegistrationV1{},
			opt:         ValidateRegistrationDuration(5 * time.Minute),
			expectedErr: ErrInvalidInput,
		}, {
			description: "failure, exipred - V2",
			in: &RegistrationV2{
				Expires: now(),
			},
			opt:         ValidateRegistrationDuration(0),
			expectedErr: ErrInvalidInput,
		},
		{
			description: "default case - invalid",
			opt:         ValidateRegistrationDuration(5 * time.Minute),
			expectedErr: ErrInvalidType,
		},
	})
}

func TestProvideTimeNowFunc(t *testing.T) {
	now := func() time.Time {
		return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	run_tests(t, []optionTest{
		{
			description: "success",
			opt:         ProvideTimeNowFunc(now),
			str:         "ProvideTimeNowFunc(func)",
		}, {
			description: "success as nil",
			opt:         ProvideTimeNowFunc(nil),
			str:         "ProvideTimeNowFunc(nil)",
		},
	})
}

func TestProvideFailureURLValidator(t *testing.T) {
	checker, err := urlegit.New(urlegit.OnlyAllowSchemes("https"))
	require.NoError(t, err)
	require.NotNil(t, checker)

	run_tests(t, []optionTest{
		{
			description: "success, no checker",
			opt:         ProvideFailureURLValidator(nil),
			str:         "ProvideFailureURLValidator(nil)",
		}, {
			description: "success, with checker - V1",
			opt:         ProvideFailureURLValidator(checker),
			in: &RegistrationV1{
				FailureURL: "https://example.com",
			},
			str: "ProvideFailureURLValidator(urlegit.Checker{ OnlyAllowSchemes('https') })",
		}, {
			description: "failure, with checker - V1",
			opt:         ProvideFailureURLValidator(checker),
			in: &RegistrationV1{
				FailureURL: "http://example.com",
			},
			expectedErr: ErrInvalidInput,
		}, {
			description: "success, with checker - V2",
			opt:         ProvideFailureURLValidator(checker),
			in: &RegistrationV2{
				FailureURL: "https://example.com",
			},
			str: "ProvideFailureURLValidator(urlegit.Checker{ OnlyAllowSchemes('https') })",
		}, {
			description: "failure, with checker - V2",
			opt:         ProvideFailureURLValidator(checker),
			in: &RegistrationV2{
				FailureURL: "http://example.com",
			},
			expectedErr: ErrInvalidInput,
		}, {
			description: "default case - invalid",
			opt:         ProvideFailureURLValidator(checker),
			expectedErr: ErrInvalidType,
		},
	})
}

func TestProvideReceiverURLValidator(t *testing.T) {
	checker, err := urlegit.New(urlegit.OnlyAllowSchemes("https"))
	require.NoError(t, err)
	require.NotNil(t, checker)

	run_tests(t, []optionTest{
		{
			description: "success, no checker",
			opt:         ProvideReceiverURLValidator(nil),
			str:         "ProvideReceiverURLValidator(nil)",
		}, {
			description: "success, with checker - V1",
			opt:         ProvideReceiverURLValidator(checker),
			in: &RegistrationV1{
				Config: DeliveryConfig{
					ReceiverURL: "https://example.com",
				},
			},
			str: "ProvideReceiverURLValidator(urlegit.Checker{ OnlyAllowSchemes('https') })",
		}, {
			description: "failure, with checker - V1",
			opt:         ProvideReceiverURLValidator(checker),
			in: &RegistrationV1{
				Config: DeliveryConfig{
					ReceiverURL: "http://example.com",
				},
			},
			expectedErr: ErrInvalidInput,
		}, {
			description: "success, with checker - V2",
			opt:         ProvideReceiverURLValidator(checker),
			in: &RegistrationV2{
				Webhooks: []Webhook{
					{
						ReceiverURLs: []string{"https://example.com",
							"https://example2.com"},
					},
				},
			},
			str: "ProvideReceiverURLValidator(urlegit.Checker{ OnlyAllowSchemes('https') })",
		}, {
			description: "failure, with checker - V2",
			opt:         ProvideReceiverURLValidator(checker),
			in: &RegistrationV2{
				Webhooks: []Webhook{
					{
						ReceiverURLs: []string{"https://example.com",
							"http://example2.com"},
					},
				},
			},
			expectedErr: ErrInvalidInput,
		}, {
			description: "default case - invalid",
			opt:         ProvideReceiverURLValidator(checker),
			expectedErr: ErrInvalidType,
		},
	})
}

func TestProvideAlternativeURLValidator(t *testing.T) {
	checker, err := urlegit.New(urlegit.OnlyAllowSchemes("https"))
	require.NoError(t, err)
	require.NotNil(t, checker)

	run_tests(t, []optionTest{
		{
			description: "success, no checker",
			opt:         ProvideAlternativeURLValidator(nil),
			str:         "ProvideAlternativeURLValidator(nil)",
		}, {
			description: "success, with checker",
			opt:         ProvideAlternativeURLValidator(checker),
			in: &RegistrationV1{
				Config: DeliveryConfig{
					AlternativeURLs: []string{"https://example.com"},
				},
			},
			str: "ProvideAlternativeURLValidator(urlegit.Checker{ OnlyAllowSchemes('https') })",
		}, {
			description: "success, with checker and multiple urls",
			opt:         ProvideAlternativeURLValidator(checker),
			in: &RegistrationV1{
				Config: DeliveryConfig{
					AlternativeURLs: []string{"https://example.com", "https://example.org"},
				},
			},
			str: "ProvideAlternativeURLValidator(urlegit.Checker{ OnlyAllowSchemes('https') })",
		}, {
			description: "failure, with checker",
			opt:         ProvideAlternativeURLValidator(checker),
			in: &RegistrationV1{
				Config: DeliveryConfig{
					AlternativeURLs: []string{"http://example.com"},
				},
			},
			expectedErr: ErrInvalidInput,
		}, {
			description: "failure, with checker with multiple urls",
			opt:         ProvideAlternativeURLValidator(checker),
			in: &RegistrationV1{
				Config: DeliveryConfig{
					AlternativeURLs: []string{"https://example.com", "http://example.com"},
				},
			},
			expectedErr: ErrInvalidInput,
		}, {
			description: "failure - RegistrationV2",
			opt:         ProvideAlternativeURLValidator(checker),
			in:          &RegistrationV2{},
			expectedErr: ErrInvalidOption,
		}, {
			description: "default case - invalid",
			opt:         ProvideAlternativeURLValidator(checker),
			expectedErr: ErrInvalidType,
		},
	})
}

func TestNoUntil(t *testing.T) {
	run_tests(t, []optionTest{
		{
			description: "success, no until set",
			in:          &RegistrationV1{},
			opt:         NoUntil(),
			str:         "NoUntil()",
		}, {
			description: "detect until set",
			opt:         NoUntil(),
			in: &RegistrationV1{
				Until: time.Now(),
			},
			expectedErr: ErrInvalidInput,
		},
		{
			description: "failure - V2",
			opt:         NoUntil(),
			in:          &RegistrationV2{},
			expectedErr: ErrInvalidOption,
		},
		{
			description: "default case - invalid",
			opt:         NoUntil(),
			expectedErr: ErrInvalidType,
		},
	})
}

func run_tests(t *testing.T, tests []optionTest) {
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			assert := assert.New(t)
			var err error
			opts := append(tc.opts, tc.opt)
			switch r := tc.in.(type) {
			case *RegistrationV1:
				err = Validate(r, opts...)
			case *RegistrationV2:
				err = Validate(r, opts...)
			default:
				for _, o := range opts {
					err = o.Validate(nil)
					assert.ErrorIs(err, tc.expectedErr)
				}
			}
			assert.ErrorIs(err, tc.expectedErr)

			if tc.str != "" && tc.opt != nil {
				assert.Equal(tc.str, tc.opt.String())
			}
		})
	}
}
