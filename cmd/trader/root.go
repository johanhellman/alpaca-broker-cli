package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

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

func printOutput(data interface{}) error {
	format := viper.GetString("output")

	switch format {
	case "table":
		// Simple table fallback. For lists, we should iterate. For single objects, print keys/values.
		// A full robust table printer should be added in a future iteration, for now we do a simple KV or array print.
		v := reflect.ValueOf(data)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
			fmt.Printf("Total Items: %d\n", v.Len())
			for i := 0; i < v.Len(); i++ {
				fmt.Printf("--- Item %d ---\n", i)
				item := v.Index(i)
				if item.Kind() == reflect.Ptr {
					item = item.Elem()
				}
				if item.Kind() == reflect.Struct {
					for j := 0; j < item.NumField(); j++ {
						fmt.Printf("%s: %v\n", item.Type().Field(j).Name, item.Field(j).Interface())
					}
				} else {
					fmt.Printf("%v\n", item.Interface())
				}
			}
		} else if v.Kind() == reflect.Struct {
			for i := 0; i < v.NumField(); i++ {
				fmt.Printf("%s: %v\n", v.Type().Field(i).Name, v.Field(i).Interface())
			}
		} else {
			fmt.Printf("%v\n", data)
		}

	case "json", "JSON":
		fallthrough
	default:
		out, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(out))
	}

	return nil
}
