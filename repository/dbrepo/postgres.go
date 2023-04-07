package dbrepo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/zahnah/study-app/internal/config"
	"github.com/zahnah/study-app/internal/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type postgresDbRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func (m *postgresDbRepo) InsertBlockForRoom(id int, startDate time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
insert into room_restrictions (restriction_id, reservation_id, room_id,
                               start_date, end_date,
                               created_at, updated_at)
values ($1, $2, $3, $4, $5, $6, $7)`
	_, err := m.DB.ExecContext(ctx, stmt,
		2,
		nil,
		id,
		startDate,
		startDate.AddDate(0, 0, 1),
		time.Now(),
		time.Now(),
	)

	return err
}

func (m *postgresDbRepo) DeleteRoomRestriction(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `delete from room_restrictions where id = $1`
	_, err := m.DB.ExecContext(ctx, stmt, id)
	return err
}

func (m *postgresDbRepo) GetRestrictionsForRoomByDate(roomID int, startDate, endDate time.Time) ([]models.RoomRestriction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var restrictions []models.RoomRestriction

	stmt := `
select rr.id, rr.room_id,
       coalesce(rr.reservation_id, 0), rr.restriction_id,
       rr.start_date, rr.end_date,
       rr.created_at, rr.updated_at
from room_restrictions rr
where rr.room_id = $1 and rr.end_date >= $2 and rr.start_date <= $3
`
	rows, err := m.DB.QueryContext(ctx, stmt, roomID, startDate, endDate)
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	if err != nil {
		return restrictions, err
	}

	for rows.Next() {
		var r models.RoomRestriction
		err := rows.Scan(
			&r.ID, &r.RoomID,
			&r.ReservationID, &r.RestrictionID,
			&r.StartDate, &r.EndDate,
			&r.CreatedAt, &r.UpdatedAt,
		)

		if err != nil {
			return restrictions, err
		}
		restrictions = append(restrictions, r)
	}

	if err = rows.Err(); err != nil {
		return restrictions, err
	}

	return restrictions, nil
}

func (m *postgresDbRepo) AllRooms() ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	stmt := `
select r.id, r.room_name,
       r.created_at, r.updated_at
from rooms r
`
	rows, err := m.DB.QueryContext(ctx, stmt)
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var r models.Room
		err := rows.Scan(
			&r.ID, &r.RoomName,
			&r.CreatedAt, &r.UpdatedAt,
		)

		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, r)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

func (m *postgresDbRepo) UpdateProcessedForReservations(id, processed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `update reservations set processed = $2 where id = $1`
	_, err := m.DB.ExecContext(ctx, stmt, id, processed)
	return err
}

func (m *postgresDbRepo) DeleteReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `delete from reservations where id = $1`
	_, err := m.DB.ExecContext(ctx, stmt, id)
	return err
}

func (m *postgresDbRepo) UpdateReservation(r models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
update reservations
set
    first_name = $2, last_name = $3,
    email = $4, phone = $5,
    updated_at = $6
where id = $1`
	_, err := m.DB.ExecContext(ctx, stmt,
		r.ID,
		r.FirstName, r.LastName,
		r.Email, r.Phone,
		time.Now(),
	)
	return err
}

