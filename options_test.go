// SPDX-FileCopyrightText: 2022 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package webhook

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var mockVal *MockValidator

func TestError(t *testing.T) {
	optError := fmt.Errorf("error")
	opt := Error(optError)
	assert.NotNil(t, opt)

	err := opt.Validate(mockVal)
	assert.NotNil(t, err)
	assert.Equal(t, optError, err)

	s := opt.String()
	assert.NotNil(t, s)
}
