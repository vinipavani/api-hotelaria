package booking

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

type BookingResponse struct {
	ID            int64         `json:"id"`
	RoomID        int64         `json:"room_id"`
	GuestName     string        `json:"guest_name"`
	GuestDocument string        `json:"guest_document"`
	Status        BookingStatus `json:"status"`
	CheckInDate   string        `json:"check_in"`
	CheckOutDate  string        `json:"check_out,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CheckIn(c *gin.Context) {
	RoomIDParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "RoomID deve ser um número inteiro."})
		return
	}

	RoomID := int64(RoomIDParam)
	var input CheckInInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	newBooking, err := h.service.CreateCheckIn(ctx, RoomID, input)
	switch {
	case err == InvalidGuestParams || err == InvalidCheckInDateFormat:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	case err == RoomNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	case err == RoomNotAvailable:
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao realizar o check-in: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, buildResponseBooking(newBooking))
}

func (h *Handler) CheckOut(c *gin.Context) {
	RoomIDParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "RoomID deve ser um número inteiro."})
		return
	}

	var input CheckOutInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	booking, err := h.service.CheckOut(ctx, int64(RoomIDParam), input)

	switch {
	case err == InvalidCheckOutDateFormat:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	case err == RoomNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	case err == RoomAvailable:
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	case err == CheckOutLesserThanCheckIn:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao realizar o check-out: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, buildResponseBooking(booking))
}

func buildResponseBooking(b *Booking) BookingResponse {
	checkoutStr := ""
	if b.CheckOutDate != nil && !b.CheckOutDate.IsZero() {
		checkoutStr = b.CheckOutDate.Format("2006-01-02")
	}

	return BookingResponse{
		ID:            b.ID,
		RoomID:        b.RoomID,
		GuestName:     b.GuestName,
		GuestDocument: b.GuestDocument,
		Status:        b.Status,
		CheckInDate:   b.CheckInDate.Format("2006-01-02"),
		CheckOutDate:  checkoutStr,
		CreatedAt:     b.CreatedAt,
	}
}
