// Package collector implements the Prometheus collector for ONT metrics.
package collector

import (
	"math"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/strayer/glasfaser-modem-2-exporter/internal/metrics"
)

// MetricsFunc is a function that returns current metrics from the ONT.
type MetricsFunc func() (metrics.Metrics, error)

// Collector implements prometheus.Collector for ONT metrics.
type Collector struct {
	metricsFunc MetricsFunc

	// Optical metrics (Gauge)
	rxPowerDBm *prometheus.Desc
	txPowerDBm *prometheus.Desc

	// Counter metrics
	bipErrors     *prometheus.Desc
	rxBytes       *prometheus.Desc
	txBytes       *prometheus.Desc
	rxPackets     *prometheus.Desc
	txPackets     *prometheus.Desc
	rxDropPackets *prometheus.Desc

	// Gauge metrics
	uptimeSeconds *prometheus.Desc
	linkSpeedMbps *prometheus.Desc
	rebootStatus  *prometheus.Desc
	ploamSuccess  *prometheus.Desc
	hardwareState *prometheus.Desc

	// Info metric
	deviceInfo *prometheus.Desc
}

// newDesc is a helper to create a prometheus.Desc with no labels.
func newDesc(name, help string) *prometheus.Desc {
	return prometheus.NewDesc(name, help, nil, nil)
}

// New creates a new Collector with the given metrics source.
func New(metricsFunc MetricsFunc) *Collector {
	return &Collector{
		metricsFunc:   metricsFunc,
		rxPowerDBm:    newDesc("ont_optical_rx_power_dbm", "Receive optical power in dBm"),
		txPowerDBm:    newDesc("ont_optical_tx_power_dbm", "Transmit optical power in dBm"),
		bipErrors:     newDesc("ont_optical_bip_errors_total", "BIP/CRC errors"),
		rxBytes:       newDesc("ont_lan_receive_bytes_total", "Bytes received"),
		txBytes:       newDesc("ont_lan_transmit_bytes_total", "Bytes transmitted"),
		rxPackets:     newDesc("ont_lan_receive_packets_total", "Packets received"),
		txPackets:     newDesc("ont_lan_transmit_packets_total", "Packets transmitted"),
		rxDropPackets: newDesc("ont_lan_packets_dropped_total", "Packets dropped"),
		uptimeSeconds: newDesc("ont_uptime_seconds", "Device uptime in seconds"),
		linkSpeedMbps: newDesc("ont_lan_link_speed_mbps", "LAN link speed in Mbps"),
		rebootStatus:  newDesc("ont_reboot_status", "Reboot status"),
		ploamSuccess:  newDesc("ont_ploam_success", "PLOAM registration success"),
		hardwareState: newDesc("ont_hardware_state", "Hardware state"),
		deviceInfo: prometheus.NewDesc(
			"ont_device_info",
			"Device information",
			[]string{"firmware_version", "hardware_revision", "model", "serial_number"},
			nil,
		),
	}
}

// Describe implements prometheus.Collector.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.rxPowerDBm

	ch <- c.txPowerDBm

	ch <- c.bipErrors

	ch <- c.rxBytes

	ch <- c.txBytes

	ch <- c.rxPackets

	ch <- c.txPackets

	ch <- c.rxDropPackets

	ch <- c.uptimeSeconds

	ch <- c.linkSpeedMbps

	ch <- c.rebootStatus

	ch <- c.ploamSuccess

	ch <- c.hardwareState

	ch <- c.deviceInfo
}

// Collect implements prometheus.Collector.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	m, err := c.metricsFunc()
	if err != nil {
		// TODO: emit scrape error metric
		return
	}

	// Optical power gauges (omit if NaN)
	if !math.IsNaN(m.RxPowerDBm) {
		ch <- prometheus.MustNewConstMetric(c.rxPowerDBm, prometheus.GaugeValue, m.RxPowerDBm)
	}

	if !math.IsNaN(m.TxPowerDBm) {
		ch <- prometheus.MustNewConstMetric(c.txPowerDBm, prometheus.GaugeValue, m.TxPowerDBm)
	}

	// Counters
	ch <- prometheus.MustNewConstMetric(c.bipErrors, prometheus.CounterValue, float64(m.BIPErrors))

	ch <- prometheus.MustNewConstMetric(c.rxBytes, prometheus.CounterValue, float64(m.RxBytes))

	ch <- prometheus.MustNewConstMetric(c.txBytes, prometheus.CounterValue, float64(m.TxBytes))

	ch <- prometheus.MustNewConstMetric(c.rxPackets, prometheus.CounterValue, float64(m.RxPackets))

	ch <- prometheus.MustNewConstMetric(c.txPackets, prometheus.CounterValue, float64(m.TxPackets))

	ch <- prometheus.MustNewConstMetric(c.rxDropPackets, prometheus.CounterValue, float64(m.RxDropPackets))

	// Gauges
	ch <- prometheus.MustNewConstMetric(c.uptimeSeconds, prometheus.GaugeValue, float64(m.UptimeSeconds))

	ch <- prometheus.MustNewConstMetric(c.linkSpeedMbps, prometheus.GaugeValue, float64(m.LinkSpeedMbps))

	ch <- prometheus.MustNewConstMetric(c.rebootStatus, prometheus.GaugeValue, float64(m.RebootStatus))

	ch <- prometheus.MustNewConstMetric(c.ploamSuccess, prometheus.GaugeValue, float64(m.PloamSuccess))

	ch <- prometheus.MustNewConstMetric(c.hardwareState, prometheus.GaugeValue, float64(m.HardwareState))

	// Info metric with labels
	ch <- prometheus.MustNewConstMetric(
		c.deviceInfo,
		prometheus.GaugeValue,
		1,
		m.FirmwareVersion,
		m.HardwareRevision,
		m.Model,
		m.SerialNumber,
	)
}
