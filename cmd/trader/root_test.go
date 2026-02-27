package cmd

import (
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
