package database

import (
	"context"
	"log"
	"time"
	"github.com/jackc/pgx/v5/pgxpool"
	migrate "github.com/golang-migrate/migrate/v4"
	source "github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	pgxmigrate "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	migrateDB "github.com/golang-migrate/migrate/v4/database"
	"github.com/jackc/pgx/v5/stdlib"
	"api-hotelaria/internal/config"
	"api-hotelaria/internal/database/migrations"
)

var DB *pgxpool.Pool

func ConnectDB() {
	dbConfig, err := pgxpool.ParseConfig(config.Env.DatabaseURL)
	if err != nil {
		log.Fatalf("Erro ao processar a string de conexão do banco: %v\n", err)
	}

	dbConfig = configDatabaseConnection(dbConfig)

	DB, err = pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Fatalf("Erro ao criar o pool de conexões do Postgres: %v\n", err)
	}

	err = DB.Ping(context.Background())
	if err != nil {
		log.Fatalf("O contêiner do Postgres respondeu, mas recusou a conexão: %v\n", err)
	}

	log.Println("Conexão com o PostgreSQL do Docker estabelecida com sucesso!")

	runMigrations()
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Pool de conexões com o banco de dados fechado com sucesso.")
	}
}

func configDatabaseConnection(dbConfig *pgxpool.Config) *pgxpool.Config {
	dbConfig.MaxConns = 10
	dbConfig.MinConns = 2
	dbConfig.MaxConnIdleTime = time.Minute * 5 // Fecha conexões inativas após 5 minutos
	
	return dbConfig
}

func runMigrations() {
	migration := migrationInstance()

	err := migration.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Erro crítico ao executar as migrações do banco: %v\n", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("✅ Banco de dados já está atualizado. Nenhuma migração pendente.")
	} else {
		log.Println("⚡ Migrações de banco de dados executadas com sucesso!")
	}
}

func migrationInstance() *migrate.Migrate {
	driver := createMigrationDriver()
	sourceDriver := getMemorizedMigrations()

	migration, err := migrate.NewWithInstance("iofs", sourceDriver, "pgx", driver)
	if err != nil {
		log.Fatalf("Erro ao inicializar o motor de migrações: %v\n", err)
	}

	return migration
}

func createMigrationDriver() migrateDB.Driver {
	dbInstance := stdlib.OpenDBFromPool(DB)
	driver, err := pgxmigrate.WithInstance(dbInstance, &pgxmigrate.Config{})
	if err != nil {
		log.Fatalf("Erro ao criar driver de migração: %v\n", err)
	}

	return driver
}

func getMemorizedMigrations() source.Driver {
	sourceDriver, err := iofs.New(migrations.FS, ".")
	if err != nil {
		log.Fatalf("Erro ao ler arquivos SQL embutidos: %v\n", err)
	}

	return sourceDriver
}
