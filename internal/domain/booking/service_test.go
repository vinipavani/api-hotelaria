package booking

import (
	"api-hotelaria/internal/domain/room"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
)

type mockRoomRepository struct {
	onFindByID func(ctx context.Context, id int64) (*room.Room, error)
}

func (m *mockRoomRepository) FindByID(ctx context.Context, id int64) (*room.Room, error) {
	return m.onFindByID(ctx, id)
}

type mockBookingRepository struct {
	onCreate               func(ctx context.Context, b *Booking) error
	onIsBookingAvailable   func(ctx context.Context, RoomID int64) (bool, error)
	onUpdateCheckOut       func(ctx context.Context, RoomID int64, checkOutDate string) (*Booking, error)
	onGetInProgressBooking func(ctx context.Context, RoomID int64) (*Booking, error)
}

func (m *mockBookingRepository) Create(ctx context.Context, b *Booking) error {
	return m.onCreate(ctx, b)
}

func (m *mockBookingRepository) IsBookingAvailable(ctx context.Context, RoomID int64) (bool, error) {
	return m.onIsBookingAvailable(ctx, RoomID)
}

func (m *mockBookingRepository) UpdateCheckOut(ctx context.Context, RoomID int64, checkOutDate string) (*Booking, error) {
	return m.onUpdateCheckOut(ctx, RoomID, checkOutDate)
}

func (m *mockBookingRepository) GetInProgressBooking(ctx context.Context, RoomID int64) (*Booking, error) {
	return m.onGetInProgressBooking(ctx, RoomID)
}

func TestCreateCheckIn_Suite(t *testing.T) {
	ctx := context.Background()

	t.Run("should create booking successfully under ideal conditions", func(t *testing.T) {
		roomMock := &mockRoomRepository{
			onFindByID: func(ctx context.Context, id int64) (*room.Room, error) {
				return &room.Room{ID: id}, nil
			},
		}
		bookingMock := &mockBookingRepository{
			onIsBookingAvailable: func(ctx context.Context, RoomID int64) (bool, error) {
				return true, nil
			},
			onCreate: func(ctx context.Context, b *Booking) error {
				b.ID = 100
				return nil
			},
		}

		service := NewService(bookingMock, roomMock)
		input := CheckInInput{GuestName: "Vinicius Dev", GuestDocument: "123.456.789-00", CheckIn: "2026-07-13"}

		result, err := service.CreateCheckIn(ctx, 1, input)

		if err != nil {
			t.Fatalf("esperava erro nulo, recebeu: %v", err)
		}
		if result.ID != 100 {
			t.Errorf("esperava ID 100, recebeu: %d", result.ID)
		}
	})

	t.Run("should fail if guest name is empty", func(t *testing.T) {
		service := NewService(&mockBookingRepository{}, &mockRoomRepository{})
		input := CheckInInput{GuestName: "", GuestDocument: "123456", CheckIn: "2026-07-13"}

		_, err := service.CreateCheckIn(ctx, 1, input)

		if !errors.Is(err, InvalidGuestParams) {
			t.Errorf("esperava erro '%v', recebeu: '%v'", InvalidGuestParams, err)
		}
	})

	t.Run("should fail if the room is already occupied", func(t *testing.T) {
		roomMock := &mockRoomRepository{
			onFindByID: func(ctx context.Context, id int64) (*room.Room, error) {
				return &room.Room{ID: id}, nil
			},
		}
		bookingMock := &mockBookingRepository{
			onIsBookingAvailable: func(ctx context.Context, RoomID int64) (bool, error) {
				return false, nil
			},
		}

		service := NewService(bookingMock, roomMock)
		input := CheckInInput{GuestName: "Hóspede Dois", GuestDocument: "654321", CheckIn: "2026-07-13"}

		_, err := service.CreateCheckIn(ctx, 1, input)

		if !errors.Is(err, RoomNotAvailable) {
			t.Errorf("esperava erro '%v', recebeu: '%v'", RoomNotAvailable, err)
		}
	})

	t.Run("should fail if the room ID does not exist", func(t *testing.T) {
		roomMock := &mockRoomRepository{
			onFindByID: func(ctx context.Context, id int64) (*room.Room, error) {
				return nil, pgx.ErrNoRows
			},
		}

		service := NewService(&mockBookingRepository{}, roomMock)
		input := CheckInInput{GuestName: "Vinicius Dev", GuestDocument: "123.456.789-00", CheckIn: "2026-07-13"}

		_, err := service.CreateCheckIn(ctx, 999, input)

		if !errors.Is(err, RoomNotFound) {
			t.Errorf("esperava erro '%v', recebeu: '%v'", RoomNotFound, err)
		}
	})
}

