# AI Agent Instructions (AGENTS.md)

Welcome! If you are an AI assistant or coding agent working on this repository, please read these instructions before starting any tasks.

## 🏛 Project Architecture
This repository contains **two distinct CLI executables** built in Go using the `cobra` and `viper` frameworks:

1. **`alpaca-broker`** (Location: `cmd/broker/`)
   - **Purpose**: Interacting with the Alpaca Broker API (B2B multi-tenant architecture).
   - **Client Implementation**: We use an auto-generated client from a canonical OpenAPI specification (`api/openapi.yaml`). The Go client is generated using `oapi-codegen` and placed in `pkg/brokerclient/`.
   - **Authentication**: Uses `ALPACA_BROKER_API_KEY` and `ALPACA_BROKER_API_SECRET`.

2. **`alpaca-trader`** (Location: `cmd/trader/`)
   - **Purpose**: Interacting with the Alpaca Trading API (Retail/Paper individual accounts).
   - **Client Implementation**: We use the official Alpaca Trading Go SDK (`github.com/alpacahq/alpaca-trade-api-go/v3`). **Do not use OpenAPI generation for the trader CLI.**
   - **Authentication**: Uses standard Alpaca environment variables: `APCA_API_KEY_ID`, `APCA_API_SECRET_KEY`, and `APCA_ENV`.

## 📜 Key Files to Review
- **`ROADMAP.md`**: Contains the planned iterative roadmap mapping the journey from MVP to a production-grade toolset. You must review this document to understand current priorities and future architectural goals.
- **`Makefile`**: Contains all build, install, and code generation scripts. Use `make install` to compile and place both binaries directly into the user's `$(go env GOPATH)/bin` directory.

## 🛠 Important Workflows
- **Docs Generation**: We auto-generate markdown documentation for both CLI command trees inside the `docs/` folder using Cobra. If you add or modify commands, you **must** update the docs by running:
  ```bash
  go run tools/broker/gendocs.go && go run tools/trader/gendocs_trader.go
  ```
- **Broker API Re-generation**: If the `api/openapi.yaml` file is updated with a newer version from Alpaca, run `make generate` to natively recreate the Go client in `pkg/brokerclient/`.
- **CI Synchronization**: The local `make ci` target must always be kept in sync with the jobs defined in `.github/workflows/ci.yml`. Run `make ci` locally to validate code before pushing.
- **Pushing and Validating**: After executing `git push`, you **must** use the GitHub CLI to monitor the CI pipeline and ensure it succeeds before moving on:
  ```bash
  gh run watch
  ```

## ✅ Coding Standards
- Do not introduce new third-party CLI or Configuration frameworks; rigorously stick to `cobra` and `viper`.
- Prefer implementing explicit, native command-line flags (e.g., `--limit 50`, `--status ACTIVE`) rather than requiring users to pass raw JSON payload files (`--file`) whenever implementing new REST API bindings.
- **Error Handling**: NEVER use `log.Fatal` or `panic` in library/helper functions or command executions. Always return errors up the stack to be handled gracefully by Cobra's `RunE`.
- **API Requests**: Always use `context.WithTimeout` when making calls to the Alpaca API; never use an unbounded `context.Background()`.
- **Output Formatting**: Ensure all CLI outputs dynamically respect the global `--output` flag (e.g., formatting as a clean table vs raw JSON). Avoid hardcoding massive JSON dumps.
