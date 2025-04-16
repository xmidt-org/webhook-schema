// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package v1

import (
	"github.com/xmidt-org/webhook-schema/stream/v1/v1schema"
)

func AddDefaults(opts ...v1schema.Option) Option {
	return OptionFunc(func(m *manifest) error {
		m.defaults = append(m.defaults, opts...)

		return nil
	})
}

func AddValidators(opts ...v1schema.Option) Option {
	return OptionFunc(func(m *manifest) error {
		m.validators = append(m.validators, opts...)

		return nil
	})
}

func Stream(s v1schema.Schema) Option {
	return OptionFunc(func(m *manifest) error {
		return m.SetStream(s)
	})
}
