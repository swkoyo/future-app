package store

import (
	"database/sql"
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
		&data.UserID,
		&data.TrainerID,
		(&data.StartedAt).Format(time.RFC3339),
		(&data.EndedAt).Format(time.RFC3339),
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
