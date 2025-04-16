// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package sink

// Manifest defines basic behaviors to handle sinks agnostically .
// For sink specific
// Meaning, this Manifest (wrpeventstream.Manifest) can't access the level of behaviors and information required
// to programmatically setup a wrpeventstream and process incoming events.
// But, by calling `Manifest.GetStream()`, that level of behaviors and information is accessible through stream.Manifest.
type Manifest interface {
	GetName() string
	GetSink() any
}
