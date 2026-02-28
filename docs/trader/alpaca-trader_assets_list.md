## alpaca-trader assets list

List available assets

```
alpaca-trader assets list [flags]
```

### Options

```
      --asset-class string   us_equity or crypto (default "us_equity")
      --exchange string      AMEX, ARCA, BATS, NYSE, NASDAQ, NYSEARCA, OTC
  -h, --help                 help for list
      --status string        active or inactive (default "active")
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

* [alpaca-trader assets](alpaca-trader_assets.md)	 - Manage tradable assets

