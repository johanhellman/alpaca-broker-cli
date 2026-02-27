package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/johanhellman/alpaca-broker-cli/pkg/client"
	"github.com/spf13/viper"
)

// Default URLs
const (
	SandboxAPIURL    = "https://broker-api.sandbox.alpaca.markets/v1"
	ProductionAPIURL = "https://broker-api.alpaca.markets/v1"
)

// AddAuthHeader adds the Basic Auth headers to the request.
type basicAuthProvider struct {
	ApiKey    string
	ApiSecret string
}

func (p *basicAuthProvider) Intercept(ctx context.Context, req *http.Request) error {
	req.SetBasicAuth(p.ApiKey, p.ApiSecret)
	return nil
}

// NewClient creates a new OpenAPI client configured with the current Viper settings.
func NewClient() (*client.ClientWithResponses, error) {
	apiKey := viper.GetString("api-key")
	apiSecret := viper.GetString("api-secret")
	env := viper.GetString("env")

	if apiKey == "" || apiSecret == "" {
		return nil, errors.New("Missing API key or secret. Use --api-key and --api-secret or set ALPACA_BROKER_API_KEY/SECRET in environment variables.")
	}

	serverURL := SandboxAPIURL
	if env == "production" || env == "prod" {
		serverURL = ProductionAPIURL
	} else if env != "sandbox" {
		return nil, errors.New("Invalid env flag. Must be 'sandbox' or 'production'")
	}

	provider := &basicAuthProvider{
		ApiKey:    apiKey,
		ApiSecret: apiSecret,
	}

	c, err := client.NewClientWithResponses(serverURL, client.WithRequestEditorFn(provider.Intercept))
	if err != nil {
		return nil, err
	}

	return c, nil
}
