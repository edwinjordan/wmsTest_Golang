package service

import (
	"context"
	"fmt"

	"github.com/edwinjordan/wmsTest_Golang/domain"
	"github.com/edwinjordan/wmsTest_Golang/repository"
)

type StockService interface {
	ProcessStockMovement(ctx context.Context, req *domain.CreateStockMovementRequest, userID int) (*domain.StockMovement, error)
	GetStockMovements(ctx context.Context, filter *domain.StockMovementFilter) ([]*domain.StockMovement, int, error)
	GetStockMovementByID(ctx context.Context, id int) (*domain.StockMovement, error)
}

type stockService struct {
	stockMovementRepo repository.StockMovementRepository
	productRepo       repository.ProductRepository
	locationRepo      repository.LocationRepository
}

func NewStockService(
	stockMovementRepo repository.StockMovementRepository,
	productRepo repository.ProductRepository,
	locationRepo repository.LocationRepository,
) StockService {
	return &stockService{
		stockMovementRepo: stockMovementRepo,
		productRepo:       productRepo,
		locationRepo:      locationRepo,
	}
}

func (s *stockService) ProcessStockMovement(ctx context.Context, req *domain.CreateStockMovementRequest, userID int) (*domain.StockMovement, error) {
	// Validate product exists and is active
	product, err := s.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Validate location exists and is active
	location, err := s.locationRepo.GetByID(ctx, req.LocationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get location: %w", err)
	}

	// Get current stock for validation from product quantity
	currentQuantity := product.Quantity

	// Business Rule 1: Stock OUT tidak boleh melebihi stok tersedia
	if req.Type == domain.StockOUT {
		if currentQuantity < req.Quantity {
			return nil, domain.ErrInsufficientStock
		}
	}

	// Create stock movement record
	movement := &domain.StockMovement{
		ProductID:  req.ProductID,
		LocationID: req.LocationID,
		UserID:     userID,
		Type:       req.Type,
		Quantity:   req.Quantity,
		Reference:  req.Reference,
		Notes:      req.Notes,
	}

	err = s.stockMovementRepo.Create(ctx, movement)
	if err != nil {
		return nil, fmt.Errorf("failed to create stock movement: %w", err)
	}

	// Business Rule 3: Quantity produk auto-update saat ada pergerakan stok
	var newQuantity int
	if req.Type == domain.StockIN {
		newQuantity = currentQuantity + req.Quantity
	} else {
		newQuantity = currentQuantity - req.Quantity
	}

	// Update product quantity
	err = s.productRepo.UpdateQuantity(ctx, req.ProductID, newQuantity)
	if err != nil {
		return nil, fmt.Errorf("failed to update product quantity: %w", err)
	}

	// Populate movement with related data
	movement.Product = product
	movement.Location = location

	return movement, nil
}

func (s *stockService) GetStockMovements(ctx context.Context, filter *domain.StockMovementFilter) ([]*domain.StockMovement, int, error) {
	movements, total, err := s.stockMovementRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get stock movements: %w", err)
	}

	return movements, total, nil
}

func (s *stockService) GetStockMovementByID(ctx context.Context, id int) (*domain.StockMovement, error) {
	movement, err := s.stockMovementRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock movement: %w", err)
	}

	return movement, nil
}
