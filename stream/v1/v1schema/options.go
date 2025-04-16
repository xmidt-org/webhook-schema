// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package v1schema

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/xmidt-org/urlegit"

	"go.uber.org/multierr"
)

var (
	ErrEmptySchema                  = errors.New("empty schema error")
	ErrEmptyEvents                  = errors.New("cannot have zero events")
	ErrExpired                      = errors.New("registration has already expired")
	ErrEventRegexCompilerFailure    = errors.New("failed to compile Event regexp")
	ErrDeviceIDRegexCompilerFailure = errors.New("failed to compile Matcher.DeviceID regexp")
	ErrInvalidReceiverURL           = errors.New("invalid ReceiverURL")
	ErrInvalidFailureURL            = errors.New("invalid FailureURL")
	ErrInvalidDuration              = errors.New("invalid Duration")
	ErrInvalidCheckUntil            = errors.New("invalid CheckUntil")
	ErrInvalidUntil                 = errors.New("invalid Until")
	ErrAlternativeURLs              = errors.New("invalid AlternativeURLs")
)

func Address(addr string) Option {
	return OptionFunc(func(s *Schema) error {
		s.Address = addr

		return nil
	})
}

func Matcher(regexps []string) Option {
	return OptionFunc(func(s *Schema) error {
		s.Matcher.DeviceID = regexps

		return nil
	})
}

func UntilNowAddDuration() Option {
	return OptionFunc(func(s *Schema) error {
		s.Until = time.Now().Add(time.Duration(s.Duration))

		return nil
	})
}
func SetStream(S Schema) Option {
	return OptionFunc(func(s *Schema) error {
		*s = S

		return nil
	})
}
func AlwaysValid() Option {
	return OptionFunc(func(s *Schema) error {
		return nil
	})
}

func NotEmptyValidator() Option {
	return OptionFunc(func(s *Schema) error {
		if s.Webhook.ReceiverURL == "" {
			return ErrEmptySchema
		}

		return nil
	})
}

// AtleastOneEventValidator makes sure there is at least one value in Events and ensures
// that all values should parse into regex.
func AtleastOneEventValidator() Option {
	return OptionFunc(func(s *Schema) error {
		if len(s.Events) == 0 {
			return ErrEmptyEvents
		}

		return nil
	})
}

// EventRegexValidator ensures that all values in Events parse into valid regex.

func EventRegexValidator() Option {
	return OptionFunc(func(s *Schema) (errs error) {
		for _, e := range s.Events {
			if _, err := regexp.Compile(e); err != nil {
				errs = multierr.Append(errs, fmt.Errorf("%w: `%s`: %s", ErrEventRegexCompilerFailure, e, err))
			}
		}

		return errs
	})
}

// DeviceIDRegexValidator ensures that all values in DeviceID parse into valid
// regex.
func DeviceIDRegexValidator() Option {
	return OptionFunc(func(s *Schema) (errs error) {
		for _, e := range s.Matcher.DeviceID {
			if _, err := regexp.Compile(e); err != nil {
				errs = multierr.Append(errs, fmt.Errorf("%w: `%s`: %s", ErrDeviceIDRegexCompilerFailure, e, err))
			}
		}

		return errs
	})
}

// DurationValidator ensures that the requsted registration duration
// of a webhook is valid.  This option checks the values set in either the
// Duration or Until fields of the webhook. If the ttl is less than or equal to
// zero, this option will not boundary check the registration duration, but will
// still ensure that the Duration or Until fields are set.
func DurationValidator(nowFunc func() time.Time, ttl time.Duration) Option {
	return OptionFunc(func(s *Schema) (errs error) {
		if ttl <= 0 {
			ttl = time.Duration(0)
		}
		if ttl != 0 && ttl < time.Duration(s.Duration) {
			errs = multierr.Append(errs, fmt.Errorf("%w: the registration is for too long", ErrInvalidDuration))

		}
		if s.Until.IsZero() && s.Duration == 0 {
			errs = multierr.Append(errs, fmt.Errorf("%w: either Duration or Until must be set", ErrInvalidDuration))
		}
		if !s.Until.IsZero() && s.Duration != 0 {
			errs = multierr.Append(errs, fmt.Errorf("%w: only one of Duration or Until may be set", ErrInvalidDuration))
		}
		if !s.Until.IsZero() {
			now := nowFunc()
			if ttl != 0 && s.Until.After(now.Add(ttl)) {
				errs = multierr.Append(errs, fmt.Errorf("%w: the registration is for too long", ErrInvalidDuration))
			}

			if s.Until.Before(now) {
				errs = multierr.Append(errs, fmt.Errorf("%w: the registration has already expired", ErrInvalidDuration))
			}
		}

		return errs
	})
}

