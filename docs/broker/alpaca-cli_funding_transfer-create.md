## alpaca-cli funding transfer-create

Create a transfer for an account

```
alpaca-cli funding transfer-create <account_id> [flags]
```

### Options

```
      --additional-info string   Optional for wire transfers
      --amount string            Amount as a string, e.g. 500.00
      --bank-id string           Required for wire transfers
      --direction string         INCOMING or OUTGOING
  -h, --help                     help for transfer-create
      --relationship-id string   Required for ACH transfers
      --transfer-type string     ach or wire
```

### Options inherited from parent commands

```
      --all                 Automatically fetch all pages for list endpoints
      --api-key string      Alpaca Broker API Key
      --api-secret string   Alpaca Broker API Secret
      --config string       config file (default is $HOME/.alpaca-cli.yaml)
      --env string          Alpaca environment (sandbox or production) (default "sandbox")
      --output string       Output format (table, json, or csv) (default "table")
      --query string        Filter output using jq-like syntax (forces json output if used)
```

### SEE ALSO

* [alpaca-cli funding](alpaca-cli_funding.md)	 - Manage funding and transfers

