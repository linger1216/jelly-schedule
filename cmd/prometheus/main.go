package main

import (
	"fmt"
	"github.com/linger1216/jelly-schedule/core"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func main() {
	fieldKeys := []string{"method"}
	progress := core.NewGaugeFrom(prometheus.GaugeOpts{
		Namespace: "lid",
		Subsystem: "progress",
		Name:      "request_latency_seconds",
		Help:      "Total duration of requests in seconds.",
	}, fieldKeys)

	//latency := core.NewHistogramFrom(prometheus.HistogramOpts{
	//	Namespace: "lid",
	//	Subsystem: "latency",
	//	Name:      "request_latency_seconds",
	//	Help:      "Total duration of requests in seconds.",
	//}, fieldKeys)

	progress.With("method", "progress1").Add(1)
	errc := make(chan error)
	go func() {
		m := http.NewServeMux()
		m.Handle("/metrics", promhttp.Handler())
		errc <- http.ListenAndServe(":9999", m)
	}()
	fmt.Println("location.thing.v1.AssetServer", <-errc)
}
