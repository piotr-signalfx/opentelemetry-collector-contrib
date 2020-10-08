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
	"testing"
	"time"

	metricspb "github.com/census-instrumentation/opencensus-proto/gen-go/metrics/v1"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/testbed/testbed"
	"go.uber.org/zap"
)

type fakeSapClient struct {
}

func (c *fakeSapClient) get() (metrics []*metricspb.Metric, err error) {
	return make([]*metricspb.Metric, 0), nil
}

func TestMetricsReceiver(t *testing.T) {

	cfg := &Config{
		CollectionInterval: 1 * time.Second,
	}

	metricsReceiver := newReceiver(
		zap.NewNop(), cfg,
		&fakeSapClient{},
		&testbed.MockMetricConsumer{},
	)
	require.NotNil(t, metricsReceiver)
	ctx := context.Background()
	err := metricsReceiver.Start(ctx, componenttest.NewNopHost())
	require.NoError(t, err)
	err = metricsReceiver.Shutdown(ctx)
	require.NoError(t, err)
}
