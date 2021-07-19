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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configcheck"
	"go.opentelemetry.io/collector/testbed/testbed"
	"go.uber.org/zap"
)

func TestCreateMetricsReceiver(t *testing.T) {
	cfg := createDefaultConfig()
	require.NotNil(t, cfg)

	r, err := createMetricsReceiver(
		context.Background(),
		component.ReceiverCreateParams{Logger: zap.NewNop()},
		cfg,
		&testbed.MockMetricConsumer{},
	)
	require.NoError(t, err)
	require.NotNil(t, r)
}

func TestCreateDefaultConfig(t *testing.T) {
	cfg := createDefaultConfig()
	assert.NotNil(t, cfg, "failed to create default config")
	assert.NoError(t, configcheck.ValidateConfig(cfg))
}

func TestCreateMetricsExporterNoConfig(t *testing.T) {
	_, err := createMetricsReceiver(context.Background(),
		component.ReceiverCreateParams{Logger: zap.NewNop()},
		nil,
		&testbed.MockMetricConsumer{},
	)
	assert.Error(t, err)
}
