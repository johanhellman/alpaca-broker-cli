package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarketDataCmd_TraderFlags(t *testing.T) {
	assert.NotNil(t, marketDataCmd)
	assert.Equal(t, "market-data", marketDataCmd.Use)

	// Check Bars Topology
	assert.NotNil(t, marketDataBarsCmd)
	bFlags := marketDataBarsCmd.Flags()
	assert.NotNil(t, bFlags.Lookup("start"))
	assert.NotNil(t, bFlags.Lookup("end"))
	assert.NotNil(t, bFlags.Lookup("timeframe"))
	assert.NotNil(t, bFlags.Lookup("adjustment"))

	// Check Quotes Topology
	assert.NotNil(t, marketDataQuotesCmd)
	qFlags := marketDataQuotesCmd.Flags()
	assert.NotNil(t, qFlags.Lookup("sort"))
	assert.NotNil(t, qFlags.Lookup("currency"))
	assert.Nil(t, qFlags.Lookup("timeframe"))

	// Check Trades Topology
	assert.NotNil(t, marketDataTradesCmd)
	tFlags := marketDataTradesCmd.Flags()
	assert.NotNil(t, tFlags.Lookup("feed"))
	assert.NotNil(t, tFlags.Lookup("as-of"))

	// Test command execution (help output)
	b := bytes.NewBufferString("")
	RootCmd.SetOut(b)
	RootCmd.SetArgs([]string{"market-data", "--help"})
	err := RootCmd.Execute()
	assert.NoError(t, err)

	out := b.String()
	assert.Contains(t, out, marketDataCmd.Short)
}
