package main

import (
	"log"
	"math/rand"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type CityStats struct {
	Zone string
}

//temperature by country
func (c *CityStats) TemperatureAndHumidity() (
	tempByCity map[string]int, humidityByCity map[string]float64,
) {
	tempByCity = map[string]int{
		"bangalore": rand.Intn(100),
		"london":    rand.Intn(1000),
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
)

func (cc CityStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(cc, ch)
}

func (cc CityStatsCollector) Collect(ch chan<- prometheus.Metric) {
	tempByCity, humidityByCity := cc.CityStats.TemperatureAndHumidity()
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
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		prometheus.NewGoCollector(),
	)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(":2112", nil))
}
