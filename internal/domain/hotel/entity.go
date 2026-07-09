package hotel

import "time"

type Hotel struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	City      string    `json:"city"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateHotelInput struct {
	Name string `json:"name" binding:"required"`
	City string `json:"city" binding:"required"`
}