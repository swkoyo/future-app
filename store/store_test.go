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

func TestGetAppointmentsByTrainerID(t *testing.T) {
	store, err := setupStore()
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	appointment := getTestAppointment()

	// INFO: Create test appointment
	createdAppointment, err := store.CreateAppointment(appointment)
	assert.NoError(t, err)
	assert.NotNil(t, createdAppointment)

	t.Run("Appointment within timeframe", func(t *testing.T) {
		appointments, err := store.GetAppointmentsByTrainerID(1, appointment.StartedAt.Add(-time.Hour), appointment.EndedAt.Add(time.Hour))
		assert.NoError(t, err)
		assert.NotNil(t, appointments)
		assert.Len(t, appointments, 1)
		assert.Equal(t, createdAppointment.ID, appointments[0].ID)
		assert.Equal(t, createdAppointment.UserID, appointments[0].UserID)
		assert.Equal(t, createdAppointment.TrainerID, appointments[0].TrainerID)
		assert.Equal(t, createdAppointment.StartedAt, appointments[0].StartedAt)
		assert.Equal(t, createdAppointment.EndedAt, appointments[0].EndedAt)
	})

	t.Run("Appointment overlaps timeframe start", func(t *testing.T) {
		appointments, err := store.GetAppointmentsByTrainerID(1, appointment.StartedAt.Add(time.Minute*15), appointment.EndedAt.Add(time.Hour))
		assert.NoError(t, err)
		assert.NotNil(t, appointments)
		assert.Len(t, appointments, 1)
		assert.Equal(t, createdAppointment.ID, appointments[0].ID)
	})

	t.Run("Appointment overlaps timeframe end", func(t *testing.T) {
		appointments, err := store.GetAppointmentsByTrainerID(1, appointment.StartedAt.Add(-time.Hour), appointment.EndedAt.Add(-time.Minute*15))
		assert.NoError(t, err)
		assert.NotNil(t, appointments)
		assert.Len(t, appointments, 1)
		assert.Equal(t, createdAppointment.ID, appointments[0].ID)
	})

	t.Run("Get trainer with no appointments", func(t *testing.T) {
		appointments, err := store.GetAppointmentsByTrainerID(2, appointment.StartedAt.Add(-time.Hour), appointment.EndedAt.Add(time.Hour))
		assert.NoError(t, err)
		assert.NotNil(t, appointments)
		assert.Len(t, appointments, 0)
	})

	t.Run("Get trainer with no appointments within timeframe", func(t *testing.T) {
		appointments, err := store.GetAppointmentsByTrainerID(1, appointment.EndedAt.Add(time.Hour), appointment.EndedAt.Add(time.Hour*2))
		assert.NoError(t, err)
		assert.NotNil(t, appointments)
		assert.Len(t, appointments, 0)
	})
}

func TestGetTrainerAvailability(t *testing.T) {
	store, err := setupStore()
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	tz := time.FixedZone(models.GLOBAL_TZ, models.GLOBAL_TZ_OFFSET)
	from := time.Date(2030, 7, 5, 0, 0, 0, 0, tz) // Friday midnight
	to := time.Date(2030, 7, 8, 0, 0, 0, 0, tz)   // Monday midnight

	t.Run("Trainer with no appointments", func(t *testing.T) {
		timeslots, err := store.GetTrainerAvailability(1, from, to)
		assert.NoError(t, err)
		assert.NotNil(t, timeslots)
		assert.NotZero(t, len(*timeslots))

		// INFO: Timeslots are in order
		for i := 1; i < len(*timeslots); i++ {
			assert.True(t, (*timeslots)[i].StartedAt.After((*timeslots)[i-1].StartedAt))
		}

		for _, timeslot := range *timeslots {
			// INFO: Timeslot is between from and to
			assert.True(t, timeslot.StartedAt.After(from) || timeslot.StartedAt.Equal(from))
			assert.True(t, timeslot.EndedAt.Before(to) || timeslot.EndedAt.Equal(to))

			// INFO: Timeslot is not on Saturday or Sunday
			assert.NotEqual(t, time.Saturday, timeslot.StartedAt.Weekday())
			assert.NotEqual(t, time.Saturday, timeslot.EndedAt.Weekday())
			assert.NotEqual(t, time.Sunday, timeslot.StartedAt.Weekday())
			assert.NotEqual(t, time.Sunday, timeslot.EndedAt.Weekday())

			// INFO: Timeslot is between 8am and 5pm
			assert.True(t, timeslot.StartedAt.Hour() >= 8 && timeslot.StartedAt.Hour() < 17)
			assert.True(t, timeslot.EndedAt.Hour() >= 8 && timeslot.EndedAt.Hour() <= 17)

			// INFO: Timeslot is on the hour or half hour
			assert.True(t, timeslot.StartedAt.Minute() == 0 || timeslot.StartedAt.Minute() == 30)
			assert.True(t, timeslot.EndedAt.Minute() == 0 || timeslot.EndedAt.Minute() == 30)

			// INFO: Timesolt is 30 minutes long
			assert.Equal(t, timeslot.EndedAt.Sub(timeslot.StartedAt), time.Minute*30)
		}
	})

	t.Run("Trainer with appointments", func(t *testing.T) {
		// INFO: Get initial availability
		timeslots, err := store.GetTrainerAvailability(1, from, to)
		assert.NoError(t, err)
		assert.NotNil(t, timeslots)
		assert.NotZero(t, len(*timeslots))

		// INFO: Create appointment on first timeslot
		createdAppointment, err := store.CreateAppointment(&models.Appointment{
			UserID:    1,
			TrainerID: 1,
			StartedAt: (*timeslots)[0].StartedAt,
			EndedAt:   (*timeslots)[0].EndedAt,
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdAppointment)

		// INFO: Get updated availability
		updatedTimeslots, err := store.GetTrainerAvailability(1, from, to)
		assert.NoError(t, err)
		assert.NotNil(t, updatedTimeslots)
		assert.Len(t, *updatedTimeslots, len(*timeslots)-1)
		assert.NotEqual(t, (*timeslots)[0].StartedAt, (*updatedTimeslots)[0].StartedAt)
		assert.NotEqual(t, (*timeslots)[0].EndedAt, (*updatedTimeslots)[0].EndedAt)
	})
}
