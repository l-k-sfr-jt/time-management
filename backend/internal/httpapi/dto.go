package httpapi

import (
	"time"

	"github.com/l-k-sfr-jt/time-management/backend/internal/db/sqlcgen"
)

type Group struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Color     *string   `json:"color,omitempty"`
	Icon      *string   `json:"icon,omitempty"`
	IsDefault bool      `json:"isDefault"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func groupFromRow(g sqlcgen.Group) Group {
	return Group{
		ID:        g.ID.String(),
		Name:      g.Name,
		Color:     g.Color,
		Icon:      g.Icon,
		IsDefault: g.IsDefault,
		CreatedAt: g.CreatedAt.Time,
		UpdatedAt: g.UpdatedAt.Time,
	}
}

type CreateGroupRequest struct {
	Name  string  `json:"name"`
	Color *string `json:"color"`
	Icon  *string `json:"icon"`
}

type UpdateGroupRequest struct {
	Name  *string `json:"name"`
	Color *string `json:"color"`
	Icon  *string `json:"icon"`
}

type ActivityType struct {
	ID                     string     `json:"id"`
	GroupID                string     `json:"groupId"`
	Name                   string     `json:"name"`
	Description            *string    `json:"description,omitempty"`
	DefaultDurationMinutes *int32     `json:"defaultDurationMinutes,omitempty"`
	Color                  *string    `json:"color,omitempty"`
	ArchivedAt             *time.Time `json:"archivedAt,omitempty"`
	CreatedAt              time.Time  `json:"createdAt"`
	UpdatedAt              time.Time  `json:"updatedAt"`
}

func activityTypeFromRow(a sqlcgen.ActivityType) ActivityType {
	dto := ActivityType{
		ID:                     a.ID.String(),
		GroupID:                a.GroupID.String(),
		Name:                   a.Name,
		Description:            a.Description,
		DefaultDurationMinutes: a.DefaultDurationMinutes,
		Color:                  a.Color,
		CreatedAt:              a.CreatedAt.Time,
		UpdatedAt:              a.UpdatedAt.Time,
	}
	if a.ArchivedAt.Valid {
		dto.ArchivedAt = &a.ArchivedAt.Time
	}
	return dto
}

type CreateActivityTypeRequest struct {
	Name                   string  `json:"name"`
	GroupID                string  `json:"groupId"`
	Description            *string `json:"description"`
	DefaultDurationMinutes *int32  `json:"defaultDurationMinutes"`
	Color                  *string `json:"color"`
}

type UpdateActivityTypeRequest struct {
	Name                   *string `json:"name"`
	Description            *string `json:"description"`
	DefaultDurationMinutes *int32  `json:"defaultDurationMinutes"`
	Color                  *string `json:"color"`
}
