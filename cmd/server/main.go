// Command server runs the GetNotified API, the check engine and the
// notification fan-out in one process.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivermigrate"

	"github.com/fajrlabs/getnotified/internal/httpapi"
	"github.com/fajrlabs/getnotified/internal/jobs"
	"github.com/fajrlabs/getnotified/internal/store"
)

const shutdownGrace = 20 * time.Second

// config is read once at boot. Nothing here has a default: a missing value is
// a misconfigured deployment, not something to guess at.
type config struct {
	databaseURL      string
	databaseAdminURL string
	port             string
	apiKey           string
	orgSlug          string
	orgName          string
}

func loadConfig() (config, error) {
	var missing []string
	get := func(key string) string {
		v := os.Getenv(key)
		if v == "" {
			missing = append(missing, key)
		}
		return v
	}

	c := config{
		databaseURL:      get("DATABASE_URL"),
		databaseAdminURL: get("DATABASE_ADMIN_URL"),
		port:             get("PORT"),
		apiKey:           get("API_KEY"),
		orgSlug:          get("ORG_SLUG"),
		orgName:          get("ORG_NAME"),
	}
	if len(missing) > 0 {
		return c, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}
	return c, nil
}

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	orgID, err := bootstrap(ctx, cfg)
	if err != nil {
		return err
	}

	pool, err := store.NewPool(ctx, cfg.databaseURL, orgID)
	if err != nil {
		return err
	}
	defer pool.Close()

	client, err := river.NewClient(riverpgxv5.New(pool), &river.Config{
		Queues:       map[string]river.QueueConfig{river.QueueDefault: {MaxWorkers: max(4, runtime.NumCPU()*4)}},
		Workers:      jobs.Workers(pool),
		PeriodicJobs: jobs.PeriodicJobs(),
	})
	if err != nil {
		return err
	}
	if err := client.Start(ctx); err != nil {
		return err
	}

	api := &httpapi.API{Pool: pool, OrgID: orgID, APIKey: cfg.apiKey}
	srv := &http.Server{
		Addr:              ":" + cfg.port,
		Handler:           api.Routes(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		slog.Info("listening", "addr", srv.Addr, "org", cfg.orgSlug)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("http server", "err", err)
			stop()
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), shutdownGrace)
	defer cancel()
	srv.Shutdown(shutdownCtx)
	return client.Stop(shutdownCtx)
}

// bootstrap runs as the admin role: schema, River's migrations, the org row
// and the grants all sit outside what row-level security lets the app role do.
func bootstrap(ctx context.Context, cfg config) (string, error) {
	admin, err := pgxpool.New(ctx, cfg.databaseAdminURL)
	if err != nil {
		return "", err
	}
	defer admin.Close()

	if err := store.Migrate(ctx, admin); err != nil {
		return "", err
	}

	// River's own migrations span more than one transaction, so they run
	// outside the schema statements above.
	migrator, err := rivermigrate.New(riverpgxv5.New(admin), nil)
	if err != nil {
		return "", err
	}
	if _, err := migrator.Migrate(ctx, rivermigrate.DirectionUp, nil); err != nil {
		return "", err
	}

	orgID, err := store.EnsureOrg(ctx, admin, cfg.orgSlug, cfg.orgName)
	if err != nil {
		return "", err
	}
	return orgID, store.Grants(ctx, admin)
}
