package logger

import (
	"fmt"
	"time"

	"github.com/navikt/aad-developer-groups-monitor/pkg/config"
	"github.com/sirupsen/logrus"
)

func GetLogger(cfg config.Log) (*logrus.Logger, error) {
	log := logrus.StandardLogger()

	switch cfg.Format {
	case "json":
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		})
	case "text":
		log.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: time.RFC3339Nano,
		})
	default:
		return nil, fmt.Errorf("invalid log format: %s", cfg.Format)
	}

	lvl, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	log.SetLevel(lvl)

	return log, nil
}
