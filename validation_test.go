// SPDX-FileCopyrightText: 2022 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package webhook

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xmidt-org/urlegit"
)

var mockNow func() time.Time = func() time.Time {
	return time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
}

func TestBuildURLChecker(t *testing.T) {
	tests := []struct {
		description string
		config      ValidatorConfig
	}{
		{
			description: "allow special use hosts and ips",
			config: ValidatorConfig{
				URL: URLVConfig{
					HTTPSOnly:            true,
					AllowLoopback:        false,
					AllowIP:              false,
					AllowSpecialUseHosts: true,
					AllowSpecialUseIPs:   true,
					InvalidHosts:         []string{"example.com", "localhost"},
					InvalidSubnets:       []string{"192.168.0.0/16"},
				},
				TTL: TTLVConfig{
					Max:    time.Hour,
					Jitter: time.Minute,
					Now:    time.Now,
				},
			},
		},
		{
			description: "!allow special use hosts and ips",
			config: ValidatorConfig{
				URL: URLVConfig{
					HTTPSOnly:            true,
					AllowLoopback:        false,
					AllowIP:              false,
					AllowSpecialUseHosts: false,
					AllowSpecialUseIPs:   false,
					InvalidHosts:         []string{"example.com", "localhost"},
					InvalidSubnets:       []string{"192.168.0.0/16"},
				},
				TTL: TTLVConfig{
					Max:    time.Hour,
					Jitter: time.Minute,
					Now:    time.Now,
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			checker, err := buildURLChecker(tc.config)
			assert.NoError(t, err)
			assert.NotNil(t, checker)
		})
	}

}

func TestBuildValidators(t *testing.T) {
	config := ValidatorConfig{
		URL: URLVConfig{
			HTTPSOnly:            true,
			AllowLoopback:        false,
			AllowIP:              false,
			AllowSpecialUseHosts: true,
			AllowSpecialUseIPs:   true,
			InvalidHosts:         []string{"example.com", "localhost"},
			InvalidSubnets:       []string{"192.168.0.0/16"},
		},
		TTL: TTLVConfig{
			Max:    time.Hour,
			Jitter: time.Minute,
			Now:    time.Now,
		},
	}

	validators, err := BuildValidators(config)
	assert.NoError(t, err)
	assert.NotNil(t, validators)

}

func TestValidatePass(t *testing.T) {
	tests := []struct {
		description string
		v           Validator
		max         time.Duration //TODO: delete max and just use the ttl.max from config
		config      ValidatorConfig
		ifChecker   bool
	}{
		{
			description: "regV1 nil checker",
			v: &RegistrationV1{
				Matcher: MetadataMatcherConfig{
					DeviceID: []string{"[a-z0-9]"},
				},
				Events:   []string{"Offline"},
				Duration: CustomDuration(2),
			},
			max:       time.Duration(2),
			ifChecker: false,
		},
		{
			description: "regV1 ttl <= 0",
			v: &RegistrationV1{
				Matcher: MetadataMatcherConfig{
					DeviceID: []string{"[a-z0-9]"},
				},
				Events:   []string{"Offline"},
				Duration: CustomDuration(2),
			},
			max:       time.Duration(0),
			ifChecker: false,
		},
		{
			description: "regV1 empty v1.nowFunc != nil",
			v: &RegistrationV1{
				nowFunc: mockNow,
				Matcher: MetadataMatcherConfig{
					DeviceID: []string{"[a-z0-9]"},
				},
				Events: []string{"Offline"},
				Until:  mockNow(),
			},
			max:       time.Duration(0),
			ifChecker: false,
		},
		{
			description: "regV1 nonnil checker",
			v: &RegistrationV1{
				Matcher: MetadataMatcherConfig{
					DeviceID: []string{"[a-z0-9]"},
				},
				Events:   []string{"Offline"},
				Duration: CustomDuration(2),
			},
			max:       time.Duration(2),
			ifChecker: true,
			config: ValidatorConfig{
				URL: URLVConfig{
					HTTPSOnly:            true,
					AllowLoopback:        false,
					AllowIP:              false,
					AllowSpecialUseHosts: true,
					AllowSpecialUseIPs:   true,
					InvalidHosts:         []string{"example.com", "localhost"},
					InvalidSubnets:       []string{"192.168.0.0/16"},
				},
			},
		},
		//TODO: need to add in validation for regV2
	}
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			var checker *urlegit.Checker
			var err error
			if tc.ifChecker {
				checker, err = buildURLChecker(tc.config)
				assert.Nil(t, err)
			} else {
				checker = nil
			}
			opts := []Option{AtLeastOneEvent(), EventRegexMustCompile(), DeviceIDRegexMustCompile(), ValidateRegistrationDuration(tc.max), ProvideAlternativeURLValidator(checker), ProvideFailureURLValidator(checker), ProvideReceiverURLValidator(checker)}
			err = Validate(tc.v, opts)
			assert.Nil(t, err, fmt.Errorf("received '%v' when nil was expected", err))

		})
	}
}

