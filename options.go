// SPDX-FileCopyrightText: 2023 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package webhook

import (
	"fmt"
	"time"

	"github.com/xmidt-org/urlegit"
	"go.uber.org/multierr"
)

// Error is an option that returns an error when called.
func Error(err error) Option {
	return errorOption{err: err}
}

type errorOption struct {
	err error
}

func (e errorOption) Validate(r Register) error {
	return error(e.err)
}

func (e errorOption) String() string {
	if e.err == nil {
		return "Error(nil)"
	}
	return "Error('" + e.err.Error() + "')"
}

func AlwaysValid() Option {
	return AlwaysValidOption{}
}

type AlwaysValidOption struct{}

func (a AlwaysValidOption) Validate(r Register) error {
	return nil
}

func (a AlwaysValidOption) String() string {
	return "alwaysValidOption"
}

func ValidateEvents() Option {
	return validateEvents{}
}

type validateEvents struct{}

// AtLeastOneEvent makes sure there is at least one value in Events and ensures
// that all values should parse into regex.
func (validateEvents) Validate(r Register) error {
	var errs error
	err := r.ValidateOneEvent()
	if err != nil {
		errs = multierr.Append(errs, err)
		err = nil
	}

	err = r.ValidateEventRegex()
	if err != nil {
		errs = multierr.Append(errs, err)
	}
	if errs != nil {
		return fmt.Errorf("%w:%w", ErrInvalidInput, errs)
	}

	return nil
}

func (validateEvents) String() string {
	return "ValidateEvents()"
}

// EventRegexMustCompile ensures that all values in Events parse into valid regex.

// DeviceIDRegexMustCompile ensures that all values in DeviceID parse into valid
// regex.
func DeviceIDRegexMustCompile() Option {
	return deviceIDRegexMustCompileOption{}
}

type deviceIDRegexMustCompileOption struct{}

func (deviceIDRegexMustCompileOption) Validate(r Register) error {
	err := r.ValidateDeviceId()
	if err != nil {
		return fmt.Errorf("%w:%w", ErrInvalidInput, err)
	}
	return nil
}

func (deviceIDRegexMustCompileOption) String() string {
	return "DeviceIDRegexMustCompile()"
}

// ValidateRegistrationDuration ensures that the requsted registration duration
// of a webhook is valid.  This option checks the values set in either the
// Duration or Until fields of the webhook. If the ttl is less than or equal to
// zero, this option will not boundary check the registration duration, but will
// still ensure that the Duration or Until fields are set.
func ValidateRegistrationDuration(ttl time.Duration) Option {
	return validateRegistrationDurationOption{ttl: ttl}
}

type validateRegistrationDurationOption struct {
	ttl time.Duration
}

func (v validateRegistrationDurationOption) Validate(r Register) error {
	if v.ttl <= 0 {
		v.ttl = time.Duration(0)
	}

	err := r.ValidateDuration(v.ttl)
	if err != nil {
		return fmt.Errorf("%w:%w", ErrInvalidInput, err)
	}

	return nil
}

func (v validateRegistrationDurationOption) String() string {
	return "ValidateRegistrationDuration(" + v.ttl.String() + ")"
}

// ProvideTimeNowFunc is an option that allows the caller to provide a function
// that returns the current time.  This is used for testing.
func ProvideTimeNowFunc(nowFunc func() time.Time) Option {
	return provideTimeNowFuncOption{nowFunc: nowFunc}
}

type provideTimeNowFuncOption struct {
	nowFunc func() time.Time
}

func (p provideTimeNowFuncOption) Validate(r Register) error {
	// r.nowFunc = p.nowFunc
	return nil
}

func (p provideTimeNowFuncOption) String() string {
	if p.nowFunc == nil {
		return "ProvideTimeNowFunc(nil)"
	}
	return "ProvideTimeNowFunc(func)"
}

// ProvideFailureURLValidator is an option that allows the caller to provide a
// URL validator that is used to validate the FailureURL.

func ProvideURLValidator(checker *urlegit.Checker) Option {
	return provideURLValidator{checker: checker}
}

type provideURLValidator struct {
	checker *urlegit.Checker
}

func (p provideURLValidator) String() string {
	if p.checker == nil {
		return "ProvideURLValidator(nil)"
	}
	return "ProvideURLValidator(" + p.checker.String() + ")"
}

func (p provideURLValidator) Validate(r Register) error {
	if p.checker == nil {
		return nil
	}
	err := r.ValidateURLs(p.checker)
	if err != nil {
		return fmt.Errorf("%w:%w", ErrInvalidInput, err)
	}

	return nil
}

// ProvideReceiverURLValidator is an option that allows the caller to provide a
// URL validator that is used to validate the ReceiverURL.

// ProvideAlternativeURLValidator is an option that allows the caller to provide
// a URL validator that is used to validate the AlternativeURL.

// NoUntil is an option that ensures that the Until field is not set.
func NoUntil() Option {
	return noUntilOption{}
}

type noUntilOption struct{}

func (noUntilOption) Validate(r Register) error {
	until := r.GetUntil()
	if !until.IsZero() {
		return fmt.Errorf("%w: Until is not allowed", ErrInvalidInput)
	}
	return nil
}

func (noUntilOption) String() string {
	return "NoUntil()"
}
