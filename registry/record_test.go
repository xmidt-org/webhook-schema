// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package registry_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	wrpeventstream "github.com/xmidt-org/webhook-schema"
	"github.com/xmidt-org/webhook-schema/registry"
	"github.com/xmidt-org/webhook-schema/stream/v1/v1schema"
	"go.uber.org/multierr"
)

func TestRecord(t *testing.T) {
	records, err := getTestRecords()
	require.NoError(t, err)
	require.NotEmpty(t, records)
	streams, err := getTestStreams()
	require.NoError(t, err)
	require.NotEmpty(t, streams)

	tcs := []struct {
		Description  string
		InputSchema  []wrpeventstream.ClientManifest
		ExpectedItem []registry.Record
		ExpectedErr  bool
	}{
		{
			Description:  "Happy path",
			InputSchema:  streams,
			ExpectedItem: records,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.Description, func(t *testing.T) {
			for i, stream := range tc.InputSchema {
				t.Run(fmt.Sprintf("stream %d", i), func(t *testing.T) {
					assert := assert.New(t)
					item, err := registry.New(stream)
					if tc.ExpectedErr {
						assert.Error(err)
						return
					}

					assert.NoError(err)
					assert.Equal(tc.ExpectedItem[i], item)
				})
			}
		})
	}
}

func getTestRecords() ([]registry.Record, error) {
	var (
		record1ExpiresInSecs int64 = 10
		record2ExpiresInSecs int64 = 20
		recordExpectedTTL    int64 = 0
		record1Duration            = time.Duration(record1ExpiresInSecs) * time.Second
		record2Duration            = time.Duration(record2ExpiresInSecs) * time.Second
	)

	refTime, err := time.Parse(time.RFC3339, "2021-01-02T15:04:00Z")
	if err != nil {
		return nil, err
	}

	return []registry.Record{
		{
			ID: "a379a6f6eeafb9a55e378c118034e2751e682fab9f2d30ab13d2125586ce1947",
			Data: map[string]any{
				"wrp_event_stream": map[string]any{
					"registered_from_address": "example.com",
					"config": map[string]any{
						"url":          "example.com",
						"content_type": "application/json",
						"secret":       "superSecretXYZ",
					},
					"events": []any{"online"},
					"matcher": map[string]any{
						"device_id": []any{"mac:aabbccddee.*"},
					},
					"failure_url": "example.com",
					"duration":    v1schema.CustomDuration(record1Duration).String(),
					"until":       refTime.Add(record1Duration).Format(time.RFC3339),
				},
			},
			TTL: &recordExpectedTTL,
		},
		{
			ID: "a379a6f6eeafb9a55e378c118034e2751e682fab9f2d30ab13d2125586ce1947",
			Data: map[string]any{
				"wrp_event_stream": map[string]any{
					"registered_from_address": "example.com",
					"config": map[string]any{
						"url":          "example.com",
						"content_type": "application/json",
						"secret":       "doNotShare:e=mc^2",
					},
					"events": []any{"online"},
					"matcher": map[string]any{
						"device_id": []any{"mac:aabbccddee.*"},
					},
					"failure_url": "example.com",
					"duration":    v1schema.CustomDuration(record2Duration).String(),
					"until":       refTime.Add(record2Duration).Format(time.RFC3339),
				},
			},
			TTL: &recordExpectedTTL,
		},
	}, nil
}

func getTestStreams() (streams []wrpeventstream.ClientManifest, errs error) {
	records, err := getTestRecords()
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		s, err := wrpeventstream.New(
			wrpeventstream.UnmarshalRecord(record))
		if err != nil {
			errs = multierr.Append(errs, err)
			continue
		}

		streams = append(streams, s)
	}

	return streams, errs
}
