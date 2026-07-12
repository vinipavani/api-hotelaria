package room

import (
	"context"
	"errors"
)

type Service struct {
	repo *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repo: repository,
	}
}

func (s *Service) findAllRooms(ctx context.Context, hotelID string, availableOnly bool) ([]*Room, error) {
	rooms, err := s.repo.FindAll(ctx, hotelID, availableOnly)
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (s *Service) CreateRoom(ctx context.Context, input CreateRoomInput) (*Room, error) {
	newRoom, err := validateParams(input.HotelID, input.Type, input.Capacity, input.PerNightValue)
	if err != nil {
		return nil, err
	}

	err = s.repo.Insert(ctx, newRoom)
	if err != nil {
		return nil, err
	}

	return newRoom, nil
}

func validateParams(hotelID int64, rt RoomType, capacity int, perNightValue float64) (*Room, error) {
	if string(rt) == "" || capacity <= 0 || perNightValue <= 0 {
		return nil, errors.New("Tipo de quarto, capacidade e diária, não pode estar em branco.")
	}

	switch rt {
	case RoomTypeSingle, RoomTypeDouble, RoomTypeSuite:
	default:
		return nil, errors.New("Tipo de quarto inválido. Os tipos válidos são: single, double, suite.")
	}

	validRoom := &Room{
		HotelID:       hotelID,
		Type:          rt,
		Capacity:      capacity,
		PerNightValue: perNightValue,
	}

	return validRoom, nil
}
