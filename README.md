# Glasfaser-Modem 2 Exporter

Prometheus exporter for the Telekom Glasfaser-Modem 2 (ONT). Scrapes the modem's status JSON endpoint and exposes metrics for monitoring optical signal levels, network traffic, and device health.

## Usage

```bash
glasfaser-modem-2-exporter [flags]
```

Configuration via flags and environment variables:
- Target URL (default: `http://192.168.100.1/ONT/client/data/Status.json`)
- Listen address
- Scrape timeout

## Prometheus Scrape Config

```yaml
scrape_configs:
  - job_name: 'ont'
    static_configs:
      - targets: ['localhost:9100']  # exporter address
```

## Metrics

All metrics are prefixed with `ont_`.

### Optical Metrics
| Metric | Type | Description |
|--------|------|-------------|
| `ont_optical_rx_power_dbm` | Gauge | Receive optical power in dBm |
| `ont_optical_tx_power_dbm` | Gauge | Transmit optical power in dBm |
| `ont_optical_bip_errors_total` | Counter | BIP/CRC errors |

### Network Traffic
| Metric | Type | Description |
|--------|------|-------------|
| `ont_lan_receive_bytes_total` | Counter | Bytes received |
| `ont_lan_transmit_bytes_total` | Counter | Bytes transmitted |
| `ont_lan_receive_packets_total` | Counter | Packets received |
| `ont_lan_transmit_packets_total` | Counter | Packets transmitted |
| `ont_lan_packets_dropped_total` | Counter | Packets dropped |

### Device Health & Status
| Metric | Type | Description |
|--------|------|-------------|
| `ont_uptime_seconds` | Gauge | Device uptime |
| `ont_lan_link_speed_mbps` | Gauge | LAN link speed |
| `ont_reboot_status` | Gauge | Reboot status |
| `ont_ploam_success` | Gauge | PLOAM registration success |
| `ont_hardware_state` | Gauge | Hardware state |

### Device Info
| Metric | Type | Labels |
|--------|------|--------|
| `ont_device_info` | Gauge (value=1) | `firmware_version`, `hardware_revision`, `serial_number`, `model` |

## Development

This project uses [mise](https://mise.jdx.dev/) for tool management and task running.

```bash
# Install dependencies (golangci-lint)
mise install

# Run tests
mise run test

# Build the binary
mise run build

# Run the exporter
mise run start

# Run linter
mise run lint

# Format code
mise run fmt
```

## AI Usage Notice

For transparency, this notice clarifies how AI is used responsibly in this project. A Large Language Model (LLM) may assist in generating code, but it operates strictly as a tool under human control. The process is governed by two principles: first, all tasks are designed and specified by a human before the LLM is engaged. Second, **all AI-generated output is subject to a mandatory human review** to verify its correctness, security, and alignment with the project's quality standards. This ensures human oversight and accountability for all code committed to the repository.

## License

AGPL-3.0
