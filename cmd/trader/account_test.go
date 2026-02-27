package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountCmd_Trader(t *testing.T) {
	// Verify command structure
	assert.NotNil(t, accountCmd)
	assert.Equal(t, "account", accountCmd.Use)

	// Test command execution (help output)
	b := bytes.NewBufferString("")
	RootCmd.SetOut(b)
	RootCmd.SetArgs([]string{"account", "--help"})
	err := RootCmd.Execute()
	assert.NoError(t, err)

	out := b.String()
	assert.Contains(t, out, accountCmd.Long)
}
