// Package prometheus provides Prometheus implementations for metrics.
// Individual metrics are mapped to their Prometheus counterparts, and
// (depending on the constructor used) may be automatically registered in the
// global Prometheus metrics registry.
package core

//// Counter describes a metric that accumulates values monotonically.
//// An example of a counter is the number of received HTTP requests.
type ICounter interface {
	With(labelValues ...string) ICounter
	Add(delta float64)
}

// Gauge describes a metric that takes specific values over time.
// An example of a gauge is the current depth of a job queue.
type IGauge interface {
	With(labelValues ...string) IGauge
	Set(value float64)
	Add(delta float64)
}

// Histogram describes a metric that takes repeated observations of the same
// kind of thing, and produces a statistical summary of those observations,
// typically expressed as quantiles or buckets. An example of a histogram is
// HTTP request latencies.
type IHistogram interface {
	With(labelValues ...string) IHistogram
	Observe(value float64)
}
