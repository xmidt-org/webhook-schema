// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package v2

import (
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"math"
	"time"

	"github.com/xmidt-org/webhook-schema/errortrace"
	"github.com/xmidt-org/webhook-schema/stream"
	"github.com/xmidt-org/webhook-schema/stream/option"
	"github.com/xmidt-org/webhook-schema/stream/sink"
	"github.com/xmidt-org/webhook-schema/stream/sink/kafka"
	"github.com/xmidt-org/webhook-schema/stream/sink/webhook"
	"github.com/xmidt-org/webhook-schema/stream/v2/v2schema"
)

var (
	ErrMarshal      = fmt.Errorf("`%T` marshal error", v2schema.Schema{})
	errMarshalTrace = errortrace.New(ErrMarshal,
		errortrace.Level(errortrace.ErrorLevel),
		errortrace.Tag(fmt.Sprintf("%T.MarshalJSON", manifest{})))
	ErrUnmarshal      = fmt.Errorf("`%T` unmarshal error", v2schema.Schema{})
	errUnmarshalTrace = errortrace.New(ErrUnmarshal,
		errortrace.Level(errortrace.ErrorLevel),
		errortrace.Tag(fmt.Sprintf("%T.UnmarshalJSON", manifest{})))
	ErrValidate       = fmt.Errorf("`%T` failed validation", v2schema.Schema{})
	errValidatorTrace = errortrace.New(ErrValidate,
		errortrace.Level(errortrace.ErrorLevel),
		errortrace.Tag(fmt.Sprintf("%T.Validate", manifest{})))
	ErrSetDefaults      = fmt.Errorf("`%T` set default failure", v2schema.Schema{})
	errSetDefaultsTrace = errortrace.New(ErrSetDefaults,
		errortrace.Level(errortrace.ErrorLevel),
		errortrace.Tag(fmt.Sprintf("%T.SetDefaults", manifest{})))
)

var (
	defaultValidators = v2schema.Options{
		// defaults validators
		option.New(v2schema.NotEmptyValidator(),
			option.Level[v2schema.Schema](errortrace.ErrorLevel),
			option.Type[v2schema.Schema](v2schema.NotEmptyValidatorType),
		),
	}
	defaultStreamValues = v2schema.Options{
		// default stream values
	}
)

// Manifest defines V2 schema specific behaviors. This interface is mainly for manifest's options.
type Manifest interface {
	stream.Manifest

	SetStream(v2schema.Schema) error
}

func New(opts ...Option) (Manifest, error) {
	m := manifest{
		defaults:   defaultStreamValues,
		validators: defaultValidators,
	}

	return &m, Options(opts).Apply(&m)
}

type manifest struct {
	stream v2schema.Schema

	defaults   v2schema.Options
	validators v2schema.Options
}

// The following implements stream.Manifest interface.

// MarshalJSON marshals a manifest into a json object.
func (m *manifest) MarshalJSON() ([]byte, error) {
	d, err := json.Marshal(m.stream)
	if err != nil {
		return nil, errMarshalTrace.AppendDetail(err)
	}

	return d, nil
}

// Unmarshal unmarshals a json object into a manifest.
func (m *manifest) UnmarshalJSON(data []byte) error {
	var stream v2schema.Schema

	if err := m.setDefaults(&stream); err != nil {
		return errUnmarshalTrace.AppendDetail(err)
	} else if err = json.Unmarshal(data, &stream); err != nil {
		return errUnmarshalTrace.AppendDetail(err)
	} else if err = m.SetStream(stream); err != nil {
		return errUnmarshalTrace.AppendDetail(err)
	}

	return nil
}

func (m manifest) GetID() string {
	return m.stream.CanonicalName
}

func (m manifest) GetFailureUrl() string {
	return m.stream.FailureURL
}

func (m manifest) GetExpired() time.Time {
	return m.stream.Expires
}

func (m manifest) GetTTLSeconds(now func() time.Time) int64 {
	return int64(math.Max(0, m.stream.Expires.Sub(now()).Seconds()))
}

func (m manifest) GetBatchMaxMesasges() int {
	return m.stream.BatchHint.MaxMesasges
}

func (m manifest) GetBatchLinger() time.Duration {
	return m.stream.BatchHint.MaxLingerDuration
}

func (m *manifest) ObfuscateSecrets() {
	for i := range m.stream.Webhooks {
		m.stream.Webhooks[i].Secret = "<obfuscated>"
	}

	for i := range m.stream.Kafkas {
		m.stream.Kafkas[i].KafkaProducer.Secret = "<obfuscated>"
	}
}

func (m *manifest) Validate() error {
	if errs := m.validators.Apply(&m.stream); errs != nil {
		return errValidatorTrace.AppendDetail(errs)
	}

	return nil
}

func (m manifest) GetStream() any {
	return m.stream
}

func (m manifest) GetSinks() (iter.Seq2[int, sink.Manifest], error) {
	if len(m.stream.Webhooks) != 0 && len(m.stream.Kafkas) != 0 {
		return nil, errors.New("")
	} else if len(m.stream.Webhooks) == 0 && len(m.stream.Kafkas) == 0 {
		return nil, errors.New("")
	}

	var i iter.Seq2[int, sink.Manifest]
	if len(m.stream.Webhooks) != 0 {
		i = webhook.IterateV2Manifest(m.stream.Webhooks)
	} else if len(m.stream.Kafkas) != 0 {
		i = kafka.IterateV1Manifest(m.stream.Kafkas)
	}

	return i, nil
}

func (m manifest) GetName() string {
	return "v2 wrpeventstream configuration"
}

// The following implements v2.Manifest interface.

func (m *manifest) SetStream(s v2schema.Schema) error {
	m.stream = s
	if errs := m.Validate(); errs != nil {
		// Validation failed, reset m.stream
		m.stream = v2schema.Schema{}

		return errs
	}

	return nil
}

func (m *manifest) SetDefaults() error {
	return m.setDefaults(&m.stream)
}

func (m manifest) setDefaults(stream *v2schema.Schema) error {
	if errs := m.defaults.Apply(stream); errs != nil {
		return errSetDefaultsTrace.AppendDetail(errs)
	}

	return nil
}
