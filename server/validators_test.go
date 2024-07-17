package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomValidator(t *testing.T) {
	cv := NewCustomValidator()

	t.Run("PostAppointmentReq Valid Input", func(t *testing.T) {
		req := PostAppointmentReq{
			UserID:    1,
			TrainerID: 1,
			StartedAt: "2030-07-08T20:00:00Z",
			EndedAt:   "2030-07-08T20:30:00Z",
		}
		err := cv.Validate(req)
		assert.NoError(t, err)
	})

	t.Run("PostAppointmentReq Invalid UserID", func(t *testing.T) {
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

	t.Run("PostAppointmentReq Invalid TrainerID", func(t *testing.T) {
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

	t.Run("PostAppointmentReq Invalid Date Format", func(t *testing.T) {
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
