// Package auth verifies Clerk session JWTs and resolves them to a local
// user row, per architecture/data-flows.md §5.
package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/l-k-sfr-jt/time-management/backend/internal/db/sqlcgen"
)

// clerkClaims are the JWT claims we read off a Clerk session token. Email
// is only present if a custom claim named "email" was added to the
// session token template in the Clerk dashboard (Clerk's default template
// does not include it) — it falls back to empty, and the users row keeps
// whatever email the Clerk webhook last synced.
type clerkClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
}

// Verifier verifies Clerk-issued JWTs against Clerk's JWKS and resolves
// the token's subject to a local users.id, upserting the user row if this
// is its first authenticated request (a fallback for when the webhook in
// architecture/data-flows.md §5 hasn't landed yet).
type Verifier struct {
	jwks    *JWKSCache
	queries *sqlcgen.Queries
}

func NewVerifier(jwks *JWKSCache, queries *sqlcgen.Queries) *Verifier {
	return &Verifier{jwks: jwks, queries: queries}
}

// Middleware rejects requests without a valid Clerk bearer token, and
// otherwise attaches the resolved internal user id to the request context
// (retrieve it with UserIDFromContext).
func (v *Verifier) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr, ok := bearerToken(r)
		if !ok {
			writeUnauthorized(w, "missing bearer token")
			return
		}

		claims := &clerkClaims{}
		_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
			kid, _ := t.Header["kid"].(string)
			if kid == "" {
				return nil, errors.New("token header missing kid")
			}
			return v.jwks.Key(r.Context(), kid)
		}, jwt.WithValidMethods([]string{"RS256"}))
		if err != nil {
			writeUnauthorized(w, "invalid token")
			return
		}

		if claims.Subject == "" {
			writeUnauthorized(w, "token missing subject")
			return
		}

		user, err := v.queries.UpsertUser(r.Context(), sqlcgen.UpsertUserParams{
			ClerkUserID: claims.Subject,
			Email:       claims.Email,
		})
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		ctx := withUserID(r.Context(), user.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func bearerToken(r *http.Request) (string, bool) {
	header := r.Header.Get("Authorization")
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return "", false
	}
	token := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	if token == "" {
		return "", false
	}
	return token, true
}

func writeUnauthorized(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusUnauthorized)
}
