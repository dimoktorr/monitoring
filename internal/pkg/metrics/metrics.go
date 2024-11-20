package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	requestCostCounter *prometheus.CounterVec
}

func New() *Metrics {
	requestCostCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "get_request_total",
			Help: "Кол-во запросов на получение продукта",
		},
		[]string{},
	)
	prometheus.MustRegister(requestCostCounter)

	return &Metrics{
		requestCostCounter: requestCostCounter,
	}
}

func (m *Metrics) IncGetProductSumCounter() {
	m.requestCostCounter.WithLabelValues().Inc()
}
