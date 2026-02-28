package cmd

import (
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/civil"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/spf13/cobra"
)

var (
	corporateActionsSymbols    []string
	corporateActionsTypes      []string
	corporateActionsStart      string
	corporateActionsEnd        string
	corporateActionsTotalLimit int
	corporateActionsPageLimit  int
	corporateActionsSort       string
)

var corporateActionsCmd = &cobra.Command{
	Use:   "corporate-actions",
	Short: "Get corporate actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getMarketDataClient()
		if err != nil {
			return err
		}

		req := marketdata.GetCorporateActionsRequest{
			Symbols:    corporateActionsSymbols,
			Types:      corporateActionsTypes,
			TotalLimit: corporateActionsTotalLimit,
			PageLimit:  corporateActionsPageLimit,
			Sort:       marketdata.Sort(strings.ToLower(corporateActionsSort)),
		}

		if corporateActionsStart != "" {
			t, err := time.Parse("2006-01-02", corporateActionsStart)
			if err != nil {
				return fmt.Errorf("invalid start date format (expected YYYY-MM-DD): %w", err)
			}
			req.Start = civil.DateOf(t)
		}

		if corporateActionsEnd != "" {
			t, err := time.Parse("2006-01-02", corporateActionsEnd)
			if err != nil {
				return fmt.Errorf("invalid end date format (expected YYYY-MM-DD): %w", err)
			}
			req.End = civil.DateOf(t)
		}

		actions, err := client.GetCorporateActions(req)
		if err != nil {
			return fmt.Errorf("failed to get corporate actions: %w", err)
		}

		return printOutput(actions)
	},
}

func init() {
	RootCmd.AddCommand(corporateActionsCmd)

	corporateActionsCmd.Flags().StringSliceVar(&corporateActionsSymbols, "symbols", nil, "Comma-separated list of company symbols")
	corporateActionsCmd.Flags().StringSliceVar(&corporateActionsTypes, "types", nil, "Comma-separated list of corporate actions types (e.g. forward_split, cash_dividend)")
	corporateActionsCmd.Flags().StringVar(&corporateActionsStart, "start", "", "Inclusive beginning of the interval (YYYY-MM-DD)")
	corporateActionsCmd.Flags().StringVar(&corporateActionsEnd, "end", "", "Inclusive end of the interval (YYYY-MM-DD)")
	corporateActionsCmd.Flags().IntVar(&corporateActionsTotalLimit, "total-limit", 0, "Limit of the total number of actions returned (0 means all)")
	corporateActionsCmd.Flags().IntVar(&corporateActionsPageLimit, "page-limit", 0, "Pagination size")
	corporateActionsCmd.Flags().StringVar(&corporateActionsSort, "sort", "asc", "Sort direction (asc or desc)")
}
