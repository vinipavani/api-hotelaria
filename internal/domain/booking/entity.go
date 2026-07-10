package booking

import "time"

type BookingStatus string

const (
	BookingStatusInProgress BookingStatus = "em_estadia"
	BookingStatusFinished   BookingStatus = "finalizada"
)

type Booking struct {
	ID            int64         `json:"id"`
	RoomID        int64         `json:"room_id"`
	GuestName     string        `json:"guest_name"`
	GuestDocument string        `json:"guest_document"`
	Status        BookingStatus `json:"status"`
	CheckInDate	  time.Time     `json:"check_in"`
	CheckOutDate  time.Time     `json:"check_out"`
	CreatedAt     time.Time     `json:"created_at"`
}

type BookingInput struct {
	RoomID        int64     `json:"room_id"`
	GuestName     string    `json:"guest_name" binding:"required"`
	GuestDocument string    `json:"guest_document" binding:"required"`
	CheckInDate	  time.Time `json:"check_in" binding:"required" format:"2006-01-02"`
	CheckOutDate  time.Time `json:"check_out" format:"2006-01-02"`
}