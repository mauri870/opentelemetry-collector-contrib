// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package jsonlogencodingextension // import "github.com/open-telemetry/opentelemetry-collector-contrib/extension/encoding/jsonlogencodingextension"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension"
)

var (
	Type      = component.MustNewType("json_log_encoding")
	ScopeName = "github.com/open-telemetry/opentelemetry-collector-contrib/extension/encoding/jsonlogencodingextension"
)

const (
	ExtensionStability = component.StabilityLevelAlpha
)

func NewFactory() extension.Factory {
	return extension.NewFactory(
		Type,
		createDefaultConfig,
		createExtension,
		ExtensionStability,
	)
}

func createExtension(_ context.Context, _ extension.Settings, config component.Config) (extension.Extension, error) {
	return &jsonLogExtension{
		config: config,
	}, nil
}

func createDefaultConfig() component.Config {
	return &Config{
		Mode: JSONEncodingModeBody,
	}
}
