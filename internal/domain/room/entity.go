package room

import "time"

type RoomType string
const (
	RoomTypeSingle RoomType = "single"
	RoomTypeDouble RoomType = "double"
	RoomTypeSuite  RoomType = "suite"
)

type Room struct {
	ID            int64     `json:"id"`
	HotelID       int64     `json:"hotel_id"`
	Number        string    `json:"number"`
	Type          RoomType  `json:"type"`
	Capacity      int       `json:"capacity"`
	PerNightValue float64   `json:"per_night_value"`
	CreatedAt     time.Time `json:"created_at"`
}

type CreateRoomInput struct {
	HotelID       int64    `json:"-"`
	Type          RoomType `json:"type" binding:"required"`
	Capacity      int      `json:"capacity" binding:"required"`
	PerNightValue float64  `json:"per_night_value" binding:"required"`
}