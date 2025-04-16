// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package wrpeventstream

import (
	"errors"

	"github.com/xmidt-org/webhook-schema/errortrace"
	"go.uber.org/multierr"
)

var (
	ErrMisconfiguredBuilder = errors.New("builder configuration error: option(s) failure(s)")
)

type BuilderOption interface {
	Apply(*Builder) error
}

type BuilderOptions []BuilderOption

func (opts BuilderOptions) Apply(b *Builder) (errs error) {
	for _, o := range opts {
		errs = multierr.Append(errs, o.Apply(b))
	}

	if errs != nil {
		errs = errortrace.New(ErrMisconfiguredBuilder,
			errortrace.Level(errortrace.ErrorLevel),
			errortrace.AppendDetail(errs),
			errortrace.Tag("Options"),
		)
	}

	return errs
}

type BuilderOptionFunc func(*Builder) error

func (f BuilderOptionFunc) Apply(b *Builder) (errs error) {
	return f(b)
}
