package Service

import "errors"

var (
	ErrRoomNotFound          = errors.New("room not found")
	ErrScheduleAlreadyExists = errors.New("schedule already exists for this room")
	ErrSlotNotFound          = errors.New("slot not found")
	ErrSlotInPast            = errors.New("cannot book a slot in the past")
	ErrSlotAlreadyBooked     = errors.New("slot is already booked")
	ErrBookingNotFound       = errors.New("booking not found")
	ErrNotOwner              = errors.New("you can only cancel your own bookings")
)
