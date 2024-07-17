package store

import (
	"database/sql"
	"errors"
	"future-app/models"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	DB *sql.DB
}

func NewStore() (*Store, error) {
	db, err := sql.Open("sqlite3", "./store.db")
	if err != nil {
		return nil, err
	}
	return &Store{DB: db}, nil
}

func NewTestStore() (*Store, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	return &Store{DB: db}, nil
}

func (s *Store) Init() error {
	return s.createAppointmentTable()
}

func (s *Store) createAppointmentTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS appointments (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        trainer_id INTEGER NOT NULL,
        starts_at DATETIME NOT NULL,
        ends_at DATETIME NOT NULL
    );
    `

	if _, err := s.DB.Exec(query); err != nil {
		return err
	}

	return nil
}

func (s *Store) Close() {
	s.DB.Close()
}

func (s *Store) CreateAppointment(data *models.Appointment) (*models.Appointment, error) {
	query := `
	INSERT INTO appointments (user_id, trainer_id, starts_at, ends_at)
	VALUES ($1, $2, $3, $4)
	`

	res, err := s.DB.Exec(
		query,
		data.UserID,
		data.TrainerID,
		data.StartsAt.Format(time.RFC3339),
		data.EndsAt.Format(time.RFC3339),
	)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return nil, err
	}

	data.ID = int(id)
	return data, nil
}

func (s *Store) ValidateAvailableTimeslot(data *models.Appointment) error {
	var count int

	query := `
	SELECT COUNT(*)
	FROM appointments
	WHERE (user_id = $1 OR trainer_id = $2) AND starts_at = $3 AND ends_at = $4
	`

	if err := s.DB.QueryRow(
		query,
		data.UserID,
		data.TrainerID,
		data.StartsAt.Format(time.RFC3339),
		data.EndsAt.Format(time.RFC3339),
	).Scan(&count); err != nil {
		return err
	}

	if count != 0 {
		return errors.New("Timeslot is not available")
	}

	return nil
}

func (s *Store) GetAppointmentsByTrainerID(trainerID int, startsAt, endsAt time.Time) ([]*models.Appointment, error) {
	appointments := make([]*models.Appointment, 0)

	var rows *sql.Rows
	var err error

	if startsAt.IsZero() || endsAt.IsZero() {
		query := `
		SELECT id, user_id, trainer_id, starts_at, ends_at
		FROM appointments
		WHERE trainer_id = $1
		ORDER BY starts_at ASC
		`
		rows, err = s.DB.Query(
			query,
			trainerID,
		)
	} else {
		query := `
		SELECT id, user_id, trainer_id, starts_at, ends_at
		FROM appointments
		WHERE trainer_id = $1
		AND (
			(starts_at >= $2 AND starts_at <= $3)
			OR (ends_at >= $2 AND ends_at <= $3)
			OR (starts_at >= $2 AND ends_at <= $3)
		)
		ORDER BY starts_at ASC
		`
		rows, err = s.DB.Query(
			query,
			trainerID,
			startsAt.Format(time.RFC3339),
			endsAt.Format(time.RFC3339),
		)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var appointment models.Appointment
		if err := rows.Scan(
			&appointment.ID,
			&appointment.UserID,
			&appointment.TrainerID,
			&appointment.StartsAt,
			&appointment.EndsAt,
		); err != nil {
			return nil, err
		}

		appointment.StartsAt = models.ConvertToFixedTZ(appointment.StartsAt)
		appointment.EndsAt = models.ConvertToFixedTZ(appointment.EndsAt)

		appointments = append(appointments, &appointment)
	}

	return appointments, nil
}

func (s *Store) GetTrainerAvailability(trainerID int, startsAt, endsAt time.Time) (*[]models.Timeslot, error) {
	appointments, err := s.GetAppointmentsByTrainerID(trainerID, startsAt, endsAt)
	if err != nil {
		return nil, err
	}

	timeslots := make([]models.Timeslot, 0)
	currAppIdx := 0

	for date := startsAt; date.Before(endsAt); date = date.Add(24 * time.Hour) {
		if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
			continue
		}

		currentDate := time.Date(date.Year(), date.Month(), date.Day(), 8, 0, 0, 0, date.Location())

		for currentDate.Before(time.Date(date.Year(), date.Month(), date.Day(), 17, 0, 0, 0, date.Location())) {
			if currAppIdx < len(appointments) && appointments[currAppIdx].StartsAt.Equal(currentDate) && appointments[currAppIdx].EndsAt.Equal(currentDate.Add(30*time.Minute)) {
				currentDate = appointments[currAppIdx].EndsAt
				currAppIdx += 1
				continue
			}

			timeslot := models.NewTimeslot(currentDate, currentDate.Add(30*time.Minute))

			timeslots = append(timeslots, timeslot)
			currentDate = currentDate.Add(30 * time.Minute)
		}
	}

	return &timeslots, nil
}
