// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package option

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestURLChecker(t *testing.T) {
	tcs := []struct {
		desc              string
		config            URLChecker
		expectedErr       error
		expectedFuncCount int
	}{
		{
			desc: "HTTPSOnly only",
			config: URLChecker{
				URL: URLConfig{
					AllowLoopback: true,
					Schemes:       []string{"https"},
				},
			},
			expectedFuncCount: 1,
		},
		{
			desc: "AllowLoopback only",
			config: URLChecker{
				URL: URLConfig{
					AllowLoopback: false,
					Schemes:       []string{"https", "http"},
				},
			},
			expectedFuncCount: 2,
		},
		{
			desc: "AllowIp Only",
			config: URLChecker{
				IP: IPConfig{
					Allow: false,
				},
			},
			expectedFuncCount: 2,
		},
		{
			desc: "AllowSpecialUseHosts Only",
			config: URLChecker{
				Domain: DomainConfig{
					AllowSpecialUseDomains: false,
				},
			},
			expectedFuncCount: 2,
		},
		{
			desc: "AllowSpecialuseIPS Only",
			config: URLChecker{
				IP: IPConfig{
					Allow: true,
				},
			},
			expectedFuncCount: 2,
		},
		{
			desc: "Forbidden Subnets",
			config: URLChecker{
				IP: IPConfig{
					Allow:            false,
					ForbiddenSubnets: []string{"10.0.0.0/8"},
				},
			},
			expectedFuncCount: 1,
		},
		{
			desc: "Forbidden Domains",
			config: URLChecker{
				Domain: DomainConfig{
					AllowSpecialUseDomains: true,
					ForbiddenDomains:       []string{"example.com."},
				},
			},
		},
		{
			desc: "Build None",
			config: URLChecker{
				URL: URLConfig{
					Schemes:       []string{"https", "http"},
					AllowLoopback: true,
				},
			},
			expectedFuncCount: 1,
		},
		{
			desc: "Build All",
			config: URLChecker{
				URL: URLConfig{
					Schemes:       []string{"https"},
					AllowLoopback: false,
				},
			},
			expectedFuncCount: 5,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			assert := assert.New(t)
			vals, err := tc.config.Build()
			if tc.expectedErr != nil {
				assert.ErrorIs(err, tc.expectedErr,
					fmt.Errorf("error [%v] doesn't contain error [%v] in its err chain",
						err, tc.expectedErr))
				assert.Nil(vals)
				return
			}
			require.NoError(t, err)
		})
	}
}
