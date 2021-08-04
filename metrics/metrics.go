package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Collector is the main prometheus data gatherer
var Collector = &Gatherer{
	AllowListSize: promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "log_reader",
		Subsystem: "allow_list",
		Name:      "size",
		Help:      "Current size of the allow list",
	}),
	AuthRequestCounter: promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "log_reader",
		Subsystem: "auth_requests",
		Name:      "total",
		Help:      "Count all auth requests",
	}),
	StatusOkRespCounter: promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "log_reader",
		Subsystem: "auth_responses",
		Name:      "status_ok",
		Help:      "Count all auth OK (200) responses",
	}),
	StatusForbiddenRespCounter: promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "log_reader",
		Subsystem: "auth_responses",
		Name:      "status_forbidden",
		Help:      "Count all auth Forbidden (403) responses",
	})}

// Gatherer contains every metric stat use by the prometheus endpoint
type Gatherer struct {
	AllowListSize              prometheus.Gauge
	AuthRequestCounter         prometheus.Counter
	StatusOkRespCounter        prometheus.Counter
	StatusForbiddenRespCounter prometheus.Counter
}
