package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
)

type seedHotel struct {
	Name string
	City string
}

type seedRoom struct {
	Number        string
	RoomType      string
	Capacity      int
	PerNightValue float64
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	log.Println("🌱 Conectando ao banco de dados para rodar a Seed em massa...")
	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("❌ Erro ao conectar no banco: %v", err)
	}
	defer db.Close()

	log.Println("🧹 Limpando dados históricos do banco...")
	_, err = db.Exec(ctx, "TRUNCATE TABLE hotels, rooms, bookings RESTART IDENTITY CASCADE;")
	if err != nil {
		log.Fatalf("❌ Erro ao limpar tabelas: %v", err)
	}

	hoteis := []seedHotel{
		{Name: "Grand Plaza Hotel", City: "Nova York"},
		{Name: "Copacabana Sea View", City: "Rio de Janeiro"},
		{Name: "Ipanema Palace", City: "Rio de Janeiro"},
		{Name: "Oceanic Resort", City: "Salvador"},
		{Name: "Mountain Retreat", City: "Gramado"},
	}

	log.Println("🚀 Iniciando a inserção em lote (5 Hotéis x 20 Quartos = 100 Quartos)...")

	for _, h := range hoteis {
		var hotelID int64
		hotelQuery := "INSERT INTO hotels (name, city) VALUES ($1, $2) RETURNING id;"

		err = db.QueryRow(ctx, hotelQuery, h.Name, h.City).Scan(&hotelID)
		if err != nil {
			log.Fatalf("❌ Erro ao criar o hotel %s: %v", h.Name, err)
		}

		log.Printf("🏨 Hotel [%s] criado com ID: %d. Gerando os 20 quartos...\n", h.Name, hotelID)

		for q := 1; q <= 20; q++ {
			room := newSeedRoomFactory(q)

			roomQuery := `
				INSERT INTO rooms (hotel_id, number, type, capacity, per_night_value) 
				VALUES ($1, $2, $3, $4, $5);
			`
			_, err = db.Exec(ctx, roomQuery, hotelID, room.Number, room.RoomType, room.Capacity, room.PerNightValue)
			if err != nil {
				log.Fatalf("❌ Erro ao criar quarto %d para o hotel %d: %v", q, hotelID, err)
			}
		}
	}

	log.Println("🚀 =====================================================")
	log.Println("✅ SEED EXECUTADA COM SUCESSO!")
	log.Println("📊 Total de Hotéis populados: 5")
	log.Println("📊 Total de Quartos injetados: 100 (20 por hotel)")
	log.Println("🚀 =====================================================")
}

func newSeedRoomFactory(index int) seedRoom {
	roomTypes := []string{"single", "double", "suite"}
	chosenType := roomTypes[index%3]

	capacity := 1
	value := 120.00

	switch chosenType {
	case "double":
		capacity = 2
		value = 220.00
	case "suite":
		capacity = 4
		value = 450.00
	}

	return seedRoom{
		Number:        fmt.Sprintf("%04d", index),
		RoomType:      chosenType,
		Capacity:      capacity,
		PerNightValue: value,
	}
}
