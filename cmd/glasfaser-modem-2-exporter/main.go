package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/strayer/glasfaser-modem-2-exporter/internal/collector"
	"github.com/strayer/glasfaser-modem-2-exporter/internal/metrics"
	"github.com/strayer/glasfaser-modem-2-exporter/internal/status"
)

const defaultONTURL = "http://192.168.100.1/ONT/client/data/Status.json"

func main() {
	listenAddr := flag.String("listen-address", ":9100", "Address to listen on for web interface and telemetry")
	ontURL := flag.String("ont-url", defaultONTURL, "URL of the ONT status JSON endpoint")
	timeout := flag.Duration("timeout", 10*time.Second, "Timeout for ONT HTTP requests")

	flag.Parse()

	handler := newMetricsHandlerWithTimeout(*ontURL, *timeout)

	http.Handle("/metrics", handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
<head><title>Glasfaser Modem 2 Exporter</title></head>
<body>
<h1>Glasfaser Modem 2 Exporter</h1>
<p><a href="/metrics">Metrics</a></p>
</body>
</html>`))
	})

	log.Printf("Listening on %s", *listenAddr)

	server := &http.Server{
		Addr:              *listenAddr,
		ReadHeaderTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting server: %v\n", err)
		os.Exit(1)
	}
}

// newMetricsHandler creates a metrics handler with default timeout.
func newMetricsHandler(ontURL string) http.Handler {
	return newMetricsHandlerWithTimeout(ontURL, 10*time.Second)
}

// newMetricsHandlerWithTimeout creates a handler that scrapes the ONT and exposes Prometheus metrics.
func newMetricsHandlerWithTimeout(ontURL string, timeout time.Duration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		registry := prometheus.NewRegistry()

		// Scrape error metric
		scrapeError := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "ont_scrape_error",
			Help: "1 if there was an error scraping the ONT, 0 otherwise",
		})
		registry.MustRegister(scrapeError)

		// Create HTTP client with timeout
		client := &http.Client{Timeout: timeout}

		// The metrics function that fetches from ONT
		metricsFunc := func() (metrics.Metrics, error) {
			req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, ontURL, nil)
			if err != nil {
				return metrics.Metrics{}, err
			}

			req.Header.Set("Accept-Language", "en")

			resp, err := client.Do(req)
			if err != nil {
				return metrics.Metrics{}, err
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return metrics.Metrics{}, fmt.Errorf("ONT returned status %d", resp.StatusCode)
			}

			data, err := status.Parse(resp.Body)
			if err != nil {
				return metrics.Metrics{}, err
			}

			return metrics.Extract(data), nil
		}

		// Create and register collector
		coll := collector.New(metricsFunc)

		// Try to get metrics to check for errors
		_, err := metricsFunc()
		if err != nil {
			log.Printf("ERROR: scrape failed: %v", err)
			scrapeError.Set(1)
		} else {
			scrapeError.Set(0)
			registry.MustRegister(coll)
		}

		// Serve metrics
		promhttp.HandlerFor(registry, promhttp.HandlerOpts{}).ServeHTTP(w, r)
	})
}
