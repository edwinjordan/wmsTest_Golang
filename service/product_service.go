package service

import (
	"context"
	"fmt"

	"github.com/edwinjordan/wmsTest_Golang/domain"
	"github.com/edwinjordan/wmsTest_Golang/repository"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req *domain.CreateProductRequest) (*domain.Product, error)
	GetProductByID(ctx context.Context, id int) (*domain.Product, error)
	GetProductBySKU(ctx context.Context, sku string) (*domain.Product, error)
	UpdateProduct(ctx context.Context, id int, req *domain.UpdateProductRequest) (*domain.Product, error)
	DeleteProduct(ctx context.Context, id int) error
	ListProducts(ctx context.Context, limit, offset int) ([]*domain.Product, int, error)
	SearchProducts(ctx context.Context, query string, limit, offset int) ([]*domain.Product, int, error)
}

type productService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productService{
		productRepo: productRepo,
	}
}

func (s *productService) CreateProduct(ctx context.Context, req *domain.CreateProductRequest) (*domain.Product, error) {
	// Check if SKU already exists
	existingProduct, err := s.productRepo.GetBySKU(ctx, req.SKU)
	if err == nil && existingProduct != nil {
		return nil, domain.ErrDuplicateEntry
	}

	product := &domain.Product{
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Weight:      req.Weight,
		Dimensions:  req.Dimensions,
		Category:    req.Category,
		Quantity:    req.Quantity,
	}

	err = s.productRepo.Create(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

func (s *productService) GetProductByID(ctx context.Context, id int) (*domain.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

func (s *productService) GetProductBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	product, err := s.productRepo.GetBySKU(ctx, sku)
	if err != nil {
		return nil, fmt.Errorf("failed to get product by SKU: %w", err)
	}

	return product, nil
}

func (s *productService) UpdateProduct(ctx context.Context, id int, req *domain.UpdateProductRequest) (*domain.Product, error) {
	// Get existing product
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Update fields if provided
	if req.SKU != nil {
		// Check if new SKU already exists (excluding current product)
		existingProduct, err := s.productRepo.GetBySKU(ctx, *req.SKU)
		if err == nil && existingProduct != nil && existingProduct.ID != id {
			return nil, domain.ErrDuplicateEntry
		}
		product.SKU = *req.SKU
	}
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Weight != nil {
		product.Weight = *req.Weight
	}
	if req.Dimensions != nil {
		product.Dimensions = *req.Dimensions
	}
	if req.Category != nil {
		product.Category = *req.Category
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	if req.Quantity != nil {
		product.Quantity = *req.Quantity
	}

	err = s.productRepo.Update(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return product, nil
}

func (s *productService) DeleteProduct(ctx context.Context, id int) error {
	err := s.productRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

func (s *productService) ListProducts(ctx context.Context, limit, offset int) ([]*domain.Product, int, error) {
	products, total, err := s.productRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}

	return products, total, nil
}

func (s *productService) SearchProducts(ctx context.Context, query string, limit, offset int) ([]*domain.Product, int, error) {
	products, total, err := s.productRepo.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search products: %w", err)
	}

	return products, total, nil
}
