// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/xmidt-org/webhook-schema/stream/option"
	"github.com/xmidt-org/webhook-schema/stream/v1/v1schema"
	"go.uber.org/multierr"
)

var (
	ErrBuildDisabled             = errors.New("v1 stream build failure: builder is disabled")
	ErrOptionConfigInvalid       = errors.New("invalid schema v1 option configuration error")
	errOptionConfigUnmarshalling = errors.New("unmarshalling OptionConfig error")
)

// Deprecated: This package should only be used for backwards compatibility
// reasons. Use v2 instead.
// TODO Add logger
type Builder struct {
	Config

	defaults   v1schema.Options
	validators v1schema.Options
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

func (b *Builder) AddDefaults(opts ...v1schema.Option) {
	b.defaults = append(b.defaults, opts...)
}

func (b *Builder) AddOption(opts ...v1schema.Option) {
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

		opts := option.Options[v1schema.Schema]{
			option.Type[v1schema.Schema](oc.Type),
			option.Level[v1schema.Schema](oc.Level),
		}

		switch oc.Type {
		case v1schema.AddressDefaultType, v1schema.MatcherType, v1schema.SetSchemaType:
			errs = fmt.Errorf("%w: option `%b` can't be unmarshalled", v1schema.ErrOptionTypeInvalid, oc.Type)
		case v1schema.AlwaysValidType:
			b.validators = append(b.validators, option.New(v1schema.AlwaysValid(), opts...))
		case v1schema.AtleastOneEventValidatorType:
			b.validators = append(b.validators, option.New(v1schema.AtleastOneEventValidator(), opts...))
		case v1schema.EventRegexValidatorType:
			b.validators = append(b.validators, option.New(v1schema.EventRegexValidator(), opts...))
		case v1schema.DeviceIDRegexValidatorType:
			b.validators = append(b.validators, option.New(v1schema.DeviceIDRegexValidator(), opts...))
		case v1schema.DurationValidatorType:
			b.validators = append(b.validators, option.New(v1schema.DurationValidator(time.Now, oc.URLChecker.TTL.Max), opts...))
		case v1schema.CheckUntilValidatorType:
			b.validators = append(b.validators, option.New(v1schema.CheckUntilValidator(time.Now, oc.URLChecker.TTL.Jitter, oc.URLChecker.TTL.Max), opts...))
		case v1schema.ReceiverURLValidatorType:
			checker, err := oc.URLChecker.Build()
			errs = multierr.Append(errs, err)
			b.validators = append(b.validators, option.New(v1schema.ReceiverURLValidator(checker), opts...))
		case v1schema.FailureURLValidatorType:
			checker, err := oc.URLChecker.Build()
			errs = multierr.Append(errs, err)
			b.validators = append(b.validators, option.New(v1schema.ReceiverURLValidator(checker), opts...))
		case v1schema.AlternativeURLValidatorType:
			checker, err := oc.URLChecker.Build()
			errs = multierr.Append(errs, err)
			b.validators = append(b.validators, option.New(v1schema.AlternativeURLValidator(checker), opts...))
		case v1schema.UntilValidatorType:
			b.validators = append(b.validators, option.New(v1schema.UntilValidator(oc.URLChecker.TTL.Jitter, oc.URLChecker.TTL.Max, time.Now), opts...))
		default:
			errs = fmt.Errorf("%w: option `%b`", v1schema.ErrOptionTypeInvalid, oc.Type)
		}
	}

	return errs
}

// Config simplifies the unmarshalling process.
type Config struct {
	Options []OptionConfig `json:"options"`
	// Disable determines whether the schema is active (`diable` is `false`)
	// or inactive (`disable` is `true`).
	// "active" means v1 schema will be supported by the client (i.e. decoding of json v1 schemas will succeed).
	// "inactive" means v1 schema will not be supported by the client (i.e. decoding of json v1 schemas will fail).
	// Default is `false`.
	Disable bool `json:"disable"`
}
