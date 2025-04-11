package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	createdPVZ = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "created_pvz_total",
		Help: "total number created pvz",
	}, []string{})

	createdProducts = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "created_products_total",
		Help: "total number created products",
	}, []string{})

	createdReceptions = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "created_receptions_total",
		Help: "total number created receptions",
	}, []string{})
)

func IncCreatedPVZ() {
	createdPVZ.With(prometheus.Labels{}).Inc()
}

func IncCreatedReceptions() {
	createdReceptions.With(prometheus.Labels{}).Inc()
}

func IncCreatedProducts() {
	createdProducts.With(prometheus.Labels{}).Inc()
}
