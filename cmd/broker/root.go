package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "alpaca-cli",
	Short: "A CLI tool for the Alpaca Broker API",
	Long: `alpaca-cli is a powerful command-line interface for interacting
with the Alpaca Broker API to manage accounts, funding, and trading.

Example:
  alpaca-cli accounts list --env sandbox
  alpaca-cli trading orders list`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// RootCmd returns the root cobra command for documentation generation
func RootCmd() *cobra.Command {
	return rootCmd
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.alpaca-cli.yaml)")
	
	rootCmd.PersistentFlags().String("api-key", "", "Alpaca Broker API Key")
	viper.BindPFlag("api-key", rootCmd.PersistentFlags().Lookup("api-key"))

	rootCmd.PersistentFlags().String("api-secret", "", "Alpaca Broker API Secret")
	viper.BindPFlag("api-secret", rootCmd.PersistentFlags().Lookup("api-secret"))

	rootCmd.PersistentFlags().String("env", "sandbox", "Alpaca environment (sandbox or production)")
	viper.BindPFlag("env", rootCmd.PersistentFlags().Lookup("env"))

	rootCmd.PersistentFlags().String("output", "table", "Output format (table or json)")
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".alpaca-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".alpaca-cli")
	}

	// Read environment variables starting with ALPACA_BROKER_
	viper.SetEnvPrefix("ALPACA_BROKER")
	// Replace hyphens with underscores in env var keys, e.g., API_KEY instead of API-KEY
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
