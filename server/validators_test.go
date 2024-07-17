package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostAppointmentReqValidator(t *testing.T) {
	cv := NewCustomValidator()

	t.Run("Valid Input", func(t *testing.T) {
		req := PostAppointmentReq{
			UserID:    1,
			TrainerID: 1,
			StartedAt: "2030-07-08T20:00:00Z",
			EndedAt:   "2030-07-08T20:30:00Z",
		}
		err := cv.Validate(req)
		assert.NoError(t, err)
	})

	t.Run("Invalid UserID", func(t *testing.T) {
		req := PostAppointmentReq{
			UserID:    -1,
			TrainerID: 1,
			StartedAt: "2030-07-08T20:00:00Z",
			EndedAt:   "2030-07-08T20:30:00Z",
		}
		err := cv.Validate(req)
		assert.Error(t, err)
		assert.Equal(t, "UserID must be 1 or greater", err.Error())
	})

	t.Run("Invalid TrainerID", func(t *testing.T) {
		req := PostAppointmentReq{
			UserID:    1,
			TrainerID: -1,
			StartedAt: "2030-07-08T20:00:00Z",
			EndedAt:   "2030-07-08T20:30:00Z",
		}
		err := cv.Validate(req)
		assert.Error(t, err)
		assert.Equal(t, "TrainerID must be 1 or greater", err.Error())
	})

	t.Run("Invalid Date Format", func(t *testing.T) {
		req := PostAppointmentReq{
			UserID:    1,
			TrainerID: 1,
			StartedAt: "2030-07-08 20:00:00",
			EndedAt:   "2030-07-08 20:30:00",
		}
		err := cv.Validate(req)
		assert.Error(t, err)
		assert.Equal(t, "StartedAt does not match the 2006-01-02T15:04:05Z07:00 format", err.Error())
	})
}

func TestGetTrainerAppointmentsReqValidator(t *testing.T) {
	cv := NewCustomValidator()

	t.Run("Valid Input", func(t *testing.T) {
		req := GetTrainerAppointmentsReq{
			TrainerID: 1,
			From:      "2030-07-08T20:00:00Z",
			To:        "2030-07-09T20:00:00Z",
		}
		err := cv.Validate(req)
		assert.NoError(t, err)
	})

	t.Run("Invalid Timeframe", func(t *testing.T) {
		req := GetTrainerAppointmentsReq{
			TrainerID: 1,
			From:      "2030-07-08T20:00:00Z",
			To:        "2030-10-09T20:00:00Z",
		}
		err := cv.Validate(req)
		assert.Error(t, err)
		assert.Equal(t, "Timeframe must be 90 days or lower", err.Error())
	})
}

func TestGetTrainerAvailabiliyReqValidator(t *testing.T) {
	cv := NewCustomValidator()

	t.Run("Invalid From Date", func(t *testing.T) {
		req := GetTrainerAvailabilityReq{
			TrainerID: 1,
			From:      "2019-07-08T20:00:00Z",
			To:        "2030-07-08T20:30:00Z",
		}
		err := cv.Validate(req)
		assert.Error(t, err)
		assert.Equal(t, "From must be a future date", err.Error())
	})
}
