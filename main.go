package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

func main() {

	liveHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "debug_service_live_response_seconds",
		Help:    "Time to respond to liveness check",
		Buckets: []float64{1, 2, 5, 6, 10},
	}, []string{"code"})

	readyHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "debug_service_ready_response_seconds",
		Help:    "Time to respond to readiness check",
		Buckets: []float64{1, 2, 5, 6, 10},
	}, []string{"code"})

	_ = prometheus.Register(liveHistogram)

	http.Handle("/live", measureHealthcheckLatency(liveHistogram))
	http.Handle("/ready", measureHealthcheckLatency(readyHistogram))
	http.Handle("/metrics", promhttp.Handler())
	_ = http.ListenAndServe(":8080", nil)
}

func measureHealthcheckLatency(histogram *prometheus.HistogramVec) http.HandlerFunc {
	time.Sleep(1 * time.Second)
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer r.Body.Close()
		code := 200

		defer func() {
			httpDuration := time.Since(start)
			histogram.WithLabelValues(fmt.Sprintf("%d", code)).Observe(httpDuration.Seconds())
		}()

		response := "ok"
		w.Write([]byte(response))
	}
}
