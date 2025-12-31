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
	OAuthConfig              OAuthConfig
}

// OAuthConfig holds settings and secrets for oauth providers.
type OAuthConfig struct {
	Github Github
	Google Google
}

// Github holds secrets and settings for github oauth.
type Github struct {
	ClientID     string `env:"OAUTH_GITHUB_CLIENT_ID"`
	ClientSecret string `env:"OAUTH_GITHUB_CLIENT_SECRET"`
	RedirectURI  string `env:"OAUTH_GITHUB_REDIRECT_URI"`
}

// Google holds secrets and settings for github oauth.
type Google struct {
	ClientID     string `env:"OAUTH_GOOGLE_CLIENT_ID"`
	ClientSecret string `env:"OAUTH_GOOGLE_CLIENT_SECRET"`
	RedirectURI  string `env:"OAUTH_GOOGLE_REDIRECT_URI"`
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
