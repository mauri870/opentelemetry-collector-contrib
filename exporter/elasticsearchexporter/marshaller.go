// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package elasticsearchexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/elasticsearchexporter"

import (
	"errors"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/plog"
)

var (
	encodingJSON = "json"
)

// Marshaler configuration used for marhsaling Protobuf
var logsMarshalers = map[string]plog.Marshaler{
	encodingJSON: &plog.JSONMarshaler{},
}

type marshaller struct {
	logsMarshaler plog.Marshaler
}

func newMarshaller(conf *Config, host component.Host) (*marshaller, error) {
	if conf.Encoding == nil {
		// default to none
		return nil, nil
	}

	exts := host.GetExtensions()
	encoding := exts[*conf.Encoding]
	if encoding == nil {
		return nil, fmt.Errorf("unknown encoding %q", conf.Encoding)
	}
	// cast with ok to avoid panics.
	logm, ok := encoding.(plog.Marshaler)
	if !ok {
		return nil, fmt.Errorf("invalid encoding %q", conf.Encoding)
	}
	return &marshaller{
		logsMarshaler: logm,
	}, nil
}

func (m *marshaller) marshalLogs(ld plog.Logs) ([]byte, error) {
	if m.logsMarshaler == nil {
		return nil, errors.New("logs are not supported by encoding")
	}
	buf, err := m.logsMarshaler.MarshalLogs(ld)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal logs: %w", err)
	}
	return buf, nil
}
