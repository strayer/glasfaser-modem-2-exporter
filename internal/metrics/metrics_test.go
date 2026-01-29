package metrics

import (
	"math"
	"testing"
)

func TestExtractMetrics_Provisioned(t *testing.T) {
	// Simulates parsed status from a provisioned ONT
	data := map[string]string{
		"device_name":       "Glasfaser-Modem 2",
		"rebooting":         "0",
		"ploam_success":     "1",
		"firmware_version":  "090144.1.0.001",
		"hardware_revision": "V1",
		"hardware_state":    "1",
		"txpackets":         "98234567",
		"txbytes":           "123456789012",
		"rxpackets":         "87654321",
		"rxbytes":           "210987654321",
		"rxdrop_packets":    "12",
		"link_status":       "1000",
		"stability":         "987654",
		"rxbip_crc":         "3",
		"serial_number":     "ONT9F3A1B2C3D4E",
		"txpower":           "2.83",
		"rxpower":           "-15.77",
	}

	m := Extract(data)

	// Float metrics (optical power)
	if m.RxPowerDBm != -15.77 {
		t.Errorf("RxPowerDBm = %v, want -15.77", m.RxPowerDBm)
	}

	if m.TxPowerDBm != 2.83 {
		t.Errorf("TxPowerDBm = %v, want 2.83", m.TxPowerDBm)
	}

	// Counter metrics (traffic)
	if m.RxBytes != 210987654321 {
		t.Errorf("RxBytes = %v, want 210987654321", m.RxBytes)
	}

	if m.TxBytes != 123456789012 {
		t.Errorf("TxBytes = %v, want 123456789012", m.TxBytes)
	}

	if m.RxPackets != 87654321 {
		t.Errorf("RxPackets = %v, want 87654321", m.RxPackets)
	}

	if m.TxPackets != 98234567 {
		t.Errorf("TxPackets = %v, want 98234567", m.TxPackets)
	}

	if m.RxDropPackets != 12 {
		t.Errorf("RxDropPackets = %v, want 12", m.RxDropPackets)
	}

	if m.BIPErrors != 3 {
		t.Errorf("BIPErrors = %v, want 3", m.BIPErrors)
	}

	// Gauge metrics
	if m.UptimeSeconds != 987654 {
		t.Errorf("UptimeSeconds = %v, want 987654", m.UptimeSeconds)
	}

	if m.LinkSpeedMbps != 1000 {
		t.Errorf("LinkSpeedMbps = %v, want 1000", m.LinkSpeedMbps)
	}

	if m.RebootStatus != 0 {
		t.Errorf("RebootStatus = %v, want 0", m.RebootStatus)
	}

	if m.PloamSuccess != 1 {
		t.Errorf("PloamSuccess = %v, want 1", m.PloamSuccess)
	}

	if m.HardwareState != 1 {
		t.Errorf("HardwareState = %v, want 1", m.HardwareState)
	}

	// Info labels
	if m.FirmwareVersion != "090144.1.0.001" {
		t.Errorf("FirmwareVersion = %q, want %q", m.FirmwareVersion, "090144.1.0.001")
	}

	if m.HardwareRevision != "V1" {
		t.Errorf("HardwareRevision = %q, want %q", m.HardwareRevision, "V1")
	}

	if m.SerialNumber != "ONT9F3A1B2C3D4E" {
		t.Errorf("SerialNumber = %q, want %q", m.SerialNumber, "ONT9F3A1B2C3D4E")
	}

	if m.Model != "Glasfaser-Modem 2" {
		t.Errorf("Model = %q, want %q", m.Model, "Glasfaser-Modem 2")
	}
}

func TestExtractMetrics_NotProvisioned(t *testing.T) {
	// Simulates parsed status from an unprovisioned ONT
	data := map[string]string{
		"device_name":       "Glasfaser-Modem 2",
		"rebooting":         "0",
		"ploam_success":     "0",
		"firmware_version":  "090144.1.0.001",
		"hardware_revision": "V1",
		"hardware_state":    "0",
		"txpackets":         "0",
		"txbytes":           "0",
		"rxpackets":         "0",
		"rxbytes":           "0",
		"rxdrop_packets":    "0",
		"link_status":       "0",
		"stability":         "0",
		"rxbip_crc":         "0",
		"serial_number":     "ONT9F3A1B2C3D4E",
		"txpower":           "--",
		"rxpower":           "--",
	}

	m := Extract(data)

	// Float metrics should be NaN for "--"
	if !math.IsNaN(m.RxPowerDBm) {
		t.Errorf("RxPowerDBm = %v, want NaN", m.RxPowerDBm)
	}

	if !math.IsNaN(m.TxPowerDBm) {
		t.Errorf("TxPowerDBm = %v, want NaN", m.TxPowerDBm)
	}

	// All counters should be 0
	if m.RxBytes != 0 {
		t.Errorf("RxBytes = %v, want 0", m.RxBytes)
	}

	if m.TxBytes != 0 {
		t.Errorf("TxBytes = %v, want 0", m.TxBytes)
	}

	// Status flags
	if m.PloamSuccess != 0 {
		t.Errorf("PloamSuccess = %v, want 0", m.PloamSuccess)
	}

	if m.HardwareState != 0 {
		t.Errorf("HardwareState = %v, want 0", m.HardwareState)
	}
}

