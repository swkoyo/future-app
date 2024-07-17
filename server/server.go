package server

import (
	"future-app/common"
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

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestID := c.Response().Header().Get(echo.HeaderXRequestID)
		logger := common.Logger.With().Str("request_id", requestID).Logger()
		c.Set("logger", logger)

		logger.Info().Fields(map[string]interface{}{
			"method": c.Request().Method,
			"uri":    c.Request().URL.Path,
			"query":  c.Request().URL.RawQuery,
		}).Msg("Incoming request")

		err := next(c)
		if err != nil {
			logger.Error().Fields(map[string]interface{}{
				"error": err.Error(),
			}).Msg("Response")
			return err
		}

		return nil
	}
}

func NewAPIServer(port string, store *s.Store) *APIServer {
	e := echo.New()
	common.NewLogger()

	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(LoggingMiddleware)

	s := &APIServer{port: port, echo: e, store: store}

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	return s
}

func (s *APIServer) Run() {
	common.Logger.Fatal().Msg(s.echo.Start(s.port).Error())
}
