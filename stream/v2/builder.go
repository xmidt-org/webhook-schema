// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package v2

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/xmidt-org/webhook-schema/stream/option"
	"github.com/xmidt-org/webhook-schema/stream/v2/v2schema"
	"go.uber.org/multierr"
)

var (
	ErrBuildDisabled             = errors.New("v2 stream build failure: builder is disabled")
	ErrOptionConfigInvalid       = errors.New("invalid schema v2 option configuration error")
	errOptionConfigUnmarshalling = errors.New("unmarshalling OptionConfig error")
)

// Config defines wrpeventstream schema options used to
// set defaults and validate registered wrpeventstream.
// (Optional)
// Available option levels: info, warning, error
// Available V1 options:
// Available V2 options:
// All options can be disabled with `disable: true`, it is false by default.
// TODO Add logger
type Builder struct {
	Config

	defaults   v2schema.Options
	validators v2schema.Options
}

// UnmarshalJSON unmarshals a json into a option.
func (b *Builder) UnmarshalJSON(data []byte) (errs error) {
	if len(data) == 0 {
		return nil
	}

	errs = multierr.Append(errs, json.Unmarshal(data, &b.Config))
	errs = multierr.Append(errs, b.gatherStreamOptions())
	if errs != nil {
		return multierr.Append(errOptionConfigUnmarshalling, errs)
	}

	return nil
}

// Build translates the configuration into a stream.
func (b Builder) Build() (Manifest, error) {
	if b.Disable {
		return nil, ErrBuildDisabled
	}

	if err := b.gatherStreamOptions(); err != nil {
		return nil, err
	}

	return New(AddDefaults(b.defaults...), AddValidators(b.validators...))
}

func (b *Builder) AddDefaults(opts ...v2schema.Option) {
	b.defaults = append(b.defaults, opts...)
}

func (b *Builder) AddOption(opts ...v2schema.Option) {
	b.validators = append(b.validators, opts...)
}

func (b *Builder) gatherStreamOptions() (errs error) {
	for _, oc := range b.Options {
		if oc.Disable {
			continue
		}
		if !oc.IsValid() {
			errs = multierr.Append(errs, fmt.Errorf("%w: `%s`", ErrOptionConfigInvalid, oc.Type))

			continue
		}

		opts := option.Options[v2schema.Schema]{
			option.Type[v2schema.Schema](oc.Type),
			option.Level[v2schema.Schema](oc.Level),
		}

		switch oc.Type {
		case v2schema.AddressDefaultType, v2schema.SetSchemaType:
			errs = multierr.Append(errs, fmt.Errorf("%w `%s` can't be unmarshalled", v2schema.ErrOptionTypeInvalid, oc.Type))
		case v2schema.AlwaysValidType:
			b.validators = append(b.validators, option.New(v2schema.AlwaysValid(), opts...))
		case v2schema.OnlyWebhooksValidatorType:
			b.validators = append(b.validators, option.New(v2schema.OnlyWebhooksValidator(), opts...))
		case v2schema.EventRegexValidatorType:
			b.validators = append(b.validators, option.New(v2schema.EventRegexValidator(), opts...))
		case v2schema.ExpiresValidatorType:
			b.validators = append(b.validators, option.New(v2schema.ExpiresValidator(), opts...))
		case v2schema.ReceiverURLValidatorType:
			checker, err := oc.URLChecker.Build()
			errs = multierr.Append(errs, err)
			b.validators = append(b.validators, option.New(v2schema.ReceiverURLValidator(checker), opts...))
		case v2schema.FailureURLValidatorType:
			checker, err := oc.URLChecker.Build()
			errs = multierr.Append(errs, err)
			b.validators = append(b.validators, option.New(v2schema.FailureURLValidator(checker), opts...))
		default:
			errs = multierr.Append(errs, fmt.Errorf("%w `%s`", v2schema.ErrOptionTypeInvalid, oc.Type))
		}
	}

	return errs
}

// Config simplifies the unmarshalling process.
type Config struct {
	// at least 1 option required
	Options []OptionConfig `json:"options"`
	// Disable determines whether the schema is active (`diable` is `false`)
	// or inactive (`disable` is `true`).
	// "active" means v2 schema will be supported by the client (i.e. decoding of json v2 schemas will succeed).
	// "inactive" means v2 schema will not be supported by the client (i.e. decoding of json v2 schemas will fail).
	// Default is `false`.
	Disable bool `json:"disable"`
}
