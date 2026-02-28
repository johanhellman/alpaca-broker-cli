package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInitConfig_Trader(t *testing.T) {
	// Clear existing viper instance for test isolation
	viper.Reset()

	// Set up environment variables
	os.Setenv("APCA_API_KEY_ID", "trader-test-key")        //nolint:errcheck
	os.Setenv("APCA_API_SECRET_KEY", "trader-test-secret") //nolint:errcheck
	os.Setenv("APCA_ENV", "paper")                         //nolint:errcheck // Assuming paper default

	// Call our setup
	initConfig()

	// Verify viper bound the env vars correctly
	assert.Equal(t, "trader-test-key", viper.GetString("api-key"))
	assert.Equal(t, "trader-test-secret", viper.GetString("api-secret"))

	// Clean up environment variables
	os.Unsetenv("APCA_API_KEY_ID")     //nolint:errcheck
	os.Unsetenv("APCA_API_SECRET_KEY") //nolint:errcheck
	os.Unsetenv("APCA_ENV")            //nolint:errcheck
}

func TestExecuteCommand_Trader(t *testing.T) {
	// Simple test to ensure Execute doesn't panic on nil arguments etc.
	// We can't easily test the full execute without mocking os.Args or exit commands,
	// but we can ensure the command structure is sound.
	assert.NotNil(t, RootCmd)
	assert.Equal(t, "alpaca-trader", RootCmd.Use)
}

func TestPrintOutput_Query(t *testing.T) {
	viper.Reset()

	// Simulate --query injection
	viper.Set("query", "nested.id")
	testData := map[string]interface{}{
		"nested": map[string]string{
			"id": "12345",
		},
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	assert.NoError(t, err)
	os.Stdout = w

	err = printOutput(testData)
	assert.NoError(t, err)

	assert.NoError(t, w.Close())
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "12345")
}

func TestPrintOutput_CSV(t *testing.T) {
	viper.Reset()
	viper.Set("output", "csv")

	type MockTrade struct {
		ID     string
		Amount float64
	}
	testData := []MockTrade{
		{ID: "trade-1", Amount: 150.5},
		{ID: "trade-2", Amount: 75.0},
	}

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	assert.NoError(t, err)
	os.Stdout = w

	err = printOutput(testData)
	assert.NoError(t, err)

	assert.NoError(t, w.Close())
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	assert.NoError(t, err)

	output := buf.String()
	// Assert headers
	assert.Contains(t, output, "ID,Amount")
	// Assert data rows
	assert.Contains(t, output, "trade-1,150.5")
	assert.Contains(t, output, "trade-2,75")
}
