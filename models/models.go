package models

import (
	"errors"
	"time"
)

type Appointment struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	TrainerID int       `json:"trainer_id"`
	StartsAt  time.Time `json:"starts_at"`
	EndsAt    time.Time `json:"ends_at"`
}

func NewAppointment(userID, trainerID int, startsAt, endsAt time.Time) (*Appointment, error) {
	if userID < 1 {
		return nil, errors.New("UserID must be greater than 0")
	}

	if trainerID < 1 {
		return nil, errors.New("TrainerID must be greater than 0")
	}

	startsAt = ConvertToFixedTZ(startsAt)
	endsAt = ConvertToFixedTZ(endsAt)

	if startsAt.Local().Before(time.Now().Add(time.Hour)) {
		return nil, errors.New("Appointments must be scheduled at least 1 hour in advance")
	}

	if startsAt.Equal(endsAt) || startsAt.After(endsAt) {
		return nil, errors.New("Appointment start time must be before end time")
	}

	if startsAt.Hour() < 8 || startsAt.Hour() >= 17 {
		return nil, errors.New("Appointment must be scheduled between 8am and 5pm PST")
	}

	if endsAt.Hour() < 8 || endsAt.Hour() > 17 {
		return nil, errors.New("Appointment must be scheduled between 8am and 5pm PST")
	}

	if int(startsAt.Weekday()) < 1 || int(startsAt.Weekday()) > 5 {
		return nil, errors.New("Appointment must be scheduled between Monday and Friday PST")
	}

	if int(endsAt.Weekday()) < 1 || int(endsAt.Weekday()) > 5 {
		return nil, errors.New("Appointment must be scheduled between Monday and Friday PST")
	}

	if startsAt.Minute() != 0 && startsAt.Minute() != 30 {
		return nil, errors.New("Appointment must be scheduled on the hour or half hour PST")
	}

	if endsAt.Minute() != 0 && endsAt.Minute() != 30 {
		return nil, errors.New("Appointment must be scheduled on the hour or half hour PST")
	}

	if !startsAt.Add(time.Minute * 30).Equal(endsAt) {
		return nil, errors.New("Appointment must be scheduled in 30-minute increments")
	}

	return &Appointment{
		UserID:    userID,
		TrainerID: trainerID,
		StartsAt:  startsAt,
		EndsAt:    endsAt,
	}, nil
}

type Timeslot struct {
	StartsAt time.Time `json:"starts_at"`
	EndsAt   time.Time `json:"ends_at"`
}

func NewTimeslot(startsAt, endsAt time.Time) Timeslot {
	return Timeslot{
		StartsAt: ConvertToFixedTZ(startsAt),
		EndsAt:   ConvertToFixedTZ(endsAt),
	}
}
