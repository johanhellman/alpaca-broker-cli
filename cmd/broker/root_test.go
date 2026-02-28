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

	// Set up environment variables
	os.Setenv("ALPACA_BROKER_API_KEY", "test-key")       //nolint:errcheck
	os.Setenv("ALPACA_BROKER_API_SECRET", "test-secret") //nolint:errcheck
	os.Setenv("ALPACA_BROKER_ENV", "production")         //nolint:errcheck

	// Call our setup
	initConfig()

	// Verify viper picked up environment variations correctly
	assert.Equal(t, "test-key", viper.GetString("api-key"))
	assert.Equal(t, "test-secret", viper.GetString("api-secret"))
	assert.Equal(t, "production", viper.GetString("env"))

	// Clean up environment variables
	os.Unsetenv("ALPACA_BROKER_API_KEY")    //nolint:errcheck
	os.Unsetenv("ALPACA_BROKER_API_SECRET") //nolint:errcheck
	os.Unsetenv("ALPACA_BROKER_ENV")        //nolint:errcheck
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

func TestPrintOutput_BrokerQuery(t *testing.T) {
	viper.Reset()

	viper.Set("query", "data.user.email")
	testData := map[string]interface{}{
		"data": map[string]interface{}{
			"user": map[string]string{
				"email": "test@example.com",
			},
		},
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := printOutput(testData)
	assert.NoError(t, err)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	assert.Contains(t, buf.String(), "test@example.com")
}
