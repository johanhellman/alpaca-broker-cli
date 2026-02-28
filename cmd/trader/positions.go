package cmd

import (
	"fmt"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

var (
	closeAllCancelOrders bool
	closeQty             float64
	closePercentage      float64
)

var positionsCmd = &cobra.Command{
	Use:   "positions",
	Short: "Manage trading positions",
}

var positionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all open positions",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		positions, err := client.GetPositions()
		if err != nil {
			return fmt.Errorf("failed to list positions: %w", err)
		}

		return printOutput(positions)
	},
}

var positionsGetCmd = &cobra.Command{
	Use:   "get <symbol>",
	Short: "Get an open position by symbol",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		symbol := args[0]
		position, err := client.GetPosition(symbol)
		if err != nil {
			return fmt.Errorf("failed to get position for %s: %w", symbol, err)
		}

		return printOutput(position)
	},
}

var positionsCloseAllCmd = &cobra.Command{
	Use:   "close-all",
	Short: "Liquidate all open positions at market price",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		req := alpaca.CloseAllPositionsRequest{
			CancelOrders: closeAllCancelOrders,
		}

		orders, err := client.CloseAllPositions(req)
		if err != nil {
			return fmt.Errorf("failed to close all positions: %w", err)
		}

		return printOutput(orders)
	},
}

var positionsCloseCmd = &cobra.Command{
	Use:   "close <symbol>",
	Short: "Close down a specific open position",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		symbol := args[0]
		var req alpaca.ClosePositionRequest

		if closeQty > 0 {
			qty := decimal.NewFromFloat(closeQty)
			req.Qty = qty
		} else if closePercentage > 0 {
			pct := decimal.NewFromFloat(closePercentage)
			req.Percentage = pct
		}

		order, err := client.ClosePosition(symbol, req)
		if err != nil {
			return fmt.Errorf("failed to close position for %s: %w", symbol, err)
		}

		return printOutput(order)
	},
}

func init() {
	RootCmd.AddCommand(positionsCmd)

	positionsCmd.AddCommand(positionsListCmd)
	positionsCmd.AddCommand(positionsGetCmd)

	positionsCloseAllCmd.Flags().BoolVar(&closeAllCancelOrders, "cancel-orders", false, "Cancel all associated open orders before closing positions")
	positionsCmd.AddCommand(positionsCloseAllCmd)

	positionsCloseCmd.Flags().Float64Var(&closeQty, "qty", 0, "Number of shares to liquidate")
	positionsCloseCmd.Flags().Float64Var(&closePercentage, "percentage", 0, "Percentage of position to liquidate (0.0 - 100.0, e.g. 50 means 50%)")
	positionsCmd.AddCommand(positionsCloseCmd)
}
