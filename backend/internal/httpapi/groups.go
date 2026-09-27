package httpapi

import (
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/l-k-sfr-jt/time-management/backend/internal/auth"
	"github.com/l-k-sfr-jt/time-management/backend/internal/db/sqlcgen"
)

func (h *Handlers) CreateGroup(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())

	var req CreateGroupRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	row, err := h.Queries.CreateGroup(r.Context(), sqlcgen.CreateGroupParams{
		UserID: userID,
		Name:   req.Name,
		Color:  req.Color,
		Icon:   req.Icon,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create group")
		return
	}
	writeJSON(w, http.StatusCreated, groupFromRow(row))
}

func (h *Handlers) ListGroups(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())

	rows, err := h.Queries.ListGroups(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list groups")
		return
	}

	out := make([]Group, len(rows))
	for i, row := range rows {
		out[i] = groupFromRow(row)
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *Handlers) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())

	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req UpdateGroupRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	row, err := h.Queries.UpdateGroup(r.Context(), sqlcgen.UpdateGroupParams{
		ID:     id,
		UserID: userID,
		Name:   req.Name,
		Color:  req.Color,
		Icon:   req.Icon,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusNotFound, "group not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update group")
		return
	}
	writeJSON(w, http.StatusOK, groupFromRow(row))
}

// DeleteGroup refuses to delete a group that still has Activity Types
// referencing it (FR-1.3) rather than cascading, so activity history never
// loses its group silently.
func (h *Handlers) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())

	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if _, err := h.Queries.GetGroup(r.Context(), sqlcgen.GetGroupParams{ID: id, UserID: userID}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "group not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to look up group")
		return
	}

	count, err := h.Queries.CountActivityTypesInGroup(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to check group")
		return
	}
	if count > 0 {
		writeError(w, http.StatusConflict, "group still has activity types; reassign or archive them first")
		return
	}

	if _, err := h.Queries.DeleteGroup(r.Context(), sqlcgen.DeleteGroupParams{ID: id, UserID: userID}); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete group")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
