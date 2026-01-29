package collector

import (
	"math"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/strayer/glasfaser-modem-2-exporter/internal/metrics"
)

func TestDescribe_ReturnsAllMetricDescriptors(t *testing.T) {
	c := New(func() (metrics.Metrics, error) {
		return metrics.Metrics{}, nil
	})

	ch := make(chan *prometheus.Desc, 100)
	c.Describe(ch)
	close(ch)

	var descs []*prometheus.Desc
	for d := range ch {
		descs = append(descs, d)
	}

	// Expected metrics: 2 optical gauges + 6 counters + 5 gauges + 1 info = 14
	expectedCount := 14
	if len(descs) != expectedCount {
		t.Errorf("expected %d descriptors, got %d", expectedCount, len(descs))
	}

	// Check that all expected metric names are present
	expectedNames := []string{
		"ont_optical_rx_power_dbm",
		"ont_optical_tx_power_dbm",
		"ont_optical_bip_errors_total",
		"ont_lan_receive_bytes_total",
		"ont_lan_transmit_bytes_total",
		"ont_lan_receive_packets_total",
		"ont_lan_transmit_packets_total",
		"ont_lan_packets_dropped_total",
		"ont_uptime_seconds",
		"ont_lan_link_speed_mbps",
		"ont_reboot_status",
		"ont_ploam_success",
		"ont_hardware_state",
		"ont_device_info",
	}

	descStrings := make([]string, len(descs))
	for i, d := range descs {
		descStrings[i] = d.String()
	}

	for _, name := range expectedNames {
		found := false

		for _, s := range descStrings {
			if strings.Contains(s, "\""+name+"\"") {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("expected metric %q not found in descriptors", name)
		}
	}
}

func TestCollect_ProvisionedMetrics(t *testing.T) {
	m := metrics.Metrics{
		RxPowerDBm:       -15.77,
		TxPowerDBm:       2.83,
		RxBytes:          210987654321,
		TxBytes:          123456789012,
		RxPackets:        1500000000,
		TxPackets:        1200000000,
		RxDropPackets:    12345,
		BIPErrors:        42,
		UptimeSeconds:    2592000,
		LinkSpeedMbps:    2500,
		RebootStatus:     0,
		PloamSuccess:     1,
		HardwareState:    1,
		FirmwareVersion:  "1.2.3",
		HardwareRevision: "HW-REV-A",
		SerialNumber:     "ONT123456",
		Model:            "Glasfaser-Modem 2",
	}

	c := New(func() (metrics.Metrics, error) {
		return m, nil
	})

	expected := `
# HELP ont_optical_rx_power_dbm Receive optical power in dBm
# TYPE ont_optical_rx_power_dbm gauge
ont_optical_rx_power_dbm -15.77
# HELP ont_optical_tx_power_dbm Transmit optical power in dBm
# TYPE ont_optical_tx_power_dbm gauge
ont_optical_tx_power_dbm 2.83
# HELP ont_optical_bip_errors_total BIP/CRC errors
# TYPE ont_optical_bip_errors_total counter
ont_optical_bip_errors_total 42
# HELP ont_lan_receive_bytes_total Bytes received
# TYPE ont_lan_receive_bytes_total counter
ont_lan_receive_bytes_total 2.10987654321e+11
# HELP ont_lan_transmit_bytes_total Bytes transmitted
# TYPE ont_lan_transmit_bytes_total counter
ont_lan_transmit_bytes_total 1.23456789012e+11
# HELP ont_lan_receive_packets_total Packets received
# TYPE ont_lan_receive_packets_total counter
ont_lan_receive_packets_total 1.5e+09
# HELP ont_lan_transmit_packets_total Packets transmitted
# TYPE ont_lan_transmit_packets_total counter
ont_lan_transmit_packets_total 1.2e+09
# HELP ont_lan_packets_dropped_total Packets dropped
# TYPE ont_lan_packets_dropped_total counter
ont_lan_packets_dropped_total 12345
# HELP ont_uptime_seconds Device uptime in seconds
# TYPE ont_uptime_seconds gauge
ont_uptime_seconds 2.592e+06
# HELP ont_lan_link_speed_mbps LAN link speed in Mbps
# TYPE ont_lan_link_speed_mbps gauge
ont_lan_link_speed_mbps 2500
# HELP ont_reboot_status Reboot status
# TYPE ont_reboot_status gauge
ont_reboot_status 0
# HELP ont_ploam_success PLOAM registration success
# TYPE ont_ploam_success gauge
ont_ploam_success 1
# HELP ont_hardware_state Hardware state
# TYPE ont_hardware_state gauge
ont_hardware_state 1
# HELP ont_device_info Device information
# TYPE ont_device_info gauge
ont_device_info{firmware_version="1.2.3",hardware_revision="HW-REV-A",model="Glasfaser-Modem 2",serial_number="ONT123456"} 1
`

	if err := testutil.CollectAndCompare(c, strings.NewReader(expected)); err != nil {
		t.Errorf("unexpected metrics:\n%s", err)
	}
}

func TestCollect_NaNValuesOmitted(t *testing.T) {
	// When power values are NaN (not provisioned), they should not be emitted
	m := metrics.Metrics{
		RxPowerDBm: math.NaN(),
		TxPowerDBm: math.NaN(),
	}

	c := New(func() (metrics.Metrics, error) {
		return m, nil
	})

	ch := make(chan prometheus.Metric, 100)
	c.Collect(ch)
	close(ch)

	for metric := range ch {
		desc := metric.Desc().String()
		if strings.Contains(desc, "ont_optical_rx_power_dbm") ||
			strings.Contains(desc, "ont_optical_tx_power_dbm") {
			t.Errorf("NaN metric should not be emitted: %s", desc)
		}
	}
}

func TestCollect_InfoMetricLabels(t *testing.T) {
	m := metrics.Metrics{
		FirmwareVersion:  "v2.0.0",
		HardwareRevision: "Rev-B",
		SerialNumber:     "SN999",
		Model:            "Test Model",
	}

	c := New(func() (metrics.Metrics, error) {
		return m, nil
	})

	expected := `
# HELP ont_device_info Device information
# TYPE ont_device_info gauge
ont_device_info{firmware_version="v2.0.0",hardware_revision="Rev-B",model="Test Model",serial_number="SN999"} 1
`

	if err := testutil.CollectAndCompare(c, strings.NewReader(expected), "ont_device_info"); err != nil {
		t.Errorf("unexpected info metric:\n%s", err)
	}
}
