package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewAppointment(t *testing.T) {
	tz := time.FixedZone(GLOBAL_TZ, GLOBAL_TZ_OFFSET)
	now := time.Now().In(tz)
	past := time.Date(2020, 7, 13, 0, 0, 0, 0, tz)  // Monday midnight
	future := time.Date(2030, 7, 8, 0, 0, 0, 0, tz) // Monday midnight

	tests := []struct {
		name      string
		userID    int
		trainerID int
		startedAt time.Time
		endedAt   time.Time
		hasErr    bool
		errMsg    string
		expected  *Appointment
	}{
		{
			name:      "valid appointment on monday at 8am",
			userID:    1,
			trainerID: 1,
			startedAt: future.Add(time.Hour * 8),
			endedAt:   future.Add(time.Hour * 8).Add(time.Minute * 30),
			hasErr:    false,
			expected: &Appointment{
				UserID:    1,
				TrainerID: 1,
				StartedAt: future.Add(time.Hour * 8),
				EndedAt:   future.Add(time.Hour * 8).Add(time.Minute * 30),
			},
		},
		{
			name:      "valid appointment on friday at 4:30pm",
			userID:    1,
			trainerID: 1,
			startedAt: future.Add(time.Hour * 24 * 4).Add(time.Hour * 16).Add(time.Minute * 30),
			endedAt:   future.Add(time.Hour * 24 * 4).Add(time.Hour * 17),
			hasErr:    false,
			expected: &Appointment{
				UserID:    1,
				TrainerID: 1,
				StartedAt: future.Add(time.Hour * 24 * 4).Add(time.Hour * 16).Add(time.Minute * 30),
				EndedAt:   future.Add(time.Hour * 24 * 4).Add(time.Hour * 17),
			},
		},
		{
			name:      "negative user ID",
			userID:    -1,
			trainerID: 1,
			startedAt: future.Add(time.Hour * 8),
			endedAt:   future.Add(time.Hour * 8).Add(time.Minute * 30),
			hasErr:    true,
			errMsg:    "UserID must be greater than 0",
		},
		{
			name:      "negative trainer ID",
			userID:    1,
			trainerID: -1,
			startedAt: future.Add(time.Hour * 8),
			endedAt:   future.Add(time.Hour * 8).Add(time.Minute * 30),
			hasErr:    true,
			errMsg:    "TrainerID must be greater than 0",
		},
		{
			name:      "start time in past",
			userID:    1,
			trainerID: 1,
			startedAt: past.Add(time.Hour * 8),
			endedAt:   past.Add(time.Hour * 8).Add(time.Minute * 30),
			hasErr:    true,
			errMsg:    "Appointments must be scheduled at least 1 hour in advance",
		},
		{
			name:      "start time less than 1 hour in advance",
			userID:    1,
			trainerID: 1,
			startedAt: now.Add(time.Minute * 30),
			endedAt:   now.Add(time.Hour),
			hasErr:    true,
			errMsg:    "Appointments must be scheduled at least 1 hour in advance",
		},
		{
			name:      "start time same as end time",
			userID:    1,
			trainerID: 1,
			startedAt: future.Add(time.Hour * 8),
			endedAt:   future.Add(time.Hour * 8),
			hasErr:    true,
			errMsg:    "Appointment start time must be before end time",
		},
		{
			name:      "start time after end time",
			userID:    1,
			trainerID: 1,
			startedAt: future.Add(time.Hour * 8).Add(time.Minute),
			endedAt:   future.Add(time.Hour * 8),
			hasErr:    true,
			errMsg:    "Appointment start time must be before end time",
		},
		{
			name:      "start time outside business hours",
			userID:    1,
			trainerID: 1,
			startedAt: future.Add(time.Hour * 7).Add(time.Minute * 30),
			endedAt:   future.Add(time.Hour * 8),
			hasErr:    true,
			errMsg:    "Appointment must be scheduled between 8am and 5pm PST",
		},
		{
			name:      "end time outside business hours",
			userID:    1,
			trainerID: 1,
			startedAt: future.Add(time.Hour * 17),
			endedAt:   future.Add(time.Hour * 17).Add(time.Minute * 30),
			hasErr:    true,
			errMsg:    "Appointment must be scheduled between 8am and 5pm PST",
		},
		{
			name:      "appointment on weekend",
			userID:    1,
			trainerID: 1,
			startedAt: future.Add(-time.Hour * 24).Add(time.Hour * 8),
			endedAt:   future.Add(-time.Hour * 24).Add(time.Hour * 8).Add(time.Minute * 30),
			hasErr:    true,
			errMsg:    "Appointment must be scheduled between Monday and Friday PST",
		},
		{
			name:      "appointment not on the hour or half hour",
			userID:    1,
			trainerID: 1,
			startedAt: future.Add(time.Hour * 8).Add(time.Minute * 15),
			endedAt:   future.Add(time.Hour * 8).Add(time.Minute * 45),
			hasErr:    true,
			errMsg:    "Appointment must be scheduled on the hour or half hour PST",
		},
		{
			name:      "appointment not in 30-minute increments",
			userID:    1,
			trainerID: 1,
			startedAt: future.Add(time.Hour * 8),
			endedAt:   future.Add(time.Hour * 9),
			hasErr:    true,
			errMsg:    "Appointment must be scheduled in 30-minute increments",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			appointment, err := NewAppointment(tc.userID, tc.trainerID, tc.startedAt, tc.endedAt)
			if tc.hasErr {
				assert.Error(t, err)
				assert.Nil(t, appointment)
				assert.Equal(t, tc.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, appointment)
			}
		})
	}
}
