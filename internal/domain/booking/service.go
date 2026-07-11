package booking

import (
	"api-hotelaria/internal/domain/room"
	"context"
	"errors"
	"time"
)

var (
	ErrRoomNotFound  = errors.New("Quarto não encontrado.")
	RoomNotAvailable = errors.New("Quarto já está ocupado.")
	GuestParamsError = errors.New("Os dados do cliente nome e número do documento são obrigatórios.")
	CheckInDateError = errors.New("A Data de checkin deve ser uma data valida no formato: 'YYYY-MM-DD'.")
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

	if err = validateRoom(RoomID, s.roomRepo, ctx); err != nil {
		return nil, err
	}

	if err := validateBooking(RoomID, s.repo, ctx); err != nil {
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

	checkInTime, err := time.Parse("2006-01-02", input.CheckIn)
	if err != nil {
		return nil, CheckInDateError
	}

	return &Booking{
		RoomID:        RoomID,
		GuestName:     input.GuestName,
		GuestDocument: input.GuestDocument,
		CheckInDate:   checkInTime,
		Status:        BookingStatusInProgress,
	}, nil
}

func validateRoom(RoomID int64, roomRepo *room.Repository, ctx context.Context) error {
	room, err := roomRepo.FindByID(ctx, RoomID)
	if room == nil || err != nil {
		return ErrRoomNotFound
	}

	return nil
}

func validateBooking(RoomID int64, r *Repository, ctx context.Context) error {
	isAvailable, err := r.isBookingAvailable(ctx, RoomID)
	if !isAvailable {
		return RoomNotAvailable
	} else if err != nil {
		return err
	}

	return nil
}
