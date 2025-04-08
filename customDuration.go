// SPDX-FileCopyrightText: 2022 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package stream

import (
	"bytes"
	"strconv"
	"strings"
	"time"
)

type InvalidDurationError struct {
	Value string
}

func (ide *InvalidDurationError) Error() string {
	var o strings.Builder
	o.WriteString("duration must be of type int or string (example:'5m'); Invalid value: ")
	o.WriteString(ide.Value)
	return o.String()
}

// CustomDuration is a custom type for time.Duration that allows for
// unmarshaling from a string or int.  If unmarshaling from a string,
// the string must be parsable by time.ParseDuration.  If unmarshaling
// from an int, the int is assumed to be in seconds. Wes
type CustomDuration time.Duration

func (cd CustomDuration) String() string {
	return time.Duration(cd).String()
}

func (cd CustomDuration) MarshalJSON() ([]byte, error) {
	d := bytes.NewBuffer(nil)
	d.WriteByte('"')
	d.WriteString(cd.String())
	d.WriteByte('"')
	return d.Bytes(), nil
}

func (cd *CustomDuration) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' {
		var d time.Duration
		d, err = time.ParseDuration(string(b[1 : len(b)-1]))
		if err == nil {
			*cd = CustomDuration(d)
			return
		}
	}

	var d int64
	d, err = strconv.ParseInt(string(b), 10, 64)
	if err == nil {
		*cd = CustomDuration(time.Duration(d) * time.Second)
		return
	}

	err = &InvalidDurationError{
		Value: string(b),
	}

	return
}
