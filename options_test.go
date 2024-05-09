// SPDX-FileCopyrightText: 2022 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package webhook

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xmidt-org/urlegit"
)

type MockRegistration struct {
	mock.Mock
	Id         string
	Until      time.Time
	Duration   CustomDuration
	Events     []string
	FailureURL string
	Config     DeliveryConfig
	Matcher    MetadataMatcherConfig
}

func (m *MockRegistration) GetId() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockRegistration) GetUntil() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

func (m *MockRegistration) ValidateOneEvent() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRegistration) ValidateEventRegex() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRegistration) ValidateDeviceId() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRegistration) ValidateDuration(ttl time.Duration) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRegistration) ValidateURLs(c *urlegit.Checker) error {
	args := m.Called()
	return args.Error(0)
}

type optionTest struct {
	description string
	in          Register
	opt         Option
	opts        []Option
	str         string
	expectedErr error
}

func TestErrorOption(t *testing.T) {
	run_tests(t, []optionTest{
		{
			description: "success",
			str:         "foo",
		}, {
			description: "simple error",
			opt:         Error(ErrInvalidInput),
			str:         "Error('invalid input')",
			expectedErr: ErrInvalidInput,
		}, {
			description: "simple nil error",
			opt:         Error(nil),
			str:         "Error(nil)",
		},
	})
}

func TestAtLeastOneEventOption(t *testing.T) {
	run_tests(t, []optionTest{
		{
			description: "there is an event",
			opt:         ValidateEvents(),
			in:          &MockRegistration{Events: []string{"foo"}},
			str:         "ValidateEvents()",
		}, {
			description: "multiple events",
			opt:         ValidateEvents(),
			in:          &MockRegistration{Events: []string{"foo", "bar"}},
			str:         "ValidateEvents()",
		}, {
			description: "there are no events",
			opt:         ValidateEvents(),
			in:          &MockRegistration{},
			expectedErr: ErrInvalidInput,
		},
	})
}

func TestEventRegexMustCompile(t *testing.T) {
	run_tests(t, []optionTest{
		{
			description: "the regex compiles",
			opt:         ValidateEvents(),
			in:          &MockRegistration{Events: []string{"event.*"}},
			str:         "ValidateEvents()",
		}, {
			description: "multiple events",
			opt:         ValidateEvents(),
			in:          &MockRegistration{Events: []string{"magic-thing", "event.*"}},
			str:         "ValidateEvents()",
		}, {
			description: "failure",
			opt:         ValidateEvents(),
			in:          &MockRegistration{Events: []string{"("}},
			expectedErr: ErrInvalidInput,
		},
	})
}

func TestDeviceIDRegexMustCompile(t *testing.T) {
	run_tests(t, []optionTest{
		{
			description: "the regex compiles",
			opt:         DeviceIDRegexMustCompile(),
			in: &MockRegistration{
				Matcher: MetadataMatcherConfig{
					DeviceID: []string{"device.*"},
				},
			},
			str: "DeviceIDRegexMustCompile()",
		}, {
			description: "multiple device ids",
			opt:         DeviceIDRegexMustCompile(),
			in: &MockRegistration{
				Matcher: MetadataMatcherConfig{
					DeviceID: []string{"device.*", "magic-thing"},
				},
			},
			str: "DeviceIDRegexMustCompile()",
		}, {
			description: "failure",
			opt:         DeviceIDRegexMustCompile(),
			in: &MockRegistration{
				Matcher: MetadataMatcherConfig{
					DeviceID: []string{"("},
				},
			},
			expectedErr: ErrInvalidInput,
		},
	})
}

