package metric

import "errors"

type AgentMetric string

const (
	Datadog       AgentMetric = "datadog"
	DatadogStatsd AgentMetric = "datadog-statsd"
	DatadogNoop   AgentMetric = "datadog-noop"
	Prome         AgentMetric = "prometheus"
)

// Metrics interface
type Metric interface {
	// Count metric for increment/decrement
	Count(metric string, value float64, tags ...string) error
	// Gauge metric for set value
	Gauge(metric string, value float64, tags ...string) error
	// Histogram metric for set value
	Histogram(metric string, value float64, tags ...string) error
	// Registry metric for register metric
	Registry(name string, metric interface{}) error
	// Client return client
	Client() interface{}
}

func New(
	agent AgentMetric,
	opts map[string]string,
	resource map[string]string,
) (Metric, error) {

	if agent == DatadogNoop {
		return newNoopMetrics()
	}

	if agent == Datadog {
		return newDatadogMetrics(opts, resource)
	}

	if agent == DatadogStatsd {
		return newDdogStatsd(opts)
	}

	if agent == Prome {
		return newProme(opts)
	}

	return nil, errors.New("invalid agent")
}
