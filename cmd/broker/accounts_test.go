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

func TestAccountsCreateCmd_BrokerFlags(t *testing.T) {
	assert.NotNil(t, accountsCreateCmd)
	
	flags := accountsCreateCmd.Flags()
	assert.NotNil(t, flags.Lookup("contact-email"))
	assert.NotNil(t, flags.Lookup("contact-phone"))
	assert.NotNil(t, flags.Lookup("id-given-name"))
	assert.NotNil(t, flags.Lookup("id-family-name"))
	assert.NotNil(t, flags.Lookup("id-dob"))
	
	// Test requirement failing constraint
	b := bytes.NewBufferString("")
	rootCmd.SetOut(b)
	rootCmd.SetArgs([]string{"accounts", "create"})
	err := rootCmd.Execute()
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "required flag")
}
