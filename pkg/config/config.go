package config

import (
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
)

type Azure struct {
	// TenantID The ID of the Azure AD tenant where the groups exist
	TenantID string `envconfig:"AZURE_APP_TENANT_ID" required:"true"`

	// ClientID The ID of the API client
	ClientID string `envconfig:"AZURE_APP_CLIENT_ID" required:"true"`

	// ClientSecret The secret of the API client
	ClientSecret string `envconfig:"AZURE_APP_CLIENT_SECRET" required:"true"`
}

type Http struct {
	// ListenAddress The host:port combination used by the http server.
	//
	// Example: "127.0.0.1:3000"
	ListenAddress string `envconfig:"MONITOR_HTTP_LISTEN_ADDRESS" default:"127.0.0.1:3000"`
}

type Log struct {
	// Format Customize the log format.
	//
	// Example: "text"
	Format string `envconfig:"MONITOR_LOG_FORMAT" default:"text"`

	// Level The log level used for logs.
	//
	// Example: "DEBUG"
	Level string `envconfig:"MONITOR_LOG_LEVEL" default:"DEBUG"`
}

type Config struct {
	Azure Azure
	Http  Http
	Log   Log

	// Groups A list of group IDs that will be included in the metrics
	Groups []uuid.UUID `envconfig:"MONITOR_GROUP_IDS" required:"true"`
}

func New() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
