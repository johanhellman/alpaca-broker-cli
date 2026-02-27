package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountsCmd_Broker(t *testing.T) {
	// Verify command structure
	assert.NotNil(t, accountsCmd)
	assert.Equal(t, "accounts", accountsCmd.Use)

	// Test command execution (help output)
	b := bytes.NewBufferString("")
	rootCmd.SetOut(b)
	rootCmd.SetArgs([]string{"accounts", "--help"})
	err := rootCmd.Execute()
	assert.NoError(t, err)

	out := b.String()
	assert.Contains(t, out, accountsCmd.Long)
}
