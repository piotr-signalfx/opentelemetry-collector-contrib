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
	"io/ioutil"
	"net/http"
	"time"

	metricspb "github.com/census-instrumentation/opencensus-proto/gen-go/metrics/v1"
	"go.opentelemetry.io/collector/consumer/consumererror"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type sapClient interface {
	get() (metrics []*metricspb.Metric, err error)
}

type sapClientImpl struct {
	url     string
	client  *http.Client
	logger  *zap.Logger
	headers map[string]string
}

type sapMetric struct {
	name        string
	value       string
	metricType  metricspb.MetricDescriptor_Type
	sampleRate  float64
	labelKeys   []*metricspb.LabelKey
	labelValues []*metricspb.LabelValue
}

func buildClient(url string, logger *zap.Logger) sapClient {
	return &sapClientImpl{
		url: url,
		client: &http.Client{
			Transport: &http.Transport{
				Proxy:              http.ProxyFromEnvironment,
				DisableCompression: true, // needed just to make resp.ContentLength return some useful value
			},
		},
		logger: logger,
		headers: map[string]string{
			"Connection": "keep-alive",
			"Accept":     "*/*",
			"User-Agent": "OpenTelemetry-Collector SAP Receiver/v0.0.1",
		},
	}
}

func (c *sapClientImpl) get() (m []*metricspb.Metric, err error) {

	req, err := http.NewRequest("GET", c.url, nil)
	if err != nil {
		return nil, consumererror.Permanent(err)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf(
			"HTTP %d %q",
			resp.StatusCode,
			http.StatusText(resp.StatusCode))
		return nil, err
	}

	var metrics []*metricspb.Metric

	metrics = append(metrics, buildResponseSizeMetric(resp.ContentLength))
	return metrics, nil
}

func buildResponseSizeMetric(size int64) *metricspb.Metric {

	labelKeys := make([]*metricspb.LabelKey, 0)
	labelValues := make([]*metricspb.LabelValue, 0)

	labelKeys = append(labelKeys, &metricspb.LabelKey{Key: "label1"})
	labelValues = append(labelValues, &metricspb.LabelValue{
		Value:    "A1",
		HasValue: true,
	})

	labelKeys = append(labelKeys, &metricspb.LabelKey{Key: "label2"})
	labelValues = append(labelValues, &metricspb.LabelValue{
		Value:    "L2",
		HasValue: true,
	})

	point := buildPoint(size)
	return buildMetric("httpResponseSize", metricspb.MetricDescriptor_GAUGE_INT64, labelKeys, labelValues, point)
}

func buildMetric(name string, metricType metricspb.MetricDescriptor_Type, labelKeys []*metricspb.LabelKey,
	labelValues []*metricspb.LabelValue, point *metricspb.Point) *metricspb.Metric {
	return &metricspb.Metric{
		MetricDescriptor: &metricspb.MetricDescriptor{
			Name:      name,
			Type:      metricType,
			LabelKeys: labelKeys,
		},
		Timeseries: []*metricspb.TimeSeries{
			{
				LabelValues: labelValues,
				Points: []*metricspb.Point{
					point,
				},
			},
		},
	}
}

var timeNowFunc = func() int64 {
	return time.Now().Unix()
}

func buildPoint(value int64) *metricspb.Point {
	now := &timestamppb.Timestamp{
		Seconds: timeNowFunc(),
	}
	point := &metricspb.Point{
		Timestamp: now,
		Value: &metricspb.Point_Int64Value{
			Int64Value: value,
		},
		// Double:
		// Value: &metricspb.Point_DoubleValue{
		// 	DoubleValue: f,
		// },
	}
	return point
}
