package repository

import (
	"github.com/zahnah/study-app/internal/models"
	"time"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)

	InsertRoomRestriction(res models.RoomRestriction) (int, error)

	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)

	SearchAvailabilityByRoomID(start, end time.Time, roomID int) (int, error)
}
