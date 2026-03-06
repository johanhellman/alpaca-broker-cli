# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Complete Iteration 3 for complex output extraction (tabular, json query, csv).
- Complete Iteration 2 commands for Trader and Broker CLI (assets, watchlists, journals, documents, events).
- Iteration 1 parameterization (refactored to native CLI flags removing json payloads).
- Native csv formatting to cli output pipelines.
- Comprehensive sandbox and paper trading verification e2e scripts (restored to green by AntiGravity).
- GoReleaser and GitHub Actions publishing pipeline.
- Automated API markdown documentation generation for both Broker and Trader pipelines.
- Dependabot configuration for go modules and GitHub Actions.

### Changed
- Project restructured to dual-CLI architecture (`alpaca-broker` and `alpaca-trader`).
- Commercial proprietary license enforcement.
- CI/CD workflow version increments (Go 1.26.0).

### Fixed
- Fixed dead code, gocyclo complexity, and errcheck violations from the technical hygiene review.
- Missing binaries in `make install`.
- [AntiGravity] Fixed Broker API `400` validation errors regarding missing `Agreements`, `Disclosures`, and strictly synthetic `tax_id`s in the `accounts create` payload.
- [AntiGravity] Fixed trailing flag conflicts (`--limit` and `--query`) breaking `test-broker-e2e.sh`.
- Goreleaser `v2` deprecation warnings.