// CheckUntilValidator is an option that ensures that the Until field is valid and not already expired.
func CheckUntilValidator(now func() time.Time, jitter, maxTTL time.Duration) Option {
	return OptionFunc(func(s *Schema) (errs error) {
		if now == nil {
			now = time.Now
		}

		if maxTTL < 0 {
			errs = multierr.Append(errs, fmt.Errorf("%w: on positive maxTTL: %s", ErrInvalidCheckUntil, maxTTL))
		}
		if jitter < 0 {
			errs = multierr.Append(errs, fmt.Errorf("%w: non positive jitter: %s", ErrInvalidCheckUntil, jitter))
		}
		if s.Until.IsZero() {
			errs = multierr.Append(errs, fmt.Errorf("%w: zero value Until: %s", ErrInvalidCheckUntil, s.Until))
		}

		limit := (now().Add(maxTTL)).Add(jitter)
		proposed := (s.Until)
		if proposed.After(limit) {
			errs = multierr.Append(errs, fmt.Errorf("%w: CheckUntil value of webhook is out of bounds: %s after %s",
				ErrInvalidCheckUntil, proposed.String(), limit.String()))
		}

		return errs
	})
}

func ReceiverURLValidator(checker *urlegit.Checker) Option {
	return OptionFunc(func(s *Schema) (errs error) {
		url := s.Webhook.ReceiverURL
		if url == "" {
			return fmt.Errorf("%w: is empty", ErrInvalidReceiverURL)
		} else if err := checker.Text(url); err != nil {
			return fmt.Errorf("%w: `%s`: %s", ErrInvalidReceiverURL, url, err)
		}

		return errs
	})
}

// FailureURLValidator is an Option that allows the caller to provide a
// URL validator that is used to validate the FailureURL.
func FailureURLValidator(checker *urlegit.Checker) Option {
	return OptionFunc(func(s *Schema) error {
		if err := checker.Text(s.FailureURL); err != nil {
			return fmt.Errorf("%w: `%s`: %s", ErrInvalidFailureURL, s.FailureURL, err)
		}

		return nil
	})
}

// AlternativeURLValidator is an option that allows the caller to provide
// a URL validator that is used to validate the AlternativeURL.
func AlternativeURLValidator(urlc *urlegit.Checker) Option {
	return OptionFunc(func(s *Schema) (errs error) {
		for _, url := range s.Webhook.AlternativeURLs {
			if err := urlc.Text(url); err != nil {
				errs = multierr.Append(errs, fmt.Errorf("%w: `%s`: %s", ErrAlternativeURLs, url, err))
			}
		}

		return errs
	})
}

func UntilValidator(jitter, maxTTL time.Duration, now func() time.Time) Option {
	return OptionFunc(func(s *Schema) (errs error) {
		if now == nil {
			now = time.Now
		}

		if maxTTL < 0 {
			errs = multierr.Append(errs, fmt.Errorf("%w: non positive maxTTL: %s", ErrInvalidUntil, maxTTL))
		}
		if jitter < 0 {
			errs = multierr.Append(errs, fmt.Errorf("%w: non positive jitter: %s", ErrInvalidUntil, jitter))
		}
		if s.Until.IsZero() {
			errs = multierr.Append(errs, fmt.Errorf("%w: zero value Until: %s", ErrInvalidUntil, s.Until))
		}

		limit := (now().Add(maxTTL)).Add(jitter)
		proposed := (s.Until)
		if proposed.After(limit) {
			errs = multierr.Append(errs,
				fmt.Errorf("%w: Until value of webhook is out of bounds: %s after %s",
					ErrInvalidUntil, proposed.String(), limit.String()))
		}

		return errs
	})
}
