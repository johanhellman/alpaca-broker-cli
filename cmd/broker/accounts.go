package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	client "github.com/johanhellman/alpaca-broker-cli/internal/brokerclient"
	"github.com/johanhellman/alpaca-broker-cli/internal/brokerclient/api"
	"github.com/oapi-codegen/runtime/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)



var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Manage Alpaca Broker accounts",
	Long:  `Create, list, and get details about Alpaca Broker accounts.`,
}

var accountsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all accounts",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := api.NewClient()
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		listQuery, _ := cmd.Flags().GetString("query")
		params := &client.GetAccountsParams{}
		if listQuery != "" {
			params.Query = &listQuery
		}

		var allAccounts []client.Account
		fetchAll := viper.GetBool("all")

		for {
			resp, err := c.GetAccountsWithResponse(ctx, params)
			if err != nil {
				return fmt.Errorf("failed to list accounts: %w", err)
			}

			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected response status: %d", resp.StatusCode())
			}

			allAccounts = append(allAccounts, *resp.JSON200...)

			if !fetchAll || resp.JSON200 == nil || len(*resp.JSON200) < 100 {
				break
			}

			// Warning: Simple naive pagination logic used. Implement explicit NextPageToken for large broker accounts if needed.
		}

		return printOutput(allAccounts)
	},
}

var accountsGetCmd = &cobra.Command{
	Use:   "get <account_id>",
	Short: "Get account details by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		accountIDStr := args[0]
		parsedUUID, err := uuid.Parse(accountIDStr)
		if err != nil {
			return fmt.Errorf("invalid account ID format: %w", err)
		}

		c, err := api.NewClient()
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		resp, err := c.GetAccountWithResponse(ctx, parsedUUID)
		if err != nil {
			return fmt.Errorf("failed to get account: %w", err)
		}

		if resp.JSON200 == nil {
			return fmt.Errorf("unexpected response status: %d", resp.StatusCode())
		}

		return printOutput(resp.JSON200)
	},
}

var accountsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new broker account",
	Long: `Create a new broker account using native CLI parameters.
Example:
  alpaca-broker accounts create \
    --id-given-name John \
    --id-family-name Doe \
    --id-dob 1990-01-01 \
    --contact-email john.doe@example.com`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := api.NewClient()
		if err != nil {
			return err
		}

		
		createIDDOB, _ := cmd.Flags().GetString("id-dob")
		createContactEmail, _ := cmd.Flags().GetString("contact-email")
		createContactPhone, _ := cmd.Flags().GetString("contact-phone")
		createContactCity, _ := cmd.Flags().GetString("contact-city")
		createContactState, _ := cmd.Flags().GetString("contact-state")
		createContactPostalCode, _ := cmd.Flags().GetString("contact-postal")
		createContactStreet1, _ := cmd.Flags().GetString("contact-street-1")
		createContactStreet2, _ := cmd.Flags().GetString("contact-street-2")
		createIDGivenName, _ := cmd.Flags().GetString("id-given-name")
		createIDFamilyName, _ := cmd.Flags().GetString("id-family-name")
		createIDTaxID, _ := cmd.Flags().GetString("id-tax-id")
		createIDCountryTax, _ := cmd.Flags().GetString("id-country-tax")
		createIDTaxType, _ := cmd.Flags().GetString("id-tax-type")
		createIDCountryBirth, _ := cmd.Flags().GetString("id-country-birth")
		createIDCountryCitizen, _ := cmd.Flags().GetString("id-country-citizen")
		createFundingSourceType, _ := cmd.Flags().GetString("funding-source-type")
	dobTime, err := time.Parse("2006-01-02", createIDDOB)
		if err != nil {
			return fmt.Errorf("invalid DOB format (expected YYYY-MM-DD): %w", err)
		}

		email := types.Email(createContactEmail)

		contact := client.Contact{
			EmailAddress: &email,
			PhoneNumber:  &createContactPhone,
			City:         &createContactCity,
			State:        &createContactState,
			PostalCode:   &createContactPostalCode,
		}

		if createContactStreet1 != "" {
			contact.StreetAddress = &[]client.StreetAddress{
				createContactStreet1,
			}
			if createContactStreet2 != "" {
				s := append(*contact.StreetAddress, createContactStreet2)
				contact.StreetAddress = &s
			}
		}

		identity := client.Identity{
			GivenName:             createIDGivenName,
			FamilyName:            createIDFamilyName,
			DateOfBirth:           types.Date{Time: dobTime},
			TaxId:                 &createIDTaxID,
			CountryOfTaxResidence: createIDCountryTax,
		}

		if createIDTaxType != "" {
			taxType := client.IdentityTaxIdType(createIDTaxType)
			identity.TaxIdType = &taxType
		}

		if createIDCountryBirth != "" {
			identity.CountryOfBirth = &createIDCountryBirth
		}
		if createIDCountryCitizen != "" {
			identity.CountryOfCitizenship = &createIDCountryCitizen
		}

		// Optional barebones funding source to bypass basic validation failures
		if createFundingSourceType != "" {
			identity.FundingSource = []client.IdentityFundingSource{
				client.IdentityFundingSource(createFundingSourceType),
			}
		} else {
			// Provide empty array instead of nil
			identity.FundingSource = []client.IdentityFundingSource{}
		}

		agreements := client.Agreements{
			{
				Agreement: client.AgreementAgreement("account_agreement"),
				IpAddress: "127.0.0.1",
				SignedAt:  time.Now().Format(time.RFC3339),
			},
			{
				Agreement: client.AgreementAgreement("customer_agreement"),
				IpAddress: "127.0.0.1",
				SignedAt:  time.Now().Format(time.RFC3339),
			},
			{
				Agreement: client.AgreementAgreement("margin_agreement"),
				IpAddress: "127.0.0.1",
				SignedAt:  time.Now().Format(time.RFC3339),
			},
		}

		disclosures := client.Disclosures{
			IsControlPerson:             false,
			IsAffiliatedExchangeOrFinra: false,
			IsPoliticallyExposed:        false,
			ImmediateFamilyExposed:      false,
		}

		reqBody := client.AccountCreationRequest{
			Contact:     &contact,
			Identity:    &identity,
			Agreements:  &agreements,
			Disclosures: &disclosures,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		resp, err := c.PostAccountsWithResponse(ctx, reqBody)
		if err != nil {
			return fmt.Errorf("failed to create account: %w", err)
		}

		if resp.JSON200 == nil {
			// Print raw body on error
			return fmt.Errorf("unexpected response status: %d, body: %s", resp.StatusCode(), string(resp.Body))
		}

		return printOutput(resp.JSON200)
	},
}

