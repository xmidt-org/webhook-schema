// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package errortrace

import "fmt"

func Level(l LevelType) Option {
	return OptionFunc(func(err *TraceableError) {
		err.level = l
	})
}

func Tag(t any) Option {
	return OptionFunc(func(err *TraceableError) {
		// var subcomponent []string
		// for _, c := range cs {
		// 	switch i := c.(type) {
		// 	case string:
		// 		subcomponent = append(subcomponent, i)
		// 	case interface{ String() string }:
		// 		subcomponent = append(subcomponent, i.String())
		// 	default:
		// 		subcomponent = append(subcomponent, fmt.Sprintf("%T", i))
		// 	}
		// }

		// err.tag = strings.Join(subcomponent, ".")

		switch i := t.(type) {
		case string:
			err.tag = i
		case interface{ String() string }:
			err.tag = i.String()
		default:
			err.tag = fmt.Sprintf("%T", i)
		}
	})
}

func AppendDetail(e error) Option {
	return OptionFunc(func(err *TraceableError) {
		*err = err.AppendDetail(e).(TraceableError)
	})
}
