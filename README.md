# Alpaca Broker CLI

`alpaca-cli` is a command-line interface for the [Alpaca Broker API](https://docs.alpaca.markets/docs/about-broker-api). It allows you to quickly manage accounts, transfers, and orders using a seamless terminal experience.

## Installation

This project requires Go. To install locally:

```bash
git clone https://github.com/johanhellman/alpaca-broker-cli.git
cd alpaca-broker-cli
make install
```

Ensure `$GOPATH/bin` is in your `PATH`.

## Configuration

You can provide your API credentials via command-line flags, environment variables, or a YAML configuration file (`~/.alpaca-cli.yaml`).

**Environment Variables:**
```bash
export ALPACA_BROKER_API_KEY="your-api-key"
export ALPACA_BROKER_API_SECRET="your-api-secret"
export ALPACA_BROKER_ENV="sandbox" # or "production"
```

**Config File (`~/.alpaca-cli.yaml`):**
```yaml
api-key: "your-api-key"
api-secret: "your-api-secret"
env: "sandbox"
```

## Quick Start

### Accounts

List all accounts:
```bash
alpaca-cli accounts list
```

Get a specific account:
```bash
alpaca-cli accounts get <account_uuid>
```

Create a new account:
```bash
alpaca-cli accounts create --file account_payload.json
```

### Funding

List transfers for an account:
```bash
alpaca-cli funding transfers <account_uuid>
```

Create a transfer:
```bash
alpaca-cli funding transfer-create <account_uuid> --file transfer_payload.json
```

### Trading

List orders for an account:
```bash
alpaca-cli trading orders <account_uuid>
```

Create an order:
```bash
alpaca-cli trading order-create <account_uuid> --file order_payload.json
```

## Documentation

Full command documentation can be found in the [docs/](./docs/) directory.
