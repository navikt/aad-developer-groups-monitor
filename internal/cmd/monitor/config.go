package monitor

import (
	"context"

	"github.com/google/uuid"
	"github.com/sethvargo/go-envconfig"
)

type Azure struct {
	// TenantID is the ID of the Azure AD tenant where the groups exist.
	TenantID string `env:"AZURE_APP_TENANT_ID,required"`

	// ClientID is the ID of the API client.
	ClientID string `env:"AZURE_APP_CLIENT_ID,required"`

	// ClientSecret is the secret of the API client.
	ClientSecret string `env:"AZURE_APP_CLIENT_SECRET,required"`
}

type Http struct {
	// ListenAddress is the host:port combination used by the http server.
	//
	// Example: "127.0.0.1:3000"
	ListenAddress string `env:"MONITOR_HTTP_LISTEN_ADDRESS,default=127.0.0.1:3000"`
}

type Log struct {
	// Format can be used to customize the log format.
	//
	// Example: "json"
	Format string `env:"MONITOR_LOG_FORMAT,default=json"`

	// Level is the log level used for logs.
	//
	// Example: "info"
	Level string `env:"MONITOR_LOG_LEVEL,default=info"`
}

type Config struct {
	Azure Azure
	Http  Http
	Log   Log

	// Groups is a list of group IDs that will be included in the metrics.
	Groups []uuid.UUID `env:"MONITOR_GROUP_IDS,required"`
}

func newConfig(ctx context.Context) (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process(ctx, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
