package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/l-k-sfr-jt/time-management/backend/internal/auth"
	"github.com/l-k-sfr-jt/time-management/backend/internal/config"
	"github.com/l-k-sfr-jt/time-management/backend/internal/db"
	"github.com/l-k-sfr-jt/time-management/backend/internal/db/sqlcgen"
	"github.com/l-k-sfr-jt/time-management/backend/internal/httpapi"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	queries := sqlcgen.New(pool)
	jwks := auth.NewJWKSCache(cfg.ClerkJWKSURL, cfg.JWKSCacheTTL)
	verifier := auth.NewVerifier(jwks, queries)

	router := httpapi.NewRouter(pool, queries, verifier, cfg.ClerkWebhookSecret)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("listening on :%s", cfg.Port)
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	case <-ctx.Done():
		log.Println("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	}
}
