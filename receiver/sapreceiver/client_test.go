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
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	metricspb "github.com/census-instrumentation/opencensus-proto/gen-go/metrics/v1"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestClient(t *testing.T) {
	srv := serverMock()
	defer srv.Close()

	fmt.Print(srv.URL)
	client := buildClient(srv.URL+"/server/api/v1/metrics", zap.NewNop())
	metrics, err := client.get()

	require.NoError(t, err)
	require.NotNil(t, metrics)
	require.NotEmpty(t, metrics)

	require.Equal(t, "httpResponseSize", metrics[0].MetricDescriptor.Name)
	require.Equal(t, metricspb.MetricDescriptor_GAUGE_INT64, metrics[0].MetricDescriptor.Type)
	require.NotEmpty(t, metrics[0].Timeseries)
	require.NotEmpty(t, metrics[0].Timeseries[0].Points)
	require.Equal(t, &metricspb.Point_Int64Value{Int64Value: 20}, metrics[0].Timeseries[0].Points[0].Value)

}

func serverMock() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/server/api/v1/metrics", handlerMock)
	srv := httptest.NewServer(handler)
	return srv
}

func handlerMock(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("mock server response"))
}
