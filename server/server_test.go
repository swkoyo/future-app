package server

import (
	"future-app/store"
	"net/http"
	"net/http/httptest"
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
