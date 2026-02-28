package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	client "github.com/johanhellman/alpaca-broker-cli/pkg/brokerclient"
	"github.com/johanhellman/alpaca-broker-cli/pkg/brokerclient/api"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/spf13/cobra"
)

var (
	docUploadFile     string
	docUploadType     string
	docUploadMimeType string
	docUploadSubType  string

	docListStartDate string
	docListEndDate   string
)

var documentsCmd = &cobra.Command{
	Use:   "documents",
	Short: "Manage account documents",
}

var documentsUploadCmd = &cobra.Command{
	Use:   "upload <account_id>",
	Short: "Upload a document for an account",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		accountIDStr := args[0]
		parsedUUID, err := uuid.Parse(accountIDStr)
		if err != nil {
			return fmt.Errorf("invalid account ID format: %w", err)
		}

		fileData, err := os.ReadFile(docUploadFile) //nolint:gosec
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", docUploadFile, err)
		}

		base64Content := base64.StdEncoding.EncodeToString(fileData)

		req := client.DocumentUploadRequest{
			Content:      base64Content,
			DocumentType: client.DocumentType(docUploadType),
			MimeType:     docUploadMimeType,
		}

		if docUploadSubType != "" {
			req.DocumentSubType = &docUploadSubType
		}

		c, err := api.NewClient()
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Some generated clients use string or UUID for custom path params
		// Depending on `AccountID` type, this might need casting. Assuming string for safety
		resp, err := c.PostAccountsAccountIdDocumentsUploadWithResponse(ctx, client.AccountID(parsedUUID), req)
		if err != nil {
			return fmt.Errorf("failed to upload document: %w", err)
		}

		if resp.StatusCode() >= 300 {
			return fmt.Errorf("unexpected response status: %d", resp.StatusCode())
		}

		fmt.Println("Document successfully uploaded.")
		return nil
	},
}

var documentsListCmd = &cobra.Command{
	Use:   "list <account_id>",
	Short: "List documents for an account",
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

		params := &client.GetAccountsAccountIdDocumentsParams{}

		if docListStartDate != "" {
			t, err := time.Parse("2006-01-02", docListStartDate)
			if err != nil {
				return fmt.Errorf("invalid start date format: %w", err)
			}
			sd := openapi_types.Date{Time: t}
			params.StartDate = &sd
		}

		if docListEndDate != "" {
			t, err := time.Parse("2006-01-02", docListEndDate)
			if err != nil {
				return fmt.Errorf("invalid end date format: %w", err)
			}
			ed := openapi_types.Date{Time: t}
			params.EndDate = &ed
		}

		resp, err := c.GetAccountsAccountIdDocumentsWithResponse(ctx, client.AccountID(parsedUUID), params)
		if err != nil {
			return fmt.Errorf("failed to list documents: %w", err)
		}

		if resp.JSON200 == nil {
			return fmt.Errorf("unexpected response status: %d", resp.StatusCode())
		}

		return printOutput(resp.JSON200)
	},
}

var documentsDownloadCmd = &cobra.Command{
	Use:   "download <account_id> <document_id>",
	Short: "Download a document by ID",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		accountIDStr := args[0]
		docIDStr := args[1]

		accountUUID, err := uuid.Parse(accountIDStr)
		if err != nil {
			return fmt.Errorf("invalid account ID format: %w", err)
		}

		docUUID, err := uuid.Parse(docIDStr)
		if err != nil {
			return fmt.Errorf("invalid document ID format: %w", err)
		}

		c, err := api.NewClient()
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		resp, err := c.GetAccountsAccountIdDocumentsDocumentIdDownloadWithResponse(ctx, client.AccountID(accountUUID), docUUID)
		if err != nil {
			return fmt.Errorf("failed to download document: %w", err)
		}

		if resp.StatusCode() >= 300 {
			return fmt.Errorf("unexpected response status: %d", resp.StatusCode())
		}

		_, err = os.Stdout.Write(resp.Body)
		return err
	},
}

func init() {
	rootCmd.AddCommand(documentsCmd)

	documentsUploadCmd.Flags().StringVar(&docUploadFile, "file", "", "Path to the document file")
	_ = documentsUploadCmd.MarkFlagRequired("file") //nolint:errcheck
	documentsUploadCmd.Flags().StringVar(&docUploadType, "document-type", "", "Type (e.g. identity_verification, tax_id_verification)")
	_ = documentsUploadCmd.MarkFlagRequired("document-type") //nolint:errcheck
	documentsUploadCmd.Flags().StringVar(&docUploadMimeType, "mime-type", "", "MIME type (e.g. image/jpeg, application/pdf)")
	_ = documentsUploadCmd.MarkFlagRequired("mime-type") //nolint:errcheck
	documentsUploadCmd.Flags().StringVar(&docUploadSubType, "document-sub-type", "", "Sub-type (optional)")
	documentsCmd.AddCommand(documentsUploadCmd)

	documentsListCmd.Flags().StringVar(&docListStartDate, "start-date", "", "Inclusive start date (YYYY-MM-DD)")
	documentsListCmd.Flags().StringVar(&docListEndDate, "end-date", "", "Inclusive end date (YYYY-MM-DD)")
	documentsCmd.AddCommand(documentsListCmd)

	documentsCmd.AddCommand(documentsDownloadCmd)
}
