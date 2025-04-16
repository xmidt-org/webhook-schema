// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package wrpeventstream

import (
	"errors"
	"fmt"
	"time"

	"github.com/xmidt-org/webhook-schema/errortrace"
	"github.com/xmidt-org/webhook-schema/registry"
	"github.com/xmidt-org/webhook-schema/stream"
	v2 "github.com/xmidt-org/webhook-schema/stream/v2"
)

var (
	ErrNew      = errors.New("failed to create and initialize client's stream configuration")
	errNewTrace = errortrace.New(ErrBuild,
		errortrace.Level(errortrace.ErrorLevel),
		errortrace.Tag("wrpeventstream.New"),
	)
)

// New returns a `ClientManifest` based on the provided builder options, where:
//   - a `ClientStream` option is expected if a client wants to use their own custom stream configuration
//     (containing a nested wrpeventstream configuration). Otherwise, defaultStream will be used, processed and return
//     as the `ClientManifest`.
//   - a `Build` option is expected if a client wants to use a configured `Builder`. A `Builder` configuration
//     determines which wrpeventstream configuration versions will be supported and used, along with which default
//     settings, validations and validation error levels are used
//   - an unmarshal option, such as , is expected if a client wants to deserialize a stream configration (custom or not)
//     and set it to `ClientManifest` via `ClientManifest.SetStream(stream)`, where `stream` is the deserialized stream configration
func New(opts ...BuilderOption) (ClientManifest, error) {
	var (
		Builder       Builder
		defaultStream defaultStream
	)

	v2stream, err := v2.New()
	if err != nil {
		return nil, errNewTrace.AppendDetail(fmt.Errorf("unexpected error occurred: %s", err))
	}

	defaultStream.SetStream(v2stream)
	defaultOpts := BuilderOptions{
		ClientStream(&defaultStream),
	}
	opts = append(defaultOpts, opts...)
	defaultVal := BuilderOptions{
		StreamValidator(),
	}
	opts = append(opts, defaultVal)
	if err := BuilderOptions(opts).Apply(&Builder); err != nil {
		return nil, errNewTrace.AppendDetail(err)
	}

	stream, err := Builder.Build()
	if err != nil {
		return nil, errNewTrace.AppendDetail(err)
	}

	return stream, nil
}

// ClientManifest defines the functionality required to handle and process stream configurations while being stream version and sink agnostic.
// ClientManifest is the starting point for all clients working with wrpeventstream configurations and will likely satisfy most basic
// use cases, e.g.: an wrpeventstream registration service, where stream configrations are moved in and out of an wrpeventstream registry
// and new configurations are validated and stored.
// While other clients with more advance use cases, such as standing up an actual stream for event processing, will need to call `ClientManifest.GetStream()`
// to upgrade to their ClientManifest to a `stream.Manifest`.
// ClientManifest must be used with `New()` due to the following:
//   - A client's ClientManifest is not usable until it's passed into `New()` using the `ClientStream` option, along with other Builder options
//     such as the stream configration deserialization.
//   - Serialized stream configration (custom or not) can only be deserialized into a ClientManifest, if a Builder stream configration deserialization
//     option is used.
//   - Default settings and validation for wrpeventstream configurations occur within `New()`.
type ClientManifest interface {
	// registry.Registrable defines the expected functionality for ClientManifest (custom client stream configurations)
	// serialization, required in order to store in an wrpeventstream registry.
	registry.Registrable

	// SetStream takes a wrpeventstream configuration and sets it in
	SetStream(stream.Manifest)
	GetStream() stream.Manifest
	ObfuscateSecrets()
	ToRecord() (registry.Record, error)
}

type defaultStream struct {
	Stream stream.Manifest `json:"wrp_event_stream"`
}

func (s *defaultStream) ToRecord() (registry.Record, error) {
	return registry.New(s)
}

func (s *defaultStream) SetStream(stream stream.Manifest) {
	s.Stream = stream
}

func (s *defaultStream) GetStream() stream.Manifest {
	return s.Stream
}

func (s *defaultStream) GetID() string {
	return s.Stream.GetID()
}

func (s *defaultStream) GetTTLSeconds(now func() time.Time) int64 {
	return s.Stream.GetTTLSeconds(now)
}

func (s *defaultStream) ObfuscateSecrets() {
	s.Stream.ObfuscateSecrets()
}

func (s *defaultStream) Validate() error {
	return s.Stream.Validate()
}
