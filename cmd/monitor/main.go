package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/navikt/aad-developer-groups-monitor/pkg/azureclient"
	"github.com/navikt/aad-developer-groups-monitor/pkg/config"
	"github.com/navikt/aad-developer-groups-monitor/pkg/logger"
	"github.com/navikt/aad-developer-groups-monitor/pkg/metrics"
	"github.com/navikt/aad-developer-groups-monitor/pkg/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

const (
	updateMetricsInterval = time.Minute * 5
	exitConfigError       = 1
	exitLoggerError       = 2
	exitRunError          = 3
)

func main() {
	ctx := context.Background()
	cfg, err := config.New(ctx)
	if err != nil {
		fmt.Printf("load config: %s", err)
		os.Exit(exitConfigError)
	}

	log, err := logger.GetLogger(cfg.Log)
	if err != nil {
		fmt.Printf("create logger: %s", err)
		os.Exit(exitLoggerError)
	}

	err = run(ctx, cfg, log)
	if err != nil {
		log.WithError(err).Errorf("error in run()")
		os.Exit(exitRunError)
	}
}

func run(ctx context.Context, cfg *config.Config, log *logrus.Logger) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	bt, _ := version.BuildTime()
	log.Infof("aad-developer-groups-monitor version %s built on %s", version.Version(), bt)

	go func() {
		addr := cfg.Http.ListenAddress
		httpServer := getHttpServer(addr)
		log.Infof("ready to accept requests at %s", addr)

		err := httpServer.ListenAndServe()
		if err != http.ErrServerClosed {
			log.WithError(err).Errorf("unexpected HTTP server error")
		}
		log.Info("HTTP server finished, terminating...")
		cancel()
	}()

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
		sig := <-signals
		log.Infof("received signal %s, terminating...", sig)
		cancel()
	}()

	azureClient := azureclient.NewFromConfig(ctx, cfg.Azure)
	metricsTimer := time.NewTimer(1 * time.Second)

	for ctx.Err() == nil {
		select {
		case <-ctx.Done():
			return nil

		case <-metricsTimer.C:
			log.Infof("update group member count metrics")

			wg := sync.WaitGroup{}
			for _, groupID := range cfg.Groups {
				wg.Add(1)
				go func(ctx context.Context, groupID uuid.UUID) {
					defer wg.Done()

					log.Debugf("get group with ID: %s", groupID)
					group, err := azureClient.GetGroup(ctx, groupID)
					if err != nil {
						log.WithError(err).Errorf("get group")
						return
					}

					log.Debugf("update group member count: %s -> %d", groupID, group.NumMembers)
					metrics.SetDeveloperCount(group.NumMembers, group.Name)
				}(ctx, groupID)
			}

			wg.Wait()
			log.Debugf("metrics updated, next update at %s...", time.Now().Add(updateMetricsInterval))
			metricsTimer.Reset(updateMetricsInterval)
		}
	}

	return nil
}

func getHttpServer(addr string) *http.Server {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {})
	r.Handle("/metrics", promhttp.Handler())

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}
