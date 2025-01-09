// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package redisexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/redisexporter"

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type redisExporter struct {
	config *Config
	client client
}

func newExporter(
	cfg *Config,
	_ exporter.Settings,
) (*redisExporter, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	opts := &redis.Options{
		Addr:     cfg.Endpoint,
		Username: cfg.Auth.Username,
		Password: string(cfg.Auth.Password),
		Network:  string(cfg.Transport),
		DB:       cfg.DB,
	}

	return &redisExporter{
		config: cfg,
		client: newRedisClient(opts),
	}, nil
}

func (e *redisExporter) Start(ctx context.Context, host component.Host) error {
	// TODO(mauri870): create redis client based on cfg
	return nil
}

func (e *redisExporter) Shutdown(ctx context.Context) error {
	doneCh := make(chan struct{})
	go func() {
		// TODO(mauri870): close redis client
		close(doneCh)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-doneCh:
		return nil
	}
}

func (e *redisExporter) pushMetricsData(
	ctx context.Context,
	metrics pmetric.Metrics,
) error {
	return nil
}

func (e *redisExporter) pushTraceData(
	ctx context.Context,
	td ptrace.Traces,
) error {
	return nil
}

func (e *redisExporter) pushLogsData(ctx context.Context, ld plog.Logs) error {
	return nil
}
