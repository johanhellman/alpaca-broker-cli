package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/johanhellman/alpaca-broker-cli/pkg/brokerclient/api"
	"github.com/johanhellman/alpaca-broker-cli/pkg/brokerclient"
	"github.com/spf13/cobra"
)

// accountsCmd represents the accounts command
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

		ctx := context.Background()
		// For MVP, we'll just fetch with empty params.
		params := &client.GetAccountsParams{}
		resp, err := c.GetAccountsWithResponse(ctx, params)
		if err != nil {
			return fmt.Errorf("failed to list accounts: %w", err)
		}

		if resp.JSON200 == nil {
			return fmt.Errorf("unexpected response status: %d", resp.StatusCode())
		}

		out, err := json.MarshalIndent(resp.JSON200, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(out))
		return nil
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

		ctx := context.Background()
		resp, err := c.GetAccountWithResponse(ctx, parsedUUID)
		if err != nil {
			return fmt.Errorf("failed to get account: %w", err)
		}

		if resp.JSON200 == nil {
			return fmt.Errorf("unexpected response status: %d", resp.StatusCode())
		}

		out, err := json.MarshalIndent(resp.JSON200, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(out))
		return nil
	},
}

var payloadFile string

var accountsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new broker account",
	Long: `Create a new broker account from a JSON payload.
Example:
  alpaca-cli accounts create --file payload.json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if payloadFile == "" {
			return fmt.Errorf("--file flag is required")
		}

		data, err := os.ReadFile(payloadFile)
		if err != nil {
			return fmt.Errorf("failed to read payload file: %w", err)
		}

		var reqBody client.AccountCreationRequest
		if err := json.Unmarshal(data, &reqBody); err != nil {
			return fmt.Errorf("failed to parse JSON payload: %w", err)
		}

		c, err := api.NewClient()
		if err != nil {
			return err
		}

		ctx := context.Background()
		resp, err := c.PostAccountsWithResponse(ctx, reqBody)
		if err != nil {
			return fmt.Errorf("failed to create account: %w", err)
		}

		if resp.JSON200 == nil {
			// Print raw body on error
			return fmt.Errorf("unexpected response status: %d, body: %s", resp.StatusCode(), string(resp.Body))
		}

		out, err := json.MarshalIndent(resp.JSON200, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(out))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(accountsCmd)
	accountsCmd.AddCommand(accountsListCmd)
	accountsCmd.AddCommand(accountsGetCmd)
	
	accountsCreateCmd.Flags().StringVarP(&payloadFile, "file", "f", "", "Path to the JSON payload file")
	accountsCreateCmd.MarkFlagRequired("file")
	accountsCmd.AddCommand(accountsCreateCmd)
}
