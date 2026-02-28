package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
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
	_ = viper.BindPFlag("api-key", rootCmd.PersistentFlags().Lookup("api-key")) //nolint:errcheck

	rootCmd.PersistentFlags().String("api-secret", "", "Alpaca Broker API Secret")
	_ = viper.BindPFlag("api-secret", rootCmd.PersistentFlags().Lookup("api-secret")) //nolint:errcheck

	rootCmd.PersistentFlags().String("env", "sandbox", "Alpaca environment (sandbox or production)")
	_ = viper.BindPFlag("env", rootCmd.PersistentFlags().Lookup("env")) //nolint:errcheck

	rootCmd.PersistentFlags().String("output", "table", "Output format (table, json, or csv)")
	_ = viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output")) //nolint:errcheck

	rootCmd.PersistentFlags().String("query", "", "Filter output using jq-like syntax (forces json output if used)")
	_ = viper.BindPFlag("query", rootCmd.PersistentFlags().Lookup("query")) //nolint:errcheck

	rootCmd.PersistentFlags().Bool("all", false, "Automatically fetch all pages for list endpoints")
	_ = viper.BindPFlag("all", rootCmd.PersistentFlags().Lookup("all")) //nolint:errcheck
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

// printOutput formats data as JSON or Table based on the `--output` flag.
func printOutput(data interface{}) error {
	outputFormat := viper.GetString("output")
	query := viper.GetString("query")

	if query != "" {
		return printQuery(data, query)
	}

	if outputFormat == "table" {
		return printTable(data)
	}

	if outputFormat == "csv" {
		return printCSV(data)
	}

	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

func printQuery(data interface{}, query string) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}
	result := gjson.GetBytes(out, query)

	if result.Type == gjson.String {
		fmt.Println(result.String())
	} else {
		fmt.Println(result.Raw)
	}
	return nil
}

func printTable(data interface{}) error {
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
			_, _ = fmt.Fprintf(w, "%-20v\t", firstElement.Type().Field(i).Name) //nolint:errcheck
		}
		_, _ = fmt.Fprintln(w) //nolint:errcheck
	}

	for i := 0; i < val.Len(); i++ {
		element := reflect.Indirect(val.Index(i))
		if element.Kind() == reflect.Struct {
			for j := 0; j < element.NumField(); j++ {
				field := element.Field(j)
				_, _ = fmt.Fprintf(w, "%-20v\t", field.Interface()) //nolint:errcheck
			}
			_, _ = fmt.Fprintln(w) //nolint:errcheck
		} else {
			_, _ = fmt.Fprintln(w, element.Interface()) //nolint:errcheck
		}
	}
	return w.Flush()
}

func printCSV(data interface{}) error {
	val := reflect.Indirect(reflect.ValueOf(data))
	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		slice := reflect.MakeSlice(reflect.SliceOf(val.Type()), 0, 1)
		val = reflect.Append(slice, val)
	}

	cw := csv.NewWriter(os.Stdout)
	if val.Len() == 0 {
		return nil
	}

	firstElement := reflect.Indirect(val.Index(0))

	if firstElement.Kind() == reflect.Struct {
		var headers []string
		for i := 0; i < firstElement.NumField(); i++ {
			headers = append(headers, firstElement.Type().Field(i).Name)
		}
		if err := cw.Write(headers); err != nil {
			return err
		}
	}

	for i := 0; i < val.Len(); i++ {
		element := reflect.Indirect(val.Index(i))
		var record []string
		if element.Kind() == reflect.Struct {
			for j := 0; j < element.NumField(); j++ {
				field := element.Field(j)
				record = append(record, fmt.Sprintf("%v", field.Interface()))
			}
		} else {
			record = append(record, fmt.Sprintf("%v", element.Interface()))
		}
		if err := cw.Write(record); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}
