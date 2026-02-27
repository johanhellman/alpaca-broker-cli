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

var fundingCmd = &cobra.Command{
	Use:   "funding",
	Short: "Manage funding and transfers",
	Long:  `List and create ACH relationships and bank transfers for an account.`,
}

var fundingTransfersCmd = &cobra.Command{
	Use:   "transfers <account_id>",
	Short: "List transfers for an account",
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
		params := &client.GetTransfersParams{}
		resp, err := c.GetTransfersWithResponse(ctx, parsedUUID, params)
		if err != nil {
			return fmt.Errorf("failed to list transfers: %w", err)
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

var transferPayloadFile string

var fundingTransferCreateCmd = &cobra.Command{
	Use:   "transfer-create <account_id>",
	Short: "Create a transfer for an account",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		accountIDStr := args[0]
		parsedUUID, err := uuid.Parse(accountIDStr)
		if err != nil {
			return fmt.Errorf("invalid account ID format: %w", err)
		}

		if transferPayloadFile == "" {
			return fmt.Errorf("--file flag is required")
		}

		data, err := os.ReadFile(transferPayloadFile)
		if err != nil {
			return fmt.Errorf("failed to read payload file: %w", err)
		}

		var reqBody client.TransferData
		if err := json.Unmarshal(data, &reqBody); err != nil {
			return fmt.Errorf("failed to parse JSON payload: %w", err)
		}

		c, err := api.NewClient()
		if err != nil {
			return err
		}

		ctx := context.Background()
		resp, err := c.PostTransfersWithResponse(ctx, parsedUUID, reqBody)
		if err != nil {
			return fmt.Errorf("failed to create transfer: %w", err)
		}

		if resp.JSON200 == nil {
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
	rootCmd.AddCommand(fundingCmd)
	fundingCmd.AddCommand(fundingTransfersCmd)

	fundingTransferCreateCmd.Flags().StringVarP(&transferPayloadFile, "file", "f", "", "Path to the JSON payload file")
	fundingTransferCreateCmd.MarkFlagRequired("file")
	fundingCmd.AddCommand(fundingTransferCreateCmd)
}
