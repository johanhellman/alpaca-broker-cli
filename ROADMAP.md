# Alpaca CLI Roadmap

This document outlines the planned iterations to evolve the `alpaca-cli` project (which contains the `alpaca-broker` and `alpaca-trader` binaries) from its current MVP state into a fully-fledged, production-ready suite of CLI tools.

## Iteration 1: Parameterization and Filtering (Next Up)
The current MVP commands rely heavily on raw JSON files (`--file payload.json`) for complex mutations and lack filtering options for lists.
- **Goal**: Add robust, native command-line flags.
- **Trader Example**: `alpaca-trader orders create --symbol AAPL --qty 1 --side buy --type market --time-in-force day`
- **Broker Example**: `alpaca-broker accounts list --status ACTIVE --limit 50`
- **Outcome**: Users rarely need to hand-write JSON or parse through excessive results to perform standard daily operations.

## Iteration 2: Feature Completeness (Breadth)
Expand the command surface area to cover the entirety of both Alpaca API specifications.
- **Trader API Additions**: 
  - Assets (`alpaca-trader assets`)
  - Watchlists (`alpaca-trader watchlists`)
  - Corporate Actions
  - Market Data (fetching historical bars/quotes natively in the terminal)
- **Broker API Additions**: 
  - Journals (moving money between sub-accounts and the firm)
  - Documents (uploading KYC / W-8BENs)
  - Events (SSE streaming of account status changes directly to `stdout`)
- **Outcome**: The CLI tools become a true 1:1 functional reflection of the Alpaca API Reference.

## Iteration 3: Complex Output & Data Extraction (UX Polish)
Standard API responses can be inherently noisy for terminal utilities. The CLI needs filtering and robust output handling capabilities.
- **Goal**: Implement `jq`-like filtering or wide/narrow table outputs. Automate pagination handling.
- **Feature**: Provide a `--format json` mode and add a global `--query "portfolio_value"` flag (using a fast JSON processor like `gjson`).
- **Feature**: Auto-follow pagination links/tokens for `list` endpoints via an `--all` flag.
- **Outcome**: The CLI becomes a powerful data plumbing tool for backend shell scripts and automations.

## Iteration 4: CI/CD & Distribution
To be production-ready and easily adopted by other developers or traders, the project must be easily distributable rather than relying strictly on the Go development toolchain (`go install`).
- **Goal**: Automated binary builds and native package manager distribution.
- **Feature**: Introduce [GoReleaser](https://goreleaser.com/) to build binaries for macOS, Linux, and Windows automatically via GitHub Actions pipelines.
- **Feature**: Create a Homebrew Tap (e.g., `brew install johanhellman/tap/alpaca-cli`).
- **Feature**: Provide pre-compiled `.tar.gz` and `.zip` releases on the GitHub Releases page.
- **Outcome**: Zero-friction installation for non-Go developers.
