package config

import (
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
)

type Azure struct {
	// TenantID The ID of the Azure AD tenant where the groups exist
	TenantID string `envconfig:"MONITOR_AZURE_TENANT_ID"`

	// ClientID The ID of the API client
	ClientID string `envconfig:"MONITOR_AZURE_CLIENT_ID"`

	// ClientSecret The secret of the API client
	ClientSecret string `envconfig:"MONITOR_AZURE_CLIENT_SECRET"`
}

type Http struct {
	// ListenAddress The host:port combination used by the http server.
	//
	// Example: "127.0.0.1:3000"
	ListenAddress string `envconfig:"MONITOR_HTTP_LISTEN_ADDRESS"`
}

type Log struct {
	// Format Customize the log format.
	//
	// Example: "text"
	Format string `envconfig:"MONITOR_LOG_FORMAT"`

	// Level The log level used for logs.
	//
	// Example: "DEBUG"
	Level string `envconfig:"MONITOR_LOG_LEVEL"`
}

type Config struct {
	Azure Azure
	Http  Http
	Log   Log

	// Groups A list of group IDs that will be included in the metrics
	Groups []uuid.UUID `envconfig:"MONITOR_GROUP_IDS"`
}

func New() (*Config, error) {
	cfg := defaults()

	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func defaults() *Config {
	return &Config{
		Http: Http{
			ListenAddress: "127.0.0.1:3000",
		},
		Log: Log{
			Format: "text",
			Level:  "DEBUG",
		},
	}
}
