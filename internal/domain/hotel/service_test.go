package hotel

import (
	"context"
	"testing"
)

type mockHotelRepository struct {
	onFindAll func(ctx context.Context) ([]*Hotel, error)
	onInsert  func(ctx context.Context, h *Hotel) error
}

func (m *mockHotelRepository) FindAll(ctx context.Context) ([]*Hotel, error) {
	return m.onFindAll(ctx)
}

func (m *mockHotelRepository) Insert(ctx context.Context, h *Hotel) error {
	return m.onInsert(ctx, h)
}

func TestCreateHotel_Suite(t *testing.T) {
	ctx := context.Background()

	t.Run("should create hotel successfully under ideal conditions", func(t *testing.T) {
		hotelRepoMock := &mockHotelRepository{
			onInsert: func(ctx context.Context, h *Hotel) error {
				return nil
			},
		}

		service := NewService(hotelRepoMock)
		input := CreateHotelInput{Name: "Copacabana Palace", City: "Rio de Janeiro"}
		_, err := service.CreateHotel(ctx, input)

		if err != nil {
			t.Fatalf("esperava erro nulo, recebeu: %v", err)
		}
	})

	t.Run("should return InvalidParams error if name or city weren't provided", func(t *testing.T) {
		hotelRepoMock := &mockHotelRepository{
			onInsert: func(ctx context.Context, h *Hotel) error {
				return nil
			},
		}
		service := NewService(hotelRepoMock)

		input := CreateHotelInput{Name: "", City: "Rio de Janeiro"}
		_, err := service.CreateHotel(ctx, input)
		if err != InvalidParams {
			t.Errorf("esperava erro '%v', recebeu: '%v'", InvalidParams, err)
		}
	})
}