func (m *postgresDbRepo) GetReservationByID(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
select res.id, res.first_name, res.last_name,
       res.email, res.phone, res.start_date, res.end_date, res.room_id,
       res.created_at, res.updated_at, res.processed,
       r.id, r.room_name
from reservations res
left join rooms r on r.id = res.room_id
where res.id = $1`
	row := m.DB.QueryRowContext(ctx, stmt, id)

	var r models.Reservation
	err := row.Scan(
		&r.ID, &r.FirstName, &r.LastName,
		&r.Email,
		&r.Phone,
		&r.StartDate,
		&r.EndDate,
		&r.RoomID,
		&r.CreatedAt,
		&r.UpdatedAt,
		&r.Processed,
		&r.Room.ID,
		&r.Room.RoomName,
	)

	if err != nil {
		return r, err
	}
	return r, err
}

func (m *postgresDbRepo) AllNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	stmt := `
select res.id, res.first_name, res.last_name,
       res.email, res.phone, res.start_date, res.end_date, res.room_id,
       res.created_at, res.updated_at, res.processed,
       r.id, r.room_name
from reservations res
left join rooms r on r.id = res.room_id
where processed = 0
`
	rows, err := m.DB.QueryContext(ctx, stmt)
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	if err != nil {
		return reservations, err
	}

	for rows.Next() {
		var r models.Reservation
		err := rows.Scan(
			&r.ID, &r.FirstName, &r.LastName,
			&r.Email,
			&r.Phone,
			&r.StartDate,
			&r.EndDate,
			&r.RoomID,
			&r.CreatedAt,
			&r.UpdatedAt,
			&r.Processed,
			&r.Room.ID,
			&r.Room.RoomName,
		)

		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, r)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

func (m *postgresDbRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	stmt := `
select res.id, res.first_name, res.last_name,
       res.email, res.phone, res.start_date, res.end_date, res.room_id,
       res.created_at, res.updated_at, res.processed,
       r.id, r.room_name
from reservations res
left join rooms r on r.id = res.room_id
`
	rows, err := m.DB.QueryContext(ctx, stmt)
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	if err != nil {
		return reservations, err
	}

	for rows.Next() {
		var r models.Reservation
		err := rows.Scan(
			&r.ID, &r.FirstName, &r.LastName,
			&r.Email,
			&r.Phone,
			&r.StartDate,
			&r.EndDate,
			&r.RoomID,
			&r.CreatedAt,
			&r.UpdatedAt,
			&r.Processed,
			&r.Room.ID,
			&r.Room.RoomName,
		)

		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, r)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

func (m *postgresDbRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User
	stmt := `
select id, first_name, last_name,
       email, password, access_level,
       created_at, updated_at
from users
where id = $1`
	row := m.DB.QueryRowContext(ctx, stmt, id)
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName,
		&user.Email, &user.Password, &user.AccessLevel,
		&user.CreatedAt, &user.UpdatedAt)
	return user, err
}

func (m *postgresDbRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
update users 
set
    first_name = $2, last_name = $3,
    email = $4, access_level = $5,
    updated_at = $6
where id = $1`
	_, err := m.DB.ExecContext(ctx, stmt,
		u.ID,
		u.FirstName, u.LastName,
		u.Email, u.AccessLevel,
		time.Now(),
	)
	return err
}

func (m *postgresDbRepo) Authenticate(email, password string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "select id, password from users where email = $1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return 0, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
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

func (m *postgresDbRepo) SearchAvailabilityByRoomID(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRows int
	stmt := `
select count(1)
from room_restrictions rr
where
    rr.room_id = $1
    AND ((rr.start_date <= $2 AND rr.end_date > $2) OR (rr.start_date <= $3 AND rr.end_date > $3)) limit 1`
	err := m.DB.QueryRowContext(ctx, stmt,
		roomID,
		start,
		end,
	).Scan(&numRows)

	return numRows == 0, err
}

func (m *postgresDbRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room
	stmt := `
select r.id, r.room_name, r.created_at, r.updated_at
from rooms r
where r.id not in (
	select distinct rr.room_id
	from room_restrictions rr
	where (rr.start_date <= $1 AND rr.end_date > $1) OR (rr.start_date <= $2 AND rr.end_date > $2)   
)`
	rows, err := m.DB.QueryContext(ctx, stmt,
		start,
		end,
	)

	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room

		err = rows.Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	return rooms, err
}

func (m *postgresDbRepo) GetRoomById(roomID int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room
	stmt := `
select r.id, r.room_name, r.created_at, r.updated_at
from rooms r
where r.id = $1`
	row := m.DB.QueryRowContext(ctx, stmt, roomID)
	err := row.Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
	return room, err
}
