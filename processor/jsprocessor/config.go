// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package jsprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/jsprocessor"

// Config holds the configuration for the JS processor.
type Config struct{}

func (*Config) Validate() error {
	return nil
}
