package dbrepo

import (
	"database/sql"
	"errors"
	"github.com/zahnah/study-app/internal/config"
	"github.com/zahnah/study-app/internal/models"
	"time"
)

type testDbRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func (t testDbRepo) GetRestrictionsForRoomByDate(roomID int, startDate, rndDate time.Time) ([]models.RoomRestriction, error) {
	var restrictions []models.RoomRestriction
	return restrictions, nil
}

func (t testDbRepo) AllRooms() ([]models.Room, error) {
	return make([]models.Room, 1), nil
}

func (t testDbRepo) UpdateProcessedForReservations(id, processed int) error {
	return nil
}

func (t testDbRepo) DeleteReservation(id int) error {
	return nil
}

func (t testDbRepo) UpdateReservation(r models.Reservation) error {
	return nil
}

func (t testDbRepo) GetReservationByID(id int) (models.Reservation, error) {
	return models.Reservation{}, nil
}

func (t testDbRepo) AllNewReservations() ([]models.Reservation, error) {
	return make([]models.Reservation, 1), nil
}

func (t testDbRepo) AllReservations() ([]models.Reservation, error) {
	return make([]models.Reservation, 1), nil
}

func (t testDbRepo) GetUserByID(id int) (models.User, error) {
	return models.User{}, nil
}

func (t testDbRepo) UpdateUser(u models.User) error {
	return nil
}

func (t testDbRepo) Authenticate(email, password string) (int, string, error) {
	return 0, "", nil
}

func (t testDbRepo) AllUsers() bool {
	return true
}

func (t testDbRepo) InsertReservation(res models.Reservation) (int, error) {
	if res.RoomID > 2 {
		return 0, errors.New("can't find the room")
	}
	return 1, nil
}

func (t testDbRepo) InsertRoomRestriction(res models.RoomRestriction) (int, error) {
	if res.RoomID > 1 {
		return 0, errors.New("can't find the room")
	}
	return 1, nil
}

func (t testDbRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	return nil, nil
}

func (t testDbRepo) SearchAvailabilityByRoomID(start, end time.Time, roomID int) (bool, error) {
	if roomID == 3 {
		return false, errors.New("can't find the room")
	}
	return false, nil
}

func (t testDbRepo) GetRoomById(roomID int) (models.Room, error) {
	var room models.Room
	if roomID > 2 {
		return room, errors.New("can't find the room")
	}
	return room, nil
}
