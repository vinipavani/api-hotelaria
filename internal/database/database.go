package database

import (
	"context"
	"log"
	"time"
	"github.com/jackc/pgx/v5/pgxpool"
	"api-hotelaria/internal/config"
)

var DB *pgxpool.Pool

func ConnectDB() {
	config, err := pgxpool.ParseConfig(config.Env.DatabaseURL)
	if err != nil {
		log.Fatalf("Erro ao processar a string de conexão do banco: %v\n", err)
	}

	config = configDatabaseConnection(config)

	DB, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Erro ao criar o pool de conexões do Postgres: %v\n", err)
	}

	err = DB.Ping(context.Background())
	if err != nil {
		log.Fatalf("O contêiner do Postgres respondeu, mas recusou a conexão: %v\n", err)
	}

	log.Println("Conexão com o PostgreSQL do Docker estabelecida com sucesso!")
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Pool de conexões com o banco de dados fechado com sucesso.")
	}
}

func configDatabaseConnection(config *pgxpool.Config) *pgxpool.Config {
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnIdleTime = time.Minute * 5 // Fecha conexões inativas após 5 minutos
	
	return config
}