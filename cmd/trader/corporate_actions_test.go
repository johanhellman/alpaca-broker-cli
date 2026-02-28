package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCorporateActionsCmd_TraderFlags(t *testing.T) {
	assert.NotNil(t, corporateActionsCmd)
	assert.Equal(t, "corporate-actions", corporateActionsCmd.Use)

	// Check topology
	flags := corporateActionsCmd.Flags()
	assert.NotNil(t, flags.Lookup("symbols"))
	assert.NotNil(t, flags.Lookup("types"))
	assert.NotNil(t, flags.Lookup("start"))
	assert.NotNil(t, flags.Lookup("end"))
	assert.NotNil(t, flags.Lookup("total-limit"))
	assert.NotNil(t, flags.Lookup("page-limit"))
	assert.NotNil(t, flags.Lookup("sort"))

	// Test command execution (help output)
	b := bytes.NewBufferString("")
	RootCmd.SetOut(b)
	RootCmd.SetArgs([]string{"corporate-actions", "--help"})
	err := RootCmd.Execute()
	assert.NoError(t, err)

	out := b.String()
	assert.Contains(t, out, corporateActionsCmd.Short)
}
