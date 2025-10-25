package domain

import "time"

// Product represents a product in the warehouse
type Product struct {
	ID          int       `json:"id"`
	SKU         string    `json:"sku"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Weight      float64   `json:"weight"`     // in kg
	Dimensions  string    `json:"dimensions"` // "LxWxH" in cm
	Category    string    `json:"category"`
	Quantity    int       `json:"quantity"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateProductRequest represents the request to create a new product
type CreateProductRequest struct {
	SKU         string  `json:"sku" validate:"required,max=50"`
	Name        string  `json:"name" validate:"required,max=255"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"min=0"`
	Weight      float64 `json:"weight" validate:"min=0"`
	Dimensions  string  `json:"dimensions"`
	Category    string  `json:"category" validate:"required,max=100"`
	Quantity    int     `json:"quantity" validate:"required,min=0"`
}

// UpdateProductRequest represents the request to update a product
type UpdateProductRequest struct {
	SKU         *string  `json:"sku,omitempty"`
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Weight      *float64 `json:"weight,omitempty"`
	Dimensions  *string  `json:"dimensions,omitempty"`
	Category    *string  `json:"category,omitempty"`
	IsActive    *bool    `json:"is_active,omitempty"`
	Quantity    *int     `json:"quantity,omitempty"`
}
