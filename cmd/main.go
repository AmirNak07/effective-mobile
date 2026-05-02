package main

import (
	"context"
	httpTrasport "effective-mobile/internal/http"
	"effective-mobile/pkg/logger"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"effective-mobile/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	cfg := config.MustLoad()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	log := logger.NewLogger(cfg.Env)
	go func() {
		if err := log.Sync(); err != nil {
			log.Error(ctx, "failed to sync logger", zap.Error(err))
		}
	}()

	log.Info(ctx, "Starting CRUD service",
		zap.String("env", cfg.Env),
	)

	pool, err := pgxpool.New(ctx, cfg.Postgres.DSN())
	if err != nil {
		log.Error(ctx, "failed to connect to postgres", zap.Error(err))
		return
	}
	defer pool.Close()

	if err = pool.Ping(ctx); err != nil {
		log.Error(ctx, "failed to ping postgres", zap.Error(err))
		return
	}
	log.Info(ctx, "connected to postgres")

	handler := httpTrasport.NewHandler()

	router := httpTrasport.NewRouter(handler)

	server := http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	go func() {
		log.Info(ctx, "CRUD service started")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error(ctx, "server error", zap.Error(err))
		}
	}()

	<-ctx.Done()
	log.Info(ctx, "shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error(ctx, "failed to shutdown server", zap.Error(err))
	}
	log.Info(ctx, "CRUD service stopped")
}
