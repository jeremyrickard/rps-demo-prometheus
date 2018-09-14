package main

import (
	"log"

	"net/http"
	"os"
	"strconv"
	"time"

	"golang.org/x/time/rate"

	"contrib.go.opencensus.io/exporter/ocagent"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// type metrics struct {
// 	currentCount int
// 	lastCount    int
// 	rps          float64
// 	lastTime     time.Time
// }

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

	// Register stats and trace exporters to export the collected data.
	serviceName := os.Getenv("SERVICE_NAME")
	if len(serviceName) == 0 {
		serviceName = "go-app"
	}

	agentHostName := os.Getenv("OCAGENT_TRACE_EXPORTER_ENDPOINT")
	if len(agentHostName) == 0 {
		agentHostName = "localhost"
	}

	exporter, err := ocagent.NewExporter(ocagent.WithInsecure(), ocagent.WithServiceName(serviceName), ocagent.WithAddress(agentHostName))
	if err != nil {
		log.Printf("Failed to create the agent exporter: %v", err)
	}

	trace.RegisterExporter(exporter)

	// Always trace for this demo. In a production application, you should
	// configure this to a trace.ProbabilitySampler set at the desired
	// probability.
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	// ctr := &metrics{
	// 	lastTime: time.Now(),
	// }
	throttledHandler := throttler(
		//	ctr,
		rpsLimit,
		sleepSeconds,
		http.FileServer(http.Dir("/app/content")),
	)
	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/", instrumentHandler(throttledHandler))
	//	go rpsTelemetryCalculator(ctr)
	log.Fatal(http.ListenAndServe(":8080", &ochttp.Handler{Propagation: &tracecontext.HTTPFormat{}}))
}

func throttler(
	//ctr *metrics,
	limit float64,
	sleepSeconds int,
	handler http.Handler,
) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(limit), 10)
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			limiter.Wait(r.Context())
			handler.ServeHTTP(w, r)
		},
	)
}
