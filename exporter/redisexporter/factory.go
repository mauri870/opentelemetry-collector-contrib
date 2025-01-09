// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:generate mdatagen metadata.yaml

package redisexporter

import (
	"context"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/redisexporter/internal/metadata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

// NewFactory creates a factory for the redis exporter.
func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		metadata.Type,
		createDefaultConfig,
		exporter.WithLogs(createLogsExporter, metadata.LogsStability),
		exporter.WithMetrics(createMetricsExporter, metadata.MetricsStability),
		exporter.WithTraces(createTracesExporter, metadata.TracesStability),
	)
}

func createDefaultConfig() component.Config {
	qs := exporterhelper.NewDefaultQueueConfig()
	qs.Enabled = false

	return &Config{
		QueueSettings: qs,
		AddrConfig: confignet.AddrConfig{
			Endpoint:  "localhost:6379",
			Transport: confignet.TransportTypeTCP,
		},
		TLS: configtls.ClientConfig{
			Insecure: true,
		},
		// TelemetrySettings: TelemetrySettings{
		// 	LogRequestBody:  false,
		// 	LogResponseBody: false,
		// },
		// Batcher: BatcherConfig{
		// 	FlushTimeout: 30 * time.Second,
		// 	MinSizeConfig: exporterbatcher.MinSizeConfig{
		// 		MinSizeItems: 5000,
		// 	},
		// 	MaxSizeConfig: exporterbatcher.MaxSizeConfig{
		// 		MaxSizeItems: 0,
		// 	},
		// },
		// Flush: FlushSettings{
		// 	Bytes:    5e+6,
		// 	Interval: 30 * time.Second,
		// },
	}
}

// createLogsExporter creates a new exporter for logs.
//
// Logs are directly indexed into Elasticsearch.
func createLogsExporter(
	ctx context.Context,
	set exporter.Settings,
	cfg component.Config,
) (exporter.Logs, error) {
	cf := cfg.(*Config)

	exporter, err := newExporter(cf, set)
	if err != nil {
		return nil, err
	}

	return exporterhelper.NewLogs(
		ctx,
		set,
		cfg,
		exporter.pushLogsData,
		exporterhelperOptions(cf, exporter.Start, exporter.Shutdown)...,
	)
}

func createMetricsExporter(
	ctx context.Context,
	set exporter.Settings,
	cfg component.Config,
) (exporter.Metrics, error) {
	cf := cfg.(*Config)

	exporter, err := newExporter(cf, set)
	if err != nil {
		return nil, err
	}

	return exporterhelper.NewMetrics(
		ctx,
		set,
		cfg,
		exporter.pushMetricsData,
		exporterhelperOptions(cf, exporter.Start, exporter.Shutdown)...,
	)
}

func createTracesExporter(ctx context.Context,
	set exporter.Settings,
	cfg component.Config,
) (exporter.Traces, error) {
	cf := cfg.(*Config)
	exporter, err := newExporter(cf, set)
	if err != nil {
		return nil, err
	}

	return exporterhelper.NewTraces(
		ctx,
		set,
		cfg,
		exporter.pushTraceData,
		exporterhelperOptions(cf, exporter.Start, exporter.Shutdown)...,
	)
}

func exporterhelperOptions(
	cfg *Config,
	start component.StartFunc,
	shutdown component.ShutdownFunc,
) []exporterhelper.Option {
	opts := []exporterhelper.Option{
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		exporterhelper.WithStart(start),
		exporterhelper.WithShutdown(shutdown),
		exporterhelper.WithQueue(cfg.QueueSettings),
	}
	// if cfg.Batcher.Enabled != nil {
	// 	batcherConfig := exporterbatcher.Config{
	// 		Enabled:       *cfg.Batcher.Enabled,
	// 		FlushTimeout:  cfg.Batcher.FlushTimeout,
	// 		MinSizeConfig: cfg.Batcher.MinSizeConfig,
	// 		MaxSizeConfig: cfg.Batcher.MaxSizeConfig,
	// 	}
	// 	opts = append(opts, exporterhelper.WithBatcher(batcherConfig))

	// 	// Effectively disable timeout_sender because timeout is enforced in bulk indexer.
	// 	//
	// 	// We keep timeout_sender enabled in the async mode (Batcher.Enabled == nil),
	// 	// to ensure sending data to the background workers will not block indefinitely.
	// 	opts = append(opts, exporterhelper.WithTimeout(exporterhelper.TimeoutConfig{Timeout: 0}))
	// }
	return opts
}
