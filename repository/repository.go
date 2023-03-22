package repository

import "github.com/zahnah/study-app/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)

	InsertRoomRestriction(res models.RoomRestriction) (int, error)
}
