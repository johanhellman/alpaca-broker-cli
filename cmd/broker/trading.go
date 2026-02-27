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

var tradingCmd = &cobra.Command{
	Use:   "trading",
	Short: "Manage trading and orders",
	Long:  `List and create trading orders for a specific account.`,
}

var tradingOrdersCmd = &cobra.Command{
	Use:   "orders <account_id>",
	Short: "List orders for an account",
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
		params := &client.GetOrdersParams{}
		resp, err := c.GetOrdersWithResponse(ctx, parsedUUID, params)
		if err != nil {
			return fmt.Errorf("failed to list orders: %w", err)
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

var orderPayloadFile string

var tradingOrderCreateCmd = &cobra.Command{
	Use:   "order-create <account_id>",
	Short: "Create an order for an account",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		accountIDStr := args[0]
		parsedUUID, err := uuid.Parse(accountIDStr)
		if err != nil {
			return fmt.Errorf("invalid account ID format: %w", err)
		}

		if orderPayloadFile == "" {
			return fmt.Errorf("--file flag is required")
		}

		data, err := os.ReadFile(orderPayloadFile)
		if err != nil {
			return fmt.Errorf("failed to read payload file: %w", err)
		}

		var reqBody client.CreateOrderRequest
		if err := json.Unmarshal(data, &reqBody); err != nil {
			return fmt.Errorf("failed to parse JSON payload: %w", err)
		}

		c, err := api.NewClient()
		if err != nil {
			return err
		}

		ctx := context.Background()
		resp, err := c.PostOrdersWithResponse(ctx, parsedUUID, reqBody)
		if err != nil {
			return fmt.Errorf("failed to submit order: %w", err)
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
	rootCmd.AddCommand(tradingCmd)
	tradingCmd.AddCommand(tradingOrdersCmd)

	tradingOrderCreateCmd.Flags().StringVarP(&orderPayloadFile, "file", "f", "", "Path to the JSON payload file")
	tradingOrderCreateCmd.MarkFlagRequired("file")
	tradingCmd.AddCommand(tradingOrderCreateCmd)
}
