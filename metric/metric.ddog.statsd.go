package metric

import (
	"github.com/DataDog/datadog-go/v5/statsd"
)

type ddogStatsd struct {
	client *statsd.Client
}

func newDdogStatsd(opts map[string]string) (Metric, error) {
	statsd, err := statsd.New(opts["host"])
	if err != nil {
		return nil, err
	}
	return &ddogStatsd{
		client: statsd,
	}, nil
}

func (d *ddogStatsd) Count(metric string, value float64, tags ...string) error {
	return d.client.Count(metric, int64(value), tags, 1)
}

func (d *ddogStatsd) Gauge(metric string, value float64, tags ...string) error {
	return d.client.Gauge(metric, value, tags, 1)
}

func (d *ddogStatsd) Histogram(metric string, value float64, tags ...string) error {
	return d.client.Histogram(metric, value, tags, 1)
}

func (d *ddogStatsd) Registry(name string, metric interface{}) error {
	return nil
}

func (d *ddogStatsd) Client() interface{} {
	return d.client
}
