// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package wrpeventstream

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/xmidt-org/webhook-schema/errortrace"
	"github.com/xmidt-org/webhook-schema/registry"
	"github.com/xmidt-org/webhook-schema/stream/v1/v1schema"
	"github.com/xmidt-org/webhook-schema/stream/v2/v2schema"
)

var (
	ErrNilStream            = fmt.Errorf("`%T` can't build a nil stream", Builder{})
	ErrUnmarshalRecord      = fmt.Errorf("failed to unmarshal `%T` into a wrpeventstream configuration", registry.Record{})
	errUnmarshalRecordTrace = errortrace.New(ErrUnmarshalRecord,
		errortrace.Level(errortrace.ErrorLevel),
		errortrace.Tag("UnmarshalRecord"),
	)
	ErrUnmarshalJson      = errors.New("failed to json unmarshal a wrpeventstream configuration")
	errUnmarshalJsonTrace = errortrace.New(ErrUnmarshalJson,
		errortrace.Level(errortrace.ErrorLevel),
		errortrace.Tag("UnmarshalJSON"),
	)
)

// Build provides the option to use an existing builders.
func Build(B Builder) BuilderOption {
	return BuilderOptionFunc(func(b *Builder) error {
		*b = B

		return nil
	})
}

// ClientStream provides the option for clients to use their own implementation of ClientManifest.
func ClientStream(s ClientManifest) BuilderOption {
	return BuilderOptionFunc(func(b *Builder) error {
		b.clientStream = s

		return nil
	})
}

// ClientStreamOptions provides the option for clients to provided their own ClientManifest's options.
// The provided options `ClientStreamOptions` will be applied at the end of `Builder.Build()`.
func ClientStreamOptions(opt ...Option) BuilderOption {
	return BuilderOptionFunc(func(b *Builder) error {
		b.clientStreamOptions = opt

		return nil
	})
}

// V1Defaults provides the option to set addition wrpeventstream v1 stream configuration defaults.
func V1Defaults(opts ...v1schema.Option) BuilderOption {
	return BuilderOptionFunc(func(b *Builder) error {
		b.V1.AddDefaults(opts...)

		return nil
	})
}

// V1Options provides the option to set addition wrpeventstream v1 stream configuration options.
func V1Options(opts ...v1schema.Option) BuilderOption {
	return BuilderOptionFunc(func(b *Builder) error {
		b.V1.AddOption(opts...)

		return nil
	})
}

// V2Defaults provides the option to set addition wrpeventstream v2 stream configuration defaults.
func V2Defaults(opts ...v2schema.Option) BuilderOption {
	return BuilderOptionFunc(func(b *Builder) error {
		b.V2.AddDefaults(opts...)

		return nil
	})
}

// V2Options provides the option to set addition wrpeventstream v2 stream configuration options.
func V2Options(opts ...v2schema.Option) BuilderOption {
	return BuilderOptionFunc(func(b *Builder) error {
		b.V2.AddOption(opts...)

		return nil
	})
}

// func UnmarshalJSON(data []byte) BuilderOption {
// 	//
// 	// Parameters:
// 	//   - clientStream: A struct containing atleast a stream and may also have json annotations.
// 	//
// 	// Returns:
// 	//   - Manifest: A Manifest interface.
// 	//   - error: An error indicating whether ClientOption(s) where applied successfully.
// 	//
// 	return unmarshal(func(clientStream ClientManifest) (errs error) {
// 		// Try to unmarshal `clientStream`.
// 		uerr := json.Unmarshal(data, clientStream)
// 		// If `data`'s nested wrpeventstream configuration stream is missing client related json tag, then
// 		// c.Validate() will return an emtpy configuration error.
// 		verr := clientStream.Validate()
// 		// If uerr and verr are nonnil, then append uerr and not both.
// 		// uerr is chosen over verr because after a successful unmarshalling, `Validate()` is called internally.
// 		// verr can be nonnil when uerr is nil, because `Validate()` is not called internally when no field or
// 		// json annotation matches are found during unmarshalling.
// 		// Meaning verr will always be nonnil when uerr is nonnil, but the inverse is not true.
// 		if uerr != nil {
// 			errs = multierr.Append(errs, uerr)
// 		} else if verr != nil {
// 			errs = multierr.Append(errs, verr)
// 		} else {
// 			// `data` contained a valid `clientStream`.
// 			return nil
// 		}

