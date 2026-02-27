package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInitConfig_Env(t *testing.T) {
	// Clear existing viper instance for test isolation
	viper.Reset()
	
	// Set environment variables
	os.Setenv("ALPACA_BROKER_API_KEY", "test-key")
	os.Setenv("ALPACA_BROKER_API_SECRET", "test-secret")
	os.Setenv("ALPACA_BROKER_ENV", "production")

	// Call our setup
	initConfig()

	// Verify viper picked up environment variations correctly
	assert.Equal(t, "test-key", viper.GetString("api-key"))
	assert.Equal(t, "test-secret", viper.GetString("api-secret"))
	assert.Equal(t, "production", viper.GetString("env"))

	// Clean up
	os.Unsetenv("ALPACA_BROKER_API_KEY")
	os.Unsetenv("ALPACA_BROKER_API_SECRET")
	os.Unsetenv("ALPACA_BROKER_ENV")
}

func TestExecuteCommand(t *testing.T) {
	cmd := RootCmd()
	b := new(bytes.Buffer)
	cmd.SetOut(b)
	cmd.SetErr(b)
	cmd.SetArgs([]string{"--help"})
	
	err := cmd.Execute()
	assert.NoError(t, err)

	out := b.String()
	assert.Contains(t, out, "Usage:")
}
