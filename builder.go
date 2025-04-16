// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package wrpeventstream

import (
	"errors"
	"fmt"

	"github.com/xmidt-org/webhook-schema/errortrace"
	"github.com/xmidt-org/webhook-schema/stream"
	v1 "github.com/xmidt-org/webhook-schema/stream/v1"
	v2 "github.com/xmidt-org/webhook-schema/stream/v2"
	"go.uber.org/multierr"
)

var (
	ErrBuild      = errors.New("failed to build stream")
	errBuildTrace = errortrace.New(ErrBuild,
		errortrace.Level(errortrace.ErrorLevel),
		errortrace.Tag(fmt.Sprintf("%T.Build", Builder{})))
	ErrAllUnmarshalAttemptsFailed      = errors.New("all supported stream version unmarshals failed")
	errAllUnmarshalAttemptsFailedTrace = errortrace.New(ErrAllUnmarshalAttemptsFailed,
		errortrace.Level(errortrace.ErrorLevel),
		errortrace.Tag(fmt.Sprintf("%T.unmarshaller", Builder{})))
	ErrAllBuilderDisabled      = errors.New("all supported stream version builders are disabled")
	errAllBuilderDisabledTrace = errortrace.New(ErrAllBuilderDisabled,
		errortrace.Level(errortrace.ErrorLevel),
		errortrace.Tag(fmt.Sprintf("%T.getSupportedVersions", Builder{})))
)

// Builder is a configuration struct used to unmarshal stream configuration objects, colloquially referred as a `record`, from an wrpeventstream registry.
// Builder is designed to unmarshal any stream configuration version based on its configuration and `New()` options.
// Meaning the stream configuration version support, validation and validation error level can be configured and embedded in a client's configuration.
// Builder makes no assumptions about records, expect for the following two requirements:
// - it is a json object
// - it contains 1 nested wrpeventstream configuration somewhere within
//   - if it does not contain at least 1, then `Builder.Build()` will return an ErrAllUnmarshalAttemptsFailed error
//   - if it contains more than 1, then `Builder.Build()` will use the first successful unmarshal among the nested configurations
//
// TODO Add logger
type Builder struct {
	// unmarshaller is a client provided option and can be set via the `UnmarshalRecord(record)` or `UnmarshalJSON(jsonBytes)` options.
	// Where `UnmarshalRecord` takes a `record` (a stream configuration object from an wrpeventstream registry) and unmarshals it into
	// `Builder.clientStream`, while `UnmarshalJSON` functions the same but it takes in json data instead of `records`.
	// `Builder.Build()` calls `Builder.unmarshaller` and returns its unmarshalled `Builder.clientStream`.
	unmarshaller func(ClientManifest) error

	// clientStream is set via the `ClientStream(cstream)` option, where cstream is some struct that implements ClientManifest.
	clientStream ClientManifest
	// clientStreamOptions is a list of client implemented options that are applied to clientStream after `Builder.Build()` succeeds.
	clientStreamOptions Options

	// v1 wrpeventstream configuration builder
	// nolint:staticcheck
	V1 v1.Builder `json:"v1"`
	// v2 wrpeventstream configuration builder
	V2 v2.Builder `json:"v2"`
}

// Build translates the configuration into a stream.
func (b Builder) Build() (ClientManifest, error) {
	clientStream, err := b.unmarshal()
	if err != nil {
		return nil, errBuildTrace.AppendDetail(err)
	}

	return clientStream, b.clientStreamOptions.Apply(clientStream)
}

// UnmarshalEnabled indicates whether or the builders has been configured to unmarshal some wrpeventstream configuration.
// To enable unmarshalling, pass a unmarshal builder option to `New()`.
func (b Builder) UnmarshalEnabled() bool { return b.unmarshaller != nil }

// unmarshal return a ClientManifest with a set and unmarshalled wrpeventstream configuration.
// The stream configuration's version is determined by the first successful wrpeventstream version unmarshal.
// unmarshal is a NoOp if the builder's unmarshal is not enabled, i.e.: an unmarshaller option was not passed into `New()`.
// An error is only return in the following cases:
//   - all supported stream versions are all disabled
//   - all supported stream version unmarshals failed
func (b Builder) unmarshal() (clientStream ClientManifest, errs error) {
	streamVersions, err := b.getSupportedVersions()
	if err != nil {
		return nil, errBuildTrace.AppendDetail(err)
	}

	if !b.UnmarshalEnabled() {
		// TODO Add debug log

		return nil, nil
	}

	clientStream = b.clientStream
	for _, sv := range streamVersions {
		// Try each supported stream version.
		clientStream.SetStream(sv)
		// A successful unmarshalling, it's a version match.
		err := b.unmarshaller(clientStream)
		if err != nil {
			// Not a match, try another version.
			err = errortrace.New(fmt.Errorf("client stream config `%T` is not a match to `%T`: unmarshalling failure", clientStream, sv.GetStream()),
				errortrace.Level(errortrace.ErrorLevel),
				errortrace.AppendDetail(err),
				errortrace.Tag(fmt.Sprintf("%T.UnmarshalJSON", sv)))
			errs = multierr.Append(errs, err)

			continue
		}

		// A version deserialization succeeded, meaning a match was found.
		// Ignore the errors from the previous versions that failed.
		errs = nil

		// Stop trying the remaining versions.
		break
	}

	if errs != nil {
		// No version match was found.
		return nil, errAllUnmarshalAttemptsFailedTrace.AppendDetail(errs)
	}

	return clientStream, nil
}

// getSupportedVersions returns a list of supported stream versions.
func (b Builder) getSupportedVersions() (streams []stream.Manifest, errs error) {
	v1stream, err := b.V1.Build()
	if err == nil {
		streams = append(streams, v1stream)
	} else if !errors.Is(err, v1.ErrBuildDisabled) {
		// v1.ErrBuildDisabled errors alone are not an issue, as long as all builders are not disabled.
		return nil, err
	}

	errs = multierr.Append(errs, err)
	v2stream, err := b.V2.Build()
	if err == nil {
		streams = append(streams, v2stream)
	} else if !errors.Is(err, v2.ErrBuildDisabled) {
		// v2.ErrBuildDisabled errors alone are not an issue, as long as all builders are not disabled.
		return nil, err
	}

	errs = multierr.Append(errs, err)
	if len(streams) == 0 {
		// all builders are not disabled.
		return nil, errAllBuilderDisabledTrace.AppendDetail(errs)
	}

	return streams, nil
}
