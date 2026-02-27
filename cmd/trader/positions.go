package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var positionsCmd = &cobra.Command{
	Use:   "positions",
	Short: "Manage trading positions",
}

var positionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List open positions",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := getClient()

		positions, err := client.GetPositions()
		if err != nil {
			return fmt.Errorf("failed to list positions: %w", err)
		}

		out, err := json.MarshalIndent(positions, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(out))
		return nil
	},
}

func init() {
	RootCmd.AddCommand(positionsCmd)
	positionsCmd.AddCommand(positionsListCmd)
}
