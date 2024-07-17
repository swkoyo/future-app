package store

import (
	"future-app/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getTestAppointment() *models.Appointment {
	tz := time.FixedZone(models.GLOBAL_TZ, models.GLOBAL_TZ_OFFSET)
	startedAt := time.Date(2030, 7, 5, 8, 0, 0, 0, tz)
	endedAt := startedAt.Add(time.Minute * 30)

	return &models.Appointment{
		UserID:    1,
		TrainerID: 1,
		StartedAt: startedAt,
		EndedAt:   endedAt,
	}
}
func setupStore() (*Store, error) {
	db, err := NewTestStore()
	if err != nil {
		return nil, err
	}

	err = db.Init()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestCreateAppointment(t *testing.T) {
	store, err := setupStore()
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	appointment := getTestAppointment()

	createdAppointment, err := store.CreateAppointment(appointment)
	assert.NoError(t, err)
	assert.NotNil(t, createdAppointment)
	assert.NotZero(t, createdAppointment.ID)
	assert.Equal(t, appointment.UserID, createdAppointment.UserID)
	assert.Equal(t, appointment.TrainerID, createdAppointment.TrainerID)
	assert.Equal(t, appointment.StartedAt, createdAppointment.StartedAt)
	assert.Equal(t, appointment.EndedAt, createdAppointment.EndedAt)
}

func TestValidateAvailableTimeslot(t *testing.T) {
	store, err := setupStore()
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	appointment := getTestAppointment()

	t.Run("Available timeslot", func(t *testing.T) {
		err = store.ValidateAvailableTimeslot(appointment)
		assert.NoError(t, err)

		createdAppointment, err := store.CreateAppointment(appointment)
		assert.NoError(t, err)
		assert.NotNil(t, createdAppointment)
	})

	t.Run("Trainer busy during timeslot", func(t *testing.T) {
		err = store.ValidateAvailableTimeslot(&models.Appointment{
			UserID:    2,
			TrainerID: 1,
			StartedAt: appointment.StartedAt,
			EndedAt:   appointment.EndedAt,
		})
		assert.Error(t, err)
		assert.Equal(t, "Timeslot is not available", err.Error())
	})

	t.Run("User busy during timeslot", func(t *testing.T) {
		err = store.ValidateAvailableTimeslot(&models.Appointment{
			UserID:    1,
			TrainerID: 2,
			StartedAt: appointment.StartedAt,
			EndedAt:   appointment.EndedAt,
		})
		assert.Error(t, err)
		assert.Equal(t, "Timeslot is not available", err.Error())
	})
}
