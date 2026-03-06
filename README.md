# Alpaca CLI Tools

A suite of powerful command-line interfaces for interacting with the Alpaca API ecosystem. This repository contains two distinct binaries:
1. `alpaca-broker` - For managing B2B broker operations (accounts, funding, sub-account trading) via the Broker API.
2. `alpaca-trader` - For managing individual retail/paper trading (account info, positions, orders) via the Trading API.

## Installation

### Via Homebrew (macOS & Linux)

The easiest way to install the tools is through the custom Homebrew tap:

```bash
brew install johanhellman/tap/alpaca-cli
```

### Via Go Toolchain (Developers)

This project requires Go 1.26+. To install locally from source:

```bash
git clone https://github.com/johanhellman/alpaca-broker-cli.git
cd alpaca-broker-cli
make install
```

This will build and install both `alpaca-broker` and `alpaca-trader` into your `$GOPATH/bin`. Ensure `$GOPATH/bin` is in your `$PATH`.

### Pre-compiled Binaries

You can also download pre-compiled `.tar.gz` or `.zip` archives mapped to your specific OS and CPU Architecture from the [GitHub Releases page](https://github.com/johanhellman/alpaca-broker-cli/releases).

## Global Configuration & Output Flags

Both tools accept configuration via environment variables, flags, or local YAML config files (`~/.alpaca-broker.yaml` or `~/.alpaca-trader.yaml`).

### Common Global Flags
To maximize composability with standard Unix tools (like `jq`, `grep`, `awk`), every command natively supports advanced output formatting via global flags.

| Flag | Feature | Example Usage |
|------|---------|---------------|
| `--output` | Renders output as `table`, `json`, or `csv`. | `alpaca-trader positions list --output csv > pf.csv` |
| `--query` | Extracts specific nodes using `gjson` syntax (forces JSON). | `alpaca-broker accounts list --query "0.account_number"` |
| `--all` | Bypasses pagination to pull all available records. | `alpaca-trader orders list --all` |
| `--env` | Overrides the environment (`sandbox`, `paper`, `live`). | `alpaca-broker funding list --env sandbox` |
| `--api-key` | Overrides the active API Key. | `alpaca-trader account get --api-key <KEY>` |
| `--api-secret` | Overrides the active API Secret. | `alpaca-trader account get --api-secret <SECRET>` |

### Authentication Fallbacks

**Broker Environment Fallbacks:**
```bash
export ALPACA_BROKER_API_KEY="your-broker-key"
export ALPACA_BROKER_API_SECRET="your-broker-secret"
export ALPACA_BROKER_ENV="sandbox" # or "production"
```

**Trader Environment Fallbacks:**
```bash
export APCA_API_KEY_ID="your-trader-key"
export APCA_API_SECRET_KEY="your-trader-secret"
export APCA_ENV="paper" # or "live"
```

## Complete Command Hierarchies

### `alpaca-broker` commands
B2B logic to manage retail sub-accounts, fund flows, and execute trades on their behalf.

* **`accounts`**: Create, edit, delete, and list underlying retail customer accounts.
* **`documents`**: Upload KYC identity documents for strict customer validation.
* **`funding`**: Execute ACH/Wire transfers and manage banking relationships.
* **`journals`**: Move capital internally between broker-tier and sub-tier accounts.
* **`trading`**: Execute trades (`orders`) on behalf of individual sub-accounts.
* **`events`**: Taps into SSE data-streams to monitor real-time lifecycle events (order fill status, journal approvals, etc).

**Broker Example**:
```bash
# Extract the first account number from the sandbox
alpaca-broker accounts list --all --query "0.account_number"

# Stream live lifecycle events
alpaca-broker events stream
```

### `alpaca-trader` commands
Retail workflows governing an individual's specific live or paper account.

* **`account`**: Retrieve core meta-data like buying power and equity.
* **`positions`**: List current active holdings.
* **`orders`**: Submit `market`, `limit`, `stop`, `trailing_stop`, and fractional trades.
* **`watchlists`**: Create, read, update, or delete lists of tracked assets.
* **`assets`**: Look up standard symbology, fractionability, and margin requirements.
* **`corporate-actions`**: Monitor dividends and splits affecting the portfolio.

**Trader Example**:
```bash
# Buy 1 share of Apple
alpaca-trader orders create --symbol AAPL --qty 1 --side buy --type market --time_in_force day

# Export all active positions to CSV format
alpaca-trader positions list --output csv > positions.csv
```

## E2E Testing (Sanity Checks)

If you are a contributor working on the CLI binaries, you can run comprehensive End-to-End (E2E) verification scripts located in the `scripts/` directory. These testing boundaries and their continuous maintenance against live API shifts are managed by AntiGravity.

### Trader CLI Tests (Paper API)
Automatically tests Accounts, Assets, Market Data, Watchlists, and complex Order permutations (Limit, Market, Fractional). 

> **Important**: You **must** provide Alpaca Paper API credentials to run this script. The script aggressively enforces `APCA_ENV="paper"` to prevent accidental live-market mutations.

```bash
export APCA_API_KEY_ID="your_paper_key"
export APCA_API_SECRET_KEY="your_paper_secret"
./scripts/test-trader-e2e.sh
```

### Broker CLI Tests (Sandbox API)
Simulates a full B2B flow: Sub-Account Creation $\rightarrow$ Bank Funding/Journaling $\rightarrow$ Sub-Account Order execution.

```bash
export ALPACA_BROKER_API_KEY="your_broker_sandbox_key"
export ALPACA_BROKER_API_SECRET="your_broker_sandbox_secret"
./scripts/test-broker-e2e.sh
```

## Documentation

Full command documentation for both CLI binaries can be found in the `docs/` folder:
- `alpaca-broker` docs are located in `docs/broker/`
- `alpaca-trader` docs are located in `docs/trader/`

These files are auto-generated. If you add new commands or flags, update them by running:
```bash
make docs
```
