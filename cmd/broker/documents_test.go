package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocumentsCmd_BrokerFlags(t *testing.T) {
	assert.NotNil(t, documentsCmd)
	assert.Equal(t, "documents", documentsCmd.Use)

	// Check upload topology
	assert.NotNil(t, documentsUploadCmd)
	uFlags := documentsUploadCmd.Flags()
	assert.NotNil(t, uFlags.Lookup("file"))
	assert.NotNil(t, uFlags.Lookup("document-type"))
	assert.NotNil(t, uFlags.Lookup("mime-type"))

	b := bytes.NewBufferString("")
	rootCmd.SetOut(b)
	rootCmd.SetArgs([]string{"documents", "upload", "123e4567-e89b-12d3-a456-426614174000"})
	err := rootCmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "required flag")

	// Check list topology
	assert.NotNil(t, documentsListCmd)
	lFlags := documentsListCmd.Flags()
	assert.NotNil(t, lFlags.Lookup("start-date"))
	assert.NotNil(t, lFlags.Lookup("end-date"))

	b.Reset()
	rootCmd.SetArgs([]string{"documents", "--help"})
	err = rootCmd.Execute()
	assert.NoError(t, err)

	out := b.String()
	assert.Contains(t, out, documentsCmd.Short)
}
