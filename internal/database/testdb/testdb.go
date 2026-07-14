package testdb

import (
	"api-hotelaria/internal/config"
	"api-hotelaria/internal/database"
	"log"
	"testing"

	migrate "github.com/golang-migrate/migrate/v4"
)

func SetupIntegrationTests(m *testing.M) int {
	config.LoadConfig()
	database.ConnectDB(config.Env.TestDatabaseURL)

	migration := database.MigrationInstance()

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("❌ testdb: Falha crítica na migração das tabelas de teste: %v", err)
	}
	migration.Close()

	code := m.Run()
	database.CloseDB()

	return code
}
