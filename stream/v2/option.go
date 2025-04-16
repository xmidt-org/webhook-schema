// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package v2

import (
	"github.com/xmidt-org/webhook-schema/stream/option"
	"github.com/xmidt-org/webhook-schema/stream/v2/v2schema"
	"go.uber.org/multierr"
)

type (
	OptionConfig = option.OptionConfig[v2schema.OptionType]
)

type Option interface {
	Apply(*manifest) error
}

type Options []Option

func (opts Options) Apply(m *manifest) (errs error) {
	for _, o := range opts {
		errs = multierr.Append(errs, o.Apply(m))
	}

	return errs
}

type OptionFunc func(*manifest) error

func (f OptionFunc) Apply(m *manifest) (errs error) {
	return f(m)
}
