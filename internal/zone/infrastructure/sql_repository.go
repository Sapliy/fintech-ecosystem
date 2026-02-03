package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/marwan562/fintech-ecosystem/internal/zone/domain"
)

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) Create(ctx context.Context, zone *domain.Zone) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO zones (id, org_id, name, mode, created_at, updated_at) 
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING created_at`,
		zone.ID, zone.OrgID, zone.Name, zone.Mode, zone.CreatedAt, zone.UpdatedAt).
		Scan(&zone.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create zone: %w", err)
	}
	return nil
}

func (r *SQLRepository) GetByID(ctx context.Context, id string) (*domain.Zone, error) {
	var zone domain.Zone
	err := r.db.QueryRowContext(ctx,
		`SELECT id, org_id, name, mode, created_at, updated_at FROM zones WHERE id = $1`,
		id).Scan(&zone.ID, &zone.OrgID, &zone.Name, &zone.Mode, &zone.CreatedAt, &zone.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get zone: %w", err)
	}
	return &zone, nil
}

func (r *SQLRepository) ListByOrgID(ctx context.Context, orgID string) ([]*domain.Zone, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, org_id, name, mode, created_at, updated_at FROM zones WHERE org_id = $1 ORDER BY created_at DESC`,
		orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to list zones: %w", err)
	}
	defer rows.Close()

	var zones []*domain.Zone
	for rows.Next() {
		var zone domain.Zone
		if err := rows.Scan(&zone.ID, &zone.OrgID, &zone.Name, &zone.Mode, &zone.CreatedAt, &zone.UpdatedAt); err != nil {
			return nil, err
		}
		zones = append(zones, &zone)
	}
	return zones, nil
}

func (r *SQLRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM zones WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete zone: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("zone not found")
	}
	return nil
}
