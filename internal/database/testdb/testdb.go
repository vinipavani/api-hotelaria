package testdb

import (
	"api-hotelaria/internal/database"
	"log"
	"testing"

	migrate "github.com/golang-migrate/migrate/v4"
)

func SetupIntegrationTests(m *testing.M) int {
	database.ConnectDB()

	migration := database.MigrationInstance()

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("❌ testdb: Falha crítica na migração das tabelas de teste: %v", err)
	}
	migration.Close()

	code := m.Run()
	database.CloseDB()

	return code
}