func TestCheckOut_Suite(t *testing.T) {
	ctx := context.Background()

	t.Run("should update booking successfully under ideal conditions", func(t *testing.T) {
		roomMock := &mockRoomRepository{
			onFindByID: func(ctx context.Context, id int64) (*room.Room, error) {
				return &room.Room{ID: id}, nil
			},
		}
		bookingMock := &mockBookingRepository{
			onGetInProgressBooking: func(ctx context.Context, RoomID int64) (*Booking, error) {
				checkInTime, _ := time.Parse("2006-01-02", "2026-07-10")
				return &Booking{RoomID: RoomID, CheckInDate: checkInTime}, nil
			},
			onIsBookingAvailable: func(ctx context.Context, RoomID int64) (bool, error) {
				return false, nil
			},
			onUpdateCheckOut: func(ctx context.Context, RoomID int64, checkOutDate string) (*Booking, error) {
				checkInTime, _ := time.Parse("2006-01-02", "2026-07-10")
				checkOutTime, _ := time.Parse("2006-01-02", checkOutDate)

				return &Booking{
					ID:            100,
					RoomID:        RoomID,
					GuestName:     "Vinicius Dev",
					GuestDocument: "123.456.789-00",
					Status:        BookingStatusFinished,
					CheckInDate:   checkInTime,
					CheckOutDate:  &checkOutTime,
				}, nil
			},
		}

		service := NewService(bookingMock, roomMock)
		input := CheckOutInput{CheckOut: "2026-07-13"}

		result, err := service.CheckOut(ctx, 1, input)

		if err != nil {
			t.Fatalf("esperava erro nulo, recebeu: %v", err)
		}
		if result == nil {
			t.Fatalf("esperava um objeto de booking preenchido, mas recebeu nil")
		}
		if result.Status != BookingStatusFinished {
			t.Errorf("esperava status %s, recebeu: %s", BookingStatusFinished, result.Status)
		}
	})

	t.Run("should block if checkout date is prior to checkin date", func(t *testing.T) {
		roomMock := &mockRoomRepository{
			onFindByID: func(ctx context.Context, id int64) (*room.Room, error) {
				return &room.Room{ID: id}, nil
			},
		}
		bookingMock := &mockBookingRepository{
			onGetInProgressBooking: func(ctx context.Context, RoomID int64) (*Booking, error) {
				checkInTime, _ := time.Parse("2006-01-02", "2026-07-15")
				return &Booking{RoomID: RoomID, CheckInDate: checkInTime}, nil
			},
			onIsBookingAvailable: func(ctx context.Context, RoomID int64) (bool, error) {
				return false, nil
			},
		}

		service := NewService(bookingMock, roomMock)
		input := CheckOutInput{CheckOut: "2026-07-10"}

		_, err := service.CheckOut(ctx, 1, input)

		if !errors.Is(err, CheckOutLesserThanCheckIn) {
			t.Errorf("esperava erro '%v', recebeu: '%v'", CheckOutLesserThanCheckIn, err)
		}
	})

	t.Run("should fail if the room is already available", func(t *testing.T) {
		roomMock := &mockRoomRepository{
			onFindByID: func(ctx context.Context, id int64) (*room.Room, error) {
				return &room.Room{ID: id}, nil
			},
		}
		bookingMock := &mockBookingRepository{
			onIsBookingAvailable: func(ctx context.Context, RoomID int64) (bool, error) {
				return true, nil
			},
			onGetInProgressBooking: func(ctx context.Context, RoomID int64) (*Booking, error) {
				checkInTime, _ := time.Parse("2006-01-02", "2026-07-10")
				return &Booking{RoomID: RoomID, CheckInDate: checkInTime}, nil
			},
		}

		service := NewService(bookingMock, roomMock)
		input := CheckOutInput{CheckOut: "2026-07-13"}

		_, err := service.CheckOut(ctx, 1, input)

		if !errors.Is(err, RoomAvailable) {
			t.Errorf("esperava erro '%v', recebeu: '%v'", RoomAvailable, err)
		}
	})

	t.Run("should fail if the room ID does not exist", func(t *testing.T) {
		roomMock := &mockRoomRepository{
			onFindByID: func(ctx context.Context, id int64) (*room.Room, error) {
				return nil, pgx.ErrNoRows
			},
		}

		service := NewService(&mockBookingRepository{}, roomMock)
		input := CheckInInput{GuestName: "Vinicius Dev", GuestDocument: "123.456.789-00", CheckIn: "2026-07-13"}

		_, err := service.CreateCheckIn(ctx, 999, input)

		if !errors.Is(err, RoomNotFound) {
			t.Errorf("esperava erro '%v', recebeu: '%v'", RoomNotFound, err)
		}
	})
}
