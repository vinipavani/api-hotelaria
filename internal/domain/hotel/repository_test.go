package hotel

import (
	"api-hotelaria/internal/database"
	"context"
	"testing"
)

func TestRepository_Hotel_Queries_Integration(t *testing.T) {
	ctx := context.Background()
	repo := NewRepository(database.DB)

	hotelName1 := "Hotel Imperial Teste"
	hotelName2 := "Hotel Fazenda Teste"

	defer func() {
		_, _ = database.DB.Exec(ctx, "DELETE FROM hotels WHERE name IN ($1, $2);", hotelName1, hotelName2)
	}()

	t.Run("should persist a new hotel query and populate generated fields", func(t *testing.T) {
		h := &Hotel{
			Name: hotelName1,
			City: "Rio de Janeiro",
		}

		err := repo.Insert(ctx, h)
		if err != nil {
			t.Fatalf("A sua query SQL de INSERT quebrou no Postgres de testes: %v", err)
		}

		if h.ID == 0 {
			t.Errorf("esperava que o RETURNING id do banco populasse a struct, mas veio zero")
		}
	})

	t.Run("should retrieve an existing hotel by its physical id", func(t *testing.T) {
		h2, err := CreateTestHotel(ctx, database.DB, hotelName2, "Socorro")
		if err != nil {
			t.Fatalf("falha ao criar hotel auxiliar via factory: %v", err)
		}

		fetchedHotel, err := repo.FindByID(ctx, h2.ID)
		if err != nil {
			t.Fatalf("A sua query SQL de FindByID falhou: %v", err)
		}

		if fetchedHotel == nil {
			t.Fatalf("esperava encontrar o hotel de ID %d, mas retornou nil", h2.ID)
		}

		if fetchedHotel.Name != hotelName2 {
			t.Errorf("esperava o nome '%s', mas recebeu: '%s'", hotelName2, fetchedHotel.Name)
		}
	})

	t.Run("should list all hotels persisted in the database table", func(t *testing.T) {
		hotels, err := repo.FindAll(ctx)
		if err != nil {
			t.Fatalf("A sua query SQL de listagem geral quebrou: %v", err)
		}

		if len(hotels) < 2 {
			t.Errorf("esperava encontrar pelo menos 2 hotéis na tabela, mas retornaram %d", len(hotels))
		}
	})
}
