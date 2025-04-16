// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package v2schema

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/xmidt-org/urlegit"
	"go.uber.org/multierr"
)

var (
	ErrEmptySchema               = errors.New("empty schema error")
	ErrEmptyWebhooks             = errors.New("invalid Webhooks: no webhooks found")
	ErrExpired                   = errors.New("registration has already expired")
	ErrEventRegexCompilerFailure = errors.New("failed to compile Matcher regexp")
	ErrInvalidReceiverURL        = errors.New("invalid ReceiverURL")
	ErrInvalidFailureURL         = errors.New("invalid FailureURL")
)

func Address(addr string) Option {
	return OptionFunc(func(s *Schema) error {
		s.Address = addr

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
	return OptionFunc(func(s *Schema) (errs error) {
		if s.CanonicalName == "" {
			return ErrEmptySchema
		}

		return nil
	})
}

func OnlyWebhooksValidator() Option {
	return OptionFunc(func(s *Schema) (errs error) {
		if len(s.Webhooks) == 0 {
			return ErrEmptyWebhooks
		}

		return nil
	})
}

func EventRegexValidator() Option {
	return OptionFunc(func(s *Schema) (errs error) {
		for _, m := range s.Matcher {
			if _, err := regexp.Compile(m.Regex); err != nil {
				errs = multierr.Append(errs, fmt.Errorf("%w: `%s`: %s", ErrEventRegexCompilerFailure, m.Regex, err))
			}
		}

		return errs
	})
}

func ExpiresValidator() Option {
	return OptionFunc(func(s *Schema) error {
		if time.Now().After(s.Expires) {
			return ErrExpired
		}

		return nil
	})
}

func ReceiverURLValidator(checker *urlegit.Checker) Option {
	return OptionFunc(func(s *Schema) (errs error) {
		for i, w := range s.Webhooks {
			for _, url := range w.ReceiverURLs {
				if url == "" {
					errs = multierr.Append(errs,
						fmt.Errorf("%w: webhook %v is empty", ErrInvalidReceiverURL, i))
				} else if err := checker.Text(url); err != nil {
					errs = multierr.Append(errs, fmt.Errorf("%w: webhook %v `%s`:%s", ErrInvalidReceiverURL, i, url, err))
				}
			}
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
