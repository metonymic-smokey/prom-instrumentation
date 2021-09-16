package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	for {
		dat, err := getTempData()
		if err != nil {
			fmt.Println(err)
			time.Sleep(2 * time.Second)
			continue
		}
		if len(dat.Data.Timestep) == 0 {
			continue
		}

		for _, interval := range dat.Data.Timestep[0].TempVal {
			jobsInQueue.Set(interval.Values.Temp)
		}

		opsProcessed.Inc()
		time.Sleep(2 * time.Second)
	}
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
	go recordMetrics()

	router := mux.NewRouter()
	router.Use(prometheusMiddleware)

	//http.Handle("/metrics", promhttp.Handler())
	router.Path("/metrics").Handler(promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
