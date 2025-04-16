// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package errortrace

import (
	"errors"
	"fmt"

	"go.uber.org/multierr"
)

type Trace interface {
	error

	AppendDetail(...error) Trace
	TraceError(int, string) string
	DetailErrs() []error
	Err() error
	Level() LevelType
	Tag() string
	Unwrap() []error
	Is(error) bool
}

// TraceError returns the trace of the input errors, showing the source of each error and their relation.
// The following is an example of a trace:
// [START OF ERROR TRACE]
// |-err0->[err="request `/api/examples` failed", level=error, tag="request"]
// |   |-err0.0->[err="failed to extract metadata from request"]
// |   |-err0.1->[err="request `/api/examples` failed", level=error, tag="processData"]
// |   |   |-err0.1.0->[err="request `/api/examples` failed", level=error, tag="transformerA"]
// |   |   |   |-err0.1.0.0->[err="error 1"]
// |   |   |   |-err0.1.0.1->[err="error 2"]
// |   |   |   |-err0.1.0.2->[err="error ..."]
// |   |   |   |-err0.1.0.3->[err="error 1020"]
// |   |   |-err0.1.1->[err="request `/api/examples` failed", level=info, tag="transformerB"]
// |   |   |   |-err0.1.1.0->[err="some randome error"]
// |   |   |-err0.1.2->[err="request `/api/examples` failed", level=info, tag="transformerB"]
// |   |   |   |-err0.1.2.0->[err="some randome error"]
// |-err.1->[err="response encoding failed"]
// |
// |
// [END OF ERROR TRACE]
func TraceError(err error) string {
	var errs []error
	switch err.(type) {
	case Trace:
		errs = append(errs, err)
	default:
		errs = multierr.Errors(err)
	}

	if len(errs) == 0 {
		return "[NO ERRORS TO TRACE]"
	}

	trace := "[START OF ERROR TRACE]"
	for i, e := range errs {
		s := fmt.Sprintf("\n|-%s.%d->[err=\"%s\"]", "err", i, e)
		if de, ok := e.(Trace); ok {
			s = de.TraceError(0, fmt.Sprintf("%s%d", "err", i))
		}

		trace += s
	}

	return trace + "\n|\n|\n[END OF ERROR TRACE]\n"
}

func New(err error, opts ...Option) Trace {
	oe := TraceableError{level: ErrorLevel, err: err}
	Options(opts).Apply(&oe)

	return oe
}

type TraceableError struct {
	// level of the error.
	// This give clients the flexibility for configuration, logging and etc.
	level LevelType
	// TODO: update tag to tags (map of strings)
	tag string
	// err is the main error.
	err error
	// detailErrs are the source error(s) leading, if not the start, to the trace.
	detailErrs []error

	unwrapped bool
}

func (oe TraceableError) AppendDetail(errs ...error) Trace {
	for errs != nil {
		var recheck []error
		for _, err := range errs {
			switch err.(type) {
			case Trace:
				oe.detailErrs = append(oe.detailErrs, err)
			// multierr.Appended-like errors are flattened and rechecked.
			case interface {
				error
				Unwrap() []error
			}:
				recheck = append(recheck, multierr.Errors(err)...)
			// Normal, combined or joined errors are treated as a single error and not unwrapped.
			default:
				oe.detailErrs = append(oe.detailErrs, err)
			}
			errs = recheck
		}
	}

	return oe
}

func (oe TraceableError) DetailErrs() []error {
	return oe.detailErrs
}

func (oe TraceableError) Err() error {
	return oe.err
}

func (oe TraceableError) Level() LevelType {
	return oe.level
}

func (oe TraceableError) Tag() string {
	return oe.tag
}

func (oe TraceableError) TraceError(layer int, id string) string {
	prefix := "\n|"
	for range layer {
		prefix += fmt.Sprintf("%*s", 4, "|")
	}

	trace := fmt.Sprintf("%s-%s->%s", prefix, id, oe)
	prefix += fmt.Sprintf("%*s", 4, "|")
	for i, e := range oe.detailErrs {
		s := fmt.Sprintf("%s-%s.%d->[err=\"%s\"]", prefix, id, i, e)
		if de, ok := e.(Trace); ok {
			s = de.TraceError(layer+1, fmt.Sprintf("%s.%d", id, i))
		}

		trace += s
	}

	return trace
}

// Implements error/multierr interface

func (oe TraceableError) Error() string {
	return fmt.Sprintf("[err=\"%s\", level=%s, tag=\"%s\"]", oe.err, oe.level, oe.tag)
}

func (oe TraceableError) Unwrap() []error {
	// Check whether or not oe has already been unwrapped once.
	// If it has, then unwrap again by removing the remaining metadata and return the error as is.
	if oe.unwrapped {
		return []error{oe.err}
	}

	// Otherwise, unwrap by uncoupling oe.err and oe.detailErrs while retaining oe's metadata for oe.err.
	oe.unwrapped = true
	// Uncouple err and detailErrs.
	errs := oe.detailErrs
	oe.detailErrs = nil
	unwrappedErrs := []error{oe}
	for errs != nil {
		var recheck []error
		for _, err := range errs {
			switch e := err.(type) {
			// Unpack all Trace errors.
			case Trace:
				unwrappedErrs = append(unwrappedErrs, e.Unwrap()...)
			// multierr.Appended-like errors are flattened and rechecked.
			case interface {
				error
				Unwrap() []error
			}:
				recheck = append(recheck, multierr.Errors(e)...)
			// Normal, combined or joined errors are treated as a single error and not unwrapped.
			default:
				unwrappedErrs = append(unwrappedErrs, e)
			}
			errs = recheck
		}
	}

	return unwrappedErrs
}

func (oe TraceableError) Is(target error) bool {
	if de, ok := target.(Trace); ok {
		target = de.Err()
	}

	// Check against all errors to determine whether target has a match.
	for _, e := range multierr.Errors(oe) {
		if de, ok := e.(Trace); ok {
			e = de.Err()
		}

		if errors.Is(e, target) {
			return true
		}
	}

	return false
}
