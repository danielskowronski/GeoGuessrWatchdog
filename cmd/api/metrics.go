package main

import (
	"net/http"

	db "github.com/danielskowronski/geoguessrwatchdog/internal/db/generated"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const fetchStatusMetricsMaxAgeDays = 7

type Metrics struct {
	fetchLastSuccess *prometheus.GaugeVec
	handler          http.Handler
}

func NewMetrics(app *App) *Metrics {
	registry := prometheus.NewRegistry()
	fetchLastSuccess := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ggwd_fetch_last_success_timestamp_seconds",
			Help: "Unix timestamp of the last successful fetch.",
		},
		[]string{"fetch_type"},
	)
	registry.MustRegister(fetchLastSuccess)

	promHandler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

	return &Metrics{
		fetchLastSuccess: fetchLastSuccess,
		handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := app.RefreshMetrics(r); err != nil {
				app.logger.Warn("failed to refresh metrics", "err", err)
				http.Error(w, "failed to refresh metrics", http.StatusInternalServerError)
				return
			}

			promHandler.ServeHTTP(w, r)
		}),
	}
}

func (a *App) RefreshMetrics(r *http.Request) error {
	q := db.New(a.db)
	fetchStatuses, err := q.GetAllStatuses(r.Context(), fetchStatusMetricsMaxAgeDays)
	if err != nil {
		return err
	}

	a.metrics.fetchLastSuccess.Reset()
	for _, row := range fetchStatuses {
		a.metrics.fetchLastSuccess.WithLabelValues(row.FetchType).Set(float64(row.LastSuccess.Time.Unix()))
	}

	return nil
}

func (m *Metrics) Handler() http.Handler {
	return m.handler
}
