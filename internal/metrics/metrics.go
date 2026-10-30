package metrics

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/AAbbey01/gonet/internal/check"
)

// Metrics holds the Prometheus series gonet exports.
type Metrics struct {
	up      *prometheus.GaugeVec
	latency *prometheus.HistogramVec
}

// New creates the metrics and registers them on reg.
func New(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		up: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "gonet_up",
			Help: "1 if the last check of the target succeeded, 0 if it failed.",
		}, []string{"target", "kind"}),
		latency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "gonet_latency_seconds",
			Help:    "How long each check took.",
			Buckets: prometheus.DefBuckets,
		}, []string{"target", "kind"}),
	}
	reg.MustRegister(m.up, m.latency)
	return m
}

// Record stores one check result.
func (m *Metrics) Record(r check.Result) {
	// TODO: m.up.WithLabelValues(r.Target.Name, r.Target.Kind).Set(1 or 0)
	// TODO: m.latency.WithLabelValues(...).Observe(r.Latency.Seconds())
	score := 0.0
	if r.Up{
		score = 1.0
	}
	m.up.WithLabelValues(r.Target.Name, r.Target.Kind).Set(score)
	m.latency.WithLabelValues(r.Target.Name, r.Target.Kind).Observe(r.Latency.Seconds())
}