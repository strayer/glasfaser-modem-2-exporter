package main

import (
	"context"
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
	"github.com/strayer/glasfaser-modem-2-exporter/internal/web"
)

const defaultONTURL = "http://192.168.100.1/ONT/client/data/Status.json"

func main() {
	listenAddr := flag.String("listen-address", ":9100", "Address to listen on for web interface and telemetry")
	ontURL := flag.String("ont-url", defaultONTURL, "URL of the ONT status JSON endpoint")
	timeout := flag.Duration("timeout", 10*time.Second, "Timeout for ONT HTTP requests")

	flag.Parse()

	// Create shared metrics function
	metricsFunc := newMetricsFunc(*ontURL, *timeout)

	// Register handlers
	http.Handle("/metrics", newMetricsHandler(*ontURL, *timeout, metricsFunc))
	http.Handle("/api/status", web.APIHandler(metricsFunc))
	http.Handle("/", web.DashboardHandler(metricsFunc))

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

// newMetricsFunc creates a function that fetches metrics from the ONT.
func newMetricsFunc(ontURL string, timeout time.Duration) func() (metrics.Metrics, error) {
	client := &http.Client{Timeout: timeout}

	return func() (metrics.Metrics, error) {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, ontURL, nil)
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
}

// newMetricsHandler creates a handler that exposes Prometheus metrics.
func newMetricsHandler(_ string, _ time.Duration, metricsFunc func() (metrics.Metrics, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		registry := prometheus.NewRegistry()

		// Scrape error metric
		scrapeError := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "ont_scrape_error",
			Help: "1 if there was an error scraping the ONT, 0 otherwise",
		})
		registry.MustRegister(scrapeError)

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
