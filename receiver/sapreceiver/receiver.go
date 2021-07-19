// Copyright 2020, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sapreceiver

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/redisreceiver/interval"
)

type sapReceiver struct {
	config         *Config
	logger         *zap.Logger
	consumer       consumer.MetricsConsumer
	intervalRunner *interval.Runner
	client         sapClient
}

func newReceiver(logger *zap.Logger, cfg *Config, client sapClient, consumer consumer.MetricsConsumer) *sapReceiver {
	return &sapReceiver{config: cfg, logger: logger, consumer: consumer, client: client}
}

// Set up and kick off the interval runner.
func (r *sapReceiver) Start(ctx context.Context, host component.Host) error {
	runnable := newRunnable(ctx, r.logger, r.client, r.consumer)
	r.intervalRunner = interval.NewRunner(r.config.CollectionInterval, runnable)

	go func() {
		if err := r.intervalRunner.Start(); err != nil {
			host.ReportFatalError(err)
		}
	}()

	return nil
}

func (r *sapReceiver) Shutdown(ctx context.Context) error {
	r.intervalRunner.Stop()
	return nil
}
