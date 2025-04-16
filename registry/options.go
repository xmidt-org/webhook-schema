// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package registry

import (
	"errors"
	"fmt"

	"go.uber.org/multierr"
)

var (
	ErrItemIDEmpty   = errors.New("item ID is required")
	ErrItemDataEmpty = errors.New("data field in item is required")
	ErrOptionFailure = fmt.Errorf("`%T` option failure", Record{})
)

func IDValidator() Option {
	return OptionFunc(func(r *Record) error {
		if len(r.ID) < 1 {
			return multierr.Append(ErrOptionFailure, ErrItemIDEmpty)
		}

		return nil
	})
}

func DataValidator() Option {
	return OptionFunc(func(r *Record) error {
		if len(r.Data) < 1 {
			return multierr.Append(ErrOptionFailure, ErrItemDataEmpty)
		}

		return nil
	})
}
