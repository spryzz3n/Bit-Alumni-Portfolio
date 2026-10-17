package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/spryzz3n/bitalumni/server/internal/config"
	"github.com/spryzz3n/bitalumni/server/internal/db"
	apphttp "github.com/spryzz3n/bitalumni/server/internal/http"
)

type dbStore struct{ db *sql.DB }

func (s dbStore) AppName(ctx context.Context) (string, error) { return db.AppName(ctx, s.db) }

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("config", "err", err)
		os.Exit(1)
	}

	// sql.Open does not connect, so the API still starts (and /health works)
	// when MySQL is down; /db-check reports the failure.
	database, err := db.Open(cfg.DSN())
	if err != nil {
		logger.Error("db open", "err", err)
		os.Exit(1)
	}
	defer database.Close()

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           apphttp.NewRouter(dbStore{database}, logger, cfg.CORSOrigin),
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		logger.Info("listening", "addr", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server", "err", err)
			stop()
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}
