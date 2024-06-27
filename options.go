// SPDX-FileCopyrightText: 2023 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package webhook

import (
	"time"

	"github.com/xmidt-org/urlegit"
)

// Error is an option that returns an error when called.
func Error(err error) Option {
	return errorOption{err: err}
}

type errorOption struct {
	err error
}

func (e errorOption) Validate(Validator) error {
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

func (a AlwaysValidOption) Validate(val Validator) error {
	return nil
}

func (a AlwaysValidOption) String() string {
	return "alwaysValidOption"
}

// AtLeastOneEvent makes sure there is at least one value in Events and ensures
// that all values should parse into regex.
func AtLeastOneEvent() Option {
	return atLeastOneEventOption{}
}

type atLeastOneEventOption struct{}

func (atLeastOneEventOption) Validate(val Validator) error {
	if err := val.ValidateOneEvent(); err != nil {
		return err
	}

	return nil
}

func (atLeastOneEventOption) String() string {
	return "AtLeastOneEvent()"
}

// EventRegexMustCompile ensures that all values in Events parse into valid regex.
func EventRegexMustCompile() Option {
	return eventRegexMustCompileOption{}
}

type eventRegexMustCompileOption struct{}

func (eventRegexMustCompileOption) Validate(val Validator) error {
	if err := val.ValidateEventRegex(); err != nil {
		return err
	}
	return nil
}

func (eventRegexMustCompileOption) String() string {
	return "EventRegexMustCompile()"
}

// DeviceIDRegexMustCompile ensures that all values in DeviceID parse into valid
// regex.
func DeviceIDRegexMustCompile() Option {
	return deviceIDRegexMustCompileOption{}
}

type deviceIDRegexMustCompileOption struct{}

func (deviceIDRegexMustCompileOption) Validate(val Validator) error {
	if err := val.ValidateDeviceId(); err != nil {
		return err
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

func (v validateRegistrationDurationOption) Validate(val Validator) error {
	if err := val.ValidateDuration(v.ttl); err != nil {
		return err
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

func (p provideTimeNowFuncOption) Validate(val Validator) error {
	val.SetNowFunc(p.nowFunc)
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
func ProvideFailureURLValidator(checker *urlegit.Checker) Option {
	return provideFailureURLValidatorOption{checker: checker}
}

type provideFailureURLValidatorOption struct {
	checker *urlegit.Checker
}

func (p provideFailureURLValidatorOption) Validate(v Validator) error {
	if p.checker == nil {
		return nil
	}

	if err := v.ValidateFailureURL(p.checker); err != nil {
		return err
	}
	return nil
}

func (p provideFailureURLValidatorOption) String() string {
	if p.checker == nil {
		return "ProvideFailureURLValidator(nil)"
	}
	return "ProvideFailureURLValidator(" + p.checker.String() + ")"
}

// ProvideReceiverURLValidator is an option that allows the caller to provide a
// URL validator that is used to validate the ReceiverURL.
func ProvideReceiverURLValidator(checker *urlegit.Checker) Option {
	return provideReceiverURLValidatorOption{checker: checker}
}

type provideReceiverURLValidatorOption struct {
	checker *urlegit.Checker
}

func (p provideReceiverURLValidatorOption) Validate(val Validator) error {
	if p.checker == nil {
		return nil
	}
	if err := val.ValidateReceiverURL(p.checker); err != nil {
		return err
	}

	return nil
}

func (p provideReceiverURLValidatorOption) String() string {
	if p.checker == nil {
		return "ProvideReceiverURLValidator(nil)"
	}
	return "ProvideReceiverURLValidator(" + p.checker.String() + ")"
}

// ProvideAlternativeURLValidator is an option that allows the caller to provide
// a URL validator that is used to validate the AlternativeURL.
func ProvideAlternativeURLValidator(checker *urlegit.Checker) Option {
	return provideAlternativeURLValidatorOption{checker: checker}
}

type provideAlternativeURLValidatorOption struct {
	checker *urlegit.Checker
}

func (p provideAlternativeURLValidatorOption) Validate(val Validator) error {
	if p.checker == nil {
		return nil
	}

	if err := val.ValidateAltURL(p.checker); err != nil {
		return err
	}
	return nil
}

func (p provideAlternativeURLValidatorOption) String() string {
	if p.checker == nil {
		return "ProvideAlternativeURLValidator(nil)"
	}
	return "ProvideAlternativeURLValidator(" + p.checker.String() + ")"
}

// NoUntil is an option that ensures that the Until field is not set.
func NoUntil() Option {
	return noUntilOption{}
}

type noUntilOption struct{}

func (noUntilOption) Validate(val Validator) error {
	if err := val.ValidateNoUntil(); err != nil {
		return err
	}
	return nil
}

func (noUntilOption) String() string {
	return "NoUntil()"
}

func Until(j time.Duration, m time.Duration, now func() time.Time) Option {
	return untilOption{
		jitter: j,
		max:    m,
		now:    now,
	}
}

type untilOption struct {
	jitter time.Duration
	max    time.Duration
	now    func() time.Time
}

func (u untilOption) Validate(val Validator) error {
	if err := val.ValidateUntil(u.jitter, u.max, u.now); err != nil {
		return err
	}
	return nil
}
func (untilOption) String() string {
	return "Until()"
}
