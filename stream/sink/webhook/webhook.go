// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package webhook

import "github.com/xmidt-org/webhook-schema/stream/sink"

type Manifest interface {
	sink.Manifest
	GetReceiverURLs() []string
	GetSecret() string
}
