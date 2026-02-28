# The Story of Alpaca CLI: From MVP to Production

This document outlines the journey of how the Alpaca CLI project (which encompasses both `alpaca-broker` and `alpaca-trader`) was conceived, built, and iteratively improved to its current robust state.

## 1. The Inspiration & Input Instructions
The Alpaca API provides immense power for trading and brokerage operations, but interacting with REST APIs via raw `curl` commands or custom Python scripts can be cumbersome. 

The explicit instructions provided to the agent were to build a first-class, natively compiled command-line interface. The `gogcli` repository by `steipete` was supplied as the structural benchmark and inspiration for what the final product should look like.

To accomplish this, AntiGravity was instructed to use **Go**, leveraging **Cobra** for the command router and **Viper** for configuration management. In addition, the agent was provided with the official Alpaca Markets public API documentation and links to public repositories, examples, and SDKs to use as the foundational ground truth.

## 2. AI Collaboration & Technical Approach
The path from those initial instructions to the mature product we have today was primarily driven by AI interacting in a dual-agent dynamic:
- **AntiGravity (Main Driver):** Acted as the lead software engineer responsible for architecting the application, writing the code, mapping the Alpaca SDKs, and executing the roadmap.
- **QA/Review Agent:** Functioned as an external code reviewer, periodically evaluating the repository against Go best practices to identify potential hygiene issues, security flaws, missing legal/community files, and CI/CD gaps.

The guiding principles for the architecture were:
- **Separation of Concerns:** The Broker API and the Trading API are fundamentally different domains. They required dedicated CLI binaries rather than one monolith.
- **Auto-Generation:** Leverage Alpaca's OpenAPI specifications to auto-generate underlying Go client bindings (using `oapi-codegen`), reducing boilerplate and manual maintenance.
- **Unix Philosophy:** Output formats should be pipeable and parsable. The CLI needed to natively support JSON, CSV, and tabular data, allowing users to pipe output directly into tools like `jq`, `grep`, or `awk`.
- **Strict Quality Control:** Every commit must be automatically linted, formatted, and tested to ensure commercial-grade stability.

## 3. The Minimum Viable Product (MVP)
The MVP started as raw API wrappers. Initially, commands relied heavily on passing raw JSON files (e.g., `--file payload.json`) to invoke mutations (like creating an account or placing an order). 
While functional, this proved to be poor UX. The MVP did, however, successfully validate the foundational routing and proved that the Alpaca OpenAPI specs could cleanly compile into native Go packages (which later matured into isolated `internal/` boundaries).

With the core routing functional, AntiGravity pivoted to a structured, iterative development approach.

## 4. Establishing the Roadmap
To transition from a "working script" to a mature product, AntiGravity defined a strict roadmap, executing features in focused, sequential Iterations:

* **Iteration 0: Foundation & Technical Quality** 
  AntiGravity implemented strict CI/CD pipelines (`.github/workflows/ci.yml`), enforced `golangci-lint` (catching unchecked errors, cyclomatic complexity, and security flaws), and upgraded legacy dependencies. AntiGravity replaced `log.Fatal` crashes with graceful Cobra error handling.

* **Iteration 1: Parameterization and Filtering**
  AntiGravity stripped away the reliance on raw JSON payloads. Commands were refactored to use native CLI flags (e.g., `alpaca-trader orders create --symbol AAPL --qty 1 --side buy`). This vastly improved the core user experience.

* **Iteration 2: Feature Completeness**
  AntiGravity mapped out the rest of the Alpaca API surface area. For the Trader CLI, watchlists, corporate actions, and market data fetching were introduced. For the Broker CLI, journals (moving money internally), documents (KYC uploads), and real-time SSE event streaming were added. 

* **Iteration 3: Complex Output & Data Extraction**
  Standard API responses are noisy. AntiGravity implemented a unified `printOutput` pipeline that allows users to seamlessly switch between human-readable tables (`--output table`), raw JSON (`--output json`), or extract specific JSON nodes on the fly using `--query` (powered by `gjson`).

* **Iteration 4: CI/CD & Distribution**
  To make the CLI accessible to non-Go developers, AntiGravity integrated **GoReleaser**. The project now automatically cross-compiles production-ready binaries for macOS, Linux, and Windows, pushing them to GitHub Releases and publishing a Homebrew Tap.

* **Iteration 5: Data Export & Offline Analytics**
  Traders need data offline for tax reconciliation and backtesting. AntiGravity added native `--output csv` formatting globally. To ensure reliability, AntiGravity built comprehensive bash-based E2E verification scripts (`scripts/test-broker-e2e.sh`, `scripts/test-trader-e2e.sh`) to integration-test the live API environments safely.

## 5. Reaching the Current Point: Security & Best Practices Audit
Having completed Iteration 5, AntiGravity paused feature development to conduct a rigorous external Repository Best Practices Audit. AntiGravity implemented:
- **Strict Module Isolation:** Shifted the generated SDK into an isolated `internal/brokerclient` architecture.
- **Commercial Hygiene:** Added a proprietary `LICENSE`, configured Dependabot for weekly automation scans, and drafted inner-source `CONTRIBUTING.md` and `CODE_OF_CONDUCT.md` guidelines.
- **Deep Linting Fixes:** Addressed remaining technical debt, including reducing cyclomatic complexity inside our output formatters and squelching all `errcheck` / `gosec` warnings.
- **Auto-Generated Docs:** Wrote a Cobra Markdown generator (`scripts/gen-docs.go`) to automatically sync CLI help text to the `docs/` folder via `make docs`.

## Next Steps
The repository currently sits at the precipice of **Iteration 6 (Multi-Environment & Secure Credential Management)**. The technical debt is zero, the CI pipeline is air-tight, the test scripts are mapped, and the project is fully commercially compliant. AntiGravity is pausing here to gather external feedback before implementing the credential Keystore and context manager.
