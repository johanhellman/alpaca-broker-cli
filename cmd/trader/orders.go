package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/spf13/cobra"
)

var ordersCmd = &cobra.Command{
	Use:   "orders",
	Short: "Manage trading orders",
}

var ordersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List orders",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		status := "open"
		req := alpaca.GetOrdersRequest{
			Status: status,
		}

		orders, err := client.GetOrders(req)
		if err != nil {
			return fmt.Errorf("failed to list orders: %w", err)
		}

		out, err := json.MarshalIndent(orders, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(out))
		return nil
	},
}

var orderPayloadFile string

var ordersCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an order",
	RunE: func(cmd *cobra.Command, args []string) error {
		if orderPayloadFile == "" {
			return fmt.Errorf("--file flag is required")
		}

		cleanFile := filepath.Clean(orderPayloadFile)
		data, err := os.ReadFile(cleanFile)
		if err != nil {
			return fmt.Errorf("failed to read payload file: %w", err)
		}

		var reqBody alpaca.PlaceOrderRequest
		if err := json.Unmarshal(data, &reqBody); err != nil {
			return fmt.Errorf("failed to parse JSON payload: %w", err)
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		order, err := client.PlaceOrder(reqBody)
		if err != nil {
			return fmt.Errorf("failed to place order: %w", err)
		}

		out, err := json.MarshalIndent(order, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(out))
		return nil
	},
}

func init() {
	RootCmd.AddCommand(ordersCmd)
	ordersCmd.AddCommand(ordersListCmd)

	ordersCreateCmd.Flags().StringVarP(&orderPayloadFile, "file", "f", "", "Path to the JSON payload file")
	ordersCreateCmd.MarkFlagRequired("file") //nolint:errcheck
	ordersCmd.AddCommand(ordersCreateCmd)
}
