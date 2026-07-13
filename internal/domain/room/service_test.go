package room

import (
	"api-hotelaria/internal/domain/hotel"
	"context"
	"errors"
	"testing"
)

type mockRoomRepository struct {
	onFindAll func(ctx context.Context, hotelID int64, availableOnly bool) ([]*Room, error)
	onInsert  func(ctx context.Context, rm *Room) error
}

func (m *mockRoomRepository) FindAll(ctx context.Context, hotelID int64, availableOnly bool) ([]*Room, error) {
	return m.onFindAll(ctx, hotelID, availableOnly)
}

func (m *mockRoomRepository) Insert(ctx context.Context, rm *Room) error {
	return m.onInsert(ctx, rm)
}

type mockHotelRepository struct {
	onFindByID func(ctx context.Context, hotelID int64) (*hotel.Hotel, error)
}

func (m *mockHotelRepository) FindByID(ctx context.Context, hotelID int64) (*hotel.Hotel, error) {
	return m.onFindByID(ctx, hotelID)
}

func TestCreateRoom_Suite(t *testing.T) {
	ctx := context.Background()

	t.Run("should create room successfully under ideal conditions", func(t *testing.T) {
		hotelMock := &mockHotelRepository{
			onFindByID: func(ctx context.Context, hotelID int64) (*hotel.Hotel, error) {
				return &hotel.Hotel{ID: hotelID, Name: "Copacabana Palace"}, nil
			},
		}
		roomMock := &mockRoomRepository{
			onInsert: func(ctx context.Context, rm *Room) error {
				rm.ID = 10
				return nil
			},
		}

		service := NewService(roomMock, hotelMock)
		input := CreateRoomInput{
			HotelID:       1,
			Type:          RoomTypeSuite,
			Capacity:      4,
			PerNightValue: 450.00,
		}

		result, err := service.CreateRoom(ctx, input)

		if err != nil {
			t.Fatalf("esperava erro nulo, recebeu: %v", err)
		}
		if result.ID != 10 {
			t.Errorf("esperava ID do quarto como 10, recebeu: %d", result.ID)
		}
		if result.Type != RoomTypeSuite {
			t.Errorf("esperava tipo suite, recebeu: %s", result.Type)
		}
	})

	t.Run("should fail if the target hotel ID does not exist", func(t *testing.T) {
		hotelMock := &mockHotelRepository{
			onFindByID: func(ctx context.Context, hotelID int64) (*hotel.Hotel, error) {
				return nil, nil
			},
		}

		service := NewService(&mockRoomRepository{}, hotelMock)
		input := CreateRoomInput{HotelID: 999, Type: RoomTypeSingle, Capacity: 1, PerNightValue: 100.00}

		_, err := service.CreateRoom(ctx, input)
		if !errors.Is(err, HotelNotFound) {
			t.Errorf("esperava erro '%v', recebeu: '%v'", HotelNotFound, err)
		}
	})

	t.Run("should fail if capacity or per night value are lesser or equal to zero", func(t *testing.T) {
		hotelMock := &mockHotelRepository{
			onFindByID: func(ctx context.Context, hotelID int64) (*hotel.Hotel, error) {
				return &hotel.Hotel{ID: hotelID}, nil
			},
		}

		service := NewService(&mockRoomRepository{}, hotelMock)
		input := CreateRoomInput{HotelID: 1, Type: RoomTypeSuite, Capacity: -2, PerNightValue: 100.00}

		_, err := service.CreateRoom(ctx, input)
		if !errors.Is(err, InvalidParams) {
			t.Errorf("esperava erro '%v', recebeu: '%v'", InvalidParams, err)
		}
	})

	t.Run("should fail if room type is invalid string value", func(t *testing.T) {
		hotelMock := &mockHotelRepository{
			onFindByID: func(ctx context.Context, hotelID int64) (*hotel.Hotel, error) {
				return &hotel.Hotel{ID: hotelID}, nil
			},
		}

		service := NewService(&mockRoomRepository{}, hotelMock)
		input := CreateRoomInput{HotelID: 1, Type: RoomType("presidencial"), Capacity: 4, PerNightValue: 900.00}
		_, err := service.CreateRoom(ctx, input)

		if !errors.Is(err, InvalidRoomType) {
			t.Errorf("esperava erro '%v', recebeu: '%v'", InvalidRoomType, err)
		}
	})
}

func TestFindAllRooms_Suite(t *testing.T) {
	ctx := context.Background()

	t.Run("should return the list of rooms for an existing hotel", func(t *testing.T) {
		hotelMock := &mockHotelRepository{
			onFindByID: func(ctx context.Context, hotelID int64) (*hotel.Hotel, error) {
				return &hotel.Hotel{ID: hotelID}, nil
			},
		}

		roomMock := &mockRoomRepository{
			onFindAll: func(ctx context.Context, hotelID int64, availableOnly bool) ([]*Room, error) {
				return []*Room{
					{ID: 1, HotelID: hotelID, Capacity: 2},
					{ID: 2, HotelID: hotelID, Capacity: 4},
				}, nil
			},
		}
		service := NewService(roomMock, hotelMock)

		result, err := service.FindAllRooms(ctx, 1, false)
		if err != nil {
			t.Fatalf("esperava erro nulo, recebeu: %v", err)
		}
		if len(result) != 2 {
			t.Errorf("esperava 2 quartos na lista, recebeu: %d", len(result))
		}
	})
}
