package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "alpaca-trader",
	Short: "A CLI tool for the Alpaca Trading API",
	Long: `alpaca-trader is a powerful command-line interface for interacting
with the Alpaca Trading API to manage your retail/paper account, positions, and orders.

Example:
  alpaca-trader account get --env paper
  alpaca-trader positions list

Set APCA_API_KEY_ID and APCA_API_SECRET_KEY in your environment to authenticate.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.alpaca-trader.yaml)")

	RootCmd.PersistentFlags().String("api-key", "", "Alpaca API Key ID")
	_ = viper.BindPFlag("api-key", RootCmd.PersistentFlags().Lookup("api-key")) //nolint:errcheck

	RootCmd.PersistentFlags().String("api-secret", "", "Alpaca API Secret Key")
	_ = viper.BindPFlag("api-secret", RootCmd.PersistentFlags().Lookup("api-secret")) //nolint:errcheck

	RootCmd.PersistentFlags().String("env", "paper", "Alpaca environment (paper or live)")
	_ = viper.BindPFlag("env", RootCmd.PersistentFlags().Lookup("env")) //nolint:errcheck

	RootCmd.PersistentFlags().String("output", "table", "Output format (table or json)")
	_ = viper.BindPFlag("output", RootCmd.PersistentFlags().Lookup("output")) //nolint:errcheck
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".alpaca-trader")
	}

	viper.SetEnvPrefix("APCA")
	// Also support the standard APCA_ environment variables directly
	viper.BindEnv("api-key", "APCA_API_KEY_ID")        //nolint:errcheck
	viper.BindEnv("api-secret", "APCA_API_SECRET_KEY") //nolint:errcheck
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
