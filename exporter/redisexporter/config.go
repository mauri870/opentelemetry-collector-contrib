// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package redisexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/redisexporter"

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/config/configopaque"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

// Config defines configuration for Elastic exporter.
// TODO(mauri870): Document this
type Config struct {
	// TODO(mauri870): Queue settings?
	QueueSettings exporterhelper.QueueConfig `mapstructure:"sending_queue"`

	confignet.AddrConfig `mapstructure:",squash"`
	Auth                 AuthSettings `mapstructure:",squash"`
	DB                   int          `mapstructure:"db"`

	// Encoding defines the encoding of the telemetry data.
	// If specified, it overrides `FormatType` and applies an encoding extension.
	// TODO(mauri870): uses an encodingextension, implement this, look at filexporter
	// Support a raw encoding, just like bodymap in elasticsearch
	Encoding *component.ID `mapstructure:"encoding"`

	TLS configtls.ClientConfig `mapstructure:"tls,omitempty"`

	// TODO(mauri870): batcher settings?
}

func (c *Config) Validate() error {
	// TODO(mauri870): implement this

	return nil
}

type AuthSettings struct {
	Username string              `mapstructure:"username"`
	Password configopaque.String `mapstructure:"password"`
}
