package booking

import (
	"api-hotelaria/internal/database"
	"api-hotelaria/internal/domain/hotel"
	"api-hotelaria/internal/domain/room"
	"context"
	"testing"
	"time"
)

func TestBookingRepository_Queries_Integration(t *testing.T) {
	ctx := context.Background()

	hotelObj, err := hotel.CreateTestHotel(ctx, database.DB, "Copacabana Palace Test", "Rio de Janeiro")
	if err != nil {
		t.Fatalf("falha ao preparar o hotel pai via factory: %v", err)
	}

	roomObj, err := room.CreateTestRoom(ctx, database.DB, hotelObj.ID, "0999", room.RoomTypeSuite)
	if err != nil {
		t.Fatalf("falha ao preparar o quarto pai via factory: %v", err)
	}

	defer func() {
		_, _ = database.DB.Exec(ctx, "DELETE FROM hotels WHERE id = $1;", hotelObj.ID)
	}()

	repo := NewRepository(database.DB)

	t.Run("should persist check-in booking query and return generated id and times", func(t *testing.T) {
		b := &Booking{
			RoomID:        roomObj.ID,
			GuestName:     "Vinicius Integração",
			GuestDocument: "123.456.789-00",
			CheckInDate:   time.Now().Truncate(time.Second),
			Status:        BookingStatusInProgress,
		}

		err := repo.Create(ctx, b)
		if err != nil {
			t.Fatalf("A sua query SQL de INSERT (Create) quebrou no Postgres de testes: %v", err)
		}

		if b.ID == 0 {
			t.Errorf("esperava que o RETURNING id do banco populasse a struct, mas veio zero")
		}

		if b.CreatedAt.IsZero() {
			t.Errorf("esperava que o DEFAULT NOW() da tabela populasse o CreatedAt, mas veio zerado")
		}

		if b.Status != BookingStatusInProgress {
			t.Errorf("esperava que a estadia fosse criada com status 'em_estadia'")
		}
	})

	t.Run("should verify that room is occupied when an active booking exists", func(t *testing.T) {
		isAvailable, err := repo.IsBookingAvailable(ctx, roomObj.ID)
		if err != nil {
			t.Fatalf("A query SQL de IsBookingAvailable falhou: %v", err)
		}

		if isAvailable {
			t.Errorf("esperava isAvailable = false porque há um hóspede no quarto, mas retornou true")
		}
	})

	t.Run("should retrieve the active booking for the target room id", func(t *testing.T) {
		activeBooking, err := repo.GetInProgressBooking(ctx, roomObj.ID)
		if err != nil {
			t.Fatalf("A query SQL de GetInProgressBooking falhou: %v", err)
		}

		if activeBooking == nil {
			t.Fatalf("esperava encontrar uma reserva ativa, mas retornou nil")
		}

		if activeBooking.GuestName != "Vinicius Integração" {
			t.Errorf("esperava o nome 'Vinicius Integração', mas recebeu: '%s'", activeBooking.GuestName)
		}
	})

	t.Run("should execute update query to finish booking and set status to finished", func(t *testing.T) {
		dataFuturaCheckOutStr := time.Now().Add(24 * time.Hour * 3).Format("2006-01-02")

		updatedBooking, err := repo.UpdateCheckOut(ctx, roomObj.ID, dataFuturaCheckOutStr)
		if err != nil {
			t.Fatalf("A sua query SQL de UPDATE (UpdateCheckOut) quebrou: %v", err)
		}

		if updatedBooking.Status != BookingStatusFinished {
			t.Errorf("esperava status alterado para '%s', mas recebeu: '%s'", BookingStatusFinished, updatedBooking.Status)
		}

		if updatedBooking.CheckOutDate == nil {
			t.Fatalf("esperava que a coluna check_out_date viesse preenchida, mas recebeu nil")
		}
	})

	t.Run("should verify that room is free again after checkout is performed", func(t *testing.T) {
		isAvailable, err := repo.IsBookingAvailable(ctx, roomObj.ID)
		if err != nil {
			t.Fatalf("A query SQL de IsBookingAvailable falhou pós checkout: %v", err)
		}

		if !isAvailable {
			t.Errorf("esperava isAvailable = true porque a estadia foi encerrada, mas retornou false")
		}
	})
}
