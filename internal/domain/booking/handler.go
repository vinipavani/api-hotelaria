package booking

import (
	"context"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CheckIn(c *gin.Context) {
	RoomID := int64(c.Param("id"))
	var input BookingInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	newBooking, err := h.service.CreateCheckIn(ctx, RoomID, input)
	switch {
	case err == GuestParamsError || err == CheckInDateError:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	case err == ErrRoomNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar o check-in: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newBooking)
}
