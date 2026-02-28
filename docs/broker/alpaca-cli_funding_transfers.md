## alpaca-cli funding transfers

List transfers for an account

```
alpaca-cli funding transfers <account_id> [flags]
```

### Options

```
      --direction string   INCOMING or OUTGOING
  -h, --help               help for transfers
      --limit int32        Maximum number of transfers to return (default 50)
      --offset int32       Pagination offset
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

