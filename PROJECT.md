# Glasfaser-Modem 2 Exporter

Prometheus exporter for the Telekom Glasfaser-Modem 2 (ONT). It scrapes the ONT status JSON and exposes a small, stable set of Prometheus metrics.

## Goals and Principles
- KISS: one binary, minimal config, minimal moving parts.
- DRY: one source of truth for each metric and parsing rule.
- TDD: core parsing and metric mapping are covered by tests.
- Avoid overengineering.
- Minimize dependencies; use the Go standard library and the official Prometheus client.
- Modern, idiomatic Go.

## Data Source
- URL: `http://192.168.100.1/ONT/client/data/Status.json`
- Method: `GET`
- Auth: none
- Format: JSON array of objects

Example element:
```json
{"vartype":"value","varid":"rxpower","varvalue":"-16.46"}
```

Important: the JSON array order is not stable. Always map `varid` -> `varvalue`.

## Parsing Rules
### Global Float Parsing Rule (Critical)
Applies to dBm, voltage/signal fields (optical power, etc.):
1) Trim whitespace
2) If value is `"--"`, `"N/A"`, or `""` then return `NaN`
3) Else parse as float

### Integer Parsing
- For counters (bytes/packets/errors), parse as unsigned 64-bit integer.
- If missing or invalid, default to `0` and continue.

## Metrics
All metrics are prefixed `ont_`.

### Optical Metrics
- `ont_optical_rx_power_dbm` (Gauge) from `rxpower` (NaN allowed)
- `ont_optical_tx_power_dbm` (Gauge) from `txpower` (NaN allowed)
- `ont_optical_bip_errors_total` (Counter) from `rxbip_crc` (int, default 0)

### Network Traffic
- `ont_lan_receive_bytes_total` (Counter) from `rxbytes`
- `ont_lan_transmit_bytes_total` (Counter) from `txbytes`
- `ont_lan_receive_packets_total` (Counter) from `rxpackets`
- `ont_lan_transmit_packets_total` (Counter) from `txpackets`
- `ont_lan_packets_dropped_total` (Counter) from `rxdrop_packets`

### Device Health & Status
- `ont_uptime_seconds` (Gauge) from `stability`
- `ont_lan_link_speed_mbps` (Gauge) from `link_status`
- `ont_reboot_status` (Gauge) from `rebooting`
- `ont_ploam_success` (Gauge) from `ploam_success`
- `ont_hardware_state` (Gauge) from `hardware_state`

### Device Metadata (Info)
Metric: `ont_device_info` (value `1`) with labels:
- `firmware_version` from `firmware_version`
- `hardware_revision` from `hardware_revision`
- `serial_number` from `serial_number`
- `model` from `device_name`

### Fields to Ignore
Ignore/discard:
- `datetime`
- `page_title`
- `ui_version`
- `service_mode`
- `save_fails`
- `fw_version_standby` (optional to ignore)

## Architecture (Planned)
- Small Go module with a single `cmd/` entrypoint.
- HTTP client fetches JSON, parses into map, updates Prometheus metrics.
- No background state beyond metrics.
- Config via flags and env (target URL, listen address, timeout).

## Testing Strategy (TDD)
- Unit tests for:
  - JSON parsing and `varid` map construction
  - Float parsing rule (including `"--"`, `"N/A"`, `""`)
  - Metric mapping and defaults
- Include a fixture JSON file derived from the sample payload.

## Non-Goals
- No UI.
- No authentication.
- No persistence beyond Prometheus counters.
- No SNMP or vendor-specific protocols.

## Repository Conventions
- `cmd/glasfaser-modem-2-exporter` for the main.
- `internal/` for parser + scraper.
- `internal/metrics` for Prometheus metric definitions.
- `internal/httpclient` (optional) if needed to keep logic clean.

