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

	SearchAvailabilityByRoomID(start, end time.Time, roomID int) (bool, error)

	GetRoomById(roomID int) (models.Room, error)

	GetUserByID(id int) (models.User, error)

	UpdateUser(u models.User) error

	Authenticate(email, password string) (int, string, error)

	AllReservations() ([]models.Reservation, error)

	AllNewReservations() ([]models.Reservation, error)

	GetReservationByID(id int) (models.Reservation, error)
}
