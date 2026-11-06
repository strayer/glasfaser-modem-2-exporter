package web_test

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/strayer/glasfaser-modem-2-exporter/internal/metrics"
	"github.com/strayer/glasfaser-modem-2-exporter/internal/web"
)

func TestDashboardHandler_ServesHTML(t *testing.T) {
	metricsFunc := func() (metrics.Metrics, error) {
		return metrics.Metrics{}, nil
	}

	handler := web.DashboardHandler(metricsFunc)
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Errorf("expected Content-Type text/html, got %s", contentType)
	}
}

func TestDashboardHandler_ContainsTitle(t *testing.T) {
	metricsFunc := func() (metrics.Metrics, error) {
		return metrics.Metrics{}, nil
	}

	handler := web.DashboardHandler(metricsFunc)
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	body := w.Body.String()
	if !strings.Contains(body, "Glasfaser Modem 2") {
		t.Error("expected dashboard to contain title 'Glasfaser Modem 2'")
	}
}

func TestAPIHandler_ReturnsJSON(t *testing.T) {
	metricsFunc := func() (metrics.Metrics, error) {
		return metrics.Metrics{
			RxPowerDBm:       -16.46,
			TxPowerDBm:       2.5,
			BIPErrors:        0,
			RxBytes:          123456789,
			TxBytes:          987654321,
			RxPackets:        12345,
			TxPackets:        54321,
			RxDropPackets:    0,
			UptimeSeconds:    123456,
			LinkSpeedMbps:    2500,
			PloamSuccess:     1,
			HardwareState:    1,
			RebootStatus:     0,
			FirmwareVersion:  "1.0.0",
			HardwareRevision: "1.0",
			Model:            "Glasfaser Modem 2",
			SerialNumber:     "ABC123",
		}, nil
	}

	handler := web.APIHandler(metricsFunc)
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/status", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	var status web.StatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	// Verify optical values
	if status.Optical.RxPower == nil || *status.Optical.RxPower != -16.46 {
		t.Errorf("expected RxPower -16.46, got %v", status.Optical.RxPower)
	}

	if status.Optical.TxPower == nil || *status.Optical.TxPower != 2.5 {
		t.Errorf("expected TxPower 2.5, got %v", status.Optical.TxPower)
	}

	if status.Optical.BIPErrors != 0 {
		t.Errorf("expected BIPErrors 0, got %d", status.Optical.BIPErrors)
	}

	// Verify network values
	if status.Network.RxBytes != 123456789 {
		t.Errorf("expected RxBytes 123456789, got %d", status.Network.RxBytes)
	}

	if status.Network.TxBytes != 987654321 {
		t.Errorf("expected TxBytes 987654321, got %d", status.Network.TxBytes)
	}

	// Verify status values
	if status.Status.UptimeSeconds != 123456 {
		t.Errorf("expected UptimeSeconds 123456, got %d", status.Status.UptimeSeconds)
	}

	if status.Status.LinkSpeedMbps != 2500 {
		t.Errorf("expected LinkSpeedMbps 2500, got %d", status.Status.LinkSpeedMbps)
	}

	if !status.Status.PloamSuccess {
		t.Error("expected PloamSuccess true")
	}

	// Verify device info
	if status.Device.FirmwareVersion != "1.0.0" {
		t.Errorf("expected FirmwareVersion 1.0.0, got %s", status.Device.FirmwareVersion)
	}

	if status.Device.Model != "Glasfaser Modem 2" {
		t.Errorf("expected Model 'Glasfaser Modem 2', got %s", status.Device.Model)
	}

	// No scrape error
	if status.ScrapeError {
		t.Error("expected ScrapeError false")
	}
}

func TestAPIHandler_ScrapeError(t *testing.T) {
	metricsFunc := func() (metrics.Metrics, error) {
		return metrics.Metrics{}, errors.New("connection refused")
	}

	handler := web.APIHandler(metricsFunc)
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/status", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer func() { _ = resp.Body.Close() }()

	// Should still return 200 with scrapeError flag
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var status web.StatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if !status.ScrapeError {
		t.Error("expected ScrapeError true when metrics func returns error")
	}
}

func TestAPIHandler_NaNHandling(t *testing.T) {
	metricsFunc := func() (metrics.Metrics, error) {
		return metrics.Metrics{
			RxPowerDBm: math.NaN(),
			TxPowerDBm: math.NaN(),
		}, nil
	}

	handler := web.APIHandler(metricsFunc)
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/status", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer func() { _ = resp.Body.Close() }()

	var status web.StatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	// NaN should be serialized as null
	if status.Optical.RxPower != nil {
		t.Error("expected RxPower to be nil for NaN")
	}

	if status.Optical.TxPower != nil {
		t.Error("expected TxPower to be nil for NaN")
	}
}
