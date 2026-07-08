package main

import (
	"log"
	"api-hotelaria/internal/config"
    "api-hotelaria/internal/database"
	"api-hotelaria/internal/server"
)

func main() {
	config.LoadConfig()
	database.ConnectDB()

	defer database.CloseDB()

	serverInstance := server.NewServer(database.DB)
	log.Println("Servidor de Hotelaria iniciando na porta :8080")

	err := serverInstance.ListenAndServe()
	if err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}