package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/edwinjordan/wmsTest_Golang/domain"
)

type stockMovementRepository struct {
	db *sql.DB
}

func NewStockMovementRepository(db *sql.DB) StockMovementRepository {
	return &stockMovementRepository{db: db}
}

func (r *stockMovementRepository) Create(ctx context.Context, movement *domain.StockMovement) error {
	query := `
		INSERT INTO stock_movements (product_id, location_id, user_id, type, quantity, reference, notes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	movement.CreatedAt = time.Now()

	err := r.db.QueryRowContext(ctx, query,
		movement.ProductID,
		movement.LocationID,
		movement.UserID,
		movement.Type,
		movement.Quantity,
		movement.Reference,
		movement.Notes,
		movement.CreatedAt,
	).Scan(&movement.ID)

	if err != nil {
		return fmt.Errorf("failed to create stock movement: %w", err)
	}

	return nil
}

func (r *stockMovementRepository) GetByID(ctx context.Context, id int) (*domain.StockMovement, error) {
	query := `
		SELECT sm.id, sm.product_id, sm.location_id, sm.user_id, sm.type, sm.quantity, sm.reference, sm.notes, sm.created_at,
		       p.id, p.sku, p.name, p.description, p.price, p.weight, p.dimensions, p.category, p.is_active, p.created_at, p.updated_at,
		       l.id, l.code, l.name, l.zone, l.aisle, l.rack, l.shelf, l.capacity, l.temperature, l.is_active, l.created_at, l.updated_at,
		       u.id, u.username, u.email, u.password, u.api_key, u.is_active, u.created_at, u.updated_at
		FROM stock_movements sm
		JOIN products p ON sm.product_id = p.id
		JOIN locations l ON sm.location_id = l.id
		JOIN users u ON sm.user_id = u.id
		WHERE sm.id = $1`

	movement := &domain.StockMovement{}
	product := &domain.Product{}
	location := &domain.Location{}
	user := &domain.User{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&movement.ID, &movement.ProductID, &movement.LocationID, &movement.UserID, &movement.Type, &movement.Quantity, &movement.Reference, &movement.Notes, &movement.CreatedAt,
		&product.ID, &product.SKU, &product.Name, &product.Description, &product.Price, &product.Weight, &product.Dimensions, &product.Category, &product.IsActive, &product.CreatedAt, &product.UpdatedAt,
		&location.ID, &location.Code, &location.Name, &location.Zone, &location.Aisle, &location.Rack, &location.Shelf, &location.Capacity, &location.Temperature, &location.IsActive, &location.CreatedAt, &location.UpdatedAt,
		&user.ID, &user.Username, &user.Email, &user.Password, &user.APIKey, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get stock movement by ID: %w", err)
	}

	movement.Product = product
	movement.Location = location
	movement.User = user

	return movement, nil
}

func (r *stockMovementRepository) List(ctx context.Context, filter *domain.StockMovementFilter) ([]*domain.StockMovement, int, error) {
	// Build WHERE conditions
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter.ProductID != nil {
		conditions = append(conditions, fmt.Sprintf("sm.product_id = $%d", argIndex))
		args = append(args, *filter.ProductID)
		argIndex++
	}

	if filter.LocationID != nil {
		conditions = append(conditions, fmt.Sprintf("sm.location_id = $%d", argIndex))
		args = append(args, *filter.LocationID)
		argIndex++
	}

	if filter.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("sm.user_id = $%d", argIndex))
		args = append(args, *filter.UserID)
		argIndex++
	}

	if filter.Type != nil {
		conditions = append(conditions, fmt.Sprintf("sm.type = $%d", argIndex))
		args = append(args, *filter.Type)
		argIndex++
	}

	if filter.DateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("sm.created_at >= $%d", argIndex))
		args = append(args, *filter.DateFrom)
		argIndex++
	}

	if filter.DateTo != nil {
		conditions = append(conditions, fmt.Sprintf("sm.created_at <= $%d", argIndex))
		args = append(args, *filter.DateTo)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total records
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM stock_movements sm 
		%s`, whereClause)

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count stock movements: %w", err)
	}

	// Add pagination parameters
	paginationClause := fmt.Sprintf("ORDER BY sm.created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	// Get paginated records
	query := fmt.Sprintf(`
		SELECT sm.id, sm.product_id, sm.location_id, sm.user_id, sm.type, sm.quantity, sm.reference, sm.notes, sm.created_at,
		       p.id, p.sku, p.name, p.description, p.price, p.weight, p.dimensions, p.category, p.is_active, p.created_at, p.updated_at,p.quantity,
		       l.id, l.code, l.name, l.zone, l.aisle, l.rack, l.shelf, l.capacity, l.temperature, l.is_active, l.created_at, l.updated_at,
		       u.id, u.username, u.email, u.password, u.api_key, u.is_active, u.created_at, u.updated_at
		FROM stock_movements sm
		JOIN products p ON sm.product_id = p.id
		JOIN locations l ON sm.location_id = l.id
		JOIN users u ON sm.user_id = u.id
		%s
		%s`, whereClause, paginationClause)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list stock movements: %w", err)
	}
	defer rows.Close()

	var movements []*domain.StockMovement
	for rows.Next() {
		movement := &domain.StockMovement{}
		product := &domain.Product{}
		location := &domain.Location{}
		user := &domain.User{}

		err := rows.Scan(
			&movement.ID, &movement.ProductID, &movement.LocationID, &movement.UserID, &movement.Type, &movement.Quantity, &movement.Reference, &movement.Notes, &movement.CreatedAt,
			&product.ID, &product.SKU, &product.Name, &product.Description, &product.Price, &product.Weight, &product.Dimensions, &product.Category, &product.IsActive, &product.CreatedAt, &product.UpdatedAt, &product.Quantity,
			&location.ID, &location.Code, &location.Name, &location.Zone, &location.Aisle, &location.Rack, &location.Shelf, &location.Capacity, &location.Temperature, &location.IsActive, &location.CreatedAt, &location.UpdatedAt,
			&user.ID, &user.Username, &user.Email, &user.Password, &user.APIKey, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan stock movement: %w", err)
		}

		movement.Product = product
		movement.Location = location
		movement.User = user
		movements = append(movements, movement)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating stock movements: %w", err)
	}

	return movements, total, nil
}

func (r *stockMovementRepository) GetByProduct(ctx context.Context, productID int, limit, offset int) ([]*domain.StockMovement, int, error) {
	filter := &domain.StockMovementFilter{
		ProductID: &productID,
		Limit:     limit,
		Offset:    offset,
	}
	return r.List(ctx, filter)
}

func (r *stockMovementRepository) GetByLocation(ctx context.Context, locationID int, limit, offset int) ([]*domain.StockMovement, int, error) {
	filter := &domain.StockMovementFilter{
		LocationID: &locationID,
		Limit:      limit,
		Offset:     offset,
	}
	return r.List(ctx, filter)
}
