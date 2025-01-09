// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package redisexporter

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/confmap/confmaptest"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/redisexporter/internal/metadata"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	defaultCfg := NewFactory().CreateDefaultConfig()

	tests := []struct {
		configFile string
		id         component.ID
		expected   func() component.Config
	}{
		{
			id:         component.NewIDWithName(metadata.Type, "default"),
			configFile: "config.yaml",
			expected:   func() component.Config { return defaultCfg },
		},
		{
			id:         component.NewIDWithName(metadata.Type, "custom"),
			configFile: "config.yaml",
			expected: func() component.Config {
				cfg := defaultCfg.(*Config)
				cfg.AddrConfig = confignet.AddrConfig{
					Endpoint:  "redis:6379",
					Transport: confignet.TransportTypeUnix,
				}
				cfg.DB = 1
				cfg.TLS = configtls.ClientConfig{
					Insecure: false,
				}
				cfg.Auth = AuthSettings{
					Username: "username",
					Password: "password",
				}

				return cfg
			},
		},
	}

	for _, tt := range tests {
		t.Run(strings.ReplaceAll(tt.id.String(), "/", "_"), func(t *testing.T) {
			factory := NewFactory()
			cfg := factory.CreateDefaultConfig()

			cm, err := confmaptest.LoadConf(filepath.Join("testdata", tt.configFile))
			require.NoError(t, err)

			sub, err := cm.Sub(tt.id.String())
			require.NoError(t, err)
			require.NoError(t, sub.Unmarshal(cfg))

			assert.NoError(t, component.ValidateConfig(cfg))
			assert.Equal(t, tt.expected(), cfg)
		})
	}
}
