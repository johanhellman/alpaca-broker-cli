package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJournalsCmd_BrokerFlags(t *testing.T) {
	assert.NotNil(t, journalsCmd)
	assert.Equal(t, "journals", journalsCmd.Use)

	// Check create topology
	assert.NotNil(t, journalsCreateCmd)
	cFlags := journalsCreateCmd.Flags()
	assert.NotNil(t, cFlags.Lookup("entry-type"))
	assert.NotNil(t, cFlags.Lookup("from-account"))
	assert.NotNil(t, cFlags.Lookup("to-account"))

	b := bytes.NewBufferString("")
	rootCmd.SetOut(b)
	rootCmd.SetArgs([]string{"journals", "create"})
	err := rootCmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "required flag")

	// Check list topology
	assert.NotNil(t, journalsListCmd)
	lFlags := journalsListCmd.Flags()
	assert.NotNil(t, lFlags.Lookup("after"))
	assert.NotNil(t, lFlags.Lookup("status"))

	b.Reset()
	rootCmd.SetArgs([]string{"journals", "--help"})
	err = rootCmd.Execute()
	assert.NoError(t, err)

	out := b.String()
	assert.Contains(t, out, journalsCmd.Short)
}
