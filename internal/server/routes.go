package server

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"api-hotelaria/internal/domain/hotel"
)

func (s *Server) RegisterRoutes() http.Handler {
	router := gin.Default()

	hotelRepo := hotel.NewRepository(s.db)
	hotelService := hotel.NewService(hotelRepo)
	hotelHandler := hotel.NewHandler(hotelService)

	router.POST("/hotels", hotelHandler.Create)

	return router
}