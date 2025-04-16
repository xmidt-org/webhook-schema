// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package v2schema

import (
	"go.uber.org/multierr"
)

type Option interface {
	Apply(*Schema) error
}

type Options []Option

func (opts Options) Apply(s *Schema) (errs error) {
	for _, o := range opts {
		errs = multierr.Append(errs, o.Apply(s))
	}

	return errs
}

type OptionFunc func(*Schema) error

func (f OptionFunc) Apply(s *Schema) (errs error) {
	return f(s)
}
