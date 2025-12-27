// Package config holds the application's configuration settings.
package config

// AppConfig defines environment-based configuration for the application.
type AppConfig struct {
	APIConfig      APIConfig
	DBConfig       DBConfig
	AuthConfig     AuthConfig
	StorageConfig  StorageConfig
	PaymentsConfig PaymentsConfig
}

// DBConfig holds settings for database.
type DBConfig struct {
	DriverName string `env:"MARKET_DB_DRIVER_NAME"`
	URI        string `env:"MARKET_DB_URI"`
}

// APIConfig holds settings for API.
type APIConfig struct {
	HTTPAddress    string `env:"MARKET_HTTP_ADDRESS"`
	MaxPayloadSize string `env:"MARKET_MAX_PAYLOAD_SIZE"`
}

// AuthConfig holds settings and secrets for auth.
type AuthConfig struct {
	Secret                   string `env:"MARKET_AUTH_SECRET"`
	TokenValidSeconds        int    `env:"MARKET_TOKEN_VALID_SECONDS"`
	RefreshTokenValidSeconds int    `env:"MARKET_REFRESH_TOKEN_VALID_SECONDS"`
	RefreshTokenLength       int    `env:"MARKET_REFRESH_TOKEN_LENGTH"`
	RefreshTokenCookieSecure bool   `env:"MARKET_REFRESH_TOKEN_COOKIE_SECURE"`
}

// StorageConfig holds settings for storage.
type StorageConfig struct {
	UploadDir           string `env:"MARKET_IMAGE_UPLOAD_DIR"`
	MaxImagesPerListing int    `env:"MARKET_MAX_IMAGES_PER_LISTING"`
}

// PaymentsConfig holds settings for payments.
type PaymentsConfig struct {
	StripeSecretKey     string `env:"STRIPE_SECRET_KEY"`
	StripeWebhookSecret string `env:"STRIPE_WEBHOOK_SECRET"`
}
