## alpaca-cli events accounts

Stream account status events

```
alpaca-cli events accounts [flags]
```

### Options

```
  -h, --help           help for accounts
      --since string   Filter events after this time (RFC3339)
      --since-id int   Filter events after this ID
      --until string   Filter events before this time (RFC3339)
      --until-id int   Filter events before this ID
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

* [alpaca-cli events](alpaca-cli_events.md)	 - Stream SSE events from the broker API

