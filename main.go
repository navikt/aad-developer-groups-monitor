package main

import (
	"context"

	"github.com/navikt/aad-developer-groups-monitor/internal/cmd/monitor"
)

func main() {
	monitor.Run(context.Background())
}
