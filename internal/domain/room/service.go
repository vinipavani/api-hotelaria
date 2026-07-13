package room

import (
	"api-hotelaria/internal/domain/hotel"
	"context"
	"errors"
)

type RoomRepository interface {
	FindAll(ctx context.Context, HotelID int64, availableOnly bool) ([]*Room, error)
	Insert(ctx context.Context, rm *Room) error
}

type HotelRepository interface {
	FindByID(ctx context.Context, HotelID int64) (*hotel.Hotel, error)
}

type Service struct {
	repo      RoomRepository
	hotelRepo HotelRepository
}

var (
	InvalidParams   = errors.New("Tipo de quarto, capacidade e diária, não pode estar em branco.")
	InvalidRoomType = errors.New("Tipo de quarto inválido. Os tipos válidos são: single, double, suite.")
	HotelNotFound   = errors.New("Hotel não encontrado.")
)

func NewService(repository RoomRepository, hotelRepository HotelRepository) *Service {
	return &Service{
		repo:      repository,
		hotelRepo: hotelRepository,
	}
}

func (s *Service) FindAllRooms(ctx context.Context, hotelID int64, availableOnly bool) ([]*Room, error) {
	err := validateHotel(s.hotelRepo, ctx, hotelID)
	if err != nil {
		return nil, err
	}

	rooms, err := s.repo.FindAll(ctx, hotelID, availableOnly)
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (s *Service) CreateRoom(ctx context.Context, input CreateRoomInput) (*Room, error) {
	err := validateHotel(s.hotelRepo, ctx, input.HotelID)
	if err != nil {
		return nil, err
	}

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
		return nil, InvalidParams
	}

	switch rt {
	case RoomTypeSingle, RoomTypeDouble, RoomTypeSuite:
	default:
		return nil, InvalidRoomType
	}

	validRoom := &Room{
		HotelID:       hotelID,
		Type:          rt,
		Capacity:      capacity,
		PerNightValue: perNightValue,
	}

	return validRoom, nil
}

func validateHotel(hr HotelRepository, ctx context.Context, hotelID int64) error {
	hotel, err := hr.FindByID(ctx, hotelID)
	if hotel == nil {
		return HotelNotFound
	} else if err != nil {
		return err
	}

	return nil
}
