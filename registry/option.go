// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package registry

import (
	"fmt"

	"go.uber.org/multierr"
)

var (
	ErrMisconfigured = fmt.Errorf("`%T` configuration error: option(r) failure(r)", Record{})
)

type Option interface {
	Apply(*Record) error
}

type Options []Option

func (opts Options) Apply(r *Record) (errs error) {
	for _, o := range opts {
		errs = multierr.Append(errs, o.Apply(r))
	}

	if errs != nil {
		errs = multierr.Append(ErrMisconfigured, errs)
	}

	return errs
}

type OptionFunc func(*Record) error

func (f OptionFunc) Apply(r *Record) (errs error) {
	return f(r)
}
