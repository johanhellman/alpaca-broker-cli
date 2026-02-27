package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getClient instantiates the Alpaca Trading API client.
func getClient() *alpaca.Client {
	apiKey := viper.GetString("api-key")
	apiSecret := viper.GetString("api-secret")
	env := viper.GetString("env")

	if apiKey == "" || apiSecret == "" {
		log.Fatal("Missing APCA_API_KEY_ID or APCA_API_SECRET_KEY")
	}

	baseURL := "https://paper-api.alpaca.markets"
	if env == "production" || env == "live" {
		baseURL = "https://api.alpaca.markets"
	}

	return alpaca.NewClient(alpaca.ClientOpts{
		APIKey:    apiKey,
		APISecret: apiSecret,
		BaseURL:   baseURL,
	})
}

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Manage your Alpaca trading account",
}

var accountGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get account details",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getClient()
		acct, err := client.GetAccount()
		if err != nil {
			return fmt.Errorf("failed to get account: %w", err)
		}

		out, err := json.MarshalIndent(acct, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(out))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(accountCmd)
	accountCmd.AddCommand(accountGetCmd)
}
