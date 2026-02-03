package domain

import (
	"context"
	"time"
)

type Mode string

const (
	ModeTest Mode = "test"
	ModeLive Mode = "live"
)

type Zone struct {
	ID        string    `json:"id"`
	OrgID     string    `json:"org_id"`
	Name      string    `json:"name"`
	Mode      Mode      `json:"mode"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateZoneParams struct {
	OrgID        string
	Name         string
	Mode         Mode
	TemplateName string
}

type Repository interface {
	Create(ctx context.Context, zone *Zone) error
	GetByID(ctx context.Context, id string) (*Zone, error)
	ListByOrgID(ctx context.Context, orgID string) ([]*Zone, error)
	Delete(ctx context.Context, id string) error
}

type Service interface {
	CreateZone(ctx context.Context, params CreateZoneParams) (*Zone, error)
	GetZone(ctx context.Context, id string) (*Zone, error)
	ListZones(ctx context.Context, orgID string) ([]*Zone, error)
	DeleteZone(ctx context.Context, id string) error
}
