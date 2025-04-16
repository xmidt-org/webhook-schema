// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package errortrace

type Option interface {
	Apply(*TraceableError)
}

type Options []Option

func (opts Options) Apply(err *TraceableError) {
	for _, o := range opts {
		o.Apply(err)
	}
}

type OptionFunc func(*TraceableError)

func (f OptionFunc) Apply(err *TraceableError) {
	f(err)
}
