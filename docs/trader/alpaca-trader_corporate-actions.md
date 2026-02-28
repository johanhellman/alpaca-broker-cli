## alpaca-trader corporate-actions

Get corporate actions

```
alpaca-trader corporate-actions [flags]
```

### Options

```
      --end string        Inclusive end of the interval (YYYY-MM-DD)
  -h, --help              help for corporate-actions
      --page-limit int    Pagination size
      --sort string       Sort direction (asc or desc) (default "asc")
      --start string      Inclusive beginning of the interval (YYYY-MM-DD)
      --symbols strings   Comma-separated list of company symbols
      --total-limit int   Limit of the total number of actions returned (0 means all)
      --types strings     Comma-separated list of corporate actions types (e.g. forward_split, cash_dividend)
```

### Options inherited from parent commands

```
      --all                 Automatically fetch all pages for list endpoints
      --api-key string      Alpaca API Key ID
      --api-secret string   Alpaca API Secret Key
      --config string       config file (default is $HOME/.alpaca-trader.yaml)
      --env string          Alpaca environment (paper or live) (default "paper")
      --output string       Output format (table, json, or csv) (default "table")
      --query string        Filter output using jq-like syntax (forces json output if used)
```

### SEE ALSO

* [alpaca-trader](alpaca-trader.md)	 - A CLI tool for the Alpaca Trading API

