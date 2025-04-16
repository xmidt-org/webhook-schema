// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package option

import "github.com/xmidt-org/webhook-schema/errortrace"

func Level[S Schema](l errortrace.LevelType) Option[S] {
	return optionFunc[S](func(o *traceable[S]) {
		o.level = l
	})
}

func Type[S Schema](t OptionType) Option[S] {
	return optionFunc[S](func(o *traceable[S]) {
		o.otype = t
	})
}
