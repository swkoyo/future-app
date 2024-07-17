package server

import (
	"future-app/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *APIServer) handlePostAppointment(c echo.Context) error {
	req := new(PostAppointmentReq)
	logger := GetEchoLogger(c)

	if err := c.Bind(req); err != nil {
		logger.Error().Err(err).Msg("Failed to bind request")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		logger.Error().Err(err).Msg("Failed to validate request")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	parsedStartsAt, _ := models.ParseDateStr(req.StartsAt)
	parsedEndsAt, _ := models.ParseDateStr(req.EndsAt)

	appointment, err := models.NewAppointment(
		req.UserID,
		req.TrainerID,
		parsedStartsAt,
		parsedEndsAt,
	)

	if err != nil {
		logger.Error().Err(err).Msg("Failed to create appointment")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := s.store.ValidateAvailableTimeslot(appointment); err != nil {
		logger.Error().Err(err).Msg("Failed to validate timeslot")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	logger.Info().Interface("appointment", appointment).Msg("Creating appointment")

	res, err := s.store.CreateAppointment(appointment)

	if err != nil {
		logger.Error().Err(err).Msg("Failed to create appointment")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	logger.Info().Int("appointment_id", res.ID).Msg("Appointment created")

	return c.JSON(http.StatusCreated, res)
}

func (s *APIServer) handleGetTrainerAppointments(c echo.Context) error {
	req := new(GetTrainerAppointmentsReq)
	logger := GetEchoLogger(c)

	if err := c.Bind(req); err != nil {
		logger.Error().Err(err).Msg("Failed to bind request")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		logger.Error().Err(err).Msg("Failed to validate request")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	parsedStartsAt, _ := models.ParseDateStr(req.StartsAt)
	parsedEndsAt, _ := models.ParseDateStr(req.EndsAt)

	appointments, err := s.store.GetAppointmentsByTrainerID(
		req.TrainerID,
		parsedStartsAt,
		parsedEndsAt,
	)

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get appointments")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, appointments)
}

func (s *APIServer) handleGetTrainerAvailability(c echo.Context) error {
	req := new(GetTrainerAvailabilityReq)
	logger := GetEchoLogger(c)

	if err := c.Bind(req); err != nil {
		logger.Error().Err(err).Msg("Failed to bind request")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		logger.Error().Err(err).Msg("Failed to validate request")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	parsedStartsAt, _ := models.ParseDateStr(req.StartsAt)
	parsedEndsAt, _ := models.ParseDateStr(req.EndsAt)

	timeSlots, err := s.store.GetTrainerAvailability(
		req.TrainerID,
		parsedStartsAt,
		parsedEndsAt,
	)

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get availability")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, timeSlots)
}
