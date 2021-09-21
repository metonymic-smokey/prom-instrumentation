package main

import (
	"fmt"
	"net/http"
	"time"

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
			tempCelsius.Set(interval.Values.Temp)
		}

		opsProcessed.Inc()
		time.Sleep(2 * time.Second)
	}
}

var opsProcessed = promauto.NewGauge(
	prometheus.GaugeOpts{
		Name: "processed_ops_total",
		Help: "The total number of processed operations",
	},
)

var tempCelsius = promauto.NewGauge(
	prometheus.GaugeOpts{
		Name: "current_temperature_api_celsius",
		Help: "Current temperature",
	},
)

func main() {
	go recordMetrics()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
