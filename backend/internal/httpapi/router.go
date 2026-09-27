package httpapi

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/l-k-sfr-jt/time-management/backend/internal/auth"
	"github.com/l-k-sfr-jt/time-management/backend/internal/db/sqlcgen"
)

// NewRouter wires every endpoint in architecture/api-design.md §3.
// /healthz and /webhooks/clerk are unauthenticated (the webhook verifies
// its own Svix signature instead of a Clerk session JWT); everything
// under /api/v1 requires a verified Clerk bearer token.
func NewRouter(pool *pgxpool.Pool, queries *sqlcgen.Queries, verifier *auth.Verifier, clerkWebhookSecret string) http.Handler {
	h := NewHandlers(queries)

	api := http.NewServeMux()
	api.HandleFunc("POST /groups", h.CreateGroup)
	api.HandleFunc("GET /groups", h.ListGroups)
	api.HandleFunc("PATCH /groups/{id}", h.UpdateGroup)
	api.HandleFunc("DELETE /groups/{id}", h.DeleteGroup)

	api.HandleFunc("POST /activity-types", h.CreateActivityType)
	api.HandleFunc("GET /activity-types", h.ListActivityTypes)
	api.HandleFunc("PATCH /activity-types/{id}", h.UpdateActivityType)
	api.HandleFunc("POST /activity-types/{id}/archive", h.ArchiveActivityType)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", HealthHandler(pool))
	mux.HandleFunc("POST /webhooks/clerk", auth.ClerkWebhookHandler(clerkWebhookSecret, queries))
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", verifier.Middleware(api)))

	return mux
}
