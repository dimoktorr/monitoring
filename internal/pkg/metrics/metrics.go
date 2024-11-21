package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	requestCostCounter  *prometheus.CounterVec
	startPayCostCounter *prometheus.CounterVec
	totalPayCostCounter *prometheus.CounterVec
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

	startPayCostCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "start_pay_total",
			Help: "Кол-во запросов на оплату",
		},
		[]string{"status"},
	)
	prometheus.MustRegister(startPayCostCounter)

	totalPayCostCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "total_pay_total",
			Help: "Сумма оплаченных продуктов",
		},
		[]string{},
	)
	prometheus.MustRegister(totalPayCostCounter)

	return &Metrics{
		requestCostCounter:  requestCostCounter,
		startPayCostCounter: startPayCostCounter,
		totalPayCostCounter: totalPayCostCounter,
	}
}

func (m *Metrics) IncGetProductSumCounter() {
	m.requestCostCounter.WithLabelValues().Inc()
}

func (m *Metrics) IncPayProductSumCounter(status string) {
	m.startPayCostCounter.WithLabelValues(status).Inc()
}

func (m *Metrics) AddAmountPayProduct(amount float64) {
	m.totalPayCostCounter.WithLabelValues().Add(amount)
}