func TestExtractMetrics_MissingFields(t *testing.T) {
	// Empty map - all fields missing
	data := map[string]string{}

	m := Extract(data)

	// Float metrics should be NaN for missing values
	if !math.IsNaN(m.RxPowerDBm) {
		t.Errorf("RxPowerDBm = %v, want NaN", m.RxPowerDBm)
	}

	if !math.IsNaN(m.TxPowerDBm) {
		t.Errorf("TxPowerDBm = %v, want NaN", m.TxPowerDBm)
	}

	// Integer counters default to 0
	if m.RxBytes != 0 {
		t.Errorf("RxBytes = %v, want 0", m.RxBytes)
	}

	if m.TxBytes != 0 {
		t.Errorf("TxBytes = %v, want 0", m.TxBytes)
	}

	if m.RxPackets != 0 {
		t.Errorf("RxPackets = %v, want 0", m.RxPackets)
	}

	if m.TxPackets != 0 {
		t.Errorf("TxPackets = %v, want 0", m.TxPackets)
	}

	if m.RxDropPackets != 0 {
		t.Errorf("RxDropPackets = %v, want 0", m.RxDropPackets)
	}

	if m.BIPErrors != 0 {
		t.Errorf("BIPErrors = %v, want 0", m.BIPErrors)
	}

	// Gauge metrics default to 0
	if m.UptimeSeconds != 0 {
		t.Errorf("UptimeSeconds = %v, want 0", m.UptimeSeconds)
	}

	if m.LinkSpeedMbps != 0 {
		t.Errorf("LinkSpeedMbps = %v, want 0", m.LinkSpeedMbps)
	}

	if m.RebootStatus != 0 {
		t.Errorf("RebootStatus = %v, want 0", m.RebootStatus)
	}

	if m.PloamSuccess != 0 {
		t.Errorf("PloamSuccess = %v, want 0", m.PloamSuccess)
	}

	if m.HardwareState != 0 {
		t.Errorf("HardwareState = %v, want 0", m.HardwareState)
	}

	// String labels default to empty
	if m.FirmwareVersion != "" {
		t.Errorf("FirmwareVersion = %q, want empty", m.FirmwareVersion)
	}

	if m.HardwareRevision != "" {
		t.Errorf("HardwareRevision = %q, want empty", m.HardwareRevision)
	}

	if m.SerialNumber != "" {
		t.Errorf("SerialNumber = %q, want empty", m.SerialNumber)
	}

	if m.Model != "" {
		t.Errorf("Model = %q, want empty", m.Model)
	}
}

func TestExtractMetrics_WhitespaceHandling(t *testing.T) {
	data := map[string]string{
		"rxpower":   "  -16.46  ",
		"txpower":   "\t2.50\n",
		"rxbytes":   "  12345  ",
		"txbytes":   "  67890  ",
		"rxbip_crc": " 5 ",
	}

	m := Extract(data)

	if m.RxPowerDBm != -16.46 {
		t.Errorf("RxPowerDBm = %v, want -16.46", m.RxPowerDBm)
	}

	if m.TxPowerDBm != 2.50 {
		t.Errorf("TxPowerDBm = %v, want 2.50", m.TxPowerDBm)
	}

	if m.RxBytes != 12345 {
		t.Errorf("RxBytes = %v, want 12345", m.RxBytes)
	}

	if m.TxBytes != 67890 {
		t.Errorf("TxBytes = %v, want 67890", m.TxBytes)
	}

	if m.BIPErrors != 5 {
		t.Errorf("BIPErrors = %v, want 5", m.BIPErrors)
	}
}

func TestExtractMetrics_NAHandling(t *testing.T) {
	data := map[string]string{
		"rxpower": "N/A",
		"txpower": "N/A",
	}

	m := Extract(data)

	if !math.IsNaN(m.RxPowerDBm) {
		t.Errorf("RxPowerDBm = %v, want NaN for N/A", m.RxPowerDBm)
	}

	if !math.IsNaN(m.TxPowerDBm) {
		t.Errorf("TxPowerDBm = %v, want NaN for N/A", m.TxPowerDBm)
	}
}
