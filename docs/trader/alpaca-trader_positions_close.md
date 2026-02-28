## alpaca-trader positions close

Close down a specific open position

```
alpaca-trader positions close <symbol> [flags]
```

### Options

```
  -h, --help               help for close
      --percentage float   Percentage of position to liquidate (0.0 - 100.0, e.g. 50 means 50%)
      --qty float          Number of shares to liquidate
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

* [alpaca-trader positions](alpaca-trader_positions.md)	 - Manage trading positions

