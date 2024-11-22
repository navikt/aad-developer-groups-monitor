package config

import (
	"context"

	"github.com/google/uuid"
	"github.com/sethvargo/go-envconfig"
)

type Azure struct {
	// TenantID The ID of the Azure AD tenant where the groups exist
	TenantID string `env:"AZURE_APP_TENANT_ID,required"`

	// ClientID The ID of the API client
	ClientID string `env:"AZURE_APP_CLIENT_ID,required"`

	// ClientSecret The secret of the API client
	ClientSecret string `env:"AZURE_APP_CLIENT_SECRET,required"`
}

type Http struct {
	// ListenAddress The host:port combination used by the http server.
	//
	// Example: "127.0.0.1:3000"
	ListenAddress string `env:"MONITOR_HTTP_LISTEN_ADDRESS,default=127.0.0.1:3000"`
}

type Log struct {
	// Format Customize the log format.
	//
	// Example: "text"
	Format string `env:"MONITOR_LOG_FORMAT,default=text"`

	// Level The log level used for logs.
	//
	// Example: "DEBUG"
	Level string `env:"MONITOR_LOG_LEVEL,default=DEBUG"`
}

type Config struct {
	Azure Azure
	Http  Http
	Log   Log

	// Groups A list of group IDs that will be included in the metrics
	Groups []uuid.UUID `env:"MONITOR_GROUP_IDS,required"`
}

func New(ctx context.Context) (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
