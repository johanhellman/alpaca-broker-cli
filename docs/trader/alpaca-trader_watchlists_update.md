## alpaca-trader watchlists update

Update a watchlist

```
alpaca-trader watchlists update <watchlist_id> [flags]
```

### Options

```
  -h, --help              help for update
      --name string       New name for the watchlist
      --symbols strings   New list of symbols for the watchlist
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

* [alpaca-trader watchlists](alpaca-trader_watchlists.md)	 - Manage watchlists

