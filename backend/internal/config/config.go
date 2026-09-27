// Package config loads runtime configuration from environment variables.
package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	// Port the HTTP server listens on.
	Port string
	// DatabaseURL is a Postgres connection string (Neon's direct, non-pooled
	// connection string for this long-running service — see
	// architecture/api-design.md §6 for why not the PgBouncer pooler).
	DatabaseURL string
	// ClerkJWKSURL is Clerk's JSON Web Key Set endpoint, used to verify
	// session JWTs. Format: https://<your-clerk-domain>/.well-known/jwks.json
	ClerkJWKSURL string
	// ClerkWebhookSecret is the signing secret for the Clerk webhook
	// endpoint (starts with "whsec_").
	ClerkWebhookSecret string
	// JWKSCacheTTL controls how long fetched JWKS keys are cached before
	// being refreshed.
	JWKSCacheTTL time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		Port:               getEnvDefault("PORT", "8080"),
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		ClerkJWKSURL:       os.Getenv("CLERK_JWKS_URL"),
		ClerkWebhookSecret: os.Getenv("CLERK_WEBHOOK_SECRET"),
		JWKSCacheTTL:       10 * time.Minute,
	}

	var missing []string
	if cfg.DatabaseURL == "" {
		missing = append(missing, "DATABASE_URL")
	}
	if cfg.ClerkJWKSURL == "" {
		missing = append(missing, "CLERK_JWKS_URL")
	}
	if cfg.ClerkWebhookSecret == "" {
		missing = append(missing, "CLERK_WEBHOOK_SECRET")
	}
	if len(missing) > 0 {
		return Config{}, fmt.Errorf("missing required environment variables: %v", missing)
	}

	return cfg, nil
}

func getEnvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
