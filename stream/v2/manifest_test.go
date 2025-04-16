// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package v2

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmidt-org/webhook-schema/errortrace"
	"github.com/xmidt-org/webhook-schema/stream/option"
	"github.com/xmidt-org/webhook-schema/stream/sink/kafka"
	"github.com/xmidt-org/webhook-schema/stream/sink/webhook"
	"github.com/xmidt-org/webhook-schema/stream/v2/v2schema"
)

func TestManifestUnmarshalJSON(t *testing.T) {
	successStream := v2schema.Schema{
		CanonicalName: "com.examples.webhook-tests",
		Matcher:       []v2schema.FieldRegex{{Field: "canonical_name", Regex: "webpa"}},
		Webhooks:      []webhook.V2Schema{{ReceiverURLs: []string{"https://example.com", "https://example.com"}}},
	}
	successJsonStream, err := json.Marshal(successStream)
	require.NoError(t, err)
	require.NotEmpty(t, successJsonStream)

	tests := []struct {
		description string
		empty       bool
		jsonStream  []byte
		stream      v2schema.Schema
		expectedErr error
	}{
		{
			description: "empty stream json object failure",
			jsonStream:  []byte(`{}`),
			expectedErr: ErrUnmarshal,
		},
		{
			description: "stream json success",
			jsonStream:  successJsonStream,
			stream:      successStream,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			stream, err := New()
			assert.NoError(err)
			require.NotNil(stream)

			err = json.Unmarshal(tc.jsonStream, &stream)
			if tc.expectedErr == nil {
				assert.NoError(err)
				assert.Equal(tc.stream, stream.GetStream())

				return
			}

			assert.ErrorIs(err, tc.expectedErr)
			assert.Empty(stream.GetStream())
		})
	}
}

func TestManifestUnmarshalWithOptions(t *testing.T) {
	successStream := v2schema.Schema{
		CanonicalName: "com.examples.webhook-tests",
		Matcher:       []v2schema.FieldRegex{{Field: "canonical_name", Regex: "webpa"}},
		Webhooks:      []webhook.V2Schema{{ReceiverURLs: []string{"https://example.com", "https://example.com"}}},
	}
	successJsonStream, err := json.Marshal(successStream)
	require.NoError(t, err)
	require.NotEmpty(t, successJsonStream)

	failureStream := v2schema.Schema{
		CanonicalName: "com.examples.webhook-tests",
		Matcher:       []v2schema.FieldRegex{{Regex: "("}},
		Webhooks:      []webhook.V2Schema{{ReceiverURLs: []string{"https://example.com", "http://example.com"}}},
	}
	failureJsonStream, err := json.Marshal(failureStream)
	require.NoError(t, err)
	require.NotEmpty(t, failureJsonStream)

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
		description  string
		options      Options
		jsonStream   []byte
		expected     manifest
		expectedErrs []error
	}{
		{
			description: "option success",
			options: Options{
				AddValidators(
					option.New(v2schema.ReceiverURLValidator(checker),
						option.Level[v2schema.Schema](errortrace.WarningLevel),
						option.Type[v2schema.Schema](v2schema.ReceiverURLValidatorType))),
			},
			jsonStream: successJsonStream,
			expected: manifest{
				stream: successStream,
				validators: append(defaultValidators,
					option.New(v2schema.ReceiverURLValidator(checker),
						option.Level[v2schema.Schema](errortrace.WarningLevel),
						option.Type[v2schema.Schema](v2schema.ReceiverURLValidatorType))),
			},
		},
		{
			description: "option failure",
			options: Options{
				AddValidators(
					option.New(v2schema.ReceiverURLValidator(checker),
						option.Level[v2schema.Schema](errortrace.WarningLevel),
						option.Type[v2schema.Schema](v2schema.ReceiverURLValidatorType))),
			},
			jsonStream: failureJsonStream,
			expected: manifest{
				validators: append(defaultValidators,
					option.New(v2schema.ReceiverURLValidator(checker),
						option.Level[v2schema.Schema](errortrace.WarningLevel),
						option.Type[v2schema.Schema](v2schema.ReceiverURLValidatorType))),
			},
			expectedErrs: []error{
				errUnmarshalTrace,
				errValidatorTrace,
				option.ErrOptionFailure,
				v2schema.ErrInvalidReceiverURL,
			},
		},
		{
			description: "multiple options success",
			options: Options{
				AddValidators(
					option.New(v2schema.ReceiverURLValidator(checker),
						option.Level[v2schema.Schema](errortrace.WarningLevel),
						option.Type[v2schema.Schema](v2schema.ReceiverURLValidatorType)),
					option.New(v2schema.EventRegexValidator(),
						option.Level[v2schema.Schema](errortrace.ErrorLevel),
						option.Type[v2schema.Schema](v2schema.EventRegexValidatorType))),
			},
			jsonStream: successJsonStream,
			expected: manifest{
				stream: successStream,
				validators: append(defaultValidators,
					option.New(v2schema.ReceiverURLValidator(checker),
						option.Level[v2schema.Schema](errortrace.WarningLevel),
						option.Type[v2schema.Schema](v2schema.ReceiverURLValidatorType)),
					option.New(v2schema.EventRegexValidator(),
						option.Level[v2schema.Schema](errortrace.ErrorLevel),
						option.Type[v2schema.Schema](v2schema.EventRegexValidatorType))),
			},
		},
		{
			description: "multiple options failure",
			options: Options{
				AddValidators(
					option.New(v2schema.ReceiverURLValidator(checker),
						option.Level[v2schema.Schema](errortrace.WarningLevel),
						option.Type[v2schema.Schema](v2schema.ReceiverURLValidatorType)),
					option.New(v2schema.EventRegexValidator(),
						option.Level[v2schema.Schema](errortrace.ErrorLevel),
						option.Type[v2schema.Schema](v2schema.EventRegexValidatorType))),
			},
			jsonStream: failureJsonStream,
			expected: manifest{
				validators: append(defaultValidators,
					option.New(v2schema.ReceiverURLValidator(checker),
						option.Level[v2schema.Schema](errortrace.WarningLevel),
						option.Type[v2schema.Schema](v2schema.ReceiverURLValidatorType)),
					option.New(v2schema.EventRegexValidator(),
						option.Level[v2schema.Schema](errortrace.ErrorLevel),
						option.Type[v2schema.Schema](v2schema.EventRegexValidatorType))),
			},
			expectedErrs: []error{
				errUnmarshalTrace,
				errValidatorTrace,
				option.ErrOptionFailure,
				v2schema.ErrInvalidReceiverURL,
				v2schema.ErrEventRegexCompilerFailure,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)
			stream, err := New(tc.options...)
			require.NoError(err)
			require.NotNil(stream)
			assert.Empty(stream.GetStream())

			smanifest, ok := stream.(*manifest)
			require.True(ok)
			require.Len(smanifest.validators, len(tc.expected.validators))
			require.Len(smanifest.defaults, len(tc.expected.defaults))

			for i, e := range tc.expected.validators {
				v := smanifest.validators[i]
				require.IsType(e, v)

				// nolint:staticcheck
				te, ok := e.(option.TracableOption[v2schema.Schema])
				require.True(ok)
				// nolint:staticcheck
				tv, ok := v.(option.TracableOption[v2schema.Schema])
				require.True(ok)

				assert.Equal(te.Level(), tv.Level())
				assert.Equal(te.Type(), tv.Type())
			}

			err = json.Unmarshal(tc.jsonStream, &stream)
			if tc.expectedErrs == nil {
				assert.NoError(err)
				assert.NotEmpty(stream.GetStream())

				return
			}

			assert.Error(err)
			// manifest is designed to `zero out` during an unmarshal related failure
			assert.Empty(stream.GetStream())
			for _, e := range tc.expectedErrs {
				assert.ErrorIs(err, e)
			}
		})
	}
}