func TestValidateFail(t *testing.T) {
	tests := []struct {
		description string
		v           Validator
		opts        []Option
		max         time.Duration
		config      ValidatorConfig
		ifChecker   bool
		expectedErr error
	}{
		{
			description: "regV1 no events",
			v: &RegistrationV1{
				Events: []string{},
			},
			opts:        []Option{AtLeastOneEvent()},
			expectedErr: ErrInvalidInput,
		},
		{
			description: "regV1 invalid event regext",
			v: &RegistrationV1{
				Events: []string{`\M`},
			},
			opts:        []Option{EventRegexMustCompile()},
			expectedErr: ErrInvalidInput,
		},
		{
			description: "regV1 invalid device id regex",
			v: &RegistrationV1{
				Matcher: MetadataMatcherConfig{
					DeviceID: []string{"", `\M`}},
			},
			opts:        []Option{DeviceIDRegexMustCompile()},
			expectedErr: ErrInvalidInput,
		},
		{
			description: "regV1 invalid duration - ttl < time.Duration",
			v: &RegistrationV1{
				Duration: CustomDuration(5),
			},
			expectedErr: ErrInvalidInput,
			opts:        []Option{ValidateRegistrationDuration(time.Duration(3))},
		},
		{
			description: "regV1 invalid duration - Duration and Until set",
			v: &RegistrationV1{
				Duration: CustomDuration(5),
				Until:    time.Now(),
			},
			expectedErr: ErrInvalidInput,
			opts:        []Option{ValidateRegistrationDuration(time.Duration(10))},
		},
		{
			description: "regV1 invalid duration - neither duration nor until set",
			v: &RegistrationV1{
				Duration: CustomDuration(0),
			},
			expectedErr: ErrInvalidInput,
			opts:        []Option{ValidateRegistrationDuration(time.Duration(10))},
		},
		{
			description: "regV1 invalid duration - Until time Before(now)",
			v: &RegistrationV1{
				Until: time.Date(2024, 06, 11, 9, 50, 0, 0, time.UTC),
			},
			expectedErr: ErrInvalidInput,
			opts:        []Option{ValidateRegistrationDuration(time.Duration(10))},
		},
		{
			description: "regV1 invalid duration - TTL < Until",
			v: &RegistrationV1{
				Until: time.Now().Add(time.Minute * 15),
			},
			expectedErr: ErrInvalidInput,
			opts:        []Option{ValidateRegistrationDuration(time.Duration(10))},
		},
		{
			description: "regV1 no until",
			v: &RegistrationV1{
				Until: time.Now().Add(time.Minute * 15),
			},
			opts:        []Option{NoUntil()},
			expectedErr: ErrInvalidInput,
		},
	}
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			err := Validate(tc.v, tc.opts)
			assert.NotNil(t, err, fmt.Errorf("expected an err - received Nil"))
			assert.True(t, errors.Is(err, tc.expectedErr),
				fmt.Errorf("error [%v] doesn't contain error [%v] in its err chain",
					err, tc.expectedErr))
		})
	}
}

