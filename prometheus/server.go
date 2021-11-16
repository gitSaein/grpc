package main

import (
	"log"
	"net/http"

	"contrib.go.opencensus.io/exporter/prometheus"
)

func main() {
	exporter, err := prometheus.NewExporter(prometheus.Options{
		Namespace: "demo",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Serve the scrape endpoint on port 9999.
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", exporter)
		if err := http.ListenAndServe(":9090", mux); err != nil {
			log.Fatalf("Failed to run Prometheus /metrics endpoint: %v", err)
		}
	}()
}
