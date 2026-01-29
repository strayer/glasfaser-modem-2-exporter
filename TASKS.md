# Tasks

## Backlog
- [ ] Add GitHub Actions CI workflow (go test, golangci-lint).

## In Progress

## Done
- [x] Set up local development tooling (golangci-lint, mise tasks).
- [x] Implement collector/exporter handler.
- [x] Define Prometheus metrics.
- [x] Write tests for metric mapping and defaults (TDD).
- [x] Implement float parsing rule and integer parsing helpers.
- [x] Write tests for float parsing rule and integer parsing helpers (TDD).
- [x] Add Prometheus client dependency.
- [x] Set up Go module and basic project layout.
- [x] Add sample fixture JSON.
- [x] Write tests for JSON parse + `varid` map construction (TDD).
- [x] Implement JSON fetch + parse (`varid` -> `varvalue` map).
- [x] Add basic README with usage and example scrape config.

## Changelog
- [x] 2026-01-29: Set up golangci-lint v2 config and mise tasks for linting/formatting.
- [x] 2026-01-29: Added collector/exporter HTTP handler with TDD tests.
- [x] 2026-01-29: Added Prometheus collector with TDD tests.
- [x] 2026-01-29: Added metric extraction with TDD tests.
- [x] 2026-01-29: Added float/uint parsing helpers with tests.
- [x] 2026-01-29: Added Prometheus client_golang dependency.
- [x] 2026-01-29: Initialized Go module and basic layout.
- [x] 2026-01-29: Added fixture JSON files (provisioned and not provisioned).
- [x] 2026-01-29: Initialized task tracker.
- [x] 2026-01-29: Added TDD tests for parsing status JSON.
- [x] 2026-01-29: Implemented status JSON parser.