func TestValidateNoUntil(t *testing.T) {
	v1 := &RegistrationV1{}
	err := NoUntil().Validate(v1)
	assert.Nil(t, err)
}
func TestValidateUnil(t *testing.T) {
	tests := []struct {
		description string
		v           Validator
		opts        []Option
		max         time.Duration
		jitter      time.Duration
		now         func() time.Time
		expectedErr error
	}{
		{
			description: "valid - empty now func",
			v: &RegistrationV1{
				Until: time.Now(),
			},
			max:         (time.Minute * 5),
			jitter:      (time.Minute * 5),
			expectedErr: nil,
		},
		{
			description: "valid - Until = 0",
			v:           &RegistrationV1{},
			max:         (time.Minute * 5),
			jitter:      (time.Minute * 5),
			expectedErr: nil,
			now:         mockNow,
		},
		{
			description: "invalid - maxTTL < 0",
			v:           &RegistrationV1{},
			max:         (time.Microsecond * -5),
			expectedErr: ErrInvalidInput,
		},
		{
			description: "invalid - jitter < 0",
			v:           &RegistrationV1{},
			max:         (time.Millisecond * 5),
			jitter:      (time.Microsecond * -5),
			expectedErr: ErrInvalidInput,
		},
		{
			description: "invalid - proposed after limit",
			v: &RegistrationV1{
				Until: time.Now().Add(time.Minute * 60),
			},
			max:         (time.Minute * 5),
			jitter:      (time.Minute * 5),
			expectedErr: ErrInvalidInput,
		},
	}
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			err := tc.v.ValidateUntil(tc.jitter, tc.max, tc.now)
			if tc.expectedErr != nil {
				assert.NotNil(t, err)
				assert.True(t, errors.Is(err, tc.expectedErr),
					fmt.Errorf("error [%v] doesn't contain error [%v] in its err chain",
						err, tc.expectedErr))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestSetNowFunc(t *testing.T) {
	v1 := &RegistrationV1{}
	// v2 := &RegistrationV2{}

	v1.SetNowFunc(mockNow)
	// v2.SetNowFunc(nowFunc)

	assert.Equal(t, mockNow(), v1.nowFunc())
	// assert.Equal(t, nowFunc, v2.nowFunc)
}
func TestValidateFailureURL(t *testing.T) {
	tests := []struct {
		description string
		v           Validator
		config      ValidatorConfig
		expectedErr bool
	}{
		{
			description: "valid failure URL",
			v: &RegistrationV1{
				FailureURL: "https://example.com/failure",
			},
			config: ValidatorConfig{
				URL: URLVConfig{
					HTTPSOnly:            true,
					AllowLoopback:        false,
					AllowIP:              false,
					AllowSpecialUseHosts: true,
					AllowSpecialUseIPs:   true,
					InvalidHosts:         []string{"localhost"},
					InvalidSubnets:       []string{"192.168.0.0/16"},
				},
			},
			expectedErr: false,
		},
		{
			description: "invalid failure URL",
			v: &RegistrationV1{
				FailureURL: "ftp://example.com/failure",
			},
			config: ValidatorConfig{
				URL: URLVConfig{
					HTTPSOnly:            true,
					AllowLoopback:        false,
					AllowIP:              false,
					AllowSpecialUseHosts: true,
					AllowSpecialUseIPs:   true,
					InvalidHosts:         []string{"localhost"},
					InvalidSubnets:       []string{"192.168.0.0/16"},
				},
			},
			expectedErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			checker, err := buildURLChecker(tc.config)
			assert.NoError(t, err)

			err = tc.v.ValidateFailureURL(checker)
			if tc.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateReceiverURL(t *testing.T) {
	tests := []struct {
		description string
		v           Validator
		config      ValidatorConfig
		expectedErr bool
	}{
		{
			description: "valid receiver URL",
			v: &RegistrationV1{
				Config: DeliveryConfig{
					ReceiverURL: "https://example.com/receiver",
				},
			},
			config: ValidatorConfig{
				URL: URLVConfig{
					HTTPSOnly:            true,
					AllowLoopback:        false,
					AllowIP:              false,
					AllowSpecialUseHosts: true,
					AllowSpecialUseIPs:   true,
					InvalidHosts:         []string{"localhost"},
					InvalidSubnets:       []string{"192.168.0.0/16"},
				},
			},
			expectedErr: false,
		},
		{
			description: "invalid receiver URL",
			v: &RegistrationV1{
				Config: DeliveryConfig{
					ReceiverURL: "ftp://example.com/receiver",
				},
			},
			config: ValidatorConfig{
				URL: URLVConfig{
					HTTPSOnly:            true,
					AllowLoopback:        false,
					AllowIP:              false,
					AllowSpecialUseHosts: true,
					AllowSpecialUseIPs:   true,
					InvalidHosts:         []string{"localhost"},
					InvalidSubnets:       []string{"192.168.0.0/16"},
				},
			},
			expectedErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			checker, err := buildURLChecker(tc.config)
			assert.NoError(t, err)

			err = tc.v.ValidateReceiverURL(checker)
			if tc.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateAltURL(t *testing.T) {
	tests := []struct {
		description string
		v           Validator
		config      ValidatorConfig
		expectedErr bool
	}{
		{
			description: "valid alternative URL",
			v: &RegistrationV1{
				Config: DeliveryConfig{
					AlternativeURLs: []string{"https://example.com/alt"},
				},
			},
			config: ValidatorConfig{
				URL: URLVConfig{
					HTTPSOnly:            true,
					AllowLoopback:        false,
					AllowIP:              false,
					AllowSpecialUseHosts: true,
					AllowSpecialUseIPs:   true,
					InvalidHosts:         []string{"localhost"},
					InvalidSubnets:       []string{"192.168.0.0/16"},
				},
			},
			expectedErr: false,
		},
		{
			description: "invalid alternative URL",
			v: &RegistrationV1{
				Config: DeliveryConfig{
					AlternativeURLs: []string{"ftp://example.com/alt"},
				},
			},
			config: ValidatorConfig{
				URL: URLVConfig{
					HTTPSOnly:            true,
					AllowLoopback:        false,
					AllowIP:              false,
					AllowSpecialUseHosts: true,
					AllowSpecialUseIPs:   true,
					InvalidHosts:         []string{"localhost"},
					InvalidSubnets:       []string{"192.168.0.0/16"},
				},
			},
			expectedErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			checker, err := buildURLChecker(tc.config)
			assert.NoError(t, err)

			err = tc.v.ValidateAltURL(checker)
			if tc.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
