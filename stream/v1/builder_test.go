// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package v1

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmidt-org/urlegit"
	"github.com/xmidt-org/webhook-schema/errortrace"
	"github.com/xmidt-org/webhook-schema/stream/option"
	"github.com/xmidt-org/webhook-schema/stream/v1/v1schema"
)

func ExampleBuilder() {
	jsonStreamBuilder := []byte(`{
	"disable": false,
	"options": [
		{
			"type": "receiver_url",
			"level": "warning",
			"url_checker": {
				"domain": {
					"allow_special_use_domains": true
				},
				"url": {
					"schemes": ["https"]
				}
			}
		},
		{
			"type": "event_regex",
			"level": "error"
		},
		{
			"type": "until",
			"level": "error",
			"disable": true
		}
	]
}`)

	var builder Builder
	if err := json.Unmarshal(jsonStreamBuilder, &builder); err != nil {
		panic(err)
	}

	stream, err := builder.Build()
	if err != nil {
		panic(err)
	}

	jsonStream := []byte(`{
		"registered_from_address": "",
		"config": {
			"url": "https://example.com",
			"content_type": ""
		},
		"failure_url": "",
		"events": [
			"event.*"
		],
		"matcher": {
			"device_id": null
		},
		"duration": "0s",
		"until": "0001-01-01T00:00:00Z"
	}`)
	err = json.Unmarshal(jsonStream, &stream)
	if err != nil {
		panic(err)
	}

	jsonStream, _ = json.MarshalIndent(stream, "", "   ")
	fmt.Printf("Valid v1 stream example\nstream:\n%s\n\n", jsonStream)

	// Output: Valid v1 stream example
	// stream:
	// {
	//    "registered_from_address": "",
	//    "config": {
	//       "url": "https://example.com",
	//       "content_type": ""
	//    },
	//    "failure_url": "",
	//    "events": [
	//       "event.*"
	//    ],
	//    "matcher": {
	//       "device_id": null
	//    },
	//    "duration": "0s",
	//    "until": "0001-01-01T00:00:00Z"
	// }
}

func TestBuilderUnmarshalJSON(t *testing.T) {
	tests := []struct {
		description string
		jsonConfig  []byte
		expectedErr error
	}{
		{
			description: "empty stream json object failure",
			jsonConfig:  []byte(`{}`),
			expectedErr: ErrUnmarshal,
		},
		{
			description: "stream json success",
			jsonConfig: []byte(`{
				"options": [
					{
						"type": "receiver_url",
						"level": "warning",
						"url_checker": {
							"domain": {
								"allow_special_use_domains": true
							},
							"url": {
								"schemes": ["https"]
							}
						}
					}
				]
			}`),
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			var builder Builder
			assert := assert.New(t)

			err := json.Unmarshal(tc.jsonConfig, &builder)
			if tc.expectedErr == nil {
				assert.NoError(err)

				return
			}

		})
	}
}

