package domain

import (
	"context"
	"time"
)

// ZoneCreatedEvent represents the event data when a zone is created
type ZoneCreatedEvent struct {
	ZoneID    string            `json:"zone_id"`
	OrgID     string            `json:"org_id"`
	Name      string            `json:"name"`
	Mode      Mode              `json:"mode"`
	Timestamp time.Time         `json:"timestamp"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// EventPublisher defines the interface for publishing domain events
type EventPublisher interface {
	PublishZoneCreated(ctx context.Context, event ZoneCreatedEvent) error
}
