# AGENTS.md

This file provides guidance to AI coding assistants when working with code in this repository.

> **Claude Code**: `CLAUDE.md` is a symlink to this file. Always edit `AGENTS.md`, never `CLAUDE.md`.

## Project Overview

Prometheus exporter for the Telekom Glasfaser-Modem 2 (ONT). Scrapes the modem's status JSON endpoint and exposes metrics for monitoring optical signal levels, network traffic, and device health. Read `README.md` for usage and the full metric list.

## Principles

- KISS: one binary, minimal config, minimal moving parts
- TDD: write tests before implementing functionality
- Minimize dependencies; use Go standard library and the official Prometheus client
- Tests use only `testing` and `net/http/httptest`. No test frameworks.
- Use `gofmt`/`goimports` formatting; don't fight it

## Build & test

mise is the single source of truth for tools (Go, goimports, golangci-lint, govulncheck) and tasks:

```sh
mise install           # install the pinned toolchain
mise run format        # goimports -w .
mise run test          # go test -race -cover ./...
mise run lint          # format-check, vet, golangci-lint, govulncheck
mise run run           # run the exporter locally
go test ./internal/status -run TestParse_NotProvisionedFixture   # single test
```

## Commits

Use [Conventional Commits](https://www.conventionalcommits.org/) for all commit messages.

## Pre-commit hooks

This project uses [prek](https://github.com/j178/prek) to run pre-commit hooks. Hooks are defined in `.pre-commit-config.yaml` as `repo: local` system hooks:

1. **goimports** — auto-formats staged `.go` files in-place (`-w`). May modify files during commit.
2. **golangci-lint** — lints against `HEAD` (`--new-from-rev=HEAD`). No auto-fix — fails the commit on violations.

If a commit fails due to goimports reformatting, re-stage the changes and commit again. If golangci-lint fails, fix the reported issues before committing.

## CI

`.github/workflows/ci.yaml` runs on pull requests and pushes to main, installing tools only via `jdx/mise-action`. Jobs (names are required status checks, don't rename): `Test`, `Vet`, `goimports`, `golangci-lint`, `govulncheck`, `Go version sync`, `Docker build`. The Go version must match in `go.mod`, `mise.toml` and the `Dockerfile`. `pr-labels.yml` (job `require-release-label`, required check) demands one of the labels `breaking`, `feature`, `enhancement`, `bug`, `chore`, `dependencies` on every PR. `build-dev.yaml` publishes a `:dev` image on every push to main. `release.yaml` (on `v*` tags) first verifies the signed annotated tag, gates the image build on the unit tests, then creates the GitHub Release (notes generated per `.github/release.yml`, `--prerelease` for `-rc` tags); `vulncheck.yaml` scans the published image weekly. Renovate config is `renovate.json5`; Go and the mise-managed dev tools are grouped as "Go toolchain".

## Releasing

Label every PR (`breaking`, `feature`, `enhancement`, `bug`, `chore` or `dependencies`); the labels drive the semver bump and the release notes. To release, run `scripts/release.sh` (`--rc` for a `vX.Y.Z-rc.N` pre-release). It creates and pushes a signed annotated tag; `release.yaml` does the rest.

## Architecture

- `cmd/glasfaser-modem-2-exporter/` - Main entrypoint
- `internal/status/` - JSON parsing for ONT status endpoint (varid → varvalue map)
- `internal/metrics/` - Metric extraction from status data
- `internal/collector/` - Prometheus collector implementation
- `internal/web/` - HTML dashboard and JSON status API
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
