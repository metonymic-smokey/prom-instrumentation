package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type CityStats struct {
	Zone string
}

//temperature by country
func (c *CityStats) TemperatureAndHumidity() (
	tempByCity map[string]float64, humidityByCity map[string]float64,
) {
	// get real time API temp data here
	tempByCity = make(map[string]float64)
	dat, err := getTempData()
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(dat.Data.Timestep) == 0 {
		fmt.Println("empty result!")
		return
	}

	cities := []string{"bangalore", "london"}

	for ind, interval := range dat.Data.Timestep[0].TempVal {
		tempByCity[cities[ind%2]] = interval.Values.Temp
	}

	humidityByCity = map[string]float64{
		"bangalore": rand.Float64(),
		"london":    rand.Float64(),
	}
	return
}

type CityStatsCollector struct {
	CityStats *CityStats
}

var (
	tempDesc = prometheus.NewDesc(
		"temperature_city_fahrenheit",
		"temperature of a city in fahrenheit",
		[]string{"city"}, nil,
	)
	humidityDesc = prometheus.NewDesc(
		"humidity_city_fraction",
		"humidity of a city as a fraction",
		[]string{"city"}, nil,
	)
	httpDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "http_response_time_seconds",
		Help:    "Duration of HTTP requests.",
		Buckets: prometheus.LinearBuckets(20, 5, 5),
	})
)

func (cc CityStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(cc, ch)
}

func (cc CityStatsCollector) Collect(ch chan<- prometheus.Metric) {
	begin := time.Now()
	tempByCity, humidityByCity := cc.CityStats.TemperatureAndHumidity()
	duration := time.Since(begin)
	httpDuration.Observe(float64(duration))

	for city, temp := range tempByCity {
		ch <- prometheus.MustNewConstMetric(
			tempDesc,
			prometheus.CounterValue,
			float64(temp),
			city,
		)
	}
	for city, humidity := range humidityByCity {
		ch <- prometheus.MustNewConstMetric(
			humidityDesc,
			prometheus.GaugeValue,
			humidity,
			city,
		)
	}
}

func NewCityStats(zone string, reg prometheus.Registerer) *CityStats {
	c := &CityStats{
		Zone: zone,
	}
	cc := CityStatsCollector{CityStats: c}
	prometheus.WrapRegistererWith(prometheus.Labels{"zone": zone}, reg).MustRegister(cc)
	return c
}

func main() {
	reg := prometheus.NewRegistry()

	NewCityStats("db", reg)
	NewCityStats("ca", reg)

	// Add the standard process and Go metrics to the custom registry.
	reg.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
		httpDuration,
	)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(":2112", nil))
}
