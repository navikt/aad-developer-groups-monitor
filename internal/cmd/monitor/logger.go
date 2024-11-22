package monitor

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

func newLogger(format, level string) (logrus.FieldLogger, error) {
	log := logrus.StandardLogger()

	switch strings.ToLower(format) {
	case "json":
		log.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	default:
		return nil, fmt.Errorf("invalid log format: %q", format)
	}

	parsedLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level %q: %w", level, err)
	}

	log.SetLevel(parsedLevel)

	return log, nil
}
