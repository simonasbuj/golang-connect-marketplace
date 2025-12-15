// Package config holds the application's configuration settings.
package config

// AppConfig defines environment-based configuration for the application.
type AppConfig struct {
	APIConfig APIConfig
}

// APIConfig holds settings for API.
type APIConfig struct {
	HTTPAddress string `env:"MARKET_HTTP_ADDRESS"`
}
