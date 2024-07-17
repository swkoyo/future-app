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
		t.Log(appointments)
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
