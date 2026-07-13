package room

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

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) List(c *gin.Context) {
	hotelID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do hotel inválido. Deve ser um número inteiro."})
		return
	}

	availableOnly := false

	if availableParam := c.Query("disponivel"); availableParam != "" {
		availableOnly, err = strconv.ParseBool(availableParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Erro no parametro 'disponivel', o valor deve ser booleano: " + err.Error()})
			return
		}
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	rooms, err := h.service.findAllRooms(ctx, hotelID, availableOnly)
	if err == HotelNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao listar os quartos: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, rooms)
}

func (h *Handler) Create(c *gin.Context) {
	hotelID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do hotel inválido. Deve ser um número inteiro."})
		return
	}

	var input CreateRoomInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	input.HotelID = int64(hotelID)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	newRoom, err := h.service.CreateRoom(ctx, input)
	switch {
	case err == HotelNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	case err == InvalidParams || err == InvalidRoomType:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newRoom)
}
