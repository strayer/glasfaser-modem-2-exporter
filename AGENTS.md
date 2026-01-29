# AGENTS.md

This file provides guidance to AI coding assistants when working with code in this repository.

## Project Overview

Prometheus exporter for the Telekom Glasfaser-Modem 2 (ONT). Scrapes the modem's status JSON endpoint and exposes metrics for monitoring optical signal levels, network traffic, and device health.

## Principles

- KISS: one binary, minimal config, minimal moving parts
- TDD: write tests before implementing functionality
- Minimize dependencies; use Go standard library and official Prometheus client

## Commands

```bash
# Run all tests
go test ./...

# Run a single test
go test ./internal/status -run TestParse_NotProvisionedFixture

# Build the binary (once main.go exists)
go build -o glasfaser-modem-2-exporter ./cmd/glasfaser-modem-2-exporter
```

## Architecture

- `cmd/glasfaser-modem-2-exporter/` - Main entrypoint
- `internal/status/` - JSON parsing for ONT status endpoint (varid → varvalue map)
- `internal/metrics/` - Prometheus metric definitions (planned)
- `testdata/` - Fixture JSON files for testing

## Data Source

The ONT exposes `http://192.168.100.1/ONT/client/data/Status.json` returning a JSON array:
```json
[{"vartype":"value","varid":"rxpower","varvalue":"-16.46"}, ...]
```

The array order is not stable. Always use the `varid` → `varvalue` map from `status.Parse()`.

## Parsing Rules

**Float parsing** (dBm, voltage fields): trim whitespace, return NaN for `"--"`, `"N/A"`, or `""`, else parse as float.

**Integer parsing** (counters): parse as uint64, default to 0 if missing/invalid.

## Metric Field Mappings

All metrics prefixed `ont_`. See README.md for full metric list.

| varid | Metric | Type |
|-------|--------|------|
| `rxpower` | `ont_optical_rx_power_dbm` | Gauge (float) |
| `txpower` | `ont_optical_tx_power_dbm` | Gauge (float) |
| `rxbip_crc` | `ont_optical_bip_errors_total` | Counter (int) |
| `rxbytes` | `ont_lan_receive_bytes_total` | Counter (int) |
| `txbytes` | `ont_lan_transmit_bytes_total` | Counter (int) |
| `rxpackets` | `ont_lan_receive_packets_total` | Counter (int) |
| `txpackets` | `ont_lan_transmit_packets_total` | Counter (int) |
| `rxdrop_packets` | `ont_lan_packets_dropped_total` | Counter (int) |
| `stability` | `ont_uptime_seconds` | Gauge |
| `link_status` | `ont_lan_link_speed_mbps` | Gauge |
| `rebooting` | `ont_reboot_status` | Gauge |
| `ploam_success` | `ont_ploam_success` | Gauge |
| `hardware_state` | `ont_hardware_state` | Gauge |

**Info metric labels** (`ont_device_info`): `firmware_version`, `hardware_revision`, `serial_number`, `model` (from `device_name`)

**Ignored fields**: `datetime`, `page_title`, `ui_version`, `service_mode`, `save_fails`, `fw_version_standby`

## Development

This project uses TDD. Write tests before implementing functionality.

## Task Tracking

Tasks are tracked in TASKS.md. When working on tasks:
- Move tasks from Backlog to In Progress when starting
- Move tasks to Done when completed (remove the task from Backlog/In Progress)
- Add a changelog entry with the date when completing significant work
