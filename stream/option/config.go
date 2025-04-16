// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package option

import errortrace "github.com/xmidt-org/webhook-schema/errortrace"

type OptionConfig[T OptionType] struct {
	// Level denotes the option's Level and should be used for client side option error handling.
	Level errortrace.LevelType `json:"level"`
	// Type assigns the option type.
	Type T `json:"type"`
	// disable determines whether the option is active (`diable` is `false`)
	// or inactive (`disable` is `true`).
	// Default is `false`.
	Disable bool `json:"disable"`

	// URLChecker is a noop for nonapplicable OptionTypes.
	URLChecker URLChecker `json:"url_checker"`
}

// IsValid returns true if the wrapped validator and its metadata are valid,
// otherwise false is returned.
func (co OptionConfig[T]) IsValid() bool {
	return co.Type.IsValid() && co.Level.IsValid()
}

// Empty returns true if the wrapped validator is nil or its metadata are consider empty,
// otherwise false is returned.
func (co OptionConfig[T]) IsEmpty() bool {
	return co.Type.IsEmpty() || co.Level.IsEmpty()
}
