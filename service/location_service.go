package service

import (
	"context"
	"fmt"

	"github.com/edwinjordan/wmsTest_Golang/domain"
	"github.com/edwinjordan/wmsTest_Golang/repository"
)

type LocationService interface {
	CreateLocation(ctx context.Context, req *domain.CreateLocationRequest) (*domain.Location, error)
	GetLocationByID(ctx context.Context, id int) (*domain.Location, error)
	GetLocationByCode(ctx context.Context, code string) (*domain.Location, error)
	UpdateLocation(ctx context.Context, id int, req *domain.UpdateLocationRequest) (*domain.Location, error)
	DeleteLocation(ctx context.Context, id int) error
	ListLocations(ctx context.Context, limit, offset int) ([]*domain.Location, int, error)
	ListLocationsByZone(ctx context.Context, zone string, limit, offset int) ([]*domain.Location, int, error)
}

type locationService struct {
	locationRepo repository.LocationRepository
}

func NewLocationService(locationRepo repository.LocationRepository) LocationService {
	return &locationService{
		locationRepo: locationRepo,
	}
}

func (s *locationService) CreateLocation(ctx context.Context, req *domain.CreateLocationRequest) (*domain.Location, error) {
	// Check if code already exists
	existingLocation, err := s.locationRepo.GetByCode(ctx, req.Code)
	if err == nil && existingLocation != nil {
		return nil, domain.ErrDuplicateEntry
	}

	location := &domain.Location{
		Code:        req.Code,
		Name:        req.Name,
		Zone:        req.Zone,
		Aisle:       req.Aisle,
		Rack:        req.Rack,
		Shelf:       req.Shelf,
		Capacity:    req.Capacity,
		Temperature: req.Temperature,
	}

	err = s.locationRepo.Create(ctx, location)
	if err != nil {
		return nil, fmt.Errorf("failed to create location: %w", err)
	}

	return location, nil
}

func (s *locationService) GetLocationByID(ctx context.Context, id int) (*domain.Location, error) {
	location, err := s.locationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get location: %w", err)
	}

	return location, nil
}

func (s *locationService) GetLocationByCode(ctx context.Context, code string) (*domain.Location, error) {
	location, err := s.locationRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get location by code: %w", err)
	}

	return location, nil
}

func (s *locationService) UpdateLocation(ctx context.Context, id int, req *domain.UpdateLocationRequest) (*domain.Location, error) {
	// Get existing location
	location, err := s.locationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get location: %w", err)
	}

	// Update fields if provided
	if req.Code != nil {
		// Check if new code already exists (excluding current location)
		existingLocation, err := s.locationRepo.GetByCode(ctx, *req.Code)
		if err == nil && existingLocation != nil && existingLocation.ID != id {
			return nil, domain.ErrDuplicateEntry
		}
		location.Code = *req.Code
	}
	if req.Name != nil {
		location.Name = *req.Name
	}
	if req.Zone != nil {
		location.Zone = *req.Zone
	}
	if req.Aisle != nil {
		location.Aisle = *req.Aisle
	}
	if req.Rack != nil {
		location.Rack = *req.Rack
	}
	if req.Shelf != nil {
		location.Shelf = *req.Shelf
	}
	if req.Capacity != nil {
		location.Capacity = *req.Capacity
	}
	if req.Temperature != nil {
		location.Temperature = req.Temperature
	}
	if req.IsActive != nil {
		location.IsActive = *req.IsActive
	}

	err = s.locationRepo.Update(ctx, location)
	if err != nil {
		return nil, fmt.Errorf("failed to update location: %w", err)
	}

	return location, nil
}

func (s *locationService) DeleteLocation(ctx context.Context, id int) error {
	err := s.locationRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete location: %w", err)
	}

	return nil
}

func (s *locationService) ListLocations(ctx context.Context, limit, offset int) ([]*domain.Location, int, error) {
	locations, total, err := s.locationRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list locations: %w", err)
	}

	return locations, total, nil
}

func (s *locationService) ListLocationsByZone(ctx context.Context, zone string, limit, offset int) ([]*domain.Location, int, error) {
	locations, total, err := s.locationRepo.ListByZone(ctx, zone, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list locations by zone: %w", err)
	}

	return locations, total, nil
}
