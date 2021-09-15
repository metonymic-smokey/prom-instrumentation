package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics(temp float64) {
	go func() {
		for {
			opsProcessed.Inc()
			jobsInQueue.Set(temp)
			time.Sleep(2 * time.Second)
		}
	}()
}

var opsProcessed = promauto.NewGauge(
	prometheus.GaugeOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	},
)

var jobsInQueue = promauto.NewGauge(
	prometheus.GaugeOpts{
		Name: "current_temperature_api_celsius",
		Help: "Current temperature",
	},
)

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path"})

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))

		timer.ObserveDuration()
	})
}

func main() {

	dat := getTempData()

	for _, interval := range dat.Data.Timestep[0].TempVal {
		recordMetrics(interval.Values.Temp)
	}

	router := mux.NewRouter()
	router.Use(prometheusMiddleware)

	//http.Handle("/metrics", promhttp.Handler())
	router.Path("/metrics").Handler(promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
