package main

import (
	"context"
	httpTrasport "effective-mobile/internal/http"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"effective-mobile/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.MustLoad()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	pool, err := pgxpool.New(ctx, cfg.Postgres.DSN())
	if err != nil {
		fmt.Errorf("failed to connect to postgres") // TODO: Добавить логи
		return
	}
	defer pool.Close()

	if err = pool.Ping(ctx); err != nil {
		fmt.Errorf("failed to ping postgres") // TODO: Добавить логи
		return
	}

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
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Errorf("Server failed to start: %v", err) // TODO: Добавить логи
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
}

// TODO: Добавить логирование
// TODO: Добавить таймауты
