package cmd

import (
	"fmt"
	"strings"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/spf13/cobra"
)

var (
	createWatchlistName    string
	createWatchlistSymbols []string

	updateWatchlistName    string
	updateWatchlistSymbols []string

	addSymbolWatchlist    string
	removeSymbolWatchlist string
)

var watchlistsCmd = &cobra.Command{
	Use:   "watchlists",
	Short: "Manage watchlists",
}

var watchlistsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all watchlists",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		watchlists, err := client.GetWatchlists()
		if err != nil {
			return fmt.Errorf("failed to list watchlists: %w", err)
		}

		return printOutput(watchlists)
	},
}

var watchlistsGetCmd = &cobra.Command{
	Use:   "get <watchlist_id>",
	Short: "Get a watchlist by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		id := args[0]
		watchlist, err := client.GetWatchlist(id)
		if err != nil {
			return fmt.Errorf("failed to get watchlist %s: %w", id, err)
		}

		return printOutput(watchlist)
	},
}

var watchlistsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new watchlist",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		req := alpaca.CreateWatchlistRequest{
			Name:    createWatchlistName,
			Symbols: createWatchlistSymbols,
		}

		watchlist, err := client.CreateWatchlist(req)
		if err != nil {
			return fmt.Errorf("failed to create watchlist: %w", err)
		}

		return printOutput(watchlist)
	},
}

var watchlistsUpdateCmd = &cobra.Command{
	Use:   "update <watchlist_id>",
	Short: "Update a watchlist",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		id := args[0]
		req := alpaca.UpdateWatchlistRequest{
			Name:    updateWatchlistName,
			Symbols: updateWatchlistSymbols,
		}

		watchlist, err := client.UpdateWatchlist(id, req)
		if err != nil {
			return fmt.Errorf("failed to update watchlist %s: %w", id, err)
		}

		return printOutput(watchlist)
	},
}

var watchlistsDeleteCmd = &cobra.Command{
	Use:   "delete <watchlist_id>",
	Short: "Delete a watchlist",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		id := args[0]
		err = client.DeleteWatchlist(id)
		if err != nil {
			return fmt.Errorf("failed to delete watchlist %s: %w", id, err)
		}

		fmt.Printf("Watchlist %s successfully deleted.\n", id)
		return nil
	},
}

var watchlistsAddAssetCmd = &cobra.Command{
	Use:   "add-asset <watchlist_id>",
	Short: "Add an asset to a watchlist",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		id := args[0]
		req := alpaca.AddSymbolToWatchlistRequest{
			Symbol: strings.ToUpper(addSymbolWatchlist),
		}

		watchlist, err := client.AddSymbolToWatchlist(id, req)
		if err != nil {
			return fmt.Errorf("failed to add symbol %s to watchlist %s: %w", req.Symbol, id, err)
		}

		return printOutput(watchlist)
	},
}

var watchlistsRemoveAssetCmd = &cobra.Command{
	Use:   "remove-asset <watchlist_id>",
	Short: "Remove an asset from a watchlist",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		id := args[0]
		sym := strings.ToUpper(removeSymbolWatchlist)
		req := alpaca.RemoveSymbolFromWatchlistRequest{
			Symbol: sym,
		}
		err = client.RemoveSymbolFromWatchlist(id, req)
		if err != nil {
			return fmt.Errorf("failed to remove symbol %s from watchlist %s: %w", sym, id, err)
		}

		fmt.Printf("Symbol %s successfully removed from watchlist %s.\n", sym, id)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(watchlistsCmd)

	watchlistsCmd.AddCommand(watchlistsListCmd)
	watchlistsCmd.AddCommand(watchlistsGetCmd)

	watchlistsCreateCmd.Flags().StringVar(&createWatchlistName, "name", "", "Name of the watchlist")
	_ = watchlistsCreateCmd.MarkFlagRequired("name") //nolint:errcheck
	watchlistsCreateCmd.Flags().StringSliceVar(&createWatchlistSymbols, "symbols", nil, "Comma-separated list of symbols to add")
	watchlistsCmd.AddCommand(watchlistsCreateCmd)

	watchlistsUpdateCmd.Flags().StringVar(&updateWatchlistName, "name", "", "New name for the watchlist")
	watchlistsUpdateCmd.Flags().StringSliceVar(&updateWatchlistSymbols, "symbols", nil, "New list of symbols for the watchlist")
	watchlistsCmd.AddCommand(watchlistsUpdateCmd)

	watchlistsCmd.AddCommand(watchlistsDeleteCmd)

	watchlistsAddAssetCmd.Flags().StringVar(&addSymbolWatchlist, "symbol", "", "Symbol to add to the watchlist")
	_ = watchlistsAddAssetCmd.MarkFlagRequired("symbol") //nolint:errcheck
	watchlistsCmd.AddCommand(watchlistsAddAssetCmd)

	watchlistsRemoveAssetCmd.Flags().StringVar(&removeSymbolWatchlist, "symbol", "", "Symbol to remove from the watchlist")
	_ = watchlistsRemoveAssetCmd.MarkFlagRequired("symbol") //nolint:errcheck
	watchlistsCmd.AddCommand(watchlistsRemoveAssetCmd)
}
