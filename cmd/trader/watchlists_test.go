package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWatchlistsCmd_TraderFlags(t *testing.T) {
	assert.NotNil(t, watchlistsCmd)
	assert.Equal(t, "watchlists", watchlistsCmd.Use)

	// Check create topology
	assert.NotNil(t, watchlistsCreateCmd)
	createFlags := watchlistsCreateCmd.Flags()
	assert.NotNil(t, createFlags.Lookup("name"))
	assert.NotNil(t, createFlags.Lookup("symbols"))

	// Create required flag failure
	b := bytes.NewBufferString("")
	RootCmd.SetOut(b)
	RootCmd.SetArgs([]string{"watchlists", "create"})
	err := RootCmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "required flag")

	// Check add-asset topology
	assert.NotNil(t, watchlistsAddAssetCmd)
	addFlags := watchlistsAddAssetCmd.Flags()
	assert.NotNil(t, addFlags.Lookup("symbol"))

	b.Reset()
	RootCmd.SetArgs([]string{"watchlists", "add-asset", "some-id"})
	err = RootCmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "required flag")

	// Check remove-asset topology
	assert.NotNil(t, watchlistsRemoveAssetCmd)
	removeFlags := watchlistsRemoveAssetCmd.Flags()
	assert.NotNil(t, removeFlags.Lookup("symbol"))

	b.Reset()
	RootCmd.SetArgs([]string{"watchlists", "remove-asset", "some-id"})
	err = RootCmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "required flag")
}
