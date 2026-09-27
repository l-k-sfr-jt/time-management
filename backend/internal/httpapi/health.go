package httpapi

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

// HealthHandler returns 200 only if the database is reachable, so it
// doubles as a readiness check.
func HealthHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := pool.Ping(r.Context()); err != nil {
			writeError(w, http.StatusServiceUnavailable, "database unreachable")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}