func TestValidateRegistrationDuration(t *testing.T) {
	now := func() time.Time {
		return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	run_tests(t, []optionTest{
		{
			description: "success with time in bounds",
			opt:         ValidateRegistrationDuration(5 * time.Minute),
			in: &MockRegistration{
				Duration: CustomDuration(4 * time.Minute),
			},
			str: "ValidateRegistrationDuration(5m0s)",
		}, {
			description: "success with time in bounds, exactly",
			opt:         ValidateRegistrationDuration(5 * time.Minute),
			in: &MockRegistration{
				Duration: CustomDuration(5 * time.Minute),
			},
		}, {
			description: "failure with time out of bounds",
			opt:         ValidateRegistrationDuration(5 * time.Minute),
			in: &MockRegistration{
				Duration: CustomDuration(6 * time.Minute),
			},
			expectedErr: ErrInvalidInput,
		}, {
			description: "success with max ttl ignored",
			opt:         ValidateRegistrationDuration(-5 * time.Minute),
			in: &MockRegistration{
				Duration: CustomDuration(1 * time.Minute),
			},
		}, {
			description: "success with max ttl ignored, 0 duration",
			opt:         ValidateRegistrationDuration(0),
			in: &MockRegistration{
				Duration: CustomDuration(1 * time.Minute),
			},
		}, {
			description: "success with until in bounds",
			opts: []Option{
				ProvideTimeNowFunc(now),
				ValidateRegistrationDuration(5 * time.Minute),
			},
			in: &MockRegistration{
				Until: time.Date(2021, 1, 1, 0, 4, 0, 0, time.UTC),
			},
		}, {
			description: "failure due to until being before now",
			opts: []Option{
				ValidateRegistrationDuration(5 * time.Minute),
				ProvideTimeNowFunc(now),
			},
			in: &MockRegistration{
				Until: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedErr: ErrInvalidInput,
		}, {
			description: "success with until exactly in bounds",
			opts: []Option{
				ProvideTimeNowFunc(now),
				ValidateRegistrationDuration(5 * time.Minute),
			},
			in: &MockRegistration{
				Until: time.Date(2021, 1, 1, 0, 5, 0, 0, time.UTC),
			},
		}, {
			description: "failure due to the options being out of order",
			opts: []Option{
				ValidateRegistrationDuration(5 * time.Minute),
				ProvideTimeNowFunc(now),
			},
			in: &MockRegistration{
				Until: time.Date(2021, 1, 1, 0, 5, 0, 0, time.UTC),
			},
			expectedErr: ErrInvalidInput,
		}, {
			description: "failure with until out of bounds",
			opts: []Option{
				ProvideTimeNowFunc(now),
				ValidateRegistrationDuration(5 * time.Minute),
			},
			in: &MockRegistration{
				Until: time.Date(2021, 1, 1, 0, 6, 0, 0, time.UTC),
			},
			expectedErr: ErrInvalidInput,
		}, {
			description: "success with until just needing to be present",
			opts: []Option{
				ProvideTimeNowFunc(now),
				ValidateRegistrationDuration(0),
			},
			in: &MockRegistration{
				Until: time.Date(2021, 1, 1, 0, 6, 0, 0, time.UTC),
			},
		}, {
			description: "failure, both expirations set",
			opt:         ValidateRegistrationDuration(5 * time.Minute),
			in: &MockRegistration{
				Duration: CustomDuration(1 * time.Minute),
				Until:    time.Date(2021, 1, 1, 0, 4, 0, 0, time.UTC),
			},
			expectedErr: ErrInvalidInput,
		}, {
			description: "failure, no expiration set",
			opt:         ValidateRegistrationDuration(5 * time.Minute),
			expectedErr: ErrInvalidInput,
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
			opt:         ProvideURLValidator(nil),
			str:         "ProvideFailureURLValidator(nil)",
		}, {
			description: "success, with checker",
			opt:         ProvideURLValidator(checker),
			in: &MockRegistration{
				FailureURL: "https://example.com",
			},
			str: "ProvideFailureURLValidator(urlegit.Checker{ OnlyAllowSchemes('https') })",
		}, {
			description: "failure, with checker",
			opt:         ProvideURLValidator(checker),
			in: &MockRegistration{
				FailureURL: "http://example.com",
			},
			expectedErr: ErrInvalidInput,
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
			opt:         ProvideURLValidator(nil),
			str:         "ProvideReceiverURLValidator(nil)",
		}, {
			description: "success, with checker",
			opt:         ProvideURLValidator(checker),
			in: &MockRegistration{
				Config: DeliveryConfig{
					ReceiverURL: "https://example.com",
				},
			},
			str: "ProvideReceiverURLValidator(urlegit.Checker{ OnlyAllowSchemes('https') })",
		}, {
			description: "failure, with checker",
			opt:         ProvideURLValidator(checker),
			in: &MockRegistration{
				Config: DeliveryConfig{
					ReceiverURL: "http://example.com",
				},
			},
			expectedErr: ErrInvalidInput,
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
			opt:         ProvideURLValidator(nil),
			str:         "ProvideAlternativeURLValidator(nil)",
		}, {
			description: "success, with checker",
			opt:         ProvideURLValidator(checker),
			in: &MockRegistration{
				Config: DeliveryConfig{
					AlternativeURLs: []string{"https://example.com"},
				},
			},
			str: "ProvideAlternativeURLValidator(urlegit.Checker{ OnlyAllowSchemes('https') })",
		}, {
			description: "success, with checker and multiple urls",
			opt:         ProvideURLValidator(checker),
			in: &MockRegistration{
				Config: DeliveryConfig{
					AlternativeURLs: []string{"https://example.com", "https://example.org"},
				},
			},
			str: "ProvideAlternativeURLValidator(urlegit.Checker{ OnlyAllowSchemes('https') })",
		}, {
			description: "failure, with checker",
			opt:         ProvideURLValidator(checker),
			in: &MockRegistration{
				Config: DeliveryConfig{
					AlternativeURLs: []string{"http://example.com"},
				},
			},
			expectedErr: ErrInvalidInput,
		}, {
			description: "failure, with checker with multiple urls",
			opt:         ProvideURLValidator(checker),
			in: &MockRegistration{
				Config: DeliveryConfig{
					AlternativeURLs: []string{"https://example.com", "http://example.com"},
				},
			},
			expectedErr: ErrInvalidInput,
		},
	})
}

func TestNoUntil(t *testing.T) {
	run_tests(t, []optionTest{
		{
			description: "success, no until set",
			opt:         NoUntil(),
			str:         "NoUntil()",
		}, {
			description: "detect until set",
			opt:         NoUntil(),
			in: &MockRegistration{
				Until: time.Now(),
			},
			expectedErr: ErrInvalidInput,
		},
	})
}
func run_tests(t *testing.T, tests []optionTest) {
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			assert := assert.New(t)

			opts := append(tc.opts, tc.opt)
			err := Validate(tc.in, opts...)

			assert.ErrorIs(err, tc.expectedErr)

			if tc.str != "" && tc.opt != nil {
				assert.Equal(tc.str, tc.opt.String())
			}
		})
	}
}
