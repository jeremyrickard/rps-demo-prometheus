package main

import (
	"log"

	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metrics struct {
	currentCount int
	lastCount    int
	rps          float64
	lastTime     time.Time
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
			now := time.Now()
			t := prometheus.NewTimer(requestDurationsHistogram)
			handler.ServeHTTP(w, r)
			defer t.ObserveDuration()
			diff := time.Since(now)
			log.Printf("Finished request : %v", diff.Seconds())
		},
	)
}

func main() {
	sleepSecondsStr := os.Getenv("SLEEP_SECONDS")
	sleepSeconds, err := strconv.Atoi(sleepSecondsStr)
	if err != nil {
		log.Fatalf("bad value for sleep seconds: %s", sleepSecondsStr)
	}

	rpsLimitStr := os.Getenv("RPS_THRESHOLD")
	rpsLimit, err := strconv.ParseFloat(rpsLimitStr, 64)
	if err != nil {
		log.Fatalf("bad value for rps limit: %s", rpsLimitStr)
	}

	ctr := &metrics{
		lastTime: time.Now(),
	}
	throttledHandler := throttler(
		ctr,
		rpsLimit,
		sleepSeconds,
		http.FileServer(http.Dir("/app/content")),
	)
	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/", instrumentHandler(throttledHandler))
	go rpsTelemetryCalculator(ctr)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func throttler(
	ctr *metrics,
	limit float64,
	sleepSeconds int,
	handler http.Handler,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if ctr.rps > limit {
				log.Printf("throttle")
				randomSleep := time.Duration((float64(sleepSeconds*1000) + (ctr.rps / limit))) * time.Millisecond
				log.Printf("random sleep: %v", randomSleep)
				time.Sleep(randomSleep)
			}
			// lets make all the requests take at least 5 ms
			time.Sleep(5 * time.Millisecond)
			handler.ServeHTTP(w, r)
			ctr.currentCount++
		},
	)
}

func rpsTelemetryCalculator(
	ctr *metrics,
) {
	tickChan := time.NewTicker(time.Second * 30).C
	var count int
	var t time.Time
	for {
		t = <-tickChan
		duration := t.Sub(ctr.lastTime)
		count = ctr.currentCount
		ctr.rps = float64(count-ctr.lastCount) / duration.Seconds()
		ctr.lastTime = t
		ctr.lastCount = count
		log.Printf("RPS: %v", ctr.rps)
	}
}
