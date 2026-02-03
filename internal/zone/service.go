package zone

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/marwan562/fintech-ecosystem/internal/zone/domain"
)

type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateZone(ctx context.Context, params domain.CreateZoneParams) (*domain.Zone, error) {
	if params.OrgID == "" {
		return nil, fmt.Errorf("org_id is required")
	}
	if params.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if params.Mode == "" {
		params.Mode = domain.ModeTest
	}

	id := fmt.Sprintf("zone_%s", strings.ReplaceAll(uuid.NewString(), "-", ""))

	zone := &domain.Zone{
		ID:        id,
		OrgID:     params.OrgID,
		Name:      params.Name,
		Mode:      params.Mode,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, zone); err != nil {
		return nil, fmt.Errorf("failed to create zone: %w", err)
	}

	return zone, nil
}

func (s *Service) GetZone(ctx context.Context, id string) (*domain.Zone, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) ListZones(ctx context.Context, orgID string) ([]*domain.Zone, error) {
	return s.repo.ListByOrgID(ctx, orgID)
}

func (s *Service) DeleteZone(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