func init() {
	rootCmd.AddCommand(accountsCmd)

	// list flags
	accountsListCmd.Flags().String("query", "", "Space-delimited partial match search (names, emails, account numbers)")
	accountsCmd.AddCommand(accountsListCmd)

	accountsCmd.AddCommand(accountsGetCmd)

	// create flags
	accountsCreateCmd.Flags().String("contact-email", "", "Contact email address")
	_ = accountsCreateCmd.MarkFlagRequired("contact-email") //nolint:errcheck
	accountsCreateCmd.Flags().String("contact-phone", "", "Contact phone number (with country code, no hyphens)")
	accountsCreateCmd.Flags().String("contact-street-1", "", "Contact street address line 1")
	accountsCreateCmd.Flags().String("contact-street-2", "", "Contact street address line 2")
	accountsCreateCmd.Flags().String("contact-city", "", "Contact city")
	accountsCreateCmd.Flags().String("contact-state", "", "Contact state/province")
	accountsCreateCmd.Flags().String("contact-postal", "", "Contact postal code")

	accountsCreateCmd.Flags().String("id-given-name", "", "Account owner's first name")
	_ = accountsCreateCmd.MarkFlagRequired("id-given-name") //nolint:errcheck
	accountsCreateCmd.Flags().String("id-family-name", "", "Account owner's last name")
	_ = accountsCreateCmd.MarkFlagRequired("id-family-name") //nolint:errcheck
	accountsCreateCmd.Flags().String("id-dob", "", "Date of birth (YYYY-MM-DD)")
	_ = accountsCreateCmd.MarkFlagRequired("id-dob") //nolint:errcheck
	accountsCreateCmd.Flags().String("id-tax-id", "", "Tax ID / SSN")
	accountsCreateCmd.Flags().String("id-tax-type", "", "USA_SSN, AUS_TFN, etc.")
	accountsCreateCmd.Flags().String("id-country-tax", "USA", "ISO 3166-1 alpha-3 country of tax residence")
	accountsCreateCmd.Flags().String("id-country-birth", "", "ISO 3166-1 alpha-3 country of birth")
	accountsCreateCmd.Flags().String("id-country-citizen", "", "ISO 3166-1 alpha-3 country of citizenship")

	accountsCreateCmd.Flags().String("funding-source-type", "", "employment_income, investments, etc.")

	accountsCmd.AddCommand(accountsCreateCmd)
}
