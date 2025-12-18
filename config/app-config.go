// Package config holds the application's configuration settings.
package config

// AppConfig defines environment-based configuration for the application.
type AppConfig struct {
	APIConfig  APIConfig
	DBConfig   DBConfig
	AuthConfig AuthConfig
}

// DBConfig holds settings for database.
type DBConfig struct {
	DriverName string `env:"MARKET_DB_DRIVER_NAME"`
	URI        string `env:"MARKET_DB_URI"`
}

// APIConfig holds settings for API.
type APIConfig struct {
	HTTPAddress string `env:"MARKET_HTTP_ADDRESS"`
}

// AuthConfig holds settings and secrets for auth.
type AuthConfig struct {
	Secret                   string `env:"MARKET_AUTH_SECRET"`
	TokenValidSeconds        int    `env:"MARKET_TOKEN_VALID_SECONDS"`
	RefreshTokenValidSeconds int    `env:"MARKET_REFRESH_TOKEN_VALID_SECONDS"`
	RefreshTokenLength       int    `env:"MARKET_REFRESH_TOKEN_LENGTH"`
}
