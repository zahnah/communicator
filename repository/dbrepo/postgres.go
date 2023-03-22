package dbrepo

import (
	"context"
	"database/sql"
	"github.com/zahnah/study-app/internal/config"
	"github.com/zahnah/study-app/internal/models"
	"time"
)

type postgresDbRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func (*postgresDbRepo) AllUsers() bool {
	return true
}

func (m *postgresDbRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	stmt := `
insert into reservations (first_name, last_name, email,
                          phone, start_date, end_date,
                          room_id, created_at, updated_at)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`
	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)
	return newID, err
}

func (m *postgresDbRepo) InsertRoomRestriction(res models.RoomRestriction) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int
	stmt := `
insert into room_restrictions (restriction_id, reservation_id, room_id,
                               start_date, end_date,
                               created_at, updated_at)
values ($1, $2, $3, $4, $5, $6, $7) returning id`
	err := m.DB.QueryRowContext(ctx, stmt,
		res.RestrictionID,
		res.ReservationID,
		res.RoomID,
		res.StartDate,
		res.EndDate,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	return newID, err
}

func (m *postgresDbRepo) SearchAvailabilityByDates(start, end time.Time, roomID int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRows int
	stmt := `
select count(1)
from room_restrictions
where
    room_id = $1
    && start_date < $2 && end_date > $3`
	err := m.DB.QueryRowContext(ctx, stmt,
		roomID,
		start,
		end,
	).Scan(&numRows)

	return numRows, err
}
