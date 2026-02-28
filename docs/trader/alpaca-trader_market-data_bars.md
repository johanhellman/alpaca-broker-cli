## alpaca-trader market-data bars

Get historical bars for a symbol

```
alpaca-trader market-data bars <symbol> [flags]
```

### Options

```
      --adjustment string   Adjustment for corporate actions (raw, split, dividend, all) (default "raw")
      --as-of string        Date when the symbols are mapped
      --currency string     Currency of displayed prices
      --end string          Inclusive end of interval (RFC3339)
      --feed string         Source of data: sip, iex, otc
  -h, --help                help for bars
      --page-limit int      Pagination size
      --sort string         Sort direction (asc or desc) (default "asc")
      --start string        Inclusive beginning of interval (RFC3339)
      --timeframe string    Aggregation size (e.g. 1Min, 1Hour, 1Day, 1Week, 1Month) (default "1Day")
      --total-limit int     Total number of items to return (0 means all)
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

* [alpaca-trader market-data](alpaca-trader_market-data.md)	 - Get historical market data

