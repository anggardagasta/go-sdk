package metric

type noopMetrics struct{}

func newNoopMetrics() (Metric, error) {
	return &noopMetrics{}, nil
}

func (m *noopMetrics) Count(metric string, value float64, tags ...string) error {
	return nil
}

func (m *noopMetrics) Gauge(metric string, value float64, tags ...string) error {
	return nil
}

func (m *noopMetrics) Histogram(metric string, value float64, tags ...string) error {
	return nil
}

func (d *noopMetrics) Registry(name string, metric interface{}) error {
	return nil
}

func (d *noopMetrics) Client() interface{} {
	return nil
}