func TestBuilderUnmarshalWithOptions(t *testing.T) {
	tests := []struct {
		description string
		jsonConfig  []byte
		disabled    bool
		disabledOpt bool
		jsonStream  []byte
		stream      v1schema.Schema
		expected    Builder
	}{
		{
			description: "disabled builder success",
			jsonConfig: []byte(`{
				"disable": true,
				"options": [
					{
						"type": "receiver_url",
						"level": "warning",
						"url_checker": {
							"domain": {
								"allow_special_use_domains": true
							},
							"url": {
								"schemes": ["https"]
							}
						}
					},
					{
						"type": "event_regex",
						"level": "error"
					}
				]
			}`),
			expected: Builder{
				validators: v1schema.Options{
					option.New(v1schema.ReceiverURLValidator(&urlegit.Checker{}),
						option.Level[v1schema.Schema](errortrace.WarningLevel),
						option.Type[v1schema.Schema](v1schema.ReceiverURLValidatorType)),
					option.New(v1schema.EventRegexValidator(),
						option.Level[v1schema.Schema](errortrace.ErrorLevel),
						option.Type[v1schema.Schema](v1schema.EventRegexValidatorType)),
				},
			},
			disabled: true,
		},
		{
			description: "disabled option success",
			jsonConfig: []byte(`{
				"options": [
					{
						"type": "receiver_url",
						"level": "warning",
						"url_checker": {
							"domain": {
								"allow_special_use_domains": true
							},
							"url": {
								"schemes": ["https"]
							}
						},
						"disable": true
					}
				]
			}`),
			disabledOpt: true,
		},
		{
			description: "option success",
			jsonConfig: []byte(`{
				"options": [
					{
						"type": "receiver_url",
						"level": "warning",
						"url_checker": {
							"domain": {
								"allow_special_use_domains": true
							},
							"url": {
								"schemes": ["https"]
							}
						}
					}
				]
			}`),
			expected: Builder{
				validators: v1schema.Options{
					option.New(v1schema.ReceiverURLValidator(&urlegit.Checker{}),
						option.Level[v1schema.Schema](errortrace.WarningLevel),
						option.Type[v1schema.Schema](v1schema.ReceiverURLValidatorType))},
			},
		},
		{
			description: "multiple options success",
			jsonConfig: []byte(`{
				"options": [
					{
						"type": "receiver_url",
						"level": "warning",
						"url_checker": {
							"domain": {
								"allow_special_use_domains": true
							},
							"url": {
								"schemes": ["https"]
							}
						}
					},
					{
						"type": "event_regex",
						"level": "error"
					}
				]
			}`),
			expected: Builder{
				validators: v1schema.Options{
					option.New(v1schema.ReceiverURLValidator(&urlegit.Checker{}),
						option.Level[v1schema.Schema](errortrace.WarningLevel),
						option.Type[v1schema.Schema](v1schema.ReceiverURLValidatorType)),
					option.New(v1schema.EventRegexValidator(),
						option.Level[v1schema.Schema](errortrace.ErrorLevel),
						option.Type[v1schema.Schema](v1schema.EventRegexValidatorType)),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			var builder Builder
			assert := assert.New(t)
			require := require.New(t)
			require.NoError(json.Unmarshal(tc.jsonConfig, &builder))
			assert.Equal(tc.disabled, builder.Disable)

			if tc.disabledOpt {
				require.Empty(builder.validators)
			}

			require.Len(builder.validators, len(tc.expected.validators))
			require.Len(builder.defaults, len(tc.expected.defaults))

			for i, e := range tc.expected.validators {
				v := builder.validators[i]
				require.IsType(e, v)

				// nolint:staticcheck
				te, ok := e.(option.TracableOption[v1schema.Schema])
				require.True(ok)
				// nolint:staticcheck
				tv, ok := v.(option.TracableOption[v1schema.Schema])
				require.True(ok)

				assert.Equal(te.Level(), tv.Level())
				assert.Equal(te.Type(), tv.Type())
			}
		})
	}
}

func TestBuilderBuild(t *testing.T) {
	urlc := option.URLChecker{
		Domain: option.DomainConfig{
			AllowSpecialUseDomains: true,
		},
		URL: option.URLConfig{
			Schemes: []string{"https"},
		},
	}
	checker, err := urlc.Build()
	require.NoError(t, err)
	require.NotNil(t, checker)

	tests := []struct {
		description string
		builder     Builder
		disabled    bool
		expected    *manifest
	}{
		{
			description: "disabled builder failure",
			builder: Builder{
				Config: Config{
					Disable: true,
					Options: []OptionConfig{
						{
							Type:       v1schema.ReceiverURLValidatorType,
							Level:      errortrace.WarningLevel,
							URLChecker: urlc,
						},
						{
							Type:  v1schema.EventRegexValidatorType,
							Level: errortrace.ErrorLevel,
						},
					},
				},
			},
			disabled: true,
		},
		{
			description: "multiple options success",
			builder: Builder{
				Config: Config{
					Options: []OptionConfig{
						{
							Type:       v1schema.ReceiverURLValidatorType,
							Level:      errortrace.WarningLevel,
							URLChecker: urlc,
						},
						{
							Type:  v1schema.EventRegexValidatorType,
							Level: errortrace.ErrorLevel,
						},
					},
				},
			},
			expected: expectedManifest(
				AddValidators(
					option.New(v1schema.ReceiverURLValidator(checker),
						option.Level[v1schema.Schema](errortrace.WarningLevel),
						option.Type[v1schema.Schema](v1schema.ReceiverURLValidatorType)),
					option.New(v1schema.EventRegexValidator(),
						option.Level[v1schema.Schema](errortrace.ErrorLevel),
						option.Type[v1schema.Schema](v1schema.EventRegexValidatorType)),
				),
			),
		},
		{
			description: "multiple options, empty webhook json Failure",
			builder: Builder{
				Config: Config{
					Options: []OptionConfig{
						{
							Type:  v1schema.EventRegexValidatorType,
							Level: errortrace.ErrorLevel,
						},
					},
				},
			},
			expected: expectedManifest(
				AddValidators(
					option.New(v1schema.EventRegexValidator(),
						option.Level[v1schema.Schema](errortrace.ErrorLevel),
						option.Type[v1schema.Schema](v1schema.EventRegexValidatorType)),
				),
			),
		},
	}
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)
			stream, err := tc.builder.Build()
			if tc.disabled {
				assert.ErrorIs(err, ErrBuildDisabled)
				assert.Nil(stream)

				return
			}

			assert.NoError(err)
			assert.NotNil(stream)
			assert.Empty(stream.GetStream())

			smanifest, ok := stream.(*manifest)
			require.True(ok)
			require.Len(smanifest.validators, len(tc.expected.validators))
			require.Len(smanifest.defaults, len(tc.expected.defaults))

			for i, e := range tc.expected.validators {
				v := smanifest.validators[i]
				require.IsType(e, v)

				// nolint:staticcheck
				te, ok := e.(option.TracableOption[v1schema.Schema])
				require.True(ok)
				// nolint:staticcheck
				tv, ok := v.(option.TracableOption[v1schema.Schema])
				require.True(ok)

				assert.Equal(te.Level(), tv.Level())
				assert.Equal(te.Type(), tv.Type())
			}
		})
	}
}

func expectedManifest(opts ...Option) *manifest {
	m, _ := New(opts...)

	return m.(*manifest)
}
