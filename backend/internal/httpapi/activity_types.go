package httpapi

import (
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/l-k-sfr-jt/time-management/backend/internal/auth"
	"github.com/l-k-sfr-jt/time-management/backend/internal/db/sqlcgen"
)

func (h *Handlers) CreateActivityType(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())

	var req CreateActivityTypeRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	groupID, err := parseUUID(req.GroupID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid groupId")
		return
	}

	// Verify the group exists and belongs to this user before attaching an
	// Activity Type to it.
	if _, err := h.Queries.GetGroup(r.Context(), sqlcgen.GetGroupParams{ID: groupID, UserID: userID}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusBadRequest, "group not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to look up group")
		return
	}

	row, err := h.Queries.CreateActivityType(r.Context(), sqlcgen.CreateActivityTypeParams{
		UserID:                 userID,
		GroupID:                groupID,
		Name:                   req.Name,
		Description:            req.Description,
		DefaultDurationMinutes: req.DefaultDurationMinutes,
		Color:                  req.Color,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create activity type")
		return
	}
	writeJSON(w, http.StatusCreated, activityTypeFromRow(row))
}

func (h *Handlers) ListActivityTypes(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())

	var groupFilter pgtype.UUID
	if raw := r.URL.Query().Get("groupId"); raw != "" {
		id, err := parseUUID(raw)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid groupId")
			return
		}
		groupFilter = id
	}
	includeArchived := r.URL.Query().Get("includeArchived") == "true"

	rows, err := h.Queries.ListActivityTypes(r.Context(), sqlcgen.ListActivityTypesParams{
		UserID:          userID,
		GroupID:         groupFilter,
		IncludeArchived: includeArchived,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list activity types")
		return
	}

	out := make([]ActivityType, len(rows))
	for i, row := range rows {
		out[i] = activityTypeFromRow(row)
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *Handlers) UpdateActivityType(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())

	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req UpdateActivityTypeRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	row, err := h.Queries.UpdateActivityType(r.Context(), sqlcgen.UpdateActivityTypeParams{
		ID:                     id,
		UserID:                 userID,
		Name:                   req.Name,
		Description:            req.Description,
		DefaultDurationMinutes: req.DefaultDurationMinutes,
		Color:                  req.Color,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusNotFound, "activity type not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update activity type")
		return
	}
	writeJSON(w, http.StatusOK, activityTypeFromRow(row))
}

// ArchiveActivityType is the soft-delete path (FR-2.3): historical Time
// Logs and Planned Entries keep referencing the Activity Type after this.
func (h *Handlers) ArchiveActivityType(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())

	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	row, err := h.Queries.ArchiveActivityType(r.Context(), sqlcgen.ArchiveActivityTypeParams{ID: id, UserID: userID})
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusNotFound, "activity type not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to archive activity type")
		return
	}
	writeJSON(w, http.StatusOK, activityTypeFromRow(row))
}
