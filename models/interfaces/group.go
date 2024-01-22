package interfaces

import (
	"context"
	"database/sql"
	"time"
)

type Group struct {
	GroupID     int64     `json:"group_id"`
	OwnerID     int64     `json:"owner_id"`
	ScopeID     int64     `json:"scope_id"`
	GroupName   string    `json:"group_name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GroupService interface {
	InsertGroup(ctx context.Context, group *Group, otx ...*sql.Tx) error
	UpdateGroup(ctx context.Context, group *Group, otx ...*sql.Tx) error
	DeleteGroup(ctx context.Context, groupID int64, otx ...*sql.Tx) error
	GetGroupByID(ctx context.Context, groupID int64, otx ...*sql.Tx) (*Group, error)
}
