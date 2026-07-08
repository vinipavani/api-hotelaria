package server

import (
	"fmt"
	"net/http"
	"time"
	"api-hotelaria/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	db *pgxpool.Pool
}

func NewServer(database *pgxpool.Pool) *http.Server {
	newServer := &Server{
		db: database,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.Env.Port), 
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	return server
}