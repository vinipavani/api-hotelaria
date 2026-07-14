package room_test

import (
	"api-hotelaria/internal/database"
	"api-hotelaria/internal/domain/booking"
	"api-hotelaria/internal/domain/hotel"
	"api-hotelaria/internal/domain/room"
	"context"
	"testing"
)

func TestRepository_Room_Queries_Integration(t *testing.T) {
	ctx := context.Background()

	hotelObj, err := hotel.CreateTestHotel(ctx, database.DB, "Hotel de Teste para Quartos", "Gramado")
	if err != nil {
		t.Fatalf("falha ao preparar o hotel pai via factory: %v", err)
	}

	defer func() {
		_, _ = database.DB.Exec(ctx, "DELETE FROM hotels WHERE id = $1;", hotelObj.ID)
	}()

	repo := room.NewRepository(database.DB)

	t.Run("should execute room insert and calculate room number sequence via database", func(t *testing.T) {
		rm1 := &room.Room{
			HotelID:       hotelObj.ID,
			Type:          room.RoomTypeSuite,
			Capacity:      4,
			PerNightValue: 350.00,
		}

		err := repo.Insert(ctx, rm1)
		if err != nil {
			t.Fatalf("A sua query SQL de INSERT de quartos falhou: %v", err)
		}

		if rm1.ID == 0 {
			t.Errorf("esperava que o RETURNING id do Postgres populasse o objeto, mas veio zero")
		}

		if rm1.Number != "0001" {
			t.Errorf("esperava que o número sequencial gerado fosse '0001', mas recebeu: '%s'", rm1.Number)
		}

		rm2 := &room.Room{
			HotelID:       hotelObj.ID,
			Type:          room.RoomTypeSingle,
			Capacity:      1,
			PerNightValue: 120.00,
		}
		err = repo.Insert(ctx, rm2)
		if err != nil {
			t.Fatalf("A sua query SQL de INSERT do segundo quarto falhou: %v", err)
		}

		if rm2.Number != "0002" {
			t.Errorf("esperava que o segundo quarto ganhasse o número '0002', mas recebeu: '%s'", rm2.Number)
		}
	})

	t.Run("should list all rooms created for the target hotel id", func(t *testing.T) {
		rooms, err := repo.FindAll(ctx, hotelObj.ID, false)
		if err != nil {
			t.Fatalf("A sua query SQL de FindAll quebrou ao listar quartos: %v", err)
		}

		if len(rooms) != 2 {
			t.Errorf("esperava encontrar 2 quartos associados ao hotel %d, mas retornaram %d", hotelObj.ID, len(rooms))
		}

		if rooms[0].Number != "0001" || rooms[1].Number != "0002" {
			t.Errorf("a ordenação dos quartos por número veio incorreta ou desalinhada")
		}
	})

	t.Run("should list only available rooms when availableOnly filter is active", func(t *testing.T) {
		hotelIsolado, err := hotel.CreateTestHotel(ctx, database.DB, "Hotel Filtro Disponibilidade", "Gramado")
		if err != nil {
			t.Fatalf("falha ao criar hotel isolado: %v", err)
		}
		defer func() {
			_, _ = database.DB.Exec(ctx, "DELETE FROM hotels WHERE id = $1;", hotelIsolado.ID)
		}()

		rmOcupado, err := room.CreateTestRoom(ctx, database.DB, hotelIsolado.ID, "0101", room.RoomTypeSuite)
		if err != nil {
			t.Fatalf("falha ao criar quarto ocupado: %v", err)
		}

		rmVago, err := room.CreateTestRoom(ctx, database.DB, hotelIsolado.ID, "0102", room.RoomTypeSingle)
		if err != nil {
			t.Fatalf("falha ao criar quarto vago: %v", err)
		}

		_, err = booking.CreateTestBooking(ctx, database.DB, rmOcupado.ID, "Hóspede Filtro", "999888", booking.BookingStatusInProgress)
		if err != nil {
			t.Fatalf("falha ao criar reserva de teste para ocupar o quarto: %v", err)
		}

		availableRooms, err := repo.FindAll(ctx, hotelIsolado.ID, true)
		if err != nil {
			t.Fatalf("A sua query SQL de FindAll com filtro de disponibilidade quebrou: %v", err)
		}

		if len(availableRooms) != 1 {
			t.Errorf("esperava encontrar exatamente 1 quarto disponível para o hotel %d, mas retornaram %d", hotelIsolado.ID, len(availableRooms))
		}

		if availableRooms[0].ID != rmVago.ID {
			t.Errorf("erro no filtro: a query trouxe o quarto de ID %d (ocupado), mas deveria trazer o ID %d (vago)", availableRooms[0].ID, rmVago.ID)
		}
	})

}
