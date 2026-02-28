package cmd

import (
	"fmt"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

var (
	// List flags
	listStatus    string
	listLimit     int
	listAfter     string
	listUntil     string
	listDirection string
	listNested    bool
	listSymbols   []string

	// Create flags
	createSymbol        string
	createQty           float64
	createNotional      float64
	createSide          string
	createType          string
	createTimeInForce   string
	createLimitPrice    float64
	createStopPrice     float64
	createExtendedHours bool
	createClientOrderID string
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

		req := alpaca.GetOrdersRequest{
			Status:    listStatus,
			Limit:     listLimit,
			Direction: listDirection,
			Nested:    listNested,
			Symbols:   listSymbols,
		}

		if listAfter != "" {
			t, err := time.Parse(time.RFC3339, listAfter)
			if err != nil {
				return fmt.Errorf("invalid after format (expected RFC3339): %w", err)
			}
			req.After = t
		}
		if listUntil != "" {
			t, err := time.Parse(time.RFC3339, listUntil)
			if err != nil {
				return fmt.Errorf("invalid until format (expected RFC3339): %w", err)
			}
			req.Until = t
		}

		orders, err := client.GetOrders(req)
		if err != nil {
			return fmt.Errorf("failed to list orders: %w", err)
		}

		return printOutput(orders)
	},
}

var ordersCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an order",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		req := alpaca.PlaceOrderRequest{
			Symbol:        createSymbol,
			Side:          alpaca.Side(createSide),
			Type:          alpaca.OrderType(createType),
			TimeInForce:   alpaca.TimeInForce(createTimeInForce),
			ExtendedHours: createExtendedHours,
			ClientOrderID: createClientOrderID,
		}

		if createQty > 0 {
			qty := decimal.NewFromFloat(createQty)
			req.Qty = &qty
		} else if createNotional > 0 {
			notional := decimal.NewFromFloat(createNotional)
			req.Notional = &notional
		} else {
			return fmt.Errorf("either --qty or --notional must be specified and greater than 0")
		}

		if createLimitPrice > 0 {
			limit := decimal.NewFromFloat(createLimitPrice)
			req.LimitPrice = &limit
		}
		if createStopPrice > 0 {
			stop := decimal.NewFromFloat(createStopPrice)
			req.StopPrice = &stop
		}

		order, err := client.PlaceOrder(req)
		if err != nil {
			return fmt.Errorf("failed to place order: %w", err)
		}

		return printOutput(order)
	},
}

func init() {
	RootCmd.AddCommand(ordersCmd)

	// orders list flags
	ordersListCmd.Flags().StringVar(&listStatus, "status", "open", "Order status to be queried (open, closed, or all)")
	ordersListCmd.Flags().IntVar(&listLimit, "limit", 50, "The maximum number of orders in response")
	ordersListCmd.Flags().StringVar(&listAfter, "after", "", "The response will include only ones submitted after this timestamp (RFC3339)")
	ordersListCmd.Flags().StringVar(&listUntil, "until", "", "The response will include only ones submitted until this timestamp (RFC3339)")
	ordersListCmd.Flags().StringVar(&listDirection, "direction", "desc", "The chronological order of response based on the submission time (asc or desc)")
	ordersListCmd.Flags().BoolVar(&listNested, "nested", false, "If true, the result will roll up multi-leg orders under the legs field of primary order")
	ordersListCmd.Flags().StringSliceVar(&listSymbols, "symbols", nil, "A comma-separated list of symbols to filter by")
	ordersCmd.AddCommand(ordersListCmd)

	// orders create flags
	ordersCreateCmd.Flags().StringVar(&createSymbol, "symbol", "", "Symbol or asset ID to identify the asset to trade")
	_ = ordersCreateCmd.MarkFlagRequired("symbol") //nolint:errcheck
	ordersCreateCmd.Flags().Float64Var(&createQty, "qty", 0, "Number of shares to trade")
	ordersCreateCmd.Flags().Float64Var(&createNotional, "notional", 0, "Dollar amount to trade")
	ordersCreateCmd.Flags().StringVar(&createSide, "side", "", "buy or sell")
	_ = ordersCreateCmd.MarkFlagRequired("side") //nolint:errcheck
	ordersCreateCmd.Flags().StringVar(&createType, "type", "", "market, limit, stop, stop_limit, trailing_stop")
	_ = ordersCreateCmd.MarkFlagRequired("type") //nolint:errcheck
	ordersCreateCmd.Flags().StringVar(&createTimeInForce, "time-in-force", "", "day, gtc, opg, cls, ioc, fok")
	_ = ordersCreateCmd.MarkFlagRequired("time-in-force") //nolint:errcheck
	ordersCreateCmd.Flags().Float64Var(&createLimitPrice, "limit-price", 0, "Required if type is limit or stop_limit")
	ordersCreateCmd.Flags().Float64Var(&createStopPrice, "stop-price", 0, "Required if type is stop or stop_limit")
	ordersCreateCmd.Flags().BoolVar(&createExtendedHours, "extended-hours", false, "Whether or not this order should be allowed to execute during extended hours")
	ordersCreateCmd.Flags().StringVar(&createClientOrderID, "client-order-id", "", "A unique identifier for the order")
	ordersCmd.AddCommand(ordersCreateCmd)
}
