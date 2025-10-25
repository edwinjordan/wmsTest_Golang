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
	if !product.IsActive {
		return nil, fmt.Errorf("product is not active")
	}

	// Validate location exists and is active
	location, err := s.locationRepo.GetByID(ctx, req.LocationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get location: %w", err)
	}
	if !location.IsActive {
		return nil, fmt.Errorf("location is not active")
	}

	// Get current stock for validation from product quantity
	currentQuantity := product.Quantity

	// Debug logging - check request type
	fmt.Printf("DEBUG: Request Type='%s', StockIN='%s', StockOUT='%s'\n", req.Type, domain.StockIN, domain.StockOUT)

	// Business Rule 1: Stock OUT tidak boleh melebihi stok tersedia
	if req.Type == domain.StockOUT {
		fmt.Printf("DEBUG: Processing Stock OUT validation\n")
		if currentQuantity < req.Quantity {
			return nil, domain.ErrInsufficientStock
		}
	}

	// Business Rule 2: Stock IN tidak boleh melebihi kapasitas lokasi
	if req.Type == domain.StockIN {
		fmt.Printf("DEBUG: Processing Stock IN validation\n")
		// Get all stock movements for this location to calculate current occupancy
		filter := &domain.StockMovementFilter{
			LocationID: &req.LocationID,
			Limit:      1000, // Get all movements for this location
			Offset:     0,
		}

		movements, _, err := s.stockMovementRepo.List(ctx, filter)
		if err != nil {
			return nil, fmt.Errorf("failed to get location stock movements: %w", err)
		}

		// Calculate current total quantity in location by summing all movements
		currentLocationQuantity := 0
		for _, movement := range movements {
			if movement.Type == domain.StockIN {
				currentLocationQuantity += movement.Quantity
			} else {
				currentLocationQuantity -= movement.Quantity
			}
		}

		// Check if adding new quantity would exceed location capacity
		newTotalQuantity := currentLocationQuantity + req.Quantity

		// Debug logging (remove in production)
		fmt.Printf("DEBUG: Location ID=%d, Current quantity=%d, Adding=%d, Total will be=%d, Capacity=%d\n",
			req.LocationID, currentLocationQuantity, req.Quantity, newTotalQuantity, location.Capacity)

		if newTotalQuantity > location.Capacity {
			fmt.Printf("DEBUG: CAPACITY EXCEEDED! Returning ErrExceedsCapacity\n")
			return nil, domain.ErrExceedsCapacity
		} else {
			fmt.Printf("DEBUG: Capacity OK - continuing with movement creation\n")
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
