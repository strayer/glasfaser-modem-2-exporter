package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMetricsHandler_Success(t *testing.T) {
	// Create a mock ONT server that returns fixture data
	fixtureData, err := os.ReadFile("../../testdata/status_provisioned.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	ontServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(fixtureData)
	}))
	defer ontServer.Close()

	// Create our exporter handler pointed at the mock
	handler := newMetricsHandler(ontServer.URL)

	// Request metrics
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()

	// Check for expected metrics
	expectedMetrics := []string{
		"ont_optical_rx_power_dbm",
		"ont_optical_tx_power_dbm",
		"ont_optical_bip_errors_total",
		"ont_lan_receive_bytes_total",
		"ont_lan_transmit_bytes_total",
		"ont_uptime_seconds",
		"ont_device_info",
	}

	for _, m := range expectedMetrics {
		if !strings.Contains(body, m) {
			t.Errorf("expected metric %q in response", m)
		}
	}
}

func TestMetricsHandler_ONTError(t *testing.T) {
	// Create a mock ONT server that returns an error
	ontServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ontServer.Close()

	handler := newMetricsHandler(ontServer.URL)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// Should still return 200 but with scrape error metric
	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "ont_scrape_error") {
		t.Errorf("expected ont_scrape_error metric in response")
	}
}

func TestMetricsHandler_InvalidJSON(t *testing.T) {
	// Create a mock ONT server that returns invalid JSON
	ontServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("not valid json"))
	}))
	defer ontServer.Close()

	handler := newMetricsHandler(ontServer.URL)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// Should still return 200 but with scrape error metric
	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "ont_scrape_error") {
		t.Errorf("expected ont_scrape_error metric in response")
	}
}

func TestMetricsHandler_ConnectionRefused(t *testing.T) {
	// Use an address that won't respond
	handler := newMetricsHandler("http://127.0.0.1:1")

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// Should still return 200 but with scrape error metric
	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "ont_scrape_error") {
		t.Errorf("expected ont_scrape_error metric in response")
	}
}

func TestMetricsHandler_SendsAcceptLanguageHeader(t *testing.T) {
	var receivedHeader string

	ontServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeader = r.Header.Get("Accept-Language")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[]`))
	}))
	defer ontServer.Close()

	handler := newMetricsHandler(ontServer.URL)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if receivedHeader == "" {
		t.Error("expected Accept-Language header to be sent, but it was empty")
	}
}
