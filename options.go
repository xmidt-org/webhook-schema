// SPDX-FileCopyrightText: 2023 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package webhook

import (
	"fmt"
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

func (e errorOption) Validate(any) error {
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

func (a AlwaysValidOption) Validate(any) error {
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

func (atLeastOneEventOption) Validate(i any) error {
	switch r := i.(type) {
	case *RegistrationV1:
		return r.ValidateOneEvent()
	case *RegistrationV2:
		return fmt.Errorf("%w: RegistrationV2 does not have an events field to validate", ErrInvalidType)
	default:
		return ErrUknownType
	}
}

func (atLeastOneEventOption) String() string {
	return "AtLeastOneEvent()"
}

// EventRegexMustCompile ensures that all values in Events parse into valid regex.
func EventRegexMustCompile() Option {
	return eventRegexMustCompileOption{}
}

type eventRegexMustCompileOption struct{}

func (eventRegexMustCompileOption) Validate(i any) error {
	switch r := i.(type) {
	case *RegistrationV1:
		return r.ValidateEventRegex()
	case *RegistrationV2:
		return r.ValidateEventRegex()
	default:
		return ErrUknownType
	}
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

func (deviceIDRegexMustCompileOption) Validate(i any) error {
	switch r := i.(type) {
	case *RegistrationV1:
		return r.ValidateDeviceId()
	case *RegistrationV2:
		//Matcher description is for Events. Are we not matching for DeviceId in Reg2?
	default:
		return ErrUknownType
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

func (v validateRegistrationDurationOption) Validate(i any) error {
	switch r := i.(type) {
	case *RegistrationV1:
		return r.ValidateDuration(v.ttl)
	case *RegistrationV2:
		return r.ValidateDuration()
	default:
		return ErrUknownType
	}
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

func (p provideTimeNowFuncOption) Validate(i any) error {
	switch r := i.(type) {
	case *RegistrationV1:
		r.SetNowFunc(p.nowFunc)
	}

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

func (p provideFailureURLValidatorOption) Validate(i any) error {
	var failureURL string
	//TODO: do we want to move this check to be inside each case statement?
	if p.checker == nil {
		return nil
	}

	switch r := i.(type) {
	case *RegistrationV1:
		failureURL = r.FailureURL
	case *RegistrationV2:
		failureURL = r.FailureURL
	default:
		return ErrUknownType
	}

	if failureURL != "" {
		if err := p.checker.Text(failureURL); err != nil {
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

func (p provideReceiverURLValidatorOption) Validate(i any) error {
	if p.checker == nil {
		return nil
	}

	switch r := i.(type) {
	case *RegistrationV1:
		return r.ValidateReceiverURL(p.checker)
	case *RegistrationV2:
		return r.ValidateReceiverURL(p.checker)
	default:
		return ErrUknownType
	}
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

func (p provideAlternativeURLValidatorOption) Validate(i any) error {
	if p.checker == nil {
		return nil
	}

	switch r := i.(type) {
	case *RegistrationV1:
		return r.ValidateAltURL(p.checker)
	case *RegistrationV2:
		return fmt.Errorf("%w: RegistrationV2 does not have an alternative urls field. Use ProvideReceiverURLValidator() to validate all non-failure urls", ErrInvalidType)
	default:
		return ErrUknownType
	}
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

func (noUntilOption) Validate(i any) error {

	switch r := i.(type) {
	case *RegistrationV1:
		return r.ValidateNoUntil()
	case *RegistrationV2:
		return fmt.Errorf("%w: RegistrationV2 does not use an Until field", ErrInvalidType)
	default:
		return ErrUknownType
	}

}

func (noUntilOption) String() string {
	return "NoUntil()"
}
