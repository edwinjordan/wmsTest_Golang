package domain

import "time"

// Location represents a storage location in the warehouse
type Location struct {
	ID          int       `json:"id"`
	Code        string    `json:"code"` // e.g., "A-01-01" (Zone-Rack-Shelf)
	Name        string    `json:"name"`
	Zone        string    `json:"zone"`        // e.g., "A", "B", "C"
	Aisle       string    `json:"aisle"`       // e.g., "01", "02"
	Rack        string    `json:"rack"`        // e.g., "01", "02"
	Shelf       string    `json:"shelf"`       // e.g., "01", "02", "03"
	Capacity    int       `json:"capacity"`    // Maximum quantity that can be stored
	Temperature *float64  `json:"temperature"` // Optional temperature requirement
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateLocationRequest represents the request to create a new location
type CreateLocationRequest struct {
	Code        string   `json:"code" validate:"required,max=20"`
	Name        string   `json:"name" validate:"required,max=255"`
	Zone        string   `json:"zone" validate:"required,max=10"`
	Aisle       string   `json:"aisle" validate:"required,max=10"`
	Rack        string   `json:"rack" validate:"required,max=10"`
	Shelf       string   `json:"shelf" validate:"required,max=10"`
	Capacity    int      `json:"capacity" validate:"min=1"`
	Temperature *float64 `json:"temperature,omitempty"`
}

// UpdateLocationRequest represents the request to update a location
type UpdateLocationRequest struct {
	Code        *string  `json:"code,omitempty"`
	Name        *string  `json:"name,omitempty"`
	Zone        *string  `json:"zone,omitempty"`
	Aisle       *string  `json:"aisle,omitempty"`
	Rack        *string  `json:"rack,omitempty"`
	Shelf       *string  `json:"shelf,omitempty"`
	Capacity    *int     `json:"capacity,omitempty"`
	Temperature *float64 `json:"temperature,omitempty"`
	IsActive    *bool    `json:"is_active,omitempty"`
}
