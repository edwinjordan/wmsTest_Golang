package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/edwinjordan/wmsTest_Golang/domain"
)

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	query := `
		INSERT INTO products (sku, name, description, price, weight, dimensions, category, quantity, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`

	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now
	product.IsActive = true

	err := r.db.QueryRowContext(ctx, query,
		product.SKU,
		product.Name,
		product.Description,
		product.Price,
		product.Weight,
		product.Dimensions,
		product.Category,
		product.Quantity,
		product.IsActive,
		product.CreatedAt,
		product.UpdatedAt,
	).Scan(&product.ID)

	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

func (r *productRepository) GetByID(ctx context.Context, id int) (*domain.Product, error) {
	query := `
		SELECT id, sku, name, description, price, weight, dimensions, category, quantity, is_active, created_at, updated_at
		FROM products 
		WHERE id = $1 AND is_active = true`

	product := &domain.Product{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.SKU,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Weight,
		&product.Dimensions,
		&product.Category,
		&product.Quantity,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get product by ID: %w", err)
	}

	return product, nil
}

func (r *productRepository) GetBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	query := `
		SELECT id, sku, name, description, price, weight, dimensions, category, quantity, is_active, created_at, updated_at
		FROM products 
		WHERE sku = $1 AND is_active = true`

	product := &domain.Product{}
	err := r.db.QueryRowContext(ctx, query, sku).Scan(
		&product.ID,
		&product.SKU,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Weight,
		&product.Dimensions,
		&product.Category,
		&product.Quantity,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get product by SKU: %w", err)
	}

	return product, nil
}

func (r *productRepository) Update(ctx context.Context, product *domain.Product) error {
	query := `
		UPDATE products 
		SET sku = $2, name = $3, description = $4, price = $5, weight = $6, 
		    dimensions = $7, category = $8, quantity = $9, is_active = $10, updated_at = $11
		WHERE id = $1`

	product.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		product.ID,
		product.SKU,
		product.Name,
		product.Description,
		product.Price,
		product.Weight,
		product.Dimensions,
		product.Category,
		product.Quantity,
		product.IsActive,
		product.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
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

func (r *productRepository) UpdateQuantity(ctx context.Context, id int, quantity int) error {
	query := `UPDATE products SET quantity = $2, updated_at = $3 WHERE id = $1 AND is_active = true`

	result, err := r.db.ExecContext(ctx, query, id, quantity, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update product quantity: %w", err)
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

func (r *productRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE products SET is_active = false, updated_at = $2 WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
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

func (r *productRepository) List(ctx context.Context, limit, offset int) ([]*domain.Product, int, error) {
	// Count total records
	countQuery := `SELECT COUNT(*) FROM products WHERE is_active = true`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// Get paginated records
	query := `
		SELECT id, sku, name, description, price, weight, dimensions, category, quantity, is_active, created_at, updated_at
		FROM products 
		WHERE is_active = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		product := &domain.Product{}
		err := rows.Scan(
			&product.ID,
			&product.SKU,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Weight,
			&product.Dimensions,
			&product.Category,
			&product.Quantity,
			&product.IsActive,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating products: %w", err)
	}

	return products, total, nil
}

func (r *productRepository) Search(ctx context.Context, query string, limit, offset int) ([]*domain.Product, int, error) {
	searchTerm := "%" + strings.ToLower(query) + "%"

	// Count total records
	countQuery := `
		SELECT COUNT(*) 
		FROM products 
		WHERE is_active = true 
		AND (LOWER(name) LIKE $1 OR LOWER(sku) LIKE $1 OR LOWER(category) LIKE $1)`

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, searchTerm).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	// Get paginated records
	searchQuery := `
		SELECT id, sku, name, description, price, weight, dimensions, category, quantity, is_active, created_at, updated_at
		FROM products 
		WHERE is_active = true 
		AND (LOWER(name) LIKE $1 OR LOWER(sku) LIKE $1 OR LOWER(category) LIKE $1)
		ORDER BY 
			CASE 
				WHEN LOWER(sku) = LOWER($4) THEN 1
				WHEN LOWER(name) = LOWER($4) THEN 2
				WHEN LOWER(sku) LIKE $1 THEN 3
				WHEN LOWER(name) LIKE $1 THEN 4
				ELSE 5
			END,
			created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, searchQuery, searchTerm, limit, offset, query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search products: %w", err)
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		product := &domain.Product{}
		err := rows.Scan(
			&product.ID,
			&product.SKU,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Weight,
			&product.Dimensions,
			&product.Category,
			&product.Quantity,
			&product.IsActive,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating search results: %w", err)
	}

	return products, total, nil
}
