package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	client "github.com/johanhellman/alpaca-broker-cli/internal/brokerclient"
	"github.com/johanhellman/alpaca-broker-cli/internal/brokerclient/api"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/spf13/cobra"
)

var (
	// Create
	journalEntryType string
	journalFrom      string
	journalTo        string
	journalAmount    string
	journalQty       string
	journalSymbol    string

	// List
	listJournalAfter       string
	listJournalBefore      string
	listJournalStatus      string
	listJournalEntryType   string
	listJournalToAccount   string
	listJournalFromAccount string
)

var journalsCmd = &cobra.Command{
	Use:   "journals",
	Short: "Manage journals between accounts",
}

var journalsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List journals",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := api.NewClient()
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		params := &client.GetJournalsParams{}

		if listJournalAfter != "" {
			t, err := time.Parse("2006-01-02", listJournalAfter)
			if err != nil {
				return fmt.Errorf("invalid after date format (expected YYYY-MM-DD): %w", err)
			}
			aft := openapi_types.Date{Time: t}
			params.After = &aft
		}
		if listJournalBefore != "" {
			t, err := time.Parse("2006-01-02", listJournalBefore)
			if err != nil {
				return fmt.Errorf("invalid before date format (expected YYYY-MM-DD): %w", err)
			}
			bef := openapi_types.Date{Time: t}
			params.Before = &bef
		}
		if listJournalStatus != "" {
			st := client.GetJournalsParamsStatus(listJournalStatus)
			params.Status = &st
		}
		if listJournalEntryType != "" {
			et := client.GetJournalsParamsEntryType(listJournalEntryType)
			params.EntryType = &et
		}
		if listJournalToAccount != "" {
			id, err := uuid.Parse(listJournalToAccount)
			if err != nil {
				return fmt.Errorf("invalid to-account format: %w", err)
			}
			params.ToAccount = &id
		}
		if listJournalFromAccount != "" {
			id, err := uuid.Parse(listJournalFromAccount)
			if err != nil {
				return fmt.Errorf("invalid from-account format: %w", err)
			}
			params.FromAccount = &id
		}

		resp, err := c.GetJournalsWithResponse(ctx, params)
		if err != nil {
			return fmt.Errorf("failed to list journals: %w", err)
		}

		if resp.JSON200 == nil {
			return fmt.Errorf("unexpected response status: %d", resp.StatusCode())
		}

		return printOutput(resp.JSON200)
	},
}

// (removed GetCmd as there is no GetJournal by ID in the SDK)

var journalsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new journal",
	RunE: func(cmd *cobra.Command, args []string) error {
		fromUUID, err := uuid.Parse(journalFrom)
		if err != nil {
			return fmt.Errorf("invalid from-account ID: %w", err)
		}

		toUUID, err := uuid.Parse(journalTo)
		if err != nil {
			return fmt.Errorf("invalid to-account ID: %w", err)
		}

		req := client.JournalData{
			EntryType:   client.JournalDataEntryType(journalEntryType),
			FromAccount: fromUUID,
			ToAccount:   toUUID,
		}

		if journalAmount != "" {
			req.Amount = &journalAmount
		}
		if journalQty != "" {
			req.Qty = &journalQty
		}
		if journalSymbol != "" {
			req.Symbol = &journalSymbol
		}

		c, err := api.NewClient()
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		resp, err := c.PostJournalsWithResponse(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to create journal: %w", err)
		}

		if resp.JSON200 == nil {
			if resp.JSON400 != nil {
				return fmt.Errorf("bad request: %v", resp.JSON400)
			}
			return fmt.Errorf("unexpected response status: %d", resp.StatusCode())
		}

		return printOutput(resp.JSON200)
	},
}

var journalsCancelCmd = &cobra.Command{
	Use:   "cancel <journal_id>",
	Short: "Cancel a pending journal",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		idStr := args[0]
		parsedUUID, err := uuid.Parse(idStr)
		if err != nil {
			return fmt.Errorf("invalid journal ID format: %w", err)
		}

		c, err := api.NewClient()
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		resp, err := c.DeleteJournalWithResponse(ctx, parsedUUID)
		if err != nil {
			return fmt.Errorf("failed to cancel journal %s: %w", idStr, err)
		}

		if resp.StatusCode() != 200 && resp.StatusCode() != 204 {
			return fmt.Errorf("unexpected response status: %d", resp.StatusCode())
		}

		fmt.Printf("Journal %s cancellation requested.\n", idStr)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(journalsCmd)

	journalsCreateCmd.Flags().StringVar(&journalEntryType, "entry-type", "", "JNLC (cash) or JNLS (shares)")
	_ = journalsCreateCmd.MarkFlagRequired("entry-type") //nolint:errcheck
	journalsCreateCmd.Flags().StringVar(&journalFrom, "from-account", "", "Source account ID")
	_ = journalsCreateCmd.MarkFlagRequired("from-account") //nolint:errcheck
	journalsCreateCmd.Flags().StringVar(&journalTo, "to-account", "", "Destination account ID")
	_ = journalsCreateCmd.MarkFlagRequired("to-account") //nolint:errcheck

	journalsCreateCmd.Flags().StringVar(&journalAmount, "amount", "", "Dollar amount for JNLC")
	journalsCreateCmd.Flags().StringVar(&journalQty, "qty", "", "Number of shares for JNLS")
	journalsCreateCmd.Flags().StringVar(&journalSymbol, "symbol", "", "Symbol for JNLS")
	journalsCmd.AddCommand(journalsCreateCmd)

	journalsListCmd.Flags().StringVar(&listJournalAfter, "after", "", "Filter by settle_date after (YYYY-MM-DD)")
	journalsListCmd.Flags().StringVar(&listJournalBefore, "before", "", "Filter by settle_date before (YYYY-MM-DD)")
	journalsListCmd.Flags().StringVar(&listJournalStatus, "status", "", "Filter by status (e.g. pending, executed, canceled)")
	journalsListCmd.Flags().StringVar(&listJournalEntryType, "entry-type", "", "Filter by entry type (JNLC, JNLS)")
	journalsListCmd.Flags().StringVar(&listJournalToAccount, "to-account", "", "Filter by destination account ID")
	journalsListCmd.Flags().StringVar(&listJournalFromAccount, "from-account", "", "Filter by source account ID")
	journalsCmd.AddCommand(journalsListCmd)

	journalsCmd.AddCommand(journalsCancelCmd)
}
