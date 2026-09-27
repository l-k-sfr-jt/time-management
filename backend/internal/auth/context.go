package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type contextKey int

const userIDKey contextKey = iota

func withUserID(ctx context.Context, id pgtype.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

// UserIDFromContext returns the authenticated request's internal user id,
// as set by Verifier.Middleware. ok is false outside an authenticated
// request — callers should treat that as a bug (the middleware wasn't
// applied), not as an unauthenticated-user case.
func UserIDFromContext(ctx context.Context) (pgtype.UUID, bool) {
	id, ok := ctx.Value(userIDKey).(pgtype.UUID)
	return id, ok
}
