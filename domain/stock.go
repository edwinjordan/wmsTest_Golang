package domain

import "time"

// StockMovementType represents the type of stock movement
type StockMovementType string

const (
	StockIN  StockMovementType = "IN"
	StockOUT StockMovementType = "OUT"
)

// StockMovement represents a stock movement transaction
type StockMovement struct {
	ID         int               `json:"id"`
	ProductID  int               `json:"product_id"`
	LocationID int               `json:"location_id"`
	UserID     int               `json:"user_id"`
	Type       StockMovementType `json:"type"`      // IN or OUT
	Quantity   int               `json:"quantity"`  // Always positive, type determines direction
	Reference  string            `json:"reference"` // Reference number (PO, SO, etc.)
	Notes      string            `json:"notes"`
	CreatedAt  time.Time         `json:"created_at"`

	// Populated relations
	Product  *Product  `json:"product,omitempty"`
	Location *Location `json:"location,omitempty"`
	User     *User     `json:"user,omitempty"`
}

// Stock represents current stock level at a location

// CreateStockMovementRequest represents the request to create a stock movement
type CreateStockMovementRequest struct {
	ProductID  int               `json:"product_id" validate:"required"`
	LocationID int               `json:"location_id" validate:"required"`
	Type       StockMovementType `json:"type" validate:"required,oneof=IN OUT"`
	Quantity   int               `json:"quantity" validate:"required,min=1"`
	Reference  string            `json:"reference" validate:"max=100"`
	Notes      string            `json:"notes"`
}

// StockMovementFilter represents filters for stock movement queries
type StockMovementFilter struct {
	ProductID  *int               `json:"product_id,omitempty"`
	LocationID *int               `json:"location_id,omitempty"`
	UserID     *int               `json:"user_id,omitempty"`
	Type       *StockMovementType `json:"type,omitempty"`
	DateFrom   *time.Time         `json:"date_from,omitempty"`
	DateTo     *time.Time         `json:"date_to,omitempty"`
	Limit      int                `json:"limit"`
	Offset     int                `json:"offset"`
}
