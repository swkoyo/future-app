package server

import (
	s "future-app/store"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type APIServer struct {
	echo  *echo.Echo
	port  string
	store *s.Store
}

func NewAPIServer(port string, store *s.Store) *APIServer {
	e := echo.New()
	NewLogger()

	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(LoggingMiddleware)

	e.Validator = NewCustomValidator()

	s := &APIServer{port: port, echo: e, store: store}

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	e.POST("/appointments", s.handlePostAppointment)

	return s
}

func (s *APIServer) Run() {
	Logger.Fatal().Msg(s.echo.Start(s.port).Error())
}
