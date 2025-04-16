// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package v2

import (
	"github.com/xmidt-org/webhook-schema/stream/v2/v2schema"
)

func AddDefaults(opts ...v2schema.Option) Option {
	return OptionFunc(func(m *manifest) error {
		m.defaults = append(m.defaults, opts...)

		return nil
	})
}

func AddValidators(opts ...v2schema.Option) Option {
	return OptionFunc(func(m *manifest) error {
		m.validators = append(m.validators, opts...)

		return nil
	})
}

func Stream(s v2schema.Schema) Option {
	return OptionFunc(func(m *manifest) error {
		return m.SetStream(s)
	})
}
