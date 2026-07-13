package booking

import (
	"api-hotelaria/internal/domain/room"
	"context"
	"errors"
	"time"
)

type BookingRepository interface {
	Create(ctx context.Context, b *Booking) error
	IsBookingAvailable(ctx context.Context, RoomID int64) (bool, error)
	UpdateCheckOut(ctx context.Context, RoomID int64, checkOutDate string) (*Booking, error)
	GetInProgressBooking(ctx context.Context, RoomID int64) (*Booking, error)
}

type RoomRepository interface {
	FindByID(ctx context.Context, id int64) (*room.Room, error)
}

type Service struct {
	repo     BookingRepository
	roomRepo RoomRepository
}

var (
	RoomNotFound              = errors.New("Quarto não encontrado.")
	RoomNotAvailable          = errors.New("Quarto já está ocupado.")
	InvalidGuestParams        = errors.New("Os dados do cliente nome e número do documento são obrigatórios.")
	InvalidCheckInDateFormat  = errors.New("A data de check-in deve ser uma data valida no formato: 'YYYY-MM-DD'.")
	InvalidCheckOutDateFormat = errors.New("A data de check-out deve ser uma data valida no formato: 'YYYY-MM-DD'.")
	CheckOutLesserThanCheckIn = errors.New("A data do check-out não pode ser anterior a do check-in.")
	RoomAvailable             = errors.New("O quarto não possui estadia ativa.")
)

func NewService(repository BookingRepository, roomRepository RoomRepository) *Service {
	return &Service{
		repo:     repository,
		roomRepo: roomRepository,
	}
}

func (s *Service) CreateCheckIn(ctx context.Context, RoomID int64, input CheckInInput) (*Booking, error) {
	newBooking, err := validateParams(RoomID, input)
	if err != nil {
		return nil, err
	}

	if err = validateRoom(RoomID, s.roomRepo, ctx); err != nil {
		return nil, err
	}

	if isAvailable, err := isBookingAvailable(RoomID, s.repo, ctx); err != nil {
		return nil, err
	} else if !isAvailable {
		return nil, RoomNotAvailable
	}

	err = s.repo.Create(ctx, newBooking)
	if err != nil {
		return nil, err
	}

	return newBooking, nil
}

func (s *Service) CheckOut(ctx context.Context, RoomID int64, input CheckOutInput) (*Booking, error) {
	var booking *Booking
	if err := validateRoom(RoomID, s.roomRepo, ctx); err != nil {
		return nil, err
	}

	checkOutTime, err := time.Parse("2006-01-02", input.CheckOut)
	if err != nil {
		return nil, InvalidCheckOutDateFormat
	}

	if err := validateCheckOutDate(checkOutTime, RoomID, s.repo, ctx); err != nil {
		return nil, err
	}

	if isAvailable, err := isBookingAvailable(RoomID, s.repo, ctx); err != nil {
		return nil, err
	} else if isAvailable {
		return nil, RoomAvailable
	}

	booking, err = s.repo.UpdateCheckOut(ctx, RoomID, input.CheckOut)

	return booking, nil
}

func validateParams(RoomID int64, input CheckInInput) (*Booking, error) {
	if input.GuestName == "" || input.GuestDocument == "" {
		return nil, InvalidGuestParams
	}

	checkInTime, err := time.Parse("2006-01-02", input.CheckIn)
	if err != nil {
		return nil, InvalidCheckInDateFormat
	}

	return &Booking{
		RoomID:        RoomID,
		GuestName:     input.GuestName,
		GuestDocument: input.GuestDocument,
		CheckInDate:   checkInTime,
		Status:        BookingStatusInProgress,
	}, nil
}

func validateRoom(RoomID int64, r RoomRepository, ctx context.Context) error {
	room, err := r.FindByID(ctx, RoomID)
	if room == nil || err != nil {
		return RoomNotFound
	}

	return nil
}

func isBookingAvailable(RoomID int64, r BookingRepository, ctx context.Context) (bool, error) {
	isAvailable, err := r.IsBookingAvailable(ctx, RoomID)
	if err != nil {
		return false, err
	}

	return isAvailable, nil
}

func validateCheckOutDate(checkOut time.Time, RoomID int64, r BookingRepository, ctx context.Context) error {
	b, err := r.GetInProgressBooking(ctx, RoomID)
	if err != nil {
		return err
	}

	if checkOut.Before(b.CheckInDate) {
		return CheckOutLesserThanCheckIn
	}

	return nil
}
