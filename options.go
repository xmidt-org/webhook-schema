// SPDX-FileCopyrightText: 2023 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package webhook

import (
	"fmt"
	"regexp"
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

func (e errorOption) Validate(*Registration) error {
	return error(e.err)
}

func (e errorOption) String() string {
	if e.err == nil {
		return "Error(nil)"
	}
	return "Error('" + e.err.Error() + "')"
}

// AtLeastOneEvent makes sure there is at least one value in Events and ensures
// that all values should parse into regex.
func AtLeastOneEvent() Option {
	return atLeastOneEventOption{}
}

type atLeastOneEventOption struct{}

func (atLeastOneEventOption) Validate(r *Registration) error {
	if len(r.Events) == 0 {
		return fmt.Errorf("%w: cannot have zero events", ErrInvalidInput)
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

func (eventRegexMustCompileOption) Validate(r *Registration) error {
	for _, e := range r.Events {
		_, err := regexp.Compile(e)
		if err != nil {
			return fmt.Errorf("%w: unable to compile matching", ErrInvalidInput)
		}
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

func (deviceIDRegexMustCompileOption) Validate(r *Registration) error {
	for _, e := range r.Matcher.DeviceID {
		_, err := regexp.Compile(e)
		if err != nil {
			return fmt.Errorf("%w: unable to compile matching", ErrInvalidInput)
		}
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

func (v validateRegistrationDurationOption) Validate(r *Registration) error {
	if v.ttl <= 0 {
		v.ttl = time.Duration(0)
	}

	if v.ttl != 0 && v.ttl < time.Duration(r.Duration) {
		return fmt.Errorf("%w: the registration is for too long", ErrInvalidInput)
	}

	if r.Until.IsZero() && r.Duration == 0 {
		return fmt.Errorf("%w: either Duration or Until must be set", ErrInvalidInput)
	}

	if !r.Until.IsZero() && r.Duration != 0 {
		return fmt.Errorf("%w: only one of Duration or Until may be set", ErrInvalidInput)
	}

	if !r.Until.IsZero() {
		nowFunc := time.Now
		if r.nowFunc != nil {
			nowFunc = r.nowFunc
		}

		now := nowFunc()
		if v.ttl != 0 && r.Until.After(now.Add(v.ttl)) {
			return fmt.Errorf("%w: the registration is for too long", ErrInvalidInput)
		}

		if r.Until.Before(now) {
			return fmt.Errorf("%w: the registration has already expired", ErrInvalidInput)
		}
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

func (p provideTimeNowFuncOption) Validate(r *Registration) error {
	r.nowFunc = p.nowFunc
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

func (p provideFailureURLValidatorOption) Validate(r *Registration) error {
	if p.checker == nil {
		return nil
	}

	if r.FailureURL != "" {
		if err := p.checker.Text(r.FailureURL); err != nil {
			return fmt.Errorf("%w: failure url is invalid", ErrInvalidInput)
		}
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

func (p provideReceiverURLValidatorOption) Validate(r *Registration) error {
	if p.checker == nil {
		return nil
	}

	if r.Config.ReceiverURL != "" {
		if err := p.checker.Text(r.Config.ReceiverURL); err != nil {
			return fmt.Errorf("%w: failure url is invalid", ErrInvalidInput)
		}
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

func (p provideAlternativeURLValidatorOption) Validate(r *Registration) error {
	if p.checker == nil {
		return nil
	}

	for _, url := range r.Config.AlternativeURLs {
		if err := p.checker.Text(url); err != nil {
			return fmt.Errorf("%w: failure url is invalid", ErrInvalidInput)
		}
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

func (noUntilOption) Validate(r *Registration) error {
	if !r.Until.IsZero() {
		return fmt.Errorf("%w: Until is not allowed", ErrInvalidInput)
	}
	return nil
}

func (noUntilOption) String() string {
	return "NoUntil()"
}
