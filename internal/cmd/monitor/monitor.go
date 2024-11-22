package monitor

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/navikt/aad-developer-groups-monitor/internal/azureclient"
	"github.com/navikt/aad-developer-groups-monitor/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

const (
	exitCodeSuccess = iota
	exitCodeEnvFileError
	exitCodeConfigError
	exitCodeLoggerError
	exitCodeRunError

	updateMetricsInterval = time.Minute * 5
)

func Run(ctx context.Context) {
	log := logrus.StandardLogger()

	if err := loadEnvFile(log); err != nil {
		log.WithError(err).Errorf("error loading .env file")
		os.Exit(exitCodeEnvFileError)
	}

	cfg, err := newConfig(ctx)
	if err != nil {
		fmt.Printf("load config: %s", err)
		os.Exit(exitCodeConfigError)
	}

	appLogger, err := newLogger(cfg.Log.Format, cfg.Log.Level)
	if err != nil {
		fmt.Printf("create logger: %s", err)
		os.Exit(exitCodeLoggerError)
	}

	if err := run(ctx, cfg, appLogger); err != nil {
		appLogger.WithError(err).Errorf("error in run()")
		os.Exit(exitCodeRunError)
	}

	os.Exit(exitCodeSuccess)
}

func run(ctx context.Context, cfg *Config, log logrus.FieldLogger) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

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

	azureClient := azureclient.NewFromConfig(ctx, cfg.Azure.TenantID, cfg.Azure.ClientID, cfg.Azure.ClientSecret)
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
