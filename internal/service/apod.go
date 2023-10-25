// Package service represents entity to manipulate with business logics of data
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/artnikel/apodservice/internal/model"
)

// ApodRepository is an interface that contains methods for apod manipulation
type ApodRepository interface {
	ApodGetAll(ctx context.Context) ([]*model.APOD, error)
	ApodGetByDate(ctx context.Context, date time.Time) (*model.APOD, error)
}

// ApodService implements methods of ApodRepository
type ApodService struct {
	apodRepo ApodRepository
}

// NewApodService accepts ApodRepository and returns an object of type *ApodService
func NewApodService(apodRepo ApodRepository) *ApodService {
	return &ApodService{apodRepo: apodRepo}
}

// GetAll is a method of ApodService that calls method of Repository
func (s *ApodService) GetAll(ctx context.Context) ([]*model.APOD, error) {
	apods, err := s.apodRepo.ApodGetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("getAll %w", err)
	}
	return apods, nil
}

// GetByDate is a method of ApodService that calls method of Repository
func (s *ApodService) GetByDate(ctx context.Context, date time.Time) (*model.APOD, error) {
	apod, err := s.apodRepo.ApodGetByDate(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("getByDate %w", err)
	}
	return apod, nil
}