// 		// Try to unmarshal the stream directly.
// 		uerr = json.Unmarshal(data, clientStream.GetStream())
// 		// Check whether or not `data` contained a valid stream.
// 		verr = clientStream.Validate()
// 		if uerr != nil {
// 			errs = multierr.Append(errs, uerr)
// 		} else if verr != nil {
// 			errs = multierr.Append(errs, verr)
// 		} else {
// 			// `data` contained a valid stream.
// 			return nil
// 		}

// 		return errortrace.New(fmt.Errorf("failed to json unmarshal `%T` directly", clientStream.GetStream()),
// 			errortrace.Level(errortrace.ErrorLevel),
// 			errortrace.AppendDetail(errs),
// 			errortrace.Tag("UnmarshalJSON"),
// 		)
// 	})
// }

// UnmarshalJSON provides the option to unmarshal a json object into a wrpeventstream configuration.
func UnmarshalJSON(data []byte) BuilderOption {
	return unmarshal(func(c ClientManifest) error {
		// Try to unmarshal as much as `ClientManifest` directly.
		if err := json.Unmarshal(data, c); err != nil {
			return errUnmarshalJsonTrace.AppendDetail(err)
		}

		// If the client expects `data`'s nested wrpeventstream configuration stream to have a certain json tag but it doesn't, then
		// `ClientManifest.Validate()` will return an emtpy configuration error.
		if err := c.Validate(); err == nil {
			// `ClientManifest` was successfully unmarshalled
			return nil
		}

		// `ClientManifest.Validate()` failed, meaning data`'s nested wrpeventstream configuration stream may be missing json tag.
		// Try to unmarshal the wrpeventstream configuration stream directly.
		s := c.GetStream()
		if err := json.Unmarshal(data, s); err != nil {
			return errUnmarshalJsonTrace.AppendDetail(err)
		}

		c.SetStream(s)

		return nil
	})
}

// UnmarshalRecord provides the option to unmarshal a wrpeventstream registry record into a wrpeventstream configuration.
func UnmarshalRecord(record registry.Record) BuilderOption {
	return unmarshal(func(c ClientManifest) error {
		data, err := json.Marshal(record.Data)
		if err != nil {
			return errUnmarshalRecordTrace.AppendDetail(err)
		}

		// Since a `record` of the client's `ClientManifest` exists, that means their `ClientManifest` was marshalled into a `record` beforehand.
		// Meaning `record` must contain all of the client's expected json tags.
		// We're not responsible for the client's `ClientManifest` backwards compatible.
		if err = json.Unmarshal(data, c); err != nil {
			return errUnmarshalRecordTrace.AppendDetail(err)
		}

		// If the client's expected json tag for the wrpeventstream configuration stream is missing,
		// then `ClientManifest.Validate()` will return an emtpy configuration error.
		if err := c.Validate(); err != nil {
			return errUnmarshalRecordTrace.AppendDetail(err)
		}

		return nil
	})
}

// StreamValidator validates whether or not `Builder.clientStream` is non-nil.
func StreamValidator() BuilderOption {
	return BuilderOptionFunc(func(b *Builder) error {
		if b.clientStream != nil {
			return nil
		}

		return errortrace.New(ErrNilStream,
			errortrace.Level(errortrace.ErrorLevel),
			errortrace.Tag("StreamValidator"),
		)
	})
}

func unmarshal(f func(ClientManifest) error) BuilderOption {
	return BuilderOptionFunc(func(b *Builder) error {
		b.unmarshaller = f

		return nil
	})
}
