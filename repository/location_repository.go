package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/edwinjordan/wmsTest_Golang/domain"
)

type locationRepository struct {
	db *sql.DB
}

func NewLocationRepository(db *sql.DB) LocationRepository {
	return &locationRepository{db: db}
}

func (r *locationRepository) Create(ctx context.Context, location *domain.Location) error {
	query := `
		INSERT INTO locations (code, name, zone, aisle, rack, shelf, capacity, temperature, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`

	now := time.Now()
	location.CreatedAt = now
	location.UpdatedAt = now
	location.IsActive = true

	err := r.db.QueryRowContext(ctx, query,
		location.Code,
		location.Name,
		location.Zone,
		location.Aisle,
		location.Rack,
		location.Shelf,
		location.Capacity,
		location.Temperature,
		location.IsActive,
		location.CreatedAt,
		location.UpdatedAt,
	).Scan(&location.ID)

	if err != nil {
		return fmt.Errorf("failed to create location: %w", err)
	}

	return nil
}

func (r *locationRepository) GetByID(ctx context.Context, id int) (*domain.Location, error) {
	query := `
		SELECT id, code, name, zone, aisle, rack, shelf, capacity, temperature, is_active, created_at, updated_at
		FROM locations 
		WHERE id = $1 AND is_active = true`

	location := &domain.Location{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&location.ID,
		&location.Code,
		&location.Name,
		&location.Zone,
		&location.Aisle,
		&location.Rack,
		&location.Shelf,
		&location.Capacity,
		&location.Temperature,
		&location.IsActive,
		&location.CreatedAt,
		&location.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get location by ID: %w", err)
	}

	return location, nil
}

func (r *locationRepository) GetByCode(ctx context.Context, code string) (*domain.Location, error) {
	query := `
		SELECT id, code, name, zone, aisle, rack, shelf, capacity, temperature, is_active, created_at, updated_at
		FROM locations 
		WHERE code = $1 AND is_active = true`

	location := &domain.Location{}
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&location.ID,
		&location.Code,
		&location.Name,
		&location.Zone,
		&location.Aisle,
		&location.Rack,
		&location.Shelf,
		&location.Capacity,
		&location.Temperature,
		&location.IsActive,
		&location.CreatedAt,
		&location.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get location by code: %w", err)
	}

	return location, nil
}

func (r *locationRepository) Update(ctx context.Context, location *domain.Location) error {
	query := `
		UPDATE locations 
		SET code = $2, name = $3, zone = $4, aisle = $5, rack = $6, shelf = $7, 
		    capacity = $8, temperature = $9, is_active = $10, updated_at = $11
		WHERE id = $1`

	location.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		location.ID,
		location.Code,
		location.Name,
		location.Zone,
		location.Aisle,
		location.Rack,
		location.Shelf,
		location.Capacity,
		location.Temperature,
		location.IsActive,
		location.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update location: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *locationRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE locations SET is_active = false, updated_at = $2 WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to delete location: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *locationRepository) List(ctx context.Context, limit, offset int) ([]*domain.Location, int, error) {
	// Count total records
	countQuery := `SELECT COUNT(*) FROM locations WHERE is_active = true`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count locations: %w", err)
	}

	// Get paginated records
	query := `
		SELECT id, code, name, zone, aisle, rack, shelf, capacity, temperature, is_active, created_at, updated_at
		FROM locations 
		WHERE is_active = true
		ORDER BY zone, aisle, rack, shelf
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list locations: %w", err)
	}
	defer rows.Close()

	var locations []*domain.Location
	for rows.Next() {
		location := &domain.Location{}
		err := rows.Scan(
			&location.ID,
			&location.Code,
			&location.Name,
			&location.Zone,
			&location.Aisle,
			&location.Rack,
			&location.Shelf,
			&location.Capacity,
			&location.Temperature,
			&location.IsActive,
			&location.CreatedAt,
			&location.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan location: %w", err)
		}
		locations = append(locations, location)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating locations: %w", err)
	}

	return locations, total, nil
}

func (r *locationRepository) ListByZone(ctx context.Context, zone string, limit, offset int) ([]*domain.Location, int, error) {
	// Count total records
	countQuery := `SELECT COUNT(*) FROM locations WHERE zone = $1 AND is_active = true`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, zone).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count locations by zone: %w", err)
	}

	// Get paginated records
	query := `
		SELECT id, code, name, zone, aisle, rack, shelf, capacity, temperature, is_active, created_at, updated_at
		FROM locations 
		WHERE zone = $1 AND is_active = true
		ORDER BY aisle, rack, shelf
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, zone, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list locations by zone: %w", err)
	}
	defer rows.Close()

	var locations []*domain.Location
	for rows.Next() {
		location := &domain.Location{}
		err := rows.Scan(
			&location.ID,
			&location.Code,
			&location.Name,
			&location.Zone,
			&location.Aisle,
			&location.Rack,
			&location.Shelf,
			&location.Capacity,
			&location.Temperature,
			&location.IsActive,
			&location.CreatedAt,
			&location.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan location: %w", err)
		}
		locations = append(locations, location)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating locations by zone: %w", err)
	}

	return locations, total, nil
}
