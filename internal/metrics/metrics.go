// Package metrics extracts typed metric values from parsed ONT status data.
package metrics

import (
	"github.com/strayer/glasfaser-modem-2-exporter/internal/status"
)

// Metrics holds the extracted metric values from ONT status data.
type Metrics struct {
	// Optical power (Gauge, float - NaN if unavailable)
	RxPowerDBm float64
	TxPowerDBm float64

	// Counters (uint64)
	RxBytes       uint64
	TxBytes       uint64
	RxPackets     uint64
	TxPackets     uint64
	RxDropPackets uint64
	BIPErrors     uint64

	// Gauges (uint64)
	UptimeSeconds uint64
	LinkSpeedMbps uint64
	RebootStatus  uint64
	PloamSuccess  uint64
	HardwareState uint64

	// Info labels (strings)
	FirmwareVersion  string
	HardwareRevision string
	SerialNumber     string
	Model            string
}

// Extract converts a varid->varvalue map into structured Metrics.
func Extract(data map[string]string) Metrics {
	return Metrics{
		// Float metrics (use ParseFloat for NaN handling)
		RxPowerDBm: status.ParseFloat(data["rxpower"]),
		TxPowerDBm: status.ParseFloat(data["txpower"]),

		// Counter metrics
		RxBytes:       status.ParseUint(data["rxbytes"]),
		TxBytes:       status.ParseUint(data["txbytes"]),
		RxPackets:     status.ParseUint(data["rxpackets"]),
		TxPackets:     status.ParseUint(data["txpackets"]),
		RxDropPackets: status.ParseUint(data["rxdrop_packets"]),
		BIPErrors:     status.ParseUint(data["rxbip_crc"]),

		// Gauge metrics
		UptimeSeconds: status.ParseUint(data["stability"]),
		LinkSpeedMbps: status.ParseUint(data["link_status"]),
		RebootStatus:  status.ParseUint(data["rebooting"]),
		PloamSuccess:  status.ParseUint(data["ploam_success"]),
		HardwareState: status.ParseUint(data["hardware_state"]),

		// Info labels (direct string values)
		FirmwareVersion:  data["firmware_version"],
		HardwareRevision: data["hardware_revision"],
		SerialNumber:     data["serial_number"],
		Model:            data["device_name"],
	}
}
