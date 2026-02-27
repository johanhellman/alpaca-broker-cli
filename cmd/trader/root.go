package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "alpaca-trader",
	Short: "A CLI tool for the Alpaca Trading API",
	Long: `alpaca-trader is a powerful command-line interface for interacting
with the Alpaca Trading API to manage your retail/paper account, positions, and orders.

Example:
  alpaca-trader account get --env paper
  alpaca-trader positions list

Set APCA_API_KEY_ID and APCA_API_SECRET_KEY in your environment to authenticate.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.alpaca-trader.yaml)")
	rootCmd.PersistentFlags().String("api-key", "", "Alpaca Trading API Key ID")
	rootCmd.PersistentFlags().String("api-secret", "", "Alpaca Trading API Secret Key")
	rootCmd.PersistentFlags().String("env", "paper", "Alpaca environment (paper or live)")
	rootCmd.PersistentFlags().String("output", "table", "Output format (table or json)")

	viper.BindPFlag("api-key", rootCmd.PersistentFlags().Lookup("api-key"))
	viper.BindPFlag("api-secret", rootCmd.PersistentFlags().Lookup("api-secret"))
	viper.BindPFlag("env", rootCmd.PersistentFlags().Lookup("env"))
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
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
	viper.BindEnv("api-key", "APCA_API_KEY_ID")
	viper.BindEnv("api-secret", "APCA_API_SECRET_KEY")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
