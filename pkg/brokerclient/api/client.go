package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/johanhellman/alpaca-broker-cli/pkg/brokerclient"
	"github.com/spf13/viper"
)

// Default URLs
const (
	SandboxAPIURL    = "https://broker-api.sandbox.alpaca.markets/v1"
	ProductionAPIURL = "https://broker-api.alpaca.markets/v1"
	PaperAPIURL      = "https://broker-api.sandbox.alpaca.markets/v1"
)

// AddAuthHeader adds the Auth headers to the request.
type authProvider struct {
	ApiKey    string
	ApiSecret string
	IsPaper   bool
}

func (p *authProvider) Intercept(ctx context.Context, req *http.Request) error {
	if p.IsPaper {
		req.Header.Set("APCA-API-KEY-ID", p.ApiKey)
		req.Header.Set("APCA-API-SECRET-KEY", p.ApiSecret)
	} else {
		req.SetBasicAuth(p.ApiKey, p.ApiSecret)
	}
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
	isPaper := false
	if env == "production" || env == "prod" {
		serverURL = ProductionAPIURL
	} else if env == "paper" {
		serverURL = PaperAPIURL
		isPaper = true
	} else if env != "sandbox" {
		return nil, errors.New("Invalid env flag. Must be 'sandbox', 'production', or 'paper'")
	}

	provider := &authProvider{
		ApiKey:    apiKey,
		ApiSecret: apiSecret,
		IsPaper:   isPaper,
	}

	c, err := client.NewClientWithResponses(serverURL, client.WithRequestEditorFn(provider.Intercept))
	if err != nil {
		return nil, err
	}

	return c, nil
}
