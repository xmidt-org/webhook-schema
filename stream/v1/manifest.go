// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package v1

import (
	"encoding/json"
	"fmt"
	"iter"
	"math"
	"time"

	"github.com/xmidt-org/webhook-schema/errortrace"
	"github.com/xmidt-org/webhook-schema/stream"
	"github.com/xmidt-org/webhook-schema/stream/option"
	"github.com/xmidt-org/webhook-schema/stream/sink"
	"github.com/xmidt-org/webhook-schema/stream/sink/webhook"
	"github.com/xmidt-org/webhook-schema/stream/v1/v1schema"
)

var (
	ErrMarshal      = fmt.Errorf("`%T` marshal error", v1schema.Schema{})
	errMarshalTrace = errortrace.New(ErrMarshal,
		errortrace.Level(errortrace.ErrorLevel),
		errortrace.Tag(fmt.Sprintf("%T.MarshalJSON", manifest{})))
	ErrUnmarshal      = fmt.Errorf("`%T` unmarshal error", v1schema.Schema{})
	errUnmarshalTrace = errortrace.New(ErrUnmarshal,
		errortrace.Level(errortrace.ErrorLevel),
		errortrace.Tag(fmt.Sprintf("%T.UnmarshalJSON", manifest{})))
	ErrValidate       = fmt.Errorf("`%T` failed validation", v1schema.Schema{})
	errValidatorTrace = errortrace.New(ErrValidate,
		errortrace.Level(errortrace.ErrorLevel),
		errortrace.Tag(fmt.Sprintf("%T.Validate", manifest{})))
	ErrSetDefaults      = fmt.Errorf("`%T` set default failure", v1schema.Schema{})
	errSetDefaultsTrace = errortrace.New(ErrSetDefaults,
		errortrace.Level(errortrace.ErrorLevel),
		errortrace.Tag(fmt.Sprintf("%T.SetDefaults", manifest{})))
)

var (
	defaultValidators = v1schema.Options{
		// defaults validators
		option.New(v1schema.NotEmptyValidator(),
			option.Level[v1schema.Schema](errortrace.ErrorLevel),
			option.Type[v1schema.Schema](v1schema.NotEmptyValidatorType),
		),
	}
	defaultStreamValues = v1schema.Options{
		// default stream values
	}
)

// Manifest defines V2 schema specific behaviors. This interface is mainly for manifest's options.
type Manifest interface {
	stream.Manifest

	GetUntil() time.Time
	SetStream(v1schema.Schema) error
}

func New(opts ...Option) (Manifest, error) {
	m := manifest{
		defaults:   defaultStreamValues,
		validators: defaultValidators,
	}

	return &m, Options(opts).Apply(&m)
}

type manifest struct {
	stream v1schema.Schema

	defaults   v1schema.Options
	validators v1schema.Options
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
	var stream v1schema.Schema

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
	return m.stream.Webhook.ReceiverURL
}

func (m manifest) GetFailureUrl() string {
	return m.stream.FailureURL
}

func (m manifest) GetUntil() time.Time {
	return m.stream.Until
}

// Note, Schema.Until is known to be poorly named and it's actually
// used as Expired.
func (m manifest) GetExpired() time.Time {
	return m.stream.Until
}

func (m manifest) GetTTLSeconds(now func() time.Time) int64 {
	return int64(math.Max(0, m.stream.Until.Sub(now()).Seconds()))
}

func (m manifest) GetBatchMaxMesasges() int {
	return 1
}

func (m manifest) GetBatchLinger() time.Duration {
	return 0
}

func (m *manifest) ObfuscateSecrets() {
	m.stream.Webhook.Secret = "<obfuscated>"
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
	// nolint:staticcheck
	return webhook.IterateV1Manifest([]webhook.V1Schema{m.stream.Webhook}), nil
}

func (m manifest) GetName() string {
	return "v1 wrpeventstream configuration"
}

// The following implements v2.Manifest interface.

func (m *manifest) SetStream(s v1schema.Schema) error {
	m.stream = s
	if errs := m.Validate(); errs != nil {
		// Validation failed, reset m.stream
		m.stream = v1schema.Schema{}

		return errs
	}

	return nil
}

func (m *manifest) SetDefaults() (errs error) {
	return m.setDefaults(&m.stream)
}

func (m manifest) setDefaults(stream *v1schema.Schema) error {
	if errs := m.defaults.Apply(stream); errs != nil {
		return errSetDefaultsTrace.AppendDetail(errs)
	}

	return nil
}
