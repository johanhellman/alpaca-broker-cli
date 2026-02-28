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

## Configuration

Both tools accept configuration via environment variables, flags, or a local YAML config file.

### Broker Config (`~/.alpaca-broker.yaml`)
```yaml
api-key: "broker-api-key"
api-secret: "broker-api-secret"
env: "sandbox" # or production
```
Environment Variable Fallbacks: `ALPACA_BROKER_API_KEY`, `ALPACA_BROKER_API_SECRET`, `ALPACA_BROKER_ENV`

### Trader Config (`~/.alpaca-trader.yaml`)
```yaml
api-key: "trader-api-key"
api-secret: "trader-api-secret"
env: "paper" # or live
```
Environment Variable Fallbacks: `APCA_API_KEY_ID`, `APCA_API_SECRET_KEY`, `APCA_ENV`

## Quick Start: Broker API

### Accounts
```bash
alpaca-broker accounts list
alpaca-broker accounts create --file account_payload.json
```

### Funding
```bash
alpaca-broker funding transfers <account_uuid>
```

### Trading (Sub-accounts)
```bash
alpaca-broker trading orders <account_uuid>
```

## Quick Start: Trading API

### Auto-authenticate with Env Vars
```bash
export APCA_API_KEY_ID="yourkey"
export APCA_API_SECRET_KEY="yoursecret"
export APCA_ENV="paper"
```

### Account & Positions
```bash
alpaca-trader account get
alpaca-trader positions list
```

### Orders
```bash
alpaca-trader orders list
alpaca-trader orders create --file order.json
```

## Documentation

Full command documentation for `alpaca-broker` can be found in the `docs/` folder (generated via `make generate`).
