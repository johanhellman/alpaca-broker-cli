package cmd

import (
	"fmt"
	"strings"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/spf13/cobra"
)

var (
	listAssetsStatus     string
	listAssetsAssetClass string
	listAssetsExchange   string
)

var assetsCmd = &cobra.Command{
	Use:   "assets",
	Short: "Manage tradable assets",
}

var assetsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available assets",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		req := alpaca.GetAssetsRequest{
			Status:     listAssetsStatus,
			AssetClass: listAssetsAssetClass,
			Exchange:   listAssetsExchange,
		}

		assets, err := client.GetAssets(req)
		if err != nil {
			return fmt.Errorf("failed to list assets: %w", err)
		}

		return printOutput(assets)
	},
}

var assetsGetCmd = &cobra.Command{
	Use:   "get <symbol_or_id>",
	Short: "Get an asset by symbol or ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		symbol := strings.ToUpper(args[0])
		asset, err := client.GetAsset(symbol)
		if err != nil {
			return fmt.Errorf("failed to get asset for %s: %w", symbol, err)
		}

		return printOutput(asset)
	},
}

func init() {
	RootCmd.AddCommand(assetsCmd)

	assetsListCmd.Flags().StringVar(&listAssetsStatus, "status", "active", "active or inactive")
	assetsListCmd.Flags().StringVar(&listAssetsAssetClass, "asset-class", "us_equity", "us_equity or crypto")
	assetsListCmd.Flags().StringVar(&listAssetsExchange, "exchange", "", "AMEX, ARCA, BATS, NYSE, NASDAQ, NYSEARCA, OTC")
	assetsCmd.AddCommand(assetsListCmd)

	assetsCmd.AddCommand(assetsGetCmd)
}
