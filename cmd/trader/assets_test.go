package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssetsCmd_TraderFlags(t *testing.T) {
	assert.NotNil(t, assetsCmd)
	assert.Equal(t, "assets", assetsCmd.Use)

	// Check assets list flags
	assert.NotNil(t, assetsListCmd)
	listFlags := assetsListCmd.Flags()
	assert.NotNil(t, listFlags.Lookup("status"))
	assert.NotNil(t, listFlags.Lookup("asset-class"))
	assert.NotNil(t, listFlags.Lookup("exchange"))

	// Test command execution (help output)
	b := bytes.NewBufferString("")
	RootCmd.SetOut(b)
	RootCmd.SetArgs([]string{"assets", "--help"})
	err := RootCmd.Execute()
	assert.NoError(t, err)

	out := b.String()
	assert.Contains(t, out, assetsCmd.Short)
}
