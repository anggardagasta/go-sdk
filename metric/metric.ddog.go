package metric

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
)

type datadogMetrics struct {
	context   context.Context
	client    *datadog.APIClient
	resources map[string]string
}

func newDatadogMetrics(params map[string]string, resources map[string]string) (Metric, error) {
	config := datadog.NewConfiguration()
	config.Debug = true
	client := datadog.NewAPIClient(config)

	apiKey, ok := params["apiKey"]
	if !ok {
		return nil, errors.New("apiKey is required")
	}

	appKey, ok := params["appKey"]
	if !ok {
		return nil, errors.New("appKey is required")
	}

	ctx := context.WithValue(
		context.Background(),
		datadog.ContextAPIKeys,
		map[string]datadog.APIKey{
			"apiKeyAuth": {
				Key: apiKey,
			},
			"appKeyAuth": {
				Key: appKey,
			},
		},
	)

	return &datadogMetrics{
		context:   ctx,
		client:    client,
		resources: map[string]string{},
	}, nil
}

func (d *datadogMetrics) getGlobalResource() []datadogV2.MetricResource {
	var res []datadogV2.MetricResource
	for kind, tag := range d.resources {
		res = append(res, datadogV2.MetricResource{
			Name: datadog.PtrString(tag),
			Type: datadog.PtrString(kind),
		})
	}
	return res
}

// Count increments a counter metric
func (d *datadogMetrics) Count(metric string, value float64, tags ...string) error {
	body := datadogV2.MetricPayload{
		Series: []datadogV2.MetricSeries{
			{
				Metric: metric,
				Type:   datadogV2.METRICINTAKETYPE_COUNT.Ptr(),
				Points: []datadogV2.MetricPoint{
					{
						Timestamp: datadog.PtrInt64(time.Now().Unix()),
						Value:     datadog.PtrFloat64(value),
					},
				},
				Resources: d.getGlobalResource(),
				Tags:      tags,
			},
		},
	}
	api := datadogV2.NewMetricsApi(d.client)
	_, _, err := api.SubmitMetrics(d.context, body, *datadogV2.NewSubmitMetricsOptionalParameters())
	if err != nil {
		return err
	}

	return nil
}

// Gauge sets or updates a gauge metric
func (d *datadogMetrics) GaugeV1(metric string, value float64, tags ...string) error {
	body := datadogV2.MetricPayload{
		Series: []datadogV2.MetricSeries{
			{
				Metric: metric,
				Type:   datadogV2.METRICINTAKETYPE_GAUGE.Ptr(),
				Points: []datadogV2.MetricPoint{
					{
						Timestamp: datadog.PtrInt64(time.Now().Unix()),
						Value:     datadog.PtrFloat64(value),
					},
				},
				Tags: tags,
			},
		},
	}
	api := datadogV2.NewMetricsApi(d.client)
	_, _, err := api.SubmitMetrics(d.context, body, *datadogV2.NewSubmitMetricsOptionalParameters())
	if err != nil {
		fmt.Println("[DEBUG] ERR: ", err)
		return err
	}

	return nil
}

// Gauge sets or updates a gauge metric
func (d *datadogMetrics) Gauge(metric string, value float64, tags ...string) error {
	body := datadogV1.MetricsPayload{
		Series: []datadogV1.Series{
			{
				Metric: metric,
				Type:   datadog.PtrString("gauge"),
				Points: [][]*float64{
					{
						datadog.PtrFloat64(float64(time.Now().Unix())),
						datadog.PtrFloat64(value),
					},
				},
				Tags: tags,
			},
		},
	}
	api := datadogV1.NewMetricsApi(d.client)
	_, _, err := api.SubmitMetrics(d.context, body, *datadogV1.NewSubmitMetricsOptionalParameters())
	if err != nil {
		return err
	}

	return nil
}

// Histogram sets or updates a histogram metric
func (d *datadogMetrics) Histogram(metric string, value float64, tags ...string) error {
	body := datadogV2.MetricPayload{
		Series: []datadogV2.MetricSeries{
			{
				Metric: metric,
				Type:   datadogV2.METRICINTAKETYPE_GAUGE.Ptr(),
				Points: []datadogV2.MetricPoint{
					{
						Timestamp: datadog.PtrInt64(time.Now().Unix()),
						Value:     datadog.PtrFloat64(value),
					},
				},
				Resources: d.getGlobalResource(),
				Tags:      tags,
			},
		},
	}
	api := datadogV2.NewMetricsApi(d.client)
	_, _, err := api.SubmitMetrics(d.context, body, *datadogV2.NewSubmitMetricsOptionalParameters())
	if err != nil {
		return err
	}

	return nil
}

func (d *datadogMetrics) Registry(name string, metric interface{}) error {
	return nil
}

func (d *datadogMetrics) Client() interface{} {
	return d.client
}
