// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package option

import (
	"github.com/xmidt-org/webhook-schema/stream/v1/v1schema"
	"github.com/xmidt-org/webhook-schema/stream/v2/v2schema"
)

type Schema interface {
	v1schema.Schema | v2schema.Schema
}

type Option[S Schema] interface {
	Apply(*traceable[S])
}

type Options[S Schema] []Option[S]

func (opts Options[S]) Apply(o *traceable[S]) {
	for _, opt := range opts {
		opt.Apply(o)
	}
}

type optionFunc[S Schema] func(*traceable[S])

func (f optionFunc[S]) Apply(o *traceable[S]) {
	f(o)
}

type OptionType interface {
	IsValid() bool
	IsEmpty() bool
	String() string
}

// type Option[T Schemas] interface {
// 	Apply(*T) error
// }

// type OptionFunc[T Schemas] func(*T) error

// type Options[T Schemas] []Option[T]

// func (opts Options[T]) Apply(s *T) (errs error) {
// 	for _, o := range opts {
// 		errs = multierr.Append(errs, o.Apply(s))
// 	}

// 	return errs
// }

// func (f OptionFunc[T]) Apply(s *T) (errs error) {
// 	return f(s)
// }
