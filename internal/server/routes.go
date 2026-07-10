package server

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"api-hotelaria/internal/domain/hotel"
	"api-hotelaria/internal/domain/room"
)

func (s *Server) RegisterRoutes() http.Handler {
	router := gin.Default()

	hotelRepo := hotel.NewRepository(s.db)
	hotelService := hotel.NewService(hotelRepo)
	hotelHandler := hotel.NewHandler(hotelService)

	router.GET("/hotels", hotelHandler.List)
	router.POST("/hotels", hotelHandler.Create)

	roomRepo := room.NewRepository(s.db)
	roomService := room.NewService(roomRepo)
	roomHandler := room.NewHandler(roomService)

	router.POST("/rooms", roomHandler.Create)

	return router
}