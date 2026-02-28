package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventsCmd_BrokerFlags(t *testing.T) {
	assert.NotNil(t, eventsCmd)
	assert.Equal(t, "events", eventsCmd.Use)

	// Check accounts topology
	assert.NotNil(t, eventsAccountsCmd)
	aFlags := eventsAccountsCmd.Flags()
	assert.NotNil(t, aFlags.Lookup("since"))
	assert.NotNil(t, aFlags.Lookup("until"))
	assert.NotNil(t, aFlags.Lookup("since-id"))
	assert.NotNil(t, aFlags.Lookup("until-id"))

	// Check journals topology
	assert.NotNil(t, eventsJournalsCmd)
	jFlags := eventsJournalsCmd.Flags()
	assert.NotNil(t, jFlags.Lookup("since-id"))

	// Check trades topology
	assert.NotNil(t, eventsTradesCmd)
	tFlags := eventsTradesCmd.Flags()
	assert.NotNil(t, tFlags.Lookup("until"))

	b := bytes.NewBufferString("")
	rootCmd.SetOut(b)
	rootCmd.SetArgs([]string{"events", "--help"})
	err := rootCmd.Execute()
	assert.NoError(t, err)

	out := b.String()
	assert.Contains(t, out, eventsCmd.Short)
}
