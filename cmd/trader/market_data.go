package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/spf13/cobra"
)

var (
	mdStart      string
	mdEnd        string
	mdTotalLimit int
	mdPageLimit  int
	mdFeed       string
	mdAsOf       string
	mdCurrency   string
	mdSort       string

	mdTimeFrame  string
	mdAdjustment string
)

var marketDataCmd = &cobra.Command{
	Use:   "market-data",
	Short: "Get historical market data",
}

var marketDataBarsCmd = &cobra.Command{
	Use:   "bars <symbol>",
	Short: "Get historical bars for a symbol",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getMarketDataClient()
		if err != nil {
			return err
		}

		symbol := strings.ToUpper(args[0])
		req := marketdata.GetBarsRequest{
			TotalLimit: mdTotalLimit,
			PageLimit:  mdPageLimit,
			Feed:       marketdata.Feed(mdFeed),
			AsOf:       mdAsOf,
			Currency:   mdCurrency,
			Sort:       marketdata.Sort(strings.ToLower(mdSort)),
			Adjustment: marketdata.Adjustment(mdAdjustment),
		}

		if mdTimeFrame != "" {
			var tf marketdata.TimeFrame
			// TimeFrame implements UnmarshalJSON for strings like "1Min", "1Day", etc.
			if err := json.Unmarshal([]byte(`"`+mdTimeFrame+`"`), &tf); err != nil {
				return fmt.Errorf("invalid timeframe format: %w", err)
			}
			req.TimeFrame = tf
		}

		if mdStart != "" {
			req.Start, err = time.Parse(time.RFC3339, mdStart)
			if err != nil {
				return fmt.Errorf("invalid start date format (expected RFC3339): %w", err)
			}
		}

		if mdEnd != "" {
			req.End, err = time.Parse(time.RFC3339, mdEnd)
			if err != nil {
				return fmt.Errorf("invalid end date format (expected RFC3339): %w", err)
			}
		}

		bars, err := client.GetBars(symbol, req)
		if err != nil {
			return fmt.Errorf("failed to get bars for %s: %w", symbol, err)
		}

		return printOutput(bars)
	},
}

var marketDataQuotesCmd = &cobra.Command{
	Use:   "quotes <symbol>",
	Short: "Get historical quotes for a symbol",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getMarketDataClient()
		if err != nil {
			return err
		}

		symbol := strings.ToUpper(args[0])
		req := marketdata.GetQuotesRequest{
			TotalLimit: mdTotalLimit,
			PageLimit:  mdPageLimit,
			Feed:       marketdata.Feed(mdFeed),
			AsOf:       mdAsOf,
			Currency:   mdCurrency,
			Sort:       marketdata.Sort(strings.ToLower(mdSort)),
		}

		if mdStart != "" {
			req.Start, err = time.Parse(time.RFC3339, mdStart)
			if err != nil {
				return fmt.Errorf("invalid start date format (expected RFC3339): %w", err)
			}
		}

		if mdEnd != "" {
			req.End, err = time.Parse(time.RFC3339, mdEnd)
			if err != nil {
				return fmt.Errorf("invalid end date format (expected RFC3339): %w", err)
			}
		}

		quotes, err := client.GetQuotes(symbol, req)
		if err != nil {
			return fmt.Errorf("failed to get quotes for %s: %w", symbol, err)
		}

		return printOutput(quotes)
	},
}

var marketDataTradesCmd = &cobra.Command{
	Use:   "trades <symbol>",
	Short: "Get historical trades for a symbol",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getMarketDataClient()
		if err != nil {
			return err
		}

		symbol := strings.ToUpper(args[0])
		req := marketdata.GetTradesRequest{
			TotalLimit: mdTotalLimit,
			PageLimit:  mdPageLimit,
			Feed:       marketdata.Feed(mdFeed),
			AsOf:       mdAsOf,
			Currency:   mdCurrency,
			Sort:       marketdata.Sort(strings.ToLower(mdSort)),
		}

		if mdStart != "" {
			req.Start, err = time.Parse(time.RFC3339, mdStart)
			if err != nil {
				return fmt.Errorf("invalid start date format (expected RFC3339): %w", err)
			}
		}

		if mdEnd != "" {
			req.End, err = time.Parse(time.RFC3339, mdEnd)
			if err != nil {
				return fmt.Errorf("invalid end date format (expected RFC3339): %w", err)
			}
		}

		trades, err := client.GetTrades(symbol, req)
		if err != nil {
			return fmt.Errorf("failed to get trades for %s: %w", symbol, err)
		}

		return printOutput(trades)
	},
}

func init() {
	RootCmd.AddCommand(marketDataCmd)

	// Add common flags to all subcommands
	for _, subCmd := range []*cobra.Command{marketDataBarsCmd, marketDataQuotesCmd, marketDataTradesCmd} {
		subCmd.Flags().StringVar(&mdStart, "start", "", "Inclusive beginning of interval (RFC3339)")
		subCmd.Flags().StringVar(&mdEnd, "end", "", "Inclusive end of interval (RFC3339)")
		subCmd.Flags().IntVar(&mdTotalLimit, "total-limit", 0, "Total number of items to return (0 means all)")
		subCmd.Flags().IntVar(&mdPageLimit, "page-limit", 0, "Pagination size")
		subCmd.Flags().StringVar(&mdFeed, "feed", "", "Source of data: sip, iex, otc")
		subCmd.Flags().StringVar(&mdAsOf, "as-of", "", "Date when the symbols are mapped")
		subCmd.Flags().StringVar(&mdCurrency, "currency", "", "Currency of displayed prices")
		subCmd.Flags().StringVar(&mdSort, "sort", "asc", "Sort direction (asc or desc)")

		marketDataCmd.AddCommand(subCmd)
	}

	// Add specific flags for bars
	marketDataBarsCmd.Flags().StringVar(&mdTimeFrame, "timeframe", "1Day", "Aggregation size (e.g. 1Min, 1Hour, 1Day, 1Week, 1Month)")
	marketDataBarsCmd.Flags().StringVar(&mdAdjustment, "adjustment", "raw", "Adjustment for corporate actions (raw, split, dividend, all)")
}
