package store

import (
	"future-app/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getTestAppointment() *models.Appointment {
	tz := time.FixedZone(models.GLOBAL_TZ, models.GLOBAL_TZ_OFFSET)
	startsAt := time.Date(2030, 7, 5, 8, 0, 0, 0, tz)
	endsAt := startsAt.Add(time.Minute * 30)

	return &models.Appointment{
		UserID:    1,
		TrainerID: 1,
		StartsAt:  startsAt,
		EndsAt:    endsAt,
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
	assert.Equal(t, appointment.StartsAt, createdAppointment.StartsAt)
	assert.Equal(t, appointment.EndsAt, createdAppointment.EndsAt)
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
			StartsAt:  appointment.StartsAt,
			EndsAt:    appointment.EndsAt,
		})
		assert.Error(t, err)
		assert.Equal(t, "Timeslot is not available", err.Error())
	})

	t.Run("User busy during timeslot", func(t *testing.T) {
		err = store.ValidateAvailableTimeslot(&models.Appointment{
			UserID:    1,
			TrainerID: 2,
			StartsAt:  appointment.StartsAt,
			EndsAt:    appointment.EndsAt,
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
		appointments, err := store.GetAppointmentsByTrainerID(1, appointment.StartsAt.Add(-time.Hour), appointment.EndsAt.Add(time.Hour))
		assert.NoError(t, err)
		assert.NotNil(t, appointments)
		assert.Len(t, appointments, 1)
		assert.Equal(t, createdAppointment.ID, appointments[0].ID)
		assert.Equal(t, createdAppointment.UserID, appointments[0].UserID)
		assert.Equal(t, createdAppointment.TrainerID, appointments[0].TrainerID)
		assert.Equal(t, createdAppointment.StartsAt, appointments[0].StartsAt)
		assert.Equal(t, createdAppointment.EndsAt, appointments[0].EndsAt)
	})

	t.Run("Appointment overlaps timeframe start", func(t *testing.T) {
		appointments, err := store.GetAppointmentsByTrainerID(1, appointment.StartsAt.Add(time.Minute*15), appointment.EndsAt.Add(time.Hour))
		assert.NoError(t, err)
		assert.NotNil(t, appointments)
		assert.Len(t, appointments, 1)
		assert.Equal(t, createdAppointment.ID, appointments[0].ID)
	})

	t.Run("Appointment overlaps timeframe end", func(t *testing.T) {
		appointments, err := store.GetAppointmentsByTrainerID(1, appointment.StartsAt.Add(-time.Hour), appointment.EndsAt.Add(-time.Minute*15))
		assert.NoError(t, err)
		assert.NotNil(t, appointments)
		assert.Len(t, appointments, 1)
		assert.Equal(t, createdAppointment.ID, appointments[0].ID)
	})

	t.Run("Get trainer with no appointments", func(t *testing.T) {
		appointments, err := store.GetAppointmentsByTrainerID(2, appointment.StartsAt.Add(-time.Hour), appointment.EndsAt.Add(time.Hour))
		assert.NoError(t, err)
		assert.NotNil(t, appointments)
		assert.Len(t, appointments, 0)
	})

	t.Run("Get trainer with no appointments within timeframe", func(t *testing.T) {
		appointments, err := store.GetAppointmentsByTrainerID(1, appointment.EndsAt.Add(time.Hour), appointment.EndsAt.Add(time.Hour*2))
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
	startsAt := time.Date(2030, 7, 5, 0, 0, 0, 0, tz) // Friday midnight
	endsAt := time.Date(2030, 7, 8, 0, 0, 0, 0, tz)   // Monday midnight

	t.Run("Trainer with no appointments", func(t *testing.T) {
		timeslots, err := store.GetTrainerAvailability(1, startsAt, endsAt)
		assert.NoError(t, err)
		assert.NotNil(t, timeslots)
		assert.NotZero(t, len(*timeslots))

		// INFO: Timeslots are in order
		for i := 1; i < len(*timeslots); i++ {
			assert.True(t, (*timeslots)[i].StartsAt.After((*timeslots)[i-1].StartsAt))
		}

		for _, timeslot := range *timeslots {
			// INFO: Timeslot is between starts_at and ends_at
			assert.True(t, timeslot.StartsAt.After(startsAt) || timeslot.StartsAt.Equal(startsAt))
			assert.True(t, timeslot.EndsAt.Before(endsAt) || timeslot.EndsAt.Equal(endsAt))

			// INFO: Timeslot is not on Saturday or Sunday
			assert.NotEqual(t, time.Saturday, timeslot.StartsAt.Weekday())
			assert.NotEqual(t, time.Saturday, timeslot.EndsAt.Weekday())
			assert.NotEqual(t, time.Sunday, timeslot.StartsAt.Weekday())
			assert.NotEqual(t, time.Sunday, timeslot.EndsAt.Weekday())

			// INFO: Timeslot is between 8am and 5pm
			assert.True(t, timeslot.StartsAt.Hour() >= 8 && timeslot.StartsAt.Hour() < 17)
			assert.True(t, timeslot.EndsAt.Hour() >= 8 && timeslot.EndsAt.Hour() <= 17)

			// INFO: Timeslot is on the hour or half hour
			assert.True(t, timeslot.StartsAt.Minute() == 0 || timeslot.StartsAt.Minute() == 30)
			assert.True(t, timeslot.EndsAt.Minute() == 0 || timeslot.EndsAt.Minute() == 30)

			// INFO: Timesolt is 30 minutes long
			assert.Equal(t, timeslot.EndsAt.Sub(timeslot.StartsAt), time.Minute*30)
		}
	})

	t.Run("Trainer with appointments", func(t *testing.T) {
		// INFO: Get initial availability
		timeslots, err := store.GetTrainerAvailability(1, startsAt, endsAt)
		assert.NoError(t, err)
		assert.NotNil(t, timeslots)
		assert.NotZero(t, len(*timeslots))

		// INFO: Create appointment on first timeslot
		createdAppointment, err := store.CreateAppointment(&models.Appointment{
			UserID:    1,
			TrainerID: 1,
			StartsAt:  (*timeslots)[0].StartsAt,
			EndsAt:    (*timeslots)[0].EndsAt,
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdAppointment)

		// INFO: Get updated availability
		updatedTimeslots, err := store.GetTrainerAvailability(1, startsAt, endsAt)
		assert.NoError(t, err)
		assert.NotNil(t, updatedTimeslots)
		assert.Len(t, *updatedTimeslots, len(*timeslots)-1)
		assert.NotEqual(t, (*timeslots)[0].StartsAt, (*updatedTimeslots)[0].StartsAt)
		assert.NotEqual(t, (*timeslots)[0].EndsAt, (*updatedTimeslots)[0].EndsAt)
	})
}
