package server

import (
	"api-hotelaria/internal/domain/booking"
	"api-hotelaria/internal/domain/hotel"
	"api-hotelaria/internal/domain/room"
	"net/http"

	"github.com/gin-gonic/gin"
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

	router.GET("/hotels/:id/rooms", roomHandler.List)
	router.POST("/hotels/:id/rooms", roomHandler.Create)

	bookingRepo := booking.NewRepository(s.db)
	bookingService := booking.NewService(bookingRepo, roomRepo)
	bookingHandler := booking.NewHandler(bookingService)

	router.POST("/rooms/:id/check-in", bookingHandler.CheckIn)
	router.PATCH("/rooms/:id/check-out", bookingHandler.CheckOut)

	return router
}
