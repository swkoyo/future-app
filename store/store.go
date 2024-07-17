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
        started_at DATETIME NOT NULL,
        ended_at DATETIME NOT NULL
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
	INSERT INTO appointments (user_id, trainer_id, started_at, ended_at)
	VALUES ($1, $2, $3, $4)
	`

	res, err := s.DB.Exec(
		query,
		data.UserID,
		data.TrainerID,
		data.StartedAt.Format(time.RFC3339),
		data.EndedAt.Format(time.RFC3339),
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
	WHERE (user_id = $1 OR trainer_id = $2) AND started_at = $3 AND ended_at = $4
	`

	if err := s.DB.QueryRow(
		query,
		data.UserID,
		data.TrainerID,
		data.StartedAt.Format(time.RFC3339),
		data.EndedAt.Format(time.RFC3339),
	).Scan(&count); err != nil {
		return err
	}

	if count != 0 {
		return errors.New("Timeslot is not available")
	}

	return nil
}

func (s *Store) GetAppointmentsByTrainerID(trainerID int, from, to time.Time) ([]*models.Appointment, error) {
	appointments := make([]*models.Appointment, 0)
	query := `
	SELECT id, user_id, trainer_id, started_at, ended_at
	FROM appointments
	WHERE trainer_id = $1
	AND (
		(started_at >= $2 AND started_at <= $3)
		OR (ended_at >= $2 AND ended_at <= $3)
		OR (started_at >= $2 AND ended_at <= $3)
	)
	ORDER BY started_at ASC
	`
	rows, err := s.DB.Query(
		query,
		trainerID,
		from.Format(time.RFC3339),
		to.Format(time.RFC3339),
	)
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
			&appointment.StartedAt,
			&appointment.EndedAt,
		); err != nil {
			return nil, err
		}

		appointment.StartedAt = models.ConvertToFixedTZ(appointment.StartedAt)
		appointment.EndedAt = models.ConvertToFixedTZ(appointment.EndedAt)

		appointments = append(appointments, &appointment)
	}

	return appointments, nil
}
