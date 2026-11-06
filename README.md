# Glasfaser-Modem 2 Exporter

Prometheus exporter for the Telekom Glasfaser-Modem 2 (ONT). Scrapes the modem's status JSON endpoint and exposes metrics for monitoring optical signal levels, network traffic, and device health.

## Usage

```sh
glasfaser-modem-2-exporter [flags]
```

| Flag | Default | Description |
|---|---|---|
| `-listen-address` | `:9100` | Address to listen on for the dashboard and metrics |
| `-ont-url` | `http://192.168.100.1/ONT/client/data/Status.json` | URL of the ONT status JSON endpoint |
| `-timeout` | `10s` | Timeout for ONT HTTP requests |

The ONT is only reachable from a host with a route to `192.168.100.1`, i.e. one that sits on the LAN side of the modem (on a typical setup the router or a machine directly behind it).

## Running with Docker

```yaml
# compose.yaml
services:
  glasfaser-modem-2-exporter:
    image: ghcr.io/strayer/glasfaser-modem-2-exporter:latest
    read_only: true
    restart: unless-stopped
    ports:
      - "9100:9100"
```

The image contains a single static binary and runs as an unprivileged user (`65534`). Pass flags via `command:` to change the listen address, ONT URL or timeout.

## Web Interface

The exporter provides a human-readable web dashboard at `/` with:

- Signal quality gauges (RX/TX power with color-coded status)
- Network traffic statistics
- Device status and uptime
- Device information (firmware, model, serial number)

The dashboard uses the Catppuccin color scheme and automatically follows your system's light/dark mode preference. Data refreshes every 5 seconds.

### Endpoints

| Path | Description |
|------|-------------|
| `/` | Web dashboard |
| `/metrics` | Prometheus metrics |
| `/api/status` | JSON API for current status |

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

This project uses [mise](https://mise.jdx.dev) for tool management and [prek](https://github.com/j178/prek) for pre-commit hooks.

```sh
mise install      # installs Go, goimports, golangci-lint, govulncheck
prek install      # installs pre-commit hooks
```

```sh
mise run test     # go test -race -cover ./...
mise run lint     # format-check, vet, golangci-lint, govulncheck
mise run format   # goimports -w .
mise run run      # run the exporter locally
```

Pre-commit hooks run automatically on `git commit`:
- **goimports** — auto-formats Go files
- **golangci-lint** — lints for errors (no auto-fix)

## Releasing

Every PR carries one of the labels `breaking`, `feature`, `enhancement`, `bug`, `chore` or `dependencies`; the labels drive the semver bump and the generated release notes. `scripts/release.sh` proposes the next version, creates a signed annotated tag and pushes it; the release workflow then builds the multi-arch image and publishes the GitHub Release.

## AI Usage Notice

For transparency, this notice clarifies how AI is used responsibly in this project. A Large Language Model (LLM) may assist in generating code, but it operates strictly as a tool under human control. The process is governed by two principles: first, all tasks are designed and specified by a human before the LLM is engaged. Second, **all AI-generated output is subject to a mandatory human review** to verify its correctness, security, and alignment with the project's quality standards. This ensures human oversight and accountability for all code committed to the repository.

## License

AGPL-3.0
