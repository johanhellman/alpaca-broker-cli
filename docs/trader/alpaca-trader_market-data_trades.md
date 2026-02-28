## alpaca-trader market-data trades

Get historical trades for a symbol

```
alpaca-trader market-data trades <symbol> [flags]
```

### Options

```
      --as-of string      Date when the symbols are mapped
      --currency string   Currency of displayed prices
      --end string        Inclusive end of interval (RFC3339)
      --feed string       Source of data: sip, iex, otc
  -h, --help              help for trades
      --page-limit int    Pagination size
      --sort string       Sort direction (asc or desc) (default "asc")
      --start string      Inclusive beginning of interval (RFC3339)
      --total-limit int   Total number of items to return (0 means all)
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

