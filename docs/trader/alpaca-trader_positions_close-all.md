## alpaca-trader positions close-all

Liquidate all open positions at market price

```
alpaca-trader positions close-all [flags]
```

### Options

```
      --cancel-orders   Cancel all associated open orders before closing positions
  -h, --help            help for close-all
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

