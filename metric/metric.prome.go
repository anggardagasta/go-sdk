package metric

import (
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

type prome struct {
	client  *prometheus.Registry
	metrics map[string]interface{}
}

func newProme(opts map[string]string) (Metric, error) {
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewBuildInfoCollector())

	return &prome{
		client:  reg,
		metrics: map[string]interface{}{},
	}, nil
}

// Count metric for increment/decrement
func (d *prome) Count(metric string, value float64, tags ...string) error {
	if met, ok := d.metrics[metric]; ok {
		if counter, ok := met.(prometheus.Counter); ok {
			counter.Add(value)
			return nil
		}
		return errors.New("metric name not match with counter")
	}
	return errors.New("Count metric not found")
}

// Gauge metric for set value (float64)
func (d *prome) Gauge(metric string, value float64, tags ...string) error {
	if met, ok := d.metrics[metric]; ok {
		if gauge, ok := met.(prometheus.Gauge); ok {
			gauge.Set(value)
			return nil
		}
		return errors.New("metric name not match with gauge")
	}
	return errors.New("Gauge metric not found")
}

// Histogram metric for set value (float64)
func (d *prome) Histogram(metric string, value float64, tags ...string) error {
	if met, ok := d.metrics[metric]; ok {
		if histogram, ok := met.(prometheus.Histogram); ok {
			histogram.Observe(value)
			return nil
		}
		return errors.New("metric name not match with histogram")
	}
	return errors.New("Histogram metric not found")
}

// Registry metric for register metric
func (d *prome) Registry(name string, metric interface{}) error {
	var collector prometheus.Collector
	switch v := metric.(type) {
	case prometheus.Counter:
		collector = v
	case prometheus.Histogram:
		collector = v
	default:
		return errors.New("metric type not found")
	}

	d.client.MustRegister(collector)
	d.metrics[name] = metric
	return nil
}

func (d *prome) Client() interface{} {
	return d.client
}
