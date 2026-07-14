package testdb

import (
	"api-hotelaria/internal/config"
	"api-hotelaria/internal/database"
	"log"
	"testing"

	migrate "github.com/golang-migrate/migrate/v4"
)

func SetupIntegrationTests(m *testing.M) int {
	config.LoadEnv()
	database.ConnectDB()

	migration := database.MigrationInstance()
	err := migration.Up()

	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("❌ testdb: Falha crítica ao erguer tabelas de teste: %v", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("✨ testdb: Tabelas de teste já estruturadas. Reutilizando esquema de forma isolada.")
	} else {
		log.Println("🚀 testdb: Estrutura inicial das tabelas erguida com sucesso no banco de testes!")
	}

	migration.Close()
	code := m.Run()

	database.CloseDB()
	return code
}
