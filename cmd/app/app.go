package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metrics struct {
	lastResponseTime float64
	currentCount     int
	lastCount        int
	lastTime         time.Time
}

var (
	requestDurationsHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "request_durations_histogram_secs",
		Buckets: prometheus.DefBuckets,
		Help:    "Requests Durations, in Seconds",
	})
)

func init() {
	prometheus.MustRegister(requestDurationsHistogram)
}

func instrumentHandler(
	handler http.Handler,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			t := prometheus.NewTimer(requestDurationsHistogram)
			handler.ServeHTTP(w, r)
			t.ObserveDuration()
		},
	)
}

func main() {
	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/", instrumentHandler(http.FileServer(http.Dir("/app/content"))))
	log.Fatal(http.ListenAndServe(":8080", nil))

}
