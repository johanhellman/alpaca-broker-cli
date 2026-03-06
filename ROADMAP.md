# Alpaca CLI Roadmap

This document outlines the planned iterations to evolve the `alpaca-cli` project (which contains the `alpaca-broker` and `alpaca-trader` binaries) from its current MVP state into a fully-fledged, production-ready suite of CLI tools.

## Iteration 0: Foundation & Technical Quality (**Completed**)
Before iterating on features, the underlying codebase needs to be brought up to standard best practices for modern Go applications.
- **CI/CD & Linting**: Introduce `.github/workflows/ci.yml` to run `go test`, `go build`, and `golangci-lint` on PRs. Add a `.golangci.yaml` configuration with robust checks (`errcheck`, `gosec`, `gocyclo`, `revive`, `gofmt`) and a `make lint` target.
- **Error Handling & Architecture**: Replace bare `log.Fatal` calls (which break testability and cause unexpected crashes) with structured error returns. Implement proper Cobra error handling and introduce API timeouts in contexts (replacing `context.Background()`).
- **Testing Coverage**: Establish a testing strategy and write initial tests for `cmd/trader` and `cmd/broker`, paving the way for fully-mocked API interactions instead of relying on a single root test file.
- **Hygiene**: Bump outdated dependencies (e.g., Viper, Validator, gRPC) to latest secure versions.
- **Outcome**: A stable, tested, and automatically validated repository ready for feature additions.

## Iteration 1: Parameterization and Filtering (**Completed**)
The current MVP commands rely heavily on raw JSON files (`--file payload.json`) for complex mutations and lack filtering options for lists.
- **Goal**: Add robust, native command-line flags.
- **Trader Example**: `alpaca-trader orders create --symbol AAPL --qty 1 --side buy --type market --time-in-force day`
- **Broker Example**: `alpaca-broker accounts list --status ACTIVE --limit 50`
- **Outcome**: Users rarely need to hand-write JSON or parse through excessive results to perform standard daily operations.

## Iteration 2: Feature Completeness (Breadth) (**Completed**)
Expand the command surface area to cover the entirety of both Alpaca API specifications.
- **Trader API Additions**: 
  - [x] Assets (`alpaca-trader assets`)
  - [x] Watchlists (`alpaca-trader watchlists`)
  - [x] Corporate Actions
  - [x] Market Data (fetching historical bars/quotes natively in the terminal)
- **Broker API Additions**: 
  - [x] Journals (moving money between sub-accounts and the firm)
  - [x] Documents (uploading KYC / W-8BENs)
  - [x] Events (SSE streaming of account status changes directly to `stdout`)
- **Outcome**: The CLI tools become a true 1:1 functional reflection of the Alpaca API Reference.

## Iteration 3: Complex Output & Data Extraction (UX Polish) (**Completed**)
Standard API responses can be inherently noisy for terminal utilities. The CLI needs filtering and robust output handling capabilities.
- **Goal**: Implement `jq`-like filtering or wide/narrow table outputs. Automate pagination handling.
- [x] **Feature**: Provide a `--format json` mode and add a global `--query "portfolio_value"` flag (using a fast JSON processor like `gjson`).
- [x] **Feature**: Auto-follow pagination links/tokens for `list` endpoints via an `--all` flag.
- **Outcome**: The CLI becomes a powerful data plumbing tool for backend shell scripts and automations.

## Iteration 4: CI/CD & Distribution (**Completed**)
To be production-ready and easily adopted by other developers or traders, the project must be easily distributable rather than relying strictly on the Go development toolchain (`go install`).
- **Goal**: Automated binary builds and native package manager distribution.
- [x] **Feature**: Introduce [GoReleaser](https://goreleaser.com/) to build binaries for macOS, Linux, and Windows automatically via GitHub Actions pipelines.
- [x] **Feature**: Create a Homebrew Tap (e.g., `brew install johanhellman/tap/alpaca-cli`).
- [x] **Feature**: Provide pre-compiled `.tar.gz` and `.zip` releases on the GitHub Releases page.
- **Outcome**: Zero-friction installation for non-Go developers.

## Iteration 5: Data Export & Offline Analytics (**Completed**)
Traders and brokers need to reconcile data offline for accounting, taxes, or strategy backtesting. This is a natural extension of the tabular output features introduced in Iteration 3.
- **Goal**: Make the CLI a first-class data ingestion tool.
- [x] **Feature**: Support `--output csv` natively across all data endpoints to inject directly into reporting pipelines.
- [x] **Feature**: Map comprehensive local validation shells `scripts/test-trader-e2e.sh` and `scripts/test-broker-e2e.sh` to track system functionality against Live parameters (maintained and patched by AntiGravity).

## Iteration 6: Multi-Environment & Secure Credential Management (**Planned**)
As users scale, they will likely manage multiple accounts (e.g., Paper vs Live, or multiple Broker sub-environments). Environment variables will always remain supported as the fundamental configuration layer, but power users need faster context switching.
- **Goal**: Introduce an optional `kubectl`-style context manager and secure API keys.
- **Feature**: `alpaca-trader config use-context live` and `alpaca-trader config use-context paper` to switch environments on the fly without changing environment variables (environment variables will still act as overrides).
- **Feature**: Integrate native OS Keychain/Keystore backing so API Secrets don't have to be stored in plaintext inside `~/.alpaca-trader.yaml`.

## Iteration 7: Advanced TUI & Live Streaming (**Planned**)
A strictly command-and-response CLI can feel static for financial tools.
- **Goal**: Implement high-fidelity Terminal User Interfaces (TUI) and real-time streaming displays.
- **Feature**: Use [Bubble Tea](https://github.com/charmbracelet/bubbletea) to build a split-pane dashboard (`alpaca-trader dashboard`) showing live portfolio PnL, active orders, and real-time market data quotes. This functions purely as a specialized sub-command within the CLI, exactly like `htop`, `lazygit`, or `k9s`.
- **Feature**: Enhance the `events` streams to use dynamic terminal spinners and colors that update in-place as trade events pour in.
