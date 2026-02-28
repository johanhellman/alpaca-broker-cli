## alpaca-cli accounts create

Create a new broker account

### Synopsis

Create a new broker account using native CLI parameters.
Example:
  alpaca-broker accounts create \
    --id-given-name John \
    --id-family-name Doe \
    --id-dob 1990-01-01 \
    --contact-email john.doe@example.com

```
alpaca-cli accounts create [flags]
```

### Options

```
      --contact-city string          Contact city
      --contact-email string         Contact email address
      --contact-phone string         Contact phone number (with country code, no hyphens)
      --contact-postal string        Contact postal code
      --contact-state string         Contact state/province
      --contact-street-1 string      Contact street address line 1
      --contact-street-2 string      Contact street address line 2
      --funding-source-type string   employment_income, investments, etc.
  -h, --help                         help for create
      --id-country-birth string      ISO 3166-1 alpha-3 country of birth
      --id-country-citizen string    ISO 3166-1 alpha-3 country of citizenship
      --id-country-tax string        ISO 3166-1 alpha-3 country of tax residence (default "USA")
      --id-dob string                Date of birth (YYYY-MM-DD)
      --id-family-name string        Account owner's last name
      --id-given-name string         Account owner's first name
      --id-tax-id string             Tax ID / SSN
      --id-tax-type string           USA_SSN, AUS_TFN, etc.
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

* [alpaca-cli accounts](alpaca-cli_accounts.md)	 - Manage Alpaca Broker accounts

