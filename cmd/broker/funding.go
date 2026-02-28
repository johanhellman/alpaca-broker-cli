package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	client "github.com/johanhellman/alpaca-broker-cli/internal/brokerclient"
	"github.com/johanhellman/alpaca-broker-cli/internal/brokerclient/api"
	"github.com/oapi-codegen/runtime/types"
	"github.com/spf13/cobra"
)

var (
	// List flags
	listTransfersDirection string
	listTransfersLimit     int32
	listTransfersOffset    int32

	// Create flags
	createTransferType           string
	createTransferAmount         string
	createTransferDirection      string
	createTransferRelationshipID string
	createTransferBankID         string
	createTransferAdditionalInfo string
)

var fundingCmd = &cobra.Command{
	Use:   "funding",
	Short: "Manage funding and transfers",
	Long:  `List and create ACH relationships and bank transfers for an account.`,
}

var fundingTransfersCmd = &cobra.Command{
	Use:   "transfers <account_id>",
	Short: "List transfers for an account",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		accountIDStr := args[0]
		parsedUUID, err := uuid.Parse(accountIDStr)
		if err != nil {
			return fmt.Errorf("invalid account ID format: %w", err)
		}

		c, err := api.NewClient()
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		params := &client.GetTransfersParams{}
		if listTransfersDirection != "" {
			dir := client.GetTransfersParamsDirection(listTransfersDirection)
			params.Direction = &dir
		}
		if listTransfersLimit > 0 {
			params.Limit = &listTransfersLimit
		}
		if listTransfersOffset > 0 {
			params.Offset = &listTransfersOffset
		}

		resp, err := c.GetTransfersWithResponse(ctx, parsedUUID, params)
		if err != nil {
			return fmt.Errorf("failed to list transfers: %w", err)
		}

		if resp.JSON200 == nil {
			return fmt.Errorf("unexpected response status: %d", resp.StatusCode())
		}

		return printOutput(resp.JSON200)
	},
}

var fundingTransferCreateCmd = &cobra.Command{
	Use:   "transfer-create <account_id>",
	Short: "Create a transfer for an account",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		accountIDStr := args[0]
		parsedUUID, err := uuid.Parse(accountIDStr)
		if err != nil {
			return fmt.Errorf("invalid account ID format: %w", err)
		}

		var reqBody client.TransferData
		reqBody.TransferType = client.TransferDataTransferType(createTransferType)

		if createTransferType == "ach" {
			if createTransferRelationshipID == "" {
				return fmt.Errorf("relationship-id is required for ACH transfers")
			}
			relUUID, err := uuid.Parse(createTransferRelationshipID)
			if err != nil {
				return fmt.Errorf("invalid relationship ID: %w", err)
			}
			achData := client.UntypedACHTransferData{
				Amount:         createTransferAmount,
				Direction:      client.UntypedACHTransferDataDirection(createTransferDirection),
				RelationshipId: types.UUID(relUUID),
			}
			if err := reqBody.FromUntypedACHTransferData(achData); err != nil {
				return fmt.Errorf("failed to set ACH transfer data: %w", err)
			}
		} else if createTransferType == "wire" {
			if createTransferBankID == "" {
				return fmt.Errorf("bank-id is required for wire transfers")
			}
			bankUUID, err := uuid.Parse(createTransferBankID)
			if err != nil {
				return fmt.Errorf("invalid bank ID: %w", err)
			}
			wireData := client.UntypedWireTransferData{
				Amount:    createTransferAmount,
				Direction: client.UntypedWireTransferDataDirection(createTransferDirection),
				BankId:    types.UUID(bankUUID),
			}
			if createTransferAdditionalInfo != "" {
				wireData.AdditionalInformation = &createTransferAdditionalInfo
			}
			if err := reqBody.FromUntypedWireTransferData(wireData); err != nil {
				return fmt.Errorf("failed to set wire transfer data: %w", err)
			}
		} else {
			return fmt.Errorf("unsupported transfer type: %s", createTransferType)
		}

		c, err := api.NewClient()
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		resp, err := c.PostTransfersWithResponse(ctx, parsedUUID, reqBody)
		if err != nil {
			return fmt.Errorf("failed to create transfer: %w", err)
		}

		if resp.JSON200 == nil {
			return fmt.Errorf("unexpected response status: %d, body: %s", resp.StatusCode(), string(resp.Body))
		}

		return printOutput(resp.JSON200)
	},
}

func init() {
	rootCmd.AddCommand(fundingCmd)

	fundingTransfersCmd.Flags().StringVar(&listTransfersDirection, "direction", "", "INCOMING or OUTGOING")
	fundingTransfersCmd.Flags().Int32Var(&listTransfersLimit, "limit", 50, "Maximum number of transfers to return")
	fundingTransfersCmd.Flags().Int32Var(&listTransfersOffset, "offset", 0, "Pagination offset")
	fundingCmd.AddCommand(fundingTransfersCmd)

	fundingTransferCreateCmd.Flags().StringVar(&createTransferType, "transfer-type", "", "ach or wire")
	_ = fundingTransferCreateCmd.MarkFlagRequired("transfer-type") //nolint:errcheck
	fundingTransferCreateCmd.Flags().StringVar(&createTransferAmount, "amount", "", "Amount as a string, e.g. 500.00")
	_ = fundingTransferCreateCmd.MarkFlagRequired("amount") //nolint:errcheck
	fundingTransferCreateCmd.Flags().StringVar(&createTransferDirection, "direction", "", "INCOMING or OUTGOING")
	_ = fundingTransferCreateCmd.MarkFlagRequired("direction") //nolint:errcheck

	fundingTransferCreateCmd.Flags().StringVar(&createTransferRelationshipID, "relationship-id", "", "Required for ACH transfers")
	fundingTransferCreateCmd.Flags().StringVar(&createTransferBankID, "bank-id", "", "Required for wire transfers")
	fundingTransferCreateCmd.Flags().StringVar(&createTransferAdditionalInfo, "additional-info", "", "Optional for wire transfers")

	fundingCmd.AddCommand(fundingTransferCreateCmd)
}
