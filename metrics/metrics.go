package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	promclient "github.com/prometheus/client_model/go"
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

// GetMetricValue read counter value from a prometheus collector
// kudos https://stackoverflow.com/questions/57952695/prometheus-counters-how-to-get-current-value-with-golang-client
func GetMetricValue(collector prometheus.Counter) float64 {
	var m promclient.Metric
	c := make(chan prometheus.Metric, 1)
	collector.Collect(c)
	_ = (<-c).Write(&m)
	return *m.Counter.Value
}
