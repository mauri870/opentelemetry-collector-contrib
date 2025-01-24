// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package redisexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/redisexporter"

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

type redisExporter struct {
	config *Config
	client client
	logger *zap.Logger
}

func newExporter(
	cfg *Config,
	set exporter.Settings,
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
		logger: set.Logger,
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
	start := time.Now()

	for i := 0; i < ld.ResourceLogs().Len(); i++ {
		logs := ld.ResourceLogs().At(i)
		res := logs.Resource()
		resAttr, err := json.Marshal(res.Attributes().AsRaw())
		if err != nil {
			return err
		}

		for j := 0; j < logs.ScopeLogs().Len(); j++ {
			rs := logs.ScopeLogs().At(j).LogRecords()
			for k := 0; k < rs.Len(); k++ {
				r := rs.At(k)
				logAttr, err := json.Marshal(r.Attributes().AsRaw())
				if err != nil {
					return err
				}
				bodyByte, err := json.Marshal(r.Body().AsRaw())
				if err != nil {
					return err
				}

				l := logPayload{
					ResourceAttributes:  resAttr,
					LogRecordAttributes: logAttr,
					LogRecordBody:       bodyByte,
				}
				err = e.client.PublishLog(ctx, l)
				if err != nil {
					e.logger.Error("failed to send logs to redis", zap.Error(err))
				}
			}
		}
	}

	duration := time.Since(start)
	e.logger.Info("published logs", zap.Int("records", ld.LogRecordCount()),
		zap.String("cost", duration.String()))
	return nil
}
