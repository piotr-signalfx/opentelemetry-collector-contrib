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

	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/consumerdata"
	"go.opentelemetry.io/collector/obsreport"
	"go.opentelemetry.io/collector/translator/internaldata"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/redisreceiver/interval"
)

var _ interval.Runnable = (*runnable)(nil)

// Runs intermittently, fetching info from Redis, creating metrics/datapoints,
// and feeding them to a metricsConsumer.
type runnable struct {
	ctx      context.Context
	client   sapClient
	logger   *zap.Logger
	consumer consumer.MetricsConsumer
	//	serviceName     string
}

func newRunnable(
	ctx context.Context,
	logger *zap.Logger,
	client sapClient,
	consumer consumer.MetricsConsumer,

) *runnable {
	return &runnable{
		ctx:      ctx,
		logger:   logger,
		client:   client,
		consumer: consumer,
	}
}

func (r *runnable) Setup() error {
	return nil
}

// Run is called periodically, querying SAP API and building Metrics to send to
// the next consumer.
func (r *runnable) Run() error {
	r.logger.Info("SAP Receiver executed...")

	const dataFormat = "sap"
	const transport = "http" // todo verify this
	ctx := obsreport.StartMetricsReceiveOp(r.ctx, dataFormat, transport)

	metrics, err := r.client.get()
	if err != nil {
		r.logger.Error("Invalid response received", zap.Error(err))
		obsreport.EndMetricsReceiveOp(ctx, dataFormat, 0, 0, err)
	} else {
		md := consumerdata.MetricsData{
			Metrics: metrics,
		}
		numTimeSeries, numPoints := obsreport.CountMetricPoints(md)
		err := r.consumer.ConsumeMetrics(ctx, internaldata.OCToMetrics(md))
		obsreport.EndMetricsReceiveOp(ctx, dataFormat, numPoints, numTimeSeries, err)
	}
	return nil
}
