package cmd

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	client "github.com/johanhellman/alpaca-broker-cli/internal/brokerclient"
	"github.com/johanhellman/alpaca-broker-cli/internal/brokerclient/api"
	"github.com/spf13/cobra"
)

var (
	// List flags
	listOrdersStatus    string
	listOrdersLimit     int
	listOrdersAfter     string
	listOrdersUntil     string
	listOrdersDirection string
	listOrdersNested    bool
	listOrdersSymbols   string

	// Create flags
	createOrderSymbol        string
	createOrderQty           float64
	createOrderNotional      float64
	createOrderSide          string
	createOrderType          string
	createOrderTimeInForce   string
	createOrderLimitPrice    float64
	createOrderStopPrice     float64
	createOrderExtendedHours bool
	createOrderClientOrderID string
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

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		params := &client.GetOrdersParams{}
		if listOrdersStatus != "" {
			st := client.GetOrdersParamsStatus(listOrdersStatus)
			params.Status = &st
		}
		if listOrdersLimit > 0 {
			lim := client.Limit(listOrdersLimit)
			params.Limit = &lim
		}
		if listOrdersDirection != "" {
			dir := client.GetOrdersParamsDirection(listOrdersDirection)
			params.Direction = &dir
		}
		if listOrdersNested {
			nest := client.Nested(listOrdersNested)
			params.Nested = &nest
		}
		if listOrdersSymbols != "" {
			syms := client.Symbols(listOrdersSymbols)
			params.Symbols = &syms
		}
		if listOrdersAfter != "" {
			t, err := time.Parse(time.RFC3339, listOrdersAfter)
			if err != nil {
				return fmt.Errorf("invalid after format (expected RFC3339): %w", err)
			}
			aft := client.After(t)
			params.After = &aft
		}
		if listOrdersUntil != "" {
			t, err := time.Parse(time.RFC3339, listOrdersUntil)
			if err != nil {
				return fmt.Errorf("invalid until format (expected RFC3339): %w", err)
			}
			utl := client.Until(t)
			params.Until = &utl
		}

		resp, err := c.GetOrdersWithResponse(ctx, parsedUUID, params)
		if err != nil {
			return fmt.Errorf("failed to list orders: %w", err)
		}

		if resp.JSON200 == nil {
			return fmt.Errorf("unexpected response status: %d", resp.StatusCode())
		}

		return printOutput(resp.JSON200)
	},
}

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

		reqBody := client.CreateOrderRequest{
			Symbol:      createOrderSymbol,
			Side:        client.CreateOrderRequestSide(createOrderSide),
			Type:        client.CreateOrderRequestType(createOrderType),
			TimeInForce: client.CreateOrderRequestTimeInForce(createOrderTimeInForce),
		}

		if createOrderQty > 0 {
			q := strconv.FormatFloat(createOrderQty, 'f', -1, 64)
			reqBody.Qty = &q
		} else if createOrderNotional > 0 {
			n := strconv.FormatFloat(createOrderNotional, 'f', -1, 64)
			reqBody.Notional = &n
		} else {
			return fmt.Errorf("either --qty or --notional must be specified and greater than 0")
		}

		if createOrderLimitPrice > 0 {
			lp := strconv.FormatFloat(createOrderLimitPrice, 'f', -1, 64)
			reqBody.LimitPrice = &lp
		}
		if createOrderStopPrice > 0 {
			sp := strconv.FormatFloat(createOrderStopPrice, 'f', -1, 64)
			reqBody.StopPrice = &sp
		}
		if createOrderExtendedHours {
			reqBody.ExtendedHours = &createOrderExtendedHours
		}
		if createOrderClientOrderID != "" {
			reqBody.ClientOrderId = &createOrderClientOrderID
		}

		c, err := api.NewClient()
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		resp, err := c.PostOrdersWithResponse(ctx, parsedUUID, reqBody)
		if err != nil {
			return fmt.Errorf("failed to submit order: %w", err)
		}

		if resp.JSON200 == nil {
			return fmt.Errorf("unexpected response status: %d, body: %s", resp.StatusCode(), string(resp.Body))
		}

		return printOutput(resp.JSON200)
	},
}

func init() {
	rootCmd.AddCommand(tradingCmd)

	tradingOrdersCmd.Flags().StringVar(&listOrdersStatus, "status", "open", "Order status to be queried (open, closed, or all)")
	tradingOrdersCmd.Flags().IntVar(&listOrdersLimit, "limit", 50, "The maximum number of orders in response")
	tradingOrdersCmd.Flags().StringVar(&listOrdersAfter, "after", "", "The response will include only ones submitted after this timestamp (RFC3339)")
	tradingOrdersCmd.Flags().StringVar(&listOrdersUntil, "until", "", "The response will include only ones submitted until this timestamp (RFC3339)")
	tradingOrdersCmd.Flags().StringVar(&listOrdersDirection, "direction", "desc", "The chronological order of response based on the submission time (asc or desc)")
	tradingOrdersCmd.Flags().BoolVar(&listOrdersNested, "nested", false, "If true, the result will roll up multi-leg orders under the legs field of primary order")
	tradingOrdersCmd.Flags().StringVar(&listOrdersSymbols, "symbols", "", "A comma-separated list of symbols to filter by")
	tradingCmd.AddCommand(tradingOrdersCmd)

	tradingOrderCreateCmd.Flags().StringVar(&createOrderSymbol, "symbol", "", "Symbol or asset ID to identify the asset to trade")
	_ = tradingOrderCreateCmd.MarkFlagRequired("symbol") //nolint:errcheck
	tradingOrderCreateCmd.Flags().Float64Var(&createOrderQty, "qty", 0, "Number of shares to trade")
	tradingOrderCreateCmd.Flags().Float64Var(&createOrderNotional, "notional", 0, "Dollar amount to trade")
	tradingOrderCreateCmd.Flags().StringVar(&createOrderSide, "side", "", "buy or sell")
	_ = tradingOrderCreateCmd.MarkFlagRequired("side") //nolint:errcheck
	tradingOrderCreateCmd.Flags().StringVar(&createOrderType, "type", "", "market, limit, stop, stop_limit, trailing_stop")
	_ = tradingOrderCreateCmd.MarkFlagRequired("type") //nolint:errcheck
	tradingOrderCreateCmd.Flags().StringVar(&createOrderTimeInForce, "time-in-force", "", "day, gtc, opg, cls, ioc, fok")
	_ = tradingOrderCreateCmd.MarkFlagRequired("time-in-force") //nolint:errcheck
	tradingOrderCreateCmd.Flags().Float64Var(&createOrderLimitPrice, "limit-price", 0, "Required if type is limit or stop_limit")
	tradingOrderCreateCmd.Flags().Float64Var(&createOrderStopPrice, "stop-price", 0, "Required if type is stop or stop_limit")
	tradingOrderCreateCmd.Flags().BoolVar(&createOrderExtendedHours, "extended-hours", false, "Whether or not this order should be allowed to execute during extended hours")
	tradingOrderCreateCmd.Flags().StringVar(&createOrderClientOrderID, "client-order-id", "", "A unique identifier for the order")
	tradingCmd.AddCommand(tradingOrderCreateCmd)
}
