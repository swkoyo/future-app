package server

import (
	"future-app/store"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var testStore *store.Store
var apiServer *APIServer

func setup() error {
	db, err := store.NewTestStore()
	if err != nil {
		return err
	}

	testStore = db
	err = testStore.Init()
	if err != nil {
		return err
	}

	apiServer = NewAPIServer(":3001", testStore)

	return nil
}

func teardown() {
	testStore.Close()
}

func TestPostAppointment(t *testing.T) {
	err := setup()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}
	defer teardown()

	e := apiServer.echo

	t.Run("Invalid dates", func(t *testing.T) {
		body := `{
        "user_id":    1,
        "trainer_id": 1,
        "started_at": "2020-07-08T20:00:00-08:00",
        "ended_at":   "2020-07-08T20:30:00-08:00"
        }`
		req := httptest.NewRequest(http.MethodPost, "/appointments", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if err := apiServer.handlePostAppointment(c); assert.NotNil(t, err) {
			he, ok := err.(*echo.HTTPError)
			if ok {
				assert.Equal(t, http.StatusBadRequest, he.Code)
			}
		}
	})

	t.Run("Valid dates, timezone should be changed to PST -8", func(t *testing.T) {
		body := `{
        "user_id":    1,
        "trainer_id": 1,
        "started_at": "2030-07-08T20:00:00Z",
        "ended_at":   "2030-07-08T20:30:00Z"
        }`
		req := httptest.NewRequest(http.MethodPost, "/appointments", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, apiServer.handlePostAppointment(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			expectedBody := `{
            "id":1,
            "user_id":1,
            "trainer_id":1,
            "started_at":"2030-07-08T12:00:00-08:00",
            "ended_at":"2030-07-08T12:30:00-08:00"
            }`
			assert.JSONEq(t, expectedBody, rec.Body.String())
		}
	})
}

func TestGetTrainerAppointments(t *testing.T) {
	err := setup()
	if err != nil {
		t.Fatalf("failed to setup test: %v", err)
	}
	defer teardown()

	e := apiServer.echo

	t.Run("Invalid date format", func(t *testing.T) {
		q := make(url.Values)
		q.Set("from", "2020-07-08")
		q.Set("to", "2020-07-08")
		req := httptest.NewRequest(http.MethodGet, "/trainers/1/appointments?"+q.Encode(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/trainers/:trainer_id/appointments")
		c.SetParamNames("trainer_id")
		c.SetParamValues("1")

		if err := apiServer.handleGetTrainerAppointments(c); assert.NotNil(t, err) {
			he, ok := err.(*echo.HTTPError)
			if ok {
				assert.Equal(t, http.StatusBadRequest, he.Code)
			}
		}
	})

	t.Run("Invalid timeframe (from after to)", func(t *testing.T) {
		q := make(url.Values)
		q.Set("from", "2020-07-08T00:00:00Z")
		q.Set("to", "2020-07-05T00:00:00Z")
		req := httptest.NewRequest(http.MethodGet, "/trainers/1/appointments?"+q.Encode(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/trainers/:trainer_id/appointments")
		c.SetParamNames("trainer_id")
		c.SetParamValues("1")

		if err := apiServer.handleGetTrainerAppointments(c); assert.NotNil(t, err) {
			he, ok := err.(*echo.HTTPError)
			if ok {
				assert.Equal(t, http.StatusBadRequest, he.Code)
			}
		}
	})

	t.Run("Invalid timeframe (more than 90 days)", func(t *testing.T) {
		q := make(url.Values)
		q.Set("from", "2020-07-01T00:00:00Z")
		q.Set("to", "2020-10-01T00:00:00Z")
		req := httptest.NewRequest(http.MethodGet, "/trainers/1/appointments?"+q.Encode(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/trainers/:trainer_id/appointments")
		c.SetParamNames("trainer_id")
		c.SetParamValues("1")

		if err := apiServer.handleGetTrainerAppointments(c); assert.NotNil(t, err) {
			he, ok := err.(*echo.HTTPError)
			if ok {
				assert.Equal(t, http.StatusBadRequest, he.Code)
			}
		}
	})

	t.Run("Valid timeframe (90 days)", func(t *testing.T) {
		q := make(url.Values)
		q.Set("from", "2020-07-01T00:00:00Z")
		q.Set("to", "2020-09-29T00:00:00Z")
		req := httptest.NewRequest(http.MethodGet, "/trainers/1/appointments?"+q.Encode(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/trainers/:trainer_id/appointments")
		c.SetParamNames("trainer_id")
		c.SetParamValues("1")

		if assert.NoError(t, apiServer.handleGetTrainerAppointments(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.JSONEq(t, `[]`, rec.Body.String())
		}
	})
}
