// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package option

import (
	"errors"
	"fmt"

	errortrace "github.com/xmidt-org/webhook-schema/errortrace"
	"go.uber.org/multierr"
)

var (
	ErrOptionFailure = errors.New("option failure")
)

type TracableOption[S Schema] interface {
	schemaOption[S]
	Level() errortrace.LevelType
	Type() OptionType
}

type schemaOption[S Schema] interface{ Apply(*S) error }

func New[S Schema](opt schemaOption[S], opts ...Option[S]) TracableOption[S] {
	o := traceable[S]{option: opt}
	Options[S](opts).Apply(&o)

	return o
}

type traceable[S Schema] struct {
	option schemaOption[S]
	level  errortrace.LevelType
	otype  OptionType
}

func (o traceable[S]) Apply(s *S) (errs error) {
	if err := o.option.Apply(s); err != nil {
		return errortrace.New(multierr.Append(ErrOptionFailure, fmt.Errorf("invaild `%T` stream, %s validator error", s, o.otype)),
			errortrace.Level(o.level),
			errortrace.Tag(o.otype),
			errortrace.AppendDetail(err),
		)
	}

	return nil
}

func (o traceable[S]) Level() errortrace.LevelType {
	return o.level
}

func (o traceable[S]) Type() OptionType {
	return o.otype
}
