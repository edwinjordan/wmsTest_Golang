package repository

import (
	"context"

	"github.com/edwinjordan/wmsTest_Golang/domain"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id int) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByAPIKey(ctx context.Context, apiKey string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]*domain.User, int, error)
}

// ProductRepository defines the interface for product data operations
type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id int) (*domain.Product, error)
	GetBySKU(ctx context.Context, sku string) (*domain.Product, error)
	Update(ctx context.Context, product *domain.Product) error
	UpdateQuantity(ctx context.Context, id int, quantity int) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]*domain.Product, int, error)
	Search(ctx context.Context, query string, limit, offset int) ([]*domain.Product, int, error)
}

// LocationRepository defines the interface for location data operations
type LocationRepository interface {
	Create(ctx context.Context, location *domain.Location) error
	GetByID(ctx context.Context, id int) (*domain.Location, error)
	GetByCode(ctx context.Context, code string) (*domain.Location, error)
	Update(ctx context.Context, location *domain.Location) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]*domain.Location, int, error)
	ListByZone(ctx context.Context, zone string, limit, offset int) ([]*domain.Location, int, error)
}

// StockMovementRepository defines the interface for stock movement data operations
type StockMovementRepository interface {
	Create(ctx context.Context, movement *domain.StockMovement) error
	GetByID(ctx context.Context, id int) (*domain.StockMovement, error)
	List(ctx context.Context, filter *domain.StockMovementFilter) ([]*domain.StockMovement, int, error)
	GetByProduct(ctx context.Context, productID int, limit, offset int) ([]*domain.StockMovement, int, error)
	GetByLocation(ctx context.Context, locationID int, limit, offset int) ([]*domain.StockMovement, int, error)
}

// Repositories aggregates all repository interfaces
type Repositories struct {
	User          UserRepository
	Product       ProductRepository
	Location      LocationRepository
	StockMovement StockMovementRepository
}
