package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrdersCmd_TraderFlags(t *testing.T) {
	assert.NotNil(t, ordersCmd)
	assert.Equal(t, "orders", ordersCmd.Use)

	// Check orders create topology
	assert.NotNil(t, ordersCreateCmd)
	createFlags := ordersCreateCmd.Flags()
	assert.NotNil(t, createFlags.Lookup("symbol"))
	assert.NotNil(t, createFlags.Lookup("qty"))
	assert.NotNil(t, createFlags.Lookup("side"))
	assert.NotNil(t, createFlags.Lookup("type"))
	assert.NotNil(t, createFlags.Lookup("time-in-force"))

	// Check required flags
	b := bytes.NewBufferString("")
	RootCmd.SetOut(b)
	RootCmd.SetArgs([]string{"orders", "create"})
	err := RootCmd.Execute()
	
	// Should fail because --symbol, --side, etc. are missing
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "required flag")
}
