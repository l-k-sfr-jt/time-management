// Package httpapi implements the HTTP handlers described in
// architecture/api-design.md.
package httpapi

import "github.com/l-k-sfr-jt/time-management/backend/internal/db/sqlcgen"

type Handlers struct {
	Queries *sqlcgen.Queries
}

func NewHandlers(queries *sqlcgen.Queries) *Handlers {
	return &Handlers{Queries: queries}
}
