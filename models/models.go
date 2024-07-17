package models

import (
	"errors"
	"time"
)

type Appointment struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	TrainerID int       `json:"trainer_id"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"`
}

func NewAppointment(userID, trainerID int, startedAt, endedAt time.Time) (*Appointment, error) {
	if userID < 1 {
		return nil, errors.New("UserID must be greater than 0")
	}

	if trainerID < 1 {
		return nil, errors.New("TrainerID must be greater than 0")
	}

	startedAt = ConvertToFixedTZ(startedAt)
	endedAt = ConvertToFixedTZ(endedAt)

	if startedAt.Local().Before(time.Now().Add(time.Hour)) {
		return nil, errors.New("Appointments must be scheduled at least 1 hour in advance")
	}

	if startedAt.Equal(endedAt) || startedAt.After(endedAt) {
		return nil, errors.New("Appointment start time must be before end time")
	}

	if startedAt.Hour() < 8 || startedAt.Hour() >= 17 {
		return nil, errors.New("Appointment must be scheduled between 8am and 5pm PST")
	}

	if endedAt.Hour() < 8 || endedAt.Hour() > 17 {
		return nil, errors.New("Appointment must be scheduled between 8am and 5pm PST")
	}

	if int(startedAt.Weekday()) < 1 || int(startedAt.Weekday()) > 5 {
		return nil, errors.New("Appointment must be scheduled between Monday and Friday PST")
	}

	if int(endedAt.Weekday()) < 1 || int(endedAt.Weekday()) > 5 {
		return nil, errors.New("Appointment must be scheduled between Monday and Friday PST")
	}

	if startedAt.Minute() != 0 && startedAt.Minute() != 30 {
		return nil, errors.New("Appointment must be scheduled on the hour or half hour PST")
	}

	if endedAt.Minute() != 0 && endedAt.Minute() != 30 {
		return nil, errors.New("Appointment must be scheduled on the hour or half hour PST")
	}

	if !startedAt.Add(time.Minute * 30).Equal(endedAt) {
		return nil, errors.New("Appointment must be scheduled in 30-minute increments")
	}

	return &Appointment{
		UserID:    userID,
		TrainerID: trainerID,
		StartedAt: startedAt,
		EndedAt:   endedAt,
	}, nil
}
