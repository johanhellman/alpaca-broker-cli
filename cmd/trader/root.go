package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
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

	RootCmd.PersistentFlags().String("query", "", "Filter output using jq-like syntax (forces json output if used)")
	_ = viper.BindPFlag("query", RootCmd.PersistentFlags().Lookup("query")) //nolint:errcheck

	RootCmd.PersistentFlags().Bool("all", false, "Automatically fetch all pages for list endpoints")
	_ = viper.BindPFlag("all", RootCmd.PersistentFlags().Lookup("all")) //nolint:errcheck
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
	query := viper.GetString("query")

	// If query is present, we force JSON marshaling and utilize gjson
	if query != "" {
		out, err := json.Marshal(data)
		if err != nil {
			return err
		}
		result := gjson.GetBytes(out, query)
		
		// If the query pulls a single string (like "0.id") don't wrap it in JSON quotes for bash parsing
		if result.Type == gjson.String {
			fmt.Println(result.String())
		} else {
			fmt.Println(result.Raw)
		}
		return nil
	}

	if format == "table" {
		val := reflect.Indirect(reflect.ValueOf(data))
		if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
			slice := reflect.MakeSlice(reflect.SliceOf(val.Type()), 0, 1)
			val = reflect.Append(slice, val)
		}

		if val.Len() == 0 {
			fmt.Println("No data found.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		firstElement := reflect.Indirect(val.Index(0))
		
		if firstElement.Kind() == reflect.Struct {
			for i := 0; i < firstElement.NumField(); i++ {
				fmt.Fprintf(w, "%v\t", firstElement.Type().Field(i).Name)
			}
			fmt.Fprintln(w)
		}

		for i := 0; i < val.Len(); i++ {
			element := reflect.Indirect(val.Index(i))
			if element.Kind() == reflect.Struct {
				for j := 0; j < element.NumField(); j++ {
					field := element.Field(j)
					fmt.Fprintf(w, "%v\t", field.Interface())
				}
				fmt.Fprintln(w)
			} else {
				fmt.Fprintln(w, element.Interface())
			}
		}
		return w.Flush()
	}

	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}
