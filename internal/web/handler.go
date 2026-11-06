// Package web serves the HTML dashboard and the JSON status API.
package web

import (
	_ "embed"
	"encoding/json"
	"math"
	"net/http"

	"github.com/strayer/glasfaser-modem-2-exporter/internal/metrics"
)

//go:embed dashboard.html
var dashboardHTML []byte

// MetricsFunc is a function that fetches current metrics from the ONT.
type MetricsFunc func() (metrics.Metrics, error)

// StatusResponse represents the JSON response from /api/status.
type StatusResponse struct {
	Optical     OpticalStatus `json:"optical"`
	Network     NetworkStatus `json:"network"`
	Status      DeviceStatus  `json:"status"`
	Device      DeviceInfo    `json:"device"`
	ScrapeError bool          `json:"scrapeError"`
}

// OpticalStatus contains optical signal metrics.
type OpticalStatus struct {
	RxPower   *float64 `json:"rxPower"`
	TxPower   *float64 `json:"txPower"`
	BIPErrors uint64   `json:"bipErrors"`
}

// NetworkStatus contains network traffic metrics.
type NetworkStatus struct {
	RxBytes        uint64 `json:"rxBytes"`
	TxBytes        uint64 `json:"txBytes"`
	RxPackets      uint64 `json:"rxPackets"`
	TxPackets      uint64 `json:"txPackets"`
	DroppedPackets uint64 `json:"droppedPackets"`
}

// DeviceStatus contains device operational status.
type DeviceStatus struct {
	UptimeSeconds uint64 `json:"uptimeSeconds"`
	LinkSpeedMbps uint64 `json:"linkSpeedMbps"`
	PloamSuccess  bool   `json:"ploamSuccess"`
	HardwareState uint64 `json:"hardwareState"`
	Rebooting     bool   `json:"rebooting"`
}

// DeviceInfo contains device identification information.
type DeviceInfo struct {
	FirmwareVersion  string `json:"firmwareVersion"`
	HardwareRevision string `json:"hardwareRevision"`
	Model            string `json:"model"`
	SerialNumber     string `json:"serialNumber"`
}

// DashboardHandler returns an HTTP handler that serves the dashboard HTML.
func DashboardHandler(_ MetricsFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(dashboardHTML)
	})
}

// APIHandler returns an HTTP handler that serves JSON status.
func APIHandler(metricsFunc MetricsFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		m, err := metricsFunc()

		var resp StatusResponse
		if err != nil {
			resp.ScrapeError = true
		} else {
			resp = metricsToResponse(&m)
		}

		_ = json.NewEncoder(w).Encode(resp)
	})
}

// metricsToResponse converts metrics.Metrics to StatusResponse.
func metricsToResponse(m *metrics.Metrics) StatusResponse {
	resp := StatusResponse{
		Optical: OpticalStatus{
			BIPErrors: m.BIPErrors,
		},
		Network: NetworkStatus{
			RxBytes:        m.RxBytes,
			TxBytes:        m.TxBytes,
			RxPackets:      m.RxPackets,
			TxPackets:      m.TxPackets,
			DroppedPackets: m.RxDropPackets,
		},
		Status: DeviceStatus{
			UptimeSeconds: m.UptimeSeconds,
			LinkSpeedMbps: m.LinkSpeedMbps,
			PloamSuccess:  m.PloamSuccess == 1,
			HardwareState: m.HardwareState,
			Rebooting:     m.RebootStatus == 1,
		},
		Device: DeviceInfo{
			FirmwareVersion:  m.FirmwareVersion,
			HardwareRevision: m.HardwareRevision,
			Model:            m.Model,
			SerialNumber:     m.SerialNumber,
		},
		ScrapeError: false,
	}

	// Handle NaN values for optical power (use nil pointer)
	if !math.IsNaN(m.RxPowerDBm) {
		resp.Optical.RxPower = &m.RxPowerDBm
	}

	if !math.IsNaN(m.TxPowerDBm) {
		resp.Optical.TxPower = &m.TxPowerDBm
	}

	return resp
}
