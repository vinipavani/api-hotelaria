package database

import (
	"api-hotelaria/internal/config"
	"api-hotelaria/internal/database/migrations"
	"context"
	"log"
	"time"

	migrate "github.com/golang-migrate/migrate/v4"
	migrateDB "github.com/golang-migrate/migrate/v4/database"
	pgxmigrate "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	source "github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
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

	checkMigrationsStatus()
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

func MigrationInstance() *migrate.Migrate {
	driver := createMigrationDriver()
	sourceDriver := getMemorizedMigrations()

	migration, err := migrate.NewWithInstance("iofs", sourceDriver, "pgx", driver)
	if err != nil {
		log.Fatalf("Erro ao inicializar o motor de migrações: %v\n", err)
	}

	return migration
}

func checkMigrationsStatus() {
	migration := MigrationInstance()
	defer migration.Close()

	version, dirty, err := migration.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			log.Println("⚠️ AVISO: Nenhuma migração foi aplicada ainda! O banco está vazio. Rode 'make db-migrate' para criar as tabelas.")
			return
		}
		log.Printf("⚠️ AVISO: Falha ao verificar versão das migrações: %v\n", err)
		return
	}

	if dirty {
		log.Printf("🚨 ALERTA CRÍTICO: A migração da versão %d está marcada como DIRTY (Corrompida). Verifique o banco de dados!\n", version)
		return
	}

	log.Printf("✅ Verificação de Tabelas: Banco de dados operando de forma saudável na Versão %d.\n", version)
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
