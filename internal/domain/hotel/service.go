package hotel

import (
	"context"
	"errors"
	"strings"
)

type HotelRepository interface {
	FindAll(ctx context.Context) ([]*Hotel, error)
	Insert(ctx context.Context, h *Hotel) error
}

type Service struct {
	repo HotelRepository
}

func NewService(repository HotelRepository) *Service {
	return &Service{
		repo: repository,
	}
}

func (s *Service) findAllHotels(ctx context.Context) ([]*Hotel, error) {
	return s.repo.FindAll(ctx)
}

func (s *Service) CreateHotel(ctx context.Context, input CreateHotelInput) (*Hotel, error) {
	newHotel, err := validateParams(input.Name, input.City)
	if err != nil {
		return nil, err
	}

	err = s.repo.Insert(ctx, newHotel)
	if err != nil {
		return nil, err
	}

	return newHotel, nil
}

func validateParams(name string, city string) (*Hotel, error) {
	nameClean := strings.TrimSpace(name)
	cityClean := strings.TrimSpace(city)

	if nameClean == "" || cityClean == "" {
		return nil, errors.New("o nome e a cidade do hotel não podem estar em branco.")
	}

	validHotel := &Hotel{
		Name: nameClean,
		City: cityClean,
	}

	return validHotel, nil
}
