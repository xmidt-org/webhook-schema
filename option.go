// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package wrpeventstream

import (
	"errors"

	"github.com/xmidt-org/webhook-schema/errortrace"
	"go.uber.org/multierr"
)

var (
	ErrOptionFailure      = errors.New("Option(s) failure(s)")
	errOptionFailureTrace = errortrace.New(ErrOptionFailure,
		errortrace.Level(errortrace.ErrorLevel),
		errortrace.Tag("ClientManifest.Options"),
	)
)

type Option interface {
	Apply(ClientManifest) error
}

type Options []Option

func (opts Options) Apply(s ClientManifest) (errs error) {
	for _, o := range opts {
		errs = multierr.Append(errs, o.Apply(s))
	}

	if errs != nil {
		return errOptionFailureTrace.AppendDetail(errs)
	}

	return nil
}

type OptionFunc func(ClientManifest) error

func (f OptionFunc) Apply(s ClientManifest) (errs error) {
	return f(s)
}