func TestManifestObfuscateSecrets(t *testing.T) {
	tests := []struct {
		description string
		stream      manifest
		expected    manifest
	}{
		{
			description: "obfuscate webhook secret success",
			stream: manifest{
				stream: v2schema.Schema{
					CanonicalName: "com.examples.webhook-tests",
					Webhooks:      []webhook.V2Schema{{Secret: "FooBar"}, {Secret: "FooBar"}},
				},
			},
			expected: manifest{
				stream: v2schema.Schema{
					CanonicalName: "com.examples.webhook-tests",
					Webhooks:      []webhook.V2Schema{{Secret: "<obfuscated>"}, {Secret: "<obfuscated>"}},
				},
			},
		},
		{
			description: "obfuscate kafka secret success",
			stream: manifest{
				stream: v2schema.Schema{
					CanonicalName: "com.examples.webhook-tests",
					Kafkas:        []kafka.V1Schema{{KafkaProducer: kafka.V1KafkaProducer{Secret: "FooBar"}}, {KafkaProducer: kafka.V1KafkaProducer{Secret: "FooBar"}}},
				},
			},
			expected: manifest{
				stream: v2schema.Schema{
					CanonicalName: "com.examples.webhook-tests",
					Kafkas:        []kafka.V1Schema{{KafkaProducer: kafka.V1KafkaProducer{Secret: "<obfuscated>"}}, {KafkaProducer: kafka.V1KafkaProducer{Secret: "<obfuscated>"}}},
				},
			},
		},
		{
			description: "obfuscate kafka & webhook secret success",
			stream: manifest{
				stream: v2schema.Schema{
					CanonicalName: "com.examples.webhook-tests",
					Webhooks:      []webhook.V2Schema{{Secret: "FooBar"}, {Secret: "FooBar"}},
					Kafkas:        []kafka.V1Schema{{KafkaProducer: kafka.V1KafkaProducer{Secret: "FooBar"}}, {KafkaProducer: kafka.V1KafkaProducer{Secret: "FooBar"}}},
				},
			},
			expected: manifest{
				stream: v2schema.Schema{
					CanonicalName: "com.examples.webhook-tests",
					Webhooks:      []webhook.V2Schema{{Secret: "<obfuscated>"}, {Secret: "<obfuscated>"}},
					Kafkas:        []kafka.V1Schema{{KafkaProducer: kafka.V1KafkaProducer{Secret: "<obfuscated>"}}, {KafkaProducer: kafka.V1KafkaProducer{Secret: "<obfuscated>"}}},
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			assert := assert.New(t)
			assert.NotEqual(tc.expected.GetStream(), tc.stream.GetStream())
			tc.stream.ObfuscateSecrets()
			assert.Equal(tc.expected.GetStream(), tc.stream.GetStream())
		})
	}
}
