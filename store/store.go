package store

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

func NewStore() (*Store, error) {
	db, err := sql.Open("sqlite3", "./store.db")
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
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

	if _, err := s.db.Exec(query); err != nil {
		return err
	}

	return nil
}

func (s *Store) Close() {
	s.db.Close()
}
