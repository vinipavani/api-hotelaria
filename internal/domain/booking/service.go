package booking

import (
	"api-hotelaria/internal/domain/room"
	"context"
	"errors"
)

var (
	ErrRoomNotFound  = errors.New("Quarto não encontrado.")
	RoomNotAvailable = errors.New("Quarto já está ocupado.")
	GuestParamsError = errors.New("Os dados do cliente nome e número do documento são obrigatórios.")
	CheckInDateError = errors.New("A Data de checkin é obrigatória.")
)

type Service struct {
	repo     *Repository
	roomRepo *room.Repository
}

func NewService(repository *Repository, roomRepository *room.Repository) *Service {
	return &Service{
		repo:     repository,
		roomRepo: roomRepository,
	}
}

func (s *Service) CreateCheckIn(ctx context.Context, RoomID int64, input BookingInput) (*Booking, error) {
	newBooking, err := validateParams(RoomID, input)
	if err != nil {
		return nil, err
	}

	err = validateRoom(RoomID, s.roomRepo, ctx)
	if err != nil {
		return nil, err
	}

	err = s.repo.Create(ctx, newBooking)
	if err != nil {
		return nil, err
	}

	return newBooking, nil
}

func validateParams(RoomID int64, input BookingInput) (*Booking, error) {
	if input.GuestName == "" || input.GuestDocument == "" {
		return nil, GuestParamsError
	}
	if input.CheckInDate.IsZero() {
		return nil, CheckInDateError
	}

	return &Booking{
		RoomID:        RoomID,
		GuestName:     input.GuestName,
		GuestDocument: input.GuestDocument,
		CheckInDate:   input.CheckInDate,
		Status:        BookingStatusInProgress,
	}, nil
}

func validateRoom(RoomID int64, roomRepo *room.Repository, ctx context.Context) error {
	room, err := roomRepo.FindByID(ctx, RoomID)
	if room == nil || err != nil {
		return ErrRoomNotFound
	}

	isAvailable, err := s.repo.isBookingAvailable(ctx, RoomID)
	if !isAvailable {
		return RoomNotAvailable
	} else if err != nil {
		return err
	}

	return nil
}
